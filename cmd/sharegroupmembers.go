package cmd

import (
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(ShareGroupMembers())
}

func ShareGroupMembers() *cobra.Command {
	ShareGroupMembers := &cobra.Command{
		Use:  "share-group-members [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command share-group-members\n\t%v", args[0])
		},
	}
	return ShareGroupMembers
}
