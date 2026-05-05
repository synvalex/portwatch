package config_test

import (
	"testing"

	"github.com/example/portwatch/internal/config"
)

func TestDefaultFingerprintConfig(t *testing.T) {
	cfg := config.DefaultFingerprintConfig()
	if !cfg.Enabled {
		t.Error("expected Enabled to be true")
	}
	if !cfg.HashExecutable {
		t.Error("expected HashExecutable to be true")
	}
	if cfg.TrustCmdline {
		t.Error("expected TrustCmdline to be false")
	}
}

func TestFingerprintConfig_Validate_Disabled(t *testing.T) {
	cfg := config.FingerprintConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFingerprintConfig_Validate_Valid(t *testing.T) {
	cfg := config.FingerprintConfig{
		Enabled:        true,
		HashExecutable: true,
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFingerprintConfig_Validate_NeitherOption(t *testing.T) {
	cfg := config.FingerprintConfig{
		Enabled:        true,
		HashExecutable: false,
		TrustCmdline:   false,
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when neither hash_executable nor trust_cmdline is set")
	}
}

func TestFingerprintConfig_Merge_FillsAllowList(t *testing.T) {
	defaults := config.FingerprintConfig{AllowList: []string{"/usr/sbin/sshd"}}
	cfg := config.FingerprintConfig{}
	merged := cfg.Merge(defaults)
	if len(merged.AllowList) != 1 || merged.AllowList[0] != "/usr/sbin/sshd" {
		t.Errorf("expected allow list from defaults, got %v", merged.AllowList)
	}
}

func TestFingerprintConfig_Merge_UserAllowListPreserved(t *testing.T) {
	defaults := config.FingerprintConfig{AllowList: []string{"/usr/sbin/sshd"}}
	cfg := config.FingerprintConfig{AllowList: []string{"/usr/bin/nginx"}}
	merged := cfg.Merge(defaults)
	if len(merged.AllowList) != 1 || merged.AllowList[0] != "/usr/bin/nginx" {
		t.Errorf("expected user allow list preserved, got %v", merged.AllowList)
	}
}
