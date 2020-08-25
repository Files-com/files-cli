package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/notification"
)

var (
	Notifications = &cobra.Command{
		Use:  "notifications [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func NotificationsInit() {
	var fieldsList string
	paramsNotificationList := files_sdk.NotificationListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsNotificationList
			params.MaxPages = MaxPagesList
			it, err := notification.List(params)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().Int64VarP(&paramsNotificationList.UserId, "user-id", "u", 0, "DEPRECATED: Show notifications for this User ID. Use `filter[user_id]` instead.")
	cmdList.Flags().IntVarP(&paramsNotificationList.Page, "page", "p", 0, "Current page number.")
	cmdList.Flags().IntVarP(&paramsNotificationList.PerPage, "per-page", "", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsNotificationList.Action, "action", "a", "", "Deprecated: If set to `count` returns a count of matching records rather than the records themselves.")
	cmdList.Flags().StringVarP(&paramsNotificationList.Cursor, "cursor", "c", "", "Send cursor to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsNotificationList.GroupId, "group-id", "r", 0, "DEPRECATED: Show notifications for this Group ID. Use `filter[group_id]` instead.")
	cmdList.Flags().StringVarP(&paramsNotificationList.Path, "path", "", "", "Show notifications for this Path.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	Notifications.AddCommand(cmdList)
	var fieldsFind string
	paramsNotificationFind := files_sdk.NotificationFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := notification.Find(paramsNotificationFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsFind)
		},
	}
	cmdFind.Flags().Int64VarP(&paramsNotificationFind.Id, "id", "i", 0, "Notification ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	Notifications.AddCommand(cmdFind)
	var fieldsCreate string
	paramsNotificationCreate := files_sdk.NotificationCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := notification.Create(paramsNotificationCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().Int64VarP(&paramsNotificationCreate.UserId, "user-id", "u", 0, "The id of the user to notify. Provide `user_id`, `username` or `group_id`.")
	cmdCreate.Flags().StringVarP(&paramsNotificationCreate.SendInterval, "send-interval", "s", "", "The time interval that notifications are aggregated by.  Can be `five_minutes`, `fifteen_minutes`, `hourly`, or `daily`.")
	cmdCreate.Flags().Int64VarP(&paramsNotificationCreate.GroupId, "group-id", "g", 0, "The ID of the group to notify.  Provide `user_id`, `username` or `group_id`.")
	cmdCreate.Flags().StringVarP(&paramsNotificationCreate.Path, "path", "p", "", "Path")
	cmdCreate.Flags().StringVarP(&paramsNotificationCreate.Username, "username", "e", "", "The username of the user to notify.  Provide `user_id`, `username` or `group_id`.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	Notifications.AddCommand(cmdCreate)
	var fieldsUpdate string
	paramsNotificationUpdate := files_sdk.NotificationUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := notification.Update(paramsNotificationUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsUpdate)
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsNotificationUpdate.Id, "id", "i", 0, "Notification ID.")
	cmdUpdate.Flags().StringVarP(&paramsNotificationUpdate.SendInterval, "send-interval", "s", "", "The time interval that notifications are aggregated by.  Can be `five_minutes`, `fifteen_minutes`, `hourly`, or `daily`.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	Notifications.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsNotificationDelete := files_sdk.NotificationDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := notification.Delete(paramsNotificationDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsNotificationDelete.Id, "id", "i", 0, "Notification ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Notifications.AddCommand(cmdDelete)
}
