# Docker Docs design migration

Product content, upstream mounts, CLI adapters, API specifications, and
synchronization scripts keep their existing locations.

## Architecture

[ARCHITECTURE.md](../../ARCHITECTURE.md) documents navigation registries, sidebar
root selection, layouts, themes, assets, search, and browser behavior. Keep
implementation guidance there; this file records migration provenance, asset
versions, and verification procedures.

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

Use Hugo `0.169.0` and the repository's npm lockfile. Build the baseline in a
separate checkout, then build this branch. Flatten both artifacts before comparing:

```console
$ hugo --destination /tmp/docs-candidate --panicOnWarning --printPathWarnings
$ node hack/flatten-and-resolve.js /tmp/docs-candidate
$ python3 hack/redesign/verify.py /tmp/docs-baseline /tmp/docs-candidate
$ node hack/api-docs/verify-output.mjs /tmp/docs-candidate
$ node --test hack/test/*.mjs
```

The comparison protects HTML routes, article heading IDs, and redirects. Eight
existing Desktop settings headings contain Hugo's internal shortcode counter;
the check normalizes that build-dependent counter while retaining the heading
prefix and shortcode slot. Those URLs are not stable anchors on the baseline.

Build Pagefind `1.5.2` into the candidate and serve it over HTTP. Use `--environment development` and a local `--baseURL` when testing
clipboard links without loading production analytics. Install Chromium with
`npx playwright install chromium`, then run:

```console
$ DOCS_TEST_URL=http://localhost:1313 npm run test:redesign
```

`CHROMIUM_PATH` can select an existing Chromium binary. The browser suite checks
both themes at desktop, tablet, and mobile sizes, navigation context, clipboard
success and failure, tab keyboards, search results and errors, focus restoration,
and navigation without JavaScript. Focused regressions cover sidebar resizing,
CLI table-of-contents sizing, pending clipboard writes, and search and chat
failure cases. `interactions.mjs` also checks shared code
copying and clipboard failures, header alignment, installation links, and the
pinned Ask row's handoff to Gordon. The page sweep rejects empty prose and
navigation links. `washes.mjs` checks category variation in both themes, consistent reading-page backgrounds,
and distinct footer gradients. Captures go to `/tmp/docs-redesign-captures`.

On ARM hosts with 64 KiB memory pages, the published musl Pagefind binary may
fail in jemalloc. Building the same release with
`cargo install pagefind --version 1.5.2 --locked` uses the native GNU allocator.
This is a local tooling choice; the production container remains unchanged.

Keep the previous deployment artifact for rollback. Before deployment, refresh
from main, repeat the checks, and review the production artifact. After launch,
monitor broken links, search failures, and browser errors.
