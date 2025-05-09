package cmd

import (
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Autos())
}

func Autos() *cobra.Command {
	Autos := &cobra.Command{
		Use:  "autos [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command autos\n\t%v", args[0])
		},
	}
	return Autos
}
