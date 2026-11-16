# Installation

Install lancher with a package manager, with the install script, or from source. Package managers keep lancher updated and install shell completion for bash, zsh and fish.

## Quick Installation

```bash
curl -fsSL https://lancher.dev/install.sh | sh
```

The script:

1. detects your system and, if a package manager is available (APT, AUR or Homebrew), shows the command to install lancher with it and asks whether to go on with the standalone install;
2. downloads the latest release for your platform and verifies its SHA256 checksum;
3. installs `lancher` in `~/.local/bin`;
4. installs bash and fish completion in your user directories and prints the line to add to `~/.zshrc` for zsh.

It never runs `sudo` and does not edit your shell files. If `~/.local/bin` is not in your `PATH`, it prints the line to add.

The script does not update lancher by itself: run it again to install a newer release.

| Variable                 | Effect                                                   |
| ------------------------ | -------------------------------------------------------- |
| `LANCHER_VERSION=v1.2.3` | Install that release instead of the latest one           |
| `LANCHER_STANDALONE=1`   | Skip the package manager suggestion and install directly |

```bash
curl -fsSL https://lancher.dev/install.sh | LANCHER_VERSION=v0.4.0 sh
```

Without a terminal (CI, scripts) the script does not ask and installs directly.

## Arch Linux (AUR)

lancher is available on the [AUR](https://aur.archlinux.org/packages/lancher). Install it with your preferred AUR helper:

```bash
yay -S lancher
# or
paru -S lancher
```

## Ubuntu / Debian

### Via APT repository

Add the lancher APT repository to get updates with your system packages:

```bash
# 1. Add the GPG public key
curl -fsSL https://lancher.dev/pubkey_BA1285237A45B3E4.asc \
  | sudo gpg --dearmor -o /usr/share/keyrings/lancher.gpg

# 2. Add the repository source
echo "deb [signed-by=/usr/share/keyrings/lancher.gpg] \
  https://repository.lancher.dev stable main" \
  | sudo tee /etc/apt/sources.list.d/lancher.list

# 3. Install
sudo apt update && sudo apt install lancher
```

### Via .deb package

Download a `.deb` package from the [GitHub Releases](https://github.com/lancher-dev/lancher/releases) page and install it:

```bash
sudo dpkg -i lancher_<version>_linux_amd64.deb
```

## Homebrew

lancher is available from its [Homebrew tap](https://github.com/lancher-dev/homebrew-lancher):

```bash
brew install lancher-dev/lancher/lancher
```

## From Source

Building lancher needs git and Go at the version set in [`go.mod`](https://github.com/lancher-dev/lancher/blob/main/go.mod).

```bash
git clone https://github.com/lancher-dev/lancher.git
cd lancher
make install    # builds and installs lancher in ~/.local/bin
```

Or build the binary and place it yourself:

```bash
go build -o lancher ./cmd/lancher
sudo mv lancher /usr/local/bin/
```

## Shell Completion

The APT, AUR and Homebrew packages install completion for bash, zsh and fish, and the install script installs it for bash and fish. Completion covers commands, flags and the names of your templates.

If you built lancher from source, generate the completion files yourself:

```bash
# bash
mkdir -p ~/.local/share/bash-completion/completions
lancher completion bash > ~/.local/share/bash-completion/completions/lancher

# fish
lancher completion fish > ~/.config/fish/completions/lancher.fish

# zsh: add this line to ~/.zshrc
source <(lancher completion zsh)
```

Bash completion needs the `bash-completion` package.

## Third Party Repositories

lancher is also distributed in other repositories.

[![Packaging status](https://repology.org/badge/vertical-allrepos/lancher.svg)](https://repology.org/project/lancher/versions)
