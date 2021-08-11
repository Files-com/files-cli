package cmd

import (
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"

	api_key "github.com/Files-com/files-sdk-go/apikey"
)

var (
	ApiKeys = &cobra.Command{}
)

func ApiKeysInit() {
	ApiKeys = &cobra.Command{
		Use:  "api-keys [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command api-keys\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsApiKeyList := files_sdk.ApiKeyListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsApiKeyList
			params.MaxPages = MaxPagesList

			client := api_key.Client{Config: *config}
			it, err := client.List(ctx, params)
			if err != nil {
				lib.ClientError(ctx, err)
			}
			err = lib.FormatIter(it, formatList, fieldsList)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdList.Flags().Int64VarP(&paramsApiKeyList.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVarP(&paramsApiKeyList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsApiKeyList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-light")
	ApiKeys.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	paramsApiKeyFind := files_sdk.ApiKeyFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := api_key.Client{Config: *config}

			result, err := client.Find(ctx, paramsApiKeyFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatFind, fieldsFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsApiKeyFind.Id, "id", "i", 0, "Api Key ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-light")
	ApiKeys.AddCommand(cmdFind)
	var fieldsFindCurrent string
	var formatFindCurrent string
	cmdFindCurrent := &cobra.Command{
		Use: "find-current",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := api_key.Client{Config: *config}

			result, err := client.FindCurrent(ctx)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatFindCurrent, fieldsFindCurrent)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}

	cmdFindCurrent.Flags().StringVarP(&fieldsFindCurrent, "fields", "", "", "comma separated list of field names")
	cmdFindCurrent.Flags().StringVarP(&formatFindCurrent, "format", "", "table", "json, csv, table, table-dark, table-light")
	ApiKeys.AddCommand(cmdFindCurrent)
	var fieldsCreate string
	var formatCreate string
	paramsApiKeyCreate := files_sdk.ApiKeyCreateParams{}
	ApiKeyCreatePermissionSet := ""

	cmdCreate := &cobra.Command{
		Use: "create [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := api_key.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsApiKeyCreate.Path = args[0]
			}
			paramsApiKeyCreate.PermissionSet = paramsApiKeyCreate.PermissionSet.Enum()[ApiKeyCreatePermissionSet]

			result, err := client.Create(ctx, paramsApiKeyCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatCreate, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().Int64VarP(&paramsApiKeyCreate.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVarP(&paramsApiKeyCreate.Name, "name", "n", "", "Internal name for the API Key.  For your use.")
	lib.TimeVarP(cmdCreate.Flags(), &paramsApiKeyCreate.ExpiresAt, "expires-at", "e")
	cmdCreate.Flags().StringVarP(&ApiKeyCreatePermissionSet, "permission-set", "r", "", fmt.Sprintf("Permissions for this API Key.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations).  Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know. %v", reflect.ValueOf(paramsApiKeyCreate.PermissionSet.Enum()).MapKeys()))
	cmdCreate.Flags().StringVarP(&paramsApiKeyCreate.Path, "path", "p", "", "Folder path restriction for this api key.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-light")
	ApiKeys.AddCommand(cmdCreate)
	var fieldsUpdate string
	var formatUpdate string
	paramsApiKeyUpdate := files_sdk.ApiKeyUpdateParams{}
	ApiKeyUpdatePermissionSet := ""

	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := api_key.Client{Config: *config}

			paramsApiKeyUpdate.PermissionSet = paramsApiKeyUpdate.PermissionSet.Enum()[ApiKeyUpdatePermissionSet]

			result, err := client.Update(ctx, paramsApiKeyUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatUpdate, fieldsUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsApiKeyUpdate.Id, "id", "i", 0, "Api Key ID.")
	cmdUpdate.Flags().StringVarP(&paramsApiKeyUpdate.Name, "name", "n", "", "Internal name for the API Key.  For your use.")
	lib.TimeVarP(cmdUpdate.Flags(), &paramsApiKeyUpdate.ExpiresAt, "expires-at", "e")
	cmdUpdate.Flags().StringVarP(&ApiKeyUpdatePermissionSet, "permission-set", "p", "", fmt.Sprintf("Permissions for this API Key.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations).  Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know. %v", reflect.ValueOf(paramsApiKeyUpdate.PermissionSet.Enum()).MapKeys()))

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-light")
	ApiKeys.AddCommand(cmdUpdate)
	var fieldsUpdateCurrent string
	var formatUpdateCurrent string
	paramsApiKeyUpdateCurrent := files_sdk.ApiKeyUpdateCurrentParams{}
	ApiKeyUpdateCurrentPermissionSet := ""

	cmdUpdateCurrent := &cobra.Command{
		Use: "update-current",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := api_key.Client{Config: *config}

			paramsApiKeyUpdateCurrent.PermissionSet = paramsApiKeyUpdateCurrent.PermissionSet.Enum()[ApiKeyUpdateCurrentPermissionSet]

			result, err := client.UpdateCurrent(ctx, paramsApiKeyUpdateCurrent)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatUpdateCurrent, fieldsUpdateCurrent)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	lib.TimeVarP(cmdUpdateCurrent.Flags(), &paramsApiKeyUpdateCurrent.ExpiresAt, "expires-at", "e")
	cmdUpdateCurrent.Flags().StringVarP(&paramsApiKeyUpdateCurrent.Name, "name", "n", "", "Internal name for the API Key.  For your use.")
	cmdUpdateCurrent.Flags().StringVarP(&ApiKeyUpdateCurrentPermissionSet, "permission-set", "p", "", fmt.Sprintf("Permissions for this API Key.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations).  Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know. %v", reflect.ValueOf(paramsApiKeyUpdateCurrent.PermissionSet.Enum()).MapKeys()))

	cmdUpdateCurrent.Flags().StringVarP(&fieldsUpdateCurrent, "fields", "", "", "comma separated list of field names")
	cmdUpdateCurrent.Flags().StringVarP(&formatUpdateCurrent, "format", "", "table", "json, csv, table, table-dark, table-light")
	ApiKeys.AddCommand(cmdUpdateCurrent)
	var fieldsDelete string
	var formatDelete string
	paramsApiKeyDelete := files_sdk.ApiKeyDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := api_key.Client{Config: *config}

			result, err := client.Delete(ctx, paramsApiKeyDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatDelete, fieldsDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsApiKeyDelete.Id, "id", "i", 0, "Api Key ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-light")
	ApiKeys.AddCommand(cmdDelete)
	var fieldsDeleteCurrent string
	var formatDeleteCurrent string
	cmdDeleteCurrent := &cobra.Command{
		Use: "delete-current",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := api_key.Client{Config: *config}

			result, err := client.DeleteCurrent(ctx)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatDeleteCurrent, fieldsDeleteCurrent)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}

	cmdDeleteCurrent.Flags().StringVarP(&fieldsDeleteCurrent, "fields", "", "", "comma separated list of field names")
	cmdDeleteCurrent.Flags().StringVarP(&formatDeleteCurrent, "format", "", "table", "json, csv, table, table-dark, table-light")
	ApiKeys.AddCommand(cmdDeleteCurrent)
}
