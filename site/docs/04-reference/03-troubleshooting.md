# Troubleshooting

Common errors and how to solve them. lancher prints errors with a `✗` and exits with a non-zero code.

## Template Not Found

```
✗ template 'myapp' not found
```

The name does not match any installed template. Check the names with `lancher list`, then add the template if it is missing.

## Template Already Exists

```
✗ template 'api-base' already exists (use --force to replace it)
```

`lancher add` does not overwrite templates. Pass another name as second argument, or `--force` to replace the template.

## Invalid Template Name

```
✗ invalid template name "my app": use letters, digits, '.', '_' and '-', starting with a letter or a digit
```

Pass a valid name as second argument to `lancher add`.

## Git Clone Failed

```
✗ failed to clone repository: exit status 128
fatal: repository 'https://github.com/user/missing/' not found
```

lancher shows the last lines of git's output. Check the URL with `git clone <url>`, and your access to the repository: SSH keys for SSH URLs, a credential helper for HTTPS. lancher stops git from asking for a password; run the command again with `--verbose` to see the whole output and to let git ask.

## Template Skipped by Sync

```
⚠ 'handmade' skipped: unknown source (re-add it with: lancher add --force <source> handmade)
⚠ 'api-base' skipped: directory not found: '/home/me/projects/api-base'
```

lancher does not know where the template comes from, or its source is gone. Add the template again from its current source:

```bash
lancher add --force ~/projects/new-place/api-base api-base
```

## Destination Is Not Empty

```
✗ destination is not empty: /home/me/my-app (use --force to write into it)
```

Choose another destination, or pass `--force` to write the template into it: files with the same path are overwritten, the others are kept.

## Hook Failed

```
✗ hook "npm install" failed: exit status 1
```

The project files are in place; the hooks after the failed one did not run. Read the hook output above the error, fix the cause (often a missing tool), and run the remaining commands by hand in the project directory. Hooks run with `sh`: bash syntax like `[[ ]]` or `source` fails on systems where `sh` is `dash`.

## Invalid YAML

```
✗ failed to load template config: invalid YAML in .lancher.yaml: yaml: line 3: did not find expected key
```

Fix the file at the line shown: use spaces, not tabs, for indentation, and quote strings that contain `:` or `#`. Open the template with `lancher cd <template>`.

## Invalid Manifest

```
✗ invalid manifest /home/me/.local/share/lancher/templates.yaml (fix or delete it): yaml: line 12: did not find expected key
```

Fix the file, or delete it: lancher rebuilds it from the template directories, and templates that are not git clones get an unknown source until you add them again.

## Permission Denied

Run lancher as your own user, not with `sudo`, and check that you own the data directory:

```bash
ls -la ~/.local/share/lancher
```

## Hooks Do Not Run

Hooks run only after you confirm, or with `--hooks`. With `-y` they need `--hooks`. Check them with `lancher info <template>`.

## Getting Help

Search the existing [issues](https://github.com/lancher-dev/lancher/issues) and [discussions](https://github.com/lancher-dev/lancher/discussions) first. If you do not find an answer, [open an issue](https://github.com/lancher-dev/lancher/issues/new) with your lancher version (`lancher --version`), your operating system, the complete error message, the steps to reproduce the problem and, if relevant, the template's `.lancher.yaml`.
