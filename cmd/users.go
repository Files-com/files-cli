package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/user"
)

var (
	Users = &cobra.Command{
		Use:  "users [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func UsersInit() {
	var fieldsList string
	paramsUserList := files_sdk.UserListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsUserList
			params.MaxPages = MaxPagesList
			it, err := user.List(params)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsUserList.Page, "page", "p", 0, "Current page number.")
	cmdList.Flags().IntVarP(&paramsUserList.PerPage, "per-page", "r", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsUserList.Action, "action", "a", "", "Deprecated: If set to `count` returns a count of matching records rather than the records themselves.")
	cmdList.Flags().StringVarP(&paramsUserList.Cursor, "cursor", "c", "", "Send cursor to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().StringVarP(&paramsUserList.Ids, "ids", "i", "", "comma-separated list of User IDs")
	cmdList.Flags().StringVarP(&paramsUserList.Search, "search", "", "", "Searches for partial matches of name, username, or email.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	Users.AddCommand(cmdList)
	var fieldsFind string
	paramsUserFind := files_sdk.UserFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := user.Find(paramsUserFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsFind)
		},
	}
	cmdFind.Flags().Int64VarP(&paramsUserFind.Id, "id", "i", 0, "User ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	Users.AddCommand(cmdFind)
	var fieldsCreate string
	paramsUserCreate := files_sdk.UserCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := user.Create(paramsUserCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().StringVarP(&paramsUserCreate.ChangePassword, "change-password", "s", "", "Used for changing a password on an existing user.")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.ChangePasswordConfirmation, "change-password-confirmation", "c", "", "Optional, but if provided, we will ensure that it matches the value sent in `change_password`.")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.Email, "email", "", "", "User's email.")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.GrantPermission, "grant-permission", "g", "", "Permission to grant on the user root.  Can be blank or `full`, `read`, `write`, `list`, or `history`.")
	cmdCreate.Flags().Int64VarP(&paramsUserCreate.GroupId, "group-id", "", 0, "Group ID to associate this user with.")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.GroupIds, "group-ids", "", "", "A list of group ids to associate this user with.  Comma delimited.")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.Password, "password", "w", "", "User password.")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.PasswordConfirmation, "password-confirmation", "", "", "Optional, but if provided, we will ensure that it matches the value sent in `password`.")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.AllowedIps, "allowed-ips", "a", "", "A list of allowed IPs if applicable.  Newline delimited")
	lib.TimeVarP(cmdCreate.Flags(), &paramsUserCreate.AuthenticateUntil, "authenticate-until", "u")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.AuthenticationMethod, "authentication-method", "e", "", "How is this user authenticated?")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.HeaderText, "header-text", "x", "", "Text to display to the user in the header of the UI")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.Language, "language", "", "", "Preferred language")
	cmdCreate.Flags().IntVarP(&paramsUserCreate.NotificationDailySendTime, "notification-daily-send-time", "y", 0, "Hour of the day at which daily notifications should be sent. Can be in range 0 to 23")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.Name, "name", "", "", "User's full name")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.Company, "company", "o", "", "User's company")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.Notes, "notes", "", "", "Any internal notes on the user")
	cmdCreate.Flags().IntVarP(&paramsUserCreate.PasswordValidityDays, "password-validity-days", "", 0, "Number of days to allow user to use the same password")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.SslRequired, "ssl-required", "", "", "SSL required setting")
	cmdCreate.Flags().Int64VarP(&paramsUserCreate.SsoStrategyId, "sso-strategy-id", "", 0, "SSO (Single Sign On) strategy ID for the user, if applicable.")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.TimeZone, "time-zone", "", "", "User time zone")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.UserRoot, "user-root", "", "", "Root folder for FTP (and optionally SFTP if the appropriate site-wide setting is set.)  Note that this is not used for API, Desktop, or Web interface.")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.Username, "username", "", "", "User's username")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	Users.AddCommand(cmdCreate)
	var fieldsUnlock string
	paramsUserUnlock := files_sdk.UserUnlockParams{}
	cmdUnlock := &cobra.Command{
		Use: "unlock",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := user.Unlock(paramsUserUnlock)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsUnlock)
		},
	}
	cmdUnlock.Flags().Int64VarP(&paramsUserUnlock.Id, "id", "i", 0, "User ID.")

	cmdUnlock.Flags().StringVarP(&fieldsUnlock, "fields", "", "", "comma separated list of field names")
	Users.AddCommand(cmdUnlock)
	var fieldsResendWelcomeEmail string
	paramsUserResendWelcomeEmail := files_sdk.UserResendWelcomeEmailParams{}
	cmdResendWelcomeEmail := &cobra.Command{
		Use: "resend-welcome-email",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := user.ResendWelcomeEmail(paramsUserResendWelcomeEmail)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsResendWelcomeEmail)
		},
	}
	cmdResendWelcomeEmail.Flags().Int64VarP(&paramsUserResendWelcomeEmail.Id, "id", "i", 0, "User ID.")

	cmdResendWelcomeEmail.Flags().StringVarP(&fieldsResendWelcomeEmail, "fields", "", "", "comma separated list of field names")
	Users.AddCommand(cmdResendWelcomeEmail)
	var fieldsUser2faReset string
	paramsUserUser2faReset := files_sdk.UserUser2faResetParams{}
	cmdUser2faReset := &cobra.Command{
		Use: "user-2fa-reset",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := user.User2faReset(paramsUserUser2faReset)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsUser2faReset)
		},
	}
	cmdUser2faReset.Flags().Int64VarP(&paramsUserUser2faReset.Id, "id", "i", 0, "User ID.")

	cmdUser2faReset.Flags().StringVarP(&fieldsUser2faReset, "fields", "", "", "comma separated list of field names")
	Users.AddCommand(cmdUser2faReset)
	var fieldsUpdate string
	paramsUserUpdate := files_sdk.UserUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := user.Update(paramsUserUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsUpdate)
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsUserUpdate.Id, "id", "", 0, "User ID.")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.ChangePassword, "change-password", "s", "", "Used for changing a password on an existing user.")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.ChangePasswordConfirmation, "change-password-confirmation", "c", "", "Optional, but if provided, we will ensure that it matches the value sent in `change_password`.")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.Email, "email", "", "", "User's email.")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.GrantPermission, "grant-permission", "g", "", "Permission to grant on the user root.  Can be blank or `full`, `read`, `write`, `list`, or `history`.")
	cmdUpdate.Flags().Int64VarP(&paramsUserUpdate.GroupId, "group-id", "", 0, "Group ID to associate this user with.")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.GroupIds, "group-ids", "", "", "A list of group ids to associate this user with.  Comma delimited.")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.Password, "password", "w", "", "User password.")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.PasswordConfirmation, "password-confirmation", "", "", "Optional, but if provided, we will ensure that it matches the value sent in `password`.")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.AllowedIps, "allowed-ips", "a", "", "A list of allowed IPs if applicable.  Newline delimited")
	lib.TimeVarP(cmdUpdate.Flags(), &paramsUserUpdate.AuthenticateUntil, "authenticate-until", "u")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.AuthenticationMethod, "authentication-method", "e", "", "How is this user authenticated?")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.HeaderText, "header-text", "x", "", "Text to display to the user in the header of the UI")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.Language, "language", "", "", "Preferred language")
	cmdUpdate.Flags().IntVarP(&paramsUserUpdate.NotificationDailySendTime, "notification-daily-send-time", "y", 0, "Hour of the day at which daily notifications should be sent. Can be in range 0 to 23")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.Name, "name", "", "", "User's full name")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.Company, "company", "o", "", "User's company")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.Notes, "notes", "", "", "Any internal notes on the user")
	cmdUpdate.Flags().IntVarP(&paramsUserUpdate.PasswordValidityDays, "password-validity-days", "", 0, "Number of days to allow user to use the same password")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.SslRequired, "ssl-required", "", "", "SSL required setting")
	cmdUpdate.Flags().Int64VarP(&paramsUserUpdate.SsoStrategyId, "sso-strategy-id", "", 0, "SSO (Single Sign On) strategy ID for the user, if applicable.")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.TimeZone, "time-zone", "", "", "User time zone")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.UserRoot, "user-root", "", "", "Root folder for FTP (and optionally SFTP if the appropriate site-wide setting is set.)  Note that this is not used for API, Desktop, or Web interface.")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.Username, "username", "", "", "User's username")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	Users.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsUserDelete := files_sdk.UserDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := user.Delete(paramsUserDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsUserDelete.Id, "id", "i", 0, "User ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Users.AddCommand(cmdDelete)
}
