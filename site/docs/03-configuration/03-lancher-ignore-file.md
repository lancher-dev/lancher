# .lancherignore

A `.lancherignore` file in the template root lists files and directories to leave out of the template:

- `lancher add` and `lancher sync` do not store them;
- `lancher new` does not copy them to the project.

Exclude dependency folders like `node_modules`, build output, logs and editor files. Never put secrets, API keys or credentials in templates.

## Configuration

Write one pattern per line. Empty lines and lines starting with `#` are skipped:

```
# Dependencies
node_modules

# Build output, only at the root
/dist
/build

# Logs everywhere
*.log
```

Patterns are case-sensitive and use `/` as separator. When a directory matches, everything inside it is left out too. The root `.git` directory is never affected: lancher handles it on its own.

## Pattern Types

### Names

A pattern without `/` matches a file or directory name at any depth:

```
node_modules        # node_modules anywhere
package-lock.json   # this file anywhere
.DS_Store           # .DS_Store in any directory
```

### Wildcards

`*` matches any characters within a name, and `?` matches a single character:

```
*.log               # any file ending with .log
test-*              # test-1, test-data, ...
file.?              # file.a, file.1, ...
```

### Paths

A pattern that contains `/` is a path from the template root. A leading `/` makes a single name root-only:

```
/config             # config at the root, not src/config
/*.log              # .log files in the root directory only
src/generated       # src/generated, not lib/src/generated
```

### Trailing Slash

A trailing slash is ignored: `dist/` is the same as `dist` and matches a file or a directory with that name.

### Recursive Patterns

In a path, `**` matches any number of directories, including none:

```
**/*.test.js        # .test.js files at any depth
docs/**             # everything under docs
src/**/test         # test directories at any depth under src
```
