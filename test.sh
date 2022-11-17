#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd "${DIR}" || exit 1
go test ./... -buildvcs=false
ERROR_CODE=$?

exit ${ERROR_CODE}
