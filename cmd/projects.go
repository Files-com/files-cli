package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/project"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = project.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
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
			it := project.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsProjectList.Page, "page", "p", 0, "List Projects")
	cmdList.Flags().IntVarP(&paramsProjectList.PerPage, "per-page", "e", 0, "List Projects")
	cmdList.Flags().StringVarP(&paramsProjectList.Action, "action", "a", "", "List Projects")
	cmdList.Flags().StringVarP(&paramsProjectList.Cursor, "cursor", "c", "", "List Projects")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
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
	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "f", "", "comma separated list of field names")
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
	cmdCreate.Flags().StringVarP(&paramsProjectCreate.GlobalAccess, "global-access", "g", "", "Create Project")
	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
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
	cmdUpdate.Flags().StringVarP(&paramsProjectUpdate.GlobalAccess, "global-access", "g", "", "Update Project")
	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "f", "", "comma separated list of field names")
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
	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
	Projects.AddCommand(cmdDelete)
}
