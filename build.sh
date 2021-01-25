#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "${DIR}" || exit 1
go mod tidy > /dev/null 2>&1
if [ -n "${DEVELOPMENT_BUILD}" ]
then
  sh test.sh
  cd ../go || exit
  sh test.sh
  cd "${DIR}" || exit 1
fi

go mod edit -replace github.com/Files-com/files-sdk-go=../go
go mod tidy
go get -d -v
go get -u github.com/mitchellh/gox
gox -os="linux darwin windows" -arch="amd64" -output="../../cli_dist/{{.OS}}/files-cli_{{.OS}}_{{.Arch}}"
go mod edit -dropreplace github.com/Files-com/files-sdk-go

cd ../../
cd cli_dist || exit
mv windows/files-cli_windows_amd64.exe windows/files-cli.exe
cd windows && zip ../files-cli_windows_amd64.zip files-cli.exe && cd ../
[[ -z "${DEVELOPMENT_BUILD}" ]] && rm -rf windows

mv darwin/files-cli_darwin_amd64 darwin/files-cli
cd darwin && zip ../files-cli_darwin_amd64.zip files-cli && cd ../
[[ -z "${DEVELOPMENT_BUILD}" ]] && rm -rf darwin

mv linux/files-cli_linux_amd64 linux/files-cli
cd linux && zip ../files-cli_linux_amd64.zip files-cli && cd ../
[[ -z "${DEVELOPMENT_BUILD}" ]] && rm -rf linux
