//go:build !windows

package alert

import (
	"net"
	"testing"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/ports"
)

func syslogListener() ports.Listener {
	return ports.Listener{
		Protocol: "tcp",
		Address:  ports.ParsedAddress{IP: net.ParseIP("0.0.0.0"), Port: 8080},
	}
}

func TestSyslogNotifier_Disabled_ReturnsNil(t *testing.T) {
	cfg := config.DefaultSyslogConfig()
	cfg.Enabled = false
	n, err := NewSyslogNotifier(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != nil {
		t.Error("expected nil notifier when disabled")
	}
}

func TestSyslogNotifier_InvalidFacility_ReturnsError(t *testing.T) {
	cfg := config.SyslogConfig{
		Enabled:  true,
		Tag:      "portwatch",
		Facility: "bogus",
	}
	_, err := NewSyslogNotifier(cfg)
	if err == nil {
		t.Error("expected error for invalid facility")
	}
}

func TestFormatSyslogMessage_WithProcess(t *testing.T) {
	event := Event{
		Kind: KindUnexpected,
		Listener: ports.Listener{
			Protocol: "tcp",
			Address:  ports.ParsedAddress{IP: net.ParseIP("127.0.0.1"), Port: 9090},
			Process:  &ports.ProcessInfo{PID: 42, Name: "nginx", Exe: "/usr/sbin/nginx"},
		},
	}
	msg := formatSyslogMessage(event)
	for _, want := range []string{"9090", "pid=42", "/usr/sbin/nginx"} {
		if !containsStr(msg, want) {
			t.Errorf("expected %q in message %q", want, msg)
		}
	}
}

func TestFormatSyslogMessage_NoProcess(t *testing.T) {
	event := Event{
		Kind:     KindAppeared,
		Listener: syslogListener(),
	}
	msg := formatSyslogMessage(event)
	if !containsStr(msg, "8080") {
		t.Errorf("expected port in message %q", msg)
	}
}

func TestParseFacility_ValidNames(t *testing.T) {
	for _, name := range []string{"daemon", "local0", "local7"} {
		_, err := parseFacility(name)
		if err != nil {
			t.Errorf("unexpected error for facility %q: %v", name, err)
		}
	}
}

func TestParseFacility_InvalidName(t *testing.T) {
	_, err := parseFacility("unknown")
	if err == nil {
		t.Error("expected error for unknown facility")
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsSubstring(s, sub))
}

func containsSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
