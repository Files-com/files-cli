package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/clickwrap"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = clickwrap.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
)

var (
	Clickwraps = &cobra.Command{
		Use:  "clickwraps [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func ClickwrapsInit() {
	var fieldsList string
	paramsClickwrapList := files_sdk.ClickwrapListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsClickwrapList
			params.MaxPages = MaxPagesList
			it := clickwrap.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsClickwrapList.Page, "page", "p", 0, "List Clickwraps")
	cmdList.Flags().IntVarP(&paramsClickwrapList.PerPage, "per-page", "e", 0, "List Clickwraps")
	cmdList.Flags().StringVarP(&paramsClickwrapList.Action, "action", "a", "", "List Clickwraps")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
	Clickwraps.AddCommand(cmdList)
	var fieldsFind string
	paramsClickwrapFind := files_sdk.ClickwrapFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := clickwrap.Find(paramsClickwrapFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsFind)
		},
	}
	cmdFind.Flags().IntVarP(&paramsClickwrapFind.Id, "id", "i", 0, "Show Clickwrap")
	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "f", "", "comma separated list of field names")
	Clickwraps.AddCommand(cmdFind)
	var fieldsCreate string
	paramsClickwrapCreate := files_sdk.ClickwrapCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := clickwrap.Create(paramsClickwrapCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().StringVarP(&paramsClickwrapCreate.Name, "name", "n", "", "Create Clickwrap")
	cmdCreate.Flags().StringVarP(&paramsClickwrapCreate.Body, "body", "b", "", "Create Clickwrap")
	cmdCreate.Flags().StringVarP(&paramsClickwrapCreate.UseWithBundles, "use-with-bundles", "u", "", "Create Clickwrap")
	cmdCreate.Flags().StringVarP(&paramsClickwrapCreate.UseWithInboxes, "use-with-inboxes", "i", "", "Create Clickwrap")
	cmdCreate.Flags().StringVarP(&paramsClickwrapCreate.UseWithUsers, "use-with-users", "s", "", "Create Clickwrap")
	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
	Clickwraps.AddCommand(cmdCreate)
	var fieldsUpdate string
	paramsClickwrapUpdate := files_sdk.ClickwrapUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := clickwrap.Update(paramsClickwrapUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsUpdate)
		},
	}
	cmdUpdate.Flags().IntVarP(&paramsClickwrapUpdate.Id, "id", "i", 0, "Update Clickwrap")
	cmdUpdate.Flags().StringVarP(&paramsClickwrapUpdate.Name, "name", "n", "", "Update Clickwrap")
	cmdUpdate.Flags().StringVarP(&paramsClickwrapUpdate.Body, "body", "b", "", "Update Clickwrap")
	cmdUpdate.Flags().StringVarP(&paramsClickwrapUpdate.UseWithBundles, "use-with-bundles", "u", "", "Update Clickwrap")
	cmdUpdate.Flags().StringVarP(&paramsClickwrapUpdate.UseWithInboxes, "use-with-inboxes", "o", "", "Update Clickwrap")
	cmdUpdate.Flags().StringVarP(&paramsClickwrapUpdate.UseWithUsers, "use-with-users", "s", "", "Update Clickwrap")
	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "f", "", "comma separated list of field names")
	Clickwraps.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsClickwrapDelete := files_sdk.ClickwrapDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := clickwrap.Delete(paramsClickwrapDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}
	cmdDelete.Flags().IntVarP(&paramsClickwrapDelete.Id, "id", "i", 0, "Delete Clickwrap")
	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
	Clickwraps.AddCommand(cmdDelete)
}
