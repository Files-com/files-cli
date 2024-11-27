# Files.com Command Line App

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

### Installation

Download the latest release for Windows, macOS, or Linux from the [CLI App Releases](https://github.com/Files-com/files-cli/releases) page.

On that page, you'll need to pick your exact operating system to download the correct version.

No installation is necessary. The app is a self contained app which can be stored anywhere on your computer.

We recommend placing the app binary into one of the folders listed in your `%PATH%` (Windows) or `$PATH` (Linux and Mac).

Here are specific instructions, grouped by OS:

#### Windows

**Download**

- AMD/Intel 64-Bit Processors [amd64](https://github.com/Files-com/files-cli/releases/latest/download/files-cli_windows_64bit.zip) *(Most Common)*
- ARM Processors [arm64](https://github.com/Files-com/files-cli/releases/latest/download/files-cli_windows_arm64.zip) *(Less Common)*

Download the Zip archive (`*.zip`), extract the files from the archive, and place the `files-cli.exe` binary file into any directory listed in your `%PATH%` environment variable.

#### Mac using Homebrew

For systems with Homebrew available, use the following commands:

```shell
brew tap Files-com/homebrew-tap
brew install files-cli
```

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

## Authentication

### Authenticate with an API Key

Authenticating with an API key is the recommended authentication method for most scenarios, and is
the method used in the examples on this site.

To use the API or SDKs with an API Key, first generate an API key from the [web
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

```shellsession title="Example Request"
files-cli --api-key=YOUR_API_KEY folders list-for ''
## After the key has been provided once it will be written to the files-cli configuration file.
## You do not need to include the same API key for future commands.
```

Don't forget to replace the placeholder, `YOUR_API_KEY`, with your actual API key.

### Authenticate with a Session

You can also authenticate to the REST API or SDKs by creating a user session using the username and
password of an active user. If the user is an administrator, the session will have full access to
the entire API. Sessions created from regular user accounts will only be able to access files that
user can access, and no access will be granted to site administration functions.

API sessions use the exact same session timeout settings as web interface sessions. When an API
session times out, simply create a new session and resume where you left off. This process is not
automatically handled by SDKs because we do not want to store password information in memory without
your explicit consent.

#### Logging In

To log in to the CLI App with your username and password, you must first configure the CLI App with information about your Files.com account.

Once you've specified your subdomain information and username, you do not need to specify it again for subsequent uses of the CLI App.

You will be prompted for the password when a command is run that requires authentication.  Once the password is entered, subsequent calls will not require a password, unless the session is no longer valid (ie, expired).

```shellsession title="Example Request"
files-cli config set --subdomain SUBDOMAIN --username motor
files-cli folders list-for ''
> password: vroom
```

#### Using a Session

Once a the password has been verified by the first run command, all subsequent commands can be run and the current verified session will be used.

```shellsession title="Example Request"
files-cli folders list-for ''
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

Setting the base URL for the API is required if your site is configured to disable global acceleration.
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

The CLI does not currently support sorting. Please see SDK documentation for other languages for
sorting.

### Filtering

Filters apply selection criteria to the underlying query that returns the results. Filters can be
applied individually to select resource fields and/or in a combination with each other.

```filter_by``` -  Client side filtering: eg. field_name=*.jpg) |

An example of cli command line filter argument:

```--filter-by="field_name=*.jpg"```

```shellsession title="Example Filter Request" hasDataFormatSelector
files-cli users list \
  --filter-by="not_site_admin=true" \
  --api-key=my-key
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

## Examples

The CLI contains commands for uploading, downloading, syncing, listing, copying, moving, and
deleting files and folders. Each CLI command supports a `--help` flag to display additional
information about the command and its options.

### Upload

To upload a file, use the command:

```shell
files-cli upload local/path/to/file.txt remote/path/to/file.txt
```

or

```shell
files-cli upload local/path/to/file.txt remote/path/to/folder/
```

To upload a folder, use the command:

```shell
files-cli upload local/path/to/folder/ remote/path/to/folder/
```

#### Creating Folders

To create a folder, use this command:

```shell
files-cli folders create path/to/folder/to/be/created
```

### Download

To download a file, use the command:

```shell
files-cli download remote/path/to/file.txt local/path/to/file.txt
```

or

```shell
files-cli download remote/path/to/file.txt local/path/to/folder/
```

To download a folder, use the command:

```shell
files-cli download remote/path/to/folder/ local/path/to/folder/
```

### Sync

To facilitate file-syncing workflows, the `--sync` flag can be used with the
upload or download command to specify that only new files be transferred.

Here is a "push" (upload) example for syncing files from a local Documents
folder to a Files.com folder of the same name:

```shell
files-cli upload Documents Documents --sync --send-logs-to-cloud
```

Here is a "pull" (download) example for syncing files to a local Documents
folder from a Files.com folder of the same name:

```shell
files-cli download Documents Documents --sync --send-logs-to-cloud
```

### List

To list files and folders in a remote directory, use the command:

```shell
files-cli folders list-for remote/path/to/folder
```

<div></div>

The response content and format can be customized. For example:

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

### Copy

To copy a file, use the command:

```shell
files-cli files copy --path remote/path/to/file.txt --destination remote/path/to/file-copy.txt
```

This also works for folders:

```shell
files-cli files copy --path remote/path/to/folder/ --destination remote/path/to/folder-copy/
```

### Move

To move a file, use the command:

```shell
files-cli files move --path remote/path/to/file.txt --destination remote/path/to/file-moved.txt
```

This also works for folders:

```shell
files-cli files move --path remote/path/to/folder/ --destination remote/path/to/folder-moved/
```

### Delete

To delete a file, use the command:

```shell
files-cli files delete remote/path/to/file.txt
```

This also works for folders:

```shell
files-cli files delete remote/path/to/folder/
```

In case the folder is not empty, you can use the `--recursive` flag:

```shell
files-cli files delete remote/path/to/folder/ --recursive
```

### Sending Operation/Run Logs to the Cloud

If you are running scripted operations, you can have the CLI send a report of
the operation including the Success/Failure status as well as a log of every
run. To do this add the flag `--send-logs-to-cloud`.

The operation logs will be made available in the web interface at
**Settings > Logs > External logs**.

### Administrator Actions

If you have administrator privileges for your Files.com account, you can use the
CLI App to perform administrator actions.

For example, you can create a user account with this command:

```shell
files-cli users create \
  --username amy \
  --authentication-method password \
  --password "S0meRea11yLongP@ssw0rd" \
  --name "Amy Anybody" \
  --company "Amy’s Company Name" \
  --user-root "/users/amy" \
  --notes "Some notes about Amy."
```

You can also configure various items, such as Folder Settings, using the CLI App.

For example, you can configure
[automatic new user folders](https://www.files.com/docs/files-and-folders/folder-settings#automatically-create-new-user-folders-here-when-users-are-created)
using the following command:

```shell
files-cli behaviors create \
  --path "/path/to/folder" \
  --behavior "create_user_folders" \
  --value '{ "permission":"full", "additional_permission":"bundle", "existing_users":false, "group_id":1, "new_folder_name":"username", "subfolders":[]}'
```

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
