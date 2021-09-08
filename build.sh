#!/usr/bin/env bash
export CGO_ENABLED=0
export GOROOT_FINAL=/usr

export GOOS=linux
export GOARCH=amd64
go build -a -trimpath -asmflags '-s -w' -ldflags '-s -w' -o 'release/stream' || exit $?

cp -f example.json release
exit 0