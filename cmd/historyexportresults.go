package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/historyexportresult"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = history_export_result.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
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
			it := history_export_result.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsHistoryExportResultList.Page, "page", "p", 0, "List History Export Results")
	cmdList.Flags().IntVarP(&paramsHistoryExportResultList.PerPage, "per-page", "e", 0, "List History Export Results")
	cmdList.Flags().StringVarP(&paramsHistoryExportResultList.Action, "action", "a", "", "List History Export Results")
	cmdList.Flags().StringVarP(&paramsHistoryExportResultList.Cursor, "cursor", "c", "", "List History Export Results")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
	HistoryExportResults.AddCommand(cmdList)
}
