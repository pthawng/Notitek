package observability_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/pthawng/notitek/internal/platform/observability"
)

func TestSecretNeverLeaks(t *testing.T) {
	s := observability.NewSecret("123456")

	for name, out := range map[string]string{
		"%v":     fmt.Sprintf("%v", s),
		"Sprint": fmt.Sprint(s),
		"%+v":    fmt.Sprintf("%+v", struct{ OTP observability.Secret }{s}),
		"%#v":    fmt.Sprintf("%#v", s),
	} {
		if strings.Contains(out, "123456") {
			t.Errorf("%s leaked secret: %s", name, out)
		}
	}

	b, err := json.Marshal(map[string]any{"otp": s})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "123456") {
		t.Errorf("json leaked secret: %s", b)
	}

	var buf bytes.Buffer
	log := observability.NewLogger(&buf, observability.LoggerConfig{Runtime: "test"})
	log.Info("sending", "otp", s)
	if strings.Contains(buf.String(), "123456") {
		t.Errorf("log leaked secret: %s", buf.String())
	}
	if s.Reveal() != "123456" {
		t.Error("Reveal must return the original value")
	}
}

func TestLoggerAddsCorrelationID(t *testing.T) {
	var buf bytes.Buffer
	log := observability.NewLogger(&buf, observability.LoggerConfig{Runtime: "api", Env: "local"})

	ctx := observability.WithCorrelationID(context.Background(), "corr-1")
	log.InfoContext(ctx, "hello")

	var rec map[string]any
	if err := json.Unmarshal(buf.Bytes(), &rec); err != nil {
		t.Fatalf("log is not JSON: %v", err)
	}
	for key, want := range map[string]string{"correlation_id": "corr-1", "runtime": "api", "service": "notitek"} {
		if rec[key] != want {
			t.Errorf("%s = %v, want %s", key, rec[key], want)
		}
	}
}

func TestMasking(t *testing.T) {
	cases := []struct{ name, got, want string }{
		{"phone", observability.MaskPhone("+84900000123"), "*********123"},
		{"short phone", observability.MaskPhone("12"), "**"},
		{"email", observability.MaskEmail("user@example.invalid"), "u***@example.invalid"},
		{"not an email", observability.MaskEmail("nope"), "****"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s: got %q, want %q", c.name, c.got, c.want)
		}
	}
}
