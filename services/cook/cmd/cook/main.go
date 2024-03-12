package main

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/virtualtam/pyroscope-profiling-demo/services/cook/cmd/cook/command"
)

func main() {
	rootCommand := command.NewRootCommand()

	commands := []*cobra.Command{
		command.NewImportCommand(),
		command.NewMigrateCommand(),
		command.NewCreateRestaurants(),
		command.NewRunCommand(),
	}

	rootCommand.AddCommand(commands...)

	if err := rootCommand.Execute(); err != nil {
		log.Error().Err(err).Msg("failed to execute command")
		os.Exit(1)
	}
}
