package config

import (
	"testing"
	"time"
)

func TestDefaultSlackConfig(t *testing.T) {
	c := DefaultSlackConfig()
	if c.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if c.WebhookURL != "" {
		t.Errorf("expected empty WebhookURL, got %q", c.WebhookURL)
	}
	if c.Timeout != 5*time.Second {
		t.Errorf("expected 5s timeout, got %v", c.Timeout)
	}
}

func TestSlackConfig_Validate_DisabledNoURL(t *testing.T) {
	c := SlackConfig{Enabled: false}
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got %v", err)
	}
}

func TestSlackConfig_Validate_EnabledWithURL(t *testing.T) {
	c := SlackConfig{Enabled: true, WebhookURL: "https://hooks.slack.com/x", Timeout: time.Second}
	if err := c.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSlackConfig_Validate_EnabledMissingURL(t *testing.T) {
	c := SlackConfig{Enabled: true, WebhookURL: ""}
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing webhook_url")
	}
}

func TestSlackConfig_Validate_NegativeTimeout(t *testing.T) {
	c := SlackConfig{Enabled: true, WebhookURL: "https://hooks.slack.com/x", Timeout: -1}
	if err := c.Validate(); err == nil {
		t.Error("expected error for negative timeout")
	}
}

func TestSlackConfig_Merge_TimeoutFromDefault(t *testing.T) {
	c := SlackConfig{Enabled: true, WebhookURL: "https://hooks.slack.com/x"}
	defaults := DefaultSlackConfig()
	c.Merge(defaults)
	if c.Timeout != defaults.Timeout {
		t.Errorf("expected timeout %v from defaults, got %v", defaults.Timeout, c.Timeout)
	}
}

func TestSlackConfig_Merge_UserTimeoutPreserved(t *testing.T) {
	c := SlackConfig{Enabled: true, WebhookURL: "https://hooks.slack.com/x", Timeout: 10 * time.Second}
	c.Merge(DefaultSlackConfig())
	if c.Timeout != 10*time.Second {
		t.Errorf("expected user timeout preserved, got %v", c.Timeout)
	}
}
