package alert

import (
	"fmt"
	"sync"
	"time"

	"github.com/wesen/portwatch/internal/ports"
)

// ThrottleFilter suppresses alerts that exceed a rate limit within a
// rolling window, optionally keyed per port.
type ThrottleFilter struct {
	mu      sync.Mutex
	rate    int
	window  time.Duration
	perPort bool
	buckets map[string][]time.Time
	now     func() time.Time
}

// NewThrottleFilter creates a ThrottleFilter with the given parameters.
func NewThrottleFilter(rate int, window time.Duration, perPort bool) *ThrottleFilter {
	return &ThrottleFilter{
		rate:    rate,
		window:  window,
		perPort: perPort,
		buckets: make(map[string][]time.Time),
		now:     time.Now,
	}
}

// Allow returns true if the event should be forwarded, false if throttled.
func (f *ThrottleFilter) Allow(l ports.Listener, eventType string) bool {
	if f.rate <= 0 {
		return true
	}
	f.mu.Lock()
	defer f.mu.Unlock()

	key := f.buildKey(l, eventType)
	now := f.now()
	cutoff := now.Add(-f.window)

	times := f.buckets[key]
	valid := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= f.rate {
		f.buckets[key] = valid
		return false
	}
	f.buckets[key] = append(valid, now)
	return true
}

func (f *ThrottleFilter) buildKey(l ports.Listener, eventType string) string {
	if f.perPort {
		return fmt.Sprintf("%s:%d:%s", l.Proto, l.Port, eventType)
	}
	return eventType
}
