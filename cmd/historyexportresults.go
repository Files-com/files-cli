package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"
	"os"

	history_export_result "github.com/Files-com/files-sdk-go/historyexportresult"
)

var (
	HistoryExportResults = &cobra.Command{
		Use:  "history-export-results [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func HistoryExportResultsInit() {
	var fieldsList string
	paramsHistoryExportResultList := files_sdk.HistoryExportResultListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsHistoryExportResultList
			params.MaxPages = MaxPagesList
			it, err := history_export_result.List(params)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().Int64VarP(&paramsHistoryExportResultList.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().IntVarP(&paramsHistoryExportResultList.Page, "page", "p", 0, "Current page number.")
	cmdList.Flags().IntVarP(&paramsHistoryExportResultList.PerPage, "per-page", "e", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsHistoryExportResultList.Action, "action", "a", "", "Deprecated: If set to `count` returns a count of matching records rather than the records themselves.")
	cmdList.Flags().StringVarP(&paramsHistoryExportResultList.Cursor, "cursor", "c", "", "Send cursor to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsHistoryExportResultList.HistoryExportId, "history-export-id", "i", 0, "ID of the associated history export.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	HistoryExportResults.AddCommand(cmdList)
}
