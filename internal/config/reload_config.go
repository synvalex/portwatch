package config

import (
	"fmt"
	"time"
)

// ReloadConfig controls hot-reload / SIGHUP behaviour.
type ReloadConfig struct {
	Enabled      bool          `yaml:"enabled"`
	Debounce     time.Duration `yaml:"debounce"`
	OnReloadFail string        `yaml:"on_reload_fail"` // "warn" | "fatal"
}

// DefaultReloadConfig returns sensible defaults.
func DefaultReloadConfig() ReloadConfig {
	return ReloadConfig{
		Enabled:      true,
		Debounce:     500 * time.Millisecond,
		OnReloadFail: "warn",
	}
}

func (r *ReloadConfig) Validate() error {
	if !r.Enabled {
		return nil
	}
	if r.Debounce < 0 {
		return fmt.Errorf("reload.debounce must be non-negative, got %s", r.Debounce)
	}
	if r.Debounce > 10*time.Second {
		return fmt.Errorf("reload.debounce too large (max 10s), got %s", r.Debounce)
	}
	switch r.OnReloadFail {
	case "warn", "fatal":
		// valid
	default:
		return fmt.Errorf("reload.on_reload_fail must be \"warn\" or \"fatal\", got %q", r.OnReloadFail)
	}
	return nil
}

func (r *ReloadConfig) Merge(other ReloadConfig) {
	if other.Debounce != 0 {
		r.Debounce = other.Debounce
	}
	if other.OnReloadFail != "" {
		r.OnReloadFail = other.OnReloadFail
	}
	r.Enabled = other.Enabled
}
