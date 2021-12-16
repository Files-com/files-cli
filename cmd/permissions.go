package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	flib "github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/Files-com/files-sdk-go/v2/permission"
)

var (
	Permissions = &cobra.Command{}
)

func PermissionsInit() {
	Permissions = &cobra.Command{
		Use:  "permissions [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command permissions\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsPermissionList := files_sdk.PermissionListParams{}
	var MaxPagesList int64
	listIncludeGroups := false

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsPermissionList
			params.MaxPages = MaxPagesList
			if listIncludeGroups {
				paramsPermissionList.IncludeGroups = flib.Bool(true)
			}

			client := permission.Client{Config: *config}
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

	cmdList.Flags().StringVar(&paramsPermissionList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsPermissionList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVar(&paramsPermissionList.Path, "path", "", "DEPRECATED: Permission path.  If provided, will scope permissions to this path. Use `filter[path]` instead.")
	cmdList.Flags().StringVar(&paramsPermissionList.GroupId, "group-id", "", "DEPRECATED: Group ID.  If provided, will scope permissions to this group. Use `filter[group_id]` instead.`")
	cmdList.Flags().StringVar(&paramsPermissionList.UserId, "user-id", "", "DEPRECATED: User ID.  If provided, will scope permissions to this user. Use `filter[user_id]` instead.`")
	cmdList.Flags().BoolVar(&listIncludeGroups, "include-groups", listIncludeGroups, "If searching by user or group, also include user's permissions that are inherited from its groups?")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Permissions.AddCommand(cmdList)
	var fieldsCreate string
	var formatCreate string
	createRecursive := false
	paramsPermissionCreate := files_sdk.PermissionCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := permission.Client{Config: *config}

			if createRecursive {
				paramsPermissionCreate.Recursive = flib.Bool(true)
			}

			if len(args) > 0 && args[0] != "" {
				paramsPermissionCreate.Path = args[0]
			}
			var result interface{}
			var err error
			result, err = client.Create(ctx, paramsPermissionCreate)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatCreate, fieldsCreate, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdCreate.Flags().Int64Var(&paramsPermissionCreate.GroupId, "group-id", 0, "Group ID")
	cmdCreate.Flags().StringVar(&paramsPermissionCreate.Path, "path", "", "Folder path")
	cmdCreate.Flags().StringVar(&paramsPermissionCreate.Permission, "permission", "", " Permission type.  Can be `admin`, `full`, `readonly`, `writeonly`, `list`, or `history`")
	cmdCreate.Flags().BoolVar(&createRecursive, "recursive", createRecursive, "Apply to subfolders recursively?")
	cmdCreate.Flags().Int64Var(&paramsPermissionCreate.UserId, "user-id", 0, "User ID.  Provide `username` or `user_id`")
	cmdCreate.Flags().StringVar(&paramsPermissionCreate.Username, "username", "", "User username.  Provide `username` or `user_id`")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Permissions.AddCommand(cmdCreate)
	var fieldsDelete string
	var formatDelete string
	paramsPermissionDelete := files_sdk.PermissionDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := permission.Client{Config: *config}

			var result interface{}
			var err error
			result, err = client.Delete(ctx, paramsPermissionDelete)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatDelete, fieldsDelete, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdDelete.Flags().Int64Var(&paramsPermissionDelete.Id, "id", 0, "Permission ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Permissions.AddCommand(cmdDelete)
}
