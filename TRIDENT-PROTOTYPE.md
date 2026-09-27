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

Manrope supplies the interface and article typeface. Article and homepage backgrounds use `background` in both themes. Article text
and quotations use `foreground`. Navigation retains the `sidebar` surface.
The separate Trident high-contrast mode stays inactive.

Home, Get started, Guides featured cards, and Manuals share the
`components/card.html` partial. Cards use `card`, `card-foreground`,
`muted-foreground`, and `elevation-raised-shadow`. Titles and descriptions use
Trident's body and small-body typography scales. The homepage uses the display,
body, label, and heading scales. Layout dimensions remain specific to the docs.

Inline and fenced code use `muted` in both themes. Trident removed its syntax
color tokens in beta.5, so code uses `foreground`, with `muted-foreground` for
comments. Syntax categories have no separate colors in this prototype.

The header uses the `sidebar` surface. Tabs and messages use `card`, and the
homepage question field uses `input-bg`.

Gordon uses the Trident header gradient, popover surface, input colors, and
semantic message, feedback, and alert states. Pagefind inherits the same typeface,
popover surface, foregrounds, accent highlights, and focus colors through its
component variables. Search ranking and chat requests retain their existing
behavior.

The homepage composes Gordon's question form, suggested questions, shared
navigation cards, and release feed with the same Trident palette. The form uses
the input surface and a visible focus outline; suggested questions use quieter
bordered controls. The navigation and release feed share a content width.
The decorative background pattern and scaling interactions are omitted.
These page-specific rules live in `assets/css/home.css`.

This is an integration prototype. Components outside these mappings still need
a design review. Loading tokens does not provide
Trident component behavior or certify accessibility.

## Snapshot provenance

- Package: `@docker/trident-tokens@2.0.0-beta.5`
- Repository: <https://github.com/docker/trident>
- Source commit: `22ff6c848be181d439fd1f5d057175dddfdc9430`
- CSS: ten unmodified generated files in `assets/css/vendor/trident/`
- Checksums: SHA-256 per file in that directory's `manifest.json`

This snapshot was generated
from the source commit using the upstream `sd.config.ts` and
`scripts/check-dark-coverage.ts`. It was not extracted from the published npm
archive. The build used Node.js and these dependency versions from the upstream
lockfile: `tsx@4.23.1`, `style-dictionary@5.5.1`, `apca-w3@0.1.9`, and
`culori@4.0.2`. Upstream checks reported 87 contrast pairs with zero failures and
explicit dark overrides for all 243 aliased tokens.

Font files come from `@fontsource-variable/manrope@5.2.8` on npmjs.org. The
unmodified WOFF2 subsets and their SIL Open Font License are under
`static/assets/fonts/manrope/`. The font declarations retain the package's
Unicode ranges and use local asset URLs.

## Removed prototype overrides

- Alternate article backgrounds: `background-primary` in light mode and `sidebar`
  in dark mode. Both use `background`.
- Gray 700/300 descriptions. Descriptions use `muted-foreground`.
- The copied dark inset card highlight and transparent light shadow. Cards and
  the homepage input use `elevation-raised-shadow`.
- Elevated paper surfaces for cards and the homepage input. They use `card` and
  `input-bg`, respectively.
- Light/dark fenced-code background switching and the 5% inline-code color mix.
  Both use `muted`.
- The Roboto Mono font override and custom card/homepage type sizes, line heights,
  and letter spacing. These use the corresponding Trident typography tokens.
- References to removed `syntax-*` tokens. Code uses semantic foreground tokens.

The dark `background` token is `oklch(0.1 0.008 245)` in both snapshots.
The beta.5 source does not contain a lighter replacement for that token.

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
