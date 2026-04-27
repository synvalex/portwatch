package alert

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/portwatch/internal/ports"
)

func webhookListener() ports.Listener {
	return ports.Listener{
		Protocol: "tcp",
		Address:  net.ParseIP("0.0.0.0"),
		Port:     9090,
		Process:  &ports.ProcessInfo{PID: 42, Name: "myapp"},
	}
}

func TestWebhookNotifier_Success(t *testing.T) {
	var received webhookPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewWebhookNotifier(ts.URL, 0)
	l := webhookListener()
	if err := n.Notify("appeared", l); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Event != "appeared" {
		t.Errorf("event = %q, want %q", received.Event, "appeared")
	}
	if received.Port != 9090 {
		t.Errorf("port = %d, want 9090", received.Port)
	}
	if received.PID != 42 {
		t.Errorf("pid = %d, want 42", received.PID)
	}
	if received.Process != "myapp" {
		t.Errorf("process = %q, want %q", received.Process, "myapp")
	}
}

func TestWebhookNotifier_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewWebhookNotifier(ts.URL, 0)
	if err := n.Notify("appeared", webhookListener()); err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
}

func TestWebhookNotifier_Unreachable(t *testing.T) {
	n := NewWebhookNotifier("http://127.0.0.1:1", 200*time.Millisecond)
	if err := n.Notify("appeared", webhookListener()); err == nil {
		t.Fatal("expected error for unreachable host, got nil")
	}
}

func TestWebhookNotifier_NoProcess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	l := ports.Listener{
		Protocol: "udp",
		Address:  net.ParseIP("0.0.0.0"),
		Port:     5353,
		Process:  nil,
	}
	n := NewWebhookNotifier(ts.URL, 0)
	if err := n.Notify("appeared", l); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
