package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/audit"
	"github.com/user/portwatch/internal/ports"
)

func makeNotifierListener() ports.Listener {
	return ports.Listener{
		Proto:   "tcp",
		Address: "0.0.0.0:9090",
		Port:    9090,
		Process: &ports.ProcessInfo{
			PID:  42,
			Name: "myapp",
			Exe:  "/usr/bin/myapp",
		},
	}
}

func TestAuditNotifier_Appeared_JSON(t *testing.T) {
	var buf bytes.Buffer
	w := audit.NewWriter(&buf, "json")
	n := audit.NewNotifier(w)

	evt := alert.Event{
		Type:     alert.EventAppeared,
		Listener: makeNotifierListener(),
	}

	if err := n.Notify(evt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var record map[string]interface{}
	if err := json.NewDecoder(&buf).Decode(&record); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}

	if record["event"] != "appeared" {
		t.Errorf("expected event=appeared, got %v", record["event"])
	}
	if record["port"] != float64(9090) {
		t.Errorf("expected port=9090, got %v", record["port"])
	}
}

func TestAuditNotifier_Disappeared_Text(t *testing.T) {
	var buf bytes.Buffer
	w := audit.NewWriter(&buf, "text")
	n := audit.NewNotifier(w)

	evt := alert.Event{
		Type:     alert.EventDisappeared,
		Listener: makeNotifierListener(),
	}

	if err := n.Notify(evt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "disappeared") {
		t.Errorf("expected 'disappeared' in output, got: %s", out)
	}
	if !strings.Contains(out, "9090") {
		t.Errorf("expected port 9090 in output, got: %s", out)
	}
}

func TestAuditNotifier_EventTypeLabel(t *testing.T) {
	tests := []struct {
		eventType alert.EventType
		want      string
	}{
		{alert.EventAppeared, "appeared"},
		{alert.EventDisappeared, "disappeared"},
	}

	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			var buf bytes.Buffer
			w := audit.NewWriter(&buf, "json")
			n := audit.NewNotifier(w)
			evt := alert.Event{Type: tc.eventType, Listener: makeNotifierListener()}
			_ = n.Notify(evt)

			var record map[string]interface{}
			_ = json.NewDecoder(&buf).Decode(&record)
			if record["event"] != tc.want {
				t.Errorf("expected event=%s, got %v", tc.want, record["event"])
			}
		})
	}
}
