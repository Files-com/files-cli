package cmd

import (
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	flib "github.com/Files-com/files-sdk-go/v2/lib"

	"fmt"

	remote_server "github.com/Files-com/files-sdk-go/v2/remoteserver"
)

var (
	RemoteServers = &cobra.Command{}
)

func RemoteServersInit() {
	RemoteServers = &cobra.Command{
		Use:  "remote-servers [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command remote-servers\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsRemoteServerList := files_sdk.RemoteServerListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsRemoteServerList
			params.MaxPages = MaxPagesList

			client := remote_server.Client{Config: *config}
			it, err := client.List(ctx, params)
			if err != nil {
				lib.ClientError(ctx, err)
			}
			err = lib.FormatIter(it, formatList, fieldsList)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdList.Flags().StringVarP(&paramsRemoteServerList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsRemoteServerList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	RemoteServers.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	paramsRemoteServerFind := files_sdk.RemoteServerFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := remote_server.Client{Config: *config}

			result, err := client.Find(ctx, paramsRemoteServerFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatFind, fieldsFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsRemoteServerFind.Id, "id", "i", 0, "Remote Server ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright")
	RemoteServers.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	createResetAuthentication := false
	createEnableDedicatedIps := false
	paramsRemoteServerCreate := files_sdk.RemoteServerCreateParams{}
	RemoteServerCreateServerCertificate := ""
	RemoteServerCreateServerType := ""
	RemoteServerCreateSsl := ""
	RemoteServerCreateOneDriveAccountType := ""

	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := remote_server.Client{Config: *config}

			if createResetAuthentication {
				paramsRemoteServerCreate.ResetAuthentication = flib.Bool(true)
			}
			if createEnableDedicatedIps {
				paramsRemoteServerCreate.EnableDedicatedIps = flib.Bool(true)
			}

			paramsRemoteServerCreate.ServerCertificate = paramsRemoteServerCreate.ServerCertificate.Enum()[RemoteServerCreateServerCertificate]
			paramsRemoteServerCreate.ServerType = paramsRemoteServerCreate.ServerType.Enum()[RemoteServerCreateServerType]
			paramsRemoteServerCreate.Ssl = paramsRemoteServerCreate.Ssl.Enum()[RemoteServerCreateSsl]
			paramsRemoteServerCreate.OneDriveAccountType = paramsRemoteServerCreate.OneDriveAccountType.Enum()[RemoteServerCreateOneDriveAccountType]

			result, err := client.Create(ctx, paramsRemoteServerCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatCreate, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.AwsAccessKey, "aws-access-key", "k", "", "AWS Access Key.")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.AwsSecretKey, "aws-secret-key", "e", "", "AWS secret key.")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.Password, "password", "p", "", "Password if needed.")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.PrivateKey, "private-key", "v", "", "Private key if needed.")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.SslCertificate, "ssl-certificate", "", "", "SSL client certificate.")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.GoogleCloudStorageCredentialsJson, "google-cloud-storage-credentials-json", "j", "", "A JSON file that contains the private key. To generate see https://cloud.google.com/storage/docs/json_api/v1/how-tos/authorizing#APIKey")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.WasabiAccessKey, "wasabi-access-key", "", "", "Wasabi access key.")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.WasabiSecretKey, "wasabi-secret-key", "", "", "Wasabi secret key.")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.BackblazeB2KeyId, "backblaze-b2-key-id", "i", "", "Backblaze B2 Cloud Storage keyID.")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.BackblazeB2ApplicationKey, "backblaze-b2-application-key", "", "", "Backblaze B2 Cloud Storage applicationKey.")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.RackspaceApiKey, "rackspace-api-key", "", "", "Rackspace API key from the Rackspace Cloud Control Panel.")
	cmdCreate.Flags().BoolVarP(&createResetAuthentication, "reset-authentication", "", createResetAuthentication, "Reset authenticated account")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.AzureBlobStorageAccessKey, "azure-blob-storage-access-key", "y", "", "Azure Blob Storage secret key.")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.Hostname, "hostname", "o", "", "Hostname or IP address")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.Name, "name", "", "", "Internal name for your reference")
	cmdCreate.Flags().Int64VarP(&paramsRemoteServerCreate.MaxConnections, "max-connections", "x", 0, "Max number of parallel connections.  Ignored for S3 connections (we will parallelize these as much as possible).")
	cmdCreate.Flags().Int64VarP(&paramsRemoteServerCreate.Port, "port", "t", 0, "Port for remote server.  Not needed for S3.")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.S3Bucket, "s3-bucket", "", "", "S3 bucket name")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.S3Region, "s3-region", "", "", "S3 region")
	cmdCreate.Flags().StringVarP(&RemoteServerCreateServerCertificate, "server-certificate", "f", "", fmt.Sprintf("Remote server certificate %v", reflect.ValueOf(paramsRemoteServerCreate.ServerCertificate.Enum()).MapKeys()))
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.ServerHostKey, "server-host-key", "", "", "Remote server SSH Host Key. If provided, we will require that the server host key matches the provided key. Uses OpenSSH format similar to what would go into ~/.ssh/known_hosts")
	cmdCreate.Flags().StringVarP(&RemoteServerCreateServerType, "server-type", "", "", fmt.Sprintf("Remote server type. %v", reflect.ValueOf(paramsRemoteServerCreate.ServerType.Enum()).MapKeys()))
	cmdCreate.Flags().StringVarP(&RemoteServerCreateSsl, "ssl", "", "", fmt.Sprintf("Should we require SSL? %v", reflect.ValueOf(paramsRemoteServerCreate.Ssl.Enum()).MapKeys()))
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.Username, "username", "", "", "Remote server username.  Not needed for S3 buckets.")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.GoogleCloudStorageBucket, "google-cloud-storage-bucket", "u", "", "Google Cloud Storage bucket name")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.GoogleCloudStorageProjectId, "google-cloud-storage-project-id", "d", "", "Google Cloud Project ID")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.BackblazeB2Bucket, "backblaze-b2-bucket", "b", "", "Backblaze B2 Cloud Storage Bucket name")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.BackblazeB2S3Endpoint, "backblaze-b2-s3-endpoint", "n", "", "Backblaze B2 Cloud Storage S3 Endpoint")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.WasabiBucket, "wasabi-bucket", "", "", "Wasabi Bucket name")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.WasabiRegion, "wasabi-region", "", "", "Wasabi region")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.RackspaceUsername, "rackspace-username", "s", "", "Rackspace username used to login to the Rackspace Cloud Control Panel.")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.RackspaceRegion, "rackspace-region", "g", "", "Three letter airport code for Rackspace region. See https://support.rackspace.com/how-to/about-regions/")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.RackspaceContainer, "rackspace-container", "", "", "The name of the container (top level directory) where files will sync.")
	cmdCreate.Flags().StringVarP(&RemoteServerCreateOneDriveAccountType, "one-drive-account-type", "r", "", fmt.Sprintf("Either personal or business_other account types %v", reflect.ValueOf(paramsRemoteServerCreate.OneDriveAccountType.Enum()).MapKeys()))
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.AzureBlobStorageAccount, "azure-blob-storage-account", "a", "", "Azure Blob Storage Account name")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.AzureBlobStorageContainer, "azure-blob-storage-container", "c", "", "Azure Blob Storage Container name")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.S3CompatibleBucket, "s3-compatible-bucket", "", "", "S3-compatible Bucket name")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.S3CompatibleRegion, "s3-compatible-region", "", "", "S3-compatible Bucket name")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.S3CompatibleEndpoint, "s3-compatible-endpoint", "", "", "S3-compatible endpoint")
	cmdCreate.Flags().BoolVarP(&createEnableDedicatedIps, "enable-dedicated-ips", "l", createEnableDedicatedIps, "`true` if remote server only accepts connections from dedicated IPs")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.S3CompatibleAccessKey, "s3-compatible-access-key", "", "", "S3-compatible access key")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.S3CompatibleSecretKey, "s3-compatible-secret-key", "", "", "S3-compatible secret key")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	RemoteServers.AddCommand(cmdCreate)
	var fieldsUpdate string
	var formatUpdate string
	updateResetAuthentication := false
	updateEnableDedicatedIps := false
	paramsRemoteServerUpdate := files_sdk.RemoteServerUpdateParams{}
	RemoteServerUpdateServerCertificate := ""
	RemoteServerUpdateServerType := ""
	RemoteServerUpdateSsl := ""
	RemoteServerUpdateOneDriveAccountType := ""

	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := remote_server.Client{Config: *config}

			if updateResetAuthentication {
				paramsRemoteServerUpdate.ResetAuthentication = flib.Bool(true)
			}
			if updateEnableDedicatedIps {
				paramsRemoteServerUpdate.EnableDedicatedIps = flib.Bool(true)
			}

			paramsRemoteServerUpdate.ServerCertificate = paramsRemoteServerUpdate.ServerCertificate.Enum()[RemoteServerUpdateServerCertificate]
			paramsRemoteServerUpdate.ServerType = paramsRemoteServerUpdate.ServerType.Enum()[RemoteServerUpdateServerType]
			paramsRemoteServerUpdate.Ssl = paramsRemoteServerUpdate.Ssl.Enum()[RemoteServerUpdateSsl]
			paramsRemoteServerUpdate.OneDriveAccountType = paramsRemoteServerUpdate.OneDriveAccountType.Enum()[RemoteServerUpdateOneDriveAccountType]

			result, err := client.Update(ctx, paramsRemoteServerUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatUpdate, fieldsUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsRemoteServerUpdate.Id, "id", "", 0, "Remote Server ID.")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.AwsAccessKey, "aws-access-key", "k", "", "AWS Access Key.")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.AwsSecretKey, "aws-secret-key", "e", "", "AWS secret key.")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.Password, "password", "p", "", "Password if needed.")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.PrivateKey, "private-key", "v", "", "Private key if needed.")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.SslCertificate, "ssl-certificate", "", "", "SSL client certificate.")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.GoogleCloudStorageCredentialsJson, "google-cloud-storage-credentials-json", "j", "", "A JSON file that contains the private key. To generate see https://cloud.google.com/storage/docs/json_api/v1/how-tos/authorizing#APIKey")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.WasabiAccessKey, "wasabi-access-key", "", "", "Wasabi access key.")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.WasabiSecretKey, "wasabi-secret-key", "", "", "Wasabi secret key.")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.BackblazeB2KeyId, "backblaze-b2-key-id", "i", "", "Backblaze B2 Cloud Storage keyID.")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.BackblazeB2ApplicationKey, "backblaze-b2-application-key", "", "", "Backblaze B2 Cloud Storage applicationKey.")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.RackspaceApiKey, "rackspace-api-key", "", "", "Rackspace API key from the Rackspace Cloud Control Panel.")
	cmdUpdate.Flags().BoolVarP(&updateResetAuthentication, "reset-authentication", "", updateResetAuthentication, "Reset authenticated account")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.AzureBlobStorageAccessKey, "azure-blob-storage-access-key", "y", "", "Azure Blob Storage secret key.")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.Hostname, "hostname", "o", "", "Hostname or IP address")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.Name, "name", "", "", "Internal name for your reference")
	cmdUpdate.Flags().Int64VarP(&paramsRemoteServerUpdate.MaxConnections, "max-connections", "x", 0, "Max number of parallel connections.  Ignored for S3 connections (we will parallelize these as much as possible).")
	cmdUpdate.Flags().Int64VarP(&paramsRemoteServerUpdate.Port, "port", "t", 0, "Port for remote server.  Not needed for S3.")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.S3Bucket, "s3-bucket", "", "", "S3 bucket name")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.S3Region, "s3-region", "", "", "S3 region")
	cmdUpdate.Flags().StringVarP(&RemoteServerUpdateServerCertificate, "server-certificate", "f", "", fmt.Sprintf("Remote server certificate %v", reflect.ValueOf(paramsRemoteServerUpdate.ServerCertificate.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.ServerHostKey, "server-host-key", "", "", "Remote server SSH Host Key. If provided, we will require that the server host key matches the provided key. Uses OpenSSH format similar to what would go into ~/.ssh/known_hosts")
	cmdUpdate.Flags().StringVarP(&RemoteServerUpdateServerType, "server-type", "", "", fmt.Sprintf("Remote server type. %v", reflect.ValueOf(paramsRemoteServerUpdate.ServerType.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVarP(&RemoteServerUpdateSsl, "ssl", "", "", fmt.Sprintf("Should we require SSL? %v", reflect.ValueOf(paramsRemoteServerUpdate.Ssl.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.Username, "username", "", "", "Remote server username.  Not needed for S3 buckets.")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.GoogleCloudStorageBucket, "google-cloud-storage-bucket", "u", "", "Google Cloud Storage bucket name")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.GoogleCloudStorageProjectId, "google-cloud-storage-project-id", "d", "", "Google Cloud Project ID")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.BackblazeB2Bucket, "backblaze-b2-bucket", "b", "", "Backblaze B2 Cloud Storage Bucket name")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.BackblazeB2S3Endpoint, "backblaze-b2-s3-endpoint", "n", "", "Backblaze B2 Cloud Storage S3 Endpoint")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.WasabiBucket, "wasabi-bucket", "", "", "Wasabi Bucket name")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.WasabiRegion, "wasabi-region", "", "", "Wasabi region")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.RackspaceUsername, "rackspace-username", "s", "", "Rackspace username used to login to the Rackspace Cloud Control Panel.")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.RackspaceRegion, "rackspace-region", "g", "", "Three letter airport code for Rackspace region. See https://support.rackspace.com/how-to/about-regions/")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.RackspaceContainer, "rackspace-container", "", "", "The name of the container (top level directory) where files will sync.")
	cmdUpdate.Flags().StringVarP(&RemoteServerUpdateOneDriveAccountType, "one-drive-account-type", "r", "", fmt.Sprintf("Either personal or business_other account types %v", reflect.ValueOf(paramsRemoteServerUpdate.OneDriveAccountType.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.AzureBlobStorageAccount, "azure-blob-storage-account", "a", "", "Azure Blob Storage Account name")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.AzureBlobStorageContainer, "azure-blob-storage-container", "c", "", "Azure Blob Storage Container name")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.S3CompatibleBucket, "s3-compatible-bucket", "", "", "S3-compatible Bucket name")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.S3CompatibleRegion, "s3-compatible-region", "", "", "S3-compatible Bucket name")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.S3CompatibleEndpoint, "s3-compatible-endpoint", "", "", "S3-compatible endpoint")
	cmdUpdate.Flags().BoolVarP(&updateEnableDedicatedIps, "enable-dedicated-ips", "l", updateEnableDedicatedIps, "`true` if remote server only accepts connections from dedicated IPs")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.S3CompatibleAccessKey, "s3-compatible-access-key", "", "", "S3-compatible access key")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.S3CompatibleSecretKey, "s3-compatible-secret-key", "", "", "S3-compatible secret key")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	RemoteServers.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	paramsRemoteServerDelete := files_sdk.RemoteServerDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := remote_server.Client{Config: *config}

			result, err := client.Delete(ctx, paramsRemoteServerDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatDelete, fieldsDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsRemoteServerDelete.Id, "id", "i", 0, "Remote Server ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright")
	RemoteServers.AddCommand(cmdDelete)
}
