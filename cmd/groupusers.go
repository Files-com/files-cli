package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"
	"os"

	group_user "github.com/Files-com/files-sdk-go/groupuser"
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
			client := group_user.Client{Config: files_sdk.GlobalConfig}
			it, err := client.List(params)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = lib.JsonMarshalIter(it, fieldsList)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdList.Flags().Int64VarP(&paramsGroupUserList.UserId, "user-id", "u", 0, "User ID.  If provided, will return group_users of this user.")
	cmdList.Flags().StringVarP(&paramsGroupUserList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&paramsGroupUserList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64VarP(&paramsGroupUserList.GroupId, "group-id", "g", 0, "Group ID.  If provided, will return group_users of this group.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	GroupUsers.AddCommand(cmdList)
	var fieldsUpdate string
	paramsGroupUserUpdate := files_sdk.GroupUserUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			client := group_user.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Update(paramsGroupUserUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = lib.JsonMarshal(result, fieldsUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsGroupUserUpdate.Id, "id", "i", 0, "Group User ID.")
	cmdUpdate.Flags().Int64VarP(&paramsGroupUserUpdate.GroupId, "group-id", "g", 0, "Group ID to add user to.")
	cmdUpdate.Flags().Int64VarP(&paramsGroupUserUpdate.UserId, "user-id", "u", 0, "User ID to add to group.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	GroupUsers.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsGroupUserDelete := files_sdk.GroupUserDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			client := group_user.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Delete(paramsGroupUserDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsGroupUserDelete.Id, "id", "i", 0, "Group User ID.")
	cmdDelete.Flags().Int64VarP(&paramsGroupUserDelete.GroupId, "group-id", "g", 0, "Group ID from which to remove user.")
	cmdDelete.Flags().Int64VarP(&paramsGroupUserDelete.UserId, "user-id", "u", 0, "User ID to remove from group.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	GroupUsers.AddCommand(cmdDelete)
}
