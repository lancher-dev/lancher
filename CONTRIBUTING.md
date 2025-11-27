# Contributing to lancher

Thank you for your interest in contributing to lancher! We welcome contributions of all kinds.

## Prerequisites

- Go 1.22 or higher
- git

## Development Setup

1. Fork and clone the repository:

```bash
git clone https://github.com/YOUR_USERNAME/lancher.git
cd lancher
```

2. Create a feature branch:

```bash
git checkout -b feature/your-feature-name
```

## Local Development

### Build

```bash
make build
```

### Run Without Installing

```bash
make run ARGS="list"
make run ARGS="create"
```

Or with Go:

```bash
go run ./cmd/lancher list
```

### Test

```bash
make test
```

### Install Locally

```bash
make install      # Install to /usr/local/bin
make uninstall    # Remove
```

## Submitting Changes

1. Test your changes: `make test && make build`
2. Commit with clear messages
3. Push to your fork
4. Open a Pull Request

## Questions?

- Check [Issues](https://github.com/Kasui92/lancher/issues)
- Visit [lancher.dev](https://lancher.dev)

Thank you for contributing!
