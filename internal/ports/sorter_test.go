package ports

import (
	"net/netip"
	"testing"
)

func makeSortListener(proto, addr string, port uint16, pid int) Listener {
	ip := netip.MustParseAddr(addr)
	l := Listener{
		Protocol: proto,
		Address:  netip.AddrPortFrom(ip, port),
	}
	if pid >= 0 {
		l.Process = &ProcessInfo{PID: pid}
	}
	return l
}

func TestParseSortField_Known(t *testing.T) {
	cases := []struct {
		input string
		want  SortField
	}{
		{"port", SortByPort},
		{"PORT", SortByPort},
		{"protocol", SortByProtocol},
		{"proto", SortByProtocol},
		{"address", SortByAddress},
		{"addr", SortByAddress},
		{"pid", SortByPID},
	}
	for _, tc := range cases {
		got, ok := ParseSortField(tc.input)
		if !ok {
			t.Errorf("ParseSortField(%q) ok=false", tc.input)
		}
		if got != tc.want {
			t.Errorf("ParseSortField(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestParseSortField_Unknown(t *testing.T) {
	got, ok := ParseSortField("unknown")
	if ok {
		t.Error("expected ok=false for unknown field")
	}
	if got != SortByPort {
		t.Errorf("expected default SortByPort, got %v", got)
	}
}

func TestSortListeners_ByPort(t *testing.T) {
	listeners := []Listener{
		makeSortListener("tcp", "0.0.0.0", 8080, -1),
		makeSortListener("tcp", "0.0.0.0", 443, -1),
		makeSortListener("tcp", "0.0.0.0", 22, -1),
	}
	SortListeners(listeners, SortByPort)
	ports := []uint16{listeners[0].Address.Port(), listeners[1].Address.Port(), listeners[2].Address.Port()}
	if ports[0] != 22 || ports[1] != 443 || ports[2] != 8080 {
		t.Errorf("unexpected order: %v", ports)
	}
}

func TestSortListeners_ByProtocol(t *testing.T) {
	listeners := []Listener{
		makeSortListener("udp", "0.0.0.0", 53, -1),
		makeSortListener("tcp", "0.0.0.0", 80, -1),
		makeSortListener("tcp", "0.0.0.0", 22, -1),
	}
	SortListeners(listeners, SortByProtocol)
	if listeners[0].Protocol != "tcp" || listeners[1].Protocol != "tcp" || listeners[2].Protocol != "udp" {
		t.Errorf("unexpected protocol order: %v %v %v",
			listeners[0].Protocol, listeners[1].Protocol, listeners[2].Protocol)
	}
	// Secondary sort by port within same protocol.
	if listeners[0].Address.Port() != 22 || listeners[1].Address.Port() != 80 {
		t.Errorf("expected secondary port sort, got %d %d",
			listeners[0].Address.Port(), listeners[1].Address.Port())
	}
}

func TestSortListeners_ByPID_NilProcessLast(t *testing.T) {
	listeners := []Listener{
		makeSortListener("tcp", "0.0.0.0", 80, -1),
		makeSortListener("tcp", "0.0.0.0", 22, 1001),
		makeSortListener("tcp", "0.0.0.0", 443, 500),
	}
	SortListeners(listeners, SortByPID)
	if listeners[0].Process == nil || listeners[0].Process.PID != 500 {
		t.Errorf("expected PID 500 first, got %v", listeners[0].Process)
	}
	if listeners[2].Process != nil {
		t.Errorf("expected nil process last")
	}
}
