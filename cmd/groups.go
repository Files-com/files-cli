package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/group"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = group.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
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
			it := group.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsGroupList.Page, "page", "p", 0, "List Groups")
	cmdList.Flags().IntVarP(&paramsGroupList.PerPage, "per-page", "r", 0, "List Groups")
	cmdList.Flags().StringVarP(&paramsGroupList.Action, "action", "a", "", "List Groups")
	cmdList.Flags().StringVarP(&paramsGroupList.Cursor, "cursor", "c", "", "List Groups")
	cmdList.Flags().StringVarP(&paramsGroupList.Ids, "ids", "d", "", "List Groups")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
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
	cmdFind.Flags().IntVarP(&paramsGroupFind.Id, "id", "i", 0, "Show Group")
	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "f", "", "comma separated list of field names")
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
	cmdCreate.Flags().StringVarP(&paramsGroupCreate.Name, "name", "n", "", "Create Group")
	cmdCreate.Flags().StringVarP(&paramsGroupCreate.Notes, "notes", "o", "", "Create Group")
	cmdCreate.Flags().StringVarP(&paramsGroupCreate.UserIds, "user-ids", "u", "", "Create Group")
	cmdCreate.Flags().StringVarP(&paramsGroupCreate.AdminIds, "admin-ids", "a", "", "Create Group")
	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
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
	cmdUpdate.Flags().IntVarP(&paramsGroupUpdate.Id, "id", "i", 0, "Update Group")
	cmdUpdate.Flags().StringVarP(&paramsGroupUpdate.Name, "name", "n", "", "Update Group")
	cmdUpdate.Flags().StringVarP(&paramsGroupUpdate.Notes, "notes", "o", "", "Update Group")
	cmdUpdate.Flags().StringVarP(&paramsGroupUpdate.UserIds, "user-ids", "u", "", "Update Group")
	cmdUpdate.Flags().StringVarP(&paramsGroupUpdate.AdminIds, "admin-ids", "a", "", "Update Group")
	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "f", "", "comma separated list of field names")
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
	cmdDelete.Flags().IntVarP(&paramsGroupDelete.Id, "id", "i", 0, "Delete Group")
	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
	Groups.AddCommand(cmdDelete)
}
