package command

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/virtualtam/pyroscope-profiling-demo/services/cook/pkg/restaurant"
)

func NewMigrateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Apply database migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Msg("applying database migrations")
			if err := db.AutoMigrate(
				&restaurant.Ingredient{},
				&restaurant.Dish{},
				&restaurant.Restaurant{},
				&restaurant.Menu{},
			); err != nil {
				log.Error().Err(err).Msg("failed to migrate database")
				return err
			}

			return nil
		},
	}
}
