package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/notification"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Notifications())
}

func Notifications() *cobra.Command {
	Notifications := &cobra.Command{
		Use:  "notifications [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command notifications\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsNotificationList := files_sdk.NotificationListParams{}
	var MaxPagesList int64
	listIncludeAncestors := true

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Notifications",
		Long:    `List Notifications`,
		Args:    cobra.RangeArgs(0, 1),
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsNotificationList
			params.MaxPages = MaxPagesList
			if len(args) > 0 && args[0] != "" {
				params.Path = args[0]
			}

			if cmd.Flags().Changed("include-ancestors") {
				params.IncludeAncestors = flib.Bool(listIncludeAncestors)
			}

			client := notification.Client{Config: config}
			it, err := client.List(params, files_sdk.WithContext(ctx))
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger)
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyList) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyList, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatList), fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdList.Flags().StringVar(&paramsNotificationList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsNotificationList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVar(&paramsNotificationList.Path, "path", "", "Show notifications for this Path.")
	cmdList.Flags().BoolVar(&listIncludeAncestors, "include-ancestors", listIncludeAncestors, "If `include_ancestors` is `true` and `path` is specified, include notifications for any parent paths. Ignored if `path` is not specified.")
	cmdList.Flags().StringVar(&paramsNotificationList.GroupId, "group-id", "", "")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
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
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := notification.Client{Config: config}

			var notification interface{}
			var err error
			notification, err = client.Find(paramsNotificationFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), notification, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsNotificationFind.Id, "id", 0, "Notification ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
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
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := notification.Client{Config: config}

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
			notification, err = client.Create(paramsNotificationCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), notification, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().Int64Var(&paramsNotificationCreate.UserId, "user-id", 0, "The id of the user to notify. Provide `user_id`, `username` or `group_id`.")
	cmdCreate.Flags().BoolVar(&createNotifyOnCopy, "notify-on-copy", createNotifyOnCopy, "If `true`, copying or moving resources into this path will trigger a notification, in addition to just uploads.")
	cmdCreate.Flags().BoolVar(&createNotifyOnDelete, "notify-on-delete", createNotifyOnDelete, "Trigger on files deleted in this path?")
	cmdCreate.Flags().BoolVar(&createNotifyOnDownload, "notify-on-download", createNotifyOnDownload, "Trigger on files downloaded in this path?")
	cmdCreate.Flags().BoolVar(&createNotifyOnMove, "notify-on-move", createNotifyOnMove, "Trigger on files moved to this path?")
	cmdCreate.Flags().BoolVar(&createNotifyOnUpload, "notify-on-upload", createNotifyOnUpload, "Trigger on files created/uploaded/updated/changed in this path?")
	cmdCreate.Flags().BoolVar(&createNotifyUserActions, "notify-user-actions", createNotifyUserActions, "If `true` actions initiated by the user will still result in a notification")
	cmdCreate.Flags().BoolVar(&createRecursive, "recursive", createRecursive, "If `true`, enable notifications for each subfolder in this path")
	cmdCreate.Flags().StringVar(&paramsNotificationCreate.SendInterval, "send-interval", "", "The time interval that notifications are aggregated by.  Can be `five_minutes`, `fifteen_minutes`, `hourly`, or `daily`.")
	cmdCreate.Flags().StringVar(&paramsNotificationCreate.Message, "message", "", "Custom message to include in notification emails")
	cmdCreate.Flags().StringSliceVar(&paramsNotificationCreate.TriggeringFilenames, "triggering-filenames", []string{}, "Array of filenames (possibly with wildcards) to scope trigger")
	cmdCreate.Flags().Int64SliceVar(&paramsNotificationCreate.TriggeringGroupIds, "triggering-group-ids", []int64{}, "If set, will only notify on actions made by a member of one of the specified groups")
	cmdCreate.Flags().Int64SliceVar(&paramsNotificationCreate.TriggeringUserIds, "triggering-user-ids", []int64{}, "If set, will only notify on actions made one of the specified users")
	cmdCreate.Flags().BoolVar(&createTriggerByShareRecipients, "trigger-by-share-recipients", createTriggerByShareRecipients, "Notify when actions are performed by a share recipient?")
	cmdCreate.Flags().Int64Var(&paramsNotificationCreate.GroupId, "group-id", 0, "The ID of the group to notify.  Provide `user_id`, `username` or `group_id`.")
	cmdCreate.Flags().StringVar(&paramsNotificationCreate.Path, "path", "", "Path")
	cmdCreate.Flags().StringVar(&paramsNotificationCreate.Username, "username", "", "The username of the user to notify.  Provide `user_id`, `username` or `group_id`.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
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
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := notification.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.NotificationUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsNotificationUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("notify-on-copy") {
				mapParams["notify_on_copy"] = updateNotifyOnCopy
			}
			if cmd.Flags().Changed("notify-on-delete") {
				mapParams["notify_on_delete"] = updateNotifyOnDelete
			}
			if cmd.Flags().Changed("notify-on-download") {
				mapParams["notify_on_download"] = updateNotifyOnDownload
			}
			if cmd.Flags().Changed("notify-on-move") {
				mapParams["notify_on_move"] = updateNotifyOnMove
			}
			if cmd.Flags().Changed("notify-on-upload") {
				mapParams["notify_on_upload"] = updateNotifyOnUpload
			}
			if cmd.Flags().Changed("notify-user-actions") {
				mapParams["notify_user_actions"] = updateNotifyUserActions
			}
			if cmd.Flags().Changed("recursive") {
				mapParams["recursive"] = updateRecursive
			}
			if cmd.Flags().Changed("send-interval") {
				lib.FlagUpdate(cmd, "send_interval", paramsNotificationUpdate.SendInterval, mapParams)
			}
			if cmd.Flags().Changed("message") {
				lib.FlagUpdate(cmd, "message", paramsNotificationUpdate.Message, mapParams)
			}
			if cmd.Flags().Changed("triggering-filenames") {
				lib.FlagUpdateLen(cmd, "triggering_filenames", paramsNotificationUpdate.TriggeringFilenames, mapParams)
			}
			if cmd.Flags().Changed("triggering-group-ids") {
				lib.FlagUpdateLen(cmd, "triggering_group_ids", paramsNotificationUpdate.TriggeringGroupIds, mapParams)
			}
			if cmd.Flags().Changed("triggering-user-ids") {
				lib.FlagUpdateLen(cmd, "triggering_user_ids", paramsNotificationUpdate.TriggeringUserIds, mapParams)
			}
			if cmd.Flags().Changed("trigger-by-share-recipients") {
				mapParams["trigger_by_share_recipients"] = updateTriggerByShareRecipients
			}

			var notification interface{}
			var err error
			notification, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), notification, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsNotificationUpdate.Id, "id", 0, "Notification ID.")
	cmdUpdate.Flags().BoolVar(&updateNotifyOnCopy, "notify-on-copy", updateNotifyOnCopy, "If `true`, copying or moving resources into this path will trigger a notification, in addition to just uploads.")
	cmdUpdate.Flags().BoolVar(&updateNotifyOnDelete, "notify-on-delete", updateNotifyOnDelete, "Trigger on files deleted in this path?")
	cmdUpdate.Flags().BoolVar(&updateNotifyOnDownload, "notify-on-download", updateNotifyOnDownload, "Trigger on files downloaded in this path?")
	cmdUpdate.Flags().BoolVar(&updateNotifyOnMove, "notify-on-move", updateNotifyOnMove, "Trigger on files moved to this path?")
	cmdUpdate.Flags().BoolVar(&updateNotifyOnUpload, "notify-on-upload", updateNotifyOnUpload, "Trigger on files created/uploaded/updated/changed in this path?")
	cmdUpdate.Flags().BoolVar(&updateNotifyUserActions, "notify-user-actions", updateNotifyUserActions, "If `true` actions initiated by the user will still result in a notification")
	cmdUpdate.Flags().BoolVar(&updateRecursive, "recursive", updateRecursive, "If `true`, enable notifications for each subfolder in this path")
	cmdUpdate.Flags().StringVar(&paramsNotificationUpdate.SendInterval, "send-interval", "", "The time interval that notifications are aggregated by.  Can be `five_minutes`, `fifteen_minutes`, `hourly`, or `daily`.")
	cmdUpdate.Flags().StringVar(&paramsNotificationUpdate.Message, "message", "", "Custom message to include in notification emails")
	cmdUpdate.Flags().StringSliceVar(&paramsNotificationUpdate.TriggeringFilenames, "triggering-filenames", []string{}, "Array of filenames (possibly with wildcards) to scope trigger")
	cmdUpdate.Flags().Int64SliceVar(&paramsNotificationUpdate.TriggeringGroupIds, "triggering-group-ids", []int64{}, "If set, will only notify on actions made by a member of one of the specified groups")
	cmdUpdate.Flags().Int64SliceVar(&paramsNotificationUpdate.TriggeringUserIds, "triggering-user-ids", []int64{}, "If set, will only notify on actions made one of the specified users")
	cmdUpdate.Flags().BoolVar(&updateTriggerByShareRecipients, "trigger-by-share-recipients", updateTriggerByShareRecipients, "Notify when actions are performed by a share recipient?")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
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
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := notification.Client{Config: config}

			var err error
			err = client.Delete(paramsNotificationDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsNotificationDelete.Id, "id", 0, "Notification ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Notifications.AddCommand(cmdDelete)
	return Notifications
}
