package monitor

import (
	"context"
	"log/slog"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/ports"
	"github.com/user/portwatch/internal/snapshot"
)

// Watcher coordinates the port scanning pipeline with snapshot diffing
// and alert dispatching. It applies debounce logic so that transient
// listeners (e.g. short-lived ephemeral ports) do not trigger spurious alerts.
type Watcher struct {
	pipeline   *ports.Pipeline
	store      *snapshot.Store
	dispatcher *alert.Dispatcher
	cfg        config.WatchConfig
	logger     *slog.Logger
}

// NewWatcher constructs a Watcher from the provided dependencies.
func NewWatcher(
	pipeline *ports.Pipeline,
	store *snapshot.Store,
	dispatcher *alert.Dispatcher,
	cfg config.WatchConfig,
	logger *slog.Logger,
) *Watcher {
	if logger == nil {
		logger = slog.Default()
	}
	return &Watcher{
		pipeline:   pipeline,
		store:      store,
		dispatcher: dispatcher,
		cfg:        cfg,
		logger:     logger,
	}
}

// Run starts the watch loop. It scans ports at the configured interval,
// diffs against the previous snapshot, and dispatches alerts for any
// appeared or disappeared listeners. The loop exits when ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	// Perform an immediate scan before waiting for the first tick.
	if err := w.scan(ctx); err != nil {
		w.logger.Warn("initial scan failed", "error", err)
	}

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("watcher stopping", "reason", ctx.Err())
			return ctx.Err()
		case <-ticker.C:
			if err := w.scan(ctx); err != nil {
				w.logger.Warn("scan failed", "error", err)
			}
		}
	}
}

// scan runs one pipeline pass, updates the snapshot store, and dispatches
// alerts for any diff entries produced by the update.
func (w *Watcher) scan(ctx context.Context) error {
	listeners, err := w.pipeline.Run(ctx)
	if err != nil {
		return err
	}

	diff := w.store.Update(listeners)
	if len(diff) == 0 {
		return nil
	}

	w.logger.Debug("snapshot diff", "entries", len(diff))

	for _, entry := range diff {
		w.dispatcher.Dispatch(ctx, entry)
	}

	return nil
}
