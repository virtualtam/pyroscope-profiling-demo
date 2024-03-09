package restaurant

type Menu struct {
	Model

	Dishes       []Dish `gorm:"many2many:menu_dishes;" json:"dishes"`
	RestaurantID uint   `json:"-"`
}
