---
name: filescom-remote-servers
description: |
  A RemoteServer is a specific type of Behavior called `remote_server_sync`.
---

# filescom-remote-servers

A RemoteServer is a specific type of Behavior called `remote_server_sync`.

Remote Servers can be either an FTP server, SFTP server, S3 bucket, Google Cloud Storage, Wasabi, Backblaze B2 Cloud Storage, Rackspace Cloud Files container, WebDAV, Box, Dropbox, OneDrive, Google Drive, Azure Blob Storage, or Files.com direct link.

Not every attribute will apply to every remote server.

FTP Servers require that you specify their `hostname`, `port`, `username`, `password`, and a value for `ssl`.  Optionally, provide `server_certificate`.

SFTP Servers require that you specify their `hostname`, `port`, `username`, `password` or `private_key`, and a value for `ssl`.  Optionally, provide `server_certificate`, `private_key_passphrase`.

S3 Buckets require that you specify their `s3_bucket` name, and `s3_region`.  Optionally provide a `aws_access_key`, and `aws_secret_key`.  If you don't provide credentials, you will need to use AWS to grant us access to your bucket.

S3-Compatible Buckets require that you specify `s3_compatible_bucket`, `s3_compatible_endpoint`, `s3_compatible_access_key`, and `s3_compatible_secret_key`. Optionally provide `s3_compatible_virtual_hosted_style` to use virtual-hosted-style URLs instead of path-style URLs.

Google Cloud Storage requires that you specify `google_cloud_storage_bucket`, and then one of the following sets of authentication credentials, selected by `google_cloud_storage_authentication_method` (defaults to `json`):
  - for JSON authentication: `google_cloud_storage_project_id`, and `google_cloud_storage_credentials_json`
  - for HMAC (S3-Compatible) authentication: `google_cloud_storage_s3_compatible_access_key`, and `google_cloud_storage_s3_compatible_secret_key`
  - for OAuth authentication: `google_cloud_storage_oauth_scope`, then follow the `auth_setup_link` and login with Google

Wasabi requires `wasabi_bucket`, `wasabi_region`, `wasabi_access_key`, and `wasabi_secret_key`.

Backblaze B2 Cloud Storage `backblaze_b2_bucket`, `backblaze_b2_s3_endpoint`, `backblaze_b2_application_key`, and `backblaze_b2_key_id`. (Requires S3 Compatible API) See https://help.backblaze.com/hc/en-us/articles/360047425453

WebDAV Servers require that you specify their `hostname`, `username`, and `password`.

OneDrive follow the `auth_setup_link` and login with Microsoft.

Sharepoint follow the `auth_setup_link` and login with Microsoft.

Box follow the `auth_setup_link` and login with Box.

Dropbox specify if `dropbox_teams` then follow the `auth_setup_link` and login with Dropbox.

Google Drive follow the `auth_setup_link` and login with Google.

Azure Blob Storage `azure_blob_storage_account`, `azure_blob_storage_container`, `azure_blob_storage_access_key`, `azure_blob_storage_sas_token`, `azure_blob_storage_dns_suffix`

Azure File Storage `azure_files_storage_account`, `azure_files_storage_access_key`, `azure_files_storage_share_name`, `azure_files_storage_dns_suffix`

Filebase requires `filebase_bucket`, `filebase_access_key`, and `filebase_secret_key`.

Cloudflare requires `cloudflare_bucket`, `cloudflare_access_key`, `cloudflare_secret_key` and `cloudflare_endpoint`.

Linode requires `linode_bucket`, `linode_access_key`, `linode_secret_key` and `linode_region`.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli remote-servers list`

List Remote Servers.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id`, `name`, `server_type`, `backblaze_b2_bucket`, `google_cloud_storage_bucket`, `wasabi_bucket`, `s3_bucket`, `azure_blob_storage_container`, `azure_files_storage_share_name`, `s3_compatible_bucket`, `filebase_bucket`, `cloudflare_bucket` or `linode_bucket`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `name`, `server_type`, `workspace_id`, `backblaze_b2_bucket`, `google_cloud_storage_bucket`, `wasabi_bucket`, `s3_bucket`, `azure_blob_storage_container`, `azure_files_storage_share_name`, `s3_compatible_bucket`, `filebase_bucket`, `cloudflare_bucket` or `linode_bucket`. Valid field combinations are `[ server_type, name ]`, `[ workspace_id, name ]`, `[ backblaze_b2_bucket, name ]`, `[ google_cloud_storage_bucket, name ]`, `[ wasabi_bucket, name ]`, `[ s3_bucket, name ]`, `[ azure_blob_storage_container, name ]`, `[ azure_files_storage_share_name, name ]`, `[ s3_compatible_bucket, name ]`, `[ filebase_bucket, name ]`, `[ cloudflare_bucket, name ]`, `[ linode_bucket, name ]`, `[ workspace_id, server_type ]` or `[ workspace_id, server_type, name ]`. |
| `--filter-prefix` | object | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `name`, `backblaze_b2_bucket`, `google_cloud_storage_bucket`, `wasabi_bucket`, `s3_bucket`, `azure_blob_storage_container`, `azure_files_storage_share_name`, `s3_compatible_bucket`, `filebase_bucket`, `cloudflare_bucket` or `linode_bucket`. Valid field combinations are `[ backblaze_b2_bucket, name ]`, `[ google_cloud_storage_bucket, name ]`, `[ wasabi_bucket, name ]`, `[ s3_bucket, name ]`, `[ azure_blob_storage_container, name ]`, `[ azure_files_storage_share_name, name ]`, `[ s3_compatible_bucket, name ]`, `[ filebase_bucket, name ]`, `[ cloudflare_bucket, name ]` or `[ linode_bucket, name ]`. |

### `files-cli remote-servers find`

Show Remote Server.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Remote Server ID. **Required.** |

### `files-cli remote-servers find-configuration-file`

Download configuration file (required for some Remote Server integrations, such as the Files.com Agent).

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Remote Server ID. **Required.** |

### `files-cli remote-servers create`

Create Remote Server.

| Flag | Type | Description |
| --- | --- | --- |
| `--password` | string | Password, if needed. |
| `--private-key` | string | Private key, if needed. |
| `--private-key-passphrase` | string | Passphrase for private key if needed. |
| `--reset-authentication` | bool | Reset authenticated account? |
| `--ssl-certificate` | string | SSL client certificate. |
| `--aws-secret-key` | string | AWS: secret key. |
| `--azure-blob-storage-access-key` | string | Azure Blob Storage: Access Key |
| `--azure-blob-storage-sas-token` | string | Azure Blob Storage: Shared Access Signature (SAS) token |
| `--azure-files-storage-access-key` | string | Azure File Storage: Access Key |
| `--azure-files-storage-sas-token` | string | Azure File Storage: Shared Access Signature (SAS) token |
| `--backblaze-b2-application-key` | string | Backblaze B2 Cloud Storage: applicationKey |
| `--backblaze-b2-key-id` | string | Backblaze B2 Cloud Storage: keyID |
| `--cloudflare-secret-key` | string | Cloudflare: Secret Key |
| `--filebase-secret-key` | string | Filebase: Secret Key |
| `--google-cloud-storage-credentials-json` | string | Google Cloud Storage: JSON file that contains the private key. To generate see https://cloud.google.com/storage/docs/json_api/v1/how-tos/authorizing#APIKey |
| `--google-cloud-storage-s3-compatible-secret-key` | string | Google Cloud Storage: S3-compatible secret key |
| `--linode-secret-key` | string | Linode: Secret Key |
| `--s3-compatible-secret-key` | string | S3-compatible: Secret Key |
| `--wasabi-secret-key` | string | Wasabi: Secret Key |
| `--allow-relative-paths` | bool | Allow relative paths in SFTP. If true, paths will not be forced to be absolute, allowing operations relative to the user's home directory. |
| `--aws-access-key` | string | AWS Access Key. |
| `--azure-blob-storage-account` | string | Azure Blob Storage: Account name |
| `--azure-blob-storage-container` | string | Azure Blob Storage: Container name |
| `--azure-blob-storage-dns-suffix` | string | Azure Blob Storage: Custom DNS suffix |
| `--azure-blob-storage-hierarchical-namespace` | bool | Azure Blob Storage: Does the storage account has hierarchical namespace feature enabled? |
| `--azure-files-storage-account` | string | Azure Files: Storage Account name |
| `--azure-files-storage-dns-suffix` | string | Azure Files: Custom DNS suffix |
| `--azure-files-storage-share-name` | string | Azure Files:  Storage Share name |
| `--backblaze-b2-bucket` | string | Backblaze B2 Cloud Storage: Bucket name |
| `--backblaze-b2-s3-endpoint` | string | Backblaze B2 Cloud Storage: S3 Endpoint |
| `--buffer-uploads` | enum | If set to always, uploads to this server will be uploaded first to Files.com before being sent to the remote server. This can improve performance in certain access patterns, such as high-latency connections.  It will cause data to be temporarily stored in Files.com. If set to auto, we will perform this optimization if we believe it to be a benefit in a given situation. One of: `auto`, `always`, `never`. |
| `--cloudflare-access-key` | string | Cloudflare: Access Key. |
| `--cloudflare-bucket` | string | Cloudflare: Bucket name |
| `--cloudflare-endpoint` | string | Cloudflare: endpoint |
| `--description` | string | Internal description for your reference |
| `--dropbox-teams` | bool | Dropbox: If true, list Team folders in root? |
| `--enable-dedicated-ips` | bool | `true` if remote server only accepts connections from dedicated IPs |
| `--filebase-access-key` | string | Filebase: Access Key. |
| `--filebase-bucket` | string | Filebase: Bucket name |
| `--files-api-key` | string | Files.com direct link: API key used once to pair the remote server. |
| `--files-agent-permission-set` | enum | Local permissions for files agent. read_only, write_only, or read_write. One of: `read_write`, `read_only`, `write_only`. |
| `--files-agent-root` | string | Agent local root path |
| `--files-agent-version` | string | Files Agent version |
| `--outbound-agent-id` | int64 | Route traffic to outbound on a files-agent |
| `--google-cloud-storage-authentication-method` | enum | Google Cloud Storage: Authentication method. Can be json, hmac, or oauth. One of: `json`, `hmac`, `oauth`. |
| `--google-cloud-storage-bucket` | string | Google Cloud Storage: Bucket Name |
| `--google-cloud-storage-oauth-scope` | string | Google Cloud Storage: OAuth scope. Can be https://www.googleapis.com/auth/devstorage.read_only or https://www.googleapis.com/auth/devstorage.read_write. |
| `--google-cloud-storage-project-id` | string | Google Cloud Storage: Project ID |
| `--google-cloud-storage-s3-compatible-access-key` | string | Google Cloud Storage: S3-compatible Access Key. |
| `--hostname` | string | Hostname or IP address |
| `--linode-access-key` | string | Linode: Access Key |
| `--linode-bucket` | string | Linode: Bucket name |
| `--linode-region` | string | Linode: region |
| `--max-connections` | int64 | Max number of parallel connections.  Ignored for S3 connections (we will parallelize these as much as possible). |
| `--name` | string | Internal name for your reference |
| `--one-drive-account-type` | enum | OneDrive: Either personal or business_other account types. One of: `personal`, `business_other`. |
| `--pin-to-site-region` | bool | If true, we will ensure that all communications with this remote server are made through the primary region of the site.  This setting can also be overridden by a site-wide setting which will force it to true. |
| `--port` | int64 | Port for remote server. |
| `--upload-staging-path` | string | Upload staging path.  Applies to SFTP only.  If a path is provided here, files will first be uploaded to this path on the remote folder and the moved into the final correct path via an SFTP move command.  This is required by some remote MFT systems to emulate atomic uploads, which are otherwise not supoprted by SFTP. |
| `--remote-server-credential-id` | int64 | ID of Remote Server Credential, if applicable. |
| `--s3-assume-role-arn` | string | AWS IAM Role ARN for AssumeRole authentication. |
| `--s3-assume-role-duration-seconds` | int64 | Session duration in seconds for AssumeRole authentication (900-43200). |
| `--s3-bucket` | string | S3 bucket name |
| `--s3-compatible-access-key` | string | S3-compatible: Access Key |
| `--s3-compatible-bucket` | string | S3-compatible: Bucket name |
| `--s3-compatible-endpoint` | string | S3-compatible: endpoint |
| `--s3-compatible-region` | string | S3-compatible: region |
| `--s3-compatible-virtual-hosted-style` | bool | S3-compatible: If true, use virtual-hosted-style URLs instead of path-style URLs |
| `--s3-region` | string | S3 region |
| `--server-certificate` | enum | Remote server certificate. One of: `require_match`, `allow_any`. |
| `--server-host-key` | string | Remote server SSH Host Key. If provided, we will require that the server host key matches the provided key. Uses OpenSSH format similar to what would go into ~/.ssh/known_hosts |
| `--server-type` | enum | Remote server type. One of: `ftp`, `sftp`, `s3`, `google_cloud_storage`, `webdav`, `wasabi`, `backblaze_b2`, `one_drive`, `box`, `dropbox`, `google_drive`, `azure`, `sharepoint`, `s3_compatible`, `azure_files`, `files_agent`, `filebase`, `cloudflare`, `linode`, `files_com`. |
| `--ssl` | enum | Should we require SSL?. One of: `if_available`, `require`, `require_implicit`, `never`. |
| `--username` | string | Remote server username. |
| `--wasabi-access-key` | string | Wasabi: Access Key. |
| `--wasabi-bucket` | string | Wasabi: Bucket name |
| `--wasabi-region` | string | Wasabi: Region |
| `--workspace-id` | int64 | Workspace ID (0 for default workspace) |

### `files-cli remote-servers agent-push-update`

Push update to Files Agent.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Remote Server ID. **Required.** |

### `files-cli remote-servers configuration-file`

Post local changes, check in, and download configuration file (used by some Remote Server integrations, such as the Files.com Agent).

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Remote Server ID. **Required.** |
| `--api-token` | string | Files Agent API Token |
| `--permission-set` | string | The permission set for the agent ['read_write', 'read_only', 'write_only'] |
| `--root` | string | The root directory for the agent |
| `--hostname` | string | (no description) |
| `--port` | int64 | Incoming port for files agent connections |
| `--status` | string | either running or shutdown |
| `--config-version` | string | agent config version |
| `--private-key` | string | The private key for the agent |
| `--public-key` | string | public key |
| `--server-host-key` | string | (no description) |
| `--subdomain` | string | Files.com subdomain site name |

### `files-cli remote-servers update`

Update Remote Server.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Remote Server ID. **Required.** |
| `--password` | string | Password, if needed. |
| `--private-key` | string | Private key, if needed. |
| `--private-key-passphrase` | string | Passphrase for private key if needed. |
| `--reset-authentication` | bool | Reset authenticated account? |
| `--ssl-certificate` | string | SSL client certificate. |
| `--aws-secret-key` | string | AWS: secret key. |
| `--azure-blob-storage-access-key` | string | Azure Blob Storage: Access Key |
| `--azure-blob-storage-sas-token` | string | Azure Blob Storage: Shared Access Signature (SAS) token |
| `--azure-files-storage-access-key` | string | Azure File Storage: Access Key |
| `--azure-files-storage-sas-token` | string | Azure File Storage: Shared Access Signature (SAS) token |
| `--backblaze-b2-application-key` | string | Backblaze B2 Cloud Storage: applicationKey |
| `--backblaze-b2-key-id` | string | Backblaze B2 Cloud Storage: keyID |
| `--cloudflare-secret-key` | string | Cloudflare: Secret Key |
| `--filebase-secret-key` | string | Filebase: Secret Key |
| `--google-cloud-storage-credentials-json` | string | Google Cloud Storage: JSON file that contains the private key. To generate see https://cloud.google.com/storage/docs/json_api/v1/how-tos/authorizing#APIKey |
| `--google-cloud-storage-s3-compatible-secret-key` | string | Google Cloud Storage: S3-compatible secret key |
| `--linode-secret-key` | string | Linode: Secret Key |
| `--s3-compatible-secret-key` | string | S3-compatible: Secret Key |
| `--wasabi-secret-key` | string | Wasabi: Secret Key |
| `--allow-relative-paths` | bool | Allow relative paths in SFTP. If true, paths will not be forced to be absolute, allowing operations relative to the user's home directory. |
| `--aws-access-key` | string | AWS Access Key. |
| `--azure-blob-storage-account` | string | Azure Blob Storage: Account name |
| `--azure-blob-storage-container` | string | Azure Blob Storage: Container name |
| `--azure-blob-storage-dns-suffix` | string | Azure Blob Storage: Custom DNS suffix |
| `--azure-blob-storage-hierarchical-namespace` | bool | Azure Blob Storage: Does the storage account has hierarchical namespace feature enabled? |
| `--azure-files-storage-account` | string | Azure Files: Storage Account name |
| `--azure-files-storage-dns-suffix` | string | Azure Files: Custom DNS suffix |
| `--azure-files-storage-share-name` | string | Azure Files:  Storage Share name |
| `--backblaze-b2-bucket` | string | Backblaze B2 Cloud Storage: Bucket name |
| `--backblaze-b2-s3-endpoint` | string | Backblaze B2 Cloud Storage: S3 Endpoint |
| `--buffer-uploads` | enum | If set to always, uploads to this server will be uploaded first to Files.com before being sent to the remote server. This can improve performance in certain access patterns, such as high-latency connections.  It will cause data to be temporarily stored in Files.com. If set to auto, we will perform this optimization if we believe it to be a benefit in a given situation. One of: `auto`, `always`, `never`. |
| `--cloudflare-access-key` | string | Cloudflare: Access Key. |
| `--cloudflare-bucket` | string | Cloudflare: Bucket name |
| `--cloudflare-endpoint` | string | Cloudflare: endpoint |
| `--description` | string | Internal description for your reference |
| `--dropbox-teams` | bool | Dropbox: If true, list Team folders in root? |
| `--enable-dedicated-ips` | bool | `true` if remote server only accepts connections from dedicated IPs |
| `--filebase-access-key` | string | Filebase: Access Key. |
| `--filebase-bucket` | string | Filebase: Bucket name |
| `--files-api-key` | string | Files.com direct link: API key used once to pair the remote server. |
| `--files-agent-permission-set` | enum | Local permissions for files agent. read_only, write_only, or read_write. One of: `read_write`, `read_only`, `write_only`. |
| `--files-agent-root` | string | Agent local root path |
| `--files-agent-version` | string | Files Agent version |
| `--outbound-agent-id` | int64 | Route traffic to outbound on a files-agent |
| `--google-cloud-storage-authentication-method` | enum | Google Cloud Storage: Authentication method. Can be json, hmac, or oauth. One of: `json`, `hmac`, `oauth`. |
| `--google-cloud-storage-bucket` | string | Google Cloud Storage: Bucket Name |
| `--google-cloud-storage-oauth-scope` | string | Google Cloud Storage: OAuth scope. Can be https://www.googleapis.com/auth/devstorage.read_only or https://www.googleapis.com/auth/devstorage.read_write. |
| `--google-cloud-storage-project-id` | string | Google Cloud Storage: Project ID |
| `--google-cloud-storage-s3-compatible-access-key` | string | Google Cloud Storage: S3-compatible Access Key. |
| `--hostname` | string | Hostname or IP address |
| `--linode-access-key` | string | Linode: Access Key |
| `--linode-bucket` | string | Linode: Bucket name |
| `--linode-region` | string | Linode: region |
| `--max-connections` | int64 | Max number of parallel connections.  Ignored for S3 connections (we will parallelize these as much as possible). |
| `--name` | string | Internal name for your reference |
| `--one-drive-account-type` | enum | OneDrive: Either personal or business_other account types. One of: `personal`, `business_other`. |
| `--pin-to-site-region` | bool | If true, we will ensure that all communications with this remote server are made through the primary region of the site.  This setting can also be overridden by a site-wide setting which will force it to true. |
| `--port` | int64 | Port for remote server. |
| `--upload-staging-path` | string | Upload staging path.  Applies to SFTP only.  If a path is provided here, files will first be uploaded to this path on the remote folder and the moved into the final correct path via an SFTP move command.  This is required by some remote MFT systems to emulate atomic uploads, which are otherwise not supoprted by SFTP. |
| `--remote-server-credential-id` | int64 | ID of Remote Server Credential, if applicable. |
| `--s3-assume-role-arn` | string | AWS IAM Role ARN for AssumeRole authentication. |
| `--s3-assume-role-duration-seconds` | int64 | Session duration in seconds for AssumeRole authentication (900-43200). |
| `--s3-bucket` | string | S3 bucket name |
| `--s3-compatible-access-key` | string | S3-compatible: Access Key |
| `--s3-compatible-bucket` | string | S3-compatible: Bucket name |
| `--s3-compatible-endpoint` | string | S3-compatible: endpoint |
| `--s3-compatible-region` | string | S3-compatible: region |
| `--s3-compatible-virtual-hosted-style` | bool | S3-compatible: If true, use virtual-hosted-style URLs instead of path-style URLs |
| `--s3-region` | string | S3 region |
| `--server-certificate` | enum | Remote server certificate. One of: `require_match`, `allow_any`. |
| `--server-host-key` | string | Remote server SSH Host Key. If provided, we will require that the server host key matches the provided key. Uses OpenSSH format similar to what would go into ~/.ssh/known_hosts |
| `--server-type` | enum | Remote server type. One of: `ftp`, `sftp`, `s3`, `google_cloud_storage`, `webdav`, `wasabi`, `backblaze_b2`, `one_drive`, `box`, `dropbox`, `google_drive`, `azure`, `sharepoint`, `s3_compatible`, `azure_files`, `files_agent`, `filebase`, `cloudflare`, `linode`, `files_com`. |
| `--ssl` | enum | Should we require SSL?. One of: `if_available`, `require`, `require_implicit`, `never`. |
| `--username` | string | Remote server username. |
| `--wasabi-access-key` | string | Wasabi: Access Key. |
| `--wasabi-bucket` | string | Wasabi: Bucket name |
| `--wasabi-region` | string | Wasabi: Region |

### `files-cli remote-servers delete`

Delete Remote Server.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Remote Server ID. **Required.** |

