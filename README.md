# lancher

Minimal local project template manager written in Go.

## Overview

`lancher` manages project templates stored locally. Add directories as templates, list them, and create new projects from them. No remote repositories, no git integration, no placeholder substitution—just directory copying.

## Installation

### Prerequisites

- Go 1.22+
- git
- curl

### Quick Install

```bash
curl -sS https://raw.githubusercontent.com/Kasui92/lancher/main/install.sh | sh
```

The script checks prerequisites, clones the repository, builds the binary, and installs to `/usr/local/bin`.

### Manual Install

```bash
git clone https://github.com/Kasui92/lancher.git
cd lancher
make install
```

### Uninstall

```bash
curl -sS https://raw.githubusercontent.com/Kasui92/lancher/main/uninstall.sh | sh
```

Or manually:

```bash
sudo rm /usr/local/bin/lancher
```

## Commands

### add

Save a directory as a template:

```bash
lancher add <name> <source_dir>
```

Example: `lancher add nextjs-starter ~/projects/my-nextjs-app`

### list

List all templates:

```bash
lancher list
lancher ls
```

### new

Create a new project from a template:

```bash
lancher new <template_name> <destination_dir>
```

Example: `lancher new nextjs-starter ~/projects/new-app`

### remove

Delete a template:

```bash
lancher remove <template_name>
lancher rm <template_name>
```

### help

Display usage information:

```bash
lancher help
```

## Storage

Templates are stored in platform-specific directories:

**Linux:**

- `$XDG_DATA_HOME/lancher/templates` (if `XDG_DATA_HOME` is set)
- `~/.local/share/lancher/templates` (fallback)

**macOS:**

- `~/Library/Application Support/lancher/templates`

## Development

### Structure

```
cmd/lancher/          # Entry point
internal/
  cli/                # Command implementations
  storage/            # Platform-specific paths
  fileutil/           # File operations
```

### Makefile

```bash
make build       # Build binary
make install     # Build and install to /usr/local/bin
make uninstall   # Remove from /usr/local/bin
make test        # Run tests
make clean       # Remove build artifacts
make build-all   # Build for Linux/macOS (amd64/arm64)
make run         # Run without installing (use ARGS='...')
```

### Testing

```bash
make test
```

Or directly:

```bash
go test ./...
```

### Local Development

```bash
make run ARGS="list"
make run ARGS="add mytemplate /path/to/source"
```

Or:

```bash
go run cmd/lancher/main.go list
```

## License

MIT
