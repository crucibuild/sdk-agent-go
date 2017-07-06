#!/bin/bash

ROOT_DIR="$(dirname "$0" | xargs readlink -f)"

check()
{
    RESULT="$($@ 2>&1)"

    if [ -n "$RESULT" ]
    then
        printf "\e[91mLight$RESULT\e[39m"
        exit 1
    fi
}

cd $ROOT_DIR
# Redirect stdout as we only want errors and linter output.
go get -u github.com/alecthomas/gometalinter 1>/dev/null
gometalinter --install --update 1>/dev/null

check gometalinter --enable misspell ./...
check gofmt -d .
check golint ./...
check go test -v -race ./...
check go vet ./...
