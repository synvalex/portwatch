package alert

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// RateLimiter suppresses alerts that exceed a maximum count within a sliding
// time window, preventing alert storms when many ports appear simultaneously.
type RateLimiter struct {
	mu       sync.Mutex
	maxBurst int
	window   time.Duration
	buckets  map[string][]time.Time
	now      func() time.Time
}

// NewRateLimiter creates a RateLimiter that allows at most maxBurst alerts of
// the same event type within window. A zero or negative maxBurst disables
// rate limiting (all events are allowed).
func NewRateLimiter(maxBurst int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		maxBurst: maxBurst,
		window:   window,
		buckets:  make(map[string][]time.Time),
		now:      time.Now,
	}
}

// Allow returns true when the event should be forwarded, false when it should
// be suppressed. The key is derived from the listener address and event type.
func (r *RateLimiter) Allow(l ports.Listener, eventType string) bool {
	if r.maxBurst <= 0 {
		return true
	}

	key := eventType + ":" + l.String()
	now := r.now()
	cutoff := now.Add(-r.window)

	r.mu.Lock()
	defer r.mu.Unlock()

	times := r.buckets[key]

	// Evict timestamps outside the window.
	valid := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= r.maxBurst {
		r.buckets[key] = valid
		return false
	}

	r.buckets[key] = append(valid, now)
	return true
}

// Reset clears all rate-limit state. Useful for testing or config reload.
func (r *RateLimiter) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.buckets = make(map[string][]time.Time)
}
