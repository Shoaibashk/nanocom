# nanocom Makefile
# Cross-platform serial communication tool

# Project metadata
BINARY_NAME := nanocom
MODULE := github.com/shoaibashk/nanocom
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Go parameters
GO := go
GOFLAGS := -trimpath
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)
CGO_ENABLED := 0

# Directories
BUILD_DIR := build
DIST_DIR := dist

# Platforms for cross-compilation
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64

# Default target
.DEFAULT_GOAL := help

# ============================================================================
# Development
# ============================================================================

.PHONY: build
build: ## Build the binary for current platform
	CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .

.PHONY: run
run: ## Run the application
	$(GO) run .

.PHONY: run-dev
run-dev: ## Run with specific port (usage: make run-dev PORT=COM3)
	$(GO) run . -p $(PORT)

# ============================================================================
# Testing & Quality
# ============================================================================

.PHONY: test
test: ## Run all tests
	$(GO) test -v ./...

.PHONY: test-short
test-short: ## Run tests in short mode
	$(GO) test -short ./...

.PHONY: test-race
test-race: ## Run tests with race detector
	$(GO) test -race ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

.PHONY: bench
bench: ## Run benchmarks
	$(GO) test -bench=. -benchmem ./...

.PHONY: lint
lint: ## Run linter (requires golangci-lint)
	golangci-lint run ./...

.PHONY: fmt
fmt: ## Format code
	$(GO) fmt ./...
	gofmt -s -w .

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

.PHONY: check
check: fmt vet lint test ## Run all checks (fmt, vet, lint, test)

# ============================================================================
# Dependencies
# ============================================================================

.PHONY: deps
deps: ## Download dependencies
	$(GO) mod download

.PHONY: deps-update
deps-update: ## Update dependencies
	$(GO) get -u ./...
	$(GO) mod tidy

.PHONY: deps-tidy
deps-tidy: ## Tidy dependencies
	$(GO) mod tidy

.PHONY: deps-verify
deps-verify: ## Verify dependencies
	$(GO) mod verify

# ============================================================================
# Build & Release
# ============================================================================

.PHONY: build-all
build-all: clean ## Build for all platforms
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} GOARCH=$${platform#*/} \
		CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(BINARY_NAME)-$${platform%/*}-$${platform#*/}$(if $(findstring windows,$${platform%/*}),.exe,) .; \
		echo "Built: $(BUILD_DIR)/$(BINARY_NAME)-$${platform%/*}-$${platform#*/}"; \
	done

.PHONY: build-linux
build-linux: ## Build for Linux (amd64)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .

.PHONY: build-darwin
build-darwin: ## Build for macOS (amd64)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .

.PHONY: build-windows
build-windows: ## Build for Windows (amd64)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

.PHONY: release
release: ## Create release with goreleaser (requires goreleaser)
	goreleaser release --clean

.PHONY: release-snapshot
release-snapshot: ## Create snapshot release (no publish)
	goreleaser release --snapshot --clean

# ============================================================================
# Installation
# ============================================================================

.PHONY: install
install: ## Install binary to GOPATH/bin
	$(GO) install $(GOFLAGS) -ldflags "$(LDFLAGS)" .

.PHONY: uninstall
uninstall: ## Uninstall binary from GOPATH/bin
	rm -f $(shell $(GO) env GOPATH)/bin/$(BINARY_NAME)

# ============================================================================
# Cleanup
# ============================================================================

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe
	rm -f coverage.out coverage.html

.PHONY: clean-all
clean-all: clean ## Clean everything including dist
	rm -rf $(DIST_DIR)
	$(GO) clean -cache -testcache

# ============================================================================
# Development Tools
# ============================================================================

.PHONY: tools
tools: ## Install development tools
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install github.com/goreleaser/goreleaser@latest

.PHONY: list-ports
list-ports: build ## List available serial ports
	./$(BINARY_NAME) list

# ============================================================================
# Docker (optional)
# ============================================================================

.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(BINARY_NAME):$(VERSION) .

# ============================================================================
# Help
# ============================================================================

.PHONY: help
help: ## Show this help message
	@echo "nanocom - Cross-platform serial communication tool"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "Examples:"
	@echo "  make build          Build for current platform"
	@echo "  make test           Run all tests"
	@echo "  make build-all      Build for all platforms"
	@echo "  make run-dev PORT=COM3  Run with specific port"

.PHONY: version
version: ## Show version information
	@echo "Version:    $(VERSION)"
	@echo "Commit:     $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Version: $(shell $(GO) version)"
