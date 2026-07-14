package cmd

import (
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(AutomationExecutionNodes())
}

func AutomationExecutionNodes() *cobra.Command {
	AutomationExecutionNodes := &cobra.Command{
		Use:  "automation-execution-nodes [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command automation-execution-nodes\n\t%v", args[0])
		},
	}
	return AutomationExecutionNodes
}
