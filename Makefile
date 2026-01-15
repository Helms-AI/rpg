.PHONY: build test run clean install lint

# Binary name
BINARY := rpg

# Build directory
BUILD_DIR := ./bin

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOMOD := $(GOCMD) mod
GOVET := $(GOCMD) vet

# Main package
MAIN_PKG := ./cmd/rpg

# Default target
all: build

# Build the binary
build:
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY) $(MAIN_PKG)
	@echo "Built $(BUILD_DIR)/$(BINARY)"

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run the server (for development)
run: build
	$(BUILD_DIR)/$(BINARY)

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Tidy dependencies
tidy:
	$(GOMOD) tidy

# Download dependencies
deps:
	$(GOMOD) download

# Vet the code
vet:
	$(GOVET) ./...

# Install the binary to $GOPATH/bin
install: build
	cp $(BUILD_DIR)/$(BINARY) $(GOPATH)/bin/$(BINARY)
	@echo "Installed to $(GOPATH)/bin/$(BINARY)"

# Install to /usr/local/bin (requires sudo)
install-global: build
	sudo cp $(BUILD_DIR)/$(BINARY) /usr/local/bin/$(BINARY)
	@echo "Installed to /usr/local/bin/$(BINARY)"

# Show help
help:
	@echo "Available targets:"
	@echo "  build          - Build the binary"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  run            - Build and run the server"
	@echo "  clean          - Remove build artifacts"
	@echo "  tidy           - Tidy go.mod"
	@echo "  deps           - Download dependencies"
	@echo "  vet            - Run go vet"
	@echo "  install        - Install to GOPATH/bin"
	@echo "  install-global - Install to /usr/local/bin (requires sudo)"
