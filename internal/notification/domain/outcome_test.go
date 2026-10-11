package domain_test

import (
	"testing"

	"github.com/pthawng/notitek/internal/notification/domain"
)

func d(group string, req domain.Requirement, ch domain.Channel, st domain.DeliveryState) domain.Delivery {
	return domain.Delivery{RequirementGroup: group, Requirement: req, Channel: ch, State: st}
}

const (
	req = domain.Required
	alt = domain.Alternative
	opt = domain.Optional
	sms = domain.ChannelSMS
	app = domain.ChannelInApp
	pus = domain.ChannelPush
	zns = domain.ChannelZNS
)

func TestAggregateOutcome(t *testing.T) {
	cases := []struct {
		name string
		in   []domain.Delivery
		want domain.Outcome
	}{
		{"no recipients", nil, domain.OutcomeSuppressed},
		{"only optional", []domain.Delivery{d("g", opt, pus, domain.StateFailed)}, domain.OutcomeSuppressed},
		{
			"otp sms accepted",
			[]domain.Delivery{d("otp", req, sms, domain.StateProviderAccepted)},
			domain.OutcomeAchieved,
		},
		{
			"in-app reaches threshold at AVAILABLE only",
			[]domain.Delivery{d("shop", req, app, domain.StateProviderAccepted)},
			domain.OutcomeFailed,
		},
		{
			"optional push failure does not block achieved",
			[]domain.Delivery{
				d("shop", req, app, domain.StateAvailable),
				d("shop-push", opt, pus, domain.StateFailed),
			},
			domain.OutcomeAchieved,
		},
		{
			"alternative group reached by sms fallback",
			[]domain.Delivery{
				d("consignee", alt, zns, domain.StateFailed),
				d("consignee", alt, sms, domain.StateDelivered),
			},
			domain.OutcomeAchieved,
		},
		{
			"some groups reached",
			[]domain.Delivery{
				d("shop", req, app, domain.StateAvailable),
				d("consignee", alt, zns, domain.StateFailed),
				d("consignee", alt, sms, domain.StateFailed),
			},
			domain.OutcomePartial,
		},
		{
			"required group needs every delivery",
			[]domain.Delivery{
				d("shop", req, app, domain.StateAvailable),
				d("shop", req, app, domain.StatePending),
			},
			domain.OutcomePending,
		},
		{
			"nothing reached and work continues",
			[]domain.Delivery{d("consignee", alt, zns, domain.StateProcessing)},
			domain.OutcomePending,
		},
		{
			"unknown wins over failed when work is done",
			[]domain.Delivery{
				d("a", req, sms, domain.StateUnknown),
				d("b", req, sms, domain.StateFailed),
			},
			domain.OutcomeUnknown,
		},
		{
			"all cancelled",
			[]domain.Delivery{
				d("a", req, app, domain.StateCancelled),
				d("b", alt, zns, domain.StateCancelled),
			},
			domain.OutcomeCancelled,
		},
		{"all expired", []domain.Delivery{d("a", req, zns, domain.StateExpired)}, domain.OutcomeExpired},
		{"all suppressed", []domain.Delivery{d("a", req, sms, domain.StateSuppressed)}, domain.OutcomeSuppressed},
		{
			"mixed terminal states",
			[]domain.Delivery{
				d("a", req, sms, domain.StateCancelled),
				d("b", req, sms, domain.StateExpired),
			},
			domain.OutcomeFailed,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := domain.AggregateOutcome(c.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != c.want {
				t.Errorf("got %s, want %s", got, c.want)
			}
		})
	}
}

func TestAggregateOutcomeRejectsMixedGroup(t *testing.T) {
	_, err := domain.AggregateOutcome([]domain.Delivery{
		d("g", req, sms, domain.StateDelivered),
		d("g", alt, zns, domain.StateFailed),
	})
	if err == nil {
		t.Fatal("expected error for a group mixing REQUIRED and ALTERNATIVE")
	}
}
