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

go mod tidy

if [ -n "${DEVELOPMENT_BUILD}" ] || [ -n "${SNAPSHOT}" ]; then
  go mod edit -replace github.com/Files-com/files-sdk-go=../go
else
  go get -u github.com/Files-com/files-sdk-go
fi

project_root="$(git rev-parse --show-toplevel)"
next_version_path="${project_root}/next_version.rb"
if [[ -f "${next_version_path}" ]] ; then
  sed -i.bu "s/_VERSION/$(ruby "${next_version_path}" cli)/g" main.go
fi

if [[ -f "main.go.bu" ]] ; then
  rm main.go.bu
fi

return_code=0

if [ -n "${DEVELOPMENT_BUILD}" ] ||  [ -n "${SNAPSHOT}" ]; then
  if goreleaser build --rm-dist --snapshot ; then
    return_code=$((return_code + $?))
    true
  else
    curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh
    ./bin/goreleaser build --rm-dist --snapshot
    return_code=$((return_code + $?))
  fi
fi

if [ -n "${SNAPSHOT}" ]; then
    go mod edit -dropreplace github.com/Files-com/files-sdk-go
fi

exit return_code
