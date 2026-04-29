package healthcheck

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/user/portwatch/internal/config"
)

// status holds the current liveness state.
type status struct {
	OK      bool   `json:"ok"`
	Version string `json:"version,omitempty"`
}

// Server is a lightweight HTTP health-check endpoint.
type Server struct {
	cfg     config.HealthCheckConfig
	log     *slog.Logger
	healthy atomic.Bool
	server  *http.Server
}

// NewServer creates a health-check server using the provided config.
func NewServer(cfg config.HealthCheckConfig, log *slog.Logger) *Server {
	s := &Server{cfg: cfg, log: log}
	s.healthy.Store(true)
	return s
}

// SetHealthy updates the liveness state reported by the endpoint.
func (s *Server) SetHealthy(ok bool) { s.healthy.Store(ok) }

// Start begins listening in the background. It returns immediately.
// The server shuts down when ctx is cancelled.
func (s *Server) Start(ctx context.Context) error {
	if !s.cfg.Enabled {
		return nil
	}

	mux := http.NewServeMux()
	mux.HandleFunc(s.cfg.Path, s.handleHealth)

	s.server = &http.Server{
		Addr:         s.cfg.Addr,
		Handler:      mux,
		ReadTimeout:  s.cfg.Timeout,
		WriteTimeout: s.cfg.Timeout,
	}

	go func() {
		s.log.Info("healthcheck server starting", "addr", s.cfg.Addr, "path", s.cfg.Path)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("healthcheck server error", "err", err)
		}
	}()

	go func() {
		<-ctx.Done()
		shutCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = s.server.Shutdown(shutCtx)
		s.log.Info("healthcheck server stopped")
	}()

	return nil
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	st := status{OK: s.healthy.Load()}
	w.Header().Set("Content-Type", "application/json")
	if !st.OK {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	_ = json.NewEncoder(w).Encode(st)
}
