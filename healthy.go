package capo

import (
	"context"
	"net/http"

	"github.com/tonygcs/gnalog"
)

// Opts provides the custom configuration for each health check. If the methods
// return an error, the server will return a 503.
type Opts struct {
	Health    func(context.Context) error
	Liveness  func(context.Context) error
	Readiness func(context.Context) error
}

// HealthyServer is a data type that implements the 'http.Handler' interface
// with the methods to reply for the server healthy requests.
type HealthyServer struct {
	mux  *http.ServeMux
	opts *Opts
}

// NewHealthyServer create a new instance of the 'HealthyServer'.
//
// Use the 'opts' to provide custom methods for each health check. The server
// will use default methods that will always return a 200 if they are not
// provided. In case a custom method returns an error, the server will return
// the status code 503.
func NewHealthyServer(opts *Opts) http.Handler {
	if opts == nil {
		opts = &Opts{}
	}

	s := &HealthyServer{
		mux:  http.NewServeMux(),
		opts: opts,
	}

	s.mux.HandleFunc("/liveness", s.liveness)
	s.mux.HandleFunc("/readiness", s.readiness)
	s.mux.HandleFunc("/health", s.health)

	return s
}

// ServeHTTP dispatches the request.
func (s *HealthyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *HealthyServer) liveness(w http.ResponseWriter, r *http.Request) {
	var err error

	if s.opts.Liveness != nil {
		ctx := r.Context()
		err = s.opts.Liveness(ctx)
	}

	if err != nil {
		gnalog.Warn("liveness check failed :: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *HealthyServer) readiness(w http.ResponseWriter, r *http.Request) {
	var err error

	if s.opts.Readiness != nil {
		ctx := r.Context()
		err = s.opts.Readiness(ctx)
	}

	if err != nil {
		gnalog.Warn("readiness check failed :: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *HealthyServer) health(w http.ResponseWriter, r *http.Request) {
	var err error

	if s.opts.Health != nil {
		ctx := r.Context()
		err = s.opts.Health(ctx)
	}

	if err != nil {
		gnalog.Warn("health failed :: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}
