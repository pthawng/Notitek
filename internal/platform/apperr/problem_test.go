package apperr_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/pthawng/notitek/internal/platform/apperr"
)

// TestEveryCodeHasProblemMembers keeps the catalogue complete: each code,
// including the internal one, has a status, a title and a unique type URI.
func TestEveryCodeHasProblemMembers(t *testing.T) {
	codes := []apperr.Code{apperr.Internal}
	for c := range apperr.Catalogue() {
		codes = append(codes, c)
	}

	types := map[string]apperr.Code{}
	for _, c := range codes {
		if s := c.Status(); s < 400 || s > 599 {
			t.Errorf("%s: status %d is not an error status", c, s)
		}
		if c.Title() == "" {
			t.Errorf("%s: empty title", c)
		}
		typ := c.Type()
		if !strings.HasPrefix(typ, apperr.TypeBase) || typ != strings.ToLower(typ) || strings.Contains(typ, "_") {
			t.Errorf("%s: type %q is not urn:...:kebab-case", c, typ)
		}
		if other, dup := types[typ]; dup {
			t.Errorf("%s and %s share type %q", c, other, typ)
		}
		types[typ] = c
	}
}

func TestOnlyTemporaryFailuresAreRetryable(t *testing.T) {
	retryable := map[apperr.Code]bool{
		apperr.RateLimited:           true,
		apperr.ServiceUnavailable:    true,
		apperr.IdempotencyInProgress: true,
	}
	for c, status := range apperr.Catalogue() {
		if got := c.Retryable(); got != retryable[c] {
			t.Errorf("%s (HTTP %d): Retryable = %v, want %v", c, status, got, retryable[c])
		}
	}
	if apperr.Internal.Retryable() {
		t.Error("INTERNAL_ERROR must not be marked retryable: look the request up by reference first")
	}
}

func TestUnknownCodeFallsBackToInternal(t *testing.T) {
	c := apperr.Code("SOMETHING_NEW")
	if c.Status() != http.StatusInternalServerError || c.Title() != apperr.Internal.Title() {
		t.Errorf("unknown code: status %d, title %q", c.Status(), c.Title())
	}
}

func TestFieldErrorBuildsAFullDetail(t *testing.T) {
	d := apperr.FieldError(apperr.LocationQuery, "page_size", apperr.ReasonOutOfRange, "Tối đa 50.")
	if d.Field != "page_size" || d.Location != apperr.LocationQuery || d.Code != apperr.ReasonOutOfRange || d.Message != "Tối đa 50." {
		t.Errorf("detail = %+v", d)
	}
}
