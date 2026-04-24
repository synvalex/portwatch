package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindConfigFile_Explicit(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "my.yaml")
	if err := os.WriteFile(path, []byte("interval: 5s\nlog_level: info\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := FindConfigFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != path {
		t.Errorf("got %q, want %q", got, path)
	}
}

func TestFindConfigFile_ExplicitMissing(t *testing.T) {
	_, err := FindConfigFile("/tmp/does-not-exist-portwatch.yaml")
	if err == nil {
		t.Error("expected error for missing explicit path")
	}
}

func TestFindConfigFile_NoCandidates(t *testing.T) {
	// Change to a temp dir so local portwatch.yaml doesn't exist.
	orig, _ := os.Getwd()
	defer os.Chdir(orig) //nolint:errcheck
	os.Chdir(t.TempDir()) //nolint:errcheck

	_, err := FindConfigFile("")
	if err == nil {
		t.Error("expected error when no config file found")
	}
}

func TestFindConfigFile_LocalFile(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	defer os.Chdir(orig) //nolint:errcheck
	os.Chdir(dir)         //nolint:errcheck

	local := filepath.Join(dir, "portwatch.yaml")
	if err := os.WriteFile(local, []byte("interval: 10s\nlog_level: info\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := FindConfigFile("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "portwatch.yaml" {
		t.Errorf("got %q, want portwatch.yaml", got)
	}
}
