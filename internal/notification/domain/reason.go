package domain

// Reason explains why a delivery stopped, failed or is unknown
// (data model v0.8 §4.5). Values are stored as text; the closed lists below
// are checked here, not by the database.
type Reason string

// Reasons of SUPPRESSED deliveries.
const (
	ReasonContactBlocked      Reason = "CONTACT_BLOCKED"
	ReasonBlockedQuota        Reason = "BLOCKED_QUOTA"
	ReasonRateLimited         Reason = "RATE_LIMITED"
	ReasonKillSwitch          Reason = "KILL_SWITCH"
	ReasonChannelNotAllowed   Reason = "CHANNEL_NOT_ALLOWED"
	ReasonRecipientUnresolved Reason = "RECIPIENT_UNRESOLVED"
	ReasonRecipientLimit      Reason = "RECIPIENT_LIMIT_EXCEEDED"
	ReasonNoValidRecipient    Reason = "NO_VALID_RECIPIENT"
	// ReasonFallbackNotNeeded: a fallback the policy never reached, because
	// an earlier delivery of its group settled without a definite rejection.
	// Not a send check: only EventFallbackNotNeeded sets it.
	ReasonFallbackNotNeeded Reason = "FALLBACK_NOT_NEEDED"
)

// Reasons of FAILED deliveries.
const (
	ReasonProviderRejected   Reason = "PROVIDER_REJECTED"
	ReasonInvalidDestination Reason = "INVALID_DESTINATION"
	ReasonTemplateRejected   Reason = "TEMPLATE_REJECTED"
	ReasonRetryExhausted     Reason = "RETRY_EXHAUSTED"
	ReasonBounceHard         Reason = "BOUNCE_HARD"
	ReasonDeliveryFailed     Reason = "DELIVERY_FAILED"
)

// Reasons of EXPIRED and UNKNOWN deliveries.
const (
	ReasonDeadlinePassed  Reason = "DEADLINE_PASSED"
	ReasonProviderTimeout Reason = "PROVIDER_TIMEOUT"
	ReasonWorkerLost      Reason = "WORKER_LOST"
)

// CANCELLED deliveries carry the reason code of the source's cancel command,
// which is not part of a closed list.

var (
	// suppressReasons are the reasons a delivery can be suppressed with. The
	// recipient reasons apply to messages, not deliveries.
	suppressReasons = set(ReasonContactBlocked, ReasonBlockedQuota, ReasonRateLimited,
		ReasonKillSwitch, ReasonChannelNotAllowed)
	// rejectReasons are what a definite rejection of an attempt maps to.
	rejectReasons = set(ReasonProviderRejected, ReasonInvalidDestination, ReasonTemplateRejected)
	// unknownReasons are why the result of an attempt is not known.
	unknownReasons = set(ReasonProviderTimeout, ReasonWorkerLost)
	// failureEvidenceReasons are what negative evidence (callback or lookup)
	// maps to. RETRY_EXHAUSTED is decided by the core, never reported.
	failureEvidenceReasons = set(ReasonBounceHard, ReasonDeliveryFailed,
		ReasonProviderRejected, ReasonInvalidDestination, ReasonTemplateRejected)
	// postAcceptReasons are the negative evidence that can follow provider
	// acceptance (matrix row 15).
	postAcceptReasons = set(ReasonBounceHard, ReasonDeliveryFailed)
)

func set(rs ...Reason) map[Reason]bool {
	m := make(map[Reason]bool, len(rs))
	for _, r := range rs {
		m[r] = true
	}
	return m
}
