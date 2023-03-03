#!/usr/bin/env bash

pushd cmd/sandstorm-rcon-client || exit 1
set -x
go mod tidy
env GOOS=linux   GOARCH=amd64 go build -ldflags "-s -w"
env GOOS=windows GOARCH=amd64 go build -ldflags "-s -w"
ls -alh sandstorm-rcon-client*
sha512sum sandstorm-rcon-client* > sandstorm-rcon-client-sha512sums.txt
set +x

popd || exit 1
