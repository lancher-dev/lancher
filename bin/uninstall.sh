#!/bin/sh
# Removes what install.sh installed and, after confirmation, the lancher data.
#
#   curl -fsSL https://lancher.dev/uninstall.sh | sh
set -eu

BIN="$HOME/.local/bin/lancher"
BASH_COMPLETION="${XDG_DATA_HOME:-$HOME/.local/share}/bash-completion/completions/lancher"
FISH_COMPLETION="${XDG_CONFIG_HOME:-$HOME/.config}/fish/completions/lancher.fish"

if [ -t 1 ] && [ -z "${NO_COLOR:-}" ]; then
    GREEN='\033[32m' YELLOW='\033[33m' CYAN='\033[36m' RESET='\033[0m'
else
    GREEN='' YELLOW='' CYAN='' RESET=''
fi

info() { printf '%b→%b %s\n' "$CYAN" "$RESET" "$1"; }
success() { printf '%b✓%b %s\n' "$GREEN" "$RESET" "$1"; }
warn() { printf '%b⚠%b %s\n' "$YELLOW" "$RESET" "$1"; }

data_dir() {
    if [ "$(uname -s)" = Darwin ]; then
        echo "$HOME/Library/Application Support/lancher"
    else
        echo "${XDG_DATA_HOME:-$HOME/.local/share}/lancher"
    fi
}

# confirm asks a yes/no question on the terminal; the answer is no without
# one.
confirm() {
    (: </dev/tty) 2>/dev/null || return 1
    printf "%s [y/N] " "$1"
    read -r answer </dev/tty || answer=""
    case "$answer" in
        [yY]*) return 0 ;;
    esac
    return 1
}

main() {
    if [ -e "$BIN" ]; then
        rm -f "$BIN"
        success "Removed $BIN"
    else
        warn "$BIN not found."
        info "If lancher came from a package manager, remove it there:"
        echo "  sudo apt remove lancher     # APT"
        echo "  sudo pacman -R lancher      # AUR"
        echo "  brew uninstall lancher      # Homebrew"
    fi

    for file in "$BASH_COMPLETION" "$FISH_COMPLETION"; do
        if [ -e "$file" ]; then
            rm -f "$file"
            success "Removed $file"
        fi
    done

    dir=$(data_dir)
    if [ -d "$dir" ]; then
        echo ""
        warn "Your templates are in $dir"
        if confirm "Remove them too?"; then
            rm -rf "$dir"
            success "Removed $dir"
        else
            info "Kept $dir"
        fi
    fi
}

main "$@"
