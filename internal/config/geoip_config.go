package config

import "fmt"

// GeoIPConfig controls optional GeoIP enrichment of listener process info.
type GeoIPConfig struct {
	Enabled  bool   `yaml:"enabled"`
	DBPath   string `yaml:"db_path"`
	CacheSize int   `yaml:"cache_size"`
}

// DefaultGeoIPConfig returns a GeoIPConfig with sensible defaults.
func DefaultGeoIPConfig() GeoIPConfig {
	return GeoIPConfig{
		Enabled:   false,
		DBPath:    "",
		CacheSize: 512,
	}
}

// Validate returns an error if the GeoIPConfig is invalid.
func (c GeoIPConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.DBPath == "" {
		return fmt.Errorf("geoip: db_path must be set when enabled")
	}
	if c.CacheSize < 0 {
		return fmt.Errorf("geoip: cache_size must be non-negative, got %d", c.CacheSize)
	}
	return nil
}

// Merge returns a new GeoIPConfig where any zero-value fields in c are filled
// from defaults.
func (c GeoIPConfig) Merge(defaults GeoIPConfig) GeoIPConfig {
	if c.CacheSize == 0 {
		c.CacheSize = defaults.CacheSize
	}
	return c
}
