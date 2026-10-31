# Managing Templates

Lancher manages project templates from multiple sources: local filesystem, Zip files, or git repositories (HTTPS/SSH). Templates are stored in platform-specific directories and can be listed, updated, or removed as needed.

## Add Templates

Templates can be sourced from local directories, remote git repositories, or GitHub/GitLab. Local templates are copied to lancher's storage, while git templates are cloned and can be kept up-to-date. Use the interactive mode to get guided through the process.

### From Local Directory

Copies the directory structure from the specified path into lancher's template storage.

```bash
lancher template add <name> <path>
```

Example:

```bash
lancher template add nextjs ~/projects/nextjs-base
```

### From Git Repository

Git-based templates maintain their repository history for pulling updates later. Authentication uses your git configuration.

```bash
lancher template add <name> <git-url>
```

Examples:

```bash
# HTTPS
lancher template add react https://github.com/user/react-template

# SSH
lancher template add vue git@github.com:user/vue-template.git
```

### From GitHub/GitLab CLI

Use the `gh:` (GitHub) and `gl:` (GitLab) aliases to add templates via [github-cli](https://cli.github.com/) or [gitlab-cli](https://gitlab.com/gitlab-org/cli). Lancher checks for CLI presence first, otherwise falls back to HTTPS URLs.

```bash
lancher template add <name> gh:<user/organization>/<repo>
lancher template add <name> gl:<user/team>/<repo>
```

Examples:

```bash
# GitHub
lancher template add react gh:user/react-template

# GitLab
lancher template add vue gl:user/vue-template.git
```

## List Templates

```bash
lancher template list
# or
lancher template ls
```

The list command provides a quick overview of all available templates, showing name, storage path, and any metadata from the template configuration (description, author, version) when a `.lancher.yaml` is present.

## Update Templates

Keep your templates synchronized with their sources. Git templates can fetch latest changes, while any template can be completely replaced.

### Git Templates

```bash
lancher template update <name>
```

Performs `git pull` to fetch latest changes from the remote repository. Only works for templates added via git.

> [!Note]
> Although a git configuration may be present in the template, all references to it are not imported when creating a new project. To keep new projects always "fresh", everything in the `.git/` folder is ignored.

### Overwrite with New Files

```bash
lancher template update <name> -d <new-path>
```

Replaces the entire template with content from the new path.

## Remove Templates

Delete templates from storage when no longer needed. This permanently removes the template directory.

```bash
lancher template remove <name>
# or
lancher template rm <name>
```

Interactive mode (multiple select from list):

```bash
lancher template remove
```

**Warning:** This action cannot be undone.
