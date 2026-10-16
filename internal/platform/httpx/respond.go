package httpx

import (
	"encoding/json"
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"github.com/pthawng/notitek/internal/platform/apperr"
	"github.com/pthawng/notitek/internal/platform/observability"
)

// envelope is the platform response shape shared with User and Order.
type envelope struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// errorData is the data of an error envelope. type, title, status and
// instance are the members of RFC 9457 Problem Details; the envelope message
// plays the role of detail. code, trace_id, retryable and details are
// extension members.
type errorData struct {
	Type              string          `json:"type"`
	Title             string          `json:"title"`
	Status            int             `json:"status"`
	Code              apperr.Code     `json:"code"`
	Instance          string          `json:"instance"`
	TraceID           string          `json:"trace_id"`
	Retryable         bool            `json:"retryable"`
	RetryAfterSeconds *int            `json:"retry_after_seconds,omitempty"`
	Details           []apperr.Detail `json:"details"`
}

// maxInstance bounds the request path echoed in an error.
const maxInstance = 200

// WriteJSON writes a successful envelope with status and data.
func WriteJSON(w http.ResponseWriter, status int, message string, data any) {
	writeEnvelope(w, status, envelope{Error: false, Message: message, Data: data})
}

// WriteError writes an error envelope. Unexpected errors are logged with
// their cause and returned to the caller as INTERNAL_ERROR without details.
func WriteError(w http.ResponseWriter, r *http.Request, logger *slog.Logger, err error) {
	e := apperr.From(err)
	if e.Code.Status() >= http.StatusInternalServerError {
		logger.ErrorContext(r.Context(), "request failed", "code", e.Code, "err", err)
	} else {
		logger.DebugContext(r.Context(), "request rejected", "code", e.Code, "err", err)
	}

	writeErrorEnvelope(w, r, e)
}

func writeErrorEnvelope(w http.ResponseWriter, r *http.Request, e *apperr.Error) {
	details := e.Details
	if details == nil {
		details = []apperr.Detail{}
	}
	message := e.Message
	if message == "" {
		message = e.Code.Title()
	}

	data := errorData{
		Type:      e.Code.Type(),
		Title:     e.Code.Title(),
		Status:    e.Code.Status(),
		Code:      e.Code,
		Instance:  instance(r),
		TraceID:   observability.CorrelationID(r.Context()),
		Retryable: e.Code.Retryable(),
		Details:   details,
	}
	if e.Code.Retryable() {
		d := e.RetryAfter
		if d <= 0 {
			d = apperr.DefaultRetryAfter
		}
		secs := max(1, int(math.Ceil(d.Seconds())))
		data.RetryAfterSeconds = &secs
		// Retry-After in delay-seconds form (RFC 9110 §10.2.3).
		w.Header().Set("Retry-After", strconv.Itoa(secs))
	}

	writeEnvelope(w, e.Code.Status(), envelope{Error: true, Message: message, Data: data})
}

// instance identifies this occurrence by the request path. The query string
// and headers are left out because they may hold personal data.
func instance(r *http.Request) string {
	p := r.URL.Path
	if len(p) > maxInstance {
		p = p[:maxInstance]
	}
	return p
}

func writeEnvelope(w http.ResponseWriter, status int, body envelope) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	// The status line is already sent; an encoding failure here can only be
	// a broken connection, which the server logs on its own.
	_ = json.NewEncoder(w).Encode(body)
}
