package healthcheck

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func newTestServer(t *testing.T, addr string) *Server {
	t.Helper()
	cfg := config.HealthCheckConfig{
		Enabled: true,
		Addr:    addr,
		Path:    "/healthz",
		Timeout: 5 * time.Second,
	}
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	return NewServer(cfg, log)
}

func TestServer_HandleHealth_OK(t *testing.T) {
	srv := newTestServer(t, ":0")
	srv.SetHealthy(true)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	srv.handleHealth(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	var st status
	if err := json.NewDecoder(rr.Body).Decode(&st); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if !st.OK {
		t.Error("expected ok=true")
	}
}

func TestServer_HandleHealth_Unhealthy(t *testing.T) {
	srv := newTestServer(t, ":0")
	srv.SetHealthy(false)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	srv.handleHealth(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", rr.Code)
	}
	var st status
	if err := json.NewDecoder(rr.Body).Decode(&st); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if st.OK {
		t.Error("expected ok=false")
	}
}

func TestServer_Disabled_StartNoOp(t *testing.T) {
	cfg := config.DefaultHealthCheckConfig()
	cfg.Enabled = false
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	srv := NewServer(cfg, log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := srv.Start(ctx); err != nil {
		t.Errorf("expected no error for disabled server, got: %v", err)
	}
}

func TestServer_StartAndShutdown(t *testing.T) {
	srv := newTestServer(t, "127.0.0.1:19110")
	ctx, cancel := context.WithCancel(context.Background())

	if err := srv.Start(ctx); err != nil {
		t.Fatalf("Start error: %v", err)
	}
	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get("http://127.0.0.1:19110/healthz")
	if err != nil {
		t.Fatalf("GET error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	cancel()
	time.Sleep(50 * time.Millisecond)
}
