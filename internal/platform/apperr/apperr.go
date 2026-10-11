// Package apperr defines the stable error codes of the Notitek API and the
// error type that carries them from use cases to the HTTP layer.
//
// The catalogue mirrors components.schemas.ApiError.x-error-codes in
// docs/api/openapi.json; a test keeps the two in sync.
package apperr

import (
	"errors"
	"fmt"
	"net/http"
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

var statusByCode = map[Code]int{
	InvalidRequest:         http.StatusBadRequest,
	ReferenceMismatch:      http.StatusBadRequest,
	InvalidCursor:          http.StatusBadRequest,
	Unauthenticated:        http.StatusUnauthorized,
	Forbidden:              http.StatusForbidden,
	NotFound:               http.StatusNotFound,
	IdempotencyConflict:    http.StatusConflict,
	IdempotencyInProgress:  http.StatusConflict,
	ReferenceRetired:       http.StatusConflict,
	PayloadTooLarge:        http.StatusRequestEntityTooLarge,
	UnsupportedMediaType:   http.StatusUnsupportedMediaType,
	UnknownType:            http.StatusUnprocessableEntity,
	UnsupportedTypeVersion: http.StatusUnprocessableEntity,
	InvalidProfile:         http.StatusUnprocessableEntity,
	InvalidTime:            http.StatusUnprocessableEntity,
	StaleRequest:           http.StatusUnprocessableEntity,
	DeadlineExpired:        http.StatusUnprocessableEntity,
	UnsupportedPurpose:     http.StatusUnprocessableEntity,
	RateLimited:            http.StatusTooManyRequests,
	ServiceUnavailable:     http.StatusServiceUnavailable,
	Internal:               http.StatusInternalServerError,
}

// Catalogue returns the public codes and their HTTP status.
func Catalogue() map[Code]int {
	out := make(map[Code]int, len(statusByCode)-1)
	for c, s := range statusByCode {
		if c != Internal {
			out[c] = s
		}
	}
	return out
}

// Status returns the HTTP status for c, or 500 for unknown codes.
func (c Code) Status() int {
	if s, ok := statusByCode[c]; ok {
		return s
	}
	return http.StatusInternalServerError
}

// Detail points at an offending field. It never echoes the input value.
type Detail struct {
	Field string `json:"field"`
	Code  string `json:"code"`
}

// Error is an API-facing error. Message is safe to show; Cause is for logs only.
type Error struct {
	Code    Code
	Message string
	Details []Detail
	Cause   error
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

// From extracts an *Error from err. Anything else becomes Internal, so
// unexpected failures never expose their text to callers.
func From(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return &Error{Code: Internal, Message: "Hệ thống gặp lỗi, vui lòng thử lại sau.", Cause: err}
}
