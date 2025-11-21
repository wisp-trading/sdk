.PHONY: generate build test clean help

# Default target
.DEFAULT_GOAL := help

## generate: Generate code from templates (version.go from version.yaml)
generate:
	@echo "Generating version.go from version.yaml..."
	@go run scripts/generate_version.go

## build: Build all packages
build: generate
	@echo "Building all packages..."
	@go build ./...

## test: Run tests
test:
	@echo "Running tests..."
	@go test ./...

## clean: Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@go clean ./...
	@rm -rf dist/

## help: Show this help message
help:
	@echo "Available targets:"
	@sed -n 's/^##//p' Makefile | column -t -s ':' | sed -e 's/^/ /'
