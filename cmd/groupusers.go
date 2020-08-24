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
			it := group_user.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsGroupUserList.Page, "page", "p", 0, "Current page number.")
	cmdList.Flags().IntVarP(&paramsGroupUserList.PerPage, "per-page", "e", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsGroupUserList.Action, "action", "a", "", "Deprecated: If set to `count` returns a count of matching records rather than the records themselves.")
	cmdList.Flags().StringVarP(&paramsGroupUserList.Cursor, "cursor", "c", "", "Send cursor to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
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

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
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

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	GroupUsers.AddCommand(cmdDelete)
}
