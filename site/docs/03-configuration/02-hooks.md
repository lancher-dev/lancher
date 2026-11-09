# Hooks

Hooks are shell commands that automate post-creation tasks, executing sequentially in the new project directory after file copying. Each command requires explicit user confirmation, displays real-time output, and halts the sequence on failure.

## Configuration

Define hooks as an array of shell command strings in `.lancher.yaml`.

```yaml
hooks:
  - npm install
  - git init
  - ./scripts/setup.sh
```

Keep your hooks simple and focused on single tasks. For complex operations, extract the logic into separate shell scripts. Always use relative paths to ensure portability, and handle errors gracefully so failures provide clear feedback.

## Execution Behavior

Hooks execute after all template files have been copied to the destination directory, running sequentially in the order you defined them. Each command executes in the newly created project directory, and before running, Lancher requires explicit user confirmation while displaying the full command text. The output of each command appears in real-time as it executes. If any command returns a non-zero exit code, the entire hook sequence stops immediately.

This behavior ensures you maintain control over potentially destructive operations like dependency installation or git initialization. Since files are already copied before hooks run, failed commands won't leave your project in a partial state. You can bypass confirmation prompts using the `--hooks` flag, but exercise caution—only use it with templates you trust and have reviewed.

## Common Patterns

### Dependency Installation

```yaml
hooks:
  - npm install
```

```yaml
hooks:
  - pnpm install
```

```yaml
hooks:
  - pip install -r requirements.txt
```

### Git Initialization

```yaml
hooks:
  - git init
  - git add .
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
  - echo "Configure .env file"
```

### Database Setup

```yaml
hooks:
  - npm run db:migrate
  - npm run db:seed
```

### Multi-Step Setup

```yaml
hooks:
  - npm install
  - npm run build
  - git init
  - npx husky install
  - chmod +x scripts/*.sh
```

## Conditional Execution

### Check Command Availability

```yaml
hooks:
  - command -v pnpm && pnpm install || npm install
```

### Platform-Specific

```yaml
hooks:
  - |
    if [[ "$OSTYPE" == "darwin"* ]]; then
      brew install something
    fi
```

### File Existence

```yaml
hooks:
  - test -f package.json && npm install || echo "No package.json"
```

## Security

Hooks execute shell commands on your system, so always review them before running. Lancher requires explicit confirmation and displays the full command text before execution.

The `--hooks` flag automatically confirms all prompts. **Never use `--hooks` with templates from untrusted sources** without first inspecting the `.lancher.yaml` file. Take standard shell security precautions and note that hooks cannot execute other `lancher` commands to prevent infinite loops—design your hooks to perform direct operations.
