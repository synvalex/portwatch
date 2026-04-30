package monitor

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/portwatch/internal/config"
)

// ReloadFunc is called with the freshly loaded config on each SIGHUP.
type ReloadFunc func(cfg *config.Config) error

// Reloader listens for SIGHUP and triggers a config reload with debounce.
type Reloader struct {
	configPath string
	cfg        config.ReloadConfig
	log        *slog.Logger
	onReload   ReloadFunc
}

// NewReloader creates a Reloader that watches configPath.
func NewReloader(configPath string, cfg config.ReloadConfig, log *slog.Logger, fn ReloadFunc) *Reloader {
	return &Reloader{
		configPath: configPath,
		cfg:        cfg,
		log:        log,
		onReload:   fn,
	}
}

// Run blocks until ctx is cancelled, reloading config on SIGHUP.
func (r *Reloader) Run(ctx context.Context) {
	if !r.cfg.Enabled {
		r.log.Info("config reload disabled")
		return
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP)
	defer signal.Stop(sigCh)

	var debounce *time.Timer
	for {
		select {
		case <-ctx.Done():
			if debounce != nil {
				debounce.Stop()
			}
			return
		case <-sigCh:
			if debounce != nil {
				debounce.Stop()
			}
			debounce = time.AfterFunc(r.cfg.Debounce, func() {
				r.reload()
			})
		}
	}
}

func (r *Reloader) reload() {
	r.log.Info("reloading config", "path", r.configPath)
	cfg, err := config.Load(r.configPath)
	if err != nil {
		r.handleError("load", err)
		return
	}
	if err := r.onReload(cfg); err != nil {
		r.handleError("apply", err)
		return
	}
	r.log.Info("config reloaded successfully")
}

func (r *Reloader) handleError(stage string, err error) {
	if r.cfg.OnReloadFail == "fatal" {
		r.log.Error("fatal config reload error", "stage", stage, "err", err)
		os.Exit(1)
	}
	r.log.Warn("config reload error", "stage", stage, "err", err)
}
