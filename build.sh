#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

echo "==================================================="
echo "Building Mobile CRUD Backend (SQLite Pure-Go)..."
echo "==================================================="

# Ensure dependencies are tidy
echo "Cleaning dependencies..."
go mod tidy

# Create bin directory if it doesn't exist
mkdir -p bin

# 1. Build for Windows
echo "Building Windows binary..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/backend-windows-amd64.exe main.go
echo "[OK] Windows build succeeded: bin/backend-windows-amd64.exe"

# 2. Build for Linux
echo "Building Linux binary..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/backend-linux-amd64 main.go
echo "[OK] Linux build succeeded: bin/backend-linux-amd64"

# 3. Build for macOS (Intel)
echo "Building macOS AMD64 binary..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/backend-darwin-amd64 main.go
echo "[OK] macOS AMD64 build succeeded: bin/backend-darwin-amd64"

# 4. Build for macOS (Apple Silicon M1/M2/M3/M4)
echo "Building macOS ARM64 binary..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/backend-darwin-arm64.app main.go
# Also build standard binary name for terminal running
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/backend-darwin-arm64 main.go
echo "[OK] macOS ARM64 build succeeded: bin/backend-darwin-arm64"

echo ""
echo "==================================================="
echo "Build complete. Binaries are stored in 'bin/' directory."
echo "==================================================="
