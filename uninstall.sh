#!/bin/bash
set -e

# lancher uninstaller script
# Usage: curl -sS https://lancher.go/uninstall.sh | sh

BINARY_NAME="lancher"
INSTALL_DIR="/usr/local/bin"

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

warn() {
    printf "${YELLOW}⚠${NC} %s\n" "$1"
}

# Detect template storage location
get_templates_dir() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        echo "$HOME/Library/Application Support/lancher"
    else
        local xdg_data="${XDG_DATA_HOME:-$HOME/.local/share}"
        echo "$xdg_data/lancher"
    fi
}

main() {
    echo ""
    echo "╔══════════════════════════════════════╗"
    echo "║   lancher Uninstaller                ║"
    echo "╚══════════════════════════════════════╝"
    echo ""

    # Remove binary
    if [ -f "${INSTALL_DIR}/${BINARY_NAME}" ]; then
        info "Removing ${BINARY_NAME} from ${INSTALL_DIR}..."
        if [ -w "${INSTALL_DIR}" ]; then
            rm "${INSTALL_DIR}/${BINARY_NAME}"
        else
            sudo rm "${INSTALL_DIR}/${BINARY_NAME}"
        fi
        success "Binary removed"
    else
        warn "${BINARY_NAME} not found in ${INSTALL_DIR}"
    fi

    # Ask about templates
    local templates_dir=$(get_templates_dir)
    if [ -d "$templates_dir" ]; then
        echo ""
        warn "Found templates directory: $templates_dir"
        read -p "Do you want to remove all stored templates? [y/N] " -n 1 -r
        echo ""
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            rm -rf "$templates_dir"
            success "Templates directory removed"
        else
            info "Templates directory preserved"
        fi
    fi

    echo ""
    success "Uninstallation complete!"
    echo ""
    info "Note: If you installed Go via the lancher installer and want to remove it:"
    echo "  sudo rm -rf /usr/local/go"
    echo "  # Then remove the PATH export from your ~/.bashrc or ~/.zshrc"
    echo ""
}

main "$@"
