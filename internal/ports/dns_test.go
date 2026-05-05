package ports

import (
	"net"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

func makeDNSListener(ip string, port uint16) Listener {
	return Listener{
		Address: ParsedAddress{
			IP:   net.ParseIP(ip),
			Port: port,
		},
		Protocol: "tcp",
	}
}

func TestDNSEnricher_SkipsNilIP(t *testing.T) {
	e := NewDNSEnricher(2*time.Second, 2, zerolog.Nop())
	listeners := []Listener{
		{Protocol: "tcp"},
	}
	out := e.Enrich(listeners)
	if out[0].Hostname != "" {
		t.Errorf("expected empty hostname for nil IP, got %q", out[0].Hostname)
	}
}

func TestDNSEnricher_LoopbackResolvesOrEmpty(t *testing.T) {
	e := NewDNSEnricher(2*time.Second, 2, zerolog.Nop())
	listeners := []Listener{
		makeDNSListener("127.0.0.1", 80),
	}
	// We don't assert a specific hostname since it depends on the host, but
	// the enricher must not panic and must return the same slice length.
	out := e.Enrich(listeners)
	if len(out) != 1 {
		t.Fatalf("expected 1 listener, got %d", len(out))
	}
}

func TestDNSEnricher_CachesResult(t *testing.T) {
	e := NewDNSEnricher(2*time.Second, 2, zerolog.Nop())
	e.mu.Lock()
	e.cache["10.0.0.1"] = "cached.example.com."
	e.mu.Unlock()

	listeners := []Listener{
		makeDNSListener("10.0.0.1", 443),
	}
	out := e.Enrich(listeners)
	if out[0].Hostname != "cached.example.com." {
		t.Errorf("expected cached hostname, got %q", out[0].Hostname)
	}
}

func TestDNSEnricher_MultipleListeners(t *testing.T) {
	e := NewDNSEnricher(2*time.Second, 4, zerolog.Nop())
	e.mu.Lock()
	e.cache["192.168.1.1"] = "router.local."
	e.cache["192.168.1.2"] = "host2.local."
	e.mu.Unlock()

	listeners := []Listener{
		makeDNSListener("192.168.1.1", 22),
		makeDNSListener("192.168.1.2", 8080),
	}
	out := e.Enrich(listeners)
	if out[0].Hostname != "router.local." {
		t.Errorf("listener 0: expected router.local., got %q", out[0].Hostname)
	}
	if out[1].Hostname != "host2.local." {
		t.Errorf("listener 1: expected host2.local., got %q", out[1].Hostname)
	}
}
