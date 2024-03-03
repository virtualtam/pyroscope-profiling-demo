package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"

	"github.com/virtualtam/pyroscope-profiling-demo/services/cook/cmd/cook/http/middleware"
)

var _ http.Handler = &Server{}

// Server exposes the HTTP API.
type Server struct {
	router *chi.Mux
}

// NewServer initializes and returns a new Server.
func NewServer() *Server {
	s := &Server{
		router: chi.NewRouter(),
	}

	s.registerHandlers()

	return s
}

// ServeHTTP satisfies the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// registerHandlers registers all middleware and handlers for the HTTP API.
func (s *Server) registerHandlers() {
	// Global middleware
	s.router.Use(chimiddleware.RequestID)
	s.router.Use(chimiddleware.RealIP)

	// Structured logging
	s.router.Use(hlog.NewHandler(log.Logger), hlog.AccessHandler(middleware.AccessLogger))

	// API
	s.router.Get("/", http.NotFound)
}
