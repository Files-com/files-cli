package cmd

import "github.com/spf13/cobra"

var (
	Statuses = &cobra.Command{
		Use:  "statuses [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func StatusesInit() {
}
