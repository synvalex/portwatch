package config

import (
	"fmt"
	"time"
)

// BaselineConfig controls how portwatch learns and persists a baseline
// of expected listeners so that only deviations are alerted.
type BaselineConfig struct {
	Enabled       bool          `yaml:"enabled"`
	File          string        `yaml:"file"`
	LearnDuration time.Duration `yaml:"learn_duration"`
	AutoSave      bool          `yaml:"auto_save"`
}

// DefaultBaselineConfig returns a safe, disabled baseline configuration.
func DefaultBaselineConfig() BaselineConfig {
	return BaselineConfig{
		Enabled:       false,
		File:          "/var/lib/portwatch/baseline.json",
		LearnDuration: 5 * time.Minute,
		AutoSave:      true,
	}
}

// Validate checks that the baseline configuration is internally consistent.
func (b *BaselineConfig) Validate() error {
	if !b.Enabled {
		return nil
	}
	if b.File == "" {
		return fmt.Errorf("baseline.file must not be empty when baseline is enabled")
	}
	if b.LearnDuration < 0 {
		return fmt.Errorf("baseline.learn_duration must be non-negative, got %s", b.LearnDuration)
	}
	return nil
}

// Merge fills zero-value fields from defaults.
func (b *BaselineConfig) Merge(defaults BaselineConfig) {
	if b.File == "" {
		b.File = defaults.File
	}
	if b.LearnDuration == 0 {
		b.LearnDuration = defaults.LearnDuration
	}
}
