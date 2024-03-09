package data

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/virtualtam/pyroscope-profiling-demo/services/cook/pkg/restaurant"
	"gorm.io/gorm"
)

type Recipe struct {
	Name        string   `json:"name"`
	Ingredients []string `json:"ingredients"`
}

func ImportRecipes(db *gorm.DB, jsonInput string) error {
	data, err := os.ReadFile(jsonInput)
	if err != nil {
		log.Error().Err(err).Msg("failed to read input file")
		return err
	}

	var recipes []Recipe

	if err := json.Unmarshal(data, &recipes); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal JSON data")
		return err
	}

	var importedRecipes int

	for _, recipe := range recipes {
		var ingredients []restaurant.Ingredient

		for _, ingredientName := range recipe.Ingredients {
			var ingredient restaurant.Ingredient

			err := db.Where("name = ?", ingredientName).First(&ingredient).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ingredient = restaurant.Ingredient{
					Name: ingredientName,
				}

				tx := db.Create(&ingredient)
				if tx.Error != nil {
					log.Error().Err(tx.Error).Msg("failed to create ingredient")
					return tx.Error
				}

			} else if err != nil {
				log.Error().Err(err).Msg("failed to retrieve ingredient")
				return err
			}

			ingredients = append(ingredients, ingredient)
		}

		price := float64(len(recipe.Name)) + float64(len(recipe.Ingredients))/100

		dish := restaurant.Dish{
			Name:        recipe.Name,
			Price:       price,
			Ingredients: ingredients,
		}

		tx := db.Create(&dish)
		if tx.Error != nil {
			log.Error().Err(tx.Error).Msg("failed to create dish")
			return tx.Error
		}

		importedRecipes++
	}

	log.Info().Int("recipe_count", importedRecipes).Msg("recipes imported")

	return nil
}
