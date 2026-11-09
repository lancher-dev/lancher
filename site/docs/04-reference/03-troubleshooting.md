# Troubleshooting

Common errors and solutions. Issues typically stem from missing templates, git authentication, file permissions, or failed hooks.

## Template Not Found

This error occurs when the specified template name doesn't exist in storage.

**Error:**

```
Error: template 'myapp' not found
```

**Solution:**

```bash
lancher template list
```

Verify template name and add if missing.

## Git Clone Failed

Git cloning can fail due to authentication, network, or invalid URLs.

**Error:**

```
Error: failed to clone repository
```

**Solutions:** Check network connectivity and verify your git credentials (SSH keys, tokens). Test the URL manually with `git clone <url>` and try switching between HTTPS and SSH if one fails.

## Hook Execution Failed

Hooks fail when commands return non-zero exit codes, usually indicating missing dependencies or configuration issues.

**Error:**

```
Error: hook 'npm install' failed with exit code 1
```

**Solutions:** Verify the required tools are installed (`npm`, `git`, etc.) and test the command manually in the project directory. Review the command syntax in `.lancher.yaml` and read the output carefully for specific error details.

## Permission Denied

**Error:**

```
Error: permission denied writing to storage directory
```

**Solutions:** Check directory permissions with `ls -la ~/.local/share/lancher`, ensure the storage directory exists, verify available disk space with `df -h`, and confirm lancher is running as your own user (avoid `sudo`).

## Destination Already Exists

**Error:**

```
Error: destination directory already exists
```

**Solution:**
Remove or rename the existing directory, or choose a different destination.

## Invalid YAML

**Error:**

```
Error: invalid .lancher.yaml syntax
```

**Solutions:** Validate the YAML syntax at https://www.yamllint.com. Use spaces (not tabs) for indentation, quote strings that contain special characters, and double-check the array/object structure.

## Build Failed

**Error:**

```
Error: failed to build lancher
```

**Solutions:** Verify your Go version (`go version`, requires 1.22+) and update if necessary. Check that `GOPATH` and `GOBIN` are configured correctly, then try `go clean -cache && go build`.

## Template Update Failed

**Error:**

```
Error: cannot update non-git template
```

**Solution:**
Use `-d` flag to specify new source:

```bash
lancher template update <name> -d <new-path>
```

Or remove and re-add template.

## Common Issues

### Templates Not Listed

```bash
lancher template ls  # Verify storage location
ls -la ~/.local/share/lancher/templates
```

### Hooks Don't Execute

Check that you confirmed the prompt (or used `--hooks` to skip it), verify the hook syntax in `.lancher.yaml`, and test the commands manually in the project directory.

### Git Updates Don't Work

Verify the template was originally added from a git URL. Check that you still have access to the remote repository, or try running `git pull` manually inside the template directory (`cd <template-dir> && git pull`).

## Getting Help

If you can't find a solution to your problem, the best way to get help is through the lancher GitHub repository. Start by searching existing [issues](https://github.com/lancher-dev/lancher/issues) and [discussions](https://github.com/lancher-dev/lancher/discussions) to see if someone has encountered the same problem. Many common questions have already been answered in these threads.

If you don't find what you're looking for, feel free to [open a new issue](https://github.com/lancher-dev/lancher/issues/new). When opening an issue, include your lancher version (`lancher --version`), operating system and version, complete error messages, steps to reproduce the problem, and the template configuration (`.lancher.yaml`) if relevant.

The community and maintainers will help diagnose and resolve your issue.
