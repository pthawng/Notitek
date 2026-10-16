package rabbitmq

import (
	"errors"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/pthawng/notitek/internal/platform/observability"
)

// Validate checks the URL without connecting. The URL is never echoed because
// it may carry credentials.
func (c Config) Validate() error {
	if _, err := amqp.ParseURI(c.URL); err != nil {
		return errors.New("invalid URL (expected amqp://user:password@host:port/vhost or amqps://...)")
	}
	return nil
}

// Secure reports whether the URL uses TLS (amqps).
func (c Config) Secure() bool {
	u, err := amqp.ParseURI(c.URL)
	return err == nil && u.Scheme == "amqps"
}

// LogValue implements slog.LogValuer without the URL password.
func (c Config) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("url", observability.RedactDSN(c.URL)),
		slog.String("connection_name", c.ConnectionName),
	)
}

// String implements fmt.Stringer so %v and %+v never print the password.
func (c Config) String() string { return c.LogValue().String() }

// GoString implements fmt.GoStringer so %#v never prints the password.
func (c Config) GoString() string { return "rabbitmq.Config{" + c.String() + "}" }
