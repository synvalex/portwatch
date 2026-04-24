package ports

import (
	"testing"
)

func TestParseAddress_Valid(t *testing.T) {
	tests := []struct {
		input   string
		wantIP  string
		wantPort int
	}{
		{"127.0.0.1:8080", "127.0.0.1", 8080},
		{"0.0.0.0:443", "0.0.0.0", 443},
		{"[::1]:22", "::1", 22},
	}
	for _, tc := range tests {
		ip, port, err := ParseAddress(tc.input)
		if err != nil {
			t.Errorf("ParseAddress(%q) unexpected error: %v", tc.input, err)
			continue
		}
		if ip != tc.wantIP {
			t.Errorf("ParseAddress(%q) ip = %q, want %q", tc.input, ip, tc.wantIP)
		}
		if port != tc.wantPort {
			t.Errorf("ParseAddress(%q) port = %d, want %d", tc.input, port, tc.wantPort)
		}
	}
}

func TestParseAddress_Invalid(t *testing.T) {
	invalids := []string{"not-a-port", "0.0.0.0:99999", "0.0.0.0:0", ""}
	for _, tc := range invalids {
		_, _, err := ParseAddress(tc)
		if err == nil {
			t.Errorf("ParseAddress(%q) expected error, got nil", tc)
		}
	}
}

func TestDeduplicateListeners(t *testing.T) {
	input := []Listener{
		{Protocol: "tcp", Address: "0.0.0.0", Port: 80},
		{Protocol: "tcp", Address: "0.0.0.0", Port: 80},
		{Protocol: "tcp", Address: "0.0.0.0", Port: 443},
		{Protocol: "udp", Address: "0.0.0.0", Port: 53},
	}
	got := DeduplicateListeners(input)
	if len(got) != 3 {
		t.Errorf("DeduplicateListeners returned %d entries, want 3", len(got))
	}
}

func TestListenerString(t *testing.T) {
	l := Listener{Protocol: "tcp", Address: "127.0.0.1", Port: 8080, PID: 1234, Process: "nginx"}
	want := "tcp 127.0.0.1:8080 (pid=1234, process=nginx)"
	if l.String() != want {
		t.Errorf("Listener.String() = %q, want %q", l.String(), want)
	}
}
