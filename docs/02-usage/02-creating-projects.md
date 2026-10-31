# Creating Projects

Lancher generates new projects by copying template files to a destination directory. You can specify templates and destinations via command-line flags or use interactive prompts to select from available options. The creation process handles validation, metadata display, file copying with pattern exclusions, and optional hook execution.

## Basic Usage

The simplest way to create a project is using the `-t` (template) and `-d` (destination) flags.

```bash
lancher create -t <template> -d <destination>
```

Example:

```bash
lancher create -t nextjs -d ./my-app
```

This command copies the `nextjs` template to `./my-app`. If the destination exists, lancher aborts to prevent overwrites.

## Interactive Mode

```bash
lancher create
```

Interactive mode prompts for template selection, destination directory, git initialization, and hook execution confirmation. It's useful when browsing available templates or when you don't remember exact template names, offering a searchable list of all stored templates.

## Process Flow

When you create a project, Lancher first validates that the template exists and the destination is available. It then displays template metadata like name, description, and version from `.lancher.yaml` if present. After that, it copies all template files to the destination directory while respecting any ignore patterns you've configured. Finally, if the template includes hooks and you confirm their execution, Lancher runs those post-creation commands.

Each step must complete successfully before proceeding. If a hook fails with a non-zero exit code, the process stops, though already-copied files remain in place.

## Examples

### Web Application

```bash
lancher create -t nextjs-starter -d ./client
lancher create -t express-api -d ./server
```

### Monorepo

```bash
lancher create -t monorepo -d ./workspace
```

### Multiple Projects

```bash
for service in auth users products; do
  lancher create -t microservice -d ./services/$service
done
```

## Flags

- `-t, --template <name>` - Template name
- `-d, --destination <path>` - Destination directory
- `--git` - Init git repository
- `--no-git` - Skip git init
- `--hook` - Execute all hooks from `.lancher.yaml`
- `--no-hooks` - Skip hooks execution
- `--keep-config` - Keep lancher config and ignore files in the new project
- `--no-keep-config` - Remove lancher config and ignore files
