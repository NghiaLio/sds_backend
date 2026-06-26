@echo off
SETLOCAL EnableDelayedExpansion

echo ===================================================
echo Building Mobile CRUD Backend (SQLite Pure-Go)...
echo ===================================================

:: Ensure dependencies are tidy
echo Cleaning dependencies...
go mod tidy

:: Create bin directory
if not exist "bin" mkdir bin

:: 1. Build for Windows (Current OS)
echo Building Windows binary...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o bin/backend-windows-amd64.exe main.go
if !ERRORLEVEL! EQU 0 (
    echo [OK] Windows build succeeded: bin/backend-windows-amd64.exe
) else (
    echo [ERROR] Windows build failed.
)

:: 2. Build for Linux (useful for servers, docker, or android termux)
echo Building Linux binary...
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w" -o bin/backend-linux-amd64 main.go
if !ERRORLEVEL! EQU 0 (
    echo [OK] Linux build succeeded: bin/backend-linux-amd64
) else (
    echo [ERROR] Linux build failed.
)

:: 3. Build for macOS (Intel)
echo Building macOS AMD64 binary...
set GOOS=darwin
set GOARCH=amd64
go build -ldflags="-s -w" -o bin/backend-darwin-amd64 main.go
if !ERRORLEVEL! EQU 0 (
    echo [OK] macOS AMD64 build succeeded: bin/backend-darwin-amd64
) else (
    echo [ERROR] macOS AMD64 build failed.
)

:: 4. Build for macOS (Apple Silicon M1/M2/M3/M4)
echo Building macOS ARM64 binary...
set GOOS=darwin
set GOARCH=arm64
go build -ldflags="-s -w" -o bin/backend-darwin-arm64 main.go
if !ERRORLEVEL! EQU 0 (
    echo [OK] macOS ARM64 build succeeded: bin/backend-darwin-arm64
) else (
    echo [ERROR] macOS ARM64 build failed.
)

echo.
echo ===================================================
echo Build complete. Binaries are stored in 'bin/' directory.
echo ===================================================
