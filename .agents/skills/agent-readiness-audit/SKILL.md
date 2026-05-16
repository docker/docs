---
name: agent-readiness-audit
description: >
  Audit a documentation site for agent-friendliness: discovery, markdown
  delivery, crawlability, semantic structure, machine-readable surfaces,
  and content legibility. Use when asked to assess docs.docker.com or any
  docs site for AI/agent readiness, produce a scored report, compare with
  external scanners, or generate a remediation list. Triggers on:
  "audit docs for agent readiness", "how agent-friendly is docs.docker.com",
  "score our docs for AI agents", "review llms.txt / markdown / crawlability",
  "create an agent-readiness remediation plan".
argument-hint: "<base-url>"
---

# Agent Readiness Audit

Audit the live site, not the source tree alone. Prefer the same fetch path
an external agent would use in the wild: direct HTTP requests, sitemap
sampling, and page-level inspection.

Do not reduce the result to a homepage-only scan or a binary checklist.

## 1. Set scope

Use `$ARGUMENTS` as the base URL when provided. Otherwise infer the base
URL from context and state the assumption.

Decide whether the host being audited is:

- a docs-only host
- an app/tool host
- a mixed host

This matters for optional checks such as MCP, plugin manifests, or other
tool discovery files. Do not penalize a docs-only host for missing
tooling manifests that belong on a separate service.

For `docs.docker.com`, treat the public docs host as docs-only. Docker's
MCP server is published separately, so missing MCP files on the docs host
should be reported as `N/A`, not as a failure.

## 2. Gather sitewide signals

Always check these resources first:

- `/llms.txt`
- `/llms-full.txt`
- `/robots.txt`
- `/sitemap.xml`

Only check host-level tool manifests when the host is an app/tool host,
mixed host, or explicitly advertises them:

- `/.well-known/ai-plugin.json`
- `/.well-known/agent.json`
- `/.well-known/agents.json`

Use the bundled script for a baseline:

```bash
bash .agents/skills/agent-readiness-audit/scripts/baseline-probes.sh \
  "$ARGUMENTS"
```

The script produces baseline evidence only. You still need to interpret
what matters for a docs property and score it with the rubric.

For docs-only hosts, you may skip tool-manifest probes to reduce noise:

```bash
CHECK_TOOL_MANIFESTS=0 \
  bash .agents/skills/agent-readiness-audit/scripts/baseline-probes.sh \
  "$ARGUMENTS"
```

## 3. Sample representative pages

Use the sitemap when available. Do not rely on the homepage alone.

If `llms.txt` exists, sample some URLs from it as well. This helps catch
stale or misleading discovery surfaces that a sitemap-only sample would miss.

Sample at least 12 pages when the site is large enough, and cover multiple
page types:

- homepage or docs landing page
- section landing pages
- task guides
- product manuals
- reference or API pages
- tutorial or learning pages

If the sitemap is missing or unusable, discover pages through internal
links and note the lower confidence.

If the site has distinct delivery patterns, sample each one. For example:

- normal content pages
- generated reference pages
- versioned docs
- localized docs

## 4. Run fetch-path checks on each sample

For each sampled page, verify:

- HTML fetch status, content type, and final URL
- `Accept: text/markdown` behavior
- direct markdown route behavior such as `<page>.md` or another stable path
- page-level markdown alternate links and whether they actually resolve
- whether page actions such as "Open Markdown" agree with the working route
- whether the HTML title or H1 matches the markdown H1 closely enough for
  retrieval parity
- whether main content is present in the initial HTML
- redirect chain length and canonical URL consistency
- obvious chrome/noise in the markdown response

Do not assume a `.md` mirror exists just because another site uses one.
Verify the actual markdown path the site exposes.

Treat these as separate signals:

- negotiated markdown works
- a stable direct markdown URL works
- the page advertises the correct markdown URL

If the page advertises dead markdown alternates but a working markdown route
exists, do not fail markdown delivery outright. Score it as a discoverability
and consistency problem instead.

For API or generated reference pages, also verify whether a machine-readable
asset such as OpenAPI YAML is directly linked and fetchable.

## 5. Judge structure and legibility

Measure structural signals:

- exactly one `h1`
- sane heading hierarchy
- `main` and `article` presence where appropriate
- canonical tags
- JSON-LD or breadcrumb structured data
- stable anchors and deep-linkable headings

Also make a qualitative judgment about agent legibility:

- markdown strips site chrome cleanly
- headings are specific and task-oriented
- code blocks stay intelligible without client-side JS
- the page is not dominated by banners, injected chat, or nav noise

Measure code block labeling explicitly when code samples are common. A page
type with many untagged fenced blocks should lose points even if the prose is
otherwise clean.

For page types that intentionally render interactive UIs with JavaScript,
judge them separately from normal docs pages. If the HTML shell is thin,
check whether the page still provides:

- a fetchable markdown summary
- a directly linked machine-readable asset
- a usable non-JS fallback

## 6. Score with the rubric

Use [references/rubric.md](references/rubric.md).

Rules:

- score only what you verified
- mark non-applicable checks as `N/A`
- normalize the final score against applicable points only
- do not let optional manifest checks dominate the grade

Apply the foundational caps from the rubric. A site with broken discovery
or broken markdown delivery should not earn a high grade because it has
clean metadata.

Do not average away a weak page type. If one major page type, such as API
reference, is materially worse than the rest of the corpus, call it out as
the weakest segment and reflect it in the category notes.

## 7. Compare with external scanners when useful

If external scanner results are available, compare them to your live
findings. Treat them as secondary evidence.

If a scanner and the live fetch disagree:

- trust the live fetch
- report the mismatch explicitly
- explain whether the scanner is testing a different assumption

## 8. Produce a remediation list

Turn findings into a short backlog:

- `P0`: fetchability or discovery blockers
- `P1`: recurring structural or parity issues
- `P2`: polish, optional manifests, or low-impact enhancements

For each remediation, include:

- the failing signal
- why it matters to agents
- a concrete fix
- whether it is sitewide or page-type-specific

## 9. Report in a stable format

Use [references/report-template.md](references/report-template.md).

Always include:

- overall score and grade
- confidence level
- sampled URLs or sample strategy
- category scores
- highest-priority findings
- remediation backlog

## Notes

- Favor docs-delivery checks over marketing-site heuristics.
- Do not fail a docs host for lacking MCP or plugin manifests unless the
  host itself is meant to expose tools.
- Treat raw byte size as supporting evidence, not as a primary scoring input.
- Prefer short evidence excerpts and commands over long copied page text.
