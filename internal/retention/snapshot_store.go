package retention

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// SnapshotEntry holds a captured set of listeners with a timestamp.
type SnapshotEntry struct {
	CapturedAt time.Time
	Listeners  []ports.Listener
}

// SnapshotStore keeps a bounded, time-limited ring of scan snapshots.
type SnapshotStore struct {
	mu           sync.RWMutex
	entries      []SnapshotEntry
	maxSnapshots int
}

// NewSnapshotStore creates a store with the given capacity.
func NewSnapshotStore(maxSnapshots int) *SnapshotStore {
	return &SnapshotStore{maxSnapshots: maxSnapshots}
}

// Add appends a new snapshot, evicting the oldest if at capacity.
func (s *SnapshotStore) Add(listeners []ports.Listener) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry := SnapshotEntry{
		CapturedAt: time.Now(),
		Listeners:  listeners,
	}
	s.entries = append(s.entries, entry)
	if len(s.entries) > s.maxSnapshots {
		s.entries = s.entries[len(s.entries)-s.maxSnapshots:]
	}
}

// Recent returns all snapshots captured within the given age.
func (s *SnapshotStore) Recent(maxAge time.Duration) []SnapshotEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cutoff := time.Now().Add(-maxAge)
	var result []SnapshotEntry
	for _, e := range s.entries {
		if e.CapturedAt.After(cutoff) {
			result = append(result, e)
		}
	}
	return result
}

// Prune removes entries older than maxAge and returns the count removed.
func (s *SnapshotStore) Prune(maxAge time.Duration) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	cutoff := time.Now().Add(-maxAge)
	var kept []SnapshotEntry
	for _, e := range s.entries {
		if e.CapturedAt.After(cutoff) {
			kept = append(kept, e)
		}
	}
	removed := len(s.entries) - len(kept)
	s.entries = kept
	return removed
}

// Len returns the current number of stored snapshots.
func (s *SnapshotStore) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}
