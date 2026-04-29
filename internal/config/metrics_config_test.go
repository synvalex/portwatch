package config

import "testing"

func TestDefaultMetricsConfig(t *testing.T) {
	c := DefaultMetricsConfig()
	if c.Enabled {
		t.Error("expected metrics disabled by default")
	}
	if c.Addr != ":9090" {
		t.Errorf("unexpected default addr: %s", c.Addr)
	}
	if c.Path != "/metrics" {
		t.Errorf("unexpected default path: %s", c.Path)
	}
}

func TestMetricsConfig_Validate_Disabled(t *testing.T) {
	c := MetricsConfig{Enabled: false}
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error for disabled metrics, got: %v", err)
	}
}

func TestMetricsConfig_Validate_EnabledValid(t *testing.T) {
	c := MetricsConfig{Enabled: true, Addr: ":9090", Path: "/metrics"}
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestMetricsConfig_Validate_EnabledMissingAddr(t *testing.T) {
	c := MetricsConfig{Enabled: true, Addr: "", Path: "/metrics"}
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing addr")
	}
}

func TestMetricsConfig_Validate_EnabledMissingPath(t *testing.T) {
	c := MetricsConfig{Enabled: true, Addr: ":9090", Path: ""}
	if err := c.Validate(); err == nil {
		t.Error("expected error for missing path")
	}
}

func TestMetricsConfig_Merge_FillsDefaults(t *testing.T) {
	c := MetricsConfig{Enabled: true}
	c.Merge(DefaultMetricsConfig())
	if c.Addr != ":9090" {
		t.Errorf("expected addr filled from default, got: %s", c.Addr)
	}
	if c.Path != "/metrics" {
		t.Errorf("expected path filled from default, got: %s", c.Path)
	}
}

func TestMetricsConfig_Merge_UserValuesPreserved(t *testing.T) {
	c := MetricsConfig{Enabled: true, Addr: ":2112", Path: "/prom"}
	c.Merge(DefaultMetricsConfig())
	if c.Addr != ":2112" {
		t.Errorf("expected user addr preserved, got: %s", c.Addr)
	}
	if c.Path != "/prom" {
		t.Errorf("expected user path preserved, got: %s", c.Path)
	}
}
