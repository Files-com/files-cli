#!/usr/bin/env bash

cd generated/cli || exit 1
go mod edit -replace github.com/Files-com/files-sdk-go=../go
go mod tidy
go get -d -v
go fmt ./...
go test ./...
go mod edit -dropreplace github.com/Files-com/files-sdk-go
