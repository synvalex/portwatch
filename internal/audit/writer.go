package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/ports"
	"github.com/user/portwatch/internal/snapshot"
)

// Event is a single audit log entry.
type Event struct {
	Timestamp time.Time        `json:"timestamp"`
	Type      snapshot.EventType `json:"type"`
	Protocol  string           `json:"protocol"`
	Address   string           `json:"address"`
	Port      uint16           `json:"port"`
	PID       int              `json:"pid,omitempty"`
	Process   string           `json:"process,omitempty"`
}

// Writer writes audit events to an io.Writer in the configured format.
type Writer struct {
	mu     sync.Mutex
	out    io.Writer
	format string // "json" or "text"
}

// NewWriter creates a Writer that writes to out using the given format.
func NewWriter(out io.Writer, format string) *Writer {
	return &Writer{out: out, format: format}
}

// NewFileWriter opens (or creates) the file at path and returns a Writer.
func NewFileWriter(path, format string) (*Writer, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o640)
	if err != nil {
		return nil, fmt.Errorf("audit: open %s: %w", path, err)
	}
	return NewWriter(f, format), nil
}

// Write records an audit event for the given listener and event type.
func (w *Writer) Write(l ports.Listener, et snapshot.EventType) error {
	ev := Event{
		Timestamp: time.Now().UTC(),
		Type:      et,
		Protocol:  l.Protocol,
		Address:   l.Address.String(),
		Port:      l.Port,
	}
	if l.Process != nil {
		ev.PID = l.Process.PID
		ev.Process = l.Process.Name
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	switch w.format {
	case "json":
		return json.NewEncoder(w.out).Encode(ev)
	default:
		_, err := fmt.Fprintf(w.out, "%s\t%s\t%s:%d\tpid=%d process=%s\n",
			ev.Timestamp.Format(time.RFC3339),
			ev.Type, ev.Address, ev.Port, ev.PID, ev.Process)
		return err
	}
}
