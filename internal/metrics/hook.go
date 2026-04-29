package metrics

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/ports"
	"github.com/user/portwatch/internal/snapshot"
)

// Hook connects portwatch runtime events to Prometheus counters.
type Hook struct {
	counters *Counters
}

// NewHook returns a Hook that updates the given Counters.
func NewHook(c *Counters) *Hook {
	return &Hook{counters: c}
}

// ObserveScan records the number of listeners found in a scan.
func (h *Hook) ObserveScan(listeners []ports.Listener) {
	h.counters.ScansTotal.Inc()
	h.counters.ListenersFound.Set(float64(len(listeners)))
}

// ObserveDiff increments alert counters for each diff entry.
func (h *Hook) ObserveDiff(diff []snapshot.Diff) {
	for _, d := range diff {
		switch d.Event {
		case snapshot.EventAppeared:
			h.counters.AlertsTotal.With(prometheus.Labels{"event": "appeared"}).Inc()
		case snapshot.EventDisappeared:
			h.counters.AlertsTotal.With(prometheus.Labels{"event": "disappeared"}).Inc()
		}
	}
}

// NotifyHook wraps a Notifier and records each notification as an alert metric.
type NotifyHook struct {
	inner    alert.Notifier
	counters *Counters
}

// NewNotifyHook wraps inner so that every Notify call increments the fired counter.
func NewNotifyHook(inner alert.Notifier, c *Counters) *NotifyHook {
	return &NotifyHook{inner: inner, counters: c}
}

// Notify delegates to the inner notifier and records the event.
func (n *NotifyHook) Notify(ctx context.Context, event alert.Event) error {
	n.counters.AlertsTotal.With(prometheus.Labels{"event": string(event.Type)}).Inc()
	return n.inner.Notify(ctx, event)
}
