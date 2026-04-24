package rules

import (
	"testing"
)

func TestRule_Matches(t *testing.T) {
	tests := []struct {
		name     string
		rule     Rule
		port     uint16
		protocol string
		address  string
		want     bool
	}{
		{
			name:     "exact port match",
			rule:     Rule{Name: "http", Port: 80, Action: ActionAlert},
			port:     80, protocol: "tcp", address: "0.0.0.0",
			want:     true,
		},
		{
			name:     "port mismatch",
			rule:     Rule{Name: "http", Port: 80, Action: ActionAlert},
			port:     443, protocol: "tcp", address: "0.0.0.0",
			want:     false,
		},
		{
			name:     "protocol filter match",
			rule:     Rule{Name: "dns", Port: 53, Protocol: "udp", Action: ActionAllow},
			port:     53, protocol: "udp", address: "127.0.0.1",
			want:     true,
		},
		{
			name:     "protocol filter mismatch",
			rule:     Rule{Name: "dns", Port: 53, Protocol: "udp", Action: ActionAllow},
			port:     53, protocol: "tcp", address: "127.0.0.1",
			want:     false,
		},
		{
			name:     "cidr address match",
			rule:     Rule{Name: "local", Port: 8080, Address: "192.168.0.0/16", Action: ActionAllow},
			port:     8080, protocol: "tcp", address: "192.168.1.5",
			want:     true,
		},
		{
			name:     "cidr address mismatch",
			rule:     Rule{Name: "local", Port: 8080, Address: "192.168.0.0/16", Action: ActionAllow},
			port:     8080, protocol: "tcp", address: "10.0.0.1",
			want:     false,
		},
		{
			name:     "wildcard rule matches anything",
			rule:     Rule{Name: "all", Action: ActionAlert},
			port:     9999, protocol: "tcp", address: "1.2.3.4",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rule.Matches(tt.port, tt.protocol, tt.address)
			if got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRule_Validate(t *testing.T) {
	valid := Rule{Name: "test", Port: 80, Protocol: "tcp", Action: ActionAlert}
	if err := valid.Validate(); err != nil {
		t.Errorf("expected valid rule, got error: %v", err)
	}

	noName := Rule{Action: ActionAllow}
	if err := noName.Validate(); err == nil {
		t.Error("expected error for empty name")
	}

	badAction := Rule{Name: "x", Action: "block"}
	if err := badAction.Validate(); err == nil {
		t.Error("expected error for unknown action")
	}

	badProto := Rule{Name: "x", Protocol: "icmp", Action: ActionAllow}
	if err := badProto.Validate(); err == nil {
		t.Error("expected error for invalid protocol")
	}

	badAddr := Rule{Name: "x", Address: "not-an-ip", Action: ActionAllow}
	if err := badAddr.Validate(); err == nil {
		t.Error("expected error for invalid address")
	}
}
