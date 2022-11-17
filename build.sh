#!/usr/bin/env bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "${DIR}" || exit 1

version=$(ruby "../../next_version.rb" cli true)
version=$(echo "$version" | sed '/^[[:space:]]*$/d')
if [ -n "${DEVELOPMENT_BUILD}" ] || [ -n "${SNAPSHOT}" ]; then
  echo "${version}-unreleased" > "_VERSION"
else
  echo "${version}" > "_VERSION"
fi

if [ -n "${DEVELOPMENT_BUILD}" ] ||  [ -n "${SNAPSHOT}" ]; then
  if goreleaser build --rm-dist --snapshot; then
    true
  else
    go install github.com/goreleaser/goreleaser@latest
    $(go env GOPATH)/bin/goreleaser build --rm-dist --snapshot || exit 1
  fi
fi
ERROR_CODE=$?

exit ${ERROR_CODE}
