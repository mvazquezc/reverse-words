.PHONY: build run test get-dependencies

build: test
	$(info Building Linux, Mac and Windows binaries)
	mkdir -p ./out/
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X 'main.gitCommit=$(shell git rev-parse --short HEAD)' -X 'main.buildTime=$(shell date +%Y-%m-%dT%H:%M:%SZ)'" -o ./out/reversewords-linux-amd64 main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-X 'main.gitCommit=$(shell git rev-parse --short HEAD)' -X 'main.buildTime=$(shell date +%Y-%m-%dT%H:%M:%SZ)'" -o ./out/reversewords-linux-arm64 main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-X 'main.gitCommit=$(shell git rev-parse --short HEAD)' -X 'main.buildTime=$(shell date +%Y-%m-%dT%H:%M:%SZ)'" -o ./out/reversewords-darwin-amd64 main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-X 'main.gitCommit=$(shell git rev-parse --short HEAD)' -X 'main.buildTime=$(shell date +%Y-%m-%dT%H:%M:%SZ)'" -o ./out/reversewords-darwin-arm64 main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-X 'main.gitCommit=$(shell git rev-parse --short HEAD)' -X 'main.buildTime=$(shell date +%Y-%m-%dT%H:%M:%SZ)'" -o ./out/reversewords-windows-amd64.exe main.go

run: get-dependencies
	go run main.go

test: get-dependencies
	go test

get-dependencies:
	$(info Downloading dependencies)
	go mod tidy
