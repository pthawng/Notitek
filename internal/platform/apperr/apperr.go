// Package apperr defines the stable error codes of the Notitek API and the
// error type that carries them from use cases to the HTTP layer.
//
// The catalogue mirrors components.schemas.ApiError.x-error-codes in
// docs/api/openapi.json; a test keeps the two in sync.
//
// The members of an error follow RFC 9457 (Problem Details for HTTP APIs):
// every code has a stable Type URI and Title, and the HTTP status comes from
// the code. The platform envelope {error, message, data} carries them in data.
package apperr

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Code is a stable, machine-readable error code. Consumers branch on Code,
// never on Message.
type Code string

// Error codes of API v1.
const (
	InvalidRequest         Code = "INVALID_REQUEST"
	ReferenceMismatch      Code = "REFERENCE_MISMATCH"
	InvalidCursor          Code = "INVALID_CURSOR"
	Unauthenticated        Code = "UNAUTHENTICATED"
	Forbidden              Code = "FORBIDDEN"
	NotFound               Code = "NOT_FOUND"
	MethodNotAllowed       Code = "METHOD_NOT_ALLOWED"
	IdempotencyConflict    Code = "IDEMPOTENCY_CONFLICT"
	IdempotencyInProgress  Code = "IDEMPOTENCY_IN_PROGRESS"
	ReferenceRetired       Code = "REFERENCE_RETIRED"
	PayloadTooLarge        Code = "PAYLOAD_TOO_LARGE"
	UnsupportedMediaType   Code = "UNSUPPORTED_MEDIA_TYPE"
	UnknownType            Code = "UNKNOWN_TYPE"
	UnsupportedTypeVersion Code = "UNSUPPORTED_TYPE_VERSION"
	InvalidProfile         Code = "INVALID_PROFILE"
	InvalidTime            Code = "INVALID_TIME"
	StaleRequest           Code = "STALE_REQUEST"
	DeadlineExpired        Code = "DEADLINE_EXPIRED"
	UnsupportedPurpose     Code = "UNSUPPORTED_PURPOSE"
	RateLimited            Code = "RATE_LIMITED"
	ServiceUnavailable     Code = "SERVICE_UNAVAILABLE"

	// Internal is never listed in the public catalogue; it marks unexpected
	// failures whose details must not reach the caller.
	Internal Code = "INTERNAL_ERROR"
)

type codeInfo struct {
	status    int
	title     string
	retryable bool
}

// info is the single source of truth per code. Titles are short, stable and
// do not change between occurrences (RFC 9457 §3.1.3); the per-occurrence
// explanation is the message.
var info = map[Code]codeInfo{
	InvalidRequest:         {http.StatusBadRequest, "Invalid request", false},
	ReferenceMismatch:      {http.StatusBadRequest, "Reference mismatch", false},
	InvalidCursor:          {http.StatusBadRequest, "Invalid cursor", false},
	Unauthenticated:        {http.StatusUnauthorized, "Authentication required", false},
	Forbidden:              {http.StatusForbidden, "Forbidden", false},
	NotFound:               {http.StatusNotFound, "Resource not found", false},
	MethodNotAllowed:       {http.StatusMethodNotAllowed, "Method not allowed", false},
	IdempotencyConflict:    {http.StatusConflict, "Idempotency conflict", false},
	IdempotencyInProgress:  {http.StatusConflict, "Request in progress", true},
	ReferenceRetired:       {http.StatusConflict, "Reference retired", false},
	PayloadTooLarge:        {http.StatusRequestEntityTooLarge, "Payload too large", false},
	UnsupportedMediaType:   {http.StatusUnsupportedMediaType, "Unsupported media type", false},
	UnknownType:            {http.StatusUnprocessableEntity, "Unknown notification type", false},
	UnsupportedTypeVersion: {http.StatusUnprocessableEntity, "Unsupported type version", false},
	InvalidProfile:         {http.StatusUnprocessableEntity, "Invalid profile data", false},
	InvalidTime:            {http.StatusUnprocessableEntity, "Invalid time", false},
	StaleRequest:           {http.StatusUnprocessableEntity, "Stale request", false},
	DeadlineExpired:        {http.StatusUnprocessableEntity, "Deadline expired", false},
	UnsupportedPurpose:     {http.StatusUnprocessableEntity, "Unsupported purpose", false},
	RateLimited:            {http.StatusTooManyRequests, "Too many requests", true},
	ServiceUnavailable:     {http.StatusServiceUnavailable, "Service unavailable", true},
	Internal:               {http.StatusInternalServerError, "Internal error", false},
}

// TypeBase prefixes the problem type URI of every code. A URN identifies the
// problem without promising a document at a URL; RFC 9457 allows it.
const TypeBase = "urn:supership:notitek:error:"

// Catalogue returns the public codes and their HTTP status.
func Catalogue() map[Code]int {
	out := make(map[Code]int, len(info)-1)
	for c, i := range info {
		if c != Internal {
			out[c] = i.status
		}
	}
	return out
}

// Status returns the HTTP status for c, or 500 for unknown codes.
func (c Code) Status() int {
	if i, ok := info[c]; ok {
		return i.status
	}
	return http.StatusInternalServerError
}

// Title returns the short, stable summary of c.
func (c Code) Title() string {
	if i, ok := info[c]; ok {
		return i.title
	}
	return info[Internal].title
}

// Type returns the problem type URI of c, for example
// urn:supership:notitek:error:idempotency-conflict.
func (c Code) Type() string {
	return TypeBase + strings.ToLower(strings.ReplaceAll(string(c), "_", "-"))
}

// Retryable reports whether the caller may repeat the same request later
// without changing it: the failure is temporary and no work was accepted.
func (c Code) Retryable() bool { return info[c].retryable }

// DefaultRetryAfter is advertised when a retryable error carries no explicit
// delay, so callers never retry in a tight loop.
const DefaultRetryAfter = time.Second

// Location says where in the request a field error was found.
type Location string

// Locations.
const (
	LocationBody   Location = "body"
	LocationHeader Location = "header"
	LocationQuery  Location = "query"
	LocationPath   Location = "path"
)

// Reasons for field errors. They are stable and machine-readable.
const (
	ReasonRequired      = "required"
	ReasonInvalidType   = "invalid_type"
	ReasonInvalidFormat = "invalid_format"
	ReasonInvalidValue  = "invalid_value"
	ReasonOutOfRange    = "out_of_range"
	ReasonTooLong       = "too_long"
	ReasonUnknownField  = "unknown_field"
	ReasonMalformedJSON = "malformed_json"
	ReasonMismatch      = "mismatch"
)

// Detail points at one offending field. Field is a JSON Pointer (RFC 6901,
// "/recipients/0/phone", "/" for the whole body) when Location is body, and
// the header, query parameter or path variable name otherwise. Message is
// safe to show. A Detail never echoes the input value.
type Detail struct {
	Field    string   `json:"field"`
	Location Location `json:"location"`
	Code     string   `json:"code"`
	Message  string   `json:"message"`
}

// FieldError builds a Detail.
func FieldError(location Location, field, code, message string) Detail {
	return Detail{Field: field, Location: location, Code: code, Message: message}
}

// Error is an API-facing error. Message is safe to show; Cause is for logs only.
type Error struct {
	Code    Code
	Message string
	Details []Detail
	// RetryAfter is how long the caller should wait before repeating a
	// retryable request. Zero means DefaultRetryAfter for retryable codes.
	RetryAfter time.Duration
	Cause      error
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error { return e.Cause }

// New returns an error with code and a caller-safe message.
func New(code Code, message string, details ...Detail) *Error {
	return &Error{Code: code, Message: message, Details: details}
}

// Wrap attaches code and message to cause.
func Wrap(cause error, code Code, message string) *Error {
	return &Error{Code: code, Message: message, Cause: cause}
}

// WithRetryAfter sets how long the caller should wait before retrying.
func (e *Error) WithRetryAfter(d time.Duration) *Error {
	e.RetryAfter = d
	return e
}

// From extracts an *Error from err. Anything else becomes Internal, so
// unexpected failures never expose their text to callers.
func From(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return &Error{Code: Internal, Message: "Hệ thống gặp lỗi, vui lòng thử lại sau.", Cause: err}
}
