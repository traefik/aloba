.PHONY: clean check test build build-crossbinary

GOFILES := $(shell go list -f '{{range $$index, $$element := .GoFiles}}{{$$.Dir}}/{{$$element}}{{"\n"}}{{end}}' ./... | grep -v '/vendor/')

default: clean check test build-crossbinary

test: clean
	go test -v -cover ./...

dependencies:
	dep ensure -v

clean:
	rm -f cover.out

build:
	go build

check:
	golangci-lint run

build-crossbinary:
	./_script/crossbinary
