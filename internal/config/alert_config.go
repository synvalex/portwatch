package config

import (
	"fmt"
	"time"
)

// AlertConfig holds settings that control how alerts are emitted.
type AlertConfig struct {
	// Level is the minimum log level for alert messages (e.g. "info", "warn", "error").
	Level string `yaml:"level"`

	// DedupWindow is the duration during which identical alerts are suppressed.
	DedupWindow time.Duration `yaml:"dedup_window"`

	// EvictInterval controls how often stale deduplication entries are purged.
	EvictInterval time.Duration `yaml:"evict_interval"`
}

// DefaultAlertConfig returns an AlertConfig populated with sensible defaults.
func DefaultAlertConfig() AlertConfig {
	return AlertConfig{
		Level:         "warn",
		DedupWindow:   5 * time.Minute,
		EvictInterval: 10 * time.Minute,
	}
}

var validAlertLevels = map[string]bool{
	"debug": true,
	"info":  true,
	"warn":  true,
	"error": true,
}

// Validate returns an error if any AlertConfig field contains an invalid value.
func (a AlertConfig) Validate() error {
	if !validAlertLevels[a.Level] {
		return fmt.Errorf("alert.level %q is not valid; choose one of debug, info, warn, error", a.Level)
	}
	if a.DedupWindow < 0 {
		return fmt.Errorf("alert.dedup_window must not be negative")
	}
	if a.DedupWindow == 0 {
		return fmt.Errorf("alert.dedup_window must be greater than zero")
	}
	if a.EvictInterval <= 0 {
		return fmt.Errorf("alert.evict_interval must be greater than zero")
	}
	return nil
}

// Merge returns a new AlertConfig where zero values are filled from defaults.
func (a AlertConfig) Merge(defaults AlertConfig) AlertConfig {
	if a.Level == "" {
		a.Level = defaults.Level
	}
	if a.DedupWindow == 0 {
		a.DedupWindow = defaults.DedupWindow
	}
	if a.EvictInterval == 0 {
		a.EvictInterval = defaults.EvictInterval
	}
	return a
}
