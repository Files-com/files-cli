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
	usePagerList := true
	paramsApiKeyList := files_sdk.ApiKeyListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List Api Keys",
		Long:  `List Api Keys`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsApiKeyList
			params.MaxPages = MaxPagesList

			client := api_key.Client{Config: *config}
			it, err := client.List(ctx, params)
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger())
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
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
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
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	ApiKeys.AddCommand(cmdList)
	var fieldsFindCurrent string
	var formatFindCurrent string
	usePagerFindCurrent := true
	cmdFindCurrent := &cobra.Command{
		Use:   "find-current",
		Short: `Show information about current API key.  (Requires current API connection to be using an API key.)`,
		Long:  `Show information about current API key.  (Requires current API connection to be using an API key.)`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := api_key.Client{Config: *config}

			var apiKey interface{}
			var err error
			apiKey, err = client.FindCurrent(ctx)
			lib.HandleResponse(ctx, apiKey, err, formatFindCurrent, fieldsFindCurrent, usePagerFindCurrent, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}

	cmdFindCurrent.Flags().StringVarP(&fieldsFindCurrent, "fields", "", "", "comma separated list of field names")
	cmdFindCurrent.Flags().StringVarP(&formatFindCurrent, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdFindCurrent.Flags().BoolVar(&usePagerFindCurrent, "use-pager", usePagerFindCurrent, "Use $PAGER (.ie less, more, etc)")

	ApiKeys.AddCommand(cmdFindCurrent)
	var fieldsFind string
	var formatFind string
	usePagerFind := true
	paramsApiKeyFind := files_sdk.ApiKeyFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Api Key`,
		Long:  `Show Api Key`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := api_key.Client{Config: *config}

			var apiKey interface{}
			var err error
			apiKey, err = client.Find(ctx, paramsApiKeyFind)
			lib.HandleResponse(ctx, apiKey, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdFind.Flags().Int64Var(&paramsApiKeyFind.Id, "id", 0, "Api Key ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	ApiKeys.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	usePagerCreate := true
	paramsApiKeyCreate := files_sdk.ApiKeyCreateParams{}
	ApiKeyCreatePermissionSet := ""

	cmdCreate := &cobra.Command{
		Use:   "create [path]",
		Short: `Create Api Key`,
		Long:  `Create Api Key`,
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
			lib.HandleResponse(ctx, apiKey, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdCreate.Flags().Int64Var(&paramsApiKeyCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVar(&paramsApiKeyCreate.Name, "name", "", "Internal name for the API Key.  For your use.")
	lib.TimeVar(cmdCreate.Flags(), paramsApiKeyCreate.ExpiresAt, "expires-at")
	cmdCreate.Flags().StringVar(&ApiKeyCreatePermissionSet, "permission-set", "", fmt.Sprintf("Permissions for this API Key.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations).  Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know. %v", reflect.ValueOf(paramsApiKeyCreate.PermissionSet.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&paramsApiKeyCreate.Path, "path", "", "Folder path restriction for this api key.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	ApiKeys.AddCommand(cmdCreate)
	var fieldsUpdateCurrent string
	var formatUpdateCurrent string
	usePagerUpdateCurrent := true
	paramsApiKeyUpdateCurrent := files_sdk.ApiKeyUpdateCurrentParams{}
	ApiKeyUpdateCurrentPermissionSet := ""

	cmdUpdateCurrent := &cobra.Command{
		Use:   "update-current",
		Short: `Update current API key.  (Requires current API connection to be using an API key.)`,
		Long:  `Update current API key.  (Requires current API connection to be using an API key.)`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := api_key.Client{Config: *config}

			var apiKey interface{}
			var err error
			paramsApiKeyUpdateCurrent.PermissionSet = paramsApiKeyUpdateCurrent.PermissionSet.Enum()[ApiKeyUpdateCurrentPermissionSet]
			apiKey, err = client.UpdateCurrent(ctx, paramsApiKeyUpdateCurrent)
			lib.HandleResponse(ctx, apiKey, err, formatUpdateCurrent, fieldsUpdateCurrent, usePagerUpdateCurrent, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	lib.TimeVar(cmdUpdateCurrent.Flags(), paramsApiKeyUpdateCurrent.ExpiresAt, "expires-at")
	cmdUpdateCurrent.Flags().StringVar(&paramsApiKeyUpdateCurrent.Name, "name", "", "Internal name for the API Key.  For your use.")
	cmdUpdateCurrent.Flags().StringVar(&ApiKeyUpdateCurrentPermissionSet, "permission-set", "", fmt.Sprintf("Permissions for this API Key.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations).  Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know. %v", reflect.ValueOf(paramsApiKeyUpdateCurrent.PermissionSet.Enum()).MapKeys()))

	cmdUpdateCurrent.Flags().StringVarP(&fieldsUpdateCurrent, "fields", "", "", "comma separated list of field names")
	cmdUpdateCurrent.Flags().StringVarP(&formatUpdateCurrent, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdUpdateCurrent.Flags().BoolVar(&usePagerUpdateCurrent, "use-pager", usePagerUpdateCurrent, "Use $PAGER (.ie less, more, etc)")

	ApiKeys.AddCommand(cmdUpdateCurrent)
	var fieldsUpdate string
	var formatUpdate string
	usePagerUpdate := true
	paramsApiKeyUpdate := files_sdk.ApiKeyUpdateParams{}
	ApiKeyUpdatePermissionSet := ""

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Api Key`,
		Long:  `Update Api Key`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := api_key.Client{Config: *config}

			var apiKey interface{}
			var err error
			paramsApiKeyUpdate.PermissionSet = paramsApiKeyUpdate.PermissionSet.Enum()[ApiKeyUpdatePermissionSet]
			apiKey, err = client.Update(ctx, paramsApiKeyUpdate)
			lib.HandleResponse(ctx, apiKey, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsApiKeyUpdate.Id, "id", 0, "Api Key ID.")
	cmdUpdate.Flags().StringVar(&paramsApiKeyUpdate.Name, "name", "", "Internal name for the API Key.  For your use.")
	lib.TimeVar(cmdUpdate.Flags(), paramsApiKeyUpdate.ExpiresAt, "expires-at")
	cmdUpdate.Flags().StringVar(&ApiKeyUpdatePermissionSet, "permission-set", "", fmt.Sprintf("Permissions for this API Key.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations).  Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know. %v", reflect.ValueOf(paramsApiKeyUpdate.PermissionSet.Enum()).MapKeys()))

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	ApiKeys.AddCommand(cmdUpdate)
	var fieldsDeleteCurrent string
	var formatDeleteCurrent string
	usePagerDeleteCurrent := true
	cmdDeleteCurrent := &cobra.Command{
		Use:   "delete-current",
		Short: `Delete current API key.  (Requires current API connection to be using an API key.)`,
		Long:  `Delete current API key.  (Requires current API connection to be using an API key.)`,
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
	cmdDeleteCurrent.Flags().StringVarP(&formatDeleteCurrent, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdDeleteCurrent.Flags().BoolVar(&usePagerDeleteCurrent, "use-pager", usePagerDeleteCurrent, "Use $PAGER (.ie less, more, etc)")

	ApiKeys.AddCommand(cmdDeleteCurrent)
	var fieldsDelete string
	var formatDelete string
	usePagerDelete := true
	paramsApiKeyDelete := files_sdk.ApiKeyDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Api Key`,
		Long:  `Delete Api Key`,
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
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	ApiKeys.AddCommand(cmdDelete)
}
