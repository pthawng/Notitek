# Notitek - HƯỚNG DẪN SỬ DỤNG BỘ TÀI LIỆU

Bộ gồm 12 tài liệu chính và bộ hợp đồng máy đọc trong thư mục `api/` (OpenAPI, ví dụ, README API và 4 bản nháp: hộp tin/Push, webhook provider/tra cứu vận hành, yêu cầu gửi đội User, yêu cầu gửi đội Order), nối phạm vi sản phẩm với công việc phát triển và nghiệm thu. Cập nhật trạng thái ngày 11/10/2026. Phạm vi sản phẩm được người yêu cầu xác nhận ngày 08/10/2026; các chi tiết kỹ thuật có trạng thái riêng để PO, BA, Tech Lead, phát triển và QA cùng sử dụng.

## Thứ tự sử dụng

BRD → phạm vi bản phát hành → use case ưu tiên → SRS, trải nghiệm và hợp đồng tích hợp → thiết kế kỹ thuật → backlog và nghiệm thu → code.

SRS, trải nghiệm và hợp đồng được hoàn thiện cùng nhau theo từng lát cắt. Đội có thể triển khai phần đã đủ quyết định; không phải đợi mọi luồng dài hạn được đặc tả.

| Tài liệu | Trả lời câu hỏi | Chủ trì | Trạng thái hiện tại |
|---|---|---|---|
| [Hướng dẫn sử dụng bộ tài liệu](<Notitek - HƯỚNG DẪN SỬ DỤNG BỘ TÀI LIỆU.md>) | Đọc tài liệu theo thứ tự nào, ai hoàn thiện và khi nào đủ điều kiện code? | PO và BA | Mục lục của bộ 12 tài liệu chính và thư mục `api/`. |
| [Đặc tả yêu cầu nghiệp vụ (BRD)](<../Notitek - ĐẶC TẢ YÊU CẦU NGHIỆP VỤ (BRD).md>) | Mục tiêu, ranh giới và yêu cầu nghiệp vụ là gì? | PO và BA | BRD 0.7; phạm vi đợt đầu đã xác nhận (gồm 4 thông báo tài khoản User và 10 mốc Order), chính sách chi tiết còn mở. |
| [Danh mục use case](<../Notitek - DANH MỤC USE CASE.md>) | Có những use case nào và ưu tiên ra sao? | PO và BA | 43 UC (thêm UC-ORD-12 đến 16 ngày 11/10/2026); các lát cắt đợt đầu đã xác nhận. |
| [Phạm vi và kế hoạch phát hành đầu tiên](<Notitek - PHẠM VI VÀ KẾ HOẠCH PHÁT HÀNH ĐẦU TIÊN.md>) | Đợt đầu phục vụ ai, làm những luồng nào, hoàn thành đến mức nào? | PO | Bản 1.2: định hướng và các lát cắt R1 đã xác nhận; mốc M0–M3 chưa có ngày, có ràng buộc theo hạn CF. |
| [Quyết định PO cho các DEC](<Notitek - QUYẾT ĐỊNH PO CHO CÁC DEC.md>) | 12 câu hỏi mở được chốt thế nào, tham số dùng giá trị nào, còn chờ ai ký? | PO | PO đã chốt DEC-01 đến DEC-12 và SRS-P01 đến SRS-P12 ngày 08/10/2026; bản 1.6 (11/10/2026) gồm A1–A8, B1 (R1 REST, 17 profile, audience vận hành, API cấp lại link), hướng xử lý B2–B11, C1–C7, D và các sửa sau review; 12 điểm chờ bên sở hữu xác nhận (CF-01 đến CF-12). |
| [Đặc tả use case ưu tiên](<Notitek - ĐẶC TẢ USE CASE ƯU TIÊN.md>) | Actor thực hiện gì, luồng chính và ngoại lệ ra sao? | BA cùng PO các module nguồn | Bản 0.4: đã phân rã các UC ưu tiên của đợt đầu; phạm vi bổ sung (thông báo tài khoản, 9 mốc Order) dùng chung mẫu phân rã. |
| [Đặc tả yêu cầu phần mềm (SRS)](<Notitek - ĐẶC TẢ YÊU CẦU PHẦN MỀM (SRS).md>) | Hệ thống phải có hành vi và chất lượng nào để đáp ứng UC? | BA và Tech Lead | SRS 1.4: baseline nghiệp vụ đã PO review và chốt, phạm vi R1 theo DEC 1.6; 21 nhóm chức năng (SRS-F22 đã loại), NFR và kiểm chứng chi tiết. Tham số/đầu ra tích hợp chưa xác nhận được quản lý riêng tại mục 11 SRS. |
| [Hợp đồng tích hợp API và sự kiện](<Notitek - HỢP ĐỒNG TÍCH HỢP API VÀ SỰ KIỆN.md>) | Các bên trao đổi dữ liệu và kết quả gì? | Tech Lead các module | Có [thiết kế API chung](api/README.md), [OpenAPI](api/openapi.json) và ví dụ cho gửi/tra cứu/hủy qua `/v1/notifications`, 17 profile R1 (7 User, 10 Order), snake_case, `subject_ref`; R1 không dùng broker. Bản nháp [hộp tin/Push](api/inbox-and-push.md), [webhook provider/tra cứu vận hành](api/provider-webhooks-and-ops.md); yêu cầu gửi [đội User](api/user-integration-requirements.md) và [đội Order](api/order-integration-requirements.md). Schema event chỉ cần khi có broker. |
| [Đặc tả hành vi giao diện (UX)](<Notitek - ĐẶC TẢ HÀNH VI GIAO DIỆN (UX).md>) | Người dùng thấy gì và thao tác ở đâu? | PO, UX và frontend | Bản 0.4: đặc tả hành vi, UX-01 gồm cả vận hành trên giao diện nội bộ; wireframe và ánh xạ màn hình cần hoàn thiện trước code UI. |
| [Thiết kế kỹ thuật sơ bộ](<Notitek - THIẾT KẾ KỸ THUẬT SƠ BỘ.md>) | Chọn giải pháp, tổ chức xử lý và lưu dữ liệu như thế nào? | Tech Lead và phát triển | Bản 0.5: mô hình logic và danh sách quyết định; stack đã chốt tại ADR-01; còn chờ nhà cung cấp (CF-12) và schema vật lý. |
| [Đề xuất công nghệ sử dụng cho dự án](<Notitek - ĐỀ XUẤT CÔNG NGHỆ SỬ DỤNG CHO DỰ ÁN.md>) | Mỗi lớp bài toán dùng công nghệ nào và khi nào đổi theo quy mô? | Tech Lead | Bản 1.2: PO chốt Go, RabbitMQ cho job gửi, PostgreSQL, SSE + Redis Pub/Sub cho In-app; R1 không dùng broker liên module (REST từ dispatcher outbox nguồn), broker để phân kỳ sau; 7 rủi ro TR cần theo dõi. |
| [Backlog triển khai và kịch bản nghiệm thu](<Notitek - BACKLOG TRIỂN KHAI VÀ KỊCH BẢN NGHIỆM THU.md>) | Ai làm phần nào, theo thứ tự nào, kiểm tra đạt bằng gì? | PO, Tech Lead và QA | Backlog 0.4: ST-01 đến ST-15, truy vết SRS 1.4 và 41 kịch bản AT; task kỹ thuật được tách sau khi chốt thiết kế. |

## Tài liệu máy đọc cần có trước tích hợp

Các tài liệu Markdown trên xác định hành vi và trách nhiệm. Đội kỹ thuật cần bổ sung các artifact có thể kiểm tra tự động:

- OpenAPI cho gửi, tra cứu và hủy đã có tại `api/openapi.json` (R1, chờ Tech Lead User/Order ký theo CF-03). Hộp tin, Push, webhook provider và tra cứu vận hành mới có bản nháp Markdown; chỉ tạo các route đó sau khi hợp đồng tương ứng được chốt.
- Schema event đầu vào/đầu ra và quy tắc tương thích (JSON Schema hoặc AsyncAPI) khi bắt đầu phân kỳ broker; R1 tích hợp hoàn toàn qua REST.
- Bộ ví dụ thành công, từ chối, trùng ID, kết quả chưa rõ, hủy và kết quả một phần; có consumer/provider contract tests.
- ERD, từ điển dữ liệu và migration theo giải pháp được chọn; phân loại dữ liệu nhạy cảm và thời gian lưu.
- Wireframe/prototype cho hộp tin và tra cứu vận hành, có các trạng thái tải, trống, lỗi và mất quyền.
- Quyết định kiến trúc về giải pháp điều phối, cơ chế phát/nhận bền vững, xử lý bí mật và ánh xạ kênh.
- Bộ kiểm thử tích hợp/UAT và hướng dẫn theo dõi, dừng gửi và khôi phục vận hành.

## Điều kiện để một story bắt đầu code

Story phải có UC, hành vi SRS, phạm vi trách nhiệm, tiêu chí nghiệm thu và bên thực hiện. Phần có API/event cần hợp đồng cùng ví dụ được các bên thống nhất. Phần UI cần hành vi màn hình và nguồn dữ liệu; phần lưu dữ liệu cần chính sách dữ liệu và thiết kế phù hợp. Các quyết định còn mở phải ghi rõ chỉ chặn phần nào.

## Điều kiện kết thúc một lát cắt

Nghiệp vụ từ module nguồn đi qua Notification tới đúng người/app/phạm vi bằng kênh thực; người dùng mở đúng hành động; kết quả có thể tra cứu; các ca trùng, lỗi, chưa rõ và hủy được kiểm chứng. Mock hỗ trợ phát triển và kiểm thử, không thay bằng chứng nghiệm thu gửi thật.

PO chịu trách nhiệm phạm vi, giá trị và nghiệm thu nghiệp vụ. Tech Lead chịu trách nhiệm lựa chọn và phê duyệt thiết kế kỹ thuật. Việc người yêu cầu chốt định hướng sản phẩm không tự phê duyệt mọi route, schema hoặc thông số trong bản nháp.
