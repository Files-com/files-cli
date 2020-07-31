package cmd

import "github.com/spf13/cobra"

var (
	Autos = &cobra.Command{
		Use:  "autos [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func AutosInit() {
}
