#!/usr/bin/env bash

if goreleaser release --rm-dist; then
  true
else
  curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh
  ./bin/goreleaser release --rm-dist || exit 1
fi
