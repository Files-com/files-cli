package cmd

import (
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(FormFields())
}

func FormFields() *cobra.Command {
	FormFields := &cobra.Command{
		Use:  "form-fields [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command form-fields\n\t%v", args[0])
		},
	}
	return FormFields
}
