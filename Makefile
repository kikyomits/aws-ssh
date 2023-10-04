#!/usr/bin/env make

SERVICE_NAME=aws-ssh
VERSION_FILE:=./VERSION
BRANCH:=$(shell git branch 2> /dev/null | sed -e '/^[^*]/d' -e 's/* \(.*\)/\1/' | tr '/' '-')
PKG_PATH=github.com/kikyonmits


## clean up dependencies
mod-tidy:
	go mod tidy
.PHONY: mod-tidy

## download dependencies
mod-dl:
	go mod download
.PHONY: mod-dl

## install golangci-lint
install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2
.PHONY: install-lint

## run linters on code
lint: install-lint code-gen
	golangci-lint run ./... --out-format checkstyle > lint.out
.PHONY: lint

## run linters in fix mode on code
lint-fix: install-lint
	golangci-lint run --fix ./...
.PHONY: lint

code-gen:
	go generate ./...
.PHONY: code-gen

## run unit tests
test: code-gen
	go test ./... -covermode=set -coverprofile=coverage.out
.PHONY: test

# test without running code-gen, to check mocks are up to date
test-ci:
	go test ./... -covermode=set -coverprofile=coverage.out -json > tests.out
.PHONY: test-ci

## build application
compile:
	go build ./...
.PHONY: compile