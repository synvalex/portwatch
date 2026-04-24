package ports

import (
	"net"
	"testing"
)

func makeChainListener(ip string, port uint16) Listener {
	return Listener{
		Addr: ParsedAddr{IP: net.ParseIP(ip), Port: port},
		Proto: "tcp",
	}
}

func TestFilterChain_Empty_PassThrough(t *testing.T) {
	chain := NewFilterChain()
	input := []Listener{
		makeChainListener("127.0.0.1", 8080),
		makeChainListener("0.0.0.0", 9090),
	}
	got := chain.Apply(input)
	if len(got) != len(input) {
		t.Fatalf("expected %d listeners, got %d", len(input), len(got))
	}
}

func TestFilterChain_NilFiltersSkipped(t *testing.T) {
	chain := NewFilterChain(nil, nil)
	if chain.Len() != 0 {
		t.Fatalf("expected 0 filters, got %d", chain.Len())
	}
}

func TestFilterChain_SingleFilter(t *testing.T) {
	f := NewFilter(&FilterOptions{ExcludeLoopback: true})
	chain := NewFilterChain(f)

	input := []Listener{
		makeChainListener("127.0.0.1", 8080),
		makeChainListener("0.0.0.0", 9090),
	}
	got := chain.Apply(input)
	if len(got) != 1 {
		t.Fatalf("expected 1 listener, got %d", len(got))
	}
	if got[0].Addr.Port != 9090 {
		t.Errorf("expected port 9090, got %d", got[0].Addr.Port)
	}
}

func TestFilterChain_MultipleFilters_ANDSemantics(t *testing.T) {
	f1 := NewFilter(&FilterOptions{ExcludeLoopback: true})
	f2 := NewFilter(&FilterOptions{ExcludePorts: []uint16{9090}})
	chain := NewFilterChain(f1, f2)

	input := []Listener{
		makeChainListener("127.0.0.1", 8080), // excluded by f1
		makeChainListener("0.0.0.0", 9090),   // excluded by f2
		makeChainListener("0.0.0.0", 3000),   // accepted by both
	}
	got := chain.Apply(input)
	if len(got) != 1 {
		t.Fatalf("expected 1 listener, got %d", len(got))
	}
	if got[0].Addr.Port != 3000 {
		t.Errorf("expected port 3000, got %d", got[0].Addr.Port)
	}
}

func TestFilterChain_Len(t *testing.T) {
	f1 := NewFilter(nil)
	f2 := NewFilter(nil)
	chain := NewFilterChain(f1, f2)
	if chain.Len() != 2 {
		t.Errorf("expected Len 2, got %d", chain.Len())
	}
}
