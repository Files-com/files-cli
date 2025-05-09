package cmd

import (
	"fmt"
	"os"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files "github.com/Files-com/files-sdk-go/v3"
	"github.com/spf13/cobra"
)

var (
	errNonInteractiveLogin = clierr.Errorf(clierr.ErrorCodeUsage, "login command does not support --%s", flagNameNonInteractive)
)

func init() {
	RootCmd.AddCommand(Login())
}

func Login() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "login",
		Args: cobra.NoArgs,
		PreRunE: func(command *cobra.Command, args []string) error {
			// the non-interactive flag is hidden, but is technically still
			// available. The logic here returns an error indicating that
			// non-interactive is not supported for the login action.
			command.Flags().MarkHidden(flagNameNonInteractive)
			if command.Flag(flagNameNonInteractive).Changed {
				return errNonInteractiveLogin
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			Profile(cmd).Overrides = lib.Overrides{Out: cmd.OutOrStdout(), In: cmd.InOrStdin()}.Init()
			err := lib.CreateSession(cmd.Context(), files.SessionCreateParams{}, Profile(cmd))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
		// the non-interactive flag also needs to be explicitly hidden
		// in the help function.
		command.Flags().MarkHidden(flagNameNonInteractive)
		command.Parent().HelpFunc()(command, strings)
	})
	return cmd
}
