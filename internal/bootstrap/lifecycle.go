package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/pthawng/notitek/internal/platform/buildinfo"
	"github.com/pthawng/notitek/internal/platform/httpx"
	"github.com/pthawng/notitek/internal/platform/observability"
)

// Component is a long-running part of a runtime (HTTP server, consumer loop,
// relay loop). Run blocks until ctx is cancelled, then stops gracefully.
type Component interface {
	Run(ctx context.Context) error
}

// ComponentFunc adapts a function to Component.
type ComponentFunc func(ctx context.Context) error

// Run implements Component.
func (f ComponentFunc) Run(ctx context.Context) error { return f(ctx) }

// assemble builds the components of a runtime from its dependencies.
type assemble func(cfg Config, deps *Dependencies, logger *slog.Logger) ([]Component, error)

// Main loads configuration, connects dependencies, runs the runtime until
// SIGINT/SIGTERM and returns the process exit code.
func Main(runtime Runtime) int {
	cfg, err := LoadConfig(runtime, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "notitek %s: invalid configuration:\n%v\n", runtime, err)
		return 2
	}
	logger := observability.NewLogger(os.Stdout, observability.LoggerConfig{
		Level: cfg.LogLevel, Runtime: string(runtime), Env: cfg.Env,
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, cfg, logger, assemblers[runtime]); err != nil {
		logger.Error("runtime stopped with error", "err", err)
		return 1
	}
	logger.Info("runtime stopped")
	return 0
}

var assemblers = map[Runtime]assemble{
	RuntimeAPI:       assembleAPI,
	RuntimeWorker:    assembleWorker,
	RuntimeRelay:     assembleRelay,
	RuntimeScheduler: assembleScheduler,
}

func run(ctx context.Context, cfg Config, logger *slog.Logger, build assemble) (err error) {
	bi := buildinfo.Read()
	logger.Info("runtime starting", "revision", bi.Revision, "modified", bi.Modified, "go", bi.GoVersion)

	deps, err := openDependencies(ctx, cfg, logger)
	if err != nil {
		return fmt.Errorf("open dependencies: %w", err)
	}
	defer func() { err = errors.Join(err, deps.Close()) }()

	components, err := build(cfg, deps, logger)
	if err != nil {
		return fmt.Errorf("assemble %s: %w", cfg.Runtime, err)
	}
	components = append(components, adminServer(cfg, deps, logger))
	if deps.RabbitMQ != nil {
		components = append(components, ComponentFunc(deps.RabbitMQ.Watch))
	}

	g, gctx := errgroup.WithContext(ctx)
	for _, c := range components {
		g.Go(func() error { return c.Run(gctx) })
	}
	return g.Wait()
}

// adminServer exposes liveness and readiness probes on the admin port.
func adminServer(cfg Config, deps *Dependencies, logger *slog.Logger) Component {
	mux := http.NewServeMux()
	httpx.NewHealth(deps.Checks()).Register(mux)
	return httpx.NewServer("admin", cfg.Admin, mux, logger)
}
