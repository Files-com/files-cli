package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/userrequest"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = user_request.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
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
	cmdList.Flags().IntVarP(&paramsUserRequestList.Page, "page", "p", 0, "List User Requests")
	cmdList.Flags().IntVarP(&paramsUserRequestList.PerPage, "per-page", "e", 0, "List User Requests")
	cmdList.Flags().StringVarP(&paramsUserRequestList.Action, "action", "a", "", "List User Requests")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
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
	cmdFind.Flags().IntVarP(&paramsUserRequestFind.Id, "id", "i", 0, "Show User Request")
	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "f", "", "comma separated list of field names")
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
	cmdCreate.Flags().StringVarP(&paramsUserRequestCreate.Name, "name", "n", "", "Create User Request")
	cmdCreate.Flags().StringVarP(&paramsUserRequestCreate.Email, "email", "e", "", "Create User Request")
	cmdCreate.Flags().StringVarP(&paramsUserRequestCreate.Details, "details", "d", "", "Create User Request")
	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
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
	cmdDelete.Flags().IntVarP(&paramsUserRequestDelete.Id, "id", "i", 0, "Delete User Request")
	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
	UserRequests.AddCommand(cmdDelete)
}
