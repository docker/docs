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
address. Run generation again after changing a specification or catalog entry.

Docker builds and Netlify deploy previews run generation before Hugo. Generated
data, validation reports, binaries, and local builds go under `tmp/api-reference/`.
Hugo reports an error if the generated data is absent.

## Pipeline

1. Parse the authoritative YAML with `libopenapi` and preserve source values and
   reference identities. Verify locked OpenAPI dialect resources.
2. Run Vacuum policy checks and JSON Schema validation of schemas and supplied
   examples, followed by documentation profile checks.
3. Generate presentation model version 1: operations, effective security and
   servers, parameters, media variants, examples, schema links, and provenance.
4. Render HTML and Markdown through the content adapter and `api-docs` templates.
5. Preserve API Markdown examples during site processing and check generated
   operation, parameter, media, schema, and link coverage.

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
`known-issues.json` records 289 inherited issues: 268 for Hub, 17 for DVP, and
four for Registry. Entries match the entire source digest, diagnostic digest,
rule, and source pointer. Parse failures, unresolved references, and unsupported
features cannot be waived. Unrecorded diagnostics fail the build.

This baseline is review debt, not approval of API behavior. Strict validation
fails until the issues are resolved. Any source edit requires deliberate review
of the affected baseline entries; the build never refreshes them automatically.
The implementation PR tracks the product decisions that require confirmation
before merge.

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
