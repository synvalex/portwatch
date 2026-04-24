package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Alert represents a detected port event.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Listener  ports.Listener
}

// String returns a human-readable representation of the alert.
func (a Alert) String() string {
	return fmt.Sprintf("[%s] %s %s — %s",
		a.Timestamp.Format(time.RFC3339),
		a.Level,
		a.Listener.String(),
		a.Message,
	)
}

// Notifier is the interface for alert output backends.
type Notifier interface {
	Notify(a Alert) error
}

// LogNotifier writes alerts as plain text lines to a writer.
type LogNotifier struct {
	Out io.Writer
}

// NewLogNotifier creates a LogNotifier that writes to stdout by default.
func NewLogNotifier(out io.Writer) *LogNotifier {
	if out == nil {
		out = os.Stdout
	}
	return &LogNotifier{Out: out}
}

// Notify writes the alert to the configured writer.
func (l *LogNotifier) Notify(a Alert) error {
	_, err := fmt.Fprintln(l.Out, a.String())
	return err
}

// New constructs an Alert with the current timestamp.
func New(level Level, listener ports.Listener, message string) Alert {
	return Alert{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Listener:  listener,
	}
}
