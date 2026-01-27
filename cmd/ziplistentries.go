package cmd

import (
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(ZipListEntries())
}

func ZipListEntries() *cobra.Command {
	ZipListEntries := &cobra.Command{
		Use:  "zip-list-entries [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command zip-list-entries\n\t%v", args[0])
		},
	}
	return ZipListEntries
}
