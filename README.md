# Notitek

Dịch vụ thông báo dùng chung của SuperPlatform. Các module (User, Order, …) gửi **ý định thông báo** qua một API chung `/v1/notifications`; Notitek chọn mẫu, kênh, người nhận, thử lại và ghi bằng chứng gửi.

- Hợp đồng API: [`docs/api/README.md`](docs/api/README.md), [`docs/api/openapi.json`](docs/api/openapi.json)
- Quyết định nghiệp vụ: [`docs/Notitek - QUYẾT ĐỊNH PO CHO CÁC DEC.md`](<docs/Notitek - QUYẾT ĐỊNH PO CHO CÁC DEC.md>)
- Cấu trúc mã và quy tắc kiến trúc: [`docs/Notitek - CẤU TRÚC MÃ NGUỒN.md`](<docs/Notitek - CẤU TRÚC MÃ NGUỒN.md>)

## Runtime

| Binary | Vai trò | Hạ tầng dùng |
|---|---|---|
| `cmd/api` | Nhận lệnh, OTP đồng bộ, hộp tin/SSE, webhook, tra cứu | PostgreSQL, Redis |
| `cmd/worker` | Gửi nền theo kênh | PostgreSQL, RabbitMQ, Redis |
| `cmd/relay` | Outbox → RabbitMQ (publisher confirms) | PostgreSQL, RabbitMQ |
| `cmd/scheduler` | Hạn gửi, khung giờ, đóng `UNKNOWN`, thời gian lưu | PostgreSQL |

Mỗi runtime mở cổng admin (mặc định `:9090`) với `/healthz` và `/readyz`. `api` phục vụ API công khai ở `:8080`.

## Chạy local

Yêu cầu: Go 1.26, Docker, Make (Git Bash hoặc WSL trên Windows), Python 3.13 cho kiểm tra hợp đồng.

```bash
make up                      # PostgreSQL, RabbitMQ, Redis
make migrate-up              # áp dụng migration
cp .env.example .env
set -a; source .env; set +a
go run ./cmd/api             # hoặc worker, relay, scheduler
curl localhost:9090/readyz
```

Chạy cả bốn runtime trong container: `make up-app`. Dừng: `make down`; xóa dữ liệu: `make reset`.

Nếu port trên máy đã bị dùng, đổi port host qua biến môi trường: `NOTITEK_POSTGRES_PORT`, `NOTITEK_RABBITMQ_PORT`, `NOTITEK_REDIS_PORT`, `NOTITEK_API_PORT`, `NOTITEK_API_ADMIN_PORT` (ví dụ `NOTITEK_POSTGRES_PORT=55432 make up`, và sửa DSN tương ứng). Mất kết nối RabbitMQ làm runtime dừng với lỗi để orchestrator khởi động lại.

## Hot reload khi phát triển

```bash
make dev                                   # hạ tầng + migration + dev-api, tự build lại khi sửa code
make dev SERVICES="dev-api dev-worker"     # thêm runtime đang làm
make dev-logs
make dev-down                              # dừng, giữ dữ liệu và cache Go
```

- Source được mount vào container; [`air`](.air.toml) theo dõi `cmd/`, `internal/`, `config/` (bỏ qua `_test.go`) và khởi động lại runtime sau khi build xong, khoảng 3 giây cho một runtime.
- Mọi runtime dùng chung một cache module và một cache build, nên dependency chỉ tải một lần và package không đổi không build lại. Binary và log build nằm trong tmpfs, không ghi vào repo.
- Dùng polling vì bind mount từ Windows/macOS không báo sự kiện file ổn định; mỗi runtime tốn khoảng 3–6% CPU khi rảnh. Chỉ bật runtime đang làm.
- Migration chạy tự động (service `migrate`) trước khi runtime khởi động. Lần đầu build image dev và tải module mất vài phút.
- Không chạy cùng lúc với `make up-app` vì trùng port. Sửa `go.mod` thì chạy lại `make dev`.

## Kiểm tra

```bash
make lint              # golangci-lint, gồm luật phụ thuộc giữa module (depguard)
make test              # unit test
make test-race         # cần cgo (Linux/macOS, hoặc Windows có gcc)
make test-integration  # cần `make up`
make contract          # OpenAPI, ví dụ và các ca dữ liệu sai
make check             # lint + test + contract
```

## Cấu hình

Toàn bộ qua biến môi trường tiền tố `NOTITEK_`, xem [`.env.example`](.env.example). Cấu hình sai hoặc thiếu làm runtime dừng ngay khi khởi động và liệt kê mọi lỗi cùng lúc. Mỗi runtime chỉ yêu cầu hạ tầng nó dùng. Secret không bao giờ nằm trong repo.

## Trạng thái

Khung dự án (M0): composition root, vòng đời và tắt mềm, cấu hình, log có che dữ liệu, envelope lỗi theo danh mục `x-error-codes`, probe, PostgreSQL `TxManager`, RabbitMQ publisher có confirms, Redis, quy tắc tổng hợp outcome và phân loại kết quả kênh.

Tiếp theo: ADR-03/04 → migration các bảng → registry (`config/registry`) → tiếp nhận `POST /v1/notifications` → outbox/relay → worker và adapter kênh (mock trước, provider thật sau CF-12).
