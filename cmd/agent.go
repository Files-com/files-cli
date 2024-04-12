//go:build !noportable

package cmd

import (
	"fmt"

	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"
)

var (
	AgentService *lib.AgentService
	AgentCmd     = &cobra.Command{
		Use:   "agent",
		Short: "Start Files.com Agent",
		Long: `use:

$ files-cli agent --config {path-to/config.json}

Please take a look at the usage below to customize the serving parameters`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStderr(), "Agent v1 is deprecated and has been removed from the CLI. Please use Agent v2.\n")
			return nil
		},
	}
)

func init() {
	if AgentService == nil {
		AgentService = &lib.AgentService{}
	}
	AgentCmd.CompletionOptions.DisableDefaultCmd = true
	AgentService.AddFlags(AgentCmd.Flags())
	RootCmd.AddCommand(AgentCmd)
	AgentCmd.Hidden = true
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, AgentCmd.Use)
}
