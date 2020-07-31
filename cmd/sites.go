package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/site"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = site.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
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
	cmdGet := &cobra.Command{
		Use: "get",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := site.Get()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsGet)
		},
	}
	cmdGet.Flags().StringVarP(&fieldsGet, "fields", "f", "", "comma separated list of field names")
	Sites.AddCommand(cmdGet)
	var fieldsGetUsage string
	cmdGetUsage := &cobra.Command{
		Use: "get-usage",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := site.GetUsage()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsGetUsage)
		},
	}
	cmdGetUsage.Flags().StringVarP(&fieldsGetUsage, "fields", "f", "", "comma separated list of field names")
	Sites.AddCommand(cmdGetUsage)
	var fieldsUpdate string
	paramsSiteUpdate := files_sdk.SiteUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := site.Update(paramsSiteUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsUpdate)
		},
	}
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.Name, "name", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.Subdomain, "subdomain", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.Domain, "domain", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.Email, "email", "", "", "Update site settings")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.BundleExpiration, "bundle-expiration", "e", 0, "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.WelcomeEmailCc, "welcome-email-cc", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.WelcomeCustomText, "welcome-custom-text", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.Language, "language", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.DefaultTimeZone, "default-time-zone", "z", "", "Update site settings")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.DesktopAppSessionLifetime, "desktop-app-session-lifetime", "", 0, "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.WelcomeScreen, "welcome-screen", "", "", "Update site settings")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.UserLockoutTries, "user-lockout-tries", "", 0, "Update site settings")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.UserLockoutWithin, "user-lockout-within", "", 0, "Update site settings")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.UserLockoutLockPeriod, "user-lockout-lock-period", "", 0, "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.AllowedIps, "allowed-ips", "i", "", "Update site settings")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.DaysToRetainBackups, "days-to-retain-backups", "d", 0, "Update site settings")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.MaxPriorPasswords, "max-prior-passwords", "", 0, "Update site settings")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.PasswordValidityDays, "password-validity-days", "", 0, "Update site settings")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.PasswordMinLength, "password-min-length", "", 0, "Update site settings")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.DisableUsersFromInactivityPeriodDays, "disable-users-from-inactivity-period-days", "", 0, "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.Require2faUserType, "require-2fa-user-type", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.Color2Top, "color2-top", "o", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.Color2Left, "color2-left", "l", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.Color2Link, "color2-link", "n", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.Color2Text, "color2-text", "x", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.Color2TopText, "color2-top-text", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.SiteHeader, "site-header", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.SiteFooter, "site-footer", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LoginHelpText, "login-help-text", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.SmtpAddress, "smtp-address", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.SmtpAuthentication, "smtp-authentication", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.SmtpFrom, "smtp-from", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.SmtpUsername, "smtp-username", "", "", "Update site settings")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.SmtpPort, "smtp-port", "", 0, "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapType, "ldap-type", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapHost, "ldap-host", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapHost2, "ldap-host-2", "2", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapHost3, "ldap-host-3", "3", "", "Update site settings")
	cmdUpdate.Flags().IntVarP(&paramsSiteUpdate.LdapPort, "ldap-port", "", 0, "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapUsername, "ldap-username", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapUsernameField, "ldap-username-field", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapDomain, "ldap-domain", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapUserAction, "ldap-user-action", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapGroupAction, "ldap-group-action", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapUserIncludeGroups, "ldap-user-include-groups", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapGroupExclusion, "ldap-group-exclusion", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapGroupInclusion, "ldap-group-inclusion", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapBaseDn, "ldap-base-dn", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapPasswordChange, "ldap-password-change", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.LdapPasswordChangeConfirmation, "ldap-password-change-confirmation", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&paramsSiteUpdate.SmtpPassword, "smtp-password", "", "", "Update site settings")
	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "f", "", "comma separated list of field names")
	Sites.AddCommand(cmdUpdate)
}
