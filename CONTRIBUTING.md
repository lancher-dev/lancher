# Contributing to lancher

Thank you for your interest in contributing to lancher! We welcome contributions of all kinds including bug reports, feature requests, documentation improvements, and code contributions.

## Getting Started

### Prerequisites

Before you begin, make sure you have the following installed:

- Go 1.22 or higher
- git
- curl

### Development Setup

1. Fork the repository on GitHub

2. Clone your fork:

```bash
git clone https://github.com/YOUR_USERNAME/lancher.git
cd lancher
```

3. Add the upstream repository:

```bash
git remote add upstream https://github.com/Kasui92/lancher.git
```

4. Create a feature branch:

```bash
git checkout -b feature/your-feature-name
```

## Development Workflow

### Building

Build the binary:

```bash
make build
```

The binary will be created in the root directory as `lancher`.

### Installing Locally

Install to `/usr/local/bin` for system-wide use:

```bash
make install
```

To uninstall:

```bash
make uninstall
```

### Running Without Installing

You can run lancher without installing it:

```bash
make run ARGS="list"
make run ARGS="template add mytemplate /path/to/source"
```

Or directly with Go:

```bash
go run ./cmd/lancher list
go run ./cmd/lancher template add mytemplate /path/to/source
```

### Testing

Run all tests:

```bash
make test
```

Or use Go directly:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

Run tests for a specific package:

```bash
go test ./internal/config
go test ./internal/storage
```

### Multi-Platform Builds

Build for all supported platforms (Linux and macOS, amd64 and arm64):

```bash
make build-all
```

Binaries will be created in the `bin/` directory.

### Cleaning Build Artifacts

Remove built binaries and temporary files:

```bash
make clean
```

## Project Structure

```
lancher/
├── cmd/
│   └── lancher/          # Main entry point
├── internal/
│   ├── cli/              # CLI routing and commands
│   │   ├── commands/     # Main commands (create, info)
│   │   ├── template/     # Template subcommands (add, list, update, remove)
│   │   └── shared/       # Shared utilities (colors, prompts, validation)
│   ├── config/           # Template configuration (.lancher.yaml)
│   ├── storage/          # Platform-specific storage paths
│   ├── fileutil/         # File operations
│   └── version/          # Version information
├── bin/                  # Installation scripts and built binaries
├── .github/workflows/    # GitHub Actions
└── Makefile              # Build automation
```

## Making Changes

1. **Make your changes** in your feature branch

2. **Test your changes** thoroughly:

   ```bash
   make test
   make build
   ./lancher --help
   ```

3. **Commit your changes** with clear, descriptive messages:

   ```bash
   git add .
   git commit -m "Add feature: description of what you added"
   ```

4. **Push to your fork**:

   ```bash
   git push origin feature/your-feature-name
   ```

5. **Open a Pull Request** on GitHub with:
   - Clear description of what the PR does
   - Any related issue numbers
   - Screenshots or examples if applicable

## Syncing Your Fork

Keep your fork up to date with the upstream repository:

```bash
git fetch upstream
git checkout main
git merge upstream/main
git push origin main
```

## Makefile Reference

| Command               | Description                                        |
| --------------------- | -------------------------------------------------- |
| `make build`          | Build the binary for current platform              |
| `make install`        | Build and install to /usr/local/bin                |
| `make uninstall`      | Remove from /usr/local/bin                         |
| `make test`           | Run all tests                                      |
| `make clean`          | Remove build artifacts                             |
| `make build-all`      | Build for all platforms (Linux/macOS, amd64/arm64) |
| `make run ARGS="..."` | Run without installing                             |

## Questions or Issues?

If you have questions or run into issues:

- Check existing [Issues](https://github.com/Kasui92/lancher/issues)
- Open a new issue with details about your problem or question
- Visit [lancher.dev](https://lancher.dev) for documentation

Thank you for contributing!
