package config

import "fmt"

// AuditNotifierConfig controls whether audit log notifier is wired into
// the alert pipeline and which writer settings it uses.
type AuditNotifierConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
	Format  string `yaml:"format"` // "json" or "text"
}

// DefaultAuditNotifierConfig returns a safe default: disabled.
func DefaultAuditNotifierConfig() AuditNotifierConfig {
	return AuditNotifierConfig{
		Enabled: false,
		Path:    "/var/log/portwatch/audit.log",
		Format:  "json",
	}
}

// Validate checks that required fields are present when enabled.
func (c AuditNotifierConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Path == "" {
		return fmt.Errorf("audit_notifier: path must not be empty when enabled")
	}
	if c.Format != "json" && c.Format != "text" {
		return fmt.Errorf("audit_notifier: invalid format %q, must be 'json' or 'text'", c.Format)
	}
	return nil
}

// Merge fills zero-value fields from defaults.
func (c AuditNotifierConfig) Merge(defaults AuditNotifierConfig) AuditNotifierConfig {
	if c.Path == "" {
		c.Path = defaults.Path
	}
	if c.Format == "" {
		c.Format = defaults.Format
	}
	return c
}
