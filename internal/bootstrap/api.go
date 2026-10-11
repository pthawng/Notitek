package bootstrap

import (
	"log/slog"
	"net/http"

	"github.com/pthawng/notitek/internal/platform/apperr"
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

	// Unmatched paths answer with the standard error envelope.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteError(w, r, logger, apperr.New(apperr.NotFound, "Không tìm thấy tài nguyên."))
	})

	handler := httpx.Chain(mux,
		httpx.Correlation(),
		httpx.Recover(logger),
		httpx.AccessLog(logger),
		httpx.NoStore(),
		httpx.MaxBody(maxRequestBody),
	)
	return []Component{httpx.NewServer("public", cfg.HTTP, handler, logger)}, nil
}
