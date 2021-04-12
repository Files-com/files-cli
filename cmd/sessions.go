package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

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
		Run:  func(cmd *cobra.Command, args []string) {},
	}
	var fieldsCreate string
	paramsSessionCreate := files_sdk.SessionCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := session.Client{Config: *ctx.GetConfig()}

			result, err := client.Create(paramsSessionCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdCreate.Flags().StringVarP(&paramsSessionCreate.Username, "username", "u", "", "Username to sign in as")
	cmdCreate.Flags().StringVarP(&paramsSessionCreate.Password, "password", "a", "", "Password for sign in")
	cmdCreate.Flags().StringVarP(&paramsSessionCreate.Otp, "otp", "o", "", "If this user has a 2FA device, provide its OTP or code here.")
	cmdCreate.Flags().StringVarP(&paramsSessionCreate.PartialSessionId, "partial-session-id", "p", "", "Identifier for a partially-completed login")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	Sessions.AddCommand(cmdCreate)
	var fieldsDelete string
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := session.Client{Config: *ctx.GetConfig()}

			result, err := client.Delete()
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Sessions.AddCommand(cmdDelete)
}
