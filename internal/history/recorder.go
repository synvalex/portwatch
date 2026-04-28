package history

import (
	"log/slog"

	"github.com/user/portwatch/internal/alert"
)

// Recorder is an alert.Notifier that writes events into a history Store.
type Recorder struct {
	store  *Store
	logger *slog.Logger
}

// NewRecorder creates a Recorder backed by the given Store.
func NewRecorder(store *Store, logger *slog.Logger) *Recorder {
	return &Recorder{store: store, logger: logger}
}

// Notify satisfies alert.Notifier by recording the event in the store.
func (r *Recorder) Notify(ev alert.Event) error {
	r.store.Add(ev)
	r.logger.Debug("history: recorded event",
		"type", ev.Type,
		"proto", ev.Listener.Proto,
		"port", ev.Listener.Address.Port,
	)
	return nil
}
