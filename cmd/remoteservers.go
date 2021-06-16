package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	flib "github.com/Files-com/files-sdk-go/lib"

	remote_server "github.com/Files-com/files-sdk-go/remoteserver"
)

var (
	RemoteServers = &cobra.Command{}
)

func RemoteServersInit() {
	RemoteServers = &cobra.Command{
		Use:  "remote-servers [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
	var fieldsList string
	paramsRemoteServerList := files_sdk.RemoteServerListParams{}
	var MaxPagesList int64
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			params := paramsRemoteServerList
			params.MaxPages = MaxPagesList
			client := remote_server.Client{Config: *ctx.GetConfig()}
			it, err := client.List(params)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
			err = lib.JsonMarshalIter(it, fieldsList)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdList.Flags().StringVarP(&paramsRemoteServerList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsRemoteServerList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	RemoteServers.AddCommand(cmdList)
	var fieldsFind string
	paramsRemoteServerFind := files_sdk.RemoteServerFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := remote_server.Client{Config: *ctx.GetConfig()}

			result, err := client.Find(paramsRemoteServerFind)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsFind)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsRemoteServerFind.Id, "id", "i", 0, "Remote Server ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	RemoteServers.AddCommand(cmdFind)
	var fieldsCreate string
	createResetAuthentication := false
	paramsRemoteServerCreate := files_sdk.RemoteServerCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := remote_server.Client{Config: *ctx.GetConfig()}

			if createResetAuthentication {
				paramsRemoteServerCreate.ResetAuthentication = flib.Bool(true)
			}

			result, err := client.Create(paramsRemoteServerCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
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
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.ServerCertificate, "server-certificate", "f", "", "Remote server certificate")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.ServerHostKey, "server-host-key", "", "", "Remote server SSH Host Key. If provided, we will require that the server host key matches the provided key. Uses OpenSSH format similar to what would go into ~/.ssh/known_hosts")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.ServerType, "server-type", "", "", "Remote server type.")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.Ssl, "ssl", "l", "", "Should we require SSL?")
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
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.OneDriveAccountType, "one-drive-account-type", "r", "", "Either personal or business_other account types")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.AzureBlobStorageAccount, "azure-blob-storage-account", "a", "", "Azure Blob Storage Account name")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.AzureBlobStorageContainer, "azure-blob-storage-container", "c", "", "Azure Blob Storage Container name")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.S3CompatibleBucket, "s3-compatible-bucket", "", "", "S3-compatible Bucket name")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.S3CompatibleRegion, "s3-compatible-region", "", "", "S3-compatible Bucket name")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.S3CompatibleEndpoint, "s3-compatible-endpoint", "", "", "S3-compatible endpoint")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.S3CompatibleAccessKey, "s3-compatible-access-key", "", "", "S3-compatible access key")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.S3CompatibleSecretKey, "s3-compatible-secret-key", "", "", "S3-compatible secret key")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	RemoteServers.AddCommand(cmdCreate)
	var fieldsUpdate string
	updateResetAuthentication := false
	paramsRemoteServerUpdate := files_sdk.RemoteServerUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := remote_server.Client{Config: *ctx.GetConfig()}

			if updateResetAuthentication {
				paramsRemoteServerUpdate.ResetAuthentication = flib.Bool(true)
			}

			result, err := client.Update(paramsRemoteServerUpdate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsUpdate)
			if err != nil {
				lib.ClientError(err, &ctx)
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
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.ServerCertificate, "server-certificate", "f", "", "Remote server certificate")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.ServerHostKey, "server-host-key", "", "", "Remote server SSH Host Key. If provided, we will require that the server host key matches the provided key. Uses OpenSSH format similar to what would go into ~/.ssh/known_hosts")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.ServerType, "server-type", "", "", "Remote server type.")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.Ssl, "ssl", "l", "", "Should we require SSL?")
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
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.OneDriveAccountType, "one-drive-account-type", "r", "", "Either personal or business_other account types")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.AzureBlobStorageAccount, "azure-blob-storage-account", "a", "", "Azure Blob Storage Account name")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.AzureBlobStorageContainer, "azure-blob-storage-container", "c", "", "Azure Blob Storage Container name")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.S3CompatibleBucket, "s3-compatible-bucket", "", "", "S3-compatible Bucket name")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.S3CompatibleRegion, "s3-compatible-region", "", "", "S3-compatible Bucket name")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.S3CompatibleEndpoint, "s3-compatible-endpoint", "", "", "S3-compatible endpoint")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.S3CompatibleAccessKey, "s3-compatible-access-key", "", "", "S3-compatible access key")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.S3CompatibleSecretKey, "s3-compatible-secret-key", "", "", "S3-compatible secret key")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	RemoteServers.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsRemoteServerDelete := files_sdk.RemoteServerDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := remote_server.Client{Config: *ctx.GetConfig()}

			result, err := client.Delete(paramsRemoteServerDelete)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsRemoteServerDelete.Id, "id", "i", 0, "Remote Server ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	RemoteServers.AddCommand(cmdDelete)
}
