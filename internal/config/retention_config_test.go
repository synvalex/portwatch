package config

import (
	"testing"
	"time"
)

func TestDefaultRetentionConfig(t *testing.T) {
	c := DefaultRetentionConfig()
	if !c.Enabled {
		t.Error("expected Enabled to be true")
	}
	if c.MaxSnapshots <= 0 {
		t.Errorf("expected MaxSnapshots > 0, got %d", c.MaxSnapshots)
	}
	if c.MaxAge <= 0 {
		t.Errorf("expected MaxAge > 0, got %s", c.MaxAge)
	}
	if c.PruneInterval <= 0 {
		t.Errorf("expected PruneInterval > 0, got %s", c.PruneInterval)
	}
}

func TestRetentionConfig_Validate_Disabled(t *testing.T) {
	c := RetentionConfig{Enabled: false}
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got %v", err)
	}
}

func TestRetentionConfig_Validate_Valid(t *testing.T) {
	c := DefaultRetentionConfig()
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestRetentionConfig_Validate_ZeroMaxSnapshots(t *testing.T) {
	c := DefaultRetentionConfig()
	c.MaxSnapshots = 0
	if err := c.Validate(); err == nil {
		t.Error("expected error for zero max_snapshots")
	}
}

func TestRetentionConfig_Validate_ZeroMaxAge(t *testing.T) {
	c := DefaultRetentionConfig()
	c.MaxAge = 0
	if err := c.Validate(); err == nil {
		t.Error("expected error for zero max_age")
	}
}

func TestRetentionConfig_Validate_PruneIntervalExceedsMaxAge(t *testing.T) {
	c := DefaultRetentionConfig()
	c.MaxAge = 1 * time.Minute
	c.PruneInterval = 10 * time.Minute
	if err := c.Validate(); err == nil {
		t.Error("expected error when prune_interval > max_age")
	}
}

func TestRetentionConfig_Merge_FillsZeroValues(t *testing.T) {
	defaults := DefaultRetentionConfig()
	user := RetentionConfig{Enabled: true}
	merged := user.Merge(defaults)
	if merged.MaxSnapshots != defaults.MaxSnapshots {
		t.Errorf("expected MaxSnapshots %d, got %d", defaults.MaxSnapshots, merged.MaxSnapshots)
	}
	if merged.MaxAge != defaults.MaxAge {
		t.Errorf("expected MaxAge %s, got %s", defaults.MaxAge, merged.MaxAge)
	}
}

func TestRetentionConfig_Merge_UserValuesPreserved(t *testing.T) {
	defaults := DefaultRetentionConfig()
	user := RetentionConfig{
		Enabled:       true,
		MaxSnapshots:  50,
		MaxAge:        2 * time.Hour,
		PruneInterval: 1 * time.Minute,
	}
	merged := user.Merge(defaults)
	if merged.MaxSnapshots != 50 {
		t.Errorf("expected MaxSnapshots 50, got %d", merged.MaxSnapshots)
	}
	if merged.MaxAge != 2*time.Hour {
		t.Errorf("expected MaxAge 2h, got %s", merged.MaxAge)
	}
}
