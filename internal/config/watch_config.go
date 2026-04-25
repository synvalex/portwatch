package config

import "fmt"

// WatchConfig controls the file-watching / inotify behaviour used
// when portwatch reloads its configuration at runtime.
type WatchConfig struct {
	// Enabled turns live config-reload on or off.
	Enabled bool `yaml:"enabled"`

	// DebounceMs is the number of milliseconds to wait after the last
	// filesystem event before actually reloading the config file.
	// This prevents rapid successive reloads when editors write files
	// in multiple steps.
	DebounceMs int `yaml:"debounce_ms"`
}

// DefaultWatchConfig returns a WatchConfig with sensible defaults.
func DefaultWatchConfig() WatchConfig {
	return WatchConfig{
		Enabled:    true,
		DebounceMs: 500,
	}
}

// Validate returns an error if the WatchConfig contains invalid values.
func (w WatchConfig) Validate() error {
	if w.DebounceMs < 0 {
		return fmt.Errorf("watch.debounce_ms must be >= 0, got %d", w.DebounceMs)
	}
	if w.DebounceMs > 60_000 {
		return fmt.Errorf("watch.debounce_ms must be <= 60000, got %d", w.DebounceMs)
	}
	return nil
}

// Merge returns a new WatchConfig where any zero-value fields in w are
// filled from defaults.
func (w WatchConfig) Merge(defaults WatchConfig) WatchConfig {
	if w.DebounceMs == 0 {
		w.DebounceMs = defaults.DebounceMs
	}
	return w
}
