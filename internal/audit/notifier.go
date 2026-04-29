package audit

import (
	"fmt"
	"log/slog"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/snapshot"
)

// Notifier implements alert.Notifier and writes each alert event to the
// audit log via a Writer.
type Notifier struct {
	w   *Writer
	log *slog.Logger
}

// NewNotifier creates an audit Notifier from the given config.
// Returns nil, nil when auditing is disabled.
func NewNotifier(cfg config.AuditConfig, log *slog.Logger) (*Notifier, error) {
	if !cfg.Enabled {
		return nil, nil
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("audit notifier: %w", err)
	}
	w, err := NewFileWriter(cfg.Path, cfg.Format)
	if err != nil {
		return nil, err
	}
	return &Notifier{w: w, log: log}, nil
}

// Notify satisfies alert.Notifier.
func (n *Notifier) Notify(ev alert.Event) error {
	if err := n.w.Write(ev.Listener, ev.Type); err != nil {
		n.log.Error("audit write failed", "err", err)
		return err
	}
	return nil
}

// Ensure Notifier satisfies the interface at compile time.
var _ interface {
	Notify(alert.Event) error
} = (*Notifier)(nil)

// eventTypeLabel returns a human-readable label for a snapshot.EventType.
func eventTypeLabel(et snapshot.EventType) string {
	switch et {
	case snapshot.EventAppeared:
		return "appeared"
	case snapshot.EventDisappeared:
		return "disappeared"
	default:
		return "unknown"
	}
}
