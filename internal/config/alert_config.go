package config

import "fmt"

// AlertConfig holds configuration for the alerting subsystem.
type AlertConfig struct {
	// Enabled controls whether alerts are dispatched at all.
	Enabled bool `yaml:"enabled"`

	// LogLevel is the log level used when emitting log-based alerts.
	// Valid values: "info", "warn", "error".
	LogLevel string `yaml:"log_level"`

	// IncludeProcess controls whether process info is included in alert messages.
	IncludeProcess bool `yaml:"include_process"`

	// DedupWindow is the number of seconds to suppress duplicate alerts
	// for the same listener/event combination.
	DedupWindow int `yaml:"dedup_window_seconds"`
}

// DefaultAlertConfig returns an AlertConfig populated with sensible defaults.
func DefaultAlertConfig() AlertConfig {
	return AlertConfig{
		Enabled:        true,
		LogLevel:       "warn",
		IncludeProcess: true,
		DedupWindow:    60,
	}
}

// Validate checks that all AlertConfig fields contain acceptable values.
func (a AlertConfig) Validate() error {
	switch a.LogLevel {
	case "info", "warn", "error":
		// valid
	default:
		return fmt.Errorf("alert.log_level %q is invalid: must be info, warn, or error", a.LogLevel)
	}
	if a.DedupWindow < 0 {
		return fmt.Errorf("alert.dedup_window_seconds must be >= 0, got %d", a.DedupWindow)
	}
	return nil
}

// Merge returns a new AlertConfig where zero/empty fields in a are filled
// from defaults d.
func (a AlertConfig) Merge(d AlertConfig) AlertConfig {
	if a.LogLevel == "" {
		a.LogLevel = d.LogLevel
	}
	if a.DedupWindow == 0 {
		a.DedupWindow = d.DedupWindow
	}
	return a
}
