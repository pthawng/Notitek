# Notitek - DANH MỤC USE CASE

**Trạng thái:** Danh mục nghiệp vụ có các lát cắt đợt đầu đã xác nhận; chính sách chi tiết, SRS và hợp đồng tích hợp có trạng thái riêng.  
**Ngày đồng bộ BRD 0.5:** 08/10/2026.  
**Tài liệu nền:** [BRD Notification hiện tại](<Notitek - ĐẶC TẢ YÊU CẦU NGHIỆP VỤ (BRD).md>).  
**Phạm vi đối chiếu:** Mã nguồn hiện có của [User](../user-spf/) và [Order](../order-spf/).

## 1. Cách đọc danh sách

Tên thành phần dùng theo BRD: Luồng thông báo, Đối tượng nhận thông báo, Nhóm nhận thông báo, Ngữ cảnh thông báo, Quy tắc chọn người nhận và Hộp thông báo. Bản 0.5 giữ nguyên 37 mã UC trước đây và bổ sung UC-USR-09 cho cảnh báo/yêu cầu thiết bị, tổng cộng 38 UC.

- **P0 / P1 / P2** là ưu tiên của danh mục; không đồng nhất với cam kết của bản phát hành. Các lát cắt đã xác nhận nằm trong bảng dưới; UC khác vẫn theo phân kỳ riêng.
- **Đã có outbox**: module nguồn ghi event vào cơ sở dữ liệu; chưa có nghĩa event đã được phát đến Notification.
- **Gửi trực tiếp**: User đang gọi `NotificationPort.send(...)`; cần quyết định cách chuyển sang hợp đồng tích hợp mới.
- **Cần event**: nghiệp vụ hoặc trạng thái đã có nhưng chưa thấy event nguồn riêng đủ rõ để kích hoạt thông báo.
- Tên event ghi chữ **đề xuất** trong tài liệu này không phải tên hợp đồng đang hoạt động.

Đợt đầu đã xác nhận đủ 5 kênh: In-app, Push, Email (bao gồm Gmail), SMS và Zalo ZNS. Ngày 08/10/2026, người yêu cầu chốt các lát cắt dưới đây, tiếp nối Shop SuperShip và giao diện nội bộ. Client Push cụ thể và các app mở rộng còn cần chốt. Điều kiện dừng do User/Order cấp/xác nhận qua hợp đồng; Notification không tự suy trạng thái mã/đơn.

| Lát cắt đã chọn | UC | Phạm vi |
|---|---|---|
| OTP | UC-USR-01 | Hợp đồng gửi có phản hồi; User sở hữu mã và xác thực. |
| Lời mời | UC-USR-03, UC-USR-04 | Email sau commit, link và hành động ở User. |
| Giao thất bại | UC-ORD-05 | In-app/Push cho Shop, ZNS/SMS theo chính sách cho liên hệ giao dịch. |
| Cảnh báo/yêu cầu thiết bị | UC-USR-09 | Tiếp nối chuông hiện có; hành động bảo mật ở User. |
| Hộp tin và tra cứu | UC-NTF-13, UC-NTF-12 | Đúng ngữ cảnh tài khoản/công việc và tra kết quả theo quyền. |

Các UC-NTF nền liên quan và cấu hình tối thiểu của UC-NTF-15 phục vụ các lát cắt. Xem [phạm vi](<docs/Notitek - PHẠM VI VÀ KẾ HOẠCH PHÁT HÀNH ĐẦU TIÊN.md>), [UC ưu tiên](<docs/Notitek - ĐẶC TẢ USE CASE ƯU TIÊN.md>) và [ma trận truy vết](<docs/Notitek - BACKLOG TRIỂN KHAI VÀ KỊCH BẢN NGHIỆM THU.md>). Không coi toàn bộ UC P0 trong danh mục là phạm vi đợt đầu.

Ranh giới nghiệp vụ: module nguồn xác nhận sự việc và sở hữu trạng thái thật. Notification nhận event hoặc yêu cầu gửi, chọn luồng và người nhận, áp dụng chính sách, gửi qua kênh, ghi kết quả và phát event kết quả. Một event nguồn có thể tạo nhiều **trường hợp gửi** cho các nhóm người nhận khác nhau. Notification không tự suy luận kết quả giao hàng từ vận đơn hoặc webhook kỹ thuật.

```text
Sự việc đã được module nguồn xác nhận
  -> tiếp nhận và chống trùng
  -> chọn luồng + trường hợp gửi + người nhận
  -> kiểm tra chính sách + mẫu + kênh
  -> gửi và nhận phản hồi nhà cung cấp
  -> lưu lịch sử + phát kết quả thông báo
```

## 2. Use case dùng chung của Notification

### Ranh giới với Module User

User/Authorization đã quản lý tài khoản, tổ chức, membership, Application/Client, Entitlement và phân quyền theo context. Notification dùng các định danh/quyết định đó, không yêu cầu tạo lại người, app hoặc role.

Notification cấu hình sự việc kích hoạt, trường hợp gửi, quy tắc chọn người nhận/nhóm nhận thông báo có căn cứ, mẫu, kênh, lịch và điều kiện dừng. Thông tin nhận thông báo là dữ liệu tối thiểu phục vụ gửi; nhóm nhận thông báo và lựa chọn nhận thông báo không cấp quyền xem đơn/app.

- UC-NTF-05 dùng định danh và dữ liệu User/nguồn, cùng quy tắc REC; không tự suy quyền theo tên vai trò.
- UC-NTF-07/08 cấu hình nội dung/kênh; app/client đích tham chiếu User, không tự cấp quyền dùng app.
- UC-NTF-13 dùng người/ngữ cảnh do backend xác nhận; đọc In-app và mở đơn vẫn kiểm tra quyền hiện tại.
- UC-NTF-15 là quản trị tự phục vụ nâng cao đề xuất P1; cấu hình tối thiểu để mọi luồng chạy thuộc phạm vi ban đầu theo BR-CAT-04.
- Các luồng chạy nền không yêu cầu phiên đăng nhập của người nhận và không dùng token của người tạo thay quyền tập nhận.

Tham khảo mô hình luồng/đối tượng nhận/ngữ cảnh thông báo của Novu được diễn giải ở mục 4.4 BRD; cách cấu hình cụ thể ở mục 6.6–6.7. Chưa có quyết định dùng Novu làm nền tảng triển khai.

| Mã | Use case | Actor chính | Kết quả cần đạt | Nhóm BRD | Ưu tiên |
|---|---|---|---|---|---|
| UC-NTF-01 | Tiếp nhận event nghiệp vụ | Module nguồn | Chỉ nhận event hợp lệ, có nguồn, phiên bản, thời gian, định danh sự việc và phạm vi tổ chức rõ ràng | TRG, INT | P0 |
| UC-NTF-02 | Tiếp nhận yêu cầu gửi có phản hồi | Module User hoặc nguồn có luồng bắt buộc | Trả trạng thái tiếp nhận/gửi phù hợp để nguồn quyết định có tiếp tục luồng hay không; đặc biệt cho OTP | TRG, INT | P0 |
| UC-NTF-03 | Chống xử lý trùng và tiếp tục sau lỗi | Module nguồn, Notification | Một sự việc gửi lại không tạo bản thông báo trùng; lỗi tạm thời được xử lý tiếp có kiểm soát | TRG | P0 |
| UC-NTF-04 | Ánh xạ event vào luồng thông báo | Notification | Xác định hành trình, luồng, mục đích và các trường hợp gửi; event không có luồng được ghi nhận là bỏ qua | CAT, CLS | P0 |
| UC-NTF-05 | Xác định người nhận | Notification, User | Chọn đúng tài khoản, liên hệ không có tài khoản hoặc đầu mối tổ chức trong đúng Shop/đơn vị | REC, INT | P0 |
| UC-NTF-06 | Kiểm tra điều kiện được gửi | Notification | Áp dụng liên hệ/quyền, danh sách chặn, phạm vi, hạn gửi, chỉ dẫn dừng và điều kiện ngân sách theo trách nhiệm được giao | PREF, SEC, TIME, COST | P0 |
| UC-NTF-07 | Tạo nội dung từ mẫu đã duyệt | Notification, quản trị viên | Chọn đúng brand/app, ngôn ngữ, kênh, phiên bản mẫu; lưu dấu nội dung thực tế đã gửi | TPL, CFG | P0 |
| UC-NTF-08 | Chọn kênh và nhà cung cấp | Notification | Chọn kênh chính, song song, dự phòng hoặc bổ sung theo điều kiện đã cấu hình | CHN, CFG | P0 |
| UC-NTF-09 | Gửi và theo dõi từng lần thử | Notification, nhà cung cấp kênh | Ghi lần gửi, phản hồi, lỗi và quyết định thử lại; phân biệt nhà cung cấp đã nhận với người dùng đã nhận | CHN | P0 |
| UC-NTF-10 | Dừng thông báo không còn giá trị | Module nguồn, Notification | Hủy tin đang chờ khi OTP/link hết hạn, đơn đã đổi trạng thái hoặc người nhận đã hoàn tất hành động | TRG, TIME | P0 |
| UC-NTF-11 | Phát event kết quả thông báo | Notification, module nguồn | Nguồn có thể biết yêu cầu đã được nhận, bị chặn, gửi thất bại hoặc được xác nhận đã giao | INT | P0 |
| UC-NTF-12 | Tra cứu một thông báo từ đầu đến cuối | Quản trị viên, vận hành | Trả lời được vì sao gửi, gửi cho ai, theo mẫu/quy tắc nào, qua kênh nào và kết quả ra sao | SEC, OPS | P0 |
| UC-NTF-13 | Trung tâm thông báo In-app | Người dùng có tài khoản | Xem đúng tài khoản/ngữ cảnh, đánh dấu đọc độc lập và mở đối tượng do nguồn kiểm tra quyền | CHN, REC, SEC, CFG | P0 |
| UC-NTF-14 | Quản lý lựa chọn nhận thông báo | Người nhận, quản trị viên | Quản lý đồng ý/từ chối Marketing và kênh ưa thích; không dùng từ chối Marketing để chặn tin bắt buộc | PREF | P1 |
| UC-NTF-15 | Quản lý luồng, mẫu và chính sách theo phạm vi | Quản trị viên được cấp quyền | Cấu hình, duyệt, phiên bản hóa, xem trước và kiểm soát kế thừa/khóa cấu hình theo tổ chức/app/brand | CAT, TPL, CFG | P1 |
| UC-NTF-16 | Vận hành lỗi và gửi lại | Vận hành Notification | Xem tồn đọng, gửi lại có kiểm soát, đổi nhà cung cấp/kênh, tạm dừng luồng và lưu nhật ký thao tác | CHN, OPS | P1 |
| UC-NTF-17 | Giới hạn tần suất, gom nhóm và chi phí | Notification, chủ ngân sách | Chống làm phiền, áp dụng thời gian gửi và theo dõi chi phí theo kênh/đơn vị chịu phí | TIME, CFG, COST | P1 |
| UC-NTF-18 | Chiến dịch gửi diện rộng | Quản trị viên chiến dịch | Duyệt tập nhận, gửi thử, kiểm tra quyền nhận tin trước giờ gửi và dừng khẩn cấp | CMP | P2 |

### Kết quả đầu ra cần phân biệt

Các tên dưới đây chỉ là **đề xuất hợp đồng**; cần chốt với module tiêu thụ trước khi triển khai.

Timeout có thể là kết quả chưa rõ, không tự coi là thất bại để gửi thêm kênh. Yêu cầu có nhiều người/kênh phải báo kết quả riêng và kết quả một phần khi cần. Bản chờ bị hủy hoặc hết hạn có kết quả riêng, không gộp thành lỗi gửi; xem mục 6.4–6.5 BRD.

| Event đề xuất | Ý nghĩa |
|---|---|
| `notification.accepted.v1` | Notification đã ghi nhận yêu cầu bền vững; **chưa** có nghĩa đã gửi tới người nhận. |
| `notification.suppressed.v1` | Chính sách hoặc trạng thái hiện tại quyết định không gửi; có mã lý do. |
| `notification.cancelled.v1` / `notification.expired.v1` | Bản chờ đã hủy hoặc hết hạn; không đồng nhất với lỗi kênh. |
| `notification.provider_accepted.v1` | Nhà cung cấp kênh đã nhận yêu cầu gửi; **chưa** có nghĩa người dùng đã nhận. |
| `notification.delivered.v1` | Có bằng chứng giao thành công theo định nghĩa của kênh. |
| `notification.failed.v1` | Đã hết phương án gửi hợp lệ hoặc lỗi cuối cùng; có mã lý do. |
| `notification.read.v1` | Người dùng mở/đánh dấu đã đọc thông báo In-app; không đồng nghĩa đã hoàn thành công việc nghiệp vụ. |

## 3. Luồng tích hợp từ Module User

Các luồng đang gửi qua `NotificationPort` và template trong User: OTP, kích hoạt tài khoản, lời mời nhân viên/thành viên, thay đổi định danh, đổi mật khẩu và nghỉ việc. Riêng OTP hiện yêu cầu **gửi thật**, không được coi việc ghi log hoặc nhận event bất đồng bộ là thành công. Xem [NotificationPort](../user-spf/services/users-core-service/src/main/java/com/supership/users/notification/application/NotificationPort.java) và [NotificationTemplates](../user-spf/services/users-core-service/src/main/java/com/supership/users/notification/application/NotificationTemplates.java).

| Mã | Luồng / sự việc nguồn | Trường hợp gửi và điều kiện dừng chính | Hiện trạng nguồn | Ưu tiên |
|---|---|---|---|---|
| UC-USR-01 | Phát hành OTP cho xác thực hoặc thao tác nhạy cảm | Gửi tới liên hệ của thử thách; không gửi mã hết hạn; nguồn cần biết gửi thất bại để dừng luồng | Gửi trực tiếp; cần hợp đồng có phản hồi | P0 |
| UC-USR-02 | Phát hành liên kết kích hoạt tài khoản | Gửi email cho chủ tài khoản; ngừng gửi lại link cũ khi token đã dùng/hết hạn | Gửi trực tiếp | P0 |
| UC-USR-03 | Mời nhân viên vào tổ chức | Gửi link đúng người được mời; lần cấp lại thay thế link cũ; hủy khi lời mời bị thu hồi/hết hạn | Gửi sau commit qua `NotificationPort` | P0 |
| UC-USR-04 | Mời thành viên vào Shop | Gửi link cho người được mời với tên Shop và người mời; tuân thủ trạng thái lời mời | Gửi sau commit qua `NotificationPort` | P0 |
| UC-USR-05 | Mật khẩu đã thay đổi/đặt lại | Cảnh báo chủ tài khoản qua liên hệ bảo mật; không chứa mật khẩu hoặc mã xác thực | Qua NotificationPort sau commit | P0 |
| UC-USR-06 | Số điện thoại/email đăng nhập đã thay đổi | Cảnh báo tới liên hệ **cũ**, chỉ hiển thị giá trị mới đã che bớt | Qua NotificationPort sau commit | P0 |
| UC-USR-07 | Chấm dứt tư cách nhân viên | Báo người bị ảnh hưởng và đầu mối có trách nhiệm theo chính sách; không gửi nếu quyết định chưa có hiệu lực | Gửi trực tiếp | P1 |
| UC-USR-08 | Yêu cầu phê duyệt và kết quả thay đổi tư cách/quyền | Nhắc đúng người duyệt; báo kết quả cho người yêu cầu; dừng nhắc khi đã xử lý | Có nghiệp vụ phê duyệt, cần chốt event và người nhận | P1 |
| UC-USR-09 | Cảnh báo và yêu cầu về thiết bị | Chủ tài khoản xem tin, mở và thực hiện hành động tại User; đọc tách kết quả bảo mật; không gắn mặc định vào Shop | Chuông Shop dùng feed device trust của User; cần hợp đồng tiếp nối vào hộp tin chung | P0 |

## 4. Luồng tích hợp từ Module Order

Order hiện ghi các event `order.created.v1`, `order.cancelled.v1`, `order.updated.v1`, `order.return.requested.v1`, `order.return.confirmed.v1` vào outbox. [OutboxWriter](../order-spf/src/main/java/vn/supership/superplatform/order/platform/outbox/OutboxWriter.java) ghi rõ dispatcher phát đi broker/HTTP **chưa được xây dựng**. Các mốc lấy/giao hàng có trạng thái nghiệp vụ nhưng cần event nguồn đủ rõ trước khi Notification sử dụng. Xem [trạng thái Order](../order-spf/src/main/java/vn/supership/superplatform/order/feature/order/internal/domain/model/OrderActionPolicy.java).

| Mã | Luồng / sự việc nguồn | Trường hợp gửi và điều kiện dừng chính | Hiện trạng nguồn | Ưu tiên |
|---|---|---|---|---|
| UC-ORD-01 | Đơn được tạo thành công | Xác nhận cho Shop/người tạo; không mặc định gửi cho người nhận hàng khi đơn chỉ vừa được tạo | `order.created.v1` đã có outbox | P0 |
| UC-ORD-02 | Tạo vận đơn NVC thất bại, cần xử lý | Cảnh báo Shop/nhân sự chịu trách nhiệm với hướng xử lý; chỉ gửi sau kết quả thất bại đã chốt | Có nghiệp vụ/trạng thái; cần event riêng | P0 |
| UC-ORD-03 | Đơn bị hủy thành công | Báo Shop và người bị ảnh hưởng theo chính sách; dừng lịch nhắc lấy/giao đang chờ | `order.cancelled.v1` đã có outbox | P0 |
| UC-ORD-04 | Bắt đầu giao hàng | Người nhận hàng nhận thông tin cần thiết; Shop có thể nhận In-app/Push; phải gắn đúng chặng/NVC | Có trạng thái; event `delivery.started` là đề xuất | P0 |
| UC-ORD-05 | Giao hàng thất bại | Người nhận biết kết quả/hướng xử lý; Shop biết đơn cần can thiệp; tránh lặp theo webhook trùng | Có trạng thái; event `delivery.failed` là đề xuất | P0 |
| UC-ORD-06 | Giao hàng thành công | Xác nhận hoàn tất cho các bên được chọn; hủy thông báo nhắc giao chưa gửi | Có trạng thái; event `delivery.completed` là đề xuất | P0 |
| UC-ORD-07 | Thay đổi thông tin quan trọng của đơn | Chỉ báo khi trường thay đổi ảnh hưởng người nhận, tiền thu hoặc việc giao; không gửi cho mọi chỉnh sửa nội bộ | `order.updated.v1` đã có outbox; cần quy tắc chọn trường | P1 |
| UC-ORD-08 | Yêu cầu chuyển hoàn | Báo Shop/người xử lý về yêu cầu và bước tiếp theo; không coi là hoàn tất chuyển hoàn | `order.return.requested.v1` đã có outbox | P1 |
| UC-ORD-09 | Xác nhận chuyển hoàn | Báo kết quả đã được module nguồn xác nhận | `order.return.confirmed.v1` đã có outbox | P1 |
| UC-ORD-10 | Đổi NVC hoặc bàn giao chặng thất bại | Báo Shop/nhân sự vận hành khi cần hành động; tránh nhân thông báo theo số vận đơn kỹ thuật | Có nghiệp vụ đa chặng; cần event riêng | P1 |
| UC-ORD-11 | COD/đối soát/thanh toán hoàn tất | Gửi đúng chủ Shop/kế toán và số tiền đã chốt; nguồn phải là module sở hữu khoản tiền | Chưa đưa vào tích hợp Order hiện tại; cần module sở hữu tài chính | P2 |

### Ví dụ phân rã một luồng

Với **UC-ORD-05 - Giao hàng thất bại**, module nguồn phát một event khi kết quả đã chốt. Notification có thể tạo hai trường hợp gửi: (1) người nhận hàng qua ZNS/SMS theo chính sách và liên hệ trên đơn; (2) Shop qua In-app/Push. Nếu Order xác nhận giao thành công hoặc chuyển hoàn trước giờ gửi, Notification hủy thông báo chưa gửi. Phản hồi từ nhà cung cấp và kết quả giao được lưu riêng cho từng người nhận/kênh.

## 5. Hợp đồng và quyết định cần chốt trước SRS

1. **Event đầu vào:** thống nhất `eventId`, `eventType`, `schemaVersion`, `occurredAt`, aggregate ID/version, tổ chức/Shop, app/brand, correlation ID, người nhận hoặc căn cứ xác định người nhận, dữ liệu nghiệp vụ tối thiểu và thời điểm hết giá trị.
2. **Event hay yêu cầu gửi:** OTP và các mã/link nhạy cảm cần hợp đồng an toàn có kết quả phản hồi; không đưa bí mật vào event bus phổ thông hoặc log. Event sau commit phù hợp với các thông báo về sự việc đã hoàn tất.
3. **Quyền sở hữu người nhận:** nguồn cung cấp người nhận phát sinh theo giao dịch; User cung cấp bản đọc định danh/tư cách tối thiểu đã phê duyệt. Khi thiếu hoặc hết hiệu lực, không tự đoán người nhận.
4. **Ý nghĩa kết quả:** `accepted`, nhà cung cấp nhận, `delivered`, `failed`, `read` và hành động nghiệp vụ đã hoàn tất là các trạng thái khác nhau; module nguồn chỉ dựa vào trạng thái đã định nghĩa rõ.
5. **Đơn hàng đa chặng:** chọn rõ sự kiện theo Order hay theo vận đơn/chặng; không gửi cho con người trên từng bản ghi tracking kỹ thuật.
6. **Phạm vi giai đoạn đầu:** BRD 0.5 đã chốt 5 kênh và các lát cắt OTP, lời mời, giao thất bại, cảnh báo thiết bị; client Push/app mở rộng và các năng lực ngoài đợt đầu còn cần chốt.
7. **Quyền và chi phí:** xác định ai có thể cấu hình mẫu/luồng theo tổ chức, ai duyệt, ai trả phí SMS/ZNS/Email và khi hết hạn mức phải xử lý ra sao.

## 6. Mẫu phân rã cho từng use case ở bước tiếp theo

Mỗi UC ở mục 3 và 4 nên được chi tiết hóa bằng cùng một biểu mẫu:

| Trường | Nội dung cần chốt |
|---|---|
| Sự kiện nguồn | Module sở hữu, tên/phiên bản event, thời điểm được coi là đã xảy ra |
| Điều kiện kích hoạt | Khi nào tạo thông báo; trường hợp nào chỉ hiển thị trên màn hình và không gửi |
| Trường hợp gửi | Nhóm người nhận, mục tiêu từng nhóm, kênh chính/song song/dự phòng/bổ sung |
| Dữ liệu và mẫu | Biến mẫu được phép dùng, ngôn ngữ, app/brand, đường dẫn hành động |
| Thời gian | Gửi ngay/hẹn giờ, hết hạn, giới hạn tần suất, điều kiện dừng |
| Bảo vệ dữ liệu | Dữ liệu nhạy cảm, che bớt thông tin, quyền xem và thời gian lưu |
| Kết quả | Event đầu ra nào nguồn cần nhận, trường hợp nào cần phản hồi đồng bộ |
| Nghiệm thu | Kịch bản thành công, trùng event, sai người nhận, nguồn đổi trạng thái, nhà cung cấp lỗi |

**Ghi chú BRD:** Bản 0.5 ghi nhận phạm vi sản phẩm đã chốt; giữ trạng thái riêng cho các chính sách/hợp đồng còn mở. Mã UC dùng cho phân rã, BR dùng cho yêu cầu. Mục 13.1 BRD liên kết nhóm; [bộ tài liệu triển khai](<docs/Notitek - HƯỚNG DẪN SỬ DỤNG BỘ TÀI LIỆU.md>) bổ sung UC chi tiết, SRS, trải nghiệm, hợp đồng, thiết kế và ma trận UC–BR–story–AT của đợt đầu.
