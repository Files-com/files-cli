package cmd

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/user"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Users())
}

func Users() *cobra.Command {
	Users := &cobra.Command{
		Use:  "users [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command users\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsUserList := files_sdk.UserListParams{}
	var MaxPagesList int64
	listIncludeParentSiteUsers := true

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Users",
		Long:    `List Users`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsUserList
			params.MaxPages = MaxPagesList

			if cmd.Flags().Changed("include-parent-site-users") {
				params.IncludeParentSiteUsers = flib.Bool(listIncludeParentSiteUsers)
			}

			client := user.Client{Config: config}
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

	cmdList.Flags().StringVar(&paramsUserList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsUserList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVar(&paramsUserList.Ids, "ids", "", "comma-separated list of User IDs")
	cmdList.Flags().BoolVar(&listIncludeParentSiteUsers, "include-parent-site-users", listIncludeParentSiteUsers, "Include users from the parent site.")
	cmdList.Flags().StringVar(&paramsUserList.Search, "search", "", "Searches for partial matches of name, username, or email.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	Users.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsUserFind := files_sdk.UserFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show User`,
		Long:  `Show User`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := user.Client{Config: config}

			var user interface{}
			var err error
			user, err = client.Find(paramsUserFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), user, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsUserFind.Id, "id", 0, "User ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	Users.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createAvatarDelete := true
	createAnnouncementsRead := true
	createAttachmentsPermission := true
	createBillingPermission := true
	createBypassUserLifecycleRules := true
	createBypassSiteAllowedIps := true
	createDavPermission := true
	createDisabled := true
	createFtpPermission := true
	createOfficeIntegrationEnabled := true
	createPartnerAdmin := true
	createReadonlySiteAdmin := true
	createReceiveAdminAlerts := true
	createRequirePasswordChange := true
	createRestapiPermission := true
	createSelfManaged := true
	createSftpPermission := true
	createSiteAdmin := true
	createSkipWelcomeScreen := true
	createSubscribeToNewsletter := true
	createWorkspaceAdmin := true
	paramsUserCreate := files_sdk.UserCreateParams{}
	UserCreateAuthenticationMethod := ""
	UserCreateFilesystemLayout := ""
	UserCreateSslRequired := ""
	UserCreateRequire2fa := ""

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create User`,
		Long:  `Create User`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := user.Client{Config: config}

			var UserCreateAuthenticationMethodErr error
			paramsUserCreate.AuthenticationMethod, UserCreateAuthenticationMethodErr = lib.FetchKey("authentication-method", paramsUserCreate.AuthenticationMethod.Enum(), UserCreateAuthenticationMethod)
			if UserCreateAuthenticationMethod != "" && UserCreateAuthenticationMethodErr != nil {
				return UserCreateAuthenticationMethodErr
			}
			var UserCreateFilesystemLayoutErr error
			paramsUserCreate.FilesystemLayout, UserCreateFilesystemLayoutErr = lib.FetchKey("filesystem-layout", paramsUserCreate.FilesystemLayout.Enum(), UserCreateFilesystemLayout)
			if UserCreateFilesystemLayout != "" && UserCreateFilesystemLayoutErr != nil {
				return UserCreateFilesystemLayoutErr
			}
			var UserCreateSslRequiredErr error
			paramsUserCreate.SslRequired, UserCreateSslRequiredErr = lib.FetchKey("ssl-required", paramsUserCreate.SslRequired.Enum(), UserCreateSslRequired)
			if UserCreateSslRequired != "" && UserCreateSslRequiredErr != nil {
				return UserCreateSslRequiredErr
			}
			var UserCreateRequire2faErr error
			paramsUserCreate.Require2fa, UserCreateRequire2faErr = lib.FetchKey("require-2fa", paramsUserCreate.Require2fa.Enum(), UserCreateRequire2fa)
			if UserCreateRequire2fa != "" && UserCreateRequire2faErr != nil {
				return UserCreateRequire2faErr
			}

			if cmd.Flags().Changed("avatar-delete") {
				paramsUserCreate.AvatarDelete = flib.Bool(createAvatarDelete)
			}
			if cmd.Flags().Changed("announcements-read") {
				paramsUserCreate.AnnouncementsRead = flib.Bool(createAnnouncementsRead)
			}
			if cmd.Flags().Changed("attachments-permission") {
				paramsUserCreate.AttachmentsPermission = flib.Bool(createAttachmentsPermission)
			}
			if cmd.Flags().Changed("billing-permission") {
				paramsUserCreate.BillingPermission = flib.Bool(createBillingPermission)
			}
			if cmd.Flags().Changed("bypass-user-lifecycle-rules") {
				paramsUserCreate.BypassUserLifecycleRules = flib.Bool(createBypassUserLifecycleRules)
			}
			if cmd.Flags().Changed("bypass-site-allowed-ips") {
				paramsUserCreate.BypassSiteAllowedIps = flib.Bool(createBypassSiteAllowedIps)
			}
			if cmd.Flags().Changed("dav-permission") {
				paramsUserCreate.DavPermission = flib.Bool(createDavPermission)
			}
			if cmd.Flags().Changed("disabled") {
				paramsUserCreate.Disabled = flib.Bool(createDisabled)
			}
			if cmd.Flags().Changed("ftp-permission") {
				paramsUserCreate.FtpPermission = flib.Bool(createFtpPermission)
			}
			if cmd.Flags().Changed("office-integration-enabled") {
				paramsUserCreate.OfficeIntegrationEnabled = flib.Bool(createOfficeIntegrationEnabled)
			}
			if cmd.Flags().Changed("partner-admin") {
				paramsUserCreate.PartnerAdmin = flib.Bool(createPartnerAdmin)
			}
			if cmd.Flags().Changed("readonly-site-admin") {
				paramsUserCreate.ReadonlySiteAdmin = flib.Bool(createReadonlySiteAdmin)
			}
			if cmd.Flags().Changed("receive-admin-alerts") {
				paramsUserCreate.ReceiveAdminAlerts = flib.Bool(createReceiveAdminAlerts)
			}
			if cmd.Flags().Changed("require-password-change") {
				paramsUserCreate.RequirePasswordChange = flib.Bool(createRequirePasswordChange)
			}
			if cmd.Flags().Changed("restapi-permission") {
				paramsUserCreate.RestapiPermission = flib.Bool(createRestapiPermission)
			}
			if cmd.Flags().Changed("self-managed") {
				paramsUserCreate.SelfManaged = flib.Bool(createSelfManaged)
			}
			if cmd.Flags().Changed("sftp-permission") {
				paramsUserCreate.SftpPermission = flib.Bool(createSftpPermission)
			}
			if cmd.Flags().Changed("site-admin") {
				paramsUserCreate.SiteAdmin = flib.Bool(createSiteAdmin)
			}
			if cmd.Flags().Changed("skip-welcome-screen") {
				paramsUserCreate.SkipWelcomeScreen = flib.Bool(createSkipWelcomeScreen)
			}
			if cmd.Flags().Changed("subscribe-to-newsletter") {
				paramsUserCreate.SubscribeToNewsletter = flib.Bool(createSubscribeToNewsletter)
			}
			if cmd.Flags().Changed("workspace-admin") {
				paramsUserCreate.WorkspaceAdmin = flib.Bool(createWorkspaceAdmin)
			}

			if paramsUserCreate.AuthenticateUntil.IsZero() {
				paramsUserCreate.AuthenticateUntil = nil
			}
			if paramsUserCreate.RequireLoginBy.IsZero() {
				paramsUserCreate.RequireLoginBy = nil
			}

			var user interface{}
			var err error
			user, err = client.Create(paramsUserCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), user, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().BoolVar(&createAvatarDelete, "avatar-delete", createAvatarDelete, "If true, the avatar will be deleted.")
	cmdCreate.Flags().StringVar(&paramsUserCreate.ChangePassword, "change-password", "", "Used for changing a password on an existing user.")
	cmdCreate.Flags().StringVar(&paramsUserCreate.ChangePasswordConfirmation, "change-password-confirmation", "", "Optional, but if provided, we will ensure that it matches the value sent in `change_password`.")
	cmdCreate.Flags().StringVar(&paramsUserCreate.Email, "email", "", "User's email.")
	cmdCreate.Flags().StringVar(&paramsUserCreate.GrantPermission, "grant-permission", "", "Permission to grant on the User Root upon user creation. Can be blank or `full`, `read`, `write`, `list`, `read+write`, or `list+write`")
	cmdCreate.Flags().Int64Var(&paramsUserCreate.GroupId, "group-id", 0, "Group ID to associate this user with.")
	cmdCreate.Flags().StringVar(&paramsUserCreate.GroupIds, "group-ids", "", "A list of group ids to associate this user with.  Comma delimited.")
	cmdCreate.Flags().StringVar(&paramsUserCreate.ImportedPasswordHash, "imported-password-hash", "", "Pre-calculated hash of the user's password. If supplied, this will be used to authenticate the user on first login. Supported hash methods are MD5, SHA1, and SHA256.")
	cmdCreate.Flags().StringVar(&paramsUserCreate.Password, "password", "", "User password.")
	cmdCreate.Flags().StringVar(&paramsUserCreate.PasswordConfirmation, "password-confirmation", "", "Optional, but if provided, we will ensure that it matches the value sent in `password`.")
	cmdCreate.Flags().BoolVar(&createAnnouncementsRead, "announcements-read", createAnnouncementsRead, "Signifies that the user has read all the announcements in the UI.")
	cmdCreate.Flags().StringVar(&paramsUserCreate.AllowedIps, "allowed-ips", "", "A list of allowed IPs if applicable.  Newline delimited")
	cmdCreate.Flags().BoolVar(&createAttachmentsPermission, "attachments-permission", createAttachmentsPermission, "DEPRECATED: If `true`, the user can user create Bundles (aka Share Links). Use the bundle permission instead.")
	paramsUserCreate.AuthenticateUntil = &time.Time{}
	lib.TimeVar(cmdCreate.Flags(), paramsUserCreate.AuthenticateUntil, "authenticate-until", "Scheduled Date/Time at which user will be deactivated")
	cmdCreate.Flags().StringVar(&UserCreateAuthenticationMethod, "authentication-method", "", fmt.Sprintf("How is this user authenticated? %v", reflect.ValueOf(paramsUserCreate.AuthenticationMethod.Enum()).MapKeys()))
	cmdCreate.Flags().BoolVar(&createBillingPermission, "billing-permission", createBillingPermission, "Allow this user to perform operations on the account, payments, and invoices?")
	cmdCreate.Flags().BoolVar(&createBypassUserLifecycleRules, "bypass-user-lifecycle-rules", createBypassUserLifecycleRules, "Exempt this user from user lifecycle rules?")
	cmdCreate.Flags().BoolVar(&createBypassSiteAllowedIps, "bypass-site-allowed-ips", createBypassSiteAllowedIps, "Allow this user to skip site-wide IP blacklists?")
	cmdCreate.Flags().BoolVar(&createDavPermission, "dav-permission", createDavPermission, "Can the user connect with WebDAV?")
	cmdCreate.Flags().BoolVar(&createDisabled, "disabled", createDisabled, "Is user disabled? Disabled users cannot log in, and do not count for billing purposes. Users can be automatically disabled after an inactivity period via a Site setting or schedule to be deactivated after specific date.")
	cmdCreate.Flags().StringVar(&UserCreateFilesystemLayout, "filesystem-layout", "", fmt.Sprintf("File system layout %v", reflect.ValueOf(paramsUserCreate.FilesystemLayout.Enum()).MapKeys()))
	cmdCreate.Flags().BoolVar(&createFtpPermission, "ftp-permission", createFtpPermission, "Can the user access with FTP/FTPS?")
	cmdCreate.Flags().StringVar(&paramsUserCreate.HeaderText, "header-text", "", "Text to display to the user in the header of the UI")
	cmdCreate.Flags().StringVar(&paramsUserCreate.Language, "language", "", "Preferred language")
	cmdCreate.Flags().Int64Var(&paramsUserCreate.NotificationDailySendTime, "notification-daily-send-time", 0, "Hour of the day at which daily notifications should be sent. Can be in range 0 to 23")
	cmdCreate.Flags().StringVar(&paramsUserCreate.Name, "name", "", "User's full name")
	cmdCreate.Flags().StringVar(&paramsUserCreate.Company, "company", "", "User's company")
	cmdCreate.Flags().StringVar(&paramsUserCreate.Notes, "notes", "", "Any internal notes on the user")
	cmdCreate.Flags().BoolVar(&createOfficeIntegrationEnabled, "office-integration-enabled", createOfficeIntegrationEnabled, "Enable integration with Office for the web?")
	cmdCreate.Flags().BoolVar(&createPartnerAdmin, "partner-admin", createPartnerAdmin, "Is this user a Partner administrator?")
	cmdCreate.Flags().Int64Var(&paramsUserCreate.PartnerId, "partner-id", 0, "Partner ID if this user belongs to a Partner")
	cmdCreate.Flags().Int64Var(&paramsUserCreate.PasswordValidityDays, "password-validity-days", 0, "Number of days to allow user to use the same password")
	cmdCreate.Flags().BoolVar(&createReadonlySiteAdmin, "readonly-site-admin", createReadonlySiteAdmin, "Is the user an allowed to view all (non-billing) site configuration for this site?")
	cmdCreate.Flags().BoolVar(&createReceiveAdminAlerts, "receive-admin-alerts", createReceiveAdminAlerts, "Should the user receive admin alerts such a certificate expiration notifications and overages?")
	paramsUserCreate.RequireLoginBy = &time.Time{}
	lib.TimeVar(cmdCreate.Flags(), paramsUserCreate.RequireLoginBy, "require-login-by", "Require user to login by specified date otherwise it will be disabled.")
	cmdCreate.Flags().BoolVar(&createRequirePasswordChange, "require-password-change", createRequirePasswordChange, "Is a password change required upon next user login?")
	cmdCreate.Flags().BoolVar(&createRestapiPermission, "restapi-permission", createRestapiPermission, "Can this user access the Web app, Desktop app, SDKs, or REST API?  (All of these tools use the API internally, so this is one unified permission set.)")
	cmdCreate.Flags().BoolVar(&createSelfManaged, "self-managed", createSelfManaged, "Does this user manage it's own credentials or is it a shared/bot user?")
	cmdCreate.Flags().BoolVar(&createSftpPermission, "sftp-permission", createSftpPermission, "Can the user access with SFTP?")
	cmdCreate.Flags().BoolVar(&createSiteAdmin, "site-admin", createSiteAdmin, "Is the user an administrator for this site?")
	cmdCreate.Flags().BoolVar(&createSkipWelcomeScreen, "skip-welcome-screen", createSkipWelcomeScreen, "Skip Welcome page in the UI?")
	cmdCreate.Flags().StringVar(&UserCreateSslRequired, "ssl-required", "", fmt.Sprintf("SSL required setting %v", reflect.ValueOf(paramsUserCreate.SslRequired.Enum()).MapKeys()))
	cmdCreate.Flags().Int64Var(&paramsUserCreate.SsoStrategyId, "sso-strategy-id", 0, "SSO (Single Sign On) strategy ID for the user, if applicable.")
	cmdCreate.Flags().BoolVar(&createSubscribeToNewsletter, "subscribe-to-newsletter", createSubscribeToNewsletter, "Is the user subscribed to the newsletter?")
	cmdCreate.Flags().StringVar(&UserCreateRequire2fa, "require-2fa", "", fmt.Sprintf("2FA required setting %v", reflect.ValueOf(paramsUserCreate.Require2fa.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&paramsUserCreate.Tags, "tags", "", "Comma-separated list of Tags for this user. Tags are used for other features, such as UserLifecycleRules, which can target specific tags.  Tags must only contain lowercase letters, numbers, and hyphens.")
	cmdCreate.Flags().StringVar(&paramsUserCreate.TimeZone, "time-zone", "", "User time zone")
	cmdCreate.Flags().StringVar(&paramsUserCreate.UserRoot, "user-root", "", "Root folder for FTP (and optionally SFTP if the appropriate site-wide setting is set).  Note that this is not used for API, Desktop, or Web interface.")
	cmdCreate.Flags().StringVar(&paramsUserCreate.UserHome, "user-home", "", "Home folder for FTP/SFTP.  Note that this is not used for API, Desktop, or Web interface.")
	cmdCreate.Flags().BoolVar(&createWorkspaceAdmin, "workspace-admin", createWorkspaceAdmin, "Is the user a Workspace administrator?  Applicable only to the workspace ID related to this user, if one is set.")
	cmdCreate.Flags().StringVar(&paramsUserCreate.Username, "username", "", "User's username")
	cmdCreate.Flags().Int64Var(&paramsUserCreate.WorkspaceId, "workspace-id", 0, "Workspace ID")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Users.AddCommand(cmdCreate)
	var fieldsUnlock []string
	var formatUnlock []string
	usePagerUnlock := true
	paramsUserUnlock := files_sdk.UserUnlockParams{}

	cmdUnlock := &cobra.Command{
		Use:   "unlock",
		Short: `Unlock user who has been locked out due to failed logins`,
		Long:  `Unlock user who has been locked out due to failed logins`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := user.Client{Config: config}

			var err error
			err = client.Unlock(paramsUserUnlock, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdUnlock.Flags().Int64Var(&paramsUserUnlock.Id, "id", 0, "User ID.")

	cmdUnlock.Flags().StringSliceVar(&fieldsUnlock, "fields", []string{}, "comma separated list of field names")
	cmdUnlock.Flags().StringSliceVar(&formatUnlock, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUnlock.Flags().BoolVar(&usePagerUnlock, "use-pager", usePagerUnlock, "Use $PAGER (.ie less, more, etc)")

	Users.AddCommand(cmdUnlock)
	var fieldsResendWelcomeEmail []string
	var formatResendWelcomeEmail []string
	usePagerResendWelcomeEmail := true
	paramsUserResendWelcomeEmail := files_sdk.UserResendWelcomeEmailParams{}

	cmdResendWelcomeEmail := &cobra.Command{
		Use:   "resend-welcome-email",
		Short: `Resend user welcome email`,
		Long:  `Resend user welcome email`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := user.Client{Config: config}

			var err error
			err = client.ResendWelcomeEmail(paramsUserResendWelcomeEmail, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdResendWelcomeEmail.Flags().Int64Var(&paramsUserResendWelcomeEmail.Id, "id", 0, "User ID.")

	cmdResendWelcomeEmail.Flags().StringSliceVar(&fieldsResendWelcomeEmail, "fields", []string{}, "comma separated list of field names")
	cmdResendWelcomeEmail.Flags().StringSliceVar(&formatResendWelcomeEmail, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdResendWelcomeEmail.Flags().BoolVar(&usePagerResendWelcomeEmail, "use-pager", usePagerResendWelcomeEmail, "Use $PAGER (.ie less, more, etc)")

	Users.AddCommand(cmdResendWelcomeEmail)
	var fieldsUser2faReset []string
	var formatUser2faReset []string
	usePagerUser2faReset := true
	paramsUserUser2faReset := files_sdk.UserUser2faResetParams{}

	cmdUser2faReset := &cobra.Command{
		Use:   "user-2fa-reset",
		Short: `Trigger 2FA Reset process for user who has lost access to their existing 2FA methods`,
		Long:  `Trigger 2FA Reset process for user who has lost access to their existing 2FA methods`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := user.Client{Config: config}

			var err error
			err = client.User2faReset(paramsUserUser2faReset, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdUser2faReset.Flags().Int64Var(&paramsUserUser2faReset.Id, "id", 0, "User ID.")

	cmdUser2faReset.Flags().StringSliceVar(&fieldsUser2faReset, "fields", []string{}, "comma separated list of field names")
	cmdUser2faReset.Flags().StringSliceVar(&formatUser2faReset, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUser2faReset.Flags().BoolVar(&usePagerUser2faReset, "use-pager", usePagerUser2faReset, "Use $PAGER (.ie less, more, etc)")

	Users.AddCommand(cmdUser2faReset)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateAvatarDelete := true
	updateAnnouncementsRead := true
	updateAttachmentsPermission := true
	updateBillingPermission := true
	updateBypassUserLifecycleRules := true
	updateBypassSiteAllowedIps := true
	updateDavPermission := true
	updateDisabled := true
	updateFtpPermission := true
	updateOfficeIntegrationEnabled := true
	updatePartnerAdmin := true
	updateReadonlySiteAdmin := true
	updateReceiveAdminAlerts := true
	updateRequirePasswordChange := true
	updateRestapiPermission := true
	updateSelfManaged := true
	updateSftpPermission := true
	updateSiteAdmin := true
	updateSkipWelcomeScreen := true
	updateSubscribeToNewsletter := true
	updateWorkspaceAdmin := true
	updateClear2fa := true
	updateConvertToPartnerUser := true
	paramsUserUpdate := files_sdk.UserUpdateParams{}
	UserUpdateAuthenticationMethod := ""
	UserUpdateFilesystemLayout := ""
	UserUpdateSslRequired := ""
	UserUpdateRequire2fa := ""

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update User`,
		Long:  `Update User`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := user.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.UserUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			var UserUpdateAuthenticationMethodErr error
			paramsUserUpdate.AuthenticationMethod, UserUpdateAuthenticationMethodErr = lib.FetchKey("authentication-method", paramsUserUpdate.AuthenticationMethod.Enum(), UserUpdateAuthenticationMethod)
			if UserUpdateAuthenticationMethod != "" && UserUpdateAuthenticationMethodErr != nil {
				return UserUpdateAuthenticationMethodErr
			}
			var UserUpdateFilesystemLayoutErr error
			paramsUserUpdate.FilesystemLayout, UserUpdateFilesystemLayoutErr = lib.FetchKey("filesystem-layout", paramsUserUpdate.FilesystemLayout.Enum(), UserUpdateFilesystemLayout)
			if UserUpdateFilesystemLayout != "" && UserUpdateFilesystemLayoutErr != nil {
				return UserUpdateFilesystemLayoutErr
			}
			var UserUpdateSslRequiredErr error
			paramsUserUpdate.SslRequired, UserUpdateSslRequiredErr = lib.FetchKey("ssl-required", paramsUserUpdate.SslRequired.Enum(), UserUpdateSslRequired)
			if UserUpdateSslRequired != "" && UserUpdateSslRequiredErr != nil {
				return UserUpdateSslRequiredErr
			}
			var UserUpdateRequire2faErr error
			paramsUserUpdate.Require2fa, UserUpdateRequire2faErr = lib.FetchKey("require-2fa", paramsUserUpdate.Require2fa.Enum(), UserUpdateRequire2fa)
			if UserUpdateRequire2fa != "" && UserUpdateRequire2faErr != nil {
				return UserUpdateRequire2faErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsUserUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("avatar-file") {
			}
			if cmd.Flags().Changed("avatar-delete") {
				mapParams["avatar_delete"] = updateAvatarDelete
			}
			if cmd.Flags().Changed("change-password") {
				lib.FlagUpdate(cmd, "change_password", paramsUserUpdate.ChangePassword, mapParams)
			}
			if cmd.Flags().Changed("change-password-confirmation") {
				lib.FlagUpdate(cmd, "change_password_confirmation", paramsUserUpdate.ChangePasswordConfirmation, mapParams)
			}
			if cmd.Flags().Changed("email") {
				lib.FlagUpdate(cmd, "email", paramsUserUpdate.Email, mapParams)
			}
			if cmd.Flags().Changed("grant-permission") {
				lib.FlagUpdate(cmd, "grant_permission", paramsUserUpdate.GrantPermission, mapParams)
			}
			if cmd.Flags().Changed("group-id") {
				lib.FlagUpdate(cmd, "group_id", paramsUserUpdate.GroupId, mapParams)
			}
			if cmd.Flags().Changed("group-ids") {
				lib.FlagUpdate(cmd, "group_ids", paramsUserUpdate.GroupIds, mapParams)
			}
			if cmd.Flags().Changed("imported-password-hash") {
				lib.FlagUpdate(cmd, "imported_password_hash", paramsUserUpdate.ImportedPasswordHash, mapParams)
			}
			if cmd.Flags().Changed("password") {
				lib.FlagUpdate(cmd, "password", paramsUserUpdate.Password, mapParams)
			}
			if cmd.Flags().Changed("password-confirmation") {
				lib.FlagUpdate(cmd, "password_confirmation", paramsUserUpdate.PasswordConfirmation, mapParams)
			}
			if cmd.Flags().Changed("announcements-read") {
				mapParams["announcements_read"] = updateAnnouncementsRead
			}
			if cmd.Flags().Changed("allowed-ips") {
				lib.FlagUpdate(cmd, "allowed_ips", paramsUserUpdate.AllowedIps, mapParams)
			}
			if cmd.Flags().Changed("attachments-permission") {
				mapParams["attachments_permission"] = updateAttachmentsPermission
			}
			if cmd.Flags().Changed("authenticate-until") {
				lib.FlagUpdate(cmd, "authenticate_until", paramsUserUpdate.AuthenticateUntil, mapParams)
			}
			if cmd.Flags().Changed("authentication-method") {
				lib.FlagUpdate(cmd, "authentication_method", paramsUserUpdate.AuthenticationMethod, mapParams)
			}
			if cmd.Flags().Changed("billing-permission") {
				mapParams["billing_permission"] = updateBillingPermission
			}
			if cmd.Flags().Changed("bypass-user-lifecycle-rules") {
				mapParams["bypass_user_lifecycle_rules"] = updateBypassUserLifecycleRules
			}
			if cmd.Flags().Changed("bypass-site-allowed-ips") {
				mapParams["bypass_site_allowed_ips"] = updateBypassSiteAllowedIps
			}
			if cmd.Flags().Changed("dav-permission") {
				mapParams["dav_permission"] = updateDavPermission
			}
			if cmd.Flags().Changed("disabled") {
				mapParams["disabled"] = updateDisabled
			}
			if cmd.Flags().Changed("filesystem-layout") {
				lib.FlagUpdate(cmd, "filesystem_layout", paramsUserUpdate.FilesystemLayout, mapParams)
			}
			if cmd.Flags().Changed("ftp-permission") {
				mapParams["ftp_permission"] = updateFtpPermission
			}
			if cmd.Flags().Changed("header-text") {
				lib.FlagUpdate(cmd, "header_text", paramsUserUpdate.HeaderText, mapParams)
			}
			if cmd.Flags().Changed("language") {
				lib.FlagUpdate(cmd, "language", paramsUserUpdate.Language, mapParams)
			}
			if cmd.Flags().Changed("notification-daily-send-time") {
				lib.FlagUpdate(cmd, "notification_daily_send_time", paramsUserUpdate.NotificationDailySendTime, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsUserUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("company") {
				lib.FlagUpdate(cmd, "company", paramsUserUpdate.Company, mapParams)
			}
			if cmd.Flags().Changed("notes") {
				lib.FlagUpdate(cmd, "notes", paramsUserUpdate.Notes, mapParams)
			}
			if cmd.Flags().Changed("office-integration-enabled") {
				mapParams["office_integration_enabled"] = updateOfficeIntegrationEnabled
			}
			if cmd.Flags().Changed("partner-admin") {
				mapParams["partner_admin"] = updatePartnerAdmin
			}
			if cmd.Flags().Changed("partner-id") {
				lib.FlagUpdate(cmd, "partner_id", paramsUserUpdate.PartnerId, mapParams)
			}
			if cmd.Flags().Changed("password-validity-days") {
				lib.FlagUpdate(cmd, "password_validity_days", paramsUserUpdate.PasswordValidityDays, mapParams)
			}
			if cmd.Flags().Changed("readonly-site-admin") {
				mapParams["readonly_site_admin"] = updateReadonlySiteAdmin
			}
			if cmd.Flags().Changed("receive-admin-alerts") {
				mapParams["receive_admin_alerts"] = updateReceiveAdminAlerts
			}
			if cmd.Flags().Changed("require-login-by") {
				lib.FlagUpdate(cmd, "require_login_by", paramsUserUpdate.RequireLoginBy, mapParams)
			}
			if cmd.Flags().Changed("require-password-change") {
				mapParams["require_password_change"] = updateRequirePasswordChange
			}
			if cmd.Flags().Changed("restapi-permission") {
				mapParams["restapi_permission"] = updateRestapiPermission
			}
			if cmd.Flags().Changed("self-managed") {
				mapParams["self_managed"] = updateSelfManaged
			}
			if cmd.Flags().Changed("sftp-permission") {
				mapParams["sftp_permission"] = updateSftpPermission
			}
			if cmd.Flags().Changed("site-admin") {
				mapParams["site_admin"] = updateSiteAdmin
			}
			if cmd.Flags().Changed("skip-welcome-screen") {
				mapParams["skip_welcome_screen"] = updateSkipWelcomeScreen
			}
			if cmd.Flags().Changed("ssl-required") {
				lib.FlagUpdate(cmd, "ssl_required", paramsUserUpdate.SslRequired, mapParams)
			}
			if cmd.Flags().Changed("sso-strategy-id") {
				lib.FlagUpdate(cmd, "sso_strategy_id", paramsUserUpdate.SsoStrategyId, mapParams)
			}
			if cmd.Flags().Changed("subscribe-to-newsletter") {
				mapParams["subscribe_to_newsletter"] = updateSubscribeToNewsletter
			}
			if cmd.Flags().Changed("require-2fa") {
				lib.FlagUpdate(cmd, "require_2fa", paramsUserUpdate.Require2fa, mapParams)
			}
			if cmd.Flags().Changed("tags") {
				lib.FlagUpdate(cmd, "tags", paramsUserUpdate.Tags, mapParams)
			}
			if cmd.Flags().Changed("time-zone") {
				lib.FlagUpdate(cmd, "time_zone", paramsUserUpdate.TimeZone, mapParams)
			}
			if cmd.Flags().Changed("user-root") {
				lib.FlagUpdate(cmd, "user_root", paramsUserUpdate.UserRoot, mapParams)
			}
			if cmd.Flags().Changed("user-home") {
				lib.FlagUpdate(cmd, "user_home", paramsUserUpdate.UserHome, mapParams)
			}
			if cmd.Flags().Changed("workspace-admin") {
				mapParams["workspace_admin"] = updateWorkspaceAdmin
			}
			if cmd.Flags().Changed("username") {
				lib.FlagUpdate(cmd, "username", paramsUserUpdate.Username, mapParams)
			}
			if cmd.Flags().Changed("clear-2fa") {
				mapParams["clear_2fa"] = updateClear2fa
			}
			if cmd.Flags().Changed("convert-to-partner-user") {
				mapParams["convert_to_partner_user"] = updateConvertToPartnerUser
			}

			if paramsUserUpdate.AuthenticateUntil.IsZero() {
				paramsUserUpdate.AuthenticateUntil = nil
			}
			if paramsUserUpdate.RequireLoginBy.IsZero() {
				paramsUserUpdate.RequireLoginBy = nil
			}

			var user interface{}
			var err error
			user, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), user, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsUserUpdate.Id, "id", 0, "User ID.")
	cmdUpdate.Flags().BoolVar(&updateAvatarDelete, "avatar-delete", updateAvatarDelete, "If true, the avatar will be deleted.")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.ChangePassword, "change-password", "", "Used for changing a password on an existing user.")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.ChangePasswordConfirmation, "change-password-confirmation", "", "Optional, but if provided, we will ensure that it matches the value sent in `change_password`.")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.Email, "email", "", "User's email.")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.GrantPermission, "grant-permission", "", "Permission to grant on the User Root upon user creation. Can be blank or `full`, `read`, `write`, `list`, `read+write`, or `list+write`")
	cmdUpdate.Flags().Int64Var(&paramsUserUpdate.GroupId, "group-id", 0, "Group ID to associate this user with.")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.GroupIds, "group-ids", "", "A list of group ids to associate this user with.  Comma delimited.")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.ImportedPasswordHash, "imported-password-hash", "", "Pre-calculated hash of the user's password. If supplied, this will be used to authenticate the user on first login. Supported hash methods are MD5, SHA1, and SHA256.")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.Password, "password", "", "User password.")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.PasswordConfirmation, "password-confirmation", "", "Optional, but if provided, we will ensure that it matches the value sent in `password`.")
	cmdUpdate.Flags().BoolVar(&updateAnnouncementsRead, "announcements-read", updateAnnouncementsRead, "Signifies that the user has read all the announcements in the UI.")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.AllowedIps, "allowed-ips", "", "A list of allowed IPs if applicable.  Newline delimited")
	cmdUpdate.Flags().BoolVar(&updateAttachmentsPermission, "attachments-permission", updateAttachmentsPermission, "DEPRECATED: If `true`, the user can user create Bundles (aka Share Links). Use the bundle permission instead.")
	paramsUserUpdate.AuthenticateUntil = &time.Time{}
	lib.TimeVar(cmdUpdate.Flags(), paramsUserUpdate.AuthenticateUntil, "authenticate-until", "Scheduled Date/Time at which user will be deactivated")
	cmdUpdate.Flags().StringVar(&UserUpdateAuthenticationMethod, "authentication-method", "", fmt.Sprintf("How is this user authenticated? %v", reflect.ValueOf(paramsUserUpdate.AuthenticationMethod.Enum()).MapKeys()))
	cmdUpdate.Flags().BoolVar(&updateBillingPermission, "billing-permission", updateBillingPermission, "Allow this user to perform operations on the account, payments, and invoices?")
	cmdUpdate.Flags().BoolVar(&updateBypassUserLifecycleRules, "bypass-user-lifecycle-rules", updateBypassUserLifecycleRules, "Exempt this user from user lifecycle rules?")
	cmdUpdate.Flags().BoolVar(&updateBypassSiteAllowedIps, "bypass-site-allowed-ips", updateBypassSiteAllowedIps, "Allow this user to skip site-wide IP blacklists?")
	cmdUpdate.Flags().BoolVar(&updateDavPermission, "dav-permission", updateDavPermission, "Can the user connect with WebDAV?")
	cmdUpdate.Flags().BoolVar(&updateDisabled, "disabled", updateDisabled, "Is user disabled? Disabled users cannot log in, and do not count for billing purposes. Users can be automatically disabled after an inactivity period via a Site setting or schedule to be deactivated after specific date.")
	cmdUpdate.Flags().StringVar(&UserUpdateFilesystemLayout, "filesystem-layout", "", fmt.Sprintf("File system layout %v", reflect.ValueOf(paramsUserUpdate.FilesystemLayout.Enum()).MapKeys()))
	cmdUpdate.Flags().BoolVar(&updateFtpPermission, "ftp-permission", updateFtpPermission, "Can the user access with FTP/FTPS?")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.HeaderText, "header-text", "", "Text to display to the user in the header of the UI")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.Language, "language", "", "Preferred language")
	cmdUpdate.Flags().Int64Var(&paramsUserUpdate.NotificationDailySendTime, "notification-daily-send-time", 0, "Hour of the day at which daily notifications should be sent. Can be in range 0 to 23")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.Name, "name", "", "User's full name")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.Company, "company", "", "User's company")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.Notes, "notes", "", "Any internal notes on the user")
	cmdUpdate.Flags().BoolVar(&updateOfficeIntegrationEnabled, "office-integration-enabled", updateOfficeIntegrationEnabled, "Enable integration with Office for the web?")
	cmdUpdate.Flags().BoolVar(&updatePartnerAdmin, "partner-admin", updatePartnerAdmin, "Is this user a Partner administrator?")
	cmdUpdate.Flags().Int64Var(&paramsUserUpdate.PartnerId, "partner-id", 0, "Partner ID if this user belongs to a Partner")
	cmdUpdate.Flags().Int64Var(&paramsUserUpdate.PasswordValidityDays, "password-validity-days", 0, "Number of days to allow user to use the same password")
	cmdUpdate.Flags().BoolVar(&updateReadonlySiteAdmin, "readonly-site-admin", updateReadonlySiteAdmin, "Is the user an allowed to view all (non-billing) site configuration for this site?")
	cmdUpdate.Flags().BoolVar(&updateReceiveAdminAlerts, "receive-admin-alerts", updateReceiveAdminAlerts, "Should the user receive admin alerts such a certificate expiration notifications and overages?")
	paramsUserUpdate.RequireLoginBy = &time.Time{}
	lib.TimeVar(cmdUpdate.Flags(), paramsUserUpdate.RequireLoginBy, "require-login-by", "Require user to login by specified date otherwise it will be disabled.")
	cmdUpdate.Flags().BoolVar(&updateRequirePasswordChange, "require-password-change", updateRequirePasswordChange, "Is a password change required upon next user login?")
	cmdUpdate.Flags().BoolVar(&updateRestapiPermission, "restapi-permission", updateRestapiPermission, "Can this user access the Web app, Desktop app, SDKs, or REST API?  (All of these tools use the API internally, so this is one unified permission set.)")
	cmdUpdate.Flags().BoolVar(&updateSelfManaged, "self-managed", updateSelfManaged, "Does this user manage it's own credentials or is it a shared/bot user?")
	cmdUpdate.Flags().BoolVar(&updateSftpPermission, "sftp-permission", updateSftpPermission, "Can the user access with SFTP?")
	cmdUpdate.Flags().BoolVar(&updateSiteAdmin, "site-admin", updateSiteAdmin, "Is the user an administrator for this site?")
	cmdUpdate.Flags().BoolVar(&updateSkipWelcomeScreen, "skip-welcome-screen", updateSkipWelcomeScreen, "Skip Welcome page in the UI?")
	cmdUpdate.Flags().StringVar(&UserUpdateSslRequired, "ssl-required", "", fmt.Sprintf("SSL required setting %v", reflect.ValueOf(paramsUserUpdate.SslRequired.Enum()).MapKeys()))
	cmdUpdate.Flags().Int64Var(&paramsUserUpdate.SsoStrategyId, "sso-strategy-id", 0, "SSO (Single Sign On) strategy ID for the user, if applicable.")
	cmdUpdate.Flags().BoolVar(&updateSubscribeToNewsletter, "subscribe-to-newsletter", updateSubscribeToNewsletter, "Is the user subscribed to the newsletter?")
	cmdUpdate.Flags().StringVar(&UserUpdateRequire2fa, "require-2fa", "", fmt.Sprintf("2FA required setting %v", reflect.ValueOf(paramsUserUpdate.Require2fa.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.Tags, "tags", "", "Comma-separated list of Tags for this user. Tags are used for other features, such as UserLifecycleRules, which can target specific tags.  Tags must only contain lowercase letters, numbers, and hyphens.")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.TimeZone, "time-zone", "", "User time zone")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.UserRoot, "user-root", "", "Root folder for FTP (and optionally SFTP if the appropriate site-wide setting is set).  Note that this is not used for API, Desktop, or Web interface.")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.UserHome, "user-home", "", "Home folder for FTP/SFTP.  Note that this is not used for API, Desktop, or Web interface.")
	cmdUpdate.Flags().BoolVar(&updateWorkspaceAdmin, "workspace-admin", updateWorkspaceAdmin, "Is the user a Workspace administrator?  Applicable only to the workspace ID related to this user, if one is set.")
	cmdUpdate.Flags().StringVar(&paramsUserUpdate.Username, "username", "", "User's username")
	cmdUpdate.Flags().BoolVar(&updateClear2fa, "clear-2fa", updateClear2fa, "If true when changing authentication_method from `password` to `sso`, remove all two-factor methods. Ignored in all other cases.")
	cmdUpdate.Flags().BoolVar(&updateConvertToPartnerUser, "convert-to-partner-user", updateConvertToPartnerUser, "If true, convert this user to a partner user by assigning the partner_id provided.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	Users.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsUserDelete := files_sdk.UserDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete User`,
		Long:  `Delete User`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := user.Client{Config: config}

			var err error
			err = client.Delete(paramsUserDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsUserDelete.Id, "id", 0, "User ID.")
	cmdDelete.Flags().Int64Var(&paramsUserDelete.NewOwnerId, "new-owner-id", 0, "Provide a User ID here to transfer ownership of certain resources such as Automations and Share Links (Bundles) to that new user.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Users.AddCommand(cmdDelete)
	return Users
}
