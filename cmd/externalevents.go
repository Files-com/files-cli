package cmd

import (
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	external_event "github.com/Files-com/files-sdk-go/v2/externalevent"
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
	usePagerList := true
	paramsExternalEventList := files_sdk.ExternalEventListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List External Events",
		Long:  `List External Events`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsExternalEventList
			params.MaxPages = MaxPagesList

			client := external_event.Client{Config: *config}
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
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}

	cmdList.Flags().StringVar(&paramsExternalEventList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsExternalEventList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	ExternalEvents.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	usePagerFind := true
	paramsExternalEventFind := files_sdk.ExternalEventFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show External Event`,
		Long:  `Show External Event`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := external_event.Client{Config: *config}

			var externalEvent interface{}
			var err error
			externalEvent, err = client.Find(ctx, paramsExternalEventFind)
			lib.HandleResponse(ctx, externalEvent, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdFind.Flags().Int64Var(&paramsExternalEventFind.Id, "id", 0, "External Event ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	ExternalEvents.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	usePagerCreate := true
	paramsExternalEventCreate := files_sdk.ExternalEventCreateParams{}
	ExternalEventCreateStatus := ""

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create External Event`,
		Long:  `Create External Event`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := external_event.Client{Config: *config}

			var externalEvent interface{}
			var err error
			paramsExternalEventCreate.Status = paramsExternalEventCreate.Status.Enum()[ExternalEventCreateStatus]
			externalEvent, err = client.Create(ctx, paramsExternalEventCreate)
			lib.HandleResponse(ctx, externalEvent, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdCreate.Flags().StringVar(&ExternalEventCreateStatus, "status", "", fmt.Sprintf("Status of event. %v", reflect.ValueOf(paramsExternalEventCreate.Status.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&paramsExternalEventCreate.Body, "body", "", "Event body")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	ExternalEvents.AddCommand(cmdCreate)
}
