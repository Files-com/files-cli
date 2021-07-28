package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Autos = &cobra.Command{}
)

func AutosInit() {
	Autos = &cobra.Command{
		Use:  "autos [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command autos\n\t%v", args[0])
		},
	}
}
