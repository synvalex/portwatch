package alert

import (
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func dedupListener(proto, ip string, port uint16) ports.Listener {
	return ports.Listener{
		Protocol: proto,
		Address:  net.ParseIP(ip),
		Port:     port,
	}
}

func TestDedupFilter_FirstEventNotDuplicate(t *testing.T) {
	f := NewDedupFilter(30 * time.Second)
	l := dedupListener("tcp", "0.0.0.0", 8080)

	if f.IsDuplicate(l, "appeared") {
		t.Fatal("expected first event to not be a duplicate")
	}
}

func TestDedupFilter_SecondEventWithinWindowIsDuplicate(t *testing.T) {
	f := NewDedupFilter(30 * time.Second)
	l := dedupListener("tcp", "0.0.0.0", 8080)

	f.IsDuplicate(l, "appeared")
	if !f.IsDuplicate(l, "appeared") {
		t.Fatal("expected second event within window to be a duplicate")
	}
}

func TestDedupFilter_ExpiredWindowNotDuplicate(t *testing.T) {
	now := time.Now()
	f := NewDedupFilter(5 * time.Second)
	f.now = func() time.Time { return now }

	l := dedupListener("tcp", "0.0.0.0", 9090)
	f.IsDuplicate(l, "appeared")

	// Advance time beyond the window.
	f.now = func() time.Time { return now.Add(10 * time.Second) }
	if f.IsDuplicate(l, "appeared") {
		t.Fatal("expected event after window expiry to not be a duplicate")
	}
}

func TestDedupFilter_DifferentEventTypesAreIndependent(t *testing.T) {
	f := NewDedupFilter(30 * time.Second)
	l := dedupListener("tcp", "0.0.0.0", 443)

	f.IsDuplicate(l, "appeared")
	if f.IsDuplicate(l, "disappeared") {
		t.Fatal("expected different event types to be treated independently")
	}
}

func TestDedupFilter_Evict_RemovesStaleEntries(t *testing.T) {
	now := time.Now()
	f := NewDedupFilter(5 * time.Second)
	f.now = func() time.Time { return now }

	l := dedupListener("udp", "127.0.0.1", 53)
	f.IsDuplicate(l, "appeared")

	if len(f.seen) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(f.seen))
	}

	f.now = func() time.Time { return now.Add(10 * time.Second) }
	f.Evict()

	if len(f.seen) != 0 {
		t.Fatalf("expected 0 entries after eviction, got %d", len(f.seen))
	}
}
