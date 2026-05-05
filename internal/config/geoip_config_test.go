package config_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/config"
)

func TestDefaultGeoIPConfig(t *testing.T) {
	cfg := config.DefaultGeoIPConfig()
	if cfg.Enabled {
		t.Error("expected Enabled=false by default")
	}
	if cfg.DBPath != "" {
		t.Errorf("expected empty DBPath, got %q", cfg.DBPath)
	}
	if cfg.CacheSize != 512 {
		t.Errorf("expected CacheSize=512, got %d", cfg.CacheSize)
	}
}

func TestGeoIPConfig_Validate_Disabled(t *testing.T) {
	cfg := config.GeoIPConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got %v", err)
	}
}

func TestGeoIPConfig_Validate_EnabledValid(t *testing.T) {
	cfg := config.GeoIPConfig{Enabled: true, DBPath: "/var/lib/GeoIP/GeoLite2-City.mmdb", CacheSize: 256}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestGeoIPConfig_Validate_EnabledMissingPath(t *testing.T) {
	cfg := config.GeoIPConfig{Enabled: true, DBPath: "", CacheSize: 256}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing db_path")
	}
}

func TestGeoIPConfig_Validate_NegativeCacheSize(t *testing.T) {
	cfg := config.GeoIPConfig{Enabled: true, DBPath: "/some/path.mmdb", CacheSize: -1}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for negative cache_size")
	}
}

func TestGeoIPConfig_Merge_FillsCacheSize(t *testing.T) {
	defaults := config.DefaultGeoIPConfig()
	user := config.GeoIPConfig{Enabled: true, DBPath: "/path.mmdb", CacheSize: 0}
	merged := user.Merge(defaults)
	if merged.CacheSize != defaults.CacheSize {
		t.Errorf("expected CacheSize=%d, got %d", defaults.CacheSize, merged.CacheSize)
	}
}

func TestGeoIPConfig_Merge_UserCacheSizePreserved(t *testing.T) {
	defaults := config.DefaultGeoIPConfig()
	user := config.GeoIPConfig{Enabled: true, DBPath: "/path.mmdb", CacheSize: 1024}
	merged := user.Merge(defaults)
	if merged.CacheSize != 1024 {
		t.Errorf("expected CacheSize=1024, got %d", merged.CacheSize)
	}
}
