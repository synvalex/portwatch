package ports

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Listener represents an open port with its associated process info.
type Listener struct {
	Protocol string
	Address  string
	Port     int
	PID      int
	Process  string
}

// String returns a human-readable representation of the listener.
func (l Listener) String() string {
	return fmt.Sprintf("%s %s:%d (pid=%d, process=%s)", l.Protocol, l.Address, l.Port, l.PID, l.Process)
}

// Scanner defines the interface for scanning open ports.
type Scanner interface {
	Scan() ([]Listener, error)
}

// ParseAddress splits a combined address:port string into its components.
func ParseAddress(addr string) (string, int, error) {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		// Try treating the whole string as just a port
		portStr = strings.TrimSpace(addr)
		host = "0.0.0.0"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid port %q: %w", portStr, err)
	}
	if port < 1 || port > 65535 {
		return "", 0, fmt.Errorf("port %d out of valid range", port)
	}
	return host, port, nil
}

// DeduplicateListeners removes duplicate entries from a listener slice.
func DeduplicateListeners(listeners []Listener) []Listener {
	seen := make(map[string]struct{})
	result := make([]Listener, 0, len(listeners))
	for _, l := range listeners {
		key := fmt.Sprintf("%s:%s:%d", l.Protocol, l.Address, l.Port)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, l)
		}
	}
	return result
}
