# Docker Docs design migration

Product content, upstream mounts, CLI adapters, API specifications, and
synchronization scripts keep their existing locations.

## Architecture

[ARCHITECTURE.md](../../ARCHITECTURE.md) documents navigation registries, sidebar
root selection, layouts, themes, assets, search, and browser behavior. Keep
implementation guidance there; this file records migration provenance, asset
versions, and a link to the test instructions.

## Theme and assets

Trident tokens `2.0.0-beta.11` and icons `0.1.0-beta.5` are pinned to source
revision `c7a3a882e34a2f1383571ba12ede3b5a3e0e885d`. The token release's
`next` mode contains the 2026 brand palette. This replaces the older token format
shown in PR #26110. Generated CSS and SVG geometry are copied without changes;
manifests record hashes. Trident custom glyphs take precedence over Lucide
`1.26.0`. The upstream icon package's Lucide license is included. These Trident
packages do not provide a separate package license notice at that revision.

`assets/css/theme.css` records the docs adaptations: Onest, JetBrains Mono,
neutral inline code backgrounds, and reading widths. Colors use Trident semantic
tokens for surfaces, actions, links, focus indicators, and status.
`category-washes.css` defines four broad, subtle compositions using the pinned
Trident palette. Each category and its reading pages share one composition.
Dark mode uses neutral colors, and footers use a softer companion gradient.
The homepage wash is taller than the reading-page wash. All wash rendering is
CSS: there is no generator, route seed, resize handler, or runtime contrast
calculation. The designs also work without JavaScript.

The fonts retain their license notices. The supplied Docker Docs SVG and PNGs
provide the header logo and favicons. Existing product artwork remains available.

`style.css` is the single Tailwind entry point. Semantic roles also drive prose,
API rendering, search, and Gordon cards. Components retain Tailwind utilities
where appropriate; the shell has dedicated CSS. The former sidebar templates,
Markdown action globals, and Pagefind component styles have been removed.

To add an interface icon, add its canonical SVG filename to the `files` object
in `assets/icons/trident/manifest.json` (an empty hash is sufficient before
refreshing). The refresh script copies each listed icon and records its hash.

For a maintainer refresh, build the pinned Trident token package in a separate
checkout, install the matching `lucide-static` version there, then run:

```console
$ node hack/vendor-trident.mjs /path/to/trident /path/to/lucide-static
```

Review version, manifest, and visual changes together. Public builds do not
access private packages or run this refresh script.

## Routes and installation

The intentional additions are `/c/`, `/c/ai/`, `/c/dev/`, `/c/admin/`,
`/get-docker/desktop/`, `/get-docker/sandboxes/`, and `/get-docker/engine/`.
`/get-docker/` becomes a canonical page. `/get-started/get-docker/` and other
installation aliases redirect directly to it.
`/manuals/` redirects to `/c/`; its source section remains for product content.
Source links target `/get-docker/_index.md`. Learning content is served by
`docker/learn` at `/learn/`; verified replacements for retired Guides and
Get started pages are recorded in `data/redirects.yml`. Other retired learning
pages have no replacement redirect. Intermediate prototype quick-install routes
are not introduced.

`go.mod` temporarily replaces BuildKit with the upstream revision containing
[moby/buildkit#7439](https://github.com/moby/buildkit/pull/7439), which removes
links to the retired glossary. Remove the replacement when a tagged release
includes that fix, then regenerate the vendor tree. The committed tree is
reproducible from the module configuration; it contains no manual vendor edits.

The installation examples are checked against the product installation pages.
Browser detection selects examples, not installation eligibility. Mobile,
unknown platforms, and browsers without JavaScript receive installation links.

## Verification

Site tests and local setup are documented in
[CONTRIBUTING.md](../../CONTRIBUTING.md#testing).
