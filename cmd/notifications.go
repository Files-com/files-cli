package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/notification"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = notification.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
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
			it := notification.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsNotificationList.Page, "page", "p", 0, "List Notifications")
	cmdList.Flags().IntVarP(&paramsNotificationList.PerPage, "per-page", "", 0, "List Notifications")
	cmdList.Flags().StringVarP(&paramsNotificationList.Action, "action", "a", "", "List Notifications")
	cmdList.Flags().StringVarP(&paramsNotificationList.Cursor, "cursor", "c", "", "List Notifications")
	cmdList.Flags().StringVarP(&paramsNotificationList.Path, "path", "", "", "List Notifications")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
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
	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "f", "", "comma separated list of field names")
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
	cmdCreate.Flags().StringVarP(&paramsNotificationCreate.SendInterval, "send-interval", "s", "", "Create Notification")
	cmdCreate.Flags().StringVarP(&paramsNotificationCreate.Path, "path", "p", "", "Create Notification")
	cmdCreate.Flags().StringVarP(&paramsNotificationCreate.Username, "username", "e", "", "Create Notification")
	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
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
	cmdUpdate.Flags().StringVarP(&paramsNotificationUpdate.SendInterval, "send-interval", "s", "", "Update Notification")
	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "f", "", "comma separated list of field names")
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
	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
	Notifications.AddCommand(cmdDelete)
}
