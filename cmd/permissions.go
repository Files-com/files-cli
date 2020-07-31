package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/permission"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = permission.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
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
	cmdList.Flags().IntVarP(&paramsPermissionList.Page, "page", "p", 0, "List Permissions")
	cmdList.Flags().IntVarP(&paramsPermissionList.PerPage, "per-page", "", 0, "List Permissions")
	cmdList.Flags().StringVarP(&paramsPermissionList.Action, "action", "a", "", "List Permissions")
	cmdList.Flags().StringVarP(&paramsPermissionList.Cursor, "cursor", "c", "", "List Permissions")
	cmdList.Flags().StringVarP(&paramsPermissionList.Path, "path", "", "", "List Permissions")
	cmdList.Flags().StringVarP(&paramsPermissionList.GroupId, "group-id", "r", "", "List Permissions")
	cmdList.Flags().StringVarP(&paramsPermissionList.UserId, "user-id", "u", "", "List Permissions")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
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
	cmdCreate.Flags().StringVarP(&paramsPermissionCreate.Path, "path", "p", "", "Create Permission")
	cmdCreate.Flags().StringVarP(&paramsPermissionCreate.Permission, "permission", "e", "", "Create Permission")
	cmdCreate.Flags().StringVarP(&paramsPermissionCreate.Username, "username", "s", "", "Create Permission")
	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
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
	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
	Permissions.AddCommand(cmdDelete)
}
