# Notitek - THIẾT KẾ KỸ THUẬT SƠ BỘ

**Phiên bản:** 0.2. **Ngày:** 08/10/2026. **Trạng thái:** Mô hình logic và hạng mục thiết kế để Tech Lead chốt.

Thiết kế hiện xác định ranh giới xử lý và dữ liệu cần giải thích các use case. Nó chưa chọn stack, engine điều phối, nhà cung cấp, broker, schema vật lý hoặc kiến trúc triển khai. Tech Lead hoàn thiện các quyết định dưới đây từ [SRS](<Notitek - ĐẶC TẢ YÊU CẦU PHẦN MỀM (SRS).md>) và [hợp đồng](<Notitek - HỢP ĐỒNG TÍCH HỢP API VÀ SỰ KIỆN.md>) trước khi code phần phụ thuộc.

## 1 Các trách nhiệm kỹ thuật

| Phần xử lý | Trách nhiệm |
|---|---|
| Tiếp nhận | Xác thực nguồn/ngữ cảnh, kiểm tra hợp đồng, reference và ghi nhận bền vững. |
| Áp dụng cấu hình | Chọn luồng/trường hợp gửi, mẫu, phiên bản và chính sách theo phạm vi đã duyệt. |
| Chuẩn bị bản tin | Dùng người nhận và dữ liệu tối thiểu được nguồn cấp, bảo vệ biến nhạy cảm và link. |
| Điều phối gửi | Kiểm tra hạn gửi/chỉ dẫn hủy, gọi kênh, ghi lần thử và xử lý phản hồi. |
| Kết quả | Cập nhật theo bằng chứng, tổng hợp theo tiêu chí luồng, cung cấp tra cứu/event theo quyền. |
| In-app | Danh sách và đọc đúng người/ngữ cảnh, tách tiến trình gửi và hành động nguồn. |
| Vận hành | Tra cứu, cấu hình tối thiểu, audit và theo dõi tồn đọng/lỗi/chưa rõ. |

Các phần này là trách nhiệm logic, không bắt buộc trở thành các microservice riêng. Tech Lead chọn cách tổ chức phù hợp với giải pháp và tải được xác nhận.

## 2 Mô hình dữ liệu logic

| Đối tượng | Điều cần biểu diễn | Ranh giới |
|---|---|---|
| Yêu cầu/sự việc tiếp nhận | Nguồn, reference, loại/phiên bản, ngữ cảnh, correlation, kết quả tiếp nhận | Không là aggregate User/Order mới. |
| Luồng và trường hợp gửi | Phạm vi, mục đích, người nhận/căn cứ, mẫu, kênh và chỉ dẫn gửi | Có phiên bản cấu hình hiệu lực. |
| Bản thông báo | Nội dung được phép lưu cho một người trong một trường hợp, ngữ cảnh và hạn gửi | Không sao chép toàn hồ sơ nguồn hoặc lưu bí mật vào lịch sử. |
| Lần thử gửi | Kênh, nhà cung cấp, reference bên ngoài, thời điểm, bằng chứng/lỗi | Thử lại không là sự việc nguồn mới. |
| Kết quả và bằng chứng | Mức áp dụng, trạng thái, lý do, thời điểm, callback liên quan | Đọc In-app và hành động nghiệp vụ được tách. |
| Trạng thái đọc | Người/bản In-app, thời điểm hoặc dấu đọc phù hợp thiết kế | Kiểm tra scope/quyền; không đổi đơn/yêu cầu bảo mật. |
| Tham chiếu đích gửi | Định danh User/nguồn, app/client và đích tối thiểu được phép dùng | Bên sở hữu hồ sơ/thiết bị gốc được xác định bằng hợp đồng. |
| Audit vận hành | Ai làm gì, phạm vi, lúc nào, lý do và phiên bản liên quan | Không chứa OTP/token nguyên văn. |

## 3 Các phương án cần ghi quyết định kiến trúc

| Mã | Quyết định | Đầu vào và tiêu chí đánh giá |
|---|---|---|
| ADR-01 | Tự xây hay dùng lại giải pháp điều phối; cách mở rộng kênh | Năm kênh, phạm vi app/Shop, bảo vệ dữ liệu, khả năng vận hành và công sức tích hợp. Tham khảo Novu không tự thành quyết định dùng Novu. |
| ADR-02 | Phương thức nhận lệnh và phát/nhận event bền vững | Hợp đồng User/Order, replay, chống trùng, cửa sổ lỗi; không mặc định phải có broker cụ thể. |
| ADR-03 | Cách lưu và bảo vệ dữ liệu, gồm bí mật tạm thời | DEC-09, thời gian xử lý/thử lại, yêu cầu tra cứu, audit và xóa. |
| ADR-04 | Mô hình kết quả, chuyển trạng thái và tổng hợp | DEC-10, callback trễ/trùng/xung đột, kết quả một phần, In-app và đọc. |
| ADR-05 | Cách nhận/tra và thu hồi đích Push | Chủ sở hữu thiết bị, client được chọn, đổi chủ thiết bị và phản hồi token lỗi. |
| ADR-06 | Cách kiểm tra quyền/ngữ cảnh tương tác và xử lý nền | Hợp đồng User/Authorization; không dùng token người tạo thay quyền tập nhận. |
| ADR-07 | Phương án chuyển feed cảnh báo User và ZNS cũ | Giữ hành động nguồn, ánh xạ ID, điểm chuyển trách nhiệm gửi và cơ chế khôi phục. |

ADR cần ghi lựa chọn, phương án đã cân nhắc, hệ quả, người quyết định và ngày. Chỉ đóng khi đã có đủ đầu vào; không dùng mô hình logic này thay phê duyệt kỹ thuật.

## 4 Các cửa sổ lỗi phải thiết kế

- Nguồn commit rồi đường phát lỗi: có thể phát lại sự việc, bên nhận chống trùng.
- Notification đã nhận nhưng phản hồi về nguồn mất: tra cứu cùng reference, không tạo yêu cầu mới vô ích.
- Kênh đã nhận nhưng Notification timeout/crash: giữ/khôi phục kết quả chưa rõ, không gửi lại mù quáng.
- Một kênh thành công, kênh khác lỗi: chỉ xử lý phần còn hợp lệ.
- Hủy và gửi diễn ra gần nhau: chỉ dừng phần còn kiểm soát được, phản hồi rõ phần đã ra ngoài.
- User gửi mã nhưng chưa kích hoạt/commit thành công: User quyết định trạng thái mã; Notification tuân thủ chỉ dẫn thử lại/hủy theo hợp đồng, không đọc session xác thực.
- Quyền/người nhận thay đổi: dùng căn cứ cập nhật theo hợp đồng trước gửi/xem, không tự suy từ bản đọc cũ.
- Callback đến sai thứ tự: giữ bằng chứng, quy tắc cập nhật không làm kết quả bị đổi sai.

## 5 Artifact kỹ thuật trước code

Tech Lead bàn giao sơ đồ thành phần và tương tác cho các lát cắt, OpenAPI/schema event, bảng chuyển trạng thái, ERD/từ điển dữ liệu, migration, thiết kế bảo vệ bí mật và quyết định kiến trúc liên quan. Cùng QA xác định kiểm thử hợp đồng, phục hồi lỗi và phân quyền; cùng vận hành chốt cấu hình môi trường, quan sát và hướng dẫn xử lý tồn đọng/chưa rõ.

## 6 Nguyên tắc mở rộng

Các module nguồn tích hợp qua hợp đồng; lõi Notification không phụ thuộc enum trạng thái Order, session OTP hoặc sổ tài chính. Luồng mới ưu tiên bổ sung hợp đồng/cấu hình/mẫu; thay lõi phải được giải thích bằng nhu cầu năng lực gửi dùng chung. Kết nối kênh không cấp quyền User và cấu hình brand không vượt ranh giới dữ liệu.

## 7 Đầu vào thiết kế sau PO review SRS 1.0

Baseline SRS đã phân rã 22 nhóm chức năng. Thiết kế cần đáp ứng yêu cầu quan sát được, còn stack/provider/schema vật lý và các ADR vẫn do Tech Lead chốt.

| Vấn đề thiết kế | Bằng chứng thiết kế và kiểm chứng cần có |
|---|---|
| Cửa sổ lệnh đồng bộ | Cách truyền/áp hạn xử lý, dừng phần chưa ra kênh khi hết cửa sổ; không biến timeout OTP thành retry nền; AT-04/23. |
| Ghi nhận và chống trùng | Ranh giới accepted, xử lý đồng thời, phục hồi sau crash, bảo vệ dấu đối chiếu bí mật, vòng đời khóa và dấu hủy; AT-02/22/23/28. |
| Kết quả | Mô hình riêng cho bằng chứng kênh, mức đạt mục tiêu và đọc; bảng chuyển trạng thái trễ/xung đột; phát lại kết quả không phát lại tin; AT-25. |
| In-app | Nguồn trạng thái đọc backend, scope hiện tại, pagination/đếm và cách tách cache frontend; chuyển trạng thái đọc feed cũ có chính sách; AT-08/15/26. |
| Đích Push | Contract chủ sở hữu, version/thu hồi/đổi chủ, độ mới và phản hồi token lỗi; AT-12/26. |
| Kết nối và bí mật | Version kết nối, xác thực callback, bảo vệ credential/biến bí mật ở cả hàng lỗi/audit/backup; AT-24/27/28. |
| Tải và vận hành | Hạn chế tải, ưu tiên, giới hạn provider, phạm vi dừng, số liệu và fault model với RPO/RTO được xác nhận; AT-27 và SRS-N01 đến SRS-N07. |

Các dòng này là đầu vào cho ADR/thiết kế, không chứng nhận đã có code hoặc đã chạy các ca AT.
