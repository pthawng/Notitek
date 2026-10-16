package observability_test

import (
	"strings"
	"testing"

	"github.com/pthawng/notitek/internal/platform/observability"
)

func TestRedactDSN(t *testing.T) {
	cases := map[string]string{
		"postgres://u:hunter2@localhost:5432/db?sslmode=require": "postgres://u:xxxxx@localhost:5432/db?sslmode=require",
		"amqps://guest:hunter2@mq:5671/vhost":                    "amqps://guest:xxxxx@mq:5671/vhost",
		"host=h user=u password=hunter2 dbname=d":                "host=h user=u password=xxxxx dbname=d",
		"host=h password='hun ter2' dbname=d":                    "host=h password=xxxxx dbname=d",
		"postgres://u@localhost/db":                              "postgres://u@localhost/db",
		"":                                                       "",
	}
	for in, want := range cases {
		if got := observability.RedactDSN(in); got != want {
			t.Errorf("RedactDSN(%q) = %q, want %q", in, got, want)
		}
	}
	if got := observability.RedactDSN("postgres://u:hunter2@bad host/db"); strings.Contains(got, "hunter2") {
		t.Errorf("unparsable DSN leaks the password: %q", got)
	}
}
