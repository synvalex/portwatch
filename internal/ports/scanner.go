package ports

import (
	"context"
	"fmt"
	"net"
	"strings"
)

// Listener represents a single open port on the host.
type Listener struct {
	Proto   string
	Address net.IP
	Port    uint16
	PID     int
}

// String returns a human-readable representation of the listener.
func (l Listener) String() string {
	return fmt.Sprintf("%s://%s:%d (pid %d)", l.Proto, l.Address, l.Port, l.PID)
}

// Scanner is the interface implemented by platform-specific scanners.
type Scanner interface {
	Scan(ctx context.Context) ([]Listener, error)
}

// ParseAddress parses a "host:port" or ":port" string into an IP and port.
func ParseAddress(addr string) (net.IP, uint16, error) {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid address %q: %w", addr, err)
	}
	var port uint16
	if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
		return nil, 0, fmt.Errorf("invalid port %q: %w", portStr, err)
	}
	var ip net.IP
	if host == "" || host == "0.0.0.0" || host == "::" {
		ip = net.IPv4zero
	} else {
		ip = net.ParseIP(strings.TrimSpace(host))
		if ip == nil {
			return nil, 0, fmt.Errorf("invalid IP address %q", host)
		}
	}
	return ip, port, nil
}

// DeduplicateListeners removes duplicate listeners from the slice.
func DeduplicateListeners(listeners []Listener) []Listener {
	seen := make(map[string]struct{}, len(listeners))
	out := make([]Listener, 0, len(listeners))
	for _, l := range listeners {
		key := fmt.Sprintf("%s|%s|%d", l.Proto, l.Address.String(), l.Port)
		if _, ok := seen[key]; !ok {
			seen[key] = struct{}{}
			out = append(out, l)
		}
	}
	return out
}
