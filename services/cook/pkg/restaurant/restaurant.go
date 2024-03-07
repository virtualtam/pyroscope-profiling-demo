package restaurant

import "gorm.io/gorm"

type Restaurant struct {
	gorm.Model

	Name string
	Menu Menu
}
