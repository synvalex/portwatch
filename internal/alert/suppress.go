package alert

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// SuppressFilter suppresses repeated alerts for the same listener
// until a configurable quiet window has elapsed since the last suppression.
type SuppressFilter struct {
	mu      sync.Mutex
	window  time.Duration
	max     int
	counts  map[string]int
	expiry  map[string]time.Time
	nowFunc func() time.Time
}

// NewSuppressFilter creates a SuppressFilter that suppresses an alert once
// the same key has fired more than max times within window.
func NewSuppressFilter(window time.Duration, max int) *SuppressFilter {
	return &SuppressFilter{
		window:  window,
		max:     max,
		counts:  make(map[string]int),
		expiry:  make(map[string]time.Time),
		nowFunc: time.Now,
	}
}

// Allow returns true when the event should be forwarded.
// Once a key exceeds max occurrences within window it is suppressed until
// the window resets.
func (s *SuppressFilter) Allow(e Event) bool {
	if s.window <= 0 || s.max <= 0 {
		return true
	}

	key := suppressKey(e.Listener, e.Type)
	now := s.nowFunc()

	s.mu.Lock()
	defer s.mu.Unlock()

	if exp, ok := s.expiry[key]; ok && now.After(exp) {
		delete(s.counts, key)
		delete(s.expiry, key)
	}

	s.counts[key]++
	if _, ok := s.expiry[key]; !ok {
		s.expiry[key] = now.Add(s.window)
	}

	return s.counts[key] <= s.max
}

func suppressKey(l ports.Listener, t EventType) string {
	return l.Address.String() + "|" + l.Protocol + "|" + string(t)
}
