package cmd

import (
	"fmt"
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	remote_server "github.com/Files-com/files-sdk-go/v3/remoteserver"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(RemoteServers())
}

func RemoteServers() *cobra.Command {
	RemoteServers := &cobra.Command{
		Use:  "remote-servers [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command remote-servers\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsRemoteServerList := files_sdk.RemoteServerListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Remote Servers",
		Long:    `List Remote Servers`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsRemoteServerList
			params.MaxPages = MaxPagesList

			client := remote_server.Client{Config: config}
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

	cmdList.Flags().StringVar(&paramsRemoteServerList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsRemoteServerList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	RemoteServers.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsRemoteServerFind := files_sdk.RemoteServerFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Remote Server`,
		Long:  `Show Remote Server`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := remote_server.Client{Config: config}

			var remoteServer interface{}
			var err error
			remoteServer, err = client.Find(paramsRemoteServerFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), remoteServer, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsRemoteServerFind.Id, "id", 0, "Remote Server ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	RemoteServers.AddCommand(cmdFind)
	var fieldsFindConfigurationFile []string
	var formatFindConfigurationFile []string
	usePagerFindConfigurationFile := true
	paramsRemoteServerFindConfigurationFile := files_sdk.RemoteServerFindConfigurationFileParams{}

	cmdFindConfigurationFile := &cobra.Command{
		Use:   "find-configuration-file",
		Short: `Download configuration file (required for some Remote Server integrations, such as the Files.com Agent)`,
		Long:  `Download configuration file (required for some Remote Server integrations, such as the Files.com Agent)`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := remote_server.Client{Config: config}

			var remoteServerConfigurationFile interface{}
			var err error
			remoteServerConfigurationFile, err = client.FindConfigurationFile(paramsRemoteServerFindConfigurationFile, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), remoteServerConfigurationFile, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFindConfigurationFile), fieldsFindConfigurationFile, usePagerFindConfigurationFile, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFindConfigurationFile.Flags().Int64Var(&paramsRemoteServerFindConfigurationFile.Id, "id", 0, "Remote Server ID.")

	cmdFindConfigurationFile.Flags().StringSliceVar(&fieldsFindConfigurationFile, "fields", []string{}, "comma separated list of field names")
	cmdFindConfigurationFile.Flags().StringSliceVar(&formatFindConfigurationFile, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFindConfigurationFile.Flags().BoolVar(&usePagerFindConfigurationFile, "use-pager", usePagerFindConfigurationFile, "Use $PAGER (.ie less, more, etc)")

	RemoteServers.AddCommand(cmdFindConfigurationFile)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createResetAuthentication := true
	createAzureBlobStorageHierarchicalNamespace := true
	createDropboxTeams := true
	createEnableDedicatedIps := true
	createPinToSiteRegion := true
	paramsRemoteServerCreate := files_sdk.RemoteServerCreateParams{}
	RemoteServerCreateBufferUploads := ""
	RemoteServerCreateFilesAgentPermissionSet := ""
	RemoteServerCreateOneDriveAccountType := ""
	RemoteServerCreateServerCertificate := ""
	RemoteServerCreateServerType := ""
	RemoteServerCreateSsl := ""

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Remote Server`,
		Long:  `Create Remote Server`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := remote_server.Client{Config: config}

			var RemoteServerCreateBufferUploadsErr error
			paramsRemoteServerCreate.BufferUploads, RemoteServerCreateBufferUploadsErr = lib.FetchKey("buffer-uploads", paramsRemoteServerCreate.BufferUploads.Enum(), RemoteServerCreateBufferUploads)
			if RemoteServerCreateBufferUploads != "" && RemoteServerCreateBufferUploadsErr != nil {
				return RemoteServerCreateBufferUploadsErr
			}
			var RemoteServerCreateFilesAgentPermissionSetErr error
			paramsRemoteServerCreate.FilesAgentPermissionSet, RemoteServerCreateFilesAgentPermissionSetErr = lib.FetchKey("files-agent-permission-set", paramsRemoteServerCreate.FilesAgentPermissionSet.Enum(), RemoteServerCreateFilesAgentPermissionSet)
			if RemoteServerCreateFilesAgentPermissionSet != "" && RemoteServerCreateFilesAgentPermissionSetErr != nil {
				return RemoteServerCreateFilesAgentPermissionSetErr
			}
			var RemoteServerCreateOneDriveAccountTypeErr error
			paramsRemoteServerCreate.OneDriveAccountType, RemoteServerCreateOneDriveAccountTypeErr = lib.FetchKey("one-drive-account-type", paramsRemoteServerCreate.OneDriveAccountType.Enum(), RemoteServerCreateOneDriveAccountType)
			if RemoteServerCreateOneDriveAccountType != "" && RemoteServerCreateOneDriveAccountTypeErr != nil {
				return RemoteServerCreateOneDriveAccountTypeErr
			}
			var RemoteServerCreateServerCertificateErr error
			paramsRemoteServerCreate.ServerCertificate, RemoteServerCreateServerCertificateErr = lib.FetchKey("server-certificate", paramsRemoteServerCreate.ServerCertificate.Enum(), RemoteServerCreateServerCertificate)
			if RemoteServerCreateServerCertificate != "" && RemoteServerCreateServerCertificateErr != nil {
				return RemoteServerCreateServerCertificateErr
			}
			var RemoteServerCreateServerTypeErr error
			paramsRemoteServerCreate.ServerType, RemoteServerCreateServerTypeErr = lib.FetchKey("server-type", paramsRemoteServerCreate.ServerType.Enum(), RemoteServerCreateServerType)
			if RemoteServerCreateServerType != "" && RemoteServerCreateServerTypeErr != nil {
				return RemoteServerCreateServerTypeErr
			}
			var RemoteServerCreateSslErr error
			paramsRemoteServerCreate.Ssl, RemoteServerCreateSslErr = lib.FetchKey("ssl", paramsRemoteServerCreate.Ssl.Enum(), RemoteServerCreateSsl)
			if RemoteServerCreateSsl != "" && RemoteServerCreateSslErr != nil {
				return RemoteServerCreateSslErr
			}

			if cmd.Flags().Changed("reset-authentication") {
				paramsRemoteServerCreate.ResetAuthentication = flib.Bool(createResetAuthentication)
			}
			if cmd.Flags().Changed("azure-blob-storage-hierarchical-namespace") {
				paramsRemoteServerCreate.AzureBlobStorageHierarchicalNamespace = flib.Bool(createAzureBlobStorageHierarchicalNamespace)
			}
			if cmd.Flags().Changed("dropbox-teams") {
				paramsRemoteServerCreate.DropboxTeams = flib.Bool(createDropboxTeams)
			}
			if cmd.Flags().Changed("enable-dedicated-ips") {
				paramsRemoteServerCreate.EnableDedicatedIps = flib.Bool(createEnableDedicatedIps)
			}
			if cmd.Flags().Changed("pin-to-site-region") {
				paramsRemoteServerCreate.PinToSiteRegion = flib.Bool(createPinToSiteRegion)
			}

			var remoteServer interface{}
			var err error
			remoteServer, err = client.Create(paramsRemoteServerCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), remoteServer, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.Password, "password", "", "Password, if needed.")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.PrivateKey, "private-key", "", "Private key, if needed.")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.PrivateKeyPassphrase, "private-key-passphrase", "", "Passphrase for private key if needed.")
	cmdCreate.Flags().BoolVar(&createResetAuthentication, "reset-authentication", createResetAuthentication, "Reset authenticated account?")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.SslCertificate, "ssl-certificate", "", "SSL client certificate.")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.AwsSecretKey, "aws-secret-key", "", "AWS: secret key.")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.AzureBlobStorageAccessKey, "azure-blob-storage-access-key", "", "Azure Blob Storage: Access Key")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.AzureBlobStorageSasToken, "azure-blob-storage-sas-token", "", "Azure Blob Storage: Shared Access Signature (SAS) token")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.AzureFilesStorageAccessKey, "azure-files-storage-access-key", "", "Azure File Storage: Access Key")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.AzureFilesStorageSasToken, "azure-files-storage-sas-token", "", "Azure File Storage: Shared Access Signature (SAS) token")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.BackblazeB2ApplicationKey, "backblaze-b2-application-key", "", "Backblaze B2 Cloud Storage: applicationKey")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.BackblazeB2KeyId, "backblaze-b2-key-id", "", "Backblaze B2 Cloud Storage: keyID")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.CloudflareSecretKey, "cloudflare-secret-key", "", "Cloudflare: Secret Key")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.FilebaseSecretKey, "filebase-secret-key", "", "Filebase: Secret Key")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.GoogleCloudStorageCredentialsJson, "google-cloud-storage-credentials-json", "", "Google Cloud Storage: JSON file that contains the private key. To generate see https://cloud.google.com/storage/docs/json_api/v1/how-tos/authorizing#APIKey")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.GoogleCloudStorageS3CompatibleSecretKey, "google-cloud-storage-s3-compatible-secret-key", "", "Google Cloud Storage: S3-compatible secret key")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.LinodeSecretKey, "linode-secret-key", "", "Linode: Secret Key")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.S3CompatibleSecretKey, "s3-compatible-secret-key", "", "S3-compatible: Secret Key")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.WasabiSecretKey, "wasabi-secret-key", "", "Wasabi: Secret Key")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.AwsAccessKey, "aws-access-key", "", "AWS Access Key.")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.AzureBlobStorageAccount, "azure-blob-storage-account", "", "Azure Blob Storage: Account name")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.AzureBlobStorageContainer, "azure-blob-storage-container", "", "Azure Blob Storage: Container name")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.AzureBlobStorageDnsSuffix, "azure-blob-storage-dns-suffix", "", "Azure Blob Storage: Custom DNS suffix")
	cmdCreate.Flags().BoolVar(&createAzureBlobStorageHierarchicalNamespace, "azure-blob-storage-hierarchical-namespace", createAzureBlobStorageHierarchicalNamespace, "Azure Blob Storage: Does the storage account has hierarchical namespace feature enabled?")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.AzureFilesStorageAccount, "azure-files-storage-account", "", "Azure Files: Storage Account name")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.AzureFilesStorageDnsSuffix, "azure-files-storage-dns-suffix", "", "Azure Files: Custom DNS suffix")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.AzureFilesStorageShareName, "azure-files-storage-share-name", "", "Azure Files:  Storage Share name")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.BackblazeB2Bucket, "backblaze-b2-bucket", "", "Backblaze B2 Cloud Storage: Bucket name")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.BackblazeB2S3Endpoint, "backblaze-b2-s3-endpoint", "", "Backblaze B2 Cloud Storage: S3 Endpoint")
	cmdCreate.Flags().StringVar(&RemoteServerCreateBufferUploads, "buffer-uploads", "", fmt.Sprintf("If set to always, uploads to this server will be uploaded first to Files.com before being sent to the remote server. This can improve performance in certain access patterns, such as high-latency connections.  It will cause data to be temporarily stored in Files.com. If set to auto, we will perform this optimization if we believe it to be a benefit in a given situation. %v", reflect.ValueOf(paramsRemoteServerCreate.BufferUploads.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.CloudflareAccessKey, "cloudflare-access-key", "", "Cloudflare: Access Key.")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.CloudflareBucket, "cloudflare-bucket", "", "Cloudflare: Bucket name")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.CloudflareEndpoint, "cloudflare-endpoint", "", "Cloudflare: endpoint")
	cmdCreate.Flags().BoolVar(&createDropboxTeams, "dropbox-teams", createDropboxTeams, "Dropbox: If true, list Team folders in root?")
	cmdCreate.Flags().BoolVar(&createEnableDedicatedIps, "enable-dedicated-ips", createEnableDedicatedIps, "`true` if remote server only accepts connections from dedicated IPs")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.FilebaseAccessKey, "filebase-access-key", "", "Filebase: Access Key.")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.FilebaseBucket, "filebase-bucket", "", "Filebase: Bucket name")
	cmdCreate.Flags().StringVar(&RemoteServerCreateFilesAgentPermissionSet, "files-agent-permission-set", "", fmt.Sprintf("Local permissions for files agent. read_only, write_only, or read_write %v", reflect.ValueOf(paramsRemoteServerCreate.FilesAgentPermissionSet.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.FilesAgentRoot, "files-agent-root", "", "Agent local root path")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.FilesAgentVersion, "files-agent-version", "", "Files Agent version")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.GoogleCloudStorageBucket, "google-cloud-storage-bucket", "", "Google Cloud Storage: Bucket Name")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.GoogleCloudStorageProjectId, "google-cloud-storage-project-id", "", "Google Cloud Storage: Project ID")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.GoogleCloudStorageS3CompatibleAccessKey, "google-cloud-storage-s3-compatible-access-key", "", "Google Cloud Storage: S3-compatible Access Key.")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.Hostname, "hostname", "", "Hostname or IP address")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.LinodeAccessKey, "linode-access-key", "", "Linode: Access Key")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.LinodeBucket, "linode-bucket", "", "Linode: Bucket name")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.LinodeRegion, "linode-region", "", "Linode: region")
	cmdCreate.Flags().Int64Var(&paramsRemoteServerCreate.MaxConnections, "max-connections", 0, "Max number of parallel connections.  Ignored for S3 connections (we will parallelize these as much as possible).")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.Name, "name", "", "Internal name for your reference")
	cmdCreate.Flags().StringVar(&RemoteServerCreateOneDriveAccountType, "one-drive-account-type", "", fmt.Sprintf("OneDrive: Either personal or business_other account types %v", reflect.ValueOf(paramsRemoteServerCreate.OneDriveAccountType.Enum()).MapKeys()))
	cmdCreate.Flags().BoolVar(&createPinToSiteRegion, "pin-to-site-region", createPinToSiteRegion, "If true, we will ensure that all communications with this remote server are made through the primary region of the site.  This setting can also be overridden by a site-wide setting which will force it to true.")
	cmdCreate.Flags().Int64Var(&paramsRemoteServerCreate.Port, "port", 0, "Port for remote server.  Not needed for S3.")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.S3Bucket, "s3-bucket", "", "S3 bucket name")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.S3CompatibleAccessKey, "s3-compatible-access-key", "", "S3-compatible: Access Key")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.S3CompatibleBucket, "s3-compatible-bucket", "", "S3-compatible: Bucket name")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.S3CompatibleEndpoint, "s3-compatible-endpoint", "", "S3-compatible: endpoint")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.S3CompatibleRegion, "s3-compatible-region", "", "S3-compatible: region")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.S3Region, "s3-region", "", "S3 region")
	cmdCreate.Flags().StringVar(&RemoteServerCreateServerCertificate, "server-certificate", "", fmt.Sprintf("Remote server certificate %v", reflect.ValueOf(paramsRemoteServerCreate.ServerCertificate.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.ServerHostKey, "server-host-key", "", "Remote server SSH Host Key. If provided, we will require that the server host key matches the provided key. Uses OpenSSH format similar to what would go into ~/.ssh/known_hosts")
	cmdCreate.Flags().StringVar(&RemoteServerCreateServerType, "server-type", "", fmt.Sprintf("Remote server type. %v", reflect.ValueOf(paramsRemoteServerCreate.ServerType.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&RemoteServerCreateSsl, "ssl", "", fmt.Sprintf("Should we require SSL? %v", reflect.ValueOf(paramsRemoteServerCreate.Ssl.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.Username, "username", "", "Remote server username.  Not needed for S3 buckets.")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.WasabiAccessKey, "wasabi-access-key", "", "Wasabi: Access Key.")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.WasabiBucket, "wasabi-bucket", "", "Wasabi: Bucket name")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCreate.WasabiRegion, "wasabi-region", "", "Wasabi: Region")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	RemoteServers.AddCommand(cmdCreate)
	var fieldsConfigurationFile []string
	var formatConfigurationFile []string
	usePagerConfigurationFile := true
	paramsRemoteServerConfigurationFile := files_sdk.RemoteServerConfigurationFileParams{}

	cmdConfigurationFile := &cobra.Command{
		Use:   "configuration-file",
		Short: `Post local changes, check in, and download configuration file (used by some Remote Server integrations, such as the Files.com Agent)`,
		Long:  `Post local changes, check in, and download configuration file (used by some Remote Server integrations, such as the Files.com Agent)`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := remote_server.Client{Config: config}

			var remoteServerConfigurationFile interface{}
			var err error
			remoteServerConfigurationFile, err = client.ConfigurationFile(paramsRemoteServerConfigurationFile, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), remoteServerConfigurationFile, err, Profile(cmd).Current().SetResourceFormat(cmd, formatConfigurationFile), fieldsConfigurationFile, usePagerConfigurationFile, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdConfigurationFile.Flags().Int64Var(&paramsRemoteServerConfigurationFile.Id, "id", 0, "Remote Server ID.")
	cmdConfigurationFile.Flags().StringVar(&paramsRemoteServerConfigurationFile.ApiToken, "api-token", "", "Files Agent API Token")
	cmdConfigurationFile.Flags().StringVar(&paramsRemoteServerConfigurationFile.PermissionSet, "permission-set", "", "The permission set for the agent ['read_write', 'read_only', 'write_only']")
	cmdConfigurationFile.Flags().StringVar(&paramsRemoteServerConfigurationFile.Root, "root", "", "The root directory for the agent")
	cmdConfigurationFile.Flags().StringVar(&paramsRemoteServerConfigurationFile.Hostname, "hostname", "", "")
	cmdConfigurationFile.Flags().Int64Var(&paramsRemoteServerConfigurationFile.Port, "port", 0, "Incoming port for files agent connections")
	cmdConfigurationFile.Flags().StringVar(&paramsRemoteServerConfigurationFile.Status, "status", "", "either running or shutdown")
	cmdConfigurationFile.Flags().StringVar(&paramsRemoteServerConfigurationFile.ConfigVersion, "config-version", "", "agent config version")
	cmdConfigurationFile.Flags().StringVar(&paramsRemoteServerConfigurationFile.PrivateKey, "private-key", "", "The private key for the agent")
	cmdConfigurationFile.Flags().StringVar(&paramsRemoteServerConfigurationFile.PublicKey, "public-key", "", "public key")
	cmdConfigurationFile.Flags().StringVar(&paramsRemoteServerConfigurationFile.ServerHostKey, "server-host-key", "", "")
	cmdConfigurationFile.Flags().StringVar(&paramsRemoteServerConfigurationFile.Subdomain, "subdomain", "", "Files.com subdomain site name")

	cmdConfigurationFile.Flags().StringSliceVar(&fieldsConfigurationFile, "fields", []string{}, "comma separated list of field names")
	cmdConfigurationFile.Flags().StringSliceVar(&formatConfigurationFile, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdConfigurationFile.Flags().BoolVar(&usePagerConfigurationFile, "use-pager", usePagerConfigurationFile, "Use $PAGER (.ie less, more, etc)")

	RemoteServers.AddCommand(cmdConfigurationFile)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateResetAuthentication := true
	updateAzureBlobStorageHierarchicalNamespace := true
	updateDropboxTeams := true
	updateEnableDedicatedIps := true
	updatePinToSiteRegion := true
	paramsRemoteServerUpdate := files_sdk.RemoteServerUpdateParams{}
	RemoteServerUpdateBufferUploads := ""
	RemoteServerUpdateFilesAgentPermissionSet := ""
	RemoteServerUpdateOneDriveAccountType := ""
	RemoteServerUpdateServerCertificate := ""
	RemoteServerUpdateServerType := ""
	RemoteServerUpdateSsl := ""

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Remote Server`,
		Long:  `Update Remote Server`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := remote_server.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.RemoteServerUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			var RemoteServerUpdateBufferUploadsErr error
			paramsRemoteServerUpdate.BufferUploads, RemoteServerUpdateBufferUploadsErr = lib.FetchKey("buffer-uploads", paramsRemoteServerUpdate.BufferUploads.Enum(), RemoteServerUpdateBufferUploads)
			if RemoteServerUpdateBufferUploads != "" && RemoteServerUpdateBufferUploadsErr != nil {
				return RemoteServerUpdateBufferUploadsErr
			}
			var RemoteServerUpdateFilesAgentPermissionSetErr error
			paramsRemoteServerUpdate.FilesAgentPermissionSet, RemoteServerUpdateFilesAgentPermissionSetErr = lib.FetchKey("files-agent-permission-set", paramsRemoteServerUpdate.FilesAgentPermissionSet.Enum(), RemoteServerUpdateFilesAgentPermissionSet)
			if RemoteServerUpdateFilesAgentPermissionSet != "" && RemoteServerUpdateFilesAgentPermissionSetErr != nil {
				return RemoteServerUpdateFilesAgentPermissionSetErr
			}
			var RemoteServerUpdateOneDriveAccountTypeErr error
			paramsRemoteServerUpdate.OneDriveAccountType, RemoteServerUpdateOneDriveAccountTypeErr = lib.FetchKey("one-drive-account-type", paramsRemoteServerUpdate.OneDriveAccountType.Enum(), RemoteServerUpdateOneDriveAccountType)
			if RemoteServerUpdateOneDriveAccountType != "" && RemoteServerUpdateOneDriveAccountTypeErr != nil {
				return RemoteServerUpdateOneDriveAccountTypeErr
			}
			var RemoteServerUpdateServerCertificateErr error
			paramsRemoteServerUpdate.ServerCertificate, RemoteServerUpdateServerCertificateErr = lib.FetchKey("server-certificate", paramsRemoteServerUpdate.ServerCertificate.Enum(), RemoteServerUpdateServerCertificate)
			if RemoteServerUpdateServerCertificate != "" && RemoteServerUpdateServerCertificateErr != nil {
				return RemoteServerUpdateServerCertificateErr
			}
			var RemoteServerUpdateServerTypeErr error
			paramsRemoteServerUpdate.ServerType, RemoteServerUpdateServerTypeErr = lib.FetchKey("server-type", paramsRemoteServerUpdate.ServerType.Enum(), RemoteServerUpdateServerType)
			if RemoteServerUpdateServerType != "" && RemoteServerUpdateServerTypeErr != nil {
				return RemoteServerUpdateServerTypeErr
			}
			var RemoteServerUpdateSslErr error
			paramsRemoteServerUpdate.Ssl, RemoteServerUpdateSslErr = lib.FetchKey("ssl", paramsRemoteServerUpdate.Ssl.Enum(), RemoteServerUpdateSsl)
			if RemoteServerUpdateSsl != "" && RemoteServerUpdateSslErr != nil {
				return RemoteServerUpdateSslErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsRemoteServerUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("password") {
				lib.FlagUpdate(cmd, "password", paramsRemoteServerUpdate.Password, mapParams)
			}
			if cmd.Flags().Changed("private-key") {
				lib.FlagUpdate(cmd, "private_key", paramsRemoteServerUpdate.PrivateKey, mapParams)
			}
			if cmd.Flags().Changed("private-key-passphrase") {
				lib.FlagUpdate(cmd, "private_key_passphrase", paramsRemoteServerUpdate.PrivateKeyPassphrase, mapParams)
			}
			if cmd.Flags().Changed("reset-authentication") {
				mapParams["reset_authentication"] = updateResetAuthentication
			}
			if cmd.Flags().Changed("ssl-certificate") {
				lib.FlagUpdate(cmd, "ssl_certificate", paramsRemoteServerUpdate.SslCertificate, mapParams)
			}
			if cmd.Flags().Changed("aws-secret-key") {
				lib.FlagUpdate(cmd, "aws_secret_key", paramsRemoteServerUpdate.AwsSecretKey, mapParams)
			}
			if cmd.Flags().Changed("azure-blob-storage-access-key") {
				lib.FlagUpdate(cmd, "azure_blob_storage_access_key", paramsRemoteServerUpdate.AzureBlobStorageAccessKey, mapParams)
			}
			if cmd.Flags().Changed("azure-blob-storage-sas-token") {
				lib.FlagUpdate(cmd, "azure_blob_storage_sas_token", paramsRemoteServerUpdate.AzureBlobStorageSasToken, mapParams)
			}
			if cmd.Flags().Changed("azure-files-storage-access-key") {
				lib.FlagUpdate(cmd, "azure_files_storage_access_key", paramsRemoteServerUpdate.AzureFilesStorageAccessKey, mapParams)
			}
			if cmd.Flags().Changed("azure-files-storage-sas-token") {
				lib.FlagUpdate(cmd, "azure_files_storage_sas_token", paramsRemoteServerUpdate.AzureFilesStorageSasToken, mapParams)
			}
			if cmd.Flags().Changed("backblaze-b2-application-key") {
				lib.FlagUpdate(cmd, "backblaze_b2_application_key", paramsRemoteServerUpdate.BackblazeB2ApplicationKey, mapParams)
			}
			if cmd.Flags().Changed("backblaze-b2-key-id") {
				lib.FlagUpdate(cmd, "backblaze_b2_key_id", paramsRemoteServerUpdate.BackblazeB2KeyId, mapParams)
			}
			if cmd.Flags().Changed("cloudflare-secret-key") {
				lib.FlagUpdate(cmd, "cloudflare_secret_key", paramsRemoteServerUpdate.CloudflareSecretKey, mapParams)
			}
			if cmd.Flags().Changed("filebase-secret-key") {
				lib.FlagUpdate(cmd, "filebase_secret_key", paramsRemoteServerUpdate.FilebaseSecretKey, mapParams)
			}
			if cmd.Flags().Changed("google-cloud-storage-credentials-json") {
				lib.FlagUpdate(cmd, "google_cloud_storage_credentials_json", paramsRemoteServerUpdate.GoogleCloudStorageCredentialsJson, mapParams)
			}
			if cmd.Flags().Changed("google-cloud-storage-s3-compatible-secret-key") {
				lib.FlagUpdate(cmd, "google_cloud_storage_s3_compatible_secret_key", paramsRemoteServerUpdate.GoogleCloudStorageS3CompatibleSecretKey, mapParams)
			}
			if cmd.Flags().Changed("linode-secret-key") {
				lib.FlagUpdate(cmd, "linode_secret_key", paramsRemoteServerUpdate.LinodeSecretKey, mapParams)
			}
			if cmd.Flags().Changed("s3-compatible-secret-key") {
				lib.FlagUpdate(cmd, "s3_compatible_secret_key", paramsRemoteServerUpdate.S3CompatibleSecretKey, mapParams)
			}
			if cmd.Flags().Changed("wasabi-secret-key") {
				lib.FlagUpdate(cmd, "wasabi_secret_key", paramsRemoteServerUpdate.WasabiSecretKey, mapParams)
			}
			if cmd.Flags().Changed("aws-access-key") {
				lib.FlagUpdate(cmd, "aws_access_key", paramsRemoteServerUpdate.AwsAccessKey, mapParams)
			}
			if cmd.Flags().Changed("azure-blob-storage-account") {
				lib.FlagUpdate(cmd, "azure_blob_storage_account", paramsRemoteServerUpdate.AzureBlobStorageAccount, mapParams)
			}
			if cmd.Flags().Changed("azure-blob-storage-container") {
				lib.FlagUpdate(cmd, "azure_blob_storage_container", paramsRemoteServerUpdate.AzureBlobStorageContainer, mapParams)
			}
			if cmd.Flags().Changed("azure-blob-storage-dns-suffix") {
				lib.FlagUpdate(cmd, "azure_blob_storage_dns_suffix", paramsRemoteServerUpdate.AzureBlobStorageDnsSuffix, mapParams)
			}
			if cmd.Flags().Changed("azure-blob-storage-hierarchical-namespace") {
				mapParams["azure_blob_storage_hierarchical_namespace"] = updateAzureBlobStorageHierarchicalNamespace
			}
			if cmd.Flags().Changed("azure-files-storage-account") {
				lib.FlagUpdate(cmd, "azure_files_storage_account", paramsRemoteServerUpdate.AzureFilesStorageAccount, mapParams)
			}
			if cmd.Flags().Changed("azure-files-storage-dns-suffix") {
				lib.FlagUpdate(cmd, "azure_files_storage_dns_suffix", paramsRemoteServerUpdate.AzureFilesStorageDnsSuffix, mapParams)
			}
			if cmd.Flags().Changed("azure-files-storage-share-name") {
				lib.FlagUpdate(cmd, "azure_files_storage_share_name", paramsRemoteServerUpdate.AzureFilesStorageShareName, mapParams)
			}
			if cmd.Flags().Changed("backblaze-b2-bucket") {
				lib.FlagUpdate(cmd, "backblaze_b2_bucket", paramsRemoteServerUpdate.BackblazeB2Bucket, mapParams)
			}
			if cmd.Flags().Changed("backblaze-b2-s3-endpoint") {
				lib.FlagUpdate(cmd, "backblaze_b2_s3_endpoint", paramsRemoteServerUpdate.BackblazeB2S3Endpoint, mapParams)
			}
			if cmd.Flags().Changed("buffer-uploads") {
				lib.FlagUpdate(cmd, "buffer_uploads", paramsRemoteServerUpdate.BufferUploads, mapParams)
			}
			if cmd.Flags().Changed("cloudflare-access-key") {
				lib.FlagUpdate(cmd, "cloudflare_access_key", paramsRemoteServerUpdate.CloudflareAccessKey, mapParams)
			}
			if cmd.Flags().Changed("cloudflare-bucket") {
				lib.FlagUpdate(cmd, "cloudflare_bucket", paramsRemoteServerUpdate.CloudflareBucket, mapParams)
			}
			if cmd.Flags().Changed("cloudflare-endpoint") {
				lib.FlagUpdate(cmd, "cloudflare_endpoint", paramsRemoteServerUpdate.CloudflareEndpoint, mapParams)
			}
			if cmd.Flags().Changed("dropbox-teams") {
				mapParams["dropbox_teams"] = updateDropboxTeams
			}
			if cmd.Flags().Changed("enable-dedicated-ips") {
				mapParams["enable_dedicated_ips"] = updateEnableDedicatedIps
			}
			if cmd.Flags().Changed("filebase-access-key") {
				lib.FlagUpdate(cmd, "filebase_access_key", paramsRemoteServerUpdate.FilebaseAccessKey, mapParams)
			}
			if cmd.Flags().Changed("filebase-bucket") {
				lib.FlagUpdate(cmd, "filebase_bucket", paramsRemoteServerUpdate.FilebaseBucket, mapParams)
			}
			if cmd.Flags().Changed("files-agent-permission-set") {
				lib.FlagUpdate(cmd, "files_agent_permission_set", paramsRemoteServerUpdate.FilesAgentPermissionSet, mapParams)
			}
			if cmd.Flags().Changed("files-agent-root") {
				lib.FlagUpdate(cmd, "files_agent_root", paramsRemoteServerUpdate.FilesAgentRoot, mapParams)
			}
			if cmd.Flags().Changed("files-agent-version") {
				lib.FlagUpdate(cmd, "files_agent_version", paramsRemoteServerUpdate.FilesAgentVersion, mapParams)
			}
			if cmd.Flags().Changed("google-cloud-storage-bucket") {
				lib.FlagUpdate(cmd, "google_cloud_storage_bucket", paramsRemoteServerUpdate.GoogleCloudStorageBucket, mapParams)
			}
			if cmd.Flags().Changed("google-cloud-storage-project-id") {
				lib.FlagUpdate(cmd, "google_cloud_storage_project_id", paramsRemoteServerUpdate.GoogleCloudStorageProjectId, mapParams)
			}
			if cmd.Flags().Changed("google-cloud-storage-s3-compatible-access-key") {
				lib.FlagUpdate(cmd, "google_cloud_storage_s3_compatible_access_key", paramsRemoteServerUpdate.GoogleCloudStorageS3CompatibleAccessKey, mapParams)
			}
			if cmd.Flags().Changed("hostname") {
				lib.FlagUpdate(cmd, "hostname", paramsRemoteServerUpdate.Hostname, mapParams)
			}
			if cmd.Flags().Changed("linode-access-key") {
				lib.FlagUpdate(cmd, "linode_access_key", paramsRemoteServerUpdate.LinodeAccessKey, mapParams)
			}
			if cmd.Flags().Changed("linode-bucket") {
				lib.FlagUpdate(cmd, "linode_bucket", paramsRemoteServerUpdate.LinodeBucket, mapParams)
			}
			if cmd.Flags().Changed("linode-region") {
				lib.FlagUpdate(cmd, "linode_region", paramsRemoteServerUpdate.LinodeRegion, mapParams)
			}
			if cmd.Flags().Changed("max-connections") {
				lib.FlagUpdate(cmd, "max_connections", paramsRemoteServerUpdate.MaxConnections, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsRemoteServerUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("one-drive-account-type") {
				lib.FlagUpdate(cmd, "one_drive_account_type", paramsRemoteServerUpdate.OneDriveAccountType, mapParams)
			}
			if cmd.Flags().Changed("pin-to-site-region") {
				mapParams["pin_to_site_region"] = updatePinToSiteRegion
			}
			if cmd.Flags().Changed("port") {
				lib.FlagUpdate(cmd, "port", paramsRemoteServerUpdate.Port, mapParams)
			}
			if cmd.Flags().Changed("s3-bucket") {
				lib.FlagUpdate(cmd, "s3_bucket", paramsRemoteServerUpdate.S3Bucket, mapParams)
			}
			if cmd.Flags().Changed("s3-compatible-access-key") {
				lib.FlagUpdate(cmd, "s3_compatible_access_key", paramsRemoteServerUpdate.S3CompatibleAccessKey, mapParams)
			}
			if cmd.Flags().Changed("s3-compatible-bucket") {
				lib.FlagUpdate(cmd, "s3_compatible_bucket", paramsRemoteServerUpdate.S3CompatibleBucket, mapParams)
			}
			if cmd.Flags().Changed("s3-compatible-endpoint") {
				lib.FlagUpdate(cmd, "s3_compatible_endpoint", paramsRemoteServerUpdate.S3CompatibleEndpoint, mapParams)
			}
			if cmd.Flags().Changed("s3-compatible-region") {
				lib.FlagUpdate(cmd, "s3_compatible_region", paramsRemoteServerUpdate.S3CompatibleRegion, mapParams)
			}
			if cmd.Flags().Changed("s3-region") {
				lib.FlagUpdate(cmd, "s3_region", paramsRemoteServerUpdate.S3Region, mapParams)
			}
			if cmd.Flags().Changed("server-certificate") {
				lib.FlagUpdate(cmd, "server_certificate", paramsRemoteServerUpdate.ServerCertificate, mapParams)
			}
			if cmd.Flags().Changed("server-host-key") {
				lib.FlagUpdate(cmd, "server_host_key", paramsRemoteServerUpdate.ServerHostKey, mapParams)
			}
			if cmd.Flags().Changed("server-type") {
				lib.FlagUpdate(cmd, "server_type", paramsRemoteServerUpdate.ServerType, mapParams)
			}
			if cmd.Flags().Changed("ssl") {
				lib.FlagUpdate(cmd, "ssl", paramsRemoteServerUpdate.Ssl, mapParams)
			}
			if cmd.Flags().Changed("username") {
				lib.FlagUpdate(cmd, "username", paramsRemoteServerUpdate.Username, mapParams)
			}
			if cmd.Flags().Changed("wasabi-access-key") {
				lib.FlagUpdate(cmd, "wasabi_access_key", paramsRemoteServerUpdate.WasabiAccessKey, mapParams)
			}
			if cmd.Flags().Changed("wasabi-bucket") {
				lib.FlagUpdate(cmd, "wasabi_bucket", paramsRemoteServerUpdate.WasabiBucket, mapParams)
			}
			if cmd.Flags().Changed("wasabi-region") {
				lib.FlagUpdate(cmd, "wasabi_region", paramsRemoteServerUpdate.WasabiRegion, mapParams)
			}

			var remoteServer interface{}
			var err error
			remoteServer, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), remoteServer, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsRemoteServerUpdate.Id, "id", 0, "Remote Server ID.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.Password, "password", "", "Password, if needed.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.PrivateKey, "private-key", "", "Private key, if needed.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.PrivateKeyPassphrase, "private-key-passphrase", "", "Passphrase for private key if needed.")
	cmdUpdate.Flags().BoolVar(&updateResetAuthentication, "reset-authentication", updateResetAuthentication, "Reset authenticated account?")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.SslCertificate, "ssl-certificate", "", "SSL client certificate.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.AwsSecretKey, "aws-secret-key", "", "AWS: secret key.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.AzureBlobStorageAccessKey, "azure-blob-storage-access-key", "", "Azure Blob Storage: Access Key")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.AzureBlobStorageSasToken, "azure-blob-storage-sas-token", "", "Azure Blob Storage: Shared Access Signature (SAS) token")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.AzureFilesStorageAccessKey, "azure-files-storage-access-key", "", "Azure File Storage: Access Key")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.AzureFilesStorageSasToken, "azure-files-storage-sas-token", "", "Azure File Storage: Shared Access Signature (SAS) token")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.BackblazeB2ApplicationKey, "backblaze-b2-application-key", "", "Backblaze B2 Cloud Storage: applicationKey")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.BackblazeB2KeyId, "backblaze-b2-key-id", "", "Backblaze B2 Cloud Storage: keyID")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.CloudflareSecretKey, "cloudflare-secret-key", "", "Cloudflare: Secret Key")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.FilebaseSecretKey, "filebase-secret-key", "", "Filebase: Secret Key")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.GoogleCloudStorageCredentialsJson, "google-cloud-storage-credentials-json", "", "Google Cloud Storage: JSON file that contains the private key. To generate see https://cloud.google.com/storage/docs/json_api/v1/how-tos/authorizing#APIKey")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.GoogleCloudStorageS3CompatibleSecretKey, "google-cloud-storage-s3-compatible-secret-key", "", "Google Cloud Storage: S3-compatible secret key")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.LinodeSecretKey, "linode-secret-key", "", "Linode: Secret Key")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.S3CompatibleSecretKey, "s3-compatible-secret-key", "", "S3-compatible: Secret Key")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.WasabiSecretKey, "wasabi-secret-key", "", "Wasabi: Secret Key")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.AwsAccessKey, "aws-access-key", "", "AWS Access Key.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.AzureBlobStorageAccount, "azure-blob-storage-account", "", "Azure Blob Storage: Account name")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.AzureBlobStorageContainer, "azure-blob-storage-container", "", "Azure Blob Storage: Container name")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.AzureBlobStorageDnsSuffix, "azure-blob-storage-dns-suffix", "", "Azure Blob Storage: Custom DNS suffix")
	cmdUpdate.Flags().BoolVar(&updateAzureBlobStorageHierarchicalNamespace, "azure-blob-storage-hierarchical-namespace", updateAzureBlobStorageHierarchicalNamespace, "Azure Blob Storage: Does the storage account has hierarchical namespace feature enabled?")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.AzureFilesStorageAccount, "azure-files-storage-account", "", "Azure Files: Storage Account name")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.AzureFilesStorageDnsSuffix, "azure-files-storage-dns-suffix", "", "Azure Files: Custom DNS suffix")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.AzureFilesStorageShareName, "azure-files-storage-share-name", "", "Azure Files:  Storage Share name")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.BackblazeB2Bucket, "backblaze-b2-bucket", "", "Backblaze B2 Cloud Storage: Bucket name")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.BackblazeB2S3Endpoint, "backblaze-b2-s3-endpoint", "", "Backblaze B2 Cloud Storage: S3 Endpoint")
	cmdUpdate.Flags().StringVar(&RemoteServerUpdateBufferUploads, "buffer-uploads", "", fmt.Sprintf("If set to always, uploads to this server will be uploaded first to Files.com before being sent to the remote server. This can improve performance in certain access patterns, such as high-latency connections.  It will cause data to be temporarily stored in Files.com. If set to auto, we will perform this optimization if we believe it to be a benefit in a given situation. %v", reflect.ValueOf(paramsRemoteServerUpdate.BufferUploads.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.CloudflareAccessKey, "cloudflare-access-key", "", "Cloudflare: Access Key.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.CloudflareBucket, "cloudflare-bucket", "", "Cloudflare: Bucket name")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.CloudflareEndpoint, "cloudflare-endpoint", "", "Cloudflare: endpoint")
	cmdUpdate.Flags().BoolVar(&updateDropboxTeams, "dropbox-teams", updateDropboxTeams, "Dropbox: If true, list Team folders in root?")
	cmdUpdate.Flags().BoolVar(&updateEnableDedicatedIps, "enable-dedicated-ips", updateEnableDedicatedIps, "`true` if remote server only accepts connections from dedicated IPs")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.FilebaseAccessKey, "filebase-access-key", "", "Filebase: Access Key.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.FilebaseBucket, "filebase-bucket", "", "Filebase: Bucket name")
	cmdUpdate.Flags().StringVar(&RemoteServerUpdateFilesAgentPermissionSet, "files-agent-permission-set", "", fmt.Sprintf("Local permissions for files agent. read_only, write_only, or read_write %v", reflect.ValueOf(paramsRemoteServerUpdate.FilesAgentPermissionSet.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.FilesAgentRoot, "files-agent-root", "", "Agent local root path")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.FilesAgentVersion, "files-agent-version", "", "Files Agent version")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.GoogleCloudStorageBucket, "google-cloud-storage-bucket", "", "Google Cloud Storage: Bucket Name")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.GoogleCloudStorageProjectId, "google-cloud-storage-project-id", "", "Google Cloud Storage: Project ID")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.GoogleCloudStorageS3CompatibleAccessKey, "google-cloud-storage-s3-compatible-access-key", "", "Google Cloud Storage: S3-compatible Access Key.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.Hostname, "hostname", "", "Hostname or IP address")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.LinodeAccessKey, "linode-access-key", "", "Linode: Access Key")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.LinodeBucket, "linode-bucket", "", "Linode: Bucket name")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.LinodeRegion, "linode-region", "", "Linode: region")
	cmdUpdate.Flags().Int64Var(&paramsRemoteServerUpdate.MaxConnections, "max-connections", 0, "Max number of parallel connections.  Ignored for S3 connections (we will parallelize these as much as possible).")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.Name, "name", "", "Internal name for your reference")
	cmdUpdate.Flags().StringVar(&RemoteServerUpdateOneDriveAccountType, "one-drive-account-type", "", fmt.Sprintf("OneDrive: Either personal or business_other account types %v", reflect.ValueOf(paramsRemoteServerUpdate.OneDriveAccountType.Enum()).MapKeys()))
	cmdUpdate.Flags().BoolVar(&updatePinToSiteRegion, "pin-to-site-region", updatePinToSiteRegion, "If true, we will ensure that all communications with this remote server are made through the primary region of the site.  This setting can also be overridden by a site-wide setting which will force it to true.")
	cmdUpdate.Flags().Int64Var(&paramsRemoteServerUpdate.Port, "port", 0, "Port for remote server.  Not needed for S3.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.S3Bucket, "s3-bucket", "", "S3 bucket name")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.S3CompatibleAccessKey, "s3-compatible-access-key", "", "S3-compatible: Access Key")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.S3CompatibleBucket, "s3-compatible-bucket", "", "S3-compatible: Bucket name")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.S3CompatibleEndpoint, "s3-compatible-endpoint", "", "S3-compatible: endpoint")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.S3CompatibleRegion, "s3-compatible-region", "", "S3-compatible: region")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.S3Region, "s3-region", "", "S3 region")
	cmdUpdate.Flags().StringVar(&RemoteServerUpdateServerCertificate, "server-certificate", "", fmt.Sprintf("Remote server certificate %v", reflect.ValueOf(paramsRemoteServerUpdate.ServerCertificate.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.ServerHostKey, "server-host-key", "", "Remote server SSH Host Key. If provided, we will require that the server host key matches the provided key. Uses OpenSSH format similar to what would go into ~/.ssh/known_hosts")
	cmdUpdate.Flags().StringVar(&RemoteServerUpdateServerType, "server-type", "", fmt.Sprintf("Remote server type. %v", reflect.ValueOf(paramsRemoteServerUpdate.ServerType.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&RemoteServerUpdateSsl, "ssl", "", fmt.Sprintf("Should we require SSL? %v", reflect.ValueOf(paramsRemoteServerUpdate.Ssl.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.Username, "username", "", "Remote server username.  Not needed for S3 buckets.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.WasabiAccessKey, "wasabi-access-key", "", "Wasabi: Access Key.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.WasabiBucket, "wasabi-bucket", "", "Wasabi: Bucket name")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerUpdate.WasabiRegion, "wasabi-region", "", "Wasabi: Region")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	RemoteServers.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsRemoteServerDelete := files_sdk.RemoteServerDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Remote Server`,
		Long:  `Delete Remote Server`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := remote_server.Client{Config: config}

			var err error
			err = client.Delete(paramsRemoteServerDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsRemoteServerDelete.Id, "id", 0, "Remote Server ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	RemoteServers.AddCommand(cmdDelete)
	return RemoteServers
}
