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
	usePagerFind := true
	paramsActionNotificationExportFind := files_sdk.ActionNotificationExportFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Action Notification Export`,
		Long:  `Show Action Notification Export`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := action_notification_export.Client{Config: *config}

			var actionNotificationExport interface{}
			var err error
			actionNotificationExport, err = client.Find(ctx, paramsActionNotificationExportFind)
			lib.HandleResponse(ctx, actionNotificationExport, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdFind.Flags().Int64Var(&paramsActionNotificationExportFind.Id, "id", 0, "Action Notification Export ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	ActionNotificationExports.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	usePagerCreate := true
	createQuerySuccess := true
	paramsActionNotificationExportCreate := files_sdk.ActionNotificationExportCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Action Notification Export`,
		Long:  `Create Action Notification Export`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := action_notification_export.Client{Config: *config}

			if cmd.Flags().Changed("query-success") {
				paramsActionNotificationExportCreate.QuerySuccess = flib.Bool(createQuerySuccess)
			}

			var actionNotificationExport interface{}
			var err error
			actionNotificationExport, err = client.Create(ctx, paramsActionNotificationExportCreate)
			lib.HandleResponse(ctx, actionNotificationExport, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdCreate.Flags().Int64Var(&paramsActionNotificationExportCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	lib.TimeVar(cmdCreate.Flags(), paramsActionNotificationExportCreate.StartAt, "start-at")
	lib.TimeVar(cmdCreate.Flags(), paramsActionNotificationExportCreate.EndAt, "end-at")
	cmdCreate.Flags().StringVar(&paramsActionNotificationExportCreate.QueryMessage, "query-message", "", "Error message associated with the request, if any.")
	cmdCreate.Flags().StringVar(&paramsActionNotificationExportCreate.QueryRequestMethod, "query-request-method", "", "The HTTP request method used by the webhook.")
	cmdCreate.Flags().StringVar(&paramsActionNotificationExportCreate.QueryRequestUrl, "query-request-url", "", "The target webhook URL.")
	cmdCreate.Flags().StringVar(&paramsActionNotificationExportCreate.QueryStatus, "query-status", "", "The HTTP status returned from the server in response to the webhook request.")
	cmdCreate.Flags().BoolVar(&createQuerySuccess, "query-success", createQuerySuccess, "true if the webhook request succeeded (i.e. returned a 200 or 204 response status). false otherwise.")
	cmdCreate.Flags().StringVar(&paramsActionNotificationExportCreate.QueryPath, "query-path", "", "Return notifications that were triggered by actions on this specific path.")
	cmdCreate.Flags().StringVar(&paramsActionNotificationExportCreate.QueryFolder, "query-folder", "", "Return notifications that were triggered by actions in this folder.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	ActionNotificationExports.AddCommand(cmdCreate)
}
