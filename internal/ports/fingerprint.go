package ports

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

// Fingerprint holds identity information for a process binary.
type Fingerprint struct {
	ExePath string
	Hash    string // SHA-256 hex digest of the executable, empty if not hashed
	Cmdline string // raw cmdline, empty if trust_cmdline is false
}

// String returns a compact representation for logging.
func (f Fingerprint) String() string {
	parts := []string{f.ExePath}
	if f.Hash != "" {
		parts = append(parts, "sha256:"+f.Hash[:12])
	}
	return strings.Join(parts, " ")
}

// FingerprintService builds Fingerprint values for process paths.
type FingerprintService struct {
	hashExe    bool
	trustCmdln bool
	openFile   func(string) (io.ReadCloser, error)
}

// NewFingerprintService constructs a FingerprintService.
func NewFingerprintService(hashExe, trustCmdline bool) *FingerprintService {
	return &FingerprintService{
		hashExe:    hashExe,
		trustCmdln: trustCmdline,
		openFile:   func(p string) (io.ReadCloser, error) { return os.Open(p) },
	}
}

// Build creates a Fingerprint for the given pid and exe path.
func (s *FingerprintService) Build(pid int, exePath string) (Fingerprint, error) {
	fp := Fingerprint{ExePath: exePath}
	if s.hashExe && exePath != "" {
		h, err := s.hashFile(exePath)
		if err != nil {
			return fp, fmt.Errorf("fingerprint hash %s: %w", exePath, err)
		}
		fp.Hash = h
	}
	if s.trustCmdln && pid > 0 {
		cmdline, err := readProcCmdline(pid)
		if err == nil {
			fp.Cmdline = cmdline
		}
	}
	return fp, nil
}

func (s *FingerprintService) hashFile(path string) (string, error) {
	f, err := s.openFile(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func readProcCmdline(pid int) (string, error) {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return "", err
	}
	// cmdline entries are NUL-separated
	return strings.ReplaceAll(string(data), "\x00", " "), nil
}
