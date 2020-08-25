package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/request"
)

var (
	Requests = &cobra.Command{
		Use:  "requests [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func RequestsInit() {
	var fieldsList string
	paramsRequestList := files_sdk.RequestListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsRequestList
			params.MaxPages = MaxPagesList
			it, err := request.List(params)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsRequestList.Page, "page", "p", 0, "Current page number.")
	cmdList.Flags().IntVarP(&paramsRequestList.PerPage, "per-page", "e", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsRequestList.Action, "action", "a", "", "Deprecated: If set to `count` returns a count of matching records rather than the records themselves.")
	cmdList.Flags().StringVarP(&paramsRequestList.Cursor, "cursor", "c", "", "Send cursor to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().StringVarP(&paramsRequestList.Path, "path", "t", "", "Path to show requests for.  If omitted, shows all paths. Send `/` to represent the root directory.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	Requests.AddCommand(cmdList)
	var fieldsGetFolder string
	paramsRequestGetFolder := files_sdk.RequestGetFolderParams{}
	cmdGetFolder := &cobra.Command{
		Use: "get-folder",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := request.GetFolder(paramsRequestGetFolder)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsGetFolder)
		},
	}
	cmdGetFolder.Flags().IntVarP(&paramsRequestGetFolder.Page, "page", "p", 0, "Current page number.")
	cmdGetFolder.Flags().IntVarP(&paramsRequestGetFolder.PerPage, "per-page", "e", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdGetFolder.Flags().StringVarP(&paramsRequestGetFolder.Action, "action", "a", "", "Deprecated: If set to `count` returns a count of matching records rather than the records themselves.")
	cmdGetFolder.Flags().StringVarP(&paramsRequestGetFolder.Cursor, "cursor", "c", "", "Send cursor to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdGetFolder.Flags().StringVarP(&paramsRequestGetFolder.Path, "path", "t", "", "Path to show requests for.  If omitted, shows all paths. Send `/` to represent the root directory.")

	cmdGetFolder.Flags().StringVarP(&fieldsGetFolder, "fields", "", "", "comma separated list of field names")
	Requests.AddCommand(cmdGetFolder)
	var fieldsCreate string
	paramsRequestCreate := files_sdk.RequestCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := request.Create(paramsRequestCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().StringVarP(&paramsRequestCreate.Path, "path", "p", "", "Folder path on which to request the file.")
	cmdCreate.Flags().StringVarP(&paramsRequestCreate.Destination, "destination", "d", "", "Destination filename (without extension) to request.")
	cmdCreate.Flags().StringVarP(&paramsRequestCreate.UserIds, "user-ids", "u", "", "A list of user IDs to request the file from. If sent as a string, it should be comma-delimited.")
	cmdCreate.Flags().StringVarP(&paramsRequestCreate.GroupIds, "group-ids", "g", "", "A list of group IDs to request the file from. If sent as a string, it should be comma-delimited.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	Requests.AddCommand(cmdCreate)
	var fieldsDelete string
	paramsRequestDelete := files_sdk.RequestDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := request.Delete(paramsRequestDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsRequestDelete.Id, "id", "i", 0, "Request ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Requests.AddCommand(cmdDelete)
}
