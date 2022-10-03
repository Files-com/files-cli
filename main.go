package main

import (
	"context"
	"strings"

	"github.com/Files-com/files-cli/cmd"
	"github.com/Files-com/files-cli/lib"
	files "github.com/Files-com/files-sdk-go/v2"
	"github.com/spf13/cobra"

	_ "embed"
	"fmt"
	"os"
)

//go:embed _VERSION
var VERSION string
var OutputPath string

func main() {
	var ignoreVersionCheck bool
	files.GlobalConfig.UserAgent = "Files.com CLI" + " " + strings.TrimSpace(VERSION)
	var rootCmd = &cobra.Command{
		Use:     "files-cli [resource]",
		Version: strings.TrimSuffix(VERSION, "\n"),
		PersistentPreRun: func(x *cobra.Command, args []string) {
			if lib.Includes(x.Use, []string{"login", "logout"}) {
				return
			}

			if len(x.Aliases) != 0 && lib.Includes(x.Aliases[0], []string{"config-set", "config-reset", "config-show"}) {
				return
			}
			config := &lib.Config{}
			err := config.Load()

			if OutputPath != "" {
				output, err := os.Create(OutputPath)
				if err != nil {
					lib.ClientError(x.Context(), err, x.ErrOrStderr())
				}
				x.SetOut(output)
			}

			if !ignoreVersionCheck {
				config.CheckVersion(VERSION, lib.FetchLatestVersionNumber(x.Context()), lib.InstalledViaBrew(), x.ErrOrStderr())
			}

			if err != nil {
				fmt.Fprintf(x.ErrOrStderr(), "%v\n", err)
				os.Exit(1)
			}

			if files.GlobalConfig.GetAPIKey() != "" {
				return
			}

			if config.ValidSession() {
				return
			}
			config.Overrides = lib.Overrides{In: x.InOrStdin(), Out: x.OutOrStdout()}
			if config.SessionExpired() {
				fmt.Fprintf(x.ErrOrStderr(), "The session has expired, you must log in again.\n")
				err = lib.CreateSession(files.SessionCreateParams{}, config)
				if err != nil {
					fmt.Fprintf(x.ErrOrStderr(), "%v\n", err)
					os.Exit(1)
				}
				return
			}

			if files.GlobalConfig.GetAPIKey() == "" {
				fmt.Fprintf(x.ErrOrStderr(), "No API Key found. Using session login.\n")
				err = lib.CreateSession(files.SessionCreateParams{}, config)
				if err != nil {
					fmt.Fprintf(x.ErrOrStderr(), "%v\n", err)
					os.Exit(1)
				}
			}
		},
	}
	rootCmd.PersistentFlags().StringVar(&files.APIKey, "api-key", "", "API Key")
	rootCmd.PersistentFlags().BoolVar(&ignoreVersionCheck, "ignore-version-check", false, "API Key")
	rootCmd.PersistentFlags().StringVarP(&OutputPath, "output", "o", "", "file path to save output")
	rootCmd.SuggestionsMinimumDistance = 1
	cmd.ConfigInit()
	rootCmd.AddCommand(cmd.Config)
	rootCmd.AddCommand(cmd.SyncInit())
	rootCmd.AddCommand(cmd.UploadCmd())
	rootCmd.AddCommand(cmd.DownloadCmd())
	rootCmd.AddCommand(cmd.VersionCmd(VERSION))
	rootCmd.AddCommand(cmd.Login)
	rootCmd.AddCommand(cmd.LogOut)
	cmd.AccountLineItemsInit()
	rootCmd.AddCommand(cmd.AccountLineItems)
	cmd.ActionsInit()
	rootCmd.AddCommand(cmd.Actions)
	cmd.ActionNotificationExportsInit()
	rootCmd.AddCommand(cmd.ActionNotificationExports)
	cmd.ActionNotificationExportResultsInit()
	rootCmd.AddCommand(cmd.ActionNotificationExportResults)
	cmd.ActionWebhookFailuresInit()
	rootCmd.AddCommand(cmd.ActionWebhookFailures)
	cmd.ApiKeysInit()
	rootCmd.AddCommand(cmd.ApiKeys)
	cmd.AppsInit()
	rootCmd.AddCommand(cmd.Apps)
	cmd.As2IncomingMessagesInit()
	rootCmd.AddCommand(cmd.As2IncomingMessages)
	cmd.As2OutgoingMessagesInit()
	rootCmd.AddCommand(cmd.As2OutgoingMessages)
	cmd.As2PartnersInit()
	rootCmd.AddCommand(cmd.As2Partners)
	cmd.As2StationsInit()
	rootCmd.AddCommand(cmd.As2Stations)
	cmd.AutosInit()
	rootCmd.AddCommand(cmd.Autos)
	cmd.AutomationsInit()
	rootCmd.AddCommand(cmd.Automations)
	cmd.AutomationRunsInit()
	rootCmd.AddCommand(cmd.AutomationRuns)
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
	cmd.FileMigrationsInit()
	rootCmd.AddCommand(cmd.FileMigrations)
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
	cmd.PrioritiesInit()
	rootCmd.AddCommand(cmd.Priorities)
	cmd.ProjectsInit()
	rootCmd.AddCommand(cmd.Projects)
	cmd.PublicIpAddressesInit()
	rootCmd.AddCommand(cmd.PublicIpAddresses)
	cmd.PublicKeysInit()
	rootCmd.AddCommand(cmd.PublicKeys)
	cmd.RemoteBandwidthSnapshotsInit()
	rootCmd.AddCommand(cmd.RemoteBandwidthSnapshots)
	cmd.RemoteServersInit()
	rootCmd.AddCommand(cmd.RemoteServers)
	cmd.RequestsInit()
	rootCmd.AddCommand(cmd.Requests)
	cmd.SessionsInit()
	rootCmd.AddCommand(cmd.Sessions)
	cmd.SettingsChangesInit()
	rootCmd.AddCommand(cmd.SettingsChanges)
	cmd.SftpHostKeysInit()
	rootCmd.AddCommand(cmd.SftpHostKeys)
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
	cmd.WebhookTestsInit()
	rootCmd.AddCommand(cmd.WebhookTests)

	rootCmd.ExecuteContext(context.WithValue(context.Background(), "config", &files.GlobalConfig))
}
