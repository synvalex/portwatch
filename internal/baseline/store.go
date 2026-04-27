package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/user/portwatch/internal/ports"
)

// entry represents a persisted baseline listener key.
type entry struct {
	Proto   string `json:"proto"`
	Address string `json:"address"`
	Port    uint16 `json:"port"`
}

func keyOf(l ports.Listener) entry {
	return entry{Proto: l.Proto, Address: l.Address, Port: l.Port}
}

// Store holds the set of expected (baseline) listeners.
type Store struct {
	mu      sync.RWMutex
	known   map[entry]struct{}
	filePath string
}

// NewStore creates an empty baseline store.
func NewStore(filePath string) *Store {
	return &Store{
		known:    make(map[entry]struct{}),
		filePath: filePath,
	}
}

// Add records a listener as part of the baseline.
func (s *Store) Add(l ports.Listener) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.known[keyOf(l)] = struct{}{}
}

// Contains reports whether a listener is in the baseline.
func (s *Store) Contains(l ports.Listener) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.known[keyOf(l)]
	return ok
}

// Save persists the baseline to disk as JSON.
func (s *Store) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entries := make([]entry, 0, len(s.known))
	for e := range s.known {
		entries = append(entries, e)
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	return os.WriteFile(s.filePath, data, 0o644)
}

// Load reads a previously saved baseline from disk.
func (s *Store) Load() error {
	data, err := os.ReadFile(s.filePath)
	if os.IsNotExist(err) {
		return nil // no baseline yet; start empty
	}
	if err != nil {
		return fmt.Errorf("baseline: read: %w", err)
	}
	var entries []entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("baseline: unmarshal: %w", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range entries {
		s.known[e] = struct{}{}
	}
	return nil
}
