# SuperShip - MODULE NOTIFICATION - ĐỀ XUẤT CÔNG NGHỆ SỬ DỤNG CHO DỰ ÁN

**Phiên bản:** 1.1. **Ngày:** 08/10/2026. **Trạng thái:** PO đã chốt ngôn ngữ (3.1), messaging (3.2), lưu trữ (3.3, 3.4), provider adapter (3.5) và In-app realtime (3.6). Broker sự kiện liên module chờ SuperPlatform; triển khai production và điều kiện đội ngũ Go được theo dõi tại mục Rủi ro.

## Mục tiêu và nguyên tắc chọn stack

Báo cáo này chọn công nghệ dựa trên **bản chất bài toán**. Nguyên tắc xuyên suốt:

- Mỗi lớp bài toán có một đặc tính kỹ thuật riêng (ACID, throughput, độ trễ, khối lượng dữ liệu) → chọn công nghệ khớp với đặc tính đó, không dùng một công nghệ gánh tất cả.

- Bắt đầu từ giải pháp **đơn giản, ít rủi ro vận hành nhất** đáp ứng đúng quy mô hiện tại; chỉ scale-out/đổi công nghệ khi có bằng chứng bottleneck cụ thể, không scale "phòng xa" quá sớm.

- Vì hệ thống chưa cố định quy mô, báo cáo đưa ra **lộ trình theo ngưỡng traffic** thay vì một con số cứng.

## Phân lớp bài toán

Hệ thống Notification được chia thành 7 lớp, mỗi lớp có yêu cầu kỹ thuật khác nhau:

|**STT**|**Lớp bài toán**|**Đặc tính kỹ thuật chính**|
|---|---|---|
|1|Tiếp nhận & xử lý nghiệp vụ (API, Workflow, Template)|Cần ACID, JOIN, ràng buộc dữ liệu rõ ràng|
|2|Hàng đợi & xử lý bất đồng bộ|Priority, retry/backoff, DLQ, độ trễ thấp|
|3|Lưu trữ giao dịch (Delivery hiện tại)|Transaction, idempotency, volume vừa phải|
|4|Lưu trữ lịch sử/audit (Delivery Event)|Append-only, volume rất lớn, truy vấn theo key|
|5|Kênh gửi & Provider Adapter|Isolation theo channel, chuẩn hóa đa nhà cung cấp|
|6|Real-time In-App|Độ trễ thấp, kết nối bền (persistent connection)|
|7|Vận hành & quan sát hệ thống|Truy vết xuyên nhiều tiến trình, cảnh báo sớm|

## Lựa chọn theo từng lớp

### 3.1. Ngôn ngữ & Backend framework

|**Lựa chọn**|**Ưu điểm**|**Nhược điểm**|
|---|---|---|
|**Go**|Concurrency tự nhiên (goroutine) khớp mô hình "mỗi channel một pool worker riêng"; binary gọn, khởi động nhanh, scale ngang nhanh khi traffic tăng đột biến (campaign gửi hàng loạt)|Hệ sinh thái ORM/ràng buộc nghiệp vụ phức tạp không mạnh bằng Java/Node|
|**Node.js (NestJS/TypeScript)**|Module hóa rõ ràng, khớp tự nhiên với các thành phần nghiệp vụ (Bộ định tuyến, Bộ tạo nội dung...); hệ sinh thái queue library phong phú|Single-thread, cần cẩn trọng khi CPU-bound (render template phức tạp)|
|**Java/Kotlin (Spring Boot)**|Hệ sinh thái enterprise mạnh, transaction phức tạp, đội ngũ quen thuộc nhiều|Nặng, khởi động chậm hơn khi cần scale nhanh theo tải đột biến|

**Chọn: Go.** Bài toán là I/O-bound (chờ Provider phản hồi) và cần isolation + scale độc lập theo từng channel — đúng thế mạnh của goroutine. Java Spring Boot là lựa chọn thay thế hợp lý nếu đội ngũ đã có sẵn năng lực Java và ưu tiên hệ sinh thái enterprise hơn tốc độ scale.

**Quyết định đã chốt với PO:** Notitek viết bằng **Go**. Go cũng phù hợp với SSE ở mục 3.6 vì mỗi kết nối giữ lâu chỉ tốn một goroutine vài KB.

User và Order đang dùng Java/Spring Boot, nên Notitek phải tự hiện thực bằng Go và giữ tương thích các thành phần dùng chung sau:

| Thành phần | Yêu cầu tương thích |
|---|---|
| Xác thực service-to-service và Access Context | Cùng cơ chế OAuth2 client credentials và cùng quyết định ALLOW/DENY/NOT_EVALUATED của User/Authorization |
| Định dạng phản hồi API | Cùng cấu trúc `ApiResponse` và catalog mã lỗi của SuperPlatform |
| Cursor pagination | Cùng định dạng cursor và quy ước phân trang với Order |
| Quan sát hệ thống | Cùng tên metric, nhãn, trace context (W3C) và định dạng log |
| Build và triển khai | Pipeline, base image và quy trình vá bảo mật cho Go được DevOps hỗ trợ |

Điều kiện đội ngũ và DevOps được theo dõi tại mục Rủi ro.

### 3.2. Hàng đợi và tích hợp bất đồng bộ với SuperPlatform

**Quyết định đã chốt với PO:** chọn **RabbitMQ cho hàng đợi công việc gửi của Notitek**. Broker sự kiện giữa các module tuân theo kiến trúc chung của SuperPlatform; **tích hợp Kafka nếu nền tảng chọn Kafka**. Không bổ sung một cụm Kafka riêng cho Notitek trong R1 chỉ để nhận sự kiện từ module khác.

Quyết định này chốt phần messaging; các lựa chọn công nghệ khác được xem xét riêng.

#### Phân vai công nghệ

| Thành phần | Lựa chọn | Trách nhiệm |
|---|---|---|
| Tích hợp sự kiện liên module | Theo chuẩn SuperPlatform; Kafka nếu nền tảng lựa chọn | Nhận sự kiện từ Order, User và các module khác; công bố kết quả khi hợp đồng yêu cầu |
| Hàng đợi công việc gửi | RabbitMQ | Phân phối công việc theo kênh/provider; hỗ trợ acknowledgement, ưu tiên và dead-lettering |
| Dữ liệu xử lý | PostgreSQL | Lưu yêu cầu, khóa chống trùng, trạng thái, outbox và hộp tin để tra cứu, phục hồi |
| OTP đồng bộ | API theo hợp đồng với User | Trả kết quả trong tối đa 5 giây; không đưa vào hàng đợi gửi nền chung |

RabbitMQ phù hợp để phân phối công việc cho worker và cô lập lỗi theo kênh. Notification chịu trách nhiệm triển khai retry/backoff, hạn gửi, giới hạn tốc độ và fallback theo nghiệp vụ; broker không tự quyết định các chính sách này.

**Công việc gửi trễ** (giữ tin đến 07:00 theo khung giờ, retry sau 1 giờ hoặc 3 giờ theo DEC-06) không được giữ dài hạn trong RabbitMQ. Thời điểm gửi dự kiến được lưu trong PostgreSQL; scheduler quét công việc đến hạn và chuyển sang RabbitMQ qua outbox. Delay ngắn dưới 15 phút có thể dùng cơ chế delay của RabbitMQ nếu Tech Lead chọn.

Kafka và RabbitMQ có thể cùng xuất hiện trong một luồng với hai vai trò khác nhau. RabbitMQ cũng có Streams hỗ trợ đọc lại; cách phân vai này là quyết định kiến trúc của phương án, không phải giới hạn tuyệt đối của sản phẩm.

#### Các luồng sử dụng

1. **Sự kiện nghiệp vụ:** module nguồn lưu nghiệp vụ và outbox trong cùng transaction, hoặc dùng cơ chế tương đương bảo đảm phát sự kiện sau commit. Sự kiện đi qua broker SuperPlatform đến Notification. Notification lưu yêu cầu cùng outbox công việc vào PostgreSQL, sau đó chuyển công việc sang RabbitMQ để worker xử lý.
2. **API gửi bất đồng bộ:** Notification lưu bền vững yêu cầu cùng outbox trước khi trả accepted, rồi chuyển công việc sang RabbitMQ. Luồng này không bắt buộc đi qua Kafka.
3. **OTP đồng bộ:** User gọi API và nhận `PROVIDER_ACCEPTED`, `FAILED` hoặc `UNKNOWN` trong cửa sổ tối đa 5 giây. Tuân thủ `sendDeadline` và `expiresAt`; không retry nền. Chỉ thử provider thay thế khi chắc chắn lần trước chưa gửi, còn trong cửa sổ đồng bộ và hợp đồng OTP cho phép. User sở hữu việc phát hành, kích hoạt và kiểm tra mã.

#### Trách nhiệm và bảo đảm giao nhận

- **SuperPlatform:** thống nhất broker sự kiện, hợp đồng và phiên bản sự kiện, quyền truy cập, retention/replay, truy vết và đầu mối vận hành hạ tầng dùng chung.
- **Module nguồn:** phát sự kiện phản ánh nghiệp vụ đã commit, cung cấp định danh ổn định để chống trùng. Module nguồn không cần biết cấu trúc queue hoặc provider của Notification.
- **Notification:** chỉ xác nhận nhận sự kiện sau khi đã lưu bền vững thông tin cần để tiếp tục xử lý. Với Kafka, commit offset theo tiến độ đã lưu bền vững của từng partition, không bỏ qua sự kiện chưa được lưu. PostgreSQL và broker không tự có transaction chung; cần outbox và cơ chế chuyển tiếp có thể phục hồi.
- **Chuyển tiếp và worker:** dùng publisher confirms, xử lý trường hợp không định tuyến được, xác định thời điểm consumer acknowledgement và chống xử lý trùng. Message có thể được giao lại; broker hoặc ACID của database không tự bảo đảm chỉ gửi đến provider đúng một lần.
- **Replay và retry:** kiểm tra chống trùng, hạn gửi, trạng thái hủy và kết quả hiện có trước khi tạo hoặc thực hiện công việc. Đọc lại sự kiện không mặc nhiên gửi lại thông báo. Provider timeout phải được xử lý theo trạng thái chưa rõ kết quả, không tự xem là chắc chắn chưa gửi.
- **Vận hành:** đầu mối nền tảng phụ trách hạ tầng sự kiện; nhóm Notification sở hữu cấu hình queue, worker và chính sách xử lý công việc. Tech Lead phối hợp vận hành xác định topology, phiên bản và HA của RabbitMQ theo yêu cầu RPO/RTO/SLO ngay từ R1; không đợi đạt ngưỡng số thông báo/ngày mới xem xét khả năng chịu lỗi.

**Phần còn phụ thuộc nền tảng:** xác nhận broker sự kiện mà SuperPlatform lựa chọn và hợp đồng kết nối tương ứng. Quyết định này không đồng nghĩa đã chọn Kafka cho toàn SuperPlatform.

Tham khảo: [Kafka — Introduction](https://kafka.apache.org/intro/), [RabbitMQ — Queues](https://www.rabbitmq.com/docs/queues), [RabbitMQ — Reliability](https://www.rabbitmq.com/docs/reliability), [RabbitMQ — Streams](https://www.rabbitmq.com/docs/streams).

### 3.3. Lưu trữ giao dịch (Delivery hiện tại, Idempotency, Workflow/Template config)

|**Lựa chọn**|**Ưu điểm**|**Nhược điểm**|
|---|---|---|
|**PostgreSQL**|ACID đầy đủ, transaction đảm bảo "ghi Delivery trước khi enqueue" không lệch trạng thái; JOIN tốt cho quan hệ Workflow→Step→Template→Provider; hỗ trợ JSONB cho payload biến động|Scale ghi ngang khó hơn NoSQL ở volume cực lớn|
|**MySQL**|Tương đương Postgres về ACID, phổ biến|JSONB/kiểu dữ liệu phong phú kém hơn Postgres|
|**MongoDB**|Schema linh hoạt theo từng channel|Thiếu transaction đa bảng mạnh, dễ race condition khi check idempotency, khó đảm bảo ràng buộc Workflow/Template|

**Chọn: PostgreSQL.** Lớp này cần ACID để tránh mất job hoặc gửi trùng OTP — đây là điểm không thể đánh đổi lấy schema linh hoạt.

### 3.4. Lưu trữ lịch sử/audit (Delivery Event)

Đây là bảng phát triển nhanh nhất hệ thống (mỗi lần đổi trạng thái tạo 1 record, append-only, gần như không update).

|**Lựa chọn**|**Ưu điểm**|**Nhược điểm**|
|---|---|---|
|**PostgreSQL (partition theo thời gian)**|Không cần thêm hệ thống mới, vẫn JOIN được với dữ liệu giao dịch khi cần tra soát|Ở khối lượng rất lớn (tỷ record/năm), vận hành index/VACUUM nặng dần|
|**Cassandra/ScyllaDB**|Scale ghi ngang tuyến tính theo node, tối ưu cho append-only + truy vấn theo partition key|Mất khả năng JOIN, cần đội vận hành riêng, chi phí ops cao hơn|
|**DynamoDB**|Quản lý vận hành (managed), scale tự động|Chi phí tăng theo throughput, vendor lock-in|
|**ClickHouse**|Truy vấn phân tích/thống kê (tỷ lệ FAILED theo Provider theo ngày) cực nhanh|Không tối ưu cho truy vấn theo 1 bản ghi đơn lẻ như tra soát 1 delivery cụ thể|

**Chọn: PostgreSQL (partition theo tháng) làm mặc định, chuyển sang Cassandra/DynamoDB khi vượt ngưỡng cụ thể** — xem Mục 4 (Lộ trình theo quy mô) để biết ngưỡng chuyển đổi.

### 3.5. Kênh gửi & Provider Adapter

**Chọn: Adapter/Strategy pattern** — mỗi Provider (VietGuys, INCOM, Firebase FCM, Zalo ZNS...) implement một interface chuẩn, đăng ký qua registry. Đây không phải lựa chọn công nghệ mà là nguyên tắc thiết kế bắt buộc, đảm bảo lõi hệ thống không phụ thuộc SDK riêng của từng Provider, dễ thêm/thay Provider mà không sửa logic core.

### 3.6. Real-time In-App

|**Lựa chọn**|**Ưu điểm**|**Nhược điểm**|
|---|---|---|
|**Polling + Web Push**|Đơn giản nhất, không giữ kết nối|Tin hiện chậm theo chu kỳ hỏi; tải server tăng theo số người mở trang|
|**SSE (Server-Sent Events) + Redis Pub/Sub**|Server đẩy một chiều qua HTTP; trình duyệt tự kết nối lại bằng `Last-Event-ID`; Redis phát tín hiệu chéo instance với throughput cao|Cần gateway hỗ trợ kết nối giữ lâu; thêm Redis cần HA|
|**SSE + PostgreSQL LISTEN/NOTIFY**|Không thêm hạ tầng; tín hiệu gắn với transaction|Throughput tín hiệu thấp hơn Redis; mỗi instance giữ một kết nối DB riêng|
|**Redis Pub/Sub tự dựng WebSocket**|Hai chiều|Phải tự xử lý reconnect, presence, scale nhiều instance; R1 không cần kênh hai chiều|
|**Centrifugo (real-time gateway chuyên dụng)**|Xử lý sẵn reconnect, presence, scale ngang|Thêm một thành phần hạ tầng cần vận hành|

**Quyết định đã chốt với PO: SSE + Redis Pub/Sub cho R1.** Không dùng WebSocket trong R1.

| Thành phần | Vai trò |
|---|---|
| SSE | Đẩy tín hiệu tin mới và số chưa đọc tới chuông Shop khi tab đang mở |
| Redis Pub/Sub | Phát tín hiệu tới mọi instance Notitek để instance giữ kết nối SSE của người nhận đẩy xuống |
| PostgreSQL | Nguồn sự thật của hộp tin và trạng thái đọc; SSE chỉ là tín hiệu để frontend tải lại |
| Polling | Dự phòng khi không mở được SSE (mạng chặn, proxy cắt kết nối) |
| Web Push (FCM) | Báo tin khi tab đã đóng; là kênh Push trong 5 kênh, không bị SSE thay thế |

Quy tắc bắt buộc:

- Chỉ publish lên Redis **sau khi transaction lưu bản In-app đã commit**. Redis không tham gia transaction của PostgreSQL; mất tín hiệu do crash giữa commit và publish được chấp nhận vì frontend tải bù qua `Last-Event-ID` và polling.
- Payload tín hiệu chỉ chứa `inboxItemId`, `userId` và ngữ cảnh; không chứa nội dung tin hay dữ liệu cá nhân.
- Một event nguồn tạo một tín hiệu gộp danh sách người nhận, không publish riêng từng người.
- Khi kết nối lại, server trả các bản tin sau `Last-Event-ID` theo quyền và ngữ cảnh hiện tại (SRS-F08.04).
- Xác thực SSE dùng cookie phiên hoặc token ngắn hạn vì `EventSource` không gửi được header `Authorization`; cơ chế cụ thể thống nhất với Shop FE.
- Gateway/proxy tắt buffering cho endpoint SSE và đặt idle timeout lớn hơn chu kỳ heartbeat 25 giây.
- Redis chạy HA (Sentinel hoặc managed) ngay từ R1. Redis lỗi không làm mất tin, chỉ làm chậm hiển thị đến chu kỳ polling.

Chuyển sang WebSocket hoặc Centrifugo khi có app mobile giữ kết nối hoặc nhu cầu tương tác hai chiều.

### 3.7. Vận hành & quan sát hệ thống

**Chọn: OpenTelemetry (tracing) + Prometheus/Grafana (metrics) + Loki (log).** Vì luồng xử lý xuyên qua nhiều tiến trình độc lập (API → Queue → Worker → Provider → Callback), bắt buộc phải có distributed tracing để xác định một request "kẹt" ở khâu nào — đây không phải lựa chọn tùy chọn mà là điều kiện cần cho hệ thống nhiều tiến trình bất đồng bộ.

## Lộ trình theo quy mô (Scaling Path)

Không có một stack "đúng cho mọi quy mô" — bảng dưới đây chốt rõ khi nào cần thay đổi:

|**Quy mô (notification/ngày)**|**Thay đổi so với baseline**|**Lý do**|
|---|---|---|
|**≤ 1 triệu (R1)**|PostgreSQL cho cả giao dịch lẫn lịch sử (có partition theo tháng); RabbitMQ cluster 3 node với quorum queue; Redis HA (Sentinel hoặc managed)|Peak ~200-400 req/giây, PostgreSQL tuning tốt xử lý dư sức; thêm NoSQL ở mức này là over-engineering. HA bắt buộc từ R1 theo SRS-P08 (khả dụng 99,9%), không phụ thuộc sản lượng|
|**1 – 10 triệu**|Thêm read replica PostgreSQL cho tra soát/dashboard; mở rộng số node RabbitMQ/Redis theo số liệu đo|Đọc lịch sử bắt đầu cạnh tranh tài nguyên với ghi giao dịch|
|**≥ 10 triệu**|Tách Delivery Event sang Cassandra/ScyllaDB hoặc DynamoDB; PostgreSQL chỉ giữ dữ liệu giao dịch (nhỏ, cần ACID)|~3.65 tỷ record/năm vượt khả năng vận hành hiệu quả của single-node Postgres (index, VACUUM, backup); đây là lúc "NoSQL cho notification" thực sự đúng — nhưng chỉ đúng ở tầng event log, không phải toàn hệ thống|

## Kết luận

Techstack không nên chọn theo một xu hướng chung mà theo **thực tế của từng lớp bài toán ở từng ngưỡng quy mô cụ thể**. Techstack khuyến nghị cho Module Notification, bao gồm:

|**Nhóm**|**Công nghệ chốt**|
|---|---|
|**Ngôn ngữ / Backend**|Go — đã chốt|
|**Hàng đợi công việc gửi của Notitek**|RabbitMQ — đã chốt|
|**Broker sự kiện liên module**|Theo chuẩn SuperPlatform; tích hợp Kafka nếu nền tảng chọn Kafka; không dựng Kafka riêng cho Notitek trong R1|
|**OTP đồng bộ**|API theo hợp đồng User; tối đa 5 giây, không retry nền|
|**Database giao dịch**|PostgreSQL|
|**Database lịch sử/audit**|PostgreSQL (partition theo tháng) → Cassandra/DynamoDB khi ≥ 10 triệu/ngày|
|**Cache / Pub-Sub**|Redis HA — phát tín hiệu realtime chéo instance|
|**Real-time In-App**|SSE + Redis Pub/Sub, polling dự phòng — đã chốt; WebSocket/Centrifugo khi có app mobile hoặc nhu cầu hai chiều|
|**Báo tin khi tab đóng**|Web Push qua FCM (DEC-01)|
|**Provider Adapter**|Adapter/Strategy pattern (tự thiết kế, không phải sản phẩm cụ thể)|
|**Hạ tầng / Triển khai**|Ở môi trường dev, sử dụng Docker Compose; production chờ Tech Lead và DevOps chốt|

## Rủi ro và điều kiện theo dõi

| Mã | Rủi ro | Điều kiện/cách kiểm soát | Bên chịu trách nhiệm | Hạn |
|---|---|---|---|---|
| TR-01 | Đội thiếu kinh nghiệm Go production làm chậm M1 | Có ít nhất 2 dev đã chạy Go production hoặc kế hoạch nhân sự/đào tạo được duyệt | Tech Lead, quản lý dự án | Trước khi bắt đầu M0 |
| TR-02 | DevOps chưa hỗ trợ stack Go | Pipeline CI/CD, base image, quét lỗ hổng và quy trình vá cho Go được DevOps xác nhận | DevOps | Trước khi bắt đầu M0 |
| TR-03 | Lệch chuẩn với các service Java | Danh sách thành phần dùng chung ở mục 3.1 có ước lượng công và contract test với User/Order. ADR-01 đã chốt tự xây (09/10/2026) | Tech Lead | Trước M0 |
| TR-04 | Gateway/proxy cắt kết nối SSE | DevOps xác nhận tắt buffering và idle timeout > 25 giây cho endpoint SSE; polling dự phòng hoạt động | DevOps, Shop FE | Trước M2 |
| TR-05 | Thêm hai hạ tầng mới (RabbitMQ, Redis) chưa có trên nền tảng | Vận hành nhận sở hữu, có HA, giám sát và runbook từ R1 | Vận hành, Tech Lead | Trước khi bật gửi thật |
| TR-06 | Broker sự kiện SuperPlatform chưa được chọn | Chốt broker và hợp đồng kết nối trước M2; M1 dùng REST nên không bị chặn | Đầu mối nền tảng | Trước M2 |
| TR-07 | Chưa có thiết kế triển khai production | Số instance, môi trường, secret management và RPO/RTO theo SRS-P08 | Tech Lead, DevOps | Trước M1 nghiệm thu |

