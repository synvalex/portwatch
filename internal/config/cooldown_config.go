package config

import (
	"fmt"
	"time"
)

// CooldownConfig controls per-listener alert cooldown behaviour.
type CooldownConfig struct {
	// Enabled toggles the cooldown filter entirely.
	Enabled bool `yaml:"enabled"`

	// Window is the minimum duration between repeated alerts for the same
	// listener+event pair.
	Window time.Duration `yaml:"window"`

	// PurgeInterval controls how often expired cooldown entries are purged
	// from memory.
	PurgeInterval time.Duration `yaml:"purge_interval"`
}

// DefaultCooldownConfig returns sensible defaults for the cooldown filter.
func DefaultCooldownConfig() CooldownConfig {
	return CooldownConfig{
		Enabled:       true,
		Window:        2 * time.Minute,
		PurgeInterval: 10 * time.Minute,
	}
}

// Validate checks that all fields are within acceptable bounds.
func (c CooldownConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Window < 0 {
		return fmt.Errorf("cooldown.window must not be negative, got %s", c.Window)
	}
	if c.Window > 24*time.Hour {
		return fmt.Errorf("cooldown.window too large (max 24h), got %s", c.Window)
	}
	if c.PurgeInterval <= 0 {
		return fmt.Errorf("cooldown.purge_interval must be positive, got %s", c.PurgeInterval)
	}
	return nil
}

// Merge fills zero-value fields in c from defaults.
func (c CooldownConfig) Merge(defaults CooldownConfig) CooldownConfig {
	if c.Window == 0 {
		c.Window = defaults.Window
	}
	if c.PurgeInterval == 0 {
		c.PurgeInterval = defaults.PurgeInterval
	}
	return c
}
