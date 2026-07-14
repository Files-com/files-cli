package cmd

import (
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(AutomationAuthoringSchemas())
}

func AutomationAuthoringSchemas() *cobra.Command {
	AutomationAuthoringSchemas := &cobra.Command{
		Use:  "automation-authoring-schemas [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command automation-authoring-schemas\n\t%v", args[0])
		},
	}
	return AutomationAuthoringSchemas
}
