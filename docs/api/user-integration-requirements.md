# Yêu cầu tích hợp Notitek cho module User

Ngày 11/10/2026. Trạng thái: **Bản nháp gửi đội User** (Tech Lead User, User/Authorization, FE nội bộ). Căn cứ: [DEC 1.6](<../Notitek - QUYẾT ĐỊNH PO CHO CÁC DEC.md>) mục 2.13–2.15, [hợp đồng API v1](README.md), [OpenAPI](openapi.json). PO đã duyệt hướng xử lý B2–B11, C6, C7 và D ngày 11/10/2026. Tài liệu liệt kê việc đội User cần làm trong repo `user-spf`; Notitek không sửa repo này.

## 1. Tổng hợp

| ID | Việc | Bên làm | CF | Hạn |
|---|---|---|---|---|
| U-01 | Adapter `NotificationPort` gọi `POST /v1/notifications` | users-core | CF-03 | 15/10/2026 (ký schema) |
| U-02 (B2) | Xử lý `sync_result=UNKNOWN` cho OTP | users-core | CF-03 | 15/10/2026 |
| U-03 (B3) | Gọi Notitek ngoài transaction DB | users-core | CF-03 | 15/10/2026 |
| U-04 (B4) | `source_ref` tất định và outbox thông báo | users-core | CF-03 | 22/10/2026 |
| U-05 (B8) | `occurred_at` có offset từ thời điểm nghiệp vụ | users-core | CF-03 | 15/10/2026 |
| U-06 (B11) | Trả `delivery_status` cho internal-fe | users-core, FE nội bộ | — | 22/10/2026 |
| U-07 (B9) | Đăng ký 9 quyền và resource type | User/Authorization | CF-07 | **15/10/2026** (kéo từ 22/10) |
| U-08 | OAuth2 client credentials cho service | User/Authorization | CF-09 | 15/10/2026 |
| U-09 (C6) | API resolve người nhận cho Notitek | User/Authorization | **CF-11** (mới) | 15/10/2026 |
| U-10 (C7) | API cấp lại link lời mời/kích hoạt | users-core | CF-10 | 22/10/2026 |
| U-11 (D, B6) | Dọn tài liệu và kênh ZALO | users-core | — | 22/10/2026 |

## 2. Chi tiết

### U-01 — Adapter `NotificationPort`

**Lý do:** mọi lần gửi đã đi qua một cổng `NotificationPort.send(Notification)`; chỉ cần thay `LoggingNotificationAdapter`/`MockNotificationAdapter` bằng adapter gọi Notitek.

**Yêu cầu:**

- Gọi bằng OAuth2 client credentials (U-08), header `Idempotency-Key` = `source_ref`, `X-Correlation-Id` lấy từ request gốc.
- JSON snake_case; context `ACCOUNT` hoặc `WORKSPACE` với `application_id`, `organization_id` khi WORKSPACE.
- Không gửi `NotificationChannel`, `template` hay `destination` dạng tự do: kênh do Notitek suy từ loại CONTACT (phone → SMS, email → Email).

**Ánh xạ template → profile:**

| `NotificationTemplates` | Profile | Context | Ánh xạ trường |
|---|---|---|---|
| `OTP_CODE` | `user.otp` (sync) | ACCOUNT | `otp_code` → `secrets.otp`; `purpose` → `data.purpose`; challenge id → `data.challenge_id`; `expiresAt` → `expires_at`; `send_deadline` = thời điểm gọi + 5s (Notitek kiểm `send_deadline ≤ received_at + 5s` và `≤ expires_at`; `occurred_at` là lúc tạo challenge, không bị ràng buộc +5s); destination → CONTACT `phone` hoặc `email` |
| `MEMBER_INVITATION_LINK` | `user.invitation` | WORKSPACE (organization loại SHOP) | `organization_name` → `data.workspace_name`; `inviter_name` → `data.inviter_name`; `link` → `secrets.invitation_url`; `data.invitation_kind=SHOP_MEMBER`; `data.invitation_id` |
| `EMPLOYEE_INVITATION_LINK` | `user.invitation` | WORKSPACE | như trên, `invitation_kind=EMPLOYEE`, không cần `inviter_name` |
| `ACCOUNT_ACTIVATION` | `user.account_activation` | ACCOUNT | `activation_token` **không gửi thô**: User dựng URL HTTPS từ base URL kích hoạt → `secrets.activation_url`; `data.activation_id`; `expires_at` = hạn token (PT24H) |
| `PASSWORD_CHANGED` | `user.password_changed` | ACCOUNT | `changed_via` → `data.changed_via`; **một request** chứa 1–2 CONTACT (phone và/hoặc email đã xác minh), thay vì gọi riêng từng liên hệ như hiện nay |
| `IDENTIFIER_CHANGED` | `user.identifier_changed` | ACCOUNT | `identifier_type`, `masked_new_value`; CONTACT là giá trị **cũ** |
| `EMPLOYEE_OFFBOARDED` | `user.employee_offboarded` | WORKSPACE | `organization_name`, `effective_at` (khai báo thêm vào `NotificationTemplates`); `data.offboarding_id` = job id |

**Ánh xạ kết quả vào `NotificationReceipt`:** giữ `notification_id` (thay `providerReference`) và `sync_result`/`outcome`. Async: HTTP 202 hoặc 200 replay là "đã tiếp nhận", không phải đã gửi.

**Kiểm tra:** contract test với bộ ví dụ `docs/api/examples/*.request.json` của 7 profile User; `scripts/validate_notification_contract.py` PASS với payload adapter sinh ra.

### U-02 (B2) — OTP nhận `UNKNOWN`

**Lý do:** `OtpDeliveryPolicy` hiện ném `OtpDeliveryException(OTP_DELIVERY_UNCONFIRMED)` khi không có `providerReference`, làm luồng thất bại dù tin có thể đã đến người dùng.

**Yêu cầu:**

| `sync_result` | Xử lý của User |
|---|---|
| `PROVIDER_ACCEPTED` | Đã gửi; challenge hiệu lực như hiện nay |
| `FAILED` | Chưa gửi; báo lỗi, cho phép gửi lại ngay theo giới hạn hiện có |
| `UNKNOWN` | **Giữ challenge hiệu lực**; trả trạng thái "chưa xác nhận đã gửi" cho client; chỉ cho gửi lại sau cooldown 3 phút; **không** tính là một lần gửi thất bại, không trừ lượt gửi lại khi chưa có lần gửi mới |

**Kiểm tra:** mô phỏng Notitek trả `UNKNOWN` → mã OTP gốc vẫn xác thực được; nút gửi lại khóa đến `resend_available_at`.

### U-03 (B3) — Không giữ transaction trong lúc gọi Notitek

**Lý do:** gọi OTP đồng bộ trong transaction có thể giữ kết nối DB tới 5s.

**Yêu cầu:** (1) transaction ngắn ghi challenge trạng thái `PENDING_DELIVERY` và commit; (2) gọi Notitek ngoài transaction; (3) transaction ngắn cập nhật kết quả gửi. Crash giữa (1) và (3): challenge `PENDING_DELIVERY` được coi như `UNKNOWN` (U-02), tra lại bằng `GET /v1/notifications?source_ref=…` (nguồn lấy từ token). Gọi lại sau crash vẫn hợp lệ vì ràng buộc 5s tính theo `received_at`, không theo `occurred_at`. OTP qua Email không có provider dự phòng trong R1.

**Kiểm tra:** không có transaction DB nào mở trong thời gian gọi HTTP (log/metric pool).

### U-04 (B4) — `source_ref` tất định và outbox

**Lý do:** chống trùng của Notitek dựa trên `source_ref`; `referenceId` ngẫu nhiên khiến retry tạo tin mới. `NotifyAfterCommit` hiện gửi best-effort và mất tin khi crash.

**Yêu cầu `source_ref`** (≤128 ký tự, `[A-Za-z0-9._:-]`, không chứa SĐT/email/token):

| Ý định | Mẫu |
|---|---|
| OTP | `otp:{challenge_id}:{delivery_version}` |
| Lời mời | `invitation:{invitation_id}:issue:{n}` |
| Kích hoạt | `activation:{activation_id}:issue:{n}` |
| Đổi mật khẩu | `password-changed:{event_id}` |
| Đổi định danh | `identifier-changed:{event_id}` |
| Nghỉ việc | `offboarded:{job_id}` |

**Yêu cầu outbox:** thông báo async (lời mời, kích hoạt, đổi mật khẩu, đổi định danh, nghỉ việc) ghi vào bảng outbox của users-core cùng transaction nghiệp vụ, theo mẫu `authorization_projection_outbox` hiện có; dispatcher gọi Notitek, chỉ đánh dấu đã gửi khi nhận 2xx; timeout/429/5xx retry **nguyên body** cùng `source_ref`; 409 `IDEMPOTENCY_CONFLICT`/4xx/422 vào dead-letter. Payload outbox chứa bí mật (link) phải mã hóa và xóa sau khi gửi được. TTL 10 phút của DEC-09 chỉ áp cho kho bí mật của Notitek, không áp cho outbox của User. Nếu link trong outbox đã cũ hoặc sắp hết hạn trước **lần gửi đầu**, dispatcher của User tự cấp link mới (cấp lại nội bộ) rồi mới gọi Notitek; U-10 chỉ dùng cho retry phía Notitek sau khi yêu cầu đã được tiếp nhận.

**Kiểm tra:** kill process sau commit nghiệp vụ → tin vẫn được gửi đúng một lần; replay cùng body trả cùng `notification_id`.

### U-05 (B8) — `occurred_at`

**Yêu cầu:** luôn là `OffsetDateTime`/`Instant` có offset (mặc định `+07:00`), lấy từ thời điểm nghiệp vụ (tạo challenge, đổi mật khẩu…), không lấy thời điểm dispatcher gọi; giữ nguyên qua các lần retry. Với OTP, dùng lúc tạo challenge là đúng: ràng buộc 5s của Notitek tính theo `received_at`. Không dùng `LocalDateTime`.

**Kiểm tra:** payload thiếu offset bị Notitek từ chối; retry giữ nguyên giá trị.

### U-06 (B11) — `delivery_status` cho internal-fe

**Lý do:** internal-fe đã chờ `delivery_status`, `delivery_channel`, `masked_destination` trên kết quả lời mời/kích hoạt nhưng backend chưa trả.

**Yêu cầu:** lưu `notification_id` cùng lời mời/kích hoạt; khi trả về FE, đọc `GET /v1/notifications/{notification_id}` (cache ngắn) và quy đổi:

| Notitek | `delivery_status` |
|---|---|
| `processing_state` `PENDING`/`PROCESSING`, không delivery nào đang chờ thử lại | `QUEUED` |
| `processing_state` `PROCESSING` và có delivery `attempt_count ≥ 1` với `next_attempt_at` khác null | `RETRYING` |
| `outcome` `ACHIEVED`, `PARTIAL` | `SENT` |
| `outcome` `FAILED`, `EXPIRED`, `CANCELLED`, `SUPPRESSED` | `FAILED` |
| `outcome` `UNKNOWN` | `UNCONFIRMED` (**FE bổ sung giá trị này**, không gộp vào `SENT`) |

**Kiểm tra:** FE hiển thị đủ 5 trạng thái; `UNKNOWN` không hiện "Đã gửi".

### U-07 (B9, CF-07) — Quyền và resource type

**Yêu cầu:** đăng ký 9 quyền DEC-12, gồm `notification.order_alert.receive` và `notification.order_ops_alert.receive`, cùng resource type tương ứng trong `authorization.resource_types`; gán mặc định theo DEC-12. Bài học từ lỗi `INVALID_ACTION_RESOURCE` (thiếu resource type `ORDER`): migration phải đăng ký resource type trước quyền.

**Kiểm tra:** kiểm tra quyền trên Shop thử trả `ALLOW` cho Chủ Shop/CSKH; không trả `INVALID_ACTION_RESOURCE`.

### U-08 (CF-09) — OAuth2 client credentials

**Lý do:** service nội bộ đang dùng `X-Service-Token` tĩnh; A3 chốt OAuth2 client credentials.

**Yêu cầu:** đăng ký client cho Order, User, Notification; token endpoint, issuer, audience; scope theo quy ước `<domain>.<tài_nguyên>.<hành_động>`: API Notitek `notification.request.send`, `notification.request.view`, `notification.request.cancel`; API User do Notitek gọi `user.notification_recipient.resolve`, `user.notification_link.reissue`; xoay vòng secret; cách service xác thực token (JWKS hoặc introspection).

**Kiểm tra:** Notitek xác thực token và ánh xạ client → `source`; token sai audience bị từ chối.

### U-09 (C6, CF-11 mới) — API resolve người nhận

**Lý do:** Notitek cần biến selector thành danh sách người nhận; authorization-service mới có gRPC projection, chưa có API liệt kê.

**Yêu cầu:**

| Selector | Truy vấn |
|---|---|
| `SHOP_ORDER_ALERT_RECEIVERS` | Identity thuộc organization X (loại SHOP), membership còn hiệu lực, có quyền `notification.order_alert.receive` |
| `ORDER_OPS_ALERT_RECEIVERS` | Thành viên tổ chức nội bộ có quyền `notification.order_ops_alert.receive` và data scope bao phủ organization Shop X |

- Gọi bằng client credentials, scope `user.notification_recipient.resolve`.
- Cursor pagination; kết quả là một snapshot nhất quán tại thời điểm truy vấn (kèm `evaluated_at`).
- Chỉ trả `identity_id` và `membership_id`; **không trả liên hệ**.
- Notitek cache ≤ 5 phút và kiểm tra lại quyền trước khi gửi/khi xem hộp tin.

**Kiểm tra:** thành viên rời Shop không còn trong kết quả; nhân viên nội bộ ngoài data scope không được trả về.

### U-10 (C7, CF-10) — API cấp lại link

**Yêu cầu** (theo [README §5](README.md)):

- Đầu vào: `source_ref`, `type` (`user.invitation` hoặc `user.account_activation`), `invitation_id`/`activation_id`, `notification_id`, số lần thử.
- Đầu ra: URL HTTPS mới và `expires_at`.
- Lời mời đã chấp nhận/thu hồi/hết hạn hoặc tài khoản đã kích hoạt → mã từ chối; Notitek kết thúc delivery với `SOURCE_INVALIDATED`.
- Không tính vào giới hạn cấp lại của người dùng (cooldown 60s, tối đa 3 lần).
- **Không vô hiệu link trước**: link cũ vẫn hiệu lực đến hạn của chính nó, link mới là bổ sung (lần gửi trước có thể đã đến người nhận muộn).
- Idempotent theo `(notification_id, attempt)`; xác thực bằng client credentials, scope `user.notification_link.reissue`.

**Kiểm tra:** gọi hai lần cùng `(notification_id, attempt)` trả cùng link; lời mời đã chấp nhận trả mã từ chối; sau khi cấp lại, link cũ vẫn dùng được đến hạn.

### U-11 (D, B6) — Dọn tài liệu và kênh ZALO

- Sửa `services/users-core-service/README.md`: `OtpNotificationPort` → `NotificationPort`.
- B6: R1 không gửi OTP qua Zalo; giữ `ZALO` chỉ ở profile dev, không mở trong môi trường thật. ZNS OTP xem lại ở R2.
- Xóa `OTP_CODE_LENGTH`, `OTP_MAX_RESENDS` không được code đọc khỏi `.env.release.example`, hoặc nối vào cấu hình thật.
