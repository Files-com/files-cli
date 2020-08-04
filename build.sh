#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "${DIR}" || exit 1
go mod tidy > /dev/null 2>&1
go mod edit -replace github.com/Files-com/files-sdk-go=../go
go mod tidy
go get -d -v
go get -u github.com/mitchellh/gox
gox -os="linux darwin windows" -arch="amd64" -output="../../cli_dist/{{.OS}}/files_cli_{{.OS}}_{{.Arch}}"
go mod edit -dropreplace github.com/Files-com/files-sdk-go

cd ../../
cd cli_dist || exit
mv windows/files_cli_windows_amd64.exe windows/files_cli.exe
cd windows && zip ../files_cli_windows_amd64.zip files_cli.exe && cd ../
rm -rf windows

mv darwin/files_cli_darwin_amd64 darwin/files_cli
cd darwin && zip ../files_cli_darwin_amd64.zip files_cli && cd ../
rm -rf darwin

mv linux/files_cli_linux_amd64 linux/files_cli
cd linux && zip ../files_cli_linux_amd64.zip files_cli && cd ../
rm -rf linux
