// Package rabbitmq wraps the AMQP connection used for Notitek's internal work
// queues. RabbitMQ is not the cross-module event bus (decision A2).
package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Config holds connection settings.
type Config struct {
	URL            string
	ConnectionName string
}

// Conn is an AMQP connection.
type Conn struct {
	conn *amqp.Connection
}

// Dial opens a connection. The URL may carry credentials, so its parse and
// dial errors are reported without the URL.
func Dial(cfg Config) (*Conn, error) {
	props := amqp.NewConnectionProperties()
	props.SetClientConnectionName(cfg.ConnectionName)
	c, err := amqp.DialConfig(cfg.URL, amqp.Config{
		Heartbeat:  10 * time.Second,
		Properties: props,
	})
	if err != nil {
		var amqpErr *amqp.Error
		if errors.As(err, &amqpErr) {
			return nil, fmt.Errorf("rabbitmq: dial: %s", amqpErr.Reason)
		}
		return nil, errors.New("rabbitmq: dial failed (check RABBITMQ_URL and broker availability)")
	}
	return &Conn{conn: c}, nil
}

// Watch blocks until ctx ends or the connection drops. A dropped connection
// is returned as an error so the runtime exits and the orchestrator restarts
// it; consumers and publishers on a dead connection cannot recover in place.
func (c *Conn) Watch(ctx context.Context) error {
	closed := c.conn.NotifyClose(make(chan *amqp.Error, 1))
	select {
	case <-ctx.Done():
		return nil
	case amqpErr, ok := <-closed:
		if ok && amqpErr != nil {
			return fmt.Errorf("rabbitmq: connection lost: %s", amqpErr.Reason)
		}
		// Closed without an error means Close was called, normally on shutdown.
		select {
		case <-ctx.Done():
			return nil
		default:
			return errors.New("rabbitmq: connection closed")
		}
	}
}

// Close closes the connection.
func (c *Conn) Close() error {
	if err := c.conn.Close(); err != nil && !errors.Is(err, amqp.ErrClosed) {
		return fmt.Errorf("rabbitmq: close: %w", err)
	}
	return nil
}

// Check is a readiness check.
func (c *Conn) Check(context.Context) error {
	if c.conn.IsClosed() {
		return errors.New("rabbitmq: connection closed")
	}
	return nil
}

// Message is what the outbox relay publishes.
type Message struct {
	Exchange    string
	RoutingKey  string
	MessageID   string // outbox job id; consumers deduplicate on it
	ContentType string
	Body        []byte
	Headers     map[string]any // trace context, correlation id
}

// ErrNotConfirmed means the broker nacked or returned the message, so the
// outbox row must stay unpublished and be retried.
var ErrNotConfirmed = errors.New("rabbitmq: message not confirmed by broker")

// Publisher publishes with publisher confirms and the mandatory flag. Publish
// returns nil only after the broker has taken responsibility for the
// message.
type Publisher struct {
	mu      sync.Mutex
	ch      *amqp.Channel
	returns chan amqp.Return
}

// NewPublisher opens a dedicated channel in confirm mode.
func (c *Conn) NewPublisher() (*Publisher, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("rabbitmq: open channel: %w", err)
	}
	if err := ch.Confirm(false); err != nil {
		_ = ch.Close()
		return nil, fmt.Errorf("rabbitmq: enable confirms: %w", err)
	}
	return &Publisher{ch: ch, returns: ch.NotifyReturn(make(chan amqp.Return, 1))}, nil
}

// Publish sends m and waits for the broker confirm. Calls are serialized so a
// return notification can be matched to the message that caused it.
func (p *Publisher) Publish(ctx context.Context, m Message) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	dc, err := p.ch.PublishWithDeferredConfirmWithContext(ctx, m.Exchange, m.RoutingKey, true, false, amqp.Publishing{
		MessageId:    m.MessageID,
		ContentType:  m.ContentType,
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now().UTC(),
		Headers:      amqp.Table(m.Headers),
		Body:         m.Body,
	})
	if err != nil {
		return fmt.Errorf("rabbitmq: publish: %w", err)
	}

	acked, err := dc.WaitContext(ctx)
	if err != nil {
		return fmt.Errorf("rabbitmq: wait confirm: %w", err)
	}
	// An unroutable mandatory message is returned before its confirm arrives.
	select {
	case r := <-p.returns:
		return fmt.Errorf("%w: returned %d %s", ErrNotConfirmed, r.ReplyCode, r.ReplyText)
	default:
	}
	if !acked {
		return ErrNotConfirmed
	}
	return nil
}

// Close closes the publisher channel.
func (p *Publisher) Close() error {
	if err := p.ch.Close(); err != nil && !errors.Is(err, amqp.ErrClosed) {
		return fmt.Errorf("rabbitmq: close channel: %w", err)
	}
	return nil
}
