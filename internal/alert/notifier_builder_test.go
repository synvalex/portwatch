package alert

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func testLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelError}))
}

func baseConfig() config.Config {
	cfg := config.DefaultConfig()
	return cfg
}

func TestNotifierBuilder_OnlyLogNotifier(t *testing.T) {
	cfg := baseConfig()
	cfg.Alert.Webhook.Enabled = false
	cfg.Alert.Slack.Enabled = false
	cfg.Alert.Email.Enabled = false

	b := NewNotifierBuilder(testLogger())
	n, err := b.Build(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNotifierBuilder_WebhookEnabled(t *testing.T) {
	cfg := baseConfig()
	cfg.Alert.Webhook.Enabled = true
	cfg.Alert.Webhook.URL = "http://example.com/hook"
	cfg.Alert.Webhook.Timeout = 5 * time.Second

	b := NewNotifierBuilder(testLogger())
	n, err := b.Build(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNotifierBuilder_SlackEnabled(t *testing.T) {
	cfg := baseConfig()
	cfg.Alert.Slack.Enabled = true
	cfg.Alert.Slack.WebhookURL = "http://hooks.slack.com/test"
	cfg.Alert.Slack.Timeout = 5 * time.Second

	b := NewNotifierBuilder(testLogger())
	n, err := b.Build(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNotifierBuilder_EmailEnabled(t *testing.T) {
	cfg := baseConfig()
	cfg.Alert.Email.Enabled = true
	cfg.Alert.Email.SMTPHost = "smtp.example.com"
	cfg.Alert.Email.SMTPPort = 587
	cfg.Alert.Email.From = "alert@example.com"
	cfg.Alert.Email.To = []string{"admin@example.com"}
	cfg.Alert.Email.Subject = "portwatch alert"
	cfg.Alert.Email.Timeout = 5 * time.Second

	b := NewNotifierBuilder(testLogger())
	n, err := b.Build(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}
