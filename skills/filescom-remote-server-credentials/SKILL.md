---
name: filescom-remote-server-credentials
description: |
  A RemoteServerCredential is a way to store a credential for Remote Servers in a centralized vault and then reference it from Remote Server definitions.
---

# filescom-remote-server-credentials

A RemoteServerCredential is a way to store a credential for Remote Servers in a centralized vault and then reference it from Remote Server definitions.

This allows you to manage your credentials in one place and avoid duplicating them across multiple Remote Server configurations. It also enhances security by allowing you to use Terraform or APIs for Remote Server management without having to worry about credential exposure.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli remote-server-credentials list`

List Remote Server Credentials.

| Flag | Type | Description |
| --- | --- | --- |
| `--sort-by` | object | If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `workspace_id` and `id`. |
| `--filter` | object | If set, return records where the specified field is equal to the supplied value. Valid fields are `workspace_id` and `name`. Valid field combinations are `[ workspace_id, name ]`. |
| `--filter-prefix` | object | If set, return records where the specified field is prefixed by the supplied value. Valid fields are `name`. |

### `files-cli remote-server-credentials find`

Show Remote Server Credential.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Remote Server Credential ID. **Required.** |

### `files-cli remote-server-credentials create`

Create Remote Server Credential.

| Flag | Type | Description |
| --- | --- | --- |
| `--name` | string | Internal name for your reference |
| `--description` | string | Internal description for your reference |
| `--server-type` | enum | Remote server type.  Remote Server Credentials are only valid for a single type of Remote Server. One of: `ftp`, `sftp`, `s3`, `google_cloud_storage`, `webdav`, `wasabi`, `backblaze_b2`, `one_drive`, `box`, `dropbox`, `google_drive`, `azure`, `sharepoint`, `s3_compatible`, `azure_files`, `files_agent`, `filebase`, `cloudflare`, `linode`. |
| `--aws-access-key` | string | AWS Access Key. |
| `--s3-assume-role-arn` | string | AWS IAM Role ARN for AssumeRole authentication. |
| `--s3-assume-role-duration-seconds` | int64 | Session duration in seconds for AssumeRole authentication (900-43200). |
| `--cloudflare-access-key` | string | Cloudflare: Access Key. |
| `--filebase-access-key` | string | Filebase: Access Key. |
| `--google-cloud-storage-s3-compatible-access-key` | string | Google Cloud Storage: S3-compatible Access Key. |
| `--linode-access-key` | string | Linode: Access Key |
| `--s3-compatible-access-key` | string | S3-compatible: Access Key |
| `--username` | string | Remote server username. |
| `--wasabi-access-key` | string | Wasabi: Access Key. |
| `--password` | string | Password, if needed. |
| `--private-key` | string | Private key, if needed. |
| `--private-key-passphrase` | string | Passphrase for private key if needed. |
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
| `--workspace-id` | int64 | Workspace ID (0 for default workspace) |
| `--copy-values-from-credential-id` | int64 | ID of Remote Server Credential to copy omitted values from. |

### `files-cli remote-server-credentials update`

Update Remote Server Credential.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Remote Server Credential ID. **Required.** |
| `--name` | string | Internal name for your reference |
| `--description` | string | Internal description for your reference |
| `--server-type` | enum | Remote server type.  Remote Server Credentials are only valid for a single type of Remote Server. One of: `ftp`, `sftp`, `s3`, `google_cloud_storage`, `webdav`, `wasabi`, `backblaze_b2`, `one_drive`, `box`, `dropbox`, `google_drive`, `azure`, `sharepoint`, `s3_compatible`, `azure_files`, `files_agent`, `filebase`, `cloudflare`, `linode`. |
| `--aws-access-key` | string | AWS Access Key. |
| `--s3-assume-role-arn` | string | AWS IAM Role ARN for AssumeRole authentication. |
| `--s3-assume-role-duration-seconds` | int64 | Session duration in seconds for AssumeRole authentication (900-43200). |
| `--cloudflare-access-key` | string | Cloudflare: Access Key. |
| `--filebase-access-key` | string | Filebase: Access Key. |
| `--google-cloud-storage-s3-compatible-access-key` | string | Google Cloud Storage: S3-compatible Access Key. |
| `--linode-access-key` | string | Linode: Access Key |
| `--s3-compatible-access-key` | string | S3-compatible: Access Key |
| `--username` | string | Remote server username. |
| `--wasabi-access-key` | string | Wasabi: Access Key. |
| `--password` | string | Password, if needed. |
| `--private-key` | string | Private key, if needed. |
| `--private-key-passphrase` | string | Passphrase for private key if needed. |
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

### `files-cli remote-server-credentials delete`

Delete Remote Server Credential.

| Flag | Type | Description |
| --- | --- | --- |
| `--id` | int64 | Remote Server Credential ID. **Required.** |

