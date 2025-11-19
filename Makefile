.PHONY: build install uninstall clean test run help

# Binary name
BINARY_NAME=lancher
INSTALL_PATH=/usr/local/bin

help:
	@echo "lancher - Development Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build       Build the binary"
	@echo "  make install     Build and install to ${INSTALL_PATH}"
	@echo "  make uninstall   Remove binary from ${INSTALL_PATH}"
	@echo "  make test        Run tests"
	@echo "  make run         Run locally (use ARGS='...' for arguments)"
	@echo "  make clean       Clean build artifacts"
	@echo "  make build-all   Build for multiple platforms"
	@echo ""
	@echo "For end users:"
	@echo "  curl -sS https://raw.githubusercontent.com/Kasui92/lancher/main/install.sh | sh"

# Build the binary
build:
	@echo "Building ${BINARY_NAME}..."
	@go build -o ${BINARY_NAME} cmd/lancher/main.go
	@echo "Build complete: ${BINARY_NAME}"

# Install the binary system-wide
install: build
	@echo "Installing ${BINARY_NAME} to ${INSTALL_PATH}..."
	@sudo cp ${BINARY_NAME} ${INSTALL_PATH}/${BINARY_NAME}
	@sudo chmod +x ${INSTALL_PATH}/${BINARY_NAME}
	@echo "Installation complete. You can now use '${BINARY_NAME}' from anywhere."

# Uninstall the binary
uninstall:
	@echo "Uninstalling ${BINARY_NAME}..."
	@sudo rm -f ${INSTALL_PATH}/${BINARY_NAME}
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
	@go run cmd/lancher/main.go $(ARGS)

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated."

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME}-linux-amd64 cmd/lancher/main.go
	@GOOS=darwin GOARCH=amd64 go build -o ${BINARY_NAME}-darwin-amd64 cmd/lancher/main.go
	@GOOS=darwin GOARCH=arm64 go build -o ${BINARY_NAME}-darwin-arm64 cmd/lancher/main.go
	@echo "Multi-platform build complete."
