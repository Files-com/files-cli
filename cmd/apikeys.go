package cmd

import (
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	api_key "github.com/Files-com/files-sdk-go/v2/apikey"
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
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(it, formatList, fieldsList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}

	cmdList.Flags().Int64Var(&paramsApiKeyList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsApiKeyList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsApiKeyList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	ApiKeys.AddCommand(cmdList)
	var fieldsFindCurrent string
	var formatFindCurrent string
	cmdFindCurrent := &cobra.Command{
		Use: "find-current",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := api_key.Client{Config: *config}

			var apiKey interface{}
			var err error
			apiKey, err = client.FindCurrent(ctx)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(apiKey, formatFindCurrent, fieldsFindCurrent, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}

	cmdFindCurrent.Flags().StringVarP(&fieldsFindCurrent, "fields", "", "", "comma separated list of field names")
	cmdFindCurrent.Flags().StringVarP(&formatFindCurrent, "format", "", "table", "json, csv, table, table-dark, table-bright")
	ApiKeys.AddCommand(cmdFindCurrent)
	var fieldsFind string
	var formatFind string
	paramsApiKeyFind := files_sdk.ApiKeyFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := api_key.Client{Config: *config}

			var apiKey interface{}
			var err error
			apiKey, err = client.Find(ctx, paramsApiKeyFind)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(apiKey, formatFind, fieldsFind, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdFind.Flags().Int64Var(&paramsApiKeyFind.Id, "id", 0, "Api Key ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright")
	ApiKeys.AddCommand(cmdFind)
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
			var apiKey interface{}
			var err error
			paramsApiKeyCreate.PermissionSet = paramsApiKeyCreate.PermissionSet.Enum()[ApiKeyCreatePermissionSet]
			apiKey, err = client.Create(ctx, paramsApiKeyCreate)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(apiKey, formatCreate, fieldsCreate, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdCreate.Flags().Int64Var(&paramsApiKeyCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVar(&paramsApiKeyCreate.Name, "name", "", "Internal name for the API Key.  For your use.")
	lib.TimeVar(cmdCreate.Flags(), paramsApiKeyCreate.ExpiresAt, "expires-at")
	cmdCreate.Flags().StringVar(&ApiKeyCreatePermissionSet, "permission-set", "", fmt.Sprintf("Permissions for this API Key.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations).  Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know. %v", reflect.ValueOf(paramsApiKeyCreate.PermissionSet.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&paramsApiKeyCreate.Path, "path", "", "Folder path restriction for this api key.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	ApiKeys.AddCommand(cmdCreate)
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

			var apiKey interface{}
			var err error
			paramsApiKeyUpdateCurrent.PermissionSet = paramsApiKeyUpdateCurrent.PermissionSet.Enum()[ApiKeyUpdateCurrentPermissionSet]
			apiKey, err = client.UpdateCurrent(ctx, paramsApiKeyUpdateCurrent)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(apiKey, formatUpdateCurrent, fieldsUpdateCurrent, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	lib.TimeVar(cmdUpdateCurrent.Flags(), paramsApiKeyUpdateCurrent.ExpiresAt, "expires-at")
	cmdUpdateCurrent.Flags().StringVar(&paramsApiKeyUpdateCurrent.Name, "name", "", "Internal name for the API Key.  For your use.")
	cmdUpdateCurrent.Flags().StringVar(&ApiKeyUpdateCurrentPermissionSet, "permission-set", "", fmt.Sprintf("Permissions for this API Key.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations).  Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know. %v", reflect.ValueOf(paramsApiKeyUpdateCurrent.PermissionSet.Enum()).MapKeys()))

	cmdUpdateCurrent.Flags().StringVarP(&fieldsUpdateCurrent, "fields", "", "", "comma separated list of field names")
	cmdUpdateCurrent.Flags().StringVarP(&formatUpdateCurrent, "format", "", "table", "json, csv, table, table-dark, table-bright")
	ApiKeys.AddCommand(cmdUpdateCurrent)
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

			var apiKey interface{}
			var err error
			paramsApiKeyUpdate.PermissionSet = paramsApiKeyUpdate.PermissionSet.Enum()[ApiKeyUpdatePermissionSet]
			apiKey, err = client.Update(ctx, paramsApiKeyUpdate)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(apiKey, formatUpdate, fieldsUpdate, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsApiKeyUpdate.Id, "id", 0, "Api Key ID.")
	cmdUpdate.Flags().StringVar(&paramsApiKeyUpdate.Name, "name", "", "Internal name for the API Key.  For your use.")
	lib.TimeVar(cmdUpdate.Flags(), paramsApiKeyUpdate.ExpiresAt, "expires-at")
	cmdUpdate.Flags().StringVar(&ApiKeyUpdatePermissionSet, "permission-set", "", fmt.Sprintf("Permissions for this API Key.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations).  Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know. %v", reflect.ValueOf(paramsApiKeyUpdate.PermissionSet.Enum()).MapKeys()))

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	ApiKeys.AddCommand(cmdUpdate)
	var fieldsDeleteCurrent string
	var formatDeleteCurrent string
	cmdDeleteCurrent := &cobra.Command{
		Use: "delete-current",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := api_key.Client{Config: *config}

			var err error
			err = client.DeleteCurrent(ctx)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}

	cmdDeleteCurrent.Flags().StringVarP(&fieldsDeleteCurrent, "fields", "", "", "comma separated list of field names")
	cmdDeleteCurrent.Flags().StringVarP(&formatDeleteCurrent, "format", "", "table", "json, csv, table, table-dark, table-bright")
	ApiKeys.AddCommand(cmdDeleteCurrent)
	var fieldsDelete string
	var formatDelete string
	paramsApiKeyDelete := files_sdk.ApiKeyDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := api_key.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsApiKeyDelete)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}
	cmdDelete.Flags().Int64Var(&paramsApiKeyDelete.Id, "id", 0, "Api Key ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright")
	ApiKeys.AddCommand(cmdDelete)
}
