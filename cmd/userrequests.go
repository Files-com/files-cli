package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	user_request "github.com/Files-com/files-sdk-go/userrequest"
)

var (
	UserRequests = &cobra.Command{
		Use:  "user-requests [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func UserRequestsInit() {
	var fieldsList string
	paramsUserRequestList := files_sdk.UserRequestListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsUserRequestList
			params.MaxPages = MaxPagesList
			client := user_request.Client{Config: files_sdk.GlobalConfig}
			it, err := client.List(params)
			if err != nil {
				lib.ClientError(err)
			}
			err = lib.JsonMarshalIter(it, fieldsList)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdList.Flags().StringVarP(&paramsUserRequestList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&paramsUserRequestList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	UserRequests.AddCommand(cmdList)
	var fieldsFind string
	paramsUserRequestFind := files_sdk.UserRequestFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			client := user_request.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Find(paramsUserRequestFind)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsFind)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsUserRequestFind.Id, "id", "i", 0, "User Request ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	UserRequests.AddCommand(cmdFind)
	var fieldsCreate string
	paramsUserRequestCreate := files_sdk.UserRequestCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			client := user_request.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Create(paramsUserRequestCreate)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdCreate.Flags().StringVarP(&paramsUserRequestCreate.Name, "name", "n", "", "Name of user requested")
	cmdCreate.Flags().StringVarP(&paramsUserRequestCreate.Email, "email", "e", "", "Email of user requested")
	cmdCreate.Flags().StringVarP(&paramsUserRequestCreate.Details, "details", "d", "", "Details of the user request")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	UserRequests.AddCommand(cmdCreate)
	var fieldsDelete string
	paramsUserRequestDelete := files_sdk.UserRequestDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			client := user_request.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Delete(paramsUserRequestDelete)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsUserRequestDelete.Id, "id", "i", 0, "User Request ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	UserRequests.AddCommand(cmdDelete)
}
