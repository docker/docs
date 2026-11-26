# Docker Docs site architecture

This document explains how the site selects navigation, renders pages, and applies
shared styles and browser behavior. Use it when changing templates, registering a
product, or investigating a layout problem. For content authoring, see
[COMPONENTS.md](COMPONENTS.md). For the API generation pipeline, see
[hack/api-docs/README.md](hack/api-docs/README.md).

## Content and build boundaries

Hugo builds the site from `content/`, mounted upstream content, and generated API
presentation data. `layouts/` renders that effective content tree. Product pages
remain under `content/manuals/` in source; `hugo.yaml` strips the `/manuals` prefix
from published product URLs. Internal source references retain that prefix.

The section at `content/manuals/_index.md` remains part of the content tree but
doesn't render a directory page. Its former `/manuals/` URL redirects to `/c/`,
which lists documentation categories. `/get-docker/` provides installation paths;
installation aliases redirect there. Tutorials and learning paths are built by
[docker/learn](https://github.com/docker/learn) and served at `/learn/`.
The former Guides and Get started trees, their layouts, and their exclusive
assets have been removed. `data/redirects.yml` maps verified equivalents and
their historical aliases directly to Learn. Pages without equivalents have no
replacement redirect. Learn owns its content, navigation, and assets; this Hugo
site links to it without importing or building its pages.

Vendored CLI data, mounted content, and upstream API specifications keep their
existing update workflows. Template changes must work with Hugo's mounted and
generated pages as well as files present directly in `content/`.

## Navigation

### Declare collection boundaries

Three files define navigation collections:

| File | Responsibility |
| --- | --- |
| [products.json](data/navigation/products.json) | Product IDs, source roots, labels, categories, and related-reference links |
| [references.json](data/navigation/references.json) | Reference-family roots and labels |
| [categories.json](data/navigation/categories.json) | Category labels, landing URLs, homepage cards, and category-page grouping |

These files identify collection boundaries. Hugo's content tree supplies the
pages within each collection. Adding a page below an existing root doesn't
require adding it to a registry.

To register a product, add its existing source root to `products.json`. Add its ID
to a category group to expose it in discovery pages. A product omitted from those
groups still has local navigation. A category entry with `page` instead of
`collection` adds a shortcut without changing ownership.

To register a reference family, add its root to `references.json` and include that
root in the appropriate products' `references` lists. The build checks that roots
and reference destinations exist, product roots aren't duplicated, categories
are valid, and reference families have product entry points. See
[validate.html](layouts/_partials/navigation/validate.html).

### Resolve the sidebar root

[context.html](layouts/_partials/navigation/context.html) finds the most specific
declared root that is the current page or one of its Hugo ancestors. Callers cache
the result per page using `partialCached`. Matching uses Hugo page identity and
ancestry, so published URL rewriting doesn't determine sidebar ownership.

For example, `/reference/api/engine` is a declared reference root. Its sidebar
starts there and contains its descendants. A nested product such as MCP Catalog
and Toolkit can have its own sidebar within the Desktop source tree because the
more specific product root takes precedence.

Reference pages derive their category from products that list the selected
reference root. [reference-category.html](layouts/_partials/navigation/reference-category.html)
caches this result per root. Multiple products can share a reference: if their
categories agree, that category supplies the active top navigation, breadcrumb,
and wash. Conflicting or missing categories leave the category unset. This
doesn't assign a single product owner to the reference or change its sidebar.

The resolver has these explicit exceptions:

- Docker CLI descendants share the complete `/reference/cli/docker` command tree,
  including commands associated with other reference families or categories.
- Discovery pages omit the local section sidebar. Generated API pages
  use the API navigation described under [Page layouts](#page-layouts).
- A page without a declared root or an applicable exception has no local sidebar.

### Share the published breadcrumb trail

[navigation/breadcrumbs.html](layouts/_partials/navigation/breadcrumbs.html)
resolves the category and published ancestors once per page. The visible
breadcrumb, structured `BreadcrumbList`, and Pagefind section metadata use this
same trail. Source-only ancestors such as Manuals don't appear in it.

### Render the content tree

[sidebar.html](layouts/_partials/navigation/sidebar.html) renders the overview
link, tree, and related-reference links.
[tree.html](layouts/_partials/navigation/tree.html) recursively walks Hugo's
`.Pages`; [children.html](layouts/_partials/navigation/children.html) filters the
children used both for rendering and for deciding whether to show an expand
control.

Front matter controls labels, weights, groups, reverse ordering, badges, link
overrides, and active-child filtering. See [COMPONENTS.md](COMPONENTS.md#front-matter).
The current page and ancestor disclosures are marked or expanded on the server.
A hidden page stays out of unrelated sidebars. Following a direct link into a
hidden section reveals that branch, without revealing unrelated hidden branches.
Keep the same navigation state and filtering rules throughout recursion.

Related-reference links come from the product config. Non-rendered API roots have
no permalink; [url.html](layouts/_partials/navigation/url.html) finds their
generated overview in `data/api-reference.json`. Missing or ambiguous destinations
fail the build rather than producing empty links.

## Page layouts

[baseof.html](layouts/baseof.html) provides the shared header, search dialog, chat
drawer, reading area, and footer. It uses the navigation context to choose whether
to render a sidebar. [reading-header.html](layouts/_partials/reading-header.html)
provides the breadcrumb, title, description, and page actions. Discovery layouts
replace the main content block; content layouts retain the shared reading header.

The homepage renders three category cards from `categories.json`. Its What's new
timeline uses `data/whats-new.json`, shows featured entries initially, and reveals
other entries in their chronological positions. The Markdown homepage retains
the category and product inventory. Reference directories and `llms.txt` use the
reference registry rather than a separate list in page front matter.

Sitemap and metadata feeds exclude hidden pages and source-only pages without a
published URL. `llms-full.txt` applies the same visibility rule to its page index.

CLI layouts render upstream command data. API references have two rendering paths:

- `layouts/api-docs.html` renders the generated API catalog, overviews, operations,
  and schemas from committed `data/api-reference.json`. API operation navigation
  comes from that model. Generation inputs and commands are documented in
  [the API build guide](hack/api-docs/README.md).
- `layouts/api.html` renders legacy Engine versions as standalone ReDoc documents.
  It loads pinned ReDoc, owns its navigation and styles, and bypasses the shared site shell.
  Keep this boundary until those references migrate to the generated format.

Engine v1.57 uses `api-docs` while its authored version page retains its title,
URL, and latest aliases. The model's `authoredOverview` flag tells the adapter to
skip overview creation; operations and schemas are still generated. The authored
page selects the shared overview layout and model ID as described in
[the API build guide](hack/api-docs/README.md#authored-overviews).

## Styles and assets

`assets/css/style.css` is the Tailwind v4 entry point, compiled by Hugo's
`css.TailwindCSS` pipeline. `docs-shell.css` defines shared layout geometry;
`components.css` styles reusable components. API references and chat have their
own scoped styles.

Trident tokens live in `assets/css/vendor/trident/`. `theme.css` applies the docs
adaptations for Onest, JetBrains Mono, reading widths, and neutral inline code.
Colors use the pinned Trident semantic roles: `primary` and `primary-foreground`
for action buttons, `primary-ink` for links, and status roles such as `info` and
`warning` for informational states. Surfaces,
text, and borders use their corresponding semantic roles. Token values follow
the active theme. Code blocks use `surface-raised`, including installation
commands and unrecognized-language blocks. `syntax.css` styles Chroma tokens
with Trident palette shades chosen for readable contrast in both themes.

Focus indicators use `--docs-focus-ring`: `brand-600` in light mode provides
contrast over the blue category wash; dark mode uses `control-focus`.

The `.prose` token mapping in `theme.css` sits outside Tailwind's layers so it wins
over Typography's default and inverted palettes. Typography supplies document
structure and spacing; Trident supplies the colors. Keep custom prose selectors
scoped away from `not-prose` components.

`theme.js` applies a stored light or dark preference before paint. Without an
explicit preference it follows the operating system; using the toggle persists
the selected theme.
`category-washes.css` selects broad header and footer gradients through the
body's `data-category`. An unset category uses the default composition. Washes
use CSS and don't require route seeds or browser generation.

`ui-icon.html` renders the local Trident icon set using canonical icon names.
`icon.html` delegates interface icons to that partial and resolves explicit SVG
paths for product logos and other local artwork. Icon names in content and
templates match the vendored filenames; no icon-library module mount or alias
registry is required. Generated vendor assets retain their source geometry and
manifests. Refresh instructions and pinned migration versions are
in [the migration notes](hack/redesign/README.md#theme-and-assets).

## Browser behavior

Hugo renders navigation links and native disclosures. JavaScript restores sidebar
scroll position per collection and collapses the mobile sidebar initially. A
breakpoint change reopens it on desktop so the navigation remains available when
its mobile disclosure toggle disappears. JavaScript doesn't build or filter the
navigation tree.

`assets/js/main.js` loads the shared browser modules. Alpine.js manages search
and chat state. The search dialog uses native dialog behavior; `docs-search.js`
loads Pagefind on demand and loads matches in batches without discarding later
results. Editing the query invalidates outstanding requests before the next
search starts. The homepage and dialog share the search composer and results
partials. Selecting Ask Docker passes the query to the Gordon store and
opens the chat drawer. Rebuilding HTML alone doesn't regenerate the Pagefind
index; build that index when checking search results.

Reading pages include `layouts/_partials/ask-dock.html`, a bottom question field
on desktop and a chat trigger on mobile. CSS sticky positioning keeps it within
the main content area and clear of the footer. It uses the existing Gordon store
and hides while search or chat is open. Discovery pages and standalone ReDoc
omit it; without JavaScript it remains hidden.

Control this component with `params.gordon.ask_dock` in `hugo.yaml`:

```yaml
params:
  gordon:
    ask_dock: true
```

Set it to `false` to omit both the desktop field and mobile trigger from the
rendered HTML. Search and the other chat entry points remain available. An
omitted setting also disables the dock.

The shared code-copy partial and module handle fenced blocks and installation
commands. `quick-install.js` selects platform examples. It doesn't determine
whether a user is eligible to install a product; installation links remain the
fallback when detection or JavaScript is unavailable.

## Validate changes

Run scoped Markdown lint for prose changes as described in [AGENTS.md](AGENTS.md).
For template and navigation changes, build with `--panicOnWarning` and check the
rendered pages, including mounted and generated content. Navigation registry
validation runs through the shared head partial during the build.

Run `docker buildx bake test-site` for Node and browser tests. The browser
suite in `hack/test/browser.spec.mjs` covers representative pages and shared
interactions. Build preparation and local commands are in
[CONTRIBUTING.md](CONTRIBUTING.md#testing). See [ACCESSIBILITY.md](ACCESSIBILITY.md)
for accessibility coverage and outstanding manual reviews. API generator and
output checks live under `hack/api-docs/`.
