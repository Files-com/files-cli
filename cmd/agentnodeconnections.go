package cmd

import (
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(AgentNodeConnections())
}

func AgentNodeConnections() *cobra.Command {
	AgentNodeConnections := &cobra.Command{
		Use:  "agent-node-connections [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command agent-node-connections\n\t%v", args[0])
		},
	}
	return AgentNodeConnections
}
