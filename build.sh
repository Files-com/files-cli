#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "${DIR}" || exit 1
go mod tidy
go get -d -v
go get -u github.com/mitchellh/gox
~/go/bin/gox -os="linux darwin windows" -arch="amd64" -output="../../cli_dist/files_cli_{{.OS}}_{{.Arch}}"
go mod edit -dropreplace github.com/Files-com/files-sdk-go
