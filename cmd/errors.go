package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Errors())
}

func Errors() *cobra.Command {
	Errors := &cobra.Command{
		Use:  "errors [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command errors\n\t%v", args[0])
		},
	}
	return Errors
}
