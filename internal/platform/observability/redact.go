package observability

import (
	"log/slog"
	"net/url"
	"regexp"
	"strings"
)

const redacted = "[REDACTED]"

// Secret holds a value that must never reach logs, traces or error messages,
// such as an OTP or an invitation link. Use Reveal only at the point where the
// value is handed to a provider.
type Secret struct {
	value string
}

// NewSecret wraps v.
func NewSecret(v string) Secret { return Secret{value: v} }

// Reveal returns the underlying value.
func (s Secret) Reveal() string { return s.value }

// IsZero reports whether no value is held.
func (s Secret) IsZero() bool { return s.value == "" }

// String implements fmt.Stringer so %v and %s never print the value.
func (Secret) String() string { return redacted }

// GoString implements fmt.GoStringer so %#v never prints the value.
func (Secret) GoString() string { return redacted }

// LogValue implements slog.LogValuer.
func (Secret) LogValue() slog.Value { return slog.StringValue(redacted) }

// MarshalText keeps the value out of JSON and other text encodings.
func (Secret) MarshalText() ([]byte, error) { return []byte(redacted), nil }

// MaskPhone keeps the last three digits of an E.164 number, per DEC-09.
func MaskPhone(phone string) string {
	return keepLast(phone, 3)
}

// MaskEmail keeps the first character of the local part and the domain.
func MaskEmail(email string) string {
	at := strings.LastIndexByte(email, '@')
	if at <= 0 {
		return keepLast(email, 0)
	}
	return email[:1] + strings.Repeat("*", at-1) + email[at:]
}

func keepLast(s string, n int) string {
	if len(s) <= n {
		return strings.Repeat("*", len(s))
	}
	return strings.Repeat("*", len(s)-n) + s[len(s)-n:]
}

var dsnPassword = regexp.MustCompile(`(?i)(password\s*=\s*)('[^']*'|\S+)`)

// RedactDSN hides the password of a connection string, in URL form
// (scheme://user:password@host) or keyword form (password=...). A value that
// cannot be parsed is replaced entirely, so it is always safe to log.
func RedactDSN(dsn string) string {
	if dsn == "" {
		return ""
	}
	if strings.Contains(dsn, "://") {
		u, err := url.Parse(dsn)
		if err != nil {
			return redacted
		}
		if u.User != nil {
			if _, ok := u.User.Password(); ok {
				u.User = url.UserPassword(u.User.Username(), "xxxxx")
			}
		}
		return u.String()
	}
	return dsnPassword.ReplaceAllString(dsn, "${1}xxxxx")
}
