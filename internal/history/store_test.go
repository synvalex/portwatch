package history

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/ports"
)

func makeEvent(port uint16) alert.Event {
	return alert.Event{
		Type: alert.EventAppeared,
		Listener: ports.Listener{
			Proto:   "tcp",
			Address: ports.Address{Port: port},
		},
	}
}

func TestStore_Add_IncreasesLen(t *testing.T) {
	s := NewStore(10, time.Hour)
	s.Add(makeEvent(8080))
	s.Add(makeEvent(9090))
	if s.Len() != 2 {
		t.Errorf("expected 2 entries, got %d", s.Len())
	}
}

func TestStore_Recent_ReturnsWithinRetention(t *testing.T) {
	s := NewStore(10, time.Hour)
	now := time.Now()
	s.now = func() time.Time { return now }
	s.Add(makeEvent(8080))
	if len(s.Recent()) != 1 {
		t.Errorf("expected 1 recent entry")
	}
}

func TestStore_Recent_ExcludesExpired(t *testing.T) {
	s := NewStore(10, time.Minute)
	base := time.Now()
	s.now = func() time.Time { return base }
	s.Add(makeEvent(8080))
	// advance time past retention
	s.now = func() time.Time { return base.Add(2 * time.Minute) }
	if len(s.Recent()) != 0 {
		t.Errorf("expected 0 recent entries after retention expired")
	}
}

func TestStore_Evict_Capacity(t *testing.T) {
	s := NewStore(3, time.Hour)
	for i := 0; i < 5; i++ {
		s.Add(makeEvent(uint16(8000 + i)))
	}
	if s.Len() > 3 {
		t.Errorf("expected at most 3 entries, got %d", s.Len())
	}
}

func TestStore_Evict_RetainsNewest(t *testing.T) {
	s := NewStore(2, time.Hour)
	s.Add(makeEvent(1111))
	s.Add(makeEvent(2222))
	s.Add(makeEvent(3333))
	entries := s.Recent()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Event.Listener.Address.Port != 2222 {
		t.Errorf("expected port 2222, got %d", entries[0].Event.Listener.Address.Port)
	}
	if entries[1].Event.Listener.Address.Port != 3333 {
		t.Errorf("expected port 3333, got %d", entries[1].Event.Listener.Address.Port)
	}
}
