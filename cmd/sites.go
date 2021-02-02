package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/site"
)

var (
	Sites = &cobra.Command{
		Use:  "sites [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func SitesInit() {
	var fieldsGet string
	paramsSiteGet := files_sdk.SiteGetParams{}
	cmdGet := &cobra.Command{
		Use: "get",
		Run: func(cmd *cobra.Command, args []string) {
			client := site.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Get(paramsSiteGet)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = lib.JsonMarshal(result, fieldsGet)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdGet.Flags().StringVarP(&paramsSiteGet.Format, "format", "f", "", "")

	cmdGet.Flags().StringVarP(&fieldsGet, "fields", "", "", "comma separated list of field names")
	Sites.AddCommand(cmdGet)
	var fieldsGetUsage string
	paramsSiteGetUsage := files_sdk.SiteGetUsageParams{}
	cmdGetUsage := &cobra.Command{
		Use: "get-usage",
		Run: func(cmd *cobra.Command, args []string) {
			client := site.Client{Config: files_sdk.GlobalConfig}
			result, err := client.GetUsage(paramsSiteGetUsage)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = lib.JsonMarshal(result, fieldsGetUsage)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdGetUsage.Flags().StringVarP(&paramsSiteGetUsage.Format, "format", "f", "", "")

	cmdGetUsage.Flags().StringVarP(&fieldsGetUsage, "fields", "", "", "comma separated list of field names")
	Sites.AddCommand(cmdGetUsage)
	var fieldsUpdate string
	paramsSiteUpdate := files_sdk.SiteUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			client := site.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Update(paramsSiteUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = lib.JsonMarshal(result, fieldsUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.Name, "name", "", "", "Site name")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.Subdomain, "subdomain", "", "", "Site subdomain")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.Domain, "domain", "", "", "Custom domain")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.Email, "email", "", "", "Main email for this site")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.ReplyToEmail, "reply-to-email", "", "", "Reply-to email for this site")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.BundleExpiration, "bundle-expiration", "e", 0, "Site-wide Bundle expiration in days")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.WelcomeEmailCc, "welcome-email-cc", "", "", "Include this email in welcome emails if enabled")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.WelcomeCustomText, "welcome-custom-text", "", "", "Custom text send in user welcome email")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.Language, "language", "", "", "Site default language")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.DefaultTimeZone, "default-time-zone", "f", "", "Site default time zone")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.DesktopAppSessionLifetime, "desktop-app-session-lifetime", "", 0, "Desktop app session lifetime (in hours)")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.WelcomeScreen, "welcome-screen", "", "", "Does the welcome screen appear?")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.UserLockoutTries, "user-lockout-tries", "", 0, "Number of login tries within `user_lockout_within` hours before users are locked out")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.UserLockoutWithin, "user-lockout-within", "", 0, "Number of hours for user lockout window")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.UserLockoutLockPeriod, "user-lockout-lock-period", "", 0, "How many hours to lock user out for failed password?")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.AllowedCountries, "allowed-countries", "c", "", "Comma seperated list of allowed Country codes")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.AllowedIps, "allowed-ips", "i", "", "List of allowed IP addresses")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.DisallowedCountries, "disallowed-countries", "w", "", "Comma seperated list of disallowed Country codes")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.DaysToRetainBackups, "days-to-retain-backups", "d", 0, "Number of days to keep deleted files")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.MaxPriorPasswords, "max-prior-passwords", "", 0, "Number of prior passwords to disallow")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.PasswordValidityDays, "password-validity-days", "", 0, "Number of days password is valid")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.PasswordMinLength, "password-min-length", "", 0, "Shortest password length for users")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.DisableUsersFromInactivityPeriodDays, "disable-users-from-inactivity-period-days", "", 0, "If greater than zero, users will unable to login if they do not show activity within this number of days.")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.Require2faUserType, "require-2fa-user-type", "", "", "What type of user is required to use two-factor authentication (when require_2fa is set to `true` for this site)?")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.Color2Top, "color2-top", "o", "", "Top bar background color")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.Color2Left, "color2-left", "l", "", "Page link and button color")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.Color2Link, "color2-link", "n", "", "Top bar link color")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.Color2Text, "color2-text", "x", "", "Page link and button color")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.Color2TopText, "color2-top-text", "", "", "Top bar text color")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.SiteHeader, "site-header", "", "", "Custom site header text")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.SiteFooter, "site-footer", "", "", "Custom site footer text")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LoginHelpText, "login-help-text", "", "", "Login help text")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.SmtpAddress, "smtp-address", "", "", "SMTP server hostname or IP")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.SmtpAuthentication, "smtp-authentication", "", "", "SMTP server authentication type")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.SmtpFrom, "smtp-from", "", "", "From address to use when mailing through custom SMTP")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.SmtpUsername, "smtp-username", "", "", "SMTP server username")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.SmtpPort, "smtp-port", "", 0, "SMTP server port")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapType, "ldap-type", "", "", "LDAP type")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapHost, "ldap-host", "", "", "LDAP host")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapHost2, "ldap-host-2", "2", "", "LDAP backup host")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapHost3, "ldap-host-3", "3", "", "LDAP backup host")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.LdapPort, "ldap-port", "", 0, "LDAP port")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapUsername, "ldap-username", "", "", "Username for signing in to LDAP server.")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapUsernameField, "ldap-username-field", "", "", "LDAP username field")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapDomain, "ldap-domain", "", "", "Domain name that will be appended to usernames")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapUserAction, "ldap-user-action", "", "", "Should we sync users from LDAP server?")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapGroupAction, "ldap-group-action", "", "", "Should we sync groups from LDAP server?")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapUserIncludeGroups, "ldap-user-include-groups", "", "", "Comma or newline separated list of group names (with optional wildcards) - if provided, only users in these groups will be added or synced.")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapGroupExclusion, "ldap-group-exclusion", "", "", "Comma or newline separated list of group names (with optional wildcards) to exclude when syncing.")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapGroupInclusion, "ldap-group-inclusion", "", "", "Comma or newline separated list of group names (with optional wildcards) to include when syncing.")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapBaseDn, "ldap-base-dn", "", "", "Base DN for looking up users in LDAP server")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapPasswordChange, "ldap-password-change", "", "", "New LDAP password.")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapPasswordChangeConfirmation, "ldap-password-change-confirmation", "", "", "Confirm new LDAP password.")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.SmtpPassword, "smtp-password", "", "", "Password for SMTP server.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	Sites.AddCommand(cmdUpdate)
}
