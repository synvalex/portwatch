package alert

import (
	"errors"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/ports"
)

func emailListener() ports.Listener {
	return ports.Listener{
		Address:  "0.0.0.0:8080",
		Protocol: "tcp",
	}
}

func newTestEmailNotifier(dialErr error, captured *[]byte) *EmailNotifier {
	cfg := EmailConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "user",
		Password: "pass",
		From:     "portwatch@example.com",
		To:       []string{"admin@example.com"},
	}
	f := NewFormatter()
	n := NewEmailNotifier(cfg, f)
	n.dialFunc = func(addr, from string, to []string, msg []byte) error {
		if captured != nil {
			*captured = append([]byte(nil), msg...)
		}
		return dialErr
	}
	return n
}

func TestEmailNotifier_Success(t *testing.T) {
	var captured []byte
	n := newTestEmailNotifier(nil, &captured)

	if err := n.Notify("appeared", emailListener()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	msg := string(captured)
	if !strings.Contains(msg, "portwatch") {
		t.Errorf("expected subject to contain 'portwatch', got: %s", msg)
	}
	if !strings.Contains(msg, "0.0.0.0:8080") {
		t.Errorf("expected body to contain listener address")
	}
}

func TestEmailNotifier_SMTPError(t *testing.T) {
	smtpErr := errors.New("connection refused")
	n := newTestEmailNotifier(smtpErr, nil)

	err := n.Notify("appeared", emailListener())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "send failed") {
		t.Errorf("expected 'send failed' in error, got: %v", err)
	}
}

func TestEmailNotifier_NoRecipients(t *testing.T) {
	cfg := EmailConfig{
		Host: "smtp.example.com",
		Port: 587,
		From: "portwatch@example.com",
		To:   []string{},
	}
	n := NewEmailNotifier(cfg, NewFormatter())

	err := n.Notify("appeared", emailListener())
	if err == nil {
		t.Fatal("expected error for no recipients")
	}
	if !strings.Contains(err.Error(), "no recipients") {
		t.Errorf("expected 'no recipients' in error, got: %v", err)
	}
}

func TestBuildMessage_ContainsHeaders(t *testing.T) {
	msg := buildMessage("from@example.com", []string{"to@example.com"}, "Test Subject", "Hello body")
	s := string(msg)

	for _, want := range []string{"From:", "To:", "Subject:", "MIME-Version:", "Hello body"} {
		if !strings.Contains(s, want) {
			t.Errorf("expected message to contain %q", want)
		}
	}
}
