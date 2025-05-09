package cmd

import (
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Statuses())
}

func Statuses() *cobra.Command {
	Statuses := &cobra.Command{
		Use:  "statuses [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command statuses\n\t%v", args[0])
		},
	}
	return Statuses
}
