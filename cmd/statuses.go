package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Statuses = &cobra.Command{}
)

func StatusesInit() {
	Statuses = &cobra.Command{
		Use:  "statuses [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command statuses\n\t%v", args[0])
		},
	}
}
