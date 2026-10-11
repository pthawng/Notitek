package channel_test

import (
	"testing"

	"github.com/pthawng/notitek/internal/channel"
)

func TestClassPolicy(t *testing.T) {
	cases := []struct {
		class    channel.Class
		retry    bool
		fallback bool
	}{
		{channel.Accepted, false, false},
		{channel.RejectedDefinite, false, true},
		{channel.TransientNotSent, true, false},
		{channel.Unknown, false, false}, // may have been sent: neither retry nor fallback
	}
	for _, c := range cases {
		if got := c.class.AllowsRetry(); got != c.retry {
			t.Errorf("%s AllowsRetry = %v, want %v", c.class, got, c.retry)
		}
		if got := c.class.AllowsFallback(); got != c.fallback {
			t.Errorf("%s AllowsFallback = %v, want %v", c.class, got, c.fallback)
		}
	}
}

func TestResultValidate(t *testing.T) {
	cases := []struct {
		name    string
		result  channel.Result
		wantErr bool
	}{
		{"accepted with id", channel.Result{Class: channel.Accepted, ProviderMessageID: "m-1"}, false},
		{"accepted without id", channel.Result{Class: channel.Accepted}, true},
		{"unknown", channel.Result{Class: channel.Unknown}, false},
		{"unset class", channel.Result{}, true},
	}
	for _, c := range cases {
		if err := c.result.Validate(); (err != nil) != c.wantErr {
			t.Errorf("%s: err = %v, wantErr %v", c.name, err, c.wantErr)
		}
	}
}
