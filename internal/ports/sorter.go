package ports

import (
	"sort"
	"strings"
)

// SortField defines the field by which listeners are sorted.
type SortField int

const (
	SortByPort SortField = iota
	SortByProtocol
	SortByAddress
	SortByPID
)

// ParseSortField parses a string into a SortField.
// Returns SortByPort and false if the field is unrecognized.
func ParseSortField(s string) (SortField, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "port":
		return SortByPort, true
	case "protocol", "proto":
		return SortByProtocol, true
	case "address", "addr":
		return SortByAddress, true
	case "pid":
		return SortByPID, true
	default:
		return SortByPort, false
	}
}

// SortListeners sorts a slice of Listener values in-place by the given field.
// Within the same primary field value, listeners are secondarily sorted by port
// to produce a stable, deterministic order.
func SortListeners(listeners []Listener, by SortField) {
	sort.SliceStable(listeners, func(i, j int) bool {
		a, b := listeners[i], listeners[j]
		switch by {
		case SortByProtocol:
			if a.Protocol != b.Protocol {
				return a.Protocol < b.Protocol
			}
		case SortByAddress:
			ai := a.Address.Addr().String()
			bi := b.Address.Addr().String()
			if ai != bi {
				return ai < bi
			}
		case SortByPID:
			apid := pidOf(a)
			bpid := pidOf(b)
			if apid != bpid {
				return apid < bpid
			}
		}
		// Secondary sort: port ascending.
		return a.Address.Port() < b.Address.Port()
	})
}

func pidOf(l Listener) int {
	if l.Process == nil {
		return -1
	}
	return l.Process.PID
}
