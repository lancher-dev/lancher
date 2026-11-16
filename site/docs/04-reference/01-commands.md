# Commands

Every command has a `--help` with its usage, flags and examples. Coming from v0.3? See the [conversion table](/docs/getting-started/upgrade#upgrading-from-v03).

## Global Flags

- `--verbose` - Show the output of git instead of a spinner
- `-v, --version` - Show the version
- `-h, --help` - Show help for any command

## `new`

Create a project from a template.

```bash
lancher new [template] [destination] [flags]
```

Without arguments, lancher asks for the template and the destination, which defaults to the template name. All questions come first; then the files are copied, git is initialized and the hooks run, in this order. See [Creating Projects](/docs/usage/creating-projects).

**Flags:**

- `-y, --yes` - Do not ask anything, use the defaults: hooks run only with `--hooks`, git is initialized only with `--git` and lancher files are kept only with `--keep-config`
- `--force` - Write into a destination that is not empty, overwriting files
- `--git`, `--no-git` - Initialize a git repository, or not
- `--hooks`, `--no-hooks` - Run the template hooks, or not
- `--keep-config`, `--no-keep-config` - Copy lancher config and ignore files to the project, or not

**Examples:**

```bash
lancher new
lancher new api ./my-api
lancher new api ./my-api --yes --git --hooks
```

## `add`

Add a template from a directory, ZIP file or git repository.

```bash
lancher add [source] [name] [flags]
```

Sources: a local path, a `.zip` file, a git URL (HTTPS or SSH), `gh:owner/repo` or `gl:owner/repo`. The name defaults to the repository, directory or ZIP name. Without arguments, lancher asks for the source and the name. See [Managing Templates](/docs/usage/templates#add-templates).

**Flags:**

- `--force` - Replace an existing template with the same name

**Examples:**

```bash
lancher add ~/projects/api-base
lancher add gh:owner/repo my-template
lancher add ./starter.zip --force
```

## `list`

List installed templates with their path, source and metadata.

```bash
lancher list [flags]
lancher ls [flags]
```

**Flags:**

- `--json` - Print templates as JSON

## `info`

Show the path, source, git commit, dates, metadata and hooks of a template.

```bash
lancher info <template>
```

## `sync`

Fetch templates again from their sources and replace them. A template is replaced only if the fetch succeeds. Without arguments, lancher syncs every template.

```bash
lancher sync [template...]
```

Templates with an unknown or missing source are skipped with a warning. See [Sync Templates](/docs/usage/templates#sync-templates).

## `remove`

Remove one or more templates. Without arguments, lancher asks which ones to remove.

```bash
lancher remove [name...]
lancher rm [name...]
```

## `cd`

Open a shell in the templates directory, or in the directory of a template. Exit the shell to return. The shell gets `LANCHER_CD` set to its directory.

```bash
lancher cd [template]
```

## `completion`

Generate the shell completion script. Package managers and the install script set it up for you; see [Shell Completion](/docs/getting-started/installation#shell-completion).

```bash
lancher completion bash
lancher completion zsh
lancher completion fish
lancher completion powershell
```
