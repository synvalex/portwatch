//go:build !windows

package alert

import (
	"fmt"

	"golang.org/x/exp/slog"

	"github.com/user/portwatch/internal/config"
)

// BuildSyslogNotifier creates a SyslogNotifier and registers it with the
// MultiNotifier if syslog is enabled in cfg. Errors are logged but do not
// abort startup — a misconfigured syslog should not prevent other notifiers.
func BuildSyslogNotifier(mn *MultiNotifier, cfg config.SyslogConfig, log *slog.Logger) error {
	if !cfg.Enabled {
		return nil
	}
	n, err := NewSyslogNotifier(cfg)
	if err != nil {
		return fmt.Errorf("syslog notifier: %w", err)
	}
	if n != nil {
		mn.Add(n)
		log.Info("syslog notifier enabled", "tag", cfg.Tag, "facility", cfg.Facility)
	}
	return nil
}
