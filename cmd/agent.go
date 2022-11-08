//go:build !noportable && !windows

package cmd

import (
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
			AgentService.Config = *Profile(cmd).Config
			err := AgentService.Init(cmd.Context())
			if err != nil {
				return err
			}
			return AgentService.Start()
		},
	}
)

func init() {
	AgentService = &lib.AgentService{}
	AgentCmd.CompletionOptions.DisableDefaultCmd = true
	AgentService.AddFlags(AgentCmd.Flags())
	RootCmd.AddCommand(AgentCmd)
	AgentCmd.Hidden = true
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, AgentCmd.Use)
}
