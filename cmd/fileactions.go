package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	FileActions = &cobra.Command{}
)

func FileActionsInit() {
	FileActions = &cobra.Command{
		Use:  "file-actions [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command file-actions\n\t%v", args[0])
		},
	}
}
