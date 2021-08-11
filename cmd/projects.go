package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/project"
)

var (
	Projects = &cobra.Command{}
)

func ProjectsInit() {
	Projects = &cobra.Command{
		Use:  "projects [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command projects\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsProjectList := files_sdk.ProjectListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsProjectList
			params.MaxPages = MaxPagesList

			client := project.Client{Config: *config}
			it, err := client.List(ctx, params)
			if err != nil {
				lib.ClientError(ctx, err)
			}
			err = lib.FormatIter(it, formatList, fieldsList)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdList.Flags().StringVarP(&paramsProjectList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsProjectList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Projects.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	paramsProjectFind := files_sdk.ProjectFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := project.Client{Config: *config}

			result, err := client.Find(ctx, paramsProjectFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatFind, fieldsFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsProjectFind.Id, "id", "i", 0, "Project ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Projects.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	paramsProjectCreate := files_sdk.ProjectCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := project.Client{Config: *config}

			result, err := client.Create(ctx, paramsProjectCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatCreate, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().StringVarP(&paramsProjectCreate.GlobalAccess, "global-access", "g", "", "Global permissions.  Can be: `none`, `anyone_with_read`, `anyone_with_full`.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Projects.AddCommand(cmdCreate)
	var fieldsUpdate string
	var formatUpdate string
	paramsProjectUpdate := files_sdk.ProjectUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := project.Client{Config: *config}

			result, err := client.Update(ctx, paramsProjectUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatUpdate, fieldsUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsProjectUpdate.Id, "id", "i", 0, "Project ID.")
	cmdUpdate.Flags().StringVarP(&paramsProjectUpdate.GlobalAccess, "global-access", "g", "", "Global permissions.  Can be: `none`, `anyone_with_read`, `anyone_with_full`.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Projects.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	paramsProjectDelete := files_sdk.ProjectDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := project.Client{Config: *config}

			result, err := client.Delete(ctx, paramsProjectDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatDelete, fieldsDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsProjectDelete.Id, "id", "i", 0, "Project ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Projects.AddCommand(cmdDelete)
}
