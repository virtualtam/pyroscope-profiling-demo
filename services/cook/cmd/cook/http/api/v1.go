package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gorm.io/gorm"
)

func (s *Server) registerV1API() {
	s.router.Route("/api/v1", func(r chi.Router) {
		r.Route("/restaurant/{restaurantID}", func(r chi.Router) {
			r.Get("/menu", s.v1getRestaurantMenu)
		})
	})
}

func (s *Server) v1getRestaurantMenu(w http.ResponseWriter, r *http.Request) {
	restaurantID := chi.URLParam(r, "restaurantID")

	menu, err := s.restaurantService.Menu(restaurantID)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, menu)
}
