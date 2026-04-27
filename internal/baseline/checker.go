package baseline

import (
	"github.com/example/portwatch/internal/ports"
)

// Checker determines whether a listener is known to the baseline store.
type Checker struct {
	store *Store
}

// NewChecker creates a Checker backed by the given Store.
func NewChecker(store *Store) *Checker {
	return &Checker{store: store}
}

// IsKnown returns true if the listener exists in the baseline.
func (c *Checker) IsKnown(l ports.Listener) bool {
	return c.store.Contains(l)
}

// FilterUnknown returns only the listeners that are NOT in the baseline.
func (c *Checker) FilterUnknown(listeners []ports.Listener) []ports.Listener {
	var unknown []ports.Listener
	for _, l := range listeners {
		if !c.store.Contains(l) {
			unknown = append(unknown, l)
		}
	}
	return unknown
}

// FilterKnown returns only the listeners that ARE in the baseline.
func (c *Checker) FilterKnown(listeners []ports.Listener) []ports.Listener {
	var known []ports.Listener
	for _, l := range listeners {
		if c.store.Contains(l) {
			known = append(known, l)
		}
	}
	return known
}
