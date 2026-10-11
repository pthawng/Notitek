// Package config reads typed configuration from environment variables.
//
// Values are collected through an Env reader that records every problem
// instead of stopping at the first one, so a misconfigured deployment reports
// all missing or malformed keys in a single startup error.
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// LookupFunc resolves a key to its raw value. os.LookupEnv satisfies it.
type LookupFunc func(key string) (string, bool)

// Env reads configuration values and accumulates validation errors.
type Env struct {
	lookup LookupFunc
	errs   []error
}

// NewEnv returns a reader backed by lookup, or by the process environment
// when lookup is nil.
func NewEnv(lookup LookupFunc) *Env {
	if lookup == nil {
		lookup = os.LookupEnv
	}
	return &Env{lookup: lookup}
}

// Err returns every problem recorded so far, or nil.
func (e *Env) Err() error {
	return errors.Join(e.errs...)
}

func (e *Env) fail(key, format string, args ...any) {
	e.errs = append(e.errs, fmt.Errorf("%s: %s", key, fmt.Sprintf(format, args...)))
}

func (e *Env) raw(key string) (string, bool) {
	v, ok := e.lookup(key)
	if !ok {
		return "", false
	}
	v = strings.TrimSpace(v)
	return v, v != ""
}

// String returns the value of key, or def when it is unset or blank.
func (e *Env) String(key, def string) string {
	if v, ok := e.raw(key); ok {
		return v
	}
	return def
}

// Required returns the value of key and records an error when it is missing.
func (e *Env) Required(key string) string {
	v, ok := e.raw(key)
	if !ok {
		e.fail(key, "is required")
	}
	return v
}

// OneOf returns the value of key (or def) and records an error when it is not
// one of allowed.
func (e *Env) OneOf(key, def string, allowed ...string) string {
	v := e.String(key, def)
	for _, a := range allowed {
		if v == a {
			return v
		}
	}
	e.fail(key, "must be one of %s, got %q", strings.Join(allowed, ", "), v)
	return def
}

// Int returns the integer value of key, or def. Values outside [lo, hi] are
// recorded as errors.
func (e *Env) Int(key string, def, lo, hi int) int {
	v, ok := e.raw(key)
	if !ok {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		e.fail(key, "must be an integer, got %q", v)
		return def
	}
	if n < lo || n > hi {
		e.fail(key, "must be between %d and %d, got %d", lo, hi, n)
		return def
	}
	return n
}

// Duration returns the duration value of key (Go syntax, e.g. "15s"), or def.
// Non-positive durations are recorded as errors.
func (e *Env) Duration(key string, def time.Duration) time.Duration {
	v, ok := e.raw(key)
	if !ok {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		e.fail(key, "must be a duration such as 15s, got %q", v)
		return def
	}
	if d <= 0 {
		e.fail(key, "must be positive, got %s", d)
		return def
	}
	return d
}
