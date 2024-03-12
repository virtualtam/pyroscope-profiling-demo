package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/grafana/pyroscope-go"
	"github.com/rs/zerolog/log"
	"github.com/virtualtam/pyroscope-profiling-demo/services/cook/pkg/restaurant"
	"gorm.io/gorm"
)

const (
	v2RedisTTL time.Duration = 30 * time.Second
)

func (s *Server) registerV2API() {
	s.router.Route("/api/v2", func(r chi.Router) {
		r.Route("/restaurant/{restaurantID}", func(r chi.Router) {
			r.Get("/menu", s.v2getRestaurantMenu)
		})
	})
}

func (s *Server) v2getRestaurantMenu(w http.ResponseWriter, r *http.Request) {
	pyroscope.TagWrapper(context.Background(), pyroscope.Labels("api", "v2"), func(ctx context.Context) {
		restaurantID := chi.URLParam(r, "restaurantID")

		var menu restaurant.Menu

		redisKey := fmt.Sprintf("cook:menu:%s", restaurantID)

		cachedMenu, err := s.redisClient.Get(ctx, redisKey).Result()
		if err == nil {
			// cache hit: return the cached object
			if err := json.Unmarshal([]byte(cachedMenu), &menu); err != nil {
				log.Error().Err(err).Str("key", redisKey).Msg("redis: failed to decode value")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}

		} else {
			// cache miss: retrieve the object from the database
			menu, err = s.restaurantService.Menu(restaurantID)

			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Warn().Str("id", restaurantID).Msg("restaurant: not found")
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}

			if err != nil {
				log.Error().Str("id", restaurantID).Msg("restaurant: failed to query database")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// cache the object (with TTL)
			menuBytes, err := json.Marshal(menu)
			if err != nil {
				log.Error().Str("id", restaurantID).Msg("restaurant: failed to encode menu")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if err := s.redisClient.Set(ctx, redisKey, menuBytes, v2RedisTTL).Err(); err != nil {
				log.Error().Str("id", restaurantID).Str("key", redisKey).Msg("restaurant: failed to cache with redis")
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		render.JSON(w, r, menu)
	})
}
