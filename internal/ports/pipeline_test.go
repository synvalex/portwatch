package ports_test

import (
	"errors"
	"net"
	"testing"

	"github.com/user/portwatch/internal/ports"
)

// stubScanner satisfies ports.Scanner for tests.
type stubScanner struct {
	listeners []ports.Listener
	err       error
}

func (s *stubScanner) Listeners() ([]ports.Listener, error) {
	return s.listeners, s.err
}

func makePipelineListener(port uint16, proto string) ports.Listener {
	return ports.Listener{
		Addr:     net.ParseIP("0.0.0.0"),
		Port:     port,
		Protocol: proto,
	}
}

func TestPipeline_Run_ReturnsListeners(t *testing.T) {
	input := []ports.Listener{
		makePipelineListener(8080, "tcp"),
		makePipelineListener(9090, "tcp"),
	}
	p := ports.NewPipeline(&stubScanner{listeners: input}, nil, nil, ports.SortOptions{})
	out, err := p.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 listeners, got %d", len(out))
	}
}

func TestPipeline_Run_ScannerError(t *testing.T) {
	p := ports.NewPipeline(&stubScanner{err: errors.New("scan failed")}, nil, nil, ports.SortOptions{})
	_, err := p.Run()
	if err == nil {
		t.Fatal("expected error from scanner, got nil")
	}
}

func TestPipeline_Run_SortApplied(t *testing.T) {
	input := []ports.Listener{
		makePipelineListener(9090, "tcp"),
		makePipelineListener(8080, "tcp"),
		makePipelineListener(443, "tcp"),
	}
	opts := ports.SortOptions{Field: ports.SortFieldPort, Ascending: true}
	p := ports.NewPipeline(&stubScanner{listeners: input}, nil, nil, opts)
	out, err := p.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Port != 443 || out[1].Port != 8080 || out[2].Port != 9090 {
		t.Errorf("unexpected sort order: %v %v %v", out[0].Port, out[1].Port, out[2].Port)
	}
}

func TestPipeline_Run_FilterApplied(t *testing.T) {
	input := []ports.Listener{
		makePipelineListener(80, "tcp"),
		makePipelineListener(8080, "tcp"),
	}
	filter := ports.NewFilter(nil, []uint16{8080})
	chain := ports.NewFilterChain(filter)
	p := ports.NewPipeline(&stubScanner{listeners: input}, nil, chain, ports.SortOptions{})
	out, err := p.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 || out[0].Port != 80 {
		t.Errorf("expected only port 80, got %v", out)
	}
}

func TestPipeline_Run_EmptyScanner(t *testing.T) {
	p := ports.NewPipeline(&stubScanner{listeners: []ports.Listener{}}, nil, nil, ports.SortOptions{})
	out, err := p.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty result, got %d listeners", len(out))
	}
}

func TestPipeline_Run_SortDescending(t *testing.T) {
	input := []ports.Listener{
		makePipelineListener(443, "tcp"),
		makePipelineListener(8080, "tcp"),
		makePipelineListener(9090, "tcp"),
	}
	opts := ports.SortOptions{Field: ports.SortFieldPort, Ascending: false}
	p := ports.NewPipeline(&stubScanner{listeners: input}, nil, nil, opts)
	out, err := p.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Port != 9090 || out[1].Port != 8080 || out[2].Port != 443 {
		t.Errorf("unexpected sort order: %v %v %v", out[0].Port, out[1].Port, out[2].Port)
	}
}
