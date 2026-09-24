# Files.com Command Line App

Files.com is the cloud-native, next-gen MFT, SFTP, and secure file-sharing platform that replaces brittle legacy servers with one always-on, secure fabric. Automate mission-critical file flows—across any cloud, protocol, or partner—while supporting human collaboration and eliminating manual work.

With universal SFTP, AS2, HTTPS, and 50+ native connectors backed by military-grade encryption, Files.com unifies governance, visibility, and compliance in a single pane of glass.

The content included here should be enough to get started, but please visit our
[Developer Documentation Website](https://developers.files.com/cli/) for the complete documentation.

## Introduction

The Files.com Command Line Interface (CLI) App is a great option for scripted or automated transfers between a local machine and Files.com.

Because it works through the standard Command Line, the CLI app is easy to script from a variety of environments without having to use our SDKs. With that said, if you are already using a programming language where we offer an SDK, the SDK may offer a higher level of integration for your application.

The CLI App is cross-platform (Windows/macOS/Linux) and supports fast, concurrent file transfers.

The CLI App uses the Files.com RESTful APIs via the HTTPS protocol (port 443) to securely communicate and transfer files so, when used interactively or from a script, no firewall changes should be required in order to allow connectivity.

### Support For All Operations, Not Just File Operations

The CLI supports all file Operations including list, download, upload, move, rename, delete, etc. But equally important is that it supports operations on every resource available in Files.com including Users, Permissions, Groups, Remote Servers, Behaviors, etc.

The available resources are listed in this documentation under the Resources menu on the left.

### Small-File Download Batching

When `files-cli download` or `files-cli sync pull` processes at least 500 eligible small files, it can batch them through the ZIP download path instead of issuing one download request per file. Files are eligible when they are smaller than 128 KiB.

The ZIP is only the transfer format. The CLI extracts each entry as it arrives, verifies it, and writes ordinary files to the requested local destination. It does not save a `.zip` archive.

ZIP batching is enabled by default and adjusts per job. The CLI starts with 16-file ZIP batches, can grow batches up to 32 files per stream, and compares ZIP throughput against normal per-file downloads. It keeps the ZIP path only when that path measures faster, and long jobs can re-check the decision as network conditions change.

Use `--no-zip-batch` to disable ZIP batching for a download or sync pull. Use `--force-zip-batch` for targeted performance testing when you want to bypass the automatic speed check.

### Installation

Download the latest release for Windows, macOS, or Linux from the [CLI App Releases](https://github.com/Files-com/files-cli/releases) page.

On that page, you'll need to pick your exact operating system to download the correct version.

On 64-bit AMD/Intel Windows systems, the MSI installer places the CLI in `C:\Program Files\Files.com-CLI` and adds that directory to the system `%PATH%`. Portable archives are also available and can be stored anywhere on your computer.

When using a portable archive, we recommend placing the app binary into one of the folders listed in your `%PATH%` (Windows) or `$PATH` (Linux and Mac).

Here are specific instructions, grouped by OS:

#### Windows

**Download**

- AMD/Intel 64-Bit Processors [MSI installer](https://github.com/Files-com/files-cli/releases/latest/download/files-cli_windows_64bit.msi) *(Recommended)*
- AMD/Intel 64-Bit Processors [portable ZIP](https://github.com/Files-com/files-cli/releases/latest/download/files-cli_windows_64bit.zip)
- ARM Processors [portable ZIP](https://github.com/Files-com/files-cli/releases/latest/download/files-cli_windows_arm64.zip) *(Less Common)*

Run the MSI installer (`*.msi`) to install `files-cli.exe` and add it to the system `%PATH%`. Open a new terminal before running `files-cli`.

For a portable installation, download the Zip archive (`*.zip`), extract it, and place `files-cli.exe` into any directory listed in your `%PATH%` environment variable.

#### Mac using Homebrew

For systems with Homebrew 6.0 or later, trust the Files.com CLI formula before installing it:

```shell
brew tap files-com/tap
brew trust --formula files-com/tap/files-cli
brew install files-cli
```

The `brew trust --formula files-com/tap/files-cli` command trusts only the `files-cli` formula. To trust the entire Files.com tap instead, run `brew trust files-com/tap` before installing. Homebrew versions earlier than 6.0 can install without the `brew trust` command.

#### Mac without Homebrew

Download the compressed Tar archive (`*.tar.gz`), extract the files from the archive, and place the `files-cli` binary file into any directory listed in your `$PATH` environment variable.

```shell
curl -L https://github.com/Files-com/files-cli/releases/latest/download/files-cli_macOS_64bit.tar.gz | tar zxv
sudo mv ./files-cli /usr/local/bin
```

#### Linux: RPM Package Manager Based Systems

For Linux systems that support RPM, such as Red Hat Linux, Fedora Linux, CentOS, openSUSE, Oracle Linux, and others, use the RPM Package Manager to install the CLI App:

```shell
curl -L https://github.com/Files-com/files-cli/releases/latest/download/files-cli_linux_{ARCH}.rpm -o files-cli.rpm

sudo rpm -i ./files-cli.rpm
```

#### Linux: Debian Based Systems

For Debian based systems, such as Ubuntu Linux, use the APT Package Manager to install the CLI App:

```shell
curl -L https://github.com/Files-com/files-cli/releases/latest/download/files-cli_linux_{ARCH}.deb -o files-cli.deb

sudo apt install ./files-cli.deb
```

Learn how to use the Files.com CLI App by reading [the
documentation](https://www.files.com/docs/client-apps/command-line-interface-cli-app).

Explore the [files-cli](https://github.com/Files-com/files-cli) code on GitHub.

### Command Help

You can get usage information for the `files-cli` program and any of its commands or subcommands by using the
`--help` option.

```shell
files-cli --help
files-cli folders --help
```

### Getting Support

The Files.com Support team provides official support for all of our official Files.com integration tools.

To initiate a support conversation, you can send an [Authenticated Support Request](https://www.files.com/docs/overview/requesting-support) or simply send an E-Mail to support@files.com.

## For Agents

If you're driving `files-cli` from an AI agent — or building an agent that uses it — start with these:

| Resource | Description |
| --- | --- |
| [CONTEXT.md](CONTEXT.md) | CLI-wide invocation contract: authentication, global flags, output format, error conventions. |
| [skills/](skills/) | `SKILL.md` packages — one per top-level CLI command (Bundles, Users, Files, …), grouped by category. |
| [skills/INDEX.md](skills/INDEX.md) | Full index of skills. |
| [agents/tool-catalog.json](agents/tool-catalog.json) | Machine-readable catalog of the API resource commands and their parameters, generated from the API schema. |
| [agents/error-catalog.json](agents/error-catalog.json) | Machine-readable catalog of known error types with HTTP codes. |

Core invocation:

```bash
files-cli <domain> <subcommand> --format json --non-interactive [flags...]
```

- Pass `--format json` so output is structured (the default `table` is for humans).
- Pass `--non-interactive` so the CLI never blocks on a prompt.
- Exit code `0` is success. Non-zero is failure; read diagnostics from stderr. `--format json` does not provide a uniform JSON error envelope.

For Claude Code, Codex, or any agent that supports filesystem-loaded skills, point the skills directory at `skills/`. For agents without skill loading, load the relevant `SKILL.md` directly into context based on the task at hand.

The installed binary can also describe itself, offline and without credentials, matched to its version:

```bash
files-cli commands                               # top-level commands and groups
files-cli commands search share link             # find commands by keyword
files-cli commands describe folders list-for --format json
files-cli workflows show recipe-searching-for-files
```

To read a long list a page at a time, add `--json-envelope` to a list command. It prints `{"has_more", "next_cursor", "data"}` for one page; pass `next_cursor` back with `--cursor` to continue. Without the flag, `--format json` output is unchanged. See CONTEXT.md for details.

## Authentication

There are two ways to authenticate: API Key authentication and Session-based authentication.

### Authenticate with an API Key

Authenticating with an API key is the recommended authentication method for most scenarios, and is
the method used in the examples on this site.

To use an API Key, first generate an API key from the [web
interface](https://www.files.com/docs/sdk-and-apis/api-keys) or [via the API or an
SDK](/cli/resources/developers/api-keys).

Note that when using a user-specific API key, if the user is an administrator, you will have full
access to the entire API. If the user is not an administrator, you will only be able to access files
that user can access, and no access will be granted to site administration functions in the API.

#### Setting by Environment Variable on MacOS or Linux

``` shell
export FILES_API_KEY="YOUR_API_KEY"
```

#### Setting by Environment Variable on Windows

``` shell
set FILES_API_KEY="YOUR_API_KEY"
```

```shell title="Example Request"
files-cli --api-key=YOUR_API_KEY folders list-for ''
## After the key has been provided once it will be written to the files-cli configuration file.
## You do not need to include the same API key for future commands.
```

Don't forget to replace the placeholder, `YOUR_API_KEY`, with your actual API key.

### Authenticate with a Session

You can also authenticate by creating a user session using the username and
password of an active user. If the user is an administrator, the session will have full access to
all capabilities of Files.com. Sessions created from regular user accounts will only be able to access files that
user can access, and no access will be granted to site administration functions.

Sessions use the exact same session timeout settings as web interface sessions. When a
session times out, simply create a new session and resume where you left off. This process is not
automatically handled by our SDKs because we do not want to store password information in memory without
your explicit consent.

#### Logging In

To log in to the CLI App with your username and password, you must first configure the CLI App with information about your Files.com account.

Once you've specified your subdomain information and username, you do not need to specify it again for subsequent uses of the CLI App.

You will be prompted for the password when a command is run that requires authentication.  Once the password is entered, subsequent calls will not require a password, unless the session is no longer valid (ie, expired).

```shell title="Example Request"
files-cli config set --subdomain SUBDOMAIN --username motor
files-cli folders list-for ''
> password: vroom
```

#### Using a Session

Once the password has been verified by the first run command, all subsequent commands can be run and the current verified session will be used.

You can also provide an existing session ID directly with `--session-id`, or store it with `files-cli config set --session-id`.

```shell title="Example Request"
files-cli folders list-for ''

## Or use an existing session ID directly.
files-cli --session-id=YOUR_SESSION_ID folders list-for ''

## You can also store the session ID for future commands.
files-cli config set --session-id YOUR_SESSION_ID
```

#### Logging Out

User sessions can be ended by calling `sessions delete`.

```bash title="Example Request"
files-cli sessions delete
```

## Configuration

The files-cli client can be configured by running `files-cli config set`.

### Using Multiple Accounts

You can use the `--profile` option to modify the configuration of a specific
profile. That same option can be used to specify the profile to use when
interacting with Files.com. This allows you to use multiple Files.com accounts
without needing to reauthenticate when switching between them.

```shell title="Example setting"
files-cli config set --profile firstaccount --username FIRSTUSERNAME
```

### Configuration Options

#### Base URL

Set this to the full https:// URL of your Files.com subdomain (e.g. `https://MY-SUBDOMAIN.files.com`).
This is not required in most cases, but one benefit of setting it is that it ensures that authentication failures will be logged to your site's API logs.  Without setting this, we won't know which site to associate the authentication failure with, and it won't be logged to your site's API logs.
This is always required if your site is configured to disable global acceleration.
This can also be set to use a mock server in development or CI.

```shell title="Example setting"
files-cli config set --endpoint https://MY-SUBDOMAIN.files.com
## alternatively
files-cli config set -e https://MY-SUBDOMAIN.files.com
```

#### Concurrent Connection Limit

Set the maximum number of concurrent connections.

```shell title="Example setting"
files-cli config set --concurrent-connection-limit 5
## alternatively
files-cli config set -c 5
```

#### Session ID

Store an existing user session ID for CLI authentication. You can also provide a session ID for a single command with `--session-id`.

```shell title="Example setting"
files-cli config set --session-id YOUR_SESSION_ID
```

#### Workspace ID

Scope every command to a specific Workspace by storing its ID. You can also scope a single command with `--workspace-id`.

```shell title="Example setting"
files-cli config set --workspace-id YOUR_WORKSPACE_ID
```

#### Default Resource Format

Set default format for displaying resources.

The supported formats are:

* json
* csv
* table

For the `json` format you can specify either `pretty` or `raw` after a comma (`,`).

For the table format you can specify a style (`interactive`, `light`, `dark`, or `bright`)
and/or a direction (`vertical` or `horizontal`) by adding those after commas (e.g. `table,dark,horizontal`).

```shell title="Example setting"
files-cli config set --format json,raw
files-cli config set -f table,bright,vertical
```

## Sort and Filter

Several of the Files.com API resources have list operations that return multiple instances of the
resource. The List operations can be sorted and filtered.

### Sorting

To sort the returned data, use the `--sort-by` parameter.

```shell title="Sort Syntax"
--sort-by="field=asc"
```

Valid directions are `asc` and `desc`.

#### Special note about the List Folder Endpoint

For historical reasons, and to maintain compatibility
with a variety of other cloud-based MFT and EFSS services, Folders will always be listed before Files
when listing a Folder. This applies regardless of the sorting parameters you provide. These *will* be
used, after the initial sort application of Folders before Files.

```shell title="Sort Example" hasDataFormatSelector
files-cli users list \
  --sort-by="username=asc"
```

### Filtering

Filters apply selection criteria to the underlying query that returns the results. They can be
applied individually or combined with other filters, and the resulting data can be sorted by a
single field.

Each resource supports a unique set of valid filter fields, filter combinations, and combinations of
filters and sort fields.

| Flag | Type | Description |
| --------- | --------- | --------- |
| `--filter="field=value"` | Exact | Find resources that have an exact field value match to a passed in value. (i.e., FIELD_VALUE = PASS_IN_VALUE). |
| `--filter-prefix="field=value"` | Pattern | Find resources where the specified field is prefixed by the supplied value. This is applicable to values that are strings. |
| `--filter-gt="field=value"` | Range | Find resources that have a field value that is greater than the passed in value. (i.e., FIELD_VALUE > PASS_IN_VALUE). |
| `--filter-gteq="field=value"` | Range | Find resources that have a field value that is greater than or equal to the passed in value. (i.e., FIELD_VALUE >= PASS_IN_VALUE). |
| `--filter-lt="field=value"` | Range | Find resources that have a field value that is less than the passed in value. (i.e., FIELD_VALUE &lt; PASS_IN_VALUE). |
| `--filter-lteq="field=value"` | Range | Find resources that have a field value that is less than or equal to the passed in value. (i.e., FIELD_VALUE &lt;= PASS_IN_VALUE). |

Use `--filter-by` for client-side wildcard filtering after records are fetched.

```--filter-by="field_name=*.jpg"```

```shell title="Exact Filter Example" hasDataFormatSelector
files-cli users list \
  --filter="not_site_admin=true"
```

```shell title="Pattern Filter Example" hasDataFormatSelector
files-cli users list \
  --filter-prefix="username=test"
```

```shell title="Range Filter Example" hasDataFormatSelector
files-cli action-logs list \
  --filter-gteq="created_at=2024-01-01"
```

```shell title="Combination Filter with Sort Example" hasDataFormatSelector
files-cli users list \
  --filter="not_site_admin=true" \
  --filter-prefix="username=test" \
  --sort-by="username=asc"
```

```shell title="Client-side Filter Example" hasDataFormatSelector
files-cli folders ls some/path \
  --filter-by="path=*.jpg"
```

## Paths

Files.com preserves the spelling of file and folder paths while comparing them using shared case and Unicode rules. Use the SDK comparison helpers when matching paths locally.
<div></div>

### Capitalization

Files.com uses case-insensitive path matching based on its fixed Unicode comparison map.

For example, the following paths have the same comparison key:

| Path Variant                          | Comparison Key              |
|---------------------------------------|------------------------------|
| `Documents/Reports/Q1.pdf`            | `documents/reports/q1.pdf`  |
| `documents/reports/q1.PDF`            | `documents/reports/q1.pdf`  |
| `DOCUMENTS/REPORTS/Q1.PDF`            | `documents/reports/q1.pdf`  |

This behavior applies across:
- API requests
- Folder and file lookup operations
- Automations and workflows

See also: [Case Sensitivity Documentation](https://www.files.com/docs/files-and-folders/case-sensitivity/)

### Slashes

Use `/` between folder and file names, without leading or trailing slashes. SDK normalization helpers convert backslashes to `/`, remove duplicate separators, and discard exact `.` and `..` components. Discarding `..` leaves the preceding folder name intact.

| Input | Normalized path |
|-------|-----------------|
| `folder/subfolder/file.txt` | `folder/subfolder/file.txt` |
| `/folder/subfolder/file.txt` | `folder/subfolder/file.txt` |
| `folder/subfolder/file.txt/` | `folder/subfolder/file.txt` |
| `//folder//file.txt` | `folder/file.txt` |
| `folder/../file.txt` | `folder/file.txt` |

<div></div>

### Unicode and Path Comparison

Files.com compares paths using a fixed mapping shared by the server and SDKs. It treats case and many accent differences as equivalent: `Résumé.txt` and `resume.txt` identify the same file, as do `q` followed by a combining acute accent and `q`. The mapping also handles other equivalences, such as Hiragana and Katakana. Lowercasing or applying a standard Unicode normalization form alone does not reproduce these rules.

SDK comparison helpers normalize path separators and dot segments, then apply the bundled [versioned comparison map](https://github.com/Files-com/files-sdk-javascript/blob/master/shared/path_comparison.json). The [shared examples](https://github.com/Files-com/files-sdk-javascript/blob/master/shared/comparison_examples.json) give exact comparison results for integrations that implement their own matching. The map uses hexadecimal Unicode scalar values as keys: a missing entry preserves the character, an empty replacement removes it, and other replacements may contain several characters. Apply each replacement once without normalizing or lowercasing the result again.

Use comparison results only for matching. Send the original path spelling in API requests and preserve it for display and local filenames; comparison results can have a different spelling or length.

Trailing whitespace is significant for comparison. `report.txt` and `report.txt ` are different file paths, and SDK helpers preserve spaces, tabs, and newlines. Folder names cannot end in whitespace. See [Unicode Normalization](https://www.files.com/docs/files-and-folders/file-system-semantics/unicode-normalization) for the complete path rules.

<div></div>

## Workspaces

A Workspace is a lightweight way to organize related resources inside a single Files.com Site.

Customers commonly group resources by project, department, client, or region. Workspaces provide a built-in structure for that grouping, so the UI can operate within a clear "workspace context" and admins can delegate management for a subset of resources without requiring full site-level isolation.

Every Site has an implicit Default workspace (ID `0`). Resources that are not explicitly assigned to a named workspace are considered part of the Default workspace.

The Files.com CLI supports workspace scoping by using `--workspace-id` for a single command. Configure subsequent commands by running `files-cli config set --workspace-id`.
```shell title="Example Request"
files-cli --workspace-id 123 folders list-for ''

files-cli config set --workspace-id 123
files-cli folders list-for ''
```

## Foreign Language Support

The Files.com CLI supports localized responses by setting the `files-cli config set -l` option.  The language can be reset by using `files-cli config reset -l`.
When configured, this guides the API in selecting a preferred language for applicable response content.

Language support currently applies to select human-facing fields only, such as notification messages
and error descriptions.

If the specified language is not supported or the value is omitted, the API defaults to English.

```shell title="Example Request"
files-cli config set -l es
```

## Errors

The Files.com CLI will detect errors coming back from the API and provide a detailed message to the current
output explaining the cause of the error.

Errors fall into two basic categories:

1. General Usage Errors - errors that result from incorrect usage of the CLI tool
2. API Response errors - errors returned from the Files.com API.

```shell title="Example General Error"

files-cli BADCOMMAND
Error: unknown command "BADCOMMAND" for "files-cli"
Run 'files-cli --help' for usage.

```

```shell title="Example Files.com API Error"

files-cli folders list-for /BADFOLDER

Error: Not Found - `Not Found.  This may be related to your permissions.`

```

### API Error Types

To understand the types of errors that come back from the Files.com API and will be displayed by the CLI, see the [Rest API Errors](/rest/overview/errors/).

## Logs

### Sending Operation/Run Logs to the Cloud

If you are running scripted operations, you can have the CLI send a report of
the operation including the Success/Failure status as well as a log of every
run. To do this add the flag `--send-logs-to-cloud`.

The operation logs will be made available in the web interface at
**Settings > Logs > External logs**.

## Mock Server

Files.com publishes a Files.com API server, which is useful for testing your use of the Files.com
SDKs and other direct integrations against the Files.com API in an integration test environment.

It is a Ruby app that operates as a minimal server for the purpose of testing basic network
operations and JSON encoding for your SDK or API client. It does not maintain state and it does not
deeply inspect your submissions for correctness.

Eventually we will add more features intended for integration testing, such as the ability to
intentionally provoke errors.

Download the server as a Docker image via [Docker Hub](https://hub.docker.com/r/filescom/files-mock-server).

The Source Code is also available on [GitHub](https://github.com/Files-com/files-mock-server).

A README is available on the GitHub link.

## Output Formatting

By default, the CLI App will output its data in table format.

You can select the output format by using the `--format` option. For example,
to specify that the output should be formatted in JSON format, use the option
`--format json`.

### Supported Formats

- table *(default)*
- table,interactive *(searchable and scrollable)*
- table,dark
- table,bright
- table,light
- table,markdown
- json *(human-readable)*
- json,raw *(compact)*
- csv

### Control Characters in Output

Values returned by the API, such as file names, can contain terminal
escape sequences and other control characters. Table output always
displays them as escaped text (for example `\x1b` for ESC) instead of
sending them to the terminal as controls. CSV and text output are escaped
the same way when written directly to a terminal; when redirected to a file
or piped to another program, values are written exactly as returned by the
API so exports stay lossless. Newline and tab characters keep their existing
formatting and are not escaped. JSON output always escapes control characters
as required by the JSON format; on a terminal the remaining DEL and C1
characters are shown as `\uXXXX` escapes as well, so terminal JSON stays valid
and decodes to the same values, while redirected JSON keeps its exact bytes.
Error messages and `--debug=STDOUT` logs are escaped the same way when written
to a terminal.

``` shell
files-cli users list --format=table,interactive
```

``` shell
files-cli folders list-for /path/to/folder --format csv
```

``` shell
files-cli folders create --path /path/to/folder/to/be/created --format=json,raw
```

``` shell
files-cli users list --format=table,dark
```

### Configuring the Default Format

You can configure a preferred format as the default for a profile.

```shell
files-cli config set --profile=interactive --format=table,interactive
```

### Customizing Output

You can customize the output by specifying the fields to include in the
response. Simply use the `--fields` option followed by a comma-separated list
of the desired field names.

```shell
files-cli folders list-for --fields path,type --format json
```

```json title="Example output"
[{
    "path": "document.docx",
    "type": "file"
},
{
    "path": "other",
    "type": "directory"
}]
```
