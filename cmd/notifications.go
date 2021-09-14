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
	paramsNotificationList := files_sdk.NotificationListParams{}
	var MaxPagesList int64
	listIncludeAncestors := false

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
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
			if err != nil {
				lib.ClientError(ctx, err)
			}
			err = lib.FormatIter(it, formatList, fieldsList)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdList.Flags().Int64VarP(&paramsNotificationList.UserId, "user-id", "u", 0, "DEPRECATED: Show notifications for this User ID. Use `filter[user_id]` instead.")
	cmdList.Flags().StringVarP(&paramsNotificationList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsNotificationList.PerPage, "per-page", "a", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64VarP(&paramsNotificationList.GroupId, "group-id", "r", 0, "DEPRECATED: Show notifications for this Group ID. Use `filter[group_id]` instead.")
	cmdList.Flags().StringVarP(&paramsNotificationList.Path, "path", "p", "", "Show notifications for this Path.")
	cmdList.Flags().BoolVarP(&listIncludeAncestors, "include-ancestors", "i", listIncludeAncestors, "If `include_ancestors` is `true` and `path` is specified, include notifications for any parent paths. Ignored if `path` is not specified.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Notifications.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	paramsNotificationFind := files_sdk.NotificationFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := notification.Client{Config: *config}

			result, err := client.Find(ctx, paramsNotificationFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatFind, fieldsFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsNotificationFind.Id, "id", "i", 0, "Notification ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Notifications.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	createNotifyOnCopy := false
	createNotifyUserActions := false
	createRecursive := false
	paramsNotificationCreate := files_sdk.NotificationCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := notification.Client{Config: *config}

			if createNotifyOnCopy {
				paramsNotificationCreate.NotifyOnCopy = flib.Bool(true)
			}
			if createNotifyUserActions {
				paramsNotificationCreate.NotifyUserActions = flib.Bool(true)
			}
			if createRecursive {
				paramsNotificationCreate.Recursive = flib.Bool(true)
			}

			if len(args) > 0 && args[0] != "" {
				paramsNotificationCreate.Path = args[0]
			}

			result, err := client.Create(ctx, paramsNotificationCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatCreate, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().Int64VarP(&paramsNotificationCreate.UserId, "user-id", "u", 0, "The id of the user to notify. Provide `user_id`, `username` or `group_id`.")
	cmdCreate.Flags().BoolVarP(&createNotifyOnCopy, "notify-on-copy", "c", createNotifyOnCopy, "If `true`, copying or moving resources into this path will trigger a notification, in addition to just uploads.")
	cmdCreate.Flags().BoolVarP(&createNotifyUserActions, "notify-user-actions", "a", createNotifyUserActions, "If `true` actions initiated by the user will still result in a notification")
	cmdCreate.Flags().BoolVarP(&createRecursive, "recursive", "r", createRecursive, "If `true`, enable notifications for each subfolder in this path")
	cmdCreate.Flags().StringVarP(&paramsNotificationCreate.SendInterval, "send-interval", "s", "", "The time interval that notifications are aggregated by.  Can be `five_minutes`, `fifteen_minutes`, `hourly`, or `daily`.")
	cmdCreate.Flags().Int64VarP(&paramsNotificationCreate.GroupId, "group-id", "g", 0, "The ID of the group to notify.  Provide `user_id`, `username` or `group_id`.")
	cmdCreate.Flags().StringVarP(&paramsNotificationCreate.Path, "path", "p", "", "Path")
	cmdCreate.Flags().StringVarP(&paramsNotificationCreate.Username, "username", "e", "", "The username of the user to notify.  Provide `user_id`, `username` or `group_id`.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Notifications.AddCommand(cmdCreate)
	var fieldsUpdate string
	var formatUpdate string
	updateNotifyOnCopy := false
	updateNotifyUserActions := false
	updateRecursive := false
	paramsNotificationUpdate := files_sdk.NotificationUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := notification.Client{Config: *config}

			if updateNotifyOnCopy {
				paramsNotificationUpdate.NotifyOnCopy = flib.Bool(true)
			}
			if updateNotifyUserActions {
				paramsNotificationUpdate.NotifyUserActions = flib.Bool(true)
			}
			if updateRecursive {
				paramsNotificationUpdate.Recursive = flib.Bool(true)
			}

			result, err := client.Update(ctx, paramsNotificationUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatUpdate, fieldsUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsNotificationUpdate.Id, "id", "i", 0, "Notification ID.")
	cmdUpdate.Flags().BoolVarP(&updateNotifyOnCopy, "notify-on-copy", "c", updateNotifyOnCopy, "If `true`, copying or moving resources into this path will trigger a notification, in addition to just uploads.")
	cmdUpdate.Flags().BoolVarP(&updateNotifyUserActions, "notify-user-actions", "a", updateNotifyUserActions, "If `true` actions initiated by the user will still result in a notification")
	cmdUpdate.Flags().BoolVarP(&updateRecursive, "recursive", "r", updateRecursive, "If `true`, enable notifications for each subfolder in this path")
	cmdUpdate.Flags().StringVarP(&paramsNotificationUpdate.SendInterval, "send-interval", "s", "", "The time interval that notifications are aggregated by.  Can be `five_minutes`, `fifteen_minutes`, `hourly`, or `daily`.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Notifications.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	paramsNotificationDelete := files_sdk.NotificationDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := notification.Client{Config: *config}

			result, err := client.Delete(ctx, paramsNotificationDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatDelete, fieldsDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsNotificationDelete.Id, "id", "i", 0, "Notification ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Notifications.AddCommand(cmdDelete)
}
