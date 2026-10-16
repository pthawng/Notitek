# API chung của Notitek — hợp đồng đề xuất v1

Ngày 13/10/2026. Trạng thái: thiết kế để triển khai và review tích hợp, chưa phải API đã chạy. Chủ hợp đồng: Notitek. Baseline: SRS 1.5 và DEC v1.9; áp dụng DB01–DB06 ngày 13/10/2026. Endpoint chung `/v1/notifications` theo DEC-02. Các quyết định PO ngày 11/10/2026 (DEC bản 1.6, mục 2.13–2.15) đã áp dụng: R1 gọi REST, broker để sau; OAuth2 client credentials; trường JSON snake_case và header `X-Correlation-Id`; context `organization_id` + `application_id`; Order gửi liên hệ người nhận (bắt buộc ở mốc giao thất bại); 10 mốc Order và 4 thông báo tài khoản User thuộc R1; audience vận hành; User cấp lại link qua API.

## 1. Một hợp đồng cho các module

Module gửi **ý định thông báo có cấu trúc**, không gửi HTML, provider, token Push hay cấu hình retry. Notitek quản lý mẫu, mục đích, kênh, định danh gửi, ngân sách, retry và bằng chứng theo loại thông báo đã đăng ký.

```text
User (OTP, sync) ───────────── OAuth2 + POST /v1/notifications ──┐
Outbox nguồn → dispatcher ──── OAuth2 + POST /v1/notifications ──┤
                                                                  ├─ kiểm tra hợp đồng/quyền/chống trùng
                                                                  ├─ sync: provider trong cửa sổ
                                                                  └─ async: DB + outbox → worker
```

R1 chỉ có đường REST: thông báo async của nguồn (Order, User) được ghi vào outbox của nguồn cùng transaction nghiệp vụ, rồi dispatcher của nguồn gọi API này với reference ổn định. Broker và adapter sự kiện để phân kỳ sau (DEC-02).

API độc lập ngôn ngữ và module. Thêm loại thông báo bằng cách đăng ký schema và chính sách, không thêm endpoint riêng cho từng module. Một endpoint không đồng nghĩa cho phép mọi service gửi mọi loại thông báo.

Phạm vi artifact này là API **ra lệnh gửi, tra cứu và hủy** dùng chung toàn hệ thống. API hộp tin, SSE, đăng ký Push, webhook provider và tra cứu vận hành có chủ thể/xác thực khác, được đặc tả ở các file riêng bên dưới. Không đưa nghiệp vụ xác nhận thiết bị hay xác thực OTP vào Notitek.

- [OpenAPI 3.1.1](openapi.json): schema envelope, 17 profile R1 (7 User, 10 Order), 5 operation, response và lỗi.
- [Ví dụ](examples/): dữ liệu giả, không gửi thử tới các địa chỉ/số minh họa. Phải thay thời điểm trước khi gọi môi trường thử.
- [Hộp tin, SSE và đăng ký Push](inbox-and-push.md) (C1, C2), [webhook provider và tra cứu vận hành](provider-webhooks-and-ops.md) (C3, C5): bản nháp.
- [Yêu cầu tích hợp gửi đội User](user-integration-requirements.md) và [gửi đội Order](order-integration-requirements.md): việc phía nguồn phải làm, gồm API resolve người nhận (CF-11), cấp lại link (CF-10) và đề xuất thay §11 hợp đồng Order.
- [Kiểm tra hợp đồng](../../scripts/validate_notification_contract.py): kiểm tra OpenAPI, ví dụ và ca dữ liệu không hợp lệ.

## 2. Các endpoint

| Method | Endpoint | Kết quả |
|---|---|---|
| POST | `/v1/notifications` | Sync: 200 với kết quả gửi; async mới: 202 sau commit bền vững |
| GET | `/v1/notifications/{notification_id}` | 200 snapshot hiện tại, không chứa dữ liệu bí mật |
| GET | `/v1/notifications?source_ref=...` | Tra đúng một yêu cầu; 404 nếu không có/không được xem. `source` mặc định là nguồn của token; chỉ vận hành được truyền nguồn khác |
| GET | `/v1/notifications?subject_ref=...&cursor=...&page_size=20` | Liệt kê yêu cầu của nguồn theo đối tượng nghiệp vụ (ví dụ đơn), mới nhất trước, cursor; tối đa 50/trang |
| POST | `/v1/notifications/{notification_id}/cancel` | Hủy phần chưa đưa ra ngoài; 200 báo số lượng |
| POST | `/v1/notifications/cancel` | Hủy theo `cancel_group` và context; đặt dấu hủy kể cả chưa có yêu cầu |

Request/response dùng JSON UTF-8 qua HTTPS; tên trường và query param dùng **snake_case** theo quy ước nền tảng. Header `X-Correlation-Id` (tùy chọn) nối log/trace giữa các service, được trả lại trong response, không tham gia chống trùng. Response giữ envelope của User/Order: `{error, message, data}`. `message` chỉ để hiển thị; consumer dùng mã và trạng thái, không phân tích câu chữ. Mọi response dùng `Cache-Control: no-store`; gateway tắt ghi body trên các endpoint này.

Ví dụ gọi trong môi trường thử sau khi thay dữ liệu và thời điểm trong file mẫu:

```http
POST /v1/notifications
Authorization: Bearer <service-access-token>
Content-Type: application/json
Idempotency-Key: user.otp-001
X-Correlation-Id: corr-demo-001
```

Body dùng [otp.request.json](examples/otp.request.json); `source_ref` phải bằng header. Hai phản hồi điển hình là [provider đã nhận](examples/otp-accepted.response.json) và [kết quả chưa rõ](examples/otp-unknown.response.json), đều HTTP 200 nhưng khác `data.sync_result`.

## 3. Xác thực và quyền

`Authorization: Bearer <access-token>` lấy qua OAuth2 client credentials của User/Authorization. URL cấp token, issuer, audience và cách kiểm tra token dùng cấu hình nền tảng; không tạo một hệ thống cấp token mới. Hiện User/Authorization **chưa cấp client credentials cho service** (service nội bộ dùng `X-Service-Token` tĩnh); đây là điều kiện tiên quyết CF-09. Cùng cơ chế dùng cho chiều Notitek gọi User (resolve người nhận, cấp lại link).

Registry của Notitek ánh xạ client được xác thực sang `source` ổn định, các `(type,type_version)` được phép và `application_id`/brand/organization được phép. Môi trường lấy từ deployment. Payload không được tự đặt `source`, môi trường, cost center hoặc quyền. Không dùng token frontend để gọi API ra lệnh gửi. Dùng chung prefix không dùng chung quyền: gateway phân biệt method/route và loại token; router ưu tiên các sub-path cố định `inbox`, `push-subscriptions`, `cancel`, `ops` trước `/{notification_id}`.

Scope OAuth theo quy ước `<domain>.<tài_nguyên>.<hành_động>` (CF-09): nguồn gọi API này cần `notification.request.send`, `notification.request.view`, `notification.request.cancel`; Notitek gọi User bằng `user.notification_recipient.resolve` và `user.notification_link.reissue`.

Quyền cần kiểm tra cho cả gửi, replay, tra cứu và hủy. `source` trong query phải khớp nguồn từ token, trừ vận hành có quyền xem đúng tài nguyên. Worker dùng service identity và kiểm tra scope/quyền người nhận hiện tại. `DENY`, `NOT_EVALUATED` và không xác minh được quyền không trở thành cho phép. ID không thuộc quyền đọc/hủy trả 404 để không lộ sự tồn tại.

## 4. Envelope gửi chung

| Trường | Quy tắc |
|---|---|
| `source_ref` | Reference ổn định của một ý định/sự việc; bằng header `Idempotency-Key`. Không chứa OTP, email, số điện thoại hay token |
| `subject_ref` | Tùy chọn; đối tượng nghiệp vụ mà thông báo nói về (B10). Bắt buộc với `order.*` và bằng `order_id`. Dùng để liệt kê và tra cứu vận hành, không phải khóa chống trùng, không chứa PII |
| `type`, `type_version` | Mã và phiên bản schema của loại thông báo đã đăng ký |
| `mode` | `sync` hoặc `async`; profile quyết định giá trị được phép, caller không tự đổi OTP sang async |
| `occurred_at` | Thời điểm nguồn tạo ý định/sự việc, giữ nguyên khi retry |
| `context` | `context_type` ACCOUNT: `application_id`, không organization. WORKSPACE: `application_id` và `organization_id` (Shop là Organization loại SHOP). Loại organization yêu cầu theo profile được kiểm tra runtime với User. `brand_id` tùy chọn nhưng phải được xác minh |
| `recipients` | Tối đa 100 selector; CONTACT là liên hệ, ACCOUNT là identity, AUDIENCE là quy tắc tập nhận đã duyệt |
| `data` | Biến không bí mật theo schema profile; string UTF-8, mỗi biến tối đa 2 KiB |
| `secrets` | Biến bí mật theo profile, chỉ ghi vào request; không xuất hiện trong response/log/trace |
| `send_deadline` | Thời điểm muộn nhất được bắt đầu gửi; bắt buộc, tính từ sự việc gốc, không kéo dài khi retry |
| `expires_at` | Hạn nghiệp vụ nguồn cung cấp, bắt buộc với OTP/lời mời/cảnh báo thiết bị; không phải quyền cho Notitek quyết định hiệu lực nghiệp vụ |
| `cancel_group` | Nhóm hủy ổn định, namespace theo nguồn/môi trường/context; dùng ID lần giao thất bại thay vì gom mọi lần theo order_id |

Không nhận `channels`, `templateCode`, HTML tự do, callback URL, priority hay retry count trong envelope v1. Mẫu được chọn từ type/version/context và ghim `config_version`. Giới hạn danh sách kênh, mục đích và cost center do registry/cấu hình đã duyệt quyết định. Kill-switch và quyền hiện tại vẫn có hiệu lực với tin đang chờ.

Với yêu cầu mới, kiểm tra `occurred_at <= send_deadline`, `received_at < send_deadline` và `send_deadline <= expires_at` nếu có (OTP có thêm ràng buộc ở mục 6); violation trả 422. Mọi timestamp phải có offset và được so theo instant UTC. Giới hạn riêng từng profile tiếp tục áp dụng. Các quan hệ giữa trường này là validation runtime, ngoài schema JSON.

Phân loại CONTACT không tạo tài khoản. `recipient_ref` là định danh người/vai trò do nguồn cấp trong yêu cầu, duy nhất trong mảng. Không hợp nhất hai identity chỉ vì trùng email/số điện thoại. AUDIENCE phải nằm trong allowlist của profile; không nhận biểu thức truy vấn/quyền tùy ý từ caller.

## 5. Các profile khởi điểm

| Type v1 | Mode/context | Đầu vào | Chính sách |
|---|---|---|---|
| `user.otp` | sync / ACCOUNT | Một CONTACT phone hoặc email; data.challenge_id, data.purpose; secrets.otp | SECURITY; phone → SMS, email → Email, kênh tương ứng bắt buộc; không lưu OTP ra đĩa, không retry nền |
| `user.invitation` | async / WORKSPACE (organization loại SHOP cho SHOP_MEMBER; loại bất kỳ cho EMPLOYEE) | Một CONTACT email hoặc phone; invitation_id, invitation_kind, workspace_name, inviter_name (bắt buộc với SHOP_MEMBER); secrets.invitation_url | TRANSACTIONAL; email → Email, phone → SMS; deadline ≤ min(occurred_at + 24h, expires_at) |
| `user.account_activation` | async / ACCOUNT | Một CONTACT email; activation_id, display_name tùy chọn; secrets.activation_url; expires_at | TRANSACTIONAL; Email bắt buộc; deadline ≤ min(occurred_at + 24h, expires_at); quy tắc URL/TTL bí mật và cấp lại link như lời mời |
| `user.password_changed` | async / ACCOUNT | 1–2 CONTACT (phone và/hoặc email đã xác minh); changed_via CHANGE/RESET | SECURITY; các CONTACT cùng một nhóm ALTERNATIVE, đạt khi một kênh đạt; deadline ≤ occurred_at + 24h |
| `user.identifier_changed` | async / ACCOUNT | Một CONTACT là giá trị **cũ**; identifier_type, masked_new_value | SECURITY; không gửi giá trị mới đầy đủ; deadline ≤ occurred_at + 24h |
| `user.employee_offboarded` | async / WORKSPACE | Một CONTACT; offboarding_id, organization_name, effective_at | TRANSACTIONAL; deadline ≤ occurred_at + 24h |
| `user.device_alert` | async / ACCOUNT | Một ACCOUNT; alert_id, device_label, location_label; expires_at | SECURITY; In-app bắt buộc; hành động và trạng thái thiết bị vẫn do User cung cấp |

**10 mốc Order** (danh sách trigger của hợp đồng Order). Tất cả là async / WORKSPACE với organization loại SHOP, TRANSACTIONAL, không có `secrets`. Mọi mốc dùng chung:

- `recipients` gồm đúng một `SHOP_ORDER_ALERT_RECEIVERS`, CONTACT người nhận hàng (`phone` bắt buộc, `email` tùy chọn) và tối đa một `ORDER_OPS_ALERT_RECEIVERS`. CONTACT **bắt buộc** ở `order.delivery_failed`, **tùy chọn** (0–1) ở 9 mốc còn lại; Order bỏ liên hệ ở mốc không nhắn người nhận để tối thiểu hóa dữ liệu.
- `source_ref`: `{order_id}:{type}` với mốc xảy ra một lần; `{order_id}:{type}:{id lần thất bại/vi phạm}` với `order.pickup_failed`, `order.delivery_failed`, `order.sla_breached` để mỗi lần là một yêu cầu riêng.
- `data` gồm `order_id`, `order_code`, `order_status` (mã `SPF-xxxx`), `status_reason_code` và `tracking_url` (HTTPS) tùy chọn.

| Type v1 | Trạng thái nguồn | Trường riêng | `cancel_group` | Kênh R1 (đã chốt) |
|---|---|---|---|---|
| `order.created` | SPF-0301 | — | Không | Shop: In-app bắt buộc, Push tùy chọn |
| `order.cancelled` | SPF-0201 | — | Không | Shop: In-app, Push tùy chọn |
| `order.pickup_failed` | SPF-0402 | pickup_failure_id, reason_code | Bắt buộc (ID lần lấy thất bại) | Shop: In-app, Push tùy chọn; vận hành: In-app nếu có selector |
| `order.picked_up` | SPF-0501 | — | Không | Shop: In-app, Push tùy chọn |
| `order.delivery_failed` | SPF-0802 | delivery_failure_id, reason_code | Bắt buộc (ID lần giao thất bại) | Shop: In-app, Push tùy chọn; người nhận hàng: ZNS, SMS thay thế; vận hành: In-app nếu có selector |
| `order.delivered` | SPF-0901 | — | Không | Shop: In-app, Push tùy chọn |
| `order.partially_delivered` | SPF-0902 | — | Không | Shop: In-app, Push tùy chọn |
| `order.return_confirmed` | SPF-1002 | — | Không | Shop: In-app, Push tùy chọn |
| `order.returned` | SPF-1201/1202/1203 | — | Không | Shop: In-app, Push tùy chọn |
| `order.sla_breached` | Không đổi trạng thái | sla_breach_id, sla_code, expected_to | Bắt buộc (ID lần vi phạm) | Shop: In-app; vận hành: In-app nếu có selector |

Order gửi liên hệ người nhận hàng trong yêu cầu (A1; bắt buộc ở giao thất bại, tùy chọn ở mốc khác theo Q2); Notitek **không gọi ngược sang Order** để lấy liên hệ. Chính sách của mốc quyết định có nhắn người nhận hay không; R1 chỉ nhắn người nhận ở `order.delivery_failed`. Mốc không dùng liên hệ thì Notitek không lưu liên hệ đó ngoài bộ nhớ xử lý yêu cầu. Thêm kênh người nhận cho mốc khác cần template ZNS được Zalo duyệt (CF-08) và ngân sách (CF-02).

Ngôn ngữ R1 chỉ tiếng Việt, envelope không có `locale` (B7); khi cần sẽ thêm `locale` tùy chọn theo quy tắc mở rộng ở mục 11. OTP qua Zalo/ZNS không thuộc R1 (B6), xem lại ở R2.

Kênh của profile User được suy từ loại CONTACT (phone → SMS, email → Email), không do caller chỉ định; ZNS cho đích phone của profile User chỉ bật khi profile được đăng ký lại. `purpose` chỉ dùng chọn nội dung tin theo mục đích OTP; Notitek không kiểm tra hiệu lực mục đích. Hạn 24h và chính sách retry của bốn profile mới đã chốt làm giá trị khởi điểm (B5), vận hành điều chỉnh khi có số liệu thật; Email retry theo bảng DEC-06 của lời mời, SMS bảo mật 24/7 thử tối đa 2 lần cách 5′.

Schema cụ thể ở `components.schemas.*Request`; `SendRequest.x-type-profiles` là bảng ánh xạ. Handler phải kiểm tra **cả envelope và profile**; extension OpenAPI này không tự chạy validation. Mã type mới phải được deploy/đăng ký trước khi caller dùng. Không dùng một schema `data: any` làm tiêu chuẩn tích hợp.

`user.*` chỉ cấp cho nguồn User. `order.*` chỉ cấp cho nguồn Order. Tên type được phân namespace theo chủ nghiệp vụ nhưng đi cùng một API.

Đối với Shop, nguồn chỉ gửi selector `SHOP_ORDER_ALERT_RECEIVERS`; Notitek resolve thành viên của organization trong context có quyền `notification.order_alert.receive`. Với vận hành/CSKH, selector `ORDER_OPS_ALERT_RECEIVERS` resolve người có quyền `notification.order_ops_alert.receive` và data scope bao phủ organization Shop đó. Cache ≤ 5 phút. Kiểm tra lại quyền trước gửi và khi xem hộp tin. Tối đa **100 người nhận sau resolve**, tính toàn bộ trường hợp trong yêu cầu; không cắt danh sách âm thầm. Nếu vượt giới hạn, dừng toàn yêu cầu trước dispatch với `RECIPIENT_LIMIT_EXCEEDED`; nguồn cần chia thành các yêu cầu có reference riêng qua profile được duyệt. Mất dịch vụ quyền thì chờ đến deadline, sau đó `RECIPIENT_UNRESOLVED`.

Hạn từng kênh là min(send_deadline, expires_at nếu có, hạn kênh theo DEC-06). Order đặt `send_deadline`; mặc định đề xuất occurred_at + 12h cho mọi mốc. Push tối đa 15 phút; ZNS/SMS chỉ 07:00–21:00 Asia/Ho_Chi_Minh. Không đổi deadline sang hôm sau để vượt hạn. Device alert dùng hạn do User cấp; giá trị 15 phút trong ví dụ chỉ minh họa, không phải chính sách mặc định mới.

URL lời mời phải là HTTPS trên host/path được đăng ký cho app, không userinfo, không redirect tùy ý. Token không xuất hiện trong GET/kết quả/preview/Push; worker không fetch URL này. Lưu URL bí mật đã mã hóa để gửi nền, TTL tối đa min(expires_at, received_at + 10 phút) theo DEC-09; xóa sớm khi không cần nữa. Sau TTL, Notitek không giữ token cũ. Khi cần retry sau TTL (trong `send_deadline`), worker gọi **API cấp lại link của User** (quyết định B1, CF-10) với `source_ref`, `type` và `invitation_id`/`activation_id` để nhận URL HTTPS mới kèm `expires_at`, rồi gửi lần thử tiếp theo; URL mới cũng theo TTL 10 phút. Yêu cầu bắt buộc với API này:

- User quyết định link còn hợp lệ hay không. Lời mời đã chấp nhận/thu hồi/hết hạn, hoặc tài khoản đã kích hoạt, thì trả mã từ chối; Notitek dừng delivery với `SOURCE_INVALIDATED`, không retry.
- Cấp lại cho Notitek **không tính** vào giới hạn cấp lại của người dùng (cooldown 60s, tối đa 3 lần) và **không vô hiệu link trước đó** (Q1): link cũ vẫn dùng được đến hạn riêng của nó, vì lần gửi trước có thể đã tới người nhận dù kết quả chưa rõ.
- Gọi bằng OAuth2 client credentials của Notitek; idempotent theo `(notification_id, attempt)` để không sinh nhiều link cho một lần thử.
- User không phản hồi hoặc lỗi tạm thời: Notitek thử lại theo lịch retry của kênh; hết `send_deadline` thì kết thúc `SECRET_EXPIRED`.

Tên endpoint và schema do User đặc tả; đến khi có, retry sau TTL kết thúc `SECRET_EXPIRED`.

## 6. Đồng bộ và bất đồng bộ

OTP: deadline hiệu lực = min(send_deadline, received_at + 5s, expires_at). `send_deadline` phải ≤ received_at + 5s và ≤ expires_at; `occurred_at` là lúc User tạo challenge, không bị ràng buộc +5s (User commit trước rồi mới gọi, theo B3); chỉ bắt đầu lần gọi provider đầu trong 2s từ tiếp nhận và còn trong deadline. Dành thời gian để ghi metadata/kết quả và trả phản hồi trong cửa sổ. Không đổi sang queue khi hết thời gian. Chỉ thử ngay một provider SMS dự phòng khi chắc chắn chưa gửi và còn thời gian, theo DEC-05. OTP qua Email không có provider dự phòng trong R1.

HTTP 200 trả `sync_result` = `PROVIDER_ACCEPTED`, `FAILED` hoặc `UNKNOWN`. `error=false` nghĩa thao tác API có kết quả hợp lệ, **không đồng nghĩa gửi thành công**. Provider từ chối chắc chắn → FAILED; timeout sau khi có thể đã gửi → UNKNOWN. UNKNOWN không được retry provider tự động. Callback muộn cập nhật bằng chứng, không kích hoạt mã OTP hoặc thay kết quả mà User đã nhận.

Cách User xử lý (B2, B3; chi tiết trong [yêu cầu gửi đội User](user-integration-requirements.md)):

| `sync_result` | User làm |
|---|---|
| `PROVIDER_ACCEPTED` | Challenge hiệu lực, báo đã gửi |
| `FAILED` | Báo lỗi kèm `reason_code`; cho gửi lại theo giới hạn của User |
| `UNKNOWN` | Challenge **vẫn hiệu lực**, báo "chưa xác nhận đã gửi"; chỉ cho gửi lại sau cooldown 3 phút; không tính là lần gửi thất bại |

User gọi API này **ngoài transaction DB**: commit challenge ở trạng thái chờ gửi, gọi Notitek, rồi cập nhật kết quả trong một transaction ngắn.

`sync_result` được ghim khi hoàn tất phản hồi lần đầu; GET/replay vẫn giữ nó. `outcome`, deliveries và revision phản ánh bằng chứng mới. Nếu process chết sau khi bắt đầu gọi provider, phục hồi ghi `sync_result=UNKNOWN`, không phát lại OTP. Hệ thống giữ khóa/metadata bền vững trước dispatch nhưng bí mật OTP chỉ ở RAM; crash trước dispatch có thể kết thúc FAILED nếu chứng minh chưa gọi provider. Không thể chứng minh thì UNKNOWN.

Async: chỉ trả 202 khi request, payload được bảo vệ, khóa chống trùng và outbox công việc đã commit trong cùng transaction (hoặc bảo đảm tương đương). Provider chưa được gọi trước ranh giới này. Mất kết nối sau commit không tạo việc mới khi replay. Replay async trả 200 snapshot hiện tại, cùng notification_id; worker tiếp tục phần còn hợp lệ, không gửi lại phần đã đạt.

`admission=ACCEPTED` ở resource là metadata yêu cầu đã ghi nhận, không phải bằng chứng gửi. Không trả HTTP 202 cho OTP. Lỗi envelope/xác thực trước admission trả 4xx; lỗi nghiệp vụ của yêu cầu đã ghi nhận thể hiện trong resource.

## 7. Chống trùng và phục hồi

Khóa tạo yêu cầu = `(source từ token, environment, source_ref)`. Header bằng source_ref; khác nhau trả 400 `REFERENCE_MISMATCH`. Mỗi lần OTP mới có source_ref mới do User cấp. `X-Correlation-Id`/trace không phải khóa chống trùng.

So sánh toàn bộ body đã parse theo JCS, gồm secrets, recipient order, deadline và occurred_at; thứ tự property không quan trọng, thứ tự array có ý nghĩa. Không tự chèn default trước khi so. Chỉ header transport/auth/trace được phép đổi. Dùng HMAC có keyId, giữ khả năng so bằng khóa cũ trong cửa sổ 7 ngày; không dùng hash OTP không khóa, không lưu canonical plaintext bí mật. Unique constraint/atomic claim phải ngăn hai request song song cùng dispatch.

- Cùng khóa/cùng body: trả cùng notification_id và snapshot; không tính thêm quota gửi cho replay.
- Khác body: 409 `IDEMPOTENCY_CONFLICT`, không trả giá trị cũ/mới chứa bí mật.
- Sync đang chạy: 409 `IDEMPOTENCY_IN_PROGRESS`, Retry-After; caller tra GET, không tạo reference mới để vượt khóa.
- Kiểm tra quyền và giới hạn tài nguyên đầu vào trước lookup; sau lookup mới xét thời hạn cho **yêu cầu mới**. Replay hợp lệ vẫn tra được khi deadline đã qua.
- Chống trùng được bảo đảm 7 ngày từ accepted; metadata giữ 180 ngày đến 180 ngày + 1 tháng theo accepted_at (DB01). Nếu còn metadata nhưng hết khóa, trả 409 `REFERENCE_RETIRED`, không tự gửi lại. Không được tái sử dụng reference.
- Đầu vào mới có occurred_at quá 7 ngày trả 422 `STALE_REQUEST` (event: `STALE_EVENT`). Thời điểm vượt tương lai > 30s trả `INVALID_TIME`; 30s là đề xuất dung sai đồng hồ v1, cần NTP. Không sửa occurred_at để vượt kiểm tra replay. Sau khi metadata bị xóa không hứa chống trùng vô hạn.

Network timeout/HTTP 5xx từ proxy không chứng minh gửi thất bại. Async: tra theo reference, rồi replay **nguyên body** nếu còn hạn. OTP: tra kết quả; không tự phát lại POST với deadline mới hoặc reference mới. GET 404 ngay sau timeout có thể xảy ra trước commit; tiếp tục đối chiếu trong cửa sổ của User. Mã OTP mới là quyết định nghiệp vụ của User.

## 8. Kết quả và tổng hợp

Response không trả request data/secrets, token, địa chỉ nhận đầy đủ hay provider credential. Dùng recipient_ref và delivery_id để đối chiếu. Một yêu cầu có tối đa 100 người nhận × 5 kênh = 500 delivery summary; lịch sử từng lần thử thuộc API vận hành riêng.

- `processing_state`: PENDING, PROCESSING, CLOSED. CLOSED nghĩa không còn chủ động gửi; callback muộn vẫn có thể cập nhật bằng chứng.
- `outcome`: PENDING, ACHIEVED, PARTIAL, FAILED, UNKNOWN, CANCELLED, EXPIRED, SUPPRESSED.
- `deliveries`: kết quả riêng từng người/kênh, requirement REQUIRED/OPTIONAL/ALTERNATIVE và requirement_group. `case_code` (tùy chọn) cho biết delivery thuộc trường hợp gửi nào (ví dụ `SHOP_ALERT`, `CONSIGNEE_NOTICE`) để consumer gom nhóm; kết quả tổng hợp theo từng trường hợp chưa có ở R1.
- `revision`: số tăng đơn điệu mỗi thay đổi quan sát được, phục vụ polling/event consumer.

Ánh xạ hiển thị cho frontend nội bộ (B11). User lưu `notification_id` cùng lời mời/kích hoạt và trả `delivery_status` qua API của User; frontend không gọi thẳng API ra lệnh gửi:

| Notitek | `delivery_status` |
|---|---|
| `outcome` PENDING, hoặc PARTIAL khi `processing_state` chưa CLOSED, chưa có lần thử lỗi | QUEUED |
| PROCESSING và có delivery với `attempt_count` ≥ 1 và `next_attempt_at` khác null | RETRYING |
| `outcome` ACHIEVED, hoặc PARTIAL khi `processing_state` CLOSED | SENT |
| `outcome` FAILED, EXPIRED, CANCELLED, SUPPRESSED | FAILED |
| `outcome` UNKNOWN | UNCONFIRMED (giá trị mới, frontend phải bổ sung; không gộp vào SENT) |

Mỗi nhóm bắt buộc phải đạt ngưỡng bằng chứng của DEC-10; nhóm thay thế đạt khi một kênh hợp lệ đạt. Tất cả nhóm bắt buộc đạt → ACHIEVED dù kênh OPTIONAL lỗi; lỗi tùy chọn vẫn hiện trong deliveries. Một số nhóm đạt → PARTIAL; chưa nhóm nào đạt và còn việc có thể tiếp tục → PENDING; không còn việc, còn kết quả chưa rõ → UNKNOWN. Nếu không đạt và không còn UNKNOWN, toàn bộ bị hủy/hết hạn/chặn tương ứng → CANCELLED/EXPIRED/SUPPRESSED; trường hợp hỗn hợp còn lại → FAILED. Không có người nhận hợp lệ → SUPPRESSED với lý do, không dùng tập rỗng để báo ACHIEVED.

PROVIDER_ACCEPTED là bằng chứng provider nhận, không phải DELIVERED. In-app đạt ở AVAILABLE. Trạng thái đọc không nằm trong kết quả gửi. Bằng chứng đã ghi nhận không bị xóa khi có phản hồi xung đột; cần lưu delivery outcome bất lợi và reason_code để tra cứu, không che hard bounce. `UNKNOWN` quá 24h đóng xử lý, vẫn là UNKNOWN.

`outcome` chỉ là kết quả chốt khi `processing_state` là CLOSED (quyết định 12/10/2026). Khi chưa CLOSED, PARTIAL nghĩa là đã đạt một phần và phần còn lại đang chạy, nên consumer xét `processing_state` trước. Hard bounce hoặc báo giao thất bại đến **sau** PROVIDER_ACCEPTED chuyển delivery sang FAILED (`reason_code` BOUNCE_HARD hoặc DELIVERY_FAILED), tính lại `outcome` và tăng `revision`; nếu đến sau DELIVERED thì delivery giữ DELIVERED và bằng chứng được đánh xung đột. Ma trận chuyển trạng thái đầy đủ ở [DATABASE v0.8 (LỊCH SỬ)](<../Notitek - DATABASE v0.8 (LỊCH SỬ).md>) mục 4.4 (mô hình v0.9 không chép lại ma trận); code hiện thực ở `internal/notification/domain/transition.go`. Kênh dự phòng chưa từng chạy (ví dụ SMS khi ZNS đã được nhận hoặc `UNKNOWN`) có state SUPPRESSED với `reason_code` FALLBACK_NOT_NEEDED, nghĩa là không cần gửi chứ không phải bị chặn (DEC 2.17).

## 9. Hủy và sự kiện đến muộn

Hủy là command riêng, có Idempotency-Key riêng; namespace `(source, environment, cancellation, key)`. Fingerprint gồm method/path/body; reuse key với target/reason khác → 409. Ghi command và fence dừng trước khi trả 200; replay trả receipt ban đầu, không thực hiện lại. Receipt là snapshot tại thời điểm áp dụng, GET request cho trạng thái mới nhất.

Worker phải atomically claim quyền dispatch và kiểm tra fence. Cancel thắng race: không gọi provider. Worker đã claim/có thể gọi provider: tính notRecallable/unknown, không hứa thu hồi. Các số đếm delivery trong receipt không chồng nhau: cancelled, notRecallable và unknown. Không sửa bản đã gửi thành cancelled.

Hủy group dùng `(source, environment, context đầy đủ, cancel_group)`; không hủy xuyên Shop/app. Mỗi context được chuẩn hóa theo registry giống lúc gửi. Ghi tombstone 72h kể cả matched_requests=0; yêu cầu cùng group đến sau trong cửa sổ được ghi nhận nhưng dừng phần chưa dispatch, không được gửi. Sau 72h không bảo đảm hủy-trước-event; nguồn phải có deadline phù hợp, không phát lại sự việc đã giải quyết với thời điểm mới. Một lần giao thất bại mới có group mới. Hủy theo notification_id chưa tồn tại trả 404; dùng group khi cần hủy trước sự kiện.

## 10. HTTP và mã lỗi

Danh mục mã ổn định nằm ở `ApiError.x-error-codes` trong OpenAPI; consumer xử lý mã chưa biết theo HTTP status.

### 10.1 Cấu trúc lỗi

Mọi lỗi dùng envelope của nền tảng `{error, message, data}` với `error: true`. `message` là câu giải thích cho người đọc (tiếng Việt ở R1, không dùng để rẽ nhánh). `data` theo **RFC 9457 (Problem Details for HTTP APIs)** cộng các trường mở rộng:

| Trường | Nguồn | Ý nghĩa |
|---|---|---|
| `type` | RFC 9457 | URI ổn định của loại lỗi, `urn:supership:notitek:error:<mã viết kebab-case>`. Là URN nên không hứa có tài liệu tại URL |
| `title` | RFC 9457 | Tóm tắt ngắn, không đổi giữa các lần xảy ra |
| `status` | RFC 9457 | Mã HTTP, lặp lại trong body cho consumer chỉ giữ body |
| `instance` | RFC 9457 | Đường dẫn của request lỗi; không có query string |
| `code` | Mở rộng | Mã máy ổn định, dùng để rẽ nhánh. Mã chưa biết thì xử lý theo `status` |
| `trace_id` | Mở rộng | Giá trị `X-Correlation-Id`; gửi kèm khi cần hỗ trợ |
| `retryable` | Mở rộng | `true` chỉ khi lặp lại đúng request sau đó có thể thành công và chưa có việc nào được nhận: 429, 503, `IDEMPOTENCY_IN_PROGRESS`. `INTERNAL_ERROR` là `false`: tra cứu theo reference trước khi gửi lại |
| `retry_after_seconds` | Mở rộng | Có khi `retryable`; bằng header `Retry-After` (RFC 9110 §10.2.3, dạng số giây) |
| `details` | Mở rộng | Lỗi theo từng trường, xem dưới |

`envelope.message` đóng vai trò `detail` của RFC 9457, nên không lặp lại trong `data`. Body vẫn là `application/json` (không phải `application/problem+json`) vì envelope là hợp đồng chung với User và Order; parser Problem Details đọc được `data`.

Mỗi phần tử `details`:

| Trường | Ý nghĩa |
|---|---|
| `field` | JSON Pointer (RFC 6901) vào body, ví dụ `/recipients/0/phone`, hoặc `/` cho cả body; với header, query, path là tên của chúng |
| `location` | `body`, `header`, `query` hoặc `path` |
| `code` | Lý do ổn định: `required`, `invalid_type`, `invalid_format`, `invalid_value`, `out_of_range`, `too_long`, `unknown_field`, `malformed_json`, `mismatch` |
| `message` | Câu giải thích an toàn để hiển thị |

**Không bao giờ lặp lại giá trị đầu vào** trong `message` hay `details`, vì body có thể chứa OTP, link hay số điện thoại. Lỗi kiểu dữ liệu khi giải mã JSON không có chỉ số mảng nếu bộ giải mã không báo (ví dụ `/recipients/phone`).

Ví dụ: [validation-error](examples/validation-error.response.json), [rate-limited](examples/rate-limited.response.json), [conflict](examples/conflict.response.json).

Đường dẫn hay phương thức không có route cũng trả envelope: 404 `NOT_FOUND`, hoặc 405 `METHOD_NOT_ALLOWED` kèm header `Allow` (RFC 9110 §15.5.6). Mọi phản hồi có `X-Content-Type-Options: nosniff` và `Cache-Control: no-store`.

| HTTP | Mã điển hình | Hành vi caller |
|---|---|---|
| 400 | INVALID_REQUEST, REFERENCE_MISMATCH | Sửa envelope/header; details chỉ có field/code, không echo input |
| 401/403 | UNAUTHENTICATED, FORBIDDEN | Sửa xác thực/quyền, không retry mù quáng |
| 404 | NOT_FOUND | Không có hoặc không được xem; không suy ra chắc chắn chưa gửi sau network timeout |
| 405 | METHOD_NOT_ALLOWED | Dùng phương thức trong header `Allow` |
| 409 | IDEMPOTENCY_CONFLICT, IDEMPOTENCY_IN_PROGRESS, REFERENCE_RETIRED | Không đổi key để né; IN_PROGRESS thì tra cứu |
| 413/415 | PAYLOAD_TOO_LARGE, UNSUPPORTED_MEDIA_TYPE | Giới hạn 65536 byte JSON đã giải nén; application/json |
| 422 | UNKNOWN_TYPE, UNSUPPORTED_TYPE_VERSION, INVALID_PROFILE, INVALID_TIME, STALE_REQUEST, DEADLINE_EXPIRED, UNSUPPORTED_PURPOSE | Sửa hợp đồng; MARKETING bị từ chối R1 |
| 429 | RATE_LIMITED | Admission chưa tạo công việc; Retry-After, vẫn tuân deadline |
| 503 | SERVICE_UNAVAILABLE | Handler không nhận bền vững; lỗi proxy/network vẫn phải đối chiếu |

Mặc định 50 request/s/source, burst theo SRS-P05; replay không tính quota gửi nhưng vẫn chịu giới hạn API. Giới hạn tần suất nghiệp vụ (OTP theo đích/mục đích, số tin theo đơn) do module nguồn sở hữu và thực thi trước khi gọi; Notitek không chặn thêm theo số điện thoại. Notitek chỉ trả 429 khi hạn chế admission, hoặc delivery FAILED/reason_code=PROVIDER_REJECTED khi provider từ chối do giới hạn của họ. Giới hạn byte mỗi biến cần kiểm tra UTF-8 thực, không chỉ maxLength ký tự trong schema.

Mã kết quả gửi như BLOCKED_QUOTA, CONTACT_BLOCKED, PROVIDER_REJECTED, PROVIDER_TIMEOUT, RECIPIENT_UNRESOLVED, RECIPIENT_LIMIT_EXCEEDED, SECRET_EXPIRED, SOURCE_INVALIDATED, EXPIRED nằm ở resource/delivery, không biến thành lỗi HTTP sau khi đã accepted. Lỗi nội bộ bất ngờ không trả stack trace hoặc secret.

### 10.1. Chính sách dữ liệu và quyền vận hành đã chốt

Theo DEC v1.9 mục 2.16: metadata gửi và inbox xóa cùng cây request, lấy accepted_at làm mốc; giữ từ 180 ngày đến 180 ngày + 1 tháng. Nội dung title/preview của inbox giữ/hiển thị 90 ngày từ created_at rồi làm NULL; nội dung và liên hệ không được kéo TTL theo partition.

Liên hệ đầy đủ trong selector tạm (DB: recipients_spec) được dọn khi PLAN xong hoặc request đóng, đồng thời không giữ quá min(accepted_at + 24 giờ, send_deadline). Job dọn độc lập với tiến trình PLAN; hết hạn không lấy payload cũ để tiếp tục lập kế hoạch. contact_value đã chuyển vào delivery vẫn theo quy tắc che sau 30 ngày của DEC-09. API inbox/ops không trả nội dung đã hết hạn dù metadata còn.

| Hành động | Quyền | Gán mặc định |
|---|---|---|
| Tạo/sửa/thử/kích hoạt/thu hồi kết nối | `notification.connection.manage` | Platform admin |
| Sửa hạn mức ngày theo ngân sách được giao | `notification.quota.manage` | Platform admin |
| Dừng/bật gửi theo loại và kênh | `notification.pause.manage` | Vận hành Notification |
| Bật/tắt Push tùy chọn của Shop | `notification.shop_channel.manage` | Vận hành Notification |

User/Authorization đăng ký/gán các quyền này qua CF-07. Tên trường JSON hiện hành không đổi theo tên cột DB v0.9; persistence ánh xạ theo từ điển mục 2/5. Thiết kế chuẩn hóa không tự thay hợp đồng GET sync đang xử lý.

## 11. Tích hợp từ outbox nguồn và mở rộng

R1: Order ghi ý định thông báo của 10 mốc vào `outbox_events` cùng transaction đổi trạng thái; dispatcher outbox của Order (CF-04) gọi `POST /v1/notifications` với profile `order.*`. Phản hồi 2xx (202 mới, 200 replay) mới đánh dấu row đã gửi. 409 `IDEMPOTENCY_CONFLICT`, 4xx hợp đồng và 422 đưa row vào dead-letter để xử lý tay; timeout, 429 và 5xx retry **nguyên body** với cùng `source_ref` trong hạn. Các thông báo async của User (lời mời, kích hoạt, đổi mật khẩu, đổi định danh, nghỉ việc) cũng nên đi qua outbox của User thay cho gửi best-effort sau commit hiện nay.

Reference của User (B4) phải tất định theo ý định, ví dụ `otp:{challenge_id}:{delivery_version}`, `invitation:{invitation_id}:issue:{n}`, `activation:{activation_id}:issue:{n}`, `password-changed:{event_id}`, `identifier-changed:{event_id}`, `offboarded:{job_id}`. `occurred_at` là thời điểm nghiệp vụ có offset (B8), không phải lúc gọi API.

Tên ở ranh giới tích hợp là `type` của Notitek; tên event nội bộ của Order chỉ dùng để ánh xạ:

| Hợp đồng Order | Outbox Order hiện có | `type` Notitek |
|---|---|---|
| ORDER_CREATED | `order.created.v1` (phải ghi khi vào SPF-0301, không ghi trước kết quả NVC) | `order.created` |
| ORDER_CANCELLED | `order.cancelled.v1` | `order.cancelled` |
| ORDER_PICKUP_FAILED | chưa có | `order.pickup_failed` |
| ORDER_PICKED_UP | chưa có | `order.picked_up` |
| ORDER_DELIVERY_FAILED | chưa có | `order.delivery_failed` |
| ORDER_DELIVERED | chưa có | `order.delivered` |
| ORDER_PARTIALLY_DELIVERED | chưa có | `order.partially_delivered` |
| ORDER_RETURN_CONFIRMED | `order.return.confirmed.v1` | `order.return_confirmed` |
| ORDER_RETURNED | chưa có | `order.returned` |
| ORDER_SLA_BREACHED | chưa có (chưa có bộ phát hiện SLA) | `order.sla_breached` |

Ánh xạ Order: `subject_ref` = `order_id`; `source_ref` theo mẫu ở mục 5, lấy từ `event_key` của outbox (Order phải đổi định dạng `event_key` hiện tại `{orderId}:{eventType}` vì UNIQUE sẽ chặn lần thất bại thứ hai; không dùng `outbox_event_id` sinh ngẫu nhiên); thời điểm nghiệp vụ → `occurred_at` (outbox phải lưu riêng, không dùng `created_at`); `shop_id` → `context.organization_id`; ID lần lấy/giao thất bại hoặc vi phạm SLA → trường riêng và `cancel_group`; `order_parties` người nhận → CONTACT. Khi lần thất bại được giải quyết (giao lại thành công, chuyển hoàn, hủy), Order gọi `POST /v1/notifications/cancel` với `cancel_group` tương ứng, cũng qua outbox.

Phân kỳ sau khi có broker: event nghiệp vụ theo CloudEvents đi qua adapter vào cùng application service; reference phải trùng với reference REST cho cùng ý định, không phát một ý định qua hai đường với hai reference khác nhau.

R1 theo DB05 (PO chốt 13/10/2026): consumer theo dõi kết quả bằng GET, polling theo revision. Mỗi thay đổi quan sát được cập nhật revision/updated_at trong cùng transaction; **không ghi notification.result_changed.v1 vào outbox R1**. outbox_job chỉ chứa PLAN/SEND nội bộ. Event kết quả, schema và consumer được thiết kế khi làm phân kỳ broker; không giữ các row event chưa sử dụng để chờ phát sau.

Thêm module/loại: đăng ký owner, allowlist caller/scope, schema data/secrets/recipients, purpose, mode, policy kênh/mẫu/ngân sách, deadline, tiêu chí đạt và ví dụ; deploy registry trước khi bật producer. Không mở rộng envelope bằng cách nhận tùy ý field không khai báo. Envelope v1 từ chối field lạ; mở rộng server/schema trước khi client gửi field mới. Thay ý nghĩa hoặc thêm field bắt buộc của type cần type_version mới; thay phá vỡ envelope cần `/v2`. Consumer response phải bỏ qua field mới và xử lý an toàn giá trị chưa biết, không mặc định thành công.

## 12. Điều kiện kiểm chứng trước khi chạy thật

| Tình huống | Kỳ vọng |
|---|---|
| Hai POST cùng key đồng thời; replay sau mất response | Một request/tập delivery, không gửi hai lần |
| Đổi OTP/người nhận/context/deadline cùng key | 409, không lộ payload trong lỗi |
| OTP provider timeout hoặc process chết sau dispatch | UNKNOWN; không worker retry; GET/callback không kích hoạt mã |
| Crash trước/sau commit async | Không trả accepted trước commit; sau commit khôi phục outbox được |
| Hủy group trước event, race worker/cancel | Tombstone chặn gửi; phần đã ra ngoài được báo riêng |
| User/Order giả type/application/organization; xem ID nguồn khác | Từ chối, không lộ metadata |
| Audience >100, mất nguồn quyền, thành viên rời Shop | Không cắt ngầm/không gửi ngoài quyền |
| ZNS chưa rõ, quota hết, callback trùng | Không fallback sai, không tạo lần gửi mới |
| Quét log/trace/DB/DLQ/response | Không có OTP; token lời mời được bảo vệ và xóa đúng TTL |
| Các request mẫu + schema sai có chủ đích | Validate được profile đúng và từ chối dữ liệu sai |

Schema validation chỉ chứng minh hình dạng dữ liệu. Những kiểm tra thời gian, byte, quyền, đồng thời, crash, quota và provider ở trên cần implementation/contract test khi code; chưa được tuyên bố đã chạy.

## 13. Cơ sở kỹ thuật

OpenAPI dùng phiên bản [3.1.1](https://spec.openapis.org/oas/v3.1.1.html). Chuẩn hóa JSON để so payload theo [RFC 8785](https://www.rfc-editor.org/rfc/rfc8785). Đây là lựa chọn biểu diễn của thiết kế, không thay các quyết định nghiệp vụ trong SRS/DEC.
