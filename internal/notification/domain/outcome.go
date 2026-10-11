package domain

import (
	"fmt"
	"sort"
)

// Outcome is the aggregated result of a notification request
// (OpenAPI Notification.outcome).
type Outcome string

// Outcomes.
const (
	OutcomePending    Outcome = "PENDING"
	OutcomeAchieved   Outcome = "ACHIEVED"
	OutcomePartial    Outcome = "PARTIAL"
	OutcomeFailed     Outcome = "FAILED"
	OutcomeUnknown    Outcome = "UNKNOWN"
	OutcomeCancelled  Outcome = "CANCELLED"
	OutcomeExpired    Outcome = "EXPIRED"
	OutcomeSuppressed Outcome = "SUPPRESSED"
)

// AggregateOutcome derives the request outcome from its deliveries, following
// API v1 §8:
//
//   - Optional deliveries never affect the outcome.
//   - A group is reached when all of its Required deliveries, or any one of its
//     Alternative deliveries, reached the evidence threshold.
//   - All groups reached → ACHIEVED; some → PARTIAL.
//   - No group reached and work can continue → PENDING.
//   - No work left and some result is unknown → UNKNOWN.
//   - Otherwise every counted delivery cancelled/expired/suppressed →
//     CANCELLED/EXPIRED/SUPPRESSED; any other mix → FAILED.
//   - No counted delivery at all (no valid recipient) → SUPPRESSED.
//
// It returns an error when a group mixes Required and Alternative deliveries,
// which the planner must never produce.
func AggregateOutcome(deliveries []Delivery) (Outcome, error) {
	groups := make(map[string][]Delivery)
	var counted []Delivery
	for _, d := range deliveries {
		if d.Requirement == Optional {
			continue
		}
		groups[d.RequirementGroup] = append(groups[d.RequirementGroup], d)
		counted = append(counted, d)
	}
	if len(counted) == 0 {
		return OutcomeSuppressed, nil
	}

	reached := 0
	for _, name := range sortedKeys(groups) {
		ok, err := groupReached(name, groups[name])
		if err != nil {
			return "", err
		}
		if ok {
			reached++
		}
	}
	switch {
	case reached == len(groups):
		return OutcomeAchieved, nil
	case reached > 0:
		return OutcomePartial, nil
	}

	states := make(map[DeliveryState]int)
	for _, d := range counted {
		if d.State.Active() {
			return OutcomePending, nil
		}
		states[d.State]++
	}
	if states[StateUnknown] > 0 {
		return OutcomeUnknown, nil
	}
	if len(states) == 1 {
		switch {
		case states[StateCancelled] > 0:
			return OutcomeCancelled, nil
		case states[StateExpired] > 0:
			return OutcomeExpired, nil
		case states[StateSuppressed] > 0:
			return OutcomeSuppressed, nil
		}
	}
	return OutcomeFailed, nil
}

func groupReached(name string, ds []Delivery) (bool, error) {
	req := ds[0].Requirement
	for _, d := range ds[1:] {
		if d.Requirement != req {
			return false, fmt.Errorf("domain: requirement group %q mixes %s and %s", name, req, d.Requirement)
		}
	}
	if req == Alternative {
		for _, d := range ds {
			if d.Achieved() {
				return true, nil
			}
		}
		return false, nil
	}
	for _, d := range ds {
		if !d.Achieved() {
			return false, nil
		}
	}
	return true, nil
}

func sortedKeys(m map[string][]Delivery) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
