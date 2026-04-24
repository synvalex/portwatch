package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestMain_VersionFlag verifies the --version flag prints version info and exits 0.
func TestMain_VersionFlag(t *testing.T) {
	if os.Getenv("PORTWATCH_RUN_MAIN") == "1" {
		os.Args = []string{"portwatch", "--version"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMain_VersionFlag")
	cmd.Env = append(os.Environ(), "PORTWATCH_RUN_MAIN=1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0, got error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(string(out), "portwatch") {
		t.Errorf("expected version output to contain 'portwatch', got: %s", out)
	}
}

// TestMain_MissingConfig verifies that a missing explicit config file causes exit 1.
func TestMain_MissingConfig(t *testing.T) {
	if os.Getenv("PORTWATCH_RUN_MAIN") == "1" {
		os.Args = []string{"portwatch", "--config", "/nonexistent/portwatch.yaml"}
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMain_MissingConfig")
	cmd.Env = append(os.Environ(), "PORTWATCH_RUN_MAIN=1")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-zero exit for missing config, got nil error\noutput: %s", out)
	}
	if !strings.Contains(string(out), "error") {
		t.Errorf("expected error message in output, got: %s", out)
	}
}
