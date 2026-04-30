package config

import (
	"testing"
	"time"
)

func TestDefaultReloadConfig(t *testing.T) {
	cfg := DefaultReloadConfig()
	if !cfg.Enabled {
		t.Error("expected Enabled=true")
	}
	if cfg.Debounce != 500*time.Millisecond {
		t.Errorf("unexpected Debounce: %s", cfg.Debounce)
	}
	if cfg.OnReloadFail != "warn" {
		t.Errorf("unexpected OnReloadFail: %s", cfg.OnReloadFail)
	}
}

func TestReloadConfig_Validate_Disabled(t *testing.T) {
	cfg := ReloadConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestReloadConfig_Validate_Valid(t *testing.T) {
	cfg := DefaultReloadConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReloadConfig_Validate_NegativeDebounce(t *testing.T) {
	cfg := DefaultReloadConfig()
	cfg.Debounce = -1 * time.Millisecond
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative debounce")
	}
}

func TestReloadConfig_Validate_TooLargeDebounce(t *testing.T) {
	cfg := DefaultReloadConfig()
	cfg.Debounce = 11 * time.Second
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for debounce > 10s")
	}
}

func TestReloadConfig_Validate_InvalidOnReloadFail(t *testing.T) {
	cfg := DefaultReloadConfig()
	cfg.OnReloadFail = "ignore"
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for invalid on_reload_fail")
	}
}

func TestReloadConfig_Merge_FillsDefaults(t *testing.T) {
	base := DefaultReloadConfig()
	other := ReloadConfig{Debounce: 2 * time.Second, OnReloadFail: "fatal", Enabled: true}
	base.Merge(other)
	if base.Debounce != 2*time.Second {
		t.Errorf("expected 2s, got %s", base.Debounce)
	}
	if base.OnReloadFail != "fatal" {
		t.Errorf("expected fatal, got %s", base.OnReloadFail)
	}
}
