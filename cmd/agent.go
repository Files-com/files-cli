//go:build !noportable

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
			err := AgentInt(cmd)
			if err != nil {
				return err
			}
			err = AgentService.Start(true)
			if err != nil {
				return err
			}
			AgentService.Wait()
			return AgentService.Error
		},
	}
)

func AgentInt(cmd *cobra.Command) error {
	AgentService.Config = Profile(cmd).Config
	return AgentService.Init(cmd.Context())
}

func init() {
	AgentService = &lib.AgentService{}
	AgentCmd.CompletionOptions.DisableDefaultCmd = true
	AgentService.AddFlags(AgentCmd.Flags())
	RootCmd.AddCommand(AgentCmd)
	AgentCmd.Hidden = true
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, AgentCmd.Use)
}
