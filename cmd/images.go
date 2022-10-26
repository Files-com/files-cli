package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Images())
}

func Images() *cobra.Command {
	Images := &cobra.Command{
		Use:  "images [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command images\n\t%v", args[0])
		},
	}
	return Images
}
