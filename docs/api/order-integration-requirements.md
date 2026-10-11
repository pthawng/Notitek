# Yêu cầu tích hợp Notitek cho module Order

Ngày 11/10/2026. Trạng thái: **Bản nháp gửi đội Order** (Tech Lead Order, PO Order). Căn cứ: [DEC 1.6](<../Notitek - QUYẾT ĐỊNH PO CHO CÁC DEC.md>) mục 2.13–2.15 (A1, A2, A7, A8, B10, Q2, Q6, Q8), [hợp đồng API v1](README.md), [OpenAPI](openapi.json). Tài liệu liệt kê việc đội Order cần làm trong repo `order-spf` và đề xuất văn bản thay thế §11 của `docs/SuperShip - MODULE ORDER - MOCK API CONTRACT CÁC MODULE TƯƠNG TÁC.md`; Notitek không sửa repo này.

## 1. Thay đổi so với §11 hiện tại

| Nội dung | §11 hiện tại | Theo DEC 1.6 |
|---|---|---|
| Endpoint | `POST /v1/notifications/events/orders` | `POST /v1/notifications` với profile `order.*` (A2) |
| Đường gọi | Order gọi sau commit | Dispatcher outbox của Order gọi REST; broker để phân kỳ sau |
| Liên hệ người nhận | Không gửi | **Bắt buộc ở `order.delivery_failed`**, tùy chọn ở 9 mốc còn lại và nên bỏ khi mốc không nhắn người nhận (Q2); Notitek không gọi ngược Order (A1) |
| Tên sự kiện | `event_type = ORDER_*` | `type = order.*`, `type_version = 1` |
| Người nhận nội bộ | Không có | Selector `ORDER_OPS_ALERT_RECEIVERS` tùy chọn (A8) |
| Hủy khi đã giải quyết | Không có | `POST /v1/notifications/cancel` theo `cancel_group` |
| Tổng quan theo đơn | `GET /v1/orders/{order_code}/notification-summary` do Notification phục vụ | Order tự phục vụ: kiểm Data Scope của người dùng trên đơn rồi gọi `GET /v1/notifications?subject_ref={order_id}` bằng client của Order (Q6), hoặc bỏ |

## 2. Ánh xạ tên sự kiện

`type` của Notitek là tên chuẩn ở ranh giới tích hợp. `event_type` trong `outbox_events` là tên nội bộ của Order; Order giữ bảng ánh xạ dưới đây.

| Trạng thái nguồn | Outbox `event_type` (đề xuất) | §11 cũ | Notitek `type` |
|---|---|---|---|
| SPF-0301 | `order.created.v1` (đã có, sai thời điểm) | `ORDER_CREATED` | `order.created` |
| SPF-0201 | `order.cancelled.v1` (đã có) | `ORDER_CANCELLED` | `order.cancelled` |
| SPF-0402 | `order.pickup_failed.v1` (mới) | `ORDER_PICKUP_FAILED` | `order.pickup_failed` |
| SPF-0501 | `order.picked_up.v1` (mới) | `ORDER_PICKED_UP` | `order.picked_up` |
| SPF-0802 | `order.delivery_failed.v1` (mới) | `ORDER_DELIVERY_FAILED` | `order.delivery_failed` |
| SPF-0901 | `order.delivered.v1` (mới) | `ORDER_DELIVERED` | `order.delivered` |
| SPF-0902 | `order.partially_delivered.v1` (mới) | `ORDER_PARTIALLY_DELIVERED` | `order.partially_delivered` |
| SPF-1002 | `order.return.confirmed.v1` (đã có) | `ORDER_RETURN_CONFIRMED` | `order.return_confirmed` |
| SPF-1201/1202/1203 | `order.returned.v1` (mới) | `ORDER_RETURNED` | `order.returned` |
| Không đổi trạng thái | `order.sla_breached.v1` (mới) | `ORDER_SLA_BREACHED` | `order.sla_breached` |

`order.return.requested.v1` hiện có không phải mốc thông báo; không gửi sang Notitek.

## 3. Việc Order cần làm (CF-04, hạn 15/10/2026)

| ID | Việc | Yêu cầu | Kiểm tra |
|---|---|---|---|
| O-01 | Dispatcher outbox | Claim row `PENDING`, gọi `POST /v1/notifications` bằng OAuth2 client credentials (CF-09). 202/200 → `SENT`. Timeout, 429, 5xx → `RETRY_WAIT`, retry **nguyên body** cùng `source_ref` trong `send_deadline`. 409 `IDEMPOTENCY_CONFLICT`, 400/401/403/422 → `DEAD_LETTER` để xử lý tay | Kill dispatcher sau khi gọi → retry trả cùng `notification_id`, không gửi trùng |
| O-02 | Outbox cho đủ 10 mốc | Ghi row cùng transaction đổi trạng thái, **kể cả mốc đến từ webhook tracking** (`applyTrackingEvent`), hiện chưa phát outbox | Mỗi lần vào trạng thái nguồn có đúng một row |
| O-03 | `order.created` đúng thời điểm | Chỉ ghi khi đơn vào SPF-0301 sau khi có kết quả tạo vận đơn; không ghi ở "TRANSACTION A" như hiện nay; tạo vận đơn lỗi (SPF-0102) không phát | Đơn tạo vận đơn lỗi không sinh thông báo |
| O-04 | `occurred_at` nghiệp vụ | Lưu thời điểm nghiệp vụ (thời điểm NVC xác nhận, không phải `created_at` của row) có offset; giữ nguyên qua retry | Payload có offset, không đổi giữa các lần retry |
| O-05 | `source_ref` ổn định | `{order_id}:{type}` cho mốc một lần; `{order_id}:{type}:{failure_or_breach_id}` cho mốc lặp lại được (`order.pickup_failed`, `order.delivery_failed`, `order.sla_breached`). **Đổi định dạng `event_key`** của outbox theo đúng mẫu này (dạng hiện tại `{orderId}:{eventType}` có UNIQUE sẽ trùng ở lần thất bại thứ hai). Không dùng `outbox_event_id` ngẫu nhiên; ≤128 ký tự, `[A-Za-z0-9._:-]` | Ghi lại row cùng sự việc không tạo yêu cầu mới; lần giao thất bại thứ hai tạo yêu cầu mới |
| O-06 | Liên hệ người nhận | Lấy `order_parties` (party_type 2) hiệu lực tại thời điểm sự việc: `phone` E.164 bắt buộc, `email` tùy chọn. **Bắt buộc** ở `order.delivery_failed`; ở 9 mốc khác chỉ gửi khi chính sách mốc nhắn người nhận (R1: không), để tối thiểu dữ liệu (DPO, CF-01) | `delivery_failed` có đúng một CONTACT; mốc khác không có CONTACT trong R1 |
| O-07 | `subject_ref` | Bắt buộc với `order.*`, giá trị `order_id`. Khi trả thông tin thông báo của đơn cho FE, Order **kiểm Data Scope của người dùng trên đơn trước**, rồi mới gọi `GET /v1/notifications?subject_ref=` bằng client của Order (không cần `source`, lấy từ token) | Tra ra đủ thông báo của đơn; người ngoài Data Scope bị Order từ chối trước khi gọi Notitek |
| O-08 | Hủy theo nhóm | Lần lấy/giao thất bại được giải quyết (giao lại thành công, chuyển hoàn, hủy đơn) hoặc SLA không còn vi phạm → ghi outbox lệnh `POST /v1/notifications/cancel` với `cancel_group` tương ứng, `Idempotency-Key` riêng | Tin nhắc lần thất bại chưa gửi bị dừng |
| O-09 | Phát hiện vi phạm SLA | Job so `expected_to` với thời điểm hiện tại, phát `order.sla_breached` một lần cho mỗi lần vi phạm (`sla_breach_id`) | Không phát lặp cho cùng lần vi phạm |
| O-10 | Selector vận hành | Gửi `ORDER_OPS_ALERT_RECEIVERS` ở mốc cần vận hành theo dõi (đã chốt: `pickup_failed`, `delivery_failed`, `sla_breached`) | Vận hành ngoài data scope Shop không nhận |
| O-11 | Tài liệu | Thay §11 theo mục 4; ghi bảng ánh xạ mục 2 vào tài liệu Order; sửa ADR 0008 ("tracking chưa phát outbox") khi O-02 xong | Review tài liệu |

## 4. Văn bản đề xuất thay §11

> ## 11. Notification
>
> Order thông báo các mốc nghiệp vụ bằng API chung của Notification `POST /v1/notifications` theo hợp đồng Notitek v1 (notification-spf: `docs/api/README.md`). Order ghi ý định thông báo vào `outbox_events` cùng transaction đổi trạng thái; dispatcher outbox gọi API sau commit. Order gửi ngữ cảnh đơn, selector người nhận và liên hệ người nhận hàng (bắt buộc ở giao thất bại, chỉ gửi ở mốc khác khi mốc có nhắn người nhận); Notification chọn mẫu, kênh, lịch gửi, retry và không gọi ngược Order. `202 Accepted` chỉ xác nhận đã tiếp nhận bền vững, không khẳng định đã gửi.
>
> **Danh sách mốc:** chỉ 10 `type` tại bảng ánh xạ mục 2. Không phát cho trạng thái trung gian (nhập kho, bàn giao, trung chuyển, tìm tài xế, đang giao) và không dùng sự kiện tổng quát kiểu `ORDER_STATUS_CHANGED`.
>
> ### 11.1. Gửi thông báo mốc Order
>
> | Thuộc tính | Nội dung |
> |---|---|
> | Phương thức | `POST` |
> | Endpoint | `/v1/notifications` |
> | Quyền | Client của Order Service (OAuth2 client credentials), chỉ được dùng `type` `order.*` |
> | Thành công | `202 Accepted` (mới), `200 OK` (replay cùng body) |
>
> Headers: `Authorization: Bearer {access_token}`, `Content-Type: application/json`, `Idempotency-Key` (bằng `source_ref`), `X-Correlation-Id` (tùy chọn).
>
> | Trường | Bắt buộc | Ý nghĩa |
> |---|---|---|
> | `source_ref` | Có | Khóa ổn định của sự việc = `event_key`: `{order_id}:{type}`, hoặc `{order_id}:{type}:{failure_or_breach_id}` với mốc lặp lại |
> | `type`, `type_version` | Có | `order.*`, `1` |
> | `mode` | Có | `async` |
> | `occurred_at` | Có | Thời điểm nghiệp vụ, ISO 8601 có offset |
> | `subject_ref` | Có | `order_id` |
> | `context` | Có | `context_type=WORKSPACE`, `application_id`, `organization_id` của Shop |
> | `recipients` | Có | Đúng một `SHOP_ORDER_ALERT_RECEIVERS`; CONTACT người nhận hàng (`phone` bắt buộc, `email` tùy chọn) **đúng một ở `order.delivery_failed`, tối đa một ở mốc khác**; tối đa một `ORDER_OPS_ALERT_RECEIVERS` (gửi ở `pickup_failed`, `delivery_failed`, `sla_breached`) |
> | `data.order_id`, `data.order_code`, `data.order_status` | Có | Mã đơn, mã 14 chữ số, trạng thái `SPF-xxxx` |
> | `data.status_reason_code`, `data.tracking_url` | Không | Mã lý do; link tracking công khai HTTPS |
> | Trường riêng | Theo mốc | `pickup_failure_id`+`reason_code`; `delivery_failure_id`+`reason_code`; `sla_breach_id`+`sla_code`+`expected_to` |
> | `send_deadline` | Có | Mặc định `occurred_at` + 12h |
> | `cancel_group` | Theo mốc | Bắt buộc với `pickup_failed`, `delivery_failed`, `sla_breached` (ID lần thất bại/vi phạm); không gửi với mốc khác |
>
> Lỗi theo [README §10](README.md): 400 `INVALID_REQUEST`/`REFERENCE_MISMATCH`, 401/403, 409 `IDEMPOTENCY_CONFLICT`, 422 `UNKNOWN_TYPE`/`INVALID_PROFILE`/`STALE_REQUEST`, 429 `RATE_LIMITED`, 503 `SERVICE_UNAVAILABLE`.
>
> ### 11.2. Hủy thông báo khi lần thất bại được giải quyết
>
> `POST /v1/notifications/cancel`, `Idempotency-Key` riêng, body `{cancel_group, context, reason_code}`. Notification dừng phần chưa gửi, giữ dấu hủy 72h kể cả khi yêu cầu gửi đến sau.
>
> ### 11.3. Tổng quan thông báo của Order
>
> FE gọi API của Order; Order kiểm Data Scope của người dùng trên đơn, rồi gọi `GET /v1/notifications?subject_ref={order_id}&cursor=&page_size=` bằng client của Order (nguồn lấy từ token) và tự tổng hợp số lượng theo `outcome`. FE không gọi thẳng Notification. Nếu không cần trên màn hình chi tiết Order ở R1 thì bỏ mục này.

## 5. Ví dụ request

Dữ liệu giả, theo [order-delivery-failed.request.json](examples/order-delivery-failed.request.json) cộng `subject_ref`:

```http
POST /v1/notifications
Authorization: Bearer <order-service-access-token>
Content-Type: application/json
Idempotency-Key: order-demo:order.delivery_failed:failure-demo
X-Correlation-Id: corr-demo-001
```

```json
{
  "source_ref": "order-demo:order.delivery_failed:failure-demo",
  "type": "order.delivery_failed",
  "type_version": 1,
  "mode": "async",
  "occurred_at": "2026-10-11T10:00:00+07:00",
  "subject_ref": "order-demo",
  "context": {
    "context_type": "WORKSPACE",
    "application_id": "app-shop-web",
    "organization_id": "org-shop-demo"
  },
  "recipients": [
    { "kind": "AUDIENCE", "recipient_ref": "shop-team", "audience_code": "SHOP_ORDER_ALERT_RECEIVERS" },
    { "kind": "CONTACT", "recipient_ref": "consignee-demo", "phone": "+84900000000", "email": "consignee@example.invalid" }
  ],
  "data": {
    "order_id": "order-demo",
    "order_code": "12345678901231",
    "order_status": "SPF-0802",
    "tracking_url": "https://tracking.example.invalid/orders/12345678901231",
    "delivery_failure_id": "failure-demo",
    "reason_code": "RECIPIENT_UNAVAILABLE"
  },
  "send_deadline": "2026-10-11T22:00:00+07:00",
  "cancel_group": "failure-demo"
}
```
