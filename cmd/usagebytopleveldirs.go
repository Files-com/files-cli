package cmd

import (
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(UsageByTopLevelDirs())
}

func UsageByTopLevelDirs() *cobra.Command {
	UsageByTopLevelDirs := &cobra.Command{
		Use:  "usage-by-top-level-dirs [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command usage-by-top-level-dirs\n\t%v", args[0])
		},
	}
	return UsageByTopLevelDirs
}
