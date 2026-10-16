# Callback nhà cung cấp và API tra cứu vận hành — hợp đồng nháp v1

Ngày 13/10/2026. Trạng thái: **Bản nháp hợp đồng, chưa triển khai**. Chủ hợp đồng: Notitek. Căn cứ: DEC-02, DEC-09, DEC-10, DEC-12 bản 1.9, gồm DB01/DB03/DB06; giải quyết C3 và C5. Quy ước chung theo [API chung](README.md).

## Phần A — Callback nhà cung cấp (C3)

### A1. Endpoint và nguyên tắc

`POST /v1/webhooks/notification-providers/{provider}/events`, với `{provider}` thuộc danh sách đăng ký trong cấu hình: `zalo-zns`, `email`, `sms-<nhà-cung-cấp>`. FCM không có callback giao, nên không có endpoint cho FCM.

- Callback **chỉ cập nhật bằng chứng** của lần thử đã có. Không bao giờ tạo lần gửi, retry hay fallback mới.
- Callback trễ chỉ nâng mức bằng chứng, không hạ mức đã đạt (DEC-10). Bằng chứng bất lợi như hard bounce vẫn được lưu kèm `reason_code`, không bị che.
- Phản hồi 2xx ngay sau khi ghi bền vững callback thô đã xác thực; xử lý cập nhật trạng thái chạy nền. Lỗi xử lý nội bộ không làm provider gửi lại vô hạn.
- Body thô không ghi log; trường liên hệ được che trước khi lưu theo DEC-09.

### A2. Xác thực theo nhà cung cấp

| Nhà cung cấp | Cách xác thực | Ghi chú |
|---|---|---|
| Zalo ZNS | Kiểm tra chữ ký theo tài liệu Zalo OA cho webhook, với secret của OA đăng ký (CF-08) | Sai chữ ký → 401, không lưu |
| Email, SMS | HMAC-SHA256 trên **raw body** và timestamp, header chữ ký do provider quy định; cửa sổ chống phát lại ±5 phút; so sánh constant-time | Nếu provider chỉ hỗ trợ IP allowlist hoặc basic auth, phải được Tech Lead duyệt ngoại lệ trước khi bật |
| FCM | Không có callback | Lỗi token lấy từ phản hồi đồng bộ khi gửi |

Secret mỗi provider lưu trong kho bí mật, hỗ trợ hai secret song song khi xoay vòng. Không dùng secret tĩnh dùng chung như webhook carrier hiện tại của Order.

**Điều kiện tiên quyết:** nhà cung cấp SMS brandname và Email **chưa được chọn**. Định dạng payload, header chữ ký và mã trạng thái của từng provider được bổ sung vào bảng ánh xạ ở A4 sau khi chọn. Endpoint của provider chưa đăng ký trả 404.

### A3. Gắn callback vào lần thử

1. Lấy `provider_message_id` (Zalo: `msg_id`) từ payload và tra lần thử có cùng provider + `provider_message_id`.
2. Khóa chống xử lý trùng = `(provider, provider_event_id)`; provider không có event id thì dùng HMAC của `(provider_message_id, trạng thái, thời điểm provider)`. Callback trùng trả 2xx, không cập nhật lại.
3. **Không tìm thấy lần thử** (callback đến trước khi kịp ghi `provider_message_id`, hoặc ID lạ) thì lưu vào hàng đối chiếu tối đa 24h, thử ghép lại định kỳ; hết hạn thì ghi nhận cho vận hành, không tạo lần gửi.
4. Cập nhật trạng thái delivery theo thứ tự bằng chứng: `PROVIDER_ACCEPTED` < `DELIVERED`; `FAILED` chỉ thay `PROVIDER_ACCEPTED`/`UNKNOWN`, không thay `DELIVERED`. Mỗi thay đổi tăng `revision` của yêu cầu.

### A4. Ánh xạ trạng thái

| Nhà cung cấp | Trạng thái provider | Delivery |
|---|---|---|
| Zalo ZNS | `user_received_message` | `DELIVERED` |
| Zalo ZNS | Số không dùng Zalo, chặn OA, template bị từ chối cho số đó | `FAILED` lỗi chắc chắn; cho phép fallback SMS theo DEC-10 nếu còn hạn, do worker quyết định, không phải callback |
| Email | delivered | `DELIVERED` |
| Email | hard bounce, complaint | `FAILED`, không retry |
| Email | soft bounce | Giữ nguyên, worker retry theo lịch DEC-06 |
| SMS | delivery report thành công | `DELIVERED` |
| SMS | delivery report thất bại cuối | `FAILED` |
| Tất cả | Không có callback sau 24h | Giữ `PROVIDER_ACCEPTED` hoặc `UNKNOWN`, đóng xử lý (DEC-10) |

Mã cụ thể của Email và SMS bổ sung sau khi chọn nhà cung cấp.

## Phần B — API tra cứu vận hành (C5)

### B1. Phân biệt với API của nguồn

| API | Bên gọi | Phạm vi |
|---|---|---|
| `GET /v1/notifications?source_ref=…` hoặc `?subject_ref=…` | Service nguồn, token service | Chỉ yêu cầu của chính nguồn đó; `source` tùy chọn, mặc định lấy từ token, chỉ vận hành được truyền nguồn khác |
| `/v1/notifications/ops/...` | Người vận hành qua internal-fe, context token | Mọi nguồn theo quyền DEC-12 |

`subject_ref` là trường tùy chọn mới của envelope (B10): định danh đối tượng nghiệp vụ mà thông báo nói tới; **bắt buộc với `order.*` và bằng `order_id`**. Router khớp các sub-path cố định `inbox`, `push-subscriptions`, `cancel`, `ops` trước `/{notification_id}`.

### B2. Endpoint

| Method | Endpoint | Quyền | Kết quả |
|---|---|---|---|
| GET | `/v1/notifications/ops/requests` | `notification.delivery.view` | Danh sách yêu cầu theo bộ lọc, cursor |
| GET | `/v1/notifications/ops/requests/{notification_id}` | `notification.delivery.view` | Chi tiết yêu cầu và delivery, liên hệ đã che |
| GET | `/v1/notifications/ops/requests/{notification_id}/attempts` | `notification.delivery.view` | Lịch sử từng lần thử và callback |
| POST | `/v1/notifications/ops/requests/{notification_id}/sensitive-view` | `notification.delivery.view_sensitive` | Xem liên hệ/nội dung gốc; body bắt buộc `reason`; có audit |
| POST | `/v1/notifications/ops/requests/{notification_id}/cancel` | `notification.delivery.cancel` | Hủy phần chưa đưa ra ngoài; `Idempotency-Key`, body `reason_code` |
| POST | `/v1/notifications/ops/requests/{notification_id}/resend` | `notification.delivery.resend` | Gửi lại delivery thất bại còn hạn; `Idempotency-Key`, body `delivery_ids`, `reason` |

### B3. Bộ lọc danh sách

| Query | Quy tắc |
|---|---|
| `source`, `type`, `outcome` | Lọc chính xác; `outcome` nhận nhiều giá trị |
| `subject_ref` | Ví dụ `order_id`; tra toàn bộ thông báo của một đơn |
| `source_ref` | Tra đúng một yêu cầu |
| `recipient` | SĐT/email đầy đủ; server chuẩn hóa rồi so bằng HMAC, không lưu hay log giá trị tìm kiếm; kết quả vẫn che |
| `organization_id` | Lọc theo Shop/tổ chức |
| `created_from`, `created_to` | Bắt buộc, tối đa 31 ngày mỗi lần, trong cửa sổ metadata 180 ngày đến 180 ngày + 1 tháng theo DB01 |
| `page_size`, `cursor` | Mặc định 20, tối đa 50; meta `page_size`, `has_next`, `next_cursor` |

p95 ≤ 1s trên dữ liệu còn giữ theo DEC-09/DB01 (DEC-06, SRS-P10). Bộ lọc bắt buộc khoảng thời gian để giữ mục tiêu này.

### B4. Quy tắc

- Mặc định trả dữ liệu đã che theo DEC-09: SĐT/email chỉ 3 ký tự cuối; OTP và token không bao giờ hiện, kể cả với `view_sensitive`.
- `sensitive-view` ghi audit gồm người xem, lý do, thời điểm, yêu cầu; kết quả không cache ở frontend.
- **Gửi lại** chỉ áp dụng cho delivery `FAILED` hoặc `EXPIRED`, còn trong hạn nghiệp vụ nguồn và không phải OTP; dùng lại bí mật đã lưu nếu còn TTL, lời mời/kích hoạt hết TTL thì xin link mới qua API cấp lại link của User (CF-10). Không gửi lại delivery `UNKNOWN` nếu chưa có bằng chứng chưa gửi.
- Hủy và gửi lại có idempotency theo `(identity vận hành, Idempotency-Key)`; replay trả receipt ban đầu.
- Quyết định quyền khác `ALLOW` đều từ chối; ID ngoài quyền trả 404.

### B4.1. Quyền thay đổi dữ liệu vận hành (DB06)

PO chốt phân quyền ngày 13/10/2026; quyền quản lý kết nối không dùng chung cho mọi thao tác.

| Hành động | Quyền | Gán mặc định |
|---|---|---|
| Tạo/sửa/thử/kích hoạt/thu hồi kết nối | `notification.connection.manage` | Platform admin |
| Sửa hạn mức ngày theo ngân sách được giao | `notification.quota.manage` | Platform admin |
| Dừng/bật gửi theo loại và kênh | `notification.pause.manage` | Vận hành Notification |
| Bật/tắt Push tùy chọn của Shop | `notification.shop_channel.manage` | Vận hành Notification |

Phải kiểm tra ALLOW đúng hành động và tài nguyên/phạm vi (cost center/loại/kênh, phạm vi dừng, hoặc Organization Shop). Mọi thay đổi ghi audit; quyền sửa quota không thay quyết định mức/cost center của chủ ngân sách CF-02. User/Authorization đăng ký quyền và resource type theo CF-07; endpoint/schema thao tác cụ thể do Tech Lead hoàn thiện.

Tra cứu chỉ trả nội dung còn trong thời hạn riêng. Sau 90 ngày từ created_at, title/preview inbox đã NULL; metadata vẫn tra theo DB01/DB02, không dựng lại nội dung đã dọn.

### B5. Câu hỏi mở

| Mã | Câu hỏi | Bên chốt |
|---|---|---|
| Q1 | Chọn nhà cung cấp SMS brandname và Email; tài liệu chữ ký webhook của họ | Vận hành + Tech Lead |
| Q2 | Thời gian giữ hàng đối chiếu callback không khớp (đề xuất 24h) | Tech Lead |
| Q3 | Có cho phép tìm theo SĐT/email đầy đủ không, hay chỉ theo `subject_ref`/`source_ref` | Pháp chế/DPO |
| Q4 | Màn hình tra cứu nằm trong internal-fe hiện có hay ứng dụng vận hành riêng | PO + FE nội bộ |
