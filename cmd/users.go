package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	flib "github.com/Files-com/files-sdk-go/lib"
	"github.com/Files-com/files-sdk-go/user"
)

var (
	Users = &cobra.Command{}
)

func UsersInit() {
	Users = &cobra.Command{
		Use:  "users [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
	var fieldsList string
	paramsUserList := files_sdk.UserListParams{}
	var MaxPagesList int64
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			params := paramsUserList
			params.MaxPages = MaxPagesList
			client := user.Client{Config: *ctx.GetConfig()}
			it, err := client.List(params)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
			err = lib.JsonMarshalIter(it, fieldsList)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdList.Flags().StringVarP(&paramsUserList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsUserList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsUserList.Ids, "ids", "i", "", "comma-separated list of User IDs")
	cmdList.Flags().StringVarP(&paramsUserList.Search, "search", "", "", "Searches for partial matches of name, username, or email.")
	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	Users.AddCommand(cmdList)
	var fieldsFind string
	paramsUserFind := files_sdk.UserFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := user.Client{Config: *ctx.GetConfig()}

			result, err := client.Find(paramsUserFind)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsFind)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsUserFind.Id, "id", "i", 0, "User ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	Users.AddCommand(cmdFind)
	var fieldsCreate string
	createAvatarDelete := false
	createAnnouncementsRead := false
	createAttachmentsPermission := false
	createBillingPermission := false
	createBypassInactiveDisable := false
	createBypassSiteAllowedIps := false
	createDavPermission := false
	createDisabled := false
	createFtpPermission := false
	createOfficeIntegrationEnabled := false
	createReceiveAdminAlerts := false
	createRequirePasswordChange := false
	createRestapiPermission := false
	createSelfManaged := false
	createSftpPermission := false
	createSiteAdmin := false
	createSkipWelcomeScreen := false
	createSubscribeToNewsletter := false
	paramsUserCreate := files_sdk.UserCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := user.Client{Config: *ctx.GetConfig()}

			if createAvatarDelete {
				paramsUserCreate.AvatarDelete = flib.Bool(true)
			}
			if createAnnouncementsRead {
				paramsUserCreate.AnnouncementsRead = flib.Bool(true)
			}
			if createAttachmentsPermission {
				paramsUserCreate.AttachmentsPermission = flib.Bool(true)
			}
			if createBillingPermission {
				paramsUserCreate.BillingPermission = flib.Bool(true)
			}
			if createBypassInactiveDisable {
				paramsUserCreate.BypassInactiveDisable = flib.Bool(true)
			}
			if createBypassSiteAllowedIps {
				paramsUserCreate.BypassSiteAllowedIps = flib.Bool(true)
			}
			if createDavPermission {
				paramsUserCreate.DavPermission = flib.Bool(true)
			}
			if createDisabled {
				paramsUserCreate.Disabled = flib.Bool(true)
			}
			if createFtpPermission {
				paramsUserCreate.FtpPermission = flib.Bool(true)
			}
			if createOfficeIntegrationEnabled {
				paramsUserCreate.OfficeIntegrationEnabled = flib.Bool(true)
			}
			if createReceiveAdminAlerts {
				paramsUserCreate.ReceiveAdminAlerts = flib.Bool(true)
			}
			if createRequirePasswordChange {
				paramsUserCreate.RequirePasswordChange = flib.Bool(true)
			}
			if createRestapiPermission {
				paramsUserCreate.RestapiPermission = flib.Bool(true)
			}
			if createSelfManaged {
				paramsUserCreate.SelfManaged = flib.Bool(true)
			}
			if createSftpPermission {
				paramsUserCreate.SftpPermission = flib.Bool(true)
			}
			if createSiteAdmin {
				paramsUserCreate.SiteAdmin = flib.Bool(true)
			}
			if createSkipWelcomeScreen {
				paramsUserCreate.SkipWelcomeScreen = flib.Bool(true)
			}
			if createSubscribeToNewsletter {
				paramsUserCreate.SubscribeToNewsletter = flib.Bool(true)
			}

			result, err := client.Create(paramsUserCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdCreate.Flags().BoolVarP(&createAvatarDelete, "avatar-delete", "d", createAvatarDelete, "If true, the avatar will be deleted.")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.ChangePassword, "change-password", "s", "", "Used for changing a password on an existing user.")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.ChangePasswordConfirmation, "change-password-confirmation", "c", "", "Optional, but if provided, we will ensure that it matches the value sent in `change_password`.")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.Email, "email", "", "", "User's email.")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.GrantPermission, "grant-permission", "g", "", "Permission to grant on the user root.  Can be blank or `full`, `read`, `write`, `list`, or `history`.")
	cmdCreate.Flags().Int64VarP(&paramsUserCreate.GroupId, "group-id", "", 0, "Group ID to associate this user with.")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.GroupIds, "group-ids", "", "", "A list of group ids to associate this user with.  Comma delimited.")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.Password, "password", "w", "", "User password.")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.PasswordConfirmation, "password-confirmation", "", "", "Optional, but if provided, we will ensure that it matches the value sent in `password`.")
	cmdCreate.Flags().BoolVarP(&createAnnouncementsRead, "announcements-read", "n", createAnnouncementsRead, "Signifies that the user has read all the announcements in the UI.")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.AllowedIps, "allowed-ips", "a", "", "A list of allowed IPs if applicable.  Newline delimited")
	cmdCreate.Flags().BoolVarP(&createAttachmentsPermission, "attachments-permission", "t", createAttachmentsPermission, "DEPRECATED: Can the user create Bundles (aka Share Links)? Use the bundle permission instead.")
	lib.TimeVarP(cmdCreate.Flags(), &paramsUserCreate.AuthenticateUntil, "authenticate-until", "u")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.AuthenticationMethod, "authentication-method", "e", "", "How is this user authenticated?")
	cmdCreate.Flags().BoolVarP(&createBillingPermission, "billing-permission", "b", createBillingPermission, "Allow this user to perform operations on the account, payments, and invoices?")
	cmdCreate.Flags().BoolVarP(&createBypassInactiveDisable, "bypass-inactive-disable", "i", createBypassInactiveDisable, "Exempt this user from being disabled based on inactivity?")
	cmdCreate.Flags().BoolVarP(&createBypassSiteAllowedIps, "bypass-site-allowed-ips", "p", createBypassSiteAllowedIps, "Allow this user to skip site-wide IP blacklists?")
	cmdCreate.Flags().BoolVarP(&createDavPermission, "dav-permission", "v", createDavPermission, "Can the user connect with WebDAV?")
	cmdCreate.Flags().BoolVarP(&createDisabled, "disabled", "l", createDisabled, "Is user disabled? Disabled users cannot log in, and do not count for billing purposes.  Users can be automatically disabled after an inactivity period via a Site setting.")
	cmdCreate.Flags().BoolVarP(&createFtpPermission, "ftp-permission", "r", createFtpPermission, "Can the user access with FTP/FTPS?")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.HeaderText, "header-text", "x", "", "Text to display to the user in the header of the UI")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.Language, "language", "", "", "Preferred language")
	cmdCreate.Flags().Int64VarP(&paramsUserCreate.NotificationDailySendTime, "notification-daily-send-time", "y", 0, "Hour of the day at which daily notifications should be sent. Can be in range 0 to 23")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.Name, "name", "", "", "User's full name")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.Company, "company", "o", "", "User's company")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.Notes, "notes", "", "", "Any internal notes on the user")
	cmdCreate.Flags().BoolVarP(&createOfficeIntegrationEnabled, "office-integration-enabled", "", createOfficeIntegrationEnabled, "Enable integration with Office for the web?")
	cmdCreate.Flags().Int64VarP(&paramsUserCreate.PasswordValidityDays, "password-validity-days", "", 0, "Number of days to allow user to use the same password")
	cmdCreate.Flags().BoolVarP(&createReceiveAdminAlerts, "receive-admin-alerts", "", createReceiveAdminAlerts, "Should the user receive admin alerts such a certificate expiration notifications and overages?")
	cmdCreate.Flags().BoolVarP(&createRequirePasswordChange, "require-password-change", "", createRequirePasswordChange, "Is a password change required upon next user login?")
	cmdCreate.Flags().BoolVarP(&createRestapiPermission, "restapi-permission", "", createRestapiPermission, "Can this user access the REST API?")
	cmdCreate.Flags().BoolVarP(&createSelfManaged, "self-managed", "", createSelfManaged, "Does this user manage it's own credentials or is it a shared/bot user?")
	cmdCreate.Flags().BoolVarP(&createSftpPermission, "sftp-permission", "", createSftpPermission, "Can the user access with SFTP?")
	cmdCreate.Flags().BoolVarP(&createSiteAdmin, "site-admin", "", createSiteAdmin, "Is the user an administrator for this site?")
	cmdCreate.Flags().BoolVarP(&createSkipWelcomeScreen, "skip-welcome-screen", "k", createSkipWelcomeScreen, "Skip Welcome page in the UI?")
	cmdCreate.Flags().StringVarP(&paramsUserCreate.SslRequired, "ssl-required", "q", "", "SSL required setting")
	cmdCreate.Flags().Int64VarP(&paramsUserCreate.SsoStrategyId, "sso-strategy-id", "", 0, "SSO (Single Sign On) strategy ID for the user, if applicable.")
	cmdCreate.Flags().BoolVarP(&createSubscribeToNewsletter, "subscribe-to-newsletter", "", createSubscribeToNewsletter, "Is the user subscribed to the newsletter?")
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
			ctx := cmd.Context().(lib.Context)
			client := user.Client{Config: *ctx.GetConfig()}

			result, err := client.Unlock(paramsUserUnlock)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsUnlock)
			if err != nil {
				lib.ClientError(err, &ctx)
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
			ctx := cmd.Context().(lib.Context)
			client := user.Client{Config: *ctx.GetConfig()}

			result, err := client.ResendWelcomeEmail(paramsUserResendWelcomeEmail)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsResendWelcomeEmail)
			if err != nil {
				lib.ClientError(err, &ctx)
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
			ctx := cmd.Context().(lib.Context)
			client := user.Client{Config: *ctx.GetConfig()}

			result, err := client.User2faReset(paramsUserUser2faReset)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsUser2faReset)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdUser2faReset.Flags().Int64VarP(&paramsUserUser2faReset.Id, "id", "i", 0, "User ID.")

	cmdUser2faReset.Flags().StringVarP(&fieldsUser2faReset, "fields", "", "", "comma separated list of field names")
	Users.AddCommand(cmdUser2faReset)
	var fieldsUpdate string
	updateAvatarDelete := false
	updateAnnouncementsRead := false
	updateAttachmentsPermission := false
	updateBillingPermission := false
	updateBypassInactiveDisable := false
	updateBypassSiteAllowedIps := false
	updateDavPermission := false
	updateDisabled := false
	updateFtpPermission := false
	updateOfficeIntegrationEnabled := false
	updateReceiveAdminAlerts := false
	updateRequirePasswordChange := false
	updateRestapiPermission := false
	updateSelfManaged := false
	updateSftpPermission := false
	updateSiteAdmin := false
	updateSkipWelcomeScreen := false
	updateSubscribeToNewsletter := false
	paramsUserUpdate := files_sdk.UserUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := user.Client{Config: *ctx.GetConfig()}

			if updateAvatarDelete {
				paramsUserUpdate.AvatarDelete = flib.Bool(true)
			}
			if updateAnnouncementsRead {
				paramsUserUpdate.AnnouncementsRead = flib.Bool(true)
			}
			if updateAttachmentsPermission {
				paramsUserUpdate.AttachmentsPermission = flib.Bool(true)
			}
			if updateBillingPermission {
				paramsUserUpdate.BillingPermission = flib.Bool(true)
			}
			if updateBypassInactiveDisable {
				paramsUserUpdate.BypassInactiveDisable = flib.Bool(true)
			}
			if updateBypassSiteAllowedIps {
				paramsUserUpdate.BypassSiteAllowedIps = flib.Bool(true)
			}
			if updateDavPermission {
				paramsUserUpdate.DavPermission = flib.Bool(true)
			}
			if updateDisabled {
				paramsUserUpdate.Disabled = flib.Bool(true)
			}
			if updateFtpPermission {
				paramsUserUpdate.FtpPermission = flib.Bool(true)
			}
			if updateOfficeIntegrationEnabled {
				paramsUserUpdate.OfficeIntegrationEnabled = flib.Bool(true)
			}
			if updateReceiveAdminAlerts {
				paramsUserUpdate.ReceiveAdminAlerts = flib.Bool(true)
			}
			if updateRequirePasswordChange {
				paramsUserUpdate.RequirePasswordChange = flib.Bool(true)
			}
			if updateRestapiPermission {
				paramsUserUpdate.RestapiPermission = flib.Bool(true)
			}
			if updateSelfManaged {
				paramsUserUpdate.SelfManaged = flib.Bool(true)
			}
			if updateSftpPermission {
				paramsUserUpdate.SftpPermission = flib.Bool(true)
			}
			if updateSiteAdmin {
				paramsUserUpdate.SiteAdmin = flib.Bool(true)
			}
			if updateSkipWelcomeScreen {
				paramsUserUpdate.SkipWelcomeScreen = flib.Bool(true)
			}
			if updateSubscribeToNewsletter {
				paramsUserUpdate.SubscribeToNewsletter = flib.Bool(true)
			}

			result, err := client.Update(paramsUserUpdate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsUpdate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsUserUpdate.Id, "id", "", 0, "User ID.")
	cmdUpdate.Flags().BoolVarP(&updateAvatarDelete, "avatar-delete", "d", updateAvatarDelete, "If true, the avatar will be deleted.")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.ChangePassword, "change-password", "s", "", "Used for changing a password on an existing user.")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.ChangePasswordConfirmation, "change-password-confirmation", "c", "", "Optional, but if provided, we will ensure that it matches the value sent in `change_password`.")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.Email, "email", "", "", "User's email.")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.GrantPermission, "grant-permission", "g", "", "Permission to grant on the user root.  Can be blank or `full`, `read`, `write`, `list`, or `history`.")
	cmdUpdate.Flags().Int64VarP(&paramsUserUpdate.GroupId, "group-id", "", 0, "Group ID to associate this user with.")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.GroupIds, "group-ids", "", "", "A list of group ids to associate this user with.  Comma delimited.")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.Password, "password", "w", "", "User password.")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.PasswordConfirmation, "password-confirmation", "", "", "Optional, but if provided, we will ensure that it matches the value sent in `password`.")
	cmdUpdate.Flags().BoolVarP(&updateAnnouncementsRead, "announcements-read", "n", updateAnnouncementsRead, "Signifies that the user has read all the announcements in the UI.")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.AllowedIps, "allowed-ips", "a", "", "A list of allowed IPs if applicable.  Newline delimited")
	cmdUpdate.Flags().BoolVarP(&updateAttachmentsPermission, "attachments-permission", "t", updateAttachmentsPermission, "DEPRECATED: Can the user create Bundles (aka Share Links)? Use the bundle permission instead.")
	lib.TimeVarP(cmdUpdate.Flags(), &paramsUserUpdate.AuthenticateUntil, "authenticate-until", "u")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.AuthenticationMethod, "authentication-method", "e", "", "How is this user authenticated?")
	cmdUpdate.Flags().BoolVarP(&updateBillingPermission, "billing-permission", "b", updateBillingPermission, "Allow this user to perform operations on the account, payments, and invoices?")
	cmdUpdate.Flags().BoolVarP(&updateBypassInactiveDisable, "bypass-inactive-disable", "i", updateBypassInactiveDisable, "Exempt this user from being disabled based on inactivity?")
	cmdUpdate.Flags().BoolVarP(&updateBypassSiteAllowedIps, "bypass-site-allowed-ips", "p", updateBypassSiteAllowedIps, "Allow this user to skip site-wide IP blacklists?")
	cmdUpdate.Flags().BoolVarP(&updateDavPermission, "dav-permission", "v", updateDavPermission, "Can the user connect with WebDAV?")
	cmdUpdate.Flags().BoolVarP(&updateDisabled, "disabled", "l", updateDisabled, "Is user disabled? Disabled users cannot log in, and do not count for billing purposes.  Users can be automatically disabled after an inactivity period via a Site setting.")
	cmdUpdate.Flags().BoolVarP(&updateFtpPermission, "ftp-permission", "r", updateFtpPermission, "Can the user access with FTP/FTPS?")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.HeaderText, "header-text", "x", "", "Text to display to the user in the header of the UI")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.Language, "language", "", "", "Preferred language")
	cmdUpdate.Flags().Int64VarP(&paramsUserUpdate.NotificationDailySendTime, "notification-daily-send-time", "y", 0, "Hour of the day at which daily notifications should be sent. Can be in range 0 to 23")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.Name, "name", "", "", "User's full name")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.Company, "company", "o", "", "User's company")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.Notes, "notes", "", "", "Any internal notes on the user")
	cmdUpdate.Flags().BoolVarP(&updateOfficeIntegrationEnabled, "office-integration-enabled", "", updateOfficeIntegrationEnabled, "Enable integration with Office for the web?")
	cmdUpdate.Flags().Int64VarP(&paramsUserUpdate.PasswordValidityDays, "password-validity-days", "", 0, "Number of days to allow user to use the same password")
	cmdUpdate.Flags().BoolVarP(&updateReceiveAdminAlerts, "receive-admin-alerts", "", updateReceiveAdminAlerts, "Should the user receive admin alerts such a certificate expiration notifications and overages?")
	cmdUpdate.Flags().BoolVarP(&updateRequirePasswordChange, "require-password-change", "", updateRequirePasswordChange, "Is a password change required upon next user login?")
	cmdUpdate.Flags().BoolVarP(&updateRestapiPermission, "restapi-permission", "", updateRestapiPermission, "Can this user access the REST API?")
	cmdUpdate.Flags().BoolVarP(&updateSelfManaged, "self-managed", "", updateSelfManaged, "Does this user manage it's own credentials or is it a shared/bot user?")
	cmdUpdate.Flags().BoolVarP(&updateSftpPermission, "sftp-permission", "", updateSftpPermission, "Can the user access with SFTP?")
	cmdUpdate.Flags().BoolVarP(&updateSiteAdmin, "site-admin", "", updateSiteAdmin, "Is the user an administrator for this site?")
	cmdUpdate.Flags().BoolVarP(&updateSkipWelcomeScreen, "skip-welcome-screen", "k", updateSkipWelcomeScreen, "Skip Welcome page in the UI?")
	cmdUpdate.Flags().StringVarP(&paramsUserUpdate.SslRequired, "ssl-required", "q", "", "SSL required setting")
	cmdUpdate.Flags().Int64VarP(&paramsUserUpdate.SsoStrategyId, "sso-strategy-id", "", 0, "SSO (Single Sign On) strategy ID for the user, if applicable.")
	cmdUpdate.Flags().BoolVarP(&updateSubscribeToNewsletter, "subscribe-to-newsletter", "", updateSubscribeToNewsletter, "Is the user subscribed to the newsletter?")
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
			ctx := cmd.Context().(lib.Context)
			client := user.Client{Config: *ctx.GetConfig()}

			result, err := client.Delete(paramsUserDelete)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsUserDelete.Id, "id", "i", 0, "User ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Users.AddCommand(cmdDelete)
}
