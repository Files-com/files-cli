package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/session"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Sessions())
}

func Sessions() *cobra.Command {
	Sessions := &cobra.Command{
		Use:  "sessions [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command sessions\n\t%v", args[0])
		},
	}
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsSessionCreate := files_sdk.SessionCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create user session (log in)`,
		Long:  `Create user session (log in)`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := session.Client{Config: config}

			var session interface{}
			var err error
			session, err = client.Create(paramsSessionCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), session, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsSessionCreate.Username, "username", "", "Username to sign in as")
	cmdCreate.Flags().StringVar(&paramsSessionCreate.Password, "password", "", "Password for sign in")
	cmdCreate.Flags().StringVar(&paramsSessionCreate.Otp, "otp", "", "If this user has a 2FA device, provide its OTP or code here.")
	cmdCreate.Flags().StringVar(&paramsSessionCreate.PartialSessionId, "partial-session-id", "", "Identifier for a partially-completed login")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Sessions.AddCommand(cmdCreate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete user session (log out)`,
		Long:  `Delete user session (log out)`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := session.Client{Config: config}

			var err error
			err = client.Delete(files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Sessions.AddCommand(cmdDelete)
	return Sessions
}
