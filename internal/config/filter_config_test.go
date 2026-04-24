package config

import (
	"testing"
)

func TestDefaultFilterConfig(t *testing.T) {
	fc := DefaultFilterConfig()
	if !fc.ExcludeLoopback {
		t.Error("expected ExcludeLoopback to be true by default")
	}
	if fc.ExcludePorts == nil {
		t.Error("expected ExcludePorts to be non-nil slice")
	}
	if len(fc.ExcludePorts) != 0 {
		t.Errorf("expected empty ExcludePorts, got %v", fc.ExcludePorts)
	}
}

func TestFilterConfig_Merge_PortsFromDefault(t *testing.T) {
	defaults := FilterConfig{
		ExcludeLoopback: true,
		ExcludePorts:    []uint16{22, 80},
	}
	userCfg := FilterConfig{
		ExcludeLoopback: false,
		ExcludePorts:    []uint16{},
	}
	result := userCfg.Merge(defaults)
	if len(result.ExcludePorts) != 2 {
		t.Errorf("expected 2 excluded ports from defaults, got %d", len(result.ExcludePorts))
	}
	if result.ExcludeLoopback {
		t.Error("expected ExcludeLoopback to remain false from user config")
	}
}

func TestFilterConfig_Merge_UserPortsPreserved(t *testing.T) {
	defaults := FilterConfig{
		ExcludeLoopback: true,
		ExcludePorts:    []uint16{22},
	}
	userCfg := FilterConfig{
		ExcludeLoopback: true,
		ExcludePorts:    []uint16{443, 8080},
	}
	result := userCfg.Merge(defaults)
	if len(result.ExcludePorts) != 2 {
		t.Errorf("expected user's 2 ports to be preserved, got %d", len(result.ExcludePorts))
	}
	if result.ExcludePorts[0] != 443 {
		t.Errorf("expected first port 443, got %d", result.ExcludePorts[0])
	}
}
