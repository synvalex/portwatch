package baseline_test

import (
	"testing"
	"time"

	"github.com/example/portwatch/internal/baseline"
	"github.com/example/portwatch/internal/ports"
)

func makeLearnerListener(port uint16, proto string) ports.Listener {
	return ports.Listener{
		Address:  ports.Address{Port: port, Proto: proto},
		Inode:    0,
		Process:  nil,
	}
}

func TestLearner_LearnAndContains(t *testing.T) {
	store := baseline.NewStore()
	learner := baseline.NewLearner(store, 50*time.Millisecond)

	l := makeLearnerListener(8080, "tcp")
	learner.Observe(l)

	if !store.Contains(l) {
		t.Error("expected listener to be in baseline after Observe")
	}
}

func TestLearner_IsLearning_WithinWindow(t *testing.T) {
	store := baseline.NewStore()
	learner := baseline.NewLearner(store, 200*time.Millisecond)

	if !learner.IsLearning() {
		t.Error("expected learner to be in learning phase")
	}
}

func TestLearner_IsLearning_AfterWindow(t *testing.T) {
	store := baseline.NewStore()
	learner := baseline.NewLearner(store, 10*time.Millisecond)

	time.Sleep(20 * time.Millisecond)

	if learner.IsLearning() {
		t.Error("expected learner to have exited learning phase")
	}
}

func TestLearner_Observe_AfterWindow_DoesNotLearn(t *testing.T) {
	store := baseline.NewStore()
	learner := baseline.NewLearner(store, 10*time.Millisecond)

	time.Sleep(20 * time.Millisecond)

	l := makeLearnerListener(9090, "tcp")
	learner.Observe(l)

	if store.Contains(l) {
		t.Error("expected listener NOT to be added after learning window closed")
	}
}

func TestLearner_MultipleListeners(t *testing.T) {
	store := baseline.NewStore()
	learner := baseline.NewLearner(store, 100*time.Millisecond)

	listeners := []ports.Listener{
		makeLearnerListener(80, "tcp"),
		makeLearnerListener(443, "tcp"),
		makeLearnerListener(53, "udp"),
	}

	for _, l := range listeners {
		learner.Observe(l)
	}

	for _, l := range listeners {
		if !store.Contains(l) {
			t.Errorf("expected listener %v to be in baseline", l.Address)
		}
	}
}
