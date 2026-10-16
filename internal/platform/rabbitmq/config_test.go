package rabbitmq_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pthawng/notitek/internal/platform/rabbitmq"
)

func TestConfigValidateAndSecure(t *testing.T) {
	cases := map[string]struct {
		url           string
		valid, secure bool
	}{
		"amqp":     {"amqp://u:p@mq:5672/notitek", true, false},
		"amqps":    {"amqps://u:p@mq:5671/notitek", true, true},
		"http":     {"http://mq", false, false},
		"garbage":  {"%%%", false, false},
		"no value": {"", false, false},
	}
	for name, c := range cases {
		cfg := rabbitmq.Config{URL: c.url}
		if got := cfg.Validate() == nil; got != c.valid {
			t.Errorf("%s: valid = %v, want %v", name, got, c.valid)
		}
		if got := cfg.Secure(); got != c.secure {
			t.Errorf("%s: Secure = %v, want %v", name, got, c.secure)
		}
	}
}

func TestConfigNeverPrintsThePassword(t *testing.T) {
	c := rabbitmq.Config{URL: "amqp://u:hunter2@mq:5672/notitek", ConnectionName: "notitek-worker"}
	outputs := map[string]string{
		"%v": fmt.Sprintf("%v", c), "%+v": fmt.Sprintf("%+v", c), "%#v": fmt.Sprintf("%#v", c),
	}
	for name, out := range outputs {
		if strings.Contains(out, "hunter2") {
			t.Errorf("%s leaks the password: %s", name, out)
		}
	}
}
