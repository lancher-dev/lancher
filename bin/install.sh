#!/bin/bash
set -e

REPO="Kasui92/lancher"
BINARY_NAME="lancher"
INSTALL_DIR="$HOME/.local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

info() {
    printf "${BLUE}ℹ${NC} %s\n" "$1"
}

success() {
    printf "${GREEN}✓${NC} %s\n" "$1"
}

error() {
    printf "${RED}✗${NC} %s\n" "$1"
}

warn() {
    printf "${YELLOW}⚠${NC} %s\n" "$1"
}

has() {
    command -v "$1" >/dev/null 2>&1
}

# Detect OS and architecture
detect_platform() {
    local os=$(uname -s | tr '[:upper:]' '[:lower:]')
    local arch=$(uname -m)

    case "$os" in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        *)
            error "Unsupported operating system: $os"
            exit 1
            ;;
    esac

    case "$arch" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            error "Unsupported architecture: $arch"
            exit 1
            ;;
    esac

    PLATFORM="${OS}-${ARCH}"
    info "Detected platform: ${PLATFORM}"
}

# Get latest release version from GitHub
get_latest_version() {
    info "Fetching latest release version..."

    if ! has curl; then
        error "curl is required but not installed."
        exit 1
    fi

    local api_url="https://api.github.com/repos/${REPO}/releases/latest"
    LATEST_VERSION=$(curl -sS "${api_url}" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "$LATEST_VERSION" ]; then
        error "Failed to fetch latest version"
        exit 1
    fi

    info "Latest version: ${LATEST_VERSION}"
}

# Download and install binary
install_binary() {
    local tmp_dir=$(mktemp -d)
    local binary_name="lancher-${PLATFORM}"
    local download_url="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${binary_name}"

    info "Downloading ${binary_name}..."

    if ! curl -sS -L -o "${tmp_dir}/${BINARY_NAME}" "${download_url}"; then
        error "Failed to download binary"
        error "URL: ${download_url}"
        rm -rf "${tmp_dir}"
        exit 1
    fi

    # Make binary executable
    chmod +x "${tmp_dir}/${BINARY_NAME}"

    # Verify binary works
    if ! "${tmp_dir}/${BINARY_NAME}" version >/dev/null 2>&1; then
        error "Downloaded binary is not working correctly"
        rm -rf "${tmp_dir}"
        exit 1
    fi

    # Install binary
    info "Installing to ${INSTALL_DIR}..."
    
    # Create directory if it doesn't exist
    mkdir -p "${INSTALL_DIR}"
    
    # Move binary to install directory
    mv "${tmp_dir}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"

    # Clean up
    rm -rf "${tmp_dir}"

    success "${BINARY_NAME} ${LATEST_VERSION} installed successfully"
}

# Main installation process
main() {
    echo ""
    echo "╔══════════════════════════════════════╗"
    echo "║   lancher Installer                  ║"
    echo "╚══════════════════════════════════════╝"
    echo ""

    # Detect platform
    detect_platform

    # Get latest version
    get_latest_version

    # Install binary
    install_binary

    echo ""
    success "Installation complete!"
    echo ""
    info "You can now use lancher:"
    echo "  ${BINARY_NAME} help"
    echo "  ${BINARY_NAME} template add mytemplate /path/to/project"
    echo "  ${BINARY_NAME} create"
    echo ""
    info "Check your version:"
    echo "  ${BINARY_NAME} version"
    echo ""
    
    # Check if ~/.local/bin is in PATH
    if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
        warn "~/.local/bin is not in your PATH"
        echo ""
        echo "Add the following to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
        echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
        echo ""
        echo "Then reload your shell or run:"
        echo "  source ~/.bashrc  # or ~/.zshrc"
        echo ""
    fi
}

main "$@"
