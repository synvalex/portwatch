package alert

import (
	"context"
	"log/slog"

	"github.com/jwhittle933/portwatch/internal/ports"
)

// Pipeline chains together a sequence of alert filters and a final notifier.
// Events flow through each filter in order; if any filter suppresses an event,
// the notifier is not called. This provides a composable way to apply dedup,
// cooldown, throttle, rate-limit, and suppress policies before dispatching.
type Pipeline struct {
	filters  []Filter
	notifier Notifier
	logger   *slog.Logger
}

// Filter is a gate that decides whether an alert event should continue
// downstream. Implementations return true to allow the event, false to suppress.
type Filter interface {
	Allow(event Event) bool
}

// PipelineOption configures a Pipeline.
type PipelineOption func(*Pipeline)

// WithFilter appends a filter to the pipeline.
func WithFilter(f Filter) PipelineOption {
	return func(p *Pipeline) {
		if f != nil {
			p.filters = append(p.filters, f)
		}
	}
}

// WithLogger sets the logger used for suppression debug messages.
func WithLogger(l *slog.Logger) PipelineOption {
	return func(p *Pipeline) {
		if l != nil {
			p.logger = l
		}
	}
}

// NewPipeline constructs a Pipeline that routes events through filters before
// forwarding to the given notifier. A nil notifier is accepted (no-op).
func NewPipeline(notifier Notifier, opts ...PipelineOption) *Pipeline {
	p := &Pipeline{
		notifier: notifier,
		logger:   slog.Default(),
	}
	for _, o := range opts {
		o(p)
	}
	return p
}

// Notify runs the event through all registered filters. The first filter that
// suppresses the event short-circuits the chain and returns nil. If all filters
// allow the event it is forwarded to the underlying notifier.
func (p *Pipeline) Notify(ctx context.Context, event Event) error {
	for _, f := range p.filters {
		if !f.Allow(event) {
			p.logger.Debug("alert suppressed by filter",
				"event_type", event.Type,
				"port", portFromListener(event.Listener),
			)
			return nil
		}
	}
	if p.notifier == nil {
		return nil
	}
	return p.notifier.Notify(ctx, event)
}

// portFromListener is a helper that safely extracts the port number from a
// listener for logging purposes, returning 0 when the listener is nil.
func portFromListener(l *ports.Listener) uint16 {
	if l == nil {
		return 0
	}
	return l.Port
}
