package cmd

import "github.com/spf13/cobra"

var (
	Previews = &cobra.Command{}
)

func PreviewsInit() {
	Previews = &cobra.Command{
		Use:  "previews [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
}
