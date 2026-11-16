# Creating Projects

`lancher new` creates a project by copying a template into a destination directory, then optionally initializes git and runs the template [hooks](/docs/configuration/hooks).

```bash
lancher new [template] [destination]
```

Example:

```bash
lancher new nextjs ./my-app
```

## Interactive Mode

```bash
lancher new
```

lancher asks for everything first, then works without further questions:

1. the template, from a list;
2. the destination, which defaults to the template name;
3. whether to run the hooks, after listing them, if the template has any;
4. whether to initialize a git repository;
5. whether to copy the lancher files (`.lancher.yaml`, `.lancherignore`) to the project, if the template has any.

The default answer to questions 3-5 is no. Flags answer the matching question in advance: `--hooks` or `--no-hooks`, `--git` or `--no-git`, `--keep-config` or `--no-keep-config`.

## Non-Interactive Mode

`-y, --yes` asks nothing. The template and the destination are then required, and every question gets its default answer unless a flag says otherwise:

```bash
lancher new api ./my-api --yes                 # copy only
lancher new api ./my-api --yes --git --hooks   # also git init and hooks
```

This is the form to use in scripts:

```bash
for service in auth users products; do
  lancher new microservice "./services/$service" --yes --git
done
```

## Destination

The destination may not exist yet or be an empty directory; its parent directory must exist. A destination that contains files is an error:

```
✗ destination is not empty: /home/me/my-app (use --force to write into it)
```

With `--force`, lancher writes the template into it: files with the same path are overwritten and the others are left alone. Nothing is deleted.

## What Gets Copied

lancher copies every file of the template except:

- the `.git` directory;
- the paths matched by the template's [`.lancherignore`](/docs/configuration/lancher-ignore-file);
- the lancher files (`.lancher.yaml`, `.lancherignore` and the other [config names](/docs/configuration/lancher-file)), unless you keep them.

## Process

1. lancher reads the template's `.lancher.yaml`, shows its metadata and asks the questions.
2. It copies the files.
3. It runs `git init`, if requested. If the destination already has a `.git` directory, it is left as it is.
4. It runs the hooks, if requested, one after another in the project directory. If a hook fails, lancher stops and exits with an error; the project files stay in place.

## Flags

- `-y, --yes` - Do not ask anything, use the defaults
- `--force` - Write into a destination that is not empty, overwriting files
- `--git` - Initialize a git repository
- `--no-git` - Do not initialize a git repository
- `--hooks` - Run the template hooks
- `--no-hooks` - Do not run the template hooks
- `--keep-config` - Keep lancher config and ignore files in the project
- `--no-keep-config` - Do not copy lancher config and ignore files to the project
