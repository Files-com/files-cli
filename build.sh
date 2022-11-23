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
  if ! command -v goreleaser &> /dev/null; then
     curl -sfL https://goreleaser.com/static/run | DISTRIBUTION=pro bash -s -- build --rm-dist --snapshot || exit 1
  else
    goreleaser build --rm-dist --snapshot
  fi
fi
ERROR_CODE=$?

exit ${ERROR_CODE}
