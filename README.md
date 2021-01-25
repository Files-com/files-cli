# Files.com Command Line App

The Files.com CLI App provides convenient access to the Files.com API.

## Installation

Download latest release for Windows, MacOS, or Linux [here](https://github.com/Files-com/files-cli/releases)

## Documentation

### Setting API Key

#### Setting by ENV 

``` sh
export FILES_API_KEY="XXXX-XXXX..."
```

#### Set Via a Flag

```sh 
files-cli folders list-for --api-key "XXXX-XXXX..."
```

### List files

*Return root folder listing*

```sh 
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

```sh 
files-cli folders list-for other
```

### Download a File/Folder

```sh
files-cli download [remote-path] [local-path]
```

### Upload a File/Folder

```sh
files-cli upload [source-path] [remote-path]
```

### Command Help

```sh
files-cli [command] --help
```

## Development

To build for testing it assumes the go package is in parallel directory. 

```sh
DEVELOPMENT_BUILD ./build.sh
```

This will build both the compressed release version and an uncompressed executable for the 3 platforms.