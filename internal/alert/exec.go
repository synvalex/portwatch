package alert

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/yourorg/portwatch/internal/config"
	"github.com/yourorg/portwatch/internal/ports"
)

// ExecNotifier runs an external command when an alert event occurs.
// The event payload is passed to the command via stdin as JSON.
type ExecNotifier struct {
	cfg config.ExecConfig
}

// NewExecNotifier creates an ExecNotifier. Returns nil if the config is disabled.
func NewExecNotifier(cfg config.ExecConfig) (*ExecNotifier, error) {
	if !cfg.Enabled {
		return nil, nil
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &ExecNotifier{cfg: cfg}, nil
}

// Notify executes the configured command with the event serialised as JSON on stdin.
func (n *ExecNotifier) Notify(event Event, l ports.Listener) error {
	payload, err := json.Marshal(map[string]any{
		"event":    string(event),
		"protocol": l.Protocol,
		"address":  l.Address.String(),
		"port":     l.Address.Port,
		"pid":      pidFromListener(l),
		"process":  processNameFromListener(l),
	})
	if err != nil {
		return fmt.Errorf("exec notifier: marshal payload: %w", err)
	}

	timeout := n.cfg.Timeout
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var cmd *exec.Cmd
	if n.cfg.Shell {
		args := append([]string{"-c", n.cfg.Command}, n.cfg.Args...)
		cmd = exec.CommandContext(ctx, "sh", args...)
	} else {
		cmd = exec.CommandContext(ctx, n.cfg.Command, n.cfg.Args...)
	}

	cmd.Stdin = bytes.NewReader(payload)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("exec notifier: command failed: %w (output: %s)", err, string(out))
	}
	return nil
}

func pidFromListener(l ports.Listener) int {
	if l.Process != nil {
		return l.Process.PID
	}
	return 0
}

func processNameFromListener(l ports.Listener) string {
	if l.Process != nil {
		return l.Process.Name
	}
	return ""
}
