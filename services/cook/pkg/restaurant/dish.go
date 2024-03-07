package restaurant

import "gorm.io/gorm"

type Dish struct {
	gorm.Model

	Name        string
	Price       float64
	Ingredients []Ingredient `gorm:"many2many:dish_ingredients;"`
}
