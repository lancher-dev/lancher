#!/bin/bash
set -e

# lancher installer script
# Usage: curl -sS https://raw.githubusercontent.com/Kasui92/lancher/main/install.sh | sh

REPO="Kasui92/lancher"
BINARY_NAME="lancher"
INSTALL_DIR="/usr/local/bin"
MIN_GO_VERSION="1.22"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

success() {
    echo -e "${GREEN}✓${NC} $1"
}

error() {
    echo -e "${RED}✗${NC} $1"
}

warn() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Compare version numbers
version_ge() {
    # Returns 0 if $1 >= $2
    printf '%s\n%s' "$2" "$1" | sort -V -C
}

# Check if Go is installed and meets minimum version
check_go() {
    if ! command -v go &> /dev/null; then
        error "Go is not installed"
        echo ""
        echo "Please install Go ${MIN_GO_VERSION} or later:"
        echo "  https://go.dev/doc/install"
        echo ""
        return 1
    fi

    local version=$(go version | awk '{print $3}' | sed 's/go//')
    info "Found Go ${version}"

    if ! version_ge "$version" "$MIN_GO_VERSION"; then
        error "Go ${version} is too old. Minimum required: ${MIN_GO_VERSION}"
        echo ""
        echo "Please upgrade Go to version ${MIN_GO_VERSION} or later:"
        echo "  https://go.dev/doc/install"
        echo ""
        return 1
    fi

    success "Go version ${version} meets requirements"
    return 0
}

# Build and install lancher
install_lancher() {
    local tmp_dir=$(mktemp -d)
    local src_dir="${tmp_dir}/lancher"

    info "Cloning lancher repository..."
    if ! git clone --depth 1 "https://github.com/${REPO}.git" "${src_dir}" 2>/dev/null; then
        error "Failed to clone repository. Make sure git is installed."
        rm -rf "${tmp_dir}"
        exit 1
    fi

    cd "${src_dir}"

    info "Building lancher..."
    if ! go build -o "${BINARY_NAME}" cmd/lancher/main.go; then
        error "Failed to build lancher"
        rm -rf "${tmp_dir}"
        exit 1
    fi

    info "Installing lancher to ${INSTALL_DIR}..."
    if [ -w "${INSTALL_DIR}" ]; then
        cp "${BINARY_NAME}" "${INSTALL_DIR}/"
        chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    else
        sudo cp "${BINARY_NAME}" "${INSTALL_DIR}/"
        sudo chmod +x "${INSTALL_DIR}/${BINARY_NAME}"
    fi

    # Clean up
    cd - > /dev/null
    rm -rf "${tmp_dir}"

    success "lancher installed successfully to ${INSTALL_DIR}/${BINARY_NAME}"
}

# Main installation process
main() {
    echo ""
    echo "╔══════════════════════════════════════╗"
    echo "║   lancher Installer                  ║"
    echo "║   Local Project Template Manager     ║"
    echo "╚══════════════════════════════════════╝"
    echo ""

    # Check for required tools
    if ! command -v git &> /dev/null; then
        error "git is required but not installed."
        echo "Please install git and try again."
        exit 1
    fi

    if ! command -v curl &> /dev/null; then
        error "curl is required but not installed."
        echo "Please install curl and try again."
        exit 1
    fi

    # Check Go version
    if ! check_go; then
        exit 1
    fi

    # Install lancher
    install_lancher

    echo ""
    success "Installation complete!"
    echo ""
    info "You can now use lancher:"
    echo "  ${BINARY_NAME} help"
    echo "  ${BINARY_NAME} add mytemplate /path/to/project"
    echo "  ${BINARY_NAME} list"
    echo "  ${BINARY_NAME} new mytemplate /path/to/new-project"
    echo ""
}

main "$@"
