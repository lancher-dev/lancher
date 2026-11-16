# Contributing to lancher

Thank you for your interest in contributing to lancher! We welcome contributions of all kinds.

## Prerequisites

- Go, at the version set in `go.mod`
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
make run ARGS="new"
```

Or with Go:

```bash
go run ./cmd/lancher list
```

On Linux, set `XDG_DATA_HOME` to try commands without touching your own templates:

```bash
XDG_DATA_HOME=/tmp/lancher-try go run ./cmd/lancher add ~/projects/api-base
```

### Test

```bash
make test
```

The tests are in `tests/`, one file per area. Most of them run commands through `cli.NewRootCmd()`, with `isolateDataDir` pointing the data directory to a temporary one.

### Install Locally

```bash
make install      # Install to ~/.local/bin
make uninstall    # Remove
```

### Website and docs

The website and documentation live in `site/`. See [site/README.md](site/README.md).

```bash
cd site
pnpm install
pnpm dev
```

## How lancher Works

### Packages

| Package             | Role                                                                            |
| ------------------- | ------------------------------------------------------------------------------- |
| `cmd/lancher`       | `main`: calls `cli.Execute` and exits with its code                             |
| `internal/cli`      | One file per command (cobra), plus `fetch.go` shared by `add` and `sync`        |
| `internal/source`   | Recognizes a source (git, ZIP, local directory), names it and fetches it        |
| `internal/storage`  | Data directory, template paths and names, atomic replace, `templates.yaml`      |
| `internal/project`  | Copies a template into a project, runs `git init` and the hooks                 |
| `internal/config`   | Reads `.lancher.yaml`                                                           |
| `internal/fileutil` | Copies directories, extracts ZIP files, applies `.lancherignore`                |
| `internal/ui`       | Colors, prompts, selection menus and the spinner                                |
| `internal/version`  | The version string, set at build time                                           |

`cli/root.go:NewRootCmd` registers the commands. `cli/root.go:Execute` prints the error of a failed command and turns it into the exit code: 1, or 130 when a prompt is cancelled.

### `lancher add <source> [name]`

1. `cli/add.go:runAdd` reads the arguments, or asks for them with `promptSource` when there are none.
2. `source/source.go:Parse` turns the source into a `Source`:
   - `gh:` and `gl:` become a GitHub or GitLab URL, and also remember the CLI to clone with (`gh`, `glab`);
   - a URL (`https://`, `ssh://`, `git@`, ...) or a path ending in `.git` is a git source;
   - a path ending in `.zip` must be an existing file;
   - anything else must be an existing directory.
3. `Source.Name` derives the name when it is not given. `storage.TemplateExists` checks it (`storage.ValidateName` rejects invalid names), and an existing template without `--force` stops here.
4. `storage/manifest.go:LoadManifest` reads `templates.yaml` and aligns it with the template directories.
5. `cli/fetch.go:fetchTemplate` runs the fetch inside `ui.Task`, which shows a spinner, or the git output with `--verbose`:
   1. `storage/storage.go:Replace` creates `templates/.staging/<name>/` and calls the fetch with the path `templates/.staging/<name>/new`;
   2. `Source.Fetch` writes the template there: `git clone --depth 1` (or `gh repo clone`/`glab repo clone`), `fileutil.UnzipToDir`, or `fileutil.CopyDir`, which leaves out `.git` and the ignored paths. For git and ZIP sources, `fileutil.ApplyIgnorePatterns` then deletes the paths matched by `.lancherignore`. On error, the new directory is deleted and the existing template is untouched;
   3. `Replace` moves the existing template to `.staging/<name>/old`, renames `new` to `templates/<name>`, and deletes `.staging`;
   4. `fetchTemplate` returns the manifest entry: kind, source, sync time and, for git, the commit from `storage.GitRef`.
6. `runAdd` sets `Added`, stores the entry and writes the manifest with `Manifest.Save`, which writes a temporary file and renames it over `templates.yaml`.

On disk, the template is now in `templates/<name>/` and has an entry in `templates.yaml`, next to `templates/`.

### `lancher new [template] [destination]`

`cli/new.go:runNew` works in two phases: `ask` collects every answer without touching the disk, then `plan.apply` does the work.

1. `ask`:
   1. `resolveTemplate` takes the template from the arguments, or with `ui.Select`, and checks that it exists;
   2. `resolveDest` takes the destination, or asks for it (default: the template name), makes it absolute and checks it with `checkDest`: missing, empty, or not empty with `--force`;
   3. `config.LoadConfigWithDetails` reads the template's `.lancher.yaml`, and its metadata is printed;
   4. `askHooks` lists the hooks and asks whether to run them. It skips them when `project.InHook` reports that lancher was started by a hook;
   5. `confirm` resolves git init and, when `project.MetaFiles` finds lancher files in the template, whether to copy them. A flag answers without asking; with `--yes` the answer is no.
2. `plan.apply`:
   1. `project.Copy` copies the template with `fileutil.CopyTree`, leaving out `.git`, the paths matched by `.lancherignore` and, unless kept, the lancher files. Existing files at the destination are overwritten;
   2. `project.InitGit` runs `git init`, unless the destination already has `.git`;
   3. `project.RunHooks` runs each hook with `sh -c` in the destination, with `LANCHER_HOOK=1`, and stops at the first failure.

### Other Commands

- `sync`: `cli/sync.go:runSync` loads the manifest, then for each template parses its stored source (`syncSource`) and calls `fetchTemplate`, like `add`.
- `info`: `cli/info.go:runInfo` prints the manifest entry, the metadata and the hooks.
- `list`: `cli/list.go` prints the manifest with `printList`, or `printJSON` with `--json`.
- `remove`: `cli/remove.go:runRemove` deletes the template directories, then `LoadManifest` drops their entries.
- `cd`: `cli/cd.go:runCd` starts `$SHELL` in the directory with `LANCHER_CD` set.
- Completion of template names: `cli/complete.go:completeTemplates`.

## Submitting Changes

1. Test your changes: `make test && make build`
2. Commit with clear messages
3. Push to your fork
4. Open a Pull Request

## Questions?

- Check [Issues](https://github.com/lancher-dev/lancher/issues)
- Visit [lancher.dev](https://lancher.dev)

Thank you for contributing!
