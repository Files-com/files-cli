package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"
	"os"

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
			it := permission.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsPermissionList.Page, "page", "p", 0, "Current page number.")
	cmdList.Flags().IntVarP(&paramsPermissionList.PerPage, "per-page", "", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsPermissionList.Action, "action", "a", "", "Deprecated: If set to `count` returns a count of matching records rather than the records themselves.")
	cmdList.Flags().StringVarP(&paramsPermissionList.Cursor, "cursor", "c", "", "Send cursor to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().StringVarP(&paramsPermissionList.Path, "path", "", "", "DEPRECATED: Permission path.  If provided, will scope permissions to this path. Use `filter[path]` instead.")
	cmdList.Flags().StringVarP(&paramsPermissionList.GroupId, "group-id", "r", "", "DEPRECATED: Group ID.  If provided, will scope permissions to this group. Use `filter[group_id]` instead.`")
	cmdList.Flags().StringVarP(&paramsPermissionList.UserId, "user-id", "u", "", "DEPRECATED: User ID.  If provided, will scope permissions to this user. Use `filter[user_id]` instead.`")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	Permissions.AddCommand(cmdList)
	var fieldsCreate string
	paramsPermissionCreate := files_sdk.PermissionCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := permission.Create(paramsPermissionCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().StringVarP(&paramsPermissionCreate.Path, "path", "p", "", "Folder path")
	cmdCreate.Flags().StringVarP(&paramsPermissionCreate.Permission, "permission", "e", "", " Permission type.  Can be `admin`, `full`, `readonly`, `writeonly`, `list`, or `history`")
	cmdCreate.Flags().StringVarP(&paramsPermissionCreate.Username, "username", "s", "", "User username.  Provide `username` or `user_id`")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	Permissions.AddCommand(cmdCreate)
	var fieldsDelete string
	paramsPermissionDelete := files_sdk.PermissionDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := permission.Delete(paramsPermissionDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Permissions.AddCommand(cmdDelete)
}
