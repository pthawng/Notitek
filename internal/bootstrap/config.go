// Package bootstrap is the composition root: the only place that builds
// infrastructure clients and wires modules together. Each runtime builds only
// the dependencies it uses.
package bootstrap

import (
	"log/slog"
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
			Addr:            e.Addr("NOTITEK_ADMIN_ADDR", ":9090"),
			ReadTimeout:     5 * time.Second,
			WriteTimeout:    5 * time.Second,
			IdleTimeout:     30 * time.Second,
			ShutdownTimeout: shutdown,
		},
	}

	if runtime == RuntimeAPI {
		cfg.HTTP = httpx.ServerConfig{
			Addr:        e.Addr("NOTITEK_HTTP_ADDR", ":8080"),
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
		if cfg.Postgres.DSN != "" {
			if err := cfg.Postgres.Validate(); err != nil {
				e.Fail("NOTITEK_POSTGRES_DSN", "%v", err)
			}
		}
	}
	if n.rabbitmq {
		cfg.RabbitMQ = rabbitmq.Config{
			URL:            e.Required("NOTITEK_RABBITMQ_URL"),
			ConnectionName: "notitek-" + string(runtime),
		}
		if cfg.RabbitMQ.URL != "" {
			if err := cfg.RabbitMQ.Validate(); err != nil {
				e.Fail("NOTITEK_RABBITMQ_URL", "%v", err)
			}
		}
	}
	if n.redis {
		cfg.Redis = redis.Config{
			Addr:     e.Required("NOTITEK_REDIS_ADDR"),
			Password: e.String("NOTITEK_REDIS_PASSWORD", ""),
			DB:       e.Int("NOTITEK_REDIS_DB", 0, 0, 15),
		}
		if cfg.Redis.Addr != "" {
			if err := cfg.Redis.Validate(); err != nil {
				e.Fail("NOTITEK_REDIS_ADDR", "%v", err)
			}
		}
	}

	validateCrossField(e, cfg)
	return cfg, e.Err()
}

// validateCrossField checks rules that involve several settings or the
// environment. Rules that need a deployment decision (for example TLS) apply
// to production only, so local and CI stacks keep working without certificates.
func validateCrossField(e *config.Env, cfg Config) {
	if cfg.Runtime == RuntimeAPI && cfg.HTTP.Addr == cfg.Admin.Addr {
		e.Fail("NOTITEK_ADMIN_ADDR", "must differ from NOTITEK_HTTP_ADDR (both %q)", cfg.Admin.Addr)
	}
	if cfg.Env != "production" {
		return
	}
	n := needs[cfg.Runtime]
	if n.postgres && cfg.Postgres.DSN != "" && cfg.Postgres.Validate() == nil && !cfg.Postgres.RequiresTLS() {
		e.Fail("NOTITEK_POSTGRES_DSN", "production requires TLS: use sslmode=require, verify-ca or verify-full")
	}
	if n.rabbitmq && cfg.RabbitMQ.URL != "" && cfg.RabbitMQ.Validate() == nil && !cfg.RabbitMQ.Secure() {
		e.Fail("NOTITEK_RABBITMQ_URL", "production requires TLS: use the amqps:// scheme")
	}
	if n.redis && cfg.Redis.Password == "" {
		e.Fail("NOTITEK_REDIS_PASSWORD", "is required in production")
	}
}

// LogValue implements slog.LogValuer. Secrets in the nested settings are
// already redacted, so the whole configuration can be logged at startup.
func (c Config) LogValue() slog.Value {
	attrs := []slog.Attr{
		slog.String("runtime", string(c.Runtime)),
		slog.String("env", c.Env),
		slog.String("log_level", c.LogLevel),
		slog.String("admin_addr", c.Admin.Addr),
		slog.Duration("shutdown_timeout", c.ShutdownTimeout),
	}
	n := needs[c.Runtime]
	if c.Runtime == RuntimeAPI {
		attrs = append(attrs, slog.String("http_addr", c.HTTP.Addr))
	}
	if n.postgres {
		attrs = append(attrs, slog.Any("postgres", c.Postgres))
	}
	if n.rabbitmq {
		attrs = append(attrs, slog.Any("rabbitmq", c.RabbitMQ))
	}
	if n.redis {
		attrs = append(attrs, slog.Any("redis", c.Redis))
	}
	return slog.GroupValue(attrs...)
}
