#!/bin/sh
# Installs lancher in ~/.local/bin, with shell completions.
#
#   curl -fsSL https://lancher.dev/install.sh | sh
#
# LANCHER_VERSION=v1.2.3  install that release instead of the latest one
# LANCHER_STANDALONE=1    skip the package manager suggestion
set -eu

REPO="lancher-dev/lancher"
BASE_URL="${LANCHER_BASE_URL:-https://github.com/$REPO/releases/download}"
INSTALL_DIR="$HOME/.local/bin"
BIN="$INSTALL_DIR/lancher"

if [ -t 1 ] && [ -z "${NO_COLOR:-}" ]; then
    RED='\033[31m' GREEN='\033[32m' YELLOW='\033[33m' CYAN='\033[36m' RESET='\033[0m'
else
    RED='' GREEN='' YELLOW='' CYAN='' RESET=''
fi

info() { printf '%b→%b %s\n' "$CYAN" "$RESET" "$1"; }
success() { printf '%b✓%b %s\n' "$GREEN" "$RESET" "$1"; }
warn() { printf '%b⚠%b %s\n' "$YELLOW" "$RESET" "$1"; }
die() {
    printf '%b✗%b %s\n' "$RED" "$RESET" "$1" >&2
    exit 1
}
has() { command -v "$1" >/dev/null 2>&1; }

detect_platform() {
    case "$(uname -s)" in
        Linux) OS=linux ;;
        Darwin) OS=darwin ;;
        *) die "Unsupported operating system: $(uname -s)" ;;
    esac
    case "$(uname -m)" in
        x86_64 | amd64) ARCH=amd64 ;;
        aarch64 | arm64) ARCH=arm64 ;;
        *) die "Unsupported architecture: $(uname -m)" ;;
    esac
}

# os_like reports whether /etc/os-release names the distribution family $1.
os_like() {
    [ -r /etc/os-release ] || return 1
    # shellcheck disable=SC1091
    ids=$(. /etc/os-release && printf ' %s %s ' "${ID:-}" "${ID_LIKE:-}")
    case "$ids" in *" $1 "*) return 0 ;; esac
    return 1
}

# suggest_package_manager prints the native install commands and returns 0
# when a package manager is available.
suggest_package_manager() {
    if [ "$OS" = darwin ] || has brew; then
        info "Homebrew installs lancher with completions and keeps it updated:"
        echo "  brew install lancher-dev/lancher/lancher"
    elif os_like debian; then
        info "The APT repository installs lancher with completions and keeps it updated:"
        echo "  curl -fsSL https://lancher.dev/pubkey_BA1285237A45B3E4.asc | sudo gpg --dearmor -o /usr/share/keyrings/lancher.gpg"
        echo "  echo 'deb [signed-by=/usr/share/keyrings/lancher.gpg] https://repository.lancher.dev stable main' | sudo tee /etc/apt/sources.list.d/lancher.list"
        echo "  sudo apt update && sudo apt install lancher"
    elif os_like arch; then
        info "The AUR package installs lancher with completions and keeps it updated:"
        echo "  yay -S lancher"
    else
        return 1
    fi
}

# confirm_standalone asks whether to go on, reading from the terminal because
# standard input is the script itself under curl | sh. Without a terminal the
# answer is yes.
confirm_standalone() {
    if ! (: </dev/tty) 2>/dev/null; then
        return 0
    fi
    printf "Continue with the standalone install in %s? [Y/n] " "$INSTALL_DIR"
    read -r answer </dev/tty || answer=""
    case "$answer" in
        [nN]*) return 1 ;;
    esac
}

resolve_version() {
    if [ -n "${LANCHER_VERSION:-}" ]; then
        VERSION="v${LANCHER_VERSION#v}"
        return
    fi
    VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" |
        sed -n 's/.*"tag_name": *"\([^"]*\)".*/\1/p')
    [ -n "$VERSION" ] || die "Failed to find the latest release"
}

sha256() {
    if has sha256sum; then
        sha256sum "$1" | cut -d ' ' -f 1
    elif has shasum; then
        shasum -a 256 "$1" | cut -d ' ' -f 1
    else
        die "sha256sum or shasum is required to verify the download"
    fi
}

install_binary() {
    asset="lancher-$OS-$ARCH"
    TMP_DIR=$(mktemp -d)
    trap 'rm -rf "$TMP_DIR"' EXIT

    info "Downloading lancher $VERSION ($OS-$ARCH)..."
    curl -fsSL -o "$TMP_DIR/$asset" "$BASE_URL/$VERSION/$asset" ||
        die "Failed to download $BASE_URL/$VERSION/$asset"
    curl -fsSL -o "$TMP_DIR/checksums.txt" "$BASE_URL/$VERSION/checksums.txt" ||
        die "Release $VERSION has no checksums.txt; this script installs v0.4.0 and later"

    expected=$(awk -v f="$asset" '$2 == f { print $1 }' "$TMP_DIR/checksums.txt")
    [ -n "$expected" ] || die "No checksum for $asset in checksums.txt"
    [ "$(sha256 "$TMP_DIR/$asset")" = "$expected" ] || die "Checksum mismatch for $asset"

    chmod +x "$TMP_DIR/$asset"
    "$TMP_DIR/$asset" --version >/dev/null 2>&1 || die "The downloaded binary does not run"

    mkdir -p "$INSTALL_DIR"
    mv "$TMP_DIR/$asset" "$BIN"
    success "Installed $("$BIN" --version) in $INSTALL_DIR"
}

install_completions() {
    if has bash; then
        dir="${XDG_DATA_HOME:-$HOME/.local/share}/bash-completion/completions"
        mkdir -p "$dir"
        "$BIN" completion bash >"$dir/lancher"
        success "bash completion installed in $dir"
    fi
    if has fish; then
        dir="${XDG_CONFIG_HOME:-$HOME/.config}/fish/completions"
        mkdir -p "$dir"
        "$BIN" completion fish >"$dir/lancher.fish"
        success "fish completion installed in $dir"
    fi
    if has zsh; then
        info "For zsh completion, add this line to ~/.zshrc:"
        echo "  source <(lancher completion zsh)"
    fi
}

main() {
    detect_platform

    if [ "${LANCHER_STANDALONE:-}" != 1 ] && suggest_package_manager; then
        echo ""
        if ! confirm_standalone; then
            info "Nothing installed."
            exit 0
        fi
    fi

    resolve_version
    install_binary
    install_completions

    case ":$PATH:" in
        *":$INSTALL_DIR:"*) ;;
        *)
            echo ""
            warn "$INSTALL_DIR is not in your PATH. Add this line to your shell profile:"
            echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
            ;;
    esac

    echo ""
    info "Get started:"
    echo "  lancher add <source>    add a template from a directory, ZIP file or git repository"
    echo "  lancher new             create a project from a template"
    echo "  lancher --help          show all commands"
}

main "$@"
