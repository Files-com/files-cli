package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	action_notification_export_result "github.com/Files-com/files-sdk-go/actionnotificationexportresult"
)

var (
	ActionNotificationExportResults = &cobra.Command{}
)

func ActionNotificationExportResultsInit() {
	ActionNotificationExportResults = &cobra.Command{
		Use:  "action-notification-export-results [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
	var fieldsList string
	paramsActionNotificationExportResultList := files_sdk.ActionNotificationExportResultListParams{}
	var MaxPagesList int64
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			params := paramsActionNotificationExportResultList
			params.MaxPages = MaxPagesList
			client := action_notification_export_result.Client{Config: *ctx.GetConfig()}
			it, err := client.List(params)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
			err = lib.JsonMarshalIter(it, fieldsList)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdList.Flags().Int64VarP(&paramsActionNotificationExportResultList.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVarP(&paramsActionNotificationExportResultList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsActionNotificationExportResultList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64VarP(&paramsActionNotificationExportResultList.ActionNotificationExportId, "action-notification-export-id", "a", 0, "ID of the associated action notification export.")
	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	ActionNotificationExportResults.AddCommand(cmdList)
}
