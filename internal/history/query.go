package history

import (
	"time"

	"github.com/yourorg/portwatch/internal/ports"
)

// EventType represents the kind of change observed for a listener.
type EventType string

const (
	EventAppeared    EventType = "appeared"
	EventDisappeared EventType = "disappeared"
)

// Event records a single port change observed by the monitor.
type Event struct {
	Listener  ports.Listener
	EventType EventType
	OccurredAt time.Time
}

// Query provides read access to historical port events.
type Query struct {
	store *Store
}

// NewQuery creates a Query backed by the given Store.
func NewQuery(s *Store) *Query {
	return &Query{store: s}
}

// Recent returns all events that occurred within the given duration.
func (q *Query) Recent(window time.Duration) []Event {
	return q.store.Recent(window)
}

// Since returns all events that occurred after the given time.
func (q *Query) Since(t time.Time) []Event {
	window := time.Since(t)
	if window <= 0 {
		return nil
	}
	return q.store.Recent(window)
}

// ByPort filters events to those matching the given port number.
func (q *Query) ByPort(port uint16) []Event {
	all := q.store.Recent(q.store.retention)
	var out []Event
	for _, e := range all {
		if e.Listener.Port == port {
			out = append(out, e)
		}
	}
	return out
}

// ByEventType filters events to those matching the given EventType.
func (q *Query) ByEventType(et EventType) []Event {
	all := q.store.Recent(q.store.retention)
	var out []Event
	for _, e := range all {
		if e.EventType == et {
			out = append(out, e)
		}
	}
	return out
}
