set CGO_ENABLED=0
cd ..

set GOPROXY=https://goproxy.cn
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-X 'main.DEV=0' -X 'main.VERSION=1.0.0' -X 'main.BUILD_TIME=$(date)' -X 'main.GO_VERSION=$(go version)'" -o release/amd64/m3u8_cli

set GOPROXY=https://goproxy.cn
set GOOS=linux
set GOARCH=arm64
go build -ldflags "-X 'main.DEV=0' -X 'main.VERSION=1.0.0' -X 'main.BUILD_TIME=$(date)' -X 'main.GO_VERSION=$(go version)'" -o release/arm64/m3u8_cli

set GOPROXY=https://goproxy.cn
set GOOS=linux
set GOARCH=arm
go build -ldflags "-X 'main.DEV=0' -X 'main.VERSION=1.0.0' -X 'main.BUILD_TIME=$(date)' -X 'main.GO_VERSION=$(go version)'" -o release/arm/m3u8_cli

set GOPROXY=https://goproxy.cn
set GOOS=darwin
set GOARCH=amd64
go build -ldflags "-X 'main.DEV=0' -X 'main.VERSION=1.0.0' -X 'main.BUILD_TIME=$(date)' -X 'main.GO_VERSION=$(go version)'" -o release/darwin/m3u8_cli

set GOPROXY=https://goproxy.cn
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-X 'main.DEV=0' -X 'main.VERSION=1.0.0' -X 'main.BUILD_TIME=$(date)' -X 'main.GO_VERSION=$(go version)'" -o release/windows/m3u8_cli.exe
