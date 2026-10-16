# Notitek - RÀ SOÁT DATABASE v0.7

**Lịch sử:** báo cáo này giữ kết quả v0.7 và quyết định áp dụng v0.8. Thiết kế hiện hành là [v0.9 theo 3NF](<Notitek - MÔ HÌNH DỮ LIỆU LOGIC VÀ TỪ ĐIỂN DỮ LIỆU.md>); tham chiếu mục/số bảng cũ đọc cùng [bản v0.8 lưu lịch sử](<Notitek - DATABASE v0.8 (LỊCH SỬ).md>).

**Ngày:** 13/10/2026. **Đối tượng:** [Mô hình dữ liệu logic và từ điển dữ liệu v0.7](<Notitek - MÔ HÌNH DỮ LIỆU LOGIC VÀ TỪ ĐIỂN DỮ LIỆU.md>).

Đã rà soát **20 bảng và 193 trường**. Mỗi bảng có tình huống giải thích vì sao cần tồn tại riêng; mỗi trường có giá trị mẫu và quy tắc/truy vấn cần giá trị đó. Tất cả tên bảng/trường dùng snake_case; tên hiện hành dài tối đa 20 ký tự. Dòng gộp nhiều trường đã được tách thành từng dòng.

Showcase là dữ liệu giả dễ đọc, không phải dữ liệu production hay chứng nhận PostgreSQL đã chạy. Repo chỉ có migration tạo schema notification, chưa có bảng nghiệp vụ. Bản cập nhật này chỉnh tài liệu; API và code giữ tên hiện hành, được nối với tên DB ở mục 7 của từ điển.

**Cập nhật sau khi người yêu cầu chốt ngày 13/10/2026:** từ điển hiện hành là v0.8, DEC là v1.9. Sáu quyết định DB01–DB06 đã được ghi nhận ở mục 5 dưới đây. Các phần 1–4 giữ kết quả review v0.7; những trạng thái “còn chốt” về TTL selector, inbox, cửa sổ metadata, outbox kết quả và phân quyền đã được thay bằng quyết định mới.

## 1 Phủ showcase từng bảng

| Bảng v0.7 | Số trường có showcase | Tình huống chứng minh cần bảng |
|---|---|---|
| `send_request` | 28 | Order báo đơn ord_42 giao thất bại lần fail_7, Notitek nhận request N1. N1 là gốc chung để tra cứu, chống dùng lại reference và tổng hợp kết quả của Shop/người nhận hàng; nếu chỉ có các dòng gửi kênh thì không biết chúng cùng một yêu cầu hay request đã accepted nhưng chưa PLAN. |
| `request_dedupe` | 9 | Order POST cùng khóa hai lần vì mất response. Dòng (order, SEND, ord_42:delivery_failed:fail_7) trả lại N1; đổi payload nhận xung đột. Tách bảng vì dấu HMAC chỉ giữ 7 ngày, còn lịch sử gửi giữ 180 ngày; lệnh CANCEL cũng cần không gian chống trùng riêng. |
| `cancel_request` | 14 | Order hủy nhóm ord_42:delivery_alert của Shop 17 trước khi một sự việc trễ đến. X1 giữ biên nhận 2 request khớp, 3 delivery hủy, 1 không thu hồi được, 1 chưa rõ; dấu nhóm có hiệu lực 72 giờ. Chỉ đổi state trên request hiện có sẽ không chặn sự việc chưa tới và không trả lại được biên nhận gốc. |
| `send_secret` | 5 | User gửi lời mời N2 có invitation_url chứa token, API accepted rồi process dừng. Worker lấy gói ENC1 để gửi tiếp. Tách khỏi request vì gói phải mã hóa và xóa sau tối đa 10 phút, trong khi metadata N2 cần giữ lâu hơn; OTP đồng bộ không có dòng ở bảng này. |
| `send_case` | 4 | N1 tạo C1/SHOP_ALERT và C2/CONSIGNEE_NOTICE. C1 đạt khi In-app cho thành viên Shop đạt; C2 đạt khi một kênh ZNS hoặc SMS đạt. Tách case để không dùng tiêu chí của người nhận hàng áp cho Shop và để giữ mã nhóm nghiệp vụ ổn định khi PLAN phát lại. |
| `recipient_message` | 7 | C1 resolve U1 và U2 thành M1 và M2; C2 tạo M3 cho người nhận hàng chưa có account_id. Tách message để trạng thái, quyền và các kênh gửi của mỗi người được theo dõi riêng, không gộp mọi người vào một dòng case. |
| `delivery` | 15 | M1 có D1/IN_APP và D2/PUSH đến S1; M3 có D3/ZNS và D4/SMS cùng nhóm ALTERNATIVE. D1 đạt vẫn có thể đi cùng D2 thất bại. Tách delivery để mỗi kênh/đích có hạn gửi, trạng thái và lịch retry riêng; hai đích Push phải là hai dòng khác nhau. |
| `send_attempt` | 13 | D3 gọi provider lần 1 thành A1; nếu adapter xác nhận chắc chắn chưa gửi và policy cho retry, lần 2 là A2. Cùng một delivery có thể dùng kết nối khác khi policy cho phép. Tách attempt để giữ lịch sử từng lần gọi, cấp/hoàn suất đúng một lần và khớp callback về đúng lần gửi. |
| `send_evidence` | 10 | Callback báo zns_msg_701 đã giao đến trước lúc response provider được ghi. E1 được lưu UNMATCHED rồi nối A1 khi khớp; callback yếu hơn đến sau vẫn giữ lịch sử. Không chỉ ghi đè state của delivery, vì sẽ mất bằng chứng, thứ tự nhận và kết quả xử lý xung đột. |
| `inbox_item` | 15 | U1 có I1 gắn D1, U2 có I2 gắn delivery riêng. U1 đọc I1 không đổi trạng thái đọc của U2. Tách inbox để có nội dung hiển thị, thứ tự SSE và trạng thái đọc từng người; đã đọc không đồng nghĩa xử lý xong sự việc Order. |
| `push_target` | 11 | U1 đăng ký hai trình duyệt S1/S2. Đăng xuất S1 thì xóa token_cipher của S1, còn S2 vẫn dùng được và lịch sử D2 vẫn tham chiếu S1. Tách đích Push khỏi tin gửi để quản lý chủ token, thu hồi và chống đăng ký cũ hồi sinh độc lập từng thông báo. |
| `blocked_contact` | 5 | Người nhận yêu cầu dừng SMS cho số kết thúc 042. Trước lần gửi tiếp, Notification tra HMAC3/SMS và chặn mục đích không phải SECURITY. Tách bảng theo liên hệ/kênh để một yêu cầu ngừng nhận áp dụng nhiều request; ai thêm/gỡ tra audit, không giữ SĐT nguyên văn tại đây. |
| `quota_counter` | 5 | Ngày 13/10, bucket (shop_ops, order.delivery_failed, SMS) đã dùng ròng 99 suất, limit 100. Hai worker tranh suất cuối, chỉ một được cấp. Cần dòng đếm cập nhật nguyên tử ở transaction cấp suất; không đếm lại toàn bộ lịch sử attempt ở mỗi lần gửi. |
| `outbox_job` | 12 | N1 và PLAN job J1 commit cùng nhau, rồi API chết trước lúc đẩy RabbitMQ. Relay khác vẫn đọc J1 và phát lại. Tách công việc phát broker để tránh DB có request nhưng broker không có job; payload broker chỉ cần ID/thứ tự, không sao chép bí mật. |
| `audit_log` | 12 | ops_9 đổi limit 100 → 200 rồi ops_10 đổi 200 → 150. Dòng hạn mức hiện hành chỉ giữ 150; L1/L2 giữ ai đổi, khi nào, trước/sau và lý do trong 2 năm. Cần audit độc lập vì lịch sử và quyền bị từ chối không thể dựng lại từ giá trị hiện hành. |
| `channel_connection` | 13 | K1 là OA ZNS của brand supership, đã gửi thử đạt, đang ACTIVE/version 3. Vận hành đổi credential thì chuyển DRAFT và phải thử lại; A1 vẫn ghi K1/version 3. Tách kết nối runtime khỏi registry để đổi tài khoản/sender có quyền và audit mà không deploy. |
| `send_pause` | 3 | Provider SMS gặp sự cố lúc 09:10. Vận hành thêm (order.delivery_failed, SMS), các lần gửi mới bị SUPPRESSED/KILL_SWITCH; xóa dòng để bật lại. Cần trạng thái dừng runtime áp dụng cả tin đang chờ, không phải sửa mẫu hay tạo lại request. |
| `shop_channel_setting` | 4 | Shop 17 tắt Push tùy chọn; Shop 18 vẫn theo mặc định registry. Dòng (shop_17, PUSH, false) là ghi đè của đúng Shop, không được tắt In-app bắt buộc. Tách thiết lập riêng Shop khỏi chính sách chung và kết nối provider. |
| `quota_limit` | 5 | Chủ ngân sách giao mức 100 SMS/ngày cho shop_ops/order.delivery_failed. Vận hành sửa thành 200 qua API có audit; worker đọc mức hiện hành còn counter vẫn là số đã dùng. Tách giới hạn khỏi bộ đếm để chỉnh cấu hình không làm mất usage của ngày đó. |
| `config_snapshot` | 3 | N1 accepted khi cấu hình H1 đang dùng; hôm sau deploy H2 đổi mẫu. Worker vẫn dựng và tính mục tiêu N1 theo nội dung H1. Tách snapshot bất biến vì chỉ lưu hash trong request không đủ để lấy lại mẫu/policy cũ; kết nối, hạn mức và dừng khẩn cấp tiếp tục đọc runtime. |

Mỗi trường ở mục 5 của từ điển trả lời hai câu: **giá trị cụ thể là gì**, và **truy vấn/quy tắc nào sẽ thiếu dữ liệu nếu bỏ trường**. ID nhãn N1/M1/D1, hash H1 và gói ENC1 chỉ dùng giải thích; định dạng vật lý do Tech Lead chốt. Những giá trị suy được vẫn không thêm cột: attempt_count, mục đích/mode/cost_center của request, tỷ lệ sử dụng quota, trạng thái nghiệp vụ Order/User.

## 2 Kết quả đặt tên

Ví dụ tên đã rút gọn: notification_request → send_request; notification_message → recipient_message; channel_delivery → delivery; push_subscription → push_target; transient_secret → send_secret; registry_snapshot → config_snapshot. Những tên đã rõ như account_id, organization_id, quota_counter và channel_connection được giữ. Không viết tắt Organization thành Shop vì hai đối tượng không đồng nghĩa.

Các trường chung chung được làm rõ: type → notice_type; data → template_data; kind ở outbox → job_type; attempts ở outbox → publish_count. Các trường dài được rút gọn mà giữ ý nghĩa: reconcile_deadline_at → reconcile_until; not_recallable_deliveries → unrecallable_count. Không gọi trường sau là sent_count vì nó biểu diễn kết quả không thể thu hồi.

Bảng ánh xạ đầy đủ nằm tại mục 7 của từ điển. Enum/header không đổi: KILL_SWITCH, ALTERNATIVE, Idempotency-Key vẫn theo hợp đồng API. Không tạo thêm các bảng đồng nghĩa với tên cũ trong tài liệu cấu trúc mã nguồn.

## 3 Tiếp nối review v0.6

| Điểm review cũ | Kết quả trên v0.7 | Còn cần hoàn thiện |
|---|---|---|
| R01 - Nhóm ZNS/SMS | Sửa ví dụ: cả hai ALTERNATIVE cùng nhóm; policy quyết định ZNS chính/SMS dự phòng. Dòng 9 không kích hoạt mọi delivery thay thế. | Chốt cách kết thúc delivery dự phòng chưa cần dùng khi ZNS đạt; chưa gọi toàn bộ điều phối fallback là đã kín. |
| R02 - Phản hồi OTP | Thêm send_request.sync_result bất biến và updated_at. Showcase: sync_result UNKNOWN nhưng callback sau làm outcome ACHIEVED. | Hiện thực transaction ghim kết quả/phục hồi và ánh xạ API; chưa có test DB. |
| R03 - TTL bí mật | Thêm send_secret.stored_at, TTL tính từ gói hiện tại. Gói mới có TTL mới, không gia hạn gói cũ; không vượt hạn nghiệp vụ nếu có. | Kết quả API cấp lại do User cung cấp theo hợp đồng; Notification không sở hữu nghiệp vụ lời mời. |
| R04 - Dữ liệu liên hệ | Làm rõ template_data không chứa liên hệ đầy đủ/bí mật. Showcase của recipients_spec chỉ rõ dữ liệu cần cho PLAN. | Dọn selector khi PLAN lỗi/đóng trước PLAN và TTL tối đa độc lập còn Q11; chưa coi là đã đóng. |
| R05 - Callback/retry/revision | Revision/updated_at phải đổi với mọi thay đổi API quan sát được, kể cả state giữ nguyên. | Callback trước Tx2, retry đến hạn và anchor UNKNOWN khi worker chết còn Q11. |
| R06 - Retention | Gỡ khẳng định 180 ngày là tối thiểu đã chốt; chỉ xóa đích Push/kết nối khi đủ tuổi giữ và không còn tham chiếu cần giữ. | Mốc cuối inbox/evidence, nội dung title/preview sau 90 ngày và khoảng giữ theo partition còn Q9/Q11. |
| R07 - Tần suất nghiệp vụ | Q2 không còn yêu cầu Notitek đếm 5 OTP/số/giờ hay 3 tin/số/đơn/ngày; User/Order sở hữu theo DEC v1.8. | Giới hạn bảo vệ admission/provider phải có scope riêng; quota chi phí vẫn thuộc Notification. |
| R08 - Reference/planning/outbox | Giữ tên API/code hiện hành và ánh xạ database minh bạch. | Đồng bộ REFERENCE_RETIRED sau 7 ngày, không tạo enum PLANNED, chốt có/không lưu event kết quả ở Q11. |
| R09 - Quota qua nửa đêm | Định nghĩa used_count là suất cấp ròng. Hoàn bucket theo send_started_at và cấu hình request đã ghim, không theo ngày callback. | Kiểm chứng hoàn đúng một lần và cấp suất đồng thời khi có schema; không thêm sổ tài chính/decision ledger. |
| R10 - Phụ thuộc đã biết | Tiếp tục Q10; showcase giải thích giới hạn theo provider và Push, không coi là lỗi chặn toàn bộ database. | Scope ID callback, lịch sử credential, nguồn registration_version và xoay HMAC phải chốt trước khi bật năng lực tương ứng. |

## 4 Kiểm tra tài liệu

Đối chiếu trực tiếp 20 bảng trong ERD với 20 mục từ điển; kiểm tên snake_case, không trùng trường trong một bảng, mỗi dòng chỉ một trường, đủ mô tả/điều kiện/giá trị mẫu/cách dùng. Số trường ban đầu sau khi tách dòng gộp là 190; bổ sung ba trường sync_result, updated_at và stored_at thành 193. Tên mới dài nhất 20 ký tự.

Thiết kế không tự thêm bản sao User/Order, bảng OTP, session, membership hoặc tài chính. Physical DB còn phải chọn kiểu PostgreSQL/index/partition và hiện thực unique có NULL, transaction, callback và job retention. Checklist này là kiểm tra tính đầy đủ của tài liệu, không thay kiểm thử thực thi.

## 5 Quyết định đã chốt và trạng thái tiếp theo

Người yêu cầu xác nhận chốt cả 6 đề xuất. Chính sách được đồng bộ trong từ điển v0.8, DEC v1.9 mục 2.16, SRS 1.5 và các tài liệu API liên quan.

| Mã | Nội dung | Quyết định PO đã chốt | Showcase |
|---|---|---|---|
| DB01 | Metadata gửi | Giữ tối thiểu 180 ngày và tối đa 180 ngày + 1 tháng để xóa theo partition tháng. Nội dung/liên hệ vẫn dọn đúng hạn riêng, không được giữ thêm theo partition. | N1 có thể còn metadata ở ngày 190, nhưng liên hệ nguyên văn đã che ở ngày 30. |
| DB02 | Mốc xóa inbox | Metadata inbox xóa cùng cây request, tính tuổi từ accepted_at của request; created_at của inbox chỉ tính thời gian hiển thị/nội dung 90 ngày. | Request nhận ngày 1, tin tạo ngày 2: xóa cùng request; inbox không có thêm một ngày lưu riêng. |
| DB03 | Nội dung inbox | title và preview giữ tối đa 90 ngày từ created_at của inbox, sau đó làm NULL. Giữ ID, trạng thái đọc và metadata tới lúc xóa cùng request. | Sau 90 ngày không xem lại nội dung tin, nhưng vẫn đối chiếu được ID và thời điểm đọc. |
| DB04 | Liên hệ tạm chờ lập kế hoạch | Xóa SĐT/email khỏi recipients_spec ngay khi PLAN xong hoặc request đóng; dù job kẹt cũng không giữ quá min(accepted_at + 24 giờ, send_deadline). TTL này khác contact_value đã chuyển vào delivery. | Resolver mất cả ngày: request vẫn còn metadata/kết quả nhưng không còn SĐT/email tạm; không tự đọc lại payload để hồi sinh liên hệ. |
| DB05 | Kết quả R1 | Nguồn chỉ GET theo revision; không ghi notification.result_changed.v1 vào outbox R1. outbox_job chỉ có PLAN/SEND; sự kiện kết quả thiết kế khi làm phân kỳ broker. | Callback nâng UNKNOWN lên ACHIEVED: tăng revision/updated_at; Order polling GET, không có row event kết quả chờ phát. |
| DB06 | Quyền vận hành | Platform admin quản lý kết nối/hạn mức. Vận hành Notification được dừng/bật gửi và đổi thiết lập Push của Shop. Tách quyền theo hành động; mọi thay đổi có audit. | Vận hành dừng SMS mà không nhận quyền sửa credential hoặc nâng hạn mức. |

**Đã đóng về chính sách:** Q7/Q9; R04 về thời hạn/dọn selector; R06 về cửa sổ và mốc inbox/nội dung; R08 về không lưu sẵn outbox kết quả. title/preview đổi từ Y sang C: bắt buộc khi còn hiển thị, NULL sau 90 ngày. Số bảng/trường và showcase vẫn là 20/193.

**Còn phần thực hiện:** job TTL/retention, ràng buộc không mồ côi, callback sớm/retry/anchor UNKNOWN, giao dịch chống trùng/ref hết hạn và triển khai quyền. CF-02 còn cost center/hạn mức cụ thể; CF-07 còn đăng ký/gán quyền; Q10 vẫn theo hợp đồng User/provider. Chưa có migration bảng nghiệp vụ hoặc nghiệm thu các job này.

**Kiểm chứng cần làm khi triển khai:** PLAN kẹt vẫn dọn liên hệ ở hạn cứng; đọc inbox ngày 91 không trả/dựng lại nội dung; request nhận và inbox tạo khác ngày vẫn xóa cùng cây; metadata ngày 190 có thể còn nhưng đúng giới hạn; quyền pause không sửa credential/quota; đổi revision không sinh event kết quả.
