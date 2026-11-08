# lancher.dev

Website and documentation of [`lancher`](https://github.com/lancher-dev/lancher), built with [Astro](https://astro.build) and [jaad](https://github.com/lancher-dev/jaad).

- `docs/`: documentation pages, served under `/docs`.
- `public/`: static files. `install.sh` and `uninstall.sh` are copied from `../bin` at build time.
- `scripts/build-apt.sh`: builds the signed APT repository from the GitHub release assets.
- `wrangler.jsonc`: Cloudflare Workers configuration.

## Development

```sh
pnpm install
pnpm dev
```

`pnpm lint`, `pnpm format` and `pnpm build` run in CI on every change under `site/`.

## Deployment

The `Site` workflow deploys to Cloudflare Workers:

- every published release: production deploy built from the release tag, APT repository included;
- manual run: production deploy of a given tag or branch (default: latest release);
- push to `main` and pull requests: preview version (aliases `main` and `pr-<number>`), without touching production.
