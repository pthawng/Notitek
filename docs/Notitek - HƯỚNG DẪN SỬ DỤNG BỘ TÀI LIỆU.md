# Notitek - HƯỚNG DẪN SỬ DỤNG BỘ TÀI LIỆU

Bộ gồm 12 tài liệu, nối phạm vi sản phẩm với công việc phát triển và nghiệm thu. Phạm vi sản phẩm được người yêu cầu xác nhận ngày 08/10/2026; các chi tiết kỹ thuật có trạng thái riêng để PO, BA, Tech Lead, phát triển và QA cùng sử dụng.

## Thứ tự sử dụng

BRD → phạm vi bản phát hành → use case ưu tiên → SRS, trải nghiệm và hợp đồng tích hợp → thiết kế kỹ thuật → backlog và nghiệm thu → code.

SRS, trải nghiệm và hợp đồng được hoàn thiện cùng nhau theo từng lát cắt. Đội có thể triển khai phần đã đủ quyết định; không phải đợi mọi luồng dài hạn được đặc tả.

| Tài liệu | Trả lời câu hỏi | Chủ trì | Trạng thái hiện tại |
|---|---|---|---|
| [Hướng dẫn sử dụng bộ tài liệu](<Notitek - HƯỚNG DẪN SỬ DỤNG BỘ TÀI LIỆU.md>) | Đọc tài liệu theo thứ tự nào, ai hoàn thiện và khi nào đủ điều kiện code? | PO và BA | Mục lục của bộ 12 tài liệu hiện hành. |
| [Đặc tả yêu cầu nghiệp vụ (BRD)](<../Notitek - ĐẶC TẢ YÊU CẦU NGHIỆP VỤ (BRD).md>) | Mục tiêu, ranh giới và yêu cầu nghiệp vụ là gì? | PO và BA | BRD 0.5; phạm vi đợt đầu đã xác nhận, chính sách chi tiết còn mở. |
| [Danh mục use case](<../Notitek - DANH MỤC USE CASE.md>) | Có những use case nào và ưu tiên ra sao? | PO và BA | 38 UC; các lát cắt đợt đầu đã xác nhận. |
| [Phạm vi và kế hoạch phát hành đầu tiên](<Notitek - PHẠM VI VÀ KẾ HOẠCH PHÁT HÀNH ĐẦU TIÊN.md>) | Đợt đầu phục vụ ai, làm những luồng nào, hoàn thành đến mức nào? | PO | Định hướng và các lát cắt đã xác nhận; bảng quyết định còn mở được giữ riêng. |
| [Quyết định PO cho các DEC](<Notitek - QUYẾT ĐỊNH PO CHO CÁC DEC.md>) | 12 câu hỏi mở được chốt thế nào, tham số dùng giá trị nào, còn chờ ai ký? | PO | PO đã chốt DEC-01 đến DEC-12 và SRS-P01 đến SRS-P12 ngày 08/10/2026; 8 điểm chờ bên sở hữu xác nhận. |
| [Đặc tả use case ưu tiên](<Notitek - ĐẶC TẢ USE CASE ƯU TIÊN.md>) | Actor thực hiện gì, luồng chính và ngoại lệ ra sao? | BA cùng PO các module nguồn | Đã phân rã các UC ưu tiên của đợt đầu; người nhận, mẫu và điều kiện cụ thể cần chốt theo bảng quyết định. |
| [Đặc tả yêu cầu phần mềm (SRS)](<Notitek - ĐẶC TẢ YÊU CẦU PHẦN MỀM (SRS).md>) | Hệ thống phải có hành vi và chất lượng nào để đáp ứng UC? | BA và Tech Lead | SRS 1.0: baseline nghiệp vụ đã PO review và chốt; 22 nhóm chức năng, NFR và kiểm chứng chi tiết. Tham số/đầu ra tích hợp chưa xác nhận được quản lý riêng tại mục 11 SRS. |
| [Hợp đồng tích hợp API và sự kiện](<Notitek - HỢP ĐỒNG TÍCH HỢP API VÀ SỰ KIỆN.md>) | Các bên trao đổi dữ liệu và kết quả gì? | Tech Lead các module | Bản nháp ngữ nghĩa; cần OpenAPI/schema event cùng bộ ví dụ tương thích trước tích hợp. |
| [Đặc tả hành vi giao diện (UX)](<Notitek - ĐẶC TẢ HÀNH VI GIAO DIỆN (UX).md>) | Người dùng thấy gì và thao tác ở đâu? | PO, UX và frontend | Đặc tả hành vi; wireframe và ánh xạ màn hình cần hoàn thiện trước code UI. |
| [Thiết kế kỹ thuật sơ bộ](<Notitek - THIẾT KẾ KỸ THUẬT SƠ BỘ.md>) | Chọn giải pháp, tổ chức xử lý và lưu dữ liệu như thế nào? | Tech Lead và phát triển | Mô hình logic và danh sách quyết định; chưa chọn stack, nhà cung cấp hoặc schema vật lý. |
| [Đề xuất công nghệ sử dụng cho dự án](<Notitek - ĐỀ XUẤT CÔNG NGHỆ SỬ DỤNG CHO DỰ ÁN.md>) | Mỗi lớp bài toán dùng công nghệ nào và khi nào đổi theo quy mô? | Tech Lead | Bản 1.1: PO chốt Go, RabbitMQ cho job gửi, PostgreSQL, SSE + Redis Pub/Sub cho In-app; broker liên module chờ SuperPlatform; 7 rủi ro TR cần theo dõi. |
| [Backlog triển khai và kịch bản nghiệm thu](<Notitek - BACKLOG TRIỂN KHAI VÀ KỊCH BẢN NGHIỆM THU.md>) | Ai làm phần nào, theo thứ tự nào, kiểm tra đạt bằng gì? | PO, Tech Lead và QA | Backlog 0.2: truy vết SRS 1.0 và 28 kịch bản AT; task kỹ thuật được tách sau khi chốt thiết kế. |

## Tài liệu máy đọc cần có trước tích hợp

Các tài liệu Markdown trên xác định hành vi và trách nhiệm. Đội kỹ thuật cần bổ sung các artifact có thể kiểm tra tự động:

- OpenAPI cho nhận yêu cầu gửi, tra cứu, hủy và hộp thông báo, gồm xác thực, dữ liệu, lỗi và ví dụ. Chỉ tạo route sau khi chốt hợp đồng.
- Schema event đầu vào/đầu ra và quy tắc tương thích; dùng JSON Schema hoặc AsyncAPI theo phương thức tích hợp được chọn.
- Bộ ví dụ thành công, từ chối, trùng ID, kết quả chưa rõ, hủy và kết quả một phần; có consumer/provider contract tests.
- ERD, từ điển dữ liệu và migration theo giải pháp được chọn; phân loại dữ liệu nhạy cảm và thời gian lưu.
- Wireframe/prototype cho hộp tin và tra cứu vận hành, có các trạng thái tải, trống, lỗi và mất quyền.
- Quyết định kiến trúc về giải pháp điều phối, cơ chế phát/nhận bền vững, xử lý bí mật và ánh xạ kênh.
- Bộ kiểm thử tích hợp/UAT và hướng dẫn chuyển đổi, theo dõi, dừng gửi và khôi phục vận hành.

## Điều kiện để một story bắt đầu code

Story phải có UC, hành vi SRS, phạm vi trách nhiệm, tiêu chí nghiệm thu và bên thực hiện. Phần có API/event cần hợp đồng cùng ví dụ được các bên thống nhất. Phần UI cần hành vi màn hình và nguồn dữ liệu; phần lưu dữ liệu cần chính sách dữ liệu và thiết kế phù hợp. Các quyết định còn mở phải ghi rõ chỉ chặn phần nào.

## Điều kiện kết thúc một lát cắt

Nghiệp vụ từ module nguồn đi qua Notification tới đúng người/app/phạm vi bằng kênh thực; người dùng mở đúng hành động; kết quả có thể tra cứu; các ca trùng, lỗi, chưa rõ và hủy được kiểm chứng. Mock hỗ trợ phát triển và kiểm thử, không thay bằng chứng nghiệm thu gửi thật.

PO chịu trách nhiệm phạm vi, giá trị và nghiệm thu nghiệp vụ. Tech Lead chịu trách nhiệm lựa chọn và phê duyệt thiết kế kỹ thuật. Việc người yêu cầu chốt định hướng sản phẩm không tự phê duyệt mọi route, schema hoặc thông số trong bản nháp.
