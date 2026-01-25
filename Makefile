# Makefile for ImageToExcel Importer

BINARY_NAME=ImageToExcel
BUILD_DIR=build/bin

.PHONY: all build clean test coverage lint

all: build

# Build the application using Wails
build:
	wails build

# Build for Windows specifically (creates .exe)
build-windows:
	wails build -platform windows/amd64

# Build release with version injection (Usage: make build-release VERSION=v1.0.0)
build-release:
	wails build -platform windows/amd64 -ldflags "-s -w -X main.CurrentVersion=$(VERSION)"

# Run unit tests
test:
	go test ./... -v

# Run tests with coverage and open report
coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report saved to coverage.html"

# Clean build artifacts
clean:
	rm -rf build/
	rm -f coverage.out coverage.html
	rm -f *.log
	rm -f *_output_*.xlsx

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod tidy
	go install github.com/wailsapp/wails/v2/cmd/wails@latest
