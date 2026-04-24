package ports

// Enricher attaches ProcessInfo to Listener values where possible.
type Enricher struct {
	lookup func(inode uint64) (*ProcessInfo, error)
}

// NewEnricher creates an Enricher using the default LookupInode implementation.
func NewEnricher() *Enricher {
	return &Enricher{lookup: LookupInode}
}

// newEnricherWithLookup creates an Enricher with a custom lookup function (for testing).
func newEnricherWithLookup(fn func(uint64) (*ProcessInfo, error)) *Enricher {
	return &Enricher{lookup: fn}
}

// Enrich attempts to resolve the owning process for each listener.
// Listeners without inode information or where lookup fails are returned unchanged.
func (e *Enricher) Enrich(listeners []Listener) []Listener {
	result := make([]Listener, len(listeners))
	for i, l := range listeners {
		if l.Inode == 0 {
			result[i] = l
			continue
		}
		info, err := e.lookup(l.Inode)
		if err != nil || info == nil {
			result[i] = l
			continue
		}
		l.Process = info
		result[i] = l
	}
	return result
}
