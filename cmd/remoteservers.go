package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/remoteserver"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = remote_server.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
)

var (
	RemoteServers = &cobra.Command{
		Use:  "remote-servers [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func RemoteServersInit() {
	var fieldsList string
	paramsRemoteServerList := files_sdk.RemoteServerListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsRemoteServerList
			params.MaxPages = MaxPagesList
			it := remote_server.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsRemoteServerList.Page, "page", "p", 0, "List Remote Servers")
	cmdList.Flags().IntVarP(&paramsRemoteServerList.PerPage, "per-page", "e", 0, "List Remote Servers")
	cmdList.Flags().StringVarP(&paramsRemoteServerList.Action, "action", "a", "", "List Remote Servers")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
	RemoteServers.AddCommand(cmdList)
	var fieldsFind string
	paramsRemoteServerFind := files_sdk.RemoteServerFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := remote_server.Find(paramsRemoteServerFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsFind)
		},
	}
	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "f", "", "comma separated list of field names")
	RemoteServers.AddCommand(cmdFind)
	var fieldsCreate string
	paramsRemoteServerCreate := files_sdk.RemoteServerCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := remote_server.Create(paramsRemoteServerCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.AwsAccessKey, "aws-access-key", "k", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.AwsSecretKey, "aws-secret-key", "e", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.Password, "password", "p", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.PrivateKey, "private-key", "v", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.GoogleCloudStorageCredentialsJson, "google-cloud-storage-credentials-json", "j", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.WasabiAccessKey, "wasabi-access-key", "", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.WasabiSecretKey, "wasabi-secret-key", "", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.BackblazeB2KeyId, "backblaze-b2-key-id", "i", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.BackblazeB2ApplicationKey, "backblaze-b2-application-key", "", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.RackspaceApiKey, "rackspace-api-key", "", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.AzureBlobStorageAccessKey, "azure-blob-storage-access-key", "y", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.Hostname, "hostname", "o", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.Name, "name", "", "", "Create Remote Server")
	cmdCreate.Flags().IntVarP(&paramsRemoteServerCreate.MaxConnections, "max-connections", "x", 0, "Create Remote Server")
	cmdCreate.Flags().IntVarP(&paramsRemoteServerCreate.Port, "port", "t", 0, "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.S3Bucket, "s3-bucket", "", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.S3Region, "s3-region", "", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.ServerCertificate, "server-certificate", "", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.ServerHostKey, "server-host-key", "", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.ServerType, "server-type", "", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.Ssl, "ssl", "l", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.Username, "username", "", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.GoogleCloudStorageBucket, "google-cloud-storage-bucket", "u", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.GoogleCloudStorageProjectId, "google-cloud-storage-project-id", "d", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.BackblazeB2Bucket, "backblaze-b2-bucket", "b", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.BackblazeB2S3Endpoint, "backblaze-b2-s3-endpoint", "n", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.WasabiBucket, "wasabi-bucket", "", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.WasabiRegion, "wasabi-region", "", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.RackspaceUsername, "rackspace-username", "s", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.RackspaceRegion, "rackspace-region", "g", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.RackspaceContainer, "rackspace-container", "", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.OneDriveAccountType, "one-drive-account-type", "r", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.AzureBlobStorageAccount, "azure-blob-storage-account", "a", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&paramsRemoteServerCreate.AzureBlobStorageContainer, "azure-blob-storage-container", "c", "", "Create Remote Server")
	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
	RemoteServers.AddCommand(cmdCreate)
	var fieldsUpdate string
	paramsRemoteServerUpdate := files_sdk.RemoteServerUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := remote_server.Update(paramsRemoteServerUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsUpdate)
		},
	}
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.AwsAccessKey, "aws-access-key", "k", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.AwsSecretKey, "aws-secret-key", "e", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.Password, "password", "p", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.PrivateKey, "private-key", "v", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.GoogleCloudStorageCredentialsJson, "google-cloud-storage-credentials-json", "j", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.WasabiAccessKey, "wasabi-access-key", "", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.WasabiSecretKey, "wasabi-secret-key", "", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.BackblazeB2KeyId, "backblaze-b2-key-id", "i", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.BackblazeB2ApplicationKey, "backblaze-b2-application-key", "", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.RackspaceApiKey, "rackspace-api-key", "", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.AzureBlobStorageAccessKey, "azure-blob-storage-access-key", "y", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.Hostname, "hostname", "o", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.Name, "name", "", "", "Update Remote Server")
	cmdUpdate.Flags().IntVarP(&paramsRemoteServerUpdate.MaxConnections, "max-connections", "x", 0, "Update Remote Server")
	cmdUpdate.Flags().IntVarP(&paramsRemoteServerUpdate.Port, "port", "t", 0, "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.S3Bucket, "s3-bucket", "", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.S3Region, "s3-region", "", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.ServerCertificate, "server-certificate", "", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.ServerHostKey, "server-host-key", "", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.ServerType, "server-type", "", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.Ssl, "ssl", "l", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.Username, "username", "", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.GoogleCloudStorageBucket, "google-cloud-storage-bucket", "u", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.GoogleCloudStorageProjectId, "google-cloud-storage-project-id", "d", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.BackblazeB2Bucket, "backblaze-b2-bucket", "b", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.BackblazeB2S3Endpoint, "backblaze-b2-s3-endpoint", "n", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.WasabiBucket, "wasabi-bucket", "", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.WasabiRegion, "wasabi-region", "", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.RackspaceUsername, "rackspace-username", "s", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.RackspaceRegion, "rackspace-region", "g", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.RackspaceContainer, "rackspace-container", "", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.OneDriveAccountType, "one-drive-account-type", "r", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.AzureBlobStorageAccount, "azure-blob-storage-account", "a", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&paramsRemoteServerUpdate.AzureBlobStorageContainer, "azure-blob-storage-container", "c", "", "Update Remote Server")
	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "f", "", "comma separated list of field names")
	RemoteServers.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsRemoteServerDelete := files_sdk.RemoteServerDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := remote_server.Delete(paramsRemoteServerDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}
	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
	RemoteServers.AddCommand(cmdDelete)
}
