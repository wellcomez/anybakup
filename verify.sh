#!/bin/bash
set -e

# Build the tool
go build -o anybakup main.go

# Clean up previous run
rm -rf test_repo
rm -rf ~/.config/anybakup
rm -f test.txt

# 1. Init
./anybakup init test_repo
if [ ! -d "test_repo" ]; then
    echo "Error: test_repo not created"
    exit 1
fi
if [ ! -f "$HOME/.config/anybakup/config.yaml" ]; then
    echo "Error: config file not created"
    exit 1
fi

# 2. Create dummy file
echo "hello world" > test.txt

# 3. Add file
./anybakup add test.txt
if [ ! -f "test_repo/test.txt" ]; then
    echo "Error: file not added to repo"
    exit 1
fi

# 4. Status
./anybakup status test.txt | grep "is tracked"

# 5. Rm file
./anybakup rm test.txt
if [ -f "test_repo/test.txt" ]; then
    echo "Error: file not removed from repo"
    exit 1
fi

# 6. Status again
./anybakup status test.txt | grep "is NOT tracked"

echo "Verification passed!"
