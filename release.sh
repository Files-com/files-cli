#!/usr/bin/env bash

version=$(ruby "../../next_version.rb" cli true)
version=$(echo "$version" | sed '/^[[:space:]]*$/d')

echo "${version}" > "_VERSION"

if GOWORK=off goreleaser release --rm-dist; then
  true
else
  go install github.com/goreleaser/goreleaser@latest
  GOWORK=off $(go env GOPATH)/bin/goreleaser release --rm-dist || exit 1
fi
