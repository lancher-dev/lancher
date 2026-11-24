import { defineJaadConfig } from '@jaamd/jaad'

export default defineJaadConfig({
    site: 'https://lancher.dev',

    title: 'lancher',
    description: 'A minimal, local project-template manager written in Go.',
    lang: 'en',

    tailwind: true,

    nav: [{ label: 'Docs', href: '/docs' }],

    docsDir: './docs',
    routeBase: '/docs',

    appearance: 'dark',

    social: { github: 'https://github.com/lancher-dev/lancher' },

    ogImage: '/og-image.png',
    editLink: 'https://github.com/lancher-dev/lancher/edit/main/site/docs/',
    footer: [
        '<span class="site-footer-links">',
        '<a href="/docs">Docs</a>',
        '<a href="/docs/getting-started/installation">Installation</a>',
        '<a href="https://github.com/lancher-dev/lancher/releases" rel="noopener">Releases</a>',
        '<a href="https://github.com/lancher-dev/lancher" rel="noopener">GitHub</a>',
        '</span>',
        '<span>',
        '<a href="https://github.com/lancher-dev/lancher/blob/main/LICENSE" rel="noopener">MIT License</a>',
        ' · Made by <a href="https://lucapattocchio.dev" rel="noopener">Luca Pattocchio</a>',
        '</span>',
    ].join(''),
})
