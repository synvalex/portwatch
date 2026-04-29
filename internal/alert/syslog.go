//go:build !windows

package alert

import (
	"fmt"
	"log/syslog"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/ports"
)

// SyslogNotifier sends alert events to syslog.
type SyslogNotifier struct {
	writer *syslog.Writer
	tag    string
}

// NewSyslogNotifier creates a SyslogNotifier from the given config.
// Returns nil and an error if syslog is unavailable or config is invalid.
func NewSyslogNotifier(cfg config.SyslogConfig) (*SyslogNotifier, error) {
	if !cfg.Enabled {
		return nil, nil
	}
	priority, err := parseFacility(cfg.Facility)
	if err != nil {
		return nil, fmt.Errorf("syslog: %w", err)
	}
	w, err := syslog.Dial(cfg.Network, cfg.Addr, priority|syslog.LOG_WARNING, cfg.Tag)
	if err != nil {
		return nil, fmt.Errorf("syslog: dial: %w", err)
	}
	return &SyslogNotifier{writer: w, tag: cfg.Tag}, nil
}

// Notify sends an alert event to syslog.
func (s *SyslogNotifier) Notify(event Event) error {
	msg := formatSyslogMessage(event)
	switch event.Kind {
	case KindDenied, KindUnexpected:
		return s.writer.Warning(msg)
	default:
		return s.writer.Info(msg)
	}
}

// Close releases the syslog connection.
func (s *SyslogNotifier) Close() error {
	return s.writer.Close()
}

func formatSyslogMessage(event Event) string {
	l := event.Listener
	base := fmt.Sprintf("[%s] %s %s:%d",
		event.Kind, l.Protocol, l.Address.IP, l.Address.Port)
	if l.Process != nil {
		base += fmt.Sprintf(" pid=%d exe=%s", l.Process.PID, l.Process.Exe)
	}
	return base
}

func parseFacility(name string) (syslog.Priority, error) {
	facilities := map[string]syslog.Priority{
		"daemon": syslog.LOG_DAEMON,
		"local0": syslog.LOG_LOCAL0,
		"local1": syslog.LOG_LOCAL1,
		"local2": syslog.LOG_LOCAL2,
		"local3": syslog.LOG_LOCAL3,
		"local4": syslog.LOG_LOCAL4,
		"local5": syslog.LOG_LOCAL5,
		"local6": syslog.LOG_LOCAL6,
		"local7": syslog.LOG_LOCAL7,
	}
	p, ok := facilities[name]
	if !ok {
		return 0, fmt.Errorf("unknown facility %q", name)
	}
	return p, nil
}

// ensure SyslogNotifier satisfies the Notifier interface.
var _ Notifier = (*SyslogNotifier)(nil)

// keep ports import used for Listener type reference in formatSyslogMessage.
var _ = ports.Listener{}
