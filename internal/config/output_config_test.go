package config

import (
	"testing"
)

func TestDefaultOutputConfig(t *testing.T) {
	cfg := DefaultOutputConfig()
	if cfg.Format != OutputFormatText {
		t.Errorf("expected format %q, got %q", OutputFormatText, cfg.Format)
	}
	if !cfg.Timestamps {
		t.Error("expected Timestamps to be true by default")
	}
	if !cfg.Color {
		t.Error("expected Color to be true by default")
	}
}

func TestOutputConfig_Validate_Valid(t *testing.T) {
	for _, f := range []OutputFormat{OutputFormatText, OutputFormatJSON} {
		cfg := OutputConfig{Format: f}
		if err := cfg.Validate(); err != nil {
			t.Errorf("format %q should be valid, got error: %v", f, err)
		}
	}
}

func TestOutputConfig_Validate_Invalid(t *testing.T) {
	cfg := OutputConfig{Format: "yaml"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}

func TestOutputConfig_Merge_EmptyFormatFilled(t *testing.T) {
	defaults := DefaultOutputConfig()
	user := OutputConfig{Timestamps: false, Color: false}
	merged := user.Merge(defaults)
	if merged.Format != OutputFormatText {
		t.Errorf("expected merged format %q, got %q", OutputFormatText, merged.Format)
	}
}

func TestOutputConfig_Merge_UserFormatPreserved(t *testing.T) {
	defaults := DefaultOutputConfig()
	user := OutputConfig{Format: OutputFormatJSON}
	merged := user.Merge(defaults)
	if merged.Format != OutputFormatJSON {
		t.Errorf("expected user format %q to be preserved, got %q", OutputFormatJSON, merged.Format)
	}
}
