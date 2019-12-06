SHELL := /bin/bash
OS = $(shell uname | tr A-Z a-z)

BUILD_DIR ?= build

GOLANGCI_VERSION = 1.21.0
GOTESTSUM_VERSION = 0.4.0

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the project
	go build -o build/ ./...

.PHONY: run
run: ## Run it
	@go run .

.PHONY: fmt
fmt: ## Format the source
	go fmt ./...

.PHONY: check
check: test lint ## Run the tests and linters

.PHONY: lint
lint: bin/golangci-lint ## Run linter
	$< run

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf $(notdir $<) $@
bin/golangci-lint-${GOLANGCI_VERSION}:
	@mkdir -p bin
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b ./bin/ v${GOLANGCI_VERSION}
	@mv bin/golangci-lint $@

.PHONY: test
TEST_PKGS ?= ./...
TEST_REPORT_NAME ?= results.xml
.PHONY: test
test: TEST_REPORT ?= main
test: TEST_FORMAT ?= short
test: SHELL = /bin/bash
test: bin/gotestsum ## Run tests
	@mkdir -p ${BUILD_DIR}/test_results/${TEST_REPORT}
	bin/gotestsum --no-summary=skipped --junitfile ${BUILD_DIR}/test_results/${TEST_REPORT}/${TEST_REPORT_NAME} --format ${TEST_FORMAT} --

.PHONY: test-all
test-all: ## Run all tests
	@${MAKE} GOARGS="${GOARGS} -run .\*" TEST_REPORT=all test

bin/gotestsum: bin/gotestsum-${GOTESTSUM_VERSION}
	@ln -sf $(notdir $<) $@
bin/gotestsum-${GOTESTSUM_VERSION}:
	@mkdir -p bin
	curl -L https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VERSION}/gotestsum_${GOTESTSUM_VERSION}_${OS}_amd64.tar.gz | tar -zOxf - gotestsum > $@ && chmod +x $@
