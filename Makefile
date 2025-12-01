# Anybakup Makefile
# This project is a backup utility written in Go

# Variables
BINARY_NAME := anybakup
BUILD_DIR := build
MAIN_PATH := .
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.0.0")
COMMIT_HASH := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commitHash=$(COMMIT_HASH) -X main.buildTime=$(BUILD_TIME) -s -w"

# Go related variables
GO := go
GOFMT := gofmt
GOLINT := golangci-lint
GOVULNCHECK := govulncheck

.PHONY: all build install clean test coverage fmt lint vet check vuln-check help

# Default target
all: check test build

# Build the project
build:
	@echo "Building $(BINARY_NAME) version $(VERSION)..."
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build completed successfully!"

# Install the project (go install)
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install $(LDFLAGS) $(MAIN_PATH)
	@echo "Installation completed successfully!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	@echo "Clean completed!"

# Run tests
test:
	@echo "Running tests..."
	$(GO) test -v ./...
	@echo "Tests completed!"

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	@echo "Code formatting completed!"

# Lint the code
lint:
	@echo "Linting code..."
	@which $(GOLINT) > /dev/null || (echo "golangci-lint not found. Install with 'go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest'" && exit 1)
	$(GOLINT) run
	@echo "Linting completed!"

# Vet the code
vet:
	@echo "Vetting code..."
	$(GO) vet ./...
	@echo "Vetting completed!"

# Run all checks
check: fmt vet lint

# Run vulnerability check
vuln-check:
	@echo "Checking for vulnerabilities..."
	@which $(GOVULNCHECK) > /dev/null || (echo "govulncheck not found. Install with 'go install golang.org/x/vuln/cmd/govulncheck@latest'" && exit 1)
	$(GOVULNCHECK) ./...
	@echo "Vulnerability check completed!"

# Run specific commands for different architectures (cross-compilation)
build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)

build-linux-arm64:
	GOOS=linux GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)

build-windows-amd64:
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

# Build for all platforms
build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64
	@echo "All builds completed!"

# Build cmd/gitcmd.go as a dynamic library
lib:
	@echo "Building gitcmd dynamic library..."
	$(GO) build -buildmode=c-shared -o $(BUILD_DIR)/libgitcmd.so ./cmd/gitcmd-lib
	@echo "gitcmd dynamic library build completed!"

# Build and run C test for the dynamic library
test-lib: lib
	@echo "Building C test program..."
	gcc -o $(BUILD_DIR)/test_gitcmd test_gitcmd.c -L$(BUILD_DIR) -lgitcmd -Wl,-rpath,$(abspath $(BUILD_DIR))
	@echo "Running C test program..."
	./$(BUILD_DIR)/test_gitcmd

# Build and run C test with real files
test-lib-real: lib
	@echo "Building C test program with real files..."
	gcc -o $(BUILD_DIR)/test_gitcmd_real test_gitcmd_real.c -L$(BUILD_DIR) -lgitcmd -Wl,-rpath,$(abspath $(BUILD_DIR))
	@echo "Running C test program with real files..."
	./$(BUILD_DIR)/test_gitcmd_real

# Clean test and library artifacts
clean-lib:
	@echo "Cleaning library and test artifacts..."
	rm -f $(BUILD_DIR)/libgitcmd.so $(BUILD_DIR)/libgitcmd.h $(BUILD_DIR)/test_gitcmd $(BUILD_DIR)/test_gitcmd_real
	@echo "Clean completed!"

# Initialize the project after cloning
init:
	@echo "Initializing project..."
	$(GO) mod tidy
	@echo "Project initialization completed!"

# Help target
help:
	@echo "Available targets:"
	@echo "  all            - Run check, test, and build (default)"
	@echo "  build          - Build the binary"
	@echo "  install        - Install the binary to GOPATH/bin"
	@echo "  clean          - Remove build artifacts"
	@echo "  test           - Run tests"
	@echo "  coverage       - Run tests with coverage report"
	@echo "  fmt            - Format code with gofmt"
	@echo "  lint           - Lint code with golangci-lint"
	@echo "  vet            - Vet code with go vet"
	@echo "  check          - Run fmt, vet, and lint"
	@echo "  vuln-check     - Check for vulnerabilities with govulncheck"
	@echo "  init           - Initialize project (go mod tidy)"
	@echo "  build-linux-amd64    - Build for Linux AMD64"
	@echo "  build-linux-arm64    - Build for Linux ARM64"
	@echo "  build-darwin-amd64   - Build for macOS AMD64"
	@echo "  build-darwin-arm64   - Build for macOS ARM64"
	@echo "  build-windows-amd64  - Build for Windows AMD64"
	@echo "  build-all      - Build for all platforms"
	@echo "  lib            - Build cmd/gitcmd.go as a dynamic library"
	@echo "  test-lib       - Build and run C test for the dynamic library"
	@echo "  test-lib-real  - Build and run C test with real files"
	@echo "  clean-lib      - Clean library and test artifacts"
	@echo "  help           - Show this help message"