#!/usr/bin/env bash

if goreleaser build --rm-dist; then
  true
else
  curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh
  ./bin/goreleaser build --rm-dist
fi
