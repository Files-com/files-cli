package cmd

import (
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Errors())
}

func Errors() *cobra.Command {
	Errors := &cobra.Command{
		Use:  "errors [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command errors\n\t%v", args[0])
		},
	}
	return Errors
}
