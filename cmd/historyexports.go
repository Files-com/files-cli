package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	history_export "github.com/Files-com/files-sdk-go/historyexport"
)

var (
	HistoryExports = &cobra.Command{}
)

func HistoryExportsInit() {
	HistoryExports = &cobra.Command{
		Use:  "history-exports [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
	var fieldsFind string
	paramsHistoryExportFind := files_sdk.HistoryExportFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := history_export.Client{Config: *ctx.GetConfig()}

			result, err := client.Find(paramsHistoryExportFind)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsFind)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsHistoryExportFind.Id, "id", "i", 0, "History Export ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	HistoryExports.AddCommand(cmdFind)
	var fieldsCreate string
	paramsHistoryExportCreate := files_sdk.HistoryExportCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := history_export.Client{Config: *ctx.GetConfig()}

			result, err := client.Create(paramsHistoryExportCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdCreate.Flags().Int64VarP(&paramsHistoryExportCreate.UserId, "user-id", "", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	lib.TimeVarP(cmdCreate.Flags(), &paramsHistoryExportCreate.StartAt, "start-at", "")
	lib.TimeVarP(cmdCreate.Flags(), &paramsHistoryExportCreate.EndAt, "end-at", "e")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryAction, "query-action", "a", "", "Filter results by this this action type. Valid values: `create`, `read`, `update`, `destroy`, `move`, `login`, `failedlogin`, `copy`, `user_create`, `user_update`, `user_destroy`, `group_create`, `group_update`, `group_destroy`, `permission_create`, `permission_destroy`, `api_key_create`, `api_key_update`, `api_key_destroy`")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryInterface, "query-interface", "n", "", "Filter results by this this interface type. Valid values: `web`, `ftp`, `robot`, `jsapi`, `webdesktopapi`, `sftp`, `dav`, `desktop`, `restapi`, `scim`, `office`")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryUserId, "query-user-id", "", "", "Return results that are actions performed by the user indiciated by this User ID")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryFileId, "query-file-id", "i", "", "Return results that are file actions related to the file indicated by this File ID")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryParentId, "query-parent-id", "", "", "Return results that are file actions inside the parent folder specified by this folder ID")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryPath, "query-path", "", "", "Return results that are file actions related to this path.")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryFolder, "query-folder", "f", "", "Return results that are file actions related to files or folders inside this folder path.")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QuerySrc, "query-src", "s", "", "Return results that are file moves originating from this path.")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryDestination, "query-destination", "d", "", "Return results that are file moves with this path as destination.")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryIp, "query-ip", "p", "", "Filter results by this IP address.")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryUsername, "query-username", "", "", "Filter results by this username.")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryFailureType, "query-failure-type", "t", "", "If searching for Histories about login failures, this parameter restricts results to failures of this specific type.  Valid values: `expired_trial`, `account_overdue`, `locked_out`, `ip_mismatch`, `password_mismatch`, `site_mismatch`, `username_not_found`, `none`, `no_ftp_permission`, `no_web_permission`, `no_directory`, `errno_enoent`, `no_sftp_permission`, `no_dav_permission`, `no_restapi_permission`, `key_mismatch`, `region_mismatch`, `expired_access`, `desktop_ip_mismatch`, `desktop_api_key_not_used_quickly_enough`, `disabled`, `country_mismatch`")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryTargetId, "query-target-id", "", "", "If searching for Histories about specific objects (such as Users, or API Keys), this paremeter restricts results to objects that match this ID.")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryTargetName, "query-target-name", "", "", "If searching for Histories about Users, Groups or other objects with names, this parameter restricts results to objects with this name/username.")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryTargetPermission, "query-target-permission", "r", "", "If searching for Histories about Permisisons, this parameter restricts results to permissions of this level.")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryTargetUserId, "query-target-user-id", "", "", "If searching for Histories about API keys, this parameter restricts results to API keys created by/for this user ID.")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryTargetUsername, "query-target-username", "u", "", "If searching for Histories about API keys, this parameter restricts results to API keys created by/for this username.")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryTargetPlatform, "query-target-platform", "l", "", "If searching for Histories about API keys, this parameter restricts results to API keys associated with this platform.")
	cmdCreate.Flags().StringVarP(&paramsHistoryExportCreate.QueryTargetPermissionSet, "query-target-permission-set", "", "", "If searching for Histories about API keys, this parameter restricts results to API keys with this permission set.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	HistoryExports.AddCommand(cmdCreate)
}
