package config

import (
	"testing"
	"time"
)

func TestDefaultEmailConfig(t *testing.T) {
	cfg := DefaultEmailConfig()
	if cfg.Enabled {
		t.Error("expected Enabled to be false by default")
	}
	if cfg.SMTPPort != 587 {
		t.Errorf("expected default SMTPPort 587, got %d", cfg.SMTPPort)
	}
	if cfg.Subject == "" {
		t.Error("expected non-empty default Subject")
	}
	if cfg.Timeout != 10*time.Second {
		t.Errorf("expected default Timeout 10s, got %v", cfg.Timeout)
	}
}

func TestEmailConfig_Validate_DisabledNoFields(t *testing.T) {
	cfg := EmailConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for disabled config, got %v", err)
	}
}

func TestEmailConfig_Validate_EnabledValid(t *testing.T) {
	cfg := EmailConfig{
		Enabled:  true,
		SMTPHost: "smtp.example.com",
		SMTPPort: 587,
		From:     "alert@example.com",
		To:       []string{"admin@example.com"},
		Timeout:  5 * time.Second,
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for valid config, got %v", err)
	}
}

func TestEmailConfig_Validate_EnabledMissingHost(t *testing.T) {
	cfg := EmailConfig{
		Enabled:  true,
		SMTPPort: 587,
		From:     "alert@example.com",
		To:       []string{"admin@example.com"},
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing smtp_host")
	}
}

func TestEmailConfig_Validate_EnabledNoRecipients(t *testing.T) {
	cfg := EmailConfig{
		Enabled:  true,
		SMTPHost: "smtp.example.com",
		SMTPPort: 587,
		From:     "alert@example.com",
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing recipients")
	}
}

func TestEmailConfig_Validate_NegativeTimeout(t *testing.T) {
	cfg := EmailConfig{
		Enabled:  true,
		SMTPHost: "smtp.example.com",
		SMTPPort: 587,
		From:     "alert@example.com",
		To:       []string{"admin@example.com"},
		Timeout:  -1 * time.Second,
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for negative timeout")
	}
}

func TestEmailConfig_Merge_FillsDefaults(t *testing.T) {
	defaults := DefaultEmailConfig()
	user := EmailConfig{Enabled: true, SMTPHost: "smtp.example.com"}
	merged := user.Merge(defaults)
	if merged.SMTPPort != 587 {
		t.Errorf("expected SMTPPort 587 from defaults, got %d", merged.SMTPPort)
	}
	if merged.Subject == "" {
		t.Error("expected Subject filled from defaults")
	}
	if merged.Timeout != 10*time.Second {
		t.Errorf("expected Timeout from defaults, got %v", merged.Timeout)
	}
}

func TestEmailConfig_Merge_UserValuesPreserved(t *testing.T) {
	defaults := DefaultEmailConfig()
	user := EmailConfig{
		Enabled:  true,
		SMTPPort: 465,
		Subject:  "custom subject",
		Timeout:  3 * time.Second,
	}
	merged := user.Merge(defaults)
	if merged.SMTPPort != 465 {
		t.Errorf("expected user SMTPPort 465, got %d", merged.SMTPPort)
	}
	if merged.Subject != "custom subject" {
		t.Errorf("expected user Subject preserved, got %q", merged.Subject)
	}
	if merged.Timeout != 3*time.Second {
		t.Errorf("expected user Timeout 3s, got %v", merged.Timeout)
	}
}
