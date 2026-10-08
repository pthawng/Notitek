# Notitek - HỢP ĐỒNG TÍCH HỢP API VÀ SỰ KIỆN

**Phiên bản:** 0.2. **Ngày:** 08/10/2026. **Trạng thái:** Bản nháp ngữ nghĩa để Tech Lead các bên hoàn thiện.

Hợp đồng xác định nguồn cung cấp gì và Notification trả gì cho các lát cắt User/Order. Các tên trường dưới đây là mô hình dữ liệu để thống nhất ý nghĩa, chưa là tên API/event đã duyệt. OpenAPI/schema event và bộ ví dụ phải được chốt trước khi các bên tích hợp.

## 1 Các thao tác tích hợp

| Thao tác | Bên gọi | Ý nghĩa |
|---|---|---|
| Yêu cầu gửi có phản hồi | User hoặc nguồn có nhu cầu tương tự | Gửi theo reference và trả mức kết quả phù hợp; OTP không coi accepted là gửi thật. |
| Nhận sự việc nguồn | Order hoặc nguồn tích hợp | Nhận sự việc đã commit và chọn luồng/trường hợp gửi theo cấu hình. |
| Tra cứu theo reference | Nguồn hoặc vận hành có quyền | Lấy kết quả ở mức yêu cầu/bản/kênh, gồm phần chưa rõ. |
| Hủy bản chờ | Nguồn/người có quyền trong phạm vi | Dừng phần chưa đưa ra ngoài theo reference/nhóm liên quan; trả phần đã gửi và phần đã dừng. |
| Nhận phản hồi kênh | Nhà cung cấp qua cơ chế được xác thực | Cập nhật kết quả đúng lần thử/bản tin, không tạo lần gửi mới từ callback. |
| Danh sách và đọc In-app | Ứng dụng với ngữ cảnh được xác nhận | Xem đúng người/phạm vi và ghi đọc độc lập. |
| Cấp/tra đích Push | Bên sở hữu thiết bị theo lựa chọn tích hợp | Cung cấp đích đúng app/người và cập nhật trạng thái hiệu lực. |
| Thực hiện hành động nghiệp vụ | Ứng dụng gọi User/Order | Nguồn kiểm tra quyền/điều kiện và trả kết quả; không gọi Notification để hoàn thành nghiệp vụ. |

## 2 Dữ liệu đầu vào và ý nghĩa

| Dữ liệu | Quy tắc nghiệp vụ | Bên xác nhận |
|---|---|---|
| Định danh nguồn | Gắn với nguồn được xác thực, không tin nhãn nguồn trong payload | Nền tảng và Tech Lead các bên. |
| Reference yêu cầu/sự việc | Ổn định cho cùng yêu cầu; phân biệt nguồn/môi trường; lần cấp mã/sự việc mới có reference phù hợp | Nguồn. |
| Loại/phiên bản đầu vào | Cho biết event hay yêu cầu gửi, hợp đồng đang dùng và luồng/mẫu cần áp dụng | Nguồn và Notification. |
| Thời điểm nguồn xác nhận | Phục vụ đánh giá độ trễ và điều kiện nguồn đã cấp; không thay phiên bản hoặc khóa chống trùng | Nguồn. |
| App và ngữ cảnh | Phân biệt tài khoản với công việc; Shop/tổ chức bắt buộc khi luồng thuộc phạm vi đó, không áp mặc định cho OTP trước đăng nhập | Nguồn/User, theo DEC-01/03. |
| Người nhận và đích | Tài khoản có định danh, liên hệ giao dịch hoặc quy tắc được duyệt; không hợp nhất theo email/số điện thoại | Nguồn/User. |
| Dữ liệu mẫu | Các biến tối thiểu đã được hợp đồng cho phép; không gửi toàn hồ sơ nguồn | Nguồn và chủ nội dung. |
| Mẫu/luồng và kênh | Tham chiếu cấu hình đã được duyệt; kênh nguồn chỉ định và kênh Notification chọn cần quy tắc ưu tiên rõ | PO nguồn/Notification, DEC-07. |
| Hạn gửi và chỉ dẫn thử lại | Giới hạn việc gửi, không là quyết định hiệu lực nghiệp vụ của mã/đơn | Nguồn và chính sách gửi. |
| Liên kết/định danh hành động | Chỉ tới chức năng nguồn đã xác nhận; nguồn tiếp tục kiểm tra quyền/token | User/Order. |
| Correlation | Nối luồng gọi, yêu cầu và kết quả; không dùng thay khóa chống trùng | Các bên. |
| Chỉ dẫn hủy/liên quan | Xác định bản chờ hoặc nhóm cần dừng; quyền hủy và phạm vi được kiểm tra | Nguồn/Notification. |
| Điều kiện ngân sách | Chỉ dẫn/ủy quyền kiểm soát theo bên sở hữu được chọn, không tự suy từ NVC/app | Bên sở hữu ngân sách, DEC-08. |

Không phải mọi đầu vào có cùng bộ trường bắt buộc. Schema riêng cho event, lệnh gửi và thao tác hộp tin phải nêu điều kiện bắt buộc theo loại/ngữ cảnh.

## 3 Các kết quả trao đổi

| Kết quả | Mức áp dụng | Ý nghĩa |
|---|---|---|
| Bị từ chối | Đầu vào | Hợp đồng/nguồn/phạm vi không hợp lệ; có mã lý do. |
| Đã tiếp nhận | Yêu cầu | Đã ghi nhận bền vững, có reference tra cứu; chưa chứng minh gửi. |
| Đang chờ/đang xử lý | Yêu cầu/bản/kênh | Có lý do và hạn xử lý theo chính sách. |
| Không gửi theo chính sách | Trường hợp/bản/kênh | Đầu vào hợp lệ nhưng điều kiện gửi không cho phép hoặc không có luồng. |
| Hủy/hết hạn | Bản/kênh chưa gửi | Phần còn chờ đã dừng; không cam kết thu hồi phần đã ra ngoài. |
| Khả dụng In-app | Bản In-app | Bản tin đã sẵn sàng cho người có quyền xem; khác với đã đọc. |
| Nhà cung cấp nhận | Lần thử/kênh | Có bằng chứng tiếp nhận gửi thật của nhà cung cấp. |
| Đã giao | Kênh | Có bằng chứng phù hợp khả năng kênh theo DEC-10. |
| Chưa rõ | Lần thử/kênh | Kết quả ngoài hệ thống chưa xác định; không tự đổi thành thất bại. |
| Thất bại cuối cùng | Kênh/bản | Không còn phương án hợp lệ theo chính sách gửi. |
| Một phần | Tổng hợp | Các người/kênh có kết quả khác nhau; trả chi tiết được phép xem. |
| Đã đọc | Bản In-app | Người dùng đọc/đánh dấu; không hoàn thành nghiệp vụ nguồn. |

Mỗi kết quả cần reference nguồn, định danh Notification, mức kết quả, thời điểm và lý do/bằng chứng khi áp dụng. Kết quả tổng hợp phụ thuộc kênh bắt buộc của từng trường hợp gửi; các bên chốt tiêu chí trước nghiệm thu. Event kết quả không chứa OTP/token hoặc liên hệ đầy đủ.

## 4 Quy tắc cần thể hiện trong schema và kiểm thử

- Cùng reference và cùng dữ liệu nghiệp vụ giữ cùng yêu cầu; tiếp tục phần chưa hoàn tất. Cùng reference nhưng dữ liệu xung đột cần kết quả riêng được các bên thống nhất; không âm thầm ghi đè bản đã gửi.
- Event nguồn trùng, đến trễ và sai thứ tự được xử lý theo sự việc/phiên bản và chỉ dẫn nguồn. Nguồn cấp việc nào còn cần thông báo; Notification không suy từ mô hình đơn/OTP.
- Callback phải được xác thực, gắn đúng lần thử và chống xử lý trùng. Chuyển trạng thái khi callback đến muộn không làm mất bằng chứng đã có; quy tắc xung đột phải được liệt kê.
- Hủy phải tách phần có thể dừng với phần đã đưa ra kênh. Hợp đồng giải thích giới hạn khi hủy và gửi xảy ra gần nhau.
- Thử lại gửi và gửi lại event kết quả là hai thao tác khác nhau. Nguồn chống trùng kết quả và không tạo vòng lặp gửi từ chính kết quả.
- Các ID người/app/Shop ở trình duyệt không tự xác nhận scope. Backend đối chiếu ngữ cảnh/quyền từ cơ chế nền tảng.
- Truyền bí mật theo đường bảo vệ; dữ liệu dùng trong log, callback, sự kiện kết quả và hàng lỗi phải được phân loại riêng.

## 5 Ánh xạ từng module

| Nguồn | Căn cứ hiện có | Điều cần hoàn thiện |
|---|---|---|
| User gửi mã/link | NotificationPort nhận Notification và trả NotificationReceipt | Ánh xạ reference, đích, template, parameters, expiresAt; bảo vệ bí mật; bổ sung ngữ cảnh đáng tin cậy khi cần; biểu diễn lỗi/chưa rõ và mức kết quả đủ cho User. |
| User cảnh báo thiết bị | Chuông Shop dùng feed và hành động device trust của User | Chọn feed/event và khóa đối chiếu để chuyển đổi; giữ trạng thái nghiệp vụ và hành động ở User; tách trạng thái đọc. |
| Order | Có các event ghi outbox; sự việc delivery.* chưa là hợp đồng được xác nhận | Định nghĩa đầu vào giao thất bại, người nhận và chỉ dẫn hủy; xây/kiểm chứng đường phát bền vững; thống nhất lỗi/tương thích. |
| User/Authorization | Application/Client, membership, Access Context và cơ chế quyền chung | Hành động/scope Notification, căn cứ xử lý nền và tách ngữ cảnh tài khoản/công việc. |
| App/năng lực thiết bị | Bên sở hữu đích Push cần xác định | Client cụ thể, đăng ký/tra/thu hồi đích, đổi chủ thiết bị và xử lý token lỗi. |

`NotificationReceipt.providerReference` hiện không tự biểu diễn đầy đủ bảng kết quả. Adapter mới phải chuyển ngữ nghĩa đã chốt vào cách gọi User sử dụng; không coi một chuỗi reference bất kỳ là bằng chứng gửi thành công.

## 6 Artifact hợp đồng cần bàn giao

Tech Lead tạo OpenAPI cho các thao tác HTTP được chọn; schema event cho các luồng bất đồng bộ; bộ ví dụ có dữ liệu giả và không chứa bí mật thật; quy tắc tương thích/version; cơ chế xác thực nguồn/callback; giới hạn kích thước, lỗi và timeout. Mỗi consumer/provider có kiểm thử hợp đồng. Tên event delivery.* và route chưa chốt không được đưa vào code như hợp đồng đang hoạt động.

## 7 Các quyết định liên quan

DEC-01/02/03/05/06/07/08/09/10/12 ở [phạm vi](<Notitek - PHẠM VI VÀ KẾ HOẠCH PHÁT HÀNH ĐẦU TIÊN.md>) điều khiển các phần hợp đồng tương ứng. Các bên có thể chốt hợp đồng của OTP/lời mời trước hợp đồng Order; đồng thời giữ thống nhất định danh, kết quả và quy tắc tương thích dùng chung.

## 8 Đồng bộ với SRS 1.0 sau PO review

[SRS](<Notitek - ĐẶC TẢ YÊU CẦU PHẦN MỀM (SRS).md>) chốt baseline hành vi. Hợp đồng này vẫn là ngữ nghĩa cần cụ thể hóa thành schema được consumer/provider xác nhận; bảng dưới là đầu ra bắt buộc cho phần liên quan, không tự tạo route hoặc event mới.

| Phần hợp đồng | Nội dung phải thể hiện | SRS liên quan |
|---|---|---|
| Gửi đồng bộ | Mức bằng chứng cần trả, thời gian chờ và hạn bắt đầu/xử lý; khi hết cửa sổ không tự chuyển sang gửi nền; callback trễ chỉ cập nhật bằng chứng | SRS-F01.05, SRS-F07.08/09, SRS-N01 |
| Chống trùng | Phạm vi khóa, dữ liệu có ý nghĩa được so, metadata được phép khác, xung đột, xử lý đồng thời và replay ngoài cửa sổ; bảo vệ dữ liệu đối chiếu bí mật | SRS-F02 |
| Hủy | Tập mục tiêu/scope, ID lệnh, phần dừng/ra ngoài/chưa rõ; hỗ trợ hoặc phản hồi rõ hủy trước event; thời gian giữ dấu hủy | SRS-F06 |
| Kết quả | Tách tiếp nhận, tiến trình, bằng chứng kênh, tiêu chí đạt và đọc; kết quả riêng theo người/kênh/lần thử, revision và ID thay đổi để phát lại | SRS-F09 |
| Kênh/provider | Loại bằng chứng thật, callback tin cậy, mã lỗi, điều kiện retry/fallback, giới hạn và version kết nối; không coi HTTP thành công là đã giao | SRS-F07/15/17/18/19 |
| Hộp tin | Loại ngữ cảnh, quyền hiện tại, phân trang/thứ tự, số chưa đọc cùng scope, ghi đọc idempotent ở backend và kết quả an toàn với ID trái quyền | SRS-F05/08, SRS-N05 |
| Push | Chủ đăng ký/thu hồi, tài khoản–app/client–đích, version/độ mới, đổi chủ, token lỗi và chính sách đa đích; không dùng dữ liệu cũ hồi sinh token | SRS-F16 |
| Phí | Bên quyết định, căn cứ cho từng lần thử/fallback, phạm vi cho phép, trạng thái chưa rõ và trách nhiệm giữ/giải phóng nếu có | SRS-F14 |
| Bí mật và lưu dữ liệu | Phân loại trường, đường truyền/lưu/quan sát được phép, che dữ liệu, TTL theo loại, hàng lỗi và dữ liệu chống trùng | SRS-F02.07, SRS-F13, SRS-N06 |
| Giới hạn/tương thích | Kích thước/tập nhận/tốc độ/phân trang, lỗi hạn chế tải, version và ví dụ consumer tests | SRS-F21, SRS-N03/07 |

Mỗi lát cắt phải chốt phần hợp đồng cần dùng và các tham số SRS-P liên quan. Nguồn không được dùng kết quả gửi/đọc để tự kết luận đã hoàn tất nghiệp vụ nếu hợp đồng nghiệp vụ của nguồn không xác nhận.
