.PHONY: clean fmt dependencies check test build

GOFILES := $(shell go list -f '{{range $$index, $$element := .GoFiles}}{{$$.Dir}}/{{$$element}}{{"\n"}}{{end}}' ./... | grep -v '/vendor/')

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse --short HEAD)
VERSION := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))
BUILD_DATE := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

VERSION_PACKAGE=github.com/containous/aloba/meta

default: clean check test build

test: clean
	go test -v -cover ./...

dependencies:
	dep ensure -v

clean:
	rm -f cover.out

build: clean
	@echo Version: $(VERSION) $(BUILD_DATE)
	go build -v -ldflags '-X "${VERSION_PACKAGE}.version=${VERSION}" -X "${VERSION_PACKAGE}.commit=${SHA}" -X "${VERSION_PACKAGE}.date=${BUILD_DATE}"'

check:
	golangci-lint run

fmt:
	@gofmt -s -l -w $(GOFILES)
