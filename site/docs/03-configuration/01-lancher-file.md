# .lancher.yaml

The `.lancher.yaml` file is an optional YAML configuration that defines template behavior with metadata (name, description, author, version) and post-creation shell commands (hooks).

## File Format

The configuration uses standard YAML syntax with predefined top-level keys.

```yaml
name: Template Name
description: Brief description
author: Author Name
version: 1.0.0
hooks:
  - command1
  - command2
```

The file can be named `.lancher.yaml`, `.lancher.yml`, `lancher.yaml`, or `lancher.yml`. If multiple configuration files exist, a warning is displayed and lancher uses the highest-priority one (`.lancher.yaml` > `.lancher.yml` > `lancher.yaml` > `lancher.yml`). View template status with `lancher template ls`.

## Fields

Each field serves a specific purpose in template configuration. Missing fields are ignored, allowing minimal configurations.

### name

**Type:** `string`
**Optional**

Template display name shown during creation.

```yaml
name: Next.js TypeScript Starter
```

### description

**Type:** `string`
**Optional**

Brief template description.

```yaml
description: Production-ready Next.js with TypeScript and Tailwind
```

### author

**Type:** `string`
**Optional**

Template author or maintainer.

```yaml
author: Development Team <dev@example.com>
```

### version

**Type:** `string`
**Optional**

Template version. Use semantic versioning.

```yaml
version: 2.1.0
```

### hooks

**Type:** `array of strings`
**Optional**

Shell commands executed in the project directory after file copy, requiring user confirmation.

```yaml
hooks:
  - npm install
  - git init
  - chmod +x scripts/setup.sh
```

Commands run sequentially in the project directory. Users must confirm before execution, and the process stops on first failure.

### ignore

> [!WARNING]
> **Deprecated**: Use `.lancherignore` file instead. This field will be removed in a future version. See [`.lancherignore`](./03-lancher-ignore-file.md) for details.

**Type:** `array of strings`
**Optional**

File patterns excluded during copy. Supports glob syntax.

```yaml
ignore:
  - node_modules
  - .git
  - "*.log"
  - .env.local
```

## Examples

<details>
<summary>Node.js Project</summary>

```yaml
name: Node.js API
description: Express.js REST API with TypeScript
author: API Team
version: 1.0.0

hooks:
  - npm install
  - npm run build
  - git init
```

</details>

<details>
<summary>Python Project</summary>

```yaml
name: Python CLI
description: Python CLI application template
version: 1.0.0

hooks:
  - python -m venv venv
  - source venv/bin/activate && pip install -r requirements.txt
```

</details>

<details>
<summary>Monorepo</summary>

```yaml
name: Full-Stack Monorepo
description: React + Node.js monorepo with workspaces
version: 2.0.0

hooks:
  - npm install
  - npm run bootstrap
  - npx husky install
```

</details>
