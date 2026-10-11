package config_test

import (
	"strings"
	"testing"
	"time"

	"github.com/pthawng/notitek/internal/platform/config"
)

func lookup(values map[string]string) config.LookupFunc {
	return func(key string) (string, bool) {
		v, ok := values[key]
		return v, ok
	}
}

func TestEnvReadsTypedValues(t *testing.T) {
	env := config.NewEnv(lookup(map[string]string{
		"NAME":    "  notitek ",
		"SIZE":    "8",
		"TIMEOUT": "15s",
		"MODE":    "dev",
	}))

	if got := env.String("NAME", "x"); got != "notitek" {
		t.Errorf("String = %q, want trimmed value", got)
	}
	if got := env.String("MISSING", "fallback"); got != "fallback" {
		t.Errorf("String default = %q", got)
	}
	if got := env.Int("SIZE", 1, 1, 10); got != 8 {
		t.Errorf("Int = %d", got)
	}
	if got := env.Duration("TIMEOUT", time.Second); got != 15*time.Second {
		t.Errorf("Duration = %s", got)
	}
	if got := env.OneOf("MODE", "local", "local", "dev"); got != "dev" {
		t.Errorf("OneOf = %q", got)
	}
	if err := env.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEnvReportsAllProblemsAtOnce(t *testing.T) {
	env := config.NewEnv(lookup(map[string]string{
		"BLANK":   "   ",
		"SIZE":    "abc",
		"RANGE":   "99",
		"TIMEOUT": "-1s",
		"MODE":    "prod-eu",
	}))

	env.Required("BLANK")
	env.Int("SIZE", 1, 1, 10)
	env.Int("RANGE", 1, 1, 10)
	env.Duration("TIMEOUT", time.Second)
	env.OneOf("MODE", "local", "local", "prod")

	err := env.Err()
	if err == nil {
		t.Fatal("expected aggregated error")
	}
	for _, key := range []string{"BLANK", "SIZE", "RANGE", "TIMEOUT", "MODE"} {
		if !strings.Contains(err.Error(), key) {
			t.Errorf("error does not mention %s: %v", key, err)
		}
	}
}
