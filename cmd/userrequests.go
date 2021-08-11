package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"

	user_request "github.com/Files-com/files-sdk-go/userrequest"
)

var (
	UserRequests = &cobra.Command{}
)

func UserRequestsInit() {
	UserRequests = &cobra.Command{
		Use:  "user-requests [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command user-requests\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsUserRequestList := files_sdk.UserRequestListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsUserRequestList
			params.MaxPages = MaxPagesList

			client := user_request.Client{Config: *config}
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
	cmdList.Flags().StringVarP(&paramsUserRequestList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsUserRequestList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	UserRequests.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	paramsUserRequestFind := files_sdk.UserRequestFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := user_request.Client{Config: *config}

			result, err := client.Find(ctx, paramsUserRequestFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatFind, fieldsFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsUserRequestFind.Id, "id", "i", 0, "User Request ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright")
	UserRequests.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	paramsUserRequestCreate := files_sdk.UserRequestCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := user_request.Client{Config: *config}

			result, err := client.Create(ctx, paramsUserRequestCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatCreate, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().StringVarP(&paramsUserRequestCreate.Name, "name", "n", "", "Name of user requested")
	cmdCreate.Flags().StringVarP(&paramsUserRequestCreate.Email, "email", "e", "", "Email of user requested")
	cmdCreate.Flags().StringVarP(&paramsUserRequestCreate.Details, "details", "d", "", "Details of the user request")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	UserRequests.AddCommand(cmdCreate)
	var fieldsDelete string
	var formatDelete string
	paramsUserRequestDelete := files_sdk.UserRequestDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := user_request.Client{Config: *config}

			result, err := client.Delete(ctx, paramsUserRequestDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatDelete, fieldsDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsUserRequestDelete.Id, "id", "i", 0, "User Request ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright")
	UserRequests.AddCommand(cmdDelete)
}
