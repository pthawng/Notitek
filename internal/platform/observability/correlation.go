package observability

import "context"

// CorrelationHeader is the platform header that ties logs and calls together
// across SuperPlatform services. It is never a deduplication key.
const CorrelationHeader = "X-Correlation-Id"

type correlationKey struct{}

// WithCorrelationID returns a context carrying id.
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationKey{}, id)
}

// CorrelationID returns the correlation id carried by ctx, or "".
func CorrelationID(ctx context.Context) string {
	id, _ := ctx.Value(correlationKey{}).(string)
	return id
}
