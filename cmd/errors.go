package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Errors = &cobra.Command{}
)

func ErrorsInit() {
	Errors = &cobra.Command{
		Use:  "errors [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command errors\n\t%v", args[0])
		},
	}
}
