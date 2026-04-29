package config

import (
	"errors"
	"time"
)

// HealthCheckConfig controls the optional HTTP health-check endpoint.
type HealthCheckConfig struct {
	Enabled bool          `yaml:"enabled"`
	Addr    string        `yaml:"addr"`
	Path    string        `yaml:"path"`
	Timeout time.Duration `yaml:"timeout"`
}

// DefaultHealthCheckConfig returns a safe default configuration.
func DefaultHealthCheckConfig() HealthCheckConfig {
	return HealthCheckConfig{
		Enabled: false,
		Addr:    ":9110",
		Path:    "/healthz",
		Timeout: 5 * time.Second,
	}
}

// Validate checks that the health-check configuration is coherent.
func (h *HealthCheckConfig) Validate() error {
	if !h.Enabled {
		return nil
	}
	if h.Addr == "" {
		return errors.New("healthcheck: addr must not be empty when enabled")
	}
	if h.Path == "" {
		return errors.New("healthcheck: path must not be empty when enabled")
	}
	if h.Timeout <= 0 {
		return errors.New("healthcheck: timeout must be positive")
	}
	return nil
}

// Merge fills zero-value fields from defaults.
func (h *HealthCheckConfig) Merge(defaults HealthCheckConfig) {
	if h.Addr == "" {
		h.Addr = defaults.Addr
	}
	if h.Path == "" {
		h.Path = defaults.Path
	}
	if h.Timeout == 0 {
		h.Timeout = defaults.Timeout
	}
}
