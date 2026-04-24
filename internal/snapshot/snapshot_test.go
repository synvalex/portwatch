package snapshot_test

import (
	"testing"

	"github.com/user/portwatch/internal/ports"
	"github.com/user/portwatch/internal/snapshot"
)

func listener(proto, addr string, port uint16) ports.Listener {
	return ports.Listener{Proto: proto, Addr: addr, Port: port}
}

func TestStore_InitialUpdate_AllAppeared(t *testing.T) {
	s := snapshot.NewStore()
	listeners := []ports.Listener{
		listener("tcp", "0.0.0.0", 80),
		listener("tcp", "0.0.0.0", 443),
	}

	diff := s.Update(listeners)

	if len(diff.Appeared) != 2 {
		t.Errorf("expected 2 appeared, got %d", len(diff.Appeared))
	}
	if len(diff.Disappeared) != 0 {
		t.Errorf("expected 0 disappeared, got %d", len(diff.Disappeared))
	}
}

func TestStore_NoChange_EmptyDiff(t *testing.T) {
	s := snapshot.NewStore()
	listeners := []ports.Listener{listener("tcp", "0.0.0.0", 8080)}

	s.Update(listeners)
	diff := s.Update(listeners)

	if !diff.IsEmpty() {
		t.Errorf("expected empty diff, got %s", diff)
	}
}

func TestStore_NewListener_Appeared(t *testing.T) {
	s := snapshot.NewStore()
	s.Update([]ports.Listener{listener("tcp", "0.0.0.0", 80)})

	diff := s.Update([]ports.Listener{
		listener("tcp", "0.0.0.0", 80),
		listener("tcp", "0.0.0.0", 9090),
	})

	if len(diff.Appeared) != 1 || diff.Appeared[0].Port != 9090 {
		t.Errorf("expected port 9090 to appear, got %+v", diff.Appeared)
	}
	if len(diff.Disappeared) != 0 {
		t.Errorf("expected 0 disappeared, got %d", len(diff.Disappeared))
	}
}

func TestStore_RemovedListener_Disappeared(t *testing.T) {
	s := snapshot.NewStore()
	s.Update([]ports.Listener{
		listener("tcp", "0.0.0.0", 80),
		listener("tcp", "0.0.0.0", 3000),
	})

	diff := s.Update([]ports.Listener{listener("tcp", "0.0.0.0", 80)})

	if len(diff.Disappeared) != 1 || diff.Disappeared[0].Port != 3000 {
		t.Errorf("expected port 3000 to disappear, got %+v", diff.Disappeared)
	}
	if len(diff.Appeared) != 0 {
		t.Errorf("expected 0 appeared, got %d", len(diff.Appeared))
	}
}

func TestStore_Current_ReturnsCopy(t *testing.T) {
	s := snapshot.NewStore()
	s.Update([]ports.Listener{listener("udp", "127.0.0.1", 53)})

	current := s.Current()
	if len(current) != 1 {
		t.Errorf("expected 1 listener, got %d", len(current))
	}
}
