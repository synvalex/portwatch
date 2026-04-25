package config

import (
	"testing"
)

func TestDefaultAlertConfig(t *testing.T) {
	c := DefaultAlertConfig()
	if !c.Enabled {
		t.Error("expected Enabled to be true by default")
	}
	if c.LogLevel != "warn" {
		t.Errorf("expected LogLevel \"warn\", got %q", c.LogLevel)
	}
	if !c.IncludeProcess {
		t.Error("expected IncludeProcess to be true by default")
	}
	if c.DedupWindow != 60 {
		t.Errorf("expected DedupWindow 60, got %d", c.DedupWindow)
	}
}

func TestAlertConfig_Validate_Valid(t *testing.T) {
	for _, level := range []string{"info", "warn", "error"} {
		c := DefaultAlertConfig()
		c.LogLevel = level
		if err := c.Validate(); err != nil {
			t.Errorf("expected no error for level %q, got: %v", level, err)
		}
	}
}

func TestAlertConfig_Validate_InvalidLevel(t *testing.T) {
	c := DefaultAlertConfig()
	c.LogLevel = "debug"
	if err := c.Validate(); err == nil {
		t.Error("expected error for invalid log level, got nil")
	}
}

func TestAlertConfig_Validate_NegativeDedupWindow(t *testing.T) {
	c := DefaultAlertConfig()
	c.DedupWindow = -1
	if err := c.Validate(); err == nil {
		t.Error("expected error for negative dedup window, got nil")
	}
}

func TestAlertConfig_Validate_ZeroDedupWindow(t *testing.T) {
	c := DefaultAlertConfig()
	c.DedupWindow = 0
	if err := c.Validate(); err != nil {
		t.Errorf("expected no error for zero dedup window, got: %v", err)
	}
}

func TestAlertConfig_Merge_EmptyLevelFilled(t *testing.T) {
	user := AlertConfig{Enabled: true, IncludeProcess: false}
	d := DefaultAlertConfig()
	merged := user.Merge(d)
	if merged.LogLevel != d.LogLevel {
		t.Errorf("expected merged LogLevel %q, got %q", d.LogLevel, merged.LogLevel)
	}
	if merged.DedupWindow != d.DedupWindow {
		t.Errorf("expected merged DedupWindow %d, got %d", d.DedupWindow, merged.DedupWindow)
	}
}

func TestAlertConfig_Merge_UserValuesPreserved(t *testing.T) {
	user := AlertConfig{LogLevel: "error", DedupWindow: 120}
	merged := user.Merge(DefaultAlertConfig())
	if merged.LogLevel != "error" {
		t.Errorf("expected user LogLevel \"error\", got %q", merged.LogLevel)
	}
	if merged.DedupWindow != 120 {
		t.Errorf("expected user DedupWindow 120, got %d", merged.DedupWindow)
	}
}
