package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/virtualtam/pyroscope-profiling-demo/services/cook/cmd/cook/http/middleware"
	"github.com/virtualtam/pyroscope-profiling-demo/services/cook/pkg/restaurant"
)

var _ http.Handler = &Server{}

// Server exposes the HTTP API.
type Server struct {
	db          *gorm.DB
	redisClient *redis.Client
	router      *chi.Mux

	restaurantService *restaurant.Service
}

// NewServer initializes and returns a new Server.
func NewServer(db *gorm.DB, redisClient *redis.Client) *Server {
	restaurantService := restaurant.NewService(db)

	s := &Server{
		db:                db,
		redisClient:       redisClient,
		router:            chi.NewRouter(),
		restaurantService: restaurantService,
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
	s.registerV1API()
	s.registerV2API()
}
