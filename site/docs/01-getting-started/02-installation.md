# Installation

Install lancher via script, package manager, or build from source.

## Quick Installation

```bash
curl -sS https://lancher.dev/install.sh | sh
```

The installer will install `lancher` in `~/.local/bin`.

> **Note:** The install script (and the corresponding [uninstall script](https://lancher.dev/uninstall.sh)) are designed to install a **specific fixed version** of lancher. They do not automatically track or pull in newer releases. If you want automatic updates, use your OS package manager instead (see below).

## Arch Linux (AUR)

lancher is available on the [AUR](https://aur.archlinux.org/packages/lancher). Install it with your preferred AUR helper, for example:

```bash
yay -S lancher
# or
paru -S lancher
```

## Ubuntu / Debian

### Via APT repository

For automatic updates alongside your system packages, add the lancher APT repository:

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

You can also download a `.deb` package directly from the [GitHub Releases](https://github.com/lancher-dev/lancher/releases) page and install it manually:

```bash
sudo dpkg -i lancher_<version>_linux_amd64.deb
```

## Homebrew

You can install lancher via [Homebrew Taps](https://docs.brew.sh/Taps) . All you need to do is:

```bash
brew tap lancher-dev/lancher
brew install lancher
```

## Manual Installation

For manual installation you need Go 1.22+ and git. Proceed with repo clone and command installation.

```bash
git clone https://github.com/lancher-dev/lancher.git
cd lancher
make install
```

### Build from Source

```bash
# Clone and build
git clone https://github.com/lancher-dev/lancher.git
cd lancher
go build -o lancher ./cmd/lancher

# Install manually
sudo mv lancher /usr/local/bin/
```

## Third party repositories

lancher is also distributed in other repositories, for easier installation.

[![Packaging status](https://repology.org/badge/vertical-allrepos/lancher.svg)](https://repology.org/project/lancher/versions)
