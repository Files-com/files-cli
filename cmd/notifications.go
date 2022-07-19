package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	flib "github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/Files-com/files-sdk-go/v2/notification"
)

var (
	Notifications = &cobra.Command{}
)

func NotificationsInit() {
	Notifications = &cobra.Command{
		Use:  "notifications [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command notifications\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	usePagerList := true
	paramsNotificationList := files_sdk.NotificationListParams{}
	var MaxPagesList int64
	listIncludeAncestors := false

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List Notifications",
		Long:  `List Notifications`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsNotificationList
			params.MaxPages = MaxPagesList
			if listIncludeAncestors {
				paramsNotificationList.IncludeAncestors = flib.Bool(true)
			}

			client := notification.Client{Config: *config}
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
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}

	cmdList.Flags().Int64Var(&paramsNotificationList.UserId, "user-id", 0, "DEPRECATED: Show notifications for this User ID. Use `filter[user_id]` instead.")
	cmdList.Flags().StringVar(&paramsNotificationList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsNotificationList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64Var(&paramsNotificationList.GroupId, "group-id", 0, "DEPRECATED: Show notifications for this Group ID. Use `filter[group_id]` instead.")
	cmdList.Flags().StringVar(&paramsNotificationList.Path, "path", "", "Show notifications for this Path.")
	cmdList.Flags().BoolVar(&listIncludeAncestors, "include-ancestors", listIncludeAncestors, "If `include_ancestors` is `true` and `path` is specified, include notifications for any parent paths. Ignored if `path` is not specified.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	Notifications.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	usePagerFind := true
	paramsNotificationFind := files_sdk.NotificationFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Notification`,
		Long:  `Show Notification`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := notification.Client{Config: *config}

			var notification interface{}
			var err error
			notification, err = client.Find(ctx, paramsNotificationFind)
			lib.HandleResponse(ctx, notification, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdFind.Flags().Int64Var(&paramsNotificationFind.Id, "id", 0, "Notification ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	Notifications.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	usePagerCreate := true
	createNotifyOnCopy := false
	createNotifyOnDelete := false
	createNotifyOnDownload := false
	createNotifyOnMove := false
	createNotifyOnUpload := false
	createNotifyUserActions := false
	createRecursive := false
	createTriggerByShareRecipients := false
	paramsNotificationCreate := files_sdk.NotificationCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create [path]",
		Short: `Create Notification`,
		Long:  `Create Notification`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := notification.Client{Config: *config}

			if createNotifyOnCopy {
				paramsNotificationCreate.NotifyOnCopy = flib.Bool(true)
			}
			if createNotifyOnDelete {
				paramsNotificationCreate.NotifyOnDelete = flib.Bool(true)
			}
			if createNotifyOnDownload {
				paramsNotificationCreate.NotifyOnDownload = flib.Bool(true)
			}
			if createNotifyOnMove {
				paramsNotificationCreate.NotifyOnMove = flib.Bool(true)
			}
			if createNotifyOnUpload {
				paramsNotificationCreate.NotifyOnUpload = flib.Bool(true)
			}
			if createNotifyUserActions {
				paramsNotificationCreate.NotifyUserActions = flib.Bool(true)
			}
			if createRecursive {
				paramsNotificationCreate.Recursive = flib.Bool(true)
			}
			if createTriggerByShareRecipients {
				paramsNotificationCreate.TriggerByShareRecipients = flib.Bool(true)
			}

			if len(args) > 0 && args[0] != "" {
				paramsNotificationCreate.Path = args[0]
			}
			var notification interface{}
			var err error
			notification, err = client.Create(ctx, paramsNotificationCreate)
			lib.HandleResponse(ctx, notification, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdCreate.Flags().Int64Var(&paramsNotificationCreate.UserId, "user-id", 0, "The id of the user to notify. Provide `user_id`, `username` or `group_id`.")
	cmdCreate.Flags().BoolVar(&createNotifyOnCopy, "notify-on-copy", createNotifyOnCopy, "If `true`, copying or moving resources into this path will trigger a notification, in addition to just uploads.")
	cmdCreate.Flags().BoolVar(&createNotifyOnDelete, "notify-on-delete", createNotifyOnDelete, "Triggers notification when deleting files from this path")
	cmdCreate.Flags().BoolVar(&createNotifyOnDownload, "notify-on-download", createNotifyOnDownload, "Triggers notification when downloading files from this path")
	cmdCreate.Flags().BoolVar(&createNotifyOnMove, "notify-on-move", createNotifyOnMove, "Triggers notification when moving files to this path")
	cmdCreate.Flags().BoolVar(&createNotifyOnUpload, "notify-on-upload", createNotifyOnUpload, "Triggers notification when uploading new files to this path")
	cmdCreate.Flags().BoolVar(&createNotifyUserActions, "notify-user-actions", createNotifyUserActions, "If `true` actions initiated by the user will still result in a notification")
	cmdCreate.Flags().BoolVar(&createRecursive, "recursive", createRecursive, "If `true`, enable notifications for each subfolder in this path")
	cmdCreate.Flags().StringVar(&paramsNotificationCreate.SendInterval, "send-interval", "", "The time interval that notifications are aggregated by.  Can be `five_minutes`, `fifteen_minutes`, `hourly`, or `daily`.")
	cmdCreate.Flags().StringVar(&paramsNotificationCreate.Message, "message", "", "Custom message to include in notification emails.")
	cmdCreate.Flags().StringSliceVar(&paramsNotificationCreate.TriggeringFilenames, "triggering-filenames", []string{}, "Array of filenames (possibly with wildcards) to match for action path")
	cmdCreate.Flags().BoolVar(&createTriggerByShareRecipients, "trigger-by-share-recipients", createTriggerByShareRecipients, "Notify when actions are performed by a share recipient?")
	cmdCreate.Flags().Int64Var(&paramsNotificationCreate.GroupId, "group-id", 0, "The ID of the group to notify.  Provide `user_id`, `username` or `group_id`.")
	cmdCreate.Flags().StringVar(&paramsNotificationCreate.Path, "path", "", "Path")
	cmdCreate.Flags().StringVar(&paramsNotificationCreate.Username, "username", "", "The username of the user to notify.  Provide `user_id`, `username` or `group_id`.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Notifications.AddCommand(cmdCreate)
	var fieldsUpdate string
	var formatUpdate string
	usePagerUpdate := true
	updateNotifyOnCopy := false
	updateNotifyOnDelete := false
	updateNotifyOnDownload := false
	updateNotifyOnMove := false
	updateNotifyOnUpload := false
	updateNotifyUserActions := false
	updateRecursive := false
	updateTriggerByShareRecipients := false
	paramsNotificationUpdate := files_sdk.NotificationUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Notification`,
		Long:  `Update Notification`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := notification.Client{Config: *config}

			if updateNotifyOnCopy {
				paramsNotificationUpdate.NotifyOnCopy = flib.Bool(true)
			}
			if updateNotifyOnDelete {
				paramsNotificationUpdate.NotifyOnDelete = flib.Bool(true)
			}
			if updateNotifyOnDownload {
				paramsNotificationUpdate.NotifyOnDownload = flib.Bool(true)
			}
			if updateNotifyOnMove {
				paramsNotificationUpdate.NotifyOnMove = flib.Bool(true)
			}
			if updateNotifyOnUpload {
				paramsNotificationUpdate.NotifyOnUpload = flib.Bool(true)
			}
			if updateNotifyUserActions {
				paramsNotificationUpdate.NotifyUserActions = flib.Bool(true)
			}
			if updateRecursive {
				paramsNotificationUpdate.Recursive = flib.Bool(true)
			}
			if updateTriggerByShareRecipients {
				paramsNotificationUpdate.TriggerByShareRecipients = flib.Bool(true)
			}

			var notification interface{}
			var err error
			notification, err = client.Update(ctx, paramsNotificationUpdate)
			lib.HandleResponse(ctx, notification, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsNotificationUpdate.Id, "id", 0, "Notification ID.")
	cmdUpdate.Flags().BoolVar(&updateNotifyOnCopy, "notify-on-copy", updateNotifyOnCopy, "If `true`, copying or moving resources into this path will trigger a notification, in addition to just uploads.")
	cmdUpdate.Flags().BoolVar(&updateNotifyOnDelete, "notify-on-delete", updateNotifyOnDelete, "Triggers notification when deleting files from this path")
	cmdUpdate.Flags().BoolVar(&updateNotifyOnDownload, "notify-on-download", updateNotifyOnDownload, "Triggers notification when downloading files from this path")
	cmdUpdate.Flags().BoolVar(&updateNotifyOnMove, "notify-on-move", updateNotifyOnMove, "Triggers notification when moving files to this path")
	cmdUpdate.Flags().BoolVar(&updateNotifyOnUpload, "notify-on-upload", updateNotifyOnUpload, "Triggers notification when uploading new files to this path")
	cmdUpdate.Flags().BoolVar(&updateNotifyUserActions, "notify-user-actions", updateNotifyUserActions, "If `true` actions initiated by the user will still result in a notification")
	cmdUpdate.Flags().BoolVar(&updateRecursive, "recursive", updateRecursive, "If `true`, enable notifications for each subfolder in this path")
	cmdUpdate.Flags().StringVar(&paramsNotificationUpdate.SendInterval, "send-interval", "", "The time interval that notifications are aggregated by.  Can be `five_minutes`, `fifteen_minutes`, `hourly`, or `daily`.")
	cmdUpdate.Flags().StringVar(&paramsNotificationUpdate.Message, "message", "", "Custom message to include in notification emails.")
	cmdUpdate.Flags().StringSliceVar(&paramsNotificationUpdate.TriggeringFilenames, "triggering-filenames", []string{}, "Array of filenames (possibly with wildcards) to match for action path")
	cmdUpdate.Flags().BoolVar(&updateTriggerByShareRecipients, "trigger-by-share-recipients", updateTriggerByShareRecipients, "Notify when actions are performed by a share recipient?")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	Notifications.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	usePagerDelete := true
	paramsNotificationDelete := files_sdk.NotificationDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Notification`,
		Long:  `Delete Notification`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := notification.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsNotificationDelete)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}
	cmdDelete.Flags().Int64Var(&paramsNotificationDelete.Id, "id", 0, "Notification ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Notifications.AddCommand(cmdDelete)
}
