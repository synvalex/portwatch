package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTempConfig(t, `
interval: 30s
log_level: debug
rules:
  - name: allow-ssh
    port: 22
    proto: tcp
    address: 0.0.0.0
    action: allow
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("interval = %v, want 30s", cfg.Interval)
	}
	if len(cfg.Rules) != 1 || cfg.Rules[0].Name != "allow-ssh" {
		t.Errorf("unexpected rules: %+v", cfg.Rules)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/portwatch.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestValidate_ShortInterval(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Interval = 500 * time.Millisecond
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for short interval")
	}
}

func TestValidate_BadLogLevel(t *testing.T) {
	cfg := DefaultConfig()
	cfg.LogLevel = "verbose"
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for unknown log level")
	}
}

func TestValidate_RuleMissingName(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Rules = []RuleConfig{{Port: 80, Action: "allow"}}
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for rule missing name")
	}
}

func TestValidate_BadAction(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Rules = []RuleConfig{{Name: "test", Port: 80, Action: "block"}}
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for bad action")
	}
}
