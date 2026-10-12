# Notitek - QUYẾT ĐỊNH PO CHO CÁC DEC

**Phiên bản:** 1.8. **Ngày:** 12/10/2026. Bản 1.8: chốt cách xử lý bounce sau `PROVIDER_ACCEPTED` và nghĩa của `outcome` khi chưa đóng xử lý (mục 2.10, 2.15 S6 đến S8). Bản 1.7: DEC-07/DEC-08 chuyển kết nối kênh (sender, OA, domain, credential), hạn mức ngày, bật/tắt Push của Shop và kill-switch sang DB và API để vận hành đổi được trên production mà không cần deploy; chính sách, mẫu, tiêu chí đạt giữ ở YAML. Bản 1.6: áp dụng kết quả review (mục 2.15): `source` mặc định theo token, ràng buộc OTP theo `received_at`, `source_ref` Order phân biệt từng lần thất bại, `attempt_count`/`next_attempt_at`, liên hệ người nhận chỉ bắt buộc ở giao thất bại, cấp lại link không vô hiệu link cũ, quy ước scope OAuth, chốt kênh 10 mốc Order, mã UC-ORD-12..16, hộp tin cho vận hành. Bản 1.5: chốt hướng xử lý B2–B11, C1–C7 và D (mục 2.14): `subject_ref` và liệt kê theo đối tượng, xử lý OTP `UNKNOWN`, reference tất định của User, ánh xạ `delivery_status`, R1 chỉ tiếng Việt và không có OTP Zalo, thêm CF-11 (API resolve người nhận) và CF-12 (chọn nhà cung cấp), kéo CF-07 về 15/10. Bản 1.4: chốt A1–A8 và B1 (mục 2.13): R1 nguồn gọi REST từ dispatcher outbox, broker để sau; Order gửi liên hệ người nhận, Notification không gọi ngược Order; OAuth2 client credentials; snake_case và `X-Correlation-Id`; context `organization_id` + `application_id`; R1 gồm 4 thông báo tài khoản User và 10 mốc Order; thêm audience/quyền vận hành; User cung cấp API cấp lại link. Bản 1.3: giới hạn tần suất nghiệp vụ (OTP, tin theo đơn) thuộc module nguồn; Notification chỉ giữ giới hạn bảo vệ hệ thống/provider; OTP nhận đích SMS hoặc Email kèm `purpose`, lời mời nhận Email hoặc SMS và có lời mời nhân viên, thêm bốn profile User (kích hoạt tài khoản, đổi mật khẩu, đổi định danh, nghỉ việc). Bản 1.2: người yêu cầu chốt endpoint chung `/v1/notifications`, thay prefix internal trong DEC-02. Bản 1.1: hủy DEC-11 (không có hệ thống cũ), chốt ADR-01 tự xây, loại ADR-07. **Chủ trì:** PO Notification.
**Trạng thái:** PO đã chốt 12 DEC của đợt đầu cùng giá trị khởi điểm cho SRS-P01 đến SRS-P12. Các điểm ghi "Chờ xác nhận" ở mục 4 cần bên sở hữu ký trước khi nghiệm thu phần liên quan; đội được dùng giá trị trong tài liệu này để code, cấu hình và kiểm thử.

Tài liệu này trả lời các câu hỏi tại mục 11 [BRD](<Notitek - ĐẶC TẢ YÊU CẦU NGHIỆP VỤ (BRD).md>), mục 4 [phạm vi phát hành](<Notitek - PHẠM VI VÀ KẾ HOẠCH PHÁT HÀNH ĐẦU TIÊN.md>) và mục 11 [SRS](<Notitek - ĐẶC TẢ YÊU CẦU PHẦN MỀM (SRS).md>). Quyết định được chốt bởi người được người yêu cầu giao vai trò PO; không thay chữ ký của Tech Lead, pháp chế, tài chính hoặc vận hành ở các điểm thuộc thẩm quyền của họ.

Nguyên tắc chốt: chọn phương án đơn giản nhất đủ cho phạm vi R1 (OTP, lời mời, bốn thông báo tài khoản User, cảnh báo thiết bị và 10 mốc Order), mọi tham số có giá trị cụ thể và có đường mở rộng. Tên endpoint/event dưới đây là tên được chốt cho hợp đồng; schema chi tiết vẫn do Tech Lead viết thành OpenAPI/JSON Schema.

## 1 Tóm tắt

| DEC | Quyết định chính | Chờ xác nhận |
|---|---|---|
| DEC-01 | Push đợt đầu là Web Push qua FCM trên SuperShip Shop; Notification giữ danh sách đích Push; app mobile sang R2 | Chủ Shop FE |
| DEC-02 | API `/v1/notifications` có `Idempotency-Key`, JSON snake_case, header `X-Correlation-Id`; R1 nguồn gọi REST (async từ dispatcher outbox của nguồn), broker để sau; OAuth2 client credentials | Tech Lead User, Order; User/Authorization (CF-09) |
| DEC-03 | Mọi bản tin có `context_type`, `application_id` và `organization_id` khi WORKSPACE; người nhận Shop theo quyền `notification.order_alert.receive`, vận hành theo `notification.order_ops_alert.receive`; Order gửi liên hệ người nhận hàng | User/Authorization |
| DEC-04 | Ba mục đích SECURITY, TRANSACTIONAL, MARKETING; R1 từ chối MARKETING; có danh sách chặn toàn cục | Pháp chế/DPO |
| DEC-05 | OTP đồng bộ ≤ 5s, ba kết quả `PROVIDER_ACCEPTED`/`FAILED`/`UNKNOWN`; không retry nền | Tech Lead User |
| DEC-06 | Bảng hạn gửi/retry/khung giờ theo luồng; mục tiêu p95 và tải thiết kế | Vận hành (số liệu tải thật) |
| DEC-07 | Ba cấp Global → App/Brand → Shop; chính sách (loại thông báo, mẫu, tiêu chí đạt) ở file YAML trong repo, dữ liệu vận hành (kết nối kênh, kill-switch, bật/tắt Push của Shop) ở DB và API; ghim version YAML, kill-switch có hiệu lực ngay | — |
| DEC-08 | Chủ ngân sách là brand/đơn vị kinh doanh; Notification áp hạn mức cứng theo ngày (đặt qua API, lưu DB), không tự giữ sổ | Tài chính, chủ brand |
| DEC-09 | Bảng thời gian lưu từng loại dữ liệu; OTP không lưu ra đĩa; xem dữ liệu gốc cần quyền riêng | Pháp chế/DPO |
| DEC-10 | Bằng chứng và tiêu chí đạt từng kênh/luồng; chỉ fallback SMS khi ZNS lỗi chắc chắn | — |
| DEC-11 | Đã hủy: không có hệ thống cũ cần chuyển đổi | — |
| DEC-12 | Catalog 9 quyền Notification; tách người sửa và người duyệt cấu hình | User/Authorization |

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
| Trường JSON và query param | snake_case | `order_code`, `page_size`, `trace_id` |
| Correlation | Header `X-Correlation-Id` | Order `RequestCorrelationFilter`, header actor của User |

Riêng Notitek, người yêu cầu chốt ngày 11/10/2026: API chung dùng `/v1/notifications` cho cả nguồn tích hợp và frontend; callback dùng `/v1/webhooks/notification-providers`. Đây là ngoại lệ được chốt so với quy ước prefix service-to-service ở bảng trên. Bỏ prefix internal không bỏ xác thực: gửi/tra cứu/hủy yêu cầu vẫn dùng danh tính service hoặc quyền vận hành; hộp tin dùng phiên người dùng. Gateway phân quyền theo method/route, không suy quyền chỉ từ prefix.

**API được chốt:**

| Thao tác | Bên gọi | Endpoint |
|---|---|---|
| Gửi thông báo (OTP `mode=sync`; lời mời `mode=async`, trả 202) | User, nguồn được phép | `POST /v1/notifications`, header `Idempotency-Key` = reference nguồn |
| Xem một yêu cầu | Nguồn, vận hành | `GET /v1/notifications/{notification_id}` |
| Tra theo reference nguồn | Nguồn, vận hành | `GET /v1/notifications?source=…&source_ref=…` |
| Liệt kê theo đối tượng nghiệp vụ | Nguồn, vận hành | `GET /v1/notifications?subject_ref=…&cursor=…&page_size=…` |
| Tra cứu vận hành | Vận hành | `/v1/notifications/ops/…`, xem [bản nháp](api/provider-webhooks-and-ops.md) |
| Hủy một yêu cầu | Nguồn, vận hành | `POST /v1/notifications/{notification_id}/cancel` |
| Hủy theo nhóm liên quan | Nguồn, vận hành | `POST /v1/notifications/cancel`, body có `cancel_group` |
| Danh sách hộp tin | Frontend | `GET /v1/notifications/inbox` |
| Số chưa đọc | Frontend | `GET /v1/notifications/inbox/unread-count` |
| Tín hiệu realtime (SSE) | Frontend | `GET /v1/notifications/inbox/stream` (`text/event-stream`, hỗ trợ `Last-Event-ID`); xem mục 3.6 tài liệu công nghệ |
| Đánh dấu đã đọc | Frontend | `POST /v1/notifications/inbox/{inbox_item_id}/read` |
| Đăng ký đích Push | Frontend | `POST /v1/notifications/push-subscriptions` (token nằm trong body, không đặt trên URL) |
| Hủy đích Push | Frontend | `DELETE /v1/notifications/push-subscriptions/{subscription_id}` |
| Callback nhà cung cấp | Provider | `POST /v1/webhooks/notification-providers/{provider}/events`, ví dụ `{provider}` = `zalo-zns`, `fcm`, `email`, `sms-<nhà-cung-cấp>` |

`notification_id` là định danh một yêu cầu gửi; `inbox_item_id` là định danh bản In-app của một người. Hai định danh không dùng thay nhau.

Thiết kế cụ thể cho API chung, schema theo loại thông báo và ví dụ nằm ở [hợp đồng API v1](api/README.md) và [OpenAPI](api/openapi.json). Router ưu tiên các sub-path cố định `inbox`, `push-subscriptions`, `cancel`, `ops` trước `/{notification_id}`.

**Đường tích hợp R1 (chốt 11/10/2026, A2):** mọi nguồn gọi REST. OTP gọi đồng bộ. Thông báo async được nguồn ghi vào outbox của chính nguồn cùng transaction nghiệp vụ, rồi **dispatcher outbox của nguồn** gọi `POST /v1/notifications` (hủy: `POST /v1/notifications/cancel`) với reference ổn định. Dispatcher chỉ đánh dấu đã gửi khi nhận 2xx; timeout/429/5xx thì retry nguyên body cùng reference; 4xx hợp đồng vào dead-letter.

- Order dùng profile `order.*` cho 10 mốc (mục 2.13); lần thất bại được giải quyết thì gọi hủy theo `cancel_group`, thay cho event `order.delivery_failure_resolved.v1`.
- Xác thực service-to-service bằng **OAuth2 client credentials** với Application/Client của User (A3). Nguồn được xác định từ token; payload không được tự khai nguồn. User/Authorization phải bổ sung cấp client credentials cho service (CF-09).
- Khóa chống trùng = nguồn + môi trường + reference. Yêu cầu có `occurred_at` cũ hơn cửa sổ chống trùng bị từ chối với lý do `STALE_REQUEST`.
- Consumer theo dõi kết quả bằng `GET` trong R1.

**Phân kỳ sau (broker):** event nghiệp vụ theo **CloudEvents 1.0**, JSON Schema, partition key `order_id`; `notification.result_changed.v1` phát cho consumer subscribe. Sản phẩm broker do Tech Lead chốt tại ADR-02 khi bắt đầu phân kỳ này. Reference qua broker phải trùng reference REST của cùng ý định.

**Lý do:** Order chưa có broker và dispatcher; REST từ outbox giữ được đảm bảo không mất tin sau commit mà không thêm hạ tầng trong R1. **Ảnh hưởng:** ST-02/03/04/11, SRS-F01/02/06/09, ADR-02.

### 2.3 DEC-03 — Người nhận và ngữ cảnh

- Mọi bản tin có `context_type` là `ACCOUNT` hoặc `WORKSPACE` và `application_id`. `WORKSPACE` bắt buộc `organization_id` (Shop là Organization loại SHOP; tổ chức nội bộ là CORPORATION…); `ACCOUNT` không gắn organization. Loại organization yêu cầu theo từng profile được kiểm tra với User.
- **Người nhận phía Shop xác định theo quyền.** Order gửi `organization_id` của Shop và selector `SHOP_ORDER_ALERT_RECEIVERS`; Notification hỏi User/Authorization danh sách thành viên có quyền `notification.order_alert.receive` trên Shop đó tại thời điểm gửi.
- **Người nhận vận hành/CSKH** (A8): selector `ORDER_OPS_ALERT_RECEIVERS` resolve thành viên tổ chức nội bộ có quyền `notification.order_ops_alert.receive` và data scope bao phủ Shop của đơn. Kết quả được cache tối đa 5 phút.
- Xử lý nền dùng service account chỉ có quyền đọc membership và quyền trên; không dùng token của người tạo sự việc.
- User/Authorization không trả lời được trong hạn gửi thì bản tin chờ đến hạn rồi dừng với lý do `RECIPIENT_UNRESOLVED`; không gửi theo dữ liệu cũ quá 5 phút.
- Người đã rời Shop không nhận và không xem tin công việc của Shop đó. Cảnh báo bảo mật đi theo tài khoản, không phụ thuộc Shop.
- Gửi tới liên hệ cũ chỉ khi User chỉ định rõ trong yêu cầu và mục đích là SECURITY.
- Người nhận hàng là liên hệ giao dịch do **Order gửi kèm trong yêu cầu** (A1): bắt buộc ở giao thất bại, tùy chọn ở 9 mốc còn lại và Order bỏ đi khi mốc không nhắn người nhận (Q2); Notification không gọi ngược sang Order để tra liên hệ. Không tạo tài khoản hay hộp In-app cho người nhận hàng; mốc không nhắn người nhận thì không lưu liên hệ.

**Ảnh hưởng:** ST-04/05/08, SRS-F05/08, ADR-06, AT-08/11/21.

### 2.4 DEC-04 — Mục đích, tính bắt buộc và kênh

| Mục đích | Luồng R1 | Bắt buộc | Người nhận tắt được | Kênh được phép |
|---|---|---|---|---|
| SECURITY | OTP, cảnh báo thiết bị, đổi mật khẩu, đổi định danh | Có | Không | SMS, Email, In-app |
| TRANSACTIONAL | Lời mời, kích hoạt tài khoản, nghỉ việc, 10 mốc Order | Có | R1 chưa hỗ trợ tắt | Email, In-app, Push, ZNS, SMS |
| MARKETING | Không có | — | — | Hệ thống từ chối mọi luồng MARKETING trong R1 |

- Tin ZNS/SMS tới người nhận hàng là tin dịch vụ; cơ sở xử lý dữ liệu là thực hiện hợp đồng giao hàng.
- Có **danh sách chặn toàn cục** cho số/email đã yêu cầu ngừng nhận; áp dụng cho mọi mục đích trừ SECURITY. Bản tin bị chặn có lý do `CONTACT_BLOCKED`.
- R1 không có màn hình cài đặt nhận tin; lựa chọn nhận tin tự phục vụ đi theo phân kỳ sau.
- Nội dung OTP không chứa quảng cáo.

**Ảnh hưởng:** SRS-F03.06/F20, BR-CLS/PREF, AT-24.

### 2.5 DEC-05 — Kết quả gửi OTP

- User gọi `mode=sync`, truyền `send_deadline` (thời điểm gọi + 5s) và `expires_at` (hạn của mã).
- Notification trả trong **tối đa 5s** một trong ba kết quả:

| Kết quả | Ý nghĩa | User xử lý |
|---|---|---|
| `PROVIDER_ACCEPTED` | Provider SMS/Email đã nhận | Coi là đã gửi, kích hoạt mã |
| `FAILED` | Lỗi chắc chắn chưa gửi, có mã lý do | Báo lỗi hoặc cho gửi lại |
| `UNKNOWN` | Hết giờ hoặc không xác định | Coi là có thể đã gửi; áp thời gian chờ của User trước khi cho gửi lại |

- **Notification không retry nền cho OTP.** Ngoại lệ duy nhất: thử ngay một lần sang provider SMS dự phòng (OTP qua Email không có dự phòng trong R1) khi lỗi chắc chắn chưa gửi (lỗi kết nối, provider từ chối trước khi nhận, 5xx có mã xác định) và vẫn trong `send_deadline`.
- Mỗi lần gửi lại là yêu cầu mới với `source_ref` mới do User cấp. Delivery report đến muộn chỉ cập nhật bằng chứng, không ảnh hưởng hiệu lực mã.
- Giới hạn tần suất OTP (cooldown, số lần gửi lại, đổi kênh, số lần mỗi ngày theo đích/mục đích) do **User sở hữu và thực thi** trước khi gọi Notification. Notification không áp thêm giới hạn nghiệp vụ theo số điện thoại; chỉ áp giới hạn bảo vệ theo SRS-P05 và giới hạn provider.
- Adapter phía User ánh xạ ba kết quả vào `NotificationReceipt`; `providerReference` không còn là tín hiệu thành công duy nhất.

**Ảnh hưởng:** ST-02, SRS-F01.05/F07.08/F07.09/F18, SRS-N01, AT-01/04/23.

### 2.6 DEC-06 — Thời gian, retry, khung giờ và hiệu năng

| Luồng/kênh | Hạn gửi | Retry khi lỗi tạm thời | Khung giờ gửi |
|---|---|---|---|
| OTP SMS/Email | 5s | 0, trừ ngoại lệ DEC-05 | 24/7 |
| Email lời mời, kích hoạt tài khoản | 24h hoặc hạn lời mời/kích hoạt nếu sớm hơn | 5 lần: 1′, 5′, 15′, 1h, 3h; sau TTL link dùng API cấp lại link của User | 24/7 |
| SMS lời mời | Như Email lời mời | 2 lần cách 5′ | 24/7 |
| Email/SMS đổi mật khẩu, đổi định danh, nghỉ việc | 24h (chốt theo B5, vận hành điều chỉnh khi có số liệu) | Email như lời mời; SMS 2 lần cách 5′ | 24/7 |
| In-app các mốc Order (Shop, vận hành) | Tức thì | Retry nội bộ đến khi lưu được | 24/7 |
| Push các mốc Order | 15′ | 3 lần trong 15′ | 24/7 |
| ZNS người nhận hàng (R1: giao thất bại) | Do Order cấp, mặc định 12h | 2 lần cách 5′ | 07:00–21:00 |
| SMS dự phòng người nhận hàng | Như ZNS | 2 lần cách 5′ | 07:00–21:00 |

- Múi giờ chuẩn `Asia/Ho_Chi_Minh`; mọi thời điểm trong hợp đồng dùng ISO 8601 có offset.
- Ngoài khung giờ, tin được giữ đến 07:00 nếu vẫn còn hạn; hết hạn thì dừng với lý do `EXPIRED`.
- Giới hạn tần suất nghiệp vụ (ví dụ số tin tới người nhận hàng theo đơn/ngày) do **module nguồn sở hữu**: nguồn chỉ phát ý định khi được phép gửi. Notification chỉ giữ giới hạn bảo vệ hệ thống và provider (SRS-P05).

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
- **App/Brand được đặt:** mẫu hiển thị theo brand (YAML) và định danh gửi (sender, OA ZNS, domain Email), nằm trong **kết nối kênh** (DB).
- **Shop chỉ được ghi đè:** bật/tắt Push (kênh tùy chọn). Lưu ở DB; R1 chưa có màn hình cài đặt (DEC-04) nên vận hành đặt theo yêu cầu.
- **Chia nơi lưu theo tính chất thay đổi:**
  - **File YAML trong repo** cho chính sách đổi chậm, cần người duyệt: loại thông báo, mẫu, tiêu chí đạt, mục đích, khóa Global, kênh được phép. Pull request là bước duyệt, git là phiên bản và audit. CI kiểm tra schema, biến mẫu và xung đột; hai cấu hình chồng nhau trong cùng cấp bị chặn khi merge.
  - **DB và API vận hành** cho dữ liệu đổi lúc chạy: kết nối kênh, kill-switch, hạn mức ngày, bật/tắt Push của Shop. Đổi qua API có kiểm tra quyền, mọi thay đổi vào audit, **không cần sửa repo hay deploy**.
- Mỗi bản tin ghim version YAML tại thời điểm tạo; thay YAML không ảnh hưởng tin đang chờ. Dữ liệu vận hành thì đọc giá trị hiện hành khi gửi; mỗi lần thử ghi kết nối đã dùng.
- **Kết nối kênh:** tạo ở trạng thái nháp, gửi thử đạt mới được kích hoạt. **Credential lưu trong DB**, mã hóa envelope trước khi lưu (PO chốt 12/10/2026); khóa gốc lấy từ secret manager qua biến môi trường. Credential không hiển thị lại qua API. Kết nối lỗi hoặc dừng không tự chuyển sang brand hay pháp nhân khác.
- **Kill-switch** dừng một luồng hoặc một kênh có hiệu lực ngay với cả tin đang chờ (kiểm tra trước mỗi lần đưa tin ra kênh); bật/tắt kill-switch được audit.

**Ảnh hưởng:** ST-01, SRS-F03/F12/F15, UX-05, ADR-03, AT-17/18/24.

### 2.8 DEC-08 — Ngân sách

- **Chủ ngân sách là từng brand/đơn vị kinh doanh**, theo cost center. Notification không giữ sổ tài chính và không tính phí theo Shop trong R1.
- Notification áp **hạn mức cứng theo ngày** cho từng cặp luồng × kênh có phí (SMS, ZNS, Email). Vận hành đặt hạn mức qua API nội bộ (có quyền, có audit, không cần deploy) theo đề nghị của chủ ngân sách; hệ thống cảnh báo ở 80%.
- Chạm hạn mức: bản tin nhận trạng thái `BLOCKED_QUOTA`, có lý do, **không chuyển sang kênh có phí khác** và không báo thành công.
- OTP có hạn mức riêng và cao, không bị chặn bởi ngân sách luồng khác; cảnh báo bất thường khi sản lượng giờ vượt 3 lần trung bình 7 ngày.
- Cuối tháng Notification xuất báo cáo số lượng theo cost center, luồng, kênh và provider để tài chính đối soát hóa đơn.

**Ảnh hưởng:** SRS-F14, AT-19.

### 2.9 DEC-09 — Lưu trữ, quyền xem và xóa

| Loại dữ liệu | Thời gian lưu và cách lưu |
|---|---|
| Giá trị OTP/token | Không ghi ra log, trace, dead-letter hay bảng chống trùng. Nếu phải xếp hàng: mã hóa, TTL bằng `expires_at`, tối đa 10′. Retry link sau TTL lấy link mới từ API cấp lại link của User (B1), không kéo dài TTL |
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
| OTP qua email | Email | — | — |
| Lời mời | Email hoặc SMS theo đích nguồn cấp | — | — |
| Kích hoạt tài khoản | Email | — | — |
| Đổi mật khẩu | — | — | SMS/Email tới các liên hệ đã xác minh, đạt khi một kênh đạt |
| Đổi định danh, nghỉ việc | SMS hoặc Email theo đích nguồn cấp | — | — |
| 10 mốc Order, phía Shop | In-app | Push | — |
| Mốc Order có selector vận hành | In-app | — | — |
| Giao thất bại, người nhận hàng | ZNS | — | SMS |
| Cảnh báo thiết bị | In-app | — | — |

- **Chỉ fallback SMS khi ZNS lỗi chắc chắn:** số không dùng Zalo, người nhận chặn OA, không liên lạc được hoặc template bị Zalo từ chối cho số đó. **Không fallback khi timeout hoặc chưa rõ.**
- Bounce cứng của Email là thất bại cuối cùng; không retry.
- Hết 24h vẫn chưa rõ thì đóng xử lý với kết quả `UNKNOWN`; không tự đổi thành đã giao hay thất bại.
- Callback trễ yếu hơn mức đã đạt không hạ mức đó (ví dụ `PROVIDER_ACCEPTED` đến sau `DELIVERED`). **Hard bounce hoặc báo giao thất bại đến sau `PROVIDER_ACCEPTED` là sự kiện tiêu cực cuối cùng, không phải bằng chứng yếu hơn:** delivery chuyển `FAILED`, `outcome` được tính lại và `revision` tăng. Nếu đến sau `DELIVERED` thì giữ `DELIVERED` và đánh dấu xung đột để vận hành đối chiếu.

**Ảnh hưởng:** SRS-F07/F09, mục 8 SRS, ADR-04, AT-09/25.

### 2.11 DEC-11 — Đã hủy

Ngày 09/10/2026 PO xác nhận Notitek là hệ thống mới, không có hệ thống thông báo hoặc ZNS cũ. Không có kiểm kê, feature flag chuyển luồng, shadow hay rollback về hệ thống cũ. OA và template ZNS được đăng ký mới và chờ Zalo duyệt (xem CF-08). Mã DEC-11 giữ lại để không đổi số, không tái sử dụng.

**Ảnh hưởng đã gỡ:** ST-10, SRS-F22, ADR-07, AT-20.

### 2.12 DEC-12 — Danh mục quyền

| Quyền | Phạm vi | Gán mặc định |
|---|---|---|
| Hộp tin của chính mình | Tài khoản | Không cần quyền trong catalog; kiểm tra quan hệ sở hữu và ngữ cảnh |
| `notification.order_alert.receive` | Shop | Chủ Shop, Quản lý, CSKH |
| `notification.order_ops_alert.receive` | Tổ chức nội bộ, giới hạn theo data scope Shop | Vận hành đơn, CSKH nội bộ |
| `notification.delivery.view` | Nền tảng | Vận hành Notification, CS nội bộ |
| `notification.delivery.view_sensitive` | Nền tảng | Trưởng vận hành |
| `notification.delivery.cancel` | Nền tảng | Vận hành Notification |
| `notification.delivery.resend` | Nền tảng | Vận hành Notification |
| `notification.config.edit` / `notification.config.approve` | Nền tảng | Tách hai người: người sửa không tự duyệt |
| `notification.connection.manage` | Nền tảng | Platform admin |
| `notification.usage.export` | Nền tảng | Vận hành và tài chính |

- Backend chỉ cho truy cập khi quyết định là `ALLOW` đúng hành động và tài nguyên; `DENY`, `NOT_EVALUATED` hoặc lỗi kiểm tra đều từ chối.
- R1 dùng PR review cho `config.approve` (cấu hình YAML); khi có màn hình quản trị thì quyền này được kiểm tra trong ứng dụng.
- Kết nối kênh, kill-switch, hạn mức ngày và thiết lập Push của Shop đổi qua API nội bộ. R1 dùng quyền `notification.connection.manage` cho các thao tác này (đề xuất, xác nhận cùng CF-07); mọi thay đổi vào audit.

**Ảnh hưởng:** SRS-F05/F08/F11/F12, SRS-N05, ADR-06, AT-16/21.

### 2.13 Quyết định bổ sung ngày 11/10/2026

| Mã | Quyết định | Áp dụng |
|---|---|---|
| A1 | Order gửi liên hệ người nhận hàng trong yêu cầu (bắt buộc ở giao thất bại, tùy chọn ở mốc khác theo Q2 mục 2.15); Notification không gọi ngược sang Order | DEC-03; profile `order.*` |
| A2 | R1 nguồn gọi REST từ dispatcher outbox; broker để phân kỳ sau | DEC-02; CF-03, CF-04 |
| A3 | Xác thực service-to-service bằng OAuth2 client credentials | DEC-02; CF-09 |
| A4 | Trường JSON và query param snake_case; correlation qua `X-Correlation-Id` | DEC-02; OpenAPI |
| A5 | Context dùng `organization_id` + `application_id` | DEC-03; OpenAPI |
| A6 | R1 gồm thêm `user.account_activation`, `user.password_changed`, `user.identifier_changed`, `user.employee_offboarded`; OTP nhận đích Email, lời mời nhận đích SMS và có lời mời nhân viên | DEC-04/06/10; ST-14 |
| A7 | R1 gồm đủ 10 mốc Order: `order.created`, `order.cancelled`, `order.pickup_failed`, `order.picked_up`, `order.delivery_failed`, `order.delivered`, `order.partially_delivered`, `order.return_confirmed`, `order.returned`, `order.sla_breached` | DEC-04/06/10; ST-15 |
| A8 | Thêm audience `ORDER_OPS_ALERT_RECEIVERS` và quyền `notification.order_ops_alert.receive` | DEC-03/12; CF-07 |
| B1 | User cung cấp API cấp lại link lời mời/kích hoạt; Notification gọi khi cần retry sau TTL bí mật | DEC-09; CF-10 |

Hợp đồng chi tiết: [API v1](api/README.md) và [OpenAPI](api/openapi.json). Kênh người nhận hàng R1 chỉ ở mốc giao thất bại; mở thêm cho mốc khác cần template ZNS (CF-08) và ngân sách (CF-02).

### 2.14 Hướng xử lý các vấn đề mở, chốt ngày 11/10/2026

| Mã | Quyết định | Bên thực hiện |
|---|---|---|
| B2 | OTP `UNKNOWN`: challenge vẫn hiệu lực, báo "chưa xác nhận đã gửi", chỉ cho gửi lại sau cooldown 3 phút, không tính là lần gửi thất bại | User |
| B3 | User gọi OTP ngoài transaction DB: commit challenge chờ gửi → gọi Notification → cập nhật trong transaction ngắn | User |
| B4 | `source_ref` của User tất định theo ý định; thông báo async đi qua outbox của User thay cho gửi best-effort sau commit | User |
| B5 | Giữ giá trị khởi điểm: hạn 24h, SMS thông báo tài khoản 2 lần cách 5′, 24/7; vận hành điều chỉnh khi có số liệu | Vận hành |
| B6 | R1 không có OTP qua Zalo/ZNS; User giữ kênh `ZALO` chỉ ở dev; xem lại ở R2 | PO |
| B7 | R1 chỉ tiếng Việt; envelope chưa có `locale`, bổ sung tùy chọn khi cần | PO |
| B8 | `occurred_at` là thời điểm nghiệp vụ, luôn có offset; phía Order thuộc CF-04 | User, Order |
| B9 | CF-07 kéo về 15/10, đăng ký kèm resource type | User/Authorization |
| B10 | Thêm trường tùy chọn `subject_ref` (bắt buộc và bằng `order_id` với `order.*`) và liệt kê `GET /v1/notifications?subject_ref=`; tổng hợp thông báo theo đơn của Order dựa trên API này | Notification, Order |
| B11 | User trả `delivery_status` cho frontend nội bộ theo bảng ánh xạ ở API v1 mục 8; thêm giá trị `UNCONFIRMED` cho `UNKNOWN` | User, frontend nội bộ |
| C1, C2 | Đặc tả hộp tin, SSE và đăng ký Push theo [bản nháp](api/inbox-and-push.md); hộp tin là ưu tiên cao nhất vì In-app là kênh bắt buộc của 11 profile | Notification, Shop FE |
| C3 | Webhook provider theo [bản nháp](api/provider-webhooks-and-ops.md); hoàn thiện sau khi chọn nhà cung cấp (CF-12) | Notification |
| C4 | JSON Schema CloudEvents hoãn cùng broker, không chặn R1 | — |
| C5 | API tra cứu vận hành theo [bản nháp](api/provider-webhooks-and-ops.md) | Notification |
| C6 | User cung cấp API resolve người nhận theo quyền và data scope cho Notification (CF-11) | User/Authorization |
| C7 | API cấp lại link (CF-10) theo yêu cầu ở API v1 mục 5 | User |
| D | Tên `type` của Notification là tên chuẩn ở ranh giới; Order sửa §11 hợp đồng mock theo [đề xuất](api/order-integration-requirements.md); User sửa README adapter theo [yêu cầu](api/user-integration-requirements.md) | Order, User |

Thứ tự ưu tiên: C6 + B9 → C1 → B10 → B2–B4 (gộp buổi ký CF-03).

### 2.15 Kết quả review hợp đồng, chốt ngày 11/10/2026

| Mã | Quyết định |
|---|---|
| S1 | `GET /v1/notifications`: `source` tùy chọn, mặc định là nguồn của token; chỉ vận hành được truyền nguồn khác |
| S2 | OTP: `send_deadline ≤ received_at + 5s` và `≤ expires_at`; `occurred_at` là lúc tạo challenge, không bị ràng buộc +5s |
| S3 | `source_ref` Order: `{order_id}:{type}` cho mốc một lần; `{order_id}:{type}:{id lần thất bại/vi phạm}` cho lấy hàng thất bại, giao thất bại, vi phạm SLA. Order đổi định dạng `event_key` của outbox |
| S4 | `Delivery` có `attempt_count` và `next_attempt_at`; RETRYING suy từ hai trường này |
| S5 | TTL 10′ của DEC-09 chỉ áp cho kho bí mật của Notification; link cũ trước lần gửi đầu do dispatcher User tự cấp mới |
| Q1 | Cấp lại link không vô hiệu link cũ; link cũ dùng được đến hạn riêng |
| Q2 | Liên hệ người nhận bắt buộc ở `order.delivery_failed`, tùy chọn ở 9 mốc còn lại; pháp chế/DPO xác nhận cùng CF-01 |
| Q3 | Thêm UC-ORD-12 Lấy hàng thất bại, UC-ORD-13 Đã lấy hàng, UC-ORD-14 Giao một phần, UC-ORD-15 Đã hoàn hàng, UC-ORD-16 Vi phạm SLA |
| Q4 | UX-01 mở rộng cho người vận hành trên giao diện nội bộ; tin vận hành hiển thị theo quyền và data scope, không theo organization của token |
| Q5 | PO gắn ngày cho M0–M3; M1 (gồm ST-14) không nghiệm thu trước 22/10/2026 vì phụ thuộc CF-10; PO Order đánh giá lại hạn CF-04 |
| Q6 | Order kiểm Data Scope của người dùng với đơn trước khi trả thông tin thông báo cho frontend |
| Q7 | Scope OAuth theo `<domain>.<tài_nguyên>.<hành_động>`: `notification.request.send/view/cancel`, `user.notification_recipient.resolve`, `user.notification_link.reissue` |
| Q8 | Chốt kênh 10 mốc Order theo API v1 mục 5 và selector vận hành ở lấy hàng thất bại, giao thất bại, vi phạm SLA |
| S6 | Bounce cứng hoặc báo giao thất bại đến sau `PROVIDER_ACCEPTED` hạ delivery xuống `FAILED` và tính lại `outcome`; đến sau `DELIVERED` thì giữ `DELIVERED` (giải quyết hai câu trái nhau ở DEC-10) |
| S7 | `outcome` chỉ là kết quả chốt khi `processing_state` là `CLOSED`; ánh xạ B11 xét `processing_state` trước |
| S8 | Hoãn kết quả theo trường hợp gửi (`cases[]`) khỏi R1; thêm trường tùy chọn `case_code` vào `Delivery`; `send_case` giữ làm thực thể nội bộ |

Các câu hỏi kỹ thuật còn mở trong bản nháp hộp tin và tra cứu vận hành do Tech Lead chốt khi triển khai, không chặn hợp đồng gửi.

## 3 Giá trị tham số SRS-P

| Mã | Giá trị được chốt | Nguồn |
|---|---|---|
| SRS-P01 | Phản hồi OTP ≤ 5s; bắt đầu gọi provider ≤ 2s; p95 ≤ 3s, p99 ≤ 5s; tối đa 1% vượt cửa sổ | DEC-05/06 |
| SRS-P02 | Hạn gửi theo bảng mục 2.6; mức bằng chứng theo mục 2.10 | DEC-06/10 |
| SRS-P03 | Mục tiêu độ trễ mục 2.6; 20 OTP/s, 50 event/s, burst ×5 trong 1 phút | DEC-06 |
| SRS-P04 | Tối đa 100 người nhận mỗi yêu cầu; payload ≤ 64KB; mỗi biến ≤ 2KB; trang hộp tin mặc định 20, tối đa 50, dùng cursor | PO + Tech Lead |
| SRS-P05 | Mỗi nguồn 50 yêu cầu/s mặc định; giới hạn provider theo hợp đồng; ưu tiên SECURITY > TRANSACTIONAL; quá tải trả 429 kèm `Retry-After` | DEC-06 |
| SRS-P06 | Retry và fallback theo mục 2.5, 2.6, 2.10; cửa sổ đối chiếu kết quả chưa rõ 24h | DEC-05/06/10 |
| SRS-P07 | Chống trùng 7 ngày; dấu hủy 72h; ID kết quả 180 ngày; yêu cầu cũ hơn 7 ngày bị từ chối `STALE_REQUEST` (`STALE_EVENT` khi có broker) | DEC-02/09 |
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
| CF-03 | Cơ chế đã chốt (A1–A5). Còn ký: schema 17 profile, ánh xạ kết quả OTP (`UNKNOWN`) và adapter `NotificationPort` | Tech Lead User, Tech Lead Order | 15/10/2026 | OpenAPI/schema, M1/M2 |
| CF-04 | Đội Order xây dispatcher outbox gọi REST Notification; ghi outbox cho đủ 10 mốc (gồm mốc từ webhook tracking), lưu `occurred_at` nghiệp vụ, đổi định dạng `event_key` (S3), liên hệ người nhận; phát hiện vi phạm SLA | Tech Lead Order, PO Order | 15/10/2026, PO Order đánh giá lại (Q5) | ST-04/15 và toàn bộ M2 |
| CF-05 | Sản lượng OTP và đơn hiện tại để đối chiếu tải thiết kế | Vận hành | 22/10/2026 | Nghiệm thu tải SRS-N03 |
| CF-06 | Triển khai Web Push trên Shop | Chủ Shop FE | 15/10/2026 | ST-06 |
| CF-07 | Đăng ký 9 quyền vào catalog (gồm `notification.order_ops_alert.receive`), resource type liên quan và gán mặc định | User/Authorization | 15/10/2026 | ST-05/09 |
| CF-08 | Đăng ký OA mới và template ZNS cho luồng giao thất bại, gửi Zalo duyệt (thời gian duyệt tính vào kế hoạch) | Vận hành | 22/10/2026 | ST-07 |
| CF-09 | Cấp OAuth2 client credentials cho service (Order, User, Notification): đăng ký client, scope theo Q7, issuer/audience, xoay vòng secret | User/Authorization | 15/10/2026 | Mọi tích hợp service-to-service |
| CF-10 | API cấp lại link lời mời/kích hoạt cho Notification: endpoint, schema, mã từ chối, không tính vào giới hạn cấp lại của người dùng | Tech Lead User | 22/10/2026 | Retry Email/SMS lời mời, kích hoạt sau 10′ |
| CF-11 | API resolve người nhận cho Notification: thành viên có quyền trên organization (audience Shop) và người nội bộ có data scope bao phủ Shop (audience vận hành); phân trang, client credentials, chỉ trả định danh | User/Authorization | 15/10/2026 | ST-05/08/15 và mọi tin gửi theo audience |
| CF-12 | Chọn nhà cung cấp SMS brandname và Email, kèm cơ chế callback/chữ ký | Vận hành, Tech Lead Notification | 22/10/2026 | C3, gửi SMS/Email thật |

Nếu bên ký đề xuất giá trị khác, PO cập nhật tài liệu này, ghi phiên bản và lý do; không sửa ngầm giá trị đã dùng trong code hoặc kiểm thử.

## 5 Việc tiếp theo

- Tech Lead viết ADR-01 đến ADR-06 dựa trên các quyết định trên (ADR-01 đã chốt tự xây, ADR-07 đã loại); ADR-02 (broker) hoãn đến phân kỳ sau theo A2.
- OpenAPI R1 đã có tại `api/openapi.json` (17 profile); Tech Lead User/Order review và ký theo CF-03.
- QA cập nhật test case AT với giá trị tại mục 3.
- Vận hành khởi động đăng ký brandname SMS, tài khoản FCM và đăng ký OA/template ZNS mới trong tuần đầu.
