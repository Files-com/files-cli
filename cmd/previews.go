package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Previews = &cobra.Command{}
)

func PreviewsInit() {
	Previews = &cobra.Command{
		Use:  "previews [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command previews\n\t%v", args[0])
		},
	}
}
