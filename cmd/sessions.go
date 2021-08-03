package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/session"
)

var (
	Sessions = &cobra.Command{}
)

func SessionsInit() {
	Sessions = &cobra.Command{
		Use:  "sessions [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command sessions\n\t%v", args[0])
		},
	}
	var fieldsCreate string
	var formatCreate string
	paramsSessionCreate := files_sdk.SessionCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := session.Client{Config: *config}

			result, err := client.Create(ctx, paramsSessionCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatCreate, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().StringVarP(&paramsSessionCreate.Username, "username", "u", "", "Username to sign in as")
	cmdCreate.Flags().StringVarP(&paramsSessionCreate.Password, "password", "a", "", "Password for sign in")
	cmdCreate.Flags().StringVarP(&paramsSessionCreate.Otp, "otp", "o", "", "If this user has a 2FA device, provide its OTP or code here.")
	cmdCreate.Flags().StringVarP(&paramsSessionCreate.PartialSessionId, "partial-session-id", "p", "", "Identifier for a partially-completed login")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-light")
	Sessions.AddCommand(cmdCreate)
	var fieldsDelete string
	var formatDelete string
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := session.Client{Config: *config}

			result, err := client.Delete(ctx)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatDelete, fieldsDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-light")
	Sessions.AddCommand(cmdDelete)
}
