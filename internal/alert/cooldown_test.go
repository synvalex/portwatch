package alert

import (
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func cooldownListener(port uint16) ports.Listener {
	return ports.Listener{
		IP:       net.ParseIP("0.0.0.0"),
		Port:     port,
		Protocol: "tcp",
	}
}

func TestCooldownFilter_ZeroWindow_AlwaysAllows(t *testing.T) {
	cf := NewCooldownFilter(0)
	l := cooldownListener(8080)
	for i := 0; i < 5; i++ {
		if !cf.Allow(l, "appeared") {
			t.Fatalf("expected allow on iteration %d with zero window", i)
		}
	}
}

func TestCooldownFilter_FirstEvent_Allowed(t *testing.T) {
	cf := NewCooldownFilter(5 * time.Second)
	l := cooldownListener(9090)
	if !cf.Allow(l, "appeared") {
		t.Fatal("expected first event to be allowed")
	}
}

func TestCooldownFilter_SecondEventWithinWindow_Suppressed(t *testing.T) {
	now := time.Now()
	cf := NewCooldownFilter(10 * time.Second)
	cf.now = func() time.Time { return now }

	l := cooldownListener(3000)
	cf.Allow(l, "appeared") // prime

	cf.now = func() time.Time { return now.Add(5 * time.Second) }
	if cf.Allow(l, "appeared") {
		t.Fatal("expected suppression within cooldown window")
	}
}

func TestCooldownFilter_AfterWindowExpiry_Allowed(t *testing.T) {
	now := time.Now()
	cf := NewCooldownFilter(10 * time.Second)
	cf.now = func() time.Time { return now }

	l := cooldownListener(3000)
	cf.Allow(l, "appeared")

	cf.now = func() time.Time { return now.Add(11 * time.Second) }
	if !cf.Allow(l, "appeared") {
		t.Fatal("expected allow after cooldown window expired")
	}
}

func TestCooldownFilter_DifferentEventTypes_Independent(t *testing.T) {
	now := time.Now()
	cf := NewCooldownFilter(10 * time.Second)
	cf.now = func() time.Time { return now }

	l := cooldownListener(4000)
	cf.Allow(l, "appeared")

	// different event type should be allowed independently
	if !cf.Allow(l, "disappeared") {
		t.Fatal("expected different event type to be allowed")
	}
}

func TestCooldownFilter_Purge_RemovesExpiredEntries(t *testing.T) {
	now := time.Now()
	cf := NewCooldownFilter(5 * time.Second)
	cf.now = func() time.Time { return now }

	l := cooldownListener(5000)
	cf.Allow(l, "appeared")

	cf.now = func() time.Time { return now.Add(10 * time.Second) }
	cf.Purge()

	if len(cf.lastSeen) != 0 {
		t.Fatalf("expected lastSeen to be empty after purge, got %d entries", len(cf.lastSeen))
	}
}
