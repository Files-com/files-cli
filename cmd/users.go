package cmd

import (
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	flib "github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/Files-com/files-sdk-go/v2/user"
)

var (
	Users = &cobra.Command{}
)

func UsersInit() {
	Users = &cobra.Command{
		Use:  "users [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command users\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsUserList := files_sdk.UserListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsUserList
			params.MaxPages = MaxPagesList

			client := user.Client{Config: *config}
			it, err := client.List(ctx, params)
			if err != nil {
				lib.ClientError(ctx, err)
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(it, formatList, fieldsList, listFilter)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}

	cmdList.Flags().StringVar(&paramsUserList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsUserList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVar(&paramsUserList.Ids, "ids", "", "comma-separated list of User IDs")
	cmdList.Flags().StringVar(&paramsUserList.Search, "search", "", "Searches for partial matches of name, username, or email.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Users.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	paramsUserFind := files_sdk.UserFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := user.Client{Config: *config}

			result, err := client.Find(ctx, paramsUserFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatFind, fieldsFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdFind.Flags().Int64Var(&paramsUserFind.Id, "id", 0, "User ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Users.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
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
	UserCreateAuthenticationMethod := ""
	UserCreateSslRequired := ""
	UserCreateRequire2fa := ""

	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := user.Client{Config: *config}

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

			paramsUserCreate.AuthenticationMethod = paramsUserCreate.AuthenticationMethod.Enum()[UserCreateAuthenticationMethod]
			paramsUserCreate.SslRequired = paramsUserCreate.SslRequired.Enum()[UserCreateSslRequired]
			paramsUserCreate.Require2fa = paramsUserCreate.Require2fa.Enum()[UserCreateRequire2fa]

			result, err := client.Create(ctx, paramsUserCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatCreate, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().BoolVar(&createAvatarDelete, "avatar-delete", createAvatarDelete, "If true, the avatar will be deleted.")
	cmdCreate.Flags().StringVar(&paramsUserCreate.ChangePassword, "change-password", "", "Used for changing a password on an existing user.")
	cmdCreate.Flags().StringVar(&paramsUserCreate.ChangePasswordConfirmation, "change-password-confirmation", "", "Optional, but if provided, we will ensure that it matches the value sent in `change_password`.")
	cmdCreate.Flags().StringVar(&paramsUserCreate.Email, "email", "", "User's email.")
	cmdCreate.Flags().StringVar(&paramsUserCreate.GrantPermission, "grant-permission", "", "Permission to grant on the user root.  Can be blank or `full`, `read`, `write`, `list`, or `history`.")
	cmdCreate.Flags().Int64Var(&paramsUserCreate.GroupId, "group-id", 0, "Group ID to associate this user with.")
	cmdCreate.Flags().StringVar(&paramsUserCreate.GroupIds, "group-ids", "", "A list of group ids to associate this user with.  Comma delimited.")
	cmdCreate.Flags().StringVar(&paramsUserCreate.ImportedPasswordHash, "imported-password-hash", "", "Pre-calculated hash of the user's password. If supplied, this will be used to authenticate the user on first login. Supported hash menthods are MD5, SHA1, and SHA256.")
	cmdCreate.Flags().StringVar(&paramsUserCreate.Password, "password", "", "User password.")
	cmdCreate.Flags().StringVar(&paramsUserCreate.PasswordConfirmation, "password-confirmation", "", "Optional, but if provided, we will ensure that it matches the value sent in `password`.")
	cmdCreate.Flags().BoolVar(&createAnnouncementsRead, "announcements-read", createAnnouncementsRead, "Signifies that the user has read all the announcements in the UI.")
	cmdCreate.Flags().StringVar(&paramsUserCreate.AllowedIps, "allowed-ips", "", "A list of allowed IPs if applicable.  Newline delimited")
	cmdCreate.Flags().BoolVar(&createAttachmentsPermission, "attachments-permission", createAttachmentsPermission, "DEPRECATED: Can the user create Bundles (aka Share Links)? Use the bundle permission instead.")
	lib.TimeVar(cmdCreate.Flags(), &paramsUserCreate.AuthenticateUntil, "authenticate-until")
	cmdCreate.Flags().StringVar(&UserCreateAuthenticationMethod, "authentication-method", "", fmt.Sprintf("How is this user authenticated? %v", reflect.ValueOf(paramsUserCreate.AuthenticationMethod.Enum()).MapKeys()))
	cmdCreate.Flags().BoolVar(&createBillingPermission, "billing-permission", createBillingPermission, "Allow this user to perform operations on the account, payments, and invoices?")
	cmdCreate.Flags().BoolVar(&createBypassInactiveDisable, "bypass-inactive-disable", createBypassInactiveDisable, "Exempt this user from being disabled based on inactivity?")
	cmdCreate.Flags().BoolVar(&createBypassSiteAllowedIps, "bypass-site-allowed-ips", createBypassSiteAllowedIps, "Allow this user to skip site-wide IP blacklists?")
	cmdCreate.Flags().BoolVar(&createDavPermission, "dav-permission", createDavPermission, "Can the user connect with WebDAV?")
	cmdCreate.Flags().BoolVar(&createDisabled, "disabled", createDisabled, "Is user disabled? Disabled users cannot log in, and do not count for billing purposes.  Users can be automatically disabled after an inactivity period via a Site setting.")
	cmdCreate.Flags().BoolVar(&createFtpPermission, "ftp-permission", createFtpPermission, "Can the user access with FTP/FTPS?")
	cmdCreate.Flags().StringVar(&paramsUserCreate.HeaderText, "header-text", "", "Text to display to the user in the header of the UI")
	cmdCreate.Flags().StringVar(&paramsUserCreate.Language, "language", "", "Preferred language")
	cmdCreate.Flags().Int64Var(&paramsUserCreate.NotificationDailySendTime, "notification-daily-send-time", 0, "Hour of the day at which daily notifications should be sent. Can be in range 0 to 23")
	cmdCreate.Flags().StringVar(&paramsUserCreate.Name, "name", "", "User's full name")
	cmdCreate.Flags().StringVar(&paramsUserCreate.Company, "company", "", "User's company")
	cmdCreate.Flags().StringVar(&paramsUserCreate.Notes, "notes", "", "Any internal notes on the user")
	cmdCreate.Flags().BoolVar(&createOfficeIntegrationEnabled, "office-integration-enabled", createOfficeIntegrationEnabled, "Enable integration with Office for the web?")
	cmdCreate.Flags().Int64Var(&paramsUserCreate.PasswordValidityDays, "password-validity-days", 0, "Number of days to allow user to use the same password")
	cmdCreate.Flags().BoolVar(&createReceiveAdminAlerts, "receive-admin-alerts", createReceiveAdminAlerts, "Should the user receive admin alerts such a certificate expiration notifications and overages?")
	cmdCreate.Flags().BoolVar(&createRequirePasswordChange, "require-password-change", createRequirePasswordChange, "Is a password change required upon next user login?")
	cmdCreate.Flags().BoolVar(&createRestapiPermission, "restapi-permission", createRestapiPermission, "Can this user access the REST API?")
	cmdCreate.Flags().BoolVar(&createSelfManaged, "self-managed", createSelfManaged, "Does this user manage it's own credentials or is it a shared/bot user?")
	cmdCreate.Flags().BoolVar(&createSftpPermission, "sftp-permission", createSftpPermission, "Can the user access with SFTP?")
	cmdCreate.Flags().BoolVar(&createSiteAdmin, "site-admin", createSiteAdmin, "Is the user an administrator for this site?")
	cmdCreate.Flags().BoolVar(&createSkipWelcomeScreen, "skip-welcome-screen", createSkipWelcomeScreen, "Skip Welcome page in the UI?")
	cmdCreate.Flags().StringVar(&UserCreateSslRequired, "ssl-required", "", fmt.Sprintf("SSL required setting %v", reflect.ValueOf(paramsUserCreate.SslRequired.Enum()).MapKeys()))
	cmdCreate.Flags().Int64Var(&paramsUserCreate.SsoStrategyId, "sso-strategy-id", 0, "SSO (Single Sign On) strategy ID for the user, if applicable.")
	cmdCreate.Flags().BoolVar(&createSubscribeToNewsletter, "subscribe-to-newsletter", createSubscribeToNewsletter, "Is the user subscribed to the newsletter?")
	cmdCreate.Flags().StringVar(&UserCreateRequire2fa, "require-2fa", "", fmt.Sprintf("2FA required setting %v", reflect.ValueOf(paramsUserCreate.Require2fa.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&paramsUserCreate.TimeZone, "time-zone", "", "User time zone")
	cmdCreate.Flags().StringVar(&paramsUserCreate.UserRoot, "user-root", "", "Root folder for FTP (and optionally SFTP if the appropriate site-wide setting is set.)  Note that this is not used for API, Desktop, or Web interface.")
	cmdCreate.Flags().StringVar(&paramsUserCreate.Username, "username", "", "User's username")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Users.AddCommand(cmdCreate)
	var fieldsUnlock string
	var formatUnlock string
	paramsUserUnlock := files_sdk.UserUnlockParams{}

	cmdUnlock := &cobra.Command{
		Use: "unlock",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := user.Client{Config: *config}

			result, err := client.Unlock(ctx, paramsUserUnlock)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatUnlock, fieldsUnlock)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdUnlock.Flags().Int64Var(&paramsUserUnlock.Id, "id", 0, "User ID.")

	cmdUnlock.Flags().StringVarP(&fieldsUnlock, "fields", "", "", "comma separated list of field names")
	cmdUnlock.Flags().StringVarP(&formatUnlock, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Users.AddCommand(cmdUnlock)
	var fieldsResendWelcomeEmail string
	var formatResendWelcomeEmail string
	paramsUserResendWelcomeEmail := files_sdk.UserResendWelcomeEmailParams{}

	cmdResendWelcomeEmail := &cobra.Command{
		Use: "resend-welcome-email",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := user.Client{Config: *config}

			result, err := client.ResendWelcomeEmail(ctx, paramsUserResendWelcomeEmail)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatResendWelcomeEmail, fieldsResendWelcomeEmail)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdResendWelcomeEmail.Flags().Int64Var(&paramsUserResendWelcomeEmail.Id, "id", 0, "User ID.")

	cmdResendWelcomeEmail.Flags().StringVarP(&fieldsResendWelcomeEmail, "fields", "", "", "comma separated list of field names")
	cmdResendWelcomeEmail.Flags().StringVarP(&formatResendWelcomeEmail, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Users.AddCommand(cmdResendWelcomeEmail)
	var fieldsUser2faReset string
	var formatUser2faReset string
	paramsUserUser2faReset := files_sdk.UserUser2faResetParams{}

	cmdUser2faReset := &cobra.Command{
		Use: "user-2fa-reset",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := user.Client{Config: *config}

			result, err := client.User2faReset(ctx, paramsUserUser2faReset)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatUser2faReset, fieldsUser2faReset)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdUser2faReset.Flags().Int64Var(&paramsUserUser2faReset.Id, "id", 0, "User ID.")

	cmdUser2faReset.Flags().StringVarP(&fieldsUser2faReset, "fields", "", "", "comma separated list of field names")
	cmdUser2faReset.Flags().StringVarP(&formatUser2faReset, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Users.AddCommand(cmdUser2faReset)
	var fieldsUpdate string
	var formatUpdate string
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
	UserUpdateAuthenticationMethod := ""
	UserUpdateSslRequired := ""
	UserUpdateRequire2fa := ""

	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := user.Client{Config: *config}

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

			paramsUserUpdate.AuthenticationMethod = paramsUserUpdate.AuthenticationMethod.Enum()[UserUpdateAuthenticationMethod]
			paramsUserUpdate.SslRequired = paramsUserUpdate.SslRequired.Enum()[UserUpdateSslRequired]
			paramsUserUpdate.Require2fa = paramsUserUpdate.Require2fa.Enum()[UserUpdateRequire2fa]

			result, err := client.Update(ctx, paramsUserUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatUpdate, fieldsUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsUserUpdate.Id, "id", 0, "User ID.")
	cmdUpdate.Flags().BoolVar(&updateAvatarDelete, "avatar-delete", updateAvatarDelete, "If true, the avatar will be deleted.")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.ChangePassword, "change-password", "", "Used for changing a password on an existing user.")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.ChangePasswordConfirmation, "change-password-confirmation", "", "Optional, but if provided, we will ensure that it matches the value sent in `change_password`.")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.Email, "email", "", "User's email.")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.GrantPermission, "grant-permission", "", "Permission to grant on the user root.  Can be blank or `full`, `read`, `write`, `list`, or `history`.")
	cmdUpdate.Flags().Int64Var(&paramsUserUpdate.GroupId, "group-id", 0, "Group ID to associate this user with.")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.GroupIds, "group-ids", "", "A list of group ids to associate this user with.  Comma delimited.")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.ImportedPasswordHash, "imported-password-hash", "", "Pre-calculated hash of the user's password. If supplied, this will be used to authenticate the user on first login. Supported hash menthods are MD5, SHA1, and SHA256.")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.Password, "password", "", "User password.")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.PasswordConfirmation, "password-confirmation", "", "Optional, but if provided, we will ensure that it matches the value sent in `password`.")
	cmdUpdate.Flags().BoolVar(&updateAnnouncementsRead, "announcements-read", updateAnnouncementsRead, "Signifies that the user has read all the announcements in the UI.")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.AllowedIps, "allowed-ips", "", "A list of allowed IPs if applicable.  Newline delimited")
	cmdUpdate.Flags().BoolVar(&updateAttachmentsPermission, "attachments-permission", updateAttachmentsPermission, "DEPRECATED: Can the user create Bundles (aka Share Links)? Use the bundle permission instead.")
	lib.TimeVar(cmdUpdate.Flags(), &paramsUserUpdate.AuthenticateUntil, "authenticate-until")
	cmdUpdate.Flags().StringVar(&UserUpdateAuthenticationMethod, "authentication-method", "", fmt.Sprintf("How is this user authenticated? %v", reflect.ValueOf(paramsUserUpdate.AuthenticationMethod.Enum()).MapKeys()))
	cmdUpdate.Flags().BoolVar(&updateBillingPermission, "billing-permission", updateBillingPermission, "Allow this user to perform operations on the account, payments, and invoices?")
	cmdUpdate.Flags().BoolVar(&updateBypassInactiveDisable, "bypass-inactive-disable", updateBypassInactiveDisable, "Exempt this user from being disabled based on inactivity?")
	cmdUpdate.Flags().BoolVar(&updateBypassSiteAllowedIps, "bypass-site-allowed-ips", updateBypassSiteAllowedIps, "Allow this user to skip site-wide IP blacklists?")
	cmdUpdate.Flags().BoolVar(&updateDavPermission, "dav-permission", updateDavPermission, "Can the user connect with WebDAV?")
	cmdUpdate.Flags().BoolVar(&updateDisabled, "disabled", updateDisabled, "Is user disabled? Disabled users cannot log in, and do not count for billing purposes.  Users can be automatically disabled after an inactivity period via a Site setting.")
	cmdUpdate.Flags().BoolVar(&updateFtpPermission, "ftp-permission", updateFtpPermission, "Can the user access with FTP/FTPS?")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.HeaderText, "header-text", "", "Text to display to the user in the header of the UI")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.Language, "language", "", "Preferred language")
	cmdUpdate.Flags().Int64Var(&paramsUserUpdate.NotificationDailySendTime, "notification-daily-send-time", 0, "Hour of the day at which daily notifications should be sent. Can be in range 0 to 23")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.Name, "name", "", "User's full name")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.Company, "company", "", "User's company")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.Notes, "notes", "", "Any internal notes on the user")
	cmdUpdate.Flags().BoolVar(&updateOfficeIntegrationEnabled, "office-integration-enabled", updateOfficeIntegrationEnabled, "Enable integration with Office for the web?")
	cmdUpdate.Flags().Int64Var(&paramsUserUpdate.PasswordValidityDays, "password-validity-days", 0, "Number of days to allow user to use the same password")
	cmdUpdate.Flags().BoolVar(&updateReceiveAdminAlerts, "receive-admin-alerts", updateReceiveAdminAlerts, "Should the user receive admin alerts such a certificate expiration notifications and overages?")
	cmdUpdate.Flags().BoolVar(&updateRequirePasswordChange, "require-password-change", updateRequirePasswordChange, "Is a password change required upon next user login?")
	cmdUpdate.Flags().BoolVar(&updateRestapiPermission, "restapi-permission", updateRestapiPermission, "Can this user access the REST API?")
	cmdUpdate.Flags().BoolVar(&updateSelfManaged, "self-managed", updateSelfManaged, "Does this user manage it's own credentials or is it a shared/bot user?")
	cmdUpdate.Flags().BoolVar(&updateSftpPermission, "sftp-permission", updateSftpPermission, "Can the user access with SFTP?")
	cmdUpdate.Flags().BoolVar(&updateSiteAdmin, "site-admin", updateSiteAdmin, "Is the user an administrator for this site?")
	cmdUpdate.Flags().BoolVar(&updateSkipWelcomeScreen, "skip-welcome-screen", updateSkipWelcomeScreen, "Skip Welcome page in the UI?")
	cmdUpdate.Flags().StringVar(&UserUpdateSslRequired, "ssl-required", "", fmt.Sprintf("SSL required setting %v", reflect.ValueOf(paramsUserUpdate.SslRequired.Enum()).MapKeys()))
	cmdUpdate.Flags().Int64Var(&paramsUserUpdate.SsoStrategyId, "sso-strategy-id", 0, "SSO (Single Sign On) strategy ID for the user, if applicable.")
	cmdUpdate.Flags().BoolVar(&updateSubscribeToNewsletter, "subscribe-to-newsletter", updateSubscribeToNewsletter, "Is the user subscribed to the newsletter?")
	cmdUpdate.Flags().StringVar(&UserUpdateRequire2fa, "require-2fa", "", fmt.Sprintf("2FA required setting %v", reflect.ValueOf(paramsUserUpdate.Require2fa.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.TimeZone, "time-zone", "", "User time zone")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.UserRoot, "user-root", "", "Root folder for FTP (and optionally SFTP if the appropriate site-wide setting is set.)  Note that this is not used for API, Desktop, or Web interface.")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.Username, "username", "", "User's username")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Users.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	paramsUserDelete := files_sdk.UserDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := user.Client{Config: *config}

			result, err := client.Delete(ctx, paramsUserDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatDelete, fieldsDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdDelete.Flags().Int64Var(&paramsUserDelete.Id, "id", 0, "User ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Users.AddCommand(cmdDelete)
}
