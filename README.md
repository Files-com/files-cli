# Files.com Command Line App

The Files.com CLI App provides convenient access to the Files.com API.

## Installation

Download the latest release for Windows, macOS, or Linux [here](https://github.com/Files-com/files-cli/releases)

### Homebrew 

```shell
brew tap Files-com/homebrew-tap
brew install files-cli
```

## Documentation

### Setting API Key

#### Setting by ENV (Linux/macOS)

``` shell
export FILES_API_KEY="XXXX-XXXX..."
```

#### Setting by ENV (Windows)

``` shell
set FILES_API_KEY="XXXX-XXXX..."
```

#### Set Via a Flag

```shell 
files-cli folders list-for --api-key "XXXX-XXXX..."
```

### Password Authentication

To log in to the CLI App with your username and password, you must first configure the CLI App with information about your Files.com account:

``` shell
files-cli config set --subdomain MYCOMPANY --username MYUSERNAME
```

When prompted, enter your password.

Once you've specified your subdomain information and username, you do not need to specify it again for subsequent uses of the CLI App. For later uses, you can login using the following command, which will prompt you for your password:

``` shell
files-cli login
```

If your account requires Two-Factor Authentication, you will be prompted for the second factor after you submit your password. Once you are logged in, subsequent uses of the CLI App will perform those actions using your credentials and permissions until you log out.

### Logging Out

Your login session will expire automatically after a period of time. The CLI App will expire your session after 6 hours or your session will expire based on the settings of your authentication system, whichever is sooner.

To log out of your session manually, use:

``` shell
files-cli logout
```

### List files

*Return root folder listing*

```shell 
files-cli folders list-for --fields path,type

[{
    "path": "document.docx",
    "type": "file"
},
{
    "path": "other",
    "type": "directory"
}]
```

*List a Folder*

```shell 
files-cli folders list-for /path/to/folder
```

### Download a File/Folder

To download a file, use the command:

```shell
files-cli download /remote/path/to/file.txt /local/path/to/file.txt
```

or

```shell
files-cli download /remote/path/to/file.txt /local/path/to/folder/
```

To download a folder, use the command:

```shell
files-cli download /remote/path/to/folder/ /local/path/to/folder/
```

### Upload a File/Folder

To upload a file, use the command:

```shell
files-cli upload /local/path/to/file.txt /remote/path/to/file.txt
```

or

```shell
files-cli upload /local/path/to/file.txt /remote/path/to/folder/
```

To upload a folder, use the command:

```shell
files-cli upload /local/path/to/folder/ /remote/path/to/folder/
```

### Creating folders

To create folders, use this command:

``` shell
files-cli folders create --path “/path/to/folder/to/be/created”
```

### Sending Operation/Run Logs to the Cloud

If you are running scripted operations, you can have the CLI send a report of the operation including the Success/Failure status as well as a log of every run. To do this add the flag `--send-logs-to-cloud`.

The operation logs will be made available in the web interface at **Settings > Logs > External logs**.

### Syncing Files

To facilitate file-syncing workflows, the `--sync` flag can be used with the upload or download command to specify that only new files be transferred.

Here is a "push" (upload) example for syncing files from a local Documents folder to a Files.com folder of the same name:

``` shell
files-cli upload Documents Documents --sync --send-logs-to-cloud
```

Here is a "pull" (download) example for syncing files to a local Documents folder from a Files.com folder of the same name:

``` shell
files-cli download Documents Documents --sync --send-logs-to-cloud
```

### Administrator actions

If you have administrator privileges for your Files.com account, you can use the CLI App to perform administrator actions.

For example, you can create a user account with this command:

``` shell
files-cli users create --username amy --password "S0meRea11yLongP@ssw0rd" --authentication-method "password" --name "Amy Anybody" --company "Amy’s Company Name" --notes "Some notes about Amy." --user-root "/users/amy"
```

You can also configure various items, such as Folder Settings, using the CLI App.

For example, you can configure [automatic new user folders](https://www.files.com/docs/topics/folder-settings#automatically-create-new-user-folders-here-when-users-are-created) using the following command:

``` shell
files-cli behaviors create --path "/path/to/folder" --behavior "create_user_folders" --value '{ "permission":"full", "additional_permission":"bundle", "existing_users":false, "group_id":1, "new_folder_name":"username", "subfolders":[]}'
```

### Formatting the output

By default, the CLI App will output its data in table format.

You can configure the output format by using the `--format` option. For example, to specify that the output should be formatted in JSON format, use the option `--format json`.

Available output formats are:

- table
- table-dark
- table-bright
- json
- csv

Here are some examples:

``` shell
files-cli folders list-for /path/to/folder --format csv
```

``` shell
files-cli folders create --path “/path/to/folder/to/be/created” --format json
```

``` shell
files-cli users list --format table-dark
```

### Command Help

```shell
files-cli [command] --help
```

## Development

To build for testing it assumes the go package is in parallel directory. 

```shell
DEVELOPMENT_BUILD ./build.sh
```

This will build both the compressed release version and an uncompressed executable for the 3 platforms.

### Docker

```shell
docker build . --tag files-cli:latest
docker run --workdir /app --volume ${PWD}:/app -it files-cli 
```

#### Run CLI

```shell
 docker run --workdir /app --volume ${PWD}:/app -it files-cli bash -c "go run main.go"
```

#### Link local Go SDK

```shell
go mod edit -replace github.com/Files-com/files-sdk-go/v2=../files-sdk-go
docker run --workdir /app --volume ${PWD}:/app --volume ${HOME}/go/src/github.com/Files-com/files-sdk-go:/files-sdk-go -it files-cli 
```

