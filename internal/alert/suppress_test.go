package alert

import (
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func suppressListener(port uint16) ports.Listener {
	return ports.Listener{
		Address:  ports.Address{IP: net.ParseIP("0.0.0.0"), Port: port},
		Protocol: "tcp",
	}
}

func TestSuppressFilter_ZeroWindow_AllowsAll(t *testing.T) {
	sf := NewSuppressFilter(0, 3)
	e := Event{Listener: suppressListener(80), Type: EventAppeared}
	for i := 0; i < 10; i++ {
		if !sf.Allow(e) {
			t.Fatalf("expected allow on iteration %d with zero window", i)
		}
	}
}

func TestSuppressFilter_WithinMax_Allowed(t *testing.T) {
	sf := NewSuppressFilter(time.Minute, 3)
	e := Event{Listener: suppressListener(8080), Type: EventAppeared}
	for i := 1; i <= 3; i++ {
		if !sf.Allow(e) {
			t.Fatalf("expected allow on call %d (max=3)", i)
		}
	}
}

func TestSuppressFilter_ExceedsMax_Suppressed(t *testing.T) {
	sf := NewSuppressFilter(time.Minute, 2)
	e := Event{Listener: suppressListener(443), Type: EventAppeared}
	sf.Allow(e)
	sf.Allow(e)
	if sf.Allow(e) {
		t.Fatal("expected suppression on 3rd call (max=2)")
	}
}

func TestSuppressFilter_WindowExpiry_Resets(t *testing.T) {
	now := time.Now()
	sf := NewSuppressFilter(50*time.Millisecond, 1)
	sf.nowFunc = func() time.Time { return now }

	e := Event{Listener: suppressListener(9000), Type: EventAppeared}
	if !sf.Allow(e) {
		t.Fatal("first call should be allowed")
	}
	if sf.Allow(e) {
		t.Fatal("second call within window should be suppressed")
	}

	// advance past window
	sf.nowFunc = func() time.Time { return now.Add(100 * time.Millisecond) }
	if !sf.Allow(e) {
		t.Fatal("call after window expiry should be allowed")
	}
}

func TestSuppressFilter_DifferentPorts_Independent(t *testing.T) {
	sf := NewSuppressFilter(time.Minute, 1)
	e1 := Event{Listener: suppressListener(80), Type: EventAppeared}
	e2 := Event{Listener: suppressListener(443), Type: EventAppeared}

	sf.Allow(e1)
	if !sf.Allow(e2) {
		t.Fatal("different port should not be affected by suppression of port 80")
	}
}
