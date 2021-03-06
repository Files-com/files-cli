package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	flib "github.com/Files-com/files-sdk-go/lib"
	"github.com/Files-com/files-sdk-go/permission"
)

var (
	Permissions = &cobra.Command{}
)

func PermissionsInit() {
	Permissions = &cobra.Command{
		Use:  "permissions [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
	var fieldsList string
	paramsPermissionList := files_sdk.PermissionListParams{}
	var MaxPagesList int64
	listIncludeGroups := false

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			params := paramsPermissionList
			params.MaxPages = MaxPagesList
			if listIncludeGroups {
				paramsPermissionList.IncludeGroups = flib.Bool(true)
			}

			client := permission.Client{Config: *ctx.GetConfig()}
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
	cmdList.Flags().StringVarP(&paramsPermissionList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsPermissionList.PerPage, "per-page", "a", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsPermissionList.Path, "path", "p", "", "DEPRECATED: Permission path.  If provided, will scope permissions to this path. Use `filter[path]` instead.")
	cmdList.Flags().StringVarP(&paramsPermissionList.GroupId, "group-id", "r", "", "DEPRECATED: Group ID.  If provided, will scope permissions to this group. Use `filter[group_id]` instead.`")
	cmdList.Flags().StringVarP(&paramsPermissionList.UserId, "user-id", "u", "", "DEPRECATED: User ID.  If provided, will scope permissions to this user. Use `filter[user_id]` instead.`")
	cmdList.Flags().BoolVarP(&listIncludeGroups, "include-groups", "i", listIncludeGroups, "If searching by user or group, also include user's permissions that are inherited from its groups?")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	Permissions.AddCommand(cmdList)
	var fieldsCreate string
	createRecursive := false
	paramsPermissionCreate := files_sdk.PermissionCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := permission.Client{Config: *ctx.GetConfig()}

			if createRecursive {
				paramsPermissionCreate.Recursive = flib.Bool(true)
			}

			if len(args) > 0 && args[0] != "" {
				paramsPermissionCreate.Path = args[0]
			}

			result, err := client.Create(paramsPermissionCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdCreate.Flags().Int64VarP(&paramsPermissionCreate.GroupId, "group-id", "g", 0, "Group ID")
	cmdCreate.Flags().StringVarP(&paramsPermissionCreate.Path, "path", "p", "", "Folder path")
	cmdCreate.Flags().StringVarP(&paramsPermissionCreate.Permission, "permission", "e", "", " Permission type.  Can be `admin`, `full`, `readonly`, `writeonly`, `list`, or `history`")
	cmdCreate.Flags().BoolVarP(&createRecursive, "recursive", "r", createRecursive, "Apply to subfolders recursively?")
	cmdCreate.Flags().Int64VarP(&paramsPermissionCreate.UserId, "user-id", "u", 0, "User ID.  Provide `username` or `user_id`")
	cmdCreate.Flags().StringVarP(&paramsPermissionCreate.Username, "username", "s", "", "User username.  Provide `username` or `user_id`")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	Permissions.AddCommand(cmdCreate)
	var fieldsDelete string
	paramsPermissionDelete := files_sdk.PermissionDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := permission.Client{Config: *ctx.GetConfig()}

			result, err := client.Delete(paramsPermissionDelete)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsPermissionDelete.Id, "id", "i", 0, "Permission ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Permissions.AddCommand(cmdDelete)
}
