package cmd

import (
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Actions())
}

func Actions() *cobra.Command {
	Actions := &cobra.Command{
		Use:  "actions [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command actions\n\t%v", args[0])
		},
	}
	return Actions
}
