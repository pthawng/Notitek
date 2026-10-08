# Notitek - ĐẶC TẢ HÀNH VI GIAO DIỆN (UX)

**Phiên bản:** 0.2. **Ngày:** 08/10/2026. **Trạng thái:** Đặc tả hành vi để PO, UX và frontend hoàn thiện.

Trải nghiệm đợt đầu tiếp nối chuông thông báo Shop, phân biệt cảnh báo tài khoản và thông báo công việc, đồng thời cung cấp tra cứu cho vận hành. Các hành động nghiệp vụ tiếp tục do User/Order xử lý. Wireframe và hợp đồng dữ liệu cần hoàn thiện trước triển khai UI.

## 1 Quy tắc ngữ cảnh

Cảnh báo tài khoản gắn với đúng chủ tài khoản và chính sách app được xác nhận, không tự gắn vào Shop đang mở. Thông báo công việc gắn với Shop/tổ chức/app của sự việc. Khi đổi Shop, danh sách công việc chuyển theo ngữ cảnh; cảnh báo tài khoản vẫn thuộc tài khoản đó. Khi đổi tài khoản/đăng xuất, dữ liệu cũ và cache liên quan phải được tách/loại theo thiết kế.

Backend kiểm tra quyền và ngữ cảnh. Giao diện dùng ngữ cảnh để trình bày và gửi yêu cầu, không tự cấp quyền từ bộ lọc hoặc ID. Phạm vi xuất hiện của cảnh báo tài khoản trên nhiều app cần hợp đồng User, không mặc định mọi app đều xem được.

## 2 Các màn hình của đợt đầu

| Mã | Màn hình/thành phần | Actor | Hành vi và nguồn dữ liệu | Story |
|---|---|---|---|---|
| UX-01 | Chuông và danh sách thông báo | Người dùng Shop | Hiển thị số chưa đọc và danh sách đúng ngữ cảnh; tiếp nối thành phần hiện có; nguồn Notification và cảnh báo User theo hợp đồng chuyển đổi | ST-05/08 |
| UX-02 | Chi tiết thông báo công việc | Người dùng có quyền trong Shop | Sự việc, đối tượng, nội dung và hành động mở đơn; đọc tách với hành động; Order kiểm tra quyền khi mở | ST-05 |
| UX-03 | Chi tiết cảnh báo/yêu cầu thiết bị | Chủ tài khoản | Hiển thị trạng thái, khả năng thao tác do User cấp; xác nhận/báo không nhận ra gọi User; cập nhật kết quả từ nguồn | ST-08 |
| UX-04 | Tra cứu yêu cầu và kết quả gửi | Vận hành có quyền | Tìm theo reference; xem các trường hợp gửi, bản/kênh/lần thử và lý do; nội dung che theo quyền | ST-09 |
| UX-05 | Cấu hình tối thiểu của luồng | Vận hành/quản trị được cấp quyền | Xem luồng, mẫu, phạm vi, kênh, phiên bản và lỗi trước kích hoạt; phương thức vận hành có thể chưa là màn hình riêng | ST-01 |

UX-05 không tự yêu cầu xây trình thiết kế workflow hoặc trang quản trị đầy đủ ở đợt đầu. Nếu dùng cấu hình do vận hành quản lý, vẫn phải có quyền, phiên bản, kiểm tra và audit theo hợp đồng.

Trạng thái đọc của In-app mới lấy backend làm nguồn sự thật theo SRS-F08.05; local cache chỉ hỗ trợ trình bày. Đánh dấu đọc được lưu theo người/bản tin và hiển thị ở lần tải sau trên thiết bị được phép khác; lỗi lưu phải có phản hồi, không tự báo thành công. Tin từ feed User cũ cần chính sách chuyển đọc riêng, không suy đã đọc từ việc yêu cầu bảo mật đã xử lý.

## 3 Trạng thái màn hình

| Trạng thái | Danh sách/hộp tin | Chi tiết/hành động | Tra cứu vận hành |
|---|---|---|---|
| Đang tải | Hiển thị trạng thái tải, không dùng dữ liệu scope cũ | Giữ thông tin tải, chưa cho hành động khi thiếu điều kiện | Hiển thị tiến trình tải. |
| Trống | Cho biết chưa có tin trong ngữ cảnh hiện tại | Không tự tạo nội dung giả | Phản hồi an toàn khi không tìm thấy/không có quyền; không tiết lộ tin ngoài phạm vi. |
| Lỗi tải | Lý do dễ hiểu và retry | Không báo hành động thành công khi chưa có phản hồi | Hiển thị lỗi tra cứu; không suy kết quả gửi. |
| Mất quyền | Không lộ nội dung công việc; cập nhật theo kết quả backend | Nguồn từ chối hành động phù hợp | Từ chối đúng phạm vi; không mở rộng dữ liệu bằng quyền xem khác. |
| Đã đọc | Bỏ đánh dấu chưa đọc theo trạng thái hộp tin | Không đổi đơn/yêu cầu bảo mật thành đã xử lý | Giữ đọc riêng với tiến trình kênh. |
| Nguồn đã xử lý/hết hiệu lực | Tin có thể giữ theo chính sách lịch sử/quyền | Hiển thị trạng thái nguồn; bỏ thao tác không còn hợp lệ | Cho biết hủy/hết hạn hoặc kết quả nguồn được phép hiển thị. |
| Chưa rõ kết quả | Không hứa người dùng đã nhận qua kênh khác | Khi là hành động nguồn, hiển thị trạng thái do nguồn trả | Hiển thị chưa rõ và hướng xử lý theo hợp đồng. |

## 4 Nội dung và hành động

Thông báo nói rõ sự việc, đối tượng người dùng nhận biết và bước tiếp theo khi cần. Nội dung Push/preview không lộ bí mật hoặc dữ liệu nhạy cảm ngoài chính sách. Link do nguồn cung cấp và được kiểm tra hợp lệ; Notification không ghép link để bỏ qua quyền của nguồn.

Email lời mời dẫn tới chức năng User. Thông báo giao thất bại dẫn tới đơn hoặc hướng liên hệ do Order/chủ nội dung cấp. Cảnh báo thiết bị tiếp tục dùng hành động User; nút cho phép đăng nhập hoặc báo không nhận ra không trở thành API nghiệp vụ của Notification.

## 5 Tiếp nối thành phần hiện có

[NotificationBell](../../shop-fe/apps/business-web/features/notification-center/components/NotificationBell.tsx) và [useNotificationCenter](../../shop-fe/apps/business-web/features/notification-center/use-notification-center.ts) hiện lấy feed cảnh báo thiết bị và thực hiện hành động qua User. Trạng thái đọc đang có phần lưu theo phiên trình duyệt. Thiết kế hộp tin mới phải xác định ánh xạ ID, nguồn dữ liệu và chuyển trạng thái đọc để không tạo hai bản cảnh báo hoặc làm mất khả năng xử lý hiện tại.

Các toast/modal báo kết quả thao tác của frontend là phản hồi trên màn hình. Chỉ đưa vào hộp thông báo khi có sự việc/yêu cầu gửi riêng được chủ nghiệp vụ xác nhận; không chuyển mọi lỗi HTTP hoặc thông báo thành công của UI thành tin bền vững.

## 6 Artifact UX cần hoàn thiện

UX/frontend bàn giao wireframe hoặc prototype cho UX-01 đến UX-04, mô tả ngữ cảnh tài khoản/công việc, trạng thái màn hình và hành động. Đính kèm dữ liệu mẫu từ hợp đồng, ánh xạ thành phần hiện có và ca kiểm thử đổi Shop/app/tài khoản/mất quyền. PO duyệt trải nghiệm; User/Order xác nhận vị trí và ngữ nghĩa các hành động nguồn.
