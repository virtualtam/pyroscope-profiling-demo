package restaurant

import "gorm.io/gorm"

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Menu(restaurantID string) (Menu, error) {
	var rest Restaurant

	err := s.db.Model(&Restaurant{}).
		Preload("Menu.Dishes.Ingredients").
		First(&rest, restaurantID).Error

	if err != nil {
		return Menu{}, err
	}

	return rest.Menu, nil
}
