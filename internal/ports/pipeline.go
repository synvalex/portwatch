package ports

// Pipeline composes scanning, enrichment, filtering, and sorting into a single
// reusable processing step that turns raw kernel listeners into a clean,
// enriched, sorted slice ready for display or alerting.
type Pipeline struct {
	scanner  Scanner
	enricher *Enricher
	chain    *FilterChain
	sortOpts SortOptions
}

// Scanner is the interface satisfied by the platform-specific scanner.
type Scanner interface {
	Listeners() ([]Listener, error)
}

// SortOptions carries the field and direction used by SortListeners.
type SortOptions struct {
	Field     SortField
	Ascending bool
}

// NewPipeline constructs a Pipeline from its constituent parts.
// Any of enricher, chain, or sortOpts may be zero-valued to skip that stage.
func NewPipeline(s Scanner, e *Enricher, c *FilterChain, opts SortOptions) *Pipeline {
	return &Pipeline{
		scanner:  s,
		enricher: e,
		chain:    c,
		sortOpts: opts,
	}
}

// Run executes every stage in order and returns the processed listeners.
func (p *Pipeline) Run() ([]Listener, error) {
	listeners, err := p.scanner.Listeners()
	if err != nil {
		return nil, err
	}

	if p.enricher != nil {
		listeners = p.enricher.Enrich(listeners)
	}

	if p.chain != nil {
		listeners = p.chain.Apply(listeners)
	}

	if p.sortOpts.Field != "" {
		SortListeners(listeners, p.sortOpts.Field, p.sortOpts.Ascending)
	}

	return listeners, nil
}
