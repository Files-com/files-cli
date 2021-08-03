package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/group"
)

var (
	Groups = &cobra.Command{}
)

func GroupsInit() {
	Groups = &cobra.Command{
		Use:  "groups [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command groups\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsGroupList := files_sdk.GroupListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsGroupList
			params.MaxPages = MaxPagesList

			client := group.Client{Config: *config}
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
	cmdList.Flags().StringVarP(&paramsGroupList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsGroupList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsGroupList.Ids, "ids", "i", "", "Comma-separated list of group ids to include in results.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-light")
	Groups.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	paramsGroupFind := files_sdk.GroupFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := group.Client{Config: *config}

			result, err := client.Find(ctx, paramsGroupFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatFind, fieldsFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsGroupFind.Id, "id", "i", 0, "Group ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-light")
	Groups.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	paramsGroupCreate := files_sdk.GroupCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := group.Client{Config: *config}

			result, err := client.Create(ctx, paramsGroupCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatCreate, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().StringVarP(&paramsGroupCreate.Name, "name", "n", "", "Group name.")
	cmdCreate.Flags().StringVarP(&paramsGroupCreate.Notes, "notes", "o", "", "Group notes.")
	cmdCreate.Flags().StringVarP(&paramsGroupCreate.UserIds, "user-ids", "u", "", "A list of user ids. If sent as a string, should be comma-delimited.")
	cmdCreate.Flags().StringVarP(&paramsGroupCreate.AdminIds, "admin-ids", "a", "", "A list of group admin user ids. If sent as a string, should be comma-delimited.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-light")
	Groups.AddCommand(cmdCreate)
	var fieldsUpdate string
	var formatUpdate string
	paramsGroupUpdate := files_sdk.GroupUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := group.Client{Config: *config}

			result, err := client.Update(ctx, paramsGroupUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatUpdate, fieldsUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsGroupUpdate.Id, "id", "i", 0, "Group ID.")
	cmdUpdate.Flags().StringVarP(&paramsGroupUpdate.Name, "name", "n", "", "Group name.")
	cmdUpdate.Flags().StringVarP(&paramsGroupUpdate.Notes, "notes", "o", "", "Group notes.")
	cmdUpdate.Flags().StringVarP(&paramsGroupUpdate.UserIds, "user-ids", "u", "", "A list of user ids. If sent as a string, should be comma-delimited.")
	cmdUpdate.Flags().StringVarP(&paramsGroupUpdate.AdminIds, "admin-ids", "a", "", "A list of group admin user ids. If sent as a string, should be comma-delimited.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-light")
	Groups.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	paramsGroupDelete := files_sdk.GroupDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := group.Client{Config: *config}

			result, err := client.Delete(ctx, paramsGroupDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatDelete, fieldsDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsGroupDelete.Id, "id", "i", 0, "Group ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-light")
	Groups.AddCommand(cmdDelete)
}
