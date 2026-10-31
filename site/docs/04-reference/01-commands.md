# Commands

Lancher provides two main commands: `create` for generating projects from templates, and `template` for managing your template library. Commands support both flags and interactive prompts.

## `create`

Create a new project from a template using flags or interactive prompts.

```bash
lancher create -t <template> -d <path_to_destination>
lancher create --template <template> --destination <path_to_destination>
```

**Flags:**

- `-t, --template <name>` - Template name
- `-d, --destination <path>` - Destination directory
- `-p, --print` - Show command output
- `--git` - Init git repository
- `--no-git` - Skip git init
- `--hook` - Execute all hooks from `.lancher.yaml`
- `--no-hooks` - Skip hooks execution
- `--keep-config` - Keep lancher config and ignore files in the new project
- `--no-keep-config` - Remove lancher config and ignore files

**Interactive:**

```bash
lancher create
```

**Examples:**

```bash
lancher create -t nextjs -d ./my-app
lancher create --template react --destination ~/projects/new-app
```

The command validates template existence before copying. If the destination exists, the operation aborts to prevent overwrites.

## `template add`

Add template from local path, Zip file, git repository, or GitHub/GitLab alias. Local templates are copied, git templates are cloned.

```bash
lancher template add <name> <source>
```

**Arguments:**

- `<name>` - Template name
- `<source>` - Local path, Zip file, git URL oh GitHub/GitLab alias

**Flags:**

- `-p, --print` - Show command output

**Interactive:**

```bash
lancher template add
```

**Examples:**

```bash
# Local
lancher template add nextjs ~/projects/nextjs-app

# Git HTTPS
lancher template add react https://github.com/user/react-template

# Git SSH
lancher template add vue git@github.com:user/vue-template.git

# GitLab CLI
lancher template add flutter gl:user/flutter-template

# GitHub CLI
lancher template add ruby-on-rails gh:user/ruby-on-rails-template
```

## `template list`

List all templates and their information.

```bash
lancher template list
lancher template ls
lancher templates
```

**Output:** for each template: name, storage path, and metadata from `.lancher.yaml` (description, author, version) when available.

## `template update`

Update template.

```bash
lancher template update <name>
lancher template update <name> -d <path>
```

**Flags:**

- `-d, --directory` - New source path (overwrites template)
- `-p, --print` - Show command output

**Examples:**

```bash
# Git pull
lancher template update react

# Overwrite
lancher template update nextjs -d ~/projects/updated-nextjs
```

## `template remove`

Remove template.

```bash
lancher template remove <name>
lancher template rm <name>
```

**Interactive:**

```bash
lancher template remove
```

**Examples:**

```bash
lancher template remove old-template
lancher template rm unused
```

## `version`

Display version.

```bash
lancher -v
lancher --version
```

## `help`

Display help.

```bash
lancher help
lancher --help
lancher -h
```

All commands and subcommands have help options.
