package config

import "fmt"

// MetricsConfig controls the optional Prometheus metrics endpoint.
type MetricsConfig struct {
	Enabled bool   `yaml:"enabled"`
	Addr    string `yaml:"addr"`
	Path    string `yaml:"path"`
}

// DefaultMetricsConfig returns a conservative default: disabled.
func DefaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		Enabled: false,
		Addr:    ":9090",
		Path:    "/metrics",
	}
}

func (c *MetricsConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Addr == "" {
		return fmt.Errorf("metrics.addr must not be empty when metrics is enabled")
	}
	if c.Path == "" {
		return fmt.Errorf("metrics.path must not be empty when metrics is enabled")
	}
	return nil
}

func (c *MetricsConfig) Merge(def MetricsConfig) {
	if c.Addr == "" {
		c.Addr = def.Addr
	}
	if c.Path == "" {
		c.Path = def.Path
	}
}
