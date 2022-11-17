#!/usr/bin/env bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "${DIR}" || exit 1

if [ -n "${DEVELOPMENT_BUILD}" ] || [ -n "${SNAPSHOT}" ]; then
  echo "GOWORK=${DIR}/go.work"
  export GOWORK="${DIR}/go.work"
else
  echo "go get -u github.com/Files-com/files-sdk-go/v2@master"
  GONOPROXY=github.com/Files-com go get -u github.com/Files-com/files-sdk-go/v2@master
fi

go install golang.org/x/tools/cmd/goimports@latest
goimports -w .
gofmt -s -w .
