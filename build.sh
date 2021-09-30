#!/usr/bin/env bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "${DIR}" || exit 1

if [ -n "${DEVELOPMENT_BUILD}" ] || [ -n "${SNAPSHOT}" ]; then
  echo "go mod edit -replace github.com/Files-com/files-sdk-go/v2=../go"
  go mod edit -replace github.com/Files-com/files-sdk-go/v2=../go
else
  echo "go get -u github.com/Files-com/files-sdk-go/v2@master"
  go get -u github.com/Files-com/files-sdk-go/v2@master
fi

go mod tidy > /dev/null 2>&1
if [ -n "${DEVELOPMENT_BUILD}" ]
then
  sh test.sh  || exit 1
  cd ../go || exit 1
  sh test.sh || exit 1
  cd "${DIR}" || exit 1
fi

go mod tidy

version=$(ruby "../../next_version.rb" cli true)
version=$(echo "$version" | sed '/^[[:space:]]*$/d')
echo "$version" > "_VERSION"

if [ -n "${DEVELOPMENT_BUILD}" ] ||  [ -n "${SNAPSHOT}" ]; then
  if goreleaser build --rm-dist --snapshot; then
    true
  else
    curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh
    ./bin/goreleaser build --rm-dist --snapshot || exit 1
  fi
fi
ERROR_CODE=$?

if [ -n "${SNAPSHOT}" ]; then
    go mod edit -dropreplace github.com/Files-com/files-sdk-go/v2
fi

exit ${ERROR_CODE}
