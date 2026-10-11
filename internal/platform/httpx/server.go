// Package httpx contains the HTTP server, middleware and response helpers
// shared by every Notitek runtime.
package httpx

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"
)

// ServerConfig holds listener settings.
type ServerConfig struct {
	Addr            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// Server is an http.Server that shuts down gracefully when its context ends.
type Server struct {
	name   string
	srv    *http.Server
	cfg    ServerConfig
	logger *slog.Logger
}

// NewServer returns a server named name (used in logs) serving h.
//
// WriteTimeout must be 0 on servers that stream SSE; the streaming handler
// then controls its own deadlines through http.ResponseController.
func NewServer(name string, cfg ServerConfig, h http.Handler, logger *slog.Logger) *Server {
	return &Server{
		name:   name,
		cfg:    cfg,
		logger: logger,
		srv: &http.Server{
			Addr:              cfg.Addr,
			Handler:           h,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       cfg.ReadTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       cfg.IdleTimeout,
			ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelWarn),
		},
	}
}

// Run serves until ctx is cancelled, then stops accepting connections and
// waits up to ShutdownTimeout for in-flight requests.
func (s *Server) Run(ctx context.Context) error {
	var lc net.ListenConfig
	ln, err := lc.Listen(ctx, "tcp", s.srv.Addr)
	if err != nil {
		return fmt.Errorf("%s server: listen %s: %w", s.name, s.srv.Addr, err)
	}
	s.srv.BaseContext = func(net.Listener) context.Context { return context.WithoutCancel(ctx) }

	errc := make(chan error, 1)
	go func() {
		s.logger.Info("http server started", "server", s.name, "addr", ln.Addr().String())
		errc <- s.srv.Serve(ln)
	}()

	select {
	case err := <-errc:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("%s server: %w", s.name, err)
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), s.cfg.ShutdownTimeout)
	defer cancel()
	s.logger.Info("http server stopping", "server", s.name)
	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("%s server: shutdown: %w", s.name, err)
	}
	return nil
}
