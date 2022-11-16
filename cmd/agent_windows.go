//go:build !noportable

package cmd

import (
	"os"

	"github.com/drakkan/sftpgo/v2/logger"

	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"
)

var (
	AgentCmdStart = &cobra.Command{
		Use:   "start",
		Short: "Start Files.com Agent Windows Service",
		Long: `use:

$ files-cli agent start

Please take a look at the usage below to customize the serving parameters`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := AgentInt(cmd)
			if err != nil {
				return err
			}
			winService := lib.WindowsService{
				Service: AgentService,
			}

			err = winService.RunService()
			logger.Debug("files-cli", "", "Start Cmd - err: %v", err)

			if err != nil {
				return fmt.Errorf("Error starting service: %v\r\n", err)
			} else {
				fmt.Printf("Service started!\r\n")
			}
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
			err := AgentInt(cmd)
			if err != nil {
				return err
			}
			winService := lib.WindowsService{
				Service: AgentService,
			}

			err = winService.Install(AgentService.ServiceArgs()...)
			if err != nil {
				return fmt.Errorf("Error installing service: %v\r\n", err)
			}
			fmt.Printf("Service installed!\r\n")

			return nil
		},
	}
	AgentCmdUninstall = &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall Files.com Agent Windows Service",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := AgentInt(cmd)
			if err != nil {
				return err
			}
			s := lib.WindowsService{
				Service: AgentService,
			}
			err = s.Uninstall()
			if err != nil {
				return fmt.Errorf("Error removing service: %v\r\n", err)
			} else {
				fmt.Printf("Service uninstalled\r\n")
			}

			return nil
		},
	}
	AgentCmdStatus = &cobra.Command{
		Use:   "status",
		Short: "Status of Files.com Agent Windows Service",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := AgentInt(cmd)
			if err != nil {
				return err
			}
			s := lib.WindowsService{
				Service: AgentService,
			}
			status, err := s.Status()
			if err != nil {
				return fmt.Errorf("Error querying service status: %v\r\n", err)
			} else {
				fmt.Printf("Service status: %#v\r\n", status.String())
			}

			return nil
		},
	}
	AgentCmdStop = &cobra.Command{
		Use:   "stop",
		Short: "Stop Files.com Agent Windows Service",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := AgentInt(cmd)
			if err != nil {
				return err
			}
			s := lib.WindowsService{
				Service: AgentService,
			}
			err = s.Stop()
			if err != nil {
				return fmt.Errorf("Error stopping service: %v\r\n", err)
			} else {
				fmt.Printf("Service stopped!\r\n")
			}
			return nil
		},
	}
	AgentCmdReload = &cobra.Command{
		Use:   "reload",
		Short: "Reload Files.com Agent Windows Service",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := AgentInt(cmd)
			if err != nil {
				return err
			}
			s := lib.WindowsService{
				Service: AgentService,
			}
			err = s.Reload()
			if err != nil {
				return fmt.Errorf("Error sending reload signal: %v\r\n", err)
			} else {
				fmt.Printf("Reload signal sent!\r\n")
			}
			return nil
		},
	}
	AgentCmdRotateLog = &cobra.Command{
		Use:   "rotatelogs",
		Short: "Signal to the running service to rotate the logs",
		Run: func(cmd *cobra.Command, args []string) {
			err := AgentInt(cmd)
			if err != nil {
				fmt.Printf("Error sending rotate log file signal to the service: %v\r\n", err)
				os.Exit(1)
			}
			s := lib.WindowsService{
				Service: AgentService,
			}
			err = s.RotateLogFile()
			if err != nil {
				fmt.Printf("Error sending rotate log file signal to the service: %v\r\n", err)
				os.Exit(1)
			} else {
				fmt.Printf("Rotate log file signal sent!\r\n")
			}
		},
	}
)

func init() {
	AgentService = &lib.AgentService{}
	AgentService.AddFlags(AgentCmdInstall.Flags())
	AgentService.AddFlags(AgentCmdStart.Flags())
	AgentCmd.AddCommand(AgentCmdUninstall)
	AgentCmd.AddCommand(AgentCmdInstall)
	AgentCmd.AddCommand(AgentCmdStop)
	AgentCmd.AddCommand(AgentCmdStart)
	AgentCmd.AddCommand(AgentCmdStatus)
	AgentCmd.AddCommand(AgentCmdReload)
	AgentCmd.AddCommand(AgentCmdRotateLog)
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, AgentCmdUninstall.Use)
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, AgentCmdInstall.Use)
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, AgentCmdStop.Use)
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, AgentCmdStart.Use)
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, AgentCmdStatus.Use)
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, AgentCmdReload.Use)
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, AgentCmdRotateLog.Use)
}
