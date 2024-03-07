package restaurant

import "gorm.io/gorm"

type Ingredient struct {
	gorm.Model

	Name string
}
