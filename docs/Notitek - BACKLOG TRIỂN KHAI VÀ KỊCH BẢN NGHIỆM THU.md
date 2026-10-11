# Notitek - BACKLOG TRIỂN KHAI VÀ KỊCH BẢN NGHIỆM THU

**Phiên bản:** 0.2. **Ngày:** 08/10/2026. **Chủ trì:** PO cùng Tech Lead và QA.

Backlog tổ chức theo các lát cắt đã chốt trong [phạm vi](<Notitek - PHẠM VI VÀ KẾ HOẠCH PHÁT HÀNH ĐẦU TIÊN.md>). Story có bên chịu trách nhiệm, phụ thuộc và kết quả nghiệm thu; đội tách task code/kiểm thử sau khi chốt hợp đồng và thiết kế. Không gán lịch sprint hoặc effort trước khi đội xác nhận năng lực và phụ thuộc.

## 1 Các story và thứ tự thực hiện

| Story | Kết quả và công việc | Bên thực hiện | Phụ thuộc và điều kiện bắt đầu | Mốc |
|---|---|---|---|---|
| ST-01 | Cấu hình tối thiểu luồng/mẫu/kênh theo phạm vi, phiên bản và kiểm tra trước kích hoạt | Notification, vận hành | DEC-07/12; mẫu/cấu hình hợp lệ; phương thức quản trị được chọn | M0 |
| ST-02 | User gửi OTP qua Notification và nhận đúng kết quả | User, Notification, vận hành SMS | DEC-02/05/06/08/09/10; hợp đồng gửi/tra/hủy và xử lý bí mật | M1 |
| ST-03 | Lời mời nhân viên/thành viên gửi email sau commit và mở đúng User | User, Notification, frontend liên quan | Hợp đồng User, mẫu/link, kết nối Email, chính sách lưu/hạn gửi | M1 |
| ST-04 | Order phát sự việc giao thất bại và chỉ dẫn hủy qua đường tích hợp bền vững | Order, Notification | DEC-02/03; event/version/reference và người nhận/căn cứ; thiết kế đường phát | M2 |
| ST-05 | Shop xem thông báo đơn đúng phạm vi và mở Order | Notification, Shop FE, User/Authorization | ST-04; hợp đồng In-app/quyền/ngữ cảnh; UX-01/02 | M2 |
| ST-06 | Push tới đúng client/app/người | Notification, chủ app/thiết bị, frontend | Client Push, chủ đích gửi, hợp đồng thiết bị và cấu hình kênh được chốt | M2 |
| ST-07 | Người nhận hàng nhận ZNS, dùng SMS dự phòng đúng chính sách | Notification, Order, vận hành | ST-04; mẫu/định danh gửi, DEC-04/08/10; OA và mẫu ZNS đăng ký mới được Zalo duyệt (CF-08); điều kiện dự phòng | M2 |
| ST-08 | Chuông Shop nối vào hộp tin Notitek, hiển thị cảnh báo User và trạng thái đọc phù hợp | Shop FE, User, Notification | Hợp đồng event và khóa đối chiếu ID; UX-01/03; chính sách ngữ cảnh/đọc | M2 |
| ST-09 | Vận hành tra cứu được tiến trình, lỗi/chưa rõ và cấu hình | Notification, frontend vận hành, User/Authorization | Hợp đồng tra cứu, quyền xem, chính sách dữ liệu/audit; UX-04 | M1 đến M3 |
| ST-10 | Đã loại bỏ: không có hệ thống cũ cần chuyển đổi (09/10/2026). Giữ mã để không đổi số story | — | — | — |
| ST-11 | Hoàn thiện và kiểm chứng hợp đồng/thiết kế cho từng lát cắt | Tech Lead các bên, BA, QA | UC/SRS; OpenAPI/schema, ví dụ, ADR và thiết kế dữ liệu phần liên quan | M0 và từng lát cắt |
| ST-12 | Nghiệm thu đầy đủ năm kênh và các ngoại lệ của đợt đầu | QA, PO nguồn/Notification, vận hành | ST-02 đến ST-09; mục tiêu NFR, môi trường và cấu hình gửi thật được xác nhận | M3 |
| ST-13 | Đo giá trị sản phẩm và công sức tích hợp luồng mới | PO, User, Order, Notification, vận hành | Mốc đo, dữ liệu nguồn, cách tính và quyền sử dụng dữ liệu được chốt | M1 đến M3 |

ST-11 được tách theo lát cắt, không là công việc làm hết kiến trúc toàn nền tảng trước khi bắt đầu. ST-09 cung cấp tra cứu cần thiết ngay khi có gửi thật; không chờ mọi màn hình quản trị nâng cao.

## 2 Ma trận truy vết cho đợt đầu

Bảng đồng bộ với SRS 1.0 sau PO review; liên kết phần được chọn, không tuyên bố đã nghiệm thu toàn bộ BRD dài hạn. Mã BR/SRS rút gọn với dấu slash giữ nguyên tiền tố của mục đầu.

| Story | UC | BR chính | SRS | UX | AT |
|---|---|---|---|---|---|
| ST-01 | UC-NTF-15, UC-NTF-04/07/08 | BR-CAT-04, BR-TPL-01/03, BR-CFG-01/02/03/05 | SRS-F03/04/12/15/20/21 | UX-05 | AT-17/18/24/27 |
| ST-02 | UC-USR-01, UC-NTF-02/03/09/10/11 | BR-INT-02/05, BR-TRG-02/05, BR-TIME-01/02, BR-SEC-01/05, BR-OPS-04 | SRS-F01/02/03/04/05/06/07/09/13/14/15/18/20/21 | Theo UI User hiện có | AT-01/02/03/04/19/21/22/23/28 |
| ST-03 | UC-USR-03/04 | BR-INT-01/03/05, BR-REC-01/03, BR-TPL-01/02/03/05, BR-CHN-08 | SRS-F01/02/04/05/06/07/09/13/14/15/17/20 | Theo UI lời mời hiện có | AT-02/03/05/06/07/08/21/24/28 |
| ST-04 | UC-ORD-05, UC-NTF-01/03/10 | BR-INT-01/05, BR-TRG-01/02/03/05, BR-TIME-02 | SRS-F01/02/03/05/06/09 | Theo UI Order hiện có | AT-02/09/10/21/22/23 |
| ST-05 | UC-ORD-05, UC-NTF-13 | BR-CHN-05, BR-REC-03/04, BR-SEC-03/06, BR-CFG-05 | SRS-F05/08/09/13 | UX-01/02 | AT-08/11/15/21/25/26 |
| ST-06 | UC-ORD-05, UC-NTF-08/09 | BR-CHN-06/09, BR-CFG-05 | SRS-F05/07/09/13/15/16/21 | UX-01/02 | AT-11/12/24/25/26/27 |
| ST-07 | UC-ORD-05, UC-NTF-08/09 | BR-CHN-03/04/07/09, BR-REC-01, BR-COST-01/03 | SRS-F05/07/09/13/14/15/18/19/20/21 | Nội dung kênh ngoài | AT-09/11/19/24/25 |
| ST-08 | UC-USR-09, UC-NTF-13 | BR-INT-01/06, BR-REC-03/06, BR-CHN-05, BR-SEC-03/06 | SRS-F02/05/08/09/10/13/22 | UX-01/03 | AT-03/08/13/14/15/21/26/28 |
| ST-09 | UC-NTF-12 | BR-INT-04, BR-OPS-01/04, BR-SEC-01/02/04/05/06 | SRS-F09/11/12/13/14/15/21 | UX-04 | AT-03/16/17/18/19/25/27/28 |
| ST-10 | Đã loại bỏ | — | — | — | — |
| ST-11 | Các UC của lát cắt tương ứng | BR-INT-04/06, BR-CFG-05 | SRS-F01 đến SRS-F21; SRS-N01 đến SRS-N07 khi liên quan | UX-01 đến UX-05 khi liên quan | AT-01 đến AT-28 khi liên quan |
| ST-12 | Các UC được chọn và UC-NTF-09/13 | BR-CHN-05/06/07/08/09, BR-OPS-04 | SRS-F01 đến SRS-F21, SRS-N01 đến SRS-N07 | UX-01 đến UX-04 | AT-01 đến AT-28; AC-17 BRD |
| ST-13 | Các UC nghiệp vụ được chọn | BR-OPS-01, BR-INT-01 | SRS-F09, SRS-N07 | Không yêu cầu dashboard riêng | Chỉ số theo phạm vi; nguồn xác nhận kết quả nghiệp vụ |

## 3 Các kịch bản nghiệm thu

AT là kịch bản của đợt đầu, bổ sung và cụ thể hóa AC trong BRD. QA tách thành test case theo schema, môi trường và thông số thật. Người phụ trách ghi bằng chứng, kết quả và lỗi; không đánh dấu đạt từ tài liệu hoặc mock.

| Mã | Tình huống | Kết quả cần đạt | AC BRD liên quan |
|---|---|---|---|
| AT-01 | User yêu cầu gửi OTP đủ dữ liệu qua kênh thực | Kết quả đáp ứng mức User đã chốt; không dùng accepted/log làm bằng chứng gửi thật | AC-01/06/17 |
| AT-02 | Nguồn gửi lại cùng reference; lỗi giữa các bước | Giữ tham chiếu và xử lý phần còn hợp lệ; không gửi lại phần đã đạt tiêu chí | AC-02 |
| AT-03 | Tra log/lịch sử/kết quả gửi mã hoặc link | Không lộ OTP/token; nội dung nhạy cảm che/lưu/xóa theo chính sách | AC-14 |
| AT-04 | OTP hết hạn gửi, timeout, phản hồi muộn hoặc User cấp lại | Notification tuân thủ chỉ dẫn gửi/thử lại/hủy; User quyết định hiệu lực mã; không tự báo thành công khi chưa có bằng chứng | AC-06/07/08 |
| AT-05 | Lời mời đã commit được gửi và mở link | Email thực và link đúng User/lời mời; kết quả gửi tách kết quả chấp nhận | AC-01/17 |
| AT-06 | Lời mời đã dùng/thu hồi/cấp lại hoặc gửi lỗi | Gửi lỗi không đảo commit; bản chờ dừng theo chỉ dẫn/hạn; không tự tạo membership | AC-07 |
| AT-07 | Email bị trả lại, một đích lỗi hoặc các email cá nhân hóa gửi cho nhiều người | Ghi kết quả theo bằng chứng; không báo giao/đọc chỉ từ tiếp nhận; không lộ địa chỉ/nội dung riêng của người nhận khác | AC-08/09 |
| AT-08 | Một người có nhiều Shop, đổi Shop/tài khoản hoặc mất quyền | Công việc tách đúng scope; cảnh báo tài khoản không lộ sang tài khoản khác; link do nguồn kiểm tra | AC-04/13/19 |
| AT-09 | Order xác nhận giao thất bại; ZNS timeout hoặc lỗi đủ điều kiện dự phòng | Dùng sự việc và dữ liệu nguồn; timeout giữ chưa rõ; SMS chỉ theo điều kiện đã chốt | AC-05/08/09 |
| AT-10 | Event trễ/sai thứ tự; Order hủy yêu cầu chờ | Theo sự việc/reference/chỉ dẫn nguồn; dừng phần kiểm soát được; không hứa thu hồi tin ngoài | AC-07 |
| AT-11 | Shop nhận In-app/Push, người nhận không có tài khoản nhận ZNS/SMS | Đúng người/app/phạm vi; không tự tạo tài khoản cho liên hệ; kết quả từng kênh riêng | AC-05/09/17/19 |
| AT-12 | Token Push lỗi, đăng xuất/đổi chủ thiết bị, callback trễ/trùng | Không gửi sang chủ/app khác; cập nhật theo hợp đồng; callback không tạo lần gửi mới | AC-16/17/19 |
| AT-13 | Người dùng đọc tin rồi xác nhận/báo thiết bị | Đọc không thực hiện bảo mật; hành động tại User, đúng điều kiện nguồn | AC-13/19 |
| AT-14 | Cảnh báo thiết bị được phát lại hoặc đến từ nhiều đường | Một cảnh báo được đối chiếu đúng ID nguồn; không tạo hai bản trong hộp tin, giữ hành động tại User | AC-02/19 |
| AT-15 | Đọc In-app khi kênh khác còn chờ; tải lỗi hoặc đổi ngữ cảnh | Đọc độc lập với gửi; trạng thái UX phù hợp; không dùng cache scope cũ để lộ dữ liệu | AC-09/13/19 |
| AT-16 | Vận hành tra reference và người không đủ quyền thử xem | Người có quyền tra đủ tiến trình; người trái quyền bị từ chối; truy cập/thao tác có audit theo chính sách | AC-01/14/18/19 |
| AT-17 | Mock/log và gửi thật cùng tồn tại khi phát triển | Có nhãn/nguồn bằng chứng riêng; mock không được tính là hoàn tất kênh thực | AC-06/17 |
| AT-18 | Mẫu/cấu hình thay đổi, phạm vi xung đột, reference cũ có dữ liệu khác | Tra được phiên bản; không kích hoạt cấu hình xung đột; xử lý payload xung đột theo hợp đồng | AC-02/11/18 |
| AT-19 | Điều kiện ngân sách không cho phép gửi | Bên sở hữu quyết định; Notification tuân thủ hợp đồng, có lý do và không báo thành công khi bị chặn | AC-12 |
| AT-20 | Đã loại bỏ: không có luồng cũ cần chuyển đổi. Giữ mã để không đổi số ca | — | — |
| AT-21 | Nguồn trái scope, payload giả người/Shop/app; người biết ID tin người khác | Không kích hoạt/gửi hoặc trả nội dung ngoài quyền; danh sách và số đếm không lộ dữ liệu; backend yêu cầu căn cứ quyền đúng hành động/tài nguyên | AC-14/19 |
| AT-22 | Lỗi trước/sau ghi nhận bền vững; mất phản hồi/phát lại; version không hỗ trợ và thêm trường tương thích | Trước ghi nhận không trả accepted; sau ghi nhận tra lại được, không mất/nhân công việc; version sai có lý do, thay đổi tương thích không làm consumer cũ lỗi; event Order có đường bàn giao thực | AC-01/02 |
| AT-23 | Hai yêu cầu cùng khóa đồng thời, crash sau provider nhận; cửa sổ đồng bộ hết; hủy đến trước event | Không nhân tập bản tin, không retry mù quáng; không bắt đầu lần gửi OTP mới ngoài cửa sổ được phép; giữ dấu hủy theo hợp đồng hoặc báo rõ chưa hỗ trợ | AC-02/07/08 |
| AT-24 | Kết nối sai scope/brand/môi trường; mẫu thiếu biến/biến có ký tự chèn; SMS dài; lựa chọn Marketing/luồng bắt buộc và liên hệ bị chặn an toàn | Không dùng cấu hình gần giống; dữ liệu không phá mẫu/cắt mất ý nghĩa; điều kiện kênh, mục đích và đích an toàn đúng chính sách; có lý do/version; thử không dùng tập khách production | AC-11/17/18 |
| AT-25 | Nhiều người/kênh bắt buộc/tùy chọn/thay thế, nhiều đích Push một người; callback trễ/trùng/trái nhau; phát kết quả lỗi | Mục tiêu tính đúng tiêu chí; phần thiếu và chưa rõ hiện riêng; bằng chứng không bị ghi đè sai; phát lại kết quả không gửi lại tin người dùng | AC-08/09 |
| AT-26 | Đọc cùng tin trên hai thiết bị, pagination; đổi quyền/ngữ cảnh; Push thu hồi rồi nhận bản đăng ký cũ | Trạng thái đọc backend nhất quán ở lần tải sau; danh sách/số đếm cùng scope; không dùng token cũ cho chủ/app khác; đọc không thực hiện nghiệp vụ | AC-13/16/19 |
| AT-27 | Provider rate-limit, burst từ nguồn khác cùng OTP, dừng luồng/kênh, đổi credential; backup/restore | Giới hạn và ưu tiên theo cấu hình; accepted không mất; phần quá hạn có kết quả; phạm vi dừng/audit rõ; credential không lộ; NFR dùng giá trị đã duyệt | AC-12/17/18 |
| AT-28 | Bí mật giả đi qua exception, trace, dead-letter, chống trùng, API hỗ trợ; đến hạn xóa và restore | Không có nguyên văn/dữ liệu dễ khôi phục bí mật ở nơi bị cấm; lưu/xóa từng loại đúng chính sách và quyền; tra cứu vẫn có metadata được phép | AC-14 |

AT-10/12/18 bổ sung ca sai thứ tự và dữ liệu xung đột chưa có AC riêng trong BRD. AT-21 đến AT-28 bổ sung kiểm chứng theo SRS 1.0; QA phải tách từng tình huống thành test case, không coi một kết quả chung là đạt mọi yêu cầu con. Test về SLA/tải/lưu dữ liệu phải dùng giá trị SRS-N đã được duyệt và ghi môi trường đo.

## 4 Điều kiện bắt đầu code

- PO xác nhận story thuộc lát cắt và actor/kết quả cần đạt.
- BA liên kết UC, BR, SRS và AT; quyết định còn mở ghi người chịu trách nhiệm và phần bị chặn.
- Tech Lead chốt hợp đồng với consumer/provider, mô hình trạng thái/dữ liệu và ADR của phần story sử dụng.
- Frontend có hành vi UX cùng nguồn dữ liệu/ví dụ; không dùng dữ liệu fixture làm hợp đồng thật.
- QA xác định ca đạt/lỗi và dữ liệu kiểm thử; các thông số cần đo có đơn vị/giá trị khi story phụ thuộc chúng.

Các việc có thể bắt đầu ngay là hoàn thiện hợp đồng/thiết kế, dữ liệu ví dụ, prototype và kế hoạch kiểm thử. Code gửi thật, quyền, dữ liệu lưu và tích hợp phải theo các quyết định phần tương ứng đã chốt.

## 5 Điều kiện hoàn thành và phát hành

Story có code được review, kiểm thử phù hợp, contract tests khi tích hợp, bằng chứng AT liên quan và tài liệu cập nhật. Lát cắt có kiểm chứng với nguồn thật và kênh thực; không chỉ kiểm tra adapter đơn lẻ. Bản phát hành đợt đầu có đủ năm kênh, các lát cắt được chọn, quyền/phạm vi, tra cứu, NFR được duyệt và quy trình vận hành.

QA và vận hành bổ sung hướng dẫn gửi thử, theo dõi tồn đọng/chưa rõ, dừng phần chưa gửi, xử lý sự cố và khôi phục. PO duyệt kết quả nghiệp vụ; Tech Lead/vận hành xác nhận điều kiện triển khai. Thời điểm bật gửi thật/chuyển luồng cũ phải được quyết định trong kế hoạch phát hành cụ thể.
