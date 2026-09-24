package cmd

import (
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(PartnerConnections())
}

func PartnerConnections() *cobra.Command {
	PartnerConnections := &cobra.Command{
		Use:   "partner-connections [command]",
		Short: "",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command partner-connections\n\t%v", args[0])
		},
	}
	return PartnerConnections
}
