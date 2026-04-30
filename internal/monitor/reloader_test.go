package monitor

import (
	"context"
	"log/slog"
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func silentReloadLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

func TestReloader_Disabled_DoesNotBlock(t *testing.T) {
	cfg := config.ReloadConfig{Enabled: false}
	var called atomic.Bool
	r := NewReloader("", cfg, silentReloadLogger(), func(_ *config.Config) error {
		called.Store(true)
		return nil
	})
	done := make(chan struct{})
	go func() {
		r.Run(context.Background())
		close(done)
	}()
	select {
	case <-done:
		// ok — disabled path returns immediately
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Run did not return for disabled reloader")
	}
	if called.Load() {
		t.Error("reload callback should not have been called")
	}
}

func TestReloader_CancelStopsRun(t *testing.T) {
	cfg := config.DefaultReloadConfig()
	r := NewReloader("", cfg, silentReloadLogger(), func(_ *config.Config) error { return nil })
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		r.Run(ctx)
		close(done)
	}()
	cancel()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Run did not stop after context cancellation")
	}
}

func TestReloader_SIGHUP_TriggersReload(t *testing.T) {
	// Write a minimal valid config to a temp file.
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString("interval: 5s\n")
	_ = f.Close()

	cfg := config.ReloadConfig{Enabled: true, Debounce: 10 * time.Millisecond, OnReloadFail: "warn"}
	var calls atomic.Int32
	r := NewReloader(f.Name(), cfg, silentReloadLogger(), func(_ *config.Config) error {
		calls.Add(1)
		return nil
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.Run(ctx)

	// Send SIGHUP to self.
	_ = syscall.Kill(os.Getpid(), syscall.SIGHUP)
	time.Sleep(150 * time.Millisecond)
	if calls.Load() == 0 {
		t.Error("expected reload callback to be called after SIGHUP")
	}
}
