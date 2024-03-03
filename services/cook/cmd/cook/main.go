package main

import (
	"github.com/spf13/cobra"

	"github.com/virtualtam/pyroscope-profiling-demo/services/cook/cmd/cook/command"
)

func main() {
	rootCommand := command.NewRootCommand()

	commands := []*cobra.Command{
		command.NewRunCommand(),
	}

	rootCommand.AddCommand(commands...)

	cobra.CheckErr(rootCommand.Execute())
}
