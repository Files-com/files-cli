package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	history_export "github.com/Files-com/files-sdk-go/v2/historyexport"
)

var (
	HistoryExports = &cobra.Command{}
)

func HistoryExportsInit() {
	HistoryExports = &cobra.Command{
		Use:  "history-exports [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command history-exports\n\t%v", args[0])
		},
	}
	var fieldsFind string
	var formatFind string
	usePagerFind := true
	paramsHistoryExportFind := files_sdk.HistoryExportFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show History Export`,
		Long:  `Show History Export`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := history_export.Client{Config: *config}

			var historyExport interface{}
			var err error
			historyExport, err = client.Find(ctx, paramsHistoryExportFind)
			lib.HandleResponse(ctx, historyExport, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdFind.Flags().Int64Var(&paramsHistoryExportFind.Id, "id", 0, "History Export ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	HistoryExports.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	usePagerCreate := true
	paramsHistoryExportCreate := files_sdk.HistoryExportCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create History Export`,
		Long:  `Create History Export`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := history_export.Client{Config: *config}

			var historyExport interface{}
			var err error
			historyExport, err = client.Create(ctx, paramsHistoryExportCreate)
			lib.HandleResponse(ctx, historyExport, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdCreate.Flags().Int64Var(&paramsHistoryExportCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	lib.TimeVar(cmdCreate.Flags(), paramsHistoryExportCreate.StartAt, "start-at")
	lib.TimeVar(cmdCreate.Flags(), paramsHistoryExportCreate.EndAt, "end-at")
	cmdCreate.Flags().StringVar(&paramsHistoryExportCreate.QueryAction, "query-action", "", "Filter results by this this action type. Valid values: `create`, `read`, `update`, `destroy`, `move`, `login`, `failedlogin`, `copy`, `user_create`, `user_update`, `user_destroy`, `group_create`, `group_update`, `group_destroy`, `permission_create`, `permission_destroy`, `api_key_create`, `api_key_update`, `api_key_destroy`")
	cmdCreate.Flags().StringVar(&paramsHistoryExportCreate.QueryInterface, "query-interface", "", "Filter results by this this interface type. Valid values: `web`, `ftp`, `robot`, `jsapi`, `webdesktopapi`, `sftp`, `dav`, `desktop`, `restapi`, `scim`, `office`, `mobile`, `as2`")
	cmdCreate.Flags().StringVar(&paramsHistoryExportCreate.QueryUserId, "query-user-id", "", "Return results that are actions performed by the user indiciated by this User ID")
	cmdCreate.Flags().StringVar(&paramsHistoryExportCreate.QueryFileId, "query-file-id", "", "Return results that are file actions related to the file indicated by this File ID")
	cmdCreate.Flags().StringVar(&paramsHistoryExportCreate.QueryParentId, "query-parent-id", "", "Return results that are file actions inside the parent folder specified by this folder ID")
	cmdCreate.Flags().StringVar(&paramsHistoryExportCreate.QueryPath, "query-path", "", "Return results that are file actions related to this path.")
	cmdCreate.Flags().StringVar(&paramsHistoryExportCreate.QueryFolder, "query-folder", "", "Return results that are file actions related to files or folders inside this folder path.")
	cmdCreate.Flags().StringVar(&paramsHistoryExportCreate.QuerySrc, "query-src", "", "Return results that are file moves originating from this path.")
	cmdCreate.Flags().StringVar(&paramsHistoryExportCreate.QueryDestination, "query-destination", "", "Return results that are file moves with this path as destination.")
	cmdCreate.Flags().StringVar(&paramsHistoryExportCreate.QueryIp, "query-ip", "", "Filter results by this IP address.")
	cmdCreate.Flags().StringVar(&paramsHistoryExportCreate.QueryUsername, "query-username", "", "Filter results by this username.")
	cmdCreate.Flags().StringVar(&paramsHistoryExportCreate.QueryFailureType, "query-failure-type", "", "If searching for Histories about login failures, this parameter restricts results to failures of this specific type.  Valid values: `expired_trial`, `account_overdue`, `locked_out`, `ip_mismatch`, `password_mismatch`, `site_mismatch`, `username_not_found`, `none`, `no_ftp_permission`, `no_web_permission`, `no_directory`, `errno_enoent`, `no_sftp_permission`, `no_dav_permission`, `no_restapi_permission`, `key_mismatch`, `region_mismatch`, `expired_access`, `desktop_ip_mismatch`, `desktop_api_key_not_used_quickly_enough`, `disabled`, `country_mismatch`")
	cmdCreate.Flags().StringVar(&paramsHistoryExportCreate.QueryTargetId, "query-target-id", "", "If searching for Histories about specific objects (such as Users, or API Keys), this paremeter restricts results to objects that match this ID.")
	cmdCreate.Flags().StringVar(&paramsHistoryExportCreate.QueryTargetName, "query-target-name", "", "If searching for Histories about Users, Groups or other objects with names, this parameter restricts results to objects with this name/username.")
	cmdCreate.Flags().StringVar(&paramsHistoryExportCreate.QueryTargetPermission, "query-target-permission", "", "If searching for Histories about Permisisons, this parameter restricts results to permissions of this level.")
	cmdCreate.Flags().StringVar(&paramsHistoryExportCreate.QueryTargetUserId, "query-target-user-id", "", "If searching for Histories about API keys, this parameter restricts results to API keys created by/for this user ID.")
	cmdCreate.Flags().StringVar(&paramsHistoryExportCreate.QueryTargetUsername, "query-target-username", "", "If searching for Histories about API keys, this parameter restricts results to API keys created by/for this username.")
	cmdCreate.Flags().StringVar(&paramsHistoryExportCreate.QueryTargetPlatform, "query-target-platform", "", "If searching for Histories about API keys, this parameter restricts results to API keys associated with this platform.")
	cmdCreate.Flags().StringVar(&paramsHistoryExportCreate.QueryTargetPermissionSet, "query-target-permission-set", "", "If searching for Histories about API keys, this parameter restricts results to API keys with this permission set.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	HistoryExports.AddCommand(cmdCreate)
}
