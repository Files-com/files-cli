package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	action_notification_export_result "github.com/Files-com/files-sdk-go/v2/actionnotificationexportresult"
)

func init() {
	RootCmd.AddCommand(ActionNotificationExportResults())
}

func ActionNotificationExportResults() *cobra.Command {
	ActionNotificationExportResults := &cobra.Command{
		Use:  "action-notification-export-results [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command action-notification-export-results\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	paramsActionNotificationExportResultList := files_sdk.ActionNotificationExportResultListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List Action Notification Export Results",
		Long:  `List Action Notification Export Results`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsActionNotificationExportResultList
			params.MaxPages = MaxPagesList

			client := action_notification_export_result.Client{Config: *config}
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
				return lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			if err != nil {
				return lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}

	cmdList.Flags().Int64Var(&paramsActionNotificationExportResultList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsActionNotificationExportResultList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsActionNotificationExportResultList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64Var(&paramsActionNotificationExportResultList.ActionNotificationExportId, "action-notification-export-id", 0, "ID of the associated action notification export.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	ActionNotificationExportResults.AddCommand(cmdList)
	return ActionNotificationExportResults
}
