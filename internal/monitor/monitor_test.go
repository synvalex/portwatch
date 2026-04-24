package monitor_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/ports"
	"github.com/user/portwatch/internal/rules"
)

// stubScanner is a Scanner that returns a fixed list of listeners.
type stubScanner struct {
	listeners []ports.Listener
	callCount int
}

func (s *stubScanner) Scan(_ context.Context) ([]ports.Listener, error) {
	s.callCount++
	return s.listeners, nil
}

func makeTestListener(port uint16) ports.Listener {
	return ports.Listener{
		Proto:   "tcp",
		Address: net.ParseIP("127.0.0.1"),
		Port:    port,
		PID:     1234,
	}
}

func TestMonitor_RunCancellation(t *testing.T) {
	scanner := &stubScanner{listeners: []ports.Listener{makeTestListener(8080)}}
	notifier := alert.NewLogNotifier(nil)
	dispatcher := alert.NewDispatcher([]rules.Rule{}, notifier, nil)

	m := monitor.New(scanner, dispatcher, 50*time.Millisecond, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	err := m.Run(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
	// Initial scan + at least one tick scan.
	if scanner.callCount < 2 {
		t.Errorf("expected at least 2 scans, got %d", scanner.callCount)
	}
}

func TestMonitor_ImmediateScan(t *testing.T) {
	scanner := &stubScanner{listeners: []ports.Listener{makeTestListener(9090)}}
	notifier := alert.NewLogNotifier(nil)
	dispatcher := alert.NewDispatcher([]rules.Rule{}, notifier, nil)

	m := monitor.New(scanner, dispatcher, 10*time.Second, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_ = m.Run(ctx)
	// Only the immediate scan should have fired within the short timeout.
	if scanner.callCount < 1 {
		t.Errorf("expected at least 1 scan, got %d", scanner.callCount)
	}
}
