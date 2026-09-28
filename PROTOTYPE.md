# Sandbox documentation prototype

This branch restructures the existing documentation for review. Start at
`/ai/sandboxes/` in a local preview. The navigation uses Hugo's existing sidebar;
no custom navigation templates are required.

The prototype is rebased onto upstream main at `f22c0e6595`. It retains the
upstream MCP Toolkit move under Desktop, refreshed sandbox recipes and release
notes, and the committed API presentation-data build. The generated API data
has been regenerated for the migrated manual links.

## Implemented

- One sandbox family entry contains shared concepts, agents, CLI tasks,
  Console tasks, API and SDK guidance, kit authoring, cloud access, and
  reference material.
- AI Governance has a separate administrator entry. Personal local network
  policy instructions live under CLI tasks. Cloud account policies remain
  distinct from organization governance.
- Shared pages explain secrets, lifecycle, files, network access, and MCP.
  Kit compatibility and cloud quotas each have a dedicated source.
- The sandbox management guide combines local and cloud commands with
  synchronized tabs. The MCP guide uses sections for different setup workflows
  and tabs for the shared attachment step.
- Upstream proxies, registry mirrors, GPU passthrough, model settings, and
  host agent skills are grouped under local configuration.
- Claude Code and Codex pages include cloud CLI setup supported by the existing
  cloud documentation. Other agent pages retain their documented local scope;
  the agent index distinguishes guide coverage from product availability.
- Console and API experimental badges remain attached to those interfaces.
  CLI cloud guidance retains its own experimental notice.

The [migration map](hack/sandboxes/migration-map.json) accounts for all 135
original Markdown files: 121 moved and 14 kept at their source paths. Two
moved files are navigation links with `render: never`, rather than published
articles. Existing diagrams and images remain at their source locations.

## CLI navigation

The CLI entry has six expandable task sections and a troubleshooting page:

```text
Use the CLI
  Get started
    Install
    Local quickstart
    Cloud quickstart
    Local and cloud differences
  Run and manage
    Local operations
    Cloud operations
    Move a sandbox
  Access and connections
    Local and cloud credentials
    Authenticate command-line tools
    Local and cloud network access
    MCP servers
  Customize environments
    Kits
    Environment files
  Development workflows
    Git
    Local development
    Automation
    Editors and apps
  Local configuration
    Settings, models, skills, GPU, proxies, and registry mirrors
  Troubleshooting
```

The Run and manage section opens directly to the combined local/cloud guide.
Published URLs still redirect to the final destinations, and page headings
retain their IDs. The unpublished prototype URLs are not added as historical
aliases.

## Review these pages

| Review question | Page |
| --- | --- |
| Does the family entry make the choices clear? | [Sandbox overview](content/manuals/ai/sandboxes/_index.md) |
| Is one explanation sufficient across interfaces? | [Secrets](content/manuals/ai/sandboxes/concepts/secrets.md) and [lifecycle](content/manuals/ai/sandboxes/concepts/lifecycle.md) |
| Do tabs keep an equivalent task readable? | [Manage sandboxes](content/manuals/ai/sandboxes/cli/manage/_index.md) |
| Do sections handle different prerequisites without duplicating the topic? | [Connect MCP servers](content/manuals/ai/sandboxes/cli/access/mcp.md) |
| Is the administrator boundary clear? | [AI Governance](content/manuals/ai/governance/_index.md) |
| Does the kit split expose the real compatibility limits? | [Kit compatibility](content/manuals/ai/sandboxes/author-kits/compatibility.md) |
| Does an agent page work across environments? | [Claude Code](content/manuals/ai/sandboxes/agents/claude-code.md) and [Codex](content/manuals/ai/sandboxes/agents/codex.md) |

## Deliberate limits

This is a structural prototype, with representative content consolidation.
It does not complete every editorial step in PLAN.md.

The detailed local and cloud operations, credentials, and network-policy
procedures remain separate under CLI tasks. They preserve complete workflows
and existing deep links while the combined management and MCP guides provide
concrete examples for reviewing page boundaries. Consolidate those remaining
procedures by task after reviewing the examples; do not create matching pages
merely to make the interface sections symmetrical.

Local architecture and security explanations retain their explicit local scope.
The shared concept index identifies that scope. Broadening them requires
verified cloud implementation details. Organization policy syntax and audit
administration retain their detailed source material in the administrator
section.

The shared secret store is a confirmed correction. Precise secret identifier
interoperability, placeholder matching, kit conversion behavior, cloud artifact
versions, and broader agent availability still need product verification.
The prototype does not claim that every interface supports the same options.
The compatibility page records the documented v2 API/SDK artifact path without
claiming that Console conversion behavior is established.

The authentication picker shortcode remains local-only and is used by the local
quickstart. Its scope has not been broadened implicitly.

## Links and generated content

Former published article URLs redirect to their replacements. Markdown `.md`
URLs have explicit redirects too. Historical heading IDs are retained, including
short sections that point to a shared explanation. The migration also updates
inbound links, manual cards, reference configuration, ownership rules, vanity
URL destinations, and the CLI release-note generator.

Cookbook recipes moved byte for byte. The importer writes to the replacement
directory and retains its pinned source/generated commit pair. Redirects for
these recipes are maintained in `data/redirects.yml`, outside imported files.

Upstream follow-ups remain:

- Three imported recipes link to the former SDK installation URL. That URL
  redirects correctly. Update those links in `docker/sbx-api`, then import the
  paired export.
- The vendored AI Governance OpenAPI introduction links to the former governance
  overview. Its redirect is preserved. Update the source in
  `docker/governor-services` before re-vendoring.

## Validation

- Hugo builds successfully with no reported path warnings.
- Original article heading IDs were compared against a build of the base
  commit. None were lost.
- All 238 expected HTML and Markdown redirects were checked. None shadow a
  canonical migrated article URL.
- Rendered links were checked across 158 pages, including migrated articles
  and sources with inbound links. No migration-related failures were found.
  An unrelated `#prerequisites` link from the security announcements page is
  also broken in the base commit and remains outside this change.
- The pinned import completed, and all 36 recipe files still match their
  original bytes. The import does not recreate the former cookbook directory.
- Scoped Markdown lint passes without errors. Existing Vale warnings in
  untouched prose on inbound-link pages remain; sandbox prototype prose
  has no warnings.
- Browser checks cover sidebar nesting and active states, cookbook navigation,
  synchronized local/cloud tabs, and a 390-pixel mobile viewport. Markdown
  output retains both labeled tab variants. Docs-owned pages use explicit
  source paths so their links also resolve in Markdown output.
- Discovery output contains the sandbox family and separate governance entry.
- API generation completed with zero diagnostics, and output verification
  passed for all 450 HTML/Markdown page pairs and API reference links.

The native Hugo build and Markdown flattening completed in this worktree.
Full Docker validation and deployment checks remain for CI. Local Pagefind
indexing failed because its allocator does not support this host's page size;
search is not available in this preview. Production analytics also report
localhost CORS errors in browser checks.

The imported recipes inherit 81 broken relative-link targets in generated
Markdown, also present in a build of the base commit. Their HTML links work.
The Markdown flattening step resolves these source-relative links against the
rendered leaf directory. Fixing that renderer behavior is separate from moving
the recipes; their imported source remains unchanged in this prototype.
