package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Images = &cobra.Command{}
)

func ImagesInit() {
	Images = &cobra.Command{
		Use:  "images [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command images\n\t%v", args[0])
		},
	}
}
