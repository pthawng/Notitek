package postgres

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/pthawng/notitek/internal/platform/observability"
)

// Validate checks the settings without connecting. The DSN is parsed by pgx
// but never echoed, because it may contain a password.
func (c Config) Validate() error {
	if _, err := pgxpool.ParseConfig(c.DSN); err != nil {
		return errors.New("invalid DSN (expected postgres://user:password@host:port/database or key=value form)")
	}
	if c.MaxConns < 1 {
		return fmt.Errorf("max connections must be at least 1, got %d", c.MaxConns)
	}
	if c.MinConns < 0 || c.MinConns > c.MaxConns {
		return fmt.Errorf("min connections must be between 0 and max connections (%d), got %d", c.MaxConns, c.MinConns)
	}
	return nil
}

// RequiresTLS reports whether the DSN forces an encrypted connection
// (sslmode require, verify-ca or verify-full). The modes disable, allow and
// prefer can fall back to plaintext, so they do not count.
func (c Config) RequiresTLS() bool {
	pc, err := pgxpool.ParseConfig(c.DSN)
	return err == nil && pc.ConnConfig.TLSConfig != nil && len(pc.ConnConfig.Fallbacks) == 0
}

// LogValue implements slog.LogValuer without the DSN password.
func (c Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("dsn", observability.RedactDSN(c.DSN)),
		slog.Int("max_conns", int(c.MaxConns)),
		slog.Int("min_conns", int(c.MinConns)),
		slog.Duration("max_conn_lifetime", c.MaxConnLifetime),
		slog.String("application_name", c.ApplicationName),
	)
}

// String implements fmt.Stringer so %v and %+v never print the password.
func (c Config) String() string { return c.LogValue().String() }

// GoString implements fmt.GoStringer so %#v never prints the password.
func (c Config) GoString() string { return "postgres.Config{" + c.String() + "}" }
