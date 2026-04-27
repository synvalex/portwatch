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

func slackListener() ports.Listener {
	return ports.Listener{
		IP:       net.ParseIP("10.0.0.1"),
		Port:     9090,
		Protocol: "tcp",
	}
}

func TestSlackNotifier_Success(t *testing.T) {
	var received slackPayload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := NewSlackNotifier(srv.URL, time.Second)
	event := Event{Type: EventAppeared, Listener: slackListener()}
	if err := n.Notify(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Text == "" {
		t.Error("expected non-empty slack message text")
	}
}

func TestSlackNotifier_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	n := NewSlackNotifier(srv.URL, time.Second)
	event := Event{Type: EventAppeared, Listener: slackListener()}
	if err := n.Notify(event); err == nil {
		t.Fatal("expected error for 500 response")
	}
}

func TestSlackNotifier_Unreachable(t *testing.T) {
	n := NewSlackNotifier("http://127.0.0.1:1", 200*time.Millisecond)
	event := Event{Type: EventDisappeared, Listener: slackListener()}
	if err := n.Notify(event); err == nil {
		t.Fatal("expected error for unreachable server")
	}
}

func TestSlackNotifier_MessageContainsPort(t *testing.T) {
	var received slackPayload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&received) //nolint:errcheck
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := NewSlackNotifier(srv.URL, time.Second)
	event := Event{Type: EventUnexpected, Listener: slackListener()}
	_ = n.Notify(event)

	if !contains(received.Text, "9090") {
		t.Errorf("expected port 9090 in message, got: %s", received.Text)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
