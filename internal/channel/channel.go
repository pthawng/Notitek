// Package channel defines the boundary between the notification core and
// provider adapters (SMS, Email, ZNS, Push, In-app).
//
// Adapters only talk to providers and classify what happened. They never
// decide retries, fallbacks or delivery state; the notification core does,
// based on the Class returned here (DEC-10).
package channel

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/pthawng/notitek/internal/platform/observability"
)

// Kind identifies a delivery channel.
type Kind string

// Channels of R1.
const (
	SMS   Kind = "SMS"
	Email Kind = "EMAIL"
	InApp Kind = "IN_APP"
	Push  Kind = "PUSH"
	ZNS   Kind = "ZNS"
)

// Class is the provider outcome of one attempt, as seen by the core.
type Class int

// Attempt classes. The zero value is invalid so an unset result is caught.
const (
	_ Class = iota
	// Accepted: the provider took the message; ProviderMessageID is set.
	Accepted
	// RejectedDefinite: certainly not sent (invalid number, blocked, template
	// rejected for this recipient). The only class that allows fallback.
	RejectedDefinite
	// TransientNotSent: certainly not sent, temporary (connection refused,
	// 429, 5xx before acceptance). Retry within the deadline.
	TransientNotSent
	// Unknown: may have been sent (timeout after the request left). Never
	// retried and never followed by fallback.
	Unknown
)

func (c Class) String() string {
	switch c {
	case Accepted:
		return "ACCEPTED"
	case RejectedDefinite:
		return "REJECTED_DEFINITE"
	case TransientNotSent:
		return "TRANSIENT_NOT_SENT"
	case Unknown:
		return "UNKNOWN"
	default:
		return fmt.Sprintf("Class(%d)", int(c))
	}
}

// AllowsFallback reports whether an alternative channel may be tried.
func (c Class) AllowsFallback() bool { return c == RejectedDefinite }

// AllowsRetry reports whether the same channel may be retried.
func (c Class) AllowsRetry() bool { return c == TransientNotSent }

// Message is a rendered message addressed to one destination.
type Message struct {
	DeliveryID string
	AttemptID  string
	Channel    Kind
	// Destination is a phone number, email address, push token or inbox
	// recipient id. It is personal data: log it only through the Mask*
	// helpers of the observability package.
	Destination string
	SenderID    string // brandname, ZNS OA, email domain; resolved from the registry
	TemplateRef string // provider template id when the channel requires one (ZNS)
	Subject     string
	Body        observability.Secret // may contain an OTP or a link
	Params      map[string]string
	Deadline    time.Time
}

// Result is the classified outcome of one attempt.
type Result struct {
	Class             Class
	ProviderMessageID string
	ProviderCode      string // raw provider status/error code, for evidence and reports
	RetryAfter        time.Duration
}

// Validate rejects results an adapter must never return.
func (r Result) Validate() error {
	switch r.Class {
	case Accepted:
		if r.ProviderMessageID == "" {
			return errors.New("channel: accepted result requires a provider message id")
		}
	case RejectedDefinite, TransientNotSent, Unknown:
	default:
		return fmt.Errorf("channel: invalid result class %s", r.Class)
	}
	return nil
}

// Sender is implemented by every provider adapter. Send must return within
// msg.Deadline. Transport errors are classified into the Result, not
// returned: a non-nil error means the adapter itself is misconfigured.
type Sender interface {
	Kind() Kind
	Send(ctx context.Context, msg Message) (Result, error)
}
