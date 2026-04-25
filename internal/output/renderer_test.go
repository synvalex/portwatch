package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/output"
)

func TestRenderer_TextOutput(t *testing.T) {
	var buf bytes.Buffer
	cfg := config.DefaultOutputConfig()
	cfg.Format = "text"
	r, err := output.NewRendererWithWriter(&buf, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	listeners := []interface{}{
		makeListener("tcp", "0.0.0.0", 22, 100, "sshd"),
	}
	_ = listeners
	if err := r.Render(nil); err != nil {
		t.Fatalf("Render error: %v", err)
	}
	if !strings.Contains(buf.String(), "PROTO") {
		t.Errorf("expected header in output")
	}
}

func TestRenderer_InvalidConfig(t *testing.T) {
	cfg := config.DefaultOutputConfig()
	cfg.Format = "xml"
	_, err := output.NewRendererWithWriter(nil, cfg)
	if err == nil {
		t.Fatal("expected error for invalid format, got nil")
	}
}

func TestRenderer_JSONOutput(t *testing.T) {
	var buf bytes.Buffer
	cfg := config.DefaultOutputConfig()
	cfg.Format = "json"
	r, err := output.NewRendererWithWriter(&buf, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := r.Render(nil); err != nil {
		t.Fatalf("Render error: %v", err)
	}
	if !strings.Contains(buf.String(), "[]") {
		t.Errorf("expected empty JSON array, got: %s", buf.String())
	}
}
