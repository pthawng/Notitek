package bootstrap

import (
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
