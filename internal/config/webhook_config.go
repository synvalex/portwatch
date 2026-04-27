package config

import (
	"fmt"
	"time"
)

// WebhookConfig holds settings for the optional HTTP webhook notifier.
type WebhookConfig struct {
	Enabled bool          `yaml:"enabled"`
	URL     string        `yaml:"url"`
	Timeout time.Duration `yaml:"timeout"`
}

// DefaultWebhookConfig returns a WebhookConfig with safe defaults.
func DefaultWebhookConfig() WebhookConfig {
	return WebhookConfig{
		Enabled: false,
		URL:     "",
		Timeout: 5 * time.Second,
	}
}

// Validate returns an error if the configuration is invalid.
func (w WebhookConfig) Validate() error {
	if !w.Enabled {
		return nil
	}
	if w.URL == "" {
		return fmt.Errorf("webhook: url must not be empty when enabled")
	}
	if w.Timeout < 0 {
		return fmt.Errorf("webhook: timeout must not be negative")
	}
	return nil
}

// Merge fills zero-value fields in w with values from defaults.
func (w WebhookConfig) Merge(defaults WebhookConfig) WebhookConfig {
	if w.Timeout == 0 {
		w.Timeout = defaults.Timeout
	}
	return w
}
