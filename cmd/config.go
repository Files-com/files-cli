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
	configParams.Load()
	subdomainCreate := &cobra.Command{
		Use:     "set",
		Aliases: []string{"config-set"},
		Run: func(cmd *cobra.Command, args []string) {
			configParams.Save()
		},
	}
	subdomainCreate.Flags().StringVarP(&configParams.Subdomain, "subdomain", "d", configParams.Subdomain, "Subdomain of site")
	subdomainCreate.Flags().StringVarP(&configParams.Username, "username", "u", configParams.Username, "Username to sign in as")
	subdomainCreate.Flags().StringVarP(&configParams.APIKey, "api-key", "a", configParams.APIKey, "API Key")
	subdomainCreate.Flags().StringVarP(&configParams.Endpoint, "endpoint", "e", configParams.Endpoint, "For testing only, example: 'https://site.files.com'\nTo change subdomain use flag instead.")

	Config.AddCommand(subdomainCreate)
	resetConfig := lib.ResetConfig{}
	resetDelete := &cobra.Command{
		Use:     "reset",
		Aliases: []string{"config-reset"},
		Run: func(cmd *cobra.Command, args []string) {
			configParams.ResetWith(resetConfig)
		},
	}
	resetDelete.Flags().BoolVarP(&resetConfig.Subdomain, "subdomain", "d", false, "Subdomain of site")
	resetDelete.Flags().BoolVarP(&resetConfig.Username, "username", "u", false, "Username to sign in as")
	resetDelete.Flags().BoolVarP(&resetConfig.APIKey, "api-key", "a", false, "API Key")
	resetDelete.Flags().BoolVarP(&resetConfig.Endpoint, "endpoint", "e", false, "For testing only, example: 'https://site.files.com'\nTo change subdomain use flag instead.")
	resetDelete.Flags().BoolVarP(&resetConfig.Session, "session", "s", false, "")

	Config.AddCommand(resetDelete)

	resetShow := &cobra.Command{
		Use:     "show",
		Aliases: []string{"config-show"},
		Run: func(cmd *cobra.Command, args []string) {
			lib.JsonMarshal(configParams, "")
		},
	}

	Config.AddCommand(resetShow)
}
