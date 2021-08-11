package cmd

import (
	"reflect"

	"github.com/Files-com/files-cli/lib"
	external_event "github.com/Files-com/files-sdk-go/externalevent"
	"github.com/spf13/cobra"

	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go"
)

var (
	ExternalEvents = &cobra.Command{}
)

func ExternalEventsInit() {
	ExternalEvents = &cobra.Command{
		Use:  "external-events [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command external-events\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsExternalEventList := files_sdk.ExternalEventListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsExternalEventList
			params.MaxPages = MaxPagesList

			client := external_event.Client{Config: *config}
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
	cmdList.Flags().StringVarP(&paramsExternalEventList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsExternalEventList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	ExternalEvents.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	paramsExternalEventFind := files_sdk.ExternalEventFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := external_event.Client{Config: *config}

			result, err := client.Find(ctx, paramsExternalEventFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatFind, fieldsFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsExternalEventFind.Id, "id", "i", 0, "External Event ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright")
	ExternalEvents.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	paramsExternalEventCreate := files_sdk.ExternalEventCreateParams{}
	ExternalEventCreateStatus := ""

	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := external_event.Client{Config: *config}

			paramsExternalEventCreate.Status = paramsExternalEventCreate.Status.Enum()[ExternalEventCreateStatus]

			result, err := client.Create(ctx, paramsExternalEventCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatCreate, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().StringVarP(&ExternalEventCreateStatus, "status", "s", "", fmt.Sprintf("Status of event. %v", reflect.ValueOf(paramsExternalEventCreate.Status.Enum()).MapKeys()))
	cmdCreate.Flags().StringVarP(&paramsExternalEventCreate.Body, "body", "b", "", "Event body")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	ExternalEvents.AddCommand(cmdCreate)
}
