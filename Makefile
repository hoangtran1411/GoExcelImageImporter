# Makefile for ImageToExcel Importer

BINARY_NAME=tool_chen_anh
BUILD_DIR=build/bin

.PHONY: all build build-windows build-release dev test coverage clean fmt lint deps

all: build

# Build the application using Wails v3
build:
	wails3 build

# Build for Windows specifically (creates .exe)
build-windows:
	wails3 task windows:build

# Build release with version injection (Usage: make build-release VERSION=v1.0.0)
build-release:
	wails3 task windows:build VERSION=$(VERSION)

# Run in development mode
dev:
	-wails3 dev -config ./build/config.yml

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
	rm -rf $(BUILD_DIR)/
	rm -f coverage.out coverage.html
	rm -f *.log
	rm -f *_output_*.xlsx
	rm -f *.syso

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod tidy
	go install github.com/wailsapp/wails/v3/cmd/wails3@latest
