package alert

import (
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/user/portwatch/internal/ports"
)

// EmailConfig holds SMTP configuration for email notifications.
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       []string
}

// EmailNotifier sends alert notifications via SMTP email.
type EmailNotifier struct {
	cfg      EmailConfig
	formatter *Formatter
	dialFunc func(addr, from string, to []string, msg []byte) error
}

// NewEmailNotifier creates an EmailNotifier with the given SMTP config.
func NewEmailNotifier(cfg EmailConfig, f *Formatter) *EmailNotifier {
	n := &EmailNotifier{cfg: cfg, formatter: f}
	n.dialFunc = func(addr, from string, to []string, msg []byte) error {
		auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
		return smtp.SendMail(addr, auth, from, to, msg)
	}
	return n
}

// Notify sends an email alert for the given event and listener.
func (e *EmailNotifier) Notify(event string, l ports.Listener) error {
	if len(e.cfg.To) == 0 {
		return fmt.Errorf("email notifier: no recipients configured")
	}

	subject := fmt.Sprintf("[portwatch] %s: %s", strings.ToUpper(event), l.Address)
	body := e.formatter.Format(event, l)

	msg := buildMessage(e.cfg.From, e.cfg.To, subject, body)
	addr := fmt.Sprintf("%s:%d", e.cfg.Host, e.cfg.Port)

	if err := e.dialFunc(addr, e.cfg.From, e.cfg.To, msg); err != nil {
		return fmt.Errorf("email notifier: send failed: %w", err)
	}
	return nil
}

func buildMessage(from string, to []string, subject, body string) []byte {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("From: %s\r\n", from))
	sb.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(to, ", ")))
	sb.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	sb.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().UTC().Format(time.RFC1123Z)))
	sb.WriteString("MIME-Version: 1.0\r\n")
	sb.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
	sb.WriteString("\r\n")
	sb.WriteString(body)
	return []byte(sb.String())
}
