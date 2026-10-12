# Notitek — Đánh giá Logical Database và Data Dictionary

Ngày đánh giá: 09/10/2026. Đối tượng: **Mô hình dữ liệu logic và từ điển dữ liệu v0.3** trong workspace hiện tại.

**Cập nhật 12/10/2026:** đã đối chiếu 15 điểm với LDB v0.6 ở mục 1.1. Phần 2 và 3 bên dưới giữ nguyên nội dung đánh giá v0.3 làm lịch sử; tên trường, trạng thái và bảng trong đó là của v0.3 (ví dụ `source_system`, `flow_code`, `goal_status`, `FAILED_FINAL`, `fallback_of_delivery_id`, `RESULT_CHANGED`, `REF`), đã được đổi theo OpenAPI ở LDB v0.5. Muốn biết hiện trạng, đọc mục 1.1.

## 1. Kết luận và phạm vi

**Thiết kế có nền tảng hợp lý, nhưng chưa đủ rõ để chốt toàn bộ Physical DB và triển khai các luồng R1.** Có 8 điểm phù hợp về hướng thiết kế và 15 điểm cần sửa hoặc bổ sung dưới đây. Đây là số nhóm nhận định, không phải điểm số chất lượng hay tỷ lệ đạt.

Phạm vi kiểm tra là tài liệu trong thư mục Notitek. Thư mục này chưa có mã nguồn ứng dụng, DDL, migration hoặc bộ kiểm thử thực thi. Vì vậy, “phù hợp” nghĩa là mô hình biểu diễn được yêu cầu; không chứng nhận hệ thống đã triển khai đúng. Các showcase là **mô phỏng bằng dữ liệu giả và suy luận từ đặc tả**, chưa phải kết quả chạy trên PostgreSQL/provider. Các mã `N1`, `D1`, `U1` là nhãn rút gọn cho định danh.

Ưu tiên trong báo cáo:

- **Cao:** cần chốt trước khi triển khai phần liên quan vì ảnh hưởng phục hồi, gửi đúng, bảo mật, bằng chứng hoặc hạn mức.
- **Vừa:** cần làm rõ trong baseline thiết kế, tránh các bên triển khai khác nghĩa.

Không xem việc chưa chọn kiểu PostgreSQL, index, UUIDv4/v7, khóa lạc quan hoặc DDL là lỗi của tài liệu logic. Tuy nhiên, ý nghĩa khóa, quan hệ, thời hạn, quy tắc duy nhất và vòng đời dữ liệu phải rõ trước khi chuyển sang thiết kế vật lý.

### Tài liệu đối chiếu

| Ký hiệu | Tài liệu | Cách sử dụng |
|---|---|---|
| LDB | [Mô hình dữ liệu logic và từ điển dữ liệu](<Notitek - MÔ HÌNH DỮ LIỆU LOGIC VÀ TỪ ĐIỂN DỮ LIỆU.md>) | v0.3; tài liệu được đánh giá; tham chiếu theo mục và tên trường. |
| DEC | [Quyết định PO cho các DEC](<Notitek - QUYẾT ĐỊNH PO CHO CÁC DEC.md>) | v1.2; quyết định cụ thể cho R1, đặc biệt DEC-01 đến DEC-10 và DEC-12. |
| SRS | [Đặc tả yêu cầu phần mềm](<Notitek - ĐẶC TẢ YÊU CẦU PHẦN MỀM (SRS).md>) | Yêu cầu hành vi, dùng mã SRS-F và mục để truy vết. |
| API | [Hợp đồng tích hợp API và sự kiện](<Notitek - HỢP ĐỒNG TÍCH HỢP API VÀ SỰ KIỆN.md>) | Bản nháp ngữ nghĩa; không xem phần còn mở là hợp đồng đã triển khai. |
| TECH | [Thiết kế kỹ thuật sơ bộ](<Notitek - THIẾT KẾ KỸ THUẬT SƠ BỘ.md>) | Ranh giới sở hữu và các cửa sổ lỗi phải xử lý. |
| AT | [Backlog triển khai và kịch bản nghiệm thu](<Notitek - BACKLOG TRIỂN KHAI VÀ KỊCH BẢN NGHIỆM THU.md>) | Gắn showcase với kịch bản nghiệm thu đã có. |

Khi các tài liệu khác ngày hoặc khác phiên bản có câu chữ chưa đồng bộ, báo cáo chỉ rõ chỗ lệch; không tự coi một phương án đề xuất là quyết định PO mới.

### 1.1 Kết quả đối chiếu với LDB v0.6

Trạng thái: **Đã đóng** (thiết kế logic đã có câu trả lời trong LDB), **Đóng một phần** (còn phần hoãn hoặc ngoài LDB), **Hoãn có chủ đích** (PO đã quyết định chờ provider/Push, ghi ở Q10). "Đã đóng" chỉ nghĩa là tài liệu đã đủ; chưa có DDL, code hay AT nào được chạy.

| Mã | Trạng thái | Nơi giải quyết trong LDB v0.6 | Còn lại |
|---|---|---|---|
| C01 | Đã đóng | Nguyên tắc 12; `notification_request.data` và `recipients_spec`; `transient_secret` (TTL ≤ 10 phút, không kéo dài); quy tắc phục hồi gọi lại API cấp link của User (CF-10) | Cần AT dừng tiến trình sau 202 (AT-02/05) |
| C02 | Đã đóng | §4.3 (`AVAILABLE` cho In-app, chờ retry là `PROCESSING` + `next_attempt_at`); §4.4 ma trận 19 dòng; quyết định A (bounce sau accepted) và B (`outcome` chốt khi `CLOSED`) | Cần AT cho dòng 15 |
| C03 | Hoãn có chủ đích | `delivery_evidence` chỉ cập nhật `attempt_id`, `applied_result` khi khớp muộn (nguyên tắc 4) nên callback đến sớm được khớp mà vẫn giữ bản gốc | Scope kết nối và khóa chống trùng theo `(connection_id, provider_event_id)` chờ provider (Q10); nên làm trước khi nhận callback thật |
| C04 | Đóng một phần | `credential_*` thành C (xóa khi `REVOKED`); `DISABLED` tách `REVOKED`; đổi credential/sender đưa `ACTIVE` về `DRAFT` | Lịch sử version/khóa xác thực callback sau khi xoay khóa chờ provider (Q10) |
| C05 | Đóng một phần | Giữ dòng `REVOKED` 180 ngày, chỉ xóa `token_encrypted`; `token_hmac` giữ để chặn đăng ký version cũ; delivery vẫn tham chiếu được | Nguồn cấp `registration_version`; bảo vệ sau 180 ngày (Q10) |
| C06 | Đã đóng | §5.13: cấp suất nguyên tử, cùng transaction với attempt, hoàn chỉ khi chắc chắn chưa gửi, không hoàn khi `ACCEPTED`/`UNKNOWN`; đơn vị đếm là attempt; báo cáo theo provider từ `delivery_attempt`; lịch sử hạn mức ở `audit_log.change` | Cần AT-19 song song và replay |
| C07 | Đã đóng | `audit_log.scope`, `change` (trước/sau theo allowlist), `reason` bắt buộc theo danh sách hành động | Allowlist trường cụ thể chốt khi thiết kế API |
| C08 | Đã đóng (chính sách) | Nguyên tắc 7: `rendered_content`, `title`, `preview` không chứa liên hệ đầy đủ, che trước khi lưu | Bước che phải được code và test; AT-03 |
| C09 | Đóng một phần | Nguyên tắc 13 (mốc `accepted_at`, xóa cùng lúc), §6 có mốc từng nhóm; v0.6 thêm retention của `contact_blocklist`; ngoại lệ của nguyên tắc 6 đã ghi | Q9: 180 ngày là tối thiểu hay tối đa khi drop partition (chờ PO) |
| C10 | Đã đóng | `cancel_marker.target_kind` tách `NOTIFICATION`/`GROUP`; phát lại trả biên nhận gốc, không kéo dài `tombstone_expires_at`; biên nhận ba số không chồng nhau | AT-10/23 |
| C11 | Đã đóng | ERD v0.5: idempotency 0..1, evidence phía attempt 0..1, `QUOTA_COUNTER` PK bốn cột, `OUTBOX_JOB` ghi chú quan hệ đa hình, `SEND_CASE` chỉ có sau lập kế hoạch | Không |
| C12 | Đã đóng | Duy nhất `(notification_id, case_code)`, `(case_id, recipient_ref)`, `(message_id, channel, push_subscription_id)`, `(delivery_id, attempt_no)`, `inbox_item.delivery_id`, `source_dedupe_key` có namespace | Chọn constraint vật lý khi làm DDL |
| C13 | Đã đóng | `reconcile_deadline_at` cố định khi vào `UNKNOWN`; `updated_at` không dùng cho cửa sổ; callback sau `CLOSED` vẫn ghi bằng chứng, không mở lại việc gửi | AT-23/25 |
| C14 | Đóng một phần | `brand_id` C (v0.6, không brand mặc định ngầm); `organization_id` đúng một theo `context_type`; `account_id` rỗng cho cả người được mời chưa có tài khoản; `hmac_key_id` ở idempotency | Xoay khóa HMAC của blocklist và token Push chờ kế hoạch (Q10) |
| C15 | Đóng một phần | LDB đã nhất quán: không có cột môi trường, hạn mức ở `quota_limit`, Push luôn `OPTIONAL` | Q6: câu chữ "môi trường" trong DEC-02, SRS-F02.01, CẤU TRÚC §11.1, API §7 và ví dụ Push bắt buộc ở SRS §8.2 chưa sửa (ngoài LDB) |

Tổng hợp: 9 đã đóng (C01, C02, C06, C07, C08, C10, C11, C12, C13), 5 đóng một phần (C04, C05, C09, C14, C15), 1 hoãn có chủ đích (C03). Không điểm mới nào phát sinh từ việc đối chiếu.

Phần "phù hợp" P01 đến P08 vẫn đúng về ý, nhưng ánh xạ tên như sau: P01 `goal_status` → `outcome`, `FAILED_FINAL` → `FAILED`; P02 `(source_system, source_ref)` → `(source, source_ref)` qua `idempotency_record`; P04 `FAILED_DEFINITE` → attempt `REJECTED_DEFINITE`, fallback → `requirement=ALTERNATIVE` cùng `requirement_group`; P06 R1 không có event `RESULT_CHANGED` (kết quả cho nguồn đi qua `GET` theo `revision`), `outbox_job` chỉ có `PLAN` và `SEND`; P07 không còn loại `REF`.

## 2. Các phần phù hợp, kèm showcase

### P01 — Phân rã theo yêu cầu, trường hợp, người nhận, kênh và lần thử

**Căn cứ:** LDB §3.1, §5.1, §5.4–5.7; SRS-F05.02/03, SRS-F09.01/02/05, SRS §8.2.

**Đánh giá:** Phân rã này giữ được kết quả riêng và không hợp nhất người nhận hàng với tài khoản Shop. Retry cũng không trở thành một sự việc nguồn mới.

**Showcase:** Event `E1` của Order tạo request `N1` và hai case:

| Case | Message | Người nhận | Delivery | Kết quả |
|---|---|---|---|---|
| SHOP_ALERT | M1 | U1 | IN_APP | Đã có inbox I1 truy cập được |
| SHOP_ALERT | M1 | U1 | PUSH | FAILED_FINAL |
| SHOP_ALERT | M2 | U2 | IN_APP | Đã có inbox I2 truy cập được |
| CONSIGNEE_NOTICE | M3 | account_id rỗng | ZNS | PROVIDER_ACCEPTED |

Tiêu chí Shop là In-app bắt buộc, Push tùy chọn: `SHOP_ALERT.goal_status=ACHIEVED`, trong khi lỗi Push vẫn tồn tại riêng. U2 không bị gộp vào U1; người nhận hàng không bị tạo tài khoản. Nếu ZNS retry, thêm attempt dưới delivery cũ. **Mô hình chứa đủ các lớp để biểu diễn tình huống này.** Việc ánh xạ trạng thái In-app còn thiếu được nêu riêng ở C02. Kiểm chứng triển khai: AT-11/25.

### P02 — Chống trùng tách khỏi dữ liệu gửi và có khả năng xoay khóa HMAC

**Căn cứ:** LDB §5.2; SRS-F02.01–07; DEC-09.

**Đánh giá:** Khóa `(source_system, source_ref)` trong DB riêng từng môi trường, kết hợp `payload_hmac`, `hmac_key_id` và `notification_id`, đáp ứng cách nhận diện replay/xung đột mà không cần lưu OTP thô để so sánh.

**Showcase:** Lưu `(USER, R1, HMAC_K1(payload_A), K1, N1)`. Gọi lại `R1` với cùng dữ liệu nghiệp vụ trả `N1`. Gọi lại `R1` nhưng đổi người nhận hoặc mã phải trả xung đột. Sau khi khóa hiện hành đổi thành K2, đọc `hmac_key_id=K1` để so với bản cũ trong cửa sổ 7 ngày. **Các trường cần cho ba nhánh đều có.** Điều kiện triển khai: chuẩn hóa payload và giữ khóa K1 đủ cửa sổ; tài liệu chưa chứng minh transaction chống tranh chấp đã tồn tại. AT-02/18/28.

### P03 — Tách trạng thái đọc khỏi trạng thái gửi và nghiệp vụ nguồn

**Căn cứ:** LDB §5.9; SRS-F08.05/06, SRS-F10.03; DEC-03.

**Đánh giá:** `inbox_item.read_at` thuộc từng người; trạng thái chấp nhận/từ chối yêu cầu thiết bị vẫn do User sở hữu. Dữ liệu scope có ngay trên inbox phục vụ truy vấn đúng người và ngữ cảnh.

**Showcase:** I1 thuộc `(U1, WORKSPACE, Shop-A)`, I2 thuộc `(U1, WORKSPACE, Shop-B)`. U1 đọc I1 lúc 09:00: chỉ `I1.read_at` đổi; I2 vẫn chưa đọc, trạng thái SMS không đổi, lời mời không tự được chấp nhận. U1 rời Shop-A: truy vấn phải kết hợp quyền hiện tại để không trả I1. **Mô hình hỗ trợ tách scope và đọc; quyền hiện tại vẫn phải được kiểm tra bằng code.** AT-08/15/26.

### P04 — Phân biệt timeout với thất bại chắc chắn trước fallback

**Căn cứ:** LDB §4.3, §5.6–5.7; DEC-05/10; SRS-F07, SRS §8.3.

**Đánh giá:** `UNKNOWN`, `FAILED_DEFINITE`, `fallback_of_delivery_id` và vai trò `FALLBACK` biểu diễn đúng các điều kiện khác nhau.

**Showcase:** ZNS D1 có attempt A1 bị timeout → `outcome=UNKNOWN`, D1 chưa rõ; SMS D2 tham chiếu D1 **chưa được gửi**. Nếu A1 trả từ chối chắc chắn, D1 hết phương án retry hợp lệ và chính sách cho fallback, D2 mới chạy sau khi kiểm tra hạn gửi/khung giờ/hạn mức. **Dữ liệu phân biệt được hai tình huống, không cần suy timeout thành thất bại.** AT-09/25.

### P05 — Tách chính sách được duyệt khỏi dữ liệu vận hành thay đổi lúc chạy

**Căn cứ:** LDB nguyên tắc 2, §5.1, §5.15–5.18; DEC-07/08.

**Đánh giá:** Ghim `config_version` cho chính sách, trong khi kết nối, kill-switch, Push của Shop và hạn mức nằm trong DB, phù hợp quyết định PO v1.2.

**Showcase:** N1 được tạo với YAML commit V1. Khi đang chờ, triển khai V2 đổi mẫu và vận hành bật `(DELIVERY_FAILED, ZNS)` trong `kill_switch`. Worker vẫn dùng mẫu V1 nếu xử lý N1, nhưng phải chặn lần gửi ZNS mới theo kill-switch hiện hành. **Hai loại thay đổi có nơi lưu và ngữ nghĩa riêng.** Điều kiện: runtime còn truy xuất được V1 và làm mới dữ liệu vận hành theo SLA được chốt. AT-17/18/24.

### P06 — Dùng outbox cho công việc và kết quả, không nhúng bí mật vào broker payload

**Căn cứ:** LDB §4.1, §5.13; SRS-F09.07/08; TECH §4.

**Đánh giá:** Có `event_kind`, định danh aggregate, `outbox_id`, thời điểm phát và revision trong payload là nền tảng phù hợp để tách phát kết quả khỏi gửi lại người dùng.

**Showcase:** Cùng một transaction ghi thay đổi kết quả N1 revision 3 và O1 loại `RESULT_CHANGED`. Dispatcher đã publish O1 nhưng crash trước khi ghi `dispatched_at`. Sau phục hồi, nó phát lại **O1/revision 3**, không tạo delivery mới. Consumer dùng ID thay đổi ổn định để loại bản trùng. **Mô hình chứa được ID và nội dung phục vụ replay.** Cần chốt ánh xạ `outbox_id` thành ID event và transaction trong thiết kế thực thi; outbox tự nó chưa chứng minh exactly-once. AT-22/25.

### P07 — Có dấu hủy độc lập để xử lý hủy đến trước sự việc

**Căn cứ:** LDB §5.3; SRS-F06.06; DEC-09.

**Đánh giá:** Không bắt buộc FK tới request là hợp lý vì request có thể chưa tồn tại.

**Showcase:** 10:00 nhận dấu hủy `(ORDER, GROUP, G1)`, hết hạn sau 72 giờ. 10:05 event E1 tới với `cancel_group=G1`. Tiếp nhận/worker đối chiếu dấu và dừng phần chưa gửi; không cần tạo request giả lúc 10:00 chỉ để chứa hủy. **Mô hình biểu diễn được sự kiện đảo thứ tự này.** Tính không mơ hồ của loại REF và replay lệnh hủy vẫn cần C10. AT-10/23.

### P08 — Nhiều đích Push của một người có kết quả riêng

**Căn cứ:** LDB §5.6, §5.10; SRS-F16.01/05; DEC-01.

**Đánh giá:** Mỗi delivery Push tham chiếu một subscription giúp giữ chi tiết từng đích, và token chỉ có một chủ ACTIVE theo quy tắc đã ghi.

**Showcase:** U1 có subscription S1 và S2. Message M1 sinh D1→S1 và D2→S2. D1 được provider nhận; D2 bị báo token hết hiệu lực, S2 bị thu hồi. Tra cứu vẫn thấy một đích nhận, một đích lỗi, không gộp D2 thành thành công. **Cấu trúc đáp ứng trường hợp đa đích.** Chính sách chọn tập đích và vòng đời sau thu hồi cần hoàn thiện; xem C05. AT-12/25/26.

## 3. Các phần chưa phù hợp hoặc chưa đủ, kèm phản ví dụ

### C01 — Chưa có dữ liệu bền vững để tiếp tục lời mời bất đồng bộ sau crash — Cao

**Loại:** Thiếu thiết kế đường dữ liệu.

**Căn cứ:** LDB §5.1 không lưu payload thô; §5.6 `rendered_content` có thể rỗng và không chứa bí mật; §5.13 outbox chỉ ID/trạng thái. SRS-F02.04, SRS-F13.02/05, SRS-F17.02; DEC-02/06/09; TECH §4.

**Showcase:** User gửi lời mời I1, email có link `https://example.invalid/invite?token=SECRET_DEMO`. Notitek trả 202 rồi tiến trình dừng trước khi gửi. Worker khởi động lại chỉ có request, template và outbox chứa ID. Nếu lưu link đầy đủ vào `rendered_content`, trái quy tắc không lưu bí mật ở đó. Nếu không lưu và không có hợp đồng lấy lại từ User, worker không khôi phục được email đúng link.

**Hệ quả:** Accepted chưa có đường phục hồi dữ liệu gửi được đặc tả đầy đủ. Tương tự, biến không bí mật phục vụ dựng mẫu sau accepted cũng cần nơi lưu nếu chưa dựng xong.

**Đề nghị:** Chốt một đường: (a) dựng và lưu nội dung không bí mật trước accepted, phần bí mật lưu tạm mã hóa riêng; hoặc (b) lưu tham chiếu và hợp đồng lấy lại dữ liệu từ nguồn. Nếu chọn (a), định nghĩa ciphertext/ref, key ID, hạn xóa và quyền giải mã. Phải giải quyết việc DEC-09 giới hạn bí mật xếp hàng tối đa 10 phút trong khi email được retry tới 24 giờ; không tự kéo TTL bí mật lên 24 giờ.

**Showcase nghiệm thu:** Dừng tiến trình ngay sau 202, khởi động lại và gửi đúng I1; không lộ SECRET_DEMO ở lịch sử/outbox. Retry sau khi hết TTL phải có hành vi được chốt. AT-02/03/05/28.

### C02 — Bảng trạng thái chưa bao phủ In-app, retry và bounce sau accepted — Cao

**Loại:** Thiếu quy tắc chuyển trạng thái.

**Căn cứ:** LDB §4.3, §5.6–5.8; SRS-F08.01, SRS-F09.03/06, SRS-F17.04 và SRS §7.2/8.3; DEC-10.

**Showcase A — In-app:** Transaction tạo I1 thành công, người nhận truy cập được. Delivery IN_APP phải dùng giá trị state nào? Danh sách hiện có không nêu `AVAILABLE`; `PROVIDER_ACCEPTED` không đúng nghĩa kênh nội bộ, còn nếu dùng `DELIVERED` thì phải ghi rõ đây là bằng chứng lưu/truy cập nội bộ, không cần provider attempt.

**Showcase B — Retry:** A1 trả `FAILED_RETRYABLE`, lần sau được lên lịch +5 phút. Delivery quay về `PENDING`, giữ `SENDING`, hay có trạng thái chờ retry? Mục 4.3 chưa mô tả cạnh chuyển này.

**Showcase C — Bounce:** Email được `PROVIDER_ACCEPTED` lúc 10:00, hard bounce lúc 10:02. Accepted và bounce có thể là hai bước hợp lệ, không nhất thiết là bằng chứng trái nhau. Quy tắc “state giữ mức cao nhất” chưa chỉ ra state mới, mức mục tiêu, revision và kết thúc retry. SRS yêu cầu thể hiện bounce.

**Đề nghị:** Thêm ma trận `(state hiện tại, sự kiện, điều kiện) → (state mới, cập nhật goal, phát revision, có retry không)`. Phân biệt bằng chứng từng nhận với kết quả giao sau đó. Chốt cách biểu diễn In-app, retry, accepted→bounce và delivered→bounce conflict; nếu một cột state vẫn đủ thì không cần thêm cột chỉ để tăng số trường.

**Showcase nghiệm thu:** Ba nhánh trên có đầu ra duy nhất được PO/Tech Lead thống nhất; lịch sử accepted vẫn còn khi bounce đến. AT-07/25.

### C03 — Callback UNMATCHED chưa có scope kết nối và cách khớp lại khi append-only — Cao

**Loại:** Thiếu cơ chế đối chiếu.

**Căn cứ:** LDB §5.8 chỉ có provider/message/event ID khi `attempt_id` rỗng; khóa chống trùng là `(provider, provider_event_id)`. §5.15 yêu cầu xác thực bằng đúng kết nối đã gửi. SRS-F09.06, API §4.

**Showcase A — Scope:** Hai kết nối C1/Brand-A và C2/Brand-B cùng provider. Nếu provider chỉ bảo đảm ID trong từng tài khoản, cả hai có thể trả `provider_message_id=42`, `provider_event_id=7`. Dòng evidence hiện không có scope kết nối khi UNMATCHED; unique theo provider/event ID sẽ loại nhầm callback thứ hai. Đây là rủi ro có điều kiện, cần đối chiếu hợp đồng provider, không khẳng định provider thực đã dùng ID như vậy.

**Showcase B — Đến sớm:** Callback E1 đến trước khi response gửi được lưu, ghi `attempt_id=NULL, applied_result=UNMATCHED`. Sau đó A1 được xác định. Sửa E1 vi phạm append-only; thêm E1 lần hai vi phạm chống trùng; giữ nguyên mà không có liên kết bổ sung thì không tra được quyết định áp dụng về sau.

**Đề nghị:** Lưu ngữ cảnh kết nối đã xác thực hoặc scope tương đương; chốt tính duy nhất của ID theo từng adapter và đường dedupe khi provider không cấp event ID. Tách evidence bất biến khỏi bản ghi matching/application, hoặc đặc tả cơ chế liên kết append-only có tham chiếu evidence cũ.

**Showcase nghiệm thu:** Hai callback cùng ID nhưng khác scope không bị trộn; callback đến sớm được khớp đúng một lần mà vẫn giữ bản gốc. AT-12/25.

### C04 — connection_version chưa giúp khôi phục lịch sử; thu hồi credential mâu thuẫn nullability — Cao

**Loại:** Thiếu dữ liệu phiên bản và mâu thuẫn vòng đời.

**Căn cứ:** LDB §5.7 có `connection_version`; §5.15 chỉ có dòng kết nối hiện hành, `credential_encrypted=Y`, `credential_key_id=Y`, nhưng yêu cầu xóa credential khi thu hồi. SRS-F15.05; DEC-07/10.

**Showcase A — Đổi khóa:** A1 gửi bằng C1/v1 và khóa callback K1. Vận hành cập nhật C1/v2 với K2, rồi callback của A1 đến. `connection_version=1` chỉ là số; không có lịch sử hoặc quy tắc giữ khóa cũ thì chưa biết lấy K1 ở đâu để xác thực.

**Showcase B — Thu hồi:** Thu hồi C1, xóa credential nhưng giữ C1 để attempt còn tham chiếu. Bản ghi sau xóa không đáp ứng hai trường bắt buộc Y; `DISABLED` cũng chưa phân biệt dừng tạm với thu hồi hẳn.

**Đề nghị:** Chốt lịch sử phiên bản kết nối hoặc tham chiếu version bất biến, chính sách khóa xác thực callback khi xoay/thu hồi, và nullability theo trạng thái. Khi cần thu hồi bí mật ngay, phải nêu cách xử lý callback cũ an toàn; không mặc định giữ credential gửi đã bị thu hồi. Đổi credential của kết nối ACTIVE cũng phải xác định việc dừng gửi cho tới khi thử lại đạt, vì hiện chỉ ghi rõ điều kiện DRAFT→ACTIVE.

**Showcase nghiệm thu:** Callback A1 được xử lý theo chính sách v1 rõ ràng sau khi có v2; thu hồi không tạo bản ghi trái dictionary; đổi khóa chưa xác minh không được vô tình tiếp tục gửi. AT-18/24/28.

### C05 — Xóa subscription làm mất dấu chống hồi sinh và vướng lịch sử delivery — Cao

**Loại:** Xung đột retention với yêu cầu vòng đời.

**Căn cứ:** LDB §3, §5.6, §5.10, §6; SRS-F16.02/03, DEC-01/09.

**Showcase:** S1 có token T, registration version 8 và được D1 tham chiếu. Người dùng logout, S1 bị thu hồi rồi xóa theo bảng retention. Request đăng ký cũ version 7 đến sau đó. Nếu không còn dấu T/version 8 thì kiểm tra version không có dữ liệu để từ chối bản cũ. Đồng thời D1 còn giữ `push_subscription_id` bắt buộc trong 180 ngày.

**Hệ quả:** Có cột version chưa đủ bảo đảm chống hồi sinh; vòng đời FK/tham chiếu chưa được mô tả.

**Đề nghị:** Xóa token mã hóa đúng hạn nhưng giữ tombstone tối thiểu không chứa token thô, hoặc dùng thế hệ/session do nguồn xác nhận với hợp đồng chống replay tương đương. Chốt thời gian lưu tombstone, ý nghĩa version, nguồn cấp version, và cách lịch sử delivery tham chiếu subscription đã thu hồi. Không dùng cascade làm mất lịch sử gửi 180 ngày.

**Showcase nghiệm thu:** Xóa bí mật của S1 rồi phát lại version 7: vẫn bị từ chối; D1 vẫn tra được; token đổi sang U2 không khiến delivery cũ gửi cho U2. AT-12/26/28.

### C06 — “used_count tăng khi gửi” chưa đủ định nghĩa hạn mức cứng và quyết định chi — Cao

**Loại:** Thiếu ngữ nghĩa nghiệp vụ, không chỉ thiếu khóa DB.

**Căn cứ:** LDB §5.12, §5.18; SRS-F14.02/04/05; DEC-08.

**Showcase A — Đồng thời:** Limit=100, used=99. Hai worker cùng thấy còn một suất rồi cùng gọi provider: thành 101 nếu không có thao tác cấp suất nguyên tử.

**Showcase B — Crash/timeout:** Worker cấp suất cho A1, provider có thể đã nhận, nhưng worker crash trước khi ghi kết quả. Replay A1 được dùng lại suất hay tăng thêm? Nếu hoàn suất khi UNKNOWN thì có thể vượt giới hạn; nếu cấp suất lại thì có thể đếm đôi.

**Showcase C — Đối soát:** Hạn mức đổi từ 100 sang 200; A1 từng bị chặn lúc mức 100. `quota_counter` và `quota_limit` hiện hành không cho biết A1 dùng quyết định nào. Báo cáo theo provider phải lấy từ attempt và cần định nghĩa đếm lần thử, lần accepted hay đơn vị có phí; riêng counter không có chiều provider.

**Đề nghị:** Chốt đơn vị đếm, thời điểm cấp/tiêu suất, xử lý lỗi chắc chắn/UNKNOWN, khóa idempotency theo attempt và tham chiếu quyết định hạn mức. Nếu dùng reservation thì ghi vòng đời tối thiểu; không cần biến Notification thành sổ tài chính. Định nghĩa join/report theo cấu hình đã ghim và connection/attempt, cùng retention của counter.

**Showcase nghiệm thu:** Với 99/100 chỉ một worker được cấp suất; replay A1 không cấp lần hai; UNKNOWN không tự hoàn suất; báo cáo tháng giải thích được số lượng theo provider. AT-19/25/27.

### C07 — Audit không giữ giá trị trước/sau dù được dùng làm lịch sử duy nhất — Cao

**Loại:** Thiếu dữ liệu chứng minh thay đổi.

**Căn cứ:** LDB §5.14, §5.16–5.18 và §7 bỏ `scope` vì suy từ target; DEC-07/08/09/12.

**Showcase:** Ngày 1 admin đổi hạn mức 100→200, ngày 2 đổi 200→150. Audit hiện chỉ ghi actor/action/target/time/reason/outcome. Đến ngày 3, DB có 150; hai audit không chứng minh được giá trị 100 và 200. Sau 180 ngày request bị xóa, audit xem dữ liệu gốc vẫn còn 2 năm nhưng scope không còn suy ra được từ request đã mất.

**Đề nghị:** Ghi diff trước/sau theo allowlist không bí mật, hoặc trỏ tới bản version bất biến có retention tương ứng. Giữ scope tối thiểu trong audit nếu cần tra độc lập sau khi target bị xóa. Đặt `target` và `reason` bắt buộc theo action; quy tắc hiện chỉ bắt buộc reason cho xem gốc/hủy/gửi lại, trong khi kết nối/kill-switch cũng mô tả phải biết vì sao thay đổi. Không ghi credential/token vào diff.

**Showcase nghiệm thu:** Từ audit dựng lại 100→200→150 và phạm vi liên quan sau khi dữ liệu đích bị xóa. AT-16/17/19/28.

### C08 — Che destination sau 30 ngày chưa bảo đảm xóa mọi bản sao liên hệ — Cao

**Loại:** Chưa khép kín chính sách dữ liệu cá nhân.

**Căn cứ:** LDB nguyên tắc 7, §5.6, §5.9 và §6; SRS-F13.04/05; DEC-09.

**Showcase:** Mẫu hợp lệ giả định có câu “Liên hệ người nhận 0900000123”. Số được lưu trong `destination_value`, `rendered_content` và có thể cả `inbox_item.preview`. Ngày 31 job chỉ che destination; số vẫn còn trong nội dung 90 ngày và preview có thể tồn tại 180 ngày. Khẳng định “liên hệ nguyên văn ở đúng một chỗ” chưa được bảo đảm bởi định nghĩa các trường tự do.

**Đề nghị:** Chốt allowlist nội dung được lưu; nếu cấm liên hệ nguyên văn trong rendered/inbox thì ghi rõ bước redaction trước lưu. Nếu cho phép, phải đưa các bản sao vào chính sách che/xóa tương ứng. Phân biệt retention của nội dung preview với metadata inbox: ghi rõ ngoại lệ được duyệt nếu preview giữ 180 ngày còn nội dung dựng giữ 90 ngày.

**Showcase nghiệm thu:** Gửi mẫu giả trên, chạy lifecycle tới ngày 31/91/181, tìm giá trị trên mọi trường có thể lưu bản sao; API hỗ trợ và restore cũng tuân cùng chính sách. Đây là phản ví dụ khi mẫu cho phép dữ liệu đó, không khẳng định đã có rò rỉ thực tế. AT-03/28.

### C09 — Retention thiếu mốc tính ở một số bảng/cột và chưa giải quyết ranh giới partition — Vừa

**Loại:** Thiếu quy tắc vòng đời trước khi thiết kế Physical DB.

**Căn cứ:** LDB §5.1, §5.6–5.8, §6 và §7; DEC-09.

**Showcase A — Mốc tính:** N1 accepted ngày 1, attempt A1 ngày 2, evidence E1 ngày 3. “180 ngày” của từng bảng tính từ accepted, started hay received? Đến ngày 181 xóa N1 mà E1 còn tới ngày 183 thì quan hệ cha/con xử lý thế nào? Có thể dùng mốc cha, nhưng tài liệu cần nói rõ; không bắt buộc thêm created_at vào mọi bảng.

**Showcase B — Drop theo tháng:** Một partition có dòng cũ nhất đã 180 ngày, dòng mới nhất mới 150 ngày. Drop ngay sẽ xóa dòng 150 ngày quá sớm; đợi toàn bộ partition đủ 180 ngày thì dòng cũ nhất được giữ thêm khoảng một tháng. Cần chốt 180 ngày là giới hạn tối đa, tối thiểu hay cửa sổ có dung sai.

**Đề nghị:** Ghi retention anchor cho từng nhóm, thứ tự hoặc chính sách xóa phụ thuộc, và cơ chế làm sạch chính xác nếu 180 ngày là trần cứng. Bổ sung retention cho `quota_counter`, `contact_blocklist` và dữ liệu phiên bản/tombstone nếu thêm. Giải thích nguyên tắc 6 “khác thời gian lưu thì khác bảng” so với channel_delivery đang giữ nhiều TTL bằng cách xóa cột; chọn ngoại lệ có lý do hoặc tách phần nhạy cảm.

**Showcase nghiệm thu:** Với bộ ngày lệch trên, chạy xóa mà không để tham chiếu mồ côi, không xóa quá sớm và không giữ quá thời hạn đã duyệt. AT-28.

### C10 — REF của dấu hủy gộp hai loại định danh; chưa định nghĩa replay lệnh hủy — Cao

**Loại:** Khóa mục tiêu chưa đủ rõ.

**Căn cứ:** LDB §5.3 ghi REF là `source_ref`/`notification_id`; SRS-F06.03/05/06; API §8 yêu cầu ID lệnh hủy.

**Showcase A — Trùng miền giá trị:** Trong cùng nguồn ORDER, N1 có `source_ref=X`, còn một request khác có `notification_id=X`. Dấu `(ORDER, REF, X)` hủy bản nào? UUID viết dạng chuỗi vẫn có thể được nguồn dùng làm source_ref; không có quy tắc cấm.

**Showcase B — Replay:** Lệnh K1 tạo dấu hủy lúc ngày 1, expiry ngày 4. Replay K1 ngày 3 mà upsert expiry thành ngày 6 có thể mở rộng cửa sổ hủy. LDB chưa chỉ ra giữ expiry cũ hay đó là lệnh mới; cũng không giữ ID lệnh trong cancel_marker để phân biệt.

**Đề nghị:** Tách `SOURCE_REF`, `NOTIFICATION_ID`, `GROUP` hoặc chuẩn hóa về một miền được định nghĩa. Ghi hợp đồng ID lệnh/idempotency và quy tắc expiry bất biến cho replay. Với hủy và gửi đua nhau, cần lưu quyết định cùng thời điểm theo từng phần; audit hoặc lịch sử riêng có thể đáp ứng nếu quy định rõ.

**Showcase nghiệm thu:** Hai request dùng giá trị X không bị hủy nhầm; replay K1 không mở rộng tập hoặc thời hạn; trả được phần dừng/đã ra kênh/chưa rõ. AT-10/23.

### C11 — ERD chưa nhất quán với dictionary và vòng đời — Vừa

**Loại:** Mâu thuẫn trực tiếp trong tài liệu.

**Căn cứ:** LDB §3 so với §5.2, §5.8, §5.12, §5.13 và §6.

**Showcase A:** ERD `IDEMPOTENCY_RECORD ||--|| NOTIFICATION_REQUEST` buộc request có một idempotency record. Ngày 8 record bị xóa, request còn tới 180 ngày: request hợp lệ khi không còn record. Cần thể hiện 0..1 ở phía record đối với request trong toàn vòng đời.

**Showcase B:** ERD `DELIVERY_ATTEMPT ||--o{ DELIVERY_EVIDENCE` cho mỗi evidence đúng một attempt, nhưng dictionary cho phép UNMATCHED với `attempt_id=NULL`. Phía attempt của evidence phải là 0..1.

**Showcase C:** ERD QUOTA_COUNTER chỉ đánh dấu PK `(counter_date, cost_center)`, dictionary nói bốn cột gồm cả flow/channel. Hai dòng cùng ngày/cost center nhưng OTP-SMS và DELIVERY_FAILED-ZNS là hợp lệ theo dictionary nhưng sẽ đụng khóa nếu triển khai theo ERD.

**Showcase D:** ERD vẽ request→outbox trực tiếp, dictionary lại dùng aggregate đa hình, ví dụ DELIVERY/delivery_id, không có `notification_id`. Job của D1 không có FK trực tiếp được mô tả tới N1 như hình có thể khiến người đọc hiểu.

**Đề nghị:** Sửa cardinality/PK; vẽ hoặc chú thích rõ quan hệ logic qua aggregate. Nêu thêm giai đoạn request accepted nhưng chưa resolve recipient có thể chưa có message/case con, thay vì chỉ vẽ cardinality ở trạng thái hoàn tất fan-out.

**Showcase nghiệm thu:** Dựng các instance ngày 8, evidence chưa khớp, hai counter khác luồng và outbox của delivery; tất cả phù hợp cả hình lẫn dictionary.

### C12 — Thiếu quy tắc duy nhất cho các đối tượng con và namespace của khóa inbox — Cao

**Loại:** Thiếu invariant logic cho phục hồi/replay.

**Căn cứ:** LDB §5.4–5.9; SRS-F02.02–05, SRS-F05.03, SRS-F10.02.

**Showcase A — Worker replay:** N1 không bị tạo trùng nhờ idempotency, nhưng hai lần xử lý job có thể tạo hai case SHOP_ALERT, mỗi case có một M1/U1 và inbox riêng nếu không có invariant phía con. UUID riêng biệt không ngăn trùng nghiệp vụ.

**Showcase B — Push:** Unique `(message_id, channel)` sẽ ngăn hai subscription hợp lệ của cùng người; không unique gì lại cho phép tạo hai delivery cho cùng subscription do job replay. Cần phân biệt đích trong identity delivery.

**Showcase C — Inbox:** U1 nhận event của ORDER với ID `42` và cảnh báo USER với ID `42`. Nếu `source_dedupe_key` chỉ là ID nguồn như ví dụ và unique `(account_id, source_dedupe_key)`, hai sự việc khác nhau bị nhập một. Cần hợp đồng namespace; không khẳng định ID thật sẽ trùng nếu nền tảng đã bảo đảm globally unique.

**Đề nghị:** Ghi invariant như một case theo `(notification_id, case_code)`; một bản người nhận theo identity trong case; một inbox cho message; số attempt duy nhất trong delivery; identity delivery có thêm subscription hoặc định danh đích khi cần. Người nhận không có account phải có cách phân biệt ổn định, không tự unique theo số/email. Chuẩn hóa `source_dedupe_key` gồm namespace nguồn và phạm vi cần thiết hoặc cam kết định danh toàn cục. Chọn constraint vật lý sau.

**Showcase nghiệm thu:** Replay đồng thời không tăng số case/message/inbox/delivery ngoài dự kiến; hai Push hợp lệ vẫn tồn tại; hai nguồn có local ID 42 không triệt tiêu nhau. AT-02/14/25.

### C13 — updated_at không phải mốc ổn định cho cửa sổ đối chiếu 24 giờ — Cao

**Loại:** Thiếu mốc nghiệp vụ và quy tắc đóng/tái cập nhật.

**Căn cứ:** LDB §4.1/4.3, §5.6 `updated_at`; DEC-10; SRS §8.3.

**Showcase:** D1 UNKNOWN từ 09:00 ngày 1. Callback/lookup hoặc chuyển trạng thái khác lúc 08:55 ngày 2 làm đổi `updated_at`. Nếu dùng `updated_at + 24h`, cửa sổ bị kéo dài gần thêm một ngày. Nếu một request có hai delivery UNKNOWN cách nhau nhiều giờ thì cũng chưa rõ đóng theo delivery đầu, cuối hay từng delivery.

**Đề nghị:** Chốt mốc bắt đầu cố định, ví dụ `reconcile_started_at`/`reconcile_deadline_at`, hoặc quy tắc suy ra từ attempt bất biến nếu đủ. Quy định sau `CLOSED/RECONCILE_TIMEOUT`, callback mới có được tăng bằng chứng/goal/revision trong khi processing vẫn CLOSED hay không; không tự mở lại việc gửi.

**Showcase nghiệm thu:** Callback lúc giờ 23:55 không reset cửa sổ; sau giờ 24 vẫn UNKNOWN nếu chưa có bằng chứng; callback DELIVERED giờ 25 được ghi nhận theo chính sách nhưng không sinh lần gửi mới. AT-23/25.

### C14 — Một số trường bắt buộc có điều kiện và vòng đời khóa chưa được định nghĩa đủ — Vừa

**Loại:** Dictionary chưa đủ chính xác để sinh validation thống nhất.

**Căn cứ:** LDB §5.1, §5.5, §5.10–5.11; đối chiếu mẫu mô tả rõ hơn ở §5.2/5.15.

**Showcase A — Điều kiện C:** Request OTP không có `brand_code`: mục 5.1 ghi C và mục đích chọn sender/cost center nhưng chưa nói khi nào bắt buộc hoặc lấy mặc định từ đâu. Request WORKSPACE có cả shop_id và org_id cũng chưa rõ được phép hay phải đúng một. Hai đội có thể tạo validator khác nhau.

**Showcase B — Người được mời:** Người nhận lời mời chưa có tài khoản. `account_id` mô tả rỗng khi người nhận hàng, nhưng chưa nêu rõ trường hợp email người được mời không có account. Cần mở nghĩa “liên hệ không có tài khoản”, không chỉ consignee.

**Showcase C — Xoay HMAC:** Blocklist lưu `HMAC_K1(contact)`; hệ thống đổi sang K2 rồi tính `HMAC_K2(contact)` để tra: không khớp, dù liên hệ vẫn phải bị chặn. Push `token_hmac` có rủi ro tương tự khi kiểm tra unique ACTIVE. Hai bảng này không có `hmac_key_id` hoặc quy tắc tra nhiều thế hệ; token mã hóa cũng chưa ghi cách xác định phiên bản khóa giải mã. Envelope có thể tự mang key ID, nhưng tài liệu chưa nói.

**Đề nghị:** Bổ sung ma trận bắt buộc theo flow/context/channel/status, quy tắc shop/org, phạm vi người nhận không có account, canonicalization và quy tắc quản lý thế hệ khóa. Có thể dùng envelope tự mô tả hoặc dịch vụ khóa, không bắt buộc thêm cột nếu cơ chế khác được đặc tả đầy đủ. Với blocklist không lưu contact thô, kế hoạch xoay khóa phải tính đến việc không thể re-HMAC chỉ từ HMAC cũ.

**Showcase nghiệm thu:** Validator cho các payload trên có đáp án thống nhất; xoay khóa không bỏ chặn liên hệ, không tạo hai chủ ACTIVE cho cùng token và vẫn giải mã được token hợp lệ. AT-08/12/24/28.

### C15 — Baseline còn câu chữ chưa đồng bộ giữa các mục và tài liệu — Vừa

**Loại:** Sai lệch tài liệu, không mặc định là sai thiết kế.

**Căn cứ:** LDB nguyên tắc 8, §7, Q6; DEC-02/07/08; SRS-F02.01; API §2/5/6; SRS §8.2.

**Showcase A — Môi trường:** LDB bỏ environment vì DB tách theo môi trường; DEC-02 và SRS-F02.01 vẫn nói khóa có environment. Hai cách có thể tương đương về cách ly, nhưng phải ghi rõ môi trường nằm ở ranh giới deployment, không phải cột. Q6 đã nhận diện việc này nhưng chưa đồng bộ văn bản.

**Showcase B — Hạn mức:** §5.12/5.18 nói hạn mức ở `quota_limit`; §7 lý do bỏ `limit_applied` vẫn nói hạn mức ở cấu hình. Người triển khai theo phần sau có thể tiếp tục đọc YAML thay vì dữ liệu vận hành DB.

**Showcase C — Push/hợp đồng:** DEC-01 và LDB quy định Push luôn OPTIONAL trong R1; SRS §8.2 có ví dụ Push bắt buộc và câu không tự chốt Push tùy chọn cho mọi Shop. API vẫn nói chủ đích Push cần xác định và route/event chưa chốt, trong khi DEC-01/02 đã chốt một phần. Cần đánh dấu ví dụ tổng quát và phần đã được DEC mới thay thế để không dùng nhầm làm cấu hình R1.

**Đề nghị:** Đồng bộ baseline và ghi quan hệ thay thế giữa các phiên bản. Giữ tách bạch “tên/ngữ nghĩa PO đã chốt” với “OpenAPI/schema chi tiết còn chờ ký”. Không cần thêm lại environment vào DB chỉ để khớp câu chữ cũ nếu cách ly triển khai đã được xác nhận.

**Showcase nghiệm thu:** Một tình huống R1 bất kỳ chỉ có một kết luận thống nhất về scope idempotency, nguồn hạn mức, vai trò Push và bên sở hữu subscription.

## 4. Thứ tự xử lý đề xuất

| Thứ tự | Nhóm | Kết quả cụ thể cần có |
|---|---|---|
| 1 | C01 | Đường phục hồi accepted cho lời mời và dữ liệu gửi, cùng TTL bí mật được chốt. |
| 2 | C02–C04, C13 | Ma trận trạng thái, callback matching/dedupe, version kết nối và cửa sổ đối chiếu. |
| 3 | C05, C08–C09, C14 | Vòng đời token/credential/PII, tombstone, mốc retention, validation và thế hệ khóa. |
| 4 | C06–C07, C10, C12 | Invariant hạn mức, audit đủ lịch sử, hủy idempotent và chống trùng các đối tượng con. |
| 5 | C11, C15 | ERD/dictionary/baseline đồng nhất để Tech Lead chuyển sang Physical DB. |

Mỗi nhóm chỉ chặn phần triển khai phụ thuộc của nó; có thể tiếp tục những phần đã đủ ngữ nghĩa. Báo cáo không sửa nội dung thiết kế gốc hoặc tự phê duyệt phương án mới.

## 5. Cách dùng các showcase khi có triển khai

Với từng showcase ở trên, QA ghi dữ liệu đầu vào, thời điểm giả lập, thao tác đồng thời/crash nếu có, bản ghi trước/sau, response/event nhận được và kết quả assertion. Đặc biệt cần chạy phục hồi sau accepted, replay song song, callback đến sớm/đến muộn, xoay khóa, thu hồi/xóa token và xóa dữ liệu tại ranh giới ngày/tháng.

Các AT tham chiếu là nơi bổ sung test case, **không có AT nào được báo cáo này đánh dấu đã chạy hoặc đã pass**. Đánh giá thiết kế có thể hoàn tất từ tài liệu; chứng nhận triển khai cần bằng chứng thực thi riêng.
