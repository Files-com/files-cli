package cmd

import (
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(AgentNodeInstances())
}

func AgentNodeInstances() *cobra.Command {
	AgentNodeInstances := &cobra.Command{
		Use:  "agent-node-instances [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command agent-node-instances\n\t%v", args[0])
		},
	}
	return AgentNodeInstances
}
