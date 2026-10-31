import { defineJaadConfig } from "@lancher-dev/jaad";

export default defineJaadConfig({
  site: "https://lancher.dev",

  title: "lancher",
  description: "A minimal, local project-template manager written in Go.",
  lang: "en",

  tailwind: true,

  nav: [{ label: "Docs", href: "/docs" }],

  docsDir: "./docs",
  routeBase: "/docs",

  appearance: "dark",

  social: { github: "https://github.com/lancher-dev/lancher" },

  ogImage: "/og-image.png",
  editLink: true,
  footer: `© ${new Date().getFullYear()} <a href="https://lucapattocchio.dev" target="_blank" rel="noopener">Luca Pattocchio</a> (Kasui)`,

  head: [
    {
      tag: "link",
      attrs: {
        rel: "stylesheet",
        href: "https://www.nerdfonts.com/assets/css/webfont.css",
      },
    },
  ],
});
