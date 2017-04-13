# Define parameters
BINARY=runp
SHELL := /bin/bash
GOPACKAGES = $(shell go list ./... | grep -v vendor)
ROOTDIR = $(pwd)

.PHONY: build env install test linux

default: build

build: runp/main.go
	go build -v -o ./build/${BINARY} runp/main.go

env:
	export GOPATH=${GOPATH}

install:
	go install  ./...

test:
	go test -race -cover ${GOPACKAGES}

clean:
	rm -rf build

linux: runp/main.go
	GOOS=linux GOARCH=amd64 go build -o ./build/linux/${BINARY} runp/main.go
