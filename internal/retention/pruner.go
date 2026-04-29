package retention

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// Prunable is implemented by any store that can remove stale entries.
type Prunable interface {
	Prune(maxAge time.Duration) int
}

// Pruner periodically calls Prune on registered stores.
type Pruner struct {
	stores   []Prunable
	maxAge   time.Duration
	interval time.Duration
	logger   *slog.Logger
	mu       sync.Mutex
}

// NewPruner creates a Pruner with the given interval and max age.
func NewPruner(interval, maxAge time.Duration, logger *slog.Logger) *Pruner {
	return &Pruner{
		interval: interval,
		maxAge:   maxAge,
		logger:   logger,
	}
}

// Register adds a Prunable store to the pruner.
func (p *Pruner) Register(store Prunable) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.stores = append(p.stores, store)
}

// Run starts the pruning loop and blocks until ctx is cancelled.
func (p *Pruner) Run(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			p.logger.Debug("retention pruner stopped")
			return
		case <-ticker.C:
			p.pruneAll()
		}
	}
}

func (p *Pruner) pruneAll() {
	p.mu.Lock()
	stores := make([]Prunable, len(p.stores))
	copy(stores, p.stores)
	p.mu.Unlock()

	total := 0
	for _, s := range stores {
		total += s.Prune(p.maxAge)
	}
	if total > 0 {
		p.logger.Debug("retention pruner evicted entries", "count", total)
	}
}
