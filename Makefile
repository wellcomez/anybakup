# Anybakup Makefile
# This project is a backup utility written in Go

# Variables
BINARY_NAME := anybakup
BUILD_DIR := build
MAIN_PATH := .

# Detect OS using Makefile built-in variables (no shell commands)
# OSTYPE is available in most Make implementations
# If not available, fall back to OS environment variable
ifdef OSTYPE
    # OSTYPE values: 'msys', 'cygwin', 'windows-gnu', 'linux-gnu', 'darwin*', etc.
    ifneq (,$(findstring $(OSTYPE),msys cygwin windows-gnu windows32))
        DETECTED_OS := Windows
    else ifneq (,$(findstring darwin,$(OSTYPE)))
        DETECTED_OS := Darwin
    else ifneq (,$(findstring linux,$(OSTYPE)))
        DETECTED_OS := Linux
    else
        DETECTED_OS := $(OSTYPE)
    endif
else
    # Fallback to OS environment variable
    ifdef OS
        ifeq ($(OS),Windows_NT)
            DETECTED_OS := Windows
        else ifneq (,$(findstring $(OS),Windows WINDOWS windows))
            DETECTED_OS := Windows
        else ifneq (,$(findstring $(OS),Darwin MAC darwin mac))
            DETECTED_OS := Darwin
        else ifneq (,$(findstring $(OS),Linux LINUX linux))
            DETECTED_OS := Linux
        else
            DETECTED_OS := Unknown
        endif
    else
        # Default fallback
        DETECTED_OS := Unknown
    endif
endif

# Cross-platform compatible version commands
ifeq ($(DETECTED_OS),Windows)
    # Try git commands first, fallback to defaults if not available
    VERSION := $(shell git describe --tags --always --dirty 2>nul || echo v0.0.0)
    COMMIT_HASH := $(shell git rev-parse --short HEAD 2>nul || echo unknown)
    # Try PowerShell for timestamp, fallback to simple format
    BUILD_TIME := $(shell powershell -Command "Get-Date -UFormat %%Y-%%m-%%dT%%H:%%M:%%SZ" 2>nul || echo "2025-01-01T00:00:00Z")
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

# Build the project - OS-specific compilation
build:
	@echo "Building $(BINARY_NAME) version $(VERSION) for $(DETECTED_OS)..."
	@echo "Detected OS: $(DETECTED_OS)"
	@echo "Building OS-specific library components..."
ifeq ($(DETECTED_OS),Windows)
	$(MAKE) build-windows-amd64
	$(MAKE) build-windows-lib
else ifeq ($(DETECTED_OS),Darwin)
	$(MAKE) build-darwin-amd64
	$(MAKE) build-darwin-lib
else ifeq ($(DETECTED_OS),Linux)
	$(MAKE) build-linux-amd64
	$(MAKE) build-unix-lib
else
	@echo "Building for unknown OS ($(DETECTED_OS)) - using generic build..."
	$(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)$(EXE_EXT) $(MAIN_PATH)
endif
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

test-darwin:
	$(GO) test -v ./...

test-linux:
	$(GO) test -v ./...

test-windows:
	@echo "Running tests for Windows..."
	@if [ "$(DETECTED_OS)" = "Windows" ]; then \
		if command -v powershell.exe >/dev/null 2>&1; then \
			powershell.exe -ExecutionPolicy Bypass -File "build-windows.ps1" -Task "test"; \
		else \
			echo "Error: PowerShell not found"; \
			exit 1; \
		fi; \
	else \
		go test -v ./...; \
	fi
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
	$(GO) vet ./...
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
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)

build-linux-arm64:
	GOOS=linux GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)

build-windows-amd64:
	@echo "Building Windows AMD64 binary..."
	powershell.exe -ExecutionPolicy Bypass -File "build-windows.ps1" -Task "exe"

# Build for all platforms
build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64 lib
	@echo "All builds completed!"

# Build all Windows components
build-windows-all:
	@if [ "$(DETECTED_OS)" = "Windows" ]; then \
		echo "Building all Windows components..."; \
		if command -v powershell.exe >/dev/null 2>&1; then \
			powershell.exe -ExecutionPolicy Bypass -File "build-windows.ps1" -Task "all"; \
		else \
			echo "Error: PowerShell not found"; \
			exit 1; \
		fi; \
	else \
		echo "Error: build-windows-all is only available on Windows"; \
		exit 1; \
	fi


build-windows-lib:
	powershell.exe -ExecutionPolicy Bypass -File "build-windows.ps1" -Task "lib"

# Build OS-specific dynamic libraries using existing targets
lib:
	@echo "Building dynamic library for $(DETECTED_OS)..."
ifeq ($(DETECTED_OS),Windows)
	$(MAKE) build-windows-lib
else ifeq ($(DETECTED_OS),Darwin)
	$(MAKE) build-darwin-lib
else ifeq ($(DETECTED_OS),Linux)
	$(MAKE) build-unix-lib
else
	@echo "Unknown OS: $(DETECTED_OS), cannot build appropriate library"
	@echo "Please use platform-specific targets: build-darwin-lib, build-unix-lib, build-windows-lib"
endif

# Legacy platform-specific targets (for backward compatibility)
build-darwin-lib:
	$(GO) build -buildmode=c-shared -o $(BUILD_DIR)/libgitcmd.dylib ./cmd/gitcmd-lib

build-unix-lib:
	$(GO) build -buildmode=c-shared -o $(BUILD_DIR)/libgitcmd.so ./cmd/gitcmd-lib


# Build and run C test for the dynamic library
# test-lib: lib
# 	@echo "Building C test program..."
# ifeq ($(DETECTED_OS),Windows)
# 	gcc -o $(BUILD_DIR)/test_gitcmd$(EXE_EXT) temp_test/test_gitcmd.c -L$(BUILD_DIR) -lgitcmd
# 	@echo "Running C test program..."
# 	./$(BUILD_DIR)/test_gitcmd$(EXE_EXT)
# else
# 	gcc -o $(BUILD_DIR)/test_gitcmd temp_test/test_gitcmd.c -L$(BUILD_DIR) -lgitcmd -Wl,-rpath,$(abspath $(BUILD_DIR))
# 	@echo "Running C test program..."
# 	./$(BUILD_DIR)/test_gitcmd
# endif

# # Build and run C test with real files
# test-lib-real: lib
# 	@echo "Building C test program with real files..."
# ifeq ($(DETECTED_OS),Windows)
# 	gcc -o $(BUILD_DIR)/test_gitcmd_real$(EXE_EXT) temp_test/test_gitcmd_real.c -L$(BUILD_DIR) -lgitcmd
# 	@echo "Running C test program with real files..."
# 	./$(BUILD_DIR)/test_gitcmd_real$(EXE_EXT)
# else
# 	gcc -o $(BUILD_DIR)/test_gitcmd_real temp_test/test_gitcmd_real.c -L$(BUILD_DIR) -lgitcmd -Wl,-rpath,$(abspath $(BUILD_DIR))
# 	@echo "Running C test program with real files..."
# 	./$(BUILD_DIR)/test_gitcmd_real
# endif

# Clean test and library artifacts
clean-lib:
	@echo "Cleaning library and test artifacts..."
	@if [ "$(DETECTED_OS)" = "Windows" ]; then \
		$(RM) $(BUILD_DIR)/gitcmd.dll $(BUILD_DIR)/gitcmd.h $(BUILD_DIR)/test_gitcmd$(EXE_EXT) $(BUILD_DIR)/test_gitcmd_real$(EXE_EXT) 2>/dev/null || true; \
	else \
		$(RM) $(BUILD_DIR)/libgitcmd.so $(BUILD_DIR)/libgitcmd.dylib $(BUILD_DIR)/gitcmd.dll $(BUILD_DIR)/gitcmd.h $(BUILD_DIR)/test_gitcmd $(BUILD_DIR)/test_gitcmd_real 2>/dev/null || true; \
	fi
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
	@echo "  build          - Build the binary with OS-specific components"
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
	@echo "  build-windows-all    - Build all Windows components (lib, exe, test)"
	@echo "  build-all      - Build for all platforms"
	@echo "  lib            - Build OS-appropriate dynamic library (Windows/Darwin/Linux)"
	@echo "  test-windows   - Run Windows tests"
	@echo "  test-lib       - Build and run C test for the dynamic library"
	@echo "  test-lib-real  - Build and run C test with real files"
	@echo "  clean-lib      - Clean library and test artifacts"
	@echo "  clean-all      - Clean all artifacts including test files"
	@echo "  help           - Show this help message"
	@echo ""
	@echo "OS Detection:"
	@echo "  Detected OS: $(DETECTED_OS)"
	@echo "  Build target automatically compiles OS-specific components:"
	@echo "    - Windows: Uses PowerShell for library build"
	@echo "    - macOS (Darwin): Builds .dylib dynamic library"
	@echo "    - Linux: Builds .so shared library"
	@echo ""
	@echo "Windows build notes:"
	@echo "  Windows builds use PowerShell script build-windows.ps1"
	@echo "  Requires MSYS2 with MinGW-w64 for CGO support"
	@echo "  PowerShell execution policy may need to be adjusted"
