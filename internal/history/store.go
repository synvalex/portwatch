package history

import (
	"sync"
	"time"
)

// Store holds a bounded, time-limited ring of port change events.
type Store struct {
	mu        sync.Mutex
	events    []Event
	maxEvents int
	retention time.Duration
}

// NewStore creates a Store with the given capacity and retention window.
func NewStore(maxEvents int, retention time.Duration) *Store {
	return &Store{
		events:    make([]Event, 0, maxEvents),
		maxEvents: maxEvents,
		retention: retention,
	}
}

// Add appends an event to the store, evicting old entries as needed.
func (s *Store) Add(e Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evict()
	if len(s.events) >= s.maxEvents {
		s.events = s.events[1:]
	}
	s.events = append(s.events, e)
}

// Recent returns all events that occurred within the given window.
func (s *Store) Recent(window time.Duration) []Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evict()
	cutoff := time.Now().Add(-window)
	var out []Event
	for _, e := range s.events {
		if e.OccurredAt.After(cutoff) {
			out = append(out, e)
		}
	}
	return out
}

// Len returns the current number of stored events.
func (s *Store) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.events)
}

// evict removes events older than the retention window. Caller must hold mu.
func (s *Store) evict() {
	cutoff := time.Now().Add(-s.retention)
	i := 0
	for i < len(s.events) && s.events[i].OccurredAt.Before(cutoff) {
		i++
	}
	if i > 0 {
		s.events = s.events[i:]
	}
}
