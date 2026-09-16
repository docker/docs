# API documentation inventory

This document records the API documentation in the Docker docs repository and
the facts that a unified documentation experience must preserve. It provides
a baseline for discussing the revamp, with HTTP APIs first and SDK
documentation reserved for follow-up.

Snapshot: September 7, 2026, at commit
`371255294cf8db2ee20206bf5cc9c165a3a1ea16`.

This is a repository audit. Statements about authentication, endpoints, and
compatibility describe the checked-in documentation and specifications.
Authenticated service behavior, private upstream sources, and the deployed
site were not tested. An absent specification means none was found in this
checkout or linked from the reviewed API guide; it does not establish that
the service has no specification.

## Most important findings

1. Discovery is incomplete. Five HTTP API families have specifications under
   Reference: Engine, Hub, Docker Verified Publisher analytics, Registry, and
   AI Governance. The [Reference landing page][reference-home] lists only the
   first four. Model Runner, Docker Agent's two HTTP servers, and Scout
   metrics have documentation in Manuals, including content mounted from
   upstream repositories.
2. Shared credentials do not imply a shared authentication procedure. Hub,
   AI Governance, and DHI document an exchange for a short-lived bearer token.
   Registry uses a challenge and a scoped registry token. Scout metrics
   documents a PAT used directly as a bearer token. SCIM has a separate
   provisioning token. Several descriptions of these flows contradict each
   other; see the [issue register](#contradictions-and-open-questions).
3. Engine belongs in a unified experience, with explicit connection and
   version selection. Its daemon access, local sockets, optional remote
   transport, version negotiation, and streaming operations need their own
   guidance. Model Runner and Docker Agent also run as user-operated
   services, so this distinction applies beyond Engine.
4. There are three presentation approaches: ReDoc, the Hugo renderer used
   by AI Governance, and handwritten Markdown. The specifications span
   Swagger 2.0, OpenAPI 3.0, and OpenAPI 3.1. The Hugo renderer has useful
   integration with the docs site, but needs broader format and operation
   support before it could serve all five specifications.
5. Content ownership is part of the design problem. Engine specifications,
   plugin protocols, and Docker Agent docs are imported from upstream. AI Governance is a
   verbatim copy from a private repository. Hub, DVP, and Registry specs live
   directly in this repository, without an upstream synchronization process
   found in the reviewed scripts and workflows.

## HTTP API inventory

These rows identify documentation surfaces, rather than distinct backend
services. Hub and AI Governance share a host; Docker Agent exposes two
different HTTP interfaces. Related protocols appear in the next table.

Published paths below are derived from Hugo configuration. Source links point
to files in this checkout. The `/manuals` prefix is removed from published
URLs but retained in Hugo source cross-references.

| ID | API and purpose | Deployment and documented address | Documentation and specification |
| --- | --- | --- | --- |
| H1 | Docker Engine: containers, images, networks, volumes, plugins, and Swarm resources | User-operated `dockerd`; HTTP over a local socket or a configured remote connection. Version prefix such as `/v1.56`. | [Overview][engine-home] at `/reference/api/engine/`; [version wrapper][engine-version] at `/reference/api/engine/version/v1.56/`; [vendored Swagger spec][engine-spec]. Connection guidance is in [daemon security][daemon-security] and [remote access][daemon-remote]. |
| H2 | Docker Hub: repositories, tags, organizations, groups, invitations, tokens, audit logs, and provisioning | Hosted at `https://hub.docker.com`; paths primarily under `/v2`. | [Reference wrapper][hub-page] at `/reference/api/hub/latest/`; [OpenAPI spec][hub-spec]; [changelog][hub-changelog] and [deprecations][hub-deprecated]. SCIM is a distinct protocol within this spec. |
| H3 | Docker Verified Publisher (DVP) Data API: discovery and export of image pull analytics | Hosted at `https://hub.docker.com/api/publisher/analytics/v1`. Authentication paths override the server to `https://hub.docker.com`. | [Reference wrapper][dvp-page] at `/reference/api/dvp/latest/`; [OpenAPI spec][dvp-spec]; [analytics guide][dvp-guide], [changelog][dvp-changelog], and [deprecations][dvp-deprecated]. |
| H4 | Supported Registry API for Docker Hub: manifests and blobs | Hosted at `https://registry-1.docker.io`, with repository paths under `/v2/{name}`. | [Reference wrapper][registry-page] at `/reference/api/registry/latest/`; [OpenAPI spec][registry-spec]; [authentication protocol][registry-auth]. The spec explicitly covers a Hub-supported subset of Registry HTTP API V2, not the complete OCI Distribution Specification. |
| H5 | Docker AI Governance: organization policy and rule management | Hosted at `https://hub.docker.com/v2`; resources under `/orgs/{org_name}/governance/policies`. | [Reference page][governance-page] at `/reference/api/ai-governance/`; [OpenAPI spec][governance-spec]. The [manuals navigation entry][governance-manual] points to this reference and does not render a second API page. |
| H6 | Docker Model Runner (DMR): inference and model management | User-operated service. Host TCP examples use `http://localhost:12434`; container addresses depend on Desktop or Engine. A Desktop Unix socket route is also documented. | [Handwritten reference][dmr-api] at `/ai/model-runner/api-reference/`. Covers OpenAI, Anthropic, and Ollama compatibility, Diffusers image generation, and native model management. No DMR OpenAPI file found in the checkout. |
| H7 | Docker Agent API server: sessions, agent execution, tool confirmation, and events | User-operated `docker agent serve api`; default `127.0.0.1:8080`, routes under `/api`. The guide also documents the control interface attached to `docker agent run --listen`. | [Vendored guide][agent-api] at `/ai/docker-agent/features/api-server/`. Endpoint tables, request examples, and server-sent events. No native API specification found in the checkout or linked from this guide. |
| H8 | Docker Agent chat server: OpenAI-compatible access to agents | User-operated `docker agent serve chat`; default `127.0.0.1:8083`, `/v1/models` and `/v1/chat/completions`. | [Vendored guide][agent-chat] at `/ai/docker-agent/features/chat-server/`. The guide documents an unauthenticated `/openapi.json` endpoint on the running server. That specification is not included in this checkout and was not fetched. |
| H9 | Docker Scout metrics exporter: organization vulnerability and policy metrics | Hosted at `https://api.scout.docker.com/v1/exporter/org/{org}/metrics`. | [Metrics guide][scout-metrics] at `/scout/explore/metrics-exporter/`. Documents Prometheus and Datadog integration, metric names, labels, and authentication. HTTP metrics exposition, with no OpenAPI file found. |

### Specification size and versions

Counts come from parsing the five specifications, counting each HTTP method
under `paths` once. They describe documented operations, not verified service
coverage. Historical Engine versions are not counted as separate APIs.

| API | Specification format | Declared `info.version` | Paths | Operations |
| --- | --- | --- | ---: | ---: |
| Engine v1.56 | Swagger 2.0 | `1.56` | 98 | 108 |
| Hub | OpenAPI 3.0.3 | `2-beta` | 35 | 54 |
| DVP | OpenAPI 3.0.0 | `1.0.0` | 12 | 12 |
| Registry | OpenAPI 3.0.3 | Missing | 4 | 11 |
| AI Governance | OpenAPI 3.1.0 | `1` | 4 | 8 |

Hub's total includes nine SCIM operations. DVP's total includes two Hub
authentication operations and 10 analytics/discovery operations. The
Governance routes are in a separate specification, not duplicated in the Hub
specification at this snapshot. None of these five specs uses external
`$ref` references to share schema definitions.

Engine has 17 local Markdown version wrappers, v1.40 through v1.56. The
imported API module also contains Swagger files for v1.25 through v1.56 and
older Markdown documents. [Hugo mounts][hugo-config] the YAML files and
changelog into the reference tree; a file in `_vendor` alone does not imply
a rendered reference page. The `/reference/api/engine/latest/` URL is an
alias on the v1.56 wrapper. [Hugo parameters][hugo-config] identify Engine
29.8.0 and API 1.56 for this snapshot.

### Related HTTP surfaces and protocols

| Surface | Evidence and role | Treatment for this exercise |
| --- | --- | --- |
| SCIM provisioning | Nine operations under `/v2/scim/2.0` in the [Hub spec][hub-spec], plus [provisioning setup][scim-guide]. Service metadata, schema discovery, and user operations use SCIM formats and a provisioning bearer token. | Include within H2, with a distinct protocol/authentication entry. Avoid counting it as an additional standalone spec. |
| Registry token service | [Registry authentication][registry-auth] documents `GET https://auth.docker.io/token`, a `WWW-Authenticate` challenge, and `service`/`scope` parameters. | Include as supporting authentication reference for H4. This is a different exchange from Hub's `/v2/auth/token`. |
| Engine plugin protocols | [Plugin API][plugin-api] specifies RPC-style JSON over HTTP, with `dockerd` calling plugin-provided endpoints using `POST`. [Volume][plugin-volume], [authorization][plugin-auth], and [logging][plugin-logging] pages describe protocols; the [network page][plugin-network] links to the upstream network protocol. | Include as an extension protocol family. Distinguish these daemon-to-plugin calls from the Engine API's plugin management endpoints. The specs are prose and examples, with upstream ownership in `docker/cli`. |
| Engine metrics | [Prometheus guide][engine-metrics] documents a separate metrics listener, using `127.0.0.1:9323` in the example. It warns that metric names may change. | Record alongside H9 as HTTP observability documentation, with its own format and stability expectations. |
| Hub webhooks | [Webhook guide][hub-webhooks] documents outbound HTTP POST delivery to a user-selected URL and an example JSON payload. It marks `callback_url` as unsupported. | Include as an event integration, distinct from request/response API operations. No standalone webhook schema found. |
| Docker Agent A2A | [Vendored A2A guide][agent-a2a] documents an HTTP server, agent-card discovery, JSON-RPC invocation, and configurable bearer authentication. | Record as HTTP protocol documentation adjacent to REST. A2A transport and protocol semantics need explicit treatment. |
| Docker Agent MCP over HTTP | [MCP mode guide][agent-mcp] documents `docker agent serve mcp --http`, default `127.0.0.1:8081`, and configurable `--auth-token` authentication. | Record as an HTTP protocol integration adjacent to REST. The same command also supports stdio, which is outside the HTTP inventory. |
| DHI GraphQL | [DHI API page][dhi-api] and [VEX query guide][dhi-vex] document `POST https://api.dso.docker.com/v1/graphql`, organization context, and `imagePackagesForImageCoords`. | Defer detailed GraphQL work as requested. No GraphQL schema, introspection JSON, or schema download link found in this checkout or these pages. Preserve the documented shared Hub token exchange in the authentication analysis. |

Searches also found service hostnames in [Desktop's allowlist][desktop-allowlist]
and integration prerequisites. Hostname mentions for services such as Offload
do not establish a documented public API contract. No separate endpoint
reference or spec was found for Build Cloud, Offload, or a general Scout REST
API beyond the surfaces above. Guides that build sample REST applications,
third-party API integrations, CLI references, and OpenAPI-import tools are
not additional Docker HTTP APIs.

The [MCP Gateway guide][mcp-gateway] also describes a user-operated Docker
protocol proxy and links to its upstream repository for complete configuration
documentation. It does not provide a standalone REST endpoint reference.

## SDK follow-up inventory

| SDK or library | Documentation location | Relationship to the HTTP inventory |
| --- | --- | --- |
| Engine Go and Python SDKs | [SDK overview][engine-sdk] and [examples][engine-examples]; full language references are linked externally. | Clients of H1. The SDK section also contains direct `curl` examples; preserve those in the HTTP-first work. |
| Docker Extensions SDK | [Interface reference][extensions-reference] under `/reference/api/extensions-sdk/`, plus [manuals][extensions-manual]. | TypeScript/JavaScript interfaces for Desktop, host, extension services, and Docker operations. `HttpService` is an SDK interface, not a hosted Docker REST API. |
| Compose Go SDK | [Compose SDK guide][compose-sdk] at `/compose/compose-sdk/`. | Embeds Compose operations and calls the Engine API. No separate Compose HTTP service is documented here. |
| Docker Agent Go SDK | [Vendored SDK guide][agent-sdk] at `/ai/docker-agent/guides/go-sdk/`. | Embeds the agent runtime; separate from H7 and H8's HTTP contracts. |
| Testcontainers libraries | [Testcontainers landing page][testcontainers] and language-specific guides, with external library references. | Container orchestration libraries using a Docker-compatible runtime. A follow-up should distinguish Docker-sponsored and community implementations. |

Model Runner and Docker Agent chat also show use of third-party OpenAI SDKs.
These are client compatibility examples for H6/H8, not additional
Docker-owned SDKs. Docker Agent's stdio ACP and MCP modes are protocol
integrations outside the HTTP inventory.

## Authentication comparison

Personal access tokens (PATs) and organization access tokens (OATs) are
credential types. A bearer header describes how a token is sent; it does
not identify its issuer, audience, permissions, or exchange requirements.

| Surface | Documented authentication | Boundary to preserve |
| --- | --- | --- |
| Hub | `POST /v2/auth/token` accepts `identifier` and `secret`, returning `access_token`; send that token as `Authorization: Bearer ...`. The identifier is a username for a PAT/password or an organization name for an OAT. | Token type, resource permissions, subscription, SSO, and endpoint support are separate concerns. See C2 for conflicting OAT prose. |
| AI Governance | Its spec points to the same Hub token exchange and prohibits using a PAT/OAT directly as the bearer token. | C1 records an inconsistent request field name. Its responses separately describe organization permissions, entitlement, and policy/rule limits. |
| DVP | The spec declares Hub JWT bearer authentication and duplicates `/v2/users/login` plus `/v2/users/2fa-login`. | A candidate for shared Hub guidance, but the token exchange and accepted credentials need reconciliation; see C3. Publisher access remains API-specific. |
| DHI GraphQL | The guide explicitly exchanges a PAT/OAT using Hub's `identifier` and `secret` fields. | Also requires `ctx.organization` in GraphQL variables. Shared authentication does not remove GraphQL authorization and error semantics. |
| Hub SCIM | The setup guide instructs administrators to copy a SCIM base URL and API token from Docker Home. The Hub spec declares `bearerSCIMAuth` separately. | Do not replace SCIM provisioning credentials with a generic PAT/OAT recipe. |
| Registry | Request receives a bearer challenge; acquire a token for the indicated service and repository scope, then retry. Anonymous token acquisition is documented for public pulls. | Registry tokens are scoped for registry access. The docs do not establish interchangeability with Hub API JWTs. |
| Scout metrics | Organization owner's PAT sent directly in the bearer header. The guide says no specific PAT permissions are required. | Direct PAT use is explicitly documented for this endpoint. Do not apply the Hub token exchange universally. |
| Engine | Local socket access is controlled by host permissions; remote protection includes SSH or mutual TLS. | `X-Registry-Auth` passes registry credentials to Engine for registry operations. It does not authenticate the caller to `dockerd`. |
| Model Runner | The OpenAI compatibility table says an API key is unnecessary and the `Authorization` header is ignored. | Address selection and listener configuration matter. Avoid generalizing this statement to every compatibility route without implementation verification. |
| Docker Agent | API server supports a configured `--auth-token`; chat server supports `--api-key`/`--api-key-env`. A2A also documents configured bearer authentication. | These are operator-configured credentials, not documented Hub token exchanges. The API guide says an attached `run --listen` control interface has no built-in bearer authentication. |

Sources: the API specifications and guides linked in the inventory, plus
[daemon security][daemon-security], [SCIM provisioning][scim-guide], and
[OAT documentation][oat-guide]. Authentication shared across products should
be organized by verified flow, with API-specific permission requirements.

## Other shared concepts and differences

| Concept | Shared opportunity | Differences evidenced in the repository |
| --- | --- | --- |
| Resource identity | Explain organizations, namespaces, repositories, tags, and digests once where definitions match. | Hub has namespace and organization routes; DVP scopes analytics to namespaces; Governance uses organizations and policy/rule IDs; DHI adds query context. Engine's local objects have their own identifiers. |
| Versions and compatibility | Give every API a place to describe version selection, support, and changes. | Engine versions depend on daemon/client compatibility. Hosted specs use different path versions and `info.version` values. DMR and Agent advertise compatibility with other API formats. An OpenAPI version is not an API product version. |
| Pagination and response envelopes | Use consistent headings and precise parameter explanations. | Hub commonly uses `page`/`page_size` with linked pages, while SCIM uses `startIndex`/`count`. Governance's list response wraps summaries in `data` without documented pagination parameters. DVP provides period-based discovery and links to downloadable analytics files. |
| Errors, limits, and retries | Offer an API-specific account of errors and retry behavior in a consistent location. | Engine describes a JSON `message`; Governance uses `error.code` and `error.message`; DHI can return field errors alongside `data`. Hub documents request rate-limit headers separately from pull limits. Scout documents a 60-minute metrics cache, which is a freshness concern. |
| Streaming and binary data | Ensure the reference can explain more than JSON request/response pairs. | Engine includes attach/log/event streams and archive transfers. Registry transfers manifests and blobs with media types and upload state. DMR and Docker Agent stream inference or runtime events. DVP analytics downloads use CSV. |

## Delivery and source ownership

| Content group | How documentation is produced | Change ownership |
| --- | --- | --- |
| Engine, Hub, DVP, Registry references | [`layout: api`][redoc-layout] renders a standalone ReDoc page using a YAML URL derived from the Markdown filename. The template loads ReDoc from a CDN `latest` URL. | Templates and local wrappers belong to `docker/docs`. Hub, DVP, and Registry YAML files live here. Engine specs are mounted from `github.com/moby/moby/api`, pinned to v1.56.0 in [go.mod][go-mod]. |
| AI Governance reference | [`layout: api-reference`][native-layout] reads an adjacent YAML page resource and renders it through Hugo templates and [partials][native-resolver]. | Template changes belong here. Spec changes belong in private `docker/governor-services`, at `governor-service-api/openapi.yaml`, then are copied using [the synchronization script][governance-sync]. Do not edit the local spec by hand. |
| Handwritten API guides | Ordinary docs Markdown, tables, examples, and links. | DMR, Scout, DHI, SCIM setup, and Hub webhooks live here. DMR's CLI reference is imported separately; its handwritten REST API page is maintained locally. |
| Docker Agent guides | [Hugo module mounts][hugo-config] expose upstream Markdown under Manuals. | `docker/docker-agent`, pinned to v1.126.0. [The sync workflow][agent-sync] updates the module and vendor snapshot. |
| Engine plugin protocols | Docker CLI documentation is mounted into `content/manuals/engine/extend/`. | `docker/cli`; the network protocol page also links to implementation documentation in `moby/moby`. Do not modify `_vendor` copies. |

The Markdown experiences differ substantially. The [ReDoc Markdown
template][redoc-markdown] provides the page title, wrapper content, and a
specification link, without the operation descriptions. The [Governance
Markdown template][native-markdown] expands authentication, operations,
parameters, responses, and schemas. Handwritten API pages use ordinary
Markdown output. This difference affects readers and tools consuming the
Markdown representation of an API reference.

ReDoc uses fragments such as `#operation/...` and `#tag/.../operation/...`.
The Hugo renderer uses `#operation-<operationId>`. A renderer migration
therefore needs a strategy for existing deep links as well as page URLs.

## Contradictions and open questions

These are findings to resolve before writing shared guidance or migrating
renderers. A documentation contradiction does not by itself establish which
behavior the service implements.

| ID | Finding and evidence | Required follow-up |
| --- | --- | --- |
| C1 | [Governance authentication][governance-spec], `components.securitySchemes.bearerAuth.description`, names a `password` field for `/v2/auth/token`. [Hub's operation][hub-spec] defines `identifier` and `secret`; [DHI][dhi-api] uses those fields. | Verify the request contract against the service owner and correct the Governance source upstream. Use the Hub operation as the reference contract when reconciling the docs. |
| C2 | [OAT's Hub API section][oat-guide] says to pass the token as a bearer token with an organization name as the username. [Hub authentication][hub-spec] says to exchange credentials first; [repository export][hub-export] demonstrates that exchange. The OAT wording leaves the raw credential versus exchanged token ambiguous. | Confirm direct-token support, if any, per endpoint. Name the credential and resulting token explicitly in shared instructions. |
| C3 | [DVP][dvp-spec] duplicates `/v2/users/login` without a deprecation marker. [Hub][hub-spec] marks that operation deprecated and directs readers to `/v2/auth/token`. | Confirm DVP's supported token flow and credential types, then replace duplicated guidance with an authoritative shared source. Retain any verified compatibility requirements. |
| C4 | [The Hugo renderer][native-layout], [its navigation][native-nav], and [Markdown template][native-markdown] enumerate only `get`, `post`, `put`, `patch`, and `delete`. Hub has three `HEAD` operations, Registry two, and Engine two. The renderer reads OpenAPI 3 `components`/`servers`, while Engine uses Swagger 2 `definitions`/`basePath`. | Treat the renderer as a candidate requiring expansion. Verify method coverage, schema handling, per-path server overrides, tag handling, and deep links before migration. DVP's authentication server overrides are a concrete case. |
| C5 | Structural checks found a missing [Hub][hub-spec] reference at `#/components/responses/team_repo`; [DVP][dvp-spec] declares security scheme `type: https`, path-level `security`, and a non-extension `features.openapi` field; [Registry][registry-spec] lacks required `info.version`. Registry authentication is described in prose but has no OpenAPI security scheme or requirements. | Resolve spec validity and machine-readable authentication coverage before relying on generated clients or a replacement renderer. These three specs are local; establish product reviewers and upstream authority. |
| C6 | The first row of [Engine's API version matrix][engine-home] labels the maximum version `1.55` for Engine 29.8 but links to v1.56. [Hugo parameters][hugo-config] and [the version spec][engine-spec] identify 1.56. | Reconcile the display label with the version source. Keep daemon release, supported minimum/maximum, latest alias, and specification selection coordinated. |
| C7 | [DMR's Anthropic table][dmr-api] lists `/anthropic/v1/messages`, while its examples call `/v1/messages`. Its OpenAI table uses `/engines/v1/...`. | Verify route aliases and document canonical paths plus supported alternatives. The differing examples alone do not establish broken endpoints. |
| C8 | No DMR/native Agent OpenAPI file or DHI GraphQL schema was found. Agent chat documents a runtime spec endpoint. No Hub/DVP/Registry spec synchronization workflow was found in the reviewed scripts. | Locate authoritative specs and maintainers before estimating conversion work. Fetch and inspect the chat server spec in a follow-up; keep DHI schema research deferred. |

Structural checks used `openapi-spec-validator` 0.9.0 after parsing YAML and
converting mapping keys to their JSON string representation. Governance had
no validation errors in that check. Hub reference resolution stopped at the
missing reference. Engine v1.56 also raised null/array example-validation
errors; these were not investigated further and belong upstream if confirmed.
The findings above are selected blockers, not a complete specification
conformance audit or evidence of runtime failures.

## Agreed direction and external dependencies

The following context was supplied after the initial inventory:

- API documentation should have two entry points: a common API catalog in
  the revamped Reference experience, and direct navigation between each
  product's manuals and its API reference.
- A separate information architecture effort is defining the site structure.
  That work will affect the API revamp, but URL and navigation restructuring
  are outside this exercise.
- The presentation will use a tailored system. Comparing established
  renderers against custom rendering is outside this work. Unification must
  also establish requirements for the source OpenAPI specifications.

These are user-supplied constraints. The presentation approach below is a
proposal for discussion, not an agreed implementation decision.

## Proposed presentation approach

This initial sketch is superseded by the research and recommendations in
[Unified API documentation architecture](API-DOCUMENTATION-ARCHITECTURE.md).
That proposal covers the tailored presentation pipeline, a common source
profile, and adoption by API owners. The [tooling review](API-DOCUMENTATION-TOOLING.md)
records the component comparison and tests behind the refined pipeline.
The [prototype findings](API-DOCUMENTATION-PROTOTYPE.md) record the subsequent
dry-run against copied sources. The inventory remains a record of the original
documentation; prototype corrections are not upstream resolutions of its issues.

Use a common documentation experience with explicit support for each API's
capabilities. The two entry points should lead to one canonical reference
for a given API version, with connections back to its product documentation.
The separate IA effort will determine the navigation and URL details.

Separate the presentation into three responsibilities:

| Responsibility | Shared behavior | API-specific content |
| --- | --- | --- |
| Documentation frame | Product identity, navigation, search, stable operation links, version context, and access to source specifications | Product relationships, available versions, and deployment type |
| Authored guidance | A consistent place for connection setup, authentication, a first request, workflows, and troubleshooting | Daemon sockets, credential exchanges, registry upload sequences, streaming protocols, and compatibility limitations |
| Generated reference | Operation signatures, parameters, request bodies, responses, schemas, and examples derived from the authoritative specification | Supported methods, media types, security requirements, schema dialect, and operation-specific server addresses |

Keep explanatory guidance adjacent to the specification, with explicit links
to operations. Reuse verified authentication guidance by flow while retaining
each API's authorization requirements. Avoid copying endpoint definitions
into a second manually maintained reference.

Prefer existing standards-aware parsing, validation, and reference resolution
over implementing those semantics in templates. Evaluate renderer components
against the site's delivery requirements before selecting a tool. Parsing in
Hugo or another build step does not establish complete rendering support.
If Swagger conversion is necessary, preserve the upstream source and validate
the conversion against representative operations. Do not assume converting
the format preserves every relevant detail.

The reference should produce meaningful HTML and Markdown from the same
source. JavaScript can enhance schema exploration and examples. Browser-based
request execution should be an optional capability: Engine sockets, remote
TLS, browser cross-origin restrictions, and streaming behavior make a
universal request console unsuitable as the core reading experience.

Use the repository's actual specifications as the evaluation corpus:

- Engine tests version selection, Swagger 2, large schema sets, streams,
  archives, and connection guidance
- Hub and Registry test authentication distinctions, `HEAD`, response headers,
  media types, pagination, and existing operation links
- DVP tests server overrides for authentication and downloadable analytics
- Governance tests OpenAPI 3.1 and integration with the existing docs layout
- Model Runner and Docker Agent guides test how references without a checked-in
  specification participate in the same experience

Evaluation must detect omitted operations, unresolved references, lost
security requirements, and missing request/response variants. Schema
composition, required properties, acceptance of `null`, and recursive references
need explicit checks. Preserve access to the source spec and expose
unsupported constructs during development instead of silently dropping them.

Use these cases to verify the tailored implementation and its source
profile. Assess semantic coverage, readable HTML and Markdown, accessible
schema navigation, search integration, stable links, performance, and
maintenance cost. Adopt the system incrementally across APIs. Detailed
GraphQL and SDK presentation remain follow-up work.

## Implications for the next discussion

A unified experience can use a common API catalog, navigation, reference
structure, source-download access, and Markdown delivery. Each API still
needs an explicit deployment type, connection procedure, authentication
flow, permission model, and version policy.

For Engine, that structure should make daemon selection, socket or remote
connection, API compatibility, and a first HTTP request visible before
endpoint browsing. Its direct HTTP examples should be discoverable without
requiring readers to enter the SDK section. Model Runner and Docker Agent
can use the same structure for user-operated services while retaining their
different addresses, credentials, and compatibility formats.

Use this inventory with the architecture proposal to assign specification
owners, reconcile the authentication contradictions, and plan adoption.
Final URL and navigation decisions remain with the separate information
architecture effort.

## Maintaining this record

Keep the H and C identifiers stable. When resolving a question, record the
verified behavior, evidence, and date in its row. Distinguish a source-doc
correction from runtime verification. Refresh the snapshot commit and
specification counts when the underlying content changes.

The inventory was built by searching `content`, `layouts`, Hugo module
mounts, `_vendor`, scripts, and workflows for API references, HTTP endpoints,
OpenAPI/Swagger declarations, GraphQL, and SDK documentation. Candidate pages
were read to distinguish Docker interfaces from sample applications and
third-party integrations. Specifications were parsed to inspect methods,
servers, security schemes, versions, and references.

[reference-home]: content/reference/_index.md
[engine-home]: content/reference/api/engine/_index.md
[engine-version]: content/reference/api/engine/version/v1.56.md
[engine-spec]: _vendor/github.com/moby/moby/api/docs/v1.56.yaml
[engine-sdk]: content/reference/api/engine/sdk/_index.md
[engine-examples]: content/reference/api/engine/sdk/examples.md
[daemon-security]: content/manuals/engine/security/protect-access.md
[daemon-remote]: content/manuals/engine/daemon/remote-access.md
[hub-page]: content/reference/api/hub/latest.md
[hub-spec]: content/reference/api/hub/latest.yaml
[hub-changelog]: content/reference/api/hub/changelog.md
[hub-deprecated]: content/reference/api/hub/deprecated.md
[hub-export]: content/manuals/docker-hub/repos/manage/export.md
[dvp-page]: content/reference/api/dvp/latest.md
[dvp-spec]: content/reference/api/dvp/latest.yaml
[dvp-guide]: content/manuals/docker-hub/repos/manage/trusted-content/insights-analytics.md
[dvp-changelog]: content/reference/api/dvp/changelog.md
[dvp-deprecated]: content/reference/api/dvp/deprecated.md
[registry-page]: content/reference/api/registry/latest.md
[registry-spec]: content/reference/api/registry/latest.yaml
[registry-auth]: content/reference/api/registry/auth.md
[governance-page]: content/reference/api/ai-governance/index.md
[governance-spec]: content/reference/api/ai-governance/api.yaml
[governance-manual]: content/manuals/ai/sandboxes/governance/reference/api.md
[governance-sync]: hack/sync-governance-api.sh
[dmr-api]: content/manuals/ai/model-runner/api-reference.md
[agent-api]: _vendor/github.com/docker/docker-agent/docs/features/api-server/index.md
[agent-chat]: _vendor/github.com/docker/docker-agent/docs/features/chat-server/index.md
[agent-a2a]: _vendor/github.com/docker/docker-agent/docs/features/a2a/index.md
[agent-mcp]: _vendor/github.com/docker/docker-agent/docs/features/mcp-mode/index.md
[mcp-gateway]: content/manuals/ai/mcp-catalog-and-toolkit/mcp-gateway.md
[agent-sdk]: _vendor/github.com/docker/docker-agent/docs/guides/go-sdk/index.md
[agent-sync]: .github/workflows/sync-docker-agent-docs.yml
[scout-metrics]: content/manuals/scout/explore/metrics-exporter.md
[scim-guide]: content/manuals/security/provisioning/scim/provision-scim.md
[oat-guide]: content/manuals/security/access-tokens/organization-access-tokens.md
[dhi-api]: content/manuals/dhi/tools/api.md
[dhi-vex]: content/manuals/dhi/how-to/vex-api.md
[plugin-api]: _vendor/github.com/docker/cli/docs/extend/plugin_api.md
[plugin-volume]: _vendor/github.com/docker/cli/docs/extend/plugins_volume.md
[plugin-auth]: _vendor/github.com/docker/cli/docs/extend/plugins_authorization.md
[plugin-logging]: _vendor/github.com/docker/cli/docs/extend/plugins_logging.md
[plugin-network]: _vendor/github.com/docker/cli/docs/extend/plugins_network.md
[engine-metrics]: content/manuals/engine/daemon/prometheus.md
[hub-webhooks]: content/manuals/docker-hub/repos/manage/webhooks.md
[extensions-reference]: content/reference/api/extensions-sdk/_index.md
[extensions-manual]: content/manuals/extensions/extensions-sdk/_index.md
[compose-sdk]: content/manuals/compose/compose-sdk.md
[testcontainers]: content/manuals/testcontainers.md
[desktop-allowlist]: content/manuals/desktop/setup/allow-list.md
[hugo-config]: hugo.yaml
[go-mod]: go.mod
[redoc-layout]: layouts/api.html
[redoc-markdown]: layouts/api.markdown.md
[native-layout]: layouts/api-reference.html
[native-markdown]: layouts/api-reference.markdown.md
[native-resolver]: layouts/_partials/api-ref/resolve.html
[native-nav]: layouts/_partials/api-ref/nav.html
