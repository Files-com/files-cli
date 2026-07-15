package cmd

import (
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(DirectConnectionInfos())
}

func DirectConnectionInfos() *cobra.Command {
	DirectConnectionInfos := &cobra.Command{
		Use:  "direct-connection-infos [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command direct-connection-infos\n\t%v", args[0])
		},
	}
	return DirectConnectionInfos
}
