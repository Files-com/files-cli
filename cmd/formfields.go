package cmd

import "github.com/spf13/cobra"

var (
	FormFields = &cobra.Command{
		Use:  "form-fields [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func FormFieldsInit() {
}
