package audit

import (
	"bytes"
	"encoding/json"
	"net"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/ports"
	"github.com/user/portwatch/internal/snapshot"
)

func makeAuditListener(port uint16, pid int, name string) ports.Listener {
	l := ports.Listener{
		Protocol: "tcp",
		Address:  net.ParseIP("0.0.0.0"),
		Port:     port,
	}
	if pid > 0 {
		l.Process = &ports.ProcessInfo{PID: pid, Name: name}
	}
	return l
}

func TestWriter_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, "json")
	l := makeAuditListener(8080, 1234, "nginx")

	if err := w.Write(l, snapshot.EventAppeared); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var ev Event
	if err := json.Unmarshal(buf.Bytes(), &ev); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if ev.Port != 8080 {
		t.Errorf("expected port 8080, got %d", ev.Port)
	}
	if ev.PID != 1234 {
		t.Errorf("expected PID 1234, got %d", ev.PID)
	}
	if ev.Process != "nginx" {
		t.Errorf("expected process \"nginx\", got %q", ev.Process)
	}
	if ev.Type != snapshot.EventAppeared {
		t.Errorf("expected type Appeared, got %v", ev.Type)
	}
}

func TestWriter_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, "text")
	l := makeAuditListener(9090, 42, "sshd")

	if err := w.Write(l, snapshot.EventDisappeared); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := buf.String()
	if !strings.Contains(line, "9090") {
		t.Errorf("expected port in output, got: %s", line)
	}
	if !strings.Contains(line, "sshd") {
		t.Errorf("expected process name in output, got: %s", line)
	}
}

func TestWriter_NoProcess(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf, "json")
	l := makeAuditListener(443, 0, "")

	if err := w.Write(l, snapshot.EventAppeared); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var ev Event
	if err := json.Unmarshal(buf.Bytes(), &ev); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if ev.PID != 0 {
		t.Errorf("expected PID 0 for no-process listener, got %d", ev.PID)
	}
}
