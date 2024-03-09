package data

import (
	"math/rand"

	"github.com/jaswdr/faker/v2"
	"github.com/rs/zerolog/log"
	"github.com/virtualtam/pyroscope-profiling-demo/services/cook/pkg/restaurant"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateRestaurants(db *gorm.DB, nRestaurants int) error {
	fake := faker.New()

	for i := 0; i < nRestaurants; i++ {
		nDishes := 5 + rand.Intn(10)

		var dishes []restaurant.Dish

		tx := db.Clauses(
			clause.OrderBy{
				Expression: clause.Expr{
					SQL: "RANDOM()",
				},
			},
		).
			Limit(nDishes).
			Find(&dishes)

		if tx.Error != nil {
			log.Error().Err(tx.Error).Int("dish_count", nDishes).Msg("failed to retrieve random dishes")
			return tx.Error
		}

		menu := restaurant.Menu{
			Dishes: dishes,
		}

		rest := restaurant.Restaurant{
			Name: fake.Company().Name(),
			Menu: menu,
		}

		tx = db.Create(&rest)
		if tx.Error != nil {
			log.Error().Err(tx.Error).Msg("failed to create menu and restaurant")
			return tx.Error
		}
	}

	log.Info().Int("restaurant_count", nRestaurants).Msg("restaurants created")

	return nil
}
