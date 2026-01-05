package cmd

import (
	"fmt"
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	remote_server_credential "github.com/Files-com/files-sdk-go/v3/remoteservercredential"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(RemoteServerCredentials())
}

func RemoteServerCredentials() *cobra.Command {
	RemoteServerCredentials := &cobra.Command{
		Use:  "remote-server-credentials [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command remote-server-credentials\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsRemoteServerCredentialList := files_sdk.RemoteServerCredentialListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Remote Server Credentials",
		Long:    `List Remote Server Credentials`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsRemoteServerCredentialList
			params.MaxPages = MaxPagesList

			client := remote_server_credential.Client{Config: config}
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

	cmdList.Flags().StringVar(&paramsRemoteServerCredentialList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsRemoteServerCredentialList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	RemoteServerCredentials.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsRemoteServerCredentialFind := files_sdk.RemoteServerCredentialFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Remote Server Credential`,
		Long:  `Show Remote Server Credential`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := remote_server_credential.Client{Config: config}

			var remoteServerCredential interface{}
			var err error
			remoteServerCredential, err = client.Find(paramsRemoteServerCredentialFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), remoteServerCredential, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsRemoteServerCredentialFind.Id, "id", 0, "Remote Server Credential ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	RemoteServerCredentials.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsRemoteServerCredentialCreate := files_sdk.RemoteServerCredentialCreateParams{}
	RemoteServerCredentialCreateServerType := ""

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Remote Server Credential`,
		Long:  `Create Remote Server Credential`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := remote_server_credential.Client{Config: config}

			var RemoteServerCredentialCreateServerTypeErr error
			paramsRemoteServerCredentialCreate.ServerType, RemoteServerCredentialCreateServerTypeErr = lib.FetchKey("server-type", paramsRemoteServerCredentialCreate.ServerType.Enum(), RemoteServerCredentialCreateServerType)
			if RemoteServerCredentialCreateServerType != "" && RemoteServerCredentialCreateServerTypeErr != nil {
				return RemoteServerCredentialCreateServerTypeErr
			}

			var remoteServerCredential interface{}
			var err error
			remoteServerCredential, err = client.Create(paramsRemoteServerCredentialCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), remoteServerCredential, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.Name, "name", "", "Internal name for your reference")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.Description, "description", "", "Internal description for your reference")
	cmdCreate.Flags().StringVar(&RemoteServerCredentialCreateServerType, "server-type", "", fmt.Sprintf("Remote server type.  Remote Server Credentials are only valid for a single type of Remote Server. %v", reflect.ValueOf(paramsRemoteServerCredentialCreate.ServerType.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.AwsAccessKey, "aws-access-key", "", "AWS Access Key.")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.AzureBlobStorageAccount, "azure-blob-storage-account", "", "Azure Blob Storage: Account name")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.AzureFilesStorageAccount, "azure-files-storage-account", "", "Azure Files: Storage Account name")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.CloudflareAccessKey, "cloudflare-access-key", "", "Cloudflare: Access Key.")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.FilebaseAccessKey, "filebase-access-key", "", "Filebase: Access Key.")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.GoogleCloudStorageS3CompatibleAccessKey, "google-cloud-storage-s3-compatible-access-key", "", "Google Cloud Storage: S3-compatible Access Key.")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.LinodeAccessKey, "linode-access-key", "", "Linode: Access Key")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.S3CompatibleAccessKey, "s3-compatible-access-key", "", "S3-compatible: Access Key")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.Username, "username", "", "Remote server username.")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.WasabiAccessKey, "wasabi-access-key", "", "Wasabi: Access Key.")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.Password, "password", "", "Password, if needed.")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.PrivateKey, "private-key", "", "Private key, if needed.")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.PrivateKeyPassphrase, "private-key-passphrase", "", "Passphrase for private key if needed.")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.AwsSecretKey, "aws-secret-key", "", "AWS: secret key.")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.AzureBlobStorageAccessKey, "azure-blob-storage-access-key", "", "Azure Blob Storage: Access Key")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.AzureBlobStorageSasToken, "azure-blob-storage-sas-token", "", "Azure Blob Storage: Shared Access Signature (SAS) token")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.AzureFilesStorageAccessKey, "azure-files-storage-access-key", "", "Azure File Storage: Access Key")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.AzureFilesStorageSasToken, "azure-files-storage-sas-token", "", "Azure File Storage: Shared Access Signature (SAS) token")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.BackblazeB2ApplicationKey, "backblaze-b2-application-key", "", "Backblaze B2 Cloud Storage: applicationKey")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.BackblazeB2KeyId, "backblaze-b2-key-id", "", "Backblaze B2 Cloud Storage: keyID")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.CloudflareSecretKey, "cloudflare-secret-key", "", "Cloudflare: Secret Key")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.FilebaseSecretKey, "filebase-secret-key", "", "Filebase: Secret Key")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.GoogleCloudStorageCredentialsJson, "google-cloud-storage-credentials-json", "", "Google Cloud Storage: JSON file that contains the private key. To generate see https://cloud.google.com/storage/docs/json_api/v1/how-tos/authorizing#APIKey")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.GoogleCloudStorageS3CompatibleSecretKey, "google-cloud-storage-s3-compatible-secret-key", "", "Google Cloud Storage: S3-compatible secret key")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.LinodeSecretKey, "linode-secret-key", "", "Linode: Secret Key")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.S3CompatibleSecretKey, "s3-compatible-secret-key", "", "S3-compatible: Secret Key")
	cmdCreate.Flags().StringVar(&paramsRemoteServerCredentialCreate.WasabiSecretKey, "wasabi-secret-key", "", "Wasabi: Secret Key")
	cmdCreate.Flags().Int64Var(&paramsRemoteServerCredentialCreate.WorkspaceId, "workspace-id", 0, "Workspace ID (0 for default workspace)")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	RemoteServerCredentials.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsRemoteServerCredentialUpdate := files_sdk.RemoteServerCredentialUpdateParams{}
	RemoteServerCredentialUpdateServerType := ""

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Remote Server Credential`,
		Long:  `Update Remote Server Credential`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := remote_server_credential.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.RemoteServerCredentialUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			var RemoteServerCredentialUpdateServerTypeErr error
			paramsRemoteServerCredentialUpdate.ServerType, RemoteServerCredentialUpdateServerTypeErr = lib.FetchKey("server-type", paramsRemoteServerCredentialUpdate.ServerType.Enum(), RemoteServerCredentialUpdateServerType)
			if RemoteServerCredentialUpdateServerType != "" && RemoteServerCredentialUpdateServerTypeErr != nil {
				return RemoteServerCredentialUpdateServerTypeErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsRemoteServerCredentialUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsRemoteServerCredentialUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("description") {
				lib.FlagUpdate(cmd, "description", paramsRemoteServerCredentialUpdate.Description, mapParams)
			}
			if cmd.Flags().Changed("server-type") {
				lib.FlagUpdate(cmd, "server_type", paramsRemoteServerCredentialUpdate.ServerType, mapParams)
			}
			if cmd.Flags().Changed("aws-access-key") {
				lib.FlagUpdate(cmd, "aws_access_key", paramsRemoteServerCredentialUpdate.AwsAccessKey, mapParams)
			}
			if cmd.Flags().Changed("azure-blob-storage-account") {
				lib.FlagUpdate(cmd, "azure_blob_storage_account", paramsRemoteServerCredentialUpdate.AzureBlobStorageAccount, mapParams)
			}
			if cmd.Flags().Changed("azure-files-storage-account") {
				lib.FlagUpdate(cmd, "azure_files_storage_account", paramsRemoteServerCredentialUpdate.AzureFilesStorageAccount, mapParams)
			}
			if cmd.Flags().Changed("cloudflare-access-key") {
				lib.FlagUpdate(cmd, "cloudflare_access_key", paramsRemoteServerCredentialUpdate.CloudflareAccessKey, mapParams)
			}
			if cmd.Flags().Changed("filebase-access-key") {
				lib.FlagUpdate(cmd, "filebase_access_key", paramsRemoteServerCredentialUpdate.FilebaseAccessKey, mapParams)
			}
			if cmd.Flags().Changed("google-cloud-storage-s3-compatible-access-key") {
				lib.FlagUpdate(cmd, "google_cloud_storage_s3_compatible_access_key", paramsRemoteServerCredentialUpdate.GoogleCloudStorageS3CompatibleAccessKey, mapParams)
			}
			if cmd.Flags().Changed("linode-access-key") {
				lib.FlagUpdate(cmd, "linode_access_key", paramsRemoteServerCredentialUpdate.LinodeAccessKey, mapParams)
			}
			if cmd.Flags().Changed("s3-compatible-access-key") {
				lib.FlagUpdate(cmd, "s3_compatible_access_key", paramsRemoteServerCredentialUpdate.S3CompatibleAccessKey, mapParams)
			}
			if cmd.Flags().Changed("username") {
				lib.FlagUpdate(cmd, "username", paramsRemoteServerCredentialUpdate.Username, mapParams)
			}
			if cmd.Flags().Changed("wasabi-access-key") {
				lib.FlagUpdate(cmd, "wasabi_access_key", paramsRemoteServerCredentialUpdate.WasabiAccessKey, mapParams)
			}
			if cmd.Flags().Changed("password") {
				lib.FlagUpdate(cmd, "password", paramsRemoteServerCredentialUpdate.Password, mapParams)
			}
			if cmd.Flags().Changed("private-key") {
				lib.FlagUpdate(cmd, "private_key", paramsRemoteServerCredentialUpdate.PrivateKey, mapParams)
			}
			if cmd.Flags().Changed("private-key-passphrase") {
				lib.FlagUpdate(cmd, "private_key_passphrase", paramsRemoteServerCredentialUpdate.PrivateKeyPassphrase, mapParams)
			}
			if cmd.Flags().Changed("aws-secret-key") {
				lib.FlagUpdate(cmd, "aws_secret_key", paramsRemoteServerCredentialUpdate.AwsSecretKey, mapParams)
			}
			if cmd.Flags().Changed("azure-blob-storage-access-key") {
				lib.FlagUpdate(cmd, "azure_blob_storage_access_key", paramsRemoteServerCredentialUpdate.AzureBlobStorageAccessKey, mapParams)
			}
			if cmd.Flags().Changed("azure-blob-storage-sas-token") {
				lib.FlagUpdate(cmd, "azure_blob_storage_sas_token", paramsRemoteServerCredentialUpdate.AzureBlobStorageSasToken, mapParams)
			}
			if cmd.Flags().Changed("azure-files-storage-access-key") {
				lib.FlagUpdate(cmd, "azure_files_storage_access_key", paramsRemoteServerCredentialUpdate.AzureFilesStorageAccessKey, mapParams)
			}
			if cmd.Flags().Changed("azure-files-storage-sas-token") {
				lib.FlagUpdate(cmd, "azure_files_storage_sas_token", paramsRemoteServerCredentialUpdate.AzureFilesStorageSasToken, mapParams)
			}
			if cmd.Flags().Changed("backblaze-b2-application-key") {
				lib.FlagUpdate(cmd, "backblaze_b2_application_key", paramsRemoteServerCredentialUpdate.BackblazeB2ApplicationKey, mapParams)
			}
			if cmd.Flags().Changed("backblaze-b2-key-id") {
				lib.FlagUpdate(cmd, "backblaze_b2_key_id", paramsRemoteServerCredentialUpdate.BackblazeB2KeyId, mapParams)
			}
			if cmd.Flags().Changed("cloudflare-secret-key") {
				lib.FlagUpdate(cmd, "cloudflare_secret_key", paramsRemoteServerCredentialUpdate.CloudflareSecretKey, mapParams)
			}
			if cmd.Flags().Changed("filebase-secret-key") {
				lib.FlagUpdate(cmd, "filebase_secret_key", paramsRemoteServerCredentialUpdate.FilebaseSecretKey, mapParams)
			}
			if cmd.Flags().Changed("google-cloud-storage-credentials-json") {
				lib.FlagUpdate(cmd, "google_cloud_storage_credentials_json", paramsRemoteServerCredentialUpdate.GoogleCloudStorageCredentialsJson, mapParams)
			}
			if cmd.Flags().Changed("google-cloud-storage-s3-compatible-secret-key") {
				lib.FlagUpdate(cmd, "google_cloud_storage_s3_compatible_secret_key", paramsRemoteServerCredentialUpdate.GoogleCloudStorageS3CompatibleSecretKey, mapParams)
			}
			if cmd.Flags().Changed("linode-secret-key") {
				lib.FlagUpdate(cmd, "linode_secret_key", paramsRemoteServerCredentialUpdate.LinodeSecretKey, mapParams)
			}
			if cmd.Flags().Changed("s3-compatible-secret-key") {
				lib.FlagUpdate(cmd, "s3_compatible_secret_key", paramsRemoteServerCredentialUpdate.S3CompatibleSecretKey, mapParams)
			}
			if cmd.Flags().Changed("wasabi-secret-key") {
				lib.FlagUpdate(cmd, "wasabi_secret_key", paramsRemoteServerCredentialUpdate.WasabiSecretKey, mapParams)
			}

			var remoteServerCredential interface{}
			var err error
			remoteServerCredential, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), remoteServerCredential, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsRemoteServerCredentialUpdate.Id, "id", 0, "Remote Server Credential ID.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.Name, "name", "", "Internal name for your reference")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.Description, "description", "", "Internal description for your reference")
	cmdUpdate.Flags().StringVar(&RemoteServerCredentialUpdateServerType, "server-type", "", fmt.Sprintf("Remote server type.  Remote Server Credentials are only valid for a single type of Remote Server. %v", reflect.ValueOf(paramsRemoteServerCredentialUpdate.ServerType.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.AwsAccessKey, "aws-access-key", "", "AWS Access Key.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.AzureBlobStorageAccount, "azure-blob-storage-account", "", "Azure Blob Storage: Account name")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.AzureFilesStorageAccount, "azure-files-storage-account", "", "Azure Files: Storage Account name")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.CloudflareAccessKey, "cloudflare-access-key", "", "Cloudflare: Access Key.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.FilebaseAccessKey, "filebase-access-key", "", "Filebase: Access Key.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.GoogleCloudStorageS3CompatibleAccessKey, "google-cloud-storage-s3-compatible-access-key", "", "Google Cloud Storage: S3-compatible Access Key.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.LinodeAccessKey, "linode-access-key", "", "Linode: Access Key")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.S3CompatibleAccessKey, "s3-compatible-access-key", "", "S3-compatible: Access Key")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.Username, "username", "", "Remote server username.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.WasabiAccessKey, "wasabi-access-key", "", "Wasabi: Access Key.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.Password, "password", "", "Password, if needed.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.PrivateKey, "private-key", "", "Private key, if needed.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.PrivateKeyPassphrase, "private-key-passphrase", "", "Passphrase for private key if needed.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.AwsSecretKey, "aws-secret-key", "", "AWS: secret key.")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.AzureBlobStorageAccessKey, "azure-blob-storage-access-key", "", "Azure Blob Storage: Access Key")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.AzureBlobStorageSasToken, "azure-blob-storage-sas-token", "", "Azure Blob Storage: Shared Access Signature (SAS) token")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.AzureFilesStorageAccessKey, "azure-files-storage-access-key", "", "Azure File Storage: Access Key")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.AzureFilesStorageSasToken, "azure-files-storage-sas-token", "", "Azure File Storage: Shared Access Signature (SAS) token")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.BackblazeB2ApplicationKey, "backblaze-b2-application-key", "", "Backblaze B2 Cloud Storage: applicationKey")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.BackblazeB2KeyId, "backblaze-b2-key-id", "", "Backblaze B2 Cloud Storage: keyID")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.CloudflareSecretKey, "cloudflare-secret-key", "", "Cloudflare: Secret Key")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.FilebaseSecretKey, "filebase-secret-key", "", "Filebase: Secret Key")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.GoogleCloudStorageCredentialsJson, "google-cloud-storage-credentials-json", "", "Google Cloud Storage: JSON file that contains the private key. To generate see https://cloud.google.com/storage/docs/json_api/v1/how-tos/authorizing#APIKey")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.GoogleCloudStorageS3CompatibleSecretKey, "google-cloud-storage-s3-compatible-secret-key", "", "Google Cloud Storage: S3-compatible secret key")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.LinodeSecretKey, "linode-secret-key", "", "Linode: Secret Key")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.S3CompatibleSecretKey, "s3-compatible-secret-key", "", "S3-compatible: Secret Key")
	cmdUpdate.Flags().StringVar(&paramsRemoteServerCredentialUpdate.WasabiSecretKey, "wasabi-secret-key", "", "Wasabi: Secret Key")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	RemoteServerCredentials.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsRemoteServerCredentialDelete := files_sdk.RemoteServerCredentialDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Remote Server Credential`,
		Long:  `Delete Remote Server Credential`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := remote_server_credential.Client{Config: config}

			var err error
			err = client.Delete(paramsRemoteServerCredentialDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsRemoteServerCredentialDelete.Id, "id", 0, "Remote Server Credential ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	RemoteServerCredentials.AddCommand(cmdDelete)
	return RemoteServerCredentials
}
