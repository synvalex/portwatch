package baseline

import (
	"context"
	"log/slog"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// Learner observes listeners for a fixed duration and adds them all to a
// baseline Store, optionally persisting the result to disk.
type Learner struct {
	store    *Store
	duration time.Duration
	autoSave bool
	logger   *slog.Logger
}

// NewLearner creates a Learner that populates store over duration.
func NewLearner(store *Store, duration time.Duration, autoSave bool, logger *slog.Logger) *Learner {
	if logger == nil {
		logger = slog.Default()
	}
	return &Learner{
		store:    store,
		duration: duration,
		autoSave: autoSave,
		logger:   logger,
	}
}

// Learn runs until ctx is cancelled or the learn duration elapses, recording
// every listener returned by scan into the store.
func (l *Learner) Learn(ctx context.Context, scan func(context.Context) ([]ports.Listener, error)) error {
	deadline := time.Now().Add(l.duration)
	l.logger.Info("baseline learning started", "duration", l.duration)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		listeners, err := scan(ctx)
		if err != nil {
			l.logger.Warn("baseline scan error", "err", err)
		} else {
			for _, li := range listeners {
				l.store.Add(li)
			}
		}

		if time.Now().After(deadline) {
			break
		}

		select {
		case <-ctx.Done():
			l.logger.Info("baseline learning cancelled")
			return ctx.Err()
		case <-ticker.C:
		}
	}

	l.logger.Info("baseline learning complete")
	if l.autoSave && l.store.filePath != "" {
		if err := l.store.Save(); err != nil {
			l.logger.Warn("baseline save failed", "err", err)
			return err
		}
		l.logger.Info("baseline saved", "file", l.store.filePath)
	}
	return nil
}
