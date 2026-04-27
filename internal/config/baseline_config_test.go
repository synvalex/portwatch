package config

import (
	"testing"
	"time"
)

func TestDefaultBaselineConfig(t *testing.T) {
	cfg := DefaultBaselineConfig()
	if cfg.Enabled {
		t.Error("expected baseline to be disabled by default")
	}
	if cfg.File == "" {
		t.Error("expected default file path to be set")
	}
	if cfg.LearnDuration <= 0 {
		t.Error("expected positive learn_duration")
	}
}

func TestBaselineConfig_Validate_Disabled(t *testing.T) {
	cfg := BaselineConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Errorf("disabled baseline should always validate, got: %v", err)
	}
}

func TestBaselineConfig_Validate_EnabledValid(t *testing.T) {
	cfg := BaselineConfig{
		Enabled:       true,
		File:          "/tmp/baseline.json",
		LearnDuration: 2 * time.Minute,
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected valid config, got: %v", err)
	}
}

func TestBaselineConfig_Validate_EnabledMissingFile(t *testing.T) {
	cfg := BaselineConfig{Enabled: true, File: ""}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing file path")
	}
}

func TestBaselineConfig_Validate_NegativeDuration(t *testing.T) {
	cfg := BaselineConfig{
		Enabled:       true,
		File:          "/tmp/b.json",
		LearnDuration: -1 * time.Second,
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for negative learn_duration")
	}
}

func TestBaselineConfig_Merge_FillsEmptyFile(t *testing.T) {
	cfg := BaselineConfig{Enabled: true}
	cfg.Merge(DefaultBaselineConfig())
	if cfg.File == "" {
		t.Error("expected file to be filled from defaults")
	}
}

func TestBaselineConfig_Merge_UserFilePreserved(t *testing.T) {
	cfg := BaselineConfig{File: "/custom/path.json"}
	cfg.Merge(DefaultBaselineConfig())
	if cfg.File != "/custom/path.json" {
		t.Errorf("expected user file to be preserved, got %s", cfg.File)
	}
}
