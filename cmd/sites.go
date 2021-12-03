package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	flib "github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/Files-com/files-sdk-go/v2/site"
)

var (
	Sites = &cobra.Command{}
)

func SitesInit() {
	Sites = &cobra.Command{
		Use:  "sites [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command sites\n\t%v", args[0])
		},
	}
	var fieldsGet string
	var formatGet string
	cmdGet := &cobra.Command{
		Use: "get",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := site.Client{Config: *config}

			result, err := client.Get(ctx)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatGet, fieldsGet, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}

	cmdGet.Flags().StringVarP(&fieldsGet, "fields", "", "", "comma separated list of field names")
	cmdGet.Flags().StringVarP(&formatGet, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Sites.AddCommand(cmdGet)
	var fieldsGetUsage string
	var formatGetUsage string
	cmdGetUsage := &cobra.Command{
		Use: "get-usage",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := site.Client{Config: *config}

			result, err := client.GetUsage(ctx)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatGetUsage, fieldsGetUsage, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}

	cmdGetUsage.Flags().StringVarP(&fieldsGetUsage, "fields", "", "", "comma separated list of field names")
	cmdGetUsage.Flags().StringVarP(&formatGetUsage, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Sites.AddCommand(cmdGetUsage)
	var fieldsUpdate string
	var formatUpdate string
	updateDomainHstsHeader := false
	updateAllowBundleNames := false
	updateOverageNotify := false
	updateWelcomeEmailEnabled := false
	updateAskAboutOverwrites := false
	updateShowRequestAccessLink := false
	updateWindowsModeFtp := false
	updateDesktopApp := false
	updateDesktopAppSessionIpPinning := false
	updateMobileApp := false
	updateMobileAppSessionIpPinning := false
	updateFolderPermissionsGroupsOnly := false
	updateOfficeIntegrationAvailable := false
	updateSslRequired := false
	updateTlsDisabled := false
	updateSftpInsecureCiphers := false
	updateUserLockout := false
	updateIncludePasswordInWelcomeEmail := false
	updatePasswordRequireLetter := false
	updatePasswordRequireMixed := false
	updatePasswordRequireSpecial := false
	updatePasswordRequireNumber := false
	updatePasswordRequireUnbreached := false
	updateSftpUserRootEnabled := false
	updateDisablePasswordReset := false
	updateImmutableFiles := false
	updateSessionPinnedByIp := false
	updateBundlePasswordRequired := false
	updateBundleRequireShareRecipient := false
	updatePasswordRequirementsApplyToBundles := false
	updateOptOutGlobal := false
	updateUseProvidedModifiedAt := false
	updateCustomNamespace := false
	updateNonSsoGroupsAllowed := false
	updateNonSsoUsersAllowed := false
	updateSharingEnabled := false
	updateUserRequestsEnabled := false
	updateFtpEnabled := false
	updateSftpEnabled := false
	updateAllowed2faMethodSms := false
	updateAllowed2faMethodU2f := false
	updateAllowed2faMethodTotp := false
	updateAllowed2faMethodYubi := false
	updateRequire2fa := false
	updateLdapEnabled := false
	updateLdapSecure := false
	updateIcon16Delete := false
	updateIcon32Delete := false
	updateIcon48Delete := false
	updateIcon128Delete := false
	updateLogoDelete := false
	updateBundleWatermarkAttachmentDelete := false
	updateDisable2faWithDelay := false
	paramsSiteUpdate := files_sdk.SiteUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := site.Client{Config: *config}

			if updateDomainHstsHeader {
				paramsSiteUpdate.DomainHstsHeader = flib.Bool(true)
			}
			if updateAllowBundleNames {
				paramsSiteUpdate.AllowBundleNames = flib.Bool(true)
			}
			if updateOverageNotify {
				paramsSiteUpdate.OverageNotify = flib.Bool(true)
			}
			if updateWelcomeEmailEnabled {
				paramsSiteUpdate.WelcomeEmailEnabled = flib.Bool(true)
			}
			if updateAskAboutOverwrites {
				paramsSiteUpdate.AskAboutOverwrites = flib.Bool(true)
			}
			if updateShowRequestAccessLink {
				paramsSiteUpdate.ShowRequestAccessLink = flib.Bool(true)
			}
			if updateWindowsModeFtp {
				paramsSiteUpdate.WindowsModeFtp = flib.Bool(true)
			}
			if updateDesktopApp {
				paramsSiteUpdate.DesktopApp = flib.Bool(true)
			}
			if updateDesktopAppSessionIpPinning {
				paramsSiteUpdate.DesktopAppSessionIpPinning = flib.Bool(true)
			}
			if updateMobileApp {
				paramsSiteUpdate.MobileApp = flib.Bool(true)
			}
			if updateMobileAppSessionIpPinning {
				paramsSiteUpdate.MobileAppSessionIpPinning = flib.Bool(true)
			}
			if updateFolderPermissionsGroupsOnly {
				paramsSiteUpdate.FolderPermissionsGroupsOnly = flib.Bool(true)
			}
			if updateOfficeIntegrationAvailable {
				paramsSiteUpdate.OfficeIntegrationAvailable = flib.Bool(true)
			}
			if updateSslRequired {
				paramsSiteUpdate.SslRequired = flib.Bool(true)
			}
			if updateTlsDisabled {
				paramsSiteUpdate.TlsDisabled = flib.Bool(true)
			}
			if updateSftpInsecureCiphers {
				paramsSiteUpdate.SftpInsecureCiphers = flib.Bool(true)
			}
			if updateUserLockout {
				paramsSiteUpdate.UserLockout = flib.Bool(true)
			}
			if updateIncludePasswordInWelcomeEmail {
				paramsSiteUpdate.IncludePasswordInWelcomeEmail = flib.Bool(true)
			}
			if updatePasswordRequireLetter {
				paramsSiteUpdate.PasswordRequireLetter = flib.Bool(true)
			}
			if updatePasswordRequireMixed {
				paramsSiteUpdate.PasswordRequireMixed = flib.Bool(true)
			}
			if updatePasswordRequireSpecial {
				paramsSiteUpdate.PasswordRequireSpecial = flib.Bool(true)
			}
			if updatePasswordRequireNumber {
				paramsSiteUpdate.PasswordRequireNumber = flib.Bool(true)
			}
			if updatePasswordRequireUnbreached {
				paramsSiteUpdate.PasswordRequireUnbreached = flib.Bool(true)
			}
			if updateSftpUserRootEnabled {
				paramsSiteUpdate.SftpUserRootEnabled = flib.Bool(true)
			}
			if updateDisablePasswordReset {
				paramsSiteUpdate.DisablePasswordReset = flib.Bool(true)
			}
			if updateImmutableFiles {
				paramsSiteUpdate.ImmutableFiles = flib.Bool(true)
			}
			if updateSessionPinnedByIp {
				paramsSiteUpdate.SessionPinnedByIp = flib.Bool(true)
			}
			if updateBundlePasswordRequired {
				paramsSiteUpdate.BundlePasswordRequired = flib.Bool(true)
			}
			if updateBundleRequireShareRecipient {
				paramsSiteUpdate.BundleRequireShareRecipient = flib.Bool(true)
			}
			if updatePasswordRequirementsApplyToBundles {
				paramsSiteUpdate.PasswordRequirementsApplyToBundles = flib.Bool(true)
			}
			if updateOptOutGlobal {
				paramsSiteUpdate.OptOutGlobal = flib.Bool(true)
			}
			if updateUseProvidedModifiedAt {
				paramsSiteUpdate.UseProvidedModifiedAt = flib.Bool(true)
			}
			if updateCustomNamespace {
				paramsSiteUpdate.CustomNamespace = flib.Bool(true)
			}
			if updateNonSsoGroupsAllowed {
				paramsSiteUpdate.NonSsoGroupsAllowed = flib.Bool(true)
			}
			if updateNonSsoUsersAllowed {
				paramsSiteUpdate.NonSsoUsersAllowed = flib.Bool(true)
			}
			if updateSharingEnabled {
				paramsSiteUpdate.SharingEnabled = flib.Bool(true)
			}
			if updateUserRequestsEnabled {
				paramsSiteUpdate.UserRequestsEnabled = flib.Bool(true)
			}
			if updateFtpEnabled {
				paramsSiteUpdate.FtpEnabled = flib.Bool(true)
			}
			if updateSftpEnabled {
				paramsSiteUpdate.SftpEnabled = flib.Bool(true)
			}
			if updateAllowed2faMethodSms {
				paramsSiteUpdate.Allowed2faMethodSms = flib.Bool(true)
			}
			if updateAllowed2faMethodU2f {
				paramsSiteUpdate.Allowed2faMethodU2f = flib.Bool(true)
			}
			if updateAllowed2faMethodTotp {
				paramsSiteUpdate.Allowed2faMethodTotp = flib.Bool(true)
			}
			if updateAllowed2faMethodYubi {
				paramsSiteUpdate.Allowed2faMethodYubi = flib.Bool(true)
			}
			if updateRequire2fa {
				paramsSiteUpdate.Require2fa = flib.Bool(true)
			}
			if updateLdapEnabled {
				paramsSiteUpdate.LdapEnabled = flib.Bool(true)
			}
			if updateLdapSecure {
				paramsSiteUpdate.LdapSecure = flib.Bool(true)
			}
			if updateIcon16Delete {
				paramsSiteUpdate.Icon16Delete = flib.Bool(true)
			}
			if updateIcon32Delete {
				paramsSiteUpdate.Icon32Delete = flib.Bool(true)
			}
			if updateIcon48Delete {
				paramsSiteUpdate.Icon48Delete = flib.Bool(true)
			}
			if updateIcon128Delete {
				paramsSiteUpdate.Icon128Delete = flib.Bool(true)
			}
			if updateLogoDelete {
				paramsSiteUpdate.LogoDelete = flib.Bool(true)
			}
			if updateBundleWatermarkAttachmentDelete {
				paramsSiteUpdate.BundleWatermarkAttachmentDelete = flib.Bool(true)
			}
			if updateDisable2faWithDelay {
				paramsSiteUpdate.Disable2faWithDelay = flib.Bool(true)
			}

			result, err := client.Update(ctx, paramsSiteUpdate)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatUpdate, fieldsUpdate, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.Name, "name", "", "Site name")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.Subdomain, "subdomain", "", "Site subdomain")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.Domain, "domain", "", "Custom domain")
	cmdUpdate.Flags().BoolVar(&updateDomainHstsHeader, "domain-hsts-header", updateDomainHstsHeader, "Send HSTS (HTTP Strict Transport Security) header when visitors access the site via a custom domain?")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.DomainLetsencryptChain, "domain-letsencrypt-chain", "", "Letsencrypt chain to use when registering SSL Certificate for domain.")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.Email, "email", "", "Main email for this site")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.ReplyToEmail, "reply-to-email", "", "Reply-to email for this site")
	cmdUpdate.Flags().BoolVar(&updateAllowBundleNames, "allow-bundle-names", updateAllowBundleNames, "Are manual Bundle names allowed?")
	cmdUpdate.Flags().Int64Var(&paramsSiteUpdate.BundleExpiration, "bundle-expiration", 0, "Site-wide Bundle expiration in days")
	cmdUpdate.Flags().BoolVar(&updateOverageNotify, "overage-notify", updateOverageNotify, "Notify site email of overages?")
	cmdUpdate.Flags().BoolVar(&updateWelcomeEmailEnabled, "welcome-email-enabled", updateWelcomeEmailEnabled, "Will the welcome email be sent to new users?")
	cmdUpdate.Flags().BoolVar(&updateAskAboutOverwrites, "ask-about-overwrites", updateAskAboutOverwrites, "If false, rename conflicting files instead of asking for overwrite confirmation.  Only applies to web interface.")
	cmdUpdate.Flags().BoolVar(&updateShowRequestAccessLink, "show-request-access-link", updateShowRequestAccessLink, "Show request access link for users without access?  Currently unused.")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.WelcomeEmailCc, "welcome-email-cc", "", "Include this email in welcome emails if enabled")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.WelcomeCustomText, "welcome-custom-text", "", "Custom text send in user welcome email")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.Language, "language", "", "Site default language")
	cmdUpdate.Flags().BoolVar(&updateWindowsModeFtp, "windows-mode-ftp", updateWindowsModeFtp, "Does FTP user Windows emulation mode?")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.DefaultTimeZone, "default-time-zone", "", "Site default time zone")
	cmdUpdate.Flags().BoolVar(&updateDesktopApp, "desktop-app", updateDesktopApp, "Is the desktop app enabled?")
	cmdUpdate.Flags().BoolVar(&updateDesktopAppSessionIpPinning, "desktop-app-session-ip-pinning", updateDesktopAppSessionIpPinning, "Is desktop app session IP pinning enabled?")
	cmdUpdate.Flags().Int64Var(&paramsSiteUpdate.DesktopAppSessionLifetime, "desktop-app-session-lifetime", 0, "Desktop app session lifetime (in hours)")
	cmdUpdate.Flags().BoolVar(&updateMobileApp, "mobile-app", updateMobileApp, "Is the mobile app enabled?")
	cmdUpdate.Flags().BoolVar(&updateMobileAppSessionIpPinning, "mobile-app-session-ip-pinning", updateMobileAppSessionIpPinning, "Is mobile app session IP pinning enabled?")
	cmdUpdate.Flags().Int64Var(&paramsSiteUpdate.MobileAppSessionLifetime, "mobile-app-session-lifetime", 0, "Mobile app session lifetime (in hours)")
	cmdUpdate.Flags().BoolVar(&updateFolderPermissionsGroupsOnly, "folder-permissions-groups-only", updateFolderPermissionsGroupsOnly, "If true, permissions for this site must be bound to a group (not a user). Otherwise, permissions must be bound to a user.")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.WelcomeScreen, "welcome-screen", "", "Does the welcome screen appear?")
	cmdUpdate.Flags().BoolVar(&updateOfficeIntegrationAvailable, "office-integration-available", updateOfficeIntegrationAvailable, "Allow users to use Office for the web?")
	cmdUpdate.Flags().BoolVar(&updateSslRequired, "ssl-required", updateSslRequired, "Is SSL required?  Disabling this is insecure.")
	cmdUpdate.Flags().BoolVar(&updateTlsDisabled, "tls-disabled", updateTlsDisabled, "Are Insecure TLS and SFTP Ciphers allowed?  Enabling this is insecure.")
	cmdUpdate.Flags().BoolVar(&updateSftpInsecureCiphers, "sftp-insecure-ciphers", updateSftpInsecureCiphers, "Are Insecure Ciphers allowed for SFTP?  Note:  Settting TLS Disabled -> True will always allow insecure ciphers for SFTP as well.  Enabling this is insecure.")
	cmdUpdate.Flags().BoolVar(&updateUserLockout, "user-lockout", updateUserLockout, "Will users be locked out after incorrect login attempts?")
	cmdUpdate.Flags().Int64Var(&paramsSiteUpdate.UserLockoutTries, "user-lockout-tries", 0, "Number of login tries within `user_lockout_within` hours before users are locked out")
	cmdUpdate.Flags().Int64Var(&paramsSiteUpdate.UserLockoutWithin, "user-lockout-within", 0, "Number of hours for user lockout window")
	cmdUpdate.Flags().Int64Var(&paramsSiteUpdate.UserLockoutLockPeriod, "user-lockout-lock-period", 0, "How many hours to lock user out for failed password?")
	cmdUpdate.Flags().BoolVar(&updateIncludePasswordInWelcomeEmail, "include-password-in-welcome-email", updateIncludePasswordInWelcomeEmail, "Include password in emails to new users?")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.AllowedCountries, "allowed-countries", "", "Comma seperated list of allowed Country codes")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.AllowedIps, "allowed-ips", "", "List of allowed IP addresses")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.DisallowedCountries, "disallowed-countries", "", "Comma seperated list of disallowed Country codes")
	cmdUpdate.Flags().Int64Var(&paramsSiteUpdate.DaysToRetainBackups, "days-to-retain-backups", 0, "Number of days to keep deleted files")
	cmdUpdate.Flags().Int64Var(&paramsSiteUpdate.MaxPriorPasswords, "max-prior-passwords", 0, "Number of prior passwords to disallow")
	cmdUpdate.Flags().Int64Var(&paramsSiteUpdate.PasswordValidityDays, "password-validity-days", 0, "Number of days password is valid")
	cmdUpdate.Flags().Int64Var(&paramsSiteUpdate.PasswordMinLength, "password-min-length", 0, "Shortest password length for users")
	cmdUpdate.Flags().BoolVar(&updatePasswordRequireLetter, "password-require-letter", updatePasswordRequireLetter, "Require a letter in passwords?")
	cmdUpdate.Flags().BoolVar(&updatePasswordRequireMixed, "password-require-mixed", updatePasswordRequireMixed, "Require lower and upper case letters in passwords?")
	cmdUpdate.Flags().BoolVar(&updatePasswordRequireSpecial, "password-require-special", updatePasswordRequireSpecial, "Require special characters in password?")
	cmdUpdate.Flags().BoolVar(&updatePasswordRequireNumber, "password-require-number", updatePasswordRequireNumber, "Require a number in passwords?")
	cmdUpdate.Flags().BoolVar(&updatePasswordRequireUnbreached, "password-require-unbreached", updatePasswordRequireUnbreached, "Require passwords that have not been previously breached? (see https://haveibeenpwned.com/)")
	cmdUpdate.Flags().BoolVar(&updateSftpUserRootEnabled, "sftp-user-root-enabled", updateSftpUserRootEnabled, "Use user FTP roots also for SFTP?")
	cmdUpdate.Flags().BoolVar(&updateDisablePasswordReset, "disable-password-reset", updateDisablePasswordReset, "Is password reset disabled?")
	cmdUpdate.Flags().BoolVar(&updateImmutableFiles, "immutable-files", updateImmutableFiles, "Are files protected from modification?")
	cmdUpdate.Flags().BoolVar(&updateSessionPinnedByIp, "session-pinned-by-ip", updateSessionPinnedByIp, "Are sessions locked to the same IP? (i.e. do users need to log in again if they change IPs?)")
	cmdUpdate.Flags().BoolVar(&updateBundlePasswordRequired, "bundle-password-required", updateBundlePasswordRequired, "Do Bundles require password protection?")
	cmdUpdate.Flags().BoolVar(&updateBundleRequireShareRecipient, "bundle-require-share-recipient", updateBundleRequireShareRecipient, "Do Bundles require recipients for sharing?")
	cmdUpdate.Flags().BoolVar(&updatePasswordRequirementsApplyToBundles, "password-requirements-apply-to-bundles", updatePasswordRequirementsApplyToBundles, "Require bundles' passwords, and passwords for other items (inboxes, public shares, etc.) to conform to the same requirements as users' passwords?")
	cmdUpdate.Flags().BoolVar(&updateOptOutGlobal, "opt-out-global", updateOptOutGlobal, "Use servers in the USA only?")
	cmdUpdate.Flags().BoolVar(&updateUseProvidedModifiedAt, "use-provided-modified-at", updateUseProvidedModifiedAt, "Allow uploaders to set `provided_modified_at` for uploaded files?")
	cmdUpdate.Flags().BoolVar(&updateCustomNamespace, "custom-namespace", updateCustomNamespace, "Is this site using a custom namespace for users?")
	cmdUpdate.Flags().Int64Var(&paramsSiteUpdate.DisableUsersFromInactivityPeriodDays, "disable-users-from-inactivity-period-days", 0, "If greater than zero, users will unable to login if they do not show activity within this number of days.")
	cmdUpdate.Flags().BoolVar(&updateNonSsoGroupsAllowed, "non-sso-groups-allowed", updateNonSsoGroupsAllowed, "If true, groups can be manually created / modified / deleted by Site Admins. Otherwise, groups can only be managed via your SSO provider.")
	cmdUpdate.Flags().BoolVar(&updateNonSsoUsersAllowed, "non-sso-users-allowed", updateNonSsoUsersAllowed, "If true, users can be manually created / modified / deleted by Site Admins. Otherwise, users can only be managed via your SSO provider.")
	cmdUpdate.Flags().BoolVar(&updateSharingEnabled, "sharing-enabled", updateSharingEnabled, "Allow bundle creation")
	cmdUpdate.Flags().BoolVar(&updateUserRequestsEnabled, "user-requests-enabled", updateUserRequestsEnabled, "Enable User Requests feature")
	cmdUpdate.Flags().BoolVar(&updateFtpEnabled, "ftp-enabled", updateFtpEnabled, "Is FTP enabled?")
	cmdUpdate.Flags().BoolVar(&updateSftpEnabled, "sftp-enabled", updateSftpEnabled, "Is SFTP enabled?")
	cmdUpdate.Flags().BoolVar(&updateAllowed2faMethodSms, "allowed-2fa-method-sms", updateAllowed2faMethodSms, "Is SMS two factor authentication allowed?")
	cmdUpdate.Flags().BoolVar(&updateAllowed2faMethodU2f, "allowed-2fa-method-u2f", updateAllowed2faMethodU2f, "Is U2F two factor authentication allowed?")
	cmdUpdate.Flags().BoolVar(&updateAllowed2faMethodTotp, "allowed-2fa-method-totp", updateAllowed2faMethodTotp, "Is TOTP two factor authentication allowed?")
	cmdUpdate.Flags().BoolVar(&updateAllowed2faMethodYubi, "allowed-2fa-method-yubi", updateAllowed2faMethodYubi, "Is yubikey two factor authentication allowed?")
	cmdUpdate.Flags().BoolVar(&updateRequire2fa, "require-2fa", updateRequire2fa, "Require two-factor authentication for all users?")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.Require2faUserType, "require-2fa-user-type", "", "What type of user is required to use two-factor authentication (when require_2fa is set to `true` for this site)?")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.Color2Top, "color2-top", "", "Top bar background color")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.Color2Left, "color2-left", "", "Page link and button color")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.Color2Link, "color2-link", "", "Top bar link color")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.Color2Text, "color2-text", "", "Page link and button color")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.Color2TopText, "color2-top-text", "", "Top bar text color")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.SiteHeader, "site-header", "", "Custom site header text")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.SiteFooter, "site-footer", "", "Custom site footer text")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.LoginHelpText, "login-help-text", "", "Login help text")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.SmtpAddress, "smtp-address", "", "SMTP server hostname or IP")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.SmtpAuthentication, "smtp-authentication", "", "SMTP server authentication type")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.SmtpFrom, "smtp-from", "", "From address to use when mailing through custom SMTP")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.SmtpUsername, "smtp-username", "", "SMTP server username")
	cmdUpdate.Flags().Int64Var(&paramsSiteUpdate.SmtpPort, "smtp-port", 0, "SMTP server port")
	cmdUpdate.Flags().BoolVar(&updateLdapEnabled, "ldap-enabled", updateLdapEnabled, "Main LDAP setting: is LDAP enabled?")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.LdapType, "ldap-type", "", "LDAP type")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.LdapHost, "ldap-host", "", "LDAP host")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.LdapHost2, "ldap-host-2", "", "LDAP backup host")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.LdapHost3, "ldap-host-3", "", "LDAP backup host")
	cmdUpdate.Flags().Int64Var(&paramsSiteUpdate.LdapPort, "ldap-port", 0, "LDAP port")
	cmdUpdate.Flags().BoolVar(&updateLdapSecure, "ldap-secure", updateLdapSecure, "Use secure LDAP?")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.LdapUsername, "ldap-username", "", "Username for signing in to LDAP server.")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.LdapUsernameField, "ldap-username-field", "", "LDAP username field")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.LdapDomain, "ldap-domain", "", "Domain name that will be appended to usernames")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.LdapUserAction, "ldap-user-action", "", "Should we sync users from LDAP server?")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.LdapGroupAction, "ldap-group-action", "", "Should we sync groups from LDAP server?")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.LdapUserIncludeGroups, "ldap-user-include-groups", "", "Comma or newline separated list of group names (with optional wildcards) - if provided, only users in these groups will be added or synced.")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.LdapGroupExclusion, "ldap-group-exclusion", "", "Comma or newline separated list of group names (with optional wildcards) to exclude when syncing.")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.LdapGroupInclusion, "ldap-group-inclusion", "", "Comma or newline separated list of group names (with optional wildcards) to include when syncing.")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.LdapBaseDn, "ldap-base-dn", "", "Base DN for looking up users in LDAP server")
	cmdUpdate.Flags().BoolVar(&updateIcon16Delete, "icon16-delete", updateIcon16Delete, "If true, will delete the file stored in icon16")
	cmdUpdate.Flags().BoolVar(&updateIcon32Delete, "icon32-delete", updateIcon32Delete, "If true, will delete the file stored in icon32")
	cmdUpdate.Flags().BoolVar(&updateIcon48Delete, "icon48-delete", updateIcon48Delete, "If true, will delete the file stored in icon48")
	cmdUpdate.Flags().BoolVar(&updateIcon128Delete, "icon128-delete", updateIcon128Delete, "If true, will delete the file stored in icon128")
	cmdUpdate.Flags().BoolVar(&updateLogoDelete, "logo-delete", updateLogoDelete, "If true, will delete the file stored in logo")
	cmdUpdate.Flags().BoolVar(&updateBundleWatermarkAttachmentDelete, "bundle-watermark-attachment-delete", updateBundleWatermarkAttachmentDelete, "If true, will delete the file stored in bundle_watermark_attachment")
	cmdUpdate.Flags().BoolVar(&updateDisable2faWithDelay, "disable-2fa-with-delay", updateDisable2faWithDelay, "If set to true, we will begin the process of disabling 2FA on this site.")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.LdapPasswordChange, "ldap-password-change", "", "New LDAP password.")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.LdapPasswordChangeConfirmation, "ldap-password-change-confirmation", "", "Confirm new LDAP password.")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.SmtpPassword, "smtp-password", "", "Password for SMTP server.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Sites.AddCommand(cmdUpdate)
}
