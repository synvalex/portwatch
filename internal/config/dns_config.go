package config

import "fmt"

// DNSConfig controls reverse-DNS enrichment for listener addresses.
type DNSConfig struct {
	Enabled    bool          `yaml:"enabled"`
	CacheSize  int           `yaml:"cache_size"`
	Timeout    Duration      `yaml:"timeout"`
	Workers    int           `yaml:"workers"`
}

// DefaultDNSConfig returns conservative defaults.
func DefaultDNSConfig() DNSConfig {
	return DNSConfig{
		Enabled:   false,
		CacheSize: 512,
		Timeout:   Duration(2_000_000_000), // 2s
		Workers:   4,
	}
}

// Validate checks that all fields are sane.
func (c DNSConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.CacheSize <= 0 {
		return fmt.Errorf("dns.cache_size must be > 0, got %d", c.CacheSize)
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("dns.timeout must be > 0")
	}
	if c.Workers <= 0 {
		return fmt.Errorf("dns.workers must be > 0, got %d", c.Workers)
	}
	return nil
}

// Merge fills zero-value fields from defaults.
func (c *DNSConfig) Merge(def DNSConfig) {
	if c.CacheSize == 0 {
		c.CacheSize = def.CacheSize
	}
	if c.Timeout == 0 {
		c.Timeout = def.Timeout
	}
	if c.Workers == 0 {
		c.Workers = def.Workers
	}
}
