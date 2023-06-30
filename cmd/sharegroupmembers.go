package cmd

import (
	"fmt"

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
			return fmt.Errorf("invalid command share-group-members\n\t%v", args[0])
		},
	}
	return ShareGroupMembers
}
