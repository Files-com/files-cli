package cmd

import (
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Previews())
}

func Previews() *cobra.Command {
	Previews := &cobra.Command{
		Use:  "previews [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command previews\n\t%v", args[0])
		},
	}
	return Previews
}
