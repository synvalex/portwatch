package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level portwatch configuration.
type Config struct {
	Interval time.Duration `yaml:"interval"`
	LogLevel string        `yaml:"log_level"`
	Rules    []RuleConfig  `yaml:"rules"`
}

// RuleConfig represents a single rule entry in the config file.
type RuleConfig struct {
	Name    string `yaml:"name"`
	Port    int    `yaml:"port"`
	Proto   string `yaml:"proto"`
	Address string `yaml:"address"`
	Action  string `yaml:"action"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Interval: 15 * time.Second,
		LogLevel: "info",
		Rules:    []RuleConfig{},
	}
}

// Load reads and parses a YAML config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}
	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return cfg, nil
}

// Validate checks that the configuration values are acceptable.
func (c *Config) Validate() error {
	if c.Interval < time.Second {
		return fmt.Errorf("interval must be at least 1s, got %s", c.Interval)
	}
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[c.LogLevel] {
		return fmt.Errorf("unknown log_level %q", c.LogLevel)
	}
	for i, r := range c.Rules {
		if r.Name == "" {
			return fmt.Errorf("rule[%d]: name is required", i)
		}
		if r.Action != "allow" && r.Action != "deny" {
			return fmt.Errorf("rule[%d] %q: action must be 'allow' or 'deny'", i, r.Name)
		}
	}
	return nil
}
