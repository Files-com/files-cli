package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	action_notification_export_result "github.com/Files-com/files-sdk-go/v2/actionnotificationexportresult"
)

var (
	ActionNotificationExportResults = &cobra.Command{}
)

func ActionNotificationExportResultsInit() {
	ActionNotificationExportResults = &cobra.Command{
		Use:  "action-notification-export-results [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command action-notification-export-results\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsActionNotificationExportResultList := files_sdk.ActionNotificationExportResultListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsActionNotificationExportResultList
			params.MaxPages = MaxPagesList

			client := action_notification_export_result.Client{Config: *config}
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
	cmdList.Flags().Int64VarP(&paramsActionNotificationExportResultList.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVarP(&paramsActionNotificationExportResultList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsActionNotificationExportResultList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64VarP(&paramsActionNotificationExportResultList.ActionNotificationExportId, "action-notification-export-id", "a", 0, "ID of the associated action notification export.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	ActionNotificationExportResults.AddCommand(cmdList)
}
