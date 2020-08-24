package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"
	"os"

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
			it := user_request.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsUserRequestList.Page, "page", "p", 0, "Current page number.")
	cmdList.Flags().IntVarP(&paramsUserRequestList.PerPage, "per-page", "e", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsUserRequestList.Action, "action", "a", "", "Deprecated: If set to `count` returns a count of matching records rather than the records themselves.")
	cmdList.Flags().StringVarP(&paramsUserRequestList.Cursor, "cursor", "c", "", "Send cursor to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	UserRequests.AddCommand(cmdList)
	var fieldsFind string
	paramsUserRequestFind := files_sdk.UserRequestFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := user_request.Find(paramsUserRequestFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsFind)
		},
	}

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	UserRequests.AddCommand(cmdFind)
	var fieldsCreate string
	paramsUserRequestCreate := files_sdk.UserRequestCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := user_request.Create(paramsUserRequestCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
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
			result, err := user_request.Delete(paramsUserRequestDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	UserRequests.AddCommand(cmdDelete)
}
