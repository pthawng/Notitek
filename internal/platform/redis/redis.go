// Package redis provides the Redis client used for realtime signals (SSE fan
// out across instances) and short-lived caches. Redis is never the source of
// truth for notification data.
package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
)

// Config holds client settings.
type Config struct {
	Addr     string
	Password string
	DB       int
}

// Open creates a client and verifies connectivity.
func Open(ctx context.Context, cfg Config) (*goredis.Client, error) {
	c := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := c.Ping(ctx).Err(); err != nil {
		_ = c.Close()
		return nil, fmt.Errorf("redis: ping %s: %w", cfg.Addr, err)
	}
	return c, nil
}

// Ping is a readiness check.
func Ping(c *goredis.Client) func(context.Context) error {
	return func(ctx context.Context) error { return c.Ping(ctx).Err() }
}
