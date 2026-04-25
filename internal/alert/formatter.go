package alert

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/yourorg/portwatch/internal/config"
	"github.com/yourorg/portwatch/internal/ports"
)

// Event is a structured alert event emitted by the dispatcher.
type Event struct {
	Timestamp time.Time      `json:"timestamp"`
	Kind      string         `json:"kind"`
	Listener  ports.Listener `json:"listener"`
	Rule      string         `json:"rule,omitempty"`
}

// Formatter converts an Event to a human- or machine-readable string.
type Formatter interface {
	Format(e Event) string
}

// NewFormatter returns a Formatter for the given output config.
func NewFormatter(cfg config.OutputConfig) Formatter {
	if cfg.Format == config.OutputFormatJSON {
		return &jsonFormatter{}
	}
	return &textFormatter{color: cfg.Color, timestamps: cfg.Timestamps}
}

type textFormatter struct {
	color      bool
	timestamps bool
}

func (f *textFormatter) Format(e Event) string {
	ts := ""
	if f.timestamps {
		ts = e.Timestamp.Format(time.RFC3339) + " "
	}
	rule := ""
	if e.Rule != "" {
		rule = fmt.Sprintf(" [rule:%s]", e.Rule)
	}
	msg := fmt.Sprintf("%s%s %s%s", ts, e.Kind, e.Listener, rule)
	if f.color {
		switch e.Kind {
		case "APPEARED":
			return "\033[33m" + msg + "\033[0m"
		case "DENIED":
			return "\033[31m" + msg + "\033[0m"
		case "ALLOWED":
			return "\033[32m" + msg + "\033[0m"
		}
	}
	return msg
}

type jsonFormatter struct{}

func (f *jsonFormatter) Format(e Event) string {
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf(`{"error":%q}`, err.Error())
	}
	return string(b)
}
