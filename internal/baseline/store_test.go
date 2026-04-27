package baseline

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/ports"
)

func makeListener(proto, addr string, port uint16) ports.Listener {
	return ports.Listener{Proto: proto, Address: addr, Port: port}
}

func TestStore_AddAndContains(t *testing.T) {
	s := NewStore("")
	l := makeListener("tcp", "0.0.0.0", 8080)
	if s.Contains(l) {
		t.Fatal("expected listener not to be in empty store")
	}
	s.Add(l)
	if !s.Contains(l) {
		t.Fatal("expected listener to be found after Add")
	}
}

func TestStore_Contains_DifferentPort(t *testing.T) {
	s := NewStore("")
	s.Add(makeListener("tcp", "0.0.0.0", 8080))
	if s.Contains(makeListener("tcp", "0.0.0.0", 9090)) {
		t.Error("different port should not match")
	}
}

func TestStore_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	s1 := NewStore(path)
	s1.Add(makeListener("tcp", "0.0.0.0", 443))
	s1.Add(makeListener("udp", "127.0.0.1", 53))
	if err := s1.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	s2 := NewStore(path)
	if err := s2.Load(); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if !s2.Contains(makeListener("tcp", "0.0.0.0", 443)) {
		t.Error("expected tcp/443 after load")
	}
	if !s2.Contains(makeListener("udp", "127.0.0.1", 53)) {
		t.Error("expected udp/53 after load")
	}
}

func TestStore_Load_MissingFile_NoError(t *testing.T) {
	s := NewStore("/nonexistent/path/baseline.json")
	if err := s.Load(); err != nil {
		t.Errorf("expected no error for missing file, got: %v", err)
	}
}

func TestStore_Load_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0o644)
	s := NewStore(path)
	if err := s.Load(); err == nil {
		t.Error("expected error for invalid JSON")
	}
}
