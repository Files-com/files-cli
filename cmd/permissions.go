package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/permission"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Permissions())
}

func Permissions() *cobra.Command {
	Permissions := &cobra.Command{
		Use:  "permissions [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command permissions\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsPermissionList := files_sdk.PermissionListParams{}
	var MaxPagesList int64
	listIncludeGroups := true

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Permissions",
		Long:    `List Permissions`,
		Args:    cobra.RangeArgs(0, 1),
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsPermissionList
			params.MaxPages = MaxPagesList
			if len(args) > 0 && args[0] != "" {
				params.Path = args[0]
			}

			if cmd.Flags().Changed("include-groups") {
				params.IncludeGroups = flib.Bool(listIncludeGroups)
			}

			client := permission.Client{Config: config}
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

	cmdList.Flags().StringVar(&paramsPermissionList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsPermissionList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVar(&paramsPermissionList.Path, "path", "", "Permission path.  If provided, will scope all permissions(including upward) to this path.")
	cmdList.Flags().BoolVar(&listIncludeGroups, "include-groups", listIncludeGroups, "If searching by user or group, also include user's permissions that are inherited from its groups?")
	cmdList.Flags().StringVar(&paramsPermissionList.GroupId, "group-id", "", "")
	cmdList.Flags().StringVar(&paramsPermissionList.PartnerId, "partner-id", "", "")
	cmdList.Flags().StringVar(&paramsPermissionList.UserId, "user-id", "", "")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	Permissions.AddCommand(cmdList)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createRecursive := true
	paramsPermissionCreate := files_sdk.PermissionCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create [path]",
		Short: `Create Permission`,
		Long:  `Create Permission`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := permission.Client{Config: config}

			if cmd.Flags().Changed("recursive") {
				paramsPermissionCreate.Recursive = flib.Bool(createRecursive)
			}

			if len(args) > 0 && args[0] != "" {
				paramsPermissionCreate.Path = args[0]
			}
			var permission interface{}
			var err error
			permission, err = client.Create(paramsPermissionCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), permission, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsPermissionCreate.Path, "path", "", "Folder path")
	cmdCreate.Flags().Int64Var(&paramsPermissionCreate.GroupId, "group-id", 0, "Group ID. Provide `group_name` or `group_id`")
	cmdCreate.Flags().StringVar(&paramsPermissionCreate.Permission, "permission", "", "Permission type.  Can be `admin`, `full`, `readonly`, `writeonly`, `list`, or `history`")
	cmdCreate.Flags().BoolVar(&createRecursive, "recursive", createRecursive, "Apply to subfolders recursively?")
	cmdCreate.Flags().Int64Var(&paramsPermissionCreate.PartnerId, "partner-id", 0, "Partner ID if this Permission belongs to a partner.")
	cmdCreate.Flags().Int64Var(&paramsPermissionCreate.UserId, "user-id", 0, "User ID.  Provide `username` or `user_id`")
	cmdCreate.Flags().StringVar(&paramsPermissionCreate.Username, "username", "", "User username.  Provide `username` or `user_id`")
	cmdCreate.Flags().StringVar(&paramsPermissionCreate.GroupName, "group-name", "", "Group name.  Provide `group_name` or `group_id`")
	cmdCreate.Flags().Int64Var(&paramsPermissionCreate.SiteId, "site-id", 0, "Site ID. If not provided, will default to current site. Used when creating a permission for a child site.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Permissions.AddCommand(cmdCreate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsPermissionDelete := files_sdk.PermissionDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Permission`,
		Long:  `Delete Permission`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := permission.Client{Config: config}

			var err error
			err = client.Delete(paramsPermissionDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsPermissionDelete.Id, "id", 0, "Permission ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Permissions.AddCommand(cmdDelete)
	return Permissions
}
