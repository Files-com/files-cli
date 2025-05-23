package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/group"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Groups())
}

func Groups() *cobra.Command {
	Groups := &cobra.Command{
		Use:  "groups [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command groups\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsGroupList := files_sdk.GroupListParams{}
	var MaxPagesList int64
	listIncludeParentSiteGroups := true

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Groups",
		Long:    `List Groups`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsGroupList
			params.MaxPages = MaxPagesList

			if cmd.Flags().Changed("include-parent-site-groups") {
				params.IncludeParentSiteGroups = flib.Bool(listIncludeParentSiteGroups)
			}

			client := group.Client{Config: config}
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
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyList) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyList, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatList), fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdList.Flags().StringVar(&paramsGroupList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsGroupList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVar(&paramsGroupList.Ids, "ids", "", "Comma-separated list of group ids to include in results.")
	cmdList.Flags().BoolVar(&listIncludeParentSiteGroups, "include-parent-site-groups", listIncludeParentSiteGroups, "Include groups from the parent site.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	Groups.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsGroupFind := files_sdk.GroupFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Group`,
		Long:  `Show Group`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := group.Client{Config: config}

			var group interface{}
			var err error
			group, err = client.Find(paramsGroupFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), group, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsGroupFind.Id, "id", 0, "Group ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	Groups.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createFtpPermission := true
	createSftpPermission := true
	createDavPermission := true
	createRestapiPermission := true
	paramsGroupCreate := files_sdk.GroupCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Group`,
		Long:  `Create Group`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := group.Client{Config: config}

			if cmd.Flags().Changed("ftp-permission") {
				paramsGroupCreate.FtpPermission = flib.Bool(createFtpPermission)
			}
			if cmd.Flags().Changed("sftp-permission") {
				paramsGroupCreate.SftpPermission = flib.Bool(createSftpPermission)
			}
			if cmd.Flags().Changed("dav-permission") {
				paramsGroupCreate.DavPermission = flib.Bool(createDavPermission)
			}
			if cmd.Flags().Changed("restapi-permission") {
				paramsGroupCreate.RestapiPermission = flib.Bool(createRestapiPermission)
			}

			var group interface{}
			var err error
			group, err = client.Create(paramsGroupCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), group, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsGroupCreate.Notes, "notes", "", "Group notes.")
	cmdCreate.Flags().StringVar(&paramsGroupCreate.UserIds, "user-ids", "", "A list of user ids. If sent as a string, should be comma-delimited.")
	cmdCreate.Flags().StringVar(&paramsGroupCreate.AdminIds, "admin-ids", "", "A list of group admin user ids. If sent as a string, should be comma-delimited.")
	cmdCreate.Flags().BoolVar(&createFtpPermission, "ftp-permission", createFtpPermission, "If true, users in this group can use FTP to login.  This will override a false value of `ftp_permission` on the user level.")
	cmdCreate.Flags().BoolVar(&createSftpPermission, "sftp-permission", createSftpPermission, "If true, users in this group can use SFTP to login.  This will override a false value of `sftp_permission` on the user level.")
	cmdCreate.Flags().BoolVar(&createDavPermission, "dav-permission", createDavPermission, "If true, users in this group can use WebDAV to login.  This will override a false value of `dav_permission` on the user level.")
	cmdCreate.Flags().BoolVar(&createRestapiPermission, "restapi-permission", createRestapiPermission, "If true, users in this group can use the REST API to login.  This will override a false value of `restapi_permission` on the user level.")
	cmdCreate.Flags().StringVar(&paramsGroupCreate.AllowedIps, "allowed-ips", "", "A list of allowed IPs if applicable.  Newline delimited")
	cmdCreate.Flags().StringVar(&paramsGroupCreate.Name, "name", "", "Group name.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Groups.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateFtpPermission := true
	updateSftpPermission := true
	updateDavPermission := true
	updateRestapiPermission := true
	paramsGroupUpdate := files_sdk.GroupUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Group`,
		Long:  `Update Group`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := group.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.GroupUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsGroupUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("notes") {
				lib.FlagUpdate(cmd, "notes", paramsGroupUpdate.Notes, mapParams)
			}
			if cmd.Flags().Changed("user-ids") {
				lib.FlagUpdate(cmd, "user_ids", paramsGroupUpdate.UserIds, mapParams)
			}
			if cmd.Flags().Changed("admin-ids") {
				lib.FlagUpdate(cmd, "admin_ids", paramsGroupUpdate.AdminIds, mapParams)
			}
			if cmd.Flags().Changed("ftp-permission") {
				mapParams["ftp_permission"] = updateFtpPermission
			}
			if cmd.Flags().Changed("sftp-permission") {
				mapParams["sftp_permission"] = updateSftpPermission
			}
			if cmd.Flags().Changed("dav-permission") {
				mapParams["dav_permission"] = updateDavPermission
			}
			if cmd.Flags().Changed("restapi-permission") {
				mapParams["restapi_permission"] = updateRestapiPermission
			}
			if cmd.Flags().Changed("allowed-ips") {
				lib.FlagUpdate(cmd, "allowed_ips", paramsGroupUpdate.AllowedIps, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsGroupUpdate.Name, mapParams)
			}

			var group interface{}
			var err error
			group, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), group, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsGroupUpdate.Id, "id", 0, "Group ID.")
	cmdUpdate.Flags().StringVar(&paramsGroupUpdate.Notes, "notes", "", "Group notes.")
	cmdUpdate.Flags().StringVar(&paramsGroupUpdate.UserIds, "user-ids", "", "A list of user ids. If sent as a string, should be comma-delimited.")
	cmdUpdate.Flags().StringVar(&paramsGroupUpdate.AdminIds, "admin-ids", "", "A list of group admin user ids. If sent as a string, should be comma-delimited.")
	cmdUpdate.Flags().BoolVar(&updateFtpPermission, "ftp-permission", updateFtpPermission, "If true, users in this group can use FTP to login.  This will override a false value of `ftp_permission` on the user level.")
	cmdUpdate.Flags().BoolVar(&updateSftpPermission, "sftp-permission", updateSftpPermission, "If true, users in this group can use SFTP to login.  This will override a false value of `sftp_permission` on the user level.")
	cmdUpdate.Flags().BoolVar(&updateDavPermission, "dav-permission", updateDavPermission, "If true, users in this group can use WebDAV to login.  This will override a false value of `dav_permission` on the user level.")
	cmdUpdate.Flags().BoolVar(&updateRestapiPermission, "restapi-permission", updateRestapiPermission, "If true, users in this group can use the REST API to login.  This will override a false value of `restapi_permission` on the user level.")
	cmdUpdate.Flags().StringVar(&paramsGroupUpdate.AllowedIps, "allowed-ips", "", "A list of allowed IPs if applicable.  Newline delimited")
	cmdUpdate.Flags().StringVar(&paramsGroupUpdate.Name, "name", "", "Group name.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	Groups.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsGroupDelete := files_sdk.GroupDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Group`,
		Long:  `Delete Group`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := group.Client{Config: config}

			var err error
			err = client.Delete(paramsGroupDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsGroupDelete.Id, "id", 0, "Group ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Groups.AddCommand(cmdDelete)
	return Groups
}
