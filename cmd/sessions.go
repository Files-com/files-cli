package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/session"
)

func init() {
	RootCmd.AddCommand(Sessions())
}

func Sessions() *cobra.Command {
	Sessions := &cobra.Command{
		Use:  "sessions [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command sessions\n\t%v", args[0])
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
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := session.Client{Config: *config}

			var session interface{}
			var err error
			session, err = client.Create(ctx, paramsSessionCreate)
			return lib.HandleResponse(ctx, Profile(cmd), session, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdCreate.Flags().StringVar(&paramsSessionCreate.Username, "username", "", "Username to sign in as")
	cmdCreate.Flags().StringVar(&paramsSessionCreate.Password, "password", "", "Password for sign in")
	cmdCreate.Flags().StringVar(&paramsSessionCreate.Otp, "otp", "", "If this user has a 2FA device, provide its OTP or code here.")
	cmdCreate.Flags().StringVar(&paramsSessionCreate.PartialSessionId, "partial-session-id", "", "Identifier for a partially-completed login")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
json-styles: {raw, pretty}`)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Sessions.AddCommand(cmdCreate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete user session (log out)`,
		Long:  `Delete user session (log out)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := session.Client{Config: *config}

			var err error
			err = client.Delete(ctx)
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
json-styles: {raw, pretty}`)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Sessions.AddCommand(cmdDelete)
	return Sessions
}
