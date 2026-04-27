package alert

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// CooldownFilter suppresses repeated alerts for the same listener+event
// combination until a cooldown period has elapsed since the last alert.
type CooldownFilter struct {
	mu       sync.Mutex
	window   time.Duration
	lastSeen map[string]time.Time
	now      func() time.Time
}

// NewCooldownFilter creates a CooldownFilter with the given cooldown window.
func NewCooldownFilter(window time.Duration) *CooldownFilter {
	return &CooldownFilter{
		window:   window,
		lastSeen: make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if the event should be forwarded (i.e. not in cooldown).
func (c *CooldownFilter) Allow(l ports.Listener, eventType string) bool {
	if c.window <= 0 {
		return true
	}

	key := cooldownKey(l, eventType)
	now := c.now()

	c.mu.Lock()
	defer c.mu.Unlock()

	if last, ok := c.lastSeen[key]; ok {
		if now.Sub(last) < c.window {
			return false
		}
	}

	c.lastSeen[key] = now
	return true
}

// Purge removes all expired entries, freeing memory for long-running daemons.
func (c *CooldownFilter) Purge() {
	now := c.now()
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, t := range c.lastSeen {
		if now.Sub(t) >= c.window {
			delete(c.lastSeen, k)
		}
	}
}

func cooldownKey(l ports.Listener, eventType string) string {
	return eventType + ":" + l.String()
}
