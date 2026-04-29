package config

import "fmt"

// SyslogConfig holds configuration for the syslog notifier.
type SyslogConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Network  string `yaml:"network"`  // "tcp", "udp", or "" for local
	Addr     string `yaml:"addr"`     // host:port for remote syslog
	Tag      string `yaml:"tag"`
	Facility string `yaml:"facility"` // "local0".."local7", "daemon", etc.
}

// DefaultSyslogConfig returns a safe default syslog configuration.
func DefaultSyslogConfig() SyslogConfig {
	return SyslogConfig{
		Enabled:  false,
		Network:  "",
		Addr:     "",
		Tag:      "portwatch",
		Facility: "daemon",
	}
}

// Validate checks the syslog configuration for correctness.
func (c SyslogConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Tag == "" {
		return fmt.Errorf("syslog: tag must not be empty")
	}
	validFacilities := map[string]bool{
		"daemon": true, "local0": true, "local1": true,
		"local2": true, "local3": true, "local4": true,
		"local5": true, "local6": true, "local7": true,
	}
	if !validFacilities[c.Facility] {
		return fmt.Errorf("syslog: unknown facility %q", c.Facility)
	}
	if c.Network != "" && c.Addr == "" {
		return fmt.Errorf("syslog: addr must be set when network is %q", c.Network)
	}
	return nil
}

// Merge fills empty fields from defaults.
func (c SyslogConfig) Merge(defaults SyslogConfig) SyslogConfig {
	if c.Tag == "" {
		c.Tag = defaults.Tag
	}
	if c.Facility == "" {
		c.Facility = defaults.Facility
	}
	return c
}
