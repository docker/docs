# Agent Readiness Rubric

Score the site on a 100-point scale before normalization. If a criterion is
not applicable, remove its points from the denominator instead of treating
it as failed.

## Grade bands

- `A`: 90-100
- `B`: 80-89
- `C`: 65-79
- `D`: 50-64
- `F`: below 50

## Confidence levels

- `High`: sitemap available and at least 12 sampled pages across at least
  four page types
- `Medium`: six to 11 sampled pages, or weaker coverage of page types
- `Low`: fewer than six sampled pages, or homepage-biased sampling

## Foundational caps

Apply these after computing the raw score:

- No `sitemap.xml` and no `llms.txt`: maximum grade `C`
- Markdown delivery fails on most sampled pages and no usable alternate
  markdown path exists: maximum grade `D`
- Main content is missing from initial HTML on more than 25% of sampled
  pages: maximum grade `D`
- `robots.txt` blocks broad crawl access to the docs site and the block is
  not clearly intentional: maximum grade `F`

Optional manifest gaps alone must not drop a docs-only host below `B`.

## Categories

### 1. Discovery and policy - 15 points

- `5` `llms.txt` exists, is fetchable, and is useful for agent discovery
- `4` `sitemap.xml` exists and includes the main docs corpus
- `4` `robots.txt` is accessible and does not unintentionally block major
  crawl agents or search agents
- `2` curated bulk-discovery aid exists, such as `llms-full.txt` or an
  equivalent machine-readable catalog

When `llms.txt` exists, sample some URLs from it. Stale or misleading
discovery links should reduce this category even if the file itself exists.

### 2. Retrieval and markdown delivery - 25 points

- `8` `Accept: text/markdown` works on sampled pages or an equivalent
  negotiated markdown response exists
- `5` a stable direct markdown route works on sampled pages
- `5` page-level markdown hints, alternates, or UI actions point to a
  working markdown URL
- `4` markdown responses strip navigation chrome and preserve headings,
  links, and code blocks cleanly
- `3` HTML and markdown stay in parity across the sampled set

### 3. Structure and semantics - 20 points

- `6` sampled pages have one `h1` and a mostly consistent heading hierarchy
- `5` `main` or `article` marks the primary content and the content is
  present in the initial HTML
- `4` canonical tags and stable final URLs are correct
- `3` structured data such as breadcrumbs or article metadata exists where
  appropriate
- `2` headings expose stable anchors or deep-link targets, and the HTML title
  or H1 stays reasonably aligned with the markdown H1

### 4. Crawlability and delivery behavior - 15 points

- `5` crawl directives are sane for a public docs property
- `4` the site does not depend on client-side rendering to expose core
  content
- `3` cache and freshness signals are reasonable for bots, such as
  `ETag`, `Last-Modified`, or useful cache headers
- `3` redirect chains are short and predictable

### 5. Machine-readable surfaces - 10 points

- `4` API or reference sections expose OpenAPI, schema, or downloadable
  machine-readable assets where relevant
- `3` pages with interactive JavaScript reference UIs still provide a usable
  non-JS fallback such as markdown, YAML, or another directly linked asset
- `3` tool manifests such as MCP, plugin, or agent descriptors exist only
  when the audited host is actually meant to expose tools

### 6. Content legibility - 15 points

- `5` markdown is clean and low-noise rather than a dump of site chrome
- `4` headings and section intros are specific enough for retrieval and
  chunking
- `3` fenced code blocks are mostly language-tagged and remain copyable and
  interpretable
- `3` repeated banners, chat chrome, consent overlays, or other boilerplate
  do not overwhelm the main content

## Scoring guidance

Use the full category only when the signal is consistently good across the
sample. Partial credit is expected.

Examples:

- A sitewide `llms.txt` that exists but is stale or too shallow may earn
  partial credit rather than full credit.
- If markdown works only on some page types, score that criterion based on
  observed coverage instead of failing or passing it outright.
- If a working markdown route exists but the page advertises a dead
  alternate URL, deduct in markdown discoverability rather than in raw
  markdown availability.
- If `llms.txt` exists but points to stale, broken, or inconsistent paths,
  deduct in discovery rather than in core fetchability.
- If tool manifests are irrelevant to the host, mark them `N/A`.
- If a major page type is weaker than the rest of the site, note that
  explicitly instead of letting stronger page types hide it in the average.

## Reporting guidance

For every category, include one line that explains the score:

- what was tested
- what passed
- what limited the score

Use evidence from live fetches. Do not score from assumptions about the
framework or source repository.
