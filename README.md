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

#### Setting by ENV 

``` shell
export FILES_API_KEY="XXXX-XXXX..."
```

#### Set Via a Flag

```shell 
files-cli folders list-for --api-key "XXXX-XXXX..."
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
files-cli folders list-for other
```

### Download a File/Folder

```shell
files-cli download [remote-path] [local-path]
```

### Upload a File/Folder

```shell
files-cli upload [source-path] [remote-path]
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

