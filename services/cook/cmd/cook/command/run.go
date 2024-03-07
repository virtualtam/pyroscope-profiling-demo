package command

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/virtualtam/pyroscope-profiling-demo/services/cook/cmd/cook/http/api"
	"github.com/virtualtam/pyroscope-profiling-demo/services/cook/pkg/restaurant"
)

const (
	defaultListenAddr string = "0.0.0.0:8080"
)

var (
	listenAddr string
)

// NewRunCommand initializes a CLI command to start the HTTP server.
func NewRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Start the HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().
				Str("log_level", logLevelValue).
				Msg("global: setting up services")

			// Database connection
			dsn := "host=localhost user=cook password=c00k dbname=restaurant port=5432 sslmode=disable"
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Error().Err(err).Msg("failed to open database connection")
			}

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

			// Cook server
			apiServer := api.NewServer(db)

			httpServer := &http.Server{
				Addr:         listenAddr,
				Handler:      apiServer,
				ReadTimeout:  15 * time.Second,
				WriteTimeout: 15 * time.Second,
			}

			log.Info().Str("http_addr", listenAddr).Msg("cook: listening for HTTP requests")
			return httpServer.ListenAndServe()
		},
	}

	cmd.Flags().StringVar(
		&listenAddr,
		"listen-addr",
		defaultListenAddr,
		"Listen to this address (host:port)",
	)

	return cmd
}
