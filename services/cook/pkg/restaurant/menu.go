package restaurant

import "gorm.io/gorm"

type Menu struct {
	gorm.Model

	Dishes       []Dish `gorm:"many2many:menu_dishes;"`
	RestaurantID uint
}
