package ports

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// DNSEnricher performs reverse-DNS lookups and caches results.
type DNSEnricher struct {
	cache   map[string]string
	mu      sync.RWMutex
	timeout time.Duration
	workers int
	log     zerolog.Logger
}

// NewDNSEnricher creates a DNSEnricher with the given cache capacity and timeout.
func NewDNSEnricher(timeout time.Duration, workers int, log zerolog.Logger) *DNSEnricher {
	return &DNSEnricher{
		cache:   make(map[string]string),
		timeout: timeout,
		workers: workers,
		log:     log,
	}
}

// Enrich performs reverse-DNS lookups on all listeners concurrently and
// populates the Hostname field.
func (e *DNSEnricher) Enrich(listeners []Listener) []Listener {
	type job struct {
		idx  int
		addr string
	}

	jobs := make(chan job, len(listeners))
	for i, l := range listeners {
		if l.Address.IP != nil {
			jobs <- job{i, l.Address.IP.String()}
		}
	}
	close(jobs)

	var wg sync.WaitGroup
	w := e.workers
	if w <= 0 {
		w = 4
	}
	for range make([]struct{}, w) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				host := e.lookup(j.addr)
				if host != "" {
					listeners[j.idx].Hostname = host
				}
			}
		}()
	}
	wg.Wait()
	return listeners
}

func (e *DNSEnricher) lookup(ip string) string {
	e.mu.RLock()
	if h, ok := e.cache[ip]; ok {
		e.mu.RUnlock()
		return h
	}
	e.mu.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	names, err := net.DefaultResolver.LookupAddr(ctx, ip)
	if err != nil || len(names) == 0 {
		e.log.Debug().Str("ip", ip).Msg("reverse DNS lookup failed or empty")
		return ""
	}
	host := names[0]

	e.mu.Lock()
	e.cache[ip] = host
	e.mu.Unlock()
	return host
}
