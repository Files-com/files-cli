package cmd

import (
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(RemoteServerConfigurationFiles())
}

func RemoteServerConfigurationFiles() *cobra.Command {
	RemoteServerConfigurationFiles := &cobra.Command{
		Use:  "remote-server-configuration-files [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command remote-server-configuration-files\n\t%v", args[0])
		},
	}
	return RemoteServerConfigurationFiles
}
