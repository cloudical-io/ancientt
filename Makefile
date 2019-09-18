PROJECTNAME ?= acntt
DESCRIPTION ?= acntt - A tool to automate network testing tools, like iperf3, in dynamic environments such as Kubernetes and more to come dynamic environments.
MAINTAINER  ?= Alexander Trost <galexrt@googlemail.com>
HOMEPAGE    ?= https://github.com/cloudical-io/acntt

GO      := go
GOFMT   := gofmt
PREFIX  ?= $(shell pwd)
VERSION ?= $(shell cat VERSION)

pkgs = $(shell go list ./... | grep -v /vendor/ | grep -v /test/)

DOCKER_IMAGE_NAME ?= acntt
DOCKER_IMAGE_TAG  ?= $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))

all: format style vet test build

docker:
	@echo ">> building docker image"
	@docker build -t "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" .

format:
	go fmt $(pkgs)

style:
	@echo ">> checking code style"
	@! $(GOFMT) -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

test:
	@$(GO) test $(pkgs)

test-short:
	@echo ">> running short tests"
	@$(GO) test -short $(pkgs)

vet:
	@echo ">> vetting code"
	@$(GO) vet $(pkgs)

docs: pkg/config/config.go
	$(GO) run ./cmd/docgen/ api pkg/config/*.go > docs/config-structure.md

.PHONY: all build crossbuild docker format promu style tarball test vet
