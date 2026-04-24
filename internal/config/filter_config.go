package config

import "slices"

// FilterConfig holds user-facing filter settings loaded from YAML.
type FilterConfig struct {
	// ExcludeLoopback controls whether loopback-bound listeners are hidden.
	ExcludeLoopback bool `yaml:"exclude_loopback"`
	// ExcludePorts lists port numbers that should never trigger alerts.
	ExcludePorts []uint16 `yaml:"exclude_ports"`
}

// DefaultFilterConfig returns safe defaults: loopback excluded, no port exclusions.
func DefaultFilterConfig() FilterConfig {
	return FilterConfig{
		ExcludeLoopback: true,
		ExcludePorts:    []uint16{},
	}
}

// Merge returns a new FilterConfig where zero values fall back to defaults.
func (f FilterConfig) Merge(defaults FilterConfig) FilterConfig {
	out := f
	if len(out.ExcludePorts) == 0 {
		out.ExcludePorts = defaults.ExcludePorts
	}
	return out
}

// IsPortExcluded reports whether the given port number appears in ExcludePorts.
func (f FilterConfig) IsPortExcluded(port uint16) bool {
	return slices.Contains(f.ExcludePorts, port)
}
