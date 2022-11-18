#!/usr/bin/env bash

if GOWORK=off goreleaser release --rm-dist; then
  true
else
  go install github.com/goreleaser/goreleaser@latest
  GOWORK=off $(go env GOPATH)/bin/goreleaser release --rm-dist || exit 1
fi
