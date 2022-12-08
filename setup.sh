#!/usr/bin/env bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "${DIR}" || exit 1

if [ -n "${DEVELOPMENT_BUILD}" ] || [ -n "${SNAPSHOT}" ]; then
  export GOWORK="${DIR}/go.work"
  echo $GOWORK
else
  echo "Updating to latest files-sdk-go"
  GONOPROXY=github.com/Files-com go get -u github.com/Files-com/files-sdk-go/v2@master
fi

go install golang.org/x/tools/cmd/goimports@latest
goimports -w .
gofmt -s -w .

if [ -n "${DEVELOPMENT_BUILD}" ] || [ -n "${SNAPSHOT}" ]; then
  echo "Skipping 'go mod tidy' for SNAPSHOT build"
else
  go mod tidy
  go mod download
fi
