package cmd

import (
	"fmt"
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	siem_http_destination "github.com/Files-com/files-sdk-go/v3/siemhttpdestination"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(SiemHttpDestinations())
}

func SiemHttpDestinations() *cobra.Command {
	SiemHttpDestinations := &cobra.Command{
		Use:  "siem-http-destinations [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command siem-http-destinations\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsSiemHttpDestinationList := files_sdk.SiemHttpDestinationListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List SIEM HTTP Destinations",
		Long:    `List SIEM HTTP Destinations`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsSiemHttpDestinationList
			params.MaxPages = MaxPagesList

			client := siem_http_destination.Client{Config: config}
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

	cmdList.Flags().StringVar(&paramsSiemHttpDestinationList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsSiemHttpDestinationList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	SiemHttpDestinations.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsSiemHttpDestinationFind := files_sdk.SiemHttpDestinationFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show SIEM HTTP Destination`,
		Long:  `Show SIEM HTTP Destination`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := siem_http_destination.Client{Config: config}

			var siemHttpDestination interface{}
			var err error
			siemHttpDestination, err = client.Find(paramsSiemHttpDestinationFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), siemHttpDestination, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsSiemHttpDestinationFind.Id, "id", 0, "Siem Http Destination ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	SiemHttpDestinations.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createSendingActive := true
	createSftpActionSendEnabled := true
	createFtpActionSendEnabled := true
	createWebDavActionSendEnabled := true
	createSyncSendEnabled := true
	createOutboundConnectionSendEnabled := true
	createAutomationSendEnabled := true
	createApiRequestSendEnabled := true
	createPublicHostingRequestSendEnabled := true
	createEmailSendEnabled := true
	createExavaultApiRequestSendEnabled := true
	createSettingsChangeSendEnabled := true
	paramsSiemHttpDestinationCreate := files_sdk.SiemHttpDestinationCreateParams{}
	SiemHttpDestinationCreateGenericPayloadType := ""
	SiemHttpDestinationCreateDestinationType := ""

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create SIEM HTTP Destination`,
		Long:  `Create SIEM HTTP Destination`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := siem_http_destination.Client{Config: config}

			var SiemHttpDestinationCreateGenericPayloadTypeErr error
			paramsSiemHttpDestinationCreate.GenericPayloadType, SiemHttpDestinationCreateGenericPayloadTypeErr = lib.FetchKey("generic-payload-type", paramsSiemHttpDestinationCreate.GenericPayloadType.Enum(), SiemHttpDestinationCreateGenericPayloadType)
			if SiemHttpDestinationCreateGenericPayloadType != "" && SiemHttpDestinationCreateGenericPayloadTypeErr != nil {
				return SiemHttpDestinationCreateGenericPayloadTypeErr
			}
			var SiemHttpDestinationCreateDestinationTypeErr error
			paramsSiemHttpDestinationCreate.DestinationType, SiemHttpDestinationCreateDestinationTypeErr = lib.FetchKey("destination-type", paramsSiemHttpDestinationCreate.DestinationType.Enum(), SiemHttpDestinationCreateDestinationType)
			if SiemHttpDestinationCreateDestinationType != "" && SiemHttpDestinationCreateDestinationTypeErr != nil {
				return SiemHttpDestinationCreateDestinationTypeErr
			}

			if cmd.Flags().Changed("sending-active") {
				paramsSiemHttpDestinationCreate.SendingActive = flib.Bool(createSendingActive)
			}
			if cmd.Flags().Changed("sftp-action-send-enabled") {
				paramsSiemHttpDestinationCreate.SftpActionSendEnabled = flib.Bool(createSftpActionSendEnabled)
			}
			if cmd.Flags().Changed("ftp-action-send-enabled") {
				paramsSiemHttpDestinationCreate.FtpActionSendEnabled = flib.Bool(createFtpActionSendEnabled)
			}
			if cmd.Flags().Changed("web-dav-action-send-enabled") {
				paramsSiemHttpDestinationCreate.WebDavActionSendEnabled = flib.Bool(createWebDavActionSendEnabled)
			}
			if cmd.Flags().Changed("sync-send-enabled") {
				paramsSiemHttpDestinationCreate.SyncSendEnabled = flib.Bool(createSyncSendEnabled)
			}
			if cmd.Flags().Changed("outbound-connection-send-enabled") {
				paramsSiemHttpDestinationCreate.OutboundConnectionSendEnabled = flib.Bool(createOutboundConnectionSendEnabled)
			}
			if cmd.Flags().Changed("automation-send-enabled") {
				paramsSiemHttpDestinationCreate.AutomationSendEnabled = flib.Bool(createAutomationSendEnabled)
			}
			if cmd.Flags().Changed("api-request-send-enabled") {
				paramsSiemHttpDestinationCreate.ApiRequestSendEnabled = flib.Bool(createApiRequestSendEnabled)
			}
			if cmd.Flags().Changed("public-hosting-request-send-enabled") {
				paramsSiemHttpDestinationCreate.PublicHostingRequestSendEnabled = flib.Bool(createPublicHostingRequestSendEnabled)
			}
			if cmd.Flags().Changed("email-send-enabled") {
				paramsSiemHttpDestinationCreate.EmailSendEnabled = flib.Bool(createEmailSendEnabled)
			}
			if cmd.Flags().Changed("exavault-api-request-send-enabled") {
				paramsSiemHttpDestinationCreate.ExavaultApiRequestSendEnabled = flib.Bool(createExavaultApiRequestSendEnabled)
			}
			if cmd.Flags().Changed("settings-change-send-enabled") {
				paramsSiemHttpDestinationCreate.SettingsChangeSendEnabled = flib.Bool(createSettingsChangeSendEnabled)
			}

			var siemHttpDestination interface{}
			var err error
			siemHttpDestination, err = client.Create(paramsSiemHttpDestinationCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), siemHttpDestination, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsSiemHttpDestinationCreate.Name, "name", "", "Name for this Destination")
	cmdCreate.Flags().BoolVar(&createSendingActive, "sending-active", createSendingActive, "Whether this SIEM HTTP Destination is currently being sent to or not")
	cmdCreate.Flags().StringVar(&SiemHttpDestinationCreateGenericPayloadType, "generic-payload-type", "", fmt.Sprintf("Applicable only for destination type: generic. Indicates the type of HTTP body. Can be json_newline or json_array. json_newline is multiple log entries as JSON separated by newlines. json_array is a single JSON array containing multiple log entries as JSON. %v", reflect.ValueOf(paramsSiemHttpDestinationCreate.GenericPayloadType.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&paramsSiemHttpDestinationCreate.SplunkToken, "splunk-token", "", "Applicable only for destination type: splunk. Authentication token provided by Splunk.")
	cmdCreate.Flags().StringVar(&paramsSiemHttpDestinationCreate.AzureDcrImmutableId, "azure-dcr-immutable-id", "", "Applicable only for destination types: azure, azure_legacy. Immutable ID of the Data Collection Rule.")
	cmdCreate.Flags().StringVar(&paramsSiemHttpDestinationCreate.AzureStreamName, "azure-stream-name", "", "Applicable only for destination type: azure. Name of the stream in the DCR that represents the destination table.")
	cmdCreate.Flags().StringVar(&paramsSiemHttpDestinationCreate.AzureOauthClientCredentialsTenantId, "azure-oauth-client-credentials-tenant-id", "", "Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Tenant ID.")
	cmdCreate.Flags().StringVar(&paramsSiemHttpDestinationCreate.AzureOauthClientCredentialsClientId, "azure-oauth-client-credentials-client-id", "", "Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Client ID.")
	cmdCreate.Flags().StringVar(&paramsSiemHttpDestinationCreate.AzureOauthClientCredentialsClientSecret, "azure-oauth-client-credentials-client-secret", "", "Applicable only for destination type: azure. Client Credentials OAuth Client Secret.")
	cmdCreate.Flags().StringVar(&paramsSiemHttpDestinationCreate.QradarUsername, "qradar-username", "", "Applicable only for destination type: qradar. Basic auth username provided by QRadar.")
	cmdCreate.Flags().StringVar(&paramsSiemHttpDestinationCreate.QradarPassword, "qradar-password", "", "Applicable only for destination type: qradar. Basic auth password provided by QRadar.")
	cmdCreate.Flags().StringVar(&paramsSiemHttpDestinationCreate.SolarWindsToken, "solar-winds-token", "", "Applicable only for destination type: solar_winds. Authentication token provided by Solar Winds.")
	cmdCreate.Flags().StringVar(&paramsSiemHttpDestinationCreate.NewRelicApiKey, "new-relic-api-key", "", "Applicable only for destination type: new_relic. API key provided by New Relic.")
	cmdCreate.Flags().StringVar(&paramsSiemHttpDestinationCreate.DatadogApiKey, "datadog-api-key", "", "Applicable only for destination type: datadog. API key provided by Datadog.")
	cmdCreate.Flags().BoolVar(&createSftpActionSendEnabled, "sftp-action-send-enabled", createSftpActionSendEnabled, "Whether or not sending is enabled for sftp_action logs.")
	cmdCreate.Flags().BoolVar(&createFtpActionSendEnabled, "ftp-action-send-enabled", createFtpActionSendEnabled, "Whether or not sending is enabled for ftp_action logs.")
	cmdCreate.Flags().BoolVar(&createWebDavActionSendEnabled, "web-dav-action-send-enabled", createWebDavActionSendEnabled, "Whether or not sending is enabled for web_dav_action logs.")
	cmdCreate.Flags().BoolVar(&createSyncSendEnabled, "sync-send-enabled", createSyncSendEnabled, "Whether or not sending is enabled for sync logs.")
	cmdCreate.Flags().BoolVar(&createOutboundConnectionSendEnabled, "outbound-connection-send-enabled", createOutboundConnectionSendEnabled, "Whether or not sending is enabled for outbound_connection logs.")
	cmdCreate.Flags().BoolVar(&createAutomationSendEnabled, "automation-send-enabled", createAutomationSendEnabled, "Whether or not sending is enabled for automation logs.")
	cmdCreate.Flags().BoolVar(&createApiRequestSendEnabled, "api-request-send-enabled", createApiRequestSendEnabled, "Whether or not sending is enabled for api_request logs.")
	cmdCreate.Flags().BoolVar(&createPublicHostingRequestSendEnabled, "public-hosting-request-send-enabled", createPublicHostingRequestSendEnabled, "Whether or not sending is enabled for public_hosting_request logs.")
	cmdCreate.Flags().BoolVar(&createEmailSendEnabled, "email-send-enabled", createEmailSendEnabled, "Whether or not sending is enabled for email logs.")
	cmdCreate.Flags().BoolVar(&createExavaultApiRequestSendEnabled, "exavault-api-request-send-enabled", createExavaultApiRequestSendEnabled, "Whether or not sending is enabled for exavault_api_request logs.")
	cmdCreate.Flags().BoolVar(&createSettingsChangeSendEnabled, "settings-change-send-enabled", createSettingsChangeSendEnabled, "Whether or not sending is enabled for settings_change logs.")
	cmdCreate.Flags().StringVar(&SiemHttpDestinationCreateDestinationType, "destination-type", "", fmt.Sprintf("Destination Type %v", reflect.ValueOf(paramsSiemHttpDestinationCreate.DestinationType.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&paramsSiemHttpDestinationCreate.DestinationUrl, "destination-url", "", "Destination Url")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	SiemHttpDestinations.AddCommand(cmdCreate)
	var fieldsSendTestEntry []string
	var formatSendTestEntry []string
	usePagerSendTestEntry := true
	sendTestEntrySendingActive := true
	sendTestEntrySftpActionSendEnabled := true
	sendTestEntryFtpActionSendEnabled := true
	sendTestEntryWebDavActionSendEnabled := true
	sendTestEntrySyncSendEnabled := true
	sendTestEntryOutboundConnectionSendEnabled := true
	sendTestEntryAutomationSendEnabled := true
	sendTestEntryApiRequestSendEnabled := true
	sendTestEntryPublicHostingRequestSendEnabled := true
	sendTestEntryEmailSendEnabled := true
	sendTestEntryExavaultApiRequestSendEnabled := true
	sendTestEntrySettingsChangeSendEnabled := true
	paramsSiemHttpDestinationSendTestEntry := files_sdk.SiemHttpDestinationSendTestEntryParams{}
	SiemHttpDestinationSendTestEntryDestinationType := ""
	SiemHttpDestinationSendTestEntryGenericPayloadType := ""

	cmdSendTestEntry := &cobra.Command{
		Use:   "send-test-entry",
		Short: `send_test_entry SIEM HTTP Destination`,
		Long:  `send_test_entry SIEM HTTP Destination`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := siem_http_destination.Client{Config: config}

			var SiemHttpDestinationSendTestEntryDestinationTypeErr error
			paramsSiemHttpDestinationSendTestEntry.DestinationType, SiemHttpDestinationSendTestEntryDestinationTypeErr = lib.FetchKey("destination-type", paramsSiemHttpDestinationSendTestEntry.DestinationType.Enum(), SiemHttpDestinationSendTestEntryDestinationType)
			if SiemHttpDestinationSendTestEntryDestinationType != "" && SiemHttpDestinationSendTestEntryDestinationTypeErr != nil {
				return SiemHttpDestinationSendTestEntryDestinationTypeErr
			}
			var SiemHttpDestinationSendTestEntryGenericPayloadTypeErr error
			paramsSiemHttpDestinationSendTestEntry.GenericPayloadType, SiemHttpDestinationSendTestEntryGenericPayloadTypeErr = lib.FetchKey("generic-payload-type", paramsSiemHttpDestinationSendTestEntry.GenericPayloadType.Enum(), SiemHttpDestinationSendTestEntryGenericPayloadType)
			if SiemHttpDestinationSendTestEntryGenericPayloadType != "" && SiemHttpDestinationSendTestEntryGenericPayloadTypeErr != nil {
				return SiemHttpDestinationSendTestEntryGenericPayloadTypeErr
			}

			if cmd.Flags().Changed("sending-active") {
				paramsSiemHttpDestinationSendTestEntry.SendingActive = flib.Bool(sendTestEntrySendingActive)
			}
			if cmd.Flags().Changed("sftp-action-send-enabled") {
				paramsSiemHttpDestinationSendTestEntry.SftpActionSendEnabled = flib.Bool(sendTestEntrySftpActionSendEnabled)
			}
			if cmd.Flags().Changed("ftp-action-send-enabled") {
				paramsSiemHttpDestinationSendTestEntry.FtpActionSendEnabled = flib.Bool(sendTestEntryFtpActionSendEnabled)
			}
			if cmd.Flags().Changed("web-dav-action-send-enabled") {
				paramsSiemHttpDestinationSendTestEntry.WebDavActionSendEnabled = flib.Bool(sendTestEntryWebDavActionSendEnabled)
			}
			if cmd.Flags().Changed("sync-send-enabled") {
				paramsSiemHttpDestinationSendTestEntry.SyncSendEnabled = flib.Bool(sendTestEntrySyncSendEnabled)
			}
			if cmd.Flags().Changed("outbound-connection-send-enabled") {
				paramsSiemHttpDestinationSendTestEntry.OutboundConnectionSendEnabled = flib.Bool(sendTestEntryOutboundConnectionSendEnabled)
			}
			if cmd.Flags().Changed("automation-send-enabled") {
				paramsSiemHttpDestinationSendTestEntry.AutomationSendEnabled = flib.Bool(sendTestEntryAutomationSendEnabled)
			}
			if cmd.Flags().Changed("api-request-send-enabled") {
				paramsSiemHttpDestinationSendTestEntry.ApiRequestSendEnabled = flib.Bool(sendTestEntryApiRequestSendEnabled)
			}
			if cmd.Flags().Changed("public-hosting-request-send-enabled") {
				paramsSiemHttpDestinationSendTestEntry.PublicHostingRequestSendEnabled = flib.Bool(sendTestEntryPublicHostingRequestSendEnabled)
			}
			if cmd.Flags().Changed("email-send-enabled") {
				paramsSiemHttpDestinationSendTestEntry.EmailSendEnabled = flib.Bool(sendTestEntryEmailSendEnabled)
			}
			if cmd.Flags().Changed("exavault-api-request-send-enabled") {
				paramsSiemHttpDestinationSendTestEntry.ExavaultApiRequestSendEnabled = flib.Bool(sendTestEntryExavaultApiRequestSendEnabled)
			}
			if cmd.Flags().Changed("settings-change-send-enabled") {
				paramsSiemHttpDestinationSendTestEntry.SettingsChangeSendEnabled = flib.Bool(sendTestEntrySettingsChangeSendEnabled)
			}

			var err error
			err = client.SendTestEntry(paramsSiemHttpDestinationSendTestEntry, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdSendTestEntry.Flags().Int64Var(&paramsSiemHttpDestinationSendTestEntry.SiemHttpDestinationId, "siem-http-destination-id", 0, "SIEM HTTP Destination ID")
	cmdSendTestEntry.Flags().StringVar(&SiemHttpDestinationSendTestEntryDestinationType, "destination-type", "", fmt.Sprintf("Destination Type %v", reflect.ValueOf(paramsSiemHttpDestinationSendTestEntry.DestinationType.Enum()).MapKeys()))
	cmdSendTestEntry.Flags().StringVar(&paramsSiemHttpDestinationSendTestEntry.DestinationUrl, "destination-url", "", "Destination Url")
	cmdSendTestEntry.Flags().StringVar(&paramsSiemHttpDestinationSendTestEntry.Name, "name", "", "Name for this Destination")
	cmdSendTestEntry.Flags().BoolVar(&sendTestEntrySendingActive, "sending-active", sendTestEntrySendingActive, "Whether this SIEM HTTP Destination is currently being sent to or not")
	cmdSendTestEntry.Flags().StringVar(&SiemHttpDestinationSendTestEntryGenericPayloadType, "generic-payload-type", "", fmt.Sprintf("Applicable only for destination type: generic. Indicates the type of HTTP body. Can be json_newline or json_array. json_newline is multiple log entries as JSON separated by newlines. json_array is a single JSON array containing multiple log entries as JSON. %v", reflect.ValueOf(paramsSiemHttpDestinationSendTestEntry.GenericPayloadType.Enum()).MapKeys()))
	cmdSendTestEntry.Flags().StringVar(&paramsSiemHttpDestinationSendTestEntry.SplunkToken, "splunk-token", "", "Applicable only for destination type: splunk. Authentication token provided by Splunk.")
	cmdSendTestEntry.Flags().StringVar(&paramsSiemHttpDestinationSendTestEntry.AzureDcrImmutableId, "azure-dcr-immutable-id", "", "Applicable only for destination types: azure, azure_legacy. Immutable ID of the Data Collection Rule.")
	cmdSendTestEntry.Flags().StringVar(&paramsSiemHttpDestinationSendTestEntry.AzureStreamName, "azure-stream-name", "", "Applicable only for destination type: azure. Name of the stream in the DCR that represents the destination table.")
	cmdSendTestEntry.Flags().StringVar(&paramsSiemHttpDestinationSendTestEntry.AzureOauthClientCredentialsTenantId, "azure-oauth-client-credentials-tenant-id", "", "Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Tenant ID.")
	cmdSendTestEntry.Flags().StringVar(&paramsSiemHttpDestinationSendTestEntry.AzureOauthClientCredentialsClientId, "azure-oauth-client-credentials-client-id", "", "Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Client ID.")
	cmdSendTestEntry.Flags().StringVar(&paramsSiemHttpDestinationSendTestEntry.AzureOauthClientCredentialsClientSecret, "azure-oauth-client-credentials-client-secret", "", "Applicable only for destination type: azure. Client Credentials OAuth Client Secret.")
	cmdSendTestEntry.Flags().StringVar(&paramsSiemHttpDestinationSendTestEntry.QradarUsername, "qradar-username", "", "Applicable only for destination type: qradar. Basic auth username provided by QRadar.")
	cmdSendTestEntry.Flags().StringVar(&paramsSiemHttpDestinationSendTestEntry.QradarPassword, "qradar-password", "", "Applicable only for destination type: qradar. Basic auth password provided by QRadar.")
	cmdSendTestEntry.Flags().StringVar(&paramsSiemHttpDestinationSendTestEntry.SolarWindsToken, "solar-winds-token", "", "Applicable only for destination type: solar_winds. Authentication token provided by Solar Winds.")
	cmdSendTestEntry.Flags().StringVar(&paramsSiemHttpDestinationSendTestEntry.NewRelicApiKey, "new-relic-api-key", "", "Applicable only for destination type: new_relic. API key provided by New Relic.")
	cmdSendTestEntry.Flags().StringVar(&paramsSiemHttpDestinationSendTestEntry.DatadogApiKey, "datadog-api-key", "", "Applicable only for destination type: datadog. API key provided by Datadog.")
	cmdSendTestEntry.Flags().BoolVar(&sendTestEntrySftpActionSendEnabled, "sftp-action-send-enabled", sendTestEntrySftpActionSendEnabled, "Whether or not sending is enabled for sftp_action logs.")
	cmdSendTestEntry.Flags().BoolVar(&sendTestEntryFtpActionSendEnabled, "ftp-action-send-enabled", sendTestEntryFtpActionSendEnabled, "Whether or not sending is enabled for ftp_action logs.")
	cmdSendTestEntry.Flags().BoolVar(&sendTestEntryWebDavActionSendEnabled, "web-dav-action-send-enabled", sendTestEntryWebDavActionSendEnabled, "Whether or not sending is enabled for web_dav_action logs.")
	cmdSendTestEntry.Flags().BoolVar(&sendTestEntrySyncSendEnabled, "sync-send-enabled", sendTestEntrySyncSendEnabled, "Whether or not sending is enabled for sync logs.")
	cmdSendTestEntry.Flags().BoolVar(&sendTestEntryOutboundConnectionSendEnabled, "outbound-connection-send-enabled", sendTestEntryOutboundConnectionSendEnabled, "Whether or not sending is enabled for outbound_connection logs.")
	cmdSendTestEntry.Flags().BoolVar(&sendTestEntryAutomationSendEnabled, "automation-send-enabled", sendTestEntryAutomationSendEnabled, "Whether or not sending is enabled for automation logs.")
	cmdSendTestEntry.Flags().BoolVar(&sendTestEntryApiRequestSendEnabled, "api-request-send-enabled", sendTestEntryApiRequestSendEnabled, "Whether or not sending is enabled for api_request logs.")
	cmdSendTestEntry.Flags().BoolVar(&sendTestEntryPublicHostingRequestSendEnabled, "public-hosting-request-send-enabled", sendTestEntryPublicHostingRequestSendEnabled, "Whether or not sending is enabled for public_hosting_request logs.")
	cmdSendTestEntry.Flags().BoolVar(&sendTestEntryEmailSendEnabled, "email-send-enabled", sendTestEntryEmailSendEnabled, "Whether or not sending is enabled for email logs.")
	cmdSendTestEntry.Flags().BoolVar(&sendTestEntryExavaultApiRequestSendEnabled, "exavault-api-request-send-enabled", sendTestEntryExavaultApiRequestSendEnabled, "Whether or not sending is enabled for exavault_api_request logs.")
	cmdSendTestEntry.Flags().BoolVar(&sendTestEntrySettingsChangeSendEnabled, "settings-change-send-enabled", sendTestEntrySettingsChangeSendEnabled, "Whether or not sending is enabled for settings_change logs.")

	cmdSendTestEntry.Flags().StringSliceVar(&fieldsSendTestEntry, "fields", []string{}, "comma separated list of field names")
	cmdSendTestEntry.Flags().StringSliceVar(&formatSendTestEntry, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdSendTestEntry.Flags().BoolVar(&usePagerSendTestEntry, "use-pager", usePagerSendTestEntry, "Use $PAGER (.ie less, more, etc)")

	SiemHttpDestinations.AddCommand(cmdSendTestEntry)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateSendingActive := true
	updateSftpActionSendEnabled := true
	updateFtpActionSendEnabled := true
	updateWebDavActionSendEnabled := true
	updateSyncSendEnabled := true
	updateOutboundConnectionSendEnabled := true
	updateAutomationSendEnabled := true
	updateApiRequestSendEnabled := true
	updatePublicHostingRequestSendEnabled := true
	updateEmailSendEnabled := true
	updateExavaultApiRequestSendEnabled := true
	updateSettingsChangeSendEnabled := true
	paramsSiemHttpDestinationUpdate := files_sdk.SiemHttpDestinationUpdateParams{}
	SiemHttpDestinationUpdateGenericPayloadType := ""
	SiemHttpDestinationUpdateDestinationType := ""

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update SIEM HTTP Destination`,
		Long:  `Update SIEM HTTP Destination`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := siem_http_destination.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.SiemHttpDestinationUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			var SiemHttpDestinationUpdateGenericPayloadTypeErr error
			paramsSiemHttpDestinationUpdate.GenericPayloadType, SiemHttpDestinationUpdateGenericPayloadTypeErr = lib.FetchKey("generic-payload-type", paramsSiemHttpDestinationUpdate.GenericPayloadType.Enum(), SiemHttpDestinationUpdateGenericPayloadType)
			if SiemHttpDestinationUpdateGenericPayloadType != "" && SiemHttpDestinationUpdateGenericPayloadTypeErr != nil {
				return SiemHttpDestinationUpdateGenericPayloadTypeErr
			}
			var SiemHttpDestinationUpdateDestinationTypeErr error
			paramsSiemHttpDestinationUpdate.DestinationType, SiemHttpDestinationUpdateDestinationTypeErr = lib.FetchKey("destination-type", paramsSiemHttpDestinationUpdate.DestinationType.Enum(), SiemHttpDestinationUpdateDestinationType)
			if SiemHttpDestinationUpdateDestinationType != "" && SiemHttpDestinationUpdateDestinationTypeErr != nil {
				return SiemHttpDestinationUpdateDestinationTypeErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsSiemHttpDestinationUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsSiemHttpDestinationUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("additional-headers") {
			}
			if cmd.Flags().Changed("sending-active") {
				mapParams["sending_active"] = updateSendingActive
			}
			if cmd.Flags().Changed("generic-payload-type") {
				lib.FlagUpdate(cmd, "generic_payload_type", paramsSiemHttpDestinationUpdate.GenericPayloadType, mapParams)
			}
			if cmd.Flags().Changed("splunk-token") {
				lib.FlagUpdate(cmd, "splunk_token", paramsSiemHttpDestinationUpdate.SplunkToken, mapParams)
			}
			if cmd.Flags().Changed("azure-dcr-immutable-id") {
				lib.FlagUpdate(cmd, "azure_dcr_immutable_id", paramsSiemHttpDestinationUpdate.AzureDcrImmutableId, mapParams)
			}
			if cmd.Flags().Changed("azure-stream-name") {
				lib.FlagUpdate(cmd, "azure_stream_name", paramsSiemHttpDestinationUpdate.AzureStreamName, mapParams)
			}
			if cmd.Flags().Changed("azure-oauth-client-credentials-tenant-id") {
				lib.FlagUpdate(cmd, "azure_oauth_client_credentials_tenant_id", paramsSiemHttpDestinationUpdate.AzureOauthClientCredentialsTenantId, mapParams)
			}
			if cmd.Flags().Changed("azure-oauth-client-credentials-client-id") {
				lib.FlagUpdate(cmd, "azure_oauth_client_credentials_client_id", paramsSiemHttpDestinationUpdate.AzureOauthClientCredentialsClientId, mapParams)
			}
			if cmd.Flags().Changed("azure-oauth-client-credentials-client-secret") {
				lib.FlagUpdate(cmd, "azure_oauth_client_credentials_client_secret", paramsSiemHttpDestinationUpdate.AzureOauthClientCredentialsClientSecret, mapParams)
			}
			if cmd.Flags().Changed("qradar-username") {
				lib.FlagUpdate(cmd, "qradar_username", paramsSiemHttpDestinationUpdate.QradarUsername, mapParams)
			}
			if cmd.Flags().Changed("qradar-password") {
				lib.FlagUpdate(cmd, "qradar_password", paramsSiemHttpDestinationUpdate.QradarPassword, mapParams)
			}
			if cmd.Flags().Changed("solar-winds-token") {
				lib.FlagUpdate(cmd, "solar_winds_token", paramsSiemHttpDestinationUpdate.SolarWindsToken, mapParams)
			}
			if cmd.Flags().Changed("new-relic-api-key") {
				lib.FlagUpdate(cmd, "new_relic_api_key", paramsSiemHttpDestinationUpdate.NewRelicApiKey, mapParams)
			}
			if cmd.Flags().Changed("datadog-api-key") {
				lib.FlagUpdate(cmd, "datadog_api_key", paramsSiemHttpDestinationUpdate.DatadogApiKey, mapParams)
			}
			if cmd.Flags().Changed("sftp-action-send-enabled") {
				mapParams["sftp_action_send_enabled"] = updateSftpActionSendEnabled
			}
			if cmd.Flags().Changed("ftp-action-send-enabled") {
				mapParams["ftp_action_send_enabled"] = updateFtpActionSendEnabled
			}
			if cmd.Flags().Changed("web-dav-action-send-enabled") {
				mapParams["web_dav_action_send_enabled"] = updateWebDavActionSendEnabled
			}
			if cmd.Flags().Changed("sync-send-enabled") {
				mapParams["sync_send_enabled"] = updateSyncSendEnabled
			}
			if cmd.Flags().Changed("outbound-connection-send-enabled") {
				mapParams["outbound_connection_send_enabled"] = updateOutboundConnectionSendEnabled
			}
			if cmd.Flags().Changed("automation-send-enabled") {
				mapParams["automation_send_enabled"] = updateAutomationSendEnabled
			}
			if cmd.Flags().Changed("api-request-send-enabled") {
				mapParams["api_request_send_enabled"] = updateApiRequestSendEnabled
			}
			if cmd.Flags().Changed("public-hosting-request-send-enabled") {
				mapParams["public_hosting_request_send_enabled"] = updatePublicHostingRequestSendEnabled
			}
			if cmd.Flags().Changed("email-send-enabled") {
				mapParams["email_send_enabled"] = updateEmailSendEnabled
			}
			if cmd.Flags().Changed("exavault-api-request-send-enabled") {
				mapParams["exavault_api_request_send_enabled"] = updateExavaultApiRequestSendEnabled
			}
			if cmd.Flags().Changed("settings-change-send-enabled") {
				mapParams["settings_change_send_enabled"] = updateSettingsChangeSendEnabled
			}
			if cmd.Flags().Changed("destination-type") {
				lib.FlagUpdate(cmd, "destination_type", paramsSiemHttpDestinationUpdate.DestinationType, mapParams)
			}
			if cmd.Flags().Changed("destination-url") {
				lib.FlagUpdate(cmd, "destination_url", paramsSiemHttpDestinationUpdate.DestinationUrl, mapParams)
			}

			var siemHttpDestination interface{}
			var err error
			siemHttpDestination, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), siemHttpDestination, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsSiemHttpDestinationUpdate.Id, "id", 0, "Siem Http Destination ID.")
	cmdUpdate.Flags().StringVar(&paramsSiemHttpDestinationUpdate.Name, "name", "", "Name for this Destination")
	cmdUpdate.Flags().BoolVar(&updateSendingActive, "sending-active", updateSendingActive, "Whether this SIEM HTTP Destination is currently being sent to or not")
	cmdUpdate.Flags().StringVar(&SiemHttpDestinationUpdateGenericPayloadType, "generic-payload-type", "", fmt.Sprintf("Applicable only for destination type: generic. Indicates the type of HTTP body. Can be json_newline or json_array. json_newline is multiple log entries as JSON separated by newlines. json_array is a single JSON array containing multiple log entries as JSON. %v", reflect.ValueOf(paramsSiemHttpDestinationUpdate.GenericPayloadType.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&paramsSiemHttpDestinationUpdate.SplunkToken, "splunk-token", "", "Applicable only for destination type: splunk. Authentication token provided by Splunk.")
	cmdUpdate.Flags().StringVar(&paramsSiemHttpDestinationUpdate.AzureDcrImmutableId, "azure-dcr-immutable-id", "", "Applicable only for destination types: azure, azure_legacy. Immutable ID of the Data Collection Rule.")
	cmdUpdate.Flags().StringVar(&paramsSiemHttpDestinationUpdate.AzureStreamName, "azure-stream-name", "", "Applicable only for destination type: azure. Name of the stream in the DCR that represents the destination table.")
	cmdUpdate.Flags().StringVar(&paramsSiemHttpDestinationUpdate.AzureOauthClientCredentialsTenantId, "azure-oauth-client-credentials-tenant-id", "", "Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Tenant ID.")
	cmdUpdate.Flags().StringVar(&paramsSiemHttpDestinationUpdate.AzureOauthClientCredentialsClientId, "azure-oauth-client-credentials-client-id", "", "Applicable only for destination types: azure, azure_legacy. Client Credentials OAuth Client ID.")
	cmdUpdate.Flags().StringVar(&paramsSiemHttpDestinationUpdate.AzureOauthClientCredentialsClientSecret, "azure-oauth-client-credentials-client-secret", "", "Applicable only for destination type: azure. Client Credentials OAuth Client Secret.")
	cmdUpdate.Flags().StringVar(&paramsSiemHttpDestinationUpdate.QradarUsername, "qradar-username", "", "Applicable only for destination type: qradar. Basic auth username provided by QRadar.")
	cmdUpdate.Flags().StringVar(&paramsSiemHttpDestinationUpdate.QradarPassword, "qradar-password", "", "Applicable only for destination type: qradar. Basic auth password provided by QRadar.")
	cmdUpdate.Flags().StringVar(&paramsSiemHttpDestinationUpdate.SolarWindsToken, "solar-winds-token", "", "Applicable only for destination type: solar_winds. Authentication token provided by Solar Winds.")
	cmdUpdate.Flags().StringVar(&paramsSiemHttpDestinationUpdate.NewRelicApiKey, "new-relic-api-key", "", "Applicable only for destination type: new_relic. API key provided by New Relic.")
	cmdUpdate.Flags().StringVar(&paramsSiemHttpDestinationUpdate.DatadogApiKey, "datadog-api-key", "", "Applicable only for destination type: datadog. API key provided by Datadog.")
	cmdUpdate.Flags().BoolVar(&updateSftpActionSendEnabled, "sftp-action-send-enabled", updateSftpActionSendEnabled, "Whether or not sending is enabled for sftp_action logs.")
	cmdUpdate.Flags().BoolVar(&updateFtpActionSendEnabled, "ftp-action-send-enabled", updateFtpActionSendEnabled, "Whether or not sending is enabled for ftp_action logs.")
	cmdUpdate.Flags().BoolVar(&updateWebDavActionSendEnabled, "web-dav-action-send-enabled", updateWebDavActionSendEnabled, "Whether or not sending is enabled for web_dav_action logs.")
	cmdUpdate.Flags().BoolVar(&updateSyncSendEnabled, "sync-send-enabled", updateSyncSendEnabled, "Whether or not sending is enabled for sync logs.")
	cmdUpdate.Flags().BoolVar(&updateOutboundConnectionSendEnabled, "outbound-connection-send-enabled", updateOutboundConnectionSendEnabled, "Whether or not sending is enabled for outbound_connection logs.")
	cmdUpdate.Flags().BoolVar(&updateAutomationSendEnabled, "automation-send-enabled", updateAutomationSendEnabled, "Whether or not sending is enabled for automation logs.")
	cmdUpdate.Flags().BoolVar(&updateApiRequestSendEnabled, "api-request-send-enabled", updateApiRequestSendEnabled, "Whether or not sending is enabled for api_request logs.")
	cmdUpdate.Flags().BoolVar(&updatePublicHostingRequestSendEnabled, "public-hosting-request-send-enabled", updatePublicHostingRequestSendEnabled, "Whether or not sending is enabled for public_hosting_request logs.")
	cmdUpdate.Flags().BoolVar(&updateEmailSendEnabled, "email-send-enabled", updateEmailSendEnabled, "Whether or not sending is enabled for email logs.")
	cmdUpdate.Flags().BoolVar(&updateExavaultApiRequestSendEnabled, "exavault-api-request-send-enabled", updateExavaultApiRequestSendEnabled, "Whether or not sending is enabled for exavault_api_request logs.")
	cmdUpdate.Flags().BoolVar(&updateSettingsChangeSendEnabled, "settings-change-send-enabled", updateSettingsChangeSendEnabled, "Whether or not sending is enabled for settings_change logs.")
	cmdUpdate.Flags().StringVar(&SiemHttpDestinationUpdateDestinationType, "destination-type", "", fmt.Sprintf("Destination Type %v", reflect.ValueOf(paramsSiemHttpDestinationUpdate.DestinationType.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&paramsSiemHttpDestinationUpdate.DestinationUrl, "destination-url", "", "Destination Url")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	SiemHttpDestinations.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsSiemHttpDestinationDelete := files_sdk.SiemHttpDestinationDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete SIEM HTTP Destination`,
		Long:  `Delete SIEM HTTP Destination`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := siem_http_destination.Client{Config: config}

			var err error
			err = client.Delete(paramsSiemHttpDestinationDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsSiemHttpDestinationDelete.Id, "id", 0, "Siem Http Destination ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	SiemHttpDestinations.AddCommand(cmdDelete)
	return SiemHttpDestinations
}
