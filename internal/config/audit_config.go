package config

import (
	"fmt"
	"time"
)

// AuditConfig controls the audit log feature that writes structured
// scan events to a rotating log file.
type AuditConfig struct {
	Enabled    bool          `yaml:"enabled"`
	Path       string        `yaml:"path"`
	MaxSizeMB  int           `yaml:"max_size_mb"`
	MaxBackups int           `yaml:"max_backups"`
	Rotation   time.Duration `yaml:"rotation"`
	Format     string        `yaml:"format"` // "json" or "text"
}

// DefaultAuditConfig returns a safe default AuditConfig.
func DefaultAuditConfig() AuditConfig {
	return AuditConfig{
		Enabled:    false,
		Path:       "/var/log/portwatch/audit.log",
		MaxSizeMB:  50,
		MaxBackups: 5,
		Rotation:   24 * time.Hour,
		Format:     "json",
	}
}

// Validate returns an error if the AuditConfig is invalid.
func (a AuditConfig) Validate() error {
	if !a.Enabled {
		return nil
	}
	if a.Path == "" {
		return fmt.Errorf("audit: path must not be empty when enabled")
	}
	if a.MaxSizeMB <= 0 {
		return fmt.Errorf("audit: max_size_mb must be positive, got %d", a.MaxSizeMB)
	}
	if a.MaxBackups < 0 {
		return fmt.Errorf("audit: max_backups must be non-negative, got %d", a.MaxBackups)
	}
	if a.Rotation <= 0 {
		return fmt.Errorf("audit: rotation must be positive, got %s", a.Rotation)
	}
	if a.Format != "json" && a.Format != "text" {
		return fmt.Errorf("audit: format must be \"json\" or \"text\", got %q", a.Format)
	}
	return nil
}

// Merge fills zero-value fields in a from defaults.
func (a AuditConfig) Merge(defaults AuditConfig) AuditConfig {
	if a.Path == "" {
		a.Path = defaults.Path
	}
	if a.MaxSizeMB == 0 {
		a.MaxSizeMB = defaults.MaxSizeMB
	}
	if a.MaxBackups == 0 {
		a.MaxBackups = defaults.MaxBackups
	}
	if a.Rotation == 0 {
		a.Rotation = defaults.Rotation
	}
	if a.Format == "" {
		a.Format = defaults.Format
	}
	return a
}
