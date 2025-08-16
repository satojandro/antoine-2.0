const MakefileContent = `# Antoine CLI Makefile

# Variables
BINARY_NAME=antoine
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags="-X 'main.Version=${VERSION}' -X 'main.BuildTime=${BUILD_TIME}'"

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

# Build the application
.PHONY: build
build:
	@echo "Building Antoine CLI..."
	@cd $(GOBASE) && go build ${LDFLAGS} -o $(GOBIN)/$(BINARY_NAME) ./cmd/antoine

# Build for multiple platforms
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(GOBIN)/dist
	@cd $(GOBASE) && GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o $(GOBIN)/dist/$(BINARY_NAME)-linux-amd64 ./cmd/antoine
	@cd $(GOBASE) && GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o $(GOBIN)/dist/$(BINARY_NAME)-darwin-amd64 ./cmd/antoine
	@cd $(GOBASE) && GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o $(GOBIN)/dist/$(BINARY_NAME)-darwin-arm64 ./cmd/antoine
	@cd $(GOBASE) && GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o $(GOBIN)/dist/$(BINARY_NAME)-windows-amd64.exe ./cmd/antoine

# Run the application
.PHONY: run
run:
	@cd $(GOBASE) && go run ./cmd/antoine

# Run with development config
.PHONY: dev
dev:
	@cd $(GOBASE) && go run ./cmd/antoine --config configs/development.yaml --verbose

# Test the application
.PHONY: test
test:
	@echo "Running tests..."
	@cd $(GOBASE) && go test -v ./...

# Test with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@cd $(GOBASE) && go test -v -coverprofile=coverage.out ./...
	@cd $(GOBASE) && go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build cache..."
	@cd $(GOBASE) && go clean
	@rm -rf $(GOBIN)

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@cd $(GOBASE) && go mod download
	@cd $(GOBASE) && go mod tidy

# Lint the code
.PHONY: lint
lint:
	@echo "Running linter..."
	@golangci-lint run

# Format the code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@cd $(GOBASE) && go fmt ./...

# Security scan
.PHONY: security
security:
	@echo "Running security scan..."
	@gosec ./...

# Install the binary
.PHONY: install
install: build
	@echo "Installing Antoine CLI..."
	@cp $(GOBIN)/$(BINARY_NAME) /usr/local/bin/

# Uninstall the binary
.PHONY: uninstall
uninstall:
	@echo "Uninstalling Antoine CLI..."
	@rm -f /usr/local/bin/$(BINARY_NAME)

# Generate documentation
.PHONY: docs
docs:
	@echo "Generating documentation..."
	@cd $(GOBASE) && go run ./cmd/antoine docs generate

# Docker build
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	@docker build -t antoine-cli:$(VERSION) .

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build         - Build the Antoine CLI binary"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  run           - Run the application"
	@echo "  dev           - Run with development config"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Install dependencies"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  security      - Run security scan"
	@echo "  install       - Install binary to system"
	@echo "  uninstall     - Remove binary from system"
	@echo "  docs          - Generate documentation"
	@echo "  docker-build  - Build Docker image"
	@echo "  help          - Show this help message"`

