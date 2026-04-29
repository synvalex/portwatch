package alert

import (
	"net"
	"testing"
	"time"

	"github.com/wesen/portwatch/internal/ports"
)

func throttleListener(port uint16) ports.Listener {
	return ports.Listener{
		Proto: "tcp",
		Port:  port,
		Addr:  net.ParseIP("0.0.0.0"),
	}
}

func TestThrottleFilter_ZeroRate_AllowsAll(t *testing.T) {
	f := NewThrottleFilter(0, time.Minute, true)
	l := throttleListener(8080)
	for i := 0; i < 20; i++ {
		if !f.Allow(l, "appeared") {
			t.Errorf("expected allow on iteration %d with zero rate", i)
		}
	}
}

func TestThrottleFilter_WithinRate_Allowed(t *testing.T) {
	now := time.Now()
	f := NewThrottleFilter(3, time.Minute, true)
	f.now = func() time.Time { return now }
	l := throttleListener(9000)
	for i := 0; i < 3; i++ {
		if !f.Allow(l, "appeared") {
			t.Errorf("expected allow on call %d", i+1)
		}
	}
}

func TestThrottleFilter_ExceedsRate_Suppressed(t *testing.T) {
	now := time.Now()
	f := NewThrottleFilter(3, time.Minute, true)
	f.now = func() time.Time { return now }
	l := throttleListener(9000)
	for i := 0; i < 3; i++ {
		f.Allow(l, "appeared")
	}
	if f.Allow(l, "appeared") {
		t.Error("expected suppression after exceeding rate")
	}
}

func TestThrottleFilter_WindowExpiry_Resets(t *testing.T) {
	now := time.Now()
	f := NewThrottleFilter(2, 30*time.Second, true)
	f.now = func() time.Time { return now }
	l := throttleListener(7070)
	f.Allow(l, "appeared")
	f.Allow(l, "appeared")
	if f.Allow(l, "appeared") {
		t.Error("expected suppression within window")
	}
	// advance past window
	f.now = func() time.Time { return now.Add(31 * time.Second) }
	if !f.Allow(l, "appeared") {
		t.Error("expected allow after window expiry")
	}
}

func TestThrottleFilter_PerPort_IndependentBuckets(t *testing.T) {
	now := time.Now()
	f := NewThrottleFilter(1, time.Minute, true)
	f.now = func() time.Time { return now }
	l1 := throttleListener(80)
	l2 := throttleListener(443)
	if !f.Allow(l1, "appeared") {
		t.Error("expected allow for port 80")
	}
	if !f.Allow(l2, "appeared") {
		t.Error("expected allow for port 443 (different bucket)")
	}
	if f.Allow(l1, "appeared") {
		t.Error("expected suppression for port 80 second call")
	}
}

func TestThrottleFilter_GlobalKey_SharedBucket(t *testing.T) {
	now := time.Now()
	f := NewThrottleFilter(2, time.Minute, false)
	f.now = func() time.Time { return now }
	l1 := throttleListener(80)
	l2 := throttleListener(443)
	f.Allow(l1, "appeared")
	f.Allow(l2, "appeared")
	if f.Allow(l1, "appeared") {
		t.Error("expected suppression: global bucket exhausted")
	}
}
