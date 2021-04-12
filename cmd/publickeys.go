package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	public_key "github.com/Files-com/files-sdk-go/publickey"
)

var (
	PublicKeys = &cobra.Command{}
)

func PublicKeysInit() {
	PublicKeys = &cobra.Command{
		Use:  "public-keys [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
	var fieldsList string
	paramsPublicKeyList := files_sdk.PublicKeyListParams{}
	var MaxPagesList int64
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			params := paramsPublicKeyList
			params.MaxPages = MaxPagesList
			client := public_key.Client{Config: *ctx.GetConfig()}
			it, err := client.List(params)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
			err = lib.JsonMarshalIter(it, fieldsList)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdList.Flags().Int64VarP(&paramsPublicKeyList.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVarP(&paramsPublicKeyList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsPublicKeyList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	PublicKeys.AddCommand(cmdList)
	var fieldsFind string
	paramsPublicKeyFind := files_sdk.PublicKeyFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := public_key.Client{Config: *ctx.GetConfig()}

			result, err := client.Find(paramsPublicKeyFind)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsFind)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsPublicKeyFind.Id, "id", "i", 0, "Public Key ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	PublicKeys.AddCommand(cmdFind)
	var fieldsCreate string
	paramsPublicKeyCreate := files_sdk.PublicKeyCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := public_key.Client{Config: *ctx.GetConfig()}

			result, err := client.Create(paramsPublicKeyCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdCreate.Flags().Int64VarP(&paramsPublicKeyCreate.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVarP(&paramsPublicKeyCreate.Title, "title", "t", "", "Internal reference for key.")
	cmdCreate.Flags().StringVarP(&paramsPublicKeyCreate.PublicKey, "public-key", "p", "", "Actual contents of SSH key.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	PublicKeys.AddCommand(cmdCreate)
	var fieldsUpdate string
	paramsPublicKeyUpdate := files_sdk.PublicKeyUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := public_key.Client{Config: *ctx.GetConfig()}

			result, err := client.Update(paramsPublicKeyUpdate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsUpdate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsPublicKeyUpdate.Id, "id", "i", 0, "Public Key ID.")
	cmdUpdate.Flags().StringVarP(&paramsPublicKeyUpdate.Title, "title", "t", "", "Internal reference for key.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	PublicKeys.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsPublicKeyDelete := files_sdk.PublicKeyDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := public_key.Client{Config: *ctx.GetConfig()}

			result, err := client.Delete(paramsPublicKeyDelete)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsPublicKeyDelete.Id, "id", "i", 0, "Public Key ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	PublicKeys.AddCommand(cmdDelete)
}
