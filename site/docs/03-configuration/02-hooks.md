# Hooks

Hooks are shell commands that `lancher new` runs in the new project after copying the files, for tasks like installing dependencies.

## Configuration

Define hooks as a list of commands in `.lancher.yaml`:

```yaml
hooks:
  - npm install
  - ./scripts/setup.sh
```

Keep each hook focused on one task, and move longer logic into a script inside the template.

## Execution

1. `lancher new` lists the hooks and asks once whether to run them. The default answer is no.
2. After copying the files and running `git init` (if requested), lancher runs the hooks in the order they are listed, each with `sh -c` in the project directory, showing their output as it comes.
3. If a hook exits with a non-zero code, lancher stops: the following hooks do not run, and the project files stay in place.

Flags skip the question:

- `--hooks` runs the hooks without asking;
- `--no-hooks` skips them;
- with `-y, --yes`, hooks run only if `--hooks` is also given.

Hooks run with `LANCHER_HOOK=1` in their environment. A `lancher new` started from a hook does not run the hooks of its own template.

Hooks run with `sh`, which is `dash` on Debian and Ubuntu: write them in POSIX shell, without bash features like `[[ ]]` or `source`.

## Common Patterns

### Dependency Installation

```yaml
hooks:
  - npm install
```

```yaml
hooks:
  - pip install -r requirements.txt
```

### First Commit

Run `lancher new` with `--git`, then commit with a hook:

```yaml
hooks:
  - git add -A
  - git commit -m "Initial commit"
```

### Script Execution

```yaml
hooks:
  - chmod +x scripts/*.sh
  - ./scripts/setup.sh
```

### Environment Setup

```yaml
hooks:
  - cp .env.example .env
```

### Database Setup

```yaml
hooks:
  - npm run db:migrate
  - npm run db:seed
```

## Conditional Execution

### Check Command Availability

```yaml
hooks:
  - if command -v pnpm >/dev/null; then pnpm install; else npm install; fi
```

### Platform-Specific

```yaml
hooks:
  - |
    if [ "$(uname -s)" = Darwin ]; then
      brew bundle
    fi
```

### File Existence

```yaml
hooks:
  - if [ -f package.json ]; then npm install; fi
```

## Security

Hooks run commands on your machine with your permissions. Before running the hooks of a template from someone else, read them:

```bash
lancher info <template>
```

`lancher new` also lists them before asking. **Never use `--hooks` with templates you have not reviewed.** A `lancher sync` can change the hooks of a template, so review them again after syncing.
