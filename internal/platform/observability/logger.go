// Package observability provides structured logging with data redaction and
// request-scoped correlation.
package observability

import (
	"context"
	"io"
	"log/slog"
	"strings"
)

// LoggerConfig selects the log level and the runtime name attached to every
// record.
type LoggerConfig struct {
	Level   string // debug, info, warn, error
	Runtime string // api, worker, relay, scheduler
	Env     string
}

// NewLogger returns a JSON logger that adds the correlation id from the
// context to each record.
func NewLogger(w io.Writer, cfg LoggerConfig) *slog.Logger {
	h := slog.NewJSONHandler(w, &slog.HandlerOptions{Level: parseLevel(cfg.Level)})
	return slog.New(contextHandler{Handler: h}).With(
		slog.String("service", "notitek"),
		slog.String("runtime", cfg.Runtime),
		slog.String("env", cfg.Env),
	)
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// contextHandler copies request-scoped values from the context into records,
// so call sites only need the *Context logging methods.
type contextHandler struct {
	slog.Handler
}

func (h contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if id := CorrelationID(ctx); id != "" {
		r.AddAttrs(slog.String("correlation_id", id))
	}
	return h.Handler.Handle(ctx, r)
}

func (h contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return contextHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h contextHandler) WithGroup(name string) slog.Handler {
	return contextHandler{Handler: h.Handler.WithGroup(name)}
}
