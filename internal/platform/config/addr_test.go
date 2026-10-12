package config_test

import (
	"strings"
	"testing"

	"github.com/pthawng/notitek/internal/platform/config"
)

func TestEnvAddr(t *testing.T) {
	good := config.NewEnv(lookup(map[string]string{"A": ":8080", "B": "0.0.0.0:9090", "C": "[::1]:80"}))
	for _, key := range []string{"A", "B", "C"} {
		if got := good.Addr(key, ":1"); got == ":1" {
			t.Errorf("Addr(%s) fell back to the default", key)
		}
	}
	if err := good.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bad := config.NewEnv(lookup(map[string]string{"NOPORT": "localhost", "ZERO": ":0", "HUGE": ":65536", "TEXT": ":http"}))
	for _, key := range []string{"NOPORT", "ZERO", "HUGE", "TEXT"} {
		if got := bad.Addr(key, ":8080"); got != ":8080" {
			t.Errorf("Addr(%s) = %q, want the default after an error", key, got)
		}
	}
	err := bad.Err()
	if err == nil {
		t.Fatal("expected errors")
	}
	for _, key := range []string{"NOPORT", "ZERO", "HUGE", "TEXT"} {
		if !strings.Contains(err.Error(), key) {
			t.Errorf("error does not mention %s: %v", key, err)
		}
	}
}

func TestEnvFailJoinsCallerErrors(t *testing.T) {
	env := config.NewEnv(lookup(nil))
	env.Required("MISSING")
	env.Fail("CUSTOM", "rule %d violated", 7)
	err := env.Err()
	if err == nil || !strings.Contains(err.Error(), "MISSING") || !strings.Contains(err.Error(), "CUSTOM: rule 7 violated") {
		t.Fatalf("unexpected error: %v", err)
	}
}
