#!/usr/bin/env bash

set -xe
pushd cmd/sandstorm-rcon-client

go mod tidy

env GOOS=linux   GOARCH=amd64 go build -ldflags "-s -w"
env GOOS=windows GOARCH=amd64 go build -ldflags "-s -w"

sha512sum sandstorm-rcon-client{,.exe} > sandstorm-rcon-client-sha512sums.txt

ls -alh sandstorm-rcon-client*

popd
