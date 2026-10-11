package bootstrap

import (
	"context"
	"log/slog"
)

// The background runtimes currently only connect their dependencies and
// expose probes. Their loops are added with the modules that own them:
//
//   - worker:    notification dispatch consumers, one pool per channel
//   - relay:     outbox claim (FOR UPDATE SKIP LOCKED) → publish with confirms
//   - scheduler: deadline/quiet-hour/UNKNOWN/retention jobs behind advisory locks

func assembleWorker(_ Config, _ *Dependencies, logger *slog.Logger) ([]Component, error) {
	return []Component{idle("worker", logger)}, nil
}

func assembleRelay(_ Config, _ *Dependencies, logger *slog.Logger) ([]Component, error) {
	return []Component{idle("relay", logger)}, nil
}

func assembleScheduler(_ Config, _ *Dependencies, logger *slog.Logger) ([]Component, error) {
	return []Component{idle("scheduler", logger)}, nil
}

// idle blocks until shutdown so a runtime without loops still serves probes.
func idle(name string, logger *slog.Logger) Component {
	return ComponentFunc(func(ctx context.Context) error {
		logger.Info("no processing loop registered yet", "component", name)
		<-ctx.Done()
		return nil
	})
}
