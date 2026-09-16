# Trident token prototype

This prototype applies Trident tokens to the existing Hugo and Tailwind docs
site. Builds use checked-in CSS and fonts, with no access to a private registry.

## Preview

```console
$ npm ci
$ hugo server --port 1313
```

Review these pages in light and dark mode, at desktop and mobile widths:

- `/desktop/setup/install/linux/`: article typography, navigation, callouts,
  tables, and code
- `/get-started/docker-concepts/running-containers/sharing-local-files/`: tabs and code blocks
- `/get-started/`: landing pages and cards

Search requires the separate Pagefind indexing step used by the existing build.

## Integration

`assets/css/style.css` imports the legacy theme, the Trident snapshot, and the
docs adapter in that order. Generated Tailwind mappings supply prefixed palette
utilities and replace shared type, spacing, and radius values. Legacy palette utilities remain
available during the migration. The adapter maps selected docs aliases and
components to Trident semantic colors.

`assets/css/trident.css` maps docs aliases to semantic tokens and styles the
header, sidebar selection, article text, cards, buttons, tabs, and callouts.
The `docs` cascade layer follows utilities so existing light and dark classes
cannot override these semantic color pairs. Token imports stay at the root
because they contain Tailwind directives as well as CSS layer declarations.

Manrope supplies the interface and article typeface. Article body text uses the
16-pixel reading scale, while navigation uses the smaller label scale. Code keeps
the existing Roboto Mono font. The existing theme switch activates Trident's
`.dark` overrides.

Article body text and ordinary quotations use Trident's primary foreground
through the local `--docs-reading-foreground` alias. The local
`--docs-reading-background` uses white `background-primary` in light mode and
`sidebar` in dark mode, matching the dark navigation surface. Subtle borders separate the navigation from the reading area. These are
local layout adaptations; the separate Trident high-contrast mode stays inactive.

Navigation cards use `background-paper-elevation-0`, compact 16-pixel titles at
weight 600, and 14-pixel descriptions in Trident gray 700 or gray 300.
The border and shadow follow the surface material in Trident core's
`tri-materials.css`: a semantic border in both themes, with a faint inset
top-edge highlight in dark mode and no visible shadow in light mode. The local
`--docs-card-shadow` mirrors this material because it is outside the token
package. The elevated background is a docs-specific surface choice. Decorative icons are
omitted; titles identify the destinations. Linked cards have a full-area target
and accented hover or focus states. Tabs retain the component surface and
semantic borders. Inactive sidebar
items use `sidebar-foreground-muted`, with stronger text on the active item.
Inactive table-of-contents links are neutral; active and hovered links use the
primary accent. Inline code uses a faint foreground tint and compact padding,
while fenced code uses Trident's default background in light mode and the muted
surface in dark mode to preserve syntax contrast. Chroma and Gordon's highlight.js
map keywords, strings, numbers, attributes, commands, and types to Trident
syntax roles. Comments retain the muted foreground.

The neutral header is a docs-specific adaptation. Trident's `AppHeader` component
uses the `header-from` and `header-to` blue gradient tokens.

Gordon uses the Trident header gradient, popover surface, input colors, and
semantic message, feedback, and alert states. Pagefind inherits the same typeface,
popover surface, foregrounds, accent highlights, and focus colors through its
component variables. Search ranking and chat requests retain their existing
behavior.

This is an integration prototype. Landing-page layouts and components outside
these mappings still need a design review. Loading tokens does not provide
Trident component behavior or certify accessibility.

Component checks cover cards, tab switching, code, search results, and Gordon
at desktop and mobile widths in both themes. Local search checks use the
preview's Pagefind assets because the local indexer fails with a native allocator
error. Gordon message and rate-limit checks use intercepted browser responses;
they verify rendering and interaction, not the live backend.

## Snapshot provenance

- Package: `@docker/trident-tokens@2.0.0-beta.3`
- Repository: <https://github.com/docker/trident>
- Release commit: `bbc0746b17209071ceddabcb14b180fcb32b2fa9`
- CSS: ten unmodified generated files in `assets/css/vendor/trident/`
- Checksums: SHA-256 per file in that directory's `manifest.json`

The registry rejected the available credentials. This snapshot was generated
from the release commit using the upstream `sd.config.ts` and
`scripts/check-dark-coverage.ts`. It was not extracted from the published npm
archive. The build used Node.js and these dependency versions from the upstream
lockfile: `tsx@4.23.1`, `style-dictionary@5.5.1`, `apca-w3@0.1.9`, and
`culori@4.0.2`. Upstream checks reported 97 contrast pairs with zero failures and
explicit dark overrides for all 286 aliased tokens.

Font files come from `@fontsource-variable/manrope@5.2.8` on npmjs.org. The
unmodified WOFF2 subsets and their SIL Open Font License are under
`static/assets/fonts/manrope/`. The font declarations retain the package's
Unicode ranges and use local asset URLs.

## Refresh the snapshot

A maintainer with Trident access can build an agreed release in a separate
checkout, or extract its published package. Then run:

```console
$ node hack/vendor-trident-tokens.mjs /path/to/trident/packages/tokens <40-character-source-commit>
```

The script validates the package name and local CSS imports, copies the
generated files without modification, and records their version and checksums.
Review the CSS diff. Keep docs adjustments in
`assets/css/trident.css`. The refresh step is separate from ordinary builds.

If the token package becomes available publicly on npmjs.org, replace the
snapshot import with the package import and pin its version in `package.json`.
The docs adapter can remain in place.
