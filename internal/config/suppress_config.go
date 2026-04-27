package config

import (
	"fmt"
	"time"
)

// SuppressConfig controls the alert suppression filter which stops
// repeated notifications for the same listener once a threshold is crossed.
type SuppressConfig struct {
	// Window is the rolling time window within which occurrences are counted.
	Window time.Duration `yaml:"window"`
	// Max is the maximum number of alerts allowed per key per window.
	// Set to 0 to disable suppression entirely.
	Max int `yaml:"max"`
}

// DefaultSuppressConfig returns conservative defaults: allow up to 5 alerts
// per listener per 10-minute window.
func DefaultSuppressConfig() SuppressConfig {
	return SuppressConfig{
		Window: 10 * time.Minute,
		Max:    5,
	}
}

// Validate checks SuppressConfig for logical consistency.
func (s SuppressConfig) Validate() error {
	if s.Max < 0 {
		return fmt.Errorf("suppress.max must be >= 0, got %d", s.Max)
	}
	if s.Max > 0 && s.Window <= 0 {
		return fmt.Errorf("suppress.window must be positive when max > 0")
	}
	if s.Window < 0 {
		return fmt.Errorf("suppress.window must not be negative")
	}
	return nil
}

// Merge fills zero-value fields from defaults.
func (s SuppressConfig) Merge(defaults SuppressConfig) SuppressConfig {
	if s.Window == 0 {
		s.Window = defaults.Window
	}
	if s.Max == 0 {
		s.Max = defaults.Max
	}
	return s
}
