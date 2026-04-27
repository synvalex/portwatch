package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// WebhookNotifier sends alert events to an HTTP endpoint as JSON payloads.
type WebhookNotifier struct {
	url     string
	client  *http.Client
	timeout time.Duration
}

// webhookPayload is the JSON body sent on each alert.
type webhookPayload struct {
	Event    string `json:"event"`
	Protocol string `json:"protocol"`
	Address  string `json:"address"`
	Port     uint16 `json:"port"`
	PID      int    `json:"pid,omitempty"`
	Process  string `json:"process,omitempty"`
	Timestamp string `json:"timestamp"`
}

// NewWebhookNotifier creates a WebhookNotifier that posts to the given URL.
// timeout controls the per-request deadline; zero uses a 5-second default.
func NewWebhookNotifier(url string, timeout time.Duration) *WebhookNotifier {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &WebhookNotifier{
		url:     url,
		timeout: timeout,
		client:  &http.Client{Timeout: timeout},
	}
}

// Notify implements Notifier by POSTing a JSON payload to the configured URL.
func (w *WebhookNotifier) Notify(event string, l ports.Listener) error {
	payload := webhookPayload{
		Event:     event,
		Protocol:  l.Protocol,
		Address:   l.Address.String(),
		Port:      l.Port,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	if l.Process != nil {
		payload.PID = l.Process.PID
		payload.Process = l.Process.Name
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post to %s: %w", w.url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: server returned %d", resp.StatusCode)
	}
	return nil
}
