PROJECT = $(shell basename $(PWD))

GOPATH ?= $(shell go env GOPATH)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

ORG = github.com/AidanFogarty/pitwall

GO_FILES = $(shell find . -type f -name '*.go' -not -path "./third_party/*")

.DEFAULT_GOAL := help

## build: Build the binary.
build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/ .

## test: Run all the tests in the project
test:
	go test ./... -v

## coverage: Run test coverage on the project
coverage:
	go test ./...  -coverpkg=./... -coverprofile ./coverage.out
	go tool cover -func ./coverage.out

## vet: Examine Go source code and reports suspicious constructs
vet:
	go vet ./...

## fmt: Formats all Go code in the project
fmt:
	go fmt ./...

## tidy: tidy the go.mod file
tidy:
	go mod tidy

# Thanks to: https://github.com/azer/go-makefile-example
.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECT)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
