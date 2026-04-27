package baseline_test

import (
	"testing"

	"github.com/example/portwatch/internal/baseline"
	"github.com/example/portwatch/internal/ports"
)

func makeCheckerListener(port uint16, proto string) ports.Listener {
	return ports.Listener{
		Address: ports.Address{Port: port, Proto: proto},
	}
}

func populatedStore(listeners []ports.Listener) *baseline.Store {
	s := baseline.NewStore()
	for _, l := range listeners {
		s.Add(l)
	}
	return s
}

func TestChecker_IsKnown_True(t *testing.T) {
	l := makeCheckerListener(80, "tcp")
	checker := baseline.NewChecker(populatedStore([]ports.Listener{l}))

	if !checker.IsKnown(l) {
		t.Error("expected listener to be known")
	}
}

func TestChecker_IsKnown_False(t *testing.T) {
	checker := baseline.NewChecker(baseline.NewStore())
	l := makeCheckerListener(9999, "tcp")

	if checker.IsKnown(l) {
		t.Error("expected listener to be unknown")
	}
}

func TestChecker_FilterUnknown(t *testing.T) {
	known := makeCheckerListener(80, "tcp")
	unknown := makeCheckerListener(8888, "tcp")

	checker := baseline.NewChecker(populatedStore([]ports.Listener{known}))
	result := checker.FilterUnknown([]ports.Listener{known, unknown})

	if len(result) != 1 || result[0].Address.Port != 8888 {
		t.Errorf("expected only unknown listener, got %v", result)
	}
}

func TestChecker_FilterKnown(t *testing.T) {
	known := makeCheckerListener(443, "tcp")
	unknown := makeCheckerListener(1234, "udp")

	checker := baseline.NewChecker(populatedStore([]ports.Listener{known}))
	result := checker.FilterKnown([]ports.Listener{known, unknown})

	if len(result) != 1 || result[0].Address.Port != 443 {
		t.Errorf("expected only known listener, got %v", result)
	}
}

func TestChecker_FilterUnknown_EmptyBaseline(t *testing.T) {
	checker := baseline.NewChecker(baseline.NewStore())
	listeners := []ports.Listener{
		makeCheckerListener(80, "tcp"),
		makeCheckerListener(443, "tcp"),
	}

	result := checker.FilterUnknown(listeners)
	if len(result) != 2 {
		t.Errorf("expected all listeners to be unknown, got %d", len(result))
	}
}
