package redis_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pthawng/notitek/internal/platform/redis"
)

func TestConfigValidate(t *testing.T) {
	if err := (redis.Config{Addr: "localhost:6379", DB: 0}).Validate(); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}
	bad := map[string]redis.Config{
		"no port":     {Addr: "localhost"},
		"bad port":    {Addr: "localhost:99999"},
		"negative db": {Addr: "localhost:6379", DB: -1},
	}
	for name, c := range bad {
		if err := c.Validate(); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
}

func TestConfigNeverPrintsThePassword(t *testing.T) {
	c := redis.Config{Addr: "localhost:6379", Password: "hunter2"}
	outputs := map[string]string{
		"%v": fmt.Sprintf("%v", c), "%+v": fmt.Sprintf("%+v", c), "%#v": fmt.Sprintf("%#v", c),
	}
	for name, out := range outputs {
		if strings.Contains(out, "hunter2") {
			t.Errorf("%s leaks the password: %s", name, out)
		}
	}
}
