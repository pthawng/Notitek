package httpx_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pthawng/notitek/internal/platform/apperr"
	"github.com/pthawng/notitek/internal/platform/httpx"
	"github.com/pthawng/notitek/internal/platform/observability"
)

var discard = slog.New(slog.NewTextHandler(io.Discard, nil))

type errorEnvelope struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    struct {
		Code    string          `json:"code"`
		TraceID string          `json:"trace_id"`
		Details []apperr.Detail `json:"details"`
	} `json:"data"`
}

func decode(t *testing.T, rec *httptest.ResponseRecorder) errorEnvelope {
	t.Helper()
	var env errorEnvelope
	if err := json.NewDecoder(rec.Body).Decode(&env); err != nil {
		t.Fatalf("decode envelope: %v", err)
	}
	return env
}

func TestCorrelationEchoesValidAndReplacesInvalid(t *testing.T) {
	var seen string
	h := httpx.Chain(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		seen = observability.CorrelationID(r.Context())
	}), httpx.Correlation())

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(observability.CorrelationHeader, "corr-demo-001")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if seen != "corr-demo-001" || rec.Header().Get(observability.CorrelationHeader) != "corr-demo-001" {
		t.Fatalf("valid id not propagated: ctx=%q header=%q", seen, rec.Header().Get(observability.CorrelationHeader))
	}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(observability.CorrelationHeader, "bad id\nwith newline")
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if got := rec.Header().Get(observability.CorrelationHeader); got == "" || strings.Contains(got, " ") {
		t.Fatalf("invalid id must be replaced, got %q", got)
	}
}

func TestWriteErrorUsesEnvelopeAndHidesInternals(t *testing.T) {
	h := httpx.Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteError(w, r, discard, errors.New("db password=secret"))
	}), httpx.Correlation())

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(observability.CorrelationHeader, "corr-1")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "secret") {
		t.Fatalf("internal error leaked: %s", rec.Body.String())
	}
	env := decode(t, rec)
	if !env.Error || env.Data.Code != string(apperr.Internal) || env.Data.TraceID != "corr-1" || env.Data.Details == nil {
		t.Fatalf("unexpected envelope: %+v", env)
	}
}

func TestWriteErrorMapsTypedErrors(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	httpx.WriteError(rec, req, discard, apperr.New(apperr.ReferenceMismatch, "Idempotency-Key phải bằng source_ref.",
		apperr.Detail{Field: "source_ref", Code: "mismatch"}))

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", rec.Code)
	}
	env := decode(t, rec)
	if env.Data.Code != string(apperr.ReferenceMismatch) || len(env.Data.Details) != 1 {
		t.Fatalf("unexpected envelope: %+v", env)
	}
}

func TestRecoverReturns500Envelope(t *testing.T) {
	h := httpx.Chain(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("boom")
	}), httpx.Recover(discard))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", rec.Code)
	}
	if env := decode(t, rec); env.Data.Code != string(apperr.Internal) {
		t.Fatalf("code = %s", env.Data.Code)
	}
}

func TestNoStoreAndMaxBody(t *testing.T) {
	h := httpx.Chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.ReadAll(r.Body); err != nil {
			w.WriteHeader(http.StatusRequestEntityTooLarge)
		}
	}), httpx.NoStore(), httpx.MaxBody(4))

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/", strings.NewReader("too long")))
	if rec.Header().Get("Cache-Control") != "no-store" {
		t.Error("Cache-Control: no-store missing")
	}
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status = %d, want 413", rec.Code)
	}
}

func TestReadinessReportsFailingDependencies(t *testing.T) {
	mux := http.NewServeMux()
	httpx.NewHealth(map[string]httpx.Check{
		"postgres": func(context.Context) error { return nil },
		"rabbitmq": func(context.Context) error { return errors.New("dial tcp: refused") },
	}).Register(mux)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "rabbitmq") || strings.Contains(body, "refused") || strings.Contains(body, "postgres") {
		t.Fatalf("unexpected body: %s", body)
	}

	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("healthz status = %d", rec.Code)
	}
}
