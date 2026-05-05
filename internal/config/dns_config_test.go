package config

import (
	"testing"
)

func TestDefaultDNSConfig(t *testing.T) {
	d := DefaultDNSConfig()
	if d.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if d.CacheSize != 512 {
		t.Errorf("expected CacheSize=512, got %d", d.CacheSize)
	}
	if d.Workers != 4 {
		t.Errorf("expected Workers=4, got %d", d.Workers)
	}
	if d.Timeout <= 0 {
		t.Error("expected positive Timeout")
	}
}

func TestDNSConfig_Validate_Disabled(t *testing.T) {
	c := DNSConfig{Enabled: false}
	if err := c.Validate(); err != nil {
		t.Errorf("disabled config should always validate, got %v", err)
	}
}

func TestDNSConfig_Validate_EnabledValid(t *testing.T) {
	c := DefaultDNSConfig()
	c.Enabled = true
	if err := c.Validate(); err != nil {
		t.Errorf("valid config should not error, got %v", err)
	}
}

func TestDNSConfig_Validate_ZeroCacheSize(t *testing.T) {
	c := DefaultDNSConfig()
	c.Enabled = true
	c.CacheSize = 0
	if err := c.Validate(); err == nil {
		t.Error("expected error for zero cache_size")
	}
}

func TestDNSConfig_Validate_ZeroWorkers(t *testing.T) {
	c := DefaultDNSConfig()
	c.Enabled = true
	c.Workers = 0
	if err := c.Validate(); err == nil {
		t.Error("expected error for zero workers")
	}
}

func TestDNSConfig_Merge_FillsDefaults(t *testing.T) {
	def := DefaultDNSConfig()
	c := DNSConfig{Enabled: true}
	c.Merge(def)
	if c.CacheSize != def.CacheSize {
		t.Errorf("expected CacheSize=%d, got %d", def.CacheSize, c.CacheSize)
	}
	if c.Workers != def.Workers {
		t.Errorf("expected Workers=%d, got %d", def.Workers, c.Workers)
	}
}

func TestDNSConfig_Merge_UserValuesPreserved(t *testing.T) {
	def := DefaultDNSConfig()
	c := DNSConfig{Enabled: true, CacheSize: 1024, Workers: 8, Timeout: 5_000_000_000}
	c.Merge(def)
	if c.CacheSize != 1024 {
		t.Errorf("user CacheSize should be preserved, got %d", c.CacheSize)
	}
	if c.Workers != 8 {
		t.Errorf("user Workers should be preserved, got %d", c.Workers)
	}
}
