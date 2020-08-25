package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/project"
)

var (
	Projects = &cobra.Command{
		Use:  "projects [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func ProjectsInit() {
	var fieldsList string
	paramsProjectList := files_sdk.ProjectListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsProjectList
			params.MaxPages = MaxPagesList
			it, err := project.List(params)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsProjectList.Page, "page", "p", 0, "Current page number.")
	cmdList.Flags().IntVarP(&paramsProjectList.PerPage, "per-page", "e", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsProjectList.Action, "action", "a", "", "Deprecated: If set to `count` returns a count of matching records rather than the records themselves.")
	cmdList.Flags().StringVarP(&paramsProjectList.Cursor, "cursor", "c", "", "Send cursor to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	Projects.AddCommand(cmdList)
	var fieldsFind string
	paramsProjectFind := files_sdk.ProjectFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := project.Find(paramsProjectFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsFind)
		},
	}
	cmdFind.Flags().Int64VarP(&paramsProjectFind.Id, "id", "i", 0, "Project ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	Projects.AddCommand(cmdFind)
	var fieldsCreate string
	paramsProjectCreate := files_sdk.ProjectCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := project.Create(paramsProjectCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().StringVarP(&paramsProjectCreate.GlobalAccess, "global-access", "g", "", "Global permissions.  Can be: `none`, `anyone_with_read`, `anyone_with_full`.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	Projects.AddCommand(cmdCreate)
	var fieldsUpdate string
	paramsProjectUpdate := files_sdk.ProjectUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := project.Update(paramsProjectUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsUpdate)
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsProjectUpdate.Id, "id", "i", 0, "Project ID.")
	cmdUpdate.Flags().StringVarP(&paramsProjectUpdate.GlobalAccess, "global-access", "g", "", "Global permissions.  Can be: `none`, `anyone_with_read`, `anyone_with_full`.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	Projects.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsProjectDelete := files_sdk.ProjectDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := project.Delete(paramsProjectDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsProjectDelete.Id, "id", "i", 0, "Project ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Projects.AddCommand(cmdDelete)
}
