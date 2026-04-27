package baseline

import (
	"log/slog"

	"github.com/dreadster3/portwatch/internal/ports"
)

// Filter wraps a Checker and a Learner to provide baseline-aware filtering
// of port listeners. During the learning window, all listeners are observed
// and passed through. After the window, only unknown listeners are returned.
type Filter struct {
	checker *Checker
	learner *Learner
	logger  *slog.Logger
}

// NewFilter constructs a Filter using the provided Checker and Learner.
func NewFilter(checker *Checker, learner *Learner, logger *slog.Logger) *Filter {
	return &Filter{
		checker: checker,
		learner: learner,
		logger:  logger,
	}
}

// Apply processes a slice of listeners. While still learning, every listener
// is observed and the full slice is returned unchanged. Once the learning
// window has closed, only listeners not present in the baseline are returned.
func (f *Filter) Apply(listeners []ports.Listener) []ports.Listener {
	if f.learner.IsLearning() {
		for _, l := range listeners {
			f.learner.Observe(l)
		}
		f.logger.Debug("baseline filter: learning phase, passing all listeners",
			"count", len(listeners))
		return listeners
	}

	unknown := f.checker.FilterUnknown(listeners)
	if len(unknown) > 0 {
		f.logger.Info("baseline filter: unknown listeners detected",
			"unknown", len(unknown),
			"total", len(listeners))
	}
	return unknown
}
