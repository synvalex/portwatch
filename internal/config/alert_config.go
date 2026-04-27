package config

import (
	"errors"
	"fmt"
	"time"
)

// AlertConfig holds all alert-related configuration.
type AlertConfig struct {
	Level       string        `yaml:"level"`
	DedupWindow time.Duration `yaml:"dedup_window"`
	Webhook     WebhookConfig `yaml:"webhook"`
	Slack       SlackConfig   `yaml:"slack"`
	Email       EmailConfig   `yaml:"email"`
}

// DefaultAlertConfig returns sensible alert defaults.
func DefaultAlertConfig() AlertConfig {
	return AlertConfig{
		Level:       "info",
		DedupWindow: 30 * time.Second,
		Webhook:     DefaultWebhookConfig(),
		Slack:       DefaultSlackConfig(),
		Email:       DefaultEmailConfig(),
	}
}

var validAlertLevels = map[string]struct{}{
	"debug": {},
	"info":  {},
	"warn":  {},
	"error": {},
}

// Validate returns an error if AlertConfig contains invalid values.
func (a AlertConfig) Validate() error {
	if _, ok := validAlertLevels[a.Level]; !ok {
		return fmt.Errorf("alert: unknown level %q", a.Level)
	}
	if a.DedupWindow < 0 {
		return errors.New("alert: dedup_window must not be negative")
	}
	if a.DedupWindow == 0 {
		return errors.New("alert: dedup_window must be greater than zero")
	}
	if err := a.Webhook.Validate(); err != nil {
		return fmt.Errorf("alert.webhook: %w", err)
	}
	if err := a.Slack.Validate(); err != nil {
		return fmt.Errorf("alert.slack: %w", err)
	}
	if err := a.Email.Validate(); err != nil {
		return fmt.Errorf("alert.email: %w", err)
	}
	return nil
}

// Merge fills zero-value fields in a from defaults.
func (a AlertConfig) Merge(defaults AlertConfig) AlertConfig {
	if a.Level == "" {
		a.Level = defaults.Level
	}
	if a.DedupWindow == 0 {
		a.DedupWindow = defaults.DedupWindow
	}
	a.Webhook = a.Webhook.Merge(defaults.Webhook)
	a.Slack = a.Slack.Merge(defaults.Slack)
	a.Email = a.Email.Merge(defaults.Email)
	return a
}
