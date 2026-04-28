package history_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/ports"
)

func makeRecorderListener(port uint16) ports.Listener {
	return ports.Listener{
		Address: ports.Address{IP: "0.0.0.0", Port: port},
		Protocol: "tcp",
	}
}

func TestRecorder_NotifyAddsEvent(t *testing.T) {
	store := history.NewStore(100, time.Hour)
	rec := history.NewRecorder(store)

	ev := alert.Event{
		Type:     alert.EventAppeared,
		Listener: makeRecorderListener(8080),
	}

	if err := rec.Notify(ev); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	events := store.Recent(10)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Listener.Address.Port != 8080 {
		t.Errorf("expected port 8080, got %d", events[0].Listener.Address.Port)
	}
}

func TestRecorder_NotifyMultipleEvents(t *testing.T) {
	store := history.NewStore(100, time.Hour)
	rec := history.NewRecorder(store)

	for _, port := range []uint16{80, 443, 8080} {
		ev := alert.Event{
			Type:     alert.EventAppeared,
			Listener: makeRecorderListener(port),
		}
		if err := rec.Notify(ev); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	events := store.Recent(10)
	if len(events) != 3 {
		t.Fatalf("expected 3 events, got %d", len(events))
	}
}

func TestRecorder_NotifyDisappearedEvent(t *testing.T) {
	store := history.NewStore(100, time.Hour)
	rec := history.NewRecorder(store)

	ev := alert.Event{
		Type:     alert.EventDisappeared,
		Listener: makeRecorderListener(9090),
	}

	if err := rec.Notify(ev); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	events := store.Recent(10)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Type != alert.EventDisappeared {
		t.Errorf("expected EventDisappeared, got %v", events[0].Type)
	}
}
