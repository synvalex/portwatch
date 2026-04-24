package ports

import "net"

// Filter holds criteria for excluding listeners from scan results.
type Filter struct {
	// ExcludeLoopback skips listeners bound only to loopback addresses.
	ExcludeLoopback bool
	// ExcludePorts is a set of port numbers to always ignore.
	ExcludePorts map[uint16]struct{}
}

// NewFilter constructs a Filter from configuration primitives.
func NewFilter(excludeLoopback bool, excludePorts []uint16) *Filter {
	set := make(map[uint16]struct{}, len(excludePorts))
	for _, p := range excludePorts {
		set[p] = struct{}{}
	}
	return &Filter{
		ExcludeLoopback: excludeLoopback,
		ExcludePorts:    set,
	}
}

// Apply returns only the listeners that pass all filter criteria.
func (f *Filter) Apply(listeners []Listener) []Listener {
	if f == nil {
		return listeners
	}
	out := listeners[:0:0]
	for _, l := range listeners {
		if f.ExcludeLoopback && isLoopback(l.IP) {
			continue
		}
		if _, skip := f.ExcludePorts[l.Port]; skip {
			continue
		}
		out = append(out, l)
	}
	return out
}

func isLoopback(ip net.IP) bool {
	if ip == nil {
		return false
	}
	return ip.IsLoopback()
}
