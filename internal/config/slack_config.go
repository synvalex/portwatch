package config

import (
	"fmt"
	"time"
)

// SlackConfig holds configuration for the Slack notifier.
type SlackConfig struct {
	Enabled    bool          `yaml:"enabled"`
	WebhookURL string        `yaml:"webhook_url"`
	Timeout    time.Duration `yaml:"timeout"`
}

// DefaultSlackConfig returns a SlackConfig with sensible defaults.
func DefaultSlackConfig() SlackConfig {
	return SlackConfig{
		Enabled:    false,
		WebhookURL: "",
		Timeout:    5 * time.Second,
	}
}

// Validate checks that the SlackConfig is consistent.
func (c *SlackConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.WebhookURL == "" {
		return fmt.Errorf("slack: webhook_url must be set when enabled")
	}
	if c.Timeout < 0 {
		return fmt.Errorf("slack: timeout must not be negative")
	}
	return nil
}

// Merge fills zero-value fields in c from defaults.
func (c *SlackConfig) Merge(defaults SlackConfig) {
	if c.Timeout == 0 {
		c.Timeout = defaults.Timeout
	}
}
