package cmd

import (
	"fmt"
	"time"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	action_notification_export "github.com/Files-com/files-sdk-go/v3/actionnotificationexport"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(ActionNotificationExports())
}

func ActionNotificationExports() *cobra.Command {
	ActionNotificationExports := &cobra.Command{
		Use:  "action-notification-exports [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command action-notification-exports\n\t%v", args[0])
		},
	}
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsActionNotificationExportFind := files_sdk.ActionNotificationExportFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Action Notification Export`,
		Long:  `Show Action Notification Export`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := action_notification_export.Client{Config: config}

			var actionNotificationExport interface{}
			var err error
			actionNotificationExport, err = client.Find(paramsActionNotificationExportFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), actionNotificationExport, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsActionNotificationExportFind.Id, "id", 0, "Action Notification Export ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	ActionNotificationExports.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createQuerySuccess := true
	paramsActionNotificationExportCreate := files_sdk.ActionNotificationExportCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Action Notification Export`,
		Long:  `Create Action Notification Export`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := action_notification_export.Client{Config: config}

			if cmd.Flags().Changed("query-success") {
				paramsActionNotificationExportCreate.QuerySuccess = flib.Bool(createQuerySuccess)
			}

			if paramsActionNotificationExportCreate.StartAt.IsZero() {
				paramsActionNotificationExportCreate.StartAt = nil
			}
			if paramsActionNotificationExportCreate.EndAt.IsZero() {
				paramsActionNotificationExportCreate.EndAt = nil
			}

			var actionNotificationExport interface{}
			var err error
			actionNotificationExport, err = client.Create(paramsActionNotificationExportCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), actionNotificationExport, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().Int64Var(&paramsActionNotificationExportCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	paramsActionNotificationExportCreate.StartAt = &time.Time{}
	lib.TimeVar(cmdCreate.Flags(), paramsActionNotificationExportCreate.StartAt, "start-at", "Start date/time of export range.")
	paramsActionNotificationExportCreate.EndAt = &time.Time{}
	lib.TimeVar(cmdCreate.Flags(), paramsActionNotificationExportCreate.EndAt, "end-at", "End date/time of export range.")
	cmdCreate.Flags().StringVar(&paramsActionNotificationExportCreate.QueryMessage, "query-message", "", "Error message associated with the request, if any.")
	cmdCreate.Flags().StringVar(&paramsActionNotificationExportCreate.QueryRequestMethod, "query-request-method", "", "The HTTP request method used by the webhook.")
	cmdCreate.Flags().StringVar(&paramsActionNotificationExportCreate.QueryRequestUrl, "query-request-url", "", "The target webhook URL.")
	cmdCreate.Flags().StringVar(&paramsActionNotificationExportCreate.QueryStatus, "query-status", "", "The HTTP status returned from the server in response to the webhook request.")
	cmdCreate.Flags().BoolVar(&createQuerySuccess, "query-success", createQuerySuccess, "true if the webhook request succeeded (i.e. returned a 200 or 204 response status). false otherwise.")
	cmdCreate.Flags().StringVar(&paramsActionNotificationExportCreate.QueryPath, "query-path", "", "Return notifications that were triggered by actions on this specific path.")
	cmdCreate.Flags().StringVar(&paramsActionNotificationExportCreate.QueryFolder, "query-folder", "", "Return notifications that were triggered by actions in this folder.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	ActionNotificationExports.AddCommand(cmdCreate)
	return ActionNotificationExports
}
