# Upgrade

If you installed lancher with a package manager (APT, AUR, Homebrew), it is updated with your other packages.

If you used the [install script](/docs/getting-started/installation#quick-installation), run it again: it installs the latest release.

```bash
curl -fsSL https://lancher.dev/install.sh | sh
```

## Upgrading from v0.3

lancher v0.4 replaces the `create` and `template` commands with flat commands. The old commands are gone, with no aliases.

| v0.3                                               | v0.4                                |
| -------------------------------------------------- | ----------------------------------- |
| `lancher create -t <template> -d <dest>`           | `lancher new <template> <dest>`     |
| `lancher create`                                   | `lancher new`                       |
| `lancher template add <name> <source>`             | `lancher add <source> [name]`       |
| `lancher template list`, `ls`, `lancher templates` | `lancher list`, `lancher ls`        |
| `lancher template update <name>`                   | `lancher sync <name>`               |
| `lancher template update <name> -d <path>`         | `lancher add --force <path> <name>` |
| `lancher template remove`, `rm`                    | `lancher remove`, `lancher rm`      |
| `-p`, `--print`                                    | `--verbose`                         |
| `--hook`                                           | `--hooks`                           |

`add` takes the source first and the name second, and the name is optional: it defaults to the repository, directory or ZIP name.

New in v0.4:

- [`lancher info`](/docs/reference/commands#info) shows the source, metadata and hooks of a template;
- [`lancher sync`](/docs/reference/commands#sync) without arguments fetches every template again, from any kind of source;
- [`lancher cd`](/docs/reference/commands#cd) opens a shell in the templates directory;
- [`lancher completion`](/docs/reference/commands#completion) generates shell completion;
- `lancher list --json` prints templates as JSON.

Behavior changes in `new`:

- a destination that is not empty is an error unless you pass `--force`; lancher no longer offers to delete it;
- `-y, --yes` creates a project without questions; hooks then run only with `--hooks`.

### Existing Templates

lancher v0.4 keeps a [manifest](/docs/reference/storage#manifest) of where each template comes from. The first time it runs, it registers the templates you already have:

- a template cloned from git gets its source from the `origin` remote, so `lancher sync` works right away;
- any other template gets an unknown source. `lancher sync` skips it until you add it again with its source:

```bash
lancher add --force ~/projects/api-base api-base
```
