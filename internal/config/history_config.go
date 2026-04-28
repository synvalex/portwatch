package config

import (
	"fmt"
	"time"
)

// HistoryConfig controls the event history ring-buffer.
type HistoryConfig struct {
	Enabled   bool          `yaml:"enabled"`
	MaxEvents int           `yaml:"max_events"`
	Retention time.Duration `yaml:"retention"`
}

// DefaultHistoryConfig returns sensible defaults.
func DefaultHistoryConfig() HistoryConfig {
	return HistoryConfig{
		Enabled:   true,
		MaxEvents: 500,
		Retention: 24 * time.Hour,
	}
}

// Validate checks that the config values are within acceptable bounds.
func (h HistoryConfig) Validate() error {
	if !h.Enabled {
		return nil
	}
	if h.MaxEvents <= 0 {
		return fmt.Errorf("history.max_events must be > 0, got %d", h.MaxEvents)
	}
	if h.Retention <= 0 {
		return fmt.Errorf("history.retention must be > 0, got %s", h.Retention)
	}
	return nil
}

// Merge fills zero-value fields from defaults.
func (h HistoryConfig) Merge(defaults HistoryConfig) HistoryConfig {
	if h.MaxEvents == 0 {
		h.MaxEvents = defaults.MaxEvents
	}
	if h.Retention == 0 {
		h.Retention = defaults.Retention
	}
	return h
}
