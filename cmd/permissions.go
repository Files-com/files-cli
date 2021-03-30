package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/permission"
)

var (
	Permissions = &cobra.Command{
		Use:  "permissions [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func PermissionsInit() {
	var fieldsList string
	paramsPermissionList := files_sdk.PermissionListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsPermissionList
			params.MaxPages = MaxPagesList
			client := permission.Client{Config: files_sdk.GlobalConfig}
			it, err := client.List(params)
			if err != nil {
				lib.ClientError(err)
			}
			err = lib.JsonMarshalIter(it, fieldsList)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdList.Flags().StringVarP(&paramsPermissionList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&paramsPermissionList.PerPage, "per-page", "a", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsPermissionList.Path, "path", "p", "", "DEPRECATED: Permission path.  If provided, will scope permissions to this path. Use `filter[path]` instead.")
	cmdList.Flags().StringVarP(&paramsPermissionList.GroupId, "group-id", "r", "", "DEPRECATED: Group ID.  If provided, will scope permissions to this group. Use `filter[group_id]` instead.`")
	cmdList.Flags().StringVarP(&paramsPermissionList.UserId, "user-id", "u", "", "DEPRECATED: User ID.  If provided, will scope permissions to this user. Use `filter[user_id]` instead.`")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	Permissions.AddCommand(cmdList)
	var fieldsCreate string
	paramsPermissionCreate := files_sdk.PermissionCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			client := permission.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Create(paramsPermissionCreate)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdCreate.Flags().Int64VarP(&paramsPermissionCreate.GroupId, "group-id", "g", 0, "Group ID")
	cmdCreate.Flags().StringVarP(&paramsPermissionCreate.Path, "path", "p", "", "Folder path")
	cmdCreate.Flags().StringVarP(&paramsPermissionCreate.Permission, "permission", "e", "", " Permission type.  Can be `admin`, `full`, `readonly`, `writeonly`, `list`, or `history`")
	cmdCreate.Flags().Int64VarP(&paramsPermissionCreate.UserId, "user-id", "u", 0, "User ID.  Provide `username` or `user_id`")
	cmdCreate.Flags().StringVarP(&paramsPermissionCreate.Username, "username", "s", "", "User username.  Provide `username` or `user_id`")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	Permissions.AddCommand(cmdCreate)
	var fieldsDelete string
	paramsPermissionDelete := files_sdk.PermissionDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			client := permission.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Delete(paramsPermissionDelete)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsPermissionDelete.Id, "id", "i", 0, "Permission ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Permissions.AddCommand(cmdDelete)
}
