package config

import "fmt"

// FingerprintConfig controls process fingerprinting for listeners.
type FingerprintConfig struct {
	Enabled       bool     `yaml:"enabled"`
	HashExecutable bool    `yaml:"hash_executable"`
	TrustCmdline  bool     `yaml:"trust_cmdline"`
	AllowList     []string `yaml:"allow_list"`
}

// DefaultFingerprintConfig returns sensible defaults.
func DefaultFingerprintConfig() FingerprintConfig {
	return FingerprintConfig{
		Enabled:        true,
		HashExecutable: true,
		TrustCmdline:   false,
		AllowList:      []string{},
	}
}

// Validate checks the fingerprint configuration for consistency.
func (c FingerprintConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if !c.HashExecutable && !c.TrustCmdline {
		return fmt.Errorf("fingerprint: at least one of hash_executable or trust_cmdline must be enabled")
	}
	return nil
}

// Merge fills zero-value fields from defaults.
func (c FingerprintConfig) Merge(defaults FingerprintConfig) FingerprintConfig {
	if len(c.AllowList) == 0 {
		c.AllowList = defaults.AllowList
	}
	return c
}
