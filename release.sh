#!/usr/bin/env bash

if goreleaser release --rm-dist; then
  true
else
  go install github.com/goreleaser/goreleaser@latest
  $(go env GOPATH)/bin/goreleaser release --rm-dist || exit 1
fi
