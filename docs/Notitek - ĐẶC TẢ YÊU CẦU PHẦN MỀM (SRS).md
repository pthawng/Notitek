# Notitek - ĐẶC TẢ YÊU CẦU PHẦN MỀM (SRS)

**Phiên bản:** 0.1. **Ngày:** 08/10/2026. **Trạng thái:** Bản nháp để BA và Tech Lead chốt theo từng lát cắt.

SRS mô tả hành vi hệ thống cần để thực hiện [use case đợt đầu](<Notitek - ĐẶC TẢ USE CASE ƯU TIÊN.md>). Các yêu cầu chức năng dưới đây chuyển từ BRD và phạm vi sản phẩm; route, cấu trúc lưu trữ, công nghệ và giá trị NFR cần hoàn thiện trong hợp đồng và thiết kế kỹ thuật trước khi triển khai phần liên quan.

## 1 Phạm vi hệ thống

Hệ thống nhận event hoặc yêu cầu gửi từ nguồn được phép, áp dụng cấu hình luồng/trường hợp gửi, dựng nội dung, điều phối năm kênh, quản lý hộp In-app và cung cấp kết quả/tra cứu. User/Order sở hữu trạng thái nghiệp vụ; Notification chỉ giữ tham chiếu và dữ liệu tối thiểu phục vụ thông báo.

## 2 Yêu cầu chức năng

| Mã | Hành vi bắt buộc của hệ thống | UC chính | BR chính |
|---|---|---|---|
| SRS-F01 | Kiểm tra nguồn, phiên bản hợp đồng, trường bắt buộc và phạm vi được phép; từ chối có lý do, không báo tiếp nhận khi chưa ghi nhận bền vững | UC-NTF-01/02 | BR-TRG-01/04, BR-INT-04 |
| SRS-F02 | Cùng yêu cầu/sự việc đã xác định trả tham chiếu cũ và tiếp tục phần chưa hoàn tất theo chính sách; không gửi lại phần đã đạt tiêu chí | UC-NTF-03 | BR-TRG-02/05 |
| SRS-F03 | Chọn cấu hình luồng/trường hợp gửi đúng phiên bản và phạm vi; event không có luồng có kết quả không gửi theo chính sách | UC-NTF-04/08 | BR-CAT-01/02/04, BR-CFG-01/02/03/05 |
| SRS-F04 | Dựng nội dung từ mẫu được duyệt và dữ liệu nguồn; thiếu biến thì chờ/dừng có lý do; dấu phiên bản và bằng chứng nội dung tuân thủ chính sách bí mật | UC-NTF-07 | BR-TPL-01/02/03/05 |
| SRS-F05 | Dùng người nhận và căn cứ do nguồn/User cung cấp; tách phạm vi tài khoản/công việc; không hợp nhất tài khoản theo số/email hoặc suy quyền từ tổ chức | UC-NTF-05/06/13 | BR-REC-01/02/03/04/05/06, BR-INT-06, BR-SEC-03/06 |
| SRS-F06 | Trước đưa tới kênh kiểm tra hạn gửi và điều kiện nguồn/chính sách đã được giao; hủy bản chờ theo lệnh hợp lệ, không cam kết thu hồi tin đã ra ngoài | UC-NTF-06/10 | BR-TIME-01/02, BR-INT-05 |
| SRS-F07 | Thực hiện kênh chính/song song/dự phòng/bổ sung theo cấu hình; thử lại hữu hạn theo hạn gửi; timeout giữ chưa rõ và không tự kích hoạt dự phòng | UC-NTF-08/09 | BR-CHN-01/02/03 |
| SRS-F08 | Tạo và cung cấp In-app đúng người/ngữ cảnh; ghi trạng thái đọc độc lập; backend kiểm tra quyền trên tải danh sách, đọc và tra cứu | UC-NTF-13 | BR-CHN-05, BR-REC-03/04, BR-SEC-03/06 |
| SRS-F09 | Ghi kết quả từng người/kênh/lần thử; tra cứu hoặc phát kết quả cho bên được phép; phân biệt tiếp nhận, khả dụng In-app, nhà cung cấp nhận, giao, chưa rõ, lỗi, hủy, hết hạn và một phần | UC-NTF-09/11/12 | BR-INT-04, BR-CHN-04/05/09 |
| SRS-F10 | Hộp tin tiếp nối cảnh báo thiết bị User; hành động gửi tới User và dùng kết quả User, không tự cho phép đăng nhập/thu hồi phiên | UC-USR-09, UC-NTF-13 | BR-INT-01/06, BR-SEC-03 |
| SRS-F11 | Cho vận hành tra từ reference nguồn tới luồng/cấu hình, người nhận được phép xem, lần thử và lý do; xem và thay đổi là quyền riêng | UC-NTF-12 | BR-OPS-01, BR-SEC-02/04/06 |
| SRS-F12 | Cung cấp cấu hình tối thiểu cho từng luồng/kênh qua phương thức vận hành được kiểm soát; kiểm tra thiếu/xung đột trước kích hoạt | UC-NTF-15 | BR-CAT-04, BR-CFG-01/02/03/05 |
| SRS-F13 | Bảo vệ OTP/token trong truyền, xử lý và lưu; không xuất hiện trong log, lịch sử hỗ trợ hoặc event kết quả; phân biệt rõ mock và gửi thật | UC-USR-01, UC-NTF-07/11/12 | BR-SEC-01/04/05, BR-TPL-03, BR-OPS-04 |
| SRS-F14 | Tôn trọng điều kiện ngân sách do bên sở hữu quyết định và hợp đồng được giao; khi bị chặn trả kết quả có lý do, không tự vượt hạn mức | UC-NTF-06 | BR-COST-01/03 |

SRS-F14 chưa chọn Notification hay dịch vụ khác làm chủ kiểm soát ngân sách. DEC-08 xác định bên quyết định và cách thực thi; thiết kế chỉ triển khai phần trách nhiệm đã được giao.

## 3 Điều kiện tối thiểu từng kênh

| Kênh | Đầu vào/điều kiện | Kết quả cần diễn giải đúng |
|---|---|---|
| In-app | Tài khoản và ngữ cảnh được xác nhận; bản thông báo lưu đúng phạm vi | Khả dụng để xem và đã đọc là hai trạng thái riêng. |
| Push | Đích thiết bị/client đúng tài khoản/app do bên sở hữu cấp; cấu hình kênh hợp lệ | Nhà cung cấp nhận không chứng minh người đọc; token không hợp lệ được xử lý theo hợp đồng thiết bị. |
| Email | Địa chỉ và định danh gửi phù hợp; mẫu được duyệt | Tiếp nhận và bounce/giao nếu có bằng chứng được phân biệt; không lộ địa chỉ các người nhận khác. |
| SMS | Số điện thoại và định danh gửi phù hợp | Chỉ nâng kết quả đến mức được bằng chứng kênh xác nhận. |
| Zalo ZNS | Liên hệ, mẫu và định danh gửi đáp ứng cấu hình đã duyệt | Kết quả và dự phòng theo hợp đồng; không suy rằng timeout là thất bại. |

## 4 Kết quả và trạng thái đọc

Trạng thái yêu cầu tổng hợp từ kết quả từng bản/kênh theo tiêu chí của trường hợp gửi. Các kết quả chưa rõ, một phần, không gửi, hủy và hết hạn phải giải thích được. In-app đã đọc có thể cùng tồn tại với SMS đang chờ; hành động đã hoàn tất do nguồn quyết định. Danh sách chuyển trạng thái và quy tắc callback đến muộn được chốt tại [hợp đồng](<Notitek - HỢP ĐỒNG TÍCH HỢP API VÀ SỰ KIỆN.md>) và DEC-10.

## 5 Quyền và dữ liệu

Backend xác định người và ngữ cảnh từ cơ chế User/Authorization. ID và permissions do frontend gửi không cấp quyền. Xử lý nền dùng căn cứ/hợp đồng được phép, không cần phiên của người nhận. Các hành động đọc, tra cứu, cấu hình, kích hoạt, hủy hoặc gửi lại phải có catalog quyền và scope được chốt ở DEC-12.

Tách dữ liệu bí mật ngắn hạn, nội dung thông báo, metadata kết quả và audit. Nếu xử lý cần lưu bí mật tạm thời, thiết kế phải có bảo vệ, giới hạn truy cập và thời điểm xóa theo DEC-09; không đưa payload chứa bí mật nguyên văn vào log, dead-letter hoặc báo cáo hỗ trợ.

## 6 Yêu cầu phi chức năng cần định lượng

| Mã | Điều cần đo/kiểm chứng | Tham số cần chốt | Bên chốt và DEC |
|---|---|---|---|
| SRS-N01 | Thời gian từ yêu cầu OTP tới phản hồi Notification; tách khỏi thời gian người dùng nhận/nhập mã | Thời gian chờ, p95/p99 phản hồi, hạn gửi | User, Notification, vận hành; DEC-05/06. |
| SRS-N02 | Độ trễ từ sự việc được nguồn xác nhận tới tiếp nhận và đưa ra kênh | Chỉ tiêu theo loại tin, cách đo khi event đến trễ | PO Order, Tech Lead; DEC-06. |
| SRS-N03 | Tải gửi và tải hộp tin theo luồng/kênh/người/phạm vi | Trung bình, đỉnh, số người/kênh mỗi yêu cầu, giới hạn và backpressure | PO, Tech Lead, vận hành; DEC-06. |
| SRS-N04 | Khả năng phục hồi sau lỗi, không mất yêu cầu đã tiếp nhận hoặc gửi lại phần thành công | RPO/RTO, cửa sổ chống trùng và thời gian xử lý tồn đọng | Tech Lead, vận hành; DEC-06/09. |
| SRS-N05 | Tách phạm vi, kiểm tra quyền và không lộ bí mật | Ca giả mạo, mất quyền, đổi tài khoản/app; chính sách truy cập dữ liệu | User/Authorization, chủ dữ liệu; DEC-03/09/12. |
| SRS-N06 | Lưu/xóa, audit và khả năng tra cứu | Thời gian theo loại dữ liệu, thời gian đáp ứng tra cứu, quyền xuất | Chủ dữ liệu, PO, vận hành; DEC-09. |
| SRS-N07 | Khả năng đưa thêm luồng/module qua hợp đồng và cấu hình | Công sức tích hợp, số thay đổi lõi, quy tắc tương thích | PO, Tech Lead; DEC-02/07. |

Chưa có con số SLA hoặc tải được duyệt. Story về hiệu năng/lưu dữ liệu thật chỉ được nghiệm thu khi bảng có giá trị, đơn vị đo, môi trường, cửa sổ đo và tiêu chí đạt cụ thể.

## 7 Điều kiện chuyển sang thiết kế và code

Tech Lead xác nhận các yêu cầu của lát cắt, chốt hợp đồng cùng đội nguồn, hoàn thiện mô hình dữ liệu/trạng thái và phương án lỗi. Frontend cần [hành vi màn hình](<Notitek - ĐẶC TẢ HÀNH VI GIAO DIỆN (UX).md>) cùng ví dụ dữ liệu. QA lấy [AT và ma trận truy vết](<Notitek - BACKLOG TRIỂN KHAI VÀ KỊCH BẢN NGHIỆM THU.md>) để bổ sung ca kỹ thuật theo schema thật.
