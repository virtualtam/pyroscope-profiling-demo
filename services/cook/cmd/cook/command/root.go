package command

import (
	"context"
	"fmt"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strings"

	"github.com/grafana/pyroscope-go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/virtualtam/venom"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/virtualtam/pyroscope-profiling-demo/services/cook/cmd/cook/config"
	"github.com/virtualtam/pyroscope-profiling-demo/services/cook/cmd/cook/http/pprofiler"
)

const (
	pyroscopeApplicationName = "demo.cook"

	defaultDatabaseHost string = "localhost"
	defaultDatabasePort uint   = 5432

	defaultRedisHost     string = "localhost"
	defaultRedisPort     uint   = 6379
	defaultRedisDatabase uint   = 0
)

var (
	defaultLogLevelValue string = zerolog.LevelInfoValue
	logLevelValue        string
	logFormat            string

	pprofAddr     string
	pyroscopeAddr string

	databaseHost string
	databasePort uint

	redisHost     string
	redisPort     uint
	redisDatabase uint

	db          *gorm.DB
	redisClient *redis.Client
)

// NewRootCommand initializes the main CLI entrypoint and common command flags.
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "cook",
		Short:        "Pyroscope Demo - Cook service",
		SilenceUsage: true,
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

			// Go pprof server
			if pprofAddr != "" {
				pprofServer := pprofiler.NewServer(pprofAddr)
				log.Info().
					Str("pprof_addr", pprofAddr).
					Msg("global: enabling pprof profiling")

				go func() {
					if err := pprofServer.ListenAndServe(); err != nil {
						log.Error().Err(err).Msg("failed to start pprof server")
					}
				}()
			}

			// Pyroscope live profiling
			if pyroscopeAddr != "" {
				log.Info().
					Str("pyroscope_addr", pyroscopeAddr).
					Str("pyroscope_app", pyroscopeApplicationName).
					Msg("global: enabling pyroscope live profiling")

				if _, err := pyroscope.Start(pyroscope.Config{
					ApplicationName: pyroscopeApplicationName,
					Logger:          &config.PyroscopeLogger{},
					ProfileTypes: []pyroscope.ProfileType{
						pyroscope.ProfileAllocObjects,
						pyroscope.ProfileAllocSpace,
						pyroscope.ProfileBlockCount,
						pyroscope.ProfileBlockDuration,
						pyroscope.ProfileCPU,
						pyroscope.ProfileGoroutines,
						pyroscope.ProfileInuseObjects,
						pyroscope.ProfileInuseSpace,
						pyroscope.ProfileMutexCount,
						pyroscope.ProfileMutexDuration,
					},
					ServerAddress: pyroscopeAddr,
				}); err != nil {
					log.Error().Err(err).Msg("global: failed to start pyroscope profiler")
					return err
				}
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

			// Redis KV store connection
			redisOpts := &redis.Options{
				Addr: fmt.Sprintf("%s:%d", redisHost, redisPort),
				DB:   int(redisDatabase),
			}
			redisClient = redis.NewClient(redisOpts)

			if err := redisClient.Ping(context.Background()).Err(); err != nil {
				log.Error().Err(err).Msg("failed to open redis connection")
				return err
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
		&pprofAddr,
		"pprof-addr",
		"",
		"pprof listen address (host:port)",
	)
	cmd.PersistentFlags().StringVar(
		&pyroscopeAddr,
		"pyroscope-addr",
		"",
		"Pyroscope server address (http://host:port)",
	)
	cmd.MarkFlagsMutuallyExclusive("pprof-addr", "pyroscope-addr")

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

	cmd.PersistentFlags().StringVar(
		&redisHost,
		"redis-host",
		defaultRedisHost,
		"Redis host",
	)
	cmd.PersistentFlags().UintVar(
		&redisPort,
		"redis-port",
		defaultRedisPort,
		"Redis port",
	)
	cmd.PersistentFlags().UintVar(
		&redisDatabase,
		"redis-db",
		defaultRedisDatabase,
		"Redis database",
	)

	return cmd
}
