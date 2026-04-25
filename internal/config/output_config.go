package config

import "fmt"

// OutputFormat controls how alerts and scan results are rendered.
type OutputFormat string

const (
	OutputFormatText OutputFormat = "text"
	OutputFormatJSON OutputFormat = "json"
)

// OutputConfig holds rendering and output preferences.
type OutputConfig struct {
	// Format is the output format: "text" or "json".
	Format OutputFormat `yaml:"format"`
	// Timestamps controls whether log lines include timestamps.
	Timestamps bool `yaml:"timestamps"`
	// Color enables ANSI color codes in text output.
	Color bool `yaml:"color"`
}

// DefaultOutputConfig returns sensible output defaults.
func DefaultOutputConfig() OutputConfig {
	return OutputConfig{
		Format:     OutputFormatText,
		Timestamps: true,
		Color:      true,
	}
}

// Merge returns a new OutputConfig where any zero-value fields in o
// are filled from defaults.
func (o OutputConfig) Merge(defaults OutputConfig) OutputConfig {
	if o.Format == "" {
		o.Format = defaults.Format
	}
	return o
}

// Validate checks that the OutputConfig holds valid values.
func (o OutputConfig) Validate() error {
	switch o.Format {
	case OutputFormatText, OutputFormatJSON:
		return nil
	default:
		return fmt.Errorf("unknown output format %q: must be \"text\" or \"json\"", o.Format)
	}
}
