# Uninstallation

If you installed lancher with the install script, run the uninstall script:

```bash
curl -fsSL https://lancher.dev/uninstall.sh | sh
```

It removes `~/.local/bin/lancher` and the completion files that the install script wrote. Then it asks whether to remove your templates too; the default answer is no, and without a terminal the templates are kept.

If you installed lancher with a package manager, remove it there:

```bash
sudo apt remove lancher     # APT
sudo pacman -R lancher      # AUR
brew uninstall lancher      # Homebrew
```

Package managers do not touch your templates. They are in the [lancher data directory](/docs/reference/storage#locations); delete it to remove them.
