package ports

// FilterChain composes multiple Filters into a single Filter that applies
// each in order. A listener is kept only if ALL filters accept it.
type FilterChain struct {
	filters []*Filter
}

// NewFilterChain creates a FilterChain from the provided filters.
// Nil filters are silently skipped.
func NewFilterChain(filters ...*Filter) *FilterChain {
	var valid []*Filter
	for _, f := range filters {
		if f != nil {
			valid = append(valid, f)
		}
	}
	return &FilterChain{filters: valid}
}

// Apply returns only the listeners accepted by every filter in the chain.
// If the chain is empty it acts as a pass-through.
func (fc *FilterChain) Apply(listeners []Listener) []Listener {
	if len(fc.filters) == 0 {
		return listeners
	}
	out := make([]Listener, 0, len(listeners))
	for _, l := range listeners {
		if fc.accept(l) {
			out = append(out, l)
		}
	}
	return out
}

// accept returns true only when every filter accepts the listener.
func (fc *FilterChain) accept(l Listener) bool {
	for _, f := range fc.filters {
		if !f.Accept(l) {
			return false
		}
	}
	return true
}

// Len returns the number of filters in the chain.
func (fc *FilterChain) Len() int { return len(fc.filters) }
