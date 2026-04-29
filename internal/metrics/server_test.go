package metrics

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func newTestRegistry(t *testing.T) *prometheus.Registry {
	t.Helper()
	return prometheus.NewRegistry()
}

func TestNewCounters_Registers(t *testing.T) {
	reg := newTestRegistry(t)
	c := NewCounters(reg)
	if c == nil {
		t.Fatal("expected non-nil counters")
	}
	c.ScansTotal.Inc()
	c.ListenersFound.Set(5)
	c.AlertsTotal.WithLabelValues("appeared").Inc()
}

func TestServer_StartAndShutdown(t *testing.T) {
	reg := newTestRegistry(t)
	NewCounters(reg)
	log := slog.Default()

	srv := NewServer("127.0.0.1:0", "/metrics", reg, log)
	// Override with a free port for testing.
	srv.httpServer.Addr = "127.0.0.1:19091"

	srv.Start()
	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get("http://127.0.0.1:19091/metrics")
	if err != nil {
		t.Fatalf("could not reach metrics endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "portwatch_scans_total") {
		t.Error("expected portwatch_scans_total in metrics output")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		t.Errorf("shutdown error: %v", err)
	}
}

func TestServer_Start_Idempotent(t *testing.T) {
	reg := newTestRegistry(t)
	log := slog.Default()
	srv := NewServer("127.0.0.1:19092", "/metrics", reg, log)
	srv.Start()
	srv.Start() // second call must not panic or double-bind
	time.Sleep(30 * time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
