package httpx

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"runtime/debug"
	"time"

	"github.com/pthawng/notitek/internal/platform/apperr"
	"github.com/pthawng/notitek/internal/platform/idgen"
	"github.com/pthawng/notitek/internal/platform/observability"
)

// Middleware wraps a handler.
type Middleware func(http.Handler) http.Handler

// Chain applies middlewares so the first one is the outermost.
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

var correlationPattern = regexp.MustCompile(`^[A-Za-z0-9._:-]{1,100}$`)

// Correlation reads X-Correlation-Id (or generates one), stores it in the
// request context and echoes it on the response. Malformed values are
// replaced rather than propagated into logs.
func Correlation() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get(observability.CorrelationHeader)
			if !correlationPattern.MatchString(id) {
				if gen, err := idgen.New(); err == nil {
					id = gen
				} else {
					id = "unavailable"
				}
			}
			w.Header().Set(observability.CorrelationHeader, id)
			next.ServeHTTP(w, r.WithContext(observability.WithCorrelationID(r.Context(), id)))
		})
	}
}

// NoStore forbids caching of every response; notification data may contain
// personal information.
func NoStore() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store")
			next.ServeHTTP(w, r)
		})
	}
}

// SecurityHeaders sets the response headers that harden a JSON API: no MIME
// sniffing, no referrer, and a policy that forbids rendering any response as
// a document.
func SecurityHeaders() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()
			h.Set("X-Content-Type-Options", "nosniff")
			h.Set("Referrer-Policy", "no-referrer")
			h.Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
			next.ServeHTTP(w, r)
		})
	}
}

// MaxBody limits request bodies to n bytes (64 KiB per SRS-P04).
func MaxBody(n int64) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, n)
			next.ServeHTTP(w, r)
		})
	}
}

// Recover turns a panic into a 500 envelope and logs the stack.
func Recover(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					handlePanic(w, r, logger, rec)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func handlePanic(w http.ResponseWriter, r *http.Request, logger *slog.Logger, rec any) {
	// http.ErrAbortHandler is the sanctioned way to abort a response; let the
	// server handle it.
	if err, ok := rec.(error); ok && errors.Is(err, http.ErrAbortHandler) {
		panic(rec)
	}
	logger.ErrorContext(r.Context(), "panic in handler", "panic", rec, "stack", string(debug.Stack()))
	writeErrorEnvelope(w, r, apperr.From(fmt.Errorf("panic: %v", rec)))
}

// AccessLog records method, route pattern, status and latency. It never logs
// bodies, query strings or headers, which may hold personal data.
func AccessLog(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rec, r)

			route := r.Pattern
			if route == "" {
				route = "unmatched"
			}
			logger.InfoContext(r.Context(), "http request",
				"method", r.Method,
				"route", route,
				"status", rec.status,
				"bytes", rec.bytes,
				"duration_ms", time.Since(start).Milliseconds(),
			)
		})
	}
}

type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int64
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.bytes += int64(n)
	return n, err
}

// Unwrap lets http.ResponseController reach the underlying writer (needed
// for flushing SSE streams).
func (r *statusRecorder) Unwrap() http.ResponseWriter { return r.ResponseWriter }
