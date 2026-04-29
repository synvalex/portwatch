package alert

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
)

// MultiNotifier fans out alert events to multiple Notifier implementations.
// All notifiers are called; errors are collected and logged but do not short-circuit.
type MultiNotifier struct {
	notifiers []Notifier
}

// NewMultiNotifier returns a MultiNotifier that dispatches to each provided Notifier.
func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{notifiers: notifiers}
}

// Notify sends the event to every registered notifier.
// It returns a combined error if one or more notifiers fail.
func (m *MultiNotifier) Notify(event Event) error {
	var errs []string

	for _, n := range m.notifiers {
		if err := n.Notify(event); err != nil {
			log.Warn().
				Err(err).
				Str("notifier", fmt.Sprintf("%T", n)).
				Msg("notifier returned error")
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("multi_notifier: %d error(s): %s", len(errs), strings.Join(errs, "; "))
	}
	return nil
}

// Add appends a Notifier to the fan-out list.
func (m *MultiNotifier) Add(n Notifier) {
	m.notifiers = append(m.notifiers, n)
}

// Len returns the number of registered notifiers.
func (m *MultiNotifier) Len() int {
	return len(m.notifiers)
}

// Remove removes the first notifier from the fan-out list that matches the
// provided value by pointer equality. It reports whether a notifier was removed.
func (m *MultiNotifier) Remove(n Notifier) bool {
	for i, existing := range m.notifiers {
		if existing == n {
			m.notifiers = append(m.notifiers[:i], m.notifiers[i+1:]...)
			return true
		}
	}
	return false
}
