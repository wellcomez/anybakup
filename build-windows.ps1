# PowerShell script for all Windows build tasks

param(
    [ValidateSet("lib", "exe", "test", "all")]
    [string]$Task = "lib"
)

Write-Host "Windows Build Script - Task: $Task" -ForegroundColor Green
Write-Host "======================================" -ForegroundColor Green

# Function to detect and setup GCC
function Setup-GCC {
    param(
        [ref]$GCC_PATH,
        [ref]$GCC_BIN_DIR
    )

    # Check for GCC in D: drive MSYS2
    if (Test-Path "D:\msys64\ucrt64\bin\gcc.exe") {
        $GCC_PATH.Value = "D:\msys64\ucrt64\bin\gcc.exe"
        $GCC_BIN_DIR.Value = "D:\msys64\ucrt64\bin"
        return $true
    }

    # Check for GCC in C: drive MSYS2
    if (Test-Path "C:\msys64\ucrt64\bin\gcc.exe") {
        $GCC_PATH.Value = "C:\msys64\ucrt64\bin\gcc.exe"
        $GCC_BIN_DIR.Value = "C:\msys64\ucrt64\bin"
        return $true
    }

    # Check in PATH
    if (Get-Command gcc -ErrorAction SilentlyContinue) {
        $GCC_PATH.Value = "gcc"
        $GCC_BIN_DIR.Value = ""
        return $true
    }

    return $false
}

# Function to build library
function Build-Library {
    Write-Host "Building gitcmd dynamic library..." -ForegroundColor Yellow
    Write-Host "Checking prerequisites..." -ForegroundColor Yellow

    $GCC_PATH = ""
    $GCC_BIN_DIR = ""

    if (-not (Setup-GCC -GCC_PATH ([ref]$GCC_PATH) -GCC_BIN_DIR ([ref]$GCC_BIN_DIR))) {
        Write-Host "Error: GCC not found. Please install MSYS2 with MinGW-w64" -ForegroundColor Red
        Write-Host "Download from: https://www.msys2.org/" -ForegroundColor Yellow
        Write-Host "Then run: pacman -S mingw-w64-ucrt-x86_64-gcc" -ForegroundColor Yellow
        exit 1
    }

    Write-Host "Found GCC at: $GCC_PATH" -ForegroundColor Green
    & $GCC_PATH --version

    # Create build directory
    if (-not (Test-Path "build")) {
        New-Item -ItemType Directory -Path "build" | Out-Null
    }

    Write-Host "Setting environment variables..." -ForegroundColor Yellow

    # Set environment variables
    $env:PATH = "$GCC_BIN_DIR;$env:PATH"
    $env:CGO_ENABLED = "1"
    $env:CC = $GCC_PATH
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"

    Write-Host "Building DLL..." -ForegroundColor Yellow

    # Build the library
    go build -buildmode=c-shared -ldflags="-s -w" -o "build\gitcmd.dll" ".\cmd\gitcmd-lib"

    if (Test-Path "build\gitcmd.dll") {
        Write-Host "Successfully built gitcmd.dll" -ForegroundColor Green
        if (Test-Path "build\gitcmd.h") {
            Write-Host "Generated header file: gitcmd.h" -ForegroundColor Green
        }
        return $true
    } else {
        Write-Host "Error: Failed to build gitcmd.dll" -ForegroundColor Red
        return $false
    }
}

# Function to build executable
function Build-Executable {
    Write-Host "Building anybakup Windows executable..." -ForegroundColor Yellow

    # Get version info
    $VERSION = "v0.0.0"
    $COMMIT_HASH = "unknown"
    $BUILD_TIME = "2025-01-01T00:00:00Z"

    try {
        $VERSION = & git describe --tags --always --dirty 2>$null
        if (-not $VERSION) { $VERSION = "v0.0.0" }
    } catch { }

    try {
        $COMMIT_HASH = & git rev-parse --short HEAD 2>$null
        if (-not $COMMIT_HASH) { $COMMIT_HASH = "unknown" }
    } catch { }

    try {
        $BUILD_TIME = Get-Date -UFormat "%Y-%m-%dT%H:%M:%SZ"
    } catch { }

    Write-Host "Version: $VERSION" -ForegroundColor Cyan
    Write-Host "Commit: $COMMIT_HASH" -ForegroundColor Cyan
    Write-Host "Build Time: $BUILD_TIME" -ForegroundColor Cyan

    # Setup GCC if available (optional for exe)
    $GCC_PATH = ""
    $GCC_BIN_DIR = ""
    $HAS_GCC = Setup-GCC -GCC_PATH ([ref]$GCC_PATH) -GCC_BIN_DIR ([ref]$GCC_BIN_DIR)

    if ($HAS_GCC) {
        Write-Host "Using GCC: $GCC_PATH" -ForegroundColor Green
        $env:PATH = "$GCC_BIN_DIR;$env:PATH"
        $env:CGO_ENABLED = "1"
        $env:CC = $GCC_PATH
    } else {
        Write-Host "Warning: GCC not found, building without CGO support" -ForegroundColor Yellow
        $env:CGO_ENABLED = "0"
    }

    $env:GOOS = "windows"
    $env:GOARCH = "amd64"

    Write-Host "Building executable..." -ForegroundColor Yellow

    # Build the executable
    $LDFLAGS = "-X main.version=$VERSION -X main.commitHash=$COMMIT_HASH -X main.buildTime=$BUILD_TIME -s -w"
    go build -ldflags $LDFLAGS -o "build\anybakup-windows-amd64.exe" "."

    if (Test-Path "build\anybakup-windows-amd64.exe") {
        Write-Host "Successfully built anybakup-windows-amd64.exe" -ForegroundColor Green
        $SIZE = (Get-Item "build\anybakup-windows-amd64.exe").Length / 1MB
        Write-Host "File size: $([math]::Round($SIZE, 2)) MB" -ForegroundColor Cyan
        return $true
    } else {
        Write-Host "Error: Failed to build executable" -ForegroundColor Red
        return $false
    }
}

# Function to run tests
function Run-Tests {
    Write-Host "Running Windows tests..." -ForegroundColor Yellow

    # Setup GCC if available
    $GCC_PATH = ""
    $GCC_BIN_DIR = ""
    $HAS_GCC = Setup-GCC -GCC_PATH ([ref]$GCC_PATH) -GCC_BIN_DIR ([ref]$GCC_BIN_DIR)

    if ($HAS_GCC) {
        Write-Host "Using GCC for CGO tests: $GCC_PATH" -ForegroundColor Green
        $env:PATH = "$GCC_BIN_DIR;$env:PATH"
        $env:CGO_ENABLED = "1"
        $env:CC = $GCC_PATH
    } else {
        Write-Host "Warning: GCC not found, running tests without CGO support" -ForegroundColor Yellow
        $env:CGO_ENABLED = "0"
    }

    $env:GOOS = "windows"
    $env:GOARCH = "amd64"

    Write-Host "Running go test..." -ForegroundColor Yellow

    # Run tests
    go test -v ./...

    if ($LASTEXITCODE -eq 0) {
        Write-Host "All tests passed!" -ForegroundColor Green
        return $true
    } else {
        Write-Host "Some tests failed!" -ForegroundColor Red
        return $false
    }
}

# Main execution
$SUCCESS = $false

switch ($Task) {
    "lib" {
        $SUCCESS = Build-Library
    }
    "exe" {
        $SUCCESS = Build-Executable
    }
    "test" {
        $SUCCESS = Run-Tests
    }
    "all" {
        $LIB_SUCCESS = Build-Library
        $EXE_SUCCESS = Build-Executable
        $TEST_SUCCESS = Run-Tests
        $SUCCESS = $LIB_SUCCESS -and $EXE_SUCCESS -and $TEST_SUCCESS
    }
}

if ($SUCCESS) {
    Write-Host "`nBuild task '$Task' completed successfully!" -ForegroundColor Green
    exit 0
} else {
    Write-Host "`nBuild task '$Task' failed!" -ForegroundColor Red
    exit 1
}