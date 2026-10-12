package httpx

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/pthawng/notitek/internal/platform/apperr"
)

// fallbackPattern is the catch-all route that RegisterFallback installs.
const fallbackPattern = "/"

// probeMethods are tried, in this order, to build the Allow header.
var probeMethods = []string{
	http.MethodGet, http.MethodHead, http.MethodPost,
	http.MethodPut, http.MethodPatch, http.MethodDelete,
}

// RegisterFallback answers every request that no route handles with the
// standard error envelope instead of net/http's plain text:
//
//   - the path exists for other methods: 405 METHOD_NOT_ALLOWED with an Allow
//     header (RFC 9110 §15.5.6)
//   - the path does not exist: 404 NOT_FOUND
//
// Register all routes first or after; the order does not matter.
func RegisterFallback(mux *http.ServeMux, logger *slog.Logger) {
	mux.HandleFunc(fallbackPattern, func(w http.ResponseWriter, r *http.Request) {
		if allowed := allowedMethods(mux, r); len(allowed) > 0 {
			w.Header().Set("Allow", strings.Join(allowed, ", "))
			WriteError(w, r, logger, apperr.New(apperr.MethodNotAllowed,
				"Phương thức không được hỗ trợ cho đường dẫn này. Phương thức được phép: "+strings.Join(allowed, ", ")+"."))
			return
		}
		WriteError(w, r, logger, apperr.New(apperr.NotFound, "Không tìm thấy tài nguyên."))
	})
}

// allowedMethods lists the methods that have a real route for r's path. The
// catch-all itself does not count.
func allowedMethods(mux *http.ServeMux, r *http.Request) []string {
	var allowed []string
	for _, m := range probeMethods {
		probe := r.Clone(r.Context())
		probe.Method = m
		if _, pattern := mux.Handler(probe); pattern != "" && pattern != fallbackPattern {
			allowed = append(allowed, m)
		}
	}
	return allowed
}
