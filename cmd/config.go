package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"
)

var (
	Config = &cobra.Command{
		Use:  "config",
		Args: cobra.ExactArgs(0),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func ConfigInit() {
	configParams := lib.Config{}
	subdomainCreate := &cobra.Command{
		Use:     "set",
		Aliases: []string{"config-set"},
		Run: func(cmd *cobra.Command, args []string) {
			configParams.Save()
		},
	}
	subdomainCreate.Flags().StringVarP(&configParams.Subdomain, "subdomain", "s", "", "Subdomain of site")
	subdomainCreate.Flags().StringVarP(&configParams.Username, "username", "u", "", "Username to sign in as")
	subdomainCreate.Flags().StringVarP(&configParams.APIKey, "api-key", "a", "", "API Key")

	Config.AddCommand(subdomainCreate)
	resetDelete := &cobra.Command{
		Use:     "reset",
		Aliases: []string{"config-reset"},
		Run: func(cmd *cobra.Command, args []string) {
			configParams.Reset()
		},
	}
	Config.AddCommand(resetDelete)
}
