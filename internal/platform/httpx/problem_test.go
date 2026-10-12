package httpx_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/pthawng/notitek/internal/platform/apperr"
	"github.com/pthawng/notitek/internal/platform/httpx"
	"github.com/pthawng/notitek/internal/platform/observability"
)

func writeErr(t *testing.T, err error, target string) (*httptest.ResponseRecorder, map[string]any) {
	t.Helper()
	h := httpx.Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteError(w, r, discard, err)
	}), httpx.Correlation())
	req := httptest.NewRequest(http.MethodPost, target, nil)
	req.Header.Set(observability.CorrelationHeader, "corr-1")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	var env map[string]any
	if e := json.Unmarshal(rec.Body.Bytes(), &env); e != nil {
		t.Fatalf("decode: %v\n%s", e, rec.Body.String())
	}
	return rec, env
}

// TestErrorEnvelopeCarriesProblemDetails checks the RFC 9457 members and the
// extensions of an error response.
func TestErrorEnvelopeCarriesProblemDetails(t *testing.T) {
	err := apperr.New(apperr.InvalidRequest, "Có trường sai kiểu dữ liệu.",
		apperr.FieldError(apperr.LocationBody, "/recipients/0/phone", apperr.ReasonInvalidType, "Trường phải có kiểu string."))
	rec, env := writeErr(t, err, "/v1/notifications?secret=abc")

	if env["error"] != true || env["message"] != "Có trường sai kiểu dữ liệu." {
		t.Fatalf("envelope = %v", env)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %q", got)
	}
	data := env["data"].(map[string]any)
	want := map[string]any{
		"type":      "urn:supership:notitek:error:invalid-request",
		"title":     "Invalid request",
		"status":    float64(400),
		"code":      "INVALID_REQUEST",
		"instance":  "/v1/notifications",
		"trace_id":  "corr-1",
		"retryable": false,
	}
	for k, v := range want {
		if data[k] != v {
			t.Errorf("data.%s = %v, want %v", k, data[k], v)
		}
	}
	if strings.Contains(rec.Body.String(), "secret=abc") {
		t.Error("the query string must not be echoed")
	}
	if _, has := data["retry_after_seconds"]; has {
		t.Error("a non-retryable error must not carry retry_after_seconds")
	}
	if rec.Header().Get("Retry-After") != "" {
		t.Error("a non-retryable error must not set Retry-After")
	}

	details := data["details"].([]any)
	if len(details) != 1 {
		t.Fatalf("details = %v", details)
	}
	d := details[0].(map[string]any)
	for k, v := range map[string]string{
		"field": "/recipients/0/phone", "location": "body", "code": "invalid_type", "message": "Trường phải có kiểu string.",
	} {
		if d[k] != v {
			t.Errorf("details[0].%s = %v, want %v", k, d[k], v)
		}
	}
}

func TestRetryableErrorsAdvertiseWhenToRetry(t *testing.T) {
	cases := map[string]struct {
		err        *apperr.Error
		wantHeader string
	}{
		"default delay":    {apperr.New(apperr.RateLimited, "Quá tải."), "1"},
		"explicit rounded": {apperr.New(apperr.ServiceUnavailable, "Tạm dừng.").WithRetryAfter(2500 * time.Millisecond), "3"},
		"in progress":      {apperr.New(apperr.IdempotencyInProgress, "Đang xử lý.").WithRetryAfter(time.Second), "1"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			rec, env := writeErr(t, c.err, "/v1/notifications")
			data := env["data"].(map[string]any)
			if data["retryable"] != true {
				t.Errorf("retryable = %v", data["retryable"])
			}
			if got := rec.Header().Get("Retry-After"); got != c.wantHeader {
				t.Errorf("Retry-After = %q, want %q", got, c.wantHeader)
			}
			if got := data["retry_after_seconds"]; got != mustFloat(t, c.wantHeader) {
				t.Errorf("retry_after_seconds = %v, want %s", got, c.wantHeader)
			}
		})
	}
}

func mustFloat(t *testing.T, s string) float64 {
	t.Helper()
	var f float64
	if err := json.Unmarshal([]byte(s), &f); err != nil {
		t.Fatal(err)
	}
	return f
}

func TestEmptyMessageFallsBackToTitle(t *testing.T) {
	_, env := writeErr(t, apperr.New(apperr.Forbidden, ""), "/v1/notifications")
	if env["message"] != "Forbidden" {
		t.Errorf("message = %v", env["message"])
	}
}

// TestRouterAnswersWithEnvelopes covers 404 and 405 for unmatched requests.
func TestFallbackAnswersWithEnvelope(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/things", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.HandleFunc("POST /v1/things", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusAccepted) })
	httpx.RegisterFallback(mux, discard)
	h := httpx.Chain(mux, httpx.Correlation())

	cases := []struct {
		method, path string
		status       int
		code         string
		allow        string
	}{
		{http.MethodGet, "/v1/things", http.StatusOK, "", ""},
		{http.MethodPost, "/v1/things", http.StatusAccepted, "", ""},
		{http.MethodDelete, "/v1/things", http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "GET, HEAD, POST"},
		{http.MethodPut, "/v1/things", http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "GET, HEAD, POST"},
		{http.MethodGet, "/v1/missing", http.StatusNotFound, "NOT_FOUND", ""},
		{http.MethodPost, "/", http.StatusNotFound, "NOT_FOUND", ""},
	}
	for _, c := range cases {
		t.Run(c.method+" "+c.path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(c.method, c.path, nil))
			if rec.Code != c.status {
				t.Fatalf("status = %d, want %d", rec.Code, c.status)
			}
			if c.code == "" {
				return
			}
			env := decode(t, rec)
			if !env.Error || env.Data.Code != c.code {
				t.Errorf("envelope = %+v", env)
			}
			if got := rec.Header().Get("Allow"); got != c.allow {
				t.Errorf("Allow = %q, want %q", got, c.allow)
			}
		})
	}
}

func TestSecurityHeaders(t *testing.T) {
	h := httpx.Chain(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}), httpx.SecurityHeaders())
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	for k, want := range map[string]string{
		"X-Content-Type-Options": "nosniff",
		"Referrer-Policy":        "no-referrer",
	} {
		if got := rec.Header().Get(k); got != want {
			t.Errorf("%s = %q, want %q", k, got, want)
		}
	}
	if !strings.Contains(rec.Header().Get("Content-Security-Policy"), "default-src 'none'") {
		t.Error("Content-Security-Policy missing")
	}
}

// TestAccessLogSeesPanicAs500 guards the middleware order: AccessLog outside
// Recover logs the 500 that a panic becomes, with the correlation id.
func TestAccessLogSeesPanicAs500(t *testing.T) {
	var logs bytes.Buffer
	logger := observability.NewLogger(&logs, observability.LoggerConfig{Level: "debug", Runtime: "api", Env: "test"})
	h := httpx.Chain(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("boom") }),
		httpx.Correlation(), httpx.AccessLog(logger), httpx.Recover(slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))))

	req := httptest.NewRequest(http.MethodGet, "/v1/x?token=hunter2", nil)
	req.Header.Set(observability.CorrelationHeader, "corr-panic")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", rec.Code)
	}
	line := logs.String()
	for _, want := range []string{`"status":500`, `"correlation_id":"corr-panic"`, `"bytes":`} {
		if !strings.Contains(line, want) {
			t.Errorf("access log lacks %s: %s", want, line)
		}
	}
	if strings.Contains(line, "hunter2") {
		t.Errorf("access log leaks the query string: %s", line)
	}
}
