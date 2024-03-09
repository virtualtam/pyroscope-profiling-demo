package command

import (
	"github.com/spf13/cobra"
	"github.com/virtualtam/pyroscope-profiling-demo/services/cook/cmd/cook/importing"
)

var (
	inputFile string
)

func NewImportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import-recipes",
		Short: "Import recipes into the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return importing.ImportRecipes(db, inputFile)
		},
	}

	cmd.Flags().StringVar(
		&inputFile,
		"input",
		"",
		"Input JSON file for recipe and ingredient data",
	)

	return cmd
}
