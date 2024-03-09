package command

import (
	"github.com/spf13/cobra"

	"github.com/virtualtam/pyroscope-profiling-demo/services/cook/cmd/cook/data"
)

const (
	defaultNRestaurants int = 10
)

var (
	nRestaurants int
)

func NewCreateRestaurants() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-restaurants",
		Short: "Create new restaurants and their menu",
		RunE: func(cmd *cobra.Command, args []string) error {
			return data.CreateRestaurants(db, nRestaurants)
		},
	}

	cmd.Flags().IntVar(
		&nRestaurants,
		"number",
		defaultNRestaurants,
		"Number of restaurants to create",
	)

	return cmd
}
