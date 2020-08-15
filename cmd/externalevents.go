package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/externalevent"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = external_event.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
)

var (
	ExternalEvents = &cobra.Command{
		Use:  "external-events [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func ExternalEventsInit() {
	var fieldsList string
	paramsExternalEventList := files_sdk.ExternalEventListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsExternalEventList
			params.MaxPages = MaxPagesList
			it := external_event.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsExternalEventList.Page, "page", "p", 0, "List External Events")
	cmdList.Flags().IntVarP(&paramsExternalEventList.PerPage, "per-page", "r", 0, "List External Events")
	cmdList.Flags().StringVarP(&paramsExternalEventList.Action, "action", "a", "", "List External Events")
	cmdList.Flags().StringVarP(&paramsExternalEventList.Cursor, "cursor", "c", "", "List External Events")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
	ExternalEvents.AddCommand(cmdList)
}
