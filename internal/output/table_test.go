package output_test

import (
	"bytes"
	"encoding/json"
	"net"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/output"
	"github.com/user/portwatch/internal/ports"
)

func makeListener(proto string, ip string, port uint16, pid int, name string) ports.Listener {
	l := ports.Listener{
		Proto: proto,
		Addr:  net.ParseIP(ip),
		Port:  port,
	}
	if pid > 0 {
		l.Process = &ports.ProcessInfo{PID: pid, Name: name}
	}
	return l
}

func TestTableWriter_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	tw := output.NewTableWriter(&buf, "text")
	listeners := []ports.Listener{
		makeListener("tcp", "0.0.0.0", 80, 1234, "nginx"),
		makeListener("tcp", "0.0.0.0", 443, 0, ""),
	}
	if err := tw.Write(listeners); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "nginx") {
		t.Errorf("expected 'nginx' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "PROTO") {
		t.Errorf("expected header in output, got:\n%s", out)
	}
}

func TestTableWriter_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	tw := output.NewTableWriter(&buf, "json")
	listeners := []ports.Listener{
		makeListener("udp", "127.0.0.1", 53, 999, "dnsmasq"),
	}
	if err := tw.Write(listeners); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0]["proto"] != "udp" {
		t.Errorf("expected proto=udp, got %v", result[0]["proto"])
	}
}

func TestTableWriter_NoProcess(t *testing.T) {
	var buf bytes.Buffer
	tw := output.NewTableWriter(&buf, "text")
	listeners := []ports.Listener{
		makeListener("tcp", "0.0.0.0", 8080, 0, ""),
	}
	if err := tw.Write(listeners); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "-") {
		t.Errorf("expected '-' placeholder for missing process")
	}
}
