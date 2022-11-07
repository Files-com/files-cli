package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	flib "github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/Files-com/files-sdk-go/v2/site"
)

func init() {
	RootCmd.AddCommand(Sites())
}

func Sites() *cobra.Command {
	Sites := &cobra.Command{
		Use:  "sites [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command sites\n\t%v", args[0])
		},
	}
	var fieldsGet string
	var formatGet string
	usePagerGet := true
	cmdGet := &cobra.Command{
		Use:   "get",
		Short: `Show site settings`,
		Long:  `Show site settings`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := site.Client{Config: *config}

			var site interface{}
			var err error
			site, err = client.Get(ctx)
			lib.HandleResponse(ctx, Profile(cmd), site, err, formatGet, fieldsGet, usePagerGet, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}

	cmdGet.Flags().StringVar(&fieldsGet, "fields", "", "comma separated list of field names")
	cmdGet.Flags().StringVar(&formatGet, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdGet.Flags().BoolVar(&usePagerGet, "use-pager", usePagerGet, "Use $PAGER (.ie less, more, etc)")

	Sites.AddCommand(cmdGet)
	var fieldsGetUsage string
	var formatGetUsage string
	usePagerGetUsage := true
	cmdGetUsage := &cobra.Command{
		Use:   "get-usage",
		Short: `Get the most recent usage snapshot (usage data for billing purposes) for a Site`,
		Long:  `Get the most recent usage snapshot (usage data for billing purposes) for a Site`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := site.Client{Config: *config}

			var usageSnapshot interface{}
			var err error
			usageSnapshot, err = client.GetUsage(ctx)
			lib.HandleResponse(ctx, Profile(cmd), usageSnapshot, err, formatGetUsage, fieldsGetUsage, usePagerGetUsage, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}

	cmdGetUsage.Flags().StringVar(&fieldsGetUsage, "fields", "", "comma separated list of field names")
	cmdGetUsage.Flags().StringVar(&formatGetUsage, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdGetUsage.Flags().BoolVar(&usePagerGetUsage, "use-pager", usePagerGetUsage, "Use $PAGER (.ie less, more, etc)")

	Sites.AddCommand(cmdGetUsage)
	var fieldsUpdate string
	var formatUpdate string
	usePagerUpdate := true
	updateDomainHstsHeader := true
	updateAllowBundleNames := true
	updateOverageNotify := true
	updateWelcomeEmailEnabled := true
	updateAskAboutOverwrites := true
	updateShowRequestAccessLink := true
	updateWindowsModeFtp := true
	updateDesktopApp := true
	updateDesktopAppSessionIpPinning := true
	updateMobileApp := true
	updateMobileAppSessionIpPinning := true
	updateFolderPermissionsGroupsOnly := true
	updateOfficeIntegrationAvailable := true
	updatePinAllRemoteServersToSiteRegion := true
	updateMotdUseForFtp := true
	updateMotdUseForSftp := true
	updateSslRequired := true
	updateTlsDisabled := true
	updateSftpInsecureCiphers := true
	updateDisableFilesCertificateGeneration := true
	updateUserLockout := true
	updateIncludePasswordInWelcomeEmail := true
	updatePasswordRequireLetter := true
	updatePasswordRequireMixed := true
	updatePasswordRequireSpecial := true
	updatePasswordRequireNumber := true
	updatePasswordRequireUnbreached := true
	updateSftpUserRootEnabled := true
	updateDisablePasswordReset := true
	updateImmutableFiles := true
	updateSessionPinnedByIp := true
	updateBundlePasswordRequired := true
	updateBundleRequireShareRecipient := true
	updatePasswordRequirementsApplyToBundles := true
	updateOptOutGlobal := true
	updateUseProvidedModifiedAt := true
	updateCustomNamespace := true
	updateNonSsoGroupsAllowed := true
	updateNonSsoUsersAllowed := true
	updateSharingEnabled := true
	updateUserRequestsEnabled := true
	updateUserRequestsNotifyAdmins := true
	updateFtpEnabled := true
	updateSftpEnabled := true
	updateAllowed2faMethodSms := true
	updateAllowed2faMethodU2f := true
	updateAllowed2faMethodTotp := true
	updateAllowed2faMethodWebauthn := true
	updateAllowed2faMethodYubi := true
	updateAllowed2faMethodBypassForFtpSftpDav := true
	updateRequire2fa := true
	updateLdapEnabled := true
	updateLdapSecure := true
	updateIcon16Delete := true
	updateIcon32Delete := true
	updateIcon48Delete := true
	updateIcon128Delete := true
	updateLogoDelete := true
	updateBundleWatermarkAttachmentDelete := true
	updateDisable2faWithDelay := true
	paramsSiteUpdate := files_sdk.SiteUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update site settings`,
		Long:  `Update site settings`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := site.Client{Config: *config}

			if cmd.Flags().Changed("domain-hsts-header") {
				paramsSiteUpdate.DomainHstsHeader = flib.Bool(updateDomainHstsHeader)
			}
			if cmd.Flags().Changed("allow-bundle-names") {
				paramsSiteUpdate.AllowBundleNames = flib.Bool(updateAllowBundleNames)
			}
			if cmd.Flags().Changed("overage-notify") {
				paramsSiteUpdate.OverageNotify = flib.Bool(updateOverageNotify)
			}
			if cmd.Flags().Changed("welcome-email-enabled") {
				paramsSiteUpdate.WelcomeEmailEnabled = flib.Bool(updateWelcomeEmailEnabled)
			}
			if cmd.Flags().Changed("ask-about-overwrites") {
				paramsSiteUpdate.AskAboutOverwrites = flib.Bool(updateAskAboutOverwrites)
			}
			if cmd.Flags().Changed("show-request-access-link") {
				paramsSiteUpdate.ShowRequestAccessLink = flib.Bool(updateShowRequestAccessLink)
			}
			if cmd.Flags().Changed("windows-mode-ftp") {
				paramsSiteUpdate.WindowsModeFtp = flib.Bool(updateWindowsModeFtp)
			}
			if cmd.Flags().Changed("desktop-app") {
				paramsSiteUpdate.DesktopApp = flib.Bool(updateDesktopApp)
			}
			if cmd.Flags().Changed("desktop-app-session-ip-pinning") {
				paramsSiteUpdate.DesktopAppSessionIpPinning = flib.Bool(updateDesktopAppSessionIpPinning)
			}
			if cmd.Flags().Changed("mobile-app") {
				paramsSiteUpdate.MobileApp = flib.Bool(updateMobileApp)
			}
			if cmd.Flags().Changed("mobile-app-session-ip-pinning") {
				paramsSiteUpdate.MobileAppSessionIpPinning = flib.Bool(updateMobileAppSessionIpPinning)
			}
			if cmd.Flags().Changed("folder-permissions-groups-only") {
				paramsSiteUpdate.FolderPermissionsGroupsOnly = flib.Bool(updateFolderPermissionsGroupsOnly)
			}
			if cmd.Flags().Changed("office-integration-available") {
				paramsSiteUpdate.OfficeIntegrationAvailable = flib.Bool(updateOfficeIntegrationAvailable)
			}
			if cmd.Flags().Changed("pin-all-remote-servers-to-site-region") {
				paramsSiteUpdate.PinAllRemoteServersToSiteRegion = flib.Bool(updatePinAllRemoteServersToSiteRegion)
			}
			if cmd.Flags().Changed("motd-use-for-ftp") {
				paramsSiteUpdate.MotdUseForFtp = flib.Bool(updateMotdUseForFtp)
			}
			if cmd.Flags().Changed("motd-use-for-sftp") {
				paramsSiteUpdate.MotdUseForSftp = flib.Bool(updateMotdUseForSftp)
			}
			if cmd.Flags().Changed("ssl-required") {
				paramsSiteUpdate.SslRequired = flib.Bool(updateSslRequired)
			}
			if cmd.Flags().Changed("tls-disabled") {
				paramsSiteUpdate.TlsDisabled = flib.Bool(updateTlsDisabled)
			}
			if cmd.Flags().Changed("sftp-insecure-ciphers") {
				paramsSiteUpdate.SftpInsecureCiphers = flib.Bool(updateSftpInsecureCiphers)
			}
			if cmd.Flags().Changed("disable-files-certificate-generation") {
				paramsSiteUpdate.DisableFilesCertificateGeneration = flib.Bool(updateDisableFilesCertificateGeneration)
			}
			if cmd.Flags().Changed("user-lockout") {
				paramsSiteUpdate.UserLockout = flib.Bool(updateUserLockout)
			}
			if cmd.Flags().Changed("include-password-in-welcome-email") {
				paramsSiteUpdate.IncludePasswordInWelcomeEmail = flib.Bool(updateIncludePasswordInWelcomeEmail)
			}
			if cmd.Flags().Changed("password-require-letter") {
				paramsSiteUpdate.PasswordRequireLetter = flib.Bool(updatePasswordRequireLetter)
			}
			if cmd.Flags().Changed("password-require-mixed") {
				paramsSiteUpdate.PasswordRequireMixed = flib.Bool(updatePasswordRequireMixed)
			}
			if cmd.Flags().Changed("password-require-special") {
				paramsSiteUpdate.PasswordRequireSpecial = flib.Bool(updatePasswordRequireSpecial)
			}
			if cmd.Flags().Changed("password-require-number") {
				paramsSiteUpdate.PasswordRequireNumber = flib.Bool(updatePasswordRequireNumber)
			}
			if cmd.Flags().Changed("password-require-unbreached") {
				paramsSiteUpdate.PasswordRequireUnbreached = flib.Bool(updatePasswordRequireUnbreached)
			}
			if cmd.Flags().Changed("sftp-user-root-enabled") {
				paramsSiteUpdate.SftpUserRootEnabled = flib.Bool(updateSftpUserRootEnabled)
			}
			if cmd.Flags().Changed("disable-password-reset") {
				paramsSiteUpdate.DisablePasswordReset = flib.Bool(updateDisablePasswordReset)
			}
			if cmd.Flags().Changed("immutable-files") {
				paramsSiteUpdate.ImmutableFiles = flib.Bool(updateImmutableFiles)
			}
			if cmd.Flags().Changed("session-pinned-by-ip") {
				paramsSiteUpdate.SessionPinnedByIp = flib.Bool(updateSessionPinnedByIp)
			}
			if cmd.Flags().Changed("bundle-password-required") {
				paramsSiteUpdate.BundlePasswordRequired = flib.Bool(updateBundlePasswordRequired)
			}
			if cmd.Flags().Changed("bundle-require-share-recipient") {
				paramsSiteUpdate.BundleRequireShareRecipient = flib.Bool(updateBundleRequireShareRecipient)
			}
			if cmd.Flags().Changed("password-requirements-apply-to-bundles") {
				paramsSiteUpdate.PasswordRequirementsApplyToBundles = flib.Bool(updatePasswordRequirementsApplyToBundles)
			}
			if cmd.Flags().Changed("opt-out-global") {
				paramsSiteUpdate.OptOutGlobal = flib.Bool(updateOptOutGlobal)
			}
			if cmd.Flags().Changed("use-provided-modified-at") {
				paramsSiteUpdate.UseProvidedModifiedAt = flib.Bool(updateUseProvidedModifiedAt)
			}
			if cmd.Flags().Changed("custom-namespace") {
				paramsSiteUpdate.CustomNamespace = flib.Bool(updateCustomNamespace)
			}
			if cmd.Flags().Changed("non-sso-groups-allowed") {
				paramsSiteUpdate.NonSsoGroupsAllowed = flib.Bool(updateNonSsoGroupsAllowed)
			}
			if cmd.Flags().Changed("non-sso-users-allowed") {
				paramsSiteUpdate.NonSsoUsersAllowed = flib.Bool(updateNonSsoUsersAllowed)
			}
			if cmd.Flags().Changed("sharing-enabled") {
				paramsSiteUpdate.SharingEnabled = flib.Bool(updateSharingEnabled)
			}
			if cmd.Flags().Changed("user-requests-enabled") {
				paramsSiteUpdate.UserRequestsEnabled = flib.Bool(updateUserRequestsEnabled)
			}
			if cmd.Flags().Changed("user-requests-notify-admins") {
				paramsSiteUpdate.UserRequestsNotifyAdmins = flib.Bool(updateUserRequestsNotifyAdmins)
			}
			if cmd.Flags().Changed("ftp-enabled") {
				paramsSiteUpdate.FtpEnabled = flib.Bool(updateFtpEnabled)
			}
			if cmd.Flags().Changed("sftp-enabled") {
				paramsSiteUpdate.SftpEnabled = flib.Bool(updateSftpEnabled)
			}
			if cmd.Flags().Changed("allowed-2fa-method-sms") {
				paramsSiteUpdate.Allowed2faMethodSms = flib.Bool(updateAllowed2faMethodSms)
			}
			if cmd.Flags().Changed("allowed-2fa-method-u2f") {
				paramsSiteUpdate.Allowed2faMethodU2f = flib.Bool(updateAllowed2faMethodU2f)
			}
			if cmd.Flags().Changed("allowed-2fa-method-totp") {
				paramsSiteUpdate.Allowed2faMethodTotp = flib.Bool(updateAllowed2faMethodTotp)
			}
			if cmd.Flags().Changed("allowed-2fa-method-webauthn") {
				paramsSiteUpdate.Allowed2faMethodWebauthn = flib.Bool(updateAllowed2faMethodWebauthn)
			}
			if cmd.Flags().Changed("allowed-2fa-method-yubi") {
				paramsSiteUpdate.Allowed2faMethodYubi = flib.Bool(updateAllowed2faMethodYubi)
			}
			if cmd.Flags().Changed("allowed-2fa-method-bypass-for-ftp-sftp-dav") {
				paramsSiteUpdate.Allowed2faMethodBypassForFtpSftpDav = flib.Bool(updateAllowed2faMethodBypassForFtpSftpDav)
			}
			if cmd.Flags().Changed("require-2fa") {
				paramsSiteUpdate.Require2fa = flib.Bool(updateRequire2fa)
			}
			if cmd.Flags().Changed("ldap-enabled") {
				paramsSiteUpdate.LdapEnabled = flib.Bool(updateLdapEnabled)
			}
			if cmd.Flags().Changed("ldap-secure") {
				paramsSiteUpdate.LdapSecure = flib.Bool(updateLdapSecure)
			}
			if cmd.Flags().Changed("icon16-delete") {
				paramsSiteUpdate.Icon16Delete = flib.Bool(updateIcon16Delete)
			}
			if cmd.Flags().Changed("icon32-delete") {
				paramsSiteUpdate.Icon32Delete = flib.Bool(updateIcon32Delete)
			}
			if cmd.Flags().Changed("icon48-delete") {
				paramsSiteUpdate.Icon48Delete = flib.Bool(updateIcon48Delete)
			}
			if cmd.Flags().Changed("icon128-delete") {
				paramsSiteUpdate.Icon128Delete = flib.Bool(updateIcon128Delete)
			}
			if cmd.Flags().Changed("logo-delete") {
				paramsSiteUpdate.LogoDelete = flib.Bool(updateLogoDelete)
			}
			if cmd.Flags().Changed("bundle-watermark-attachment-delete") {
				paramsSiteUpdate.BundleWatermarkAttachmentDelete = flib.Bool(updateBundleWatermarkAttachmentDelete)
			}
			if cmd.Flags().Changed("disable-2fa-with-delay") {
				paramsSiteUpdate.Disable2faWithDelay = flib.Bool(updateDisable2faWithDelay)
			}

			var site interface{}
			var err error
			site, err = client.Update(ctx, paramsSiteUpdate)
			lib.HandleResponse(ctx, Profile(cmd), site, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
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
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.WelcomeEmailSubject, "welcome-email-subject", "", "Include this email subject in welcome emails if enabled")
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
	cmdUpdate.Flags().BoolVar(&updatePinAllRemoteServersToSiteRegion, "pin-all-remote-servers-to-site-region", updatePinAllRemoteServersToSiteRegion, "If true, we will ensure that all internal communications with any remote server are made through the primary region of the site. This setting overrides individual remote server settings.")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.MotdText, "motd-text", "", "A message to show users when they connect via FTP or SFTP.")
	cmdUpdate.Flags().BoolVar(&updateMotdUseForFtp, "motd-use-for-ftp", updateMotdUseForFtp, "Show message to users connecting via FTP")
	cmdUpdate.Flags().BoolVar(&updateMotdUseForSftp, "motd-use-for-sftp", updateMotdUseForSftp, "Show message to users connecting via SFTP")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.SessionExpiry, "session-expiry", "", "Session expiry in hours")
	cmdUpdate.Flags().BoolVar(&updateSslRequired, "ssl-required", updateSslRequired, "Is SSL required?  Disabling this is insecure.")
	cmdUpdate.Flags().BoolVar(&updateTlsDisabled, "tls-disabled", updateTlsDisabled, "Are Insecure TLS and SFTP Ciphers allowed?  Enabling this is insecure.")
	cmdUpdate.Flags().BoolVar(&updateSftpInsecureCiphers, "sftp-insecure-ciphers", updateSftpInsecureCiphers, "Are Insecure Ciphers allowed for SFTP?  Note:  Settting TLS Disabled -> True will always allow insecure ciphers for SFTP as well.  Enabling this is insecure.")
	cmdUpdate.Flags().BoolVar(&updateDisableFilesCertificateGeneration, "disable-files-certificate-generation", updateDisableFilesCertificateGeneration, "If set, Files.com will not set the CAA records required to generate future SSL certificates for this domain.")
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
	cmdUpdate.Flags().BoolVar(&updateUserRequestsNotifyAdmins, "user-requests-notify-admins", updateUserRequestsNotifyAdmins, "Send email to site admins when a user request is received?")
	cmdUpdate.Flags().BoolVar(&updateFtpEnabled, "ftp-enabled", updateFtpEnabled, "Is FTP enabled?")
	cmdUpdate.Flags().BoolVar(&updateSftpEnabled, "sftp-enabled", updateSftpEnabled, "Is SFTP enabled?")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.SftpHostKeyType, "sftp-host-key-type", "", "Sftp Host Key Type")
	cmdUpdate.Flags().Int64Var(&paramsSiteUpdate.ActiveSftpHostKeyId, "active-sftp-host-key-id", 0, "Id of the currently selected custom SFTP Host Key")
	cmdUpdate.Flags().BoolVar(&updateAllowed2faMethodSms, "allowed-2fa-method-sms", updateAllowed2faMethodSms, "Is SMS two factor authentication allowed?")
	cmdUpdate.Flags().BoolVar(&updateAllowed2faMethodU2f, "allowed-2fa-method-u2f", updateAllowed2faMethodU2f, "Is U2F two factor authentication allowed?")
	cmdUpdate.Flags().BoolVar(&updateAllowed2faMethodTotp, "allowed-2fa-method-totp", updateAllowed2faMethodTotp, "Is TOTP two factor authentication allowed?")
	cmdUpdate.Flags().BoolVar(&updateAllowed2faMethodWebauthn, "allowed-2fa-method-webauthn", updateAllowed2faMethodWebauthn, "Is WebAuthn two factor authentication allowed?")
	cmdUpdate.Flags().BoolVar(&updateAllowed2faMethodYubi, "allowed-2fa-method-yubi", updateAllowed2faMethodYubi, "Is yubikey two factor authentication allowed?")
	cmdUpdate.Flags().BoolVar(&updateAllowed2faMethodBypassForFtpSftpDav, "allowed-2fa-method-bypass-for-ftp-sftp-dav", updateAllowed2faMethodBypassForFtpSftpDav, "Are users allowed to configure their two factor authentication to be bypassed for FTP/SFTP/WebDAV?")
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

	cmdUpdate.Flags().StringVar(&fieldsUpdate, "fields", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVar(&formatUpdate, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	Sites.AddCommand(cmdUpdate)
	return Sites
}
