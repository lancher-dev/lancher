import { defineConfig } from 'oxfmt'

export default defineConfig({
    printWidth: 120,
    tabWidth: 4,
    singleQuote: true,
    semi: false,
    sortTailwindcss: true,
    sortImports: {
        groups: ['type', 'builtin', 'external', 'internal', ['parent', 'sibling', 'index'], 'side_effect', 'unknown'],
        newlinesBetween: false,
        order: 'asc',
        sortSideEffects: false,
    },
    astro: {},
    ignorePatterns: ['pnpm-lock.yaml'],
})
