package config

// NotifierConfig aggregates all notifier-specific configurations
// and controls which notification channels are active.
type NotifierConfig struct {
	Log     LogNotifierConfig     `yaml:"log"`
	Webhook WebhookConfig         `yaml:"webhook"`
	Slack   SlackConfig           `yaml:"slack"`
	Email   EmailConfig           `yaml:"email"`
}

// LogNotifierConfig controls the built-in log notifier.
type LogNotifierConfig struct {
	Enabled bool   `yaml:"enabled"`
	Level   string `yaml:"level"`
}

// DefaultLogNotifierConfig returns sensible defaults for the log notifier.
func DefaultLogNotifierConfig() LogNotifierConfig {
	return LogNotifierConfig{
		Enabled: true,
		Level:   "info",
	}
}

// Validate checks that the log notifier config is consistent.
func (l LogNotifierConfig) Validate() error {
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if l.Enabled && !validLevels[l.Level] {
		return &ValidationError{Field: "log.level", Message: "must be one of: debug, info, warn, error"}
	}
	return nil
}

// Merge fills zero-value fields in n from defaults d.
func (n NotifierConfig) Merge(d NotifierConfig) NotifierConfig {
	if n.Log.Level == "" {
		n.Log.Level = d.Log.Level
	}
	return n
}

// DefaultNotifierConfig returns a NotifierConfig with all defaults applied.
func DefaultNotifierConfig() NotifierConfig {
	return NotifierConfig{
		Log:     DefaultLogNotifierConfig(),
		Webhook: DefaultWebhookConfig(),
		Slack:   DefaultSlackConfig(),
		Email:   DefaultEmailConfig(),
	}
}

// Validate checks all nested notifier configurations.
func (n NotifierConfig) Validate() error {
	if err := n.Log.Validate(); err != nil {
		return err
	}
	if err := n.Webhook.Validate(); err != nil {
		return err
	}
	if err := n.Slack.Validate(); err != nil {
		return err
	}
	if err := n.Email.Validate(); err != nil {
		return err
	}
	return nil
}
