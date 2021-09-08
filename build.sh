#!/usr/bin/env bash
COMMITHASH=$(git describe --long --dirty --always --match '')

export CGO_ENABLED=0
export GOROOT_FINAL=/usr

export GOOS=linux
export GOARCH=amd64
go build -a -trimpath -asmflags '-s -w' -ldflags "-s -w -X github.com/aiocloud/stream.CommitHash=${COMMITHASH}" -o 'release/stream' || exit $?

cp -f example.json release
exit 0