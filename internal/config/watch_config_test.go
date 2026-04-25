package config

import "testing"

func TestDefaultWatchConfig(t *testing.T) {
	cfg := DefaultWatchConfig()
	if !cfg.Enabled {
		t.Error("expected Enabled to be true by default")
	}
	if cfg.DebounceMs != 500 {
		t.Errorf("expected DebounceMs=500, got %d", cfg.DebounceMs)
	}
}

func TestWatchConfig_Validate_Valid(t *testing.T) {
	cfg := DefaultWatchConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestWatchConfig_Validate_NegativeDebounce(t *testing.T) {
	cfg := WatchConfig{Enabled: true, DebounceMs: -1}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for negative debounce_ms")
	}
}

func TestWatchConfig_Validate_TooLargeDebounce(t *testing.T) {
	cfg := WatchConfig{Enabled: true, DebounceMs: 99_999}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for debounce_ms > 60000")
	}
}

func TestWatchConfig_Validate_ZeroDebounce(t *testing.T) {
	// Zero is explicitly allowed (disables debouncing).
	cfg := WatchConfig{Enabled: false, DebounceMs: 0}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error for zero debounce_ms: %v", err)
	}
}

func TestWatchConfig_Merge_ZeroFilledFromDefaults(t *testing.T) {
	user := WatchConfig{Enabled: false, DebounceMs: 0}
	defaults := DefaultWatchConfig()
	merged := user.Merge(defaults)
	if merged.DebounceMs != defaults.DebounceMs {
		t.Errorf("expected DebounceMs=%d after merge, got %d",
			defaults.DebounceMs, merged.DebounceMs)
	}
}

func TestWatchConfig_Merge_UserValuePreserved(t *testing.T) {
	user := WatchConfig{Enabled: true, DebounceMs: 250}
	defaults := DefaultWatchConfig()
	merged := user.Merge(defaults)
	if merged.DebounceMs != 250 {
		t.Errorf("expected user DebounceMs=250 to be preserved, got %d", merged.DebounceMs)
	}
}
