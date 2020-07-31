package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/groupuser"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = group_user.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
)

var (
	GroupUsers = &cobra.Command{
		Use:  "group-users [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func GroupUsersInit() {
	var fieldsList string
	paramsGroupUserList := files_sdk.GroupUserListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsGroupUserList
			params.MaxPages = MaxPagesList
			it := group_user.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsGroupUserList.UserId, "user-id", "u", 0, "List Group Users")
	cmdList.Flags().IntVarP(&paramsGroupUserList.Page, "page", "p", 0, "List Group Users")
	cmdList.Flags().IntVarP(&paramsGroupUserList.PerPage, "per-page", "e", 0, "List Group Users")
	cmdList.Flags().StringVarP(&paramsGroupUserList.Action, "action", "a", "", "List Group Users")
	cmdList.Flags().IntVarP(&paramsGroupUserList.GroupId, "group-id", "g", 0, "List Group Users")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
	GroupUsers.AddCommand(cmdList)
	var fieldsUpdate string
	paramsGroupUserUpdate := files_sdk.GroupUserUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := group_user.Update(paramsGroupUserUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsUpdate)
		},
	}
	cmdUpdate.Flags().IntVarP(&paramsGroupUserUpdate.Id, "id", "i", 0, "Update Group User")
	cmdUpdate.Flags().IntVarP(&paramsGroupUserUpdate.GroupId, "group-id", "g", 0, "Update Group User")
	cmdUpdate.Flags().IntVarP(&paramsGroupUserUpdate.UserId, "user-id", "u", 0, "Update Group User")
	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "f", "", "comma separated list of field names")
	GroupUsers.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsGroupUserDelete := files_sdk.GroupUserDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := group_user.Delete(paramsGroupUserDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}
	cmdDelete.Flags().IntVarP(&paramsGroupUserDelete.Id, "id", "i", 0, "Delete Group User")
	cmdDelete.Flags().IntVarP(&paramsGroupUserDelete.GroupId, "group-id", "g", 0, "Delete Group User")
	cmdDelete.Flags().IntVarP(&paramsGroupUserDelete.UserId, "user-id", "u", 0, "Delete Group User")
	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
	GroupUsers.AddCommand(cmdDelete)
}
