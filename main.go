package main

import (
	"github.com/Files-com/files-cli/cmd"
	"github.com/Files-com/files-cli/lib"
	files "github.com/Files-com/files-sdk-go"
	"github.com/spf13/cobra"

	"fmt"
	"os"
)

var VERSION = "1.0.267"

func main() {
	var rootCmd = &cobra.Command{
		Use:     "files-cli [resource]",
		Version: VERSION,
		PersistentPreRun: func(x *cobra.Command, args []string) {
			if len(x.Aliases) != 0 && (x.Aliases[0] == "config-set" || x.Aliases[0] == "config-reset") {
				return
			}
			config := &lib.Config{}
			err := config.Load()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if files.GlobalConfig.GetAPIKey() != "" {
				return
			}

			if config.ValidSession() {
				return
			}

			if config.SessionExpired() {
				fmt.Println("The session has expired, you must log in again.")
				err = lib.CreateSession(files.SessionCreateParams{}, *config)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				return
			}

			if files.GlobalConfig.GetAPIKey() == "" {
				fmt.Println("No API Key found. Using session login.")
				err = lib.CreateSession(files.SessionCreateParams{}, *config)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		},
	}
	rootCmd.PersistentFlags().StringVar(&files.APIKey, "api-key", "", "API Key")
	rootCmd.SuggestionsMinimumDistance = 1
	cmd.ConfigInit()
	rootCmd.AddCommand(cmd.Config)
	rootCmd.AddCommand(cmd.UploadCmd())
	rootCmd.AddCommand(cmd.DownloadCmd())
	rootCmd.AddCommand(cmd.VersionCmd(VERSION))
	rootCmd.AddCommand(cmd.Login)
	rootCmd.AddCommand(cmd.LogOut)
	cmd.AccountLineItemsInit()
	rootCmd.AddCommand(cmd.AccountLineItems)
	cmd.ActionsInit()
	rootCmd.AddCommand(cmd.Actions)
	cmd.ApiKeysInit()
	rootCmd.AddCommand(cmd.ApiKeys)
	cmd.AppsInit()
	rootCmd.AddCommand(cmd.Apps)
	cmd.As2KeysInit()
	rootCmd.AddCommand(cmd.As2Keys)
	cmd.AutosInit()
	rootCmd.AddCommand(cmd.Autos)
	cmd.AutomationsInit()
	rootCmd.AddCommand(cmd.Automations)
	cmd.BandwidthSnapshotsInit()
	rootCmd.AddCommand(cmd.BandwidthSnapshots)
	cmd.BehaviorsInit()
	rootCmd.AddCommand(cmd.Behaviors)
	cmd.BundlesInit()
	rootCmd.AddCommand(cmd.Bundles)
	cmd.BundleDownloadsInit()
	rootCmd.AddCommand(cmd.BundleDownloads)
	cmd.BundleRecipientsInit()
	rootCmd.AddCommand(cmd.BundleRecipients)
	cmd.BundleRegistrationsInit()
	rootCmd.AddCommand(cmd.BundleRegistrations)
	cmd.ClickwrapsInit()
	rootCmd.AddCommand(cmd.Clickwraps)
	cmd.DnsRecordsInit()
	rootCmd.AddCommand(cmd.DnsRecords)
	cmd.ErrorsInit()
	rootCmd.AddCommand(cmd.Errors)
	cmd.ExternalEventsInit()
	rootCmd.AddCommand(cmd.ExternalEvents)
	cmd.FilesInit()
	rootCmd.AddCommand(cmd.Files)
	cmd.FileActionsInit()
	rootCmd.AddCommand(cmd.FileActions)
	cmd.FileCommentsInit()
	rootCmd.AddCommand(cmd.FileComments)
	cmd.FileCommentReactionsInit()
	rootCmd.AddCommand(cmd.FileCommentReactions)
	cmd.FileUploadPartsInit()
	rootCmd.AddCommand(cmd.FileUploadParts)
	cmd.FoldersInit()
	rootCmd.AddCommand(cmd.Folders)
	cmd.FormFieldsInit()
	rootCmd.AddCommand(cmd.FormFields)
	cmd.FormFieldSetsInit()
	rootCmd.AddCommand(cmd.FormFieldSets)
	cmd.GroupsInit()
	rootCmd.AddCommand(cmd.Groups)
	cmd.GroupUsersInit()
	rootCmd.AddCommand(cmd.GroupUsers)
	cmd.HistoriesInit()
	rootCmd.AddCommand(cmd.Histories)
	cmd.HistoryExportsInit()
	rootCmd.AddCommand(cmd.HistoryExports)
	cmd.HistoryExportResultsInit()
	rootCmd.AddCommand(cmd.HistoryExportResults)
	cmd.ImagesInit()
	rootCmd.AddCommand(cmd.Images)
	cmd.InboxRecipientsInit()
	rootCmd.AddCommand(cmd.InboxRecipients)
	cmd.InboxRegistrationsInit()
	rootCmd.AddCommand(cmd.InboxRegistrations)
	cmd.InboxUploadsInit()
	rootCmd.AddCommand(cmd.InboxUploads)
	cmd.InvoicesInit()
	rootCmd.AddCommand(cmd.Invoices)
	cmd.InvoiceLineItemsInit()
	rootCmd.AddCommand(cmd.InvoiceLineItems)
	cmd.IpAddressesInit()
	rootCmd.AddCommand(cmd.IpAddresses)
	cmd.LocksInit()
	rootCmd.AddCommand(cmd.Locks)
	cmd.MessagesInit()
	rootCmd.AddCommand(cmd.Messages)
	cmd.MessageCommentsInit()
	rootCmd.AddCommand(cmd.MessageComments)
	cmd.MessageCommentReactionsInit()
	rootCmd.AddCommand(cmd.MessageCommentReactions)
	cmd.MessageReactionsInit()
	rootCmd.AddCommand(cmd.MessageReactions)
	cmd.NotificationsInit()
	rootCmd.AddCommand(cmd.Notifications)
	cmd.PaymentsInit()
	rootCmd.AddCommand(cmd.Payments)
	cmd.PaymentLineItemsInit()
	rootCmd.AddCommand(cmd.PaymentLineItems)
	cmd.PermissionsInit()
	rootCmd.AddCommand(cmd.Permissions)
	cmd.PreviewsInit()
	rootCmd.AddCommand(cmd.Previews)
	cmd.ProjectsInit()
	rootCmd.AddCommand(cmd.Projects)
	cmd.PublicIpAddressesInit()
	rootCmd.AddCommand(cmd.PublicIpAddresses)
	cmd.PublicKeysInit()
	rootCmd.AddCommand(cmd.PublicKeys)
	cmd.RemoteServersInit()
	rootCmd.AddCommand(cmd.RemoteServers)
	cmd.RequestsInit()
	rootCmd.AddCommand(cmd.Requests)
	cmd.SessionsInit()
	rootCmd.AddCommand(cmd.Sessions)
	cmd.SettingsChangesInit()
	rootCmd.AddCommand(cmd.SettingsChanges)
	cmd.SitesInit()
	rootCmd.AddCommand(cmd.Sites)
	cmd.SsoStrategiesInit()
	rootCmd.AddCommand(cmd.SsoStrategies)
	cmd.StatusesInit()
	rootCmd.AddCommand(cmd.Statuses)
	cmd.StylesInit()
	rootCmd.AddCommand(cmd.Styles)
	cmd.UsageDailySnapshotsInit()
	rootCmd.AddCommand(cmd.UsageDailySnapshots)
	cmd.UsageSnapshotsInit()
	rootCmd.AddCommand(cmd.UsageSnapshots)
	cmd.UsersInit()
	rootCmd.AddCommand(cmd.Users)
	cmd.UserCipherUsesInit()
	rootCmd.AddCommand(cmd.UserCipherUses)
	cmd.UserRequestsInit()
	rootCmd.AddCommand(cmd.UserRequests)

	rootCmd.Execute()
}
