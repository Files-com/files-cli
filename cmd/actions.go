package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Actions = &cobra.Command{}
)

func ActionsInit() {
	Actions = &cobra.Command{
		Use:  "actions [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command actions\n\t%v", args[0])
		},
	}
}
