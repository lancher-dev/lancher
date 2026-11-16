# Storage

lancher keeps its data in a per-user directory: the templates and a manifest of where each template comes from.

## Locations

### Linux

```
$XDG_DATA_HOME/lancher
```

Fallback, when `XDG_DATA_HOME` is not set:

```
~/.local/share/lancher
```

### macOS

```
~/Library/Application Support/lancher
```

## Structure

```
lancher/
├── templates.yaml        # manifest
└── templates/
    ├── api-base/
    │   ├── .lancher.yaml
    │   └── [template files]
    └── web/
        ├── .git/         # git templates only
        └── [template files]
```

Each directory in `templates/` is a template, named after the directory. Directories whose name is not a valid template name are ignored.

## Manifest

`templates.yaml` records the source of each template, for `lancher sync`, `lancher info` and `lancher list`:

```yaml
version: 1
templates:
  api-base:
    kind: local
    source: /home/me/projects/api-base
    added: 2026-11-16T17:12:14Z
    synced: 2026-11-16T17:12:14Z
  web:
    kind: git
    source: https://github.com/user/web-template
    ref: 6f9a967c2cd11dee2620939b26f98068f0744181
    added: 2026-11-16T17:12:14Z
    synced: 2026-11-16T17:12:14Z
```

| Field    | Meaning                                                 |
| -------- | ------------------------------------------------------- |
| `kind`   | `git`, `local`, `zip`, or `unknown`                     |
| `source` | The git URL (or `gh:`/`gl:` alias) or the absolute path |
| `ref`    | The commit of a git template after the last fetch       |
| `added`  | When the template was added, in UTC                     |
| `synced` | When the template was last fetched, in UTC              |

The directories are the source of truth. Each time lancher reads the manifest, it:

- registers a template directory that has no entry: a git clone gets its source from the `origin` remote, any other directory gets the `unknown` kind;
- drops the entry of a template whose directory is gone.

lancher writes the manifest itself; you do not need to edit it. If it becomes invalid, lancher stops with an error: fix the file or delete it, and lancher rebuilds it from the directories, with unknown sources for templates that are not git clones.

## Manual Changes

Open a shell in the templates directory, or in a template, with:

```bash
lancher cd [template]
```

You can edit template files there. Changes are lost at the next `lancher sync` of that template, which replaces it with a fresh copy of its source.

## Cleanup

Remove templates with `lancher remove`. To remove all of lancher's data, delete its directory:

```bash
rm -rf ~/.local/share/lancher                       # Linux
rm -rf ~/Library/Application\ Support/lancher       # macOS
```
