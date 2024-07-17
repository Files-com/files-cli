package cmd

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	api_key "github.com/Files-com/files-sdk-go/v3/apikey"
	"github.com/spf13/cobra"
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
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsApiKeyList := files_sdk.ApiKeyListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List API Keys",
		Long:    `List API Keys`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsApiKeyList
			params.MaxPages = MaxPagesList

			client := api_key.Client{Config: config}
			it, err := client.List(params, files_sdk.WithContext(ctx))
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger)
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyList) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyList, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatList), fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdList.Flags().Int64Var(&paramsApiKeyList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsApiKeyList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsApiKeyList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVar(&paramsApiKeyList.Action, "action", "", "")
	cmdList.Flags().Int64Var(&paramsApiKeyList.Page, "page", 0, "")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	ApiKeys.AddCommand(cmdList)
	var fieldsFindCurrent []string
	var formatFindCurrent []string
	usePagerFindCurrent := true
	cmdFindCurrent := &cobra.Command{
		Use:   "find-current",
		Short: `Show information about current API key.  (Requires current API connection to be using an API key.)`,
		Long:  `Show information about current API key.  (Requires current API connection to be using an API key.)`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := api_key.Client{Config: config}

			var apiKey interface{}
			var err error
			apiKey, err = client.FindCurrent(files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), apiKey, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFindCurrent), fieldsFindCurrent, usePagerFindCurrent, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}

	cmdFindCurrent.Flags().StringSliceVar(&fieldsFindCurrent, "fields", []string{}, "comma separated list of field names")
	cmdFindCurrent.Flags().StringSliceVar(&formatFindCurrent, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFindCurrent.Flags().BoolVar(&usePagerFindCurrent, "use-pager", usePagerFindCurrent, "Use $PAGER (.ie less, more, etc)")

	ApiKeys.AddCommand(cmdFindCurrent)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsApiKeyFind := files_sdk.ApiKeyFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show API Key`,
		Long:  `Show API Key`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := api_key.Client{Config: config}

			var apiKey interface{}
			var err error
			apiKey, err = client.Find(paramsApiKeyFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), apiKey, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsApiKeyFind.Id, "id", 0, "Api Key ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	ApiKeys.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsApiKeyCreate := files_sdk.ApiKeyCreateParams{}
	ApiKeyCreatePermissionSet := ""

	cmdCreate := &cobra.Command{
		Use:   "create [path]",
		Short: `Create API Key`,
		Long:  `Create API Key`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := api_key.Client{Config: config}

			var ApiKeyCreatePermissionSetErr error
			paramsApiKeyCreate.PermissionSet, ApiKeyCreatePermissionSetErr = lib.FetchKey("permission-set", paramsApiKeyCreate.PermissionSet.Enum(), ApiKeyCreatePermissionSet)
			if ApiKeyCreatePermissionSet != "" && ApiKeyCreatePermissionSetErr != nil {
				return ApiKeyCreatePermissionSetErr
			}

			if paramsApiKeyCreate.ExpiresAt.IsZero() {
				paramsApiKeyCreate.ExpiresAt = nil
			}

			if len(args) > 0 && args[0] != "" {
				paramsApiKeyCreate.Path = args[0]
			}
			var apiKey interface{}
			var err error
			apiKey, err = client.Create(paramsApiKeyCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), apiKey, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().Int64Var(&paramsApiKeyCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVar(&paramsApiKeyCreate.Description, "description", "", "User-supplied description of API key.")
	paramsApiKeyCreate.ExpiresAt = &time.Time{}
	lib.TimeVar(cmdCreate.Flags(), paramsApiKeyCreate.ExpiresAt, "expires-at", "API Key expiration date")
	cmdCreate.Flags().StringVar(&ApiKeyCreatePermissionSet, "permission-set", "", fmt.Sprintf("Permissions for this API Key. It must be full for site-wide API Keys.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations).  Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know. %v", reflect.ValueOf(paramsApiKeyCreate.PermissionSet.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&paramsApiKeyCreate.Name, "name", "", "Internal name for the API Key.  For your use.")
	cmdCreate.Flags().StringVar(&paramsApiKeyCreate.Path, "path", "", "Folder path restriction for this API key.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	ApiKeys.AddCommand(cmdCreate)
	var fieldsUpdateCurrent []string
	var formatUpdateCurrent []string
	usePagerUpdateCurrent := true
	paramsApiKeyUpdateCurrent := files_sdk.ApiKeyUpdateCurrentParams{}
	ApiKeyUpdateCurrentPermissionSet := ""

	cmdUpdateCurrent := &cobra.Command{
		Use:   "update-current",
		Short: `Update current API key.  (Requires current API connection to be using an API key.)`,
		Long:  `Update current API key.  (Requires current API connection to be using an API key.)`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := api_key.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.ApiKeyUpdateCurrentParams{})
			if convertErr != nil {
				return convertErr
			}

			var ApiKeyUpdateCurrentPermissionSetErr error
			paramsApiKeyUpdateCurrent.PermissionSet, ApiKeyUpdateCurrentPermissionSetErr = lib.FetchKey("permission-set", paramsApiKeyUpdateCurrent.PermissionSet.Enum(), ApiKeyUpdateCurrentPermissionSet)
			if ApiKeyUpdateCurrentPermissionSet != "" && ApiKeyUpdateCurrentPermissionSetErr != nil {
				return ApiKeyUpdateCurrentPermissionSetErr
			}

			if cmd.Flags().Changed("expires-at") {
				lib.FlagUpdate(cmd, "expires_at", paramsApiKeyUpdateCurrent.ExpiresAt, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsApiKeyUpdateCurrent.Name, mapParams)
			}
			if cmd.Flags().Changed("permission-set") {
				lib.FlagUpdate(cmd, "permission_set", paramsApiKeyUpdateCurrent.PermissionSet, mapParams)
			}

			if paramsApiKeyUpdateCurrent.ExpiresAt.IsZero() {
				paramsApiKeyUpdateCurrent.ExpiresAt = nil
			}

			var apiKey interface{}
			var err error
			apiKey, err = client.UpdateCurrentWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), apiKey, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdateCurrent), fieldsUpdateCurrent, usePagerUpdateCurrent, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	paramsApiKeyUpdateCurrent.ExpiresAt = &time.Time{}
	lib.TimeVar(cmdUpdateCurrent.Flags(), paramsApiKeyUpdateCurrent.ExpiresAt, "expires-at", "API Key expiration date")
	cmdUpdateCurrent.Flags().StringVar(&paramsApiKeyUpdateCurrent.Name, "name", "", "Internal name for the API Key.  For your use.")
	cmdUpdateCurrent.Flags().StringVar(&ApiKeyUpdateCurrentPermissionSet, "permission-set", "", fmt.Sprintf("Permissions for this API Key. It must be full for site-wide API Keys.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations).  Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know. %v", reflect.ValueOf(paramsApiKeyUpdateCurrent.PermissionSet.Enum()).MapKeys()))

	cmdUpdateCurrent.Flags().StringSliceVar(&fieldsUpdateCurrent, "fields", []string{}, "comma separated list of field names")
	cmdUpdateCurrent.Flags().StringSliceVar(&formatUpdateCurrent, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdateCurrent.Flags().BoolVar(&usePagerUpdateCurrent, "use-pager", usePagerUpdateCurrent, "Use $PAGER (.ie less, more, etc)")

	ApiKeys.AddCommand(cmdUpdateCurrent)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsApiKeyUpdate := files_sdk.ApiKeyUpdateParams{}
	ApiKeyUpdatePermissionSet := ""

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update API Key`,
		Long:  `Update API Key`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := api_key.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.ApiKeyUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			var ApiKeyUpdatePermissionSetErr error
			paramsApiKeyUpdate.PermissionSet, ApiKeyUpdatePermissionSetErr = lib.FetchKey("permission-set", paramsApiKeyUpdate.PermissionSet.Enum(), ApiKeyUpdatePermissionSet)
			if ApiKeyUpdatePermissionSet != "" && ApiKeyUpdatePermissionSetErr != nil {
				return ApiKeyUpdatePermissionSetErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsApiKeyUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("description") {
				lib.FlagUpdate(cmd, "description", paramsApiKeyUpdate.Description, mapParams)
			}
			if cmd.Flags().Changed("expires-at") {
				lib.FlagUpdate(cmd, "expires_at", paramsApiKeyUpdate.ExpiresAt, mapParams)
			}
			if cmd.Flags().Changed("permission-set") {
				lib.FlagUpdate(cmd, "permission_set", paramsApiKeyUpdate.PermissionSet, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsApiKeyUpdate.Name, mapParams)
			}

			if paramsApiKeyUpdate.ExpiresAt.IsZero() {
				paramsApiKeyUpdate.ExpiresAt = nil
			}

			var apiKey interface{}
			var err error
			apiKey, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), apiKey, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsApiKeyUpdate.Id, "id", 0, "Api Key ID.")
	cmdUpdate.Flags().StringVar(&paramsApiKeyUpdate.Description, "description", "", "User-supplied description of API key.")
	paramsApiKeyUpdate.ExpiresAt = &time.Time{}
	lib.TimeVar(cmdUpdate.Flags(), paramsApiKeyUpdate.ExpiresAt, "expires-at", "API Key expiration date")
	cmdUpdate.Flags().StringVar(&ApiKeyUpdatePermissionSet, "permission-set", "", fmt.Sprintf("Permissions for this API Key. It must be full for site-wide API Keys.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations).  Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know. %v", reflect.ValueOf(paramsApiKeyUpdate.PermissionSet.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&paramsApiKeyUpdate.Name, "name", "", "Internal name for the API Key.  For your use.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	ApiKeys.AddCommand(cmdUpdate)
	var fieldsDeleteCurrent []string
	var formatDeleteCurrent []string
	usePagerDeleteCurrent := true
	cmdDeleteCurrent := &cobra.Command{
		Use:   "delete-current",
		Short: `Delete current API key.  (Requires current API connection to be using an API key.)`,
		Long:  `Delete current API key.  (Requires current API connection to be using an API key.)`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := api_key.Client{Config: config}

			var err error
			err = client.DeleteCurrent(files_sdk.WithContext(ctx))
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}

	cmdDeleteCurrent.Flags().StringSliceVar(&fieldsDeleteCurrent, "fields", []string{}, "comma separated list of field names")
	cmdDeleteCurrent.Flags().StringSliceVar(&formatDeleteCurrent, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDeleteCurrent.Flags().BoolVar(&usePagerDeleteCurrent, "use-pager", usePagerDeleteCurrent, "Use $PAGER (.ie less, more, etc)")

	ApiKeys.AddCommand(cmdDeleteCurrent)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsApiKeyDelete := files_sdk.ApiKeyDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete API Key`,
		Long:  `Delete API Key`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := api_key.Client{Config: config}

			var err error
			err = client.Delete(paramsApiKeyDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsApiKeyDelete.Id, "id", 0, "Api Key ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	ApiKeys.AddCommand(cmdDelete)
	return ApiKeys
}
