#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "${DIR}" || exit 1
go mod edit -replace github.com/Files-com/files-sdk-go=../go
go build -o files main.go
GOOS=windows GOARCH=386 go build -o files.exe main.go
go mod edit -dropreplace github.com/Files-com/files-sdk-go
cp files.exe ../.. && cd ../..
