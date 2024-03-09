package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/grafana/pyroscope-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/virtualtam/venom"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/virtualtam/pyroscope-profiling-demo/services/cook/cmd/cook/config"
)

const (
	pyroscopeApplicationName = "demo.cook"

	defaultDatabaseHost string = "localhost"
	defaultDatabasePort uint   = 5432
)

var (
	defaultLogLevelValue string = zerolog.LevelInfoValue
	logLevelValue        string
	logFormat            string

	pyroscopeAddr string

	databaseHost string
	databasePort uint

	db *gorm.DB
)

// NewRootCommand initializes the main CLI entrypoint and common command flags.
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cook",
		Short: "Pyroscope Demo - Cook service",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Configuration file lookup paths
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			homeConfigPath := filepath.Join(home, ".config")

			configPaths := []string{config.DefaultConfigPath, homeConfigPath, "."}

			// Inject global configuration as a pre-run hook
			//
			// This is required to let Viper load environment variables and
			// configuration entries before invoking nested commands.
			v := viper.New()
			if err := venom.InjectTo(v, cmd, config.EnvPrefix, configPaths, config.ConfigName, false); err != nil {
				return err
			}

			// Global logger configuration
			if err := config.SetupGlobalLogger(logFormat, logLevelValue); err != nil {
				return err
			}

			if configFileUsed := v.ConfigFileUsed(); configFileUsed != "" {
				log.Info().Str("config_file", v.ConfigFileUsed()).Msg("configuration: using file")
			} else {
				log.Info().Strs("config_paths", configPaths).Msg("configuration: no file found")
			}

			// Pyroscope live profiling
			if pyroscopeAddr != "" {
				log.Info().
					Str("pyroscope_addr", pyroscopeAddr).
					Str("pyroscope_app", pyroscopeApplicationName).
					Msg("global: enabling live profiling")
				pyroscope.Start(pyroscope.Config{
					ApplicationName: pyroscopeApplicationName,
					Logger:          &config.PyroscopeLogger{},
					ServerAddress:   pyroscopeAddr,
				})
			}

			// Database connection
			dsn := fmt.Sprintf(
				"host=%s user=cook password=c00k dbname=restaurant port=%d sslmode=disable",
				databaseHost,
				databasePort,
			)
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Error().Err(err).Msg("failed to open database connection")
			}

			return nil
		},
	}

	cmd.PersistentFlags().StringVar(
		&logFormat,
		"log-format",
		config.LogFormatConsole,
		fmt.Sprintf("Log format (%s, %s)", config.LogFormatJSON, config.LogFormatConsole),
	)
	cmd.PersistentFlags().StringVar(
		&logLevelValue,
		"log-level",
		defaultLogLevelValue,
		fmt.Sprintf(
			"Log level (%s)",
			strings.Join(config.LogLevelValues, ", "),
		),
	)

	cmd.PersistentFlags().StringVar(
		&pyroscopeAddr,
		"pyroscope-addr",
		"",
		"Pyroscope server address (http://host:port)",
	)

	cmd.PersistentFlags().StringVar(
		&databaseHost,
		"db-host",
		defaultDatabaseHost,
		"Database host",
	)
	cmd.PersistentFlags().UintVar(
		&databasePort,
		"db-port",
		defaultDatabasePort,
		"Database port",
	)

	return cmd
}
