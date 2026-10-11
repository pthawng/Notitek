package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/pthawng/notitek/internal/platform/apperr"
	"github.com/pthawng/notitek/internal/platform/observability"
)

// envelope is the platform response shape shared with User and Order.
type envelope struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type errorData struct {
	Code    apperr.Code     `json:"code"`
	TraceID string          `json:"trace_id"`
	Details []apperr.Detail `json:"details"`
}

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
	writeEnvelope(w, e.Code.Status(), envelope{
		Error:   true,
		Message: e.Message,
		Data: errorData{
			Code:    e.Code,
			TraceID: observability.CorrelationID(r.Context()),
			Details: details,
		},
	})
}

func writeEnvelope(w http.ResponseWriter, status int, body envelope) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	// The status line is already sent; an encoding failure here can only be
	// a broken connection, which the server logs on its own.
	_ = json.NewEncoder(w).Encode(body)
}
