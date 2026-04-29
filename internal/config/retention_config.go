package config

import (
	"fmt"
	"time"
)

// RetentionConfig controls how long scan snapshots and events are kept in memory.
type RetentionConfig struct {
	Enabled       bool          `yaml:"enabled"`
	MaxSnapshots  int           `yaml:"max_snapshots"`
	MaxAge        time.Duration `yaml:"max_age"`
	PruneInterval time.Duration `yaml:"prune_interval"`
}

// DefaultRetentionConfig returns a RetentionConfig with sensible defaults.
func DefaultRetentionConfig() RetentionConfig {
	return RetentionConfig{
		Enabled:       true,
		MaxSnapshots:  100,
		MaxAge:        24 * time.Hour,
		PruneInterval: 5 * time.Minute,
	}
}

// Validate checks that all fields are within acceptable ranges.
func (c RetentionConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.MaxSnapshots <= 0 {
		return fmt.Errorf("retention: max_snapshots must be > 0, got %d", c.MaxSnapshots)
	}
	if c.MaxAge <= 0 {
		return fmt.Errorf("retention: max_age must be > 0, got %s", c.MaxAge)
	}
	if c.PruneInterval <= 0 {
		return fmt.Errorf("retention: prune_interval must be > 0, got %s", c.PruneInterval)
	}
	if c.PruneInterval > c.MaxAge {
		return fmt.Errorf("retention: prune_interval (%s) must not exceed max_age (%s)", c.PruneInterval, c.MaxAge)
	}
	return nil
}

// Merge fills zero-value fields from defaults.
func (c RetentionConfig) Merge(defaults RetentionConfig) RetentionConfig {
	if c.MaxSnapshots == 0 {
		c.MaxSnapshots = defaults.MaxSnapshots
	}
	if c.MaxAge == 0 {
		c.MaxAge = defaults.MaxAge
	}
	if c.PruneInterval == 0 {
		c.PruneInterval = defaults.PruneInterval
	}
	return c
}
