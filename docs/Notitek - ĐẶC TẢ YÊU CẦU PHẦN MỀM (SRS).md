# Notitek - ĐẶC TẢ YÊU CẦU PHẦN MỀM (SRS)

**Phiên bản:** 1.1 — baseline nghiệp vụ sau PO review, kèm giá trị tham số PO đã chốt.
**Ngày:** 08/10/2026.
**Chủ trì:** PO/BA Notification; phối hợp User, Order, User/Authorization, chủ ứng dụng, vận hành và QA.
**Trạng thái:** Đã PO review và chốt nội dung yêu cầu nghiệp vụ trong vai trò PO được giao. Giá trị định lượng và đầu ra tích hợp cần bên sở hữu xác nhận được quản lý riêng tại mục 11; chưa ghi nhận nghiệm thu phần mềm.

## 1. Mục đích và cách sử dụng

Notification là dịch vụ thông báo dùng chung của SuperPlatform. Một module gọi Notification cần biết phải cung cấp gì, Notification cam kết xử lý đến đâu, khi nào có thể coi việc gửi đạt yêu cầu và làm gì khi kết quả chưa rõ. Người nhận cần thấy đúng thông tin, đúng tài khoản và đúng phạm vi. Người vận hành cần giải thích được vì sao một tin được gửi, bị chặn, thất bại hoặc chưa có kết quả.

SRS này đặc tả hành vi phần mềm quan sát và kiểm chứng được. Căn cứ nghiệp vụ là [BRD](<../Notitek - ĐẶC TẢ YÊU CẦU NGHIỆP VỤ (BRD).md>), [danh mục use case](<../Notitek - DANH MỤC USE CASE.md>), [phạm vi phát hành](<Notitek - PHẠM VI VÀ KẾ HOẠCH PHÁT HÀNH ĐẦU TIÊN.md>) và [UC ưu tiên](<Notitek - ĐẶC TẢ USE CASE ƯU TIÊN.md>).

[Hợp đồng tích hợp](<Notitek - HỢP ĐỒNG TÍCH HỢP API VÀ SỰ KIỆN.md>) phải cụ thể hóa các yêu cầu thành schema, lỗi, cơ chế xác thực và ví dụ. [Thiết kế kỹ thuật](<Notitek - THIẾT KẾ KỸ THUẬT SƠ BỘ.md>) quyết định cách triển khai. SRS không chọn sẵn broker, cơ sở dữ liệu, engine, nhà cung cấp hoặc tên endpoint.

### 1.1. Quy ước yêu cầu

- “Phải” là hành vi bắt buộc trong phạm vi được ghi; “không được” là điều hệ thống phải ngăn.
- SRS-F01 đến SRS-F14 giữ ý nghĩa nhóm của bản 0.1. Mã con, ví dụ SRS-F01.01, là yêu cầu cụ thể để phát triển và QA truy vết.
- SRS-F15 đến SRS-F22 bổ sung các năng lực còn thiếu. Tất cả các nhóm chức năng dưới đây thuộc đợt đầu ở mức tối thiểu được mô tả.
- SRS-N01 đến SRS-N07 giữ mã yêu cầu chất lượng. Các tham số chưa có giá trị được quản lý ở mục 11; không được coi là đã nghiệm thu.
- Mã BR/UC/AT liên kết đến tài liệu nguồn. Một nhóm SRS có nhiều yêu cầu con; đạt một ca AT chưa chứng minh đạt cả nhóm.
- Tên dữ liệu và trạng thái trong SRS diễn đạt ý nghĩa nghiệp vụ. Tên trường, enum và mã HTTP chính thức thuộc hợp đồng tích hợp.

### 1.2. Lịch sử phiên bản

| Phiên bản | Nội dung |
|---|---|
| 0.1 | Tóm tắt 14 nhóm chức năng và 7 nhóm phi chức năng. |
| 0.2 | Viết lại từ nhu cầu của nguồn, người nhận và vận hành; phân rã hành vi, bổ sung giao diện ngoài, năm kênh, kết quả và tiêu chí kiểm chứng. |
| 1.0 | Hoàn thành vòng PO review ngày 08/10/2026; sửa cửa sổ đồng bộ/bảo vệ dữ liệu chống trùng, đồng bộ UC–UX–hợp đồng–thiết kế–backlog và chốt baseline nghiệp vụ. |
| 1.1 | Ngày 08/10/2026, PO chốt DEC-01 đến DEC-12 và giá trị khởi điểm SRS-P01 đến SRS-P12; mục 11 dẫn chiếu tài liệu Quyết định PO cho các DEC. Không đổi yêu cầu chức năng. |

### 1.3. Các từ dùng trong tích hợp

| Từ | Cách hiểu trong tài liệu |
|---|---|
| Reference | Mã nguồn dùng để nhận biết một yêu cầu/sự việc và tra lại khi phát lại. |
| Provider | Nhà cung cấp kênh gửi; nhận yêu cầu chưa luôn chứng minh người nhận đã nhận. |
| Payload/schema | Dữ liệu trao đổi và quy tắc về trường, kiểu, điều kiện hợp lệ. |
| Retry | Thử lại phần gửi còn được phép sau lỗi; không phải nguồn phát hành nghiệp vụ mới. |
| Fallback | Chuyển sang kênh dự phòng khi có điều kiện đã duyệt. |
| Callback | Phản hồi kênh gửi về sau để cập nhật bằng chứng. |
| Scope | Phạm vi người, app, Shop/tổ chức và tài nguyên được phép xử lý/xem. |
| Idempotent | Gọi lại cùng thao tác hợp lệ không nhân công việc hay thay kết quả theo cách sai. |
| Outbox/dead-letter | Nơi nguồn giữ event cần bàn giao/nơi giữ công việc lỗi; có bản ghi chưa có nghĩa người nhận đã nhận tin. |
| Mock/shadow | Giả lập/chạy đối chiếu; không được tính là gửi thật. |

## 2. Phạm vi và ranh giới của hệ thống

### 2.1. Nhu cầu thực tế mà đợt đầu phải đáp ứng

| Bên có yêu cầu | Yêu cầu đối với Notification | Kết quả bên đó cần |
|---|---|---|
| User — xác thực | “Tôi đã phát hành một mã và chọn liên hệ. Hãy gửi qua SMS thật, trong hạn gửi tôi cấp, và trả đúng mức kết quả để tôi quyết định tiếp tục.” | Phân biệt đã tiếp nhận, nhà cung cấp nhận, lỗi và chưa rõ; không tự kích hoạt hoặc xác minh mã. |
| User — lời mời | “Lời mời đã được lưu. Hãy gửi email đúng tổ chức/Shop và link tôi cấp; nếu gửi lỗi, lời mời vẫn tồn tại.” | Truy vết được email; không tạo membership hoặc đảo ngược lời mời đã commit. |
| Order | “Tôi xác nhận một lần giao thất bại. Hãy báo đúng Shop và liên hệ của đơn, với nội dung và kênh khác nhau; dừng phần chờ khi tôi yêu cầu.” | Kết quả riêng theo trường hợp gửi; phân biệt lần phát lại với sự việc giao thất bại mới. |
| Người dùng Shop | “Tôi muốn xem tin của tài khoản và công việc được phép, biết tin nào chưa đọc và mở đúng đơn.” | Không lẫn Shop, app hoặc tài khoản; đọc tin không hoàn tất công việc. |
| Chủ tài khoản | “Tôi cần tiếp tục thấy cảnh báo thiết bị và xử lý tại User, kể cả khi chuyển Shop.” | Không mất hành động bảo mật đang có; trạng thái nguồn tách với trạng thái đọc. |
| Vận hành | “Tôi cần tìm từ reference nguồn, biết đã gửi cho ai, bằng cấu hình nào và có thể dừng phần nào.” | Lịch sử có bằng chứng và lý do, che dữ liệu theo quyền; không cần xem OTP/token. |
| Chủ ngân sách | “Kênh có phí chỉ được gửi theo điều kiện tôi đã cho phép.” | Không tự vượt hạn mức, không báo thành công khi bị chặn, không tự lập sổ tài chính. |

### 2.2. Nội dung thuộc đợt đầu

- Tiếp nhận lệnh gửi và sự việc từ nguồn được phép; chống trùng; chọn luồng, mẫu và người nhận theo hợp đồng.
- SMS cho OTP; Email cho lời mời; In-app/Push cho Shop; Zalo ZNS và SMS dự phòng theo chính sách cho liên hệ của đơn.
- Hộp tin tiếp nối cảnh báo thiết bị User, tách ngữ cảnh tài khoản/công việc.
- Cấu hình tối thiểu, kết nối năm kênh, kết quả, tra cứu, quyền, bảo vệ dữ liệu, điều kiện chi phí và chuyển đổi luồng cũ được chọn.
- Đường tích hợp và việc phối hợp User, Order, frontend, chủ đích Push, vận hành là phụ thuộc có người chịu trách nhiệm.

Đợt đầu phải kiểm chứng đủ năm kênh trên client và môi trường đã chọn. Không yêu cầu mỗi luồng phải dùng cả năm kênh.

### 2.3. Ngoài phạm vi cam kết của đợt đầu

Trình thiết kế workflow kéo thả, trang quản trị tự phục vụ đầy đủ, quản lý lựa chọn nhận tin tự phục vụ, gom nhóm/nhắc định kỳ nâng cao, báo cáo chi phí nâng cao, chiến dịch Marketing và A/B testing đi theo phân kỳ riêng. Tài chính, Support, Pricing và các app khác là hướng mở rộng, chưa có cam kết tích hợp trong release này.

### 2.4. Ai sở hữu quyết định nào?

| Nội dung | Chủ sở hữu | Trách nhiệm của Notification |
|---|---|---|
| Mã/link, hiệu lực, cấp lại, xác minh | User | Gửi theo yêu cầu và hạn gửi; dừng phần chờ theo chỉ dẫn; trả bằng chứng gửi. |
| Đơn, chặng, NVC, giao đủ/giao một phần, việc còn cần nhắc | Order hoặc module nguồn được giao | Dùng sự việc và chỉ dẫn nguồn; không suy kết quả giao từ tracking thô. |
| Identity, membership, app/client, quyền | User/Authorization và module sở hữu tài nguyên | Dùng định danh và quyết định được xác nhận; không tạo bản quản lý quyền riêng. |
| Đăng ký, đổi chủ, thu hồi đích Push | Bên sở hữu thiết bị/app được chọn | Nhận/tra đích có hiệu lực; ngừng sử dụng đích bị thu hồi; báo phản hồi kênh. |
| Ngân sách và quyền chi | Bên được chốt tại DEC-08 | Thực thi phần trách nhiệm được giao; ghi lý do khi bị chặn. |
| Nội dung, điều phối, In-app và kết quả gửi | Notification, theo chính sách chủ nghiệp vụ duyệt | Chịu trách nhiệm hành vi được đặc tả trong SRS này. |

Notification không cần biết session OTP hoặc toàn bộ trạng thái đơn để thực hiện chức năng. Notification cần hợp đồng chung đủ dữ liệu và chỉ dẫn từ bên sở hữu.

## 3. Căn cứ hiện trạng đã đối chiếu

Đối chiếu mã nguồn ngày 08/10/2026. Đây là căn cứ xác định phụ thuộc, không phải bằng chứng hệ thống đã chạy tích hợp.

| Căn cứ | Điều đã thấy | Hệ quả đối với SRS |
|---|---|---|
| [NotificationPort](../../user-spf/services/users-core-service/src/main/java/com/supership/users/notification/application/NotificationPort.java), [Notification](../../user-spf/services/users-core-service/src/main/java/com/supership/users/notification/application/Notification.java), [Receipt](../../user-spf/services/users-core-service/src/main/java/com/supership/users/notification/application/NotificationReceipt.java) | User truyền reference, channel, destination, template, parameters, expiresAt; receipt hiện chỉ có providerReference. | Adapter cần ánh xạ hợp đồng mới; một chuỗi reference không tự chứng minh gửi thật. |
| [OtpDeliveryPolicy](../../user-spf/services/users-core-service/src/main/java/com/supership/users/authentication/application/OtpDeliveryPolicy.java), [OtpSessionService](../../user-spf/services/users-core-service/src/main/java/com/supership/users/authentication/application/OtpSessionService.java) | Mỗi lần gửi có deliveryId riêng; User gửi trước khi cập nhật mã/counters; timeout không tự retry, mã trước có thể còn dùng được. | Không kích hoạt retry nền OTP ngoài chỉ dẫn; kết quả đến muộn không làm Notification quyết định mã có hiệu lực. |
| [InvitationLinkNotifier](../../user-spf/services/users-core-service/src/main/java/com/supership/users/invitation/InvitationLinkNotifier.java), [NotifyAfterCommit](../../user-spf/services/users-core-service/src/main/java/com/supership/users/notification/application/NotifyAfterCommit.java) | Tin sau sự việc được gửi sau commit; lỗi gửi không đảo thay đổi nguồn. | Nguồn cần cơ chế bàn giao bền vững nếu phải phát lại sau lỗi; callback trong bộ nhớ chưa đủ chứng minh không mất tin. |
| [OutboxWriter](../../order-spf/src/main/java/vn/supership/superplatform/order/platform/outbox/OutboxWriter.java) | Order ghi event cùng giao dịch; chú thích ghi dispatcher chưa được xây dựng. | Có outbox chưa có nghĩa Notification đã nhận; hợp đồng giao thất bại và đường phát là công việc của kế hoạch tích hợp. |
| [Chuông Shop](../../shop-fe/apps/business-web/features/notification-center/use-notification-center.ts) | Feed và hành động device trust thực từ User; đọc có phần lưu theo browser session. | Cần đối chiếu ID khi chuyển nguồn; trạng thái đọc mới phải có nguồn sự thật rõ ràng. |
| [Thông báo lời mời nội bộ](../../internal-fe/features/employees/queries/use-employee-invitation-queries.ts) | Giao diện hiện thông báo hệ thống chưa gửi email. | Hoàn thành tích hợp phải sửa thông điệp theo kết quả thật, không giữ lời hứa sai. |
| [ResolveDtos](../../user-spf/services/users-core-service/src/main/java/com/supership/users/accesscontext/api/ResolveDtos.java) | Quyết định ALLOW/DENY/NOT_EVALUATED theo hành động/tài nguyên; permissions là gợi ý UI; scope có thể là bộ lọc rộng. | Backend cần quyết định phù hợp; bộ lọc hoặc danh sách permissions không thay kiểm tra quyền trên tài nguyên. |

Mốc HEAD tham khảo: User 3363aa17; Order 3eb3111f; Shop FE ab6198b; Internal FE 701fe56. Tệp local có thể tiếp tục thay đổi; các hành vi cần được consumer xác nhận khi chốt hợp đồng.

## 4. Mô hình nghiệp vụ và dữ liệu tối thiểu

### 4.1. Các đối tượng cần phân biệt

| Đối tượng | Ý nghĩa | Ví dụ |
|---|---|---|
| Yêu cầu/sự việc nguồn | Một lần nguồn đề nghị thông báo | Một lần phát hành mã; một lần giao thất bại đã chốt. |
| Luồng thông báo | Mục tiêu và chính sách xử lý | Báo giao thất bại. |
| Trường hợp gửi | Nhóm người nhận, nội dung và tiêu chí riêng trong luồng | Shop xem In-app/Push; người nhận hàng nhận ZNS/SMS. |
| Bản thông báo | Nội dung cho một người trong một trường hợp gửi | Tin của thành viên Shop A. |
| Đích kênh | Địa chỉ/token hoặc hộp tin dùng để đưa tin | Email lời mời; một đích Push của app đã chọn. |
| Lần thử gửi | Một lần giao yêu cầu cho nhà cung cấp | Lần gọi SMS thứ nhất. |
| Bằng chứng kênh | Phản hồi có thể kiểm chứng | Provider nhận; xác nhận giao; bounce; kết quả tra cứu. |
| Trạng thái đọc | Trạng thái In-app của người nhận | Đã đánh dấu đọc trên hộp tin. |
| Trạng thái nghiệp vụ nguồn | Kết quả ở User/Order | Lời mời đã chấp nhận; yêu cầu thiết bị đã từ chối. |

“Giao tin thành công” là kết quả kênh Notification. “Giao hàng thành công” là nghiệp vụ Order. Hai kết quả này không dùng chung trạng thái.

### 4.2. Dữ liệu nguồn cần cung cấp

| Nhóm dữ liệu | Khi nào bắt buộc | Quy tắc |
|---|---|---|
| Nguồn và môi trường | Mọi lệnh/event | Ràng buộc với danh tính tích hợp đã xác thực; nhãn tự khai không cấp quyền. |
| Reference ổn định, loại và phiên bản hợp đồng | Mọi lệnh/event | Cùng lần gửi lại dùng cùng reference; sự việc mới dùng reference mới theo hợp đồng. |
| Luồng hoặc mẫu được phép | Theo cách tích hợp | Lệnh trực tiếp phải xác định được mục tiêu; event được ánh xạ vào luồng đã duyệt. |
| Thời điểm sự việc, phiên bản nguồn nếu có | Event; các luồng cần kiểm tra độ mới | Không dùng thời gian nhận làm thời gian sự việc. Thiếu version không được tự tạo version nghiệp vụ. |
| Loại ngữ cảnh: tài khoản hoặc công việc | Mọi bản tin | Ngữ cảnh tài khoản không tự nhận Shop đang mở làm scope. |
| App/client đích và brand cần dùng | Theo luồng/kênh | Tham chiếu danh mục đã xác nhận; app quản trị và app nhận tin có thể khác nhau. |
| Shop/tổ chức, phạm vi đối tượng | Tin công việc | Theo nguồn và căn cứ được phép; OTP trước đăng nhập không bắt buộc membership Shop. |
| Người nhận hoặc quy tắc chọn đã duyệt | Luồng có gửi | Tài khoản ổn định, liên hệ giao dịch hoặc tham chiếu bộ chọn; không suy identity từ số/email. |
| Liên hệ/đích kênh | Khi kênh ngoài cần dùng | Đích do nguồn/User/chủ thiết bị xác nhận, đủ điều kiện của kênh. |
| Biến mẫu | Theo schema mẫu | Chỉ dữ liệu cần thiết; phân loại biến bí mật, dữ liệu cá nhân và dữ liệu được hiển thị. |
| Hạn gửi, thời điểm bắt đầu nếu hẹn | OTP/link và luồng có thời hạn | Hạn gửi là giới hạn bắt đầu đưa tin ra kênh; không phải cam kết người nhận sẽ nhận trước hạn. |
| Chỉ dẫn thử lại/chuyển kênh | Khi áp dụng | Theo hợp đồng và cấu hình; bí mật không tự được phép retry nền. |
| Link hoặc hành động nguồn | Khi có hành động | Nguồn cấp đích hợp lệ; nguồn kiểm tra quyền/token khi mở. |
| Điều kiện chi phí | Kênh phát sinh phí | Có đơn vị chịu phí và căn cứ cho phép theo DEC-08. |
| Correlation và tham chiếu đối tượng | Theo hợp đồng tra cứu | Nối luồng gọi và hỗ trợ tìm kiếm; không thay khóa chống trùng. |
| Chỉ dẫn hủy/nhóm liên quan | Khi nguồn yêu cầu dừng | Chỉ rõ tập mục tiêu và scope; không hủy mọi tin của một đơn vì một lệnh mơ hồ. |

Không bắt mọi đầu vào dùng cùng một schema. Event, lệnh gửi, lệnh hủy, phản hồi kênh và thao tác hộp tin có bộ điều kiện riêng.

## 5. Giao diện ngoài và trách nhiệm phối hợp

| Giao diện nghiệp vụ | Bên gọi/cung cấp | Đầu vào chính | Đầu ra quan sát được |
|---|---|---|---|
| Nhận lệnh gửi có phản hồi | User/nguồn được phép | Reference, mục tiêu, người/đích, biến, hạn gửi và chính sách | Kết quả tiếp nhận; mức gửi đã có bằng chứng hoặc chưa rõ; định danh tra cứu. |
| Nhận sự việc đã chốt | Order/nguồn | Loại/version, event reference, ngữ cảnh, dữ liệu tối thiểu | Xác nhận bàn giao bền vững hoặc lý do không nhận; kết quả xử lý sau đó. |
| Tra cứu kết quả | Nguồn/vận hành | Reference hoặc ID Notification, scope được phép | Tiến trình, kết quả riêng và tổng hợp; bằng chứng và lý do đã che theo quyền. |
| Hủy phần chờ | Nguồn/vận hành có quyền | ID/reference, phạm vi/tập mục tiêu, lý do | Phần đã dừng, phần đã ra ngoài, phần chưa rõ; trạng thái lệnh hủy. |
| Phát kết quả | Notification tới consumer | Thay đổi kết quả có định danh và revision | Consumer nhận/phát lại theo hợp đồng; lỗi phát kết quả không gửi lại tin người dùng. |
| Danh sách/chi tiết/chưa đọc/đánh dấu đọc | App với backend xác thực | Người, ngữ cảnh, cursor hoặc ID tin | Chỉ dữ liệu được phép; kết quả đọc độc lập với nguồn. |
| Tra/cập nhật hiệu lực đích Push | Chủ app/thiết bị | Đích, chủ tài khoản, app/client, phiên bản/hiệu lực | Đích có thể dùng hoặc lý do loại; phản hồi token lỗi tới bên sở hữu. |
| Nhận callback/tra kết quả kênh | Nhà cung cấp qua đường tin cậy | Provider reference, trạng thái, bằng chứng | Cập nhật đúng lần thử; không tạo lần gửi mới. |
| Quản lý cấu hình/kết nối | Vận hành được cấp quyền | Luồng, mẫu, kênh, scope, phiên bản, phê duyệt | Kiểm tra hợp lệ, kích hoạt/tạm dừng có audit; bí mật không hiển thị. |

Hợp đồng phải liệt kê lỗi có thể sửa, lỗi có thể retry và lỗi chưa rõ. Phải có quy tắc tương thích phiên bản, giới hạn dữ liệu, timeout, xác thực và ví dụ cho từng giao diện được triển khai. Không yêu cầu User đăng nhập thay người nhận để chạy một event nền.

## 6. Yêu cầu chức năng chi tiết

### 6.1. SRS-F01 — Tiếp nhận đúng nguồn và trả đúng cam kết

**Nhu cầu:** “Là module gọi, tôi cần biết yêu cầu có được nhận thật hay không, và sửa gì nếu bị từ chối.”
**Truy vết:** UC-NTF-01/02; BR-TRG-01/04, BR-INT-02/04/06, BR-SEC-06; AT-01/02/21/22.

| Mã | Yêu cầu bắt buộc |
|---|---|
| SRS-F01.01 | Hệ thống phải xác thực nguồn và quyền kích hoạt luồng trong scope yêu cầu trước khi xử lý gửi. Nguồn không được đổi scope hoặc nguồn gửi bằng cách sửa payload. |
| SRS-F01.02 | Hệ thống phải kiểm tra loại/version, reference, trường bắt buộc theo loại đầu vào, kiểu dữ liệu và giới hạn đã công bố; trả lý do cụ thể khi không hợp lệ. |
| SRS-F01.03 | Hệ thống chỉ được trả “đã tiếp nhận” khi đã ghi nhận yêu cầu và thông tin đủ để xử lý/tra cứu bền vững theo thiết kế được duyệt. Nếu ghi nhận không thành công, không được trả accepted. |
| SRS-F01.04 | Hệ thống phải trả ID Notification hoặc reference tra cứu, correlation, thời điểm và mức kết quả; việc mất phản hồi phải có cách tra cứu theo reference nguồn. |
| SRS-F01.05 | Với lệnh cần phản hồi gửi, hệ thống phải trả mức bằng chứng đã đạt trong thời gian chờ được thỏa thuận. Hết thời gian chờ không được biến thành “đã giao” hoặc mặc định thất bại. |
| SRS-F01.06 | Event thông báo về sự việc phải theo hợp đồng nguồn đã commit. Lệnh gửi OTP có thể được User gọi trước khi kích hoạt mã; Notification không được áp quy tắc “mọi tin phải đợi commit nguồn” cho OTP. |
| SRS-F01.07 | Đầu vào thiếu dữ liệu gửi phải bị từ chối hoặc chuyển chờ bổ sung có thời hạn, theo hợp đồng. Hệ thống không được tự lấy Shop mặc định, người nhận gần giống hoặc nội dung giả để tiếp tục. |

**Kiểm chứng:** Nguồn trái quyền không có bản gửi; dữ liệu sai có lý do; sự cố trước ghi nhận không trả accepted; mất phản hồi sau ghi nhận tra lại được cùng yêu cầu.

### 6.2. SRS-F02 — Chống trùng và tiếp tục sau lỗi

**Nhu cầu:** “Tôi phát lại vì không biết lần trước đã nhận chưa; đừng gửi khách thêm một tin giống hệt.”
**Truy vết:** UC-NTF-03; BR-TRG-02/03/05; AT-02/18/23.

| Mã | Yêu cầu bắt buộc |
|---|---|
| SRS-F02.01 | Hệ thống phải xác định khóa chống trùng từ nguồn, môi trường và reference theo hợp đồng; không dùng correlation hoặc số điện thoại làm khóa duy nhất. |
| SRS-F02.02 | Cùng khóa và cùng dữ liệu có ý nghĩa nghiệp vụ phải trả cùng yêu cầu, giữ tiến trình đã có, kể cả khi hai bản đến đồng thời. |
| SRS-F02.03 | Cùng khóa nhưng khác người nhận, scope, mẫu, bí mật hoặc dữ liệu gửi phải trả xung đột; không âm thầm ghi đè. Các trường được phép khác như metadata vận chuyển phải được hợp đồng liệt kê. |
| SRS-F02.04 | Khi phục hồi, hệ thống phải tiếp tục phần còn hợp lệ; không gửi lại người/kênh đã đạt tiêu chí chỉ vì phần khác lỗi hoặc consumer chưa nhận event kết quả. |
| SRS-F02.05 | Reference cho sự việc mới do nguồn cấp phải được xử lý độc lập, kể cả cùng đơn, người nhận hoặc loại sự việc. Không loại sự việc độc lập chỉ vì version thấp hơn một event khác. |
| SRS-F02.06 | Hệ thống phải giữ thông tin chống trùng trong cửa sổ được duyệt. Hợp đồng phải quy định phát lại ngoài cửa sổ; không được hứa chống trùng vô thời hạn khi dữ liệu đã xóa. |
| SRS-F02.07 | Cơ chế đối chiếu dữ liệu có bí mật phải bảo vệ thông tin đối chiếu và tuân thủ thời gian lưu; không lưu mã/token nguyên văn vào bảng chống trùng hoặc dùng biểu diễn dễ khôi phục bí mật chỉ để so payload. |

**Kiểm chứng:** Hai lệnh song song cùng khóa chỉ tạo một tập bản tin; đổi payload trả xung đột; In-app đã tạo không tạo lại khi Push lỗi; lần giao thất bại mới vẫn được nhận.

### 6.3. SRS-F03 — Chọn luồng, trường hợp gửi và cấu hình đúng phạm vi

**Nhu cầu:** “Một sự việc có thể cần báo Shop và khách khác nhau; hãy dùng đúng chính sách đã duyệt.”
**Truy vết:** UC-NTF-04/08; BR-CAT-01/02/04, BR-CLS-01/02/03, BR-CFG-01/02/03/05; AT-18/24.

| Mã | Yêu cầu bắt buộc |
|---|---|
| SRS-F03.01 | Hệ thống phải ánh xạ loại/version đầu vào vào luồng hợp lệ trong scope được phép; không dùng luồng gần giống khi không có ánh xạ. |
| SRS-F03.02 | Event hợp lệ không có luồng phải có kết quả “không có luồng áp dụng”, không tạo bản gửi. Lệnh chỉ định luồng không tồn tại/không được phép phải bị từ chối theo hợp đồng. |
| SRS-F03.03 | Hệ thống phải tách trường hợp gửi khi khác tập nhận, nội dung, kênh hoặc thời điểm; mỗi trường hợp có tiêu chí đạt riêng và thuộc cùng yêu cầu nguồn. |
| SRS-F03.04 | Hệ thống phải chọn cấu hình bằng thứ tự ưu tiên phạm vi đã duyệt; giữ phần bị khóa. Nếu các chiều Shop, app, brand, NVC xung đột chưa được giải quyết thì chặn kích hoạt hoặc chờ có lý do. |
| SRS-F03.05 | Hệ thống phải ghi luồng, trường hợp gửi và phiên bản cấu hình thực tế dùng cho từng bản; thay cấu hình không sửa lịch sử. |
| SRS-F03.06 | Mục đích, tính bắt buộc và mức khẩn phải là thuộc tính riêng được duyệt; không tự coi mọi tin Order hoặc “hệ thống” là bắt buộc, không trộn quảng cáo vào OTP. |

**Kiểm chứng:** Một event tạo hai trường hợp đúng phạm vi; luồng không có ánh xạ không gửi; xung đột cấu hình không chọn ngẫu nhiên; lịch sử vẫn chỉ ra phiên bản cũ.

### 6.4. SRS-F04 — Tạo nội dung dễ hiểu từ mẫu được duyệt

**Nhu cầu:** “Người nhận phải biết chuyện gì xảy ra, liên quan đối tượng nào và cần làm gì tiếp.”
**Truy vết:** UC-NTF-07; BR-TPL-01/02/03/04/05; AT-03/18/24.

| Mã | Yêu cầu bắt buộc |
|---|---|
| SRS-F04.01 | Tin thật phải dùng mẫu được duyệt đúng luồng, người nhận, ngôn ngữ, app/brand, kênh và định danh gửi; không nhận nội dung tự do để vượt duyệt. |
| SRS-F04.02 | Hệ thống phải kiểm tra biến bắt buộc và điều kiện của mẫu; thiếu dữ liệu thì chờ/dừng có lý do. Không gửi tên Shop, mã đơn hoặc giá trị quan trọng trống. |
| SRS-F04.03 | Nội dung phải phản ánh dữ liệu nguồn đã cấp; hệ thống không được tự tính lại tiền, kết luận giao đủ hoặc hứa giao lại khi nguồn chưa cấp căn cứ. |
| SRS-F04.04 | Hệ thống phải dùng link/app đích đã được xác nhận và kiểm tra chính sách đích; không tự ghép link mở rộng quyền. Nguồn vẫn xác thực khi người nhận mở. |
| SRS-F04.05 | Biến dữ liệu phải được xử lý theo ngữ cảnh hiển thị để không chèn mã hoặc làm thay cấu trúc nội dung; preview/Push không lộ bí mật hay dữ liệu ngoài chính sách. |
| SRS-F04.06 | Hệ thống phải lưu phiên bản và bằng chứng nội dung ở mức được phép. Thử mẫu phải dùng dữ liệu giả/đích thử được phép, có nhãn thử và không nhắm tập khách thật. |

**Kiểm chứng:** Thiếu mã đơn không gửi mẫu lỗi; hai lời mời khác Shop giữ đúng tên/link; biến có ký tự đặc biệt không phá email; hỗ trợ không đọc được token trong bằng chứng nội dung.

### 6.5. SRS-F05 — Đúng người nhận, đúng tài khoản và ngữ cảnh

**Nhu cầu:** “Tôi thuộc nhiều Shop; hãy gửi và cho tôi xem đúng việc được phép, không ghép người theo số điện thoại.”
**Truy vết:** UC-NTF-05/06/13; BR-REC-01/02/03/04/05/06, BR-INT-06, BR-SEC-03/06; AT-08/11/21/26.

| Mã | Yêu cầu bắt buộc |
|---|---|
| SRS-F05.01 | Hệ thống phải dùng người nhận do nguồn cấp hoặc quy tắc chọn người nhận đã duyệt dựa trên hợp đồng User/nguồn; không duyệt cây tổ chức để đoán quản lý/kế toán. |
| SRS-F05.02 | Hệ thống phải hỗ trợ tài khoản định danh và liên hệ giao dịch không có tài khoản. Liên hệ không có tài khoản không được tự tạo identity, membership hoặc hộp In-app. |
| SRS-F05.03 | Hệ thống không được hợp nhất hai người/lời mời/ngữ cảnh vì trùng số điện thoại hoặc email. Khử trùng người nhận trong cùng trường hợp phải theo định danh và chính sách đã duyệt. |
| SRS-F05.04 | Mỗi bản tin phải có loại ngữ cảnh và scope nguồn rõ ràng. Tin công việc gắn đúng Shop/tổ chức/app; cảnh báo tài khoản gắn đúng chủ tài khoản và app được phép. |
| SRS-F05.05 | Xử lý nền phải dùng căn cứ nguồn/dịch vụ được phép, không cần phiên đăng nhập người nhận và không dùng token người tạo thay quyền tập nhận. |
| SRS-F05.06 | Khi gửi chậm/thử lại, hệ thống phải kiểm tra căn cứ người nhận còn phù hợp theo hợp đồng. Không xác nhận được dữ liệu cần thiết thì chờ/dừng theo chính sách, không tiếp tục từ bản đọc hết hiệu lực. |
| SRS-F05.07 | Ngoại lệ gửi tới liên hệ cũ hoặc người đã mất membership phải dựa trên chỉ dẫn User và mục đích bảo mật được duyệt; không áp ngoại lệ đó cho mọi tin công việc. |
| SRS-F05.08 | Khi xem tin, hệ thống phải kiểm tra quyền hiện tại; việc từng được nhận tin không cấp quyền vĩnh viễn. Link nguồn phải tiếp tục qua quyền/token của nguồn. |

**Kiểm chứng:** Cùng email được mời vào hai Shop không bị gộp; khách không có User vẫn nhận ZNS; mất quyền Shop không xem lịch sử công việc đó; cảnh báo liên hệ cũ chỉ theo chỉ dẫn User.

### 6.6. SRS-F06 — Hạn gửi và hủy phần còn kiểm soát được

**Nhu cầu:** “Hãy dừng bản đang chờ khi tôi bảo dừng; cho tôi biết phần nào đã ra ngoài.”
**Truy vết:** UC-NTF-06/10; BR-TIME-01/02, BR-INT-05; AT-04/06/10/23.

| Mã | Yêu cầu bắt buộc |
|---|---|
| SRS-F06.01 | Hệ thống phải kiểm tra hạn gửi và điều kiện bắt đầu ngay trước từng lần đưa tin ra kênh, kể cả lần retry/fallback; hết hạn thì không bắt đầu lần gửi mới. Thời điểm phải có múi giờ/quy ước thống nhất theo hợp đồng để so sánh đúng hạn. |
| SRS-F06.02 | Hệ thống phải thực hiện chỉ dẫn hủy/điều kiện dừng được nguồn cấp hoặc xác nhận qua hợp đồng. Không tự đọc session OTP hay enum Order để suy điều kiện. |
| SRS-F06.03 | Lệnh hủy phải xác thực quyền, scope, tập bản/reference mục tiêu và lý do. Gửi lại cùng lệnh hủy không mở lại yêu cầu hoặc hủy rộng hơn. |
| SRS-F06.04 | Hệ thống phải trả riêng phần đã dừng, phần đã giao cho kênh, phần đang có kết quả chưa rõ và mục tiêu không tìm thấy; không hứa thu hồi tin ngoài khả năng của kênh. |
| SRS-F06.05 | Khi gửi và hủy xảy ra gần nhau, hệ thống phải ghi thời điểm và kết quả từng phần. Phần đã kiểm soát dừng không có lần gửi mới; phần đã ra ngoài giữ bằng chứng riêng. |
| SRS-F06.06 | Hủy có thể đến trước event mục tiêu. Nếu hợp đồng hỗ trợ trường hợp này, hệ thống phải giữ dấu dừng đúng reference/version trong cửa sổ được duyệt để event đến muộn không khởi tạo lại phần đã hủy; nếu chưa hỗ trợ phải trả rõ, không báo đã hủy thành công. |
| SRS-F06.07 | Hạn gửi/hủy không được xóa bằng chứng gửi đã có. Callback trễ được ghi cho lần thử cũ, không kích hoạt bản mới và không đổi quyết định nghiệp vụ của nguồn. |

**Kiểm chứng:** Hết hạn khi đang chờ không retry; hủy một lần giao không hủy mọi cảnh báo của đơn; callback sau hủy không gây tin mới; nguồn thấy rõ giới hạn thu hồi.

### 6.7. SRS-F07 — Điều phối kênh, retry và fallback có kiểm soát

**Nhu cầu:** “Một kênh lỗi không làm cả luồng lặp lại; timeout phải được giải thích đúng.”
**Truy vết:** UC-NTF-08/09; BR-CHN-01/02/03/04, BR-TRG-05; AT-04/09/11/23/25.

| Mã | Yêu cầu bắt buộc |
|---|---|
| SRS-F07.01 | Hệ thống phải thực hiện kênh chính, song song, dự phòng và bổ sung theo cấu hình; mỗi kênh có mục tiêu, điều kiện và mức bằng chứng cần đạt. |
| SRS-F07.02 | Hệ thống phải phân loại lỗi tạm thời, lỗi không thể retry và kết quả chưa rõ; ghi mã lý do, lần thử, thời điểm và phần có thể xử lý tiếp. |
| SRS-F07.03 | Retry phải có số lần tối đa, thời gian chờ và hạn xử lý được duyệt; không retry vô hạn với liên hệ sai, mẫu sai hoặc không được phép. |
| SRS-F07.04 | Timeout sau khi có khả năng provider đã nhận phải giữ “chưa rõ”, ưu tiên đối chiếu/tra cứu theo khả năng kênh. Không tự retry hoặc fallback chỉ vì thiếu phản hồi. |
| SRS-F07.05 | Fallback chỉ được chạy theo điều kiện đã duyệt, kiểm tra lại đích, hạn gửi và ngân sách. Nếu chính sách cho phép chuyển khi chưa rõ, phải ghi rõ người duyệt và rủi ro trùng chấp nhận; không có chính sách thì không chuyển. |
| SRS-F07.06 | Mỗi lần thử phải gắn đúng bản tin, người nhận, kênh, kết nối, phiên bản và provider reference nếu có. Chuyển provider không được làm mất lịch sử lần trước. |
| SRS-F07.07 | Callback/tra cứu phải được xác thực và đối chiếu đúng lần thử; phản hồi trùng không tạo thêm lần gửi hay event kết quả cho cùng thay đổi. |
| SRS-F07.08 | Với yêu cầu chứa bí mật, mặc định không retry nền khi kết quả chưa rõ; chỉ thực hiện retry theo chỉ dẫn nguồn/hợp đồng đã duyệt còn hiệu lực. |
| SRS-F07.09 | Với lệnh đồng bộ mà nguồn chỉ kích hoạt nghiệp vụ sau phản hồi gửi, hệ thống phải tuân thủ hạn bắt đầu/xử lý đã thỏa thuận; không bắt đầu lần gửi mới sau cửa sổ này hoặc tự chuyển sang hàng gửi nền khi nguồn chưa cho phép. Phần đã đưa ra ngoài vẫn có thể trả kết quả muộn, không được hứa thu hồi. |

**Kiểm chứng:** ZNS timeout không tự gửi SMS; lỗi được phép fallback chỉ tạo một nhánh dự phòng; lỗi kênh khác không gửi lại Email đã đạt; OTP chưa rõ không nằm trong retry nền chung.


### 6.8. SRS-F08 — Hộp thông báo In-app và trạng thái đọc

**Nhu cầu:** “Tôi muốn xem lại tin đúng ngữ cảnh; đọc trên một thiết bị không làm công việc tự hoàn tất.”
**Truy vết:** UC-NTF-13; BR-CHN-05, BR-REC-03/04, BR-SEC-03/06, BR-CFG-05; AT-08/11/15/26.

| Mã | Yêu cầu bắt buộc |
|---|---|
| SRS-F08.01 | Hệ thống phải lưu bản In-app đúng người và scope; chỉ ghi “khả dụng In-app” khi bản tin đã có thể truy cập qua giao diện được phép. |
| SRS-F08.02 | Hệ thống phải cung cấp danh sách có phân trang, chi tiết và số chưa đọc theo cùng điều kiện quyền/ngữ cảnh. Số đếm không được làm lộ tin ngoài scope. |
| SRS-F08.03 | Danh sách phải có thứ tự ổn định theo thời gian tạo bản tin và khóa phân xử; phân trang không bỏ/mở lại tin do thứ tự không xác định. Thời điểm sự việc nguồn hiển thị riêng khi khác thời điểm tạo tin. |
| SRS-F08.04 | Backend phải kiểm tra người/ngữ cảnh và quyền trên tải danh sách, lấy chi tiết và đánh dấu đọc. ID tin hoặc Shop do trình duyệt gửi không tự cấp quyền. |
| SRS-F08.05 | Hệ thống phải ghi trạng thái đọc ở backend theo người và bản tin; đánh dấu đọc lặp lại có cùng kết quả, không tạo tin mới. Trạng thái đọc đã lưu phải được thấy ở lần tải sau trên thiết bị được phép khác. |
| SRS-F08.06 | Đã đọc không được đổi kết quả SMS/Push, hoàn tất đơn hoặc chấp nhận lời mời/yêu cầu thiết bị. Ghi đọc lỗi không được báo là đã lưu thành công. |
| SRS-F08.07 | Khi đổi Shop/app/tài khoản hoặc đăng xuất, giao diện phải tách cache và bỏ dữ liệu của ngữ cảnh không còn hợp lệ. Tin tài khoản chỉ hiện trên app được chính sách cho phép, không mặc định mọi app đều xem được. |
| SRS-F08.08 | Tin không có, bị ẩn/xóa hoặc mất quyền phải có phản hồi an toàn theo hợp đồng. UI phải phân biệt tải, trống và lỗi mà không tiết lộ sự tồn tại/nội dung của tin trái quyền. |

**Kiểm chứng:** Danh sách và số chưa đọc cùng scope; đọc lặp không sai số; thiết bị khác tải được trạng thái đã lưu; mất membership không xem lại tin Shop dù biết ID; đổi tài khoản không thấy cache cũ.

### 6.9. SRS-F09 — Kết quả theo bằng chứng và kết quả một phần

**Nhu cầu:** “Tôi cần biết từng người/kênh đã đến đâu, thay vì một chữ success chung cho cả yêu cầu.”
**Truy vết:** UC-NTF-09/11/12; BR-INT-04, BR-CHN-04/05/09, BR-OPS-01; AT-01/07/09/11/25.

| Mã | Yêu cầu bắt buộc |
|---|---|
| SRS-F09.01 | Hệ thống phải phân biệt kết quả tiếp nhận, tiến trình xử lý, bằng chứng từng kênh, mức đạt mục tiêu và trạng thái đọc; không gộp các nghĩa vào một cờ thành công. |
| SRS-F09.02 | Hệ thống phải ghi kết quả riêng cho trường hợp gửi, người nhận, kênh và lần thử; ghi thời điểm, nguyên nhân, nguồn bằng chứng và revision thay đổi. |
| SRS-F09.03 | “Provider nhận”, “đã giao” và “đã đọc” chỉ được ghi khi có bằng chứng tương ứng được duyệt cho kênh. Nếu provider không có xác nhận giao/đọc thì không được tự suy từ thời gian chờ hoặc mã HTTP thành công. |
| SRS-F09.04 | Mỗi trường hợp gửi phải có tiêu chí đạt được cấu hình trước kích hoạt: kênh bắt buộc, kênh thay thế và mức bằng chứng. Mục tiêu đạt chỉ khi đủ tiêu chí; bị chặn/hủy/hết hạn không được tính là gửi thành công. |
| SRS-F09.05 | Kết quả tổng hợp phải trả phần đạt, chưa đạt, bị chặn/hủy/hết hạn và chưa rõ. Một người hoặc một kênh thành công không được che phần bắt buộc chưa đạt. |
| SRS-F09.06 | Callback đến muộn/sai thứ tự phải được đối chiếu với loại bằng chứng và chính sách xung đột. Phản hồi yếu hơn không làm mất bằng chứng mạnh đã xác nhận; bằng chứng trái nhau phải có trạng thái đối chiếu/lý do, không ghi đè im lặng. |
| SRS-F09.07 | Hệ thống phải cung cấp tra cứu và phát kết quả theo hợp đồng cho consumer được phép; phát lại kết quả có cùng ID thay đổi. Lỗi phát kết quả không được kích hoạt gửi lại tin cho người nhận. |
| SRS-F09.08 | Kết quả phải có định danh nguồn/Notification, mức áp dụng, thời điểm và lý do phù hợp; event kết quả không chứa OTP/token, payload bí mật hoặc liên hệ đầy đủ. |

**Kiểm chứng:** Một Shop có In-app nhưng Push lỗi hiện đủ hai kết quả; chỉ một trong hai người nhận đạt không báo cả yêu cầu hoàn tất; callback cũ không đảo bằng chứng; consumer mất kết nối chỉ nhận lại kết quả.

### 6.10. SRS-F10 — Tiếp nối cảnh báo và yêu cầu thiết bị User

**Nhu cầu:** “Tôi vẫn cần cảnh báo thiết bị hiện có và các nút xử lý đúng trạng thái User.”
**Truy vết:** UC-USR-09, UC-NTF-13; BR-INT-01/06, BR-REC-03/06, BR-SEC-03/06; AT-13/14/15/26.

| Mã | Yêu cầu bắt buộc |
|---|---|
| SRS-F10.01 | Hệ thống phải tiếp nối cảnh báo/yêu cầu của User qua feed/event theo hợp đồng được chọn; chỉ hiển thị cho đúng chủ tài khoản trên app được phép. |
| SRS-F10.02 | Mỗi cảnh báo nguồn phải có khóa đối chiếu ổn định. Cùng cảnh báo xuất hiện qua feed cũ và đường mới không được thành hai tin độc lập trong hộp tin. |
| SRS-F10.03 | Trạng thái yêu cầu bảo mật, khả năng thao tác và mô tả hành động phải do User xác nhận; Notification không tự kết luận từ đã đọc hoặc số lần nhấn nút. |
| SRS-F10.04 | Ứng dụng phải gọi User để cho phép/báo không nhận ra thiết bị; User kiểm tra điều kiện và thực hiện. Notification không mở phiên, tin cậy thiết bị hoặc thu hồi đăng nhập. |
| SRS-F10.05 | Khi nguồn đã xử lý/hết hạn hoặc không lấy được trạng thái cần thiết, UI không được cho thao tác dựa trên bản cũ; phải hiển thị kết quả/lỗi phù hợp và có cách tải lại. |
| SRS-F10.06 | Chuyển Shop không làm cảnh báo tài khoản thành tin của Shop khác; đổi tài khoản/đăng xuất phải loại dữ liệu người trước. Chuyển nguồn phải có chính sách đọc rõ, không dùng trạng thái APPROVED/REJECTED thay trạng thái đọc mới. |

**Kiểm chứng:** Đọc không cho phép đăng nhập; hai nguồn chỉ hiện một cảnh báo; User từ chối hành động thì UI không báo thành công; chuyển Shop vẫn giữ đúng cảnh báo tài khoản.

### 6.11. SRS-F11 — Tra cứu và vận hành tối thiểu

**Nhu cầu:** “Khách nói chưa nhận được tin; tôi cần tìm nguyên nhân bằng reference, không cần đọc bí mật.”
**Truy vết:** UC-NTF-12, phần tối thiểu UC-NTF-16; BR-OPS-01/03/04, BR-SEC-02/04/06; AT-03/16/17/27.

| Mã | Yêu cầu bắt buộc |
|---|---|
| SRS-F11.01 | Vận hành có quyền phải tìm được theo nguồn + reference hoặc ID Notification; thấy luồng, scope, người nhận được phép xem, mẫu/config version, lần thử và kết quả. |
| SRS-F11.02 | Hệ thống phải giải thích “không gửi” bằng lý do cụ thể như thiếu dữ liệu, không có luồng, không đủ quyền nhận, hết hạn, hủy, cấu hình lỗi hoặc ngân sách chặn. |
| SRS-F11.03 | Hệ thống phải phân biệt quyền xem, sửa cấu hình, kích hoạt, dừng, gửi lại và xuất dữ liệu. Quyền xem lịch sử không bao hàm quyền thao tác gửi. |
| SRS-F11.04 | Tra cứu/thao tác nhạy cảm phải có audit: chủ thể, thời điểm, scope, đối tượng, hành động, kết quả và lý do theo chính sách. Audit không chứa bí mật và không được sửa trái quyền. |
| SRS-F11.05 | Hệ thống phải cung cấp số liệu tồn đọng, lỗi, chưa rõ, hết hạn, tiếp nhận và bằng chứng giao theo luồng/kênh/phạm vi; nhãn mock/test tách với gửi thật. |
| SRS-F11.06 | Khi vận hành dừng luồng/kênh/scope, hệ thống phải ghi rõ ảnh hưởng tin mới, tin chờ hoặc cả hai và phần đã ra ngoài; không diễn giải “dừng” thành thu hồi tin đã giao. |
| SRS-F11.07 | Gửi lại thủ công nếu được triển khai phải kiểm tra quyền, lý do, hạn, người nhận và ngân sách; không gửi lại bản đã đạt hay bí mật hết hạn. Gửi mới có chủ đích cần reference riêng liên kết bản cũ. |

**Kiểm chứng:** Người chỉ được xem không gửi lại được; tra reference đủ nguyên nhân; mọi thao tác dừng có phạm vi/tác động; mock không được tính vào tỷ lệ gửi thật. Đợt đầu không bắt buộc xây bộ UI xử lý hàng loạt.

### 6.12. SRS-F12 — Cấu hình tối thiểu có phiên bản và phê duyệt

**Nhu cầu:** “Vận hành phải biết luồng có thể chạy an toàn trước khi bật.”
**Truy vết:** UC-NTF-15 ở mức tối thiểu; BR-CAT-04, BR-CFG-01/02/03/05, BR-TPL-01/03; AT-18/24/27.

| Mã | Yêu cầu bắt buộc |
|---|---|
| SRS-F12.01 | Cấu hình kích hoạt phải có chủ trách nhiệm, scope, mục đích, trigger, trường hợp gửi, người nhận/căn cứ, mẫu, thời điểm/hạn, kênh, retry/fallback, tiêu chí đạt và điều kiện phí khi áp dụng. |
| SRS-F12.02 | Hệ thống phải kiểm tra các thành phần bắt buộc, quyền dùng mẫu/kết nối/app và xung đột phạm vi trước kích hoạt; cấu hình không đạt không được chạy gửi thật. |
| SRS-F12.03 | Cấu hình/mẫu đã dùng phải có phiên bản, thời điểm hiệu lực, chủ thể duyệt/kích hoạt và lịch sử. Thay đổi phải tạo phiên bản mới và chỉ rõ tác động tin chờ. |
| SRS-F12.04 | Tin chờ phải giữ được phiên bản đã chọn hoặc áp phiên bản mới theo chính sách chuyển đổi được duyệt. Dù giữ mẫu cũ, quyền, đích, hạn gửi và điều kiện dừng hiện hành vẫn phải kiểm tra lại. |
| SRS-F12.05 | Cấu hình app/client phải tham chiếu User; brand/định danh gửi và môi trường provider được ánh xạ riêng. Không tạo app hoặc cấp entitlement qua kích hoạt luồng. |
| SRS-F12.06 | Đợt đầu có thể quản trị qua công cụ cấu hình được kiểm soát thay vì màn hình đầy đủ; vẫn phải có quyền, kiểm tra hợp lệ, phiên bản, phê duyệt và audit. |

**Kiểm chứng:** Không bật luồng thiếu tiêu chí đạt hoặc template; thay mẫu không sửa lịch sử; chính sách khóa không bị ghi đè bởi brand; tra được config thực tế dùng trước và sau thay đổi.

### 6.13. SRS-F13 — Bảo vệ bí mật, dữ liệu cá nhân và môi trường

**Nhu cầu:** “Tôi phải gửi mã/link, nhưng hỗ trợ và hệ thống quan sát không được đọc được chúng.”
**Truy vết:** UC-USR-01/03/04, UC-NTF-07/11/12; BR-SEC-01/04/05, BR-TPL-03, BR-OPS-04; AT-03/17/28.

| Mã | Yêu cầu bắt buộc |
|---|---|
| SRS-F13.01 | Bí mật phải truyền qua đường được bảo vệ, chỉ cấp cho thành phần cần xử lý; không đưa OTP/token/mật khẩu vào log, trace, analytics, event kết quả hoặc dữ liệu hỗ trợ. |
| SRS-F13.02 | Nếu phải lưu tạm để gửi, hệ thống phải phân loại, bảo vệ truy cập và xóa theo thời hạn được duyệt; lịch sử lâu dài chỉ giữ metadata/bằng chứng được phép, không giữ bản bí mật nguyên văn. |
| SRS-F13.03 | Hàng lỗi, lỗi exception, callback lưu lại và bản sao dự phòng phải tuân thủ cùng chính sách bảo vệ; chuyển payload sang dead-letter không được làm mất bảo vệ bí mật. |
| SRS-F13.04 | Liên hệ, số tiền và nội dung riêng tư phải che theo quyền và mục đích xem; việc có quyền tra lỗi không tự cho phép xem toàn nội dung người nhận. |
| SRS-F13.05 | Link mang token được coi là bí mật; bản hiển thị hỗ trợ và preview không chứa token. Thông tin cần cho người nhận chỉ đưa qua kênh được chính sách cho phép. |
| SRS-F13.06 | Mock/test phải có nhãn và cấu hình riêng; production không được dùng phản hồi mock/log để chứng minh gửi thật. Thử nghiệm không được tự sử dụng tập đích production. |

**Kiểm chứng:** QA dùng mã/token giả có dấu nhận biết và tìm trên toàn bộ đường quan sát, hàng lỗi, API hỗ trợ; không có bản nguyên văn. Kiểm tra cả exception và lỗi provider, không chỉ log đường thành công.

### 6.14. SRS-F14 — Điều kiện chi phí và hạn mức được giao

**Nhu cầu:** “Hệ thống phải tôn trọng quyền chi nhưng không tự biến mình thành sổ ngân sách.”
**Truy vết:** UC-NTF-06, phần kiểm soát tối thiểu UC-NTF-17; BR-COST-01/03; AT-19/25.

| Mã | Yêu cầu bắt buộc |
|---|---|
| SRS-F14.01 | Kênh phát sinh phí phải có đơn vị chịu phí và chính sách/căn cứ cho phép đã được bên sở hữu xác nhận; hệ thống không tự chọn người trả phí theo NVC/app. |
| SRS-F14.02 | Hệ thống phải thực thi điều kiện ngân sách ở thời điểm theo hợp đồng, gồm retry/fallback có phí; không có căn cứ đủ tin cậy thì chờ/chặn, không tự coi là được phép. |
| SRS-F14.03 | Quyết định chặn phải có lý do, phạm vi và phần ảnh hưởng. Không âm thầm bỏ OTP/cảnh báo rồi trả thành công hoặc tự vượt hạn mức vì tin được gắn khẩn. |
| SRS-F14.04 | Hệ thống phải ghi tham chiếu quyết định chi và kết nối/lần thử liên quan theo trách nhiệm được giao; không tự ghi sổ COD, thanh toán hoặc đối soát tài chính. |
| SRS-F14.05 | Nếu dùng cơ chế cấp phép/giữ hạn mức bên ngoài, hợp đồng phải phân biệt dùng lại cùng lần thử với xin phép cho lần mới; giải phóng/quyết toán theo kết quả và trách nhiệm bên sở hữu, không tự hoàn phí khi kênh chưa rõ. |

**Kiểm chứng:** Fallback không vượt điều kiện chi; ngân sách chặn trả lý do; mất kết nối kiểm soát phí không mặc định cho phép. Sổ ngân sách và đơn giá chưa được chọn ở SRS.


### 6.15. SRS-F15 — Thiết lập và quản lý kết nối kênh

**Nhu cầu:** “Tôi cần biết kết nối nào dùng cho app/brand nào và nó đã đủ điều kiện gửi thật chưa.”
**Truy vết:** UC-NTF-08/15 ở mức tối thiểu; BR-CHN-01/06/07/08/09, BR-CFG-05, BR-SEC-02/04; AT-17/24/27.

| Mã | Yêu cầu bắt buộc |
|---|---|
| SRS-F15.01 | Hệ thống phải quản lý tham chiếu kết nối gồm kênh/provider, môi trường, scope/app/brand được phép, định danh gửi, khả năng trả bằng chứng và trạng thái sử dụng. |
| SRS-F15.02 | Bí mật kết nối phải được bảo vệ và chỉ cấp cho thành phần cần dùng; không hiển thị nguyên văn trong API tra cứu, audit hoặc màn hình. |
| SRS-F15.03 | Trước cho luồng gửi thật, hệ thống phải kiểm tra kết nối, định danh gửi, ánh xạ mẫu và phạm vi sử dụng. Kết nối thiếu/hết hiệu lực không được báo là gửi thành công. |
| SRS-F15.04 | Hệ thống phải hỗ trợ kiểm tra/gửi thử bằng đích được phép và ghi bằng chứng môi trường; kết quả thử không được trộn vào báo cáo production. |
| SRS-F15.05 | Thay/thu hồi thông tin kết nối phải có quyền, audit và ảnh hưởng tin chờ rõ ràng; mỗi lần thử giữ được tham chiếu kết nối thực tế đã dùng. |
| SRS-F15.06 | Kết nối bị dừng/lỗi không được tự chuyển sang tài khoản provider của brand, pháp nhân hoặc môi trường khác. Chuyển kết nối phải qua chính sách đã duyệt. |

**Kiểm chứng:** Provider đúng kênh nhưng sai brand không được dùng; đổi credential có audit nhưng không lộ key; callback từ kết nối khác không cập nhật được lần thử này.

### 6.16. SRS-F16 — Push và vòng đời đích thiết bị

**Nhu cầu:** “Push phải đến đúng app và đúng chủ thiết bị hiện tại; đừng gửi tin người cũ sau đăng xuất.”
**Truy vết:** UC-NTF-08/09, UC-ORD-05; BR-CHN-06/09, BR-REC-06, BR-CFG-05; AT-11/12/26.

| Mã | Yêu cầu bắt buộc |
|---|---|
| SRS-F16.01 | Hệ thống phải gửi Push chỉ tới đích được chủ thiết bị/app xác nhận, có quan hệ tài khoản–app/client–đích và trạng thái hiệu lực theo hợp đồng. Token một mình không đủ xác nhận người nhận. |
| SRS-F16.02 | Hợp đồng phải xác định bên đăng ký, tra, thay và thu hồi đích; Notification phải tiêu thụ thay đổi theo version/hiệu lực và không dùng bản cũ hơn để hồi sinh token. |
| SRS-F16.03 | Trước gửi, hệ thống phải kiểm tra đích còn hợp lệ theo hợp đồng độ mới. Khi đã nhận thu hồi/đổi chủ, không bắt đầu lần gửi mới cho chủ cũ qua đích đó. Phần đã ra ngoài phải được báo đúng giới hạn. |
| SRS-F16.04 | Token bị provider xác nhận mất hiệu lực phải bị loại khỏi tập có thể dùng và phản hồi cho bên sở hữu theo hợp đồng; không retry vô hạn cùng token lỗi. |
| SRS-F16.05 | Khi một người có nhiều đích, chính sách phải nêu gửi mọi đích hay tập được chọn, mức đạt tính theo người hay đích; kết quả giữ chi tiết để một token lỗi không bị che bởi token khác. |
| SRS-F16.06 | Push phải dùng preview được phép và đường mở đúng app/ngữ cảnh; provider nhận không tự chứng minh thiết bị hiển thị hay người đã đọc. |

**Kiểm chứng:** Thu hồi rồi nhận bản đăng ký cũ không làm token sống lại; cùng token đổi chủ không nhận tin người trước; test trên client thực đã chọn. Thiết bị tin cậy User và đích nhận Push là hai loại dữ liệu khác nhau.

### 6.17. SRS-F17 — Email thực cho lời mời và nội dung được phép

**Nhu cầu:** “Tôi cần gửi lời mời qua email tới người chưa là thành viên, và biết nếu thư bị trả lại.”
**Truy vết:** UC-USR-03/04; BR-CHN-08/09, BR-TPL-01/05, BR-SEC-01; AT-05/06/07/24.

| Mã | Yêu cầu bắt buộc |
|---|---|
| SRS-F17.01 | Email phải dùng địa chỉ do nguồn xác nhận, sender/brand và mẫu được duyệt. Gmail là loại địa chỉ nhận Email, không phải kênh thứ sáu. |
| SRS-F17.02 | Email lời mời phải giữ đúng tổ chức/Shop, người mời và link nguồn; không yêu cầu người nhận đã có membership trước khi gửi lời mời. |
| SRS-F17.03 | Gửi cho nhiều người phải bảo vệ địa chỉ và dữ liệu riêng của từng người; không lộ danh sách người nhận khác qua To/Cc hoặc nội dung cá nhân hóa. |
| SRS-F17.04 | Hệ thống phải ghi phản hồi nhận thư, bounce/chặn và bằng chứng giao theo khả năng provider; không dùng việc mở pixel hoặc provider nhận thay cho chấp nhận lời mời. |
| SRS-F17.05 | Email chứa nội dung nhạy cảm phải tuân thủ chính sách link/đính kèm đã duyệt; link đã hết hạn/thu hồi vẫn do User kiểm tra khi mở. |

**Kiểm chứng:** Email thật mở đúng lời mời; hard bounce dừng retry không hợp lệ; địa chỉ người A không lộ cho B; thư lỗi không xóa lời mời.

### 6.18. SRS-F18 — SMS thực và hạn của mã xác thực

**Nhu cầu:** “Tôi cần nhận đúng mã qua SMS; hệ thống không được hứa rằng tôi đã nhận chỉ từ một biên nhận.”
**Truy vết:** UC-USR-01, UC-ORD-05; BR-CHN-07/09, BR-INT-02, BR-TIME-01; AT-01/04/09/19.

| Mã | Yêu cầu bắt buộc |
|---|---|
| SRS-F18.01 | SMS phải có số/định dạng theo hợp đồng, định danh gửi và kết nối hợp lệ. Hệ thống không tự sửa một số mơ hồ thành số của người khác. |
| SRS-F18.02 | Nội dung sau dựng mẫu phải được kiểm tra giới hạn kênh; chính sách phải nêu việc tách đoạn và tác động phí nếu có. Không cắt mất mã, hạn hoặc hướng dẫn quan trọng. |
| SRS-F18.03 | OTP chỉ dùng mã và đích User cấp; hệ thống không phát sinh, thay mã hoặc xác minh. SMS hết hạn gửi không được bắt đầu thử mới. |
| SRS-F18.04 | Receipt phải dựa trên phản hồi gửi thật. Đã giao chỉ ghi khi có bằng chứng đã duyệt; provider reference, log hoặc accepted nội bộ không đủ để tự nâng mức. |
| SRS-F18.05 | SMS dự phòng phải qua cùng kiểm tra hạn, đích và ngân sách như SMS chính; không tự dùng cho mọi lỗi ZNS. |

**Kiểm chứng:** Nội dung dài có kết quả rõ; expired không ra kênh; biên nhận mock không được User coi là gửi thực; fallback có phí phải được cho phép.

### 6.19. SRS-F19 — Zalo ZNS đúng mẫu, định danh và chính sách

**Nhu cầu:** “Tin giao thất bại phải dùng mẫu ZNS và đích do Order xác nhận; lỗi ZNS được giải thích đúng.”
**Truy vết:** UC-ORD-05; BR-CHN-03/07/09, BR-TPL-01/02, BR-CFG-05; AT-09/11/20/24.

| Mã | Yêu cầu bắt buộc |
|---|---|
| SRS-F19.01 | ZNS phải có kết nối/định danh gửi và ánh xạ mẫu provider được duyệt cho mục đích, brand và scope; không thay bằng một mẫu gần giống để vượt điều kiện. |
| SRS-F19.02 | Hệ thống phải kiểm tra biến và điều kiện kênh/provider theo hợp đồng cấu hình hiện hành; không ghi “giao” khi provider từ chối mẫu hoặc đích. |
| SRS-F19.03 | Người nhận hàng có thể nhận ZNS theo liên hệ giao dịch được Order cấp mà không có tài khoản User; việc nhận tin không tạo quyền truy cập Shop. |
| SRS-F19.04 | Hệ thống phải giữ nguyên mã lỗi provider cần tra cứu đồng thời ánh xạ sang lý do an toàn; timeout giữ chưa rõ, SMS chỉ theo chính sách fallback đã duyệt. |
| SRS-F19.05 | Mẫu/định danh/đích từ luồng ZNS cũ chỉ được chuyển sau kiểm kê, xác nhận và thử thật theo kế hoạch; không giả định mọi cấu hình cũ dùng được ngay. |

**Kiểm chứng:** Sai mẫu dừng đúng lý do; liên hệ không có User vẫn gửi được khi hợp lệ; timeout không tự nhân SMS; chuyển luồng có bằng chứng mẫu thật.

### 6.20. SRS-F20 — Áp dụng mục đích và lựa chọn nhận tin

**Nhu cầu:** “Việc tôi từ chối quảng cáo không được vô tình chặn mã xác thực; liên hệ không an toàn cũng không được nhận bí mật.”
**Truy vết:** UC-NTF-06, phần áp dụng chính sách của UC-NTF-14; BR-CLS-01/02/03, BR-PREF-01/02/04; AT-08/19/24.

| Mã | Yêu cầu bắt buộc |
|---|---|
| SRS-F20.01 | Hệ thống phải áp dụng mục đích, danh sách chặn và căn cứ lựa chọn do bên sở hữu cấp trước gửi theo hợp đồng; không suy consent từ việc có số/email. |
| SRS-F20.02 | Từ chối Marketing không tự chặn luồng bắt buộc đã duyệt; luồng bắt buộc vẫn phải kiểm tra đích an toàn, bảo mật, điều kiện provider và kênh được phép. |
| SRS-F20.03 | Hệ thống phải phân biệt chặn mục đích, chặn kênh và liên hệ bị báo sai/xâm phạm. Gửi bí mật tới liên hệ bị chặn an toàn chỉ xử lý theo chỉ dẫn User được duyệt. |
| SRS-F20.04 | Khi thiếu căn cứ cần thiết, hệ thống phải chờ/chặn có lý do. Đợt đầu không được tự bật gửi Marketing hoặc xây consent bằng cách mặc định mọi người đồng ý. |

**Kiểm chứng:** Opt-out quảng cáo không chặn OTP hợp lệ; liên hệ sai bị chặn không nhận bí mật. Trang người dùng tự quản lý lựa chọn chưa thuộc đợt đầu.

### 6.21. SRS-F21 — Giới hạn gửi và xử lý quá tải tối thiểu

**Nhu cầu:** “Kênh hoặc nguồn tăng tải không được làm OTP hết hạn âm thầm hay gửi vượt giới hạn.”
**Truy vết:** UC-NTF-06/09, phần tối thiểu UC-NTF-17; BR-CHN-02, BR-TIME-01, BR-COST-03; AT-04/19/27.

| Mã | Yêu cầu bắt buộc |
|---|---|
| SRS-F21.01 | Hệ thống phải áp dụng giới hạn kênh/provider và nguồn theo cấu hình được duyệt; không tự tăng gửi vượt giới hạn khi backlog lớn. |
| SRS-F21.02 | Khi chưa tiếp nhận được vì tải, hệ thống phải trả kết quả hạn chế/tạm không sẵn sàng có hướng retry theo hợp đồng; không trả accepted cho công việc chưa giữ được. |
| SRS-F21.03 | Công việc đã accepted phải có trạng thái chờ, hạn xử lý và khả năng tra cứu; khi hết hạn thì ghi expired/chặn phù hợp, không bỏ khỏi hệ thống mà không có kết quả. |
| SRS-F21.04 | Ưu tiên theo mục đích/mức khẩn đã duyệt phải được áp dụng khi điều phối; traffic luồng khác không được làm OTP vượt hạn xử lý đã cam kết. Không cần gom nhóm hoặc lịch nhắc nâng cao để đáp ứng yêu cầu này. |

**Kiểm chứng:** Burst có kết quả nhận/chặn rõ; backlog quá hạn có lý do; provider rate-limit không gây retry storm. Mức tải và thuật toán điều phối thuộc tham số/thiết kế.

### 6.22. SRS-F22 — Chuyển đổi nguồn cũ có thể kiểm soát

**Nhu cầu:** “Bật Notification mới không được gửi gấp đôi hoặc làm mất cảnh báo đang dùng.”
**Truy vết:** UC-USR-09, UC-ORD-05 và luồng cũ được chọn; BR-TRG-02/05, BR-OPS-04; AT-14/20/23.

| Mã | Yêu cầu bắt buộc |
|---|---|
| SRS-F22.01 | Mỗi luồng chuyển đổi phải có chủ sở hữu, nguồn cũ/mới, mẫu/đích, khóa đối chiếu và điểm chuyển trách nhiệm gửi rõ ràng. |
| SRS-F22.02 | Không được bật hai nguồn gửi thật cho cùng tập sự việc mà không có cơ chế phân chia/chống trùng đã kiểm chứng. Chạy đối chiếu không gửi thật phải có nhãn riêng. |
| SRS-F22.03 | Chuyển feed cảnh báo phải giữ ID nguồn và hành động tại User; chính sách giữ/khởi tạo trạng thái đọc phải được ghi và kiểm tra, không tự suy từ trạng thái nghiệp vụ. |
| SRS-F22.04 | Kế hoạch khôi phục phải xác định phần đã gửi, chưa gửi và chưa rõ; quay lại nguồn cũ không được replay mù quáng phần đã đạt. |
| SRS-F22.05 | Chỉ coi luồng chuyển đổi hoàn tất khi có bằng chứng gửi/hiển thị thật, đối chiếu số liệu, ngoại lệ và cách dừng/khôi phục trong phạm vi được chọn. |

**Kiểm chứng:** Chuyển một luồng thử và khôi phục không nhân tin; cảnh báo vẫn xử lý ở User; mock/shadow không trở thành bằng chứng gửi thật.

## 7. Cách áp dụng vào bốn luồng ưu tiên

Phần này diễn giải các yêu cầu ở mục 6 bằng tình huống người dùng/module thực tế. Không tạo thêm trách nhiệm nghiệp vụ cho Notification.

### 7.1. R1-OTP — User gửi mã xác thực

1. User chọn liên hệ, phát hành mã và reference cho lần gửi, cấp hạn gửi và cửa sổ xử lý có phản hồi.
2. Notification kiểm tra nguồn, chống trùng, mẫu, đích, kết nối SMS và điều kiện phí; không cần người nhận đã có session/membership.
3. Notification gửi theo chính sách đồng bộ đã thống nhất và trả đúng mức bằng chứng. Cấu hình R1 phải có phản hồi gửi thật tối thiểu do provider tiếp nhận; nếu User cần mức cao hơn thì phải có hợp đồng và bằng chứng kênh tương ứng.
4. User quyết định có kích hoạt mã và cập nhật counters hay không. Notification không xác minh mã hoặc cập nhật session.
5. Timeout có thể là provider đã nhận nhưng phản hồi bị mất. Notification giữ chưa rõ và reference tra cứu; không tự xếp thêm lần gửi nền làm mã ứng viên đến sau khi User đã bỏ lần cấp đó.
6. Nếu User cấp lại mã, đây là yêu cầu gửi mới và User chỉ định cách dừng bản cũ. Notification không dùng session ID để suy mã nào mới nhất.
7. Callback tới sau chỉ bổ sung bằng chứng lần gửi; không hồi sinh mã, đảo quyết định User hoặc tự gửi thêm.

| Tình huống kiểm chứng | Kết quả mong đợi | Truy vết |
|---|---|---|
| Mất phản hồi sau provider nhận | Tra cùng reference thấy kết quả/bằng chứng đã có hoặc chưa rõ; không tạo lần cấp mã hay SMS mới do replay | SRS-F01/02/07/09/18; AT-01/02/04 |
| User chưa kích hoạt mã vì timeout | Không có retry nền tự phát; nguồn quyết định mã cũ/mới, Notification không can thiệp | SRS-F06/07/13; AT-04/23 |
| Mã hết hạn khi còn chờ | Không bắt đầu SMS mới; có expired và lý do | SRS-F06/18/21; AT-04 |

Cửa sổ chờ phản hồi, hạn bắt đầu thử, TTL mã và thời điểm người thực sự nhận SMS là bốn khái niệm khác nhau; hợp đồng phải nói rõ. Không cam kết thu hồi SMS đã ra ngoài hoặc đảm bảo mã đã gửi chắc chắn còn được User chấp nhận.

### 7.2. R1-INVITE — Lời mời nhân viên/thành viên

User chỉ bàn giao lời mời sau commit. Email có đúng tên tổ chức/Shop, người mời và link User cấp. Người được mời có thể chưa có membership. Notification báo tiến trình gửi riêng; lời mời vẫn tồn tại khi email lỗi.

Hai lời mời cho cùng email ở hai Shop là hai đối tượng độc lập. User cấp lại/thu hồi/chấp nhận lời mời thì User cấp chỉ dẫn dừng bản chờ; Notification không tạo membership và không kết luận “đã chấp nhận” từ email đã giao.

| Tình huống kiểm chứng | Kết quả mong đợi | Truy vết |
|---|---|---|
| Email bounce sau khi provider nhận | Lời mời không bị xóa; kết quả kênh thể hiện bounce và lý do | SRS-F09/17; AT-05/07 |
| Hai Shop mời cùng email | Hai bản giữ đúng dữ liệu/link, không bị hợp nhất theo email | SRS-F02/04/05/17; AT-05/08 |
| Lời mời thu hồi trước retry | Dừng phần chờ theo reference; email đã ra ngoài không được hứa thu hồi | SRS-F06/07; AT-06 |

Internal FE phải cập nhật thông điệp theo kết quả thật: “đã tạo lời mời” khác “đã yêu cầu gửi email”, khác “provider đã nhận” và khác “người được mời đã chấp nhận”.

### 7.3. R1-ORDER — Order xác nhận giao thất bại

Order bàn giao sự việc đã chốt, người nhận/căn cứ, liên hệ giao dịch, chặng/lần giao khi cần, dữ liệu công bố và chỉ dẫn còn cần gửi. Notification không tiêu thụ tracking NVC thô để kết luận thất bại.

- Trường hợp Shop: bản tin In-app/Push đúng Shop/app, link mở Order có kiểm tra quyền.
- Trường hợp liên hệ của đơn: ZNS theo mẫu được duyệt, SMS dự phòng theo điều kiện cụ thể.
- Thêm lần giao thất bại thực sự: Order cấp reference mới; webhook lặp cho cùng sự việc dùng reference cũ.
- Tin không còn cần gửi: Order chỉ định tập cần dừng. Giao đủ hay giao một phần là quyết định Order; Notification chỉ áp lệnh/chỉ dẫn.

| Tình huống kiểm chứng | Kết quả mong đợi | Truy vết |
|---|---|---|
| Shop nhận In-app, Push lỗi; khách nhận ZNS | Hiển thị kết quả từng trường hợp/kênh; mức đạt theo cấu hình, không giấu lỗi Push | SRS-F03/07/08/09/16/19; AT-11/25 |
| ZNS timeout | Giữ chưa rõ; SMS chưa được chạy nếu không có chính sách cho trường hợp chưa rõ | SRS-F07/09/19; AT-09 |
| Order hủy mục tiêu trong lúc event cũ đến trễ | Dừng đúng tập theo hợp đồng; không tự hủy sự việc độc lập khác của cùng đơn | SRS-F02/06; AT-10/23 |

Order chịu trách nhiệm đường phát bền vững và ngữ nghĩa sự việc. Việc ghi một dòng outbox không được dùng làm bằng chứng Notification đã tiếp nhận.

### 7.4. R1-SECURITY — Cảnh báo và yêu cầu về thiết bị

Cảnh báo thuộc tài khoản và xuất hiện trên app được chính sách cho phép. ID cảnh báo User là căn cứ đối chiếu feed cũ/mới. Trạng thái đọc lưu riêng; trạng thái yêu cầu và khả năng hành động lấy từ User.

| Tình huống kiểm chứng | Kết quả mong đợi | Truy vết |
|---|---|---|
| Người dùng chỉ đánh dấu đọc | Không cho phép đăng nhập, không tin cậy/thu hồi thiết bị | SRS-F08/10; AT-13 |
| Đổi Shop khi cảnh báo đang chờ | Vẫn thuộc chủ tài khoản, không bị chuyển thành tin công việc của Shop mới | SRS-F05/08/10; AT-08/15 |
| User xử lý yêu cầu trên thiết bị khác | Tải lại lấy trạng thái nguồn; nút không còn hợp lệ bị vô hiệu, trạng thái đọc vẫn riêng | SRS-F08/10; AT-13/26 |
| Feed cũ và event mới cùng có cảnh báo | Một mục theo ID nguồn, giữ hành động User | SRS-F02/10/22; AT-14 |

## 8. Ý nghĩa trạng thái và quy tắc tổng hợp

### 8.1. Trạng thái theo từng mức

| Mức | Kết quả cần biểu diễn | Ý nghĩa |
|---|---|---|
| Đầu vào | Từ chối; xung đột; cùng yêu cầu đã có | Đầu vào không được nhận, khóa bị dùng sai, hoặc đây là replay hợp lệ. |
| Tiếp nhận | Đã tiếp nhận | Công việc đã được giữ bền vững; chưa chứng minh gửi. |
| Xử lý | Chờ; đang xử lý; đang đối chiếu; đã kết thúc xử lý theo chính sách | Tiến trình xử lý nội bộ; kết thúc xử lý không luôn có nghĩa đã giao. |
| Bản/kênh chưa gửi | Không gửi theo chính sách; hủy; hết hạn | Có lý do riêng; không được tính là thành công kênh. |
| In-app | Khả dụng | Bản tin được người có quyền truy cập; chưa chứng minh đã đọc. |
| Lần thử ngoài | Provider nhận; lỗi chắc chắn; chưa rõ | Bằng chứng tiếp nhận/thất bại, hoặc chưa biết việc gửi ngoài hệ thống. |
| Kênh | Đã giao; bounce/không giao; thất bại cuối cùng | Chỉ theo bằng chứng kênh và chính sách retry đã hết phương án. |
| Người đọc | Chưa đọc; đã đọc | Chỉ thuộc In-app/ngữ nghĩa đọc đã thống nhất. |
| Mục tiêu gửi | Đạt; chưa đạt; đạt một phần; không còn áp dụng theo chỉ dẫn | Đánh giá các điều kiện bắt buộc; trả kèm bằng chứng và phần chi tiết. |
| Nghiệp vụ | Trạng thái User/Order cung cấp | Dữ liệu nguồn để hiển thị/điều kiện; không do Notification quyết định. |

### 8.2. Tiêu chí đạt phải được cấu hình

Mỗi trường hợp phải định nghĩa tập nhận, kênh bắt buộc/tùy chọn/thay thế và mức bằng chứng. Với nhiều người bắt buộc, chỉ một người đạt không làm cả trường hợp đạt. Trường hợp thay thế hợp lệ có thể đạt qua kênh dự phòng, nhưng lần thử kênh chính vẫn giữ kết quả.

| Cấu hình minh họa — không phải mặc định mọi luồng | Kết quả thực tế | Cách báo đúng |
|---|---|---|
| In-app bắt buộc; Push bổ sung | In-app khả dụng, Push lỗi | Mục tiêu In-app đạt; kết quả các kênh một phần và Push lỗi vẫn hiện rõ. |
| In-app và Push đều bắt buộc | In-app khả dụng, Push chưa rõ | Mục tiêu chưa đạt đủ/đạt một phần; còn Push chưa rõ. |
| ZNS hoặc SMS dự phòng ở mức provider nhận | ZNS từ chối đủ điều kiện; SMS provider nhận | Mục tiêu đạt qua SMS; ZNS thất bại giữ nguyên; không báo đã giao nếu chưa có bằng chứng. |
| Hai người đều bắt buộc | Người A đạt, người B bị chặn ngân sách | Mục tiêu đạt một phần; người B có lý do chặn; không gộp cả hai thành success. |

Tiêu chí của từng luồng thật phải được PO nguồn và Notification xác nhận ở DEC-10. Ví dụ chỉ giúp kiểm tra thuật toán; không tự chốt Push tùy chọn cho mọi Shop.

### 8.3. Callback trễ, xung đột và đóng xử lý

- Callback đã xác thực chỉ cập nhật lần thử liên quan; không tạo thông báo mới.
- Đã có bằng chứng giao, callback cũ chỉ báo provider nhận không được hạ thành chưa giao.
- Bằng chứng trái nhau phải có dấu đối chiếu và chính sách giải quyết; giữ được bằng chứng trước/sau.
- Khi đã hủy phần chờ, callback của phần ra ngoài vẫn có thể báo giao; phải hiện cả quyết định hủy và giới hạn thực tế.
- Nếu hết cửa sổ đối chiếu mà vẫn chưa rõ, có thể kết thúc xử lý theo chính sách nhưng phải giữ kết quả ngoài “chưa rõ”; không tự đổi thành đã giao hay thất bại chắc chắn.
- Kết quả Notification không dùng để xác nhận Order đã giao đủ hoặc User đã chấp nhận lời mời.

## 9. Lỗi và phản hồi mà nguồn cần hiểu

| Nhóm lý do | Ví dụ | Nguồn/vận hành cần biết | Có được tự retry? |
|---|---|---|---|
| Không được phép | Nguồn trái scope; không đủ quyền xem/hủy | Không nhận hoặc không cung cấp dữ liệu; không lộ thông tin ngoài quyền | Không, phải sửa quyền/căn cứ. |
| Đầu vào không hợp lệ | Thiếu biến, version không hỗ trợ, dữ liệu quá giới hạn | Trường/điều kiện cần sửa; không có gửi bằng dữ liệu đoán | Không retry nguyên payload. |
| Xung đột reference | Cùng khóa nhưng khác đích/nội dung | Yêu cầu cũ giữ nguyên; phải xác định replay hay yêu cầu mới | Không ghi đè. |
| Không đủ điều kiện gửi | Không có luồng, lựa chọn chặn, hết hạn, đã hủy | Phần ảnh hưởng và lý do | Chỉ theo quyết định/chỉ dẫn mới hợp lệ. |
| Kết nối/mẫu/đích không hợp lệ | Token thu hồi, sender sai, mẫu chưa duyệt | Cần sửa cấu hình hoặc dữ liệu bởi bên sở hữu | Không retry vô hạn. |
| Hạn mức/chi phí | Bên kiểm soát không cho phép | Căn cứ quyết định và phần bị chặn/chờ | Chỉ khi có phép/chính sách mới. |
| Lỗi tạm thời chắc chắn chưa gửi | Tạm không nhận, provider rate-limit có phản hồi xác định | Khi nào có thể thử lại theo hợp đồng | Hữu hạn, còn hạn và được phép. |
| Kết quả chưa rõ | Timeout sau khả năng provider nhận; crash giữa gửi và ghi kết quả | Reference đối chiếu và hướng tra/đợi | Không retry mù quáng; theo SRS-F07. |
| Hệ thống không giữ được yêu cầu | Lỗi ghi nhận trước accepted | Không có cam kết đã nhận | Nguồn có thể phát lại cùng khóa theo hợp đồng. |

Mỗi phản hồi lỗi phải có mã ổn định, mức áp dụng, correlation, thông điệp an toàn và hướng xử lý khi có. Thông điệp người dùng không chứa stack trace, credential, payload bí mật hoặc liên hệ người khác. Catalog mã chính thức do hợp đồng tích hợp quy định.


## 10. Yêu cầu phi chức năng và cách đo

Các yêu cầu định tính dưới đây bắt buộc trong đợt đầu. Một mục tiêu định lượng chỉ có hiệu lực nghiệm thu khi có giá trị, đơn vị, phạm vi, môi trường và cửa sổ đo được chủ trách nhiệm xác nhận ở mục 11.

### 10.1. SRS-N01 — Thời gian phản hồi luồng xác thực

| Mã | Yêu cầu và bằng chứng |
|---|---|
| SRS-N01.01 | Hệ thống phải đo riêng thời gian nhận lệnh → phản hồi Notification, nhận lệnh → bắt đầu gọi kênh và gọi kênh → bằng chứng nhận; cùng correlation/reference để đối chiếu. |
| SRS-N01.02 | Luồng OTP phải có ngân sách thời gian chờ/cửa sổ xử lý theo hợp đồng User; không bắt đầu lần gửi mới ngoài cửa sổ được phép để bù timeout. |
| SRS-N01.03 | Báo cáo phải có số mẫu, p95/p99, lỗi/chưa rõ và tỷ lệ vượt cửa sổ theo tải được duyệt; không dùng thời gian phản hồi API thay thời gian người dùng nhận/nhập mã. |

**Cách kiểm chứng:** Đo end-to-end với User và SMS thật trên môi trường được chọn, gồm phản hồi chậm, mất phản hồi và backlog. Tham số SRS-P01/02; DEC-05/06.

### 10.2. SRS-N02 — Độ trễ thông báo sau sự việc

| Mã | Yêu cầu và bằng chứng |
|---|---|
| SRS-N02.01 | Hệ thống phải lưu/đo thời điểm nguồn xác nhận, Notification tiếp nhận, tạo bản và đưa ra kênh để tách độ trễ nguồn, hàng chờ và kênh. |
| SRS-N02.02 | Độ trễ phải được báo theo luồng/trường hợp/kênh và đỉnh tải; event trễ từ nguồn vẫn giữ thời điểm sự việc, không đổi thành sự việc mới xảy ra. |
| SRS-N02.03 | Chỉ tiêu từ nguồn tới người nhận chỉ được công bố khi có bằng chứng kênh tương ứng. Nếu chỉ đo tới provider nhận thì tên chỉ số phải nói đúng mức đó. |

**Cách kiểm chứng:** Phát sự việc Order/lời mời từ nguồn thật; đo từng đoạn và trường hợp đến trễ. Tham số SRS-P02/03; DEC-06/10.

### 10.3. SRS-N03 — Tải, giới hạn và hộp tin

| Mã | Yêu cầu và bằng chứng |
|---|---|
| SRS-N03.01 | Hệ thống phải công bố giới hạn đầu vào, số người/đích/kênh, tốc độ nhận và phân trang; vượt giới hạn có kết quả an toàn trước phát sinh công việc không kiểm soát. |
| SRS-N03.02 | Tải thử phải bao gồm traffic OTP đồng thời với event hàng loạt và đọc In-app; đo độ trễ, lỗi, backlog, tỷ lệ hết hạn và mức sử dụng kênh. |
| SRS-N03.03 | Quá tải không được làm mất yêu cầu đã accepted hoặc vượt giới hạn provider/ngân sách; có phản hồi hạn chế tải và số liệu đủ để vận hành điều chỉnh. |

**Cách kiểm chứng:** Load/burst test theo cơ cấu tải được nguồn và vận hành xác nhận, không suy quy mô chỉ từ số người dùng nền tảng. Tham số SRS-P03/04/05; DEC-06.

### 10.4. SRS-N04 — Bền vững, phục hồi và khả dụng

| Mã | Yêu cầu và bằng chứng |
|---|---|
| SRS-N04.01 | Trong mô hình lỗi đã duyệt, restart/crash tiến trình sau accepted phải phục hồi yêu cầu và tiến trình; không mất thông tin khiến phần đã đạt bị gửi lại tự động. |
| SRS-N04.02 | Crash sau đưa tin ra ngoài nhưng trước lưu receipt phải được biểu diễn/đối chiếu như kết quả có thể chưa rõ; không tuyên bố exactly-once qua provider khi chưa có cơ chế/bằng chứng tương ứng. |
| SRS-N04.03 | Kế hoạch backup/restore phải nêu RPO/RTO và cách giữ/chuyển cửa sổ chống trùng; khôi phục không replay mù quáng toàn lịch sử. |
| SRS-N04.04 | Hệ thống phải có chỉ tiêu khả dụng cho nhận lệnh, tra cứu và hộp tin, cùng cách đo và phạm vi loại trừ đã duyệt; lỗi một kênh phải quan sát riêng với lỗi toàn dịch vụ. |

**Cách kiểm chứng:** Fault injection ở các cửa sổ trước/sau accepted, trước/sau provider nhận, phát kết quả và restore. Tham số SRS-P06/07/08; DEC-06/09/10.

### 10.5. SRS-N05 — Quyền và bảo mật

| Mã | Yêu cầu và bằng chứng |
|---|---|
| SRS-N05.01 | Backend phải yêu cầu quyết định quyền đủ căn cứ cho hành động/tài nguyên; DENY, NOT_EVALUATED hoặc không kiểm tra được không được cấp truy cập. Scope rộng chỉ hỗ trợ lọc, không thay kiểm tra tài nguyên cần thiết. |
| SRS-N05.02 | Bộ kiểm chứng phải bao gồm giả mạo identity/Shop/app, mất membership/entitlement, ID tin người khác, nguồn trái scope, callback giả và token thiết bị cũ. Bất kỳ ca lộ dữ liệu hoặc vượt quyền nào đều không đạt. |
| SRS-N05.03 | Mọi đường quan sát/hàng lỗi/tra cứu phải kiểm chứng bằng bí mật giả có dấu nhận biết; không được tìm thấy nguyên văn trong nơi bị cấm theo SRS-F13. |
| SRS-N05.04 | Truyền dữ liệu và bí mật kết nối phải được bảo vệ theo cơ chế nền tảng và thiết kế được duyệt; quyền service chạy nền giới hạn vào công việc được giao. |

**Cách kiểm chứng:** Negative authorization tests, kiểm tra quyền hiện tại và rà đường chứa bí mật. DEC-03/09/12; không chờ SLA mới kiểm chứng được các điều cấm.

### 10.6. SRS-N06 — Lưu, xóa, audit và khả năng hỗ trợ

| Mã | Yêu cầu và bằng chứng |
|---|---|
| SRS-N06.01 | Hệ thống phải có thời gian lưu riêng cho bí mật tạm, nội dung, kết quả, chống trùng, đích nhận và audit; không dùng một TTL chung cho mọi dữ liệu. |
| SRS-N06.02 | Xóa/ẩn phải theo quyền và chính sách; có bằng chứng thực hiện và ảnh hưởng backup/restore rõ ràng. Không khôi phục dữ liệu đã xóa trái chính sách qua bản sao cũ. |
| SRS-N06.03 | Audit và metadata cần tra cứu phải giữ liên kết đủ dùng trong cửa sổ được duyệt mà không giữ bí mật nguyên văn; quyền xuất dữ liệu phải tách với quyền xem. |
| SRS-N06.04 | Tra cứu vận hành và tải hộp tin phải được đo theo khối lượng dữ liệu/lịch sử được duyệt; lỗi tra cứu không được suy thành gửi thất bại. |

**Cách kiểm chứng:** Kiểm tra đến hạn lưu/xóa, truy vết sau xóa bí mật, backup/restore và truy vấn trên dữ liệu chuẩn. Tham số SRS-P09/10; DEC-09.

### 10.7. SRS-N07 — Tương thích và khả năng tích hợp module mới

| Mã | Yêu cầu và bằng chứng |
|---|---|
| SRS-N07.01 | Hệ thống phải công bố version và quy tắc tương thích cho API/event; thay đổi phá vỡ hợp đồng cần version/chuyển đổi rõ và consumer tests. |
| SRS-N07.02 | Cấu hình luồng mới phải dùng định danh, scope, mẫu và kết quả dùng chung; lõi không phụ thuộc enum trạng thái Order, session OTP hoặc sổ tài chính của nguồn. |
| SRS-N07.03 | Mỗi tích hợp phải có bộ ví dụ thành công, từ chối, trùng, xung đột, timeout, hủy và một phần cùng contract tests; thêm trường tùy chọn tương thích không được làm consumer cũ lỗi. |
| SRS-N07.04 | Đội phải ghi thời gian/công việc và thay đổi lõi khi thêm luồng; PO dùng bằng chứng này để đánh giá công sức tích hợp. Không cam kết “không cần sửa code” cho mọi module/provider chưa biết. |

**Cách kiểm chứng:** Chạy bộ consumer hiện hành và thử thêm một luồng mẫu theo hợp đồng; đo việc cần bổ sung. DEC-02/07; chỉ tiêu công sức nếu muốn dùng làm gate phải có giá trị riêng.

## 11. Tham số cần có trước nghiệm thu phần liên quan

### 11.1. Bảng tham số định lượng

Ngày 08/10/2026, PO đã chốt giá trị khởi điểm cho cả 12 tham số tại mục 3 [Quyết định PO cho các DEC](<Notitek - QUYẾT ĐỊNH PO CHO CÁC DEC.md>). Đội dùng các giá trị đó để code, cấu hình và kiểm thử. Giá trị nào phụ thuộc bên sở hữu (pháp chế, tài chính, vận hành, Tech Lead nguồn) được ghi ở mục 4 của tài liệu đó và phải được ký trước khi nghiệm thu phần liên quan. Mỗi thay đổi giá trị phải có đơn vị, phạm vi, môi trường/cửa sổ đo, người xác nhận và ngày.

| Mã | Tham số/đầu ra phải điền | Trạng thái | Chủ trì và người phối hợp | Phần bị chặn khi chưa có |
|---|---|---|---|---|
| SRS-P01 | Thời gian chờ phản hồi OTP; hạn bắt đầu/xử lý lần gửi; p95/p99 và tỷ lệ được phép vượt | PO đã chốt giá trị khởi điểm | User + Notification + vận hành | SLA và hợp đồng đồng bộ R1-OTP. |
| SRS-P02 | Hạn gửi từng luồng; thời điểm bắt đầu nếu hẹn; mức bằng chứng phải đạt trong hạn | PO đã chốt giá trị khởi điểm | PO nguồn + Notification | Lịch gửi và nghiệm thu hạn/giá trị tin thật. |
| SRS-P03 | Độ trễ mục tiêu từng đoạn và từng luồng; trung bình/đỉnh/burst, cơ cấu tải | PO đã chốt giá trị khởi điểm | PO nguồn + Tech Lead + vận hành | Nghiệm thu hiệu năng/tải. |
| SRS-P04 | Số người/đích/kênh mỗi yêu cầu; kích thước payload/biến; giới hạn trang/cursor | PO đã chốt giá trị khởi điểm | Tech Lead + consumer + frontend | Schema giới hạn và tải hộp tin. |
| SRS-P05 | Rate-limit nguồn/provider, mức ưu tiên, thời gian chờ và cách hạn chế tải | PO đã chốt giá trị khởi điểm | Vận hành + Tech Lead + PO | Điều phối dưới tải và giới hạn kênh thực. |
| SRS-P06 | Max attempts, khoảng chờ, cửa sổ retry/đối chiếu, điều kiện fallback theo kênh | PO đã chốt giá trị khởi điểm | Notification + PO nguồn + provider/vận hành | Gửi thật có retry/fallback. |
| SRS-P07 | Cửa sổ chống trùng, giữ dấu hủy, giữ ID kết quả; replay ngoài cửa sổ | PO đã chốt giá trị khởi điểm | Nguồn + Notification + chủ dữ liệu | Replay, chuyển đổi và chính sách xóa liên quan. |
| SRS-P08 | Mô hình lỗi, RPO/RTO, mục tiêu khả dụng và kế hoạch phục hồi backlog | PO đã chốt giá trị khởi điểm | Tech Lead + vận hành | Gate phục hồi/khả dụng. |
| SRS-P09 | Thời gian lưu/xóa từng loại dữ liệu; backup, audit, quyền xuất và xóa | PO đã chốt giá trị khởi điểm | Chủ dữ liệu + PO + vận hành | Lưu dữ liệu thật và xuất lịch sử. |
| SRS-P10 | p95/p99 tải danh sách/chi tiết/tra cứu dưới dữ liệu/tải cụ thể | PO đã chốt giá trị khởi điểm | PO UX + Tech Lead + vận hành | Nghiệm thu đáp ứng hộp tin/tra cứu. |
| SRS-P11 | Độ mới căn cứ quyền/người nhận/đích Push; độ trễ nhận thu hồi và xử lý mất nguồn | PO đã chốt giá trị khởi điểm | User/Authorization + chủ thiết bị + Notification | Gửi nền, Push và kiểm tra dữ liệu cập nhật. |
| SRS-P12 | Chủ ngân sách, đơn vị chịu phí, căn cứ cấp phép/hạn mức và xử lý chưa rõ | PO đã chốt giá trị khởi điểm | Chủ ngân sách + PO + Tech Lead | Kênh phát sinh phí. |

### 11.2. Quyết định nghiệp vụ/tích hợp còn cần đầu ra

| Quyết định trong phạm vi | Điều SRS đã xác định | Đầu ra còn cần chốt |
|---|---|---|
| DEC-01 | Đủ 5 kênh; app/client phải xác nhận, không tự tạo app | Client Push thực và app được nghiệm thu. |
| DEC-02 | Nhận/tra/hủy/kết quả có nghĩa rõ, chống trùng và lỗi an toàn | OpenAPI/schema, xác thực, tên sự việc và đường phát User/Order. |
| DEC-03 | Tách tài khoản/công việc, xử lý nền không dùng phiên người nhận | Tập nhận Shop, căn cứ resolver/service và độ mới. |
| DEC-04 | Mục đích tách tính bắt buộc; không suy consent, Marketing ngoài R1 | Danh sách luồng bắt buộc và kênh được phép thật. |
| DEC-05 | OTP cần bằng chứng gửi thật; timeout không kích hoạt retry nền; User sở hữu mã | Consumer User xác nhận mức phản hồi và SRS-P01, ánh xạ receipt/lỗi. |
| DEC-06 | Retry hữu hạn, kiểm tra hạn, đo từng đoạn và tải | Các giá trị SRS-P02 đến SRS-P08/P10 khi liên quan. |
| DEC-07 | Phiên bản, thứ tự ưu tiên rõ, giữ khóa, xung đột không tự chọn | Thứ tự từng chiều và chính sách tin chờ theo luồng. |
| DEC-08 | Không tự vượt phí, không tự làm sổ ngân sách | Bên sở hữu/cấp phép, chính sách SRS-P12. |
| DEC-09 | Tách bí mật/nội dung/metadata/audit, che và xóa | Chính sách SRS-P07/09, quyền xem/xuất/xóa. |
| DEC-10 | Kết quả theo bằng chứng, mục tiêu và một phần tách riêng | Bằng chứng provider thực, kênh bắt buộc/thay thế và bảng chuyển trạng thái máy đọc. |
| DEC-11 | Chuyển đổi có điểm bàn giao, đối chiếu và khôi phục | Kiểm kê ZNS cũ, mẫu/định danh, tập luồng và thời điểm bật thật. |
| DEC-12 | Hành động có quyền riêng; ALLOW đúng tài nguyên | Catalog action/resource/scope, người được cấp và contract tests User/Authorization. |

PO đã chốt kết quả cho cả 12 DEC tại [Quyết định PO cho các DEC](<Notitek - QUYẾT ĐỊNH PO CHO CÁC DEC.md>); cột "Đầu ra còn cần chốt" nay là các artifact kỹ thuật (schema, ADR, contract test) và xác nhận của bên sở hữu tại mục 4 tài liệu đó. Các giá trị/đầu ra này chặn phần tương ứng. Có thể làm hợp đồng, thiết kế và kiểm thử phần không phụ thuộc chúng; không coi toàn bộ SRS chưa có giá trị khi một provider chưa được chọn.

## 12. Truy vết và tiêu chí nghiệm thu

### 12.1. Ma trận theo lát cắt

| Phần triển khai | Nhóm SRS chức năng chính | UC | Story | AT chính |
|---|---|---|---|---|
| Cấu hình, mẫu, kết nối | SRS-F03/04/12/15/20/21 | UC-NTF-04/07/08/15 | ST-01/11 | AT-17/18/24/27 |
| OTP SMS | SRS-F01/02/03/04/05/06/07/09/13/14/15/18/20/21 | UC-USR-01 | ST-02 | AT-01/02/03/04/19/21/22/23/28 |
| Email lời mời | SRS-F01/02/04/05/06/07/09/13/14/15/17/20 | UC-USR-03/04 | ST-03 | AT-02/03/05/06/07/08/21/24/28 |
| Sự việc Order và dừng | SRS-F01/02/03/05/06/09 | UC-ORD-05, UC-NTF-01/03/10 | ST-04 | AT-02/09/10/21/22/23 |
| In-app Shop | SRS-F05/08/09/13 | UC-ORD-05, UC-NTF-13 | ST-05 | AT-08/11/15/21/25/26 |
| Push | SRS-F05/07/09/13/15/16/21 | UC-ORD-05, UC-NTF-08/09 | ST-06 | AT-11/12/24/25/26/27 |
| ZNS/SMS dự phòng | SRS-F05/07/09/13/14/15/18/19/20/21 | UC-ORD-05 | ST-07 | AT-09/11/19/24/25 |
| Cảnh báo thiết bị | SRS-F02/05/08/09/10/13/22 | UC-USR-09 | ST-08 | AT-03/08/13/14/15/21/26/28 |
| Tra cứu/vận hành | SRS-F09/11/12/13/14/15/21 | UC-NTF-12 | ST-09 | AT-03/16/17/18/19/25/27/28 |
| Chuyển đổi | SRS-F02/09/10/22 | UC-USR-09, UC-ORD-05 và luồng cũ được chọn | ST-10 | AT-14/20/23 |
| Hợp đồng, nghiệm thu và đo giá trị | Nhóm liên quan trong SRS-F01 đến SRS-F22; SRS-N01 đến SRS-N07 | Các UC được chọn | ST-11/12/13 | Ca liên quan trong AT-01 đến AT-28 và chỉ số nguồn. |

Chi tiết AT nằm ở [backlog và nghiệm thu](<Notitek - BACKLOG TRIỂN KHAI VÀ KỊCH BẢN NGHIỆM THU.md>). Phạm vi SRS mở rộng năng lực tối thiểu của các story hiện có, không tự bổ sung release Marketing hoặc module tài chính.

### 12.2. Điều kiện chấp nhận một yêu cầu

Mỗi yêu cầu con phải được QA chọn ít nhất ca đạt hoặc ca bị ngăn/lỗi phù hợp; các yêu cầu có quyền, retry, timeout hoặc trạng thái phải có ca biên liên quan. Bằng chứng cần reference giả, scope, phiên bản cấu hình/hợp đồng, kết quả thực tế và mức provider xác nhận. Không cần xuất bí mật để làm bằng chứng.

Một lát cắt chỉ đạt khi:

1. Có hợp đồng nguồn, cấu hình và thông số liên quan đã được xác nhận.
2. Nguồn thật bàn giao đúng dữ liệu; Notification trả đúng mức cam kết.
3. Kênh thật/hiển thị thật đến đúng người và client đã chọn.
4. Hành động mở đúng User/Order; quyền hiện tại được kiểm tra.
5. Trùng, xung đột, timeout, hủy, hết hạn và một phần được kiểm chứng khi liên quan.
6. Tra cứu giải thích được kết quả; không lộ bí mật, không vượt quyền hoặc điều kiện chi.
7. Mục tiêu NFR phần liên quan có giá trị và đạt trên môi trường đo được duyệt.

### 12.3. Điều kiện để xuống code

PO chốt hành vi/ưu tiên; BA chọn yêu cầu con của story. Tech Lead và các consumer/provider chốt schema, version, lỗi, cơ chế quyền, mô hình trạng thái/dữ liệu và ADR liên quan. QA có ví dụ và ca kiểm chứng. Frontend có hành vi màn hình và nguồn dữ liệu. Phần lưu dữ liệu/gửi thật phải có chính sách dữ liệu, phí và kết nối hợp lệ.

SRS chốt yêu cầu quan sát được; việc chọn giải pháp hoặc bật production được thực hiện theo kế hoạch phát hành và bằng chứng tích hợp.

## 13. Các quy tắc được chốt ở mức yêu cầu sản phẩm

| Mã | Quy tắc của baseline |
|---|---|
| SRS-PO01 | Notification là năng lực gửi dùng chung; không sở hữu OTP, trạng thái đơn, membership, quyết định bảo mật hoặc sổ tài chính. |
| SRS-PO02 | Release đầu giữ bốn lát cắt và đủ năm kênh; không bắt mỗi luồng dùng mọi kênh và không mở rộng Marketing. |
| SRS-PO03 | Accepted chỉ sau ghi nhận bền vững; provider nhận/giao/đọc/nghiệp vụ hoàn tất là các mức khác nhau. |
| SRS-PO04 | Replay đúng dữ liệu giữ cùng yêu cầu; cùng khóa dữ liệu xung đột bị từ chối; không tự ghi đè và không retry phần đã đạt. |
| SRS-PO05 | Bí mật không được tự retry nền sau cửa sổ/chỉ dẫn nguồn; timeout không chứng minh thất bại và không tự kích hoạt fallback. |
| SRS-PO06 | Mọi tin có người và ngữ cảnh rõ; liên hệ trùng không hợp nhất identity; quyền hiện tại kiểm tra khi xem/hành động. |
| SRS-PO07 | Trạng thái đọc backend là nguồn sự thật cho In-app mới; đọc không thực hiện nghiệp vụ. Chuyển feed phải có chính sách trạng thái đọc riêng. |
| SRS-PO08 | Kết quả mục tiêu và kết quả từng kênh tách riêng; kênh tùy chọn lỗi vẫn được hiển thị, phần bắt buộc thiếu không bị che bằng success tổng. |
| SRS-PO09 | Cấu hình/kết nối phải được kiểm tra, duyệt, phiên bản hóa và audit ngay cả khi quản trị chưa có UI đầy đủ. |
| SRS-PO10 | Chưa đủ căn cứ quyền/đích an toàn/ngân sách thì không mặc định được phép gửi; thông số chưa xác nhận không được ghi là nghiệm thu đạt. |

Các quy tắc trên được dùng để PO review ở mục 14. Đầu ra DEC và tham số tại mục 11 vẫn phải có bên sở hữu thực tế xác nhận.

## 14. PO review và chốt tài liệu

### 14.1. Phạm vi và người thực hiện

**Vòng review:** 1, ngày 08/10/2026. **Người thực hiện:** Trợ lý được người yêu cầu giao đóng vai PO dự án Notitek. Review này là đánh giá và quyết định baseline tài liệu trong vai trò được giao; không ghi thay chữ ký/xác nhận của PO nguồn, Tech Lead, chủ ngân sách hoặc vận hành.

Review đối chiếu BRD 0.5, phạm vi đã chốt, UC ưu tiên, hợp đồng, UX, thiết kế, backlog và căn cứ mã nguồn ở mục 3. Tiêu chí gồm đúng vấn đề người dùng/module, đúng ranh giới, không mở scope ngầm, yêu cầu rõ và kiểm chứng được, có truy vết, có cách xử lý ngoại lệ và phân biệt điều đã chốt với tham số chưa có căn cứ.

### 14.2. Phát hiện và xử lý trong vòng review

| Mã | Mức | Phát hiện cần xử lý trước khi chốt | Cách xử lý trong bản 1.0 | Kết quả review tài liệu |
|---|---|---|---|---|
| PRV-01 | Cao | Bản 0.1 là nhóm yêu cầu tóm tắt, thiếu đầu vào/ngoại lệ/đầu ra để phát triển và QA dùng | Mục 4–9 có dữ liệu, giao diện, 22 nhóm/136 yêu cầu cụ thể và ca kiểm chứng; giữ mã nhóm cũ | Đã xử lý. |
| PRV-02 | Cao | Chỉ có TTL mã/hạn gửi chưa ngăn việc gửi nền sau khi User bỏ ứng viên vì timeout | SRS-F01.06, SRS-F07.08/09, SRS-N01 và R1-OTP phân biệt cửa sổ đồng bộ, TTL và callback trễ; AT-04/23 | Đã xử lý; giá trị cửa sổ còn cần User xác nhận ở SRS-P01. |
| PRV-03 | Cao | Dễ đưa quyết định OTP/đơn/thiết bị tin cậy/ngân sách sang Notification khi mô tả chi tiết | Mục 2.4, SRS-F05/06/10/14 và bốn lát cắt chỉ yêu cầu căn cứ/chỉ dẫn/hợp đồng từ chủ sở hữu | Đã xử lý; không mở trách nhiệm nguồn. |
| PRV-04 | Cao | Một cờ thành công có thể che người/kênh bắt buộc chưa đạt hoặc nhầm đọc với hoàn tất nghiệp vụ | SRS-F09 và mục 8 tách tiếp nhận, bằng chứng, mục tiêu, đọc; có ví dụ kênh bắt buộc/tùy chọn/thay thế, AT-25 | Đã xử lý; tiêu chí từng luồng thực cần DEC-10. |
| PRV-05 | Cao | Có tên năm kênh nhưng thiếu kết nối thực, sender/mẫu và vòng đời đích Push | SRS-F15 đến SRS-F19; kiểm tra scope/brand/môi trường, đăng ký/thu hồi qua chủ thiết bị, token lỗi, AT-12/24/26 | Đã xử lý; client/provider thật vẫn cần hợp đồng. |
| PRV-06 | Vừa | UC cũ để mở hành vi đồng bộ đọc trong khi SRS mới cần nguồn trạng thái đọc rõ | SRS-F08.05 và UC/UX 0.2 thống nhất backend là nguồn sự thật cho In-app mới, chuyển feed cũ có chính sách riêng; AT-26 | Đã đồng bộ. |
| PRV-07 | Cao | So payload chống trùng có thể tạo chỗ lưu/khôi phục OTP ngoài đường gửi được phép | SRS-F02.07, SRS-F13 và AT-28 yêu cầu bảo vệ dữ liệu đối chiếu, hàng lỗi, trace, xóa và restore | Đã xử lý ở yêu cầu; cơ chế bảo vệ do thiết kế chốt. |
| PRV-08 | Vừa | Các yêu cầu mở rộng cần ca AT và ma trận tương ứng, không chỉ thêm ID trong SRS | Backlog 0.2 bổ sung AT-21 đến AT-28, cập nhật ST-01 đến ST-12; hợp đồng/thiết kế nhận đầu vào SRS 1.0 | Đã đồng bộ; tham chiếu được kiểm tra. |
| PRV-09 | Cao | Ghi “chốt SRS” dễ bị hiểu là SLA, tải, phí, retention và tích hợp đã được các bên duyệt | Mục 10–11 nêu cách đo, 12 nhóm tham số, chủ trì và phần bị chặn; trạng thái đầu tài liệu nêu đúng phạm vi chốt | Đã phân định; không điền số hoặc xác nhận thay bên sở hữu. |
| PRV-10 | Vừa | UX phân biệt không tìm thấy và không có quyền có thể tiết lộ tin của người khác | SRS-F08.08 và UX 0.2 yêu cầu phản hồi an toàn, không tiết lộ tồn tại/nội dung tin ngoài scope; AT-21 | Đã xử lý. |

### 14.3. Kết luận PO và baseline được chốt

**Chấp nhận SRS phiên bản 1.0 làm baseline yêu cầu nghiệp vụ của Notification cho đợt đầu.** Các quy tắc SRS-PO01 đến SRS-PO10 và các yêu cầu chức năng/định tính trong phạm vi đã mô tả được chốt để BA, phát triển và QA phân rã công việc, cụ thể hóa hợp đồng và thiết kế.

Các phát hiện của vòng review đã được xử lý hoặc đưa đúng phần cần bên sở hữu xác nhận ở mục 11. Không mở thêm Marketing, sổ tài chính, app mới hoặc quyền sở hữu nghiệp vụ của User/Order. Bộ vẫn gồm 10 tài liệu; biên bản review được lưu ngay trong SRS.

**Chưa chốt thay các bên:** giá trị SRS-P01 đến SRS-P12, client/provider, tập nhận thật, mẫu/sender, catalog quyền và schema/ADR. Những đầu ra này tiếp tục chặn nghiệm thu phần liên quan; chúng không làm thay đổi các nguyên tắc xử lý đã chốt. Thay đổi hành vi hoặc scope sau baseline phải ghi phiên bản, tác động và cập nhật truy vết; không sửa ngầm.

### 14.4. Bằng chứng rà soát tài liệu

- 22 nhóm chức năng, 136 yêu cầu con có mã duy nhất; 7 nhóm phi chức năng với 25 yêu cầu/tiêu chí kiểm chứng.
- Giữ các mã SRS-F01 đến SRS-F14 và SRS-N01 đến SRS-N07; bổ sung có mã riêng, không dùng lại mã cho nghĩa khác.
- 28 kịch bản AT được định nghĩa trong backlog; các tham chiếu BR, UC, story, SRS và AT hợp lệ.
- Mục lục, tiêu đề và liên kết local được kiểm tra; bảng Markdown có header và số cột thống nhất.
- Đây là bằng chứng chất lượng/tính nhất quán của tài liệu. Chưa có code Notification hoặc bằng chứng thực thi AT/NFR để ghi nhận phần mềm đã đạt.
