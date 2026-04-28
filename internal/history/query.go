package history

import (
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Query holds filter parameters for searching history events.
type Query struct {
	// Since filters events to those recorded at or after this time.
	// Zero value means no lower bound.
	Since time.Time

	// EventType filters by event type. Empty string means all types.
	EventType alert.EventType

	// Port filters by listener port. 0 means all ports.
	Port uint16

	// Limit caps the number of results. 0 means use store default.
	Limit int
}

// Search returns events from the store that match the given query.
func (s *Store) Search(q Query) []alert.Event {
	limit := q.Limit
	if limit <= 0 {
		limit = 100
	}

	all := s.Recent(limit * 10) // fetch extra to allow filtering

	var results []alert.Event
	for _, ev := range all {
		if !q.Since.IsZero() && ev.Time.Before(q.Since) {
			continue
		}
		if q.EventType != "" && ev.Type != q.EventType {
			continue
		}
		if q.Port != 0 && ev.Listener.Address.Port != q.Port {
			continue
		}
		results = append(results, ev)
		if len(results) >= limit {
			break
		}
	}
	return results
}
