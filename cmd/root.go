package cmd

import (
	"context"
	"log"
	"strings"

	"github.com/Files-com/files-cli/lib"
	files "github.com/Files-com/files-sdk-go/v2"
	"github.com/spf13/cobra"
	cobracompletefig "github.com/withfig/autocomplete-tools/integrations/cobra"

	_ "embed"
	"fmt"
	"os"
)

var (
	commit                 string
	date                   string
	IgnoreCredentialsCheck []string
	Version                string
	ProfileValue           string
	Environment            string
	APIKey                 string
	debug                  string
	ignoreVersionCheck     bool
	OutputPath             string
	Reauthentication       bool
	RootCmd                = &cobra.Command{
		Use: "files-cli [resource]",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			sdkConfig := cmd.Context().Value("config").(*files.Config)
			if APIKey != "" {
				sdkConfig.APIKey = APIKey
			}

			sdkConfig.Environment = files.NewEnvironment(Environment)
			debugFlag := cmd.Flag("debug")
			if debugFlag.Changed {
				logFile, err := os.Create(debug)
				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "%v\n", err)
					os.Exit(1)
				}
				sdkConfig.Debug = true
				sdkConfig.SetLogger(log.New(logFile, "", log.LstdFlags))
			}

			profile := &lib.Profiles{}
			err := profile.Load(sdkConfig, ProfileValue)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "%v\n", err)
				return lib.ClientError(cmd.Context(), Profile(cmd), err, cmd.ErrOrStderr())
			}
			profile.Overrides = lib.Overrides{In: cmd.InOrStdin(), Out: cmd.OutOrStdout()}
			cmd.SetContext(context.WithValue(cmd.Context(), "profile", profile))

			if OutputPath != "" {
				output, err := os.Create(OutputPath)
				if err != nil {
					return lib.ClientError(cmd.Context(), Profile(cmd), err, cmd.ErrOrStderr())
				}
				cmd.SetOut(output)
			}

			if lib.Includes(cmd.Use, []string{"login", "logout"}) {
				return nil
			}

			if lib.Includes(cmd.Use, IgnoreCredentialsCheck) || lib.Includes(cmd.Parent().Use, IgnoreCredentialsCheck) {
				return nil
			}

			if len(cmd.Aliases) != 0 && lib.Includes(cmd.Aliases[0], []string{"config-set", "config-reset", "config-show", "version", "agent"}) {
				return nil
			}

			if !ignoreVersionCheck {
				Profile(cmd).CheckVersion(Version, lib.FetchLatestVersionNumber(*sdkConfig, cmd.Context()), lib.InstalledViaBrew(), cmd.ErrOrStderr())
			}

			if Profile(cmd).Config.GetAPIKey() != "" {
				return nil
			}

			if Profile(cmd).ValidSession() {
				if Reauthentication {
					err = lib.Reauthenicate(Profile(cmd))
					if err != nil {
						return err
					}
				}
				return nil
			}
			if Profile(cmd).SessionExpired() {
				fmt.Fprintf(cmd.ErrOrStderr(), "The session has expired, you must log in again.\n")
				err = lib.CreateSession(cmd.Context(), files.SessionCreateParams{}, Profile(cmd))
				if err != nil {
					return lib.ClientError(cmd.Context(), Profile(cmd), err, cmd.ErrOrStderr())
				}
				return nil
			}

			if Profile(cmd).Config.GetAPIKey() == "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "No API Key found. Using session login.\n")
				err = lib.CreateSession(cmd.Context(), files.SessionCreateParams{}, Profile(cmd))
				if err != nil {
					return lib.ClientError(cmd.Context(), Profile(cmd), err, cmd.ErrOrStderr())
				}
			}
			return nil
		},
	}
)

func Init(version string, _commit string, _date string, config *files.Config) {
	commit = _commit
	date = _date
	Version = version
	RootCmd.Version = strings.TrimSuffix(Version, "\n")
	config.UserAgent = "Files.com CLI" + " " + strings.TrimSpace(Version)
	RootCmd.ExecuteContext(context.WithValue(context.Background(), "config", config))
}

func init() {
	RootCmd.PersistentFlags().StringVar(&debug, "debug", "", "verbose logging")
	RootCmd.PersistentFlags().BoolVar(&ignoreVersionCheck, "ignore-version-check", false, "API Key")
	RootCmd.PersistentFlags().StringVar(&ProfileValue, "profile", ProfileValue, "setup a connection profile")
	RootCmd.PersistentFlags().StringVar(&Environment, "environment", Environment, "Set connection to an environment or site")
	RootCmd.PersistentFlags().Lookup("environment").Hidden = true
	RootCmd.PersistentFlags().StringVar(&APIKey, "api-key", "", "API Key")
	RootCmd.PersistentFlags().StringVarP(&OutputPath, "output", "o", "", "file path to save output")
	RootCmd.PersistentFlags().BoolVar(&Reauthentication, "reauthentication", Reauthentication, "If authenticating to the API via a session ID (as opposed to an API key), we require that you provide the session user’s password again in a X-Files-Reauthentication header for certain types of requests where we want to add an additional level of security. We call this process Reauthentication.")
	RootCmd.SuggestionsMinimumDistance = 1
	RootCmd.AddCommand(cobracompletefig.CreateCompletionSpecCommand())
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, "completion")
}

func Profile(cmd *cobra.Command) *lib.Profiles {
	return cmd.Context().Value("profile").(*lib.Profiles)
}
