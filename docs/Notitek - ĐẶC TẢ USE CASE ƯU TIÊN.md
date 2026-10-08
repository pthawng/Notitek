# Notitek - ĐẶC TẢ USE CASE ƯU TIÊN

**Phiên bản:** 0.1. **Ngày:** 08/10/2026. **Chủ trì:** BA cùng PO Notification và PO module nguồn.

Các use case ưu tiên của đợt đầu dưới đây phân rã [phạm vi đã xác nhận](<Notitek - PHẠM VI VÀ KẾ HOẠCH PHÁT HÀNH ĐẦU TIÊN.md>) thành actor, điều kiện, luồng chính, ngoại lệ và hậu điều kiện. Đây là đặc tả nghiệp vụ cho đợt đầu; schema, route và công nghệ thuộc [hợp đồng](<Notitek - HỢP ĐỒNG TÍCH HỢP API VÀ SỰ KIỆN.md>) và [thiết kế kỹ thuật](<Notitek - THIẾT KẾ KỸ THUẬT SƠ BỘ.md>).

## 1 Quy tắc dùng chung

- Module nguồn quyết định sự việc, người nhận theo giao dịch và các yêu cầu cần dừng. Notification sử dụng chỉ dẫn gửi/hủy và dữ liệu nguồn; không đọc mô hình nghiệp vụ riêng để suy kết quả.
- Một yêu cầu có thể tạo nhiều trường hợp gửi, bản thông báo và lần thử. Kết quả tra cứu giữ rõ các mức này.
- Yêu cầu không hợp lệ bị từ chối; tiếp nhận bền vững chưa phải gửi thành công. Timeout có thể là chưa rõ; dự phòng phải theo điều kiện được duyệt.
- Hành động sau khi mở thông báo do nguồn xác thực và kiểm tra quyền. Đã đọc chỉ là trạng thái trải nghiệm hộp tin.
- Người vận hành xem dữ liệu theo quyền; OTP/token không hiển thị trong lịch sử hoặc log.
- Chi tiết còn mở được liên kết với DEC trong [phạm vi](<Notitek - PHẠM VI VÀ KẾ HOẠCH PHÁT HÀNH ĐẦU TIÊN.md>); không tự chọn người nhận, mẫu hoặc hạn mức bằng dữ liệu giả định.

## 2 Gửi mã xác thực do User phát hành

**Mã:** UC-USR-01. **Actor chính:** User; người hưởng lợi là người đang xác thực.

**Kích hoạt:** User yêu cầu gửi mã đã phát hành tới liên hệ được chọn. **Tiền điều kiện:** Nguồn được phép; yêu cầu có reference, đích, mẫu/biến và hạn gửi; kênh và chính sách áp dụng hợp lệ.

**Luồng chính:**

1. User cấp yêu cầu gửi cùng chỉ dẫn thời hạn và khả năng thử lại theo hợp đồng.
2. Notification kiểm tra hợp đồng, chống trùng và điều kiện gửi được giao.
3. Notification dựng nội dung từ mẫu phù hợp, gửi SMS theo lát cắt đầu và ghi kết quả ở mức được kênh xác nhận.
4. Notification trả kết quả/định danh tra cứu theo hợp đồng. User quyết định kích hoạt mã, tiếp tục xác thực hoặc hiển thị phương án tiếp theo.
5. Người dùng nhập mã tại User; Notification không xác minh mã.

**Ngoại lệ:** Yêu cầu hết hạn trước gửi được dừng; lỗi kênh được trả rõ; timeout giữ kết quả chưa rõ và tuân thủ chỉ dẫn thử lại. Cùng reference không tạo lần cấp mã mới. User cấp lại mã mới bằng yêu cầu có định danh phù hợp và chỉ dẫn xử lý bản cũ; Notification không tự suy mã nào đang hợp lệ.

**Hậu điều kiện:** Có kết quả gửi có thể tra cứu; không chứa bí mật trong log/lịch sử/kết quả công khai; User vẫn sở hữu toàn bộ trạng thái xác thực.

**Truy vết:** BR-INT-02/05, BR-TRG-02/05, BR-TIME-01/02, BR-SEC-01/05, BR-OPS-04; SRS-F01/02/03/04/06/07/09/13; story ST-02; AT-01/02/03/04. **DEC:** 02, 05, 06, 08, 09, 10.

## 3 Gửi lời mời nhân viên vào tổ chức

**Mã:** UC-USR-03. **Actor chính:** User; người mời và người được mời tham gia trải nghiệm.

**Kích hoạt:** Lời mời đã được User lưu thành công. **Tiền điều kiện:** Có email/link/thời hạn và dữ liệu tổ chức được nguồn xác nhận; mẫu đã duyệt; nguồn gọi gửi sau commit.

**Luồng chính:**

1. User cung cấp yêu cầu gửi lời mời đã commit.
2. Notification chống trùng, kiểm tra điều kiện gửi và dựng mẫu Email.
3. Email cho biết tổ chức mời, mục đích và link do User cấp; không thêm quyền hoặc thay đổi lời mời.
4. Kết quả được ghi và cung cấp cho bên được phép. Người được mời mở link để xử lý tại User.

**Ngoại lệ:** Gửi lỗi không xóa lời mời. Cấp lại/thu hồi/đã dùng link do User quyết định; Notification dừng bản chờ theo reference/chỉ dẫn hoặc hạn gửi. Email bị trả lại được ghi nhận và vận hành tra cứu theo quyền. Kết quả gửi không chứng minh người được mời đã chấp nhận.

**Hậu điều kiện:** Lời mời vẫn do User quản lý; tin gửi có mẫu, phiên bản và kết quả rõ.

**Truy vết:** BR-INT-01/03/05, BR-TPL-01/02/03/05, BR-CHN-08, BR-TIME-01/02; SRS-F01/02/04/06/07/09/13; ST-03; AT-02/05/06/07. **DEC:** 02, 04, 06, 08, 09, 10.

## 4 Gửi lời mời thành viên vào Shop

**Mã:** UC-USR-04. **Actor chính:** User; người mời và người được mời.

Áp dụng luồng UC-USR-03 cho lời mời thành viên Shop. Dữ liệu mẫu gồm tên Shop, người mời và link do User cấp, đúng phạm vi lời mời. Notification không yêu cầu người nhận đã có membership trong Shop trước khi nhận email; việc tạo/xác nhận membership thuộc User.

**Ngoại lệ riêng:** Link của Shop A không bị thay bằng dữ liệu Shop B khi cùng email được mời vào hai Shop. Hai lời mời độc lập không bị hợp nhất chỉ vì cùng liên hệ.

**Hậu điều kiện và truy vết:** Theo UC-USR-03, bổ sung BR-REC-01/03 và SRS-F05/08; ST-03; AT-05/06/08. **DEC:** 02, 03, 04, 06, 08, 09.

## 5 Thông báo giao thất bại

**Mã:** UC-ORD-05. **Actor chính:** Order; người hưởng lợi là Shop và người nhận hàng.

**Kích hoạt:** Order xác nhận một sự việc giao thất bại cần thông báo. **Tiền điều kiện:** Order cấp định danh sự việc, dữ liệu dựng mẫu, chặng/lần thực hiện khi cần, liên hệ giao dịch và người nhận Shop/căn cứ đã duyệt. Tên event và đường phát phải được chốt trong hợp đồng.

**Luồng chính:**

1. Order phát sự việc chuẩn hoặc yêu cầu gửi đã xác nhận; Notification không lấy webhook NVC thô để tự kết luận.
2. Notification chọn luồng và các trường hợp gửi theo cấu hình hiệu lực.
3. Trường hợp Shop tạo In-app và Push cho tài khoản đúng Shop/app; nội dung dẫn tới đơn tại Order.
4. Trường hợp người nhận hàng gửi ZNS theo liên hệ trên đơn; SMS dự phòng chỉ chạy khi có điều kiện được duyệt.
5. Notification ghi kết quả riêng từng người/kênh và cung cấp kết quả tổng hợp phù hợp.

**Ngoại lệ:** Event trùng không nhân tin; lần giao thất bại mới có thể là sự việc mới nếu Order cấp reference tương ứng. Thiếu người nhận/dữ liệu/mẫu không được đoán. Timeout không tự chuyển SMS. Khi Order xác nhận thông báo/nhắc không còn cần thiết, Notification hủy bản chờ theo chỉ dẫn được xác thực; không tự đọc trạng thái đơn để tính lại điều kiện.

**Hậu điều kiện:** Nội dung phản ánh sự việc nguồn xác nhận; không hứa giao lại nếu nguồn chưa cấp căn cứ; kết quả theo trường hợp gửi tra được. Không cam kết thu hồi tin đã giao cho kênh bên ngoài.

**Truy vết:** BR-CAT-02, BR-INT-01/04/05, BR-TRG-02/03/05, BR-REC-01/03/04, BR-TPL-02/05, BR-CHN-03/04/05/06/07, BR-TIME-02; SRS-F01 đến SRS-F09, SRS-F13; ST-04/05/06/07; AT-02/08/09/10/11/12. **DEC:** 02, 03, 04, 06, 07, 08, 10, 11.

## 6 Xem cảnh báo và yêu cầu về thiết bị

**Mã:** UC-USR-09. **Actor chính:** Chủ tài khoản; User cung cấp dữ liệu và xử lý hành động.

**Kích hoạt:** User có cảnh báo/yêu cầu bảo mật cần hiển thị. **Tiền điều kiện:** Người dùng được xác thực; nguồn xác nhận cảnh báo thuộc tài khoản này; hộp tin phân biệt ngữ cảnh tài khoản với ngữ cảnh công việc.

**Luồng chính:**

1. User cung cấp cảnh báo/yêu cầu theo hợp đồng. Phương án nhận event hoặc tiếp nối feed User hiện có do Tech Lead chốt.
2. Chuông/hộp tin Shop hiển thị cảnh báo tài khoản và trạng thái có thể hành động do User cấp.
3. Người dùng mở tin; trạng thái đọc được ghi độc lập với trạng thái yêu cầu bảo mật.
4. Khi xác nhận hoặc báo không nhận ra thiết bị, ứng dụng gọi User. User kiểm tra điều kiện và thực hiện quyết định.
5. Giao diện lấy kết quả cập nhật từ User; Notification không tự cho phép đăng nhập hoặc thu hồi phiên.

**Ngoại lệ:** Yêu cầu đã xử lý/hết hiệu lực vẫn hiển thị kết quả phù hợp, không còn nút hành động trái điều kiện nguồn. Đổi Shop không chuyển cảnh báo sang tài khoản khác. Đăng xuất/đổi tài khoản không giữ dữ liệu người trước.

**Hậu điều kiện:** Hành động bảo mật nằm tại User; hộp tin có trạng thái đọc riêng; không phát lại cảnh báo trùng từ feed và event trong thời gian chuyển đổi.

**Truy vết:** BR-INT-01/06, BR-REC-03/06, BR-CHN-05, BR-SEC-03/06; SRS-F05/08/10/13; ST-08; AT-08/13/14. **DEC:** 02, 03, 09, 12.

## 7 Xem hộp thông báo

**Mã:** UC-NTF-13. **Actor chính:** Người dùng có tài khoản.

**Tiền điều kiện:** Người dùng và ngữ cảnh được backend xác nhận; quyền hành động phù hợp. **Luồng chính:** Tải danh sách và số chưa đọc trong phạm vi; mở/đánh dấu đọc; ghi trạng thái đọc; mở liên kết do nguồn cung cấp và để nguồn kiểm tra quyền hành động.

**Ngoại lệ:** Mất quyền không lộ nội dung công việc cũ; giả mạo người/app/Shop bị từ chối; tải lỗi có retry; đổi ngữ cảnh cập nhật danh sách đúng phạm vi. Cảnh báo tài khoản được xác định riêng, không bị gắn mặc định vào Shop đang mở. Việc đồng bộ trạng thái đọc giữa thiết bị được Tech Lead đặc tả, không dùng trạng thái yêu cầu bảo mật thay thế.

**Hậu điều kiện:** Đọc tin không hoàn thành nghiệp vụ; trạng thái đọc tách kết quả gửi từng kênh.

**Truy vết:** BR-CHN-05, BR-REC-03/04, BR-SEC-03/06, BR-CFG-05; SRS-F05/08/10; ST-05/08; AT-08/13/14/15. **DEC:** 03, 09, 12.

## 8 Tra cứu vận hành

**Mã:** UC-NTF-12. **Actor chính:** Vận hành Notification có quyền.

**Tiền điều kiện:** Có quyền tra cứu trên hành động/phạm vi; có reference nguồn hoặc định danh Notification. **Luồng chính:** Tìm yêu cầu; xem luồng, cấu hình/mẫu hiệu lực, bản tin được phép xem, lần thử và kết quả; xác định phần chưa gửi/lỗi/chưa rõ; theo dõi hướng xử lý đã được giao.

**Ngoại lệ:** Không có quyền thì từ chối; OTP/token không được hiện; tin đã hủy/hết hạn không được biến thành thành công. Gửi lại hoặc thay cấu hình là hành động riêng cần quyền/chính sách, không tự nằm trong quyền xem.

**Hậu điều kiện:** Giải thích được vì sao gửi/không gửi; truy cập nhạy cảm và thao tác thay đổi có audit theo chính sách.

**Truy vết:** BR-INT-04, BR-OPS-01/04, BR-SEC-01/02/04/05/06; SRS-F09/11/13; ST-09; AT-03/16/17. **DEC:** 08, 09, 12.
