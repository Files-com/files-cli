package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	flib "github.com/Files-com/files-sdk-go/v2/lib"

	"fmt"

	bundle_notification "github.com/Files-com/files-sdk-go/v2/bundlenotification"
)

func init() {
	RootCmd.AddCommand(BundleNotifications())
}

func BundleNotifications() *cobra.Command {
	BundleNotifications := &cobra.Command{
		Use:  "bundle-notifications [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command bundle-notifications\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	paramsBundleNotificationList := files_sdk.BundleNotificationListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List Bundle Notifications",
		Long:  `List Bundle Notifications`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsBundleNotificationList
			params.MaxPages = MaxPagesList

			client := bundle_notification.Client{Config: *config}
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
				return lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			if err != nil {
				return lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}

	cmdList.Flags().Int64Var(&paramsBundleNotificationList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsBundleNotificationList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsBundleNotificationList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64Var(&paramsBundleNotificationList.BundleId, "bundle-id", 0, "Bundle ID to notify on")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	BundleNotifications.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsBundleNotificationFind := files_sdk.BundleNotificationFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Bundle Notification`,
		Long:  `Show Bundle Notification`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := bundle_notification.Client{Config: *config}

			var bundleNotification interface{}
			var err error
			bundleNotification, err = client.Find(ctx, paramsBundleNotificationFind)
			lib.HandleResponse(ctx, Profile(cmd), bundleNotification, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdFind.Flags().Int64Var(&paramsBundleNotificationFind.Id, "id", 0, "Bundle Notification ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	BundleNotifications.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createNotifyOnRegistration := true
	createNotifyOnUpload := true
	paramsBundleNotificationCreate := files_sdk.BundleNotificationCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Bundle Notification`,
		Long:  `Create Bundle Notification`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := bundle_notification.Client{Config: *config}

			if cmd.Flags().Changed("notify-on-registration") {
				paramsBundleNotificationCreate.NotifyOnRegistration = flib.Bool(createNotifyOnRegistration)
			}
			if cmd.Flags().Changed("notify-on-upload") {
				paramsBundleNotificationCreate.NotifyOnUpload = flib.Bool(createNotifyOnUpload)
			}

			var bundleNotification interface{}
			var err error
			bundleNotification, err = client.Create(ctx, paramsBundleNotificationCreate)
			lib.HandleResponse(ctx, Profile(cmd), bundleNotification, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdCreate.Flags().Int64Var(&paramsBundleNotificationCreate.UserId, "user-id", 0, "The id of the user to notify.")
	cmdCreate.Flags().BoolVar(&createNotifyOnRegistration, "notify-on-registration", createNotifyOnRegistration, "Triggers bundle notification when a registration action occurs for it.")
	cmdCreate.Flags().BoolVar(&createNotifyOnUpload, "notify-on-upload", createNotifyOnUpload, "Triggers bundle notification when a upload action occurs for it.")
	cmdCreate.Flags().Int64Var(&paramsBundleNotificationCreate.BundleId, "bundle-id", 0, "Bundle ID to notify on")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	BundleNotifications.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateNotifyOnRegistration := true
	updateNotifyOnUpload := true
	paramsBundleNotificationUpdate := files_sdk.BundleNotificationUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Bundle Notification`,
		Long:  `Update Bundle Notification`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := bundle_notification.Client{Config: *config}

			mapParams, convertErr := lib.StructToMap(files_sdk.BundleNotificationUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsBundleNotificationUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("notify-on-registration") {
				mapParams["notify_on_registration"] = updateNotifyOnRegistration
			}
			if cmd.Flags().Changed("notify-on-upload") {
				mapParams["notify_on_upload"] = updateNotifyOnUpload
			}

			var bundleNotification interface{}
			var err error
			bundleNotification, err = client.UpdateWithMap(ctx, mapParams)
			lib.HandleResponse(ctx, Profile(cmd), bundleNotification, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsBundleNotificationUpdate.Id, "id", 0, "Bundle Notification ID.")
	cmdUpdate.Flags().BoolVar(&updateNotifyOnRegistration, "notify-on-registration", updateNotifyOnRegistration, "Triggers bundle notification when a registration action occurs for it.")
	cmdUpdate.Flags().BoolVar(&updateNotifyOnUpload, "notify-on-upload", updateNotifyOnUpload, "Triggers bundle notification when a upload action occurs for it.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	BundleNotifications.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsBundleNotificationDelete := files_sdk.BundleNotificationDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Bundle Notification`,
		Long:  `Delete Bundle Notification`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := bundle_notification.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsBundleNotificationDelete)
			if err != nil {
				return lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsBundleNotificationDelete.Id, "id", 0, "Bundle Notification ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	BundleNotifications.AddCommand(cmdDelete)
	return BundleNotifications
}