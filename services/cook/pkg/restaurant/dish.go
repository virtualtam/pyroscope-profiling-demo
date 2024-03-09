package restaurant

type Dish struct {
	Model

	Name        string       `json:"name"`
	Price       float64      `json:"price"`
	Ingredients []Ingredient `gorm:"many2many:dish_ingredients;" json:"ingredients"`
}
