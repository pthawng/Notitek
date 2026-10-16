# Registry thông báo

Nguồn sự thật của mọi loại thông báo. Thay đổi đi qua pull request (DEC-07); `config_version` là hash của toàn bộ thư mục này và mỗi phiên bản được lưu snapshot bất biến trong DB.

```text
types/        <type>.v<n>.yaml   owner, source được phép, mode, purpose, context,
                                 JSON Schema của recipients/data/secrets, hạn gửi
policies/     kênh bắt buộc/tùy chọn/thay thế, retry, khung giờ, tiêu chí đạt
templates/vi/<kênh>/             nội dung theo type và kênh
audiences.yaml                   selector được phép và quyền tương ứng
sources.yaml                     client → source, type/application/brand được phép
overrides/{application,brand,organization}/   ba cấp Global → App/Brand → Shop
```

Dữ liệu vận hành không nằm ở đây mà ở DB, đổi qua API vận hành có audit và không cần deploy (DEC-07 v1.7): kết nối kênh (provider, sender/OA/domain, credential), kill-switch, hạn mức ngày, bật/tắt Push của Shop. Secret (khóa gốc, khóa HMAC) không bao giờ nằm trong thư mục này.

Định dạng chi tiết được chốt khi hiện thực module `internal/registry`; 17 profile R1 hiện được mô tả trong [`docs/api/openapi.json`](../../docs/api/openapi.json) và sẽ được chuyển sang đây, sau đó OpenAPI sinh `x-type-profiles` từ registry.
