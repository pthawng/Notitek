package httpx

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"time"
)

// Check reports whether a dependency is usable.
type Check func(ctx context.Context) error

// Health serves liveness and readiness probes.
//
// /healthz answers 200 while the process runs. /readyz runs every check in
// parallel with a short timeout and answers 503 listing failing dependency
// names only; error text stays in logs.
type Health struct {
	checks  map[string]Check
	timeout time.Duration
}

// NewHealth returns probes for the named checks.
func NewHealth(checks map[string]Check) *Health {
	return &Health{checks: checks, timeout: 2 * time.Second}
}

// Register mounts the probe endpoints on mux.
func (h *Health) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeProbe(w, http.StatusOK, probeBody{Status: "ok"})
	})
	mux.HandleFunc("GET /readyz", h.ready)
}

type probeBody struct {
	Status string   `json:"status"`
	Failed []string `json:"failed,omitempty"`
}

func (h *Health) ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.timeout)
	defer cancel()

	var (
		mu     sync.Mutex
		failed []string
		wg     sync.WaitGroup
	)
	for name, check := range h.checks {
		wg.Go(func() {
			if err := check(ctx); err != nil {
				mu.Lock()
				failed = append(failed, name)
				mu.Unlock()
			}
		})
	}
	wg.Wait()

	if len(failed) > 0 {
		sort.Strings(failed)
		writeProbe(w, http.StatusServiceUnavailable, probeBody{Status: "unavailable", Failed: failed})
		return
	}
	writeProbe(w, http.StatusOK, probeBody{Status: "ok"})
}

func writeProbe(w http.ResponseWriter, status int, body probeBody) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
