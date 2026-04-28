package config

import (
	"testing"
	"time"
)

func TestDefaultHistoryConfig(t *testing.T) {
	c := DefaultHistoryConfig()
	if !c.Enabled {
		t.Error("expected Enabled to be true")
	}
	if c.MaxEvents != 500 {
		t.Errorf("expected MaxEvents 500, got %d", c.MaxEvents)
	}
	if c.Retention != 24*time.Hour {
		t.Errorf("expected Retention 24h, got %s", c.Retention)
	}
}

func TestHistoryConfig_Validate_Disabled(t *testing.T) {
	c := HistoryConfig{Enabled: false}
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got %v", err)
	}
}

func TestHistoryConfig_Validate_Valid(t *testing.T) {
	c := HistoryConfig{Enabled: true, MaxEvents: 100, Retention: time.Hour}
	if err := c.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHistoryConfig_Validate_ZeroMaxEvents(t *testing.T) {
	c := HistoryConfig{Enabled: true, MaxEvents: 0, Retention: time.Hour}
	if err := c.Validate(); err == nil {
		t.Error("expected error for MaxEvents=0")
	}
}

func TestHistoryConfig_Validate_ZeroRetention(t *testing.T) {
	c := HistoryConfig{Enabled: true, MaxEvents: 10, Retention: 0}
	if err := c.Validate(); err == nil {
		t.Error("expected error for Retention=0")
	}
}

func TestHistoryConfig_Merge_FillsDefaults(t *testing.T) {
	user := HistoryConfig{Enabled: true}
	defaults := DefaultHistoryConfig()
	merged := user.Merge(defaults)
	if merged.MaxEvents != defaults.MaxEvents {
		t.Errorf("expected MaxEvents %d, got %d", defaults.MaxEvents, merged.MaxEvents)
	}
	if merged.Retention != defaults.Retention {
		t.Errorf("expected Retention %s, got %s", defaults.Retention, merged.Retention)
	}
}

func TestHistoryConfig_Merge_UserValuesPreserved(t *testing.T) {
	user := HistoryConfig{Enabled: true, MaxEvents: 50, Retention: 2 * time.Hour}
	merged := user.Merge(DefaultHistoryConfig())
	if merged.MaxEvents != 50 {
		t.Errorf("expected MaxEvents 50, got %d", merged.MaxEvents)
	}
	if merged.Retention != 2*time.Hour {
		t.Errorf("expected Retention 2h, got %s", merged.Retention)
	}
}
