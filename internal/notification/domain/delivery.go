// Package domain holds the notification lifecycle rules: delivery states,
// evidence thresholds and outcome aggregation. It performs no I/O and
// imports no other Notitek module.
package domain

// Channel identifies a delivery channel. It mirrors channel.Kind; the
// application layer converts between the two so the domain stays isolated.
type Channel string

// Channels of R1.
const (
	ChannelSMS   Channel = "SMS"
	ChannelEmail Channel = "EMAIL"
	ChannelInApp Channel = "IN_APP"
	ChannelPush  Channel = "PUSH"
	ChannelZNS   Channel = "ZNS"
)

// DeliveryState is the state of one recipient × channel delivery
// (OpenAPI Delivery.state).
type DeliveryState string

// Delivery states.
const (
	StatePending          DeliveryState = "PENDING"
	StateProcessing       DeliveryState = "PROCESSING"
	StateAvailable        DeliveryState = "AVAILABLE"
	StateProviderAccepted DeliveryState = "PROVIDER_ACCEPTED"
	StateDelivered        DeliveryState = "DELIVERED"
	StateFailed           DeliveryState = "FAILED"
	StateUnknown          DeliveryState = "UNKNOWN"
	StateCancelled        DeliveryState = "CANCELLED"
	StateExpired          DeliveryState = "EXPIRED"
	StateSuppressed       DeliveryState = "SUPPRESSED"
)

// Active reports whether work on the delivery can still continue.
func (s DeliveryState) Active() bool {
	return s == StatePending || s == StateProcessing
}

// Requirement says how a delivery counts toward the outcome
// (OpenAPI Delivery.requirement).
type Requirement string

// Requirements.
const (
	// Required: every delivery of the group must reach its threshold.
	Required Requirement = "REQUIRED"
	// Alternative: the group is reached when any one delivery reaches it.
	Alternative Requirement = "ALTERNATIVE"
	// Optional: never affects the outcome; failures stay visible.
	Optional Requirement = "OPTIONAL"
)

// Delivery is the part of a delivery the outcome and the transitions depend on.
type Delivery struct {
	ID               string
	Channel          Channel
	Requirement      Requirement
	RequirementGroup string
	State            DeliveryState
	// AttemptInFlight is true while an attempt is STARTED: the request may be
	// at the provider. A PROCESSING delivery without one is waiting to retry.
	AttemptInFlight bool
}

// Achieved reports whether the delivery reached the evidence threshold of its
// channel (DEC-10): provider acceptance for SMS, Email, ZNS and Push;
// availability to the recipient for In-app.
func (d Delivery) Achieved() bool {
	if d.Channel == ChannelInApp {
		return d.State == StateAvailable
	}
	return d.State == StateProviderAccepted || d.State == StateDelivered
}
