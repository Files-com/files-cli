package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func VersionCmd(Version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "version",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("files-cli version %v\n", Version)
		},
	}

	return cmd
}
