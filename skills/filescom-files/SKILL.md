---
name: filescom-files
description: |
  Files.com Files via files-cli.
---

# filescom-files

Files.com Files management via files-cli.

All subcommands also accept the global flags documented in [`CONTEXT.md`](../../CONTEXT.md) (`--api-key`, `--format`, `--workspace-id`, `--debug`, and the pagination flags `--cursor` / `--per-page` / `--max-pages` on `list`). Those are not repeated below.

## Commands

### `files-cli files download [path]`

Download File.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Path to operate on. **Required.** |
| `--preview-size` | string | Request a preview size.  Can be `small` (default), `large`, `xlarge`, or `pdf`. |
| `--with-previews` | bool | Include file preview information? |
| `--with-priority-color` | bool | Include file priority color information? |

### `files-cli files create [path]`

Upload File.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Path to operate on. **Required.** |
| `--length` | int64 | Length of file. |
| `--mkdir-parents` | bool | Create parent directories if they do not exist? |
| `--part` | int64 | Part if uploading a part. |
| `--parts` | int64 | How many parts to fetch? |
| `--provided-mtime` | datetime | User provided modification time. |
| `--ref` | string | (no description) |
| `--restart` | int64 | File byte offset to restart from. |
| `--size` | int64 | Size of file. |
| `--copy-behaviors` | bool | If copying a folder, also copy supported behaviors to the destination folder tree? |
| `--structure` | string | If copying folder, copy just the structure? |
| `--with-rename` | bool | Allow file rename instead of overwrite? |
| `--buffered-upload` | bool | If true, and the path refers to a destination not stored on Files.com (such as a remote server mount), the upload will be uploaded first to Files.com before being sent to the remote server mount. This can allow clients to upload using parallel parts to a remote server destination that does not offer parallel parts support natively. |

### `files-cli files update [path]`

Update File/Folder Metadata.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Path to operate on. **Required.** |
| `--custom-metadata` | object | Custom metadata map of keys and values. Limited to 32 keys, 256 characters per key and 1024 characters per value. |
| `--provided-mtime` | datetime | Modified time of file. |
| `--priority-color` | string | Priority/Bookmark color of file. |

### `files-cli files delete [path]`

Delete File/Folder.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Path to operate on. **Required.** |
| `--recursive` | bool | If true, will recursively delete folders.  Otherwise, will error on non-empty folders. |

### `files-cli files find [path]`

Find File/Folder by Path.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Path to operate on. **Required.** |
| `--preview-size` | string | Request a preview size.  Can be `small` (default), `large`, `xlarge`, or `pdf`. |
| `--with-previews` | bool | Include file preview information? |
| `--with-priority-color` | bool | Include file priority color information? |

### `files-cli files zip-list-contents [path]`

List the contents of a ZIP file.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Path to operate on. **Required.** |

### `files-cli files copy [path]`

Copy File/Folder.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Path to operate on. **Required.** |
| `--destination` | string | Copy destination path. **Required.** |
| `--copy-behaviors` | bool | If copying a folder, also copy supported behaviors to the destination folder tree? |
| `--structure` | bool | Copy structure only? |
| `--overwrite` | bool | Overwrite existing file(s) in the destination? |

### `files-cli files move [path]`

Move File/Folder.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Path to operate on. **Required.** |
| `--destination` | string | Move destination path. **Required.** |
| `--overwrite` | bool | Overwrite existing file(s) in the destination? |

### `files-cli files transform [path]`

Transform a file and save the output to a destination path.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Path to operate on. **Required.** |
| `--destination` | string | Destination file path for the transformed output. **Required.** |
| `--transform-type` | string | Transform type. Supported values are `image_convert` and `document_convert`. **Required.** |
| `--target-format` | string | Destination format to create. **Required.** |
| `--width` | int64 | Maximum output width for image_convert. |
| `--height` | int64 | Maximum output height for image_convert. |
| `--overwrite` | bool | Overwrite existing file in the destination? |

### `files-cli files gpg-decrypt [path]`

Decrypt a GPG-encrypted file and save it to a destination path.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Path to operate on. **Required.** |
| `--destination` | string | Destination file path for the decrypted file. **Required.** |
| `--gpg-key-ids` | []int64 | GPG Key IDs to decrypt with. If omitted, every accessible private GPG key in the source workspace is used. |
| `--gpg-key-partner-id` | int64 | Partner ID whose GPG keys should be used for decryption. |
| `--use-all-private-keys` | bool | Use every accessible private GPG key in the source workspace for decryption. |
| `--ignore-mdc-error` | bool | Ignore errors from the MDC (modification detection code) check. |
| `--overwrite` | bool | Overwrite existing file in the destination? |

### `files-cli files gpg-encrypt [path]`

Encrypt a file with GPG and save it to a destination path.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Path to operate on. **Required.** |
| `--destination` | string | Destination file path for the encrypted file. **Required.** |
| `--gpg-key-ids` | []int64 | GPG Key IDs to encrypt with. |
| `--gpg-key-partner-id` | int64 | Partner ID whose GPG keys should be used for encryption. |
| `--signing-key-id` | int64 | Optional GPG Key ID to sign with. |
| `--armor` | bool | Output ASCII-armored encrypted data. |
| `--overwrite` | bool | Overwrite existing file in the destination? |

### `files-cli files unzip [path]`

Extract a ZIP file to a destination folder.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | ZIP file path to extract. **Required.** |
| `--destination` | string | Destination folder path for extracted files. **Required.** |
| `--filename` | string | Optional single entry filename to extract. |
| `--overwrite` | bool | Overwrite existing files in the destination? |

### `files-cli files zip`

Create a ZIP from one or more paths and save it to a destination path.

| Flag | Type | Description |
| --- | --- | --- |
| `--paths` | []string | Paths to include in the ZIP. **Required.** |
| `--destination` | string | Destination file path for the ZIP. **Required.** |
| `--overwrite` | bool | Overwrite existing file in the destination? |

### `files-cli files begin-upload [path]`

Begin File Upload.

| Flag | Type | Description |
| --- | --- | --- |
| `--path` | string | Path to operate on. **Required.** |
| `--mkdir-parents` | bool | Create parent directories if they do not exist? |
| `--part` | int64 | Part if uploading a part. |
| `--parts` | int64 | How many parts to fetch? |
| `--ref` | string | (no description) |
| `--restart` | int64 | File byte offset to restart from. |
| `--size` | int64 | Total bytes of file being uploaded (include bytes being retained if appending/restarting). |
| `--with-rename` | bool | Allow file rename instead of overwrite? |
| `--buffered-upload` | bool | If true, and the path refers to a destination not stored on Files.com (such as a remote server mount), the upload will be uploaded first to Files.com before being sent to the remote server mount. This can allow clients to upload using parallel parts to a remote server destination that does not offer parallel parts support natively. |

