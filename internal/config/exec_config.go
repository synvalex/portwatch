package config

import (
	"fmt"
	"time"
)

// ExecConfig holds configuration for the exec (run-command) notifier.
type ExecConfig struct {
	Enabled bool          `yaml:"enabled"`
	Command string        `yaml:"command"`
	Args    []string      `yaml:"args"`
	Timeout time.Duration `yaml:"timeout"`
	Shell   bool          `yaml:"shell"`
}

// DefaultExecConfig returns a safe default ExecConfig.
func DefaultExecConfig() ExecConfig {
	return ExecConfig{
		Enabled: false,
		Command: "",
		Args:    nil,
		Timeout: 5 * time.Second,
		Shell:   false,
	}
}

// Validate checks that the ExecConfig is internally consistent.
func (c ExecConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Command == "" {
		return fmt.Errorf("exec notifier: command must not be empty when enabled")
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("exec notifier: timeout must be positive, got %s", c.Timeout)
	}
	return nil
}

// Merge fills zero-value fields in c with values from defaults.
func (c ExecConfig) Merge(defaults ExecConfig) ExecConfig {
	if c.Timeout == 0 {
		c.Timeout = defaults.Timeout
	}
	return c
}
