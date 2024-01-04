package apiserver

import (
	"bufio"
	"context"
	"expvar"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

const (
	defaultNamespace = "mev_commit_oracle"
)

// Service wraps http.Server with additional functionality for metrics and
// other common middlewares.
type Service struct {
	metricsRegistry *prometheus.Registry
	router          *http.ServeMux
	srv             *http.Server
}

// New creates a new Service.
func New() *Service {
	srv := &Service{
		router:          http.NewServeMux(),
		metricsRegistry: newMetrics(),
	}

	srv.registerDebugEndpoints()
	return srv
}

func (a *Service) registerDebugEndpoints() {
	// register metrics handler
	a.router.Handle("/metrics", promhttp.HandlerFor(a.metricsRegistry, promhttp.HandlerOpts{}))

	// register pprof handlers
	a.router.Handle(
		"/debug/pprof",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u := r.URL
			u.Path += "/"
			http.Redirect(w, r, u.String(), http.StatusPermanentRedirect)
		}),
	)
	a.router.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	a.router.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	a.router.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	a.router.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	a.router.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	a.router.Handle("/debug/pprof/{profile}", http.HandlerFunc(pprof.Index))
	a.router.Handle("/debug/vars", expvar.Handler())
}

func newMetrics() (r *prometheus.Registry) {
	r = prometheus.NewRegistry()

	// register standard metrics
	r.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{
			Namespace: defaultNamespace,
		}),
		collectors.NewGoCollector(),
	)

	return r
}

func (a *Service) Start(addr string) <-chan struct{} {
	log.Info().Msg("starting api server")

	srv := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			recorder := &responseStatusRecorder{ResponseWriter: w}

			start := time.Now()
			a.router.ServeHTTP(recorder, req)
			log.Info().
				Int("status", recorder.status).
				Str("method", req.Method).
				Str("path", req.URL.Path).
				Dur("duration", time.Since(start)).
				Msg("api access")
		}),
	}
	a.srv = srv

	done := make(chan struct{})
	go func() {
		defer close(done)

		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("api server failed")
		}
	}()

	return done
}

func (a *Service) Stop() error {
	log.Info().Msg("stopping api server")
	if a.srv == nil {
		return nil
	}
	return a.srv.Shutdown(context.Background())
}

// RegisterMetricsCollectors registers prometheus collectors.
func (a *Service) RegisterMetricsCollectors(cs ...prometheus.Collector) {
	a.metricsRegistry.MustRegister(cs...)
}

type responseStatusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *responseStatusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// Hijack implements http.Hijacker.
func (r *responseStatusRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.ResponseWriter.(http.Hijacker).Hijack()
}

// Flush implements http.Flusher.
func (r *responseStatusRecorder) Flush() {
	r.ResponseWriter.(http.Flusher).Flush()
}

// Push implements http.Pusher.
func (r *responseStatusRecorder) Push(target string, opts *http.PushOptions) error {
	return r.ResponseWriter.(http.Pusher).Push(target, opts)
}
