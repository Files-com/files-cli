#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "${DIR}" || exit 1
# Changes files_sdk_go latest to a version number
go mod tidy > /dev/null 2>&1
go mod edit -replace github.com/Files-com/files-sdk-go=../go
go mod tidy
go get -d -v
go fmt ./...
go test ./...
go mod edit -dropreplace github.com/Files-com/files-sdk-go
