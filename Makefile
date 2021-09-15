PROJECTNAME ?= ancientt
DESCRIPTION ?= ancientt - A tool to automate network testing tools, like iperf3, in dynamic environments such as Kubernetes and more to come dynamic environments.
HOMEPAGE    ?= https://github.com/cloudical-io/ancientt

GO_SUPPORTED_VERSIONS ?= 1.15|1.16|1.17

DOCKER  := docker
GO      := go
GOFMT   := gofmt
PREFIX  ?= $(shell pwd)
VERSION ?= $(shell cat VERSION)

GO_TEST_FLAGS ?= -race

pkgs = $(shell go list ./... | grep -v /vendor/ | grep -v /test/)

DOCKER_IMAGE_NAME ?= ancientt
DOCKER_IMAGE_TAG  ?= $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))

go.check:
ifneq ($(shell $(GO) version | grep -q -E '\bgo($(GO_SUPPORTED_VERSIONS))\b' && echo 0 || echo 1), 0)
	$(error unsupported go version. Please make install one of the following supported version: '$(GO_SUPPORTED_VERSIONS)')
endif

all: format style vet test build

ancientt: go.check
	$(GO) build -o ancientt $(PREFIX)/cmd/ancientt/

build: ancientt

docker:
	@echo ">> building docker image"
	$(DOCKER) build -t "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" .

format: go.check
	$(GO) fmt $(pkgs)

style: go.check
	@echo ">> checking code style"
	$(GOFMT) -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

test: go.check
	@echo ">> running tests"
	$(GO) test $(GO_TEST_FLAGS) $(pkgs)

test-short: go.check
	@echo ">> running short tests"
	$(GO) test -short $(pkgs)

vet: go.check
	@echo ">> vetting code"
	$(GO) vet $(pkgs)

docs: pkg/config/config.go
	@echo ">> generating docs"
	$(GO) run ./cmd/docgen/ api pkg/config/*.go > docs/config-structure.md

.PHONY: all build docker docs format go.check style test test-short vet
