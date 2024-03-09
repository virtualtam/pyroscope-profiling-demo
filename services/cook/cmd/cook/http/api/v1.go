package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"

	"github.com/virtualtam/pyroscope-profiling-demo/services/cook/pkg/restaurant"
)

func (s *Server) registerV1API() {
	s.router.Route("/api/v1", func(r chi.Router) {
		r.Route("/restaurant/{restaurantID}", func(r chi.Router) {
			r.Get("/menu", s.getMenu)
		})
	})
}

func (s *Server) getMenu(w http.ResponseWriter, r *http.Request) {
	restaurantID := chi.URLParam(r, "restaurantID")

	var rest restaurant.Restaurant

	err := s.db.Model(&restaurant.Restaurant{}).
		Preload("Menu.Dishes.Ingredients").
		First(&rest, restaurantID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, rest.Menu)
}
