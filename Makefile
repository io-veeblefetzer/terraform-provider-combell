# Copyright (c) Veeblefetzer
# SPDX-License-Identifier: MPL-2.0

# ==============================================================================
# Variables
# ==============================================================================

# Provider information
HOSTNAME := registry.terraform.io
NAMESPACE := veeblefetzer
NAME := combell
BINARY := terraform-provider-$(NAME)
VERSION ?= 0.1.0

# Go settings
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOBIN ?= $(shell go env GOPATH)/bin

# Build settings
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

# Local install directory for Terraform development
# Terraform >= 0.14 uses this path for local provider installs
OS_ARCH := $(GOOS)_$(GOARCH)
INSTALL_DIR := ~/.terraform.d/plugins/$(HOSTNAME)/$(NAMESPACE)/$(NAME)/$(VERSION)/$(OS_ARCH)

# Test settings
TEST_TIMEOUT := 120s
ACCTEST_TIMEOUT := 10m
ACCTEST_PARALLELISM ?= 4

# ==============================================================================
# Default target
# ==============================================================================

.PHONY: default
default: build

# ==============================================================================
# Build targets
# ==============================================================================

.PHONY: build
build: ## Build the provider binary
	go build $(LDFLAGS) -o $(BINARY)

.PHONY: install
install: build ## Install the provider locally for Terraform development
	@mkdir -p $(INSTALL_DIR)
	@cp $(BINARY) $(INSTALL_DIR)/$(BINARY)
	@echo "Installed $(BINARY) to $(INSTALL_DIR)"

.PHONY: release
release: ## Build release binaries for all supported platforms
	@mkdir -p dist
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)_darwin_amd64
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)_darwin_arm64
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)_linux_amd64
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY)_linux_arm64
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY)_windows_amd64.exe
	@echo "Release binaries built in dist/"

# ==============================================================================
# Test targets
# ==============================================================================

.PHONY: test
test: ## Run unit tests
	go test ./... -v -timeout $(TEST_TIMEOUT)

.PHONY: testacc
testacc: ## Run acceptance tests (requires TF_ACC=1 and API credentials)
	@echo "Running acceptance tests..."
	@echo "Ensure COMBELL_API_KEY and COMBELL_API_SECRET are set"
	TF_ACC=1 go test ./... -v -timeout $(ACCTEST_TIMEOUT) -parallel $(ACCTEST_PARALLELISM)

.PHONY: testcover
testcover: ## Run tests with coverage report
	go test ./... -v -timeout $(TEST_TIMEOUT) -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: testclient
testclient: ## Run client tests only
	go test ./internal/client/... -v -timeout $(TEST_TIMEOUT)

# ==============================================================================
# Code quality targets
# ==============================================================================

.PHONY: fmt
fmt: ## Format Go source code
	gofmt -s -w .
	@echo "Code formatted"

.PHONY: fmtcheck
fmtcheck: ## Check if code is formatted
	@files=$$(gofmt -l .); \
	if [ -n "$$files" ]; then \
		echo "The following files are not formatted:"; \
		echo "$$files"; \
		exit 1; \
	fi
	@echo "All files are properly formatted"

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint (if installed)
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install with: brew install golangci-lint"; \
		exit 1; \
	fi

.PHONY: check
check: fmtcheck vet ## Run all code quality checks

# ==============================================================================
# Dependency management
# ==============================================================================

.PHONY: deps
deps: ## Download and tidy dependencies
	go mod download
	go mod tidy
	@echo "Dependencies updated"

.PHONY: verify
verify: ## Verify dependencies
	go mod verify
	@echo "Dependencies verified"

# ==============================================================================
# Documentation
# ==============================================================================

.PHONY: docs
docs: ## Generate provider documentation
	@if command -v tfplugindocs >/dev/null 2>&1; then \
		tfplugindocs generate; \
	else \
		echo "tfplugindocs not installed. Install with: go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest"; \
		exit 1; \
	fi

.PHONY: docscheck
docscheck: ## Validate provider documentation
	@if command -v tfplugindocs >/dev/null 2>&1; then \
		tfplugindocs validate; \
	else \
		echo "tfplugindocs not installed. Install with: go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest"; \
		exit 1; \
	fi

# ==============================================================================
# Utility targets
# ==============================================================================

.PHONY: clean
clean: ## Remove build artifacts
	rm -f $(BINARY)
	rm -rf dist/
	rm -f coverage.out coverage.html
	@echo "Cleaned build artifacts"

.PHONY: uninstall
uninstall: ## Remove locally installed provider
	rm -rf ~/.terraform.d/plugins/$(HOSTNAME)/$(NAMESPACE)/$(NAME)
	@echo "Uninstalled provider from local Terraform plugins"

.PHONY: tools
tools: ## Install development tools
	go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
	@echo "Tools installed"

.PHONY: help
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Variables:"
	@echo "  VERSION          Provider version (default: $(VERSION))"
	@echo "  GOOS             Target OS (default: $(GOOS))"
	@echo "  GOARCH           Target architecture (default: $(GOARCH))"
	@echo ""
	@echo "Examples:"
	@echo "  make build                    # Build provider"
	@echo "  make install                  # Install locally for testing"
	@echo "  make test                     # Run unit tests"
	@echo "  make testacc                  # Run acceptance tests"
	@echo "  make VERSION=1.0.0 release    # Build release with specific version"