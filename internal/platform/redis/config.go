package redis

import (
	"fmt"
	"log/slog"

	"github.com/pthawng/notitek/internal/platform/config"
	"github.com/pthawng/notitek/internal/platform/observability"
)

// Validate checks the settings without connecting.
func (c Config) Validate() error {
	if err := config.ValidateAddr(c.Addr); err != nil {
		return fmt.Errorf("address %w", err)
	}
	if c.DB < 0 {
		return fmt.Errorf("db must not be negative, got %d", c.DB)
	}
	return nil
}

// LogValue implements slog.LogValuer without the password.
func (c Config) LogValue() slog.Value {
	password := ""
	if c.Password != "" {
		password = observability.NewSecret(c.Password).String()
	}
	return slog.GroupValue(
		slog.String("addr", c.Addr),
		slog.String("password", password),
		slog.Int("db", c.DB),
	)
}

// String implements fmt.Stringer so %v and %+v never print the password.
func (c Config) String() string { return c.LogValue().String() }

// GoString implements fmt.GoStringer so %#v never prints the password.
func (c Config) GoString() string { return "redis.Config{" + c.String() + "}" }
