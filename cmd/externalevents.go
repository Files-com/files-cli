package cmd

import (
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	external_event "github.com/Files-com/files-sdk-go/v2/externalevent"
)

func init() {
	RootCmd.AddCommand(ExternalEvents())
}

func ExternalEvents() *cobra.Command {
	ExternalEvents := &cobra.Command{
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
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsExternalEventList
			params.MaxPages = MaxPagesList

			client := external_event.Client{Config: *config}
			it, err := client.List(ctx, params)
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger())
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}

	cmdList.Flags().StringVar(&paramsExternalEventList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsExternalEventList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVar(&fieldsList, "fields", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVar(&formatList, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
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
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := external_event.Client{Config: *config}

			var externalEvent interface{}
			var err error
			externalEvent, err = client.Find(ctx, paramsExternalEventFind)
			lib.HandleResponse(ctx, Profile(cmd), externalEvent, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdFind.Flags().Int64Var(&paramsExternalEventFind.Id, "id", 0, "External Event ID.")

	cmdFind.Flags().StringVar(&fieldsFind, "fields", "", "comma separated list of field names")
	cmdFind.Flags().StringVar(&formatFind, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
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
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := external_event.Client{Config: *config}

			var externalEvent interface{}
			var err error
			var ExternalEventCreateStatusOk bool
			paramsExternalEventCreate.Status, ExternalEventCreateStatusOk = paramsExternalEventCreate.Status.Enum()[ExternalEventCreateStatus]
			if ExternalEventCreateStatus != "" && !ExternalEventCreateStatusOk {
				return fmt.Errorf("invalid %v flag value: '%v'", "status", ExternalEventCreateStatus)
			}
			externalEvent, err = client.Create(ctx, paramsExternalEventCreate)
			lib.HandleResponse(ctx, Profile(cmd), externalEvent, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdCreate.Flags().StringVar(&ExternalEventCreateStatus, "status", "", fmt.Sprintf("Status of event. %v", reflect.ValueOf(paramsExternalEventCreate.Status.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&paramsExternalEventCreate.Body, "body", "", "Event body")

	cmdCreate.Flags().StringVar(&fieldsCreate, "fields", "", "comma separated list of field names")
	cmdCreate.Flags().StringVar(&formatCreate, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	ExternalEvents.AddCommand(cmdCreate)
	return ExternalEvents
}
