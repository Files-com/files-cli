package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	flib "github.com/Files-com/files-sdk-go/lib"

	action_notification_export "github.com/Files-com/files-sdk-go/actionnotificationexport"
)

var (
	ActionNotificationExports = &cobra.Command{}
)

func ActionNotificationExportsInit() {
	ActionNotificationExports = &cobra.Command{
		Use:  "action-notification-exports [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
	var fieldsFind string
	paramsActionNotificationExportFind := files_sdk.ActionNotificationExportFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := action_notification_export.Client{Config: *ctx.GetConfig()}

			result, err := client.Find(paramsActionNotificationExportFind)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsFind)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsActionNotificationExportFind.Id, "id", "i", 0, "Action Notification Export ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	ActionNotificationExports.AddCommand(cmdFind)
	var fieldsCreate string
	createQuerySuccess := false
	paramsActionNotificationExportCreate := files_sdk.ActionNotificationExportCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := action_notification_export.Client{Config: *ctx.GetConfig()}

			if createQuerySuccess {
				paramsActionNotificationExportCreate.QuerySuccess = flib.Bool(true)
			}

			result, err := client.Create(paramsActionNotificationExportCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
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
	ActionNotificationExports.AddCommand(cmdCreate)
}
