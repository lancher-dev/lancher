# .lancherignore

Define which files and directories to exclude during template copy operations. Patterns use glob syntax and are case-sensitive, evaluated against each file's relative path within the template.

You should always exclude version control directories like `.git`, dependency folders like `node_modules`, and build output directories. Never include secrets, API keys, or credentials in templates. Platform-specific files (`.DS_Store`, `Thumbs.db`) and IDE configurations (`.vscode`, `.idea`) should typically be excluded as well.

## Configuration

Create a `.lancherignore` file in your template root directory. Each line represents a pattern:

```
node_modules
.git
*.log
dist/
build/
```

Lines starting with `#` are treated as comments and help document why specific patterns are excluded:

```
# Dependencies
node_modules/
pnpm-lock.yaml

# Build output (only at root)
/dist
/build

# Logs everywhere
*.log
**/*.log
```

Patterns are evaluated against each file's relative path from the template root, matching against the basename, the full relative path, and any individual path component.

Patterns use the system's filepath matching. An empty `.lancherignore` file will include _ALL_ files in the template.

> **Note**: Path separators are normalized automatically. Both `/` and `\` work on all platforms.

## Pattern Types

### Exact Match

Matches the exact file or directory name anywhere in the template:

```
node_modules        # Matches node_modules at any depth
package-lock.json   # Matches this file anywhere
.DS_Store           # Matches .DS_Store in any directory
```

### Wildcards

Use `*` to match any characters (except path separator) and `?` to match a single character:

```
*.log               # Matches any file ending with .log
*.tmp               # Matches any .tmp file
test-*              # Matches test-1, test-data, etc.
file.?              # Matches file.a, file.1, etc.
```

### Root-Relative Patterns

Patterns starting with `/` only match at the template root:

```
/config             # Matches config/ at root, NOT src/config
/build              # Matches build/ at root only
/*.log              # Matches .log files in root directory only
```

### Directories

Trailing slash indicates a directory (optional but clarifies intent):

```
dist/               # Matches dist directory (same as 'dist')
build/              # Matches build directory
.cache/             # Matches .cache directory
```

### Recursive Patterns

Use `**` for recursive matching across any directory depth:

```
**/*.test.js        # Matches .test.js files at any depth
**/node_modules     # Matches node_modules at any level
**/vendor           # Matches vendor directory anywhere
dist/**             # Matches everything under dist/
src/**/test         # Matches test directories under src/
```

## Deprecated Method

> [!WARNING]
> The `ignore` field in `.lancher.yaml` is deprecated and will be removed in a future version. Migrate to `.lancherignore`.

If using the old method, define patterns as an array in `.lancher.yaml`:

```yaml
ignore:
  - node_modules
  - "*.log"
```

**Precedence**: If both `.lancherignore` and the `ignore` field in `.lancher.yaml` exist, `.lancherignore` takes precedence and the YAML ignore field is ignored with a deprecation warning.
