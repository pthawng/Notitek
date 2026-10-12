package bootstrap

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"
)

func env(values map[string]string) func(string) (string, bool) {
	return func(k string) (string, bool) {
		v, ok := values[k]
		return v, ok
	}
}

func TestLoadConfigRequiresOnlyWhatTheRuntimeUses(t *testing.T) {
	cases := []struct {
		runtime Runtime
		missing []string
		absent  []string
	}{
		{RuntimeAPI, []string{"NOTITEK_POSTGRES_DSN", "NOTITEK_REDIS_ADDR"}, []string{"NOTITEK_RABBITMQ_URL"}},
		{RuntimeWorker, []string{"NOTITEK_POSTGRES_DSN", "NOTITEK_RABBITMQ_URL", "NOTITEK_REDIS_ADDR"}, nil},
		{RuntimeRelay, []string{"NOTITEK_POSTGRES_DSN", "NOTITEK_RABBITMQ_URL"}, []string{"NOTITEK_REDIS_ADDR"}},
		{RuntimeScheduler, []string{"NOTITEK_POSTGRES_DSN"}, []string{"NOTITEK_RABBITMQ_URL", "NOTITEK_REDIS_ADDR"}},
	}
	for _, c := range cases {
		t.Run(string(c.runtime), func(t *testing.T) {
			_, err := LoadConfig(c.runtime, env(nil))
			if err == nil {
				t.Fatal("expected missing configuration error")
			}
			for _, k := range c.missing {
				if !strings.Contains(err.Error(), k) {
					t.Errorf("error should mention %s: %v", k, err)
				}
			}
			for _, k := range c.absent {
				if strings.Contains(err.Error(), k) {
					t.Errorf("%s runtime must not require %s", c.runtime, k)
				}
			}
		})
	}
}

func TestLoadConfigAPI(t *testing.T) {
	cfg, err := LoadConfig(RuntimeAPI, env(map[string]string{
		"NOTITEK_ENV":              "dev",
		"NOTITEK_POSTGRES_DSN":     "postgres://u:p@localhost:5432/notitek",
		"NOTITEK_REDIS_ADDR":       "localhost:6379",
		"NOTITEK_SHUTDOWN_TIMEOUT": "20s",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.HTTP.Addr != ":8080" || cfg.Admin.Addr != ":9090" {
		t.Errorf("addresses = %q / %q", cfg.HTTP.Addr, cfg.Admin.Addr)
	}
	if cfg.HTTP.WriteTimeout != 0 {
		t.Error("public server must not set a write timeout (SSE)")
	}
	if cfg.HTTP.ShutdownTimeout != 20*time.Second || cfg.Admin.ShutdownTimeout != 20*time.Second {
		t.Error("shutdown timeout not propagated")
	}
	if cfg.Postgres.ApplicationName != "notitek-api" {
		t.Errorf("application name = %q", cfg.Postgres.ApplicationName)
	}
}

// valid returns a complete configuration for every runtime; tests override
// single keys.
func valid(overrides map[string]string) func(string) (string, bool) {
	values := map[string]string{
		"NOTITEK_POSTGRES_DSN": "postgres://u:s3cr3t-pg@localhost:5432/notitek?sslmode=disable",
		"NOTITEK_RABBITMQ_URL": "amqp://u:s3cr3t-mq@localhost:5672/notitek",
		"NOTITEK_REDIS_ADDR":   "localhost:6379",
	}
	for k, v := range overrides {
		values[k] = v
	}
	return env(values)
}

func TestLoadConfigRejectsMalformedValues(t *testing.T) {
	cases := map[string]struct {
		overrides map[string]string
		wantKey   string
	}{
		"dsn":              {map[string]string{"NOTITEK_POSTGRES_DSN": "postgres://u:p@localhost:notaport/db"}, "NOTITEK_POSTGRES_DSN"},
		"rabbitmq url":     {map[string]string{"NOTITEK_RABBITMQ_URL": "http://localhost"}, "NOTITEK_RABBITMQ_URL"},
		"redis addr":       {map[string]string{"NOTITEK_REDIS_ADDR": "localhost"}, "NOTITEK_REDIS_ADDR"},
		"http addr":        {map[string]string{"NOTITEK_HTTP_ADDR": "8080"}, "NOTITEK_HTTP_ADDR"},
		"admin port range": {map[string]string{"NOTITEK_ADMIN_ADDR": ":70000"}, "NOTITEK_ADMIN_ADDR"},
		"min above max":    {map[string]string{"NOTITEK_POSTGRES_MAX_CONNS": "5", "NOTITEK_POSTGRES_MIN_CONNS": "9"}, "min connections"},
		"same listen addr": {map[string]string{"NOTITEK_HTTP_ADDR": ":9090"}, "NOTITEK_ADMIN_ADDR"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			// The worker never uses the public listener, so test the API for
			// listener rules and the worker for the rest through one runtime.
			_, err := LoadConfig(RuntimeAPI, valid(c.overrides))
			if name == "rabbitmq url" {
				_, err = LoadConfig(RuntimeWorker, valid(c.overrides))
			}
			if err == nil {
				t.Fatal("expected a validation error")
			}
			if !strings.Contains(err.Error(), c.wantKey) {
				t.Errorf("error should mention %q: %v", c.wantKey, err)
			}
		})
	}
}

func TestLoadConfigNeverEchoesSecretsInErrors(t *testing.T) {
	_, err := LoadConfig(RuntimeWorker, valid(map[string]string{
		"NOTITEK_POSTGRES_DSN": "postgres://u:s3cr3t-pg@localhost:notaport/notitek",
		"NOTITEK_RABBITMQ_URL": "http://u:s3cr3t-mq@localhost",
	}))
	if err == nil {
		t.Fatal("expected a validation error")
	}
	for _, secret := range []string{"s3cr3t-pg", "s3cr3t-mq"} {
		if strings.Contains(err.Error(), secret) {
			t.Errorf("error leaks %q: %v", secret, err)
		}
	}
}

func TestLoadConfigProductionRules(t *testing.T) {
	secure := map[string]string{
		"NOTITEK_ENV":            "production",
		"NOTITEK_POSTGRES_DSN":   "postgres://u:p@db:5432/notitek?sslmode=verify-full",
		"NOTITEK_RABBITMQ_URL":   "amqps://u:p@mq:5671/notitek",
		"NOTITEK_REDIS_PASSWORD": "redis-secret",
	}
	if _, err := LoadConfig(RuntimeWorker, valid(secure)); err != nil {
		t.Fatalf("secure production configuration rejected: %v", err)
	}

	cases := map[string]struct{ key, value string }{
		"postgres without TLS":      {"NOTITEK_POSTGRES_DSN", "postgres://u:p@db:5432/notitek?sslmode=disable"},
		"postgres prefer is no TLS": {"NOTITEK_POSTGRES_DSN", "postgres://u:p@db:5432/notitek?sslmode=prefer"},
		"rabbitmq without TLS":      {"NOTITEK_RABBITMQ_URL", "amqp://u:p@mq:5672/notitek"},
		"redis without password":    {"NOTITEK_REDIS_PASSWORD", ""},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			o := map[string]string{}
			for k, v := range secure {
				o[k] = v
			}
			o[c.key] = c.value
			_, err := LoadConfig(RuntimeWorker, valid(o))
			if err == nil || !strings.Contains(err.Error(), c.key) {
				t.Fatalf("want an error for %s, got %v", c.key, err)
			}
		})
	}

	// The same insecure values are fine outside production.
	if _, err := LoadConfig(RuntimeWorker, valid(map[string]string{"NOTITEK_ENV": "staging"})); err != nil {
		t.Errorf("staging should not require TLS: %v", err)
	}
	// A runtime that does not use a dependency is not held to its rules.
	if _, err := LoadConfig(RuntimeScheduler, valid(map[string]string{
		"NOTITEK_ENV":          "production",
		"NOTITEK_POSTGRES_DSN": "postgres://u:p@db:5432/notitek?sslmode=require",
	})); err != nil {
		t.Errorf("scheduler uses only Postgres: %v", err)
	}
}

func TestConfigNeverPrintsSecrets(t *testing.T) {
	cfg, err := LoadConfig(RuntimeWorker, valid(map[string]string{"NOTITEK_REDIS_PASSWORD": "redis-secret"}))
	if err != nil {
		t.Fatal(err)
	}
	var logged bytes.Buffer
	slog.New(slog.NewJSONHandler(&logged, nil)).Info("configuration loaded", "config", cfg)

	outputs := map[string]string{
		"%v":   fmt.Sprintf("%v", cfg),
		"%+v":  fmt.Sprintf("%+v", cfg),
		"%#v":  fmt.Sprintf("%#v", cfg),
		"slog": logged.String(),
	}
	for name, out := range outputs {
		for _, secret := range []string{"s3cr3t-pg", "s3cr3t-mq", "redis-secret"} {
			if strings.Contains(out, secret) {
				t.Errorf("%s leaks %q: %s", name, secret, out)
			}
		}
	}
	if !strings.Contains(logged.String(), "localhost:5432") {
		t.Errorf("startup log should still show non-secret settings: %s", logged.String())
	}
}

// TestEnvExampleLoadsForEveryRuntime keeps .env.example in step with the
// code: a copied example must start every runtime.
func TestEnvExampleLoadsForEveryRuntime(t *testing.T) {
	raw, err := os.ReadFile("../../.env.example")
	if err != nil {
		t.Fatal(err)
	}
	values := map[string]string{}
	for _, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			t.Fatalf("malformed line in .env.example: %q", line)
		}
		values[k] = v
	}
	for _, rt := range []Runtime{RuntimeAPI, RuntimeWorker, RuntimeRelay, RuntimeScheduler} {
		if _, err := LoadConfig(rt, env(values)); err != nil {
			t.Errorf("%s does not start with .env.example: %v", rt, err)
		}
	}
}
