package config

import (
	"fmt"

	"github.com/user/portwatch/internal/ports"
)

// PipelineConfig groups all settings that influence the scanning pipeline.
type PipelineConfig struct {
	Filter FilterConfig `yaml:"filter"`
	Sort   SortConfig   `yaml:"sort"`
	Output OutputConfig `yaml:"output"`
}

// DefaultPipelineConfig returns a PipelineConfig populated with safe defaults.
func DefaultPipelineConfig() PipelineConfig {
	return PipelineConfig{
		Filter: DefaultFilterConfig(),
		Sort:   DefaultSortConfig(),
		Output: DefaultOutputConfig(),
	}
}

// Validate checks that each sub-config is valid.
func (p PipelineConfig) Validate() error {
	if err := p.Output.Validate(); err != nil {
		return fmt.Errorf("output: %w", err)
	}
	return nil
}

// Merge fills zero-value fields from defaults, delegating to each sub-config.
func (p PipelineConfig) Merge(defaults PipelineConfig) PipelineConfig {
	p.Filter = p.Filter.Merge(defaults.Filter)
	p.Sort = p.Sort.Merge(defaults.Sort)
	p.Output = p.Output.Merge(defaults.Output)
	return p
}

// SortOptions converts the SortConfig into the ports.SortOptions type.
func (p PipelineConfig) SortOptions() (ports.SortOptions, error) {
	field, err := ports.ParseSortField(p.Sort.Field)
	if err != nil {
		return ports.SortOptions{}, fmt.Errorf("sort field: %w", err)
	}
	return ports.SortOptions{
		Field:     field,
		Ascending: p.Sort.Ascending,
	}, nil
}
