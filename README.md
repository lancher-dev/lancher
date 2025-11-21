<h1 align="center" class="b">
  <br>
  lancher 🚀
  <br>
  <br>
</h1>

A minimal, local project-template manager written in Go.

`lancher` lets you organize and use project templates stored on your machine. You can register directories as templates, list them, and generate new projects from them. Templates can come from local paths or from git repositories (via HTTPS or SSH). Each template may include a `.lancher.yaml` file for metadata and post-creation hooks, allowing you to customize how new projects are initialized.

## Installation

The easiest way to install `lancher` is via the dedicated script:

```bash
curl -sS https:/lancher.dev/install.sh | sh
```

The script checks prerequisites, clones the repository, builds the binary, and installs to `/usr/local/bin`.

### Prerequisites

- Go 1.22+
- git
- curl

## Usage

The following commands are provided to make it easier to manage templates and start new projects.

#### create

Create a new project from a template. Can be used interactively or with flags:

```bash
# Interactive mode
lancher create

# With flags
lancher create -t <template_name> -d <destination_dir>
lancher create --template myapp --destination ./new-project
```

#### template

Manage templates with subcommands:

**add** - Add a template from local path or git repository:

```bash
# From local path
lancher template add <name> <source_dir>
lancher template add nextjs-starter ~/projects/my-nextjs-app

# From git repository (https or ssh)
lancher template add <name> <git_url>
lancher template add nextjs https://github.com/user/nextjs-template
lancher template add myapp git@github.com:user/my-template.git

# Interactive mode (prompts for name and source)
lancher template add
```

**list** - List all templates:

```bash
lancher template list
lancher template ls
```

**update** - Update a template:

```bash
# Git pull (only for templates added from git)
lancher template update <template_name>

# Overwrite with new files from path
lancher template update <template_name> -d <new_path>
```

**remove** - Delete a template:

```bash
# Interactive selection
lancher template remove

# Direct removal
lancher template remove <template_name>
lancher template rm <template_name>
```

#### info

Display storage information and list all templates with their paths:

```bash
lancher info
```

### version

Print version information:

```bash
lancher version
lancher -v
lancher --version
```

### help

Display usage information:

```bash
lancher help
```

## Storage

Templates are stored in platform-specific directories. Use `lancher info` to see the storage path on your system.

**Linux:**

- `$XDG_DATA_HOME/lancher/templates` (if `XDG_DATA_HOME` is set)
- `~/.local/share/lancher/templates` (fallback)

**macOS:**

- `~/Library/Application Support/lancher/templates`

## Template Configuration

Templates can include a `.lancher.yaml` file for metadata and post-create hooks.

### .lancher.yaml Format

```yaml
name: My Project Template
description: A template for building awesome projects
author: Your Name
version: 1.0.0

# Commands to run after project creation (executed in project directory)
hooks:
  - npm install
  - git init
  - chmod +x scripts/setup.sh

# Files/patterns to ignore during project creation
ignore:
  - node_modules
  - .git
  - "*.log"
  - .env.local
```

### Configuration Fields

- **name**: Template display name (shown during creation)
- **description**: Brief template description
- **author**: Template author
- **version**: Template version
- **hooks**: Array of shell commands to execute after project creation (requires interactive confirmation)
- **ignore**: File patterns to exclude when copying template (supports glob patterns)

### Hooks

When creating a project from a template with `hooks` defined:

1. Template metadata is displayed
2. Project files are copied (respecting `ignore` patterns)
3. Hooks are listed and require confirmation before execution
4. Each hook executes in the project directory with output shown

## Contributing

Contributions are welcome! Whether you want to fix bugs, add features, or improve documentation, here's how to get started with local development.

### Development Setup

Clone the repository and set up your development environment:

```bash
git clone https://github.com/Kasui92/lancher.git
cd lancher
```

#### Makefile

```bash
make build       # Build binary
make install     # Build and install to /usr/local/bin
make uninstall   # Remove from /usr/local/bin
make test        # Run tests
make clean       # Remove build artifacts
make build-all   # Build for Linux/macOS (amd64/arm64)
make run         # Run without installing (use ARGS='...')
```

#### Testing

```bash
make test
```

Or directly:

```bash
go test ./...
```

#### Local Development

```bash
make run ARGS="list"
make run ARGS="add mytemplate /path/to/source"
```

Or:

```bash
go run ./cmd/lancher list
```

## License

MIT
