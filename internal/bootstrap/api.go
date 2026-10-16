package bootstrap

import (
	"log/slog"
	"net/http"

	"github.com/pthawng/notitek/internal/platform/httpx"
)

// maxRequestBody is the 64 KiB JSON limit of SRS-P04 / API v1 §10.
const maxRequestBody = 64 << 10

func assembleAPI(cfg Config, _ *Dependencies, logger *slog.Logger) ([]Component, error) {
	// Probes live on the admin port (lifecycle.go), not on the public listener.
	mux := http.NewServeMux()

	// Module routes are registered here as modules are implemented, e.g.
	// notification.NewModule(...).RegisterRoutes(mux). Fixed sub-paths
	// (inbox, push-subscriptions, cancel, ops) must be registered before
	// /v1/notifications/{notification_id}; net/http picks the most specific
	// pattern, which preserves that precedence.

	// Unmatched paths and methods answer with the standard error envelope
	// (404, or 405 with Allow).
	httpx.RegisterFallback(mux, logger)

	// Order matters: Correlation first so every log line and error carries the
	// id; AccessLog outside Recover so a panic is logged as the 500 it becomes.
	handler := httpx.Chain(mux,
		httpx.Correlation(),
		httpx.AccessLog(logger),
		httpx.Recover(logger),
		httpx.SecurityHeaders(),
		httpx.NoStore(),
		httpx.MaxBody(maxRequestBody),
	)
	return []Component{httpx.NewServer("public", cfg.HTTP, handler, logger)}, nil
}
