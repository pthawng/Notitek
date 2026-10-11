//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/pthawng/notitek/internal/platform/postgres"
)

// Run with: NOTITEK_TEST_POSTGRES_DSN=... go test -tags integration ./internal/platform/postgres/
func openTest(t *testing.T) *postgres.TxManager {
	t.Helper()
	dsn := os.Getenv("NOTITEK_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("NOTITEK_TEST_POSTGRES_DSN not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pool, err := postgres.Open(ctx, postgres.Config{DSN: dsn, MaxConns: 4, MaxConnLifetime: time.Minute, ApplicationName: "notitek-test"})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)

	tm := postgres.NewTxManager(pool)
	if _, err := tm.Conn(ctx).Exec(ctx, `CREATE TABLE IF NOT EXISTS public.tx_probe (v int)`); err != nil {
		t.Fatal(err)
	}
	if _, err := tm.Conn(ctx).Exec(ctx, `TRUNCATE public.tx_probe`); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_, _ = tm.Conn(context.Background()).Exec(context.Background(), `DROP TABLE IF EXISTS public.tx_probe`)
	})
	return tm
}

func TestWithinTx(t *testing.T) {
	tm := openTest(t)
	ctx := context.Background()
	if err := tm.WithinTx(ctx, func(ctx context.Context) error {
		insert := func(ctx context.Context) error {
			_, err := tm.Conn(ctx).Exec(ctx, `INSERT INTO public.tx_probe VALUES (1)`)
			return err
		}

		if err := insert(ctx); err != nil {
			return err
		}
		// Nested call joins the outer transaction.
		if err := tm.WithinTx(ctx, insert); err != nil {
			return err
		}
		var n int
		if err := tm.Conn(ctx).QueryRow(ctx, `SELECT count(*) FROM public.tx_probe`).Scan(&n); err != nil {
			return err
		}
		if n != 2 {
			t.Errorf("rows inside tx = %d, want 2", n)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestWithinTxRollsBackOnError(t *testing.T) {
	tm := openTest(t)
	ctx := context.Background()
	sentinel := errors.New("boom")

	err := tm.WithinTx(ctx, func(ctx context.Context) error {
		if _, err := tm.Conn(ctx).Exec(ctx, `CREATE TABLE IF NOT EXISTS public.tx_probe_rb (v int)`); err != nil {
			return err
		}
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("err = %v, want sentinel", err)
	}
	var exists bool
	if err := tm.Conn(ctx).QueryRow(ctx, `SELECT to_regclass('public.tx_probe_rb') IS NOT NULL`).Scan(&exists); err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Error("table created inside a failed transaction must not exist")
	}
}

func TestWithinTxRollsBackOnPanic(t *testing.T) {
	tm := openTest(t)
	ctx := context.Background()

	func() {
		defer func() {
			if recover() == nil {
				t.Error("panic must propagate")
			}
		}()
		_ = tm.WithinTx(ctx, func(ctx context.Context) error {
			if _, err := tm.Conn(ctx).Exec(ctx, `CREATE TABLE IF NOT EXISTS public.tx_probe_panic (v int)`); err != nil {
				return err
			}
			panic("boom")
		})
	}()

	var exists bool
	if err := tm.Conn(ctx).QueryRow(ctx, `SELECT to_regclass('public.tx_probe_panic') IS NOT NULL`).Scan(&exists); err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Error("table created before a panic must be rolled back")
	}
}
