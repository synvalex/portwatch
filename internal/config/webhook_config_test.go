package config

import (
	"testing"
	"time"
)

func TestDefaultWebhookConfig(t *testing.T) {
	cfg := DefaultWebhookConfig()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.URL != "" {
		t.Errorf("expected empty URL, got %q", cfg.URL)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("expected 5s timeout, got %v", cfg.Timeout)
	}
}

func TestWebhookConfig_Validate_DisabledNoURL(t *testing.T) {
	cfg := WebhookConfig{Enabled: false, URL: ""}
	if err := cfg.Validate(); err != nil {
		t.Errorf("disabled config should be valid, got: %v", err)
	}
}

func TestWebhookConfig_Validate_EnabledWithURL(t *testing.T) {
	cfg := WebhookConfig{Enabled: true, URL: "http://example.com/hook", Timeout: 3 * time.Second}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected valid, got: %v", err)
	}
}

func TestWebhookConfig_Validate_EnabledMissingURL(t *testing.T) {
	cfg := WebhookConfig{Enabled: true, URL: ""}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for enabled webhook with no URL")
	}
}

func TestWebhookConfig_Validate_NegativeTimeout(t *testing.T) {
	cfg := WebhookConfig{Enabled: true, URL: "http://example.com", Timeout: -1 * time.Second}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for negative timeout")
	}
}

func TestWebhookConfig_Validate_ZeroTimeout(t *testing.T) {
	// A zero timeout on an enabled webhook should be rejected, as it would
	// cause all webhook requests to time out immediately.
	cfg := WebhookConfig{Enabled: true, URL: "http://example.com", Timeout: 0}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero timeout on enabled webhook")
	}
}

func TestWebhookConfig_Merge_TimeoutFromDefault(t *testing.T) {
	defaults := DefaultWebhookConfig()
	cfg := WebhookConfig{Enabled: true, URL: "http://example.com"}
	merged := cfg.Merge(defaults)
	if merged.Timeout != defaults.Timeout {
		t.Errorf("expected timeout %v from defaults, got %v", defaults.Timeout, merged.Timeout)
	}
}

func TestWebhookConfig_Merge_UserTimeoutPreserved(t *testing.T) {
	defaults := DefaultWebhookConfig()
	cfg := WebhookConfig{Enabled: true, URL: "http://example.com", Timeout: 10 * time.Second}
	merged := cfg.Merge(defaults)
	if merged.Timeout != 10*time.Second {
		t.Errorf("expected user timeout 10s preserved, got %v", merged.Timeout)
	}
}
