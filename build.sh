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

project_root="$(git rev-parse --show-toplevel)"
next_version_path="${project_root}/next_version.rb"
if [[ -f "${next_version_path}" ]] ; then
  sed -i.bu "s/_VERSION/$(ruby "${next_version_path}" cli)/g" main.go
fi

go get -u github.com/mitchellh/gox
gox -os="linux darwin windows" -arch="amd64" -output="../../cli_dist/{{.OS}}/files-cli_{{.OS}}_{{.Arch}}"
go mod edit -dropreplace github.com/Files-com/files-sdk-go

if [[ -f "main.go.bu" ]] ; then
  rm main.go
  cp main.go.bu main.go
fi

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
