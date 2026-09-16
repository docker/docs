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
through the local `--docs-reading-foreground` alias. The reading surface uses
`background-primary`: white in light mode and a lifted dark surface in dark
mode. Subtle borders separate the navigation from the reading area. These are
local layout adaptations; the separate Trident high-contrast mode stays inactive.

Card descriptions and metadata retain their muted colors. Inactive sidebar
items use `sidebar-foreground-muted`, with stronger text on the active item.
Inactive table-of-contents links are neutral; active and hovered links use the
primary accent. Inline code uses a faint foreground tint and compact padding,
while fenced code retains its existing syntax highlighting.

The neutral header is a docs-specific adaptation. Trident's `AppHeader` component
uses the `header-from` and `header-to` blue gradient tokens.

This is an integration prototype. Syntax highlighting, search internals, Gordon,
landing-page layouts, and other components still need a design review. Loading
tokens does not provide Trident component behavior or certify accessibility.

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
