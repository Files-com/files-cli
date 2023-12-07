package cmd

import (
	"fmt"

	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func init() {
	RootCmd.AddCommand(Config())
}

func Config() *cobra.Command {
	Config := &cobra.Command{
		Use:  "config",
		Args: cobra.ExactArgs(0),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
	configParams := &lib.Profile{}
	configSet := &cobra.Command{
		Use:     "set",
		Aliases: []string{"config-set"},
		Run: func(cmd *cobra.Command, args []string) {
			if configParams.Subdomain != "" {
				Profile(cmd).Current().Subdomain = configParams.Subdomain
			}
			if configParams.Username != "" {
				Profile(cmd).Current().Username = configParams.Username
			}
			if configParams.APIKey != "" {
				Profile(cmd).Current().APIKey = configParams.APIKey
			}
			if configParams.Endpoint != "" {
				Profile(cmd).Current().Endpoint = configParams.Endpoint
			}

			if configParams.ConcurrentConnectionLimit != 0 {
				Profile(cmd).Current().ConcurrentConnectionLimit = configParams.ConcurrentConnectionLimit
			}

			if len(configParams.ResourceFormat) > 0 {
				Profile(cmd).Current().ResourceFormat = configParams.ResourceFormat
			}

			Profile(cmd).Save()
		},
	}
	configSet.Flags().StringVarP(&configParams.Subdomain, "subdomain", "d", configParams.Subdomain, "Subdomain of site")
	configSet.Flags().StringVarP(&configParams.Username, "username", "u", configParams.Username, "Username to sign in as")
	configSet.Flags().StringVarP(&configParams.APIKey, "api-key", "a", configParams.APIKey, "Store API Key in plain text config file `~/.config/files-cli`")
	configSet.Flags().StringVarP(&configParams.Endpoint, "endpoint", "e", configParams.Endpoint, "For use with custom domains. Example: 'https://myfilescustomdomain.com'")
	configSet.Flags().IntVarP(&configParams.ConcurrentConnectionLimit, "concurrent-connection-limit", "c", configParams.ConcurrentConnectionLimit, "Set the maximum number of concurrent connections.")
	configSet.Flags().StringSliceVarP(&configParams.ResourceFormat, "format", "f", configParams.ResourceFormat, fmt.Sprintf("Set default resource format: %v", lib.FormatHelpText))

	Config.AddCommand(configSet)
	resetConfig := lib.ResetConfig{}
	var resetDelete *cobra.Command
	resetDelete = &cobra.Command{
		Use:     "reset",
		Aliases: []string{"config-reset"},
		RunE: func(cmd *cobra.Command, args []string) error {
			anyFlagSet := false
			resetDelete.Flags().Visit(func(flag *pflag.Flag) {
				anyFlagSet = true
			})

			if anyFlagSet {
				return Profile(cmd).ResetWith(resetConfig)
			} else {
				return Profile(cmd).Reset()
			}
		},
	}
	resetDelete.Flags().BoolVarP(&resetConfig.Subdomain, "subdomain", "d", false, "Subdomain of site")
	resetDelete.Flags().BoolVarP(&resetConfig.Username, "username", "u", false, "Username to sign in as")
	resetDelete.Flags().BoolVarP(&resetConfig.APIKey, "api-key", "a", false, "API Key")
	resetDelete.Flags().BoolVarP(&resetConfig.Endpoint, "endpoint", "e", false, "For testing only, example: 'https://site.files.com'\nTo change subdomain use flag instead.")
	resetDelete.Flags().BoolVarP(&resetConfig.Session, "session", "s", false, "")
	resetDelete.Flags().BoolVarP(&resetConfig.VersionCheck, "version-check", "v", false, "")
	resetDelete.Flags().BoolVarP(&resetConfig.ConcurrentConnectionLimit, "concurrent-connection-limit", "c", false, "Reset the maximum number of concurrent connections.")
	resetDelete.Flags().BoolVarP(&resetConfig.ResourceFormat, "format", "f", false, "Reset default resource format")

	Config.AddCommand(resetDelete)

	configShow := &cobra.Command{
		Use:     "show",
		Aliases: []string{"config-show"},
		Run: func(cmd *cobra.Command, args []string) {
			fields := ""
			if len(args) > 0 {
				fields = args[0]
			}
			lib.JsonMarshal(Profile(cmd), []string{fields}, false, "")
		},
	}
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, "config")
	Config.AddCommand(configShow)
	return Config
}
