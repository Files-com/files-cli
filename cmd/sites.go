package cmd

import (
	"fmt"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/site"
	"github.com/spf13/cobra"
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
	var fieldsGet []string
	var formatGet []string
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
			return lib.HandleResponse(ctx, Profile(cmd), site, err, formatGet, fieldsGet, usePagerGet, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}

	cmdGet.Flags().StringSliceVar(&fieldsGet, "fields", []string{}, "comma separated list of field names")
	cmdGet.Flags().StringSliceVar(&formatGet, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
json-styles: {raw, pretty}`)
	cmdGet.Flags().BoolVar(&usePagerGet, "use-pager", usePagerGet, "Use $PAGER (.ie less, more, etc)")

	Sites.AddCommand(cmdGet)
	var fieldsGetUsage []string
	var formatGetUsage []string
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
			return lib.HandleResponse(ctx, Profile(cmd), usageSnapshot, err, formatGetUsage, fieldsGetUsage, usePagerGetUsage, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}

	cmdGetUsage.Flags().StringSliceVar(&fieldsGetUsage, "fields", []string{}, "comma separated list of field names")
	cmdGetUsage.Flags().StringSliceVar(&formatGetUsage, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
json-styles: {raw, pretty}`)
	cmdGetUsage.Flags().BoolVar(&usePagerGetUsage, "use-pager", usePagerGetUsage, "Use $PAGER (.ie less, more, etc)")

	Sites.AddCommand(cmdGetUsage)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateDomainHstsHeader := true
	updateAllowBundleNames := true
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
	updateGroupAdminsCanSetUserPassword := true
	updateAllowed2faMethodSms := true
	updateAllowed2faMethodU2f := true
	updateAllowed2faMethodTotp := true
	updateAllowed2faMethodWebauthn := true
	updateAllowed2faMethodYubi := true
	updateAllowed2faMethodBypassForFtpSftpDav := true
	updateRequire2fa := true
	updateLdapEnabled := true
	updateLdapSecure := true
	updateUploadsViaEmailAuthentication := true
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

			mapParams, convertErr := lib.StructToMap(files_sdk.SiteUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsSiteUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("subdomain") {
				lib.FlagUpdate(cmd, "subdomain", paramsSiteUpdate.Subdomain, mapParams)
			}
			if cmd.Flags().Changed("domain") {
				lib.FlagUpdate(cmd, "domain", paramsSiteUpdate.Domain, mapParams)
			}
			if cmd.Flags().Changed("domain-hsts-header") {
				mapParams["domain_hsts_header"] = updateDomainHstsHeader
			}
			if cmd.Flags().Changed("domain-letsencrypt-chain") {
				lib.FlagUpdate(cmd, "domain_letsencrypt_chain", paramsSiteUpdate.DomainLetsencryptChain, mapParams)
			}
			if cmd.Flags().Changed("email") {
				lib.FlagUpdate(cmd, "email", paramsSiteUpdate.Email, mapParams)
			}
			if cmd.Flags().Changed("reply-to-email") {
				lib.FlagUpdate(cmd, "reply_to_email", paramsSiteUpdate.ReplyToEmail, mapParams)
			}
			if cmd.Flags().Changed("allow-bundle-names") {
				mapParams["allow_bundle_names"] = updateAllowBundleNames
			}
			if cmd.Flags().Changed("bundle-expiration") {
				lib.FlagUpdate(cmd, "bundle_expiration", paramsSiteUpdate.BundleExpiration, mapParams)
			}
			if cmd.Flags().Changed("welcome-email-enabled") {
				mapParams["welcome_email_enabled"] = updateWelcomeEmailEnabled
			}
			if cmd.Flags().Changed("ask-about-overwrites") {
				mapParams["ask_about_overwrites"] = updateAskAboutOverwrites
			}
			if cmd.Flags().Changed("show-request-access-link") {
				mapParams["show_request_access_link"] = updateShowRequestAccessLink
			}
			if cmd.Flags().Changed("welcome-email-cc") {
				lib.FlagUpdate(cmd, "welcome_email_cc", paramsSiteUpdate.WelcomeEmailCc, mapParams)
			}
			if cmd.Flags().Changed("welcome-email-subject") {
				lib.FlagUpdate(cmd, "welcome_email_subject", paramsSiteUpdate.WelcomeEmailSubject, mapParams)
			}
			if cmd.Flags().Changed("welcome-custom-text") {
				lib.FlagUpdate(cmd, "welcome_custom_text", paramsSiteUpdate.WelcomeCustomText, mapParams)
			}
			if cmd.Flags().Changed("language") {
				lib.FlagUpdate(cmd, "language", paramsSiteUpdate.Language, mapParams)
			}
			if cmd.Flags().Changed("windows-mode-ftp") {
				mapParams["windows_mode_ftp"] = updateWindowsModeFtp
			}
			if cmd.Flags().Changed("default-time-zone") {
				lib.FlagUpdate(cmd, "default_time_zone", paramsSiteUpdate.DefaultTimeZone, mapParams)
			}
			if cmd.Flags().Changed("desktop-app") {
				mapParams["desktop_app"] = updateDesktopApp
			}
			if cmd.Flags().Changed("desktop-app-session-ip-pinning") {
				mapParams["desktop_app_session_ip_pinning"] = updateDesktopAppSessionIpPinning
			}
			if cmd.Flags().Changed("desktop-app-session-lifetime") {
				lib.FlagUpdate(cmd, "desktop_app_session_lifetime", paramsSiteUpdate.DesktopAppSessionLifetime, mapParams)
			}
			if cmd.Flags().Changed("mobile-app") {
				mapParams["mobile_app"] = updateMobileApp
			}
			if cmd.Flags().Changed("mobile-app-session-ip-pinning") {
				mapParams["mobile_app_session_ip_pinning"] = updateMobileAppSessionIpPinning
			}
			if cmd.Flags().Changed("mobile-app-session-lifetime") {
				lib.FlagUpdate(cmd, "mobile_app_session_lifetime", paramsSiteUpdate.MobileAppSessionLifetime, mapParams)
			}
			if cmd.Flags().Changed("folder-permissions-groups-only") {
				mapParams["folder_permissions_groups_only"] = updateFolderPermissionsGroupsOnly
			}
			if cmd.Flags().Changed("welcome-screen") {
				lib.FlagUpdate(cmd, "welcome_screen", paramsSiteUpdate.WelcomeScreen, mapParams)
			}
			if cmd.Flags().Changed("office-integration-available") {
				mapParams["office_integration_available"] = updateOfficeIntegrationAvailable
			}
			if cmd.Flags().Changed("office-integration-type") {
				lib.FlagUpdate(cmd, "office_integration_type", paramsSiteUpdate.OfficeIntegrationType, mapParams)
			}
			if cmd.Flags().Changed("pin-all-remote-servers-to-site-region") {
				mapParams["pin_all_remote_servers_to_site_region"] = updatePinAllRemoteServersToSiteRegion
			}
			if cmd.Flags().Changed("motd-text") {
				lib.FlagUpdate(cmd, "motd_text", paramsSiteUpdate.MotdText, mapParams)
			}
			if cmd.Flags().Changed("motd-use-for-ftp") {
				mapParams["motd_use_for_ftp"] = updateMotdUseForFtp
			}
			if cmd.Flags().Changed("motd-use-for-sftp") {
				mapParams["motd_use_for_sftp"] = updateMotdUseForSftp
			}
			if cmd.Flags().Changed("session-expiry") {
				lib.FlagUpdate(cmd, "session_expiry", paramsSiteUpdate.SessionExpiry, mapParams)
			}
			if cmd.Flags().Changed("ssl-required") {
				mapParams["ssl_required"] = updateSslRequired
			}
			if cmd.Flags().Changed("tls-disabled") {
				mapParams["tls_disabled"] = updateTlsDisabled
			}
			if cmd.Flags().Changed("sftp-insecure-ciphers") {
				mapParams["sftp_insecure_ciphers"] = updateSftpInsecureCiphers
			}
			if cmd.Flags().Changed("disable-files-certificate-generation") {
				mapParams["disable_files_certificate_generation"] = updateDisableFilesCertificateGeneration
			}
			if cmd.Flags().Changed("user-lockout") {
				mapParams["user_lockout"] = updateUserLockout
			}
			if cmd.Flags().Changed("user-lockout-tries") {
				lib.FlagUpdate(cmd, "user_lockout_tries", paramsSiteUpdate.UserLockoutTries, mapParams)
			}
			if cmd.Flags().Changed("user-lockout-within") {
				lib.FlagUpdate(cmd, "user_lockout_within", paramsSiteUpdate.UserLockoutWithin, mapParams)
			}
			if cmd.Flags().Changed("user-lockout-lock-period") {
				lib.FlagUpdate(cmd, "user_lockout_lock_period", paramsSiteUpdate.UserLockoutLockPeriod, mapParams)
			}
			if cmd.Flags().Changed("include-password-in-welcome-email") {
				mapParams["include_password_in_welcome_email"] = updateIncludePasswordInWelcomeEmail
			}
			if cmd.Flags().Changed("allowed-countries") {
				lib.FlagUpdate(cmd, "allowed_countries", paramsSiteUpdate.AllowedCountries, mapParams)
			}
			if cmd.Flags().Changed("allowed-ips") {
				lib.FlagUpdate(cmd, "allowed_ips", paramsSiteUpdate.AllowedIps, mapParams)
			}
			if cmd.Flags().Changed("disallowed-countries") {
				lib.FlagUpdate(cmd, "disallowed_countries", paramsSiteUpdate.DisallowedCountries, mapParams)
			}
			if cmd.Flags().Changed("days-to-retain-backups") {
				lib.FlagUpdate(cmd, "days_to_retain_backups", paramsSiteUpdate.DaysToRetainBackups, mapParams)
			}
			if cmd.Flags().Changed("max-prior-passwords") {
				lib.FlagUpdate(cmd, "max_prior_passwords", paramsSiteUpdate.MaxPriorPasswords, mapParams)
			}
			if cmd.Flags().Changed("password-validity-days") {
				lib.FlagUpdate(cmd, "password_validity_days", paramsSiteUpdate.PasswordValidityDays, mapParams)
			}
			if cmd.Flags().Changed("password-min-length") {
				lib.FlagUpdate(cmd, "password_min_length", paramsSiteUpdate.PasswordMinLength, mapParams)
			}
			if cmd.Flags().Changed("password-require-letter") {
				mapParams["password_require_letter"] = updatePasswordRequireLetter
			}
			if cmd.Flags().Changed("password-require-mixed") {
				mapParams["password_require_mixed"] = updatePasswordRequireMixed
			}
			if cmd.Flags().Changed("password-require-special") {
				mapParams["password_require_special"] = updatePasswordRequireSpecial
			}
			if cmd.Flags().Changed("password-require-number") {
				mapParams["password_require_number"] = updatePasswordRequireNumber
			}
			if cmd.Flags().Changed("password-require-unbreached") {
				mapParams["password_require_unbreached"] = updatePasswordRequireUnbreached
			}
			if cmd.Flags().Changed("sftp-user-root-enabled") {
				mapParams["sftp_user_root_enabled"] = updateSftpUserRootEnabled
			}
			if cmd.Flags().Changed("disable-password-reset") {
				mapParams["disable_password_reset"] = updateDisablePasswordReset
			}
			if cmd.Flags().Changed("immutable-files") {
				mapParams["immutable_files"] = updateImmutableFiles
			}
			if cmd.Flags().Changed("session-pinned-by-ip") {
				mapParams["session_pinned_by_ip"] = updateSessionPinnedByIp
			}
			if cmd.Flags().Changed("bundle-password-required") {
				mapParams["bundle_password_required"] = updateBundlePasswordRequired
			}
			if cmd.Flags().Changed("bundle-require-share-recipient") {
				mapParams["bundle_require_share_recipient"] = updateBundleRequireShareRecipient
			}
			if cmd.Flags().Changed("bundle-registration-notifications") {
				lib.FlagUpdate(cmd, "bundle_registration_notifications", paramsSiteUpdate.BundleRegistrationNotifications, mapParams)
			}
			if cmd.Flags().Changed("bundle-activity-notifications") {
				lib.FlagUpdate(cmd, "bundle_activity_notifications", paramsSiteUpdate.BundleActivityNotifications, mapParams)
			}
			if cmd.Flags().Changed("bundle-upload-receipt-notifications") {
				lib.FlagUpdate(cmd, "bundle_upload_receipt_notifications", paramsSiteUpdate.BundleUploadReceiptNotifications, mapParams)
			}
			if cmd.Flags().Changed("password-requirements-apply-to-bundles") {
				mapParams["password_requirements_apply_to_bundles"] = updatePasswordRequirementsApplyToBundles
			}
			if cmd.Flags().Changed("opt-out-global") {
				mapParams["opt_out_global"] = updateOptOutGlobal
			}
			if cmd.Flags().Changed("use-provided-modified-at") {
				mapParams["use_provided_modified_at"] = updateUseProvidedModifiedAt
			}
			if cmd.Flags().Changed("custom-namespace") {
				mapParams["custom_namespace"] = updateCustomNamespace
			}
			if cmd.Flags().Changed("disable-users-from-inactivity-period-days") {
				lib.FlagUpdate(cmd, "disable_users_from_inactivity_period_days", paramsSiteUpdate.DisableUsersFromInactivityPeriodDays, mapParams)
			}
			if cmd.Flags().Changed("non-sso-groups-allowed") {
				mapParams["non_sso_groups_allowed"] = updateNonSsoGroupsAllowed
			}
			if cmd.Flags().Changed("non-sso-users-allowed") {
				mapParams["non_sso_users_allowed"] = updateNonSsoUsersAllowed
			}
			if cmd.Flags().Changed("sharing-enabled") {
				mapParams["sharing_enabled"] = updateSharingEnabled
			}
			if cmd.Flags().Changed("user-requests-enabled") {
				mapParams["user_requests_enabled"] = updateUserRequestsEnabled
			}
			if cmd.Flags().Changed("user-requests-notify-admins") {
				mapParams["user_requests_notify_admins"] = updateUserRequestsNotifyAdmins
			}
			if cmd.Flags().Changed("ftp-enabled") {
				mapParams["ftp_enabled"] = updateFtpEnabled
			}
			if cmd.Flags().Changed("sftp-enabled") {
				mapParams["sftp_enabled"] = updateSftpEnabled
			}
			if cmd.Flags().Changed("sftp-host-key-type") {
				lib.FlagUpdate(cmd, "sftp_host_key_type", paramsSiteUpdate.SftpHostKeyType, mapParams)
			}
			if cmd.Flags().Changed("active-sftp-host-key-id") {
				lib.FlagUpdate(cmd, "active_sftp_host_key_id", paramsSiteUpdate.ActiveSftpHostKeyId, mapParams)
			}
			if cmd.Flags().Changed("bundle-watermark-value") {
				lib.FlagUpdateLen(cmd, "bundle_watermark_value", paramsSiteUpdate.BundleWatermarkValue, mapParams)
			}
			if cmd.Flags().Changed("group-admins-can-set-user-password") {
				mapParams["group_admins_can_set_user_password"] = updateGroupAdminsCanSetUserPassword
			}
			if cmd.Flags().Changed("allowed-2fa-method-sms") {
				mapParams["allowed_2fa_method_sms"] = updateAllowed2faMethodSms
			}
			if cmd.Flags().Changed("allowed-2fa-method-u2f") {
				mapParams["allowed_2fa_method_u2f"] = updateAllowed2faMethodU2f
			}
			if cmd.Flags().Changed("allowed-2fa-method-totp") {
				mapParams["allowed_2fa_method_totp"] = updateAllowed2faMethodTotp
			}
			if cmd.Flags().Changed("allowed-2fa-method-webauthn") {
				mapParams["allowed_2fa_method_webauthn"] = updateAllowed2faMethodWebauthn
			}
			if cmd.Flags().Changed("allowed-2fa-method-yubi") {
				mapParams["allowed_2fa_method_yubi"] = updateAllowed2faMethodYubi
			}
			if cmd.Flags().Changed("allowed-2fa-method-bypass-for-ftp-sftp-dav") {
				mapParams["allowed_2fa_method_bypass_for_ftp_sftp_dav"] = updateAllowed2faMethodBypassForFtpSftpDav
			}
			if cmd.Flags().Changed("require-2fa") {
				mapParams["require_2fa"] = updateRequire2fa
			}
			if cmd.Flags().Changed("require-2fa-user-type") {
				lib.FlagUpdate(cmd, "require_2fa_user_type", paramsSiteUpdate.Require2faUserType, mapParams)
			}
			if cmd.Flags().Changed("color2-top") {
				lib.FlagUpdate(cmd, "color2_top", paramsSiteUpdate.Color2Top, mapParams)
			}
			if cmd.Flags().Changed("color2-left") {
				lib.FlagUpdate(cmd, "color2_left", paramsSiteUpdate.Color2Left, mapParams)
			}
			if cmd.Flags().Changed("color2-link") {
				lib.FlagUpdate(cmd, "color2_link", paramsSiteUpdate.Color2Link, mapParams)
			}
			if cmd.Flags().Changed("color2-text") {
				lib.FlagUpdate(cmd, "color2_text", paramsSiteUpdate.Color2Text, mapParams)
			}
			if cmd.Flags().Changed("color2-top-text") {
				lib.FlagUpdate(cmd, "color2_top_text", paramsSiteUpdate.Color2TopText, mapParams)
			}
			if cmd.Flags().Changed("site-header") {
				lib.FlagUpdate(cmd, "site_header", paramsSiteUpdate.SiteHeader, mapParams)
			}
			if cmd.Flags().Changed("site-footer") {
				lib.FlagUpdate(cmd, "site_footer", paramsSiteUpdate.SiteFooter, mapParams)
			}
			if cmd.Flags().Changed("login-help-text") {
				lib.FlagUpdate(cmd, "login_help_text", paramsSiteUpdate.LoginHelpText, mapParams)
			}
			if cmd.Flags().Changed("smtp-address") {
				lib.FlagUpdate(cmd, "smtp_address", paramsSiteUpdate.SmtpAddress, mapParams)
			}
			if cmd.Flags().Changed("smtp-authentication") {
				lib.FlagUpdate(cmd, "smtp_authentication", paramsSiteUpdate.SmtpAuthentication, mapParams)
			}
			if cmd.Flags().Changed("smtp-from") {
				lib.FlagUpdate(cmd, "smtp_from", paramsSiteUpdate.SmtpFrom, mapParams)
			}
			if cmd.Flags().Changed("smtp-username") {
				lib.FlagUpdate(cmd, "smtp_username", paramsSiteUpdate.SmtpUsername, mapParams)
			}
			if cmd.Flags().Changed("smtp-port") {
				lib.FlagUpdate(cmd, "smtp_port", paramsSiteUpdate.SmtpPort, mapParams)
			}
			if cmd.Flags().Changed("ldap-enabled") {
				mapParams["ldap_enabled"] = updateLdapEnabled
			}
			if cmd.Flags().Changed("ldap-type") {
				lib.FlagUpdate(cmd, "ldap_type", paramsSiteUpdate.LdapType, mapParams)
			}
			if cmd.Flags().Changed("ldap-host") {
				lib.FlagUpdate(cmd, "ldap_host", paramsSiteUpdate.LdapHost, mapParams)
			}
			if cmd.Flags().Changed("ldap-host-2") {
				lib.FlagUpdate(cmd, "ldap_host_2", paramsSiteUpdate.LdapHost2, mapParams)
			}
			if cmd.Flags().Changed("ldap-host-3") {
				lib.FlagUpdate(cmd, "ldap_host_3", paramsSiteUpdate.LdapHost3, mapParams)
			}
			if cmd.Flags().Changed("ldap-port") {
				lib.FlagUpdate(cmd, "ldap_port", paramsSiteUpdate.LdapPort, mapParams)
			}
			if cmd.Flags().Changed("ldap-secure") {
				mapParams["ldap_secure"] = updateLdapSecure
			}
			if cmd.Flags().Changed("ldap-username") {
				lib.FlagUpdate(cmd, "ldap_username", paramsSiteUpdate.LdapUsername, mapParams)
			}
			if cmd.Flags().Changed("ldap-username-field") {
				lib.FlagUpdate(cmd, "ldap_username_field", paramsSiteUpdate.LdapUsernameField, mapParams)
			}
			if cmd.Flags().Changed("ldap-domain") {
				lib.FlagUpdate(cmd, "ldap_domain", paramsSiteUpdate.LdapDomain, mapParams)
			}
			if cmd.Flags().Changed("ldap-user-action") {
				lib.FlagUpdate(cmd, "ldap_user_action", paramsSiteUpdate.LdapUserAction, mapParams)
			}
			if cmd.Flags().Changed("ldap-group-action") {
				lib.FlagUpdate(cmd, "ldap_group_action", paramsSiteUpdate.LdapGroupAction, mapParams)
			}
			if cmd.Flags().Changed("ldap-user-include-groups") {
				lib.FlagUpdate(cmd, "ldap_user_include_groups", paramsSiteUpdate.LdapUserIncludeGroups, mapParams)
			}
			if cmd.Flags().Changed("ldap-group-exclusion") {
				lib.FlagUpdate(cmd, "ldap_group_exclusion", paramsSiteUpdate.LdapGroupExclusion, mapParams)
			}
			if cmd.Flags().Changed("ldap-group-inclusion") {
				lib.FlagUpdate(cmd, "ldap_group_inclusion", paramsSiteUpdate.LdapGroupInclusion, mapParams)
			}
			if cmd.Flags().Changed("ldap-base-dn") {
				lib.FlagUpdate(cmd, "ldap_base_dn", paramsSiteUpdate.LdapBaseDn, mapParams)
			}
			if cmd.Flags().Changed("uploads-via-email-authentication") {
				mapParams["uploads_via_email_authentication"] = updateUploadsViaEmailAuthentication
			}
			if cmd.Flags().Changed("icon16-file") {
			}
			if cmd.Flags().Changed("icon16-delete") {
				mapParams["icon16_delete"] = updateIcon16Delete
			}
			if cmd.Flags().Changed("icon32-file") {
			}
			if cmd.Flags().Changed("icon32-delete") {
				mapParams["icon32_delete"] = updateIcon32Delete
			}
			if cmd.Flags().Changed("icon48-file") {
			}
			if cmd.Flags().Changed("icon48-delete") {
				mapParams["icon48_delete"] = updateIcon48Delete
			}
			if cmd.Flags().Changed("icon128-file") {
			}
			if cmd.Flags().Changed("icon128-delete") {
				mapParams["icon128_delete"] = updateIcon128Delete
			}
			if cmd.Flags().Changed("logo-file") {
			}
			if cmd.Flags().Changed("logo-delete") {
				mapParams["logo_delete"] = updateLogoDelete
			}
			if cmd.Flags().Changed("bundle-watermark-attachment-file") {
			}
			if cmd.Flags().Changed("bundle-watermark-attachment-delete") {
				mapParams["bundle_watermark_attachment_delete"] = updateBundleWatermarkAttachmentDelete
			}
			if cmd.Flags().Changed("disable-2fa-with-delay") {
				mapParams["disable_2fa_with_delay"] = updateDisable2faWithDelay
			}
			if cmd.Flags().Changed("ldap-password-change") {
				lib.FlagUpdate(cmd, "ldap_password_change", paramsSiteUpdate.LdapPasswordChange, mapParams)
			}
			if cmd.Flags().Changed("ldap-password-change-confirmation") {
				lib.FlagUpdate(cmd, "ldap_password_change_confirmation", paramsSiteUpdate.LdapPasswordChangeConfirmation, mapParams)
			}
			if cmd.Flags().Changed("smtp-password") {
				lib.FlagUpdate(cmd, "smtp_password", paramsSiteUpdate.SmtpPassword, mapParams)
			}
			if cmd.Flags().Changed("session-expiry-minutes") {
				lib.FlagUpdate(cmd, "session_expiry_minutes", paramsSiteUpdate.SessionExpiryMinutes, mapParams)
			}

			var site interface{}
			var err error
			site, err = client.UpdateWithMap(ctx, mapParams)
			return lib.HandleResponse(ctx, Profile(cmd), site, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
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
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.OfficeIntegrationType, "office-integration-type", "", "Office integration application used to edit and view the MS Office documents")
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
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.BundleRegistrationNotifications, "bundle-registration-notifications", "", "Do Bundle owners receive registration notification?")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.BundleActivityNotifications, "bundle-activity-notifications", "", "Do Bundle owners receive activity notifications?")
	cmdUpdate.Flags().StringVar(&paramsSiteUpdate.BundleUploadReceiptNotifications, "bundle-upload-receipt-notifications", "", "Do Bundle uploaders receive upload confirmation notifications?")
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
	cmdUpdate.Flags().BoolVar(&updateGroupAdminsCanSetUserPassword, "group-admins-can-set-user-password", updateGroupAdminsCanSetUserPassword, "Allow group admins set password authentication method")
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
	cmdUpdate.Flags().BoolVar(&updateUploadsViaEmailAuthentication, "uploads-via-email-authentication", updateUploadsViaEmailAuthentication, "Do incoming emails in the Inboxes require checking for SPF/DKIM/DMARC?")
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
	cmdUpdate.Flags().Int64Var(&paramsSiteUpdate.SessionExpiryMinutes, "session-expiry-minutes", 0, "Session expiry in minutes")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
json-styles: {raw, pretty}`)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	Sites.AddCommand(cmdUpdate)
	return Sites
}
