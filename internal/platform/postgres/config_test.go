package postgres_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pthawng/notitek/internal/platform/postgres"
)

func TestConfigValidate(t *testing.T) {
	ok := postgres.Config{DSN: "postgres://u:p@localhost:5432/notitek", MaxConns: 10, MinConns: 1}
	if err := ok.Validate(); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}
	keyword := postgres.Config{DSN: "host=localhost user=u password=p dbname=notitek", MaxConns: 2}
	if err := keyword.Validate(); err != nil {
		t.Fatalf("keyword DSN rejected: %v", err)
	}

	bad := map[string]postgres.Config{
		"unparsable dsn": {DSN: "postgres://u:p@localhost:notaport/db", MaxConns: 1},
		"no max conns":   {DSN: ok.DSN, MaxConns: 0},
		"min above max":  {DSN: ok.DSN, MaxConns: 2, MinConns: 3},
		"negative min":   {DSN: ok.DSN, MaxConns: 2, MinConns: -1},
	}
	for name, c := range bad {
		if err := c.Validate(); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
	if err := bad["unparsable dsn"].Validate(); err != nil && strings.Contains(err.Error(), "notaport") {
		t.Errorf("error echoes the DSN: %v", err)
	}
}

func TestConfigRequiresTLS(t *testing.T) {
	modes := map[string]bool{
		"disable":     false,
		"allow":       false,
		"prefer":      false,
		"require":     true,
		"verify-ca":   true,
		"verify-full": true,
	}
	for mode, want := range modes {
		c := postgres.Config{DSN: "postgres://u:p@db:5432/notitek?sslmode=" + mode, MaxConns: 1}
		if got := c.RequiresTLS(); got != want {
			t.Errorf("sslmode=%s: RequiresTLS = %v, want %v", mode, got, want)
		}
	}
}

func TestConfigNeverPrintsThePassword(t *testing.T) {
	for _, dsn := range []string{
		"postgres://u:hunter2@localhost:5432/notitek",
		"host=localhost user=u password=hunter2 dbname=notitek",
	} {
		c := postgres.Config{DSN: dsn, MaxConns: 1}
		outputs := map[string]string{
			"%v": fmt.Sprintf("%v", c), "%+v": fmt.Sprintf("%+v", c), "%#v": fmt.Sprintf("%#v", c), "LogValue": c.LogValue().String(),
		}
		for name, out := range outputs {
			if strings.Contains(out, "hunter2") {
				t.Errorf("%s leaks the password for %q: %s", name, dsn, out)
			}
		}
	}
}
