package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	FormFields = &cobra.Command{}
)

func FormFieldsInit() {
	FormFields = &cobra.Command{
		Use:  "form-fields [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command form-fields\n\t%v", args[0])
		},
	}
}
