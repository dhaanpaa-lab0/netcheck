# Makefile for netcheck

# Application name
APP_NAME := netcheck

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

# Build output directory
BUILD_DIR := build
DIST_DIR := dist

# Platforms for cross-compilation
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# Build flags
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

.PHONY: all build clean test deps help install cross dist

# Default target
all: clean build

## help: Display this help message
help:
	@echo "Available targets:"
	@echo "  make build         - Build for current platform"
	@echo "  make install       - Build and install to GOPATH/bin"
	@echo "  make cross         - Build for all platforms"
	@echo "  make dist          - Create distribution packages for all platforms"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make deps          - Download dependencies"
	@echo "  make run           - Build and run the application"
	@echo ""
	@echo "Cross-compilation targets:"
	@echo "  make linux-amd64   - Build for Linux AMD64"
	@echo "  make linux-arm64   - Build for Linux ARM64"
	@echo "  make darwin-amd64  - Build for macOS AMD64 (Intel)"
	@echo "  make darwin-arm64  - Build for macOS ARM64 (Apple Silicon)"
	@echo "  make windows-amd64 - Build for Windows AMD64"

## build: Build for current platform
build: deps
	@echo "Building $(APP_NAME) v$(VERSION) for current platform..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

## install: Build and install to GOPATH/bin
install: deps
	@echo "Installing $(APP_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(GOPATH)/bin/$(APP_NAME) .
	@echo "Installed to $(GOPATH)/bin/$(APP_NAME)"

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR) $(DIST_DIR)
	$(GOCLEAN)
	@echo "Clean complete"

## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

## run: Build and run the application
run: build
	@echo "Running $(APP_NAME)..."
	@$(BUILD_DIR)/$(APP_NAME)

## cross: Build for all platforms
cross: clean
	@echo "Cross-compiling for all platforms..."
	@mkdir -p $(BUILD_DIR)
	@$(MAKE) linux-amd64
	@$(MAKE) linux-arm64
	@$(MAKE) darwin-amd64
	@$(MAKE) darwin-arm64
	@$(MAKE) windows-amd64
	@echo "Cross-compilation complete"

## Platform-specific targets
linux-amd64: deps
	@echo "Building for Linux AMD64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 .

linux-arm64: deps
	@echo "Building for Linux ARM64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64 .

darwin-amd64: deps
	@echo "Building for macOS AMD64 (Intel)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 .

darwin-arm64: deps
	@echo "Building for macOS ARM64 (Apple Silicon)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 .

windows-amd64: deps
	@echo "Building for Windows AMD64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe .

## dist: Create distribution packages
dist: cross
	@echo "Creating distribution packages..."
	@mkdir -p $(DIST_DIR)

	# Linux AMD64
	@echo "Packaging Linux AMD64..."
	@mkdir -p $(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-amd64
	@cp $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-amd64/$(APP_NAME)
	@cp README.md $(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-amd64/
	@cp netcheck.txt $(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-amd64/netcheck.txt.example
	@cd $(DIST_DIR) && tar czf $(APP_NAME)-$(VERSION)-linux-amd64.tar.gz $(APP_NAME)-$(VERSION)-linux-amd64
	@rm -rf $(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-amd64

	# Linux ARM64
	@echo "Packaging Linux ARM64..."
	@mkdir -p $(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-arm64
	@cp $(BUILD_DIR)/$(APP_NAME)-linux-arm64 $(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-arm64/$(APP_NAME)
	@cp README.md $(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-arm64/
	@cp netcheck.txt $(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-arm64/netcheck.txt.example
	@cd $(DIST_DIR) && tar czf $(APP_NAME)-$(VERSION)-linux-arm64.tar.gz $(APP_NAME)-$(VERSION)-linux-arm64
	@rm -rf $(DIST_DIR)/$(APP_NAME)-$(VERSION)-linux-arm64

	# macOS AMD64
	@echo "Packaging macOS AMD64..."
	@mkdir -p $(DIST_DIR)/$(APP_NAME)-$(VERSION)-darwin-amd64
	@cp $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(DIST_DIR)/$(APP_NAME)-$(VERSION)-darwin-amd64/$(APP_NAME)
	@cp README.md $(DIST_DIR)/$(APP_NAME)-$(VERSION)-darwin-amd64/
	@cp netcheck.txt $(DIST_DIR)/$(APP_NAME)-$(VERSION)-darwin-amd64/netcheck.txt.example
	@cd $(DIST_DIR) && tar czf $(APP_NAME)-$(VERSION)-darwin-amd64.tar.gz $(APP_NAME)-$(VERSION)-darwin-amd64
	@rm -rf $(DIST_DIR)/$(APP_NAME)-$(VERSION)-darwin-amd64

	# macOS ARM64
	@echo "Packaging macOS ARM64..."
	@mkdir -p $(DIST_DIR)/$(APP_NAME)-$(VERSION)-darwin-arm64
	@cp $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 $(DIST_DIR)/$(APP_NAME)-$(VERSION)-darwin-arm64/$(APP_NAME)
	@cp README.md $(DIST_DIR)/$(APP_NAME)-$(VERSION)-darwin-arm64/
	@cp netcheck.txt $(DIST_DIR)/$(APP_NAME)-$(VERSION)-darwin-arm64/netcheck.txt.example
	@cd $(DIST_DIR) && tar czf $(APP_NAME)-$(VERSION)-darwin-arm64.tar.gz $(APP_NAME)-$(VERSION)-darwin-arm64
	@rm -rf $(DIST_DIR)/$(APP_NAME)-$(VERSION)-darwin-arm64

	# Windows AMD64
	@echo "Packaging Windows AMD64..."
	@mkdir -p $(DIST_DIR)/$(APP_NAME)-$(VERSION)-windows-amd64
	@cp $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(DIST_DIR)/$(APP_NAME)-$(VERSION)-windows-amd64/$(APP_NAME).exe
	@cp README.md $(DIST_DIR)/$(APP_NAME)-$(VERSION)-windows-amd64/
	@cp netcheck.txt $(DIST_DIR)/$(APP_NAME)-$(VERSION)-windows-amd64/netcheck.txt.example
	@cd $(DIST_DIR) && zip -r $(APP_NAME)-$(VERSION)-windows-amd64.zip $(APP_NAME)-$(VERSION)-windows-amd64
	@rm -rf $(DIST_DIR)/$(APP_NAME)-$(VERSION)-windows-amd64

	@echo ""
	@echo "Distribution packages created in $(DIST_DIR):"
	@ls -lh $(DIST_DIR)
