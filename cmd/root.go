package cmd

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Files-com/files-cli/lib"
	files "github.com/Files-com/files-sdk-go/v3"
	lib2 "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/spf13/cobra"
	cobracompletefig "github.com/withfig/autocomplete-tools/integrations/cobra"
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
	featureFlags           []string
	RootCmd                = &cobra.Command{
		Use: "files-cli [resource]",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			sdkConfig := cmd.Context().Value("config").(files.Config)
			for _, flag := range featureFlags {
				sdkConfig.FeatureFlag(flag) // panic unknown flag
				sdkConfig.FeatureFlags[flag] = true
			}
			if APIKey != "" {
				sdkConfig.APIKey = APIKey
			}

			sdkConfig.Environment = files.NewEnvironment(Environment)
			debugFlag := cmd.Flag("debug")
			if debugFlag.Changed {
				if debug == "files-cli_[command]_[timestamp].log" {
					debug = fmt.Sprintf("files-cli_%v_%v.log", cmd.CalledAs(), strings.Replace(time.Now().Format(time.DateTime), " ", "_", 1))
				}
				if strings.ToLower(debug) == "stdout" {
					sdkConfig.Logger = log.New(os.Stdout, "", log.LstdFlags)
				} else {
					logFile, err := os.Create(debug)
					if err != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "%v\n", err)
						os.Exit(1)
					}
					sdkConfig.Logger = log.New(logFile, "", log.LstdFlags)
				}
				sdkConfig.Debug = true
				sdkConfig.Logger.Printf("Command: %v", strings.Join(os.Args, " "))
			}

			profile := &lib.Profiles{}
			err := profile.Load(&sdkConfig, ProfileValue)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "%v\n", err)
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			profile.Overrides = lib.Overrides{In: cmd.InOrStdin(), Out: cmd.OutOrStdout()}
			cmd.SetContext(context.WithValue(cmd.Context(), "profile", profile))
			cmd.SetContext(context.WithValue(cmd.Context(), "config", sdkConfig))
			if sdkConfig.Debug {
				sdkConfig.Logger.Printf("Environment: %v", sdkConfig.Environment.String())
				if transport, ok := sdkConfig.Client.HTTPClient.Transport.(*lib2.Transport); ok {
					sdkConfig.Logger.Printf("MaxIdleConnsPerHost: %v", transport.MaxIdleConnsPerHost)
					sdkConfig.Logger.Printf("MaxIdleConns: %v", transport.MaxIdleConns)
					sdkConfig.Logger.Printf("MaxConnsPerHost: %v", transport.MaxConnsPerHost)
				}
			}

			if OutputPath != "" {
				output, err := os.Create(OutputPath)
				if err != nil {
					return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
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
				Profile(cmd).CheckVersion(Version, lib.FetchLatestVersionNumber(sdkConfig, cmd.Context()), lib.InstalledViaBrew(), cmd.ErrOrStderr())
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
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}

			if Profile(cmd).Config.GetAPIKey() == "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "No API Key found. Using session login.\n")
				err = lib.CreateSession(cmd.Context(), files.SessionCreateParams{}, Profile(cmd))
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
)

func Init(version string, _commit string, _date string, config files.Config) {
	commit = _commit
	date = _date
	Version = version
	RootCmd.Version = strings.TrimSuffix(Version, "\n")
	config.UserAgent = "Files.com CLI" + " " + strings.TrimSpace(Version)
	RootCmd.ExecuteContext(context.WithValue(context.Background(), "config", config))
}

func init() {
	RootCmd.PersistentFlags().StringVar(&debug, "debug", "", "Enable verbose logging. Use --debug=[LOG FILE NAME] to specify a log file or --debug=STDOUT to display logs directly on the screen.")
	RootCmd.PersistentFlags().Lookup("debug").NoOptDefVal = "files-cli_[command]_[timestamp].log"
	RootCmd.PersistentFlags().BoolVar(&ignoreVersionCheck, "ignore-version-check", false, "Do not check for a new version of the CLI")
	RootCmd.PersistentFlags().StringVar(&ProfileValue, "profile", ProfileValue, "Setup a connection profile")
	RootCmd.PersistentFlags().StringVar(&Environment, "environment", Environment, "Set connection to an environment or site")
	RootCmd.PersistentFlags().Lookup("environment").Hidden = true
	RootCmd.PersistentFlags().StringVar(&APIKey, "api-key", "", "API Key")
	RootCmd.PersistentFlags().StringVarP(&OutputPath, "output", "o", "", "File path to save output")
	RootCmd.PersistentFlags().BoolVar(&Reauthentication, "reauthentication", Reauthentication, "For enhanced security during specific types of requests, we mandate reauthentication when using a session ID for authentication. In such cases, please supply the session user's password again using the --reauthentication flag.")
	RootCmd.PersistentFlags().StringSliceVar(&featureFlags, "feature-flag", featureFlags, "Enable feature flags")
	RootCmd.SuggestionsMinimumDistance = 1
	RootCmd.AddCommand(cobracompletefig.CreateCompletionSpecCommand())
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, "completion")
}

func Profile(cmd *cobra.Command) *lib.Profiles {
	return cmd.Context().Value("profile").(*lib.Profiles)
}
