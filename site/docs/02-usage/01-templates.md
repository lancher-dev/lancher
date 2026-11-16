# Managing Templates

lancher stores templates in its [data directory](/docs/reference/storage). A template can come from a local directory, a ZIP file or a git repository, and lancher remembers its source so that `lancher sync` can fetch it again.

## Add Templates

```bash
lancher add <source> [name]
```

The name defaults to the repository, directory or ZIP file name. Names use letters, digits, `.`, `_` and `-`, and start with a letter or a digit.

Without arguments, lancher asks for the source and the name:

```bash
lancher add
```

If a template with the same name exists, `add` stops. Use `--force` to replace it: the old template is replaced only after the new one has been fetched.

```bash
lancher add --force ~/projects/api-base
```

### From a Local Directory

lancher copies the directory, without its `.git` directory and without the paths listed in its [`.lancherignore`](/docs/configuration/lancher-ignore-file).

```bash
lancher add ~/projects/api-base
lancher add ~/projects/api-base api
```

### From a ZIP File

lancher extracts the file. If every entry is inside a single top-level directory, that directory becomes the template root.

```bash
lancher add ./starter.zip
```

### From a Git Repository

lancher clones the latest commit of the default branch. Authentication uses your git configuration: SSH keys or a credential helper. git does not prompt for a password unless you pass `--verbose`, which also shows its output.

```bash
# HTTPS
lancher add https://github.com/user/react-template

# SSH
lancher add git@github.com:user/vue-template.git
```

### From GitHub or GitLab

`gh:` and `gl:` are short forms for GitHub and GitLab repositories. lancher clones them with [GitHub CLI](https://cli.github.com/) or [GitLab CLI](https://gitlab.com/gitlab-org/cli) when installed, so their login is used, and with git over HTTPS otherwise.

```bash
lancher add gh:user/react-template
lancher add gl:group/vue-template
```

> [!NOTE]
> A git template keeps its `.git` directory, which `lancher sync` and `lancher info` use. `lancher new` never copies it to new projects.

## List Templates

```bash
lancher list
# or
lancher ls
```

For each template, `list` shows its path, its source and the metadata from its [`.lancher.yaml`](/docs/configuration/lancher-file), with any warning about the config file.

```
Available Templates:

  • api-base
    Path: ~/.local/share/lancher/templates/api-base
    Source: ~/projects/api-base (local)
    Name: API base
    Description: Go HTTP API
```

`--json` prints the templates as JSON, for scripts:

```bash
lancher list --json
```

## Inspect a Template

```bash
lancher info <template>
```

`info` shows the path, the source, the git commit, when the template was added and synced, its metadata and its hooks. Read the hooks here before running them with `lancher new`.

```
api-base
    Path: ~/.local/share/lancher/templates/api-base
    Source: ~/projects/api-base (local)
    Added: 2026-11-16 18:12
    Synced: 2026-11-16 18:12
    Name: API base
    Description: Go HTTP API
    Hooks:
      1. go mod tidy
      2. git add -A
```

## Sync Templates

```bash
lancher sync [template...]
```

`sync` fetches templates again from their sources: a new clone for git repositories, a new copy for directories and ZIP files. A template is replaced only if the fetch succeeds. Without arguments, lancher syncs every template.

```bash
lancher sync           # every template
lancher sync api web   # only these
```

`sync` skips, with a warning, a template whose source is unknown or no longer exists, and goes on with the others. It exits with an error if a fetch fails.

```
⚠ 'handmade' skipped: unknown source (re-add it with: lancher add --force <source> handmade)
```

To change the source of a template, add it again with `--force`.

## Remove Templates

```bash
lancher remove <name...>
# or
lancher rm <name...>
```

Without arguments, lancher lists the templates and lets you select the ones to remove:

```bash
lancher remove
```

> [!WARNING]
> Removing a template deletes its directory. It cannot be undone.

## Open the Templates Directory

```bash
lancher cd [template]
```

`cd` opens a shell in the templates directory, or in the directory of a template. Edit the files there, then `exit` to return. It needs a terminal, and cannot be started from a shell that `cd` opened. The shell has `LANCHER_CD` set to its directory, which you can show in your prompt.

Changes made there are lost at the next `lancher sync` of that template.
