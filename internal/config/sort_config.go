package config

import "strings"

// SortConfig controls how listeners are ordered in output and reports.
type SortConfig struct {
	// Field is the primary sort field: "port", "protocol", "address", or "pid".
	// Defaults to "port".
	Field string `yaml:"field"`

	// Ascending controls sort direction. Defaults to true.
	Ascending bool `yaml:"ascending"`
}

// DefaultSortConfig returns a SortConfig with sensible defaults.
func DefaultSortConfig() SortConfig {
	return SortConfig{
		Field:     "port",
		Ascending: true,
	}
}

// Validate checks that the SortConfig fields are valid.
// It returns a non-nil error for unrecognized sort fields.
func (s *SortConfig) Validate() error {
	valid := map[string]bool{
		"port":     true,
		"protocol": true,
		"proto":    true,
		"address":  true,
		"addr":     true,
		"pid":      true,
	}
	norm := strings.ToLower(strings.TrimSpace(s.Field))
	if norm == "" {
		s.Field = "port"
		return nil
	}
	if !valid[norm] {
		return &ValidationError{Field: "sort.field", Message: "unrecognized sort field: " + s.Field}
	}
	s.Field = norm
	return nil
}

// Merge returns a SortConfig where zero values in s are replaced by defaults.
func (s SortConfig) Merge(defaults SortConfig) SortConfig {
	if strings.TrimSpace(s.Field) == "" {
		s.Field = defaults.Field
	}
	return s
}
