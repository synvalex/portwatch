package alert

import (
	"fmt"
	"log/slog"

	"github.com/user/portwatch/internal/config"
)

// NotifierBuilder constructs a MultiNotifier from the full application config.
type NotifierBuilder struct {
	logger *slog.Logger
}

// NewNotifierBuilder creates a NotifierBuilder using the given logger.
func NewNotifierBuilder(logger *slog.Logger) *NotifierBuilder {
	return &NotifierBuilder{logger: logger}
}

// Build assembles all configured notifiers into a single MultiNotifier.
// Enabled notifiers are added; disabled ones are silently skipped.
func (b *NotifierBuilder) Build(cfg config.Config) (Notifier, error) {
	mn := NewMultiNotifier()

	// Always include the log notifier.
	mn.Add(NewLogNotifier(b.logger))

	if cfg.Alert.Webhook.Enabled {
		wh, err := NewWebhookNotifier(cfg.Alert.Webhook.URL, cfg.Alert.Webhook.Timeout, b.logger)
		if err != nil {
			return nil, fmt.Errorf("notifier_builder: webhook: %w", err)
		}
		mn.Add(wh)
		b.logger.Info("webhook notifier enabled", "url", cfg.Alert.Webhook.URL)
	}

	if cfg.Alert.Slack.Enabled {
		sl := NewSlackNotifier(cfg.Alert.Slack.WebhookURL, cfg.Alert.Slack.Timeout, b.logger)
		mn.Add(sl)
		b.logger.Info("slack notifier enabled")
	}

	if cfg.Alert.Email.Enabled {
		em := NewEmailNotifier(
			cfg.Alert.Email.SMTPHost,
			cfg.Alert.Email.SMTPPort,
			cfg.Alert.Email.Username,
			cfg.Alert.Email.Password,
			cfg.Alert.Email.From,
			cfg.Alert.Email.To,
			cfg.Alert.Email.Subject,
			b.logger,
		)
		mn.Add(em)
		b.logger.Info("email notifier enabled", "recipients", cfg.Alert.Email.To)
	}

	return mn, nil
}
