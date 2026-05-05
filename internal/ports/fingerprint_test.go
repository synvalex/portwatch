package ports

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestFingerprint_String_WithHash(t *testing.T) {
	fp := Fingerprint{
		ExePath: "/usr/sbin/sshd",
		Hash:    "abcdef123456789000000000000000000000000000000000000000000000000",
	}
	s := fp.String()
	if !strings.Contains(s, "/usr/sbin/sshd") {
		t.Errorf("expected exe path in string, got %q", s)
	}
	if !strings.Contains(s, "sha256:abcdef123456") {
		t.Errorf("expected truncated hash in string, got %q", s)
	}
}

func TestFingerprint_String_NoHash(t *testing.T) {
	fp := Fingerprint{ExePath: "/usr/bin/nginx"}
	if fp.String() != "/usr/bin/nginx" {
		t.Errorf("unexpected string: %q", fp.String())
	}
}

func TestFingerprintService_Build_HashesFile(t *testing.T) {
	svc := NewFingerprintService(true, false)
	// Use a real temp file so we can hash it.
	tmp, err := os.CreateTemp(t.TempDir(), "fp-test")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = tmp.WriteString("hello portwatch")
	tmp.Close()

	fp, err := svc.Build(0, tmp.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fp.Hash == "" {
		t.Error("expected non-empty hash")
	}
	if fp.ExePath != tmp.Name() {
		t.Errorf("expected ExePath %s, got %s", tmp.Name(), fp.ExePath)
	}
}

func TestFingerprintService_Build_NoHash(t *testing.T) {
	svc := NewFingerprintService(false, false)
	fp, err := svc.Build(0, "/some/path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fp.Hash != "" {
		t.Errorf("expected empty hash, got %q", fp.Hash)
	}
}

func TestFingerprintService_Build_MissingFile_ReturnsError(t *testing.T) {
	svc := NewFingerprintService(true, false)
	_, err := svc.Build(0, "/nonexistent/path/binary")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestFingerprintService_Build_CustomOpener(t *testing.T) {
	svc := NewFingerprintService(true, false)
	svc.openFile = func(p string) (io.ReadCloser, error) {
		return io.NopCloser(strings.NewReader("fake binary content")), nil
	}
	fp, err := svc.Build(0, "/fake/binary")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fp.Hash == "" {
		t.Error("expected hash from custom opener")
	}
}
