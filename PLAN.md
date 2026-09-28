# Sandbox documentation migration plan

## Purpose

Organize Docker sandbox documentation around shared behavior, practical tasks,
agent configuration, kit authoring, and organization administration. Reduce
drift by giving each behavioral explanation one authoritative home. Preserve
complete workflows and existing public links throughout the migration.

This plan covers the full restructure discussed with the maintainers. It does
not implement content moves. Destination paths below are proposals to settle
in the mapping phase, rather than an approved redirect map.

## Editorial rules

- Concepts explain behavior and qualify differences by environment or
  interface. They do not contain four copies of a procedure.
- CLI task pages cover local and cloud together where practical. Explain a
  flag difference once, use tabs for equivalent alternatives, and use sections
  for different prerequisites or sequences. Split pages only when the
  workflows become substantially independent.
- Console and SDK guides retain interface-specific procedures. Do not require
  a matching page for every topic in every interface.
- Agent pages own agent-specific authentication, configuration, interaction,
  and limitations across verified environments. Generic secret management and
  sandbox lifecycle instructions belong in task guides.
- Kit concepts explain the model; authoring guides explain building and
  distribution; interface guides explain applying kits. Maintain one version
  and interface compatibility matrix.
- Developer controls belong in task guides. Organization policies, sign-in
  enforcement, audit administration, and SIEM belong in AI Governance.
- Local-only capabilities, such as upstream proxy configuration, have scoped
  task pages without empty cloud alternatives.

Short contextual explanations and launch examples can remain in procedures.
Link to the authoritative explanation for detailed behavior. Use includes only
for stable snippets that genuinely need repetition; do not assemble whole
procedures from shared fragments.

## Target organization

Navigation labels do not require separate pages for every leaf. Retain existing
reference and billing sources where practical and link to them.

```text
Docker Sandboxes
  Overview
  Choose an environment and interface
  Concepts
    Isolation and security
    Files, workspaces, and storage
    Lifecycle and persistence
    Secrets and credential injection
    Network access and policy evaluation
    MCP connections and gateways
    Kits, workloads, and mixins
  Agents
    Overview and verified availability
    One page per supported agent
    Shell
  Use the CLI
    Install and sign in
    Get started locally
    Get started in the cloud
    Create and manage sandboxes
    Run commands and attach sessions
    Work with files and repositories
    Configure credentials
    Configure network access
    Connect MCP servers
    Use kits and mixins
    Use environment files
    Expose services
    Manage storage and templates
    Move between local and cloud
    Automate sandbox workflows
    Connect editors and applications
    Local configuration
      Runtime settings and model providers
      Agent skills
      GPU and host hardware
      Upstream proxy and registry mirror
    Troubleshooting
  Use the Console
    Get started with Docker Agentic Platform
    Create and manage sandboxes
    Select and run kits
    Work with terminals, files, and services
    Manage secrets
    Configure network access
    Connect MCP servers
    Troubleshooting
  Develop with the API and SDK
    Get started and install
    Authentication and authorization
    Resources, endpoints, and asynchronous operations
    Launch kits and images
    Cookbook
    Errors, retries, and request limits
  Author kits
    Get started
    Version and interface compatibility
    Workloads, mixins, and kit sets
    Base images, credentials, network requirements, and setup
    Test, publish, and distribute
    Maintain and migrate earlier formats
  Cloud access and billing
    Activate access
    Compute sizes and account quotas
    Usage and billing
  Reference and releases
    CLI, REST API, SDK, and kit references
    Release notes
AI Governance
  Overview and supported products
  Subscription and prerequisites
  Set up organization governance
  Organization and team policies
    Scope and assignment
    Network, filesystem, and MCP controls
  Enforce organization sign-in
  Inspect and troubleshoot enforcement
  Audit logs
    Coverage, delivery, and retention
    Local collection, hosted events, and SIEM
  Governance API, MCP policy, and audit references
```

Proposed source roots:

| Destination | Proposed source location |
| --- | --- |
| Sandbox family landing page | Keep `content/manuals/ai/sandboxes/_index.md` |
| Shared explanations | `content/manuals/ai/sandboxes/concepts/` |
| Agent guides | Keep `content/manuals/ai/sandboxes/agents/` |
| CLI tasks | `content/manuals/ai/sandboxes/cli/` |
| Console tasks | `content/manuals/ai/sandboxes/console/` |
| API and SDK guides | `content/manuals/ai/sandboxes/api/` |
| Kit authoring | `content/manuals/ai/sandboxes/author-kits/` |
| Shared cloud account guidance | `content/manuals/ai/sandboxes/cloud/` |
| Organization administration | `content/manuals/ai/governance/` |

The proposed API root overlaps historical aliases such as
`/ai/sandboxes/api/concepts/`. Check existing aliases before assigning canonical
URLs, and remove any alias that becomes its page's canonical URL. Reusing the
cloud directory also requires mapping its existing CLI pages explicitly.

## Repository findings

The baseline contains 135 Markdown pages, including section indexes:

| Source | Pages | Migration role |
| --- | ---: | --- |
| `content/manuals/ai/sandboxes/` | 80 | Split shared concepts, CLI tasks, agent guides, authoring, and administration |
| `content/manuals/agentic-platform/` | 10 | Console procedures plus shared cloud account information |
| `content/manuals/ai/sandboxes-api/` | 45 | Docs-owned guides and an imported cookbook |

The sandbox section includes 20 governance pages, six cloud CLI pages, ten kit
pages, and 12 agent pages. The SDK cookbook contains 36 imported recipes and a
docs-owned index. Refresh these counts against the implementation base commit.

### Navigation and rendering

- `layouts/_partials/sidebar/sections.html` recursively walks Hugo sections
  from `.FirstSection`. Nesting Console and SDK content under Sandboxes requires
  the corresponding content hierarchy. Weights, section indexes, and
  `params.sidebar` control ordering, grouping, labels, and badges. Start with
  content and front matter changes; no sidebar rewrite is planned.
- `content/manuals/_index.md` separately defines product cards. Its group data
  also feeds `layouts/home.llms.txt`, so updating discovery requires more than
  moving pages.
- The actual shortcode directory is `layouts/_shortcodes/`. Tabs have HTML and
  Markdown implementations. The Markdown implementation emits every variant;
  ensure each variant remains clearly labeled and meaningful outside the UI.
- `layouts/_shortcodes/sandbox-auth.html` and its Markdown counterpart contain
  local authentication instructions. Audit their consumers before broadening
  agent or quickstart scope. Update both outputs if the component changes.
- Cookbook discovery uses `recipe-list`, `params.recipeCatalog`, recipe
  `sidebar.group` values, and `sidebar.activeChildrenOnly`. The recipe list reads
  the index's direct regular pages. Keep recipes as direct children unless
  changing and validating that behavior deliberately.
- `sitemap: false` hides content from search and discovery, not from direct
  access. It is not a general migration staging mechanism. Section values do
  not automatically hide descendants without a cascade.

### URLs and dependencies

- `hugo.yaml` strips `/manuals` from published URLs, preserves path case, and
  treats failed references as errors. Internal source references must retain
  the appropriate `/manuals/` path.
- `layouts/home.redirects.json` combines page aliases and `data/redirects.yml`
  into production redirect data. It rejects self-aliases. Check collisions
  across both sources; a successful build alone is not a complete collision
  audit.
- Preserve stable vanity URLs, including `/go/dap/`, `/go/mcp-policy/`, and
  `/go/sbx-mcp-gateway/`. Update their destination keys in `data/redirects.yml`
  rather than changing the vanity URLs.
- Search `data/whats-new.json`, `data/summary.yaml`, guides, subscription pages,
  front matter grids, shortcodes, layouts, and configuration for inbound links.
  What's new includes sandbox links with fragments.
- Architecture diagrams and images live with some sandbox content. Preserve
  their references when moving or splitting the owning pages.

### Imported content ownership

`hack/sandboxes/README.md` identifies the cookbook recipes and
`content/reference/api/sandboxes/api.yaml` as generated from `docker/sbx-api`.
The importer copies them byte for byte. Prose and example corrections belong
upstream, followed by an import using paired source and generated commits in
`hack/sandboxes/source.json`.

`hack/sync-sandboxes-docs.py` hardcodes the cookbook destination. Update that
destination and its README as part of moving the cookbook. Coordinate upstream
links, front matter, and aliases so later imports preserve the migration.
Do not hand-edit imported recipes to attach redirects that the next import
would erase. If upstream alias changes cannot land with the move, use the
docs-owned redirect mechanism and verify its generated output.

The cookbook index and the SDK overview, installation, authentication, concepts,
errors, and limits pages are docs-owned. Keep generated API reference routes
under `content/reference/`; `hack/api-docs/sources.json` and its renderer manage
those outputs. Likewise, link to CLI reference rather than relocating generated
pages or editing their YAML to make this migration work. The `sbx` renderer uses
`data/sbx_cli/`; inspect its source workflow if reference changes become needed.
Do not edit `_vendor/`, `data/cli/`, or the vendored Governance API spec.

## Phase 1: map content and verify behavior

- [ ] Record the implementation base commit and inventory every source page,
  including indexes, FAQs, release notes, includes, diagrams, and imported files.
- [ ] Give every page a disposition: keep, move, split, merge, or retire. Map
  substantial sections when one page feeds several destinations.
- [ ] Create the source-path, public-URL, alias, and fragment mapping before
  making moves. Record ownership, destination, validation status, and dependent
  inbound links. Resolve final slugs and navigation labels in this phase.
- [ ] Assign an authoritative page to each shared behavioral claim. Identify
  repeated explanations to remove or shorten in task guides.
- [ ] Verify disputed behavior against implementation, supported releases, or
  product-owner confirmation. Record evidence and version scope. Documentation
  disagreement alone is not evidence of a product difference.

Carry these established corrections and open questions into verification:

| Topic | Required treatment |
| --- | --- |
| Cloud secrets | Maintainer confirms CLI and Console share a store. Correct the separate-store claim; verify naming, selection, attachment, and SDK interoperability details. |
| Secret injection | Treat differing descriptions as inconsistent coverage until verified. Explain placeholder and destination matching precisely. |
| Kit compatibility | Verify accepted descriptor versions, OCI references, artifact preparation, conversions, and mixins for each interface. Console v3-to-v2 conversion is a hypothesis, not a settled fact. |
| Lifecycle | Explain execution, memory, disk, connections, expiration, and billing separately. Different emphasis is not necessarily contradictory behavior. |
| Network policy | Distinguish explicit deny from default deny and document how account, sandbox, kit, and organization rules combine. |
| Organization governance | Preserve the documented local-only scope. Cloud audit delivery does not establish cloud sandbox governance coverage. |
| Agent support | Verify each agent's environment, interface, authentication, and kit availability before extending local instructions to cloud. |
| MCP | Verify local registration versus cloud connection and authorization, creation-time configuration, later attachment, and policy coverage. |

Exit condition: every source has a destination or explicit retention decision,
and each pilot topic has enough verified behavior to write accurately. Unresolved
topics remain scoped and do not block unrelated work.

## Phase 2: pilot shared explanations

- [ ] Introduce shared secrets and lifecycle explanations while the existing
  interface sections remain usable.
- [ ] Revise their docs-owned procedures to link to those explanations, keeping
  enough context to complete a task without reading several prerequisite pages.
- [ ] Prototype one combined local/cloud CLI MCP page. Use sections for
  different workflows and tabs only for equivalent alternatives.
- [ ] Coordinate equivalent cookbook edits in `docker/sbx-api` and import the
  approved export. Track pending upstream work explicitly.
- [ ] Review the pilots in HTML and Markdown. Check whether readers can find
  their environment, follow the procedure continuously, and understand shared
  behavior without repeating explanations.

Exit condition: maintainers can assess real pages and settle the page-boundary
rules before applying them across the family. Adapt the remaining map based on
the pilot; do not require equal page counts across interfaces.

## Phase 3: migrate content in coherent batches

Each batch includes complete pages, inbound link repairs, and redirects. Keep
every intermediate state navigable; do not publish empty section scaffolding.

| Batch | High-level changes |
| --- | --- |
| CLI and concepts | Combine relevant parts of root usage, cloud usage, configuration, workflows, security, and architecture. Retain local/cloud quickstarts. Put host-only tasks in local configuration. |
| Agents and integrations | Broaden agent pages only where verified. Keep agent-specific guidance there and generic operations in task pages. Audit integration scope separately. |
| Kits | Separate concepts, compatibility, authoring, and interface usage. Preserve v2 guidance and verified v3 requirements. Avoid implying universal OCI-reference support. |
| Governance | Extract organization/team policy administration, sign-in enforcement, audits, and references. Keep individual CLI policies and network troubleshooting in task guides. Split mixed policy-concept content by audience. |
| Console and shared cloud information | Move Console procedures beneath Sandboxes. Extract shared access, quotas, and compute explanations; keep canonical subscription and billing sources. Retain the Docker Agentic Platform name in the entry page. |
| API and SDK | Move docs-owned guides and coordinate the cookbook/importer migration. Preserve recipe discovery, upstream ownership, and stable API reference routes. |

Distribute FAQ answers to their authoritative concept or task where appropriate;
retain useful FAQ entry points. Keep release histories identifiable by product
and release stream rather than merging unrelated version sequences.

## Phase 4: navigation and discovery

- [ ] Update `_index.md` pages, weights, `linkTitle` values, grids, and sidebar
  groups as each destination becomes complete. Keep sidebar labels under
  30 characters.
- [ ] Put AI Governance alongside Sandboxes in the AI and agents group. Remove
  redundant standalone Console and SDK navigation entries after their
  replacements are ready.
- [ ] Update the manuals cards and reference entry links. Preserve access to
  Console, SDK, and governance from the Sandboxes overview and relevant tasks.
- [ ] Scope experimental badges to applicable interfaces and capabilities;
  do not mark the entire family experimental by moving a badge upward.
- [ ] Review breadcrumbs, active-page highlighting, expanded ancestors, and
  cookbook navigation on desktop and mobile. Retain existing hidden-page rules.
- [ ] Inspect search, `sitemap.xml`, `metadata.json`, `llms.txt`, `llms-full.txt`,
  and page Markdown for canonical destinations and coherent descriptions.

## Redirect and link procedure

Apply this procedure to every batch, not as cleanup after the full migration:

1. Use the batch's old-to-target mapping to inventory inbound links across the
   repository. Separate path-only references from references with fragments.
   The batch diff identifies moved content, not all inbound link sources.
2. Add former published URLs to destination aliases, preserving older aliases.
   For imported files, coordinate upstream or use docs-owned redirect data.
   Reject collisions, self-aliases, loops, and unnecessary redirect chains.
3. Update path-only links in all editable sources, including relative links
   inside moved pages, grids, includes, and configuration. Update vanity-link
   destinations without changing their stable source URLs.
4. Resolve fragment-bearing links against the explicit mapping and verify
   generated heading IDs. Preserve old IDs where practical. Follow the IA
   migration skill for unresolved mappings; do not invent replacement anchors.
5. For splits and merges, choose a destination that preserves the old page's
   purpose. An HTTP redirect cannot select a destination from the browser's
   fragment. Preserve compatibility anchors or retain a useful transition page
   when one destination cannot serve established deep links.
6. Sweep again for old source paths and public URLs. Only intentional alias or
   redirect sources and explicitly recorded deferrals may remain. Search
   deleted slugs and heading text throughout content and YAML/TOML configuration.

Review generated redirect data as well as rendered pages. Include HTML and
Markdown URL behavior in the compatibility check; do not assume page aliases
also cover every alternate representation.

## Validation and completion

For each batch:

```console
$ npx --no-install rumdl fmt <changed-markdown-files>
$ scripts/lint.sh <changed-markdown-files>
```

If the pinned formatter is unavailable, run `npm ci`. Read Vale warnings and
suggestions as well as errors. Run Prettier on changed non-Markdown sources.

In a regular checkout or CI, run the release build and full validation,
including Vale separately:

```console
$ docker buildx bake
$ docker buildx bake validate
$ docker buildx bake vale
```

`docker-bake.hcl` does not include Vale in its `validate` group; the build
workflow runs it separately. In git worktrees, follow AGENTS.md: do not change
`hugo.yaml` to work around Git path resolution. Use scoped lint locally and a
regular checkout or CI for complete build, HTML/link, redirect, and vendor
validation. Inspect `.github/workflows/build.yml` for the actual CI checks.

Before each PR, review the full branch against its target using the
review-changes skill. Stage files explicitly and exclude dependency-install
noise. Cookbook destination changes also need an import reproducibility check:
the pinned import must not recreate the old tree or overwrite migration fixes.

Completion requires:

- Every baseline page and substantial section has an accounted-for destination.
- Shared behavior has one authoritative explanation; procedures retain complete
  task paths and explicit scope. No unsupported parity claims were introduced.
- Every migrated public URL and mapped fragment works or has a documented,
  maintainer-accepted exception. Inbound link sweeps are clean.
- Navigation, badges, redirects, HTML, Markdown, and discovery outputs pass
  review. Required CI checks pass without migration-related failures.
- Imported content remains reproducible, and unresolved upstream corrections
  are recorded rather than hidden by local patches.

After deployment, check representative old links, vanity URLs, deep links, and
the main user journeys on the live site. If a batch must be reverted, revert
its content, navigation, redirect, and importer changes together.

## Preventing drift after migration

Add a concise content-ownership convention to the repository instructions once
the pilot establishes real canonical pages: behavior changes update the concept
and affected procedures; agent changes update the agent page; kit compatibility
changes update the shared matrix. Include upstream cookbook ownership in that
convention so future edits survive imports.

Use the migration map as the implementation checklist and close it only after
the final URL and content audit. Product parity is not a prerequisite for
completion; accurate scope and consistent explanations are.

## Prototype implementation

A prototype applies the proposed hierarchy to the existing documentation.
See [PROTOTYPE.md](PROTOTYPE.md) for review entry points, completed checks, and
remaining editorial and product-verification work. The complete source mapping
is recorded in [migration-map.json](hack/sandboxes/migration-map.json).
