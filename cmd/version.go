package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	VersionCmd *cobra.Command
)

func init() {
	VersionCmd = &cobra.Command{
		Use:     "version",
		Aliases: []string{"version"},
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("files-cli version %v\n", RootCmd.Version)
		},
	}

	RootCmd.AddCommand(VersionCmd)
}
