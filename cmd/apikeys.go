package cmd

import (
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	api_key "github.com/Files-com/files-sdk-go/v2/apikey"
)

func init() {
	RootCmd.AddCommand(ApiKeys())
}

func ApiKeys() *cobra.Command {
	ApiKeys := &cobra.Command{
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
		RunE: func(cmd *cobra.Command, args []string) error {
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
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}

	cmdList.Flags().Int64Var(&paramsApiKeyList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsApiKeyList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsApiKeyList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVar(&fieldsList, "fields", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVar(&formatList, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	ApiKeys.AddCommand(cmdList)
	var fieldsFindCurrent string
	var formatFindCurrent string
	usePagerFindCurrent := true
	cmdFindCurrent := &cobra.Command{
		Use:   "find-current",
		Short: `Show information about current API key.  (Requires current API connection to be using an API key.)`,
		Long:  `Show information about current API key.  (Requires current API connection to be using an API key.)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := api_key.Client{Config: *config}

			var apiKey interface{}
			var err error
			apiKey, err = client.FindCurrent(ctx)
			lib.HandleResponse(ctx, Profile(cmd), apiKey, err, formatFindCurrent, fieldsFindCurrent, usePagerFindCurrent, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}

	cmdFindCurrent.Flags().StringVar(&fieldsFindCurrent, "fields", "", "comma separated list of field names")
	cmdFindCurrent.Flags().StringVar(&formatFindCurrent, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
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
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := api_key.Client{Config: *config}

			var apiKey interface{}
			var err error
			apiKey, err = client.Find(ctx, paramsApiKeyFind)
			lib.HandleResponse(ctx, Profile(cmd), apiKey, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdFind.Flags().Int64Var(&paramsApiKeyFind.Id, "id", 0, "Api Key ID.")

	cmdFind.Flags().StringVar(&fieldsFind, "fields", "", "comma separated list of field names")
	cmdFind.Flags().StringVar(&formatFind, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
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
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := api_key.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsApiKeyCreate.Path = args[0]
			}
			var apiKey interface{}
			var err error
			var ApiKeyCreatePermissionSetErr error
			paramsApiKeyCreate.PermissionSet, ApiKeyCreatePermissionSetErr = lib.FetchKey("permission-set", paramsApiKeyCreate.PermissionSet.Enum(), ApiKeyCreatePermissionSet)
			if ApiKeyCreatePermissionSet != "" && ApiKeyCreatePermissionSetErr != nil {
				return ApiKeyCreatePermissionSetErr
			}
			apiKey, err = client.Create(ctx, paramsApiKeyCreate)
			lib.HandleResponse(ctx, Profile(cmd), apiKey, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdCreate.Flags().Int64Var(&paramsApiKeyCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVar(&paramsApiKeyCreate.Name, "name", "", "Internal name for the API Key.  For your use.")
	cmdCreate.Flags().StringVar(&paramsApiKeyCreate.Description, "description", "", "User-supplied description of API key.")
	lib.TimeVar(cmdCreate.Flags(), paramsApiKeyCreate.ExpiresAt, "expires-at")
	cmdCreate.Flags().StringVar(&ApiKeyCreatePermissionSet, "permission-set", "", fmt.Sprintf("Permissions for this API Key.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations).  Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know. %v", reflect.ValueOf(paramsApiKeyCreate.PermissionSet.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&paramsApiKeyCreate.Path, "path", "", "Folder path restriction for this api key.")

	cmdCreate.Flags().StringVar(&fieldsCreate, "fields", "", "comma separated list of field names")
	cmdCreate.Flags().StringVar(&formatCreate, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
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
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := api_key.Client{Config: *config}

			var apiKey interface{}
			var err error
			var ApiKeyUpdateCurrentPermissionSetErr error
			paramsApiKeyUpdateCurrent.PermissionSet, ApiKeyUpdateCurrentPermissionSetErr = lib.FetchKey("permission-set", paramsApiKeyUpdateCurrent.PermissionSet.Enum(), ApiKeyUpdateCurrentPermissionSet)
			if ApiKeyUpdateCurrentPermissionSet != "" && ApiKeyUpdateCurrentPermissionSetErr != nil {
				return ApiKeyUpdateCurrentPermissionSetErr
			}
			apiKey, err = client.UpdateCurrent(ctx, paramsApiKeyUpdateCurrent)
			lib.HandleResponse(ctx, Profile(cmd), apiKey, err, formatUpdateCurrent, fieldsUpdateCurrent, usePagerUpdateCurrent, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	lib.TimeVar(cmdUpdateCurrent.Flags(), paramsApiKeyUpdateCurrent.ExpiresAt, "expires-at")
	cmdUpdateCurrent.Flags().StringVar(&paramsApiKeyUpdateCurrent.Name, "name", "", "Internal name for the API Key.  For your use.")
	cmdUpdateCurrent.Flags().StringVar(&ApiKeyUpdateCurrentPermissionSet, "permission-set", "", fmt.Sprintf("Permissions for this API Key.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations).  Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know. %v", reflect.ValueOf(paramsApiKeyUpdateCurrent.PermissionSet.Enum()).MapKeys()))

	cmdUpdateCurrent.Flags().StringVar(&fieldsUpdateCurrent, "fields", "", "comma separated list of field names")
	cmdUpdateCurrent.Flags().StringVar(&formatUpdateCurrent, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
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
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := api_key.Client{Config: *config}

			var apiKey interface{}
			var err error
			var ApiKeyUpdatePermissionSetErr error
			paramsApiKeyUpdate.PermissionSet, ApiKeyUpdatePermissionSetErr = lib.FetchKey("permission-set", paramsApiKeyUpdate.PermissionSet.Enum(), ApiKeyUpdatePermissionSet)
			if ApiKeyUpdatePermissionSet != "" && ApiKeyUpdatePermissionSetErr != nil {
				return ApiKeyUpdatePermissionSetErr
			}
			apiKey, err = client.Update(ctx, paramsApiKeyUpdate)
			lib.HandleResponse(ctx, Profile(cmd), apiKey, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsApiKeyUpdate.Id, "id", 0, "Api Key ID.")
	cmdUpdate.Flags().StringVar(&paramsApiKeyUpdate.Name, "name", "", "Internal name for the API Key.  For your use.")
	cmdUpdate.Flags().StringVar(&paramsApiKeyUpdate.Description, "description", "", "User-supplied description of API key.")
	lib.TimeVar(cmdUpdate.Flags(), paramsApiKeyUpdate.ExpiresAt, "expires-at")
	cmdUpdate.Flags().StringVar(&ApiKeyUpdatePermissionSet, "permission-set", "", fmt.Sprintf("Permissions for this API Key.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations).  Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know. %v", reflect.ValueOf(paramsApiKeyUpdate.PermissionSet.Enum()).MapKeys()))

	cmdUpdate.Flags().StringVar(&fieldsUpdate, "fields", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVar(&formatUpdate, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	ApiKeys.AddCommand(cmdUpdate)
	var fieldsDeleteCurrent string
	var formatDeleteCurrent string
	usePagerDeleteCurrent := true
	cmdDeleteCurrent := &cobra.Command{
		Use:   "delete-current",
		Short: `Delete current API key.  (Requires current API connection to be using an API key.)`,
		Long:  `Delete current API key.  (Requires current API connection to be using an API key.)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := api_key.Client{Config: *config}

			var err error
			err = client.DeleteCurrent(ctx)
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}

	cmdDeleteCurrent.Flags().StringVar(&fieldsDeleteCurrent, "fields", "", "comma separated list of field names")
	cmdDeleteCurrent.Flags().StringVar(&formatDeleteCurrent, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
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
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := api_key.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsApiKeyDelete)
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsApiKeyDelete.Id, "id", 0, "Api Key ID.")

	cmdDelete.Flags().StringVar(&fieldsDelete, "fields", "", "comma separated list of field names")
	cmdDelete.Flags().StringVar(&formatDelete, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	ApiKeys.AddCommand(cmdDelete)
	return ApiKeys
}
