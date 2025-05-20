package cmd

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files "github.com/Files-com/files-sdk-go/v3"
	lib2 "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/spf13/cobra"
	cobracompletefig "github.com/withfig/autocomplete-tools/integrations/cobra"
)

const (
	flagNameApiKey             = "api-key"
	flagNameDebug              = "debug"
	flagNameEnvironment        = "environment"
	flagNameFeatreFlag         = "feature-flag"
	flagNameFormat             = "format"
	flagNameIgnoreVersionCheck = "ignore-version-check"
	flagNameInteractive        = "interactive"
	flagNameNonInteractive     = "non-interactive"
	flagNameOutputPath         = "output"
	flagNameOutputPathShort    = "o"
	flagNameProfile            = "profile"
	flagNameReauthentication   = "reauthentication"
	flagNameUsePager           = "use-pager"
)

const (
	flagValueStdout = "stdout"
)

const (
	commandNameLogin       = "login"
	commandNameLogout      = "logout"
	commandNameVersion     = "version"
	commandNameAgent       = "agent"
	commandNameConfig      = "config"
	commandNameConfigSet   = "config-set"
	commandNameConfigReset = "config-reset"
	commandNameConfigShow  = "config-show"
	commandNameCompletion  = "completion"
)

const (
	contextKeyProfile = "profile"
	contextKeyConfig  = "config"
)

const (
	userAgentPattern = "Files.com CLI %s"
)

var (
	errNonInteractiveRequiresInput = clierr.Errorf(clierr.ErrorCodeUsage, "--%s provided without valid profile or API key", flagNameNonInteractive)
)

var (
	noAuthCmds = []string{commandNameConfigSet, commandNameConfigReset, commandNameConfigShow, commandNameVersion, commandNameAgent}
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
	nonInteractive         bool
	OutputPath             string
	Reauthentication       bool
	featureFlags           []string
	RootCmd                = &cobra.Command{
		Use: "files-cli [resource]",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// configure non-interactive flag combinations before anything else
			// because non-interactive influences the format of the output if
			// there are any errors.
			if err := configureNonInteractiveMode(cmd); err != nil {
				return err
			}

			sdkConfig := cmd.Context().Value(contextKeyConfig).(files.Config)
			for _, flag := range featureFlags {
				sdkConfig.FeatureFlag(flag) // panic unknown flag
				sdkConfig.FeatureFlags[flag] = true
			}
			if APIKey != "" {
				sdkConfig.APIKey = APIKey
			}

			sdkConfig.Environment = files.NewEnvironment(Environment)
			debugFlag := cmd.Flag(flagNameDebug)
			if debugFlag.Changed {
				if debug == "files-cli_[command]_[timestamp].log" {
					debug = fmt.Sprintf(
						"files-cli_%v_%v.log",
						cmd.CalledAs(),
						time.Now().Format("20060102_150405"),
					)
				}
				if strings.ToLower(debug) == flagValueStdout {
					sdkConfig.Logger = log.New(os.Stdout, "", log.LstdFlags)
				} else {
					logFile, err := os.Create(debug)
					if err != nil {
						return clierr.New(clierr.ErrorCodeFatal, err)
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
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			profile.Overrides = lib.Overrides{In: cmd.InOrStdin(), Out: cmd.OutOrStdout()}
			cmd.SetContext(context.WithValue(cmd.Context(), contextKeyProfile, profile))
			cmd.SetContext(context.WithValue(cmd.Context(), contextKeyConfig, sdkConfig))
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
					return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
				}
				cmd.SetOut(output)
			}

			// the login and logout commands don't need any further validation.
			if slices.Contains([]string{commandNameLogin, commandNameLogout}, cmd.Use) {
				return nil
			}

			// if the command is in the list of commands that ignore credentials, there's no need to do further validation.
			if slices.Contains(IgnoreCredentialsCheck, cmd.Use) || slices.Contains(IgnoreCredentialsCheck, cmd.Parent().Use) {
				return nil
			}

			// if the command has an alias that is in the list of commands that don't require credentials, there's no need to do further validation.
			if len(cmd.Aliases) != 0 && slices.Contains(noAuthCmds, cmd.Aliases[0]) {
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
					if err = lib.Reauthenicate(Profile(cmd)); err != nil {
						return clierr.New(clierr.ErrorCodeFatal, err)
					}
				}
				return nil
			}

			// check the non-interactive flag alone, because at this point the only
			// option is to prompt for input.
			if nonInteractive {
				return errNonInteractiveRequiresInput
			}

			if Profile(cmd).SessionExpired() {
				fmt.Fprintf(cmd.ErrOrStderr(), "The session has expired, you must log in again.\n")
				err = lib.CreateSession(cmd.Context(), files.SessionCreateParams{}, Profile(cmd))
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}

			if Profile(cmd).Config.GetAPIKey() == "" {
				fmt.Fprintf(cmd.ErrOrStderr(), "No API Key found.\nTo use an API key, provide it with the --api-key=FILES_API_KEY flag.\nFalling back to session login...\n")
				err = lib.CreateSession(cmd.Context(), files.SessionCreateParams{}, Profile(cmd))
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
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
	config.UserAgent = fmt.Sprintf(userAgentPattern, strings.TrimSpace(Version))
	if err := RootCmd.ExecuteContext(context.WithValue(context.Background(), contextKeyConfig, config)); err != nil {
		checkErr(err)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&debug, flagNameDebug, "", "Enable verbose logging. Use --debug=[LOG FILE NAME] to specify a log file or --debug=STDOUT to display logs directly on the screen.")
	RootCmd.PersistentFlags().Lookup(flagNameDebug).NoOptDefVal = "files-cli_[command]_[timestamp].log"
	RootCmd.PersistentFlags().BoolVar(&ignoreVersionCheck, flagNameIgnoreVersionCheck, false, "Do not check for a new version of the CLI")
	RootCmd.PersistentFlags().BoolVar(&nonInteractive, flagNameNonInteractive, false, "Do not prompt for user input")
	RootCmd.PersistentFlags().StringVar(&ProfileValue, flagNameProfile, ProfileValue, "Setup a connection profile")
	RootCmd.PersistentFlags().StringVar(&Environment, flagNameEnvironment, Environment, "Set connection to an environment or site")
	RootCmd.PersistentFlags().Lookup(flagNameEnvironment).Hidden = true
	RootCmd.PersistentFlags().StringVar(&APIKey, flagNameApiKey, "", "Set API Key for single use")
	RootCmd.PersistentFlags().StringVarP(&OutputPath, flagNameOutputPath, flagNameOutputPathShort, "", "File path to save output")
	RootCmd.PersistentFlags().BoolVar(&Reauthentication, flagNameReauthentication, Reauthentication, "For enhanced security during specific types of requests, we mandate reauthentication when using a session ID for authentication. In such cases, please supply the session user's password again using the --reauthentication flag.")
	RootCmd.PersistentFlags().StringSliceVar(&featureFlags, flagNameFeatreFlag, featureFlags, "Enable feature flags")
	RootCmd.SuggestionsMinimumDistance = 1
	RootCmd.AddCommand(cobracompletefig.CreateCompletionSpecCommand())
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, commandNameCompletion)
}

func Profile(cmd *cobra.Command) *lib.Profiles {
	if profile, ok := cmd.Context().Value(contextKeyProfile).(*lib.Profiles); ok {
		return profile
	}
	return &lib.Profiles{}
}

// checkErr exits with the appropriate status code if available.
func checkErr(err error) {
	if err == nil {
		return
	}
	statusError := clierr.From(err)
	os.Exit(int(statusError.Code))
}
