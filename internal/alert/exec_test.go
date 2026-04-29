package alert

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/config"
	"github.com/yourorg/portwatch/internal/ports"
)

func execListener() ports.Listener {
	l, _ := ports.ParseAddress("127.0.0.1:9090")
	return ports.Listener{
		Protocol: "tcp",
		Address:  l,
		Process:  &ports.ProcessInfo{PID: 42, Name: "myapp"},
	}
}

func TestExecNotifier_Disabled_ReturnsNil(t *testing.T) {
	cfg := config.DefaultExecConfig()
	cfg.Enabled = false
	n, err := NewExecNotifier(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != nil {
		t.Fatal("expected nil notifier when disabled")
	}
}

func TestExecNotifier_MissingCommand_ReturnsError(t *testing.T) {
	cfg := config.DefaultExecConfig()
	cfg.Enabled = true
	cfg.Command = ""
	_, err := NewExecNotifier(cfg)
	if err == nil {
		t.Fatal("expected error for missing command")
	}
}

func TestExecNotifier_EchoCommand_Success(t *testing.T) {
	cfg := config.ExecConfig{
		Enabled: true,
		Command: "cat",
		Timeout: 2 * time.Second,
	}
	n, err := NewExecNotifier(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := n.Notify(EventAppeared, execListener()); err != nil {
		t.Fatalf("Notify returned error: %v", err)
	}
}

func TestExecNotifier_BadCommand_ReturnsError(t *testing.T) {
	cfg := config.ExecConfig{
		Enabled: true,
		Command: "false",
		Timeout: 2 * time.Second,
	}
	n, err := NewExecNotifier(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := n.Notify(EventAppeared, execListener()); err == nil {
		t.Fatal("expected error from failing command")
	}
}

func TestExecNotifier_Timeout_ReturnsError(t *testing.T) {
	cfg := config.ExecConfig{
		Enabled:  true,
		Command:  "sleep",
		Args:     []string{"10"},
		Timeout:  50 * time.Millisecond,
	}
	n, err := NewExecNotifier(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := n.Notify(EventAppeared, execListener()); err == nil {
		t.Fatal("expected timeout error")
	}
}
