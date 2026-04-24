// Package snapshot provides functionality to track and diff port listener state
// between successive scans, enabling detection of new or removed listeners.
package snapshot

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/ports"
)

// Key uniquely identifies a listener by protocol, address, and port.
type Key struct {
	Proto string
	Addr  string
	Port  uint16
}

// keyOf constructs a Key from a Listener.
func keyOf(l ports.Listener) Key {
	return Key{
		Proto: l.Proto,
		Addr:  l.Addr,
		Port:  l.Port,
	}
}

// Diff holds the result of comparing two snapshots.
type Diff struct {
	Appeared []ports.Listener
	Disappeared []ports.Listener
}

// String returns a human-readable summary of the diff.
func (d Diff) String() string {
	return fmt.Sprintf("appeared=%d disappeared=%d", len(d.Appeared), len(d.Disappeared))
}

// IsEmpty reports whether the diff contains no changes.
func (d Diff) IsEmpty() bool {
	return len(d.Appeared) == 0 && len(d.Disappeared) == 0
}

// Store holds the most recent set of listeners and can compute diffs.
type Store struct {
	mu      sync.Mutex
	current map[Key]ports.Listener
}

// NewStore creates an empty Store.
func NewStore() *Store {
	return &Store{
		current: make(map[Key]ports.Listener),
	}
}

// Update replaces the stored snapshot with next and returns the Diff.
func (s *Store) Update(next []ports.Listener) Diff {
	s.mu.Lock()
	defer s.mu.Unlock()

	nextMap := make(map[Key]ports.Listener, len(next))
	for _, l := range next {
		nextMap[keyOf(l)] = l
	}

	var diff Diff

	for k, l := range nextMap {
		if _, ok := s.current[k]; !ok {
			diff.Appeared = append(diff.Appeared, l)
		}
	}

	for k, l := range s.current {
		if _, ok := nextMap[k]; !ok {
			diff.Disappeared = append(diff.Disappeared, l)
		}
	}

	s.current = nextMap
	return diff
}

// Current returns a copy of the current snapshot.
func (s *Store) Current() []ports.Listener {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]ports.Listener, 0, len(s.current))
	for _, l := range s.current {
		out = append(out, l)
	}
	return out
}
