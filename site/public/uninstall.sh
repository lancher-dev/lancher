#!/bin/bash
set -e

# lancher uninstaller script proxy
# Usage: curl -sS https://lancher.dev/uninstall.sh | sh
# This script downloads and executes the main uninstaller from the lancher repository

REPO_URL="https://raw.githubusercontent.com/lancher-dev/lancher/main/bin/uninstall.sh"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

info() {
    printf "${BLUE}ℹ${NC} %s\n" "$1"
}

error() {
    printf "${RED}✗${NC} %s\n" "$1"
}

has() {
    command -v "$1" >/dev/null 2>&1
}

# Check for curl
if ! has curl; then
    error "curl is required but not installed."
    echo "Please install curl and try again."
    exit 1
fi

info "Downloading lancher uninstaller from ${REPO_URL}..."

# Download and execute the main uninstaller
if ! curl -fsSL "${REPO_URL}" | sh; then
    error "Failed to download or execute the uninstaller."
    echo ""
    echo "Please check:"
    echo "  - Your internet connection"
    echo "  - The repository URL: ${REPO_URL}"
    echo ""
    exit 1
fi
