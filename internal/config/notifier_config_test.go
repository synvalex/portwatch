package config

import (
	"testing"
)

func TestDefaultNotifierConfig(t *testing.T) {
	cfg := DefaultNotifierConfig()

	if !cfg.Log.Enabled {
		t.Error("expected log notifier to be enabled by default")
	}
	if cfg.Log.Level != "info" {
		t.Errorf("expected default log level 'info', got %q", cfg.Log.Level)
	}
	if cfg.Webhook.Enabled {
		t.Error("expected webhook to be disabled by default")
	}
	if cfg.Slack.Enabled {
		t.Error("expected slack to be disabled by default")
	}
	if cfg.Email.Enabled {
		t.Error("expected email to be disabled by default")
	}
}

func TestLogNotifierConfig_Validate_Valid(t *testing.T) {
	for _, level := range []string{"debug", "info", "warn", "error"} {
		cfg := LogNotifierConfig{Enabled: true, Level: level}
		if err := cfg.Validate(); err != nil {
			t.Errorf("expected valid level %q to pass, got: %v", level, err)
		}
	}
}

func TestLogNotifierConfig_Validate_InvalidLevel(t *testing.T) {
	cfg := LogNotifierConfig{Enabled: true, Level: "verbose"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for invalid log level")
	}
}

func TestLogNotifierConfig_Validate_DisabledSkipsLevelCheck(t *testing.T) {
	cfg := LogNotifierConfig{Enabled: false, Level: "verbose"}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected disabled notifier to skip level validation, got: %v", err)
	}
}

func TestNotifierConfig_Merge_FillsEmptyLevel(t *testing.T) {
	user := NotifierConfig{Log: LogNotifierConfig{Enabled: true, Level: ""}}
	defaults := DefaultNotifierConfig()
	merged := user.Merge(defaults)
	if merged.Log.Level != "info" {
		t.Errorf("expected merged level 'info', got %q", merged.Log.Level)
	}
}

func TestNotifierConfig_Merge_PreservesUserLevel(t *testing.T) {
	user := NotifierConfig{Log: LogNotifierConfig{Enabled: true, Level: "warn"}}
	defaults := DefaultNotifierConfig()
	merged := user.Merge(defaults)
	if merged.Log.Level != "warn" {
		t.Errorf("expected user level 'warn' to be preserved, got %q", merged.Log.Level)
	}
}

func TestNotifierConfig_Validate_Valid(t *testing.T) {
	cfg := DefaultNotifierConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected default config to be valid, got: %v", err)
	}
}
