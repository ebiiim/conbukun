VERSION?=$(shell git describe --tags --match "v*")
GO_BUILD_ARGS=-trimpath -ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: all
all: gen test build

.PHONY: gen
gen:
	go generate ./...

.PHONY: test
test: gen
	go test -v -cover -race ./...

.PHONY: build
build: gen conbukun

.PHONY: conbukun
conbukun: gen
	CGO_ENABLED=0 go build $(GO_BUILD_ARGS) -o bin/conbukun cmd/conbukun/main.go
