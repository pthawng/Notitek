package domain_test

import (
	"errors"
	"testing"

	"github.com/pthawng/notitek/internal/notification/domain"
)

const (
	eml = domain.ChannelEmail

	pending   = domain.StatePending
	process   = domain.StateProcessing
	available = domain.StateAvailable
	accepted  = domain.StateProviderAccepted
	delivered = domain.StateDelivered
	failed    = domain.StateFailed
	unknown   = domain.StateUnknown
	cancelled = domain.StateCancelled
	expired   = domain.StateExpired
	suppress  = domain.StateSuppressed
)

var allStates = []domain.DeliveryState{
	pending, process, available, accepted, delivered, failed, unknown, cancelled, expired, suppress,
}

// dl builds a delivery; inFlight marks an attempt STARTED.
func dl(ch domain.Channel, r domain.Requirement, st domain.DeliveryState, inFlight bool) domain.Delivery {
	return domain.Delivery{ID: "d1", Channel: ch, Requirement: r, RequirementGroup: "g", State: st, AttemptInFlight: inFlight}
}

func ev(k domain.EventKind, r domain.Reason) domain.Event { return domain.Event{Kind: k, Reason: r} }

// TestTransitionMatrix checks every row of the R1 matrix (logical data model
// §4.4) with its side effects.
func TestTransitionMatrix(t *testing.T) {
	cases := []struct {
		name string
		d    domain.Delivery
		e    domain.Event
		want domain.Change
	}{
		{"1 claim pending", dl(sms, req, pending, false), ev(domain.EventClaimed, ""),
			domain.Change{Row: 1, From: pending, To: process, StartAttempt: true}},
		{"1 claim a retry", dl(zns, alt, process, false), ev(domain.EventClaimed, ""),
			domain.Change{Row: 1, From: process, To: process, StartAttempt: true}},
		{"2 kill-switch on pending", dl(sms, req, pending, false), ev(domain.EventSuppressed, domain.ReasonKillSwitch),
			domain.Change{Row: 2, From: pending, To: suppress, Reason: domain.ReasonKillSwitch}},
		{"2 quota on a waiting retry", dl(zns, alt, process, false), ev(domain.EventSuppressed, domain.ReasonBlockedQuota),
			domain.Change{Row: 2, From: process, To: suppress, Reason: domain.ReasonBlockedQuota}},
		{"2 in-app kill-switch", dl(app, req, pending, false), ev(domain.EventSuppressed, domain.ReasonKillSwitch),
			domain.Change{Row: 2, From: pending, To: suppress, Reason: domain.ReasonKillSwitch}},
		{"3 cancel pending", dl(app, req, pending, false), ev(domain.EventCancelled, "ORDER_DELIVERED"),
			domain.Change{Row: 3, From: pending, To: cancelled, Reason: "ORDER_DELIVERED"}},
		{"3 cancel a waiting retry", dl(sms, req, process, false), ev(domain.EventCancelled, "ORDER_DELIVERED"),
			domain.Change{Row: 3, From: process, To: cancelled, Reason: "ORDER_DELIVERED"}},
		{"4 deadline pending", dl(eml, req, pending, false), ev(domain.EventDeadlinePassed, ""),
			domain.Change{Row: 4, From: pending, To: expired, Reason: domain.ReasonDeadlinePassed}},
		{"4 deadline waiting retry", dl(eml, req, process, false), ev(domain.EventDeadlinePassed, ""),
			domain.Change{Row: 4, From: process, To: expired, Reason: domain.ReasonDeadlinePassed}},
		{"5 inbox stored", dl(app, req, pending, false), ev(domain.EventInboxStored, ""),
			domain.Change{Row: 5, From: pending, To: available}},
		{"6 attempt accepted", dl(pus, opt, process, true), ev(domain.EventAttemptAccepted, ""),
			domain.Change{Row: 6, From: process, To: accepted}},
		{"7 transient, retry left", dl(eml, req, process, true), domain.Event{Kind: domain.EventAttemptTransient, RetryAllowed: true},
			domain.Change{Row: 7, From: process, To: process, ScheduleRetry: true, RefundQuota: true}},
		{"8 transient, no retry left", dl(eml, req, process, true), domain.Event{Kind: domain.EventAttemptTransient},
			domain.Change{Row: 8, From: process, To: failed, Reason: domain.ReasonRetryExhausted, RefundQuota: true}},
		{"9 rejected, alternative group", dl(zns, alt, process, true), ev(domain.EventAttemptRejected, domain.ReasonTemplateRejected),
			domain.Change{Row: 9, From: process, To: failed, Reason: domain.ReasonTemplateRejected, RefundQuota: true, ActivateAlternatives: true}},
		{"9 rejected, required delivery", dl(sms, req, process, true), ev(domain.EventAttemptRejected, domain.ReasonInvalidDestination),
			domain.Change{Row: 9, From: process, To: failed, Reason: domain.ReasonInvalidDestination, RefundQuota: true}},
		{"10 provider timeout", dl(zns, alt, process, true), ev(domain.EventAttemptUnknown, domain.ReasonProviderTimeout),
			domain.Change{Row: 10, From: process, To: unknown, Reason: domain.ReasonProviderTimeout, StartReconcileWindow: true}},
		{"10 worker lost its lease", dl(sms, req, process, true), ev(domain.EventAttemptUnknown, domain.ReasonWorkerLost),
			domain.Change{Row: 10, From: process, To: unknown, Reason: domain.ReasonWorkerLost, StartReconcileWindow: true}},
		{"11 unknown then accepted", dl(zns, alt, unknown, false), ev(domain.EventEvidenceAccepted, ""),
			domain.Change{Row: 11, From: unknown, To: accepted, Evidence: domain.EvidenceApplied}},
		{"12 unknown then delivered", dl(sms, req, unknown, false), ev(domain.EventEvidenceDelivered, ""),
			domain.Change{Row: 12, From: unknown, To: delivered, Evidence: domain.EvidenceApplied}},
		{"13 unknown then lookup says not sent", dl(zns, alt, unknown, false), ev(domain.EventEvidenceFailed, domain.ReasonProviderRejected),
			domain.Change{Row: 13, From: unknown, To: failed, Reason: domain.ReasonProviderRejected, Evidence: domain.EvidenceApplied}},
		{"14 accepted then delivered", dl(sms, req, accepted, false), ev(domain.EventEvidenceDelivered, ""),
			domain.Change{Row: 14, From: accepted, To: delivered, Evidence: domain.EvidenceApplied}},
		{"15 accepted then hard bounce", dl(eml, req, accepted, false), ev(domain.EventEvidenceFailed, domain.ReasonBounceHard),
			domain.Change{Row: 15, From: accepted, To: failed, Reason: domain.ReasonBounceHard, Evidence: domain.EvidenceApplied}},
		{"15 accepted then delivery failed", dl(zns, alt, accepted, false), ev(domain.EventEvidenceFailed, domain.ReasonDeliveryFailed),
			domain.Change{Row: 15, From: accepted, To: failed, Reason: domain.ReasonDeliveryFailed, Evidence: domain.EvidenceApplied}},
		{"16 delivered then bounce", dl(eml, req, delivered, false), ev(domain.EventEvidenceFailed, domain.ReasonBounceHard),
			domain.Change{Row: 16, From: delivered, To: delivered, Evidence: domain.EvidenceConflict}},
		{"17 accepted after delivered", dl(sms, req, delivered, false), ev(domain.EventEvidenceAccepted, ""),
			domain.Change{Row: 17, From: delivered, To: delivered, Evidence: domain.EvidenceSuperseded}},
		{"17 duplicate acceptance", dl(sms, req, accepted, false), ev(domain.EventEvidenceAccepted, ""),
			domain.Change{Row: 17, From: accepted, To: accepted, Evidence: domain.EvidenceSuperseded}},
		{"18 reconcile window ends", dl(zns, alt, unknown, false), ev(domain.EventReconcileDeadline, ""),
			domain.Change{Row: 18, From: unknown, To: unknown}},
		{"19 failed then bounce again", dl(eml, req, failed, false), ev(domain.EventEvidenceFailed, domain.ReasonBounceHard),
			domain.Change{Row: 19, From: failed, To: failed, Evidence: domain.EvidenceSuperseded}},
		{"19 failed then delivered", dl(eml, req, failed, false), ev(domain.EventEvidenceDelivered, ""),
			domain.Change{Row: 19, From: failed, To: failed, Evidence: domain.EvidenceConflict}},
		{"19 cancelled then accepted", dl(sms, req, cancelled, false), ev(domain.EventEvidenceAccepted, ""),
			domain.Change{Row: 19, From: cancelled, To: cancelled, Evidence: domain.EvidenceConflict}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := domain.Transition(c.d, c.e)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != c.want {
				t.Errorf("got  %+v\nwant %+v", got, c.want)
			}
		})
	}
}

// TestBounceAfterAcceptanceLowersTheOutcome follows decision A end to end:
// an invite email is accepted (ACHIEVED), then hard-bounces. The request must
// drop to FAILED with a new revision so User can issue a new link.
func TestBounceAfterAcceptanceLowersTheOutcome(t *testing.T) {
	email := dl(eml, req, accepted, false)
	before, err := domain.AggregateOutcome([]domain.Delivery{email})
	if err != nil || before != domain.OutcomeAchieved {
		t.Fatalf("before bounce: %s, %v", before, err)
	}

	c, err := domain.Transition(email, ev(domain.EventEvidenceFailed, domain.ReasonBounceHard))
	if err != nil {
		t.Fatal(err)
	}
	if !c.BumpsRevision() {
		t.Error("the bounce must bump the revision")
	}
	email.State = c.To
	after, err := domain.AggregateOutcome([]domain.Delivery{email})
	if err != nil || after != domain.OutcomeFailed {
		t.Errorf("after bounce: %s, %v; want FAILED", after, err)
	}
	if c.ActivateAlternatives || c.RefundQuota {
		t.Errorf("a bounce after acceptance neither falls back nor refunds: %+v", c)
	}

	// Once DELIVERED, a later bounce no longer lowers the outcome (row 16).
	c, err = domain.Transition(dl(eml, req, delivered, false), ev(domain.EventEvidenceFailed, domain.ReasonBounceHard))
	if err != nil || c.StateChanged() || c.Evidence != domain.EvidenceConflict {
		t.Errorf("bounce after DELIVERED: %+v, %v", c, err)
	}

	// Only post-acceptance reasons apply after acceptance.
	_, err = domain.Transition(dl(eml, req, accepted, false), ev(domain.EventEvidenceFailed, domain.ReasonTemplateRejected))
	if !errors.Is(err, domain.ErrInvalidReason) {
		t.Errorf("template rejection after acceptance: %v, want ErrInvalidReason", err)
	}
}

// TestUnknownIsNeverResent follows a ZNS timeout (row 10): the message may
// already be at the recipient, so nothing may send it again or send SMS
// instead; only evidence can settle it, even after the window (rows 11-13, 18).
func TestUnknownIsNeverResent(t *testing.T) {
	zalo := dl(zns, alt, process, true)
	c, err := domain.Transition(zalo, ev(domain.EventAttemptUnknown, domain.ReasonProviderTimeout))
	if err != nil {
		t.Fatal(err)
	}
	if c.RefundQuota || c.ScheduleRetry || c.ActivateAlternatives {
		t.Errorf("UNKNOWN must not refund, retry or fall back: %+v", c)
	}
	if !c.StartReconcileWindow {
		t.Error("UNKNOWN must start the reconcile window")
	}

	zalo.State, zalo.AttemptInFlight = c.To, false
	for _, e := range []domain.Event{
		ev(domain.EventClaimed, ""),
		{Kind: domain.EventAttemptTransient, RetryAllowed: true},
		ev(domain.EventAttemptRejected, domain.ReasonProviderRejected),
		ev(domain.EventAttemptAccepted, ""),
		ev(domain.EventDeadlinePassed, ""),
		ev(domain.EventCancelled, "ORDER_DELIVERED"),
		ev(domain.EventSuppressed, domain.ReasonKillSwitch),
	} {
		if _, err := domain.Transition(zalo, e); !errors.Is(err, domain.ErrInvalidTransition) {
			t.Errorf("%s on UNKNOWN: %v, want ErrInvalidTransition", e.Kind, err)
		}
	}

	// The window closes and the delivery stays UNKNOWN (row 18).
	c, err = domain.Transition(zalo, ev(domain.EventReconcileDeadline, ""))
	if err != nil || c.StateChanged() || c.BumpsRevision() {
		t.Errorf("window end: %+v, %v", c, err)
	}

	// Lookup says it was never sent: FAILED, still without fallback (row 13).
	c, err = domain.Transition(zalo, ev(domain.EventEvidenceFailed, domain.ReasonProviderRejected))
	if err != nil || c.To != failed || c.ActivateAlternatives {
		t.Errorf("late negative evidence: %+v, %v", c, err)
	}
}

func TestTransitionChecksReasons(t *testing.T) {
	cases := []struct {
		name string
		d    domain.Delivery
		e    domain.Event
	}{
		{"suppress with a failure reason", dl(sms, req, pending, false), ev(domain.EventSuppressed, domain.ReasonRetryExhausted)},
		{"suppress with a message-level reason", dl(sms, req, pending, false), ev(domain.EventSuppressed, domain.ReasonNoValidRecipient)},
		{"suppress without reason", dl(sms, req, pending, false), ev(domain.EventSuppressed, "")},
		{"cancel without reason", dl(sms, req, pending, false), ev(domain.EventCancelled, "")},
		{"rejection with a bounce reason", dl(sms, req, process, true), ev(domain.EventAttemptRejected, domain.ReasonBounceHard)},
		{"unknown with a failure reason", dl(sms, req, process, true), ev(domain.EventAttemptUnknown, domain.ReasonProviderRejected)},
		{"evidence with retry exhausted", dl(sms, req, unknown, false), ev(domain.EventEvidenceFailed, domain.ReasonRetryExhausted)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := domain.Transition(c.d, c.e); !errors.Is(err, domain.ErrInvalidReason) {
				t.Errorf("got %v, want ErrInvalidReason", err)
			}
		})
	}
}

// TestTransitionInvariants sweeps every state × in-flight flag × channel ×
// event. Whatever the matrix allows must keep the safety rules, and every row
// must be reachable, so no row is dead and no pair falls through silently.
func TestTransitionInvariants(t *testing.T) {
	events := []domain.Event{
		ev(domain.EventClaimed, ""),
		ev(domain.EventSuppressed, domain.ReasonKillSwitch),
		ev(domain.EventCancelled, "SOURCE_CANCELLED"),
		ev(domain.EventDeadlinePassed, ""),
		ev(domain.EventInboxStored, ""),
		ev(domain.EventAttemptAccepted, ""),
		{Kind: domain.EventAttemptTransient, RetryAllowed: true},
		{Kind: domain.EventAttemptTransient},
		ev(domain.EventAttemptRejected, domain.ReasonProviderRejected),
		ev(domain.EventAttemptUnknown, domain.ReasonProviderTimeout),
		ev(domain.EventEvidenceAccepted, ""),
		ev(domain.EventEvidenceDelivered, ""),
		ev(domain.EventEvidenceFailed, domain.ReasonBounceHard),
		ev(domain.EventReconcileDeadline, ""),
		ev("SOMETHING_NEW", ""),
	}
	evidence := map[domain.EventKind]bool{
		domain.EventEvidenceAccepted: true, domain.EventEvidenceDelivered: true, domain.EventEvidenceFailed: true,
	}
	reached := map[int]bool{}

	for _, st := range allStates {
		for _, inFlight := range []bool{false, true} {
			if inFlight && st != process {
				continue // only PROCESSING can have an attempt STARTED
			}
			for _, ch := range []domain.Channel{sms, eml, zns, pus, app} {
				for _, r := range []domain.Requirement{req, alt, opt} {
					for _, e := range events {
						d := dl(ch, r, st, inFlight)
						c, err := domain.Transition(d, e)
						if err != nil {
							if !errors.Is(err, domain.ErrInvalidTransition) {
								t.Errorf("%s on %+v: unexpected error %v", e.Kind, d, err)
							}
							continue
						}
						reached[c.Row] = true
						checkInvariants(t, d, e, c, evidence[e.Kind])
					}
				}
			}
		}
	}
	for row := 1; row <= 19; row++ {
		if !reached[row] {
			t.Errorf("matrix row %d is never reached", row)
		}
	}
}

func checkInvariants(t *testing.T, d domain.Delivery, e domain.Event, c domain.Change, isEvidence bool) {
	t.Helper()
	where := func() string { return string(e.Kind) + " on " + string(d.State) + "/" + string(d.Channel) }

	if c.From != d.State || c.Row < 1 || c.Row > 19 {
		t.Errorf("%s: malformed change %+v", where(), c)
	}
	if !d.State.Active() && c.To.Active() {
		t.Errorf("%s: a finished delivery became active again (%s)", where(), c.To)
	}
	if (d.State == accepted || d.State == delivered) && c.To == unknown {
		t.Errorf("%s: evidence was lost", where())
	}
	if d.State == delivered && c.To != delivered {
		t.Errorf("%s: DELIVERED must stay DELIVERED", where())
	}
	if isEvidence != (c.Evidence != "") {
		t.Errorf("%s: evidence result %q", where(), c.Evidence)
	}
	if d.Channel == app && (c.StartAttempt || isEvidence || c.To == accepted || c.To == unknown) {
		t.Errorf("%s: In-app has no provider attempt", where())
	}
	if d.Channel != app && c.To == available {
		t.Errorf("%s: AVAILABLE is for In-app only", where())
	}
	if c.ActivateAlternatives && (c.Row != 9 || d.Requirement != alt) {
		t.Errorf("%s: only a definite rejection of an ALTERNATIVE delivery falls back", where())
	}
	if c.RefundQuota && c.Row != 7 && c.Row != 8 && c.Row != 9 {
		t.Errorf("%s: quota refunded although the message may have been sent", where())
	}
	if c.To == unknown && (c.RefundQuota || c.ScheduleRetry || c.ActivateAlternatives) {
		t.Errorf("%s: UNKNOWN must not refund, retry or fall back", where())
	}
	if c.StartReconcileWindow != (c.Row == 10) {
		t.Errorf("%s: reconcile window started on row %d", where(), c.Row)
	}
	needsReason := c.StateChanged() && c.To != process && c.To != accepted && c.To != delivered && c.To != available
	if needsReason != (c.Reason != "") {
		t.Errorf("%s: reason %q for %s", where(), c.Reason, c.To)
	}
}
