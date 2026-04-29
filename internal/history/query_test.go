package history_test

import (
	"net"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/history"
	"github.com/yourorg/portwatch/internal/ports"
)

func makeQueryListener(port uint16) ports.Listener {
	return ports.Listener{
		IP:       net.ParseIP("127.0.0.1"),
		Port:     port,
		Protocol: "tcp",
	}
}

func populatedQuery(t *testing.T) (*history.Store, *history.Query) {
	t.Helper()
	s := history.NewStore(100, time.Hour)
	s.Add(history.Event{Listener: makeQueryListener(80), EventType: history.EventAppeared, OccurredAt: time.Now()})
	s.Add(history.Event{Listener: makeQueryListener(443), EventType: history.EventAppeared, OccurredAt: time.Now()})
	s.Add(history.Event{Listener: makeQueryListener(80), EventType: history.EventDisappeared, OccurredAt: time.Now()})
	return s, history.NewQuery(s)
}

func TestQuery_Recent_ReturnsAll(t *testing.T) {
	_, q := populatedQuery(t)
	events := q.Recent(time.Hour)
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
}

func TestQuery_Since_FutureTime_Empty(t *testing.T) {
	_, q := populatedQuery(t)
	events := q.Since(time.Now().Add(time.Hour))
	if len(events) != 0 {
		t.Fatalf("expected 0 events for future time, got %d", len(events))
	}
}

func TestQuery_ByPort_Filters(t *testing.T) {
	_, q := populatedQuery(t)
	events := q.ByPort(80)
	if len(events) != 2 {
		t.Fatalf("expected 2 events for port 80, got %d", len(events))
	}
	for _, e := range events {
		if e.Listener.Port != 80 {
			t.Errorf("unexpected port %d in result", e.Listener.Port)
		}
	}
}

func TestQuery_ByEventType_Appeared(t *testing.T) {
	_, q := populatedQuery(t)
	events := q.ByEventType(history.EventAppeared)
	if len(events) != 2 {
		t.Fatalf("expected 2 appeared events, got %d", len(events))
	}
}

func TestQuery_ByEventType_Disappeared(t *testing.T) {
	_, q := populatedQuery(t)
	events := q.ByEventType(history.EventDisappeared)
	if len(events) != 1 {
		t.Fatalf("expected 1 disappeared event, got %d", len(events))
	}
}

func TestQuery_ByPort_NoMatch(t *testing.T) {
	_, q := populatedQuery(t)
	events := q.ByPort(9999)
	if len(events) != 0 {
		t.Fatalf("expected 0 events for unknown port, got %d", len(events))
	}
}
