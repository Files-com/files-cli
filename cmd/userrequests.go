package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	user_request "github.com/Files-com/files-sdk-go/v2/userrequest"
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
		Short: "List User Requests",
		Long:  `List User Requests`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsUserRequestList
			params.MaxPages = MaxPagesList

			client := user_request.Client{Config: *config}
			it, err := client.List(ctx, params)
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, formatList, config.Logger())
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(it, formatList, fieldsList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}

	cmdList.Flags().StringVar(&paramsUserRequestList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsUserRequestList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	UserRequests.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	paramsUserRequestFind := files_sdk.UserRequestFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show User Request`,
		Long:  `Show User Request`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := user_request.Client{Config: *config}

			var userRequest interface{}
			var err error
			userRequest, err = client.Find(ctx, paramsUserRequestFind)
			lib.HandleResponse(ctx, userRequest, err, formatFind, fieldsFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdFind.Flags().Int64Var(&paramsUserRequestFind.Id, "id", 0, "User Request ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright")
	UserRequests.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	paramsUserRequestCreate := files_sdk.UserRequestCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create User Request`,
		Long:  `Create User Request`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := user_request.Client{Config: *config}

			var userRequest interface{}
			var err error
			userRequest, err = client.Create(ctx, paramsUserRequestCreate)
			lib.HandleResponse(ctx, userRequest, err, formatCreate, fieldsCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdCreate.Flags().StringVar(&paramsUserRequestCreate.Name, "name", "", "Name of user requested")
	cmdCreate.Flags().StringVar(&paramsUserRequestCreate.Email, "email", "", "Email of user requested")
	cmdCreate.Flags().StringVar(&paramsUserRequestCreate.Details, "details", "", "Details of the user request")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	UserRequests.AddCommand(cmdCreate)
	var fieldsDelete string
	var formatDelete string
	paramsUserRequestDelete := files_sdk.UserRequestDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete User Request`,
		Long:  `Delete User Request`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := user_request.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsUserRequestDelete)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}
	cmdDelete.Flags().Int64Var(&paramsUserRequestDelete.Id, "id", 0, "User Request ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright")
	UserRequests.AddCommand(cmdDelete)
}
