package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	public_key "github.com/Files-com/files-sdk-go/v2/publickey"
)

var (
	PublicKeys = &cobra.Command{}
)

func PublicKeysInit() {
	PublicKeys = &cobra.Command{
		Use:  "public-keys [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command public-keys\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsPublicKeyList := files_sdk.PublicKeyListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List Public Keys",
		Long:  `List Public Keys`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsPublicKeyList
			params.MaxPages = MaxPagesList

			client := public_key.Client{Config: *config}
			it, err := client.List(ctx, params)
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, formatList, config.Logger())
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(it, formatList, fieldsList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}

	cmdList.Flags().Int64Var(&paramsPublicKeyList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsPublicKeyList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsPublicKeyList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	PublicKeys.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	paramsPublicKeyFind := files_sdk.PublicKeyFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Public Key`,
		Long:  `Show Public Key`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := public_key.Client{Config: *config}

			var publicKey interface{}
			var err error
			publicKey, err = client.Find(ctx, paramsPublicKeyFind)
			lib.HandleResponse(ctx, publicKey, err, formatFind, fieldsFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdFind.Flags().Int64Var(&paramsPublicKeyFind.Id, "id", 0, "Public Key ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright")
	PublicKeys.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	paramsPublicKeyCreate := files_sdk.PublicKeyCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Public Key`,
		Long:  `Create Public Key`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := public_key.Client{Config: *config}

			var publicKey interface{}
			var err error
			publicKey, err = client.Create(ctx, paramsPublicKeyCreate)
			lib.HandleResponse(ctx, publicKey, err, formatCreate, fieldsCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdCreate.Flags().Int64Var(&paramsPublicKeyCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVar(&paramsPublicKeyCreate.Title, "title", "", "Internal reference for key.")
	cmdCreate.Flags().StringVar(&paramsPublicKeyCreate.PublicKey, "public-key", "", "Actual contents of SSH key.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	PublicKeys.AddCommand(cmdCreate)
	var fieldsUpdate string
	var formatUpdate string
	paramsPublicKeyUpdate := files_sdk.PublicKeyUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Public Key`,
		Long:  `Update Public Key`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := public_key.Client{Config: *config}

			var publicKey interface{}
			var err error
			publicKey, err = client.Update(ctx, paramsPublicKeyUpdate)
			lib.HandleResponse(ctx, publicKey, err, formatUpdate, fieldsUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsPublicKeyUpdate.Id, "id", 0, "Public Key ID.")
	cmdUpdate.Flags().StringVar(&paramsPublicKeyUpdate.Title, "title", "", "Internal reference for key.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	PublicKeys.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	paramsPublicKeyDelete := files_sdk.PublicKeyDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Public Key`,
		Long:  `Delete Public Key`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := public_key.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsPublicKeyDelete)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}
	cmdDelete.Flags().Int64Var(&paramsPublicKeyDelete.Id, "id", 0, "Public Key ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright")
	PublicKeys.AddCommand(cmdDelete)
}
