package domain

import (
	"errors"
	"fmt"
)

// Transition errors. Callers match them with errors.Is.
var (
	// ErrInvalidTransition: the event does not apply to the delivery in its
	// current state, for example a claim of an UNKNOWN delivery or provider
	// evidence for an In-app one. The caller must not change the delivery.
	ErrInvalidTransition = errors.New("domain: transition not allowed")
	// ErrInvalidReason: the event carries a reason that the target state does
	// not accept (data model v0.8 §4.5).
	ErrInvalidReason = errors.New("domain: reason not allowed")
)

// EventKind is something that happens to a delivery.
type EventKind string

// Events, with the rows of the transition matrix (§4.4 of the data model v0.8)
// they drive.
const (
	// EventClaimed: a worker claimed the delivery and every send check passed
	// (row 1). Also the claim of a retry.
	EventClaimed EventKind = "CLAIMED"
	// EventSuppressed: a send check failed: kill-switch, blocked contact,
	// quota, rate limit or channel not allowed (row 2).
	EventSuppressed EventKind = "SUPPRESSED"
	// EventCancelled: a cancel command won the race (row 3).
	EventCancelled EventKind = "CANCELLED"
	// EventDeadlinePassed: the send deadline passed before an attempt (row 4).
	EventDeadlinePassed EventKind = "DEADLINE_PASSED"
	// EventInboxStored: the In-app item was stored (row 5).
	EventInboxStored EventKind = "INBOX_STORED"
	// EventAttemptAccepted: the provider accepted the attempt (row 6).
	EventAttemptAccepted EventKind = "ATTEMPT_ACCEPTED"
	// EventAttemptTransient: the attempt was certainly not sent and the error
	// is temporary (rows 7 and 8).
	EventAttemptTransient EventKind = "ATTEMPT_TRANSIENT_NOT_SENT"
	// EventAttemptRejected: the attempt was certainly not sent and will never
	// be (row 9).
	EventAttemptRejected EventKind = "ATTEMPT_REJECTED_DEFINITE"
	// EventAttemptUnknown: the attempt may have been sent: timeout after the
	// request left, or the worker lost its lease (row 10).
	EventAttemptUnknown EventKind = "ATTEMPT_UNKNOWN"
	// EventEvidenceAccepted: a callback or lookup says the provider accepted
	// the message (rows 11, 17, 19).
	EventEvidenceAccepted EventKind = "EVIDENCE_ACCEPTED"
	// EventEvidenceDelivered: a callback or lookup says the message was
	// delivered (rows 12, 14, 17, 19).
	EventEvidenceDelivered EventKind = "EVIDENCE_DELIVERED"
	// EventEvidenceFailed: a callback or lookup says the message was not
	// delivered: bounce, delivery failure or not sent (rows 13, 15, 16, 19).
	EventEvidenceFailed EventKind = "EVIDENCE_FAILED"
	// EventReconcileDeadline: the 24-hour reconcile window of an UNKNOWN
	// delivery ended (row 18).
	EventReconcileDeadline EventKind = "RECONCILE_DEADLINE"
	// EventFallbackNotNeeded: the fallback never ran and never will, because
	// an earlier delivery of its group settled without a definite rejection
	// (row 20). Sent for each delivery a Change with SkipFallback names.
	EventFallbackNotNeeded EventKind = "FALLBACK_NOT_NEEDED"
)

// Event is one input of Transition.
type Event struct {
	Kind EventKind
	// Reason is required by EventSuppressed, EventCancelled (the source's
	// reason code), EventAttemptRejected, EventAttemptUnknown and
	// EventEvidenceFailed; ignored otherwise.
	Reason Reason
	// RetryAllowed, for EventAttemptTransient only: the retry policy has an
	// attempt left and the next one starts before the deadline.
	RetryAllowed bool
}

// EvidenceResult is how a piece of evidence was applied
// (delivery_evidence.applied_result). UNMATCHED is decided before
// Transition, when no attempt matches.
type EvidenceResult string

// Evidence results.
const (
	// EvidenceApplied: the evidence changed the state.
	EvidenceApplied EvidenceResult = "APPLIED"
	// EvidenceSuperseded: the evidence is not stronger than the current state.
	EvidenceSuperseded EvidenceResult = "SUPERSEDED"
	// EvidenceConflict: the evidence contradicts a state that stays.
	EvidenceConflict EvidenceResult = "CONFLICT"
)

// Change is the result of a transition. The caller persists it with a
// conditional UPDATE on From, so a concurrent transition makes it a no-op.
type Change struct {
	// Row is the row of the transition matrix that applied.
	Row  int
	From DeliveryState
	To   DeliveryState
	// Reason is the reason_code after the change; empty when To needs none.
	// Only meaningful when StateChanged.
	Reason Reason
	// Evidence is set for evidence events only.
	Evidence EvidenceResult

	// StartAttempt: create an attempt STARTED and grant one quota unit in the
	// same transaction, then commit before calling the provider (row 1).
	StartAttempt bool
	// ScheduleRetry: set next_attempt_at for another attempt (row 7).
	ScheduleRetry bool
	// RefundQuota: give back the unit granted at claim; the attempt certainly
	// did not send (rows 7 to 9).
	RefundQuota bool
	// ActivateFallback: start the next fallback delivery of the ALTERNATIVE
	// group that has not run yet, as the pinned policy orders them, for
	// example SMS after ZNS. Never all alternatives at once (row 9 only).
	ActivateFallback bool
	// SkipFallback: in a FALLBACK group, close every fallback delivery that
	// has not run yet with EventFallbackNotNeeded. Set when an ALTERNATIVE
	// delivery stops being active for any reason other than a definite
	// rejection (which activates the fallback instead) or a cancel (which
	// cancels its siblings itself): it was accepted, became UNKNOWN, expired,
	// was suppressed or ran out of retries. Without it the fallback would stay
	// PENDING and the request would never close. Ignore it in PARALLEL groups,
	// where every alternative runs on its own.
	SkipFallback bool
	// StartReconcileWindow: set reconcile_until once, to the end of the
	// attempt plus 24 hours (row 10).
	StartReconcileWindow bool
}

// StateChanged reports whether the delivery state changes.
func (c Change) StateChanged() bool { return c.From != c.To }

// BumpsRevision reports whether the request revision must increase: every
// change the API shows, including a new attempt (attempt_count) or a
// scheduled retry (next_attempt_at, shown as RETRYING) while the state stays.
// The outcome is recomputed when the state changes.
func (c Change) BumpsRevision() bool { return c.StateChanged() || c.StartAttempt || c.ScheduleRetry }

// Transition applies an event to a delivery according to the R1 transition
// matrix (ADR-04; §4.4 of the data model v0.8, kept in "DATABASE v0.8
// (LỊCH SỬ)"), rows 1 to 20. It is pure: the caller has done the checks that
// need I/O (deadline, kill-switch, quota, retry policy) and passes their
// result in the event.
//
// Rows 1 to 4 apply to a delivery that waits for work: PENDING, or
// PROCESSING with no attempt in flight (waiting to retry). A retry is claimed
// and checked like a first attempt.
func Transition(d Delivery, e Event) (Change, error) {
	c := Change{From: d.State, To: d.State}
	inApp := d.Channel == ChannelInApp
	waiting := d.State == StatePending || (d.State == StateProcessing && !d.AttemptInFlight)
	// In-app is stored, never sent to a provider: it has no attempt.
	inFlight := d.State == StateProcessing && d.AttemptInFlight && !inApp

	switch e.Kind {
	case EventClaimed:
		if inApp || !waiting {
			return invalid(d, e)
		}
		c.Row, c.To, c.StartAttempt = 1, StateProcessing, true

	case EventSuppressed:
		if !waiting {
			return invalid(d, e)
		}
		if !suppressReasons[e.Reason] {
			return badReason(d, e)
		}
		c.Row, c.To, c.Reason = 2, StateSuppressed, e.Reason

	case EventCancelled:
		if !waiting {
			return invalid(d, e)
		}
		if e.Reason == "" {
			return badReason(d, e)
		}
		c.Row, c.To, c.Reason = 3, StateCancelled, e.Reason

	case EventDeadlinePassed:
		if !waiting {
			return invalid(d, e)
		}
		c.Row, c.To, c.Reason = 4, StateExpired, ReasonDeadlinePassed

	case EventInboxStored:
		if !inApp || d.State != StatePending {
			return invalid(d, e)
		}
		c.Row, c.To = 5, StateAvailable

	case EventAttemptAccepted:
		if !inFlight {
			return invalid(d, e)
		}
		c.Row, c.To = 6, StateProviderAccepted

	case EventAttemptTransient:
		if !inFlight {
			return invalid(d, e)
		}
		c.RefundQuota = true
		if e.RetryAllowed {
			c.Row, c.To, c.ScheduleRetry = 7, StateProcessing, true
		} else {
			c.Row, c.To, c.Reason = 8, StateFailed, ReasonRetryExhausted
		}

	case EventAttemptRejected:
		if !inFlight {
			return invalid(d, e)
		}
		if !rejectReasons[e.Reason] {
			return badReason(d, e)
		}
		c.Row, c.To, c.Reason = 9, StateFailed, e.Reason
		c.RefundQuota = true
		c.ActivateFallback = d.Requirement == Alternative

	case EventAttemptUnknown:
		if !inFlight {
			return invalid(d, e)
		}
		if !unknownReasons[e.Reason] {
			return badReason(d, e)
		}
		c.Row, c.To, c.Reason = 10, StateUnknown, e.Reason
		c.StartReconcileWindow = true

	case EventEvidenceAccepted, EventEvidenceDelivered, EventEvidenceFailed:
		return applyEvidence(c, d, e)

	case EventReconcileDeadline:
		if d.State != StateUnknown || inApp {
			return invalid(d, e)
		}
		// The delivery stays UNKNOWN; evidence arriving later still applies
		// through rows 11 to 13.
		c.Row = 18

	case EventFallbackNotNeeded:
		// Only a fallback that never ran: PENDING, in an ALTERNATIVE group.
		if d.State != StatePending || d.Requirement != Alternative {
			return invalid(d, e)
		}
		c.Row, c.To, c.Reason = 20, StateSuppressed, ReasonFallbackNotNeeded

	default:
		return invalid(d, e)
	}

	c.SkipFallback = d.Requirement == Alternative && d.State.Active() && !c.To.Active() &&
		c.Row != 3 && c.Row != 9 && c.Row != 20
	return c, nil
}

// applyEvidence covers rows 11 to 17 and 19. Evidence only exists for
// provider channels after an attempt.
func applyEvidence(c Change, d Delivery, e Event) (Change, error) {
	if d.Channel == ChannelInApp {
		return invalid(d, e)
	}
	if e.Kind == EventEvidenceFailed && !failureEvidenceReasons[e.Reason] {
		return badReason(d, e)
	}

	switch d.State {
	case StateUnknown:
		c.Evidence = EvidenceApplied
		switch e.Kind {
		case EventEvidenceAccepted:
			c.Row, c.To = 11, StateProviderAccepted
		case EventEvidenceDelivered:
			c.Row, c.To = 12, StateDelivered
		default:
			// Not a definite rejection before sending: no alternative.
			c.Row, c.To, c.Reason = 13, StateFailed, e.Reason
		}

	case StateProviderAccepted:
		switch e.Kind {
		case EventEvidenceAccepted:
			c.Row, c.Evidence = 17, EvidenceSuperseded
		case EventEvidenceDelivered:
			c.Row, c.To, c.Evidence = 14, StateDelivered, EvidenceApplied
		default:
			// Decision A: a bounce or delivery failure after acceptance is
			// the final word and can lower the outcome.
			if !postAcceptReasons[e.Reason] {
				return badReason(d, e)
			}
			c.Row, c.To, c.Reason, c.Evidence = 15, StateFailed, e.Reason, EvidenceApplied
		}

	case StateDelivered:
		if e.Kind == EventEvidenceFailed {
			c.Row, c.Evidence = 16, EvidenceConflict
		} else {
			c.Row, c.Evidence = 17, EvidenceSuperseded
		}

	case StateFailed, StateCancelled, StateExpired, StateSuppressed:
		// Row 19: the state stays. Negative evidence repeats what is known;
		// positive evidence contradicts it and goes to operations.
		c.Row = 19
		if e.Kind == EventEvidenceFailed {
			c.Evidence = EvidenceSuperseded
		} else {
			c.Evidence = EvidenceConflict
		}

	default:
		// PENDING and PROCESSING have no finished attempt to match.
		return invalid(d, e)
	}
	return c, nil
}

func invalid(d Delivery, e Event) (Change, error) {
	return Change{}, fmt.Errorf("%w: %s on %s %s delivery (attempt in flight: %t)",
		ErrInvalidTransition, e.Kind, d.State, d.Channel, d.AttemptInFlight)
}

func badReason(d Delivery, e Event) (Change, error) {
	return Change{}, fmt.Errorf("%w: %s with reason %q on %s delivery",
		ErrInvalidReason, e.Kind, e.Reason, d.State)
}
