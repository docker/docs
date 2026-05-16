# Agent Readiness Report Template

Use this structure for final audit output.

```markdown
## Agent Readiness Audit

**Site:** <base-url>
**Date:** <YYYY-MM-DD>
**Overall score:** <score>/100
**Grade:** <A-F>
**Confidence:** <High|Medium|Low>

### Summary

<2-4 sentence verdict focused on what an external agent can actually
discover, fetch, and interpret on this site.>

### Category Scores

| Category | Score | Notes |
| --- | ---: | --- |
| Discovery and policy | <x>/<y> | <short note> |
| Retrieval and markdown delivery | <x>/<y> | <short note> |
| Structure and semantics | <x>/<y> | <short note> |
| Crawlability and delivery behavior | <x>/<y> | <short note> |
| Machine-readable surfaces | <x>/<y> | <short note or N/A> |
| Content legibility | <x>/<y> | <short note> |

### Sample

- Sample strategy: <sitemap / internal links / explicit URLs>
- Sampled pages: <count>
- Page types covered: <landing, guide, manual, reference, ...>
- Weakest page type: <if any>

### Findings

- `P0`: <highest-priority blocker with evidence>
- `P1`: <important recurring issue with evidence>
- `P2`: <lower-priority or optional improvement>

### Remediation

- `P0`: <fix>, because <why it matters to agents>
- `P1`: <fix>, because <why it matters to agents>
- `P2`: <fix>, because <why it matters to agents>

### Evidence

- Sitewide checks: <llms.txt, robots.txt, sitemap.xml, manifests>
- Fetch-path checks: <markdown negotiation, direct markdown routes,
  advertised alternates, parity>
- Structural checks: <h1/main/article/canonical/json-ld/title-h1 parity>
- Code block checks: <fence count, language-tag coverage>
- Scanner comparison: <optional>
```

## Notes

- Keep the summary short and outcome-oriented.
- Findings should refer to concrete URLs or page types.
- If a criterion is `N/A`, say why instead of leaving it blank.
