package cmd

import (
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(AccountLineItems())
}

func AccountLineItems() *cobra.Command {
	AccountLineItems := &cobra.Command{
		Use:  "account-line-items [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command account-line-items\n\t%v", args[0])
		},
	}
	return AccountLineItems
}
