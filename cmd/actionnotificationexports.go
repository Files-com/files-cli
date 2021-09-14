package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	flib "github.com/Files-com/files-sdk-go/v2/lib"

	"fmt"

	action_notification_export "github.com/Files-com/files-sdk-go/v2/actionnotificationexport"
)

var (
	ActionNotificationExports = &cobra.Command{}
)

func ActionNotificationExportsInit() {
	ActionNotificationExports = &cobra.Command{
		Use:  "action-notification-exports [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command action-notification-exports\n\t%v", args[0])
		},
	}
	var fieldsFind string
	var formatFind string
	paramsActionNotificationExportFind := files_sdk.ActionNotificationExportFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := action_notification_export.Client{Config: *config}

			result, err := client.Find(ctx, paramsActionNotificationExportFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatFind, fieldsFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsActionNotificationExportFind.Id, "id", "i", 0, "Action Notification Export ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright")
	ActionNotificationExports.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	createQuerySuccess := false
	paramsActionNotificationExportCreate := files_sdk.ActionNotificationExportCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := action_notification_export.Client{Config: *config}

			if createQuerySuccess {
				paramsActionNotificationExportCreate.QuerySuccess = flib.Bool(true)
			}

			result, err := client.Create(ctx, paramsActionNotificationExportCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatCreate, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().Int64VarP(&paramsActionNotificationExportCreate.UserId, "user-id", "i", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	lib.TimeVarP(cmdCreate.Flags(), &paramsActionNotificationExportCreate.StartAt, "start-at", "r")
	lib.TimeVarP(cmdCreate.Flags(), &paramsActionNotificationExportCreate.EndAt, "end-at", "e")
	cmdCreate.Flags().StringVarP(&paramsActionNotificationExportCreate.QueryMessage, "query-message", "s", "", "Error message associated with the request, if any.")
	cmdCreate.Flags().StringVarP(&paramsActionNotificationExportCreate.QueryRequestMethod, "query-request-method", "t", "", "The HTTP request method used by the webhook.")
	cmdCreate.Flags().StringVarP(&paramsActionNotificationExportCreate.QueryRequestUrl, "query-request-url", "u", "", "The target webhook URL.")
	cmdCreate.Flags().StringVarP(&paramsActionNotificationExportCreate.QueryStatus, "query-status", "a", "", "The HTTP status returned from the server in response to the webhook request.")
	cmdCreate.Flags().BoolVarP(&createQuerySuccess, "query-success", "c", createQuerySuccess, "true if the webhook request succeeded (i.e. returned a 200 or 204 response status). false otherwise.")
	cmdCreate.Flags().StringVarP(&paramsActionNotificationExportCreate.QueryPath, "query-path", "p", "", "Return notifications that were triggered by actions on this specific path.")
	cmdCreate.Flags().StringVarP(&paramsActionNotificationExportCreate.QueryFolder, "query-folder", "f", "", "Return notifications that were triggered by actions in this folder.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	ActionNotificationExports.AddCommand(cmdCreate)
}
