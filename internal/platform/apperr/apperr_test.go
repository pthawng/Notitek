package apperr_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/pthawng/notitek/internal/platform/apperr"
)

// TestCatalogueMatchesOpenAPI keeps the Go catalogue and the published
// contract identical: same codes, same HTTP status.
func TestCatalogueMatchesOpenAPI(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "docs", "api", "openapi.json"))
	if err != nil {
		t.Fatalf("read openapi: %v", err)
	}
	var spec struct {
		Components struct {
			Schemas struct {
				APIError struct {
					Codes map[string]int `json:"x-error-codes"`
				} `json:"ApiError"`
			} `json:"schemas"`
		} `json:"components"`
	}
	if err := json.Unmarshal(raw, &spec); err != nil {
		t.Fatalf("parse openapi: %v", err)
	}
	want := spec.Components.Schemas.APIError.Codes
	if len(want) == 0 {
		t.Fatal("openapi has no x-error-codes")
	}

	got := apperr.Catalogue()
	for code, status := range want {
		if s, ok := got[apperr.Code(code)]; !ok {
			t.Errorf("code %s is in OpenAPI but not in apperr", code)
		} else if s != status {
			t.Errorf("code %s: apperr status %d, OpenAPI %d", code, s, status)
		}
	}
	for code := range got {
		if _, ok := want[string(code)]; !ok {
			t.Errorf("code %s is in apperr but not in OpenAPI", code)
		}
	}
}

func TestFromHidesUnexpectedErrors(t *testing.T) {
	cause := errors.New("pq: relation does not exist")
	e := apperr.From(fmt.Errorf("load: %w", cause))

	if e.Code != apperr.Internal || e.Code.Status() != http.StatusInternalServerError {
		t.Fatalf("got %s/%d, want INTERNAL_ERROR/500", e.Code, e.Code.Status())
	}
	if !errors.Is(e, cause) {
		t.Error("cause must stay reachable for logging")
	}
}

func TestFromKeepsTypedErrors(t *testing.T) {
	orig := apperr.New(apperr.IdempotencyConflict, "Khác nội dung", apperr.Detail{Field: "source_ref", Code: "conflict"})
	e := apperr.From(fmt.Errorf("admit: %w", orig))

	if e != orig {
		t.Fatalf("From returned %v, want the original error", e)
	}
	if e.Code.Status() != http.StatusConflict {
		t.Errorf("status = %d", e.Code.Status())
	}
}
