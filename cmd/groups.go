package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/group"
)

var (
	Groups = &cobra.Command{
		Use:  "groups [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func GroupsInit() {
	var fieldsList string
	paramsGroupList := files_sdk.GroupListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsGroupList
			params.MaxPages = MaxPagesList
			it, err := group.List(params)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().StringVarP(&paramsGroupList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&paramsGroupList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsGroupList.Ids, "ids", "i", "", "Comma-separated list of group ids to include in results.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	Groups.AddCommand(cmdList)
	var fieldsFind string
	paramsGroupFind := files_sdk.GroupFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := group.Find(paramsGroupFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsFind)
		},
	}
	cmdFind.Flags().Int64VarP(&paramsGroupFind.Id, "id", "i", 0, "Group ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	Groups.AddCommand(cmdFind)
	var fieldsCreate string
	paramsGroupCreate := files_sdk.GroupCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := group.Create(paramsGroupCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().StringVarP(&paramsGroupCreate.Name, "name", "n", "", "Group name.")
	cmdCreate.Flags().StringVarP(&paramsGroupCreate.Notes, "notes", "o", "", "Group notes.")
	cmdCreate.Flags().StringVarP(&paramsGroupCreate.UserIds, "user-ids", "u", "", "A list of user ids. If sent as a string, should be comma-delimited.")
	cmdCreate.Flags().StringVarP(&paramsGroupCreate.AdminIds, "admin-ids", "a", "", "A list of group admin user ids. If sent as a string, should be comma-delimited.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	Groups.AddCommand(cmdCreate)
	var fieldsUpdate string
	paramsGroupUpdate := files_sdk.GroupUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := group.Update(paramsGroupUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsUpdate)
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsGroupUpdate.Id, "id", "i", 0, "Group ID.")
	cmdUpdate.Flags().StringVarP(&paramsGroupUpdate.Name, "name", "n", "", "Group name.")
	cmdUpdate.Flags().StringVarP(&paramsGroupUpdate.Notes, "notes", "o", "", "Group notes.")
	cmdUpdate.Flags().StringVarP(&paramsGroupUpdate.UserIds, "user-ids", "u", "", "A list of user ids. If sent as a string, should be comma-delimited.")
	cmdUpdate.Flags().StringVarP(&paramsGroupUpdate.AdminIds, "admin-ids", "a", "", "A list of group admin user ids. If sent as a string, should be comma-delimited.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	Groups.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsGroupDelete := files_sdk.GroupDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := group.Delete(paramsGroupDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsGroupDelete.Id, "id", "i", 0, "Group ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Groups.AddCommand(cmdDelete)
}
