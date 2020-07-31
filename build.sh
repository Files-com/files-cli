#!/usr/bin/env bash

cd generated/cli || exit 1
go mod edit -replace github.com/Files-com/files-sdk-go=../go
GOOS=windows GOARCH=386 go build -o files.exe main.go
go mod edit -dropreplace github.com/Files-com/files-sdk-go
cp files.exe ../.. && cd ../..
