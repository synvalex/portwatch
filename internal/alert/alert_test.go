package alert_test

import (
	"bytes"
	"net"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/ports"
	"github.com/user/portwatch/internal/rules"
)

func makeListener(ip string, port uint16, proto string) ports.Listener {
	return ports.Listener{
		IP:       net.ParseIP(ip),
		Port:     port,
		Protocol: proto,
	}
}

func TestLogNotifier_Notify(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewLogNotifier(&buf)
	l := makeListener("127.0.0.1", 8080, "tcp")
	a := alert.New(alert.LevelAlert, l, "test message")

	if err := n.Notify(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ALERT") {
		t.Errorf("expected ALERT in output, got: %s", out)
	}
	if !strings.Contains(out, "test message") {
		t.Errorf("expected message in output, got: %s", out)
	}
}

func TestDispatcher_AllowedRule(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewLogNotifier(&buf)

	r := rules.Rule{
		Name:     "allow-ssh",
		Port:     22,
		Protocol: "tcp",
		Action:   rules.ActionAllow,
	}

	d := alert.NewDispatcher([]rules.Rule{r}, n)
	d.Evaluate([]ports.Listener{makeListener("0.0.0.0", 22, "tcp")})

	if buf.Len() != 0 {
		t.Errorf("expected no alert for allowed listener, got: %s", buf.String())
	}
}

func TestDispatcher_DeniedRule(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewLogNotifier(&buf)

	r := rules.Rule{
		Name:     "deny-telnet",
		Port:     23,
		Protocol: "tcp",
		Action:   rules.ActionDeny,
	}

	d := alert.NewDispatcher([]rules.Rule{r}, n)
	d.Evaluate([]ports.Listener{makeListener("0.0.0.0", 23, "tcp")})

	if !strings.Contains(buf.String(), "ALERT") {
		t.Errorf("expected ALERT for denied listener, got: %s", buf.String())
	}
}

func TestDispatcher_UnexpectedListener(t *testing.T) {
	var buf bytes.Buffer
	n := alert.NewLogNotifier(&buf)

	d := alert.NewDispatcher([]rules.Rule{}, n)
	d.Evaluate([]ports.Listener{makeListener("0.0.0.0", 9999, "tcp")})

	if !strings.Contains(buf.String(), "WARN") {
		t.Errorf("expected WARN for unexpected listener, got: %s", buf.String())
	}
}
