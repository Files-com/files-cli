package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

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
			client := user.Client{Config: files_sdk.GlobalConfig}
			it, err := client.List(params)
			if err != nil {
				lib.ClientError(err)
			}
			err = lib.JsonMarshalIter(it, fieldsList)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdList.Flags().StringVarP(&paramsUserList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&paramsUserList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsUserList.Ids, "ids", "i", "", "comma-separated list of User IDs")
	cmdList.Flags().StringVarP(&paramsUserList.Search, "search", "", "", "Searches for partial matches of name, username, or email.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	Users.AddCommand(cmdList)
	var fieldsFind string
	paramsUserFind := files_sdk.UserFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			client := user.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Find(paramsUserFind)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsFind)
			if err != nil {
				lib.ClientError(err)
			}
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
			client := user.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Create(paramsUserCreate)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(err)
			}
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
	cmdCreate.Flags().StringVarP(&paramsUserCreate.SslRequired, "ssl-required", "q", "", "SSL required setting")
	cmdCreate.Flags().Int64VarP(&paramsUserCreate.SsoStrategyId, "sso-strategy-id", "", 0, "SSO (Single Sign On) strategy ID for the user, if applicable.")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.Require2fa, "require-2fa", "2", "", "2FA required setting")
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
			client := user.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Unlock(paramsUserUnlock)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsUnlock)
			if err != nil {
				lib.ClientError(err)
			}
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
			client := user.Client{Config: files_sdk.GlobalConfig}
			result, err := client.ResendWelcomeEmail(paramsUserResendWelcomeEmail)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsResendWelcomeEmail)
			if err != nil {
				lib.ClientError(err)
			}
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
			client := user.Client{Config: files_sdk.GlobalConfig}
			result, err := client.User2faReset(paramsUserUser2faReset)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsUser2faReset)
			if err != nil {
				lib.ClientError(err)
			}
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
			client := user.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Update(paramsUserUpdate)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsUpdate)
			if err != nil {
				lib.ClientError(err)
			}
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
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.SslRequired, "ssl-required", "q", "", "SSL required setting")
	cmdUpdate.Flags().Int64VarP(&paramsUserUpdate.SsoStrategyId, "sso-strategy-id", "", 0, "SSO (Single Sign On) strategy ID for the user, if applicable.")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.Require2fa, "require-2fa", "2", "", "2FA required setting")
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
			client := user.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Delete(paramsUserDelete)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsUserDelete.Id, "id", "i", 0, "User ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Users.AddCommand(cmdDelete)
}
