package cmd

import (
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(AgentNodes())
}

func AgentNodes() *cobra.Command {
	AgentNodes := &cobra.Command{
		Use:  "agent-nodes [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command agent-nodes\n\t%v", args[0])
		},
	}
	return AgentNodes
}
