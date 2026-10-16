# Notitek - THIẾT KẾ DATABASE 3NF VÀ AN TOÀN MIGRATION

**Phiên bản:** 0.9. **Ngày:** 13/10/2026. **Phạm vi:** thiết kế logic và kế hoạch triển khai; chưa chạy migration bảng nghiệp vụ.

Yêu cầu là thiết kế lại theo 3NF và tránh mất dữ liệu ở phase đầu. Repository hiện chỉ có migration 00001_create_notification_schema.sql tạo schema notification. Vì vậy đường triển khai chính là khởi tạo schema nghiệp vụ mới, không dựng một dự án chuyển dữ liệu legacy không tồn tại. Nếu kiểm kê môi trường triển khai phát hiện dữ liệu thử/ngoài repo, phải dùng đường chuyển đổi có đối chiếu ở mục 3.

[Từ điển v0.9](<Notitek - MÔ HÌNH DỮ LIỆU LOGIC VÀ TỪ ĐIỂN DỮ LIỆU.md>) có 56 bảng/260 trường với showcase và 63 FK. [Bản v0.8](<Notitek - DATABASE v0.8 (LỊCH SỬ).md>) giữ nguyên để kiểm ma trận 193 trường cũ. [Manifest](db/notitek_3nf.schema.json) là dữ liệu thiết kế máy đọc, không phải DDL đã áp.

## 1 Căn cứ 3NF và phép tách

Xét mọi FD không tầm thường X → A: X là siêu khóa hoặc A là thuộc tính prime. Thêm UUID không tự sửa phụ thuộc bộ phận/bắc cầu. Khóa/FD được công bố ở mục 3 từ điển; kết luận 3NF có điều kiện theo FD nghiệp vụ đã biết, không phải chứng minh tự động từ sample.

| Trước | Phép tách | Giao/khóa bảo đảm tái dựng | Bảo vệ ý nghĩa lịch sử |
|---|---|---|---|
| Request lặp type/version/hash/context | profile/schema/type/scope/context | request.profile_id → profile; profile.schema_id/scope_id → parent | Profile, schema, template và scope ghim bất biến; không đọc YAML mới |
| Delivery lặp kênh/template/nhóm/tiêu chí | rule → group/template → channel | delivery.rule_id → rule; rule FK tới group/template | Constraint bảo đảm cùng case/hash/scope |
| Attempt lặp provider, kết nối runtime bị sửa | connection identity + revision + secret | attempt.(connection_id, connection_version) → revision | Version sender bất biến, credential tách TTL |
| Evidence bị cập nhật để gắn attempt | evidence + evidence_link | evidence_id → tối đa một link; link.attempt_id → attempt | Receipt chưa có link vẫn tồn tại; apply idempotent |
| JSON tập người/biến/header/audit | bảng child có PK của từng phần tử | request_id/spec_id/outbox_id/audit_id + tên/key phần tử | View gom phần tử về JSON hợp đồng, không nhét secret vào projection |
| Nội dung và metadata TTL khác nhau | delivery_text/inbox_text/spec_contact | PK child là FK parent | LEFT JOIN sau hết TTL trả NULL, không xóa metadata |
| Token/owner đổi và HMAC xoay khóa | token/hash/secret/target epoch | target → token, token/key → hash | Không sửa owner epoch cũ; giữ khóa khi còn fingerprint cần tra |

Với phép tách thuộc tính dựa FD, intersection xác định một phía giúp JOIN không sinh dòng giả. Các tập child đa trị không tái dựng bằng JOIN phẳng tất cả child: phải gom từng tập theo parent rồi dựng DTO; nếu JOIN vars × recipients × headers sẽ nhân dòng. Các bảng runtime/event có khóa và FD của sự kiện riêng; quota, sync/cancel receipt, rendered snapshot và action_ref không bị xóa chỉ vì có thể tính một giá trị hiện tại khác.

Không đồng nhất mất thuộc tính với mất thông tin: outcome, provider, context, attempt_count là view chỉ khi còn đủ nguồn đúng tuổi và đọc snapshot nhất quán. 193 ánh xạ trong từ điển chỉ định nơi lưu/view thay thế; retention cũng là phần bắt buộc của ánh xạ. Không dựng lại sync_result hay số đếm cancel từ callback hiện tại.

## 2 Đường khởi tạo phase đầu

1. Kiểm kê DB đích bằng catalog trước deploy; xác nhận chưa có bảng/dữ liệu nghiệp vụ. Ghi phiên bản schema, migration ledger và kiểm kê; không suy DB trống chỉ từ repo.
2. Review manifest với Tech Lead; viết DDL version mới với PK/unique/FK, row checks, trigger và protocol khóa. Không sửa migration 00001 đã chạy; tạo migration mới. Tạo bảng trước rồi FK để xử lý thứ tự dependency/vòng connection.
3. Tạo các FK vòng NO ACTION DEFERRABLE INITIALLY DEFERRED; FK còn lại RESTRICT. Dùng unique NULLS NOT DISTINCT và partial unique đúng manifest; không dùng partial index làm target FK.
4. Nạp snapshot/profile/schema/template theo hash trong transaction; đối chiếu artifact với cấu hình đã biên dịch. Không seed provider, cost center, credential hoặc quyền khi CF chưa có dữ liệu thật đã chốt.
5. Viết adapter map tên API và views, aggregate scope; giữ đúng transaction mục 8 từ điển. Chốt GET sync đang in-progress và đồng bộ OpenAPI trước tích hợp.
6. Chạy ma trận nghiệm thu mục 4 trên PostgreSQL 17 dùng schema tạo từ migration và restore từ backup. Chỉ mở nhận request sau khi bảo đảm commit accepted có đủ dữ liệu phục hồi.

Phase đầu dùng bảng thường và xóa lô. Partition chưa cần để đáp ứng DB01; xóa ở ngày 180 nằm trong cửa sổ đã chốt. Khi partition sau phải chứng minh unique/reference toàn cục và FK/receipt vẫn đúng, không đưa partition key vào mọi PK rồi vô tình nới chống trùng. PostgreSQL yêu cầu unique trên partitioned table chứa partition key; xem [tài liệu partitioning](https://www.postgresql.org/docs/17/ddl-partitioning.html).

## 3 Nếu kiểm kê phát hiện dữ liệu có sẵn

Đây là đường có điều kiện, không mặc định tồn tại hệ thống legacy.

1. Chụp backup nhất quán và thử restore cô lập; giữ manifest, migration ledger, thời điểm snapshot và key reference. Snapshot/backup nhạy cảm phải mã hóa và tuân TTL. Không log liên hệ/secret để đối chiếu.
2. Tạo schema/bảng đích theo cách mở rộng, giữ nguồn nguyên vẹn. Với phase đầu ưu tiên cửa sổ tạm dừng writer/worker/relay/TTL, thoát lease đang gửi và lấy snapshot; nếu cần online, phải thiết kế cơ chế delta/CDC và đối chiếu riêng, không dual-write best-effort. Không gọi provider khi backfill.
3. Dùng ID/ref hiện có; tạo bảng mapping ID mới nếu cần chuẩn hóa identity, với unique và checkpoint. Phân biệt tài khoản nguồn bằng source/namespace, scope app/org/brand, đúng config và connection version đã dùng. Không ép external ID sang UUID.
4. Backfill theo thứ tự parent → child bằng transaction từng lô có checkpoint và upsert kiểm nội dung. Tách JSON theo key ổn định; không âm thầm gộp hai selector hoặc overwrite biến/audit field trùng. Header, recipient_ref hoặc template trùng không khớp nghĩa phải đưa vào báo cáo lỗi và dừng cutover.
5. Thiếu hash/snapshot, scope/template, connection version, key giải mã hoặc epoch Push tin cậy: không suy từ giá trị hiện hành. Ghi lỗi đối chiếu với ID an toàn; cần nguồn xác minh hoặc quyết định xử lý có audit. Không giữ bản raw payload trong quarantine vượt TTL để chờ vô hạn.
6. Đối chiếu mọi dòng/193 thuộc tính còn hợp lệ theo ma trận. So số lượng theo source/type/ngày/trạng thái và child cardinality; anti-join tìm ID/ref mất hoặc thừa; so view/DTO, byte ciphertext/envelope và key_id còn hạn; đối chiếu HMAC trong môi trường được phép, không export secret. Count tổng bằng nhau chưa đủ.
7. Nếu snapshot đã hết TTL lúc cutover, áp xóa/che theo giờ hiện tại ở cả kế hoạch đối chiếu; ghi lý do loại bỏ hợp lệ và số lượng. Không phục hồi phone đã che, secret đã xóa hoặc inbox hết 90 ngày từ nguồn/backup.
8. Validate constraints, zero orphan, zero conflict chưa giải quyết, zero mismatch ngoài TTL đã giải thích. Khóa writer, áp delta cuối nếu có, so high-water mark rồi đổi adapter đọc/ghi; không ACK công việc chỉ vì migration đã copy ID.
9. Giữ nguồn read-only trong cửa sổ rollback đã định trước và đúng retention; không DROP bảng nguồn ở phase đầu. Rollback trước có write mới có thể đổi route về nguồn. Sau có write/receipt mới phải đồng bộ toàn bộ ý định/attempt/evidence/job/revision mới về schema tương thích hoặc forward-fix; không restore snapshot cũ rồi mất accepted.

Chưa có nguồn dữ liệu thì không tạo backfill giả hay báo đã chứng minh migration không mất dòng. Ma trận thiết kế là điều kiện để viết backfill và bằng chứng đối chiếu trên nguồn thật.

## 4 Ma trận nghiệm thu DDL và persistence

| Nhóm | Ca kiểm chứng | Kết quả bắt buộc |
|---|---|---|
| Khóa/NULL | Hai connection ACTIVE cùng scope/priority, hai delivery non-Push cùng rule, context/account organization NULL, event NULL/non-NULL | Unique đúng điều kiện; event NULL không bị gom nhầm; FK target toàn bảng |
| Case/scope | Rule khác hash/Shop; case không thuộc profile request; hai writer thêm kênh trùng | Commit bị chặn, không write skew; override Shop giữ đúng mẫu |
| Accepted | Kill API trước commit, sau commit/trước response, kill relay trước/sau publisher confirm | Trước commit không có ý định nửa vời; sau commit cây đầu vào và outbox đủ; redelivery không gửi thêm |
| PLAN | Resolver transient/terminal, trả rỗng/vượt 100; kill giữa chuyển contact và tạo SEND jobs | Retry có nghĩa vụ bền vững; terminal có receipt/lý do; rollback/commit nguyên tử |
| OTP | Kill trước marker, sau marker/trước provider, sau provider/trước receipt; callback muộn | OTP không vào DB/log; không gửi lại khi không chắc; receipt đầu tiên bất biến |
| SEND/quota | Hai worker tranh suất cuối; claim/cancel cạnh tranh; refund qua nửa đêm | Một claim hợp lệ; đúng bucket ngày cấp; refund một lần khi chắc chưa gửi |
| Callback | Đến trước response, duplicate, sai connection, accepted sau delivered, hard bounce sau accepted/delivered | Receipt immutable/link một lần; revision đúng; mức/xung đột theo ADR-04 |
| Aggregate | Hai recipient ALTERNATIVE, CASE SMS/Email, cây chưa PLAN, PLAN lỗi, callback và GET cạnh tranh | Không gộp mục tiêu người khác; PENDING đúng; snapshot revision/outcome nhất quán |
| Inbox | Hai recipient cùng source_ref; created_at muộn hơn accepted_at; ngày 90 và request ngày 180 | Không lẫn owner/context; nội dung/hiển thị hết hạn riêng; metadata xóa cùng request |
| Push/key | Đổi owner, replay version cũ, token thu hồi sau rotate HMAC; restart với ciphertext cũ | Epoch cũ không đổi chủ; tombstone vẫn tra được; đủ envelope/key, không tái ACTIVE từ version không tin cậy |
| Retention | TTL job đua PLAN/callback/secret renew; xóa matched evidence; xóa vòng connection | Không resurrect PII, không biến matched thành unmatched; không orphan; vòng FK xóa nguyên tử |
| Restore/failover | Restore backup 30 ngày tuổi; failover sau response accepted | Dọn TTL trước mở; không mất accepted, không tự phát lại UNKNOWN |
| Mapping | So ID/ref và 193 trường với projection v0.9 | Không bỏ thuộc tính lịch sử, mismatch có giải thích TTL; không chỉ so count |

Các ca này là tiêu chí cần triển khai và chạy, chưa phải kết quả test đã đạt. Lệnh node scripts/validate-db-design.mjs chỉ kiểm manifest/dictionary/legacy mapping và FK/type/key. Nó không chứng minh 3NF của FD chưa biết hoặc crash safety của runtime.

## 5 Trạng thái bàn giao

Thiết kế v0.9 và manifest đã lưu; giữ DB01–DB06. Chưa tạo bảng nghiệp vụ, đổi API hay chạy migration. Các quyết định còn theo Tech Lead/CF được liệt kê ở mục 9 từ điển. Phân tích normalization tham chiếu [Microsoft Learn](https://learn.microsoft.com/en-us/previous-versions/troubleshoot/microsoft-365/microsoft-365-apps/access/database-normalization-description); chi tiết unique/NULL/FK theo [PostgreSQL 17](https://www.postgresql.org/docs/17/ddl-constraints.html).
