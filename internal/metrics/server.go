package metrics

import (
	"context"
	"log/slog"
	"net/http"
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Counters holds the Prometheus counters exposed by portwatch.
type Counters struct {
	ScansTotal     prometheus.Counter
	ListenersFound prometheus.Gauge
	AlertsTotal    *prometheus.CounterVec
}

// Server wraps an HTTP server that exposes Prometheus metrics.
type Server struct {
	httpServer *http.Server
	log        *slog.Logger
	running    atomic.Bool
}

// NewCounters registers and returns the default portwatch metric set.
func NewCounters(reg prometheus.Registerer) *Counters {
	c := &Counters{
		ScansTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "portwatch_scans_total",
			Help: "Total number of port scans performed.",
		}),
		ListenersFound: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "portwatch_listeners_found",
			Help: "Number of listeners observed in the last scan.",
		}),
		AlertsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "portwatch_alerts_total",
			Help: "Total alerts fired, partitioned by event type.",
		}, []string{"event"}),
	}
	reg.MustRegister(c.ScansTotal, c.ListenersFound, c.AlertsTotal)
	return c
}

// NewServer creates a metrics HTTP server on the given address and path.
func NewServer(addr, path string, reg prometheus.Gatherer, log *slog.Logger) *Server {
	mux := http.NewServeMux()
	mux.Handle(path, promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	return &Server{
		httpServer: &http.Server{Addr: addr, Handler: mux},
		log:        log,
	}
}

// Start begins serving metrics in a background goroutine.
func (s *Server) Start() {
	if s.running.Swap(true) {
		return
	}
	go func() {
		s.log.Info("metrics server listening", "addr", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("metrics server error", "err", err)
		}
	}()
}

// Shutdown gracefully stops the metrics server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
