package monitor

import (
	"context"
	"log/slog"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/ports"
)

// Monitor periodically scans open ports and dispatches alerts via the dispatcher.
type Monitor struct {
	scanner    ports.Scanner
	dispatcher *alert.Dispatcher
	interval   time.Duration
	logger     *slog.Logger
}

// New creates a new Monitor.
func New(scanner ports.Scanner, dispatcher *alert.Dispatcher, interval time.Duration, logger *slog.Logger) *Monitor {
	if logger == nil {
		logger = slog.Default()
	}
	return &Monitor{
		scanner:    scanner,
		dispatcher: dispatcher,
		interval:   interval,
		logger:     logger,
	}
}

// Run starts the monitoring loop, blocking until ctx is cancelled.
func (m *Monitor) Run(ctx context.Context) error {
	m.logger.Info("monitor started", "interval", m.interval)
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	// Run an immediate scan before waiting for the first tick.
	if err := m.scan(ctx); err != nil {
		m.logger.Error("initial scan failed", "err", err)
	}

	for {
		select {
		case <-ctx.Done():
			m.logger.Info("monitor stopped")
			return ctx.Err()
		case <-ticker.C:
			if err := m.scan(ctx); err != nil {
				m.logger.Error("scan failed", "err", err)
			}
		}
	}
}

// scan performs a single port scan and dispatches results.
func (m *Monitor) scan(ctx context.Context) error {
	listeners, err := m.scanner.Scan(ctx)
	if err != nil {
		return err
	}
	m.logger.Debug("scan complete", "listeners", len(listeners))
	for _, l := range listeners {
		m.dispatcher.Dispatch(l)
	}
	return nil
}
