.PHONY: build install uninstall clean test run help

# Binary name
BINARY_NAME=lancher
INSTALL_PATH=$(HOME)/.local/bin

# Version info
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X github.com/lancher-dev/lancher/internal/version.Version=$(VERSION) -X github.com/lancher-dev/lancher/internal/version.Commit=$(COMMIT)"

help:
	@echo "lancher - Development Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build       Build the binary"
	@echo "  make install     Build and install to ${INSTALL_PATH}"
	@echo "  make uninstall   Remove binary from ${INSTALL_PATH}"
	@echo "  make test        Run all tests"
	@echo "  make run         Run locally (use ARGS='...' for arguments)"
	@echo "  make clean       Clean build artifacts"
	@echo "  make build-all   Build for multiple platforms"
	@echo ""
	@echo "For end users:"
	@echo "  curl -sS https://lancher.dev/install.sh | sh"

# Build the binary
build:
	@echo "Building ${BINARY_NAME}..."
	@go build $(LDFLAGS) -o ${BINARY_NAME} ./cmd/lancher
	@echo "Build complete: ${BINARY_NAME}"

# Install the binary to user's local bin
install: build
	@echo "Installing ${BINARY_NAME} to ${INSTALL_PATH}..."
	@mkdir -p ${INSTALL_PATH}
	@cp ${BINARY_NAME} ${INSTALL_PATH}/${BINARY_NAME}
	@chmod +x ${INSTALL_PATH}/${BINARY_NAME}
	@echo "Installation complete. You can now use '${BINARY_NAME}' from anywhere."

# Uninstall the binary
uninstall:
	@echo "Uninstalling ${BINARY_NAME}..."
	@rm -f ${INSTALL_PATH}/${BINARY_NAME}
	@echo "Uninstall complete."

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f ${BINARY_NAME}
	@go clean
	@echo "Clean complete."

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run the application locally (for testing)
run:
	@go run ./cmd/lancher $(ARGS)

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated."

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o ${BINARY_NAME}-linux-amd64 ./cmd/lancher
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o ${BINARY_NAME}-linux-arm64 ./cmd/lancher
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o ${BINARY_NAME}-darwin-amd64 ./cmd/lancher
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o ${BINARY_NAME}-darwin-arm64 ./cmd/lancher
	@echo "Multi-platform build complete."
