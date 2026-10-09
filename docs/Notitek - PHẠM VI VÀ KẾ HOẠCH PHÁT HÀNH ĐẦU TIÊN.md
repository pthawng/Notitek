# Notitek - PHẠM VI VÀ KẾ HOẠCH PHÁT HÀNH ĐẦU TIÊN

**Phiên bản:** 1.0. **Ngày:** 08/10/2026. **Chủ trì:** PO Notification.

Bản phát hành đầu tiên cung cấp năng lực thông báo dùng chung cho User và Order, tiếp nối trải nghiệm Shop và ứng dụng nội bộ. Mục tiêu là các luồng đã chọn hoạt động từ nghiệp vụ thật đến người nhận và tra cứu kết quả, đồng thời kiểm chứng khả năng hỗ trợ đủ năm kênh của SuperPlatform.

## 1 Quyết định sản phẩm đã xác nhận

Người yêu cầu xác nhận định hướng PO trong cuộc trao đổi ngày 08/10/2026. Phạm vi đã chốt gồm:

| Mã | Quyết định |
|---|---|
| PO-01 | Notification quản lý năng lực thông báo dùng chung: luồng, nội dung, điều phối kênh, hộp tin và kết quả gửi. Module nguồn quyết định sự việc, đối tượng nghiệp vụ và hành động. |
| PO-02 | Đợt đầu hỗ trợ In-app, Push, Email, SMS và Zalo ZNS. Các lát cắt được làm lần lượt; nghiệm thu đợt đầu vẫn đủ năm kênh. |
| PO-03 | Chọn OTP, lời mời nhân viên/thành viên, giao thất bại và cảnh báo bảo mật/thiết bị làm các lát cắt đầu tiên. |
| PO-04 | Tách thông báo tài khoản với thông báo công việc theo Shop/tổ chức/app. Đổi Shop không làm mất cảnh báo tài khoản; thông báo đơn không lộ sang Shop khác. |
| PO-05 | Dùng lại định danh, tư cách, Application/Client và phân quyền của User/Authorization; nối chuông thông báo Shop hiện có vào hộp tin Notitek. |
| PO-06 | Các phụ thuộc User, Order, frontend và vận hành trở thành công việc được giao trong cùng kế hoạch tích hợp. |
| PO-07 | Hoàn thành phải có bằng chứng từ đầu đến cuối, gồm gửi thật, mở đúng hành động, truy vết và xử lý ngoại lệ. |
| PO-08 | Đã loại bỏ (09/10/2026): Notitek là hệ thống mới, không có hệ thống thông báo/ZNS cũ cần kiểm kê hoặc chuyển đổi. OA và mẫu ZNS được đăng ký mới. |
| PO-09 | Đo khả năng hoàn thành xác thực, tiếp nhận lời mời, thời gian xử lý vấn đề giao hàng, chất lượng gửi và công sức tích hợp luồng mới. |

Các quyết định này chốt phạm vi và cách tổ chức sản phẩm. Chúng không tự xác nhận tên event, người nhận vận hành, nhà cung cấp, mẫu kênh, SLA hoặc ngân sách cụ thể.

## 2 Các lát cắt của đợt đầu

| Lát cắt | UC nghiệp vụ | Người hưởng lợi | Kết quả cần đạt | Kênh và phụ thuộc |
|---|---|---|---|---|
| R1-OTP | UC-USR-01 | Người xác thực tài khoản/thao tác | Nhận mã qua đường gửi thật; User nhận đúng mức kết quả để quyết định luồng xác thực | SMS cho lát cắt đầu; User sở hữu mã, reference, đích và hạn gửi. Chính sách phản hồi/thử lại phải chốt. |
| R1-INVITE | UC-USR-03, UC-USR-04 | Người mời và người được mời | Lời mời đã lưu được gửi bằng email; link mở đúng chức năng User; gửi lỗi không xóa lời mời | Email; User cấp link, thời hạn và dữ liệu lời mời. |
| R1-ORDER | UC-ORD-05 | Shop và người nhận hàng | Nhận đúng thông tin giao thất bại và hướng xử lý, theo các trường hợp gửi khác nhau | In-app/Push cho Shop; ZNS và SMS dự phòng cho liên hệ đơn theo chính sách được duyệt. Order xác nhận sự việc và tin cần dừng. |
| R1-SECURITY | UC-USR-09 | Chủ tài khoản | Xem cảnh báo/yêu cầu về thiết bị và thực hiện hành động tại User | Tiếp nối chuông Shop; User tiếp tục sở hữu yêu cầu và quyết định bảo mật. Phương án đưa dữ liệu vào hộp tin chung cần hợp đồng. |

UC-NTF-01 đến UC-NTF-13 là các năng lực nền liên quan; cấu hình tối thiểu của UC-NTF-15 phải đủ để chạy các lát cắt. Màn hình tự phục vụ nâng cao, gom nhóm, báo cáo chi phí nâng cao và chiến dịch vẫn đi theo phân kỳ riêng. P0 ở danh mục dài hạn không đồng nghĩa mọi UC P0 đều thuộc bản phát hành này.

SuperShip Shop và giao diện nội bộ là điểm tích hợp ưu tiên. Client Push cụ thể và việc áp dụng cho app khác cần xác định trước nghiệm thu Push; không suy ra đã có ứng dụng mobile chỉ từ danh mục app.

## 3 Ranh giới trách nhiệm

| Phần | Chủ quyết định nghiệp vụ | Notification thực hiện |
|---|---|---|
| OTP | User phát hành, kích hoạt, xác minh, cấp lại và vô hiệu hóa mã | Gửi theo hợp đồng chung, hạn gửi và chỉ dẫn thử lại/hủy; trả kết quả gửi. |
| Lời mời | User tạo, cấp lại, thu hồi, xác nhận tư cách | Gửi mẫu/link do nguồn cấp; dừng bản chờ theo chỉ dẫn hoặc hạn gửi. |
| Giao thất bại, giao đủ hay giao một phần | Order xác nhận sự việc, đối tượng liên quan và việc còn cần nhắc | Ánh xạ sự việc vào luồng/mẫu, điều phối kênh và hủy đúng tham chiếu. |
| Quyền và người nhận | Nguồn/User/Authorization cung cấp căn cứ | Áp dụng căn cứ được xác nhận, không tạo quyền hoặc suy từ cây tổ chức. |
| Thiết bị Push | App/User/năng lực thiết bị theo hợp đồng được chọn | Nhận hoặc tra đích gửi được xác nhận, gửi đúng app và xử lý phản hồi token không hợp lệ. |
| Ngân sách | Bên sở hữu ngân sách cần chốt tại DEC-08 | Chỉ thực thi điều kiện được giao; không tự tính sổ sách tài chính. |
| Hành động trong thông báo | User/Order kiểm tra quyền và thực hiện | Dẫn tới đúng chức năng; trạng thái đọc không hoàn tất nghiệp vụ. |

## 4 Các quyết định theo từng lát cắt

Ngày 08/10/2026, PO đã chốt cả 12 DEC dưới đây tại [Quyết định PO cho các DEC](<Notitek - QUYẾT ĐỊNH PO CHO CÁC DEC.md>). Bảng giữ nội dung câu hỏi và phần bị chặn để truy vết; các điểm cần pháp chế, tài chính, Tech Lead nguồn hoặc vận hành ký được theo dõi tại mục 4 tài liệu đó và chỉ chặn nghiệm thu phần liên quan.

| BRD DEC | Nội dung cần đầu ra cụ thể | Bên chịu trách nhiệm | Phần bị chặn |
|---|---|---|---|
| DEC-01 | Client Push và app mở rộng ngoài các giao diện ưu tiên | PO, chủ app, frontend | Nghiệm thu Push trên client được chọn và mở rộng app. |
| DEC-02 | Hợp đồng gửi/tra cứu/hủy, event giao thất bại và đường phát outbox | Tech Lead User, Order, Notification | Tích hợp backend các lát cắt tương ứng. |
| DEC-03 | Người nhận Shop, ngữ cảnh tài khoản/công việc, căn cứ cho xử lý nền | PO Order, User/Authorization, Notification | Phát tin Shop và hộp tin theo phạm vi. |
| DEC-04 | Mục đích, luồng bắt buộc, kênh được dùng và dữ liệu lựa chọn nhận tin | PO các nguồn, chủ chính sách | Gửi thật theo chính sách; Marketing chưa thuộc đợt đầu. |
| DEC-05 | Mức kết quả User cần, thời gian chờ và chỉ dẫn thử lại OTP | User, Notification, vận hành | Nghiệm thu R1-OTP. |
| DEC-06 | Hạn gửi, số lần thử, thời gian xử lý và mục tiêu độ trễ/sản lượng | PO, vận hành, Tech Lead | Nghiệm thu thời gian và tải. |
| DEC-07 | Cấu hình hiệu lực, xung đột phạm vi và ảnh hưởng bản chờ | PO nền tảng, Notification | Kích hoạt cấu hình nhiều phạm vi. |
| DEC-08 | Chủ sở hữu ngân sách, bên kiểm soát, hạn mức và phản hồi khi bị chặn | Chủ ngân sách, PO, Tech Lead | Gửi thật qua kênh phát sinh phí. |
| DEC-09 | Phân loại dữ liệu, thời gian lưu, quyền xem/xóa và audit | Chủ dữ liệu, PO, vận hành | Lưu dữ liệu thật và xuất lịch sử. |
| DEC-10 | Bằng chứng từng kênh, kết quả tổng hợp và điều kiện dự phòng | Notification, chủ nguồn, vận hành | Nghiệm thu kết quả và chuyển kênh. |
| DEC-11 | Đã hủy: không có hệ thống cũ cần chuyển đổi | — | — |
| DEC-12 | Hành động Notification trong catalog quyền và người được cấp | User/Authorization, PO | Truy cập hộp tin/quản trị/vận hành. |

## 5 Các mốc hoàn thành

| Mốc | Nội dung | Bằng chứng để đi tiếp |
|---|---|---|
| M0 | Hoàn thiện UC, hợp đồng, trải nghiệm và thiết kế cho lát cắt đầu | Story đạt điều kiện bắt đầu trong [backlog](<Notitek - BACKLOG TRIỂN KHAI VÀ KỊCH BẢN NGHIỆM THU.md>); bộ ví dụ hợp đồng thống nhất. |
| M1 | OTP và lời mời User chạy từ đầu đến cuối | SMS/Email thực, phản hồi đúng, link tới User, bí mật không lộ và lỗi không đảo nghiệp vụ đã commit. |
| M2 | Giao thất bại và cảnh báo thiết bị trong trải nghiệm đã chọn | Có đường phát Order, In-app/Push/ZNS thực, đúng phạm vi, hành động vẫn do nguồn xử lý. |
| M3 | Nghiệm thu đợt đầu | Đủ năm kênh, ngoại lệ được kiểm chứng, tra cứu được, có quy trình vận hành. |

## 6 Chỉ số sản phẩm

User cung cấp kết quả xác thực/tiếp nhận lời mời; Order cung cấp kết quả xử lý đơn. Notification cung cấp tiến trình gửi, kết quả theo bằng chứng, lỗi/chưa rõ, thời gian xử lý và lần thử. PO kết hợp các nguồn để đo giá trị; không dùng tỷ lệ đọc tin làm tỷ lệ hoàn thành nghiệp vụ. Mốc đo ban đầu và mục tiêu định lượng được chốt ở DEC-05/06/08 trước nghiệm thu phần liên quan.

## 7 Phạm vi mở rộng

Tài chính, Support, Pricing, chính sách/sự cố và Marketing được giữ trong phạm vi dài hạn của SuperPlatform. Đội ghi nhu cầu/hợp đồng dự kiến để kiểm tra khả năng mở rộng; các module này chưa tự trở thành cam kết tích hợp trong đợt đầu. Việc bổ sung một luồng thông báo từ nguồn mới phải được đánh giá theo công sức tích hợp và lượng thay đổi ở lõi Notification.
