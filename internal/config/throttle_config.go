package config

import (
	"fmt"
	"time"
)

// ThrottleConfig controls per-port alert throttling.
type ThrottleConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Rate     int           `yaml:"rate"`      // max alerts per window
	Window   time.Duration `yaml:"window"`    // rolling window duration
	PerPort  bool          `yaml:"per_port"`  // throttle per port vs globally
}

// DefaultThrottleConfig returns sensible defaults.
func DefaultThrottleConfig() ThrottleConfig {
	return ThrottleConfig{
		Enabled: false,
		Rate:    10,
		Window:  1 * time.Minute,
		PerPort: true,
	}
}

// Validate checks that ThrottleConfig fields are sensible.
func (c ThrottleConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Rate <= 0 {
		return fmt.Errorf("throttle rate must be > 0, got %d", c.Rate)
	}
	if c.Window <= 0 {
		return fmt.Errorf("throttle window must be > 0, got %s", c.Window)
	}
	return nil
}

// Merge fills zero-value fields from defaults.
func (c ThrottleConfig) Merge(defaults ThrottleConfig) ThrottleConfig {
	if c.Rate == 0 {
		c.Rate = defaults.Rate
	}
	if c.Window == 0 {
		c.Window = defaults.Window
	}
	return c
}
