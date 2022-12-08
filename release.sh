#!/usr/bin/env bash

if ! command -v goreleaser &> /dev/null; then
  curl -sfL https://goreleaser.com/static/run | GOWORK=off DISTRIBUTION=pro bash -s -- release --rm-dist --skip-sign || exit 1
else
  GOWORK=off goreleaser release --rm-dist --skip-sign
fi
