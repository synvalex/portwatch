package history

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Entry is a single recorded alert event.
type Entry struct {
	Event     alert.Event
	Timestamp time.Time
}

// Store is a thread-safe ring-buffer of alert events.
type Store struct {
	mu        sync.RWMutex
	entries   []Entry
	maxEvents int
	retention time.Duration
	now       func() time.Time
}

// NewStore creates a Store with the given capacity and retention window.
func NewStore(maxEvents int, retention time.Duration) *Store {
	return &Store{
		maxEvents: maxEvents,
		retention: retention,
		now:       time.Now,
	}
}

// Add appends a new event, evicting oldest entries beyond capacity or retention.
func (s *Store) Add(ev alert.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries = append(s.entries, Entry{Event: ev, Timestamp: s.now()})
	s.evict()
}

// Recent returns all entries that fall within the retention window.
func (s *Store) Recent() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cutoff := s.now().Add(-s.retention)
	var out []Entry
	for _, e := range s.entries {
		if e.Timestamp.After(cutoff) {
			out = append(out, e)
		}
	}
	return out
}

// Len returns the current number of stored entries.
func (s *Store) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}

// evict removes entries that exceed capacity or retention. Must be called with lock held.
func (s *Store) evict() {
	cutoff := s.now().Add(-s.retention)
	start := 0
	for start < len(s.entries) && s.entries[start].Timestamp.Before(cutoff) {
		start++
	}
	s.entries = s.entries[start:]
	if len(s.entries) > s.maxEvents {
		s.entries = s.entries[len(s.entries)-s.maxEvents:]
	}
}
