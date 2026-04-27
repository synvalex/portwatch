package alert

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// dedupKey uniquely identifies an alert event for deduplication purposes.
type dedupKey struct {
	protocol string
	address  string
	port     uint16
	eventType string
}

// DedupFilter suppresses repeated alerts for the same listener within a
// configurable time window.
type DedupFilter struct {
	mu     sync.Mutex
	seen   map[dedupKey]time.Time
	window time.Duration
	now    func() time.Time
}

// NewDedupFilter creates a DedupFilter with the given deduplication window.
func NewDedupFilter(window time.Duration) *DedupFilter {
	return &DedupFilter{
		seen:   make(map[dedupKey]time.Time),
		window: window,
		now:    time.Now,
	}
}

// IsDuplicate returns true if an identical alert was already seen within the
// deduplication window. It records the event if it is not a duplicate.
func (d *DedupFilter) IsDuplicate(l ports.Listener, eventType string) bool {
	key := dedupKey{
		protocol:  l.Protocol,
		address:   l.Address.String(),
		port:      l.Port,
		eventType: eventType,
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	if last, ok := d.seen[key]; ok && now.Sub(last) < d.window {
		return true
	}
	d.seen[key] = now
	return false
}

// Evict removes stale entries older than the deduplication window.
// Call periodically to prevent unbounded memory growth.
func (d *DedupFilter) Evict() {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	for k, t := range d.seen {
		if now.Sub(t) >= d.window {
			delete(d.seen, k)
		}
	}
}
