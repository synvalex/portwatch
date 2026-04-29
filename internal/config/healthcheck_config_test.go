package config

import (
	"testing"
	"time"
)

func TestDefaultHealthCheckConfig(t *testing.T) {
	cfg := DefaultHealthCheckConfig()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.Addr != ":9110" {
		t.Errorf("unexpected default addr: %s", cfg.Addr)
	}
	if cfg.Path != "/healthz" {
		t.Errorf("unexpected default path: %s", cfg.Path)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("unexpected default timeout: %v", cfg.Timeout)
	}
}

func TestHealthCheckConfig_Validate_Disabled(t *testing.T) {
	cfg := HealthCheckConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got: %v", err)
	}
}

func TestHealthCheckConfig_Validate_EnabledValid(t *testing.T) {
	cfg := HealthCheckConfig{Enabled: true, Addr: ":9110", Path: "/healthz", Timeout: 5 * time.Second}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestHealthCheckConfig_Validate_MissingAddr(t *testing.T) {
	cfg := HealthCheckConfig{Enabled: true, Path: "/healthz", Timeout: 5 * time.Second}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing addr")
	}
}

func TestHealthCheckConfig_Validate_NegativeTimeout(t *testing.T) {
	cfg := HealthCheckConfig{Enabled: true, Addr: ":9110", Path: "/healthz", Timeout: -1}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for negative timeout")
	}
}

func TestHealthCheckConfig_Merge_FillsDefaults(t *testing.T) {
	cfg := HealthCheckConfig{Enabled: true}
	cfg.Merge(DefaultHealthCheckConfig())
	if cfg.Addr != ":9110" {
		t.Errorf("expected addr to be filled from defaults, got %s", cfg.Addr)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("expected timeout to be filled from defaults, got %v", cfg.Timeout)
	}
}

func TestHealthCheckConfig_Merge_UserValuesPreserved(t *testing.T) {
	cfg := HealthCheckConfig{Enabled: true, Addr: ":8888", Path: "/ping", Timeout: 2 * time.Second}
	cfg.Merge(DefaultHealthCheckConfig())
	if cfg.Addr != ":8888" {
		t.Errorf("expected user addr to be preserved, got %s", cfg.Addr)
	}
	if cfg.Path != "/ping" {
		t.Errorf("expected user path to be preserved, got %s", cfg.Path)
	}
}
