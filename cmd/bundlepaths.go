package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(BundlePaths())
}

func BundlePaths() *cobra.Command {
	BundlePaths := &cobra.Command{
		Use:  "bundle-paths [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command bundle-paths\n\t%v", args[0])
		},
	}
	return BundlePaths
}
