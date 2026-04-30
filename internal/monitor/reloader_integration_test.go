package monitor_test

import (
	"context"
	"log/slog"
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
)

func TestReloader_DebounceCoalesces(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString("interval: 5s\n")
	_ = f.Close()

	cfg := config.ReloadConfig{
		Enabled:      true,
		Debounce:     80 * time.Millisecond,
		OnReloadFail: "warn",
	}
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))

	var calls atomic.Int32
	r := monitor.NewReloader(f.Name(), cfg, log, func(_ *config.Config) error {
		calls.Add(1)
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go r.Run(ctx)

	// Fire three rapid SIGHUPs — debounce should coalesce them into one reload.
	for i := 0; i < 3; i++ {
		_ = syscall.Kill(os.Getpid(), syscall.SIGHUP)
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(300 * time.Millisecond)

	if n := calls.Load(); n != 1 {
		t.Errorf("expected 1 coalesced reload, got %d", n)
	}
}
