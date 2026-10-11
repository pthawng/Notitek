# Hộp tin, SSE và đăng ký Push — hợp đồng nháp v1

Ngày 11/10/2026. Trạng thái: **Bản nháp hợp đồng, chưa triển khai**. Chủ hợp đồng: Notitek. Căn cứ: DEC-01, DEC-02 (danh sách endpoint), DEC-03, DEC-06, DEC-09, DEC-12 bản 1.6; giải quyết C1 và C2. Quy ước chung (envelope, snake_case, `X-Correlation-Id`) theo [API chung](README.md).

## 1. Phạm vi và chủ thể

API này phục vụ **frontend với phiên người dùng**, khác API ra lệnh gửi dành cho service. Hai nhóm không dùng chung token hay quyền.

| Nhóm | Bên gọi | Token | Không được dùng |
|---|---|---|---|
| Hộp tin, SSE | shop-fe (business-web), internal-fe | Context token của User | Token service, token API của tổ chức |
| Đăng ký Push | business-web (R1) | Context token của User | Token service |

Mỗi request gửi `Authorization: Bearer <context token>`. Notitek resolve qua User `POST /v1/access-context/resolve` để lấy identity, membership, `organization_id`, `application_id` và client; kết quả cache theo thời hạn token, tối đa 60s. Payload và query không được tự khai identity, organization hay application. Token không resolve được trả 401; resolve trả `DENY` hoặc `NOT_EVALUATED` trả 403.

## 2. Endpoint

| Method | Endpoint | Kết quả |
|---|---|---|
| GET | `/v1/notifications/inbox` | Danh sách bản In-app trong ngữ cảnh hiện tại, cursor |
| GET | `/v1/notifications/inbox/unread-count` | Số chưa đọc cùng phạm vi với danh sách |
| GET | `/v1/notifications/inbox/stream` | SSE `text/event-stream`, chỉ báo tín hiệu |
| POST | `/v1/notifications/inbox/{inbox_item_id}/read` | Đánh dấu đã đọc, idempotent |
| POST | `/v1/notifications/push-subscriptions` | Đăng ký hoặc làm mới đích Web Push |
| DELETE | `/v1/notifications/push-subscriptions/{subscription_id}` | Hủy đích Push |

Router khớp các sub-path cố định `inbox`, `push-subscriptions`, `cancel`, `ops` trước `/{notification_id}`. Gateway tách route frontend và route service theo loại token.

## 3. Phạm vi hiển thị

- **Ngữ cảnh ACCOUNT:** chỉ bản tin gắn với chính identity, như cảnh báo thiết bị, đổi mật khẩu.
- **Ngữ cảnh WORKSPACE (bản tin Shop):** bản tin của identity trong đúng `organization_id` của context token. Đổi workspace trên frontend là lấy token mới; danh sách đổi theo.
- **Bản tin vận hành** (gửi qua `ORDER_OPS_ALERT_RECEIVERS`, context là Shop của đơn): hiển thị cho người dùng nội bộ khi đồng thời (1) họ là người nhận đã resolve của bản tin, (2) hiện còn quyền `notification.order_ops_alert.receive`, (3) data scope hiện tại còn bao phủ organization Shop của bản tin. **Không lọc theo `organization_id` của token** (token mang tổ chức nội bộ, không phải Shop).
- **Kiểm tra quyền lúc đọc, không chỉ lúc gửi.** Bản tin Shop yêu cầu quyền `notification.order_alert.receive` hiện tại trên Shop đó; bản tin vận hành yêu cầu `notification.order_ops_alert.receive` và data scope vẫn bao phủ Shop của đơn. Thành viên đã rời Shop hoặc bị gỡ quyền không thấy bản tin của Shop, kể cả bản đã nhận trước đó.
- Danh sách mặc định gồm bản ACCOUNT của identity, bản WORKSPACE của organization hiện tại và bản tin vận hành đạt điều kiện ở trên. Query `scope=ACCOUNT|WORKSPACE` lọc thêm, không mở rộng.
- `inbox_item_id` không thuộc người gọi hoặc ngoài phạm vi trả 404, không lộ sự tồn tại.
- Thời gian theo DEC-09: hiển thị 90 ngày, xóa sau 180 ngày.

Người dùng internal-fe có quyền vận hành thấy bản tin `ORDER_OPS_ALERT_RECEIVERS` theo quy tắc bản tin vận hành ở trên, trong ngữ cảnh tổ chức nội bộ; không cần chọn Shop. Bị gỡ quyền hoặc data scope thu hẹp thì bản tin của Shop ngoài phạm vi biến mất khỏi danh sách và số chưa đọc.

## 4. Danh sách

`GET /v1/notifications/inbox?page_size=20&cursor=...&scope=...&unread_only=true`

| Query | Quy tắc |
|---|---|
| `page_size` | Mặc định 20, tối đa 50 (SRS-P04); vượt → 400 |
| `cursor` | Opaque, ký bởi server, gắn với identity và ngữ cảnh; sai hoặc của người khác → 400 `INVALID_CURSOR` |
| `scope` | Tùy chọn `ACCOUNT` hoặc `WORKSPACE` |
| `unread_only` | Tùy chọn, mặc định `false` |

Thứ tự `created_at` giảm dần, phụ `inbox_item_id`. Meta theo cursor của nền tảng: `page_size`, `has_next`, `next_cursor`. Không trả tổng số.

| Trường item | Ý nghĩa |
|---|---|
| `inbox_item_id` | Định danh bản In-app của một người; khác `notification_id` |
| `notification_id` | Yêu cầu gửi gốc, để đối chiếu |
| `type` | Type nguồn, ví dụ `order.delivery_failed`, `user.device_alert` |
| `title`, `body` | Nội dung đã dựng theo mẫu và `config_version`; không chứa bí mật, liên hệ đầy đủ hay token |
| `action` | `{kind, target_ref}` trỏ về chức năng của module nguồn, ví dụ `kind=ORDER_DETAIL`, `target_ref=order_code`; `kind=DEVICE_TRUST_REVIEW`, `target_ref=alert_id`. Frontend điều hướng; module nguồn kiểm tra quyền và trạng thái nghiệp vụ |
| `context` | `{context_type, organization_id?}` của bản tin |
| `created_at` | Thời điểm bản In-app khả dụng |
| `read_at` | `null` nếu chưa đọc |
| `visible_until` | `created_at` + 90 ngày |

Notitek không lưu và không trả trạng thái nghiệp vụ, ví dụ thiết bị đã được tin cậy hay đơn đã giao lại. Frontend lấy trạng thái đó từ User/Order khi mở hành động.

## 5. Số chưa đọc và đánh dấu đọc

`GET /v1/notifications/inbox/unread-count` trả `{unread_count, as_of}` với đúng tập lọc phạm vi và quyền ở mục 3; nhận cùng `scope`. Tối đa trả 99 kèm `capped=true` để giữ p95.

`POST /v1/notifications/inbox/{inbox_item_id}/read` không cần body và không cần `Idempotency-Key`. Gọi lại khi đã đọc trả 200 với `read_at` ban đầu. Trạng thái đọc lưu ở backend theo người; frontend bỏ lưu `sessionStorage`. Đọc không hoàn thành nghiệp vụ nguồn và không thay kết quả gửi. Đánh dấu đọc hàng loạt chưa có trong v1 (câu hỏi mở Q3).

## 6. SSE

`GET /v1/notifications/inbox/stream`, `Accept: text/event-stream`.

- **Chỉ là tín hiệu.** Event `inbox.changed` với `data: {"unread_count": n, "latest_created_at": "..."}`; không gửi nội dung tin. Frontend nhận tín hiệu thì gọi lại danh sách và số chưa đọc.
- `id` của event tăng đơn điệu theo người. Kết nối lại gửi `Last-Event-ID`; server phát một tín hiệu tổng hợp nếu có thay đổi sau id đó, không phát lại từng event.
- Heartbeat comment mỗi 25s. Server đóng kết nối khi context token hết hạn; frontend lấy token mới rồi kết nối lại với backoff 1s → 30s.
- Mỗi người tối đa 5 kết nối đồng thời; vượt thì đóng kết nối cũ nhất.
- **Dự phòng:** SSE lỗi liên tục hoặc bị proxy chặn thì frontend poll `unread-count` mỗi 30s. Tín hiệu mất không làm mất tin vì danh sách là nguồn sự thật.

## 7. Thay thế cơ chế hiện có của frontend

| Hiện có | Sau khi tích hợp |
|---|---|
| shop-fe và internal-fe poll `GET /v1/me/device-trust-requests` mỗi 10s cho `NotificationBell` | Chuông đọc hộp tin Notitek và SSE; cảnh báo thiết bị đến từ profile `user.device_alert` |
| `NotificationItem.kind` = `NEW_DEVICE_LOGIN` hoặc `SYSTEM`, dữ liệu `MOCK_NOTIFICATIONS` | Dùng `type` và `action` của item; adapter frontend ánh xạ `user.device_alert` → hiển thị cảnh báo thiết bị hiện có |
| Trạng thái đọc lưu `sessionStorage` | `read_at` từ backend, đồng bộ giữa thiết bị |
| Phê duyệt thiết bị gọi User | Giữ nguyên: hành động `DEVICE_TRUST_REVIEW` mở màn hình User; Notitek không xử lý phê duyệt |

## 8. Đăng ký Push (Web Push qua FCM)

R1 chỉ có business-web (DEC-01). internal-fe không dùng Push.

`POST /v1/notifications/push-subscriptions`

| Trường body | Quy tắc |
|---|---|
| `fcm_token` | Bắt buộc, chỉ nằm trong body, không đặt trên URL hay log |
| `client_label` | Tùy chọn, tối đa 100 ký tự, ví dụ trình duyệt; chỉ để người dùng nhận diện |

Identity, `application_id` và client lấy từ Access Context. Phản hồi 201 khi tạo mới, 200 khi làm mới đích sẵn có của cùng người–app, kèm `{subscription_id, created_at, last_seen_at}`; không trả lại token.

- **Đổi chủ token:** token đã gắn với identity khác thì chuyển sang identity hiện tại và thu hồi ràng buộc cũ ngay; không gửi Push của người cũ tới thiết bị này nữa.
- **Đăng xuất:** frontend gọi `DELETE /v1/notifications/push-subscriptions/{subscription_id}`; có hiệu lực ≤ 1 phút (SRS-P11). Xóa đích không thuộc người gọi trả 404. Xóa lại trả 204.
- **Không hoạt động 60 ngày** thì loại đích.
- FCM trả `UNREGISTERED` hoặc `INVALID_ARGUMENT` cho token thì vô hiệu ngay và không dùng lại; đăng ký lại phải qua POST mới.
- Mỗi identity–app tối đa 10 đích; vượt thì thay đích cũ nhất.
- Push **không bao giờ là kênh bắt buộc**; Shop chỉ được tắt Push (DEC-07). Nội dung Push không chứa bí mật, liên hệ đầy đủ hay token.

## 9. HTTP và mã lỗi

Envelope `{error, message, data}`; lỗi có `data.code`, `data.trace_id`, `data.details[]`. Mọi response có `Cache-Control: no-store` và trả lại `X-Correlation-Id` (đã khai báo trong OpenAPI cho mọi response).

| HTTP | Mã | Khi nào |
|---|---|---|
| 400 | `INVALID_REQUEST`, `INVALID_CURSOR` | Query/body sai, cursor sai hoặc không thuộc người gọi |
| 401 | `UNAUTHENTICATED` | Thiếu hoặc hết hạn context token |
| 403 | `FORBIDDEN` | Access Context từ chối; dùng token service gọi route frontend |
| 404 | `NOT_FOUND` | Item hoặc subscription không có hoặc ngoài phạm vi |
| 429 | `RATE_LIMITED` | Vượt giới hạn theo người; kèm `Retry-After` |
| 503 | `SERVICE_UNAVAILABLE` | Không resolve được Access Context hoặc quyền; không trả dữ liệu cũ thay thế |

## 10. Mục tiêu và kiểm chứng

- p95 ≤ 300ms cho danh sách và số chưa đọc (DEC-06, SRS-P10).
- Kiểm chứng tối thiểu: thành viên rời Shop không còn thấy bản tin Shop; đổi workspace đổi danh sách; ID của người khác trả 404; đọc lại idempotent; SSE mất kết nối vẫn hội tụ nhờ danh sách; token Push đổi chủ không nhận tin của người cũ; FCM `UNREGISTERED` vô hiệu đích ngay.

## 11. Câu hỏi mở

| Mã | Câu hỏi | Bên chốt |
|---|---|---|
| Q1 | Cache kết quả resolve Access Context tối đa 60s có chấp nhận được với yêu cầu gỡ quyền có hiệu lực ngay không | User/Authorization |
| Q2 | Danh mục `action.kind` và `target_ref` cho 17 profile; route frontend tương ứng ở shop-fe/internal-fe | PO + Shop FE + FE nội bộ |
| Q3 | Có cần "đánh dấu tất cả đã đọc" trong R1 không | PO |
| Q4 | Giới hạn 5 kết nối SSE và 10 đích Push mỗi người | Tech Lead + vận hành |
| Q5 | Gateway hiện có hỗ trợ giữ kết nối SSE dài và tắt buffer không | Platform |
