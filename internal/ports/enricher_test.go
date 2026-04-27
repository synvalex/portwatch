package ports

import (
	"errors"
	"net"
	"testing"
)

func makeListenerWithInode(port uint16, inode uint64) Listener {
	return Listener{
		Addr:     net.TCPAddr{IP: net.ParseIP("0.0.0.0"), Port: int(port)},
		Protocol: "tcp",
		Inode:    inode,
	}
}

func TestEnricher_SkipsZeroInode(t *testing.T) {
	called := false
	e := newEnricherWithLookup(func(inode uint64) (*ProcessInfo, error) {
		called = true
		return nil, nil
	})

	listeners := []Listener{makeListenerWithInode(80, 0)}
	out := e.Enrich(listeners)

	if called {
		t.Error("lookup should not be called for inode=0")
	}
	if out[0].Process != nil {
		t.Error("expected nil Process for zero inode")
	}
}

func TestEnricher_AttachesProcess(t *testing.T) {
	expected := &ProcessInfo{PID: 123, Name: "nginx"}
	e := newEnricherWithLookup(func(inode uint64) (*ProcessInfo, error) {
		if inode == 55 {
			return expected, nil
		}
		return nil, nil
	})

	listeners := []Listener{makeListenerWithInode(443, 55)}
	out := e.Enrich(listeners)

	if out[0].Process == nil {
		t.Fatal("expected Process to be set")
	}
	if out[0].Process.PID != 123 {
		t.Errorf("got PID %d, want 123", out[0].Process.PID)
	}
}

func TestEnricher_LookupError_LeavesProcessNil(t *testing.T) {
	e := newEnricherWithLookup(func(inode uint64) (*ProcessInfo, error) {
		return nil, errors.New("permission denied")
	})

	listeners := []Listener{makeListenerWithInode(22, 77)}
	out := e.Enrich(listeners)

	if out[0].Process != nil {
		t.Error("expected Process to remain nil on lookup error")
	}
}

func TestEnricher_MultipleListeners(t *testing.T) {
	procs := map[uint64]*ProcessInfo{
		10: {PID: 1, Name: "sshd"},
		20: {PID: 2, Name: "nginx"},
	}
	e := newEnricherWithLookup(func(inode uint64) (*ProcessInfo, error) {
		return procs[inode], nil
	})

	listeners := []Listener{
		makeListenerWithInode(22, 10),
		makeListenerWithInode(80, 20),
		makeListenerWithInode(9000, 0),
	}
	out := e.Enrich(listeners)

	if out[0].Process == nil || out[0].Process.Name != "sshd" {
		t.Error("expected sshd for port 22")
	}
	if out[1].Process == nil || out[1].Process.Name != "nginx" {
		t.Error("expected nginx for port 80")
	}
	if out[2].Process != nil {
		t.Error("expected nil process for zero inode")
	}
}

func TestEnricher_EmptyListeners(t *testing.T) {
	called := false
	e := newEnricherWithLookup(func(inode uint64) (*ProcessInfo, error) {
		called = true
		return nil, nil
	})

	out := e.Enrich([]Listener{})

	if called {
		t.Error("lookup should not be called for empty listener list")
	}
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}
