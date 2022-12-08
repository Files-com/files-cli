#!/usr/bin/env bash

GONOPROXY=github.com/Files-com go get -u github.com/Files-com/files-sdk-go/v2

if ! command -v goreleaser &> /dev/null; then
   GOWORK=off curl -sfL https://goreleaser.com/static/run | DISTRIBUTION=pro bash -s -- release --rm-dist --skip-sign  || exit 1
else
  GOWORK=off goreleaser release --rm-dist --skip-sign
fi
