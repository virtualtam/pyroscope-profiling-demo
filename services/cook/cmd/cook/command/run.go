package command

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/virtualtam/pyroscope-profiling-demo/services/cook/cmd/cook/http/api"
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

			// Cook server
			apiServer := api.NewServer()

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
