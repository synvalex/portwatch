package baseline_test

import (
	"log/slog"
	"net"
	"os"
	"testing"
	"time"

	"github.com/dreadster3/portwatch/internal/baseline"
	"github.com/dreadster3/portwatch/internal/ports"
)

func makeFilterListener(port uint16) ports.Listener {
	return ports.Listener{
		IP:       net.ParseIP("0.0.0.0"),
		Port:     port,
		Protocol: "tcp",
	}
}

func silentLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Discard, nil))
}

func TestFilter_DuringLearning_PassesAll(t *testing.T) {
	store := baseline.NewStore()
	learner := baseline.NewLearner(store, 10*time.Minute)
	checker := baseline.NewChecker(store)
	filter := baseline.NewFilter(checker, learner, silentLogger())

	listeners := []ports.Listener{
		makeFilterListener(80),
		makeFilterListener(443),
	}

	result := filter.Apply(listeners)
	if len(result) != len(listeners) {
		t.Fatalf("expected %d listeners during learning, got %d", len(listeners), len(result))
	}
}

func TestFilter_DuringLearning_ObservesListeners(t *testing.T) {
	store := baseline.NewStore()
	learner := baseline.NewLearner(store, 10*time.Minute)
	checker := baseline.NewChecker(store)
	filter := baseline.NewFilter(checker, learner, silentLogger())

	listeners := []ports.Listener{makeFilterListener(8080)}
	filter.Apply(listeners)

	if !checker.IsKnown(listeners[0]) {
		t.Error("expected listener to be known after observation during learning")
	}
}

func TestFilter_AfterLearning_FiltersKnown(t *testing.T) {
	store := baseline.NewStore()
	learner := baseline.NewLearner(store, 0) // zero window — already expired
	checker := baseline.NewChecker(store)
	filter := baseline.NewFilter(checker, learner, silentLogger())

	known := makeFilterListener(80)
	store.Add(known)

	listeners := []ports.Listener{known, makeFilterListener(9999)}
	result := filter.Apply(listeners)

	if len(result) != 1 {
		t.Fatalf("expected 1 unknown listener, got %d", len(result))
	}
	if result[0].Port != 9999 {
		t.Errorf("expected port 9999, got %d", result[0].Port)
	}
}

func TestFilter_AfterLearning_AllKnown_ReturnsEmpty(t *testing.T) {
	store := baseline.NewStore()
	learner := baseline.NewLearner(store, 0)
	checker := baseline.NewChecker(store)
	filter := baseline.NewFilter(checker, learner, silentLogger())

	l := makeFilterListener(443)
	store.Add(l)

	result := filter.Apply([]ports.Listener{l})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d listeners", len(result))
	}
}
