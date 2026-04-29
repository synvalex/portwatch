package config

import (
	"testing"
	"time"
)

func TestDefaultThrottleConfig(t *testing.T) {
	cfg := DefaultThrottleConfig()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.Rate != 10 {
		t.Errorf("expected Rate=10, got %d", cfg.Rate)
	}
	if cfg.Window != time.Minute {
		t.Errorf("expected Window=1m, got %s", cfg.Window)
	}
	if !cfg.PerPort {
		t.Error("expected PerPort=true by default")
	}
}

func TestThrottleConfig_Validate_Disabled(t *testing.T) {
	cfg := ThrottleConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got %v", err)
	}
}

func TestThrottleConfig_Validate_Valid(t *testing.T) {
	cfg := ThrottleConfig{Enabled: true, Rate: 5, Window: 30 * time.Second}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestThrottleConfig_Validate_ZeroRate(t *testing.T) {
	cfg := ThrottleConfig{Enabled: true, Rate: 0, Window: time.Minute}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero rate")
	}
}

func TestThrottleConfig_Validate_ZeroWindow(t *testing.T) {
	cfg := ThrottleConfig{Enabled: true, Rate: 5, Window: 0}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero window")
	}
}

func TestThrottleConfig_Merge_FillsZeroValues(t *testing.T) {
	defaults := DefaultThrottleConfig()
	cfg := ThrottleConfig{Enabled: true}
	merged := cfg.Merge(defaults)
	if merged.Rate != defaults.Rate {
		t.Errorf("expected Rate=%d, got %d", defaults.Rate, merged.Rate)
	}
	if merged.Window != defaults.Window {
		t.Errorf("expected Window=%s, got %s", defaults.Window, merged.Window)
	}
}

func TestThrottleConfig_Merge_UserValuesPreserved(t *testing.T) {
	defaults := DefaultThrottleConfig()
	cfg := ThrottleConfig{Enabled: true, Rate: 3, Window: 10 * time.Second}
	merged := cfg.Merge(defaults)
	if merged.Rate != 3 {
		t.Errorf("expected Rate=3, got %d", merged.Rate)
	}
	if merged.Window != 10*time.Second {
		t.Errorf("expected Window=10s, got %s", merged.Window)
	}
}
