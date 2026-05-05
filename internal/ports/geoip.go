package ports

import (
	"fmt"
	"net"
	"sync"
)

// GeoInfo holds geographic metadata resolved for an IP address.
type GeoInfo struct {
	CountryCode string
	CountryName string
	City        string
}

// GeoLookup is a function type that resolves an IP to GeoInfo.
type GeoLookup func(ip net.IP) (*GeoInfo, error)

// GeoEnricher attaches GeoInfo to Listener values using a pluggable lookup
// function and an in-process LRU-style cache.
type GeoEnricher struct {
	lookup GeoLookup
	mu    sync.Mutex
	cache map[string]*GeoInfo
	max   int
}

// NewGeoEnricher creates a GeoEnricher backed by the provided lookup function.
// cacheSize controls the maximum number of cached entries (0 disables caching).
func NewGeoEnricher(lookup GeoLookup, cacheSize int) *GeoEnricher {
	if cacheSize < 0 {
		cacheSize = 0
	}
	return &GeoEnricher{
		lookup: lookup,
		cache:  make(map[string]*GeoInfo, cacheSize),
		max:    cacheSize,
	}
}

// Enrich resolves GeoInfo for each Listener whose address is a non-loopback
// IP and attaches it via Listener.Geo. Listeners without a parseable IP are
// left unchanged.
func (e *GeoEnricher) Enrich(listeners []Listener) []Listener {
	out := make([]Listener, len(listeners))
	copy(out, listeners)
	for i, l := range out {
		ip := net.ParseIP(l.Address)
		if ip == nil || ip.IsLoopback() {
			continue
		}
		geo, err := e.resolve(ip)
		if err != nil {
			continue
		}
		out[i].Geo = geo
	}
	return out
}

func (e *GeoEnricher) resolve(ip net.IP) (*GeoInfo, error) {
	key := ip.String()
	e.mu.Lock()
	if cached, ok := e.cache[key]; ok {
		e.mu.Unlock()
		return cached, nil
	}
	e.mu.Unlock()

	geo, err := e.lookup(ip)
	if err != nil {
		return nil, fmt.Errorf("geoip lookup %s: %w", key, err)
	}

	if e.max > 0 {
		e.mu.Lock()
		if len(e.cache) >= e.max {
			// simple eviction: clear all when full
			e.cache = make(map[string]*GeoInfo, e.max)
		}
		e.cache[key] = geo
		e.mu.Unlock()
	}
	return geo, nil
}
