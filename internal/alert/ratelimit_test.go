package alert

import (
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func rlListener(port uint16) ports.Listener {
	return ports.Listener{
		Proto:   "tcp",
		Address: net.TCPAddr{IP: net.ParseIP("0.0.0.0"), Port: int(port)},
	}
}

func TestRateLimiter_ZeroBurst_AllowsAll(t *testing.T) {
	rl := NewRateLimiter(0, time.Second)
	l := rlListener(8080)
	for i := 0; i < 20; i++ {
		if !rl.Allow(l, "appeared") {
			t.Fatalf("expected allow on iteration %d with zero burst", i)
		}
	}
}

func TestRateLimiter_WithinBurst_Allowed(t *testing.T) {
	rl := NewRateLimiter(3, time.Minute)
	l := rlListener(9090)
	for i := 0; i < 3; i++ {
		if !rl.Allow(l, "appeared") {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
}

func TestRateLimiter_ExceedsBurst_Suppressed(t *testing.T) {
	rl := NewRateLimiter(3, time.Minute)
	l := rlListener(9090)
	for i := 0; i < 3; i++ {
		rl.Allow(l, "appeared")
	}
	if rl.Allow(l, "appeared") {
		t.Fatal("expected suppression after burst exceeded")
	}
}

func TestRateLimiter_WindowExpiry_AllowsAgain(t *testing.T) {
	window := 50 * time.Millisecond
	rl := NewRateLimiter(2, window)
	l := rlListener(7070)

	rl.Allow(l, "appeared")
	rl.Allow(l, "appeared")

	// Advance fake clock past the window.
	base := time.Now().Add(window + time.Millisecond)
	rl.now = func() time.Time { return base }

	if !rl.Allow(l, "appeared") {
		t.Fatal("expected allow after window expired")
	}
}

func TestRateLimiter_DifferentEventTypes_IndependentBuckets(t *testing.T) {
	rl := NewRateLimiter(1, time.Minute)
	l := rlListener(3000)

	if !rl.Allow(l, "appeared") {
		t.Fatal("first appeared should be allowed")
	}
	if !rl.Allow(l, "disappeared") {
		t.Fatal("first disappeared should be allowed (different bucket)")
	}
	if rl.Allow(l, "appeared") {
		t.Fatal("second appeared should be suppressed")
	}
}

func TestRateLimiter_Reset_ClearsState(t *testing.T) {
	rl := NewRateLimiter(1, time.Minute)
	l := rlListener(4000)

	rl.Allow(l, "appeared")
	if rl.Allow(l, "appeared") {
		t.Fatal("expected suppression before reset")
	}

	rl.Reset()

	if !rl.Allow(l, "appeared") {
		t.Fatal("expected allow after reset")
	}
}
