# Notitek - QUYẾT ĐỊNH PO CHO CÁC DEC

**Phiên bản:** 1.1. **Ngày:** 09/10/2026. Bản 1.1: hủy DEC-11 (không có hệ thống cũ), chốt ADR-01 tự xây, loại ADR-07. **Chủ trì:** PO Notification.
**Trạng thái:** PO đã chốt 12 DEC của đợt đầu cùng giá trị khởi điểm cho SRS-P01 đến SRS-P12. Các điểm ghi "Chờ xác nhận" ở mục 3 cần bên sở hữu ký trước khi nghiệm thu phần liên quan; đội được dùng giá trị trong tài liệu này để code, cấu hình và kiểm thử.

Tài liệu này trả lời các câu hỏi tại mục 11 [BRD](<../Notitek - ĐẶC TẢ YÊU CẦU NGHIỆP VỤ (BRD).md>), mục 4 [phạm vi phát hành](<Notitek - PHẠM VI VÀ KẾ HOẠCH PHÁT HÀNH ĐẦU TIÊN.md>) và mục 11 [SRS](<Notitek - ĐẶC TẢ YÊU CẦU PHẦN MỀM (SRS).md>). Quyết định được chốt bởi người được người yêu cầu giao vai trò PO; không thay chữ ký của Tech Lead, pháp chế, tài chính hoặc vận hành ở các điểm thuộc thẩm quyền của họ.

Nguyên tắc chốt: chọn phương án đơn giản nhất đủ cho bốn lát cắt đợt đầu, mọi tham số có giá trị cụ thể và có đường mở rộng. Tên endpoint/event dưới đây là tên được chốt cho hợp đồng; schema chi tiết vẫn do Tech Lead viết thành OpenAPI/JSON Schema.

## 1 Tóm tắt

| DEC | Quyết định chính | Chờ xác nhận |
|---|---|---|
| DEC-01 | Push đợt đầu là Web Push qua FCM trên SuperShip Shop; Notification giữ danh sách đích Push; app mobile sang R2 | Chủ Shop FE |
| DEC-02 | API `/v1/internal/notifications` có `Idempotency-Key`; event Order `order.delivery_failed.v1` qua outbox → broker theo CloudEvents; OAuth2 client credentials | Tech Lead User, Order |
| DEC-03 | Mọi bản tin có `contextType`; người nhận Shop xác định theo quyền `notification.order_alert.receive` | User/Authorization |
| DEC-04 | Ba mục đích SECURITY, TRANSACTIONAL, MARKETING; R1 từ chối MARKETING; có danh sách chặn toàn cục | Pháp chế/DPO |
| DEC-05 | OTP đồng bộ ≤ 5s, ba kết quả `PROVIDER_ACCEPTED`/`FAILED`/`UNKNOWN`; không retry nền | Tech Lead User |
| DEC-06 | Bảng hạn gửi/retry/khung giờ theo luồng; mục tiêu p95 và tải thiết kế | Vận hành (số liệu tải thật) |
| DEC-07 | Ba cấp Global → App/Brand → Shop; cấu hình lưu dạng file trong repo; ghim version, kill-switch có hiệu lực ngay | — |
| DEC-08 | Chủ ngân sách là brand/đơn vị kinh doanh; Notification áp hạn mức cứng theo ngày, không tự giữ sổ | Tài chính, chủ brand |
| DEC-09 | Bảng thời gian lưu từng loại dữ liệu; OTP không lưu ra đĩa; xem dữ liệu gốc cần quyền riêng | Pháp chế/DPO |
| DEC-10 | Bằng chứng và tiêu chí đạt từng kênh/luồng; chỉ fallback SMS khi ZNS lỗi chắc chắn | — |
| DEC-11 | Đã hủy: không có hệ thống cũ cần chuyển đổi | — |
| DEC-12 | Catalog 8 quyền Notification; tách người sửa và người duyệt cấu hình | User/Authorization |

## 2 Nội dung quyết định

### 2.1 DEC-01 — Client Push và phạm vi app

- Push đợt đầu là **Web Push qua Firebase Cloud Messaging (FCM)** trên SuperShip Shop (business-web). Đây là client đã có và là điểm tích hợp ưu tiên; tài liệu hiện có không chứng minh có app mobile. App mobile và app khác đi theo R2.
- **Notification giữ danh sách đích Push.** Frontend gọi API đăng ký đích bằng phiên đăng nhập; tài khoản, app và client lấy từ Access Context, không lấy từ payload.
- Đăng xuất: frontend gọi hủy đăng ký. Token không hoạt động 60 ngày bị loại. Khi FCM trả `UNREGISTERED` hoặc `INVALID_ARGUMENT` cho token, đích bị vô hiệu ngay và không dùng lại.
- Push không bao giờ là kênh bắt buộc, vì FCM không cung cấp bằng chứng đã giao.
- Giao diện nội bộ dùng In-app, không có Push trong R1.

**Lý do:** đủ năm kênh mà không phải chờ app mobile; giữ đích Push tại Notification tránh thêm một hợp đồng thiết bị với User ở đợt đầu. **Ảnh hưởng:** ST-06, SRS-F16, ADR-05, AT-12/26.

### 2.2 DEC-02 — Hợp đồng tích hợp và đường phát

**Quy ước đặt tên** theo code hiện có của nền tảng (User `authorization-service`, Order `order-spf`):

| Loại | Quy ước | Ví dụ đang có |
|---|---|---|
| API cho frontend (phiên người dùng) | `/v1/<tài-nguyên-số-nhiều-kebab>` | `/v1/orders`, `/v1/shipping-labels` |
| API service-to-service và vận hành | `/v1/internal/<tài-nguyên>` | `/v1/internal/orders`, `/v1/internal/roles` |
| Callback từ bên ngoài | `/v1/webhooks/<nguồn>/<loại-event>` | `/v1/webhooks/carrier/events` |
| Hành động trên tài nguyên | Sub-path động từ: `/{id}/<hành-động>` | `/{roleId}/activate`, `/{roleId}/retire` |
| Event | `<domain>.<sự_việc_quá_khứ>.v<major>` | `order.created.v1` |
| Quyền | `<domain>.<tài_nguyên>.<hành_động>` | `user.status.manage` |

Mọi đường dẫn của module bắt đầu bằng `/v1/notifications`, `/v1/internal/notifications` hoặc `/v1/webhooks/notification-providers` để gateway định tuyến theo prefix.

**API được chốt:**

| Thao tác | Bên gọi | Endpoint |
|---|---|---|
| Gửi thông báo (OTP `mode=sync`; lời mời `mode=async`, trả 202) | User, nguồn được phép | `POST /v1/internal/notifications`, header `Idempotency-Key` = reference nguồn |
| Xem một yêu cầu | Nguồn, vận hành | `GET /v1/internal/notifications/{notificationId}` |
| Tra theo reference nguồn | Nguồn, vận hành | `GET /v1/internal/notifications?source=…&sourceRef=…` |
| Hủy một yêu cầu | Nguồn, vận hành | `POST /v1/internal/notifications/{notificationId}/cancel` |
| Hủy theo nhóm liên quan | Nguồn, vận hành | `POST /v1/internal/notifications/cancel`, body có `cancelGroup` |
| Danh sách hộp tin | Frontend | `GET /v1/notifications/inbox` |
| Số chưa đọc | Frontend | `GET /v1/notifications/inbox/unread-count` |
| Tín hiệu realtime (SSE) | Frontend | `GET /v1/notifications/inbox/stream` (`text/event-stream`, hỗ trợ `Last-Event-ID`); xem mục 3.6 tài liệu công nghệ |
| Đánh dấu đã đọc | Frontend | `POST /v1/notifications/inbox/{inboxItemId}/read` |
| Đăng ký đích Push | Frontend | `POST /v1/notifications/push-subscriptions` (token nằm trong body, không đặt trên URL) |
| Hủy đích Push | Frontend | `DELETE /v1/notifications/push-subscriptions/{subscriptionId}` |
| Callback nhà cung cấp | Provider | `POST /v1/webhooks/notification-providers/{provider}/events`, ví dụ `{provider}` = `zalo-zns`, `fcm`, `email`, `sms-<nhà-cung-cấp>` |

`notificationId` là định danh một yêu cầu gửi; `inboxItemId` là định danh bản In-app của một người. Hai định danh không dùng thay nhau.

**Event được chốt:**

| Event | Bên phát | Ý nghĩa |
|---|---|---|
| `order.delivery_failed.v1` | Order | Order xác nhận một lần giao thất bại; partition key `orderId` |
| `order.delivery_failure_resolved.v1` | Order | Lần giao thất bại không còn cần nhắc; Notification ánh xạ thành hủy phần chờ. Cùng partition key để giữ thứ tự |
| `notification.result_changed.v1` | Notification | Kết quả một yêu cầu thay đổi, có `revision`; consumer tự chọn subscribe |

Order phát **sự việc nghiệp vụ**, không phát lệnh dành cho Notification. Nhờ vậy Order không cần biết Notification và các consumer khác dùng lại được cùng event.

- Event dùng envelope **CloudEvents 1.0**; schema viết bằng JSON Schema; trong cùng major version chỉ được thêm trường tùy chọn.
- Xác thực service-to-service bằng **OAuth2 client credentials** với Application/Client của User. Nguồn được xác định từ token; trường nguồn trong payload chỉ để đối chiếu.
- Khóa chống trùng = nguồn + môi trường + reference. Event có `occurredAt` cũ hơn cửa sổ chống trùng bị từ chối với lý do `STALE_EVENT`.
- Sản phẩm broker cụ thể (Kafka, RabbitMQ…) do Tech Lead chốt tại ADR-02, ưu tiên broker nền tảng đã vận hành.

**Lý do:** lệnh cần phản hồi dùng REST; sự việc nghiệp vụ dùng outbox để không mất tin khi nguồn đã commit. **Ảnh hưởng:** ST-02/03/04/11, SRS-F01/02/06/09, ADR-02.

### 2.3 DEC-03 — Người nhận và ngữ cảnh

- Mọi bản tin có `contextType` là `ACCOUNT` hoặc `WORKSPACE`. `WORKSPACE` bắt buộc `shopId` (hoặc `orgId`); `ACCOUNT` không gắn Shop.
- **Người nhận phía Shop xác định theo quyền.** Order chỉ gửi `shopId`; Notification hỏi User/Authorization danh sách thành viên có quyền `notification.order_alert.receive` trên Shop đó tại thời điểm gửi. Kết quả được cache tối đa 5 phút.
- Xử lý nền dùng service account chỉ có quyền đọc membership và quyền trên; không dùng token của người tạo sự việc.
- User/Authorization không trả lời được trong hạn gửi thì bản tin chờ đến hạn rồi dừng với lý do `RECIPIENT_UNRESOLVED`; không gửi theo dữ liệu cũ quá 5 phút.
- Người đã rời Shop không nhận và không xem tin công việc của Shop đó. Cảnh báo bảo mật đi theo tài khoản, không phụ thuộc Shop.
- Gửi tới liên hệ cũ chỉ khi User chỉ định rõ trong yêu cầu và mục đích là SECURITY.
- Người nhận hàng là liên hệ giao dịch do Order cấp; không tạo tài khoản hay hộp In-app cho họ.

**Ảnh hưởng:** ST-04/05/08, SRS-F05/08, ADR-06, AT-08/11/21.

### 2.4 DEC-04 — Mục đích, tính bắt buộc và kênh

| Mục đích | Luồng R1 | Bắt buộc | Người nhận tắt được | Kênh được phép |
|---|---|---|---|---|
| SECURITY | OTP, cảnh báo thiết bị | Có | Không | SMS, Email, In-app |
| TRANSACTIONAL | Lời mời, giao thất bại | Có | R1 chưa hỗ trợ tắt | Email, In-app, Push, ZNS, SMS |
| MARKETING | Không có | — | — | Hệ thống từ chối mọi luồng MARKETING trong R1 |

- Tin ZNS/SMS tới người nhận hàng là tin dịch vụ; cơ sở xử lý dữ liệu là thực hiện hợp đồng giao hàng.
- Có **danh sách chặn toàn cục** cho số/email đã yêu cầu ngừng nhận; áp dụng cho mọi mục đích trừ SECURITY. Bản tin bị chặn có lý do `CONTACT_BLOCKED`.
- R1 không có màn hình cài đặt nhận tin; lựa chọn nhận tin tự phục vụ đi theo phân kỳ sau.
- Nội dung OTP không chứa quảng cáo.

**Ảnh hưởng:** SRS-F03.06/F20, BR-CLS/PREF, AT-24.

### 2.5 DEC-05 — Kết quả gửi OTP

- User gọi `mode=sync`, truyền `sendDeadline` (thời điểm gọi + 5s) và `expiresAt` (hạn của mã).
- Notification trả trong **tối đa 5s** một trong ba kết quả:

| Kết quả | Ý nghĩa | User xử lý |
|---|---|---|
| `PROVIDER_ACCEPTED` | Provider SMS đã nhận | Coi là đã gửi, kích hoạt mã |
| `FAILED` | Lỗi chắc chắn chưa gửi, có mã lý do | Báo lỗi hoặc cho gửi lại |
| `UNKNOWN` | Hết giờ hoặc không xác định | Coi là có thể đã gửi; áp thời gian chờ của User trước khi cho gửi lại |

- **Notification không retry nền cho OTP.** Ngoại lệ duy nhất: thử ngay một lần sang provider SMS dự phòng khi lỗi chắc chắn chưa gửi (lỗi kết nối, provider từ chối trước khi nhận, 5xx có mã xác định) và vẫn trong `sendDeadline`.
- Mỗi lần gửi lại là yêu cầu mới với `deliveryId` mới do User cấp. Delivery report đến muộn chỉ cập nhật bằng chứng, không ảnh hưởng hiệu lực mã.
- Notification giới hạn tối đa 5 OTP mỗi số điện thoại mỗi giờ, bổ sung cho giới hạn hiện có của User; vượt giới hạn trả `FAILED` với lý do `RATE_LIMITED`.
- Adapter phía User ánh xạ ba kết quả vào `NotificationReceipt`; `providerReference` không còn là tín hiệu thành công duy nhất.

**Ảnh hưởng:** ST-02, SRS-F01.05/F07.08/F07.09/F18, SRS-N01, AT-01/04/23.

### 2.6 DEC-06 — Thời gian, retry, khung giờ và hiệu năng

| Luồng/kênh | Hạn gửi | Retry khi lỗi tạm thời | Khung giờ gửi |
|---|---|---|---|
| OTP SMS | 5s | 0, trừ ngoại lệ DEC-05 | 24/7 |
| Email lời mời | 24h hoặc hạn lời mời nếu sớm hơn | 5 lần: 1′, 5′, 15′, 1h, 3h | 24/7 |
| In-app giao thất bại | Tức thì | Retry nội bộ đến khi lưu được | 24/7 |
| Push giao thất bại | 15′ | 3 lần trong 15′ | 24/7 |
| ZNS người nhận hàng | Do Order cấp, mặc định 12h | 2 lần cách 5′ | 07:00–21:00 |
| SMS dự phòng người nhận hàng | Như ZNS | 2 lần cách 5′ | 07:00–21:00 |

- Múi giờ chuẩn `Asia/Ho_Chi_Minh`; mọi thời điểm trong hợp đồng dùng ISO 8601 có offset.
- Ngoài khung giờ, tin được giữ đến 07:00 nếu vẫn còn hạn; hết hạn thì dừng với lý do `EXPIRED`.
- Giới hạn tần suất: tối đa 3 tin ZNS/SMS mỗi số điện thoại mỗi đơn mỗi ngày.

**Mục tiêu hiệu năng (p95):**

| Chỉ số | Mục tiêu |
|---|---|
| OTP: nhận lệnh → provider nhận | ≤ 3s; p99 ≤ 5s |
| Event Order → In-app khả dụng | ≤ 30s |
| Event Order → ZNS provider nhận (trong khung giờ) | ≤ 60s |
| API hộp tin (danh sách, số chưa đọc) | ≤ 300ms |
| API tra cứu vận hành | ≤ 1s trên dữ liệu 180 ngày |

**Tải thiết kế:** 20 OTP/s và 50 event/s ổn định; chịu burst gấp 5 lần trong 1 phút mà không mất yêu cầu đã accepted. OTP được ưu tiên hơn TRANSACTIONAL khi quá tải.

**Ảnh hưởng:** SRS-F06/F07/F21, SRS-N01/N02/N03, AT-04/09/27.

### 2.7 DEC-07 — Thứ tự ưu tiên cấu hình

- R1 có ba cấp: **Global → App/Brand → Shop**. Chiều NVC chưa dùng trong R1.
- **Khóa ở Global:** mục đích, tính bắt buộc, mẫu SECURITY và danh sách kênh được phép của luồng.
- **App/Brand được đặt:** định danh gửi (sender, OA ZNS, domain Email) và mẫu hiển thị theo brand.
- **Shop chỉ được ghi đè:** bật/tắt Push (kênh tùy chọn).
- Cấu hình lưu dạng **file trong repo** (YAML). Pull request là bước duyệt, git là phiên bản và audit. CI kiểm tra schema, biến mẫu và xung đột; hai cấu hình chồng nhau trong cùng cấp bị chặn khi merge.
- Mỗi bản tin ghim version cấu hình tại thời điểm tạo; thay cấu hình không ảnh hưởng tin đang chờ.
- **Kill-switch** dừng một luồng hoặc một kênh có hiệu lực ngay với cả tin đang chờ; bật/tắt kill-switch được audit.

**Ảnh hưởng:** ST-01, SRS-F03/F12, UX-05, AT-18.

### 2.8 DEC-08 — Ngân sách

- **Chủ ngân sách là từng brand/đơn vị kinh doanh**, theo cost center. Notification không giữ sổ tài chính và không tính phí theo Shop trong R1.
- Notification áp **hạn mức cứng theo ngày** cho từng cặp luồng × kênh có phí (SMS, ZNS, Email). Vận hành đặt hạn mức theo đề nghị của chủ ngân sách; hệ thống cảnh báo ở 80%.
- Chạm hạn mức: bản tin nhận trạng thái `BLOCKED_QUOTA`, có lý do, **không chuyển sang kênh có phí khác** và không báo thành công.
- OTP có hạn mức riêng và cao, không bị chặn bởi ngân sách luồng khác; cảnh báo bất thường khi sản lượng giờ vượt 3 lần trung bình 7 ngày.
- Cuối tháng Notification xuất báo cáo số lượng theo cost center, luồng, kênh và provider để tài chính đối soát hóa đơn.

**Ảnh hưởng:** SRS-F14, AT-19.

### 2.9 DEC-09 — Lưu trữ, quyền xem và xóa

| Loại dữ liệu | Thời gian lưu và cách lưu |
|---|---|
| Giá trị OTP/token | Không ghi ra log, trace, dead-letter hay bảng chống trùng. Nếu phải xếp hàng: mã hóa, TTL bằng `expiresAt`, tối đa 10′ |
| Dấu đối chiếu payload có bí mật | HMAC với khóa xoay vòng; không lưu nguyên văn |
| Nội dung đã dựng (không bí mật) | 90 ngày; với OTP chỉ lưu mẫu và biến đã che |
| Metadata gửi, lần thử, kết quả | 180 ngày |
| SĐT/email trong lần gửi | Che sau 30 ngày, chỉ giữ 3 ký tự cuối |
| Tin In-app | Hiển thị 90 ngày, xóa sau 180 ngày |
| Khóa chống trùng | 7 ngày |
| Dấu hủy | 72h |
| Đích Push | Đến khi hủy đăng ký hoặc 60 ngày không hoạt động |
| Audit | 2 năm |
| Backup | 30 ngày; dữ liệu đã xóa không được khôi phục ra ngoài cửa sổ lưu |

- Vận hành mặc định xem dữ liệu đã che. Xem dữ liệu gốc cần quyền `notification.delivery.view_sensitive`, phải nhập lý do và có audit. Xuất dữ liệu là quyền riêng.

**Ảnh hưởng:** SRS-F02.07/F13, SRS-N05/N06, ADR-03, AT-03/28.

### 2.10 DEC-10 — Bằng chứng và tiêu chí đạt

| Kênh | Mức coi là gửi thành công | Bằng chứng giao |
|---|---|---|
| SMS | Provider nhận | Delivery report khi nhà mạng trả |
| Email | Provider nhận | Webhook delivery, bounce, complaint |
| ZNS | Zalo trả `msg_id` | Webhook `user_received_message` |
| Push | FCM nhận | Không có |
| In-app | Đã lưu và người nhận truy cập được | — |

| Luồng | Kênh bắt buộc | Kênh tùy chọn | Kênh thay thế |
|---|---|---|---|
| OTP | SMS | — | — |
| Lời mời | Email | — | — |
| Giao thất bại, phía Shop | In-app | Push | — |
| Giao thất bại, người nhận hàng | ZNS | — | SMS |
| Cảnh báo thiết bị | In-app | — | — |

- **Chỉ fallback SMS khi ZNS lỗi chắc chắn:** số không dùng Zalo, người nhận chặn OA, không liên lạc được hoặc template bị Zalo từ chối cho số đó. **Không fallback khi timeout hoặc chưa rõ.**
- Bounce cứng của Email là thất bại cuối cùng; không retry.
- Hết 24h vẫn chưa rõ thì đóng xử lý với kết quả `UNKNOWN`; không tự đổi thành đã giao hay thất bại.
- Callback trễ chỉ nâng mức bằng chứng, không hạ mức đã đạt.

**Ảnh hưởng:** SRS-F07/F09, mục 8 SRS, ADR-04, AT-09/25.

### 2.11 DEC-11 — Đã hủy

Ngày 09/10/2026 PO xác nhận Notitek là hệ thống mới, không có hệ thống thông báo hoặc ZNS cũ. Không có kiểm kê, feature flag chuyển luồng, shadow hay rollback về hệ thống cũ. OA và template ZNS được đăng ký mới và chờ Zalo duyệt (xem CF-08). Mã DEC-11 giữ lại để không đổi số, không tái sử dụng.

**Ảnh hưởng đã gỡ:** ST-10, SRS-F22, ADR-07, AT-20.

### 2.12 DEC-12 — Danh mục quyền

| Quyền | Phạm vi | Gán mặc định |
|---|---|---|
| Hộp tin của chính mình | Tài khoản | Không cần quyền trong catalog; kiểm tra quan hệ sở hữu và ngữ cảnh |
| `notification.order_alert.receive` | Shop | Chủ Shop, Quản lý, CSKH |
| `notification.delivery.view` | Nền tảng | Vận hành Notification, CS nội bộ |
| `notification.delivery.view_sensitive` | Nền tảng | Trưởng vận hành |
| `notification.delivery.cancel` | Nền tảng | Vận hành Notification |
| `notification.delivery.resend` | Nền tảng | Vận hành Notification |
| `notification.config.edit` / `notification.config.approve` | Nền tảng | Tách hai người: người sửa không tự duyệt |
| `notification.connection.manage` | Nền tảng | Platform admin |
| `notification.usage.export` | Nền tảng | Vận hành và tài chính |

- Backend chỉ cho truy cập khi quyết định là `ALLOW` đúng hành động và tài nguyên; `DENY`, `NOT_EVALUATED` hoặc lỗi kiểm tra đều từ chối.
- R1 dùng PR review cho `config.approve`; khi có màn hình quản trị thì quyền này được kiểm tra trong ứng dụng.

**Ảnh hưởng:** SRS-F05/F08/F11/F12, SRS-N05, ADR-06, AT-16/21.

## 3 Giá trị tham số SRS-P

| Mã | Giá trị được chốt | Nguồn |
|---|---|---|
| SRS-P01 | Phản hồi OTP ≤ 5s; bắt đầu gọi provider ≤ 2s; p95 ≤ 3s, p99 ≤ 5s; tối đa 1% vượt cửa sổ | DEC-05/06 |
| SRS-P02 | Hạn gửi theo bảng mục 2.6; mức bằng chứng theo mục 2.10 | DEC-06/10 |
| SRS-P03 | Mục tiêu độ trễ mục 2.6; 20 OTP/s, 50 event/s, burst ×5 trong 1 phút | DEC-06 |
| SRS-P04 | Tối đa 100 người nhận mỗi yêu cầu; payload ≤ 64KB; mỗi biến ≤ 2KB; trang hộp tin mặc định 20, tối đa 50, dùng cursor | PO + Tech Lead |
| SRS-P05 | Mỗi nguồn 50 yêu cầu/s mặc định; giới hạn provider theo hợp đồng; ưu tiên SECURITY > TRANSACTIONAL; quá tải trả 429 kèm `Retry-After` | DEC-06 |
| SRS-P06 | Retry và fallback theo mục 2.5, 2.6, 2.10; cửa sổ đối chiếu kết quả chưa rõ 24h | DEC-05/06/10 |
| SRS-P07 | Chống trùng 7 ngày; dấu hủy 72h; ID kết quả 180 ngày; event cũ hơn 7 ngày bị từ chối `STALE_EVENT` | DEC-02/09 |
| SRS-P08 | RPO = 0 với yêu cầu đã accepted; RTO 30 phút; khả dụng 99,9%/tháng cho nhận lệnh, tra cứu và hộp tin | PO + Tech Lead |
| SRS-P09 | Thời gian lưu mục 2.9 | DEC-09 |
| SRS-P10 | Hộp tin p95 ≤ 300ms; tra cứu p95 ≤ 1s trên dữ liệu 180 ngày | DEC-06 |
| SRS-P11 | Cache người nhận ≤ 5 phút; hủy đăng ký Push có hiệu lực ≤ 1 phút; mất nguồn quyền thì chờ đến hạn rồi dừng | DEC-01/03 |
| SRS-P12 | Hạn mức ngày theo luồng × kênh, cảnh báo 80%, chặn cứng khi chạm | DEC-08 |

## 4 Điểm chờ bên sở hữu xác nhận

PO đã chốt hướng xử lý; các điểm sau thuộc thẩm quyền bên khác và chặn **nghiệm thu** phần liên quan, không chặn việc code.

| Mã | Nội dung cần ký | Bên ký | Hạn đề xuất | Chặn |
|---|---|---|---|---|
| CF-01 | Thời gian lưu mục 2.9 và cơ sở pháp lý gửi ZNS/SMS người nhận hàng theo Nghị định 13/2023/NĐ-CP | Pháp chế/DPO | 22/10/2026 | Lưu dữ liệu thật, gửi ZNS/SMS thật |
| CF-02 | Cost center và mức hạn mức ngày cụ thể | Tài chính, chủ brand | 22/10/2026 | Gửi thật qua kênh có phí |
| CF-03 | Cơ chế REST/outbox, tên event và kết quả OTP | Tech Lead User, Tech Lead Order | 15/10/2026 | OpenAPI/schema, M1/M2 |
| CF-04 | Đội Order xây dispatcher outbox trong kế hoạch | Tech Lead Order, PO Order | 15/10/2026 | ST-04 và toàn bộ M2 |
| CF-05 | Sản lượng OTP và đơn hiện tại để đối chiếu tải thiết kế | Vận hành | 22/10/2026 | Nghiệm thu tải SRS-N03 |
| CF-06 | Triển khai Web Push trên Shop | Chủ Shop FE | 15/10/2026 | ST-06 |
| CF-07 | Đăng ký 8 quyền vào catalog và gán mặc định | User/Authorization | 22/10/2026 | ST-05/09 |
| CF-08 | Đăng ký OA mới và template ZNS cho luồng giao thất bại, gửi Zalo duyệt (thời gian duyệt tính vào kế hoạch) | Vận hành | 22/10/2026 | ST-07 |

Nếu bên ký đề xuất giá trị khác, PO cập nhật tài liệu này, ghi phiên bản và lý do; không sửa ngầm giá trị đã dùng trong code hoặc kiểm thử.

## 5 Việc tiếp theo

- Tech Lead viết ADR-01 đến ADR-06 dựa trên các quyết định trên (ADR-01 đã chốt tự xây, ADR-07 đã loại); ADR-02 chọn sản phẩm broker.
- Tech Lead viết OpenAPI và JSON Schema theo mục 2.2, kèm danh mục mã lỗi và bộ ví dụ.
- QA cập nhật test case AT với giá trị tại mục 3.
- Vận hành khởi động đăng ký brandname SMS, tài khoản FCM và đăng ký OA/template ZNS mới trong tuần đầu.
