package config

import (
	"errors"
	"time"
)

// EmailConfig holds configuration for email alert notifications.
type EmailConfig struct {
	Enabled    bool          `yaml:"enabled"`
	SMTPHost   string        `yaml:"smtp_host"`
	SMTPPort   int           `yaml:"smtp_port"`
	Username   string        `yaml:"username"`
	Password   string        `yaml:"password"`
	From       string        `yaml:"from"`
	To         []string      `yaml:"to"`
	Subject    string        `yaml:"subject"`
	Timeout    time.Duration `yaml:"timeout"`
}

// DefaultEmailConfig returns a sane default EmailConfig.
func DefaultEmailConfig() EmailConfig {
	return EmailConfig{
		Enabled:  false,
		SMTPPort: 587,
		Subject:  "portwatch alert",
		Timeout:  10 * time.Second,
	}
}

// Validate returns an error if the EmailConfig is invalid.
func (e EmailConfig) Validate() error {
	if !e.Enabled {
		return nil
	}
	if e.SMTPHost == "" {
		return errors.New("email: smtp_host is required when enabled")
	}
	if e.SMTPPort <= 0 || e.SMTPPort > 65535 {
		return errors.New("email: smtp_port must be between 1 and 65535")
	}
	if e.From == "" {
		return errors.New("email: from address is required when enabled")
	}
	if len(e.To) == 0 {
		return errors.New("email: at least one recipient is required when enabled")
	}
	if e.Timeout < 0 {
		return errors.New("email: timeout must not be negative")
	}
	return nil
}

// Merge fills zero-value fields in e from defaults.
func (e EmailConfig) Merge(defaults EmailConfig) EmailConfig {
	if e.SMTPPort == 0 {
		e.SMTPPort = defaults.SMTPPort
	}
	if e.Subject == "" {
		e.Subject = defaults.Subject
	}
	if e.Timeout == 0 {
		e.Timeout = defaults.Timeout
	}
	return e
}
