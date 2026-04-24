package ports

import (
	"net"
	"testing"
)

func makeListener(ip string, port uint16) Listener {
	return Listener{IP: net.ParseIP(ip), Port: port, Protocol: "tcp"}
}

func TestFilter_NilFilter(t *testing.T) {
	listeners := []Listener{
		makeListener("127.0.0.1", 8080),
		makeListener("0.0.0.0", 9090),
	}
	var f *Filter
	got := f.Apply(listeners)
	if len(got) != 2 {
		t.Fatalf("expected 2 listeners, got %d", len(got))
	}
}

func TestFilter_ExcludeLoopback(t *testing.T) {
	listeners := []Listener{
		makeListener("127.0.0.1", 8080),
		makeListener("0.0.0.0", 9090),
		makeListener("::1", 7070),
	}
	f := NewFilter(true, nil)
	got := f.Apply(listeners)
	if len(got) != 1 {
		t.Fatalf("expected 1 listener, got %d", len(got))
	}
	if got[0].Port != 9090 {
		t.Errorf("expected port 9090, got %d", got[0].Port)
	}
}

func TestFilter_ExcludePorts(t *testing.T) {
	listeners := []Listener{
		makeListener("0.0.0.0", 22),
		makeListener("0.0.0.0", 80),
		makeListener("0.0.0.0", 443),
	}
	f := NewFilter(false, []uint16{22, 443})
	got := f.Apply(listeners)
	if len(got) != 1 {
		t.Fatalf("expected 1 listener, got %d", len(got))
	}
	if got[0].Port != 80 {
		t.Errorf("expected port 80, got %d", got[0].Port)
	}
}

func TestFilter_Combined(t *testing.T) {
	listeners := []Listener{
		makeListener("127.0.0.1", 5432),
		makeListener("0.0.0.0", 22),
		makeListener("0.0.0.0", 8080),
	}
	f := NewFilter(true, []uint16{22})
	got := f.Apply(listeners)
	if len(got) != 1 {
		t.Fatalf("expected 1 listener, got %d", len(got))
	}
	if got[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", got[0].Port)
	}
}

func TestFilter_EmptyInput(t *testing.T) {
	f := NewFilter(true, []uint16{80})
	got := f.Apply(nil)
	if len(got) != 0 {
		t.Fatalf("expected 0 listeners, got %d", len(got))
	}
}
