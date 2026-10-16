package httpx_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pthawng/notitek/internal/platform/apperr"
	"github.com/pthawng/notitek/internal/platform/httpx"
)

type recipient struct {
	Phone string `json:"phone"`
}

type request struct {
	SourceRef  string      `json:"source_ref"`
	Count      int         `json:"count"`
	Recipients []recipient `json:"recipients"`
}

func decodeBody(contentType, body string, limit int64) (request, error) {
	req := httptest.NewRequest(http.MethodPost, "/v1/notifications", strings.NewReader(body))
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if limit > 0 {
		req.Body = http.MaxBytesReader(httptest.NewRecorder(), req.Body, limit)
	}
	var out request
	err := httpx.DecodeJSON(req, &out)
	return out, err
}

func TestDecodeJSONAcceptsValidBody(t *testing.T) {
	for _, ct := range []string{"application/json", "application/json; charset=utf-8", "Application/JSON"} {
		got, err := decodeBody(ct, `{"source_ref":"r-1","count":2,"recipients":[{"phone":"x"}]}`+"\n  ", 0)
		if err != nil {
			t.Fatalf("%s: %v", ct, err)
		}
		if got.SourceRef != "r-1" || got.Count != 2 || len(got.Recipients) != 1 {
			t.Errorf("%s: decoded %+v", ct, got)
		}
	}
}

func TestDecodeJSONErrors(t *testing.T) {
	cases := map[string]struct {
		contentType, body string
		limit             int64
		code              apperr.Code
		field, location   string
		reason            string
	}{
		"missing content type": {"", `{}`, 0, apperr.UnsupportedMediaType, "Content-Type", "header", apperr.ReasonInvalidValue},
		"wrong content type":   {"text/plain", `{}`, 0, apperr.UnsupportedMediaType, "Content-Type", "header", apperr.ReasonInvalidValue},
		"empty body":           {"application/json", ``, 0, apperr.InvalidRequest, "/", "body", apperr.ReasonRequired},
		"truncated":            {"application/json", `{"source_ref":`, 0, apperr.InvalidRequest, "/", "body", apperr.ReasonMalformedJSON},
		"syntax error":         {"application/json", `{source_ref: 1}`, 0, apperr.InvalidRequest, "/", "body", apperr.ReasonMalformedJSON},
		"trailing value":       {"application/json", `{} {}`, 0, apperr.InvalidRequest, "/", "body", apperr.ReasonMalformedJSON},
		"unknown field":        {"application/json", `{"source_ref":"a","channels":["SMS"]}`, 0, apperr.InvalidRequest, "/channels", "body", apperr.ReasonUnknownField},
		"wrong type":           {"application/json", `{"count":"two"}`, 0, apperr.InvalidRequest, "/count", "body", apperr.ReasonInvalidType},
		"wrong nested type":    {"application/json", `{"recipients":[{"phone":123}]}`, 0, apperr.InvalidRequest, "/recipients/phone", "body", apperr.ReasonInvalidType},
		"too large":            {"application/json", `{"source_ref":"` + strings.Repeat("a", 100) + `"}`, 16, apperr.PayloadTooLarge, "/", "body", apperr.ReasonTooLong},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := decodeBody(c.contentType, c.body, c.limit)
			var e *apperr.Error
			if !errors.As(err, &e) {
				t.Fatalf("want *apperr.Error, got %v", err)
			}
			if e.Code != c.code {
				t.Errorf("code = %s, want %s", e.Code, c.code)
			}
			if len(e.Details) != 1 {
				t.Fatalf("details = %+v", e.Details)
			}
			d := e.Details[0]
			if d.Field != c.field || string(d.Location) != c.location || d.Code != c.reason || d.Message == "" {
				t.Errorf("detail = %+v, want field=%s location=%s code=%s", d, c.field, c.location, c.reason)
			}
		})
	}
}

// TestDecodeJSONNeverEchoesInput: a body can hold an OTP, a link or a phone
// number; neither the message nor any detail may repeat it.
func TestDecodeJSONNeverEchoesInput(t *testing.T) {
	const secret = "PII-0900000123"
	for _, body := range []string{
		`{"count":"` + secret + `"}`,
		`{"` + "x" + `":1,"source_ref":"` + secret + `",}`,
		`{"recipients":[{"phone":` + secret + `}]}`,
		`{"source_ref":"` + secret,
	} {
		_, err := decodeBody("application/json", body, 0)
		var e *apperr.Error
		if !errors.As(err, &e) {
			t.Fatalf("want *apperr.Error for %q, got %v", body, err)
		}
		all := e.Message
		for _, d := range e.Details {
			all += d.Field + d.Code + d.Message
		}
		if strings.Contains(all, secret) {
			t.Errorf("error echoes the input for %q: %s", body, all)
		}
	}
}
