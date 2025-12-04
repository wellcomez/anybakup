# Anybakup Makefile
# This project is a backup utility written in Go

# Variables
BINARY_NAME := anybakup
BUILD_DIR := build
MAIN_PATH := .

# Detect OS
UNAME_S := $(shell uname -s)
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
else
    DETECTED_OS := $(UNAME_S)
endif

# Cross-platform compatible version commands
ifeq ($(DETECTED_OS),Windows)
    VERSION := $(shell git describe --tags --always --dirty 2>nul || echo v0.0.0)
    COMMIT_HASH := $(shell git rev-parse --short HEAD 2>nul || echo unknown)
    BUILD_TIME := $(shell powershell -Command "Get-Date -UFormat %%Y-%%m-%%dT%%H:%%M:%%SZ")
    RM := del /Q
    RMDIR := rmdir /S /Q
    MKDIR := mkdir
    MOVE := move
    EXE_EXT := .exe
    LIB_EXT := .dll
else
    VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "v0.0.0")
    COMMIT_HASH := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
    RM := rm -rf
    RMDIR := rm -rf
    MKDIR := mkdir -p
    MOVE := mv
    EXE_EXT :=
    LIB_EXT := .so
endif

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
build: lib
	@echo "Building $(BINARY_NAME) version $(VERSION)..."
	@$(MKDIR) temp_test
	@$(MOVE) test_*.c temp_test/ 2>/dev/null || true
	@$(MOVE) test_*.py temp_test/ 2>/dev/null || true
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)$(EXE_EXT) $(MAIN_PATH)
	@$(MOVE) temp_test/* . 2>/dev/null || true
	@$(RMDIR) temp_test/ 2>/dev/null || true
	@echo "Build completed successfully!"

# Install the project (go install)
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install $(LDFLAGS) $(MAIN_PATH)
	@echo "Installation completed successfully!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@$(RMDIR) $(BUILD_DIR) 2>/dev/null || true
	@echo "Clean completed!"

# Clean all including test artifacts
clean-all: clean clean-lib
	@echo "Cleaning all artifacts including test files..."

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
	@if command -v $(GOLINT) >/dev/null 2>&1; then \
		$(GOLINT) run; \
	else \
		echo "golangci-lint not found. Install with 'go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest'"; \
		exit 1; \
	fi
	@echo "Linting completed!"

# Vet the code
vet:
	@echo "Vetting code..."
	@$(MKDIR) temp_test
	@$(MOVE) test_*.c temp_test/ 2>/dev/null || true
	@$(MOVE) test_*.py temp_test/ 2>/dev/null || true
	$(GO) vet ./...
	@$(MOVE) temp_test/* . 2>/dev/null || true
	@$(RMDIR) temp_test/ 2>/dev/null || true
	@echo "Vetting completed!"

# Run all checks
check: fmt vet lint

# Run vulnerability check
vuln-check:
	@echo "Checking for vulnerabilities..."
	@if command -v $(GOVULNCHECK) >/dev/null 2>&1; then \
		$(GOVULNCHECK) ./...; \
	else \
		echo "govulncheck not found. Install with 'go install golang.org/x/vuln/cmd/govulncheck@latest'"; \
		exit 1; \
	fi
	@echo "Vulnerability check completed!"

# Run specific commands for different architectures (cross-compilation)
build-linux-amd64:
	@$(MKDIR) temp_test
	@$(MOVE) test_*.c temp_test/ 2>/dev/null || true
	@$(MOVE) test_*.py temp_test/ 2>/dev/null || true
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@$(MOVE) temp_test/* . 2>/dev/null || true
	@$(RMDIR) temp_test/ 2>/dev/null || true

build-linux-arm64:
	@$(MKDIR) temp_test
	@$(MOVE) test_*.c temp_test/ 2>/dev/null || true
	@$(MOVE) test_*.py temp_test/ 2>/dev/null || true
	GOOS=linux GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	@$(MOVE) temp_test/* . 2>/dev/null || true
	@$(RMDIR) temp_test/ 2>/dev/null || true

build-darwin-amd64:
	@$(MKDIR) temp_test
	@$(MOVE) test_*.c temp_test/ 2>/dev/null || true
	@$(MOVE) test_*.py temp_test/ 2>/dev/null || true
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@$(MOVE) temp_test/* . 2>/dev/null || true
	@$(RMDIR) temp_test/ 2>/dev/null || true

build-darwin-arm64:
	@$(MKDIR) temp_test
	@$(MOVE) test_*.c temp_test/ 2>/dev/null || true
	@$(MOVE) test_*.py temp_test/ 2>/dev/null || true
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@$(MOVE) temp_test/* . 2>/dev/null || true
	@$(RMDIR) temp_test/ 2>/dev/null || true

build-windows-amd64:
	@$(MKDIR) temp_test
	@$(MOVE) test_*.c temp_test/ 2>/dev/null || true
	@$(MOVE) test_*.py temp_test/ 2>/dev/null || true
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@$(MOVE) temp_test/* . 2>/dev/null || true
	@$(RMDIR) temp_test/ 2>/dev/null || true

# Build for all platforms
build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64 lib
	@echo "All builds completed!"

# Build cmd/gitcmd.go as a dynamic library
lib:
	@echo "Building gitcmd dynamic library..."
ifeq ($(DETECTED_OS),Windows)
	$(GO) build -buildmode=c-shared -o $(BUILD_DIR)/gitcmd.dll ./cmd/gitcmd-lib
else ifeq ($(DETECTED_OS),Darwin)
	$(GO) build -buildmode=c-shared -o $(BUILD_DIR)/libgitcmd.dylib ./cmd/gitcmd-lib
else
	$(GO) build -buildmode=c-shared -o $(BUILD_DIR)/libgitcmd.so ./cmd/gitcmd-lib
endif
	@echo "gitcmd dynamic library build completed!"

# Build and run C test for the dynamic library
test-lib: lib
	@echo "Building C test program..."
ifeq ($(DETECTED_OS),Windows)
	gcc -o $(BUILD_DIR)/test_gitcmd$(EXE_EXT) temp_test/test_gitcmd.c -L$(BUILD_DIR) -lgitcmd
	@echo "Running C test program..."
	./$(BUILD_DIR)/test_gitcmd$(EXE_EXT)
else
	gcc -o $(BUILD_DIR)/test_gitcmd temp_test/test_gitcmd.c -L$(BUILD_DIR) -lgitcmd -Wl,-rpath,$(abspath $(BUILD_DIR))
	@echo "Running C test program..."
	./$(BUILD_DIR)/test_gitcmd
endif

# Build and run C test with real files
test-lib-real: lib
	@echo "Building C test program with real files..."
ifeq ($(DETECTED_OS),Windows)
	gcc -o $(BUILD_DIR)/test_gitcmd_real$(EXE_EXT) temp_test/test_gitcmd_real.c -L$(BUILD_DIR) -lgitcmd
	@echo "Running C test program with real files..."
	./$(BUILD_DIR)/test_gitcmd_real$(EXE_EXT)
else
	gcc -o $(BUILD_DIR)/test_gitcmd_real temp_test/test_gitcmd_real.c -L$(BUILD_DIR) -lgitcmd -Wl,-rpath,$(abspath $(BUILD_DIR))
	@echo "Running C test program with real files..."
	./$(BUILD_DIR)/test_gitcmd_real
endif

# Clean test and library artifacts
clean-lib:
	@echo "Cleaning library and test artifacts..."
	@$(RMDIR) temp_test 2>/dev/null || true
ifeq ($(DETECTED_OS),Windows)
	@$(RM) $(BUILD_DIR)\gitcmd.dll $(BUILD_DIR)\gitcmd.h $(BUILD_DIR)\test_gitcmd$(EXE_EXT) $(BUILD_DIR)\test_gitcmd_real$(EXE_EXT) 2>nul
else
	@$(RM) $(BUILD_DIR)/libgitcmd.so $(BUILD_DIR)/libgitcmd.dylib $(BUILD_DIR)/gitcmd.dll $(BUILD_DIR)/gitcmd.h $(BUILD_DIR)/test_gitcmd $(BUILD_DIR)/test_gitcmd_real 2>/dev/null || true
endif
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
	@echo "  clean-all      - Clean all artifacts including test files"
	@echo "  help           - Show this help message"