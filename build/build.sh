#!/bin/sh

export GOPATH=""

export CGO_ENABLED=0

set -e
set -x

cd ..
rm -rf release

GOPROXY=https://goproxy.cn GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.DEV=0' -X 'main.VERSION=1.0.0' -X 'main.BUILD_TIME=$(date)' -X 'main.GO_VERSION=$(go version)'" -o release/amd64/m3u8_cli
GOPROXY=https://goproxy.cn GOOS=linux GOARCH=arm64 go build -ldflags "-X 'main.DEV=0' -X 'main.VERSION=1.0.0' -X 'main.BUILD_TIME=$(date)' -X 'main.GO_VERSION=$(go version)'" -o release/arm64/m3u8_cli
GOPROXY=https://goproxy.cn GOOS=linux GOARCH=arm go build -ldflags "-X 'main.DEV=0' -X 'main.VERSION=1.0.0' -X 'main.BUILD_TIME=$(date)' -X 'main.GO_VERSION=$(go version)'" -o release/arm/m3u8_cli
GOPROXY=https://goproxy.cn GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.DEV=0' -X 'main.VERSION=1.0.0' -X 'main.BUILD_TIME=$(date)' -X 'main.GO_VERSION=$(go version)'" -o release/darwin/m3u8_cli
GOPROXY=https://goproxy.cn GOOS=windows GOARCH=amd64 go build -ldflags "-X 'main.DEV=0' -X 'main.VERSION=1.0.0' -X 'main.BUILD_TIME=$(date)' -X 'main.GO_VERSION=$(go version)'" -o release/windows/m3u8_cli.exe
