.PHONY: build clean test install help

# Binary name
BINARY_NAME=domainbadger

# Build directory
BUILD_DIR=.

# Source directory (root now contains main.go)
CMD_DIR=.

# Version information
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)
	@echo "✓ Binary created at $(BUILD_DIR)/$(BINARY_NAME)"

clean: ## Remove build artifacts
	@echo "Cleaning..."
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@rm -rf dist/
	@echo "✓ Cleaned"

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

install: ## Install the binary to $GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	@go install $(LDFLAGS) $(CMD_DIR)
	@echo "✓ Installed to $(GOPATH)/bin/$(BINARY_NAME)"

run: build ## Build and run the binary with example arguments
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME) --help

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✓ Dependencies updated"

all: clean deps build test ## Clean, download deps, build, and test
