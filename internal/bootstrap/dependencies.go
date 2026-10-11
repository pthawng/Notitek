package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	goredis "github.com/redis/go-redis/v9"

	"github.com/pthawng/notitek/internal/platform/httpx"
	"github.com/pthawng/notitek/internal/platform/postgres"
	"github.com/pthawng/notitek/internal/platform/rabbitmq"
	"github.com/pthawng/notitek/internal/platform/redis"
)

// Dependencies are the infrastructure clients of one runtime. Fields a
// runtime does not need stay nil.
type Dependencies struct {
	Postgres *pgxpool.Pool
	Tx       *postgres.TxManager
	RabbitMQ *rabbitmq.Conn
	Redis    *goredis.Client

	closers []func() error
}

// openDependencies connects to what cfg.Runtime needs. On failure it closes
// whatever was already opened.
func openDependencies(ctx context.Context, cfg Config, logger *slog.Logger) (deps *Dependencies, err error) {
	deps = &Dependencies{}
	defer func() {
		if err != nil {
			err = errors.Join(err, deps.Close())
			deps = nil
		}
	}()

	n := needs[cfg.Runtime]
	if n.postgres {
		pool, err := postgres.Open(ctx, cfg.Postgres)
		if err != nil {
			return deps, err
		}
		deps.Postgres = pool
		deps.Tx = postgres.NewTxManager(pool)
		deps.closers = append(deps.closers, func() error { pool.Close(); return nil })
		logger.Info("connected", "dependency", "postgres")
	}
	if n.redis {
		rc, err := redis.Open(ctx, cfg.Redis)
		if err != nil {
			return deps, err
		}
		deps.Redis = rc
		deps.closers = append(deps.closers, rc.Close)
		logger.Info("connected", "dependency", "redis")
	}
	if n.rabbitmq {
		conn, err := rabbitmq.Dial(cfg.RabbitMQ)
		if err != nil {
			return deps, err
		}
		deps.RabbitMQ = conn
		deps.closers = append(deps.closers, conn.Close)
		logger.Info("connected", "dependency", "rabbitmq")
	}
	return deps, nil
}

// Checks returns readiness checks for the connected dependencies.
func (d *Dependencies) Checks() map[string]httpx.Check {
	checks := map[string]httpx.Check{}
	if d.Postgres != nil {
		checks["postgres"] = postgres.Ping(d.Postgres)
	}
	if d.Redis != nil {
		checks["redis"] = redis.Ping(d.Redis)
	}
	if d.RabbitMQ != nil {
		checks["rabbitmq"] = d.RabbitMQ.Check
	}
	return checks
}

// Close releases dependencies in reverse order of opening.
func (d *Dependencies) Close() error {
	var errs []error
	for i := len(d.closers) - 1; i >= 0; i-- {
		if err := d.closers[i](); err != nil {
			errs = append(errs, fmt.Errorf("close dependency: %w", err))
		}
	}
	d.closers = nil
	return errors.Join(errs...)
}
