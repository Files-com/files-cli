package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	api_key "github.com/Files-com/files-sdk-go/apikey"
)

var (
	ApiKeys = &cobra.Command{}
)

func ApiKeysInit() {
	ApiKeys = &cobra.Command{
		Use:  "api-keys [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
	var fieldsList string
	paramsApiKeyList := files_sdk.ApiKeyListParams{}
	var MaxPagesList int64
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			params := paramsApiKeyList
			params.MaxPages = MaxPagesList
			client := api_key.Client{Config: *ctx.GetConfig()}
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
	cmdList.Flags().Int64VarP(&paramsApiKeyList.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVarP(&paramsApiKeyList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsApiKeyList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	ApiKeys.AddCommand(cmdList)
	var fieldsFindCurrent string
	cmdFindCurrent := &cobra.Command{
		Use: "find-current",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := api_key.Client{Config: *ctx.GetConfig()}

			result, err := client.FindCurrent()
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsFindCurrent)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}

	cmdFindCurrent.Flags().StringVarP(&fieldsFindCurrent, "fields", "", "", "comma separated list of field names")
	ApiKeys.AddCommand(cmdFindCurrent)
	var fieldsFind string
	paramsApiKeyFind := files_sdk.ApiKeyFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := api_key.Client{Config: *ctx.GetConfig()}

			result, err := client.Find(paramsApiKeyFind)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsFind)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsApiKeyFind.Id, "id", "i", 0, "Api Key ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	ApiKeys.AddCommand(cmdFind)
	var fieldsCreate string
	paramsApiKeyCreate := files_sdk.ApiKeyCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := api_key.Client{Config: *ctx.GetConfig()}

			if len(args) > 0 && args[0] != "" {
				paramsApiKeyCreate.Path = args[0]
			}
			result, err := client.Create(paramsApiKeyCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdCreate.Flags().Int64VarP(&paramsApiKeyCreate.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVarP(&paramsApiKeyCreate.Name, "name", "n", "", "Internal name for the API Key.  For your use.")
	lib.TimeVarP(cmdCreate.Flags(), &paramsApiKeyCreate.ExpiresAt, "expires-at", "e")
	cmdCreate.Flags().StringVarP(&paramsApiKeyCreate.PermissionSet, "permission-set", "r", "", "Permissions for this API Key.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations).  Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know.")
	cmdCreate.Flags().StringVarP(&paramsApiKeyCreate.Path, "path", "p", "", "Folder path restriction for this api key.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	ApiKeys.AddCommand(cmdCreate)
	var fieldsUpdateCurrent string
	paramsApiKeyUpdateCurrent := files_sdk.ApiKeyUpdateCurrentParams{}
	cmdUpdateCurrent := &cobra.Command{
		Use: "update-current",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := api_key.Client{Config: *ctx.GetConfig()}

			result, err := client.UpdateCurrent(paramsApiKeyUpdateCurrent)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsUpdateCurrent)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	lib.TimeVarP(cmdUpdateCurrent.Flags(), &paramsApiKeyUpdateCurrent.ExpiresAt, "expires-at", "e")
	cmdUpdateCurrent.Flags().StringVarP(&paramsApiKeyUpdateCurrent.Name, "name", "n", "", "Internal name for the API Key.  For your use.")
	cmdUpdateCurrent.Flags().StringVarP(&paramsApiKeyUpdateCurrent.PermissionSet, "permission-set", "p", "", "Permissions for this API Key.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations).  Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know.")

	cmdUpdateCurrent.Flags().StringVarP(&fieldsUpdateCurrent, "fields", "", "", "comma separated list of field names")
	ApiKeys.AddCommand(cmdUpdateCurrent)
	var fieldsUpdate string
	paramsApiKeyUpdate := files_sdk.ApiKeyUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := api_key.Client{Config: *ctx.GetConfig()}

			result, err := client.Update(paramsApiKeyUpdate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsUpdate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsApiKeyUpdate.Id, "id", "i", 0, "Api Key ID.")
	cmdUpdate.Flags().StringVarP(&paramsApiKeyUpdate.Name, "name", "n", "", "Internal name for the API Key.  For your use.")
	lib.TimeVarP(cmdUpdate.Flags(), &paramsApiKeyUpdate.ExpiresAt, "expires-at", "e")
	cmdUpdate.Flags().StringVarP(&paramsApiKeyUpdate.PermissionSet, "permission-set", "p", "", "Permissions for this API Key.  Keys with the `desktop_app` permission set only have the ability to do the functions provided in our Desktop App (File and Share Link operations).  Additional permission sets may become available in the future, such as for a Site Admin to give a key with no administrator privileges.  If you have ideas for permission sets, please let us know.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	ApiKeys.AddCommand(cmdUpdate)
	var fieldsDeleteCurrent string
	cmdDeleteCurrent := &cobra.Command{
		Use: "delete-current",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := api_key.Client{Config: *ctx.GetConfig()}

			result, err := client.DeleteCurrent()
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsDeleteCurrent)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}

	cmdDeleteCurrent.Flags().StringVarP(&fieldsDeleteCurrent, "fields", "", "", "comma separated list of field names")
	ApiKeys.AddCommand(cmdDeleteCurrent)
	var fieldsDelete string
	paramsApiKeyDelete := files_sdk.ApiKeyDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := api_key.Client{Config: *ctx.GetConfig()}

			result, err := client.Delete(paramsApiKeyDelete)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsApiKeyDelete.Id, "id", "i", 0, "Api Key ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	ApiKeys.AddCommand(cmdDelete)
}
