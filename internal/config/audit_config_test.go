package config

import (
	"testing"
	"time"
)

func TestDefaultAuditConfig(t *testing.T) {
	cfg := DefaultAuditConfig()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.Path == "" {
		t.Error("expected non-empty default Path")
	}
	if cfg.MaxSizeMB <= 0 {
		t.Errorf("expected positive MaxSizeMB, got %d", cfg.MaxSizeMB)
	}
	if cfg.Format != "json" {
		t.Errorf("expected default format \"json\", got %q", cfg.Format)
	}
}

func TestAuditConfig_Validate_Disabled(t *testing.T) {
	cfg := AuditConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Errorf("disabled config should always be valid, got: %v", err)
	}
}

func TestAuditConfig_Validate_EnabledValid(t *testing.T) {
	cfg := AuditConfig{
		Enabled:    true,
		Path:       "/tmp/audit.log",
		MaxSizeMB:  10,
		MaxBackups: 3,
		Rotation:   time.Hour,
		Format:     "json",
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected valid config, got: %v", err)
	}
}

func TestAuditConfig_Validate_MissingPath(t *testing.T) {
	cfg := AuditConfig{Enabled: true, MaxSizeMB: 10, MaxBackups: 2, Rotation: time.Hour, Format: "json"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestAuditConfig_Validate_InvalidFormat(t *testing.T) {
	cfg := AuditConfig{
		Enabled: true, Path: "/tmp/a.log", MaxSizeMB: 5,
		MaxBackups: 1, Rotation: time.Hour, Format: "csv",
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for invalid format")
	}
}

func TestAuditConfig_Validate_NegativeMaxSize(t *testing.T) {
	cfg := AuditConfig{
		Enabled: true, Path: "/tmp/a.log", MaxSizeMB: -1,
		MaxBackups: 1, Rotation: time.Hour, Format: "json",
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for negative MaxSizeMB")
	}
}

func TestAuditConfig_Merge_FillsEmptyFields(t *testing.T) {
	user := AuditConfig{Enabled: true, Path: "/custom/audit.log"}
	defaults := DefaultAuditConfig()
	merged := user.Merge(defaults)
	if merged.Path != "/custom/audit.log" {
		t.Errorf("user path should be preserved, got %q", merged.Path)
	}
	if merged.MaxSizeMB != defaults.MaxSizeMB {
		t.Errorf("expected default MaxSizeMB %d, got %d", defaults.MaxSizeMB, merged.MaxSizeMB)
	}
	if merged.Format != defaults.Format {
		t.Errorf("expected default format %q, got %q", defaults.Format, merged.Format)
	}
}
