//go:build !noportable

package cmd

import (
	"fmt"

	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"
)

var (
	AgentCmdStart = &cobra.Command{
		Use:   "start",
		Short: "Start Files.com Agent Windows Service",
		Long: `use:

$ files-cli agent start

Please take a look at the usage below to customize the serving parameters`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStderr(), "Agent v1 is deprecated and has been removed from the CLI. Please use Agent v2.\n")
			return nil
		},
	}
	AgentCmdInstall = &cobra.Command{
		Use:   "install",
		Short: "Install Files.com Agent Windows Service",
		Long: `use:

$ files-cli agent install --config {path-to/config.json}

Please take a look at the usage below to customize the serving parameters`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStderr(), "Agent v1 is deprecated and has been removed from the CLI. Please use Agent v2.\n")
			return nil
		},
	}
	AgentCmdUninstall = &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall Files.com Agent Windows Service",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStderr(), "Agent v1 is deprecated and has been removed from the CLI. Please use Agent v2.\n")
			return nil
		},
	}
	AgentCmdStatus = &cobra.Command{
		Use:   "status",
		Short: "Status of Files.com Agent Windows Service",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStderr(), "Agent v1 is deprecated and has been removed from the CLI. Please use Agent v2.\n")
			return nil
		},
	}
	AgentCmdStop = &cobra.Command{
		Use:   "stop",
		Short: "Stop Files.com Agent Windows Service",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStderr(), "Agent v1 is deprecated and has been removed from the CLI. Please use Agent v2.\n")
			return nil
		},
	}
	AgentCmdRotateLog = &cobra.Command{
		Use:   "rotatelogs",
		Short: "Signal to the running service to rotate the logs",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStderr(), "Agent v1 is deprecated and has been removed from the CLI. Please use Agent v2.\n")
		},
	}
)

func init() {
	if AgentService == nil {
		AgentService = &lib.AgentService{}
	}
	AgentService.AddFlags(AgentCmdInstall.Flags())
	AgentService.AddFlags(AgentCmdStart.Flags())
	AgentCmd.AddCommand(AgentCmdUninstall)
	AgentCmd.AddCommand(AgentCmdInstall)
	AgentCmd.AddCommand(AgentCmdStop)
	AgentCmd.AddCommand(AgentCmdStart)
	AgentCmd.AddCommand(AgentCmdStatus)
	AgentCmd.AddCommand(AgentCmdRotateLog)
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, AgentCmdUninstall.Use)
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, AgentCmdInstall.Use)
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, AgentCmdStop.Use)
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, AgentCmdStart.Use)
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, AgentCmdStatus.Use)
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, AgentCmdRotateLog.Use)
}
