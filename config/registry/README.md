# Registry thông báo

Nguồn sự thật của mọi loại thông báo. Thay đổi đi qua pull request (DEC-07); `config_version` là hash của toàn bộ thư mục này và mỗi phiên bản được lưu snapshot bất biến trong DB.

```text
types/        <type>.v<n>.yaml   owner, source được phép, mode, purpose, context,
                                 JSON Schema của recipients/data/secrets, hạn gửi
policies/     kênh bắt buộc/tùy chọn/thay thế, retry, khung giờ, tiêu chí đạt
templates/vi/<kênh>/             nội dung theo type và kênh
audiences.yaml                   selector được phép và quyền tương ứng
sources.yaml                     client → source, type/application/brand được phép
channels.yaml                    định tuyến provider, sender/OA/domain theo brand
overrides/{application,brand,organization}/   ba cấp Global → App/Brand → Shop
```

Kill-switch không nằm ở đây: đó là trạng thái runtime trong DB, đổi qua API vận hành có audit. Secret (khóa provider, khóa HMAC) không bao giờ nằm trong thư mục này.

Định dạng chi tiết được chốt khi hiện thực module `internal/registry`; 17 profile R1 hiện được mô tả trong [`docs/api/openapi.json`](../../docs/api/openapi.json) và sẽ được chuyển sang đây, sau đó OpenAPI sinh `x-type-profiles` từ registry.
