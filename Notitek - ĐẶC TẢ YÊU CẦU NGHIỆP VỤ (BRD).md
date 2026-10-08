# Notitek - ĐẶC TẢ YÊU CẦU NGHIỆP VỤ (BRD)

**Phiên bản:** 0.6 - chốt định hướng sản phẩm, các lát cắt đợt đầu và DEC-01 đến DEC-12 ngày 08/10/2026.  
**Trạng thái:** Phạm vi sản phẩm đợt đầu đã xác nhận: đủ 5 kênh, OTP, lời mời nhân viên/thành viên, giao thất bại và cảnh báo bảo mật/thiết bị; tiếp nối trải nghiệm Shop và nội bộ. Các DEC và thông số khởi điểm đã được PO chốt tại docs/Notitek - QUYẾT ĐỊNH PO CHO CÁC DEC.md; schema hợp đồng và một số xác nhận của bên sở hữu còn cần hoàn thiện.  
**Đối tượng đọc:** Chủ nghiệp vụ, người dùng, quản trị viên, nhóm User/Order/các module nguồn, BA, phát triển và kiểm thử.

Tài liệu mô tả người dùng cần gì và các module phải phối hợp thế nào. API, cấu trúc event, cơ chế xử lý và công nghệ được đặc tả trong SRS sau khi các quyết định nghiệp vụ liên quan được chốt.

Tài liệu liên quan: [Danh mục use case](<Notitek - DANH MỤC USE CASE.md>).

Từ nghiệp vụ xuống code: [Bộ tài liệu triển khai](<docs/Notitek - HƯỚNG DẪN SỬ DỤNG BỘ TÀI LIỆU.md>), [phạm vi bản phát hành đã chốt](<docs/Notitek - PHẠM VI VÀ KẾ HOẠCH PHÁT HÀNH ĐẦU TIÊN.md>), [use case ưu tiên](<docs/Notitek - ĐẶC TẢ USE CASE ƯU TIÊN.md>) và [backlog/nghiệm thu](<docs/Notitek - BACKLOG TRIỂN KHAI VÀ KỊCH BẢN NGHIỆM THU.md>).

## 1. Cách đọc và lịch sử thay đổi

### 1.1. Quy ước

- **BR-<NHÓM>-<SỐ>** là mã yêu cầu để liên kết BRD, use case, SRS và nghiệm thu. Mã đã gán cho một yêu cầu không được dùng lại cho nội dung khác.
- **P0:** Nền tảng và luồng cần ưu tiên để tích hợp User/Order an toàn. **P1:** Mở rộng trải nghiệm và vận hành. **P2:** Năng lực mở rộng theo module và chiến dịch. Đây là ưu tiên triển khai đề xuất, không phải mức khẩn của tin gửi.
- **Đối chiếu hiện trạng:** Có bằng chứng trong mã nguồn đã đọc; không đồng nghĩa đã chạy tích hợp hoặc gửi thật ở môi trường sản xuất.
- **Đề xuất:** Quy tắc nghiệp vụ được đưa ra trong tài liệu, cần chủ nghiệp vụ duyệt.
- **Cần chốt:** Chưa đủ căn cứ để lựa chọn chính sách; được ghi tại mục 11. Các yêu cầu BR là đề xuất trừ phần được ghi rõ đã xác nhận; phạm vi kênh và các lát cắt sản phẩm đợt đầu đã được xác nhận. Việc chốt định hướng không tự phê duyệt mọi chính sách BR hoặc chi tiết kỹ thuật.

Một use case có thể cần nhiều nhóm BR. Đợt đầu đã xác nhận đủ 5 kênh tại mục 5.2; các yêu cầu bắt buộc để sử dụng từng kênh nằm trong phạm vi nghiệm thu đợt đầu. Màn hình quản trị nâng cao và từng luồng nghiệp vụ vẫn cần xác định phạm vi riêng.

### 1.2. Lịch sử thay đổi

| Phiên bản | Ngày | Nội dung |
|---|---|---|
| Bản đầu vào | Chưa có thông tin trong bản gốc | Định hướng Notification dùng chung, các nhóm yêu cầu và phạm vi SuperPlatform; chưa có yêu cầu đánh mã cụ thể. |
| 0.2 | 06/10/2026 | Điền mục tiêu, actor, quy tắc, kết quả, nghiệm thu và quyết định còn mở; đối chiếu User/Order; tách phạm vi mục tiêu với lộ trình đề xuất. Lưu bản đầu vào tại thời điểm rà soát. |
| 0.2 - cập nhật phạm vi | 06/10/2026 | Người yêu cầu xác nhận đợt đầu hỗ trợ In-app, Push, Email, SMS và Zalo; cập nhật yêu cầu kênh và trạng thái DEC-01. |
| 0.3 | 06/10/2026 | Đối chiếu Application, Membership, Entitlement và Access Context của User; tham khảo mô hình Novu; bổ sung cách cấu hình luồng, dữ liệu người nhận và app đích, dẫn chiếu phân quyền dùng chung. |
| 0.4 | 06/10/2026 | Chuẩn hóa tên gọi các thành phần theo nghiệp vụ SPF; giữ nguyên mô hình, chức năng, phạm vi và mã BR/UC/AC. |
| 0.5 | 08/10/2026 | Người yêu cầu chốt định hướng PO và các lát cắt OTP, lời mời, giao thất bại, cảnh báo thiết bị; giữ 5 kênh, tách ngữ cảnh tài khoản/công việc, làm rõ ranh giới nguồn và Notification. Bổ sung UC-USR-09 và bộ tài liệu xuống triển khai; lưu bản 0.4 tại thời điểm chốt phạm vi. |
| 0.6 | 08/10/2026 | PO chốt DEC-01 đến DEC-12 và giá trị khởi điểm SRS-P01 đến SRS-P12 tại tài liệu Quyết định PO cho các DEC; mục 11 dẫn chiếu kết quả, giữ câu hỏi gốc. |

## 2. Bối cảnh, mục tiêu và căn cứ

### 2.1. Notification giải quyết việc gì?

Khi một sự việc xảy ra, người liên quan cần biết **chuyện gì đã xảy ra, có cần làm gì tiếp và xem chi tiết ở đâu**. Notification là đầu mối nhận sự việc hoặc yêu cầu gửi từ các module, chuyển thành thông báo phù hợp, gửi đến đúng người và trả kết quả xử lý.

SuperPlatform phục vụ nhiều ứng dụng/thương hiệu như SuperShip, SuperAI, ứng dụng vận hành và cổng đối tác; nhiều nhà vận chuyển đứng ngang hàng về nghiệp vụ; nhiều tổ chức, Shop, bưu cục và pháp nhân độc lập. Vì vậy, cùng một người hoặc số điện thoại xuất hiện ở nhiều nơi không có nghĩa được nhận mọi thông tin của các nơi đó.

BRD đầu vào nêu SuperShip đã sử dụng ZNS cho các mốc giao hàng, chuyển hoàn và đánh giá dịch vụ. Đây là thông tin nghiệp vụ kế thừa cần chủ vận hành xác nhận về mẫu, người nhận, tần suất và chi phí; chưa được coi là bằng chứng tích hợp Notification mới.

### 2.2. Mục tiêu nghiệp vụ

| Mã | Mục tiêu | Kết quả người dùng hoặc module nguồn cần thấy |
|---|---|---|
| OBJ-01 | Đúng người, đúng phạm vi | Shop chỉ nhận thông tin của Shop được phép; người nhận hàng chỉ nhận thông tin cần thiết của đơn liên quan. |
| OBJ-02 | Thông tin có ích và còn hiệu lực | Tin giúp hiểu sự việc/hành động tiếp theo; mã hết hạn hoặc lời nhắc không còn cần thiết không tiếp tục được gửi. |
| OBJ-03 | Tích hợp thống nhất | Module nguồn biết yêu cầu đã được nhận hay bị từ chối, và có thể theo dõi kết quả theo đúng ý nghĩa từng trạng thái. |
| OBJ-04 | Trải nghiệm và thương hiệu nhất quán | Nội dung rõ, đúng app/brand, không bị lặp vô ích khi nguồn gửi lại event. |
| OBJ-05 | Vận hành có trách nhiệm | Tra được vì sao gửi/không gửi, ai nhận, phiên bản cấu hình, lỗi và đơn vị chịu phí. |
| OBJ-06 | Mở rộng dùng chung | Bổ sung module, nhà vận chuyển, tổ chức hoặc kênh bằng hợp đồng và chính sách được duyệt; không chuyển nghiệp vụ gốc sang Notification. |

### 2.3. Hiện trạng đã đối chiếu và giới hạn hiểu biết

| Phạm vi | Đã thấy trong mã nguồn | Hệ quả với BRD |
|---|---|---|
| User | NotificationPort cho OTP, kích hoạt, lời mời, đổi định danh, đổi mật khẩu và nghỉ việc; hiện có adapter giả lập/ghi log. | Phải phân biệt gửi thật với giả lập. OTP cần phản hồi gửi; các cảnh báo sau sự việc không được đảo ngược nghiệp vụ đã hoàn tất khi gửi lỗi. |
| User / quyền | Quản lý định danh, tư cách tham gia tổ chức/Shop, ứng dụng và phân quyền theo phạm vi. | Không mặc định người ở cấp trên hoặc chủ một Shop được quản trị mọi Notification. |
| User / ứng dụng và ngữ cảnh truy cập | Có Application/Client, Entitlement của tổ chức, Access Context gắn app và membership; Resolve có đánh giá hành động/phạm vi. | Notification tham chiếu danh mục và quyết định quyền từ User/Authorization; không tạo lại app, membership, role hoặc quyền truy cập riêng. Contract/mã nguồn chưa tự chứng minh mọi consumer đã tích hợp. |
| Order | Có outbox cho order.created.v1, order.cancelled.v1, order.updated.v1, order.return.requested.v1, order.return.confirmed.v1. OutboxWriter ghi dispatcher chưa được xây dựng. | Có bản ghi outbox chưa có nghĩa Notification đã nhận event; cần tích hợp phát/nhận và hợp đồng cụ thể. |
| Order / vận chuyển | Có trạng thái chuẩn và nghiệp vụ nhiều chặng/vận đơn. Tạo đơn được ghi nhận trước bước tạo vận đơn NVC. | Không dùng “đơn đã tạo” để báo “vận đơn đã tạo” hoặc “đang giao”. Trạng thái chuẩn do Order/module sở hữu xác nhận. |
| SuperPlatform mở rộng | BRD đầu vào mô tả tài chính, hỗ trợ, giá/sản phẩm, chiến dịch, nhiều app/brand/NVC. | Giữ làm phạm vi mục tiêu. Chưa có đủ bằng chứng để khẳng định mọi module, event, người nhận hoặc kênh đã sẵn sàng. |

Không coi dữ liệu phòng ban, đồng ý Marketing, thiết bị Push, hợp đồng nhà cung cấp kênh gửi hay mẫu kênh đã có sẵn nếu chưa được xác nhận. Các tham chiếu mã nguồn ở mục 13 là căn cứ tại thời điểm rà soát, không phải cam kết về hệ thống đã triển khai.

## 3. Thuật ngữ và mô hình dễ hiểu

Trong SPF, tên nghiệp vụ được dùng xuyên suốt tài liệu. Tên Novu chỉ giữ tại bảng đối chiếu và tài liệu tham khảo. Việc đổi tên không làm thay đổi mô hình hoặc trách nhiệm của các module; tên API, event và các định danh kỹ thuật hiện có giữ nguyên.

| Thuật ngữ | Nghĩa trong tài liệu |
|---|---|
| Module nguồn | Nơi sở hữu sự việc và trạng thái thật: User, Order hoặc module tài chính/hỗ trợ tương ứng. |
| Event nghiệp vụ | Thông tin về sự việc đã được nguồn xác nhận, thường phát sau khi nghiệp vụ được lưu thành công. |
| Yêu cầu gửi | Lệnh gửi có mục đích rõ; nguồn có thể cần phản hồi ngay, như yêu cầu gửi OTP do User vừa phát hành. |
| Hành trình | Nhóm công việc lớn như gia nhập tổ chức hoặc giao hàng; dùng để sắp xếp, không tự gửi tin. |
| Luồng thông báo | Một mục tiêu thông báo, ví dụ báo giao thất bại; có điều kiện kích hoạt và điều kiện dừng. |
| Trường hợp gửi | Cách thông báo cho một nhóm người nhận của luồng; quy định nội dung, thời gian, kênh và quyền nhận. |
| Bản thông báo | Nội dung dành cho một người nhận trong một trường hợp gửi; có thể được đưa qua nhiều kênh. |
| Lần thử gửi | Một lần đưa bản thông báo tới một kênh/nhà cung cấp kênh gửi. Thử lại không phải sự việc nghiệp vụ mới. |
| Kênh chính | Kênh được chọn đầu tiên trong một trường hợp gửi. |
| Kênh song song | Gửi cùng thời điểm theo chính sách, như Push báo ngay và In-app để xem lại. |
| Kênh dự phòng | Chỉ dùng khi có điều kiện chuyển kênh đã được duyệt; ví dụ SMS khi ZNS xác nhận không gửi được. |
| Kênh bổ sung | Gửi thêm vì mục tiêu khác, như Email báo cáo chi tiết sau thông báo ngắn. Không phải thử lại. |
| Phạm vi | Ngữ cảnh tổ chức/Shop, ứng dụng, thương hiệu, pháp nhân và đối tượng nghiệp vụ được phép sử dụng. |
| Quy tắc chọn người nhận | Quy tắc đã được duyệt dùng dữ liệu do nguồn/User cung cấp, như đầu mối vận hành của đúng Shop. |
| Đối tượng nhận thông báo | Người/tài khoản hoặc liên hệ được chỉ định nhận tin. Có thể có hoặc không có tài khoản User. |
| Thông tin nhận thông báo | Dữ liệu tối thiểu về định danh, liên hệ và đích gửi của đối tượng nhận; tham chiếu User/nguồn, không thay hồ sơ chính thức. |
| Nhóm nhận thông báo | Tập đối tượng nhận cho cùng mục đích; không phải vai trò hoặc nhóm quyền trong User. |
| Ngữ cảnh thông báo | Tổ chức/Shop, app và phạm vi nghiệp vụ mà tin thuộc về; dùng để chọn cấu hình và tách hộp thông báo. Không đồng nhất với Access Context của phiên đăng nhập. |
| Dữ liệu sự việc | Dữ liệu của lần phát sinh cụ thể như đơn, chặng, kết quả và thời hạn; do module nguồn cung cấp. |
| Lựa chọn nhận thông báo | Thiết lập nhận/từ chối theo luồng, kênh và phạm vi; không thay quyền truy cập User. |
| Kết nối kênh gửi | Cấu hình liên kết Notification với kênh/nhà cung cấp và định danh gửi phù hợp. |
| Nhà cung cấp kênh gửi | Đơn vị/dịch vụ tiếp nhận yêu cầu gửi qua một kênh và trả kết quả hỗ trợ. |
| Hộp thông báo | Nơi người dùng xem thông báo In-app theo đúng người và ngữ cảnh được phép. |
| Bước xử lý thông báo | Một bước trong luồng: kiểm tra điều kiện, chờ, gom nhóm hoặc gửi qua kênh. |
| Mẫu thông báo | Nội dung và biến dữ liệu được duyệt để tạo tin cho từng kênh/đối tượng. |
| Đồng ý / từ chối | Lựa chọn cho một mục đích nhận tin và phạm vi; không đồng nhất với chọn kênh ưa thích. |
| Nhà vận chuyển (NVC) | Đơn vị thực hiện vận chuyển; vận đơn NVC không đồng nhất với đơn SuperPlatform. |

Ví dụ: Order xác nhận **giao thất bại**. Notification áp dụng luồng “Giao thất bại”, có thể tạo một trường hợp gửi cho người nhận hàng và một trường hợp cho Shop. Hai trường hợp này khác người nhận, nội dung và kênh, nhưng đều xuất phát từ cùng một sự việc.

Các tên delivery.started, delivery.failed, delivery.completed chỉ là ví dụ/đề xuất, chưa phải event đã tích hợp. BRD không bắt buộc phải có một module Delivery riêng; hiện trạng giao hàng phải do Order hoặc module được giao quyền sở hữu xác nhận.

## 4. Người sử dụng, actor và trách nhiệm

### 4.1. Nhu cầu theo góc nhìn người sử dụng

| Actor | Tôi cần gì? | Giới hạn |
|---|---|---|
| Người dùng tài khoản | Nhận mã/link/cảnh báo bảo mật đúng liên hệ, biết lý do và bước tiếp theo. | Không thấy mã, liên hệ hoặc dữ liệu của người khác. |
| Chủ/thành viên Shop | Biết đơn nào cần xử lý, xem lại thông báo và mở đúng đơn của Shop. | Quyền nhận/xem gắn với tư cách và quyền được cấp. |
| Người nhận hàng không có tài khoản | Biết đơn liên quan đang giao hoặc gặp vấn đề, có hướng xử lý dễ hiểu. | Không tự có hộp In-app hoặc quyền xem toàn bộ hồ sơ Shop. |
| Nhân sự vận hành/người duyệt | Nhận việc thuộc trách nhiệm, không bị nhắc sau khi việc đã xử lý. | Nhóm chịu trách nhiệm do module nghiệp vụ cung cấp/xác nhận. |
| Đầu mối doanh nghiệp/kế toán | Nhận thông tin vận hành/tài chính đúng chức trách và pháp nhân. | Không suy ra quyền nhận tiền/COD chỉ vì cùng tổ chức. |
| Quản trị nội dung/chính sách | Soạn, xem trước, thử, trình duyệt và kích hoạt cấu hình đúng phạm vi. | Chỉ được chỉnh phần được cấp quyền; không vượt quy tắc khóa. |
| Nhân sự vận hành Notification | Tìm nguyên nhân, xử lý tồn đọng, gửi lại hoặc dừng có kiểm soát. | Không sửa trạng thái đơn/tài khoản; không được đọc bí mật để xử lý lỗi. |
| Chủ ngân sách | Biết chi phí, hạn mức, nguồn phát sinh và cơ chế xử lý khi hết hạn mức. | Không được vô hiệu hóa yêu cầu bảo vệ dữ liệu vì mục tiêu tiết kiệm. |

### 4.2. Module nào chịu trách nhiệm phần nào?

| Bên | Trách nhiệm | Notification cần từ bên này |
|---|---|---|
| User | Định danh, liên hệ, tư cách, phát hành/xác minh mã và link, trạng thái lời mời/tài khoản. | Đích nhận đã chọn, thời hạn, dữ liệu tối thiểu; bản đọc người nhận được duyệt nếu dùng quy tắc chọn người nhận. |
| Order | Đơn/chặng/vận đơn, trạng thái chuẩn, kết quả thao tác và các liên hệ phát sinh theo đơn. | Sự việc đã chốt, đơn/chặng liên quan, người nhận/Shop và điều kiện còn hiệu lực. |
| Authorization / chức năng phân quyền của nền tảng | Quyết định quyền truy cập theo người, hành động và phạm vi. | Căn cứ cấp/thu hồi quyền để cấu hình, xem lịch sử và mở nội dung được kiểm soát. |
| User + chủ quản ứng dụng/thương hiệu | User sở hữu danh mục app/client; chủ quản chịu trách nhiệm nhận diện thương hiệu, định danh gửi và liên kết hợp lệ. | Tham chiếu app/client từ User và cấu hình thương hiệu đã được duyệt. Không yêu cầu tạo lại danh mục app ở Notification. |
| Tài chính, hỗ trợ, giá/sản phẩm khi tích hợp | Số tiền, kết quả nghiệp vụ và người chịu trách nhiệm của chính module đó. | Hợp đồng nguồn riêng; không lấy dữ liệu Order để tự suy ra kết quả đối soát. |
| Notification | Chọn luồng/trường hợp gửi, kiểm tra chính sách, dựng nội dung, điều phối kênh, theo dõi và trả kết quả. | Không thay đổi trạng thái nghiệp vụ của bên khác. |
| Nhà cung cấp kênh gửi | Tiếp nhận gửi và trả kết quả trong phạm vi kênh hỗ trợ. | Mã tham chiếu, trạng thái và lỗi; không mặc định có xác nhận đã đọc/đã giao. |
| Chủ nghiệp vụ / bảo vệ dữ liệu / vận hành | Duyệt mục đích, tập nhận, nội dung, chi phí, thời gian lưu và tiêu chí nghiệm thu. | Người chịu trách nhiệm cho từng luồng và quyết định còn mở. |

Chức danh trong bảng là trách nhiệm cần phân công, không phải khẳng định các vai trò hoặc quyền đã tồn tại trong User.

### 4.3. Phần nào cấu hình ở User, phần nào ở Notification?

BRD cần mô tả cách thông báo tới đúng người, đúng app và đúng phạm vi. Các yêu cầu này sử dụng năng lực User/Authorization đã có. Người quản trị Notification không phải tạo lại tài khoản, tổ chức, ứng dụng hoặc vai trò.

| Nội dung | Nơi quản lý dữ liệu/quyết định gốc | Notification làm gì? |
|---|---|---|
| Người dùng, định danh và liên hệ chính thức | User | Tham chiếu định danh ổn định, dùng bản đọc tối thiểu hoặc dữ liệu được nguồn cấp; không sửa hồ sơ gốc. |
| Tổ chức/Shop và tư cách thành viên | User | Dùng đúng ngữ cảnh sự việc và tư cách còn hiệu lực để áp dụng quy tắc chọn người nhận. |
| Danh mục Application/Client, domain tin cậy và điều kiện tổ chức được dùng app | User | Chọn app/client đã đăng ký; áp dụng điều kiện truy cập do User cung cấp, không tự cấp Entitlement. |
| Role, Permission, Scope và quyết định thao tác | User/Authorization | Đăng ký các hành động/tài nguyên Notification theo cơ chế catalog chung, kiểm tra quyết định tại backend; không tạo bộ vai trò độc lập. |
| Ai liên quan một đơn, ai được giao xử lý/duyệt một việc | Module sở hữu nghiệp vụ | Nhận người cụ thể hoặc quy tắc chọn người nhận có căn cứ; không coi mọi người có quyền xem đơn đều cần nhận mọi tin. |
| Event nào kích hoạt, trường hợp gửi, mẫu, thời điểm, điều kiện và phối hợp kênh | Notification, chủ nghiệp vụ duyệt | Đây là phần cấu hình luồng thông báo. |
| Dữ liệu phục vụ gửi: thông tin nhận thông báo tối thiểu, đăng ký thiết bị Push theo app, hộp thông báo, lịch sử và lựa chọn nhận thông báo | Notification hoặc năng lực nền tảng được chỉ định | Quản lý dữ liệu giao tin trong phạm vi được giao; tham chiếu User, đồng bộ thay đổi, không thay hồ sơ/quyền gốc. Nơi đăng ký thiết bị phải được chốt khi đặc tả tích hợp. |

Ba điều kiện phải được phân biệt:

1. **Cần nhận tin:** Người này liên quan sự việc theo chính sách luồng.
2. **Được nhận nội dung/app này:** Thông tin và đích gửi phù hợp với quyền/phạm vi và lựa chọn nhận thông báo.
3. **Được thực hiện hành động sau khi mở:** Module nguồn kiểm tra quyền hiện tại khi người dùng xem đơn, duyệt việc hoặc thao tác.

Có quyền xem không mặc định phải nhận tin. Được nhận tin không cấp thêm quyền thao tác. Lựa chọn nhận thông báo của Notification cũng không thay đổi quyền do User cấp.

Trong User, Access Context phản ánh app và membership của phiên tương tác. Notification dùng ngữ cảnh đã xác thực khi người dùng xem hộp thông báo hoặc quản trị. Với event chạy nền, không yêu cầu người nhận đang đăng nhập hoặc giữ phiên của họ; nguồn và hợp đồng User/Authorization phải cung cấp căn cứ phù hợp cho xử lý nền. Không dùng token của người tạo đơn để suy quyền của tất cả người nhận.

### 4.4. Tham khảo mô hình Novu

Novu tổ chức gửi tin qua luồng và các bước xử lý; khi kích hoạt, ứng dụng cung cấp luồng, đối tượng nhận thông báo và dữ liệu sự việc. Đây là mô hình tham khảo phù hợp với định hướng nhận event rồi điều phối gửi của Notification. [Nguồn: Novu Workflows](https://docs.novu.co/platform/concepts/workflows).

| Tên tham khảo từ Novu | Tên gọi trong SPF | Ý nghĩa giữ nguyên |
|---|---|---|
| Workflow | Luồng thông báo | Các trường hợp gửi và bước kênh/thời gian/điều kiện do Notification quản lý. |
| Subscriber | Đối tượng nhận thông báo | Có thông tin nhận tối thiểu gắn định danh ổn định từ User/nguồn; không tự tạo tài khoản User cho liên hệ trên đơn. |
| Payload | Dữ liệu sự việc | Đơn, chặng, kết quả, thời hạn của sự việc; do module nguồn cung cấp. |
| Context | Ngữ cảnh thông báo | Tổ chức/Shop, app và phạm vi để định tuyến, chọn nội dung, tách hộp thông báo. |
| Topic | Nhóm nhận thông báo | Nhóm phân phối cho một mục đích; nguồn cung cấp hoặc đồng bộ; không phải Role/quyền nghiệp vụ. |
| Preferences | Lựa chọn nhận thông báo | Theo luồng/kênh/phạm vi; độc lập với Role/Permission/Scope của User. |
| Integration | Kết nối kênh gửi | Kết nối và định danh gửi; không thay Application/Client hoặc Entitlement của User. |
| Provider | Nhà cung cấp kênh gửi | Dịch vụ nhận yêu cầu gửi và trả kết quả của kênh. |
| Inbox | Hộp thông báo | Thông báo In-app của đúng người trong ngữ cảnh được backend xác nhận. |
| Step | Bước xử lý thông báo | Bước kênh, thời gian hoặc điều kiện trong luồng. |
| Template | Mẫu thông báo | Nội dung/biến được duyệt để dựng tin. |

Novu dùng định danh đối tượng nhận ổn định và cho phép lưu dữ liệu/đích phục vụ gửi tin; ngữ cảnh thông báo hỗ trợ phân tách theo tổ chức/app. Với SuperPlatform, các bản đọc này phải dẫn chiếu User hoặc nguồn sở hữu dữ liệu. Đây là cách áp dụng đề xuất, không khẳng định đã có tích hợp Novu. [Nguồn: Subscribers](https://docs.novu.co/platform/concepts/subscribers), [Contexts](https://docs.novu.co/platform/concepts/contexts).

Nhóm nhận thông báo chỉ xác định tập nhận; lựa chọn nhận thông báo kiểm soát việc nhận tin. Vì vậy nhóm nhận thông báo hoặc lựa chọn bật một luồng không được dùng làm bằng chứng quyền xem dữ liệu của Shop. [Nguồn: Topics](https://docs.novu.co/platform/concepts/topics), [Preferences](https://docs.novu.co/platform/concepts/preferences).

Phải xác thực người và ngữ cảnh hộp thông báo tại backend; không tin các ID người/app/Shop do trình duyệt tự gửi. Novu cũng có cơ chế bảo vệ đối tượng nhận/ngữ cảnh của hộp thông báo. Cơ chế cụ thể của SuperPlatform cần đặc tả dựa trên User/Access Context đang có. [Nguồn: Inbox production](https://docs.novu.co/platform/inbox/prepare-for-production), [Contexts](https://docs.novu.co/platform/concepts/contexts).

Application của User là ứng dụng nghiệp vụ; applicationIdentifier của Novu Inbox là định danh kết nối môi trường Novu, không tự chứng minh quyền vào SuperShip/SuperAI. Ngữ cảnh app nghiệp vụ phải được ánh xạ rõ, không đồng nhất hai loại định danh.

Tham khảo Novu không đồng nghĩa chọn Novu làm nền tảng triển khai, dùng nguyên API của Novu hoặc mặc định Novu có tích hợp Zalo phù hợp. Phạm vi đợt đầu vẫn là 5 kênh đã chốt; lựa chọn giải pháp/nhà cung cấp kênh gửi thuộc đặc tả sau.

## 5. Phạm vi và lộ trình đề xuất

### 5.1. Phạm vi mục tiêu của module

Giữ phạm vi dài hạn từ BRD đầu vào:

- Tài khoản, tổ chức và bảo mật; đơn hàng, lấy/giao/chuyển hoàn; tạo vận đơn, đổi NVC và bàn giao chặng.
- COD, ví, đối soát, thanh toán, công nợ, hóa đơn; giá/sản phẩm; hỗ trợ/khiếu nại; chính sách, sự cố và Marketing.
- In-app, Push, SMS, Email, Zalo ZNS; định danh gửi theo app/brand/pháp nhân.
- Quản lý luồng, quy tắc chọn người nhận, mẫu/phiên bản, lịch, chống trùng, giới hạn tần suất, gom nhóm, điều kiện dừng.
- Lựa chọn nhận thông báo, chiến dịch, cấu hình theo phạm vi, chi phí/hạn mức, tra cứu và nhật ký quản trị.

Mỗi phạm vi chỉ được kích hoạt khi có module/chủ nghiệp vụ sở hữu, hợp đồng đầu vào, người nhận, mẫu, kênh và nghiệm thu tương ứng.

### 5.2. Lộ trình để lựa chọn phạm vi triển khai

**Đã xác nhận:** Đợt đầu bao gồm đầy đủ 5 kênh dưới đây, giữ phạm vi kênh của Phase 1 bản đầu.

| Kênh | Cách hiểu trong BRD | Điều kiện sử dụng |
|---|---|---|
| In-app | Thông báo trong ứng dụng, có danh sách để xem lại và trạng thái đã đọc/chưa đọc. | Đúng tài khoản, phạm vi và quyền truy cập. |
| Push | Thông báo đẩy tới thiết bị của người dùng. | Đúng tài khoản, thiết bị, ứng dụng đích và cấu hình Push hợp lệ. |
| Email | Gửi tới địa chỉ email, bao gồm Gmail. Gmail là dịch vụ email; nhà cung cấp gửi chưa được chỉ định. | Địa chỉ hợp lệ và định danh gửi đã được duyệt. |
| SMS | Gửi tin nhắn tới số điện thoại. | Số điện thoại hợp lệ, định danh gửi và điều kiện kênh được duyệt. |
| Zalo | Gửi thông báo qua Zalo ZNS theo phạm vi kênh đang mô tả trong BRD. | Liên hệ, mẫu, định danh gửi và điều kiện nhà cung cấp kênh gửi phù hợp. |

Hỗ trợ đủ 5 kênh không có nghĩa mỗi thông báo phải gửi qua cả 5. Mỗi trường hợp gửi vẫn chọn kênh chính, song song, dự phòng hoặc bổ sung theo mục đích. Ngày 08/10/2026, người yêu cầu xác nhận các lát cắt OTP (UC-USR-01), lời mời nhân viên/thành viên (UC-USR-03/04), giao thất bại (UC-ORD-05), cảnh báo bảo mật/thiết bị (UC-USR-09), cùng hộp tin và tra cứu phục vụ các luồng này. Shop SuperShip và giao diện nội bộ là điểm tích hợp ưu tiên; client Push cụ thể và app mở rộng còn cần chốt.

Thông báo tài khoản và thông báo công việc được tách ngữ cảnh: đổi Shop không làm mất cảnh báo tài khoản, còn thông báo đơn phải đúng Shop/app/phạm vi. Notification tiếp nối chuông hiện có; User/Order vẫn sở hữu và kiểm tra các hành động sau khi mở. Phạm vi cụ thể và các mốc hoàn thành được quản lý tại [phạm vi bản phát hành](<docs/Notitek - PHẠM VI VÀ KẾ HOẠCH PHÁT HÀNH ĐẦU TIÊN.md>).

| Đợt đề xuất | Nội dung | Điều kiện hoàn thành |
|---|---|---|
| P0 - nền tảng và tích hợp đầu | Năng lực nhận event/yêu cầu gửi, chống trùng, người nhận/phạm vi, mẫu, đầy đủ 5 kênh, hộp In-app, kết quả, bảo vệ dữ liệu và tra cứu cơ bản cho các lát cắt đã chọn. Các UC P0 khác trong danh mục dài hạn không tự thuộc bản phát hành này. | Có hợp đồng nguồn; nghiệm thu các lát cắt đã chọn và gửi thật của cả 5 kênh theo mục 9. Luồng thiếu đường phát/event nguồn chưa được tuyên bố hoàn tất. |
| P1 - trải nghiệm và vận hành | Lựa chọn nhận thông báo, quản trị tự phục vụ, lịch/gom nhóm, quản lý hạn mức và luồng P1. | Có nguồn quyền/tư cách, dữ liệu lựa chọn nhận thông báo, quyền quản trị và chính sách chi phí đã chốt. |
| P2 - mở rộng | Module tài chính/hỗ trợ/giá, chiến dịch, đánh giá dịch vụ, A/B Testing và mở rộng app/kênh theo nhu cầu. | Có chủ nghiệp vụ và hợp đồng riêng; Marketing có căn cứ đồng ý, tập nhận và phê duyệt. |

Phạm vi 5 kênh và các lát cắt đợt đầu đã được xác nhận. Các luồng còn lại và quản trị nâng cao được phân kỳ riêng; danh mục dài hạn không là cam kết hoàn thành đồng thời. Hợp đồng, chính sách và thời điểm bật gửi thật/chuyển đổi được chốt theo từng phần.

### 5.3. Ngoài trách nhiệm của Notification

- Không tạo/xác thực OTP, kích hoạt tài khoản, duyệt lời mời, sửa đơn, chốt COD hoặc tự giải quyết hỗ trợ.
- Không chuyển lỗi nhập liệu, HTTP 500 hoặc thao tác chưa hoàn tất trên màn hình thành thông báo cho mọi người. Cảnh báo vận hành có thể là luồng riêng khi chủ nghiệp vụ xác nhận cần xử lý.
- Không thay thế User, phân quyền, HRM/CRM, trung tâm hội thoại hoặc hệ thống webhook nghiệp vụ.
- Không tự vận hành hạ tầng cung cấp SMS/Email/ZNS/Push thay nhà cung cấp kênh gửi; không tự viết nội dung bằng AI hoặc tự quyết định hành động nghiệp vụ.
- Không lưu chứng từ hoặc hồ sơ nhạy cảm không cần thiết trong thông báo.

Notification vẫn nhận đầu vào và phát event kết quả cho module khác. Việc này phục vụ theo dõi gửi tin, không thay thế event/webhook trao đổi trạng thái đơn hoặc thanh toán.

## 6. Luồng xử lý và hợp đồng nghiệp vụ

### 6.1. Hai cách module khác gọi Notification

**Cách A - báo sự việc đã xảy ra:** Module nguồn lưu nghiệp vụ thành công rồi phát event. Notification tiếp nhận và xử lý; nếu gửi lỗi, không đảo ngược việc đổi mật khẩu, hủy đơn hoặc lời mời đã lưu.

**Cách B - yêu cầu gửi có phản hồi:** User phát hành mã/link rồi yêu cầu gửi qua hợp đồng an toàn. Nguồn cần biết yêu cầu gửi thật đã thành công đến mức nào hoặc thất bại/chưa rõ để quyết định bước tiếp theo. User vẫn sở hữu mã, thời hạn và việc xác minh.

OTP không được coi thành công chỉ vì Notification nhận event hoặc ghi log. “Gửi thật thành công” tối thiểu phải có bằng chứng gửi qua kênh thực theo tiêu chí hợp đồng đã chốt; không được hứa người nhận đã đọc. Khi phản hồi chưa rõ, nguồn tra cứu/cấp lại theo chính sách, không tự tạo nhiều mã hoặc gửi lại vô hạn.

User quyết định kích hoạt, cấp lại và vô hiệu hóa mã; Notification chỉ biết yêu cầu gửi, reference, hạn gửi và chỉ dẫn thử lại/hủy theo hợp đồng chung. Tương tự, Order xác định giao đủ/giao một phần và việc còn cần thông báo; Notification không đọc mô hình đơn để tự tính điều kiện nghiệp vụ. Các tình huống này là kịch bản tích hợp, không chuyển quyền sở hữu nghiệp vụ vào Notification.

TOTP do ứng dụng xác thực sinh tại thiết bị không phải luồng gửi OTP qua Notification. Mã/link nhạy cảm chỉ đi qua đường tích hợp được bảo vệ, không đưa vào bus phổ thông, log hoặc event kết quả.

### 6.2. Nguồn cần cung cấp gì?

| Thông tin nghiệp vụ | Mục đích |
|---|---|
| Nguồn được phép, loại/phiên bản sự việc, định danh duy nhất, thời điểm xảy ra | Hiểu sự việc, kiểm tra hợp đồng và chống nhận trùng. |
| Đối tượng nghiệp vụ, phiên bản/mốc thay đổi nếu liên quan, đơn/chặng/vận đơn khi cần | Phân biệt các lần giao thực sự, kiểm tra tin cũ và truy vết. |
| Tổ chức/Shop, app/brand/pháp nhân và bối cảnh áp dụng | Chọn cấu hình, nội dung, người nhận và đơn vị chịu phí đúng phạm vi. |
| Người nhận cụ thể hoặc quy tắc chọn người nhận và căn cứ được duyệt | Không suy đoán người nhận từ số điện thoại hoặc quan hệ tổ chức. |
| Dữ liệu tối thiểu để dựng mẫu và liên kết hợp lệ | Nói đúng sự việc, không truyền cả hồ sơ nguồn. |
| Hạn dùng, điều kiện dừng và cách xác nhận còn hiệu lực | Tránh gửi mã/link cũ hoặc lời nhắc sau khi việc đã xử lý. |
| Tham chiếu yêu cầu/correlation và nhu cầu phản hồi | Nguồn tra được từ sự việc đến kết quả gửi. |

Không yêu cầu mọi trường giống nhau cho mọi luồng; hợp đồng từng luồng xác định trường bắt buộc. Đây là nội dung nghiệp vụ cần trao đổi, chưa phải schema event/API.

### 6.3. Các bước người gọi có thể theo dõi

1. Xác thực nguồn, kiểm tra hợp đồng và phạm vi. Yêu cầu không hợp lệ được từ chối có lý do, chưa được coi là đã tiếp nhận thành công.
2. Ghi nhận yêu cầu và định danh theo dõi; xử lý nhận trùng bằng tham chiếu cũ.
3. Chọn luồng và các trường hợp gửi; xác định người nhận, quyền nhận và cấu hình có hiệu lực.
4. Kiểm tra thời hạn/trạng thái hiện tại, dữ liệu mẫu, kênh và ngân sách áp dụng.
5. Tạo hoặc xếp lịch bản thông báo; kiểm tra lại điều kiện có thể thay đổi ngay trước gửi.
6. Gửi, nhận phản hồi, xử lý thử lại/chuyển kênh theo quy tắc.
7. Lưu kết quả từng người nhận/kênh, cập nhật kết quả tổng hợp và phát kết quả cho bên được phép.

Thiếu người nhận, mẫu hay dữ liệu không được biến thành trạng thái “đã gửi”. Phải biết đang chờ bổ sung hay đã dừng, ai xử lý và hết thời hạn nào.

### 6.4. Ý nghĩa kết quả

| Kết quả | Nghĩa với người gọi/người vận hành |
|---|---|
| Bị từ chối | Đầu vào/nguồn không hợp lệ; có lý do, không có cam kết xử lý gửi. |
| Đã tiếp nhận | Đã ghi nhận bền vững và có tham chiếu; chưa chứng minh đã gửi. |
| Đang chờ / đang xử lý | Chờ lịch, dữ liệu hoặc đang đưa tới kênh; có lý do và hạn xử lý. |
| Không gửi theo chính sách | Yêu cầu hợp lệ nhưng bị loại vì quy tắc, không có luồng áp dụng, quyền nhận hoặc không còn cần thông báo. |
| Đã hủy trước gửi / hết hạn | Bản chờ đã dừng hoặc mất hiệu lực; không tiếp tục gửi mã/link cũ. |
| Nhà cung cấp đã nhận | Nhà cung cấp kênh gửi xác nhận tiếp nhận gửi thật; chưa chứng minh đã tới thiết bị/hộp thư. |
| Chưa rõ kết quả | Có thể nhà cung cấp kênh gửi đã nhận nhưng phản hồi mất/chậm. Không được ghi chắc chắn thất bại hoặc thành công. |
| Đã giao theo bằng chứng của kênh | Có xác nhận phù hợp với kênh; phải nêu loại bằng chứng, không tự suy ra người đọc. |
| Thất bại cuối cùng | Không còn phương án hợp lệ sau quy tắc thử lại; có nguyên nhân, không đồng nhất với bị chặn/hết hạn. |
| Đã đọc | Thao tác đọc In-app được ghi nhận; không có nghĩa đã xử lý đơn hoặc phê duyệt công việc. |

Một yêu cầu có nhiều người nhận/kênh phải có kết quả riêng. Nếu có người thành công, người thất bại, báo **kết quả một phần** kèm số lượng/chi tiết được phép xem; không ghi chung “thành công” cho tất cả. Tiêu chí thành công tổng hợp phụ thuộc kênh bắt buộc của từng trường hợp gửi.

### 6.5. Event kết quả Notification

Notification phát kết quả để nguồn biết tiến trình gửi, không phát lệnh sửa đơn/tài khoản. Các tên dưới đây là đề xuất và phải chốt hợp đồng trước tích hợp:

| Event đề xuất | Tình huống |
|---|---|
| notification.accepted.v1 | Yêu cầu đã được ghi nhận. |
| notification.suppressed.v1 | Quy tắc quyết định không gửi, có mã lý do. |
| notification.cancelled.v1 / notification.expired.v1 | Dừng bản chờ hoặc hết hạn. |
| notification.provider_accepted.v1 | Có bằng chứng nhà cung cấp kênh gửi nhận gửi. |
| notification.delivered.v1 | Có bằng chứng giao phù hợp với kênh. |
| notification.failed.v1 | Thất bại cuối cùng. |
| notification.read.v1 | Người dùng đọc In-app. |

Mỗi kết quả có định danh riêng, tham chiếu yêu cầu gốc, mức kết quả (yêu cầu/người nhận/kênh), thời điểm và lý do khi cần. Không phát OTP, token, nội dung nhạy cảm hoặc liên hệ đầy đủ. Nguồn phải xử lý được event kết quả gửi lại; không tạo vòng lặp gửi tin từ chính event kết quả. Bên không cần tích hợp realtime có thể tra cứu thay vì đăng ký mọi loại kết quả.

### 6.6. Người quản trị cấu hình một luồng như thế nào?

Người quản trị mở Notification trong app/phạm vi được User xác nhận, sau đó thực hiện:

1. **Chọn sự việc và mục tiêu:** Chọn module nguồn, sự việc đã có hợp đồng và luồng cần thông báo; ví dụ Order xác nhận giao thất bại.
2. **Chọn phạm vi/app đích:** Tham chiếu tổ chức/Shop và Application đã đăng ký ở User. App đích của tin có thể khác app đang mở màn hình cấu hình nếu người quản trị được cấp quyền tương ứng.
3. **Chọn trường hợp gửi và người nhận:** Người cụ thể do nguồn cấp hoặc quy tắc chọn người nhận/nhóm nhận thông báo được duyệt. Không tạo người dùng/vai trò mới tại bước này.
4. **Chọn nội dung và kênh:** Mẫu được duyệt, các biến nguồn, định danh gửi, kênh chính và điều kiện song song/dự phòng/bổ sung.
5. **Chọn thời gian và điều kiện:** Gửi ngay/nhắc, hạn dùng, điều kiện dừng, quyền nhận tin, giới hạn tần suất và ngân sách.
6. **Xem trước và kiểm tra:** Hiển thị cấu hình kế thừa/thực tế, nguồn dữ liệu người nhận, app/client đích, biến mẫu và lỗi còn thiếu. Thử bằng dữ liệu/tập nhận thử được phép.
7. **Duyệt và kích hoạt:** Theo quyền và quy trình chung; lưu phiên bản. Nguồn phát event hoặc gọi yêu cầu gửi, Notification áp dụng cấu hình và trả kết quả.

| Phần của một cấu hình | Ví dụ giao thất bại |
|---|---|
| Nguồn / sự việc | Order xác nhận một lần giao thất bại của đơn/chặng; tên event cần hợp đồng nguồn. |
| Phạm vi / app đích | Shop A, app SUPERSHIP trong danh mục User. |
| Trường hợp Shop | Đầu mối Shop do nguồn cung cấp hoặc quy tắc chọn người nhận được duyệt; In-app và Push; chỉ nội dung có quyền nhận. |
| Trường hợp người nhận hàng | Liên hệ trên đơn do Order xác nhận; Zalo ZNS chính, SMS dự phòng theo điều kiện. Không yêu cầu tài khoản/app access. |
| Điều kiện dừng | Order xác nhận mốc mới làm lời nhắc không còn giá trị; không lấy “đã đọc” thay kết quả xử lý. |
| Kết quả | Theo từng người/kênh và tham chiếu sự việc; không thay đổi trạng thái đơn. |

Ví dụ minh họa cách cấu hình, chưa duyệt quy tắc chọn người nhận/kênh cho luồng giao thất bại. Màn hình quản trị nâng cao có thể triển khai sau; cấu hình tối thiểu để luồng chạy vẫn phải đủ các nội dung trên ngay khi đưa luồng vào sử dụng.

### 6.7. Đúng người, đúng app và đúng phạm vi khi gửi/xem

- Nguồn gửi một sự việc kèm ngữ cảnh nghiệp vụ đáng tin cậy. Notification kiểm tra và áp dụng app đích của luồng, không tự chọn app mà người nhận đang mở trên trình duyệt.
- Với In-app/Push của tài khoản, căn cứ định danh, tư cách và điều kiện app do User cung cấp; áp dụng chính sách nhận nội dung. Không tự viết lại quy tắc Entitlement, MFA hoặc các ngoại lệ truy cập app của User.
- Một người có thể nhận tin ở nhiều Shop/app. Hộp thông báo của Shop A/SuperShip không được trộn dữ liệu Shop B/SuperAI khi đổi ngữ cảnh thông báo; Push phải tới thiết bị/client đúng app.
- Khi người dùng mở hộp thông báo, đọc hoặc mở liên kết, backend kiểm tra người/ngữ cảnh hiện tại. Đường dẫn tới đơn vẫn do Order kiểm tra quyền trên đơn.
- Liên hệ không có tài khoản nhận SMS/Email/Zalo theo căn cứ giao dịch, không phải có Entitlement hay phiên đăng nhập. OTP/cảnh báo liên hệ cũ đi theo chính sách bảo mật User đã chỉ định.

## 7. Yêu cầu nghiệp vụ

Các mã dưới đây là yêu cầu cụ thể của bản nháp hiện tại. Giá trị thời gian, hạn mức và danh sách được phép chưa chốt được quản lý ở mục 11.

### 7.1. CAT - Danh mục luồng

- **BR-CAT-01 [P0]:** Mỗi luồng có module/chủ nghiệp vụ chịu trách nhiệm, mục đích, điều kiện bắt đầu, trường hợp gửi, thời hạn và điều kiện dừng. Luồng thiếu các nội dung bắt buộc không được kích hoạt.
- **BR-CAT-02 [P0]:** Tách trường hợp gửi khi khác người nhận, nội dung, kênh hoặc thời điểm; không tạo luồng riêng chỉ vì đổi nhà cung cấp kênh gửi.
- **BR-CAT-03 [P1]:** Luồng có trạng thái nháp, chờ duyệt, hoạt động, tạm dừng, ngừng sử dụng; lưu phiên bản và người duyệt. Tạm dừng phải chỉ rõ chặn tin mới, tin đang chờ hay cả hai.
- **BR-CAT-04 [P0]:** Cấu hình tối thiểu của luồng phải có các thành phần ở mục 6.6. Người quản trị chỉ chọn/tham chiếu đối tượng từ User và dữ liệu nguồn; không tạo lại tài khoản, tổ chức, app hoặc role tại Notification.

### 7.2. TRG - Tiếp nhận và chống trùng

- **BR-TRG-01 [P0]:** Chỉ nguồn được phép mới được kích hoạt luồng trong phạm vi đã cấp. Event báo sự việc phải dựa trên nghiệp vụ đã được xác nhận; yêu cầu gửi mã do User phát hành đi theo hợp đồng riêng.
- **BR-TRG-02 [P0]:** Nhận lại cùng sự việc/yêu cầu trả về tham chiếu đã có, không tạo bản gửi mới vô ích. Khóa chống trùng phải phân biệt nguồn, trường hợp gửi, người nhận và chu kỳ nghiệp vụ; một lần giao lại thực sự có thể tạo thông báo mới.
- **BR-TRG-03 [P0]:** Event đến chậm/sai thứ tự được đánh giá theo thời điểm, phiên bản và điều kiện còn hiệu lực của chính luồng. Không loại toàn bộ event có phiên bản thấp hơn nếu đó là sự việc độc lập vẫn cần thông báo.
- **BR-TRG-04 [P0]:** Khi thiếu dữ liệu bắt buộc hoặc phiên bản hợp đồng không hỗ trợ, ghi lý do từ chối/chờ xử lý; không gửi bằng dữ liệu tự đoán. Nếu chờ, phải có thời hạn và hướng xử lý.
- **BR-TRG-05 [P0]:** Thử lại sau lỗi hệ thống phải tiếp tục từ kết quả đã biết; không gửi lại người/kênh đã xác nhận đạt tiêu chí chỉ vì kênh khác thất bại.

### 7.3. INT - Ranh giới và phản hồi

- **BR-INT-01 [P0]:** Module nguồn sở hữu sự việc, người nhận theo giao dịch và điều kiện nghiệp vụ. Notification chỉ áp dụng chính sách gửi đã duyệt, không ghi ngược trạng thái nghiệp vụ.
- **BR-INT-02 [P0]:** Hợp đồng yêu cầu gửi có phản hồi phải nói rõ mức thành công cần cho nguồn, thời hạn chờ và cách xử lý “chưa rõ”. OTP không dùng tiếp nhận/giả lập/ghi log làm bằng chứng gửi thật.
- **BR-INT-03 [P0]:** Lỗi thông báo sau một nghiệp vụ đã hoàn tất không tự hủy nghiệp vụ đó. Nguồn được biết lỗi để hiển thị phương án phù hợp, như cấp lại lời mời hoặc sao chép link bằng chức năng được phép.
- **BR-INT-04 [P0]:** Ghi nhận và cung cấp kết quả theo mục 6.4–6.5; có thể tra từ yêu cầu nguồn đến từng người nhận/kênh. Event kết quả chỉ gửi cho bên được phép, không mang bí mật.
- **BR-INT-05 [P0]:** Nếu luồng cần dữ liệu/trạng thái cập nhật, phải có nguồn và hợp đồng đồng bộ/kiểm tra được duyệt. Khi dữ liệu chưa đủ tin cậy, chờ hoặc dừng theo chính sách; không dựa vào bản đọc cũ để gửi tin nhạy cảm.
- **BR-INT-06 [P0]:** Tích hợp người dùng/app/quyền phải dùng cơ chế User/Authorization hiện có. Resolve cho phiên tương tác không mặc định là API liệt kê người nhận hoặc hợp đồng phân quyền cho xử lý nền; các luồng nền phải có căn cứ nguồn và hợp đồng phù hợp, không cần phiên đăng nhập của người nhận.

### 7.4. REC - Người nhận

- **BR-REC-01 [P0]:** Hỗ trợ tài khoản, liên hệ không có tài khoản và đầu mối tổ chức. Nguồn cung cấp người nhận theo sự việc hoặc quy tắc chọn người nhận được duyệt; không gắn số điện thoại của người nhận hàng vào tài khoản chỉ vì trùng số.
- **BR-REC-02 [P0]:** Quy tắc chọn người nhận chỉ dùng dữ liệu tối thiểu từ User/nguồn đúng phạm vi và còn hiệu lực. Không tự duyệt cây tổ chức để đoán quản lý hoặc kế toán; đầu mối nhận thông báo phải có căn cứ được module nguồn hoặc User cung cấp.
- **BR-REC-03 [P0]:** Một người thuộc nhiều Shop/tổ chức phải nhận và xem theo đúng ngữ cảnh từng sự việc. Cùng một người xuất hiện ở nhiều quy tắc chọn người nhận được hợp nhất trong cùng trường hợp gửi theo chính sách, không hợp nhất nội dung khác phạm vi.
- **BR-REC-04 [P0]:** Trước gửi lại/nhắc việc, kiểm tra điều kiện người nhận còn phù hợp. Việc xem hộp thông báo và mở đối tượng nguồn luôn tuân thủ quyền hiện tại, kể cả thông báo được tạo trước khi quyền bị thu hồi.
- **BR-REC-05 [P0]:** Luồng bảo mật có thể gửi đến liên hệ cũ do User chỉ định, như cảnh báo đổi số/email, hoặc gửi thông báo nghỉ việc cho người bị ảnh hưởng. Đây là ngoại lệ được duyệt theo mục đích, không tự yêu cầu tư cách còn hoạt động rồi loại cảnh báo cần thiết.
- **BR-REC-06 [P0]:** Thông tin nhận thông báo chỉ là bản đọc/đích phục vụ gửi với định danh ổn định; User vẫn sở hữu hồ sơ chính thức. Nhóm nhận thông báo không tạo quyền; thay tư cách/liên hệ phải cập nhật bản đọc theo hợp đồng, không tự hợp nhất người theo email/số điện thoại.

### 7.5. CLS - Mục đích và mức cần thiết

- **BR-CLS-01 [P0]:** Mỗi luồng có một mục đích chính: xác thực/bảo mật, giao dịch, yêu cầu hành động/nhắc việc, vận hành, thông tin dịch vụ hoặc Marketing. Mức khẩn và có bắt buộc gửi hay không là thuộc tính riêng, có người duyệt.
- **BR-CLS-02 [P0]:** Nội dung quảng cáo không được trộn vào mã xác thực/cảnh báo hoặc gắn nhãn giao dịch để bỏ qua lựa chọn Marketing.
- **BR-CLS-03 [P0]:** Danh sách luồng bắt buộc và ngoại lệ giờ hạn chế/hạn mức phải được duyệt riêng. Không mặc định mọi tin Order, COD hoặc “hệ thống” đều bắt buộc.

### 7.6. TPL - Mẫu và nội dung

- **BR-TPL-01 [P0]:** Tin gửi thật dùng nội dung/mẫu được duyệt cho luồng, người nhận, kênh, ngôn ngữ, app/brand và định danh gửi. Tin thủ công/chiến dịch vẫn phải đi qua soạn và duyệt, không dùng nội dung tự do để vượt kiểm soát.
- **BR-TPL-02 [P0]:** Mẫu thể hiện sự việc, đối tượng dễ nhận biết, hành động tiếp theo và đầu mối/đường dẫn khi cần. Thiếu biến quan trọng thì chờ/dừng có lý do; không gửi ô trống hoặc số tiền tự tính lại.
- **BR-TPL-03 [P0]:** Lưu phiên bản và bằng chứng nội dung đã gửi ở mức phù hợp để tra cứu; thay mẫu không sửa lịch sử. OTP/token không được lưu nguyên văn trong lịch sử hay màn hình hỗ trợ.
- **BR-TPL-04 [P1]:** Có xem trước/gửi thử với dữ liệu giả hoặc danh sách thử được phép; không gửi thử đến tập khách thật. Thay mẫu đã duyệt phải tạo phiên bản mới, chỉ rõ ảnh hưởng tin đang chờ.
- **BR-TPL-05 [P0]:** Chỉ dùng đường dẫn/app đích hợp lệ; bấm thông báo vẫn kiểm tra quyền hoặc token được nguồn cấp. Không nhúng dữ liệu tài chính/chứng từ nhạy cảm trong phần preview.

### 7.7. CHN - Kênh, gửi và kết quả

- **BR-CHN-01 [P0]:** Mỗi trường hợp gửi có một kênh chính; kênh song song, dự phòng, bổ sung phải có mục tiêu/điều kiện riêng. Kênh thiếu liên hệ, cấu hình hoặc điều kiện nhà cung cấp kênh gửi không được ghi là gửi thành công.
- **BR-CHN-02 [P0]:** Phân biệt lỗi tạm thời, lỗi không thể thử lại và kết quả chưa rõ. Giới hạn số lần/thời gian thử lại theo luồng, luôn nằm trong hạn dùng. Lỗi liên hệ/mẫu/quyền không được thử lại vô hạn.
- **BR-CHN-03 [P0]:** Chỉ chuyển kênh theo điều kiện đã duyệt. Timeout hoặc thiếu xác nhận giao không tự chứng minh thất bại; tra cứu/đợi kết quả hoặc áp dụng chính sách chấp nhận rủi ro trùng được duyệt trước khi chuyển.
- **BR-CHN-04 [P0]:** Ghi kết quả từng người nhận, kênh và lần thử; kết quả tổng hợp phải nêu một phần/chưa rõ khi có. Phản hồi nhà cung cấp kênh gửi gửi lại không tạo thêm lần gửi hay nhiều event kết quả cùng một thay đổi.
- **BR-CHN-05 [P0]:** In-app có danh sách đúng người/phạm vi, chưa đọc/đã đọc và liên kết hành động. Bản được ghi vào hộp thư chỉ chứng minh khả dụng để xem, không chứng minh đã đọc.
- **BR-CHN-06 [P0]:** Push dùng đích được xác nhận đúng tài khoản, thiết bị và app; xử lý phản hồi token mất hiệu lực và thay đổi đăng xuất/đổi chủ theo hợp đồng với bên sở hữu thiết bị. Không mặc định Notification sở hữu đăng ký/thu hồi thiết bị. Không hiển thị bí mật trên màn hình khóa; nhà cung cấp nhận không chứng minh người dùng đã xem.
- **BR-CHN-07 [P0]:** SMS và Zalo ZNS có liên hệ hợp lệ, định danh gửi, mẫu và điều kiện nhà cung cấp kênh gửi được duyệt. Người nhận hàng không cần có tài khoản; SMS không được coi tương đương ZNS về khả năng xác nhận giao.
- **BR-CHN-08 [P0]:** Email, bao gồm địa chỉ Gmail, dùng địa chỉ và định danh gửi đúng phạm vi; thư bị trả lại/chặn được ghi nhận. Không lộ các người nhận khác qua danh sách địa chỉ; tài liệu nhạy cảm dùng liên kết có kiểm soát.
- **BR-CHN-09 [P0]:** Chỉ ghi “đã giao” khi có loại bằng chứng đã được chốt cho kênh. Kênh không cung cấp bằng chứng giữ trạng thái tương ứng, không tự nâng thành đã giao/đã đọc.

### 7.8. TIME - Thời điểm, tần suất và điều kiện dừng

- **BR-TIME-01 [P0]:** Mỗi luồng quy định gửi ngay/hẹn giờ và hạn dùng. Không bắt đầu thử lại hoặc lịch nhắc bằng mã/link hết hạn; khi hết hạn cần kết quả rõ cho nguồn.
- **BR-TIME-02 [P0]:** Ngay trước gửi áp dụng hạn gửi và chỉ dẫn/điều kiện dừng do nguồn cung cấp hoặc xác nhận qua hợp đồng, như nguồn báo đơn/lời mời không còn cần nhắc. Notification không tự suy trạng thái nghiệp vụ. Chỉ hủy được tin chưa đưa tới nhà cung cấp kênh gửi; tin đã ra ngoài không được cam kết thu hồi.
- **BR-TIME-03 [P1]:** Khung giờ hạn chế áp dụng theo múi giờ đã xác định. Ngoại lệ cho xác thực/bảo mật phải được duyệt; không dời tin sang thời điểm sau hạn dùng.
- **BR-TIME-04 [P1]:** Giới hạn tần suất/gom nhóm tính theo người, luồng và phạm vi phù hợp. Các sự việc mới thật sự có thể được gom hoặc trì hoãn; chúng không phải event trùng. Tin yêu cầu hành động khẩn không bị gom làm mất thông tin cần xử lý.
- **BR-TIME-05 [P1]:** Nhắc lại có chu kỳ, giới hạn và điều kiện dừng rõ; việc hoàn thành do nguồn xác nhận, không dựa riêng vào “đã đọc”.

### 7.9. PREF - Lựa chọn và điều kiện nhận tin

- **BR-PREF-01 [P0]:** Kiểm tra mục đích, phạm vi đồng ý/từ chối, liên hệ và danh sách không gửi ngay trước gửi. Nếu Marketing chưa có dữ liệu đồng ý đủ căn cứ, không mặc định được gửi.
- **BR-PREF-02 [P0]:** Từ chối Marketing không tự chặn luồng bắt buộc đã được duyệt. Luồng bắt buộc vẫn phải tuân thủ liên hệ hợp lệ, bảo mật và điều kiện nhà cung cấp kênh gửi; không mặc định gửi qua mọi kênh.
- **BR-PREF-03 [P1]:** Người nhận quản lý lựa chọn nhận thông báo/kênh trong phạm vi cho phép và thấy lựa chọn nào không thể tắt cùng lý do. Lưu nguồn, thời điểm và phiên bản lựa chọn.
- **BR-PREF-04 [P0]:** Tách chặn Marketing, chặn một kênh và chặn vì liên hệ bị báo sai/xâm phạm. Với liên hệ không còn an toàn, luồng bảo mật làm theo chỉ dẫn được User xác nhận, không bỏ qua chặn để gửi bí mật.

### 7.10. SEC - Quyền và bảo vệ dữ liệu

- **BR-SEC-01 [P0]:** Chỉ nhận, dựng và lưu dữ liệu tối thiểu cho mục đích gửi. OTP, token, mật khẩu không xuất hiện trong log/event kết quả; thông tin liên hệ và số tiền được che theo quyền xem.
- **BR-SEC-02 [P0]:** Phân quyền riêng cho xem, soạn, duyệt, kích hoạt, cấu hình kênh, gửi lại, dừng và xuất lịch sử, theo phạm vi. Không suy ra quyền quản trị Notification từ quan hệ cha-con tổ chức.
- **BR-SEC-03 [P0]:** Người nhận chỉ xem dữ liệu được phép; gửi đúng địa chỉ không đủ để chứng minh quyền nhận nội dung tài chính. Thông báo/link không được tạo quyền mới ở module nguồn.
- **BR-SEC-04 [P0]:** Ghi nhật ký thay đổi cấu hình, phê duyệt, gửi lại, dừng và truy cập dữ liệu nhạy cảm, có người/thời điểm/phạm vi/lý do. Không sửa/xóa trái phép; lưu theo chính sách đã duyệt.
- **BR-SEC-05 [P0]:** Chính sách lưu, ẩn, xóa và tra cứu dữ liệu phải phân biệt bí mật ngắn hạn, nội dung thông báo, metadata và audit. Không dùng yêu cầu lưu lịch sử để lưu nguyên văn bí mật vô thời hạn.
- **BR-SEC-06 [P0]:** Khi quản trị hoặc dùng hộp thông báo, backend phải kiểm tra người/ngữ cảnh và quyền trên hành động/tài nguyên qua cơ chế phân quyền chung. Danh sách permissions, bộ lọc scope hoặc ID do frontend gửi không thay quyết định ALLOW; DENY, chưa đánh giá hoặc không kiểm tra được không cho phép thao tác.

### 7.11. CFG - Cấu hình và kế thừa

- **BR-CFG-01 [P0]:** Cấu hình áp dụng từ mặc định và các phạm vi được phép; chính sách bảo mật/bắt buộc bị khóa luôn có hiệu lực. Cấp dưới chỉ ghi đè phần được cho phép và do người có quyền thực hiện.
- **BR-CFG-02 [P0]:** Mỗi luồng có thứ tự ưu tiên phạm vi được duyệt cho tổ chức/Shop, app/brand, NVC, dịch vụ và các thuộc tính nguồn. Không tự cho rằng “cụ thể hơn” luôn thắng khi nhiều chiều xung đột. Xung đột chưa giải quyết phải chặn kích hoạt hoặc chờ có lý do.
- **BR-CFG-03 [P0]:** Tra được cấu hình thực tế và phiên bản được áp dụng cho từng bản gửi. Tin đang chờ phải kiểm tra lại quyền/điều kiện dừng; việc dùng mẫu cũ hay mới khi cấu hình đổi có chính sách rõ.
- **BR-CFG-04 [P1]:** Người quản trị thấy được phạm vi của mình, phần kế thừa, phần khóa và tác động dự kiến trước khi duyệt thay đổi. Không dùng cấu hình brand để vượt ranh giới pháp nhân/quyền dữ liệu.
- **BR-CFG-05 [P0]:** App/client đích và phạm vi nghiệp vụ phải dẫn chiếu danh mục User, có ánh xạ rõ với cấu hình kênh/nhà cung cấp kênh gửi. App đang dùng để quản trị, app nhận tin và định danh môi trường nhà cung cấp kênh gửi là các thông tin riêng; không tự tạo app hoặc cấp quyền khi kích hoạt luồng.

### 7.12. COST - Chi phí và hạn mức

- **BR-COST-01 [P0]:** Kênh phát sinh phí phải có đơn vị chịu phí và chính sách hạn mức được bên sở hữu ngân sách xác nhận trước gửi. Notification nhận/thực thi căn cứ theo phần trách nhiệm được giao trong hợp đồng; không mặc định sở hữu sổ ngân sách hoặc tự suy đơn vị trả phí từ NVC/app.
- **BR-COST-02 [P1]:** Theo dõi phí ước tính và phí nhà cung cấp kênh gửi xác nhận riêng, bao gồm thử lại/kênh dự phòng. Phí gửi tin không thay thế sổ COD/đối soát/thanh toán do module tài chính sở hữu.
- **BR-COST-03 [P0]:** Khi bên kiểm soát ngân sách không cho phép gửi, Notification tuân thủ quyết định/chỉ dẫn dừng/chờ/chuyển kênh theo hợp đồng và trả lý do. Không âm thầm loại OTP/cảnh báo rồi báo thành công hoặc tự vượt hạn mức. Bên sở hữu và cách kiểm soát được chốt tại DEC-08.

### 7.13. CMP - Chiến dịch

- **BR-CMP-01 [P2]:** Chiến dịch có mục đích, chủ trách nhiệm, tập nhận, phiên bản nội dung, thời gian, ngân sách và phê duyệt; không gửi thật từ bản nháp.
- **BR-CMP-02 [P2]:** Kiểm tra lại quyền nhận, loại trừ và tần suất sát thời điểm gửi. Tập nhận chốt từ trước không được bỏ qua người đã từ chối sau đó.
- **BR-CMP-03 [P2]:** Có gửi thử và dừng khẩn cấp tin chưa gửi; báo phần đã gửi/chờ/dừng. A/B Testing chỉ kích hoạt khi có phương án thử được duyệt, không là điều kiện bắt buộc của tích hợp User/Order đầu tiên.

### 7.14. OPS - Tra cứu và vận hành

- **BR-OPS-01 [P0]:** Từ tham chiếu nguồn, tra được luồng, phạm vi, tập nhận được phép xem, quy tắc/mẫu/kênh áp dụng, lần thử, kết quả, chi phí và lý do không gửi. Người dùng thấy thông báo lỗi dễ hiểu; mã lỗi kỹ thuật dùng cho hỗ trợ.
- **BR-OPS-02 [P1]:** Gửi lại cần quyền và lý do; kiểm tra lại hiệu lực, người nhận, đồng ý và ngân sách. Không gửi lại bản thành công hoặc mã/link đã hết hạn; trường hợp gửi lại có chủ đích được duyệt phải có tham chiếu riêng, vẫn nối với bản gốc.
- **BR-OPS-03 [P1]:** Theo dõi tồn đọng, lỗi, kết quả chưa rõ, tỷ lệ giao theo bằng chứng, hạn mức và nhà cung cấp kênh gửi; tạm dừng theo luồng/kênh/phạm vi có tác động rõ. Không dùng chỉ số tiếp nhận để báo tỷ lệ người nhận đã nhận.
- **BR-OPS-04 [P0]:** Phân biệt môi trường thử với gửi thật. Giả lập/ghi log được gắn nhãn riêng, không đưa vào kết quả hoặc báo cáo thành công của môi trường thật.

## 8. Kịch bản từ góc nhìn nguồn và người nhận

### 8.1. User gửi OTP

1. Người dùng yêu cầu xác thực. User tạo thử thách, chọn liên hệ, phát hành mã và thời hạn.
2. User yêu cầu Notification gửi theo tham chiếu ổn định qua đường tích hợp được bảo vệ.
3. Notification kiểm tra hiệu lực và gửi qua kênh thực được chọn; trả kết quả theo hợp đồng. Không lưu mã trong log/lịch sử.
4. Nếu gửi thất bại/chưa rõ, User quyết định hiển thị lỗi, tra cứu hoặc cho cấp lại theo chính sách. Không thông báo “đã gửi mã” chỉ vì event đã được nhận.
5. Người dùng nhập mã tại User; User xác minh. Notification không xác minh hoặc tự mở khóa tài khoản.

### 8.2. User mời thành viên vào Shop

Lời mời đã lưu mới được gửi. Người được mời nhận tên Shop/người mời và link đúng thời hạn. Gửi lỗi không làm lời mời đã lưu biến mất; người mời có cách gửi lại hoặc sao chép link theo quyền. Nếu lời mời bị thu hồi/đã dùng, Notification dừng nhắc link cũ dựa trên thông tin User cung cấp.

### 8.3. Order xác nhận giao thất bại

| Người nhận | Nội dung và kênh minh họa | Điều kiện |
|---|---|---|
| Người nhận hàng | Mã nhận biết đơn, kết quả và hướng liên hệ; ZNS chính, SMS dự phòng nếu được duyệt. | Order xác nhận liên hệ, chặng và sự việc; đủ điều kiện kênh. Không hứa “sẽ giao lại” khi Order chưa xác nhận. |
| Shop | Đơn/chặng cần theo dõi, lý do được phép công bố và đường dẫn đơn; In-app/Push nếu đã triển khai. | Thành viên/đầu mối đúng Shop được phép nhận. |
| Vận hành | Cảnh báo cần can thiệp. | Chỉ tạo khi có điều kiện và đầu mối chịu trách nhiệm, không gửi cho mọi nhân viên. |

Webhook NVC trùng chỉ tạo một lần xử lý cho cùng sự việc. Một lần giao lại thực sự thất bại có thể tạo sự việc mới. Khi Order xác nhận giao thành công/chuyển hoàn trước giờ nhắc, dừng lời nhắc không còn giá trị; không tự suy ra từ việc người dùng đọc Push.

### 8.4. Đơn mới và chuyển hoàn

- order.created.v1 xác nhận đơn SuperPlatform đã được tạo; chưa đủ để báo tạo vận đơn thành công hoặc bắt đầu giao.
- order.return.requested.v1 và order.return.confirmed.v1 lần lượt là yêu cầu/xác nhận chuyển hoàn theo nghiệp vụ nguồn; không phải hàng đã hoàn về Shop.
- Thông báo “đã hoàn hàng” cần mốc hoàn thực tế do Order/module sở hữu xác nhận, với chặng/loại hoàn rõ ràng.

### 8.5. Module tài chính gọi Notification trong tương lai

Module tài chính chốt kết quả đối soát và cung cấp người nhận có quyền, số tiền, pháp nhân và đường dẫn. Notification không tự cộng COD từ các đơn để tạo thông báo “đã thanh toán”. Người không có quyền xem số tiền chỉ nhận nội dung tối thiểu hoặc không nhận trường hợp gửi này.

## 9. Tiêu chí nghiệm thu

Các kịch bản dưới đây là đề xuất kiểm chứng nghiệp vụ. Đợt đầu phải nghiệm thu đủ 5 kênh đã xác nhận và các luồng được duyệt, cùng quy tắc nền tảng liên quan. Không tuyên bố tất cả kênh hoặc module đã sẵn sàng chỉ từ một kênh thử thành công.

| Mã | Tình huống | Kết quả cần đạt | BR chính |
|---|---|---|---|
| AC-01 | Nguồn hợp lệ gửi sự việc với dữ liệu đủ | Có tham chiếu, đúng luồng/người/phạm vi, tra được tới kết quả từng kênh. | TRG-01, INT-04, OPS-01 |
| AC-02 | Gửi lại cùng event hoặc lỗi giữa các bước | Không nhân bản tin; tiếp tục phần chưa hoàn tất, không gửi lại phần đã đạt tiêu chí. | TRG-02, TRG-05 |
| AC-03 | Nguồn trái quyền, sai phiên bản hoặc thiếu biến mẫu | Từ chối/chờ có lý do, không tự đoán hoặc báo đã gửi. | TRG-01, TRG-04, TPL-02 |
| AC-04 | Một người thuộc hai Shop, sau đó mất quyền một Shop | Không lộ chéo dữ liệu; quyền xem và nhắc tiếp được kiểm tra trong đúng ngữ cảnh. | REC-03, REC-04, SEC-03 |
| AC-05 | Người nhận hàng không có tài khoản | Gửi kênh liên hệ đủ điều kiện; không tự gắn tài khoản/hộp thư hoặc cấp quyền xem Shop. | REC-01, CHN-07 |
| AC-06 | OTP hết hạn, nhà cung cấp kênh gửi lỗi hoặc chỉ dùng adapter ghi log | Không gửi mã hết hạn; User nhận đúng lỗi/chưa rõ; log không được coi gửi thật thành công. | INT-02, TIME-01, OPS-04 |
| AC-07 | Lời mời/đơn đã đổi trạng thái trước lịch nhắc | Tin chưa gửi được dừng theo nguồn; lịch sử vẫn giải thích được vì sao. | INT-05, TIME-02 |
| AC-08 | Nhà cung cấp kênh gửi xác nhận nhận, timeout, hoặc không có xác nhận giao | Tách đúng các trạng thái; không tự đổi timeout thành thất bại rồi gửi nhiều kênh. | CHN-03, CHN-09 |
| AC-09 | Một người/kênh thành công, một người/kênh lỗi | Kết quả một phần và chi tiết đúng; chỉ xử lý phần còn hợp lệ. | CHN-04, TRG-05 |
| AC-10 | Từ chối Marketing / liên hệ bị báo sai | Chặn đúng mục đích/kênh; không gộp mọi lệnh chặn, không lách bằng nhãn giao dịch. | CLS-02, PREF-01, PREF-04 |
| AC-11 | Mẫu/cấu hình chồng phạm vi hoặc bị sửa sau khi gửi | Không áp dụng ngẫu nhiên; truy ra phiên bản thực tế, lịch sử không bị đổi. | CFG-02, CFG-03, TPL-03 |
| AC-12 | Gửi lại tin, hết ngân sách hoặc tạm dừng | Kiểm tra quyền/hiệu lực/chi phí; có lý do và audit; không báo thành công khi không gửi. | OPS-02, COST-03, SEC-04 |
| AC-13 | Đọc In-app rồi mở link khi đã mất quyền | Đọc không hoàn thành nghiệp vụ; link bị kiểm soát bởi nguồn, không cấp quyền từ thông báo. | CHN-05, TPL-05, SEC-03 |
| AC-14 | Tra log, event kết quả và màn hình hỗ trợ | Không lộ OTP/token; liên hệ/nội dung nhạy cảm chỉ hiển thị theo quyền, có chính sách lưu. | SEC-01, SEC-05, INT-04 |
| AC-15 | Nhận order.created hoặc return.confirmed | Nội dung đúng mốc đã chốt, không báo đã tạo vận đơn/đang giao/đã hoàn hàng sai sự việc. | INT-01, TPL-02 |
| AC-16 | Nhà cung cấp kênh gửi hoặc bên tiêu thụ gửi lại phản hồi/event kết quả | Không phát sinh lần gửi mới/vòng lặp; bên nguồn phân biệt đúng mức kết quả. | CHN-04, INT-04 |
| AC-17 | Kiểm chứng đủ 5 kênh đợt đầu | In-app ghi và xem được đúng hộp thông báo; Push gửi tới thiết bị/app đích; Email gửi thật tới địa chỉ hợp lệ, gồm Gmail; SMS và Zalo ZNS gửi thật theo cấu hình. Mỗi kênh có kết quả/bằng chứng đúng khả năng, không dùng giả lập làm bằng chứng hoàn tất. | CHN-05, CHN-06, CHN-07, CHN-08, CHN-09, OPS-04 |
| AC-18 | Cấu hình luồng dùng người/app/quyền đã có trong User | Không yêu cầu tạo lại tài khoản/role/app. Chỉ chọn phạm vi và app được phép, có cấu hình gửi riêng; kết quả áp dụng truy vết được. | CAT-04, CFG-05, SEC-06 |
| AC-19 | Giả mạo ID người/app/Shop, đổi ngữ cảnh thông báo hoặc xử lý event khi người nhận không đăng nhập | Hộp thông báo/Push không lộ chéo; backend kiểm tra ngữ cảnh/quyền. Event nền không phụ thuộc phiên người nhận; không dùng token người tạo thay quyền cả tập nhận. | INT-06, REC-06, SEC-06, CFG-05 |

Mục tiêu tốc độ phản hồi OTP, độ trễ tin giao dịch, sản lượng, thời gian phục hồi, thời gian lưu và hạn mức phải được định lượng ở DEC-05, DEC-06, DEC-08. Chưa có giá trị được duyệt, nên bản BRD này không tự đặt SLA hoặc tỷ lệ giao cam kết.

## 10. Giả định, phụ thuộc và giới hạn

### 10.1. Giả định cần xác nhận

- **ASM-01:** User/Order cung cấp được hợp đồng nguồn, tham chiếu ổn định và dữ liệu tối thiểu. Hiện trạng port/outbox chưa chứng minh đường tích hợp mới đã hoạt động.
- **ASM-02:** Có cơ chế nhận thay đổi tư cách/quyền, liên hệ và điều kiện dừng đủ tin cậy cho các luồng được chọn. Nếu chưa có, phải chọn cách nguồn cung cấp/kiểm tra trước gửi.
- **ASM-03:** Chủ nghiệp vụ cung cấp danh sách đầu mối và quyết định nội dung. Notification không coi “quản lý cấp trên”, “kế toán” là quy tắc chọn người nhận đã sẵn sàng.
- **ASM-04:** Định danh gửi, mẫu nhà cung cấp kênh gửi và ngân sách các kênh sẽ được phê duyệt trước gửi thật.
- **ASM-05:** Nền tảng cung cấp căn cứ phân quyền và ngữ cảnh app/brand/pháp nhân; chưa xác nhận thì không tự lấy SuperShip làm mặc định cho mọi yêu cầu.

### 10.2. Giới hạn phải giải thích cho người sử dụng

- Nhà cung cấp kênh gửi nhận gửi không bảo đảm người nhận đọc; có kênh không có xác nhận giao. Không cam kết “gửi đúng một lần tuyệt đối” khi timeout ngoài hệ thống; phải có chính sách chống trùng và xử lý chưa rõ.
- Dữ liệu nguồn đến muộn làm tăng nguy cơ tin cũ; cần điều kiện còn hiệu lực và cơ chế kiểm tra đã thống nhất. Không thể hứa thu hồi SMS/Email/Push đã đưa ra ngoài.
- Chưa có tài khoản/thiết bị thì không sử dụng In-app/Push cho người nhận đó.
- Chất lượng liên hệ, nhà cung cấp kênh gửi, quy tắc mẫu và thời gian phản hồi ảnh hưởng kết quả; Notification cần hiển thị giới hạn và nguyên nhân.
- Không tự bảo đảm dữ liệu người nhận/quyền chính xác khi User hoặc module nguồn chưa đồng bộ.

## 11. Quyết định cần chốt trước đặc tả/triển khai

DEC-01 đã **chốt 5 kênh và các lát cắt đợt đầu** ngày 08/10/2026. Cũng ngày 08/10/2026, **PO đã chốt DEC-01 đến DEC-12** cho đợt đầu tại [Quyết định PO cho các DEC](<docs/Notitek - QUYẾT ĐỊNH PO CHO CÁC DEC.md>); các điểm thuộc thẩm quyền pháp chế, tài chính, Tech Lead và vận hành được ghi là chờ xác nhận tại mục 4 của tài liệu đó. Bảng dưới giữ nguyên câu hỏi gốc để truy vết.

| Mã | Câu hỏi cần quyết định | Bên chủ trì | Phạm vi bị ảnh hưởng |
|---|---|---|---|
| DEC-01 | Đã xác nhận 5 kênh ngày 06/10/2026. Ngày 08/10/2026, người yêu cầu chốt OTP, lời mời nhân viên/thành viên, giao thất bại và cảnh báo thiết bị; hộp tin/tra cứu và các năng lực nền phục vụ các lát cắt. Shop SuperShip và giao diện nội bộ là điểm tích hợp ưu tiên. Còn cần chốt client Push cụ thể, app mở rộng và phân kỳ các năng lực ngoài đợt đầu. | Chủ sản phẩm + User/Order/Notification | Phạm vi bản phát hành tại docs/Notitek - PHẠM VI VÀ KẾ HOẠCH PHÁT HÀNH ĐẦU TIÊN.md; không coi mọi UC P0 dài hạn là cam kết đợt đầu. |
| DEC-02 | Event nguồn, thời điểm commit, mức đơn/chặng và cách phát outbox; yêu cầu gửi trực tiếp giữ hay thay bằng hợp đồng nào? | Chủ từng module nguồn + Notification | Mọi luồng tích hợp; không tự dùng tên event đề xuất như hợp đồng thật. |
| DEC-03 | Ai nhận mỗi trường hợp; hợp đồng bản đọc/nhóm nhận thông báo từ User/nguồn, độ mới, điều kiện quyền cho xử lý nền và ngoại lệ liên hệ cũ/người nghỉ việc? | Chủ nghiệp vụ + User/phân quyền | Dùng mô hình User hiện có; chốt dữ liệu tích hợp, không tạo lại hệ thống người dùng/quyền. |
| DEC-04 | Luồng nào bắt buộc, kênh nào cho phép, đồng ý Marketing lấy ở đâu và chặn liên hệ xử lý thế nào? | Chủ nghiệp vụ + bảo vệ dữ liệu | CLS/PREF, không mặc định mọi tin giao dịch bắt buộc. |
| DEC-05 | Mức kết quả gửi đủ cho User, thời gian chờ và chỉ dẫn thử lại/hủy/tra cứu khi chưa rõ là gì? User tiếp tục quyết định phát hành, kích hoạt, xác minh và cấp lại mã. | User + Notification + vận hành kênh | Hợp đồng gửi dùng chung và tích hợp xác thực, không đưa nghiệp vụ OTP vào Notification. |
| DEC-06 | Thời hạn, nhắc lại, khung giờ/múi giờ, giới hạn tần suất, số lần thử và thời gian xử lý từng luồng? | Chủ nghiệp vụ + vận hành | TIME/CHN và nghiệm thu hiệu năng. |
| DEC-07 | Thứ tự ưu tiên cấu hình giữa Shop/tổ chức/app/brand/NVC; khóa gì; thay cấu hình ảnh hưởng tin chờ thế nào? | Chủ nền tảng + quản trị Notification | CFG, TPL và phạm vi quyền. |
| DEC-08 | Ai sở hữu/kiểm soát ngân sách, Notification được giao thực thi phần nào; đơn vị chịu phí, hạn mức và chỉ dẫn khi không được gửi là gì? | Chủ ngân sách + tài chính + vận hành | COST và khả năng gửi thật; không mặc định Notification quản lý sổ ngân sách. |
| DEC-09 | Giữ loại dữ liệu nào bao lâu; ai xem/xuất, khi nào ẩn/xóa và audit giữ thế nào? | Bảo vệ dữ liệu + vận hành + chủ sản phẩm | SEC, lịch sử và hộp In-app. |
| DEC-10 | Loại bằng chứng giao từng kênh; tiêu chí kết quả tổng hợp; khi nào chuyển dự phòng với kết quả chưa rõ? | Notification + nhà cung cấp kênh gửi + chủ nghiệp vụ | CHN và event kết quả. |
| DEC-11 | Mẫu ZNS/SMS kế thừa nào đang dùng; nội dung, đầu mối, brand/pháp nhân và điều kiện đã được duyệt? | Vận hành hiện tại + chủ brand | Chuyển đổi luồng cũ sang module mới. |
| DEC-12 | Hành động/tài nguyên Notification nào cần đăng ký catalog User/Authorization; ai được cấp quyền soạn/duyệt/gửi lại/dừng; có tách người soạn và người duyệt không? | Chủ nền tảng + User/phân quyền | Dùng phân quyền chung; không xây role/quyền độc lập trong Notification. |

Khi quyết định được duyệt, ghi nội dung, người duyệt, ngày và các BR/UC bị ảnh hưởng; không xóa câu hỏi mà không lưu kết quả. Một DEC chưa chốt chỉ chặn phần phụ thuộc vào nó, không mặc định chặn toàn bộ module.

## 12. Rủi ro và kiểm soát

| Rủi ro | Hậu quả | Kiểm soát nghiệp vụ |
|---|---|---|
| Gửi sai người hoặc chéo tổ chức | Lộ thông tin, sai hành động, thiệt hại tài chính | Người nhận có căn cứ; quyền/phạm vi hiện tại; ngoại lệ bảo mật có User xác nhận. |
| Event trùng, trễ, nhiều vận đơn kỹ thuật | Spam, thông báo trạng thái sai | Chống trùng theo sự việc; mốc/chặng rõ; kiểm tra còn hiệu lực. |
| Timeout rồi chuyển kênh ngay | Người nhận nhận nhiều bản, tăng phí | Theo dõi chưa rõ, tra kết quả và chính sách dự phòng đã duyệt. |
| Ghi log/tiếp nhận được báo thành gửi thành công | OTP không đến, nguồn quyết định sai | Tách kết quả, bằng chứng gửi thật và môi trường giả lập. |
| Gắn Marketing thành giao dịch/bắt buộc | Gửi ngoài lựa chọn của người nhận | Duyệt mục đích, quản lý đồng ý/từ chối và kiểm tra trước gửi. |
| Phạm vi Phase 1 quá rộng so với phụ thuộc | Không chốt được SRS/ nghiệm thu | Chọn đợt theo UC và hợp đồng nguồn; giữ phạm vi dài hạn riêng. |
| Mẫu hoặc cấu hình nhiều chiều xung đột | Brand, kênh, nội dung áp dụng sai | Chốt ưu tiên, kiểm tra xung đột, lưu cấu hình thực tế. |
| Lưu bí mật hoặc đưa vào event kết quả | Lộ OTP/token/dữ liệu nhạy cảm | Dữ liệu tối thiểu, đường truyền được bảo vệ, che theo quyền và giới hạn lưu. |
| Hết ngân sách/không có người chịu phí | Dừng tin cần thiết hoặc phí không kiểm soát | Chốt chủ ngân sách, hạn mức và ngoại lệ trước gửi. |

## 13. Truy vết use case và tài liệu nguồn

### 13.1. Liên kết với danh sách use case

| Nhóm UC | Nhóm BR áp dụng chính |
|---|---|
| UC-NTF-01 đến UC-NTF-04 | CAT, TRG, INT, CLS |
| UC-NTF-05 đến UC-NTF-08 | REC, PREF, SEC, TPL, CHN, CFG |
| UC-NTF-09 đến UC-NTF-12 | CHN, TIME, INT, OPS |
| UC-NTF-13 đến UC-NTF-18 | CHN, REC, PREF, CAT, TPL, CFG, TIME, COST, CMP, OPS |
| UC-USR-01 đến UC-USR-09 | INT, REC, TPL, TIME, PREF, SEC, CHN; đặc biệt INT-02 cho OTP, INT-01/06 và CHN-05 cho cảnh báo thiết bị |
| UC-ORD-01 đến UC-ORD-11 | TRG, INT, REC, TPL, TIME, CHN, CFG; thêm SEC/COST và nguồn tài chính khi tích hợp |

Mã BR rút gọn trong bảng nghiệm thu/truy vết được đọc với tiền tố BR-, ví dụ INT-02 là BR-INT-02. Bảng trên là liên kết nhóm; bước phân rã từng UC phải chọn mã BR cụ thể, không mặc định mọi yêu cầu trong nhóm đều áp dụng.

### 13.2. Căn cứ đọc mã nguồn tại thời điểm rà soát

- [NotificationPort của User](../user-spf/services/users-core-service/src/main/java/com/supership/users/notification/application/NotificationPort.java): ranh giới gửi tin, OTP cần gửi thật và xử lý lỗi của bên gọi.
- [Dữ liệu yêu cầu gửi User](../user-spf/services/users-core-service/src/main/java/com/supership/users/notification/application/Notification.java) và [NotificationReceipt](../user-spf/services/users-core-service/src/main/java/com/supership/users/notification/application/NotificationReceipt.java): tham chiếu, đích, mẫu, tham số, thời hạn và biên nhận nhà cung cấp kênh gửi.
- [Template User](../user-spf/services/users-core-service/src/main/java/com/supership/users/notification/application/NotificationTemplates.java): OTP, kích hoạt, lời mời, đổi mật khẩu/định danh và nghỉ việc.
- [LoggingNotificationAdapter](../user-spf/services/users-core-service/src/main/java/com/supership/users/notification/infrastructure/LoggingNotificationAdapter.java): chưa phải adapter gửi thật, không được dùng log làm bằng chứng đã giao.
- [OutboxWriter của Order](../order-spf/src/main/java/vn/supership/superplatform/order/platform/outbox/OutboxWriter.java): ghi event cùng giao dịch, dispatcher chưa được xây dựng theo mô tả tại mã nguồn.
- [CreateOrderHandler](../order-spf/src/main/java/vn/supership/superplatform/order/feature/order/internal/application/handler/CreateOrderHandler.java): mốc order.created.v1 trước bước tạo vận đơn.
- [OrderActionPolicy](../order-spf/src/main/java/vn/supership/superplatform/order/feature/order/internal/domain/model/OrderActionPolicy.java): các trạng thái nghiệp vụ hiện có.
- [Danh mục Application/Client/Entitlement của User](../user-spf/database-docs/data/user/business/tables/application.yaml): trách nhiệm dữ liệu ứng dụng, điều kiện sử dụng và domain.
- [AccessContextEligibilityService](../user-spf/services/users-core-service/src/main/java/com/supership/users/authentication/application/AccessContextEligibilityService.java): điều kiện tư cách/app và đặc thù truy cập do User sở hữu.
- [Hợp đồng Access Context](../user-spf/contracts/access-context/README.md), [ResolveDtos](../user-spf/services/users-core-service/src/main/java/com/supership/users/accesscontext/api/ResolveDtos.java), [AccessContextResolver](../user-spf/services/users-core-service/src/main/java/com/supership/users/accesscontext/application/AccessContextResolver.java): ngữ cảnh, quyết định theo hành động, scope và giới hạn sử dụng.

Tên event, chính sách kênh cụ thể và phạm vi mở rộng chưa xác nhận vẫn được đánh dấu đề xuất. Phạm vi hỗ trợ 5 kênh và các lát cắt đợt đầu đã xác nhận tại DEC-01. Truy vết UC–BR–SRS–story–AT của đợt đầu nằm ở [backlog/nghiệm thu](<docs/Notitek - BACKLOG TRIỂN KHAI VÀ KỊCH BẢN NGHIỆM THU.md>). Khi mã nguồn hoặc hợp đồng thay đổi, cập nhật phần đối chiếu; không dùng mô tả hiện trạng thay phê duyệt nghiệp vụ.
