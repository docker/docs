# API reference build

This tool validates the Hub, DVP, and Registry OpenAPI sources in
`content/reference/api/` and generates presentation data for Hugo. Engine and
AI Governance retain their existing sources and renderers.

## Commands

Use the Go version declared in `go.mod`, Hugo, and the site's Node dependencies.
From the repository root:

```console
$ npm ci
$ ./hack/api-docs/run.sh test
$ ./hack/api-docs/run.sh generate
$ hugo server
```

For a static build with HTML/Markdown checks, run
`./hack/api-docs/run.sh build`. To serve that build on port 1314, run
`./hack/api-docs/run.sh serve`. Set `DOCS_URL` and `DOCS_PORT` when using another
address. Run generation again after changing a specification or source manifest entry.

Docker builds and Netlify deploy previews run generation before Hugo. Generated
data, validation reports, binaries, and local builds go under `tmp/api-reference/`.
Hugo reports an error if the generated data is absent.

## Processor inputs

`sources.json` registers input specifications and their product/manual
relationships. Both generation and the validation wrapper use this manifest.
Hugo reads the generated presentation data, rather than this file. Authentication
descriptions come from the specifications and linked guides.

`validation/` contains the Vacuum rules, known-issue baseline, and locked official
schema resources. `testdata/` contains validation fixtures. The presentation model
version and supported dialect are defined by the processor, not configurable
manifest fields.

## Pipeline

1. Parse the authoritative YAML with `libopenapi` and preserve source values and
   reference identities. Verify locked OpenAPI dialect resources.
2. Run Vacuum policy checks and JSON Schema validation of schemas and supplied
   examples, followed by documentation profile checks.
3. Generate presentation model version 1: operations, effective security and
   servers, parameters, media variants, examples, schema links, and provenance.
4. Render HTML and Markdown through the content adapter and `api-docs` templates.
5. Run the shared site flattening script, then check generated operation,
   parameter, media, schema, and link coverage.

The shared `hack/flatten-and-resolve.js` script resolves Markdown link destinations
relative to each original file, then moves `index.md` files to flattened paths.
It preserves code examples and other text. API pages use the same processing as
other pages; links that already use published URLs remain unchanged.

The published YAML URLs still serve the source files directly. There is no
conversion step, snapshot dependency, Node migration package, or source archive
in the build. The source diff and
[prototype PR](https://github.com/docker/docs/pull/26043) provide migration context.

## Validation baseline

`check` is strict and reports every documentation profile failure:

```console
$ ./hack/api-docs/run.sh check
```

Generation explicitly uses `--allow-known-issues`. The checked-in
`validation/known-issues.json` records 168 remaining issues: 160 for Hub and eight
for DVP. Registry has no remaining baseline entries. Entries match the entire source digest, diagnostic digest,
rule, and source pointer. Parse failures, unresolved references, and unsupported
features cannot be waived. Unrecorded diagnostics fail the build.

This baseline is review debt, not approval of API behavior. Strict validation
fails until the issues are resolved. Any source edit requires deliberate review
of the affected baseline entries; the build never refreshes them automatically.
Resolve the baseline before merging the implementation PR, then remove the
exception file and the generation override. A passing documentation check does
not settle the separate product decisions about authentication and response
contracts.

### Issue triage

The first pass resolved 121 of the original 289 entries:

| Straightforward fix | Entries resolved |
| --- | ---: |
| Parameter descriptions derived from endpoint context and existing documentation | 43 |
| Quote the Hub and DVP 2FA examples to match their string schemas | 2 |
| Complete request and response examples using documented fields and values | 76 |

Shared parameters and responses account for multiple entries. Examples use
existing field annotations where available; team requests and DVP namespace/year
responses use illustrative values. SCIM error examples contain the declared
schema identifier and status, without inventing service error messages. These
are documentation examples, not captured service responses. Example lookup also
needed a fix to retain annotations on intermediate schema references.

The remaining entries need investigation before selecting a correction:

| Area | Entries | Evidence or decision needed |
| --- | ---: | --- |
| Hub error responses | 128 | Review payloads for `Error`, `error`, `ValueError`, and `rpcStatus`, plus two responses without schemas. Establish which fields each error returns. Object schemas using `items` leave error-map values unconstrained. |
| Hub pagination | 10 | Resolve six null/string conflicts and supply four list examples. Confirm page-boundary values and keep counts, links, and result arrays consistent. |
| Hub repositories and tags | 8 | Confirm the immutable-tag regex contract and two repository examples missing required `user` and `permissions` fields. Complete five media examples after resolving those contracts. |
| Hub teams, members, and invitations | 8 | Check four team responses, the member response and list wrapper, the bulk-invitation wrapper, and CSV export. The export requires `Role` but defines `Permission`; its array schema also needs a representation suitable for CSV. |
| Hub personal token update | 1 | Confirm whether the update response returns or redacts the token. The shared schema has a token value, while the retrieval endpoint documents an empty string. |
| Hub SCIM | 5 | Verify list-envelope casing, service-provider capabilities, and update semantics before completing examples. The source uses `resources` in list responses and `enabled` in updates, while user objects use `active`. |
| DVP analytics | 8 | Obtain representative metadata, pull, and export-download payloads. Resolve overlapping month/week `oneOf` branches: neither branch requires its distinguishing property. |

Counts include both missing examples and example/schema mismatches. The exact
locations remain in `validation/known-issues.json`; `check` writes diagnostics to
`tmp/api-reference/validation.json`. Repeated error responses are the largest
group, so investigate their shared schemas first. After each group is resolved,
validate its examples and remove its baseline entries. Do not replace a
conflicting example merely to make it pass, or loosen a schema without evidence.

## Tests and scope

Go fixtures cover dialects, references, recursion, boolean schemas, examples,
security overrides, server and parameter precedence, and request generation.
`verify-output.mjs` checks all 181 generated HTML/Markdown pairs and retention of
Engine v1.40–v1.56 in ReDoc, unchanged Governance rendering, and byte-identical
published specifications. `browser-checks.mjs` exports a Playwright check for
navigation, page aliases, filtering, requests, and narrow screens.

Callbacks and webhook navigation are unsupported and fail validation. Request
examples are POSIX shell templates; they do not make service calls. Specification
conversion and source-owner adoption remain separate from page rendering.
