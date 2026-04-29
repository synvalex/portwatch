package retention

import (
	"context"
	"log/slog"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

type mockStore struct {
	calls atomic.Int32
	ret   int
}

func (m *mockStore) Prune(_ time.Duration) int {
	m.calls.Add(1)
	return m.ret
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

func TestPruner_Register(t *testing.T) {
	p := NewPruner(10*time.Millisecond, time.Hour, testLogger())
	s := &mockStore{}
	p.Register(s)
	if len(p.stores) != 1 {
		t.Fatalf("expected 1 store, got %d", len(p.stores))
	}
}

func TestPruner_Run_CallsPrune(t *testing.T) {
	p := NewPruner(20*time.Millisecond, time.Hour, testLogger())
	s := &mockStore{ret: 0}
	p.Register(s)

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		p.Run(ctx)
		close(done)
	}()
	<-done

	if s.calls.Load() < 2 {
		t.Errorf("expected at least 2 prune calls, got %d", s.calls.Load())
	}
}

func TestPruner_Run_StopsOnCancel(t *testing.T) {
	p := NewPruner(50*time.Millisecond, time.Hour, testLogger())
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan struct{})
	go func() {
		p.Run(ctx)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Error("pruner did not stop after context cancellation")
	}
}

func TestPruner_MultipleStores(t *testing.T) {
	p := NewPruner(20*time.Millisecond, time.Hour, testLogger())
	s1 := &mockStore{ret: 1}
	s2 := &mockStore{ret: 2}
	p.Register(s1)
	p.Register(s2)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	p.Run(ctx)

	if s1.calls.Load() == 0 {
		t.Error("expected s1.Prune to be called")
	}
	if s2.calls.Load() == 0 {
		t.Error("expected s2.Prune to be called")
	}
}
