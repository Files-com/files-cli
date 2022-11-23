package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	flib "github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/Files-com/files-sdk-go/v2/notification"
)

func init() {
	RootCmd.AddCommand(Notifications())
}

func Notifications() *cobra.Command {
	Notifications := &cobra.Command{
		Use:  "notifications [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command notifications\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	paramsNotificationList := files_sdk.NotificationListParams{}
	var MaxPagesList int64
	listIncludeAncestors := true

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List Notifications",
		Long:  `List Notifications`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsNotificationList
			params.MaxPages = MaxPagesList
			if cmd.Flags().Changed("include-ancestors") {
				paramsNotificationList.IncludeAncestors = flib.Bool(listIncludeAncestors)
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

	cmdList.Flags().Int64Var(&paramsNotificationList.UserId, "user-id", 0, "DEPRECATED: Show notifications for this User ID. Use `filter[user_id]` instead.")
	cmdList.Flags().StringVar(&paramsNotificationList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsNotificationList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64Var(&paramsNotificationList.GroupId, "group-id", 0, "DEPRECATED: Show notifications for this Group ID. Use `filter[group_id]` instead.")
	cmdList.Flags().StringVar(&paramsNotificationList.Path, "path", "", "Show notifications for this Path.")
	cmdList.Flags().BoolVar(&listIncludeAncestors, "include-ancestors", listIncludeAncestors, "If `include_ancestors` is `true` and `path` is specified, include notifications for any parent paths. Ignored if `path` is not specified.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	Notifications.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsNotificationFind := files_sdk.NotificationFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Notification`,
		Long:  `Show Notification`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := notification.Client{Config: *config}

			var notification interface{}
			var err error
			notification, err = client.Find(ctx, paramsNotificationFind)
			lib.HandleResponse(ctx, Profile(cmd), notification, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdFind.Flags().Int64Var(&paramsNotificationFind.Id, "id", 0, "Notification ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	Notifications.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createNotifyOnCopy := true
	createNotifyOnDelete := true
	createNotifyOnDownload := true
	createNotifyOnMove := true
	createNotifyOnUpload := true
	createNotifyUserActions := true
	createRecursive := true
	createTriggerByShareRecipients := true
	paramsNotificationCreate := files_sdk.NotificationCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create [path]",
		Short: `Create Notification`,
		Long:  `Create Notification`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := notification.Client{Config: *config}

			if cmd.Flags().Changed("notify-on-copy") {
				paramsNotificationCreate.NotifyOnCopy = flib.Bool(createNotifyOnCopy)
			}
			if cmd.Flags().Changed("notify-on-delete") {
				paramsNotificationCreate.NotifyOnDelete = flib.Bool(createNotifyOnDelete)
			}
			if cmd.Flags().Changed("notify-on-download") {
				paramsNotificationCreate.NotifyOnDownload = flib.Bool(createNotifyOnDownload)
			}
			if cmd.Flags().Changed("notify-on-move") {
				paramsNotificationCreate.NotifyOnMove = flib.Bool(createNotifyOnMove)
			}
			if cmd.Flags().Changed("notify-on-upload") {
				paramsNotificationCreate.NotifyOnUpload = flib.Bool(createNotifyOnUpload)
			}
			if cmd.Flags().Changed("notify-user-actions") {
				paramsNotificationCreate.NotifyUserActions = flib.Bool(createNotifyUserActions)
			}
			if cmd.Flags().Changed("recursive") {
				paramsNotificationCreate.Recursive = flib.Bool(createRecursive)
			}
			if cmd.Flags().Changed("trigger-by-share-recipients") {
				paramsNotificationCreate.TriggerByShareRecipients = flib.Bool(createTriggerByShareRecipients)
			}

			if len(args) > 0 && args[0] != "" {
				paramsNotificationCreate.Path = args[0]
			}
			var notification interface{}
			var err error
			notification, err = client.Create(ctx, paramsNotificationCreate)
			lib.HandleResponse(ctx, Profile(cmd), notification, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
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

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Notifications.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateNotifyOnCopy := true
	updateNotifyOnDelete := true
	updateNotifyOnDownload := true
	updateNotifyOnMove := true
	updateNotifyOnUpload := true
	updateNotifyUserActions := true
	updateRecursive := true
	updateTriggerByShareRecipients := true
	paramsNotificationUpdate := files_sdk.NotificationUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Notification`,
		Long:  `Update Notification`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := notification.Client{Config: *config}

			if cmd.Flags().Changed("notify-on-copy") {
				paramsNotificationUpdate.NotifyOnCopy = flib.Bool(updateNotifyOnCopy)
			}
			if cmd.Flags().Changed("notify-on-delete") {
				paramsNotificationUpdate.NotifyOnDelete = flib.Bool(updateNotifyOnDelete)
			}
			if cmd.Flags().Changed("notify-on-download") {
				paramsNotificationUpdate.NotifyOnDownload = flib.Bool(updateNotifyOnDownload)
			}
			if cmd.Flags().Changed("notify-on-move") {
				paramsNotificationUpdate.NotifyOnMove = flib.Bool(updateNotifyOnMove)
			}
			if cmd.Flags().Changed("notify-on-upload") {
				paramsNotificationUpdate.NotifyOnUpload = flib.Bool(updateNotifyOnUpload)
			}
			if cmd.Flags().Changed("notify-user-actions") {
				paramsNotificationUpdate.NotifyUserActions = flib.Bool(updateNotifyUserActions)
			}
			if cmd.Flags().Changed("recursive") {
				paramsNotificationUpdate.Recursive = flib.Bool(updateRecursive)
			}
			if cmd.Flags().Changed("trigger-by-share-recipients") {
				paramsNotificationUpdate.TriggerByShareRecipients = flib.Bool(updateTriggerByShareRecipients)
			}

			var notification interface{}
			var err error
			notification, err = client.Update(ctx, paramsNotificationUpdate)
			lib.HandleResponse(ctx, Profile(cmd), notification, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
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

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	Notifications.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsNotificationDelete := files_sdk.NotificationDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Notification`,
		Long:  `Delete Notification`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := notification.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsNotificationDelete)
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsNotificationDelete.Id, "id", 0, "Notification ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Notifications.AddCommand(cmdDelete)
	return Notifications
}
