package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// SlackNotifier sends alert events to a Slack incoming webhook.
type SlackNotifier struct {
	webhookURL string
	client     *http.Client
}

type slackPayload struct {
	Text string `json:"text"`
}

// NewSlackNotifier creates a SlackNotifier that posts to the given Slack webhook URL.
func NewSlackNotifier(webhookURL string, timeout time.Duration) *SlackNotifier {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &SlackNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: timeout},
	}
}

// Notify sends a formatted Slack message for the given event.
func (s *SlackNotifier) Notify(event Event) error {
	msg := formatSlackMessage(event)
	payload := slackPayload{Text: msg}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("slack: marshal payload: %w", err)
	}

	resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("slack: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack: unexpected status %d", resp.StatusCode)
	}
	return nil
}

func formatSlackMessage(event Event) string {
	l := event.Listener
	icon := ":warning:"
	if event.Type == EventAppeared {
		icon = ":large_green_circle:"
	} else if event.Type == EventDisappeared {
		icon = ":red_circle:"
	}

	base := fmt.Sprintf("%s *portwatch* `%s` %s on `%s`",
		icon, event.Type, l.Protocol, ports.FormatAddress(l.IP, l.Port))

	if l.Process != nil {
		base += fmt.Sprintf(" — process *%s* (pid %d)", l.Process.Name, l.Process.PID)
	}
	return base
}
