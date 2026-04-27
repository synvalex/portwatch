package baseline

import (
	"sync"
	"time"

	"github.com/example/portwatch/internal/ports"
)

// Learner observes listeners during a learning window and adds them to the
// baseline store. Once the window expires, Observe becomes a no-op.
type Learner struct {
	mu       sync.Mutex
	store    *Store
	deadline time.Time
}

// NewLearner creates a Learner that records observations for the given duration.
func NewLearner(store *Store, window time.Duration) *Learner {
	return &Learner{
		store:    store,
		deadline: time.Now().Add(window),
	}
}

// Observe records the listener in the baseline if the learning window is still
// open. It is safe to call concurrently.
func (l *Learner) Observe(listener ports.Listener) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if time.Now().Before(l.deadline) {
		l.store.Add(listener)
	}
}

// IsLearning reports whether the learning window is still active.
func (l *Learner) IsLearning() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return time.Now().Before(l.deadline)
}

// Remaining returns how much time is left in the learning window.
// Returns zero if the window has already closed.
func (l *Learner) Remaining() time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()

	rem := time.Until(l.deadline)
	if rem < 0 {
		return 0
	}
	return rem
}
