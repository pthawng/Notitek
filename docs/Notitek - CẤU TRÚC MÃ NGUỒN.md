# Notitek - CẤU TRÚC MÃ NGUỒN

**Phiên bản:** 0.4. **Ngày:** 13/10/2026. Bản 0.4: đồng bộ DB01–DB06 (DEC 1.9): mốc lưu/xóa, TTL selector 24 giờ, R1 chỉ polling kết quả và tách quyền vận hành. Bản 0.3: kết nối kênh, hạn mức, thiết lập Push của Shop và kill-switch chuyển sang DB và API vận hành (module `connection`, DEC-07 v1.7); credential mã hóa trong DB. Bản 0.2: thêm kiến trúc mã (mục 10), đồng thời và transaction (mục 11), nền tảng production (mục 12), bắt buộc tuân thủ bằng CI (mục 13). **Trạng thái:** Đề xuất để Tech Lead chốt thành ADR-08 trước M0. **Căn cứ:** ADR-01 (Go, RabbitMQ nội bộ, PostgreSQL, Redis), [DEC 1.6](<Notitek - QUYẾT ĐỊNH PO CHO CÁC DEC.md>), [API v1](api/README.md).

Tài liệu quy định bố cục thư mục, ranh giới module, hướng phụ thuộc và cách mở rộng. Mục tiêu là khi User, Order hoặc module mới của SuperPlatform cần thêm thông báo, phần lớn thay đổi chỉ nằm ở cấu hình registry, không phải code lõi.

## 1 Đặc tính của Notification trong SPF quyết định cấu trúc

| Đặc tính | Hệ quả lên cấu trúc |
|---|---|
| Một API chung cho mọi module, mở rộng bằng type đã đăng ký | Registry là nguồn sự thật của từng type; code không có struct riêng cho từng type |
| OTP đồng bộ ≤ 5s, phần còn lại async qua outbox | Tách tiến trình `api` (nhận lệnh, OTP) và `worker` (gửi nền); `relay` chuyển outbox sang RabbitMQ |
| Chống trùng, kết quả theo người/kênh, `revision` | Yêu cầu và các delivery là một aggregate, cập nhật trong một transaction, do module `notification` sở hữu |
| Fallback chỉ khi chắc chắn chưa gửi; `UNKNOWN` không retry | Adapter kênh phải phân loại kết quả chuẩn; quyết định retry/fallback nằm ở lõi |
| Bí mật (OTP, link) không ra đĩa/log; TTL ≤ 10′ | Module riêng cho bí mật tạm thời; logger có lớp che dữ liệu |
| Người nhận theo quyền (Shop, vận hành) từ User | Module `recipient` dùng port; client User nằm ở `integration/user` |
| Cấu hình 3 cấp, ghim `config_version`, kill-switch tức thì, thêm kết nối kênh không cần deploy | Chính sách dạng file có snapshot bất biến trong DB; kết nối kênh, kill-switch, thiết lập Push của Shop và hạn mức là dữ liệu runtime trong DB, đổi qua API vận hành có audit |
| In-app + SSE, Push, callback provider, tra cứu vận hành | Mỗi năng lực một module, dùng chung use case của lõi khi đổi kết quả gửi |
| Ngân sách theo ngày, danh sách chặn, thời gian lưu | Module `budget`, `suppression`, `retention` riêng, không rải vào lõi |

## 2 Cây thư mục

```text
notification-spf/
├── cmd/                                  # Mỗi thư mục là một binary, chỉ gọi bootstrap
│   ├── api/main.go                       # HTTP: /v1/notifications, inbox, SSE, push-subscriptions, webhooks, ops
│   ├── worker/main.go                    # Consumer RabbitMQ, pool riêng mỗi kênh
│   ├── relay/main.go                     # Outbox Postgres → RabbitMQ (publisher confirms)
│   ├── scheduler/main.go                 # Job định kỳ: hạn gửi, giữ tin đến 07:00, đóng UNKNOWN 24h, lưu trữ
│   └── registryctl/main.go               # CLI cho CI: kiểm tra registry, sinh x-type-profiles cho OpenAPI
│
├── internal/
│   ├── bootstrap/                        # Nơi DUY NHẤT nối dependency; mỗi runtime chỉ dựng thứ nó cần
│   │   ├── config.go
│   │   ├── dependencies.go               # Postgres, Redis, RabbitMQ, OTel, logger
│   │   ├── api.go  worker.go  relay.go  scheduler.go
│   │   └── lifecycle.go                  # graceful shutdown, health/readiness
│   │
│   ├── platform/                         # Hạ tầng dùng chung, không biết gì về nghiệp vụ
│   │   ├── config/
│   │   ├── postgres/                     # pool, TxManager, advisory lock
│   │   ├── rabbitmq/                     # publisher confirms, consumer ack, DLQ
│   │   ├── redis/
│   │   ├── httpx/                        # server, middleware: X-Correlation-Id, envelope {error,message,data},
│   │   │                                 #   Cache-Control: no-store, recover, giới hạn body 64KB
│   │   ├── auth/                         # xác thực token service (OAuth2) và context token (qua User)
│   │   ├── observability/                # OTel, metrics, logger có che OTP/token/liên hệ
│   │   ├── apperr/                       # danh mục mã lỗi = ApiError.x-error-codes
│   │   ├── jcs/                          # chuẩn hóa JSON RFC 8785 cho so sánh chống trùng
│   │   ├── crypto/                       # HMAC keyring xoay vòng, AEAD cho bí mật
│   │   ├── pagination/                   # cursor (page_size, has_next, next_cursor)
│   │   ├── clock/                        # đồng hồ inject được để test hạn gửi/khung giờ
│   │   └── idgen/
│   │
│   ├── notification/                     # LÕI: vòng đời yêu cầu và delivery
│   │   ├── domain/
│   │   │   ├── request.go                # aggregate: yêu cầu + deliveries + revision
│   │   │   ├── delivery.go  attempt.go
│   │   │   ├── outcome.go                # tổng hợp ACHIEVED/PARTIAL/... theo DEC-10
│   │   │   ├── deadline.go               # min(send_deadline, expires_at, hạn kênh, khung giờ)
│   │   │   ├── cancellation.go           # hủy, tombstone 72h
│   │   │   └── errors.go
│   │   ├── application/
│   │   │   ├── ports.go                  # interface: Repository, Registry, RecipientResolver, Sender,
│   │   │   │                             #   SecretStore, Budget, Suppression, Outbox
│   │   │   ├── admit.go                  # xác thực nguồn, profile, thời gian, chống trùng
│   │   │   ├── send_sync.go              # OTP trong cửa sổ 5s
│   │   │   ├── send_async.go             # ghi yêu cầu + outbox trong một transaction
│   │   │   ├── dispatch.go               # worker gọi: claim, gửi, phân loại, retry/fallback
│   │   │   ├── record_result.go          # callback provider gọi: nâng bằng chứng, cập nhật outcome
│   │   │   ├── query.go                  # GET theo id/source_ref, liệt kê theo subject_ref
│   │   │   ├── cancel.go
│   │   │   └── jobs.go                   # hết hạn, đóng UNKNOWN 24h (scheduler gọi)
│   │   ├── adapter/
│   │   │   ├── http/                     # cài đặt server sinh từ OpenAPI
│   │   │   └── postgres/
│   │   └── module.go                     # API công khai của module cho module khác
│   │
│   ├── registry/                         # Type profile, policy, template, snapshot cấu hình
│   │   ├── domain/                       # TypeProfile, ChannelPolicy, Template, Snapshot
│   │   ├── application/                  # nạp, kiểm tra payload bằng JSON Schema, resolve theo context
│   │   ├── adapter/fs/                   # đọc config/registry
│   │   ├── adapter/postgres/             # registry_snapshots bất biến; kill-switch; thiết lập kênh của Shop
│   │   └── module.go
│   │
│   ├── connection/                       # Kết nối kênh: provider, sender/OA/domain, credential mã hóa
│   │   ├── domain/                       # Connection, trạng thái DRAFT/ACTIVE/DISABLED
│   │   ├── application/                  # tạo, gửi thử, kích hoạt, xoay credential, thu hồi
│   │   ├── adapter/postgres/             # channel_connection
│   │   └── module.go
│   │
│   ├── render/                           # Dựng nội dung từ template + biến; che dữ liệu cho preview
│   │
│   ├── recipient/                        # CONTACT / ACCOUNT / AUDIENCE → danh sách nhận
│   │   ├── audience/                     # shop_order_alert.go, ops_order_alert.go (mỗi audience một file)
│   │   ├── cache.go                      # ≤ 5 phút
│   │   └── module.go
│   │
│   ├── channel/                          # Chỉ nói chuyện với provider, không quyết định nghiệp vụ
│   │   ├── channel.go                    # Sender, Message, Result{ACCEPTED|REJECTED_DEFINITE|
│   │   │                                 #   TRANSIENT_NOT_SENT|UNKNOWN}, provider_message_id
│   │   ├── sms/<nhà-cung-cấp>/  sms/mock/
│   │   ├── email/<nhà-cung-cấp>/  email/mock/
│   │   ├── zns/
│   │   ├── push/fcm/
│   │   └── inapp/                        # ghi bản In-app qua port của inbox
│   │
│   ├── secretstore/                      # OTP chỉ trong bộ nhớ; link mã hóa, TTL ≤ 10′, xóa sớm;
│   │                                     #   lấy link mới qua integration/user (cấp lại link)
│   ├── budget/                           # hạn mức ngày luồng × kênh, cảnh báo 80%, BLOCKED_QUOTA
│   ├── suppression/                      # danh sách chặn toàn cục, CONTACT_BLOCKED (trừ SECURITY)
│   │
│   ├── inbox/                            # hộp tin, số chưa đọc, đọc idempotent, SSE (Redis chỉ báo tín hiệu)
│   │   ├── domain/  application/
│   │   ├── adapter/{http,postgres,sse}/
│   │   └── module.go
│   ├── pushsub/                          # đăng ký/hủy đích Push, chuyển chủ, vô hiệu token lỗi
│   ├── webhook/                          # xác thực + chuẩn hóa callback từng provider → notification.RecordResult
│   ├── ops/                              # API tra cứu vận hành, xem dữ liệu gốc có lý do + audit
│   ├── audit/                            # ghi audit 2 năm (xem dữ liệu nhạy cảm, kill-switch, hủy/gửi lại)
│   ├── retention/                        # selector ≤ 24h/hạn gửi; inbox nội dung 90d; metadata 180d + tối đa 1 tháng
│   ├── outbox/                           # outbox chung: ghi trong transaction, relay claim/publish/confirm
│   │   ├── application/
│   │   └── adapter/{postgres,rabbitmq}/
│   │
│   ├── integration/
│   │   └── user/                         # access-context resolve, resolve người nhận (CF-11),
│   │                                     #   cấp lại link (CF-10); có bản mock cho dev/test
│   │                                     # KHÔNG có integration/order: Notitek không gọi ngược Order (A1)
│   └── gen/
│       └── openapi/                      # code sinh từ docs/api/openapi.json; không sửa tay
│
├── config/
│   ├── registry/                         # Nguồn sự thật của loại thông báo; thay đổi qua PR (DEC-07)
│   │   ├── types/                        # <type>.v<n>.yaml: owner, source được phép, mode, purpose,
│   │   │                                 #   context, JSON Schema recipients/data/secrets, deadline
│   │   ├── policies/                     # kênh bắt buộc/tùy chọn/thay thế, retry, khung giờ, tiêu chí đạt
│   │   ├── templates/vi/<kênh>/          # nội dung theo type và kênh, có phiên bản
│   │   ├── audiences.yaml                # selector được phép và quyền tương ứng
│   │   ├── sources.yaml                  # client → source, type được phép, application/brand được phép
│   │   └── overrides/{application,brand,organization}/   # 3 cấp Global → App/Brand → Shop
│   └── app.example.yaml                  # cấu hình runtime mẫu (không chứa secret)
│
├── db/
│   ├── migrations/                       # NNNN_<mô_tả>.sql, chỉ tiến, có phiên bản
│   └── queries/<module>/*.sql            # nếu dùng sqlc; mỗi module chỉ truy vấn bảng của mình
│
├── deploy/
│   ├── compose.yaml                      # postgres, rabbitmq, redis, mock provider, mock User
│   └── mocks/                            # mock SMS/Email/ZNS/FCM và mock User cho dev
│
├── test/
│   ├── contract/                         # chạy docs/api/examples vào API thật; so OpenAPI ↔ registry
│   ├── integration/                      # Postgres/RabbitMQ/Redis thật (testcontainers)
│   ├── e2e/                              # kịch bản AT trong backlog
│   └── fixtures/
│
├── docs/                                 # tài liệu hiện có; docs/api/openapi.json là nguồn sinh code
├── scripts/                              # validate_notification_contract.py, gen.sh
├── .github/workflows/ci.yml
├── .gitattributes                        # cố định kết thúc dòng (repo đang lẫn CRLF/LF)
├── .golangci.yml                         # gồm luật phụ thuộc giữa module
├── Dockerfile                            # multi-stage, một image, chọn binary theo target
├── Makefile                              # gen, lint, test, migrate, run-<runtime>
├── go.mod  go.sum
└── README.md
```

## 3 Quy ước module

Module phức tạp dùng đủ bốn phần; module nhỏ (`render`, `budget`, `suppression`, `retention`, `audit`) được giữ phẳng. **Không tạo thư mục rỗng chỉ để đủ hình thức.**

| Phần | Chứa | Không được |
|---|---|---|
| `domain/` | Quy tắc nghiệp vụ thuần, không I/O | Import adapter, platform (trừ `clock`), module khác |
| `application/` | Use case và **port** (interface) nó cần | Gọi SQL, HTTP, AMQP trực tiếp |
| `adapter/` | Cài đặt port: HTTP handler, repository, client, consumer | Chứa quyết định nghiệp vụ |
| `module.go` | API công khai cho module khác và hàm dựng module | Lộ kiểu nội bộ của adapter |

**Hướng phụ thuộc** (CI chặn vi phạm bằng `depguard` hoặc `go-arch-lint`):

```text
cmd → bootstrap → các module → platform
adapter → application → domain
Module A chỉ gọi module B qua B/module.go hoặc qua port do A khai báo.
channel/* không import notification. notification phụ thuộc interface channel.Sender.
webhook → notification (RecordResult) + chuẩn hóa của channel.
integration/user cài đặt port do recipient, secretstore và platform/auth khai báo.
platform không import bất kỳ module nghiệp vụ nào.
```

## 4 Quyền sở hữu

| Trách nhiệm | Module sở hữu |
|---|---|
| Tiếp nhận, chống trùng, vòng đời yêu cầu | `notification` |
| Delivery, lần thử, retry, fallback, `revision` | `notification` (cùng aggregate; chỉ tách thành module riêng khi đủ lớn) |
| Tổng hợp `outcome` | `notification` |
| Profile, policy, template, `config_version`, kill-switch, thiết lập kênh của Shop | `registry` |
| Kết nối kênh (provider, sender/OA/domain, credential) | `connection` |
| Dựng nội dung | `render` |
| Resolve người nhận | `recipient` (gọi User qua `integration/user`) |
| Gọi provider và phân loại kết quả | `channel/*` |
| Xác thực và chuẩn hóa callback | `webhook`; ghi kết quả qua `notification` |
| Bí mật tạm thời | `secretstore` |
| Hạn mức (giới hạn và bộ đếm) và danh sách chặn | `budget`, `suppression` |
| Hộp tin, đọc, SSE | `inbox` |
| Đích Push | `pushsub` |
| Đưa công việc vào RabbitMQ | `outbox` |

**Quy tắc then chốt:** chỉ `notification` được thay đổi trạng thái delivery và outcome. `channel` trả kết quả đã phân loại, `webhook` gọi use case `RecordResult`, `scheduler` gọi job của `notification`. Không module nào ghi thẳng vào bảng của module khác.

## 5 Runtime dùng module nào

| Module | api | worker | relay | scheduler |
|---|---|---|---|---|
| notification | ● | ● | | ● |
| registry, render | ● | ● | | ● |
| connection | ● | ● | | |
| recipient, integration/user | ● | ● | | |
| channel/* | ● (OTP) | ● | | |
| secretstore, budget, suppression | ● | ● | | ● |
| inbox | ● | ● (ghi In-app) | | ● |
| pushsub, webhook, ops | ● | | | |
| outbox | ● (ghi) | ● (ghi) | ● (publish) | |
| retention, audit | ● (audit) | | | ● |

`relay` không nạp HTTP handler hay provider; `scheduler` không nạp provider. Mỗi file trong `bootstrap/` chỉ dựng đúng cột của nó.

## 6 Registry và phiên bản cấu hình

- `config/registry` là **nguồn duy nhất** của type. `x-type-profiles` trong `docs/api/openapi.json` được `registryctl` sinh ra; CI báo lỗi nếu hai bên lệch nhau.
- Payload được kiểm tra **lúc chạy** bằng JSON Schema trong registry. Thêm type không cần thêm struct Go.
- Khi khởi động, `api` tính hash toàn bộ registry làm `config_version` và ghi snapshot đã resolve vào `registry_snapshots` (chỉ thêm, không sửa). Yêu cầu ghim `config_version`; `worker` xử lý theo snapshot đó.
- **Dữ liệu vận hành nằm trong DB, không nằm trong file** (DEC-07 v1.7): kết nối kênh (sender/OA/domain, credential), kill-switch, bật/tắt Push của Shop, hạn mức ngày. Đổi qua API vận hành có kiểm tra quyền và audit, có hiệu lực không cần deploy. Kill-switch có hiệu lực ngay với cả tin đang chờ.
- **Kết nối kênh** tạo ở `DRAFT`, gửi thử đạt mới `ACTIVE`. Mỗi `delivery_attempt` ghi kết nối và version đã dùng. Kết nối lỗi hoặc dừng không tự chuyển sang brand khác.
- **Credential** của kết nối mã hóa envelope trước khi lưu DB, chỉ giải mã trong tiến trình gửi và xác thực callback, không bao giờ trả ra API. Khóa gốc và khóa HMAC lấy từ secret manager qua biến môi trường; không nằm trong `config/`, repo hay log.

## 7 Dữ liệu

- Một database, schema `notification`. Mỗi bảng thuộc đúng một module; module khác chỉ truy cập qua `module.go`.
- Bảng dự kiến theo module: `notification` (request, delivery, attempt, send_case, notification_message, idempotency_key, cancel_marker), `outbox` (outbox_job), `registry` (registry_snapshot, kill_switch, shop_channel_setting), `connection` (channel_connection), `secretstore` (transient_secret), `inbox` (inbox_item), `pushsub` (push_subscription), `webhook` (provider_callback, tức delivery_evidence), `budget` (quota_counter, quota_limit), `suppression` (contact_blocklist), `audit` (audit_log). Từ điển từng trường: [Mô hình dữ liệu logic](<Notitek - MÔ HÌNH DỮ LIỆU LOGIC VÀ TỪ ĐIỂN DỮ LIỆU.md>).
- Schema vật lý chốt sau ADR-03 (lưu trữ, bí mật) và ADR-04 (trạng thái, tổng hợp).

## 8 Cách mở rộng

| Nhu cầu | Làm gì | Sửa code Go |
|---|---|---|
| User/Order thêm loại thông báo hoặc mốc mới | Thêm `types/`, `policies/`, `templates/`; chạy `registryctl`; PR | Không |
| Module mới của SPF gửi thông báo | Thêm client vào `sources.yaml` và các type của module | Không |
| Đổi nội dung, kênh, retry của type có sẵn | Sửa policy/template → `config_version` mới | Không |
| Thêm nhà cung cấp cho kênh đã có | Thêm `channel/<kênh>/<nhà-cung-cấp>/` | Chỉ trong adapter |
| Thêm kết nối mới cho provider đã có (OA Zalo, sender SMS, domain Email cho brand mới) | Gọi API vận hành tạo kết nối, gửi thử, kích hoạt; không sửa repo, không deploy | Không |
| Thêm kênh mới | Thêm `channel/<kênh>/`, tạo kết nối qua API vận hành | Cục bộ |
| Thêm kiểu người nhận mới | Thêm `recipient/audience/<tên>.go` + API trong `integration/` + `audiences.yaml` | Cục bộ |
| Thay đổi phá tương thích envelope | `/v2` | Có, hiếm |

## 9 Quy ước khác

- Tên package Go viết thường, không gạch dưới; một khái niệm một package. Dùng `application/`, không dùng lẫn `app/`.
- JSON và cột DB dùng snake_case; mã lỗi theo `ApiError.x-error-codes`; thông điệp cho người dùng bằng tiếng Việt (R1).
- Thời gian lưu và so sánh theo UTC, hiển thị và khung giờ gửi theo `Asia/Ho_Chi_Minh`; luôn lấy giờ qua `platform/clock`.
- Không log payload, `secrets`, liên hệ đầy đủ hay token; lớp che dữ liệu nằm ở `platform/observability`.
- Code sinh tự động (`internal/gen/`, sqlc) không sửa tay; `make gen` phải sạch trong CI.

## 10 Kiến trúc mã: ranh giới, interface và composition root

**Composition root.** Chỉ `internal/bootstrap` tạo đối tượng và nối dependency bằng constructor injection. Không dùng biến toàn cục, `init()` có side effect hay service locator. Mỗi runtime có một hàm dựng riêng (`bootstrap.API`, `bootstrap.Worker`…) chỉ nạp module nó cần (mục 5).

**Interface do bên dùng khai báo.** Port nằm ở `application/ports.go` của module cần nó, nhỏ, theo hành vi (ví dụ `RecipientResolver.Resolve`, `Sender.Send`). Module cung cấp không định nghĩa interface "cho người khác dùng". Ví dụ: `notification` khai báo `RecipientResolver`; `recipient/module.go` cài đặt; `bootstrap` nối hai bên.

**API công khai của module.** `module.go` export một struct facade và các kiểu dữ liệu trao đổi (DTO). Không export kiểu của `adapter/` hay `domain/` nội bộ ra ngoài module. Gọi giữa module luôn đồng bộ trong cùng process; không dùng event nội bộ giữa module trong R1.

**Transaction xuyên module.** `platform/postgres.TxManager` mang transaction qua `context.Context`. Repository của nhiều module tham gia cùng một transaction khi use case cần (ví dụ `notification` + `outbox` + `secretstore` khi nhận yêu cầu async). Module dùng chung transaction nhưng **không** đọc ghi bảng của nhau.

**Code sinh tự động.** `internal/gen/openapi` chỉ được import bởi `notification/adapter/http`, `inbox/adapter/http`, `ops`… (lớp adapter). Domain và application không phụ thuộc kiểu sinh từ OpenAPI.

## 11 Đồng thời và transaction

Nguyên tắc chung: **mọi chuyển trạng thái là câu UPDATE có điều kiện trên trạng thái hiện tại**, khóa theo thứ tự cố định *request → delivery → attempt* để tránh deadlock, và không giữ transaction khi gọi mạng (provider, User, RabbitMQ).

**Chính sách dữ liệu và quyền vận hành đã chốt:**

Áp dụng DEC v1.9 mục 2.16, DB01–DB06; tên persistence theo từ điển v0.9 (56 bảng/260 trường), ánh xạ API/view ở mục 2 và đối chiếu 193 trường v0.8 ở mục 5. Các tên bảng v0.8/bản cũ trong mô tả luồng bên dưới phải được adapter ánh xạ, không tạo bảng đồng nghĩa. Ranh giới transaction/TTL v0.9 ở mục 7–8; [kế hoạch migration 3NF](<Notitek - THIẾT KẾ DATABASE 3NF VÀ AN TOÀN MIGRATION.md>) chưa được áp thành migration/persistence.

- Metadata request và bản con, gồm inbox/evidence đã khớp, xóa cùng cây theo accepted_at; giữ 180 ngày đến 180 ngày + 1 tháng. Evidence chưa khớp dùng received_at, chuyển theo tuổi request khi khớp.
- title/preview inbox làm NULL sau 90 ngày từ created_at; ID/trạng thái đọc/metadata còn tới lúc xóa request. Mọi đường đọc/làm mới cache áp dụng thời hạn hiển thị; job và backup/restore không hồi sinh nội dung đã dọn.
- TTL selector liên hệ tối đa min(accepted_at + 24 giờ, send_deadline), đồng thời dọn sớm khi PLAN xong/request đóng. Quy tắc che contact_value sau 30 ngày và bí mật TTL 10 phút vẫn độc lập.

| Hành động | Quyền | Gán mặc định |
|---|---|---|
| Tạo/sửa/thử/kích hoạt/thu hồi kết nối | `notification.connection.manage` | Platform admin |
| Sửa hạn mức ngày theo ngân sách được giao | `notification.quota.manage` | Platform admin |
| Dừng/bật gửi theo loại và kênh | `notification.pause.manage` | Vận hành Notification |
| Bật/tắt Push tùy chọn của Shop | `notification.shop_channel.manage` | Vận hành Notification |

Các quyền mới là đầu vào triển khai CF-07 của User/Authorization, không phải quyền đã cấp trên production.

### 11.1 Tiếp nhận và chống trùng

| Bước | Cách làm |
|---|---|
| Khóa chống trùng | Bảng `idempotency_key` khóa chính `(source, environment, source_ref)`, lưu `fingerprint_hmac`, `key_id`, `notification_id`, `state`, `expires_at` |
| Yêu cầu mới | `INSERT … ON CONFLICT DO NOTHING RETURNING`. Không có dòng trả về nghĩa là đã tồn tại: đọc lại, so fingerprint → replay (200) hoặc `IDEMPOTENCY_CONFLICT` (409) |
| Async | Một transaction: khóa chống trùng + request + delivery dự kiến + bí mật mã hóa + `outbox_job` → commit → 202 |
| OTP (sync) | Tx1 ghi khóa `IN_PROGRESS` + request rồi commit (yêu cầu song song nhận `IDEMPOTENCY_IN_PROGRESS`). Ghi `dispatch_started_at` và commit **trước** khi gọi provider. Gọi provider ngoài transaction. Tx2 ghi `sync_result`, chuyển khóa sang `COMPLETED`. OTP chỉ nằm trong bộ nhớ |
| Phục hồi OTP | Scheduler quét khóa `IN_PROGRESS` quá cửa sổ: có `dispatch_started_at` → `UNKNOWN`; không có → `FAILED` (chứng minh được chưa gọi provider) |
| Hủy theo nhóm và yêu cầu mới | Cả hai lấy `pg_advisory_xact_lock(hash(source, environment, organization_id, cancel_group))` trong transaction, nên tombstone và yêu cầu mới không lọt qua nhau |

### 11.2 Lập kế hoạch người nhận

Audience được mở rộng ở bước **plan** do worker chạy trước mọi lần gửi của yêu cầu đó: gọi `recipient` (ngoài transaction), rồi trong một transaction tạo delivery cho từng người/kênh và ghi planned_at; processing_state chỉ dùng PENDING/PROCESSING/CLOSED theo API, không lưu PLANNED. Vượt 100 người nhận → dừng cả yêu cầu với `RECIPIENT_LIMIT_EXCEEDED` trước khi gửi bất kỳ kênh nào. Mất User → retry plan đến `send_deadline`, sau đó `RECIPIENT_UNRESOLVED`.

**Dọn liên hệ (DB04):** transaction PLAN thành công ghi planned_at và xóa phone/email khỏi recipients_spec sau khi đã tạo các delivery cần gửi. Transaction đóng request cũng dọn selector liên hệ. Job retention độc lập dọn trước min(accepted_at + 24 giờ, send_deadline), kể cả resolver/worker bị kẹt; sau hạn không dùng lại payload để resolve/gửi. Selector không bí mật/reference và metadata có thể còn để tra cứu.

### 11.3 Outbox và relay

- **DB05:** outbox_job R1 chỉ chứa PLAN/SEND; thay đổi kết quả cập nhật revision/updated_at và nguồn polling GET. Không ghi sẵn RESULT_CHANGED hoặc notification.result_changed.v1.
- `outbox_job` có `status`, `available_at`, `claimed_by`, `claim_until`, `attempts`, `headers` (mang trace context và correlation id).
- Relay claim theo lô: `UPDATE … SET claimed_by, claim_until = now() + lease WHERE id IN (SELECT id … WHERE status = 'PENDING' AND available_at <= now() ORDER BY id LIMIT n FOR UPDATE SKIP LOCKED) RETURNING …`. Nhiều replica chạy song song không lấy trùng dòng.
- Publish với **publisher confirms**; chỉ đánh dấu `PUBLISHED` sau khi broker xác nhận. Không xác nhận hoặc bị trả về (`mandatory`) → để lease hết hạn và claim lại.
- Retry có trễ dùng `available_at` của outbox, không dựa vào TTL hay dead-letter queue của RabbitMQ để hẹn giờ.
- Vì có thể publish lặp, **consumer phải idempotent** (mục 11.4). Theo dõi số job chờ và tuổi job chờ lâu nhất.

### 11.4 Worker gửi một delivery

1. Nhận message `{delivery_id, attempt_seq}`.
2. **Tx1 (claim):** khóa request rồi delivery; kiểm tra trạng thái cho phép gửi (`PENDING`, hoặc `PROCESSING` đang chờ retry với `next_attempt_at` đến hạn), chưa bị hủy (fence), còn hạn, khung giờ, kill-switch, danh sách chặn, hạn mức (`UPDATE budget_usage SET used = used + 1 WHERE used < limit RETURNING`); chuyển delivery sang `PROCESSING`, tạo attempt có `dispatch_started_at`; commit. Không thỏa điều kiện → ghi kết quả tương ứng, ack, kết thúc.
3. Gọi provider ngoài transaction, timeout không vượt hạn còn lại của delivery.
4. **Tx2 (ghi kết quả):** `UPDATE attempt … WHERE id = ? AND state = 'STARTED'`; quyết định theo phân loại của channel: `ACCEPTED` ghi bằng chứng; `TRANSIENT_NOT_SENT` tạo `outbox_job` mới với `available_at` theo lịch retry; `REJECTED_DEFINITE` chuyển sang kênh thay thế nếu policy cho phép; `UNKNOWN` không retry, không fallback. Tính lại `outcome`, tăng `revision`; commit.
5. Ack message **sau** Tx2. Message bị giao lại khi delivery đang `PROCESSING` với attempt còn lease → bỏ qua và ack.
6. Worker chết giữa bước 2 và 4: scheduler thấy attempt `STARTED` quá lease → `UNKNOWN` (có thể đã gửi), không tự gửi lại. Riêng In-app là ghi nội bộ idempotent theo `delivery_id`, nên được thử lại.

### 11.5 Hủy, callback, In-app

- **Hủy so với gửi:** transaction hủy khóa request rồi delivery và chỉ chuyển `PENDING` hoặc `PROCESSING` đang chờ retry (chưa có attempt `STARTED`) sang `CANCELLED`; delivery đang `PROCESSING` được tính `not_recallable` hoặc `unknown`. Vì cùng khóa dòng với Tx1 của worker, bên nào giữ khóa trước thì thắng, không có trạng thái lửng.
- **Callback provider:** `INSERT provider_callback … ON CONFLICT (provider, provider_event_id) DO NOTHING`, trả 2xx ngay, xử lý qua `notification.RecordResult`. Mức bằng chứng chỉ tăng, không giảm (DEC-10). Không khớp `provider_message_id` → giữ để đối chiếu.
- **In-app:** `inbox_item` có ràng buộc duy nhất theo `delivery_id`; phát tín hiệu Redis **sau** commit. Redis mất tín hiệu thì client vẫn lấy được dữ liệu qua API danh sách.

### 11.6 Scheduler nhiều replica

Mỗi job lấy `pg_try_advisory_lock(job_id)`; replica không lấy được khóa thì bỏ lượt. Job xử lý theo lô nhỏ, idempotent, chạy lại được sau khi crash.

## 12 Nền tảng vận hành production

| Chủ đề | Quy định |
|---|---|
| **Lỗi** | Domain trả lỗi có kiểu và mã; `platform/apperr` ánh xạ sang HTTP status và mã trong `x-error-codes`. Dùng `%w` để giữ nguyên chuỗi lỗi; response không có stack trace, SQL hay dữ liệu đầu vào. Lỗi provider luôn được phân loại về 4 nhóm của `channel.Result` |
| **Log** | `log/slog` dạng JSON. Trường bắt buộc: `correlation_id`, `trace_id`, `source`, `notification_id`, `delivery_id`, `type`, `runtime`. Kiểu `Secret` cài đặt `LogValue()` trả giá trị đã che; liên hệ chỉ log dạng che. Không log body request |
| **Cấu hình** | Biến môi trường (có file YAML mẫu cho dev) nạp vào struct có kiểu, kiểm tra khi khởi động và dừng ngay nếu sai. Mỗi runtime chỉ yêu cầu cấu hình nó dùng |
| **Secret** | Khóa gốc, khóa HMAC và khóa mã hóa lấy từ secret manager qua môi trường; không nằm trong repo, `config/` hay log. Credential kết nối kênh mã hóa envelope và lưu ở DB, thêm kết nối mới không đụng hạ tầng. Khóa có `key_id`, hỗ trợ xoay vòng và đọc bằng khóa cũ trong cửa sổ 7 ngày |
| **Trace** | OpenTelemetry xuyên `api → outbox → relay → RabbitMQ → worker → provider`; trace context đi qua `outbox_job.headers` và header AMQP |
| **Metric** | Độ trễ tiếp nhận; OTP p95/p99; số job outbox chờ và tuổi cũ nhất; độ sâu queue theo kênh; lần thử theo nhóm kết quả; độ trễ provider; tỷ lệ `UNKNOWN`; mức dùng hạn mức; số kết nối SSE |
| **Cảnh báo** | OTP vượt p99 5s; job outbox chờ quá 1 phút; tăng đột biến `UNKNOWN` hoặc `REJECTED_DEFINITE`; hạn mức 80%; callback không khớp tăng |
| **Vòng đời** | `/healthz` (process sống) và `/readyz` (dependency của runtime đó sẵn sàng). Tắt mềm: `api` ngừng nhận và chờ request đang chạy (OTP tối đa 5s); `worker` hủy consume, xử lý xong message đang giữ hoặc nack để giao lại; `relay` xong lô hiện tại; scheduler nhả advisory lock |
| **Migration** | Chạy như bước triển khai riêng, không tự chạy khi ứng dụng khởi động; theo kiểu expand → migrate → contract để rolling deploy không lỗi |

**Kiểm thử:**

| Tầng | Kiểm tra gì |
|---|---|
| Unit | Domain thuần: tổng hợp outcome, tính hạn gửi, phân loại kết quả, quy tắc fallback (table-driven) |
| Application | Use case với fake port: chống trùng, hủy, retry, phục hồi OTP |
| Integration | Adapter với Postgres, RabbitMQ, Redis thật (testcontainers): claim `SKIP LOCKED`, publisher confirms, advisory lock |
| Đồng thời | Hai POST cùng key; hủy so với worker; hai relay cùng lúc; crash giữa Tx1 và Tx2 (điểm chèn lỗi) |
| Contract | Ví dụ trong `docs/api/examples` vào API thật; OpenAPI khớp registry |
| E2E | Các kịch bản AT trong backlog trên Docker Compose với mock provider và mock User |

Mọi test chạy với `go test -race`.

## 13 Bắt buộc tuân thủ bằng CI

| Kiểm tra | Công cụ hoặc cách làm | Chặn merge |
|---|---|---|
| Hướng phụ thuộc giữa module và lớp (mục 3) | `depguard` trong `golangci-lint` hoặc `go-arch-lint` | Có |
| Không import `internal/gen` ngoài adapter | Luật `depguard` | Có |
| Code sinh không lệch | `make gen` rồi `git diff --exit-code` | Có |
| Registry hợp lệ và khớp OpenAPI | `registryctl validate` + `scripts/validate_notification_contract.py` | Có |
| Lint, vet, race | `golangci-lint`, `go vet`, `go test -race ./...` | Có |
| Migration an toàn | Kiểm tra migration chỉ tiến, có review cho thay đổi phá tương thích | Có |
| Không lộ secret | Quét secret trong diff; test che log | Có |
| Contract và integration | Chạy trên Docker Compose ở pipeline | Có |

## 14 Việc cần chốt trước khi dựng

| Việc | Người chốt |
|---|---|
| Chốt tài liệu này thành ADR-08 | Tech Lead |
| Thư viện: router HTTP, `pgx` + `sqlc`, công cụ migration, client AMQP, JSON Schema validator, `oapi-codegen` | Tech Lead |
| ADR-03 và ADR-04 (trước khi viết migration); ADR-03 mô tả chi tiết mã hóa credential kết nối trong DB (đã chốt hướng), gồm thuật toán, `key_id`, xoay khóa, quyền giải mã | Tech Lead, DevOps |
| TR-01/TR-02: nhân sự Go và pipeline DevOps cho Go | Tech Lead, DevOps |
| `.gitattributes` cho kết thúc dòng | Tech Lead |
