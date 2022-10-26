package cmd

import (
	"fmt"

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
			return fmt.Errorf("invalid command previews\n\t%v", args[0])
		},
	}
	return Previews
}
