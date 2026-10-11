// Package bootstrap is the composition root: the only place that builds
// infrastructure clients and wires modules together. Each runtime builds only
// the dependencies it uses.
package bootstrap

import (
	"time"

	"github.com/pthawng/notitek/internal/platform/config"
	"github.com/pthawng/notitek/internal/platform/httpx"
	"github.com/pthawng/notitek/internal/platform/postgres"
	"github.com/pthawng/notitek/internal/platform/rabbitmq"
	"github.com/pthawng/notitek/internal/platform/redis"
)

// Runtime names one of the Notitek binaries.
type Runtime string

// Runtimes.
const (
	RuntimeAPI       Runtime = "api"
	RuntimeWorker    Runtime = "worker"
	RuntimeRelay     Runtime = "relay"
	RuntimeScheduler Runtime = "scheduler"
)

// needs lists the infrastructure each runtime connects to (see
// docs/Notitek - CẤU TRÚC MÃ NGUỒN.md §5). The API writes jobs to the outbox
// in Postgres and never talks to RabbitMQ directly.
var needs = map[Runtime]struct{ postgres, rabbitmq, redis bool }{
	RuntimeAPI:       {postgres: true, redis: true},
	RuntimeWorker:    {postgres: true, rabbitmq: true, redis: true},
	RuntimeRelay:     {postgres: true, rabbitmq: true},
	RuntimeScheduler: {postgres: true},
}

// Config is the configuration of one runtime.
type Config struct {
	Runtime  Runtime
	Env      string
	LogLevel string

	// HTTP is the public listener (api only). Admin serves health probes for
	// every runtime.
	HTTP  httpx.ServerConfig
	Admin httpx.ServerConfig

	Postgres postgres.Config
	RabbitMQ rabbitmq.Config
	Redis    redis.Config

	ShutdownTimeout time.Duration
}

// LoadConfig reads the configuration of runtime from the environment and
// returns every problem at once.
func LoadConfig(runtime Runtime, lookup config.LookupFunc) (Config, error) {
	e := config.NewEnv(lookup)
	n := needs[runtime]

	shutdown := e.Duration("NOTITEK_SHUTDOWN_TIMEOUT", 15*time.Second)
	cfg := Config{
		Runtime:         runtime,
		Env:             e.OneOf("NOTITEK_ENV", "local", "local", "dev", "staging", "production"),
		LogLevel:        e.OneOf("NOTITEK_LOG_LEVEL", "info", "debug", "info", "warn", "error"),
		ShutdownTimeout: shutdown,
		Admin: httpx.ServerConfig{
			Addr:            e.String("NOTITEK_ADMIN_ADDR", ":9090"),
			ReadTimeout:     5 * time.Second,
			WriteTimeout:    5 * time.Second,
			IdleTimeout:     30 * time.Second,
			ShutdownTimeout: shutdown,
		},
	}

	if runtime == RuntimeAPI {
		cfg.HTTP = httpx.ServerConfig{
			Addr:        e.String("NOTITEK_HTTP_ADDR", ":8080"),
			ReadTimeout: e.Duration("NOTITEK_HTTP_READ_TIMEOUT", 10*time.Second),
			// 0 = no server-wide write deadline: SSE streams are long-lived and
			// manage their own deadlines.
			WriteTimeout:    0,
			IdleTimeout:     e.Duration("NOTITEK_HTTP_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: shutdown,
		}
	}
	if n.postgres {
		cfg.Postgres = postgres.Config{
			DSN:             e.Required("NOTITEK_POSTGRES_DSN"),
			MaxConns:        int32(e.Int("NOTITEK_POSTGRES_MAX_CONNS", 10, 1, 500)),
			MinConns:        int32(e.Int("NOTITEK_POSTGRES_MIN_CONNS", 1, 0, 100)),
			MaxConnLifetime: e.Duration("NOTITEK_POSTGRES_MAX_CONN_LIFETIME", time.Hour),
			ApplicationName: "notitek-" + string(runtime),
		}
	}
	if n.rabbitmq {
		cfg.RabbitMQ = rabbitmq.Config{
			URL:            e.Required("NOTITEK_RABBITMQ_URL"),
			ConnectionName: "notitek-" + string(runtime),
		}
	}
	if n.redis {
		cfg.Redis = redis.Config{
			Addr:     e.Required("NOTITEK_REDIS_ADDR"),
			Password: e.String("NOTITEK_REDIS_PASSWORD", ""),
			DB:       e.Int("NOTITEK_REDIS_DB", 0, 0, 15),
		}
	}

	return cfg, e.Err()
}
