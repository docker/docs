# API reference build

This tool validates the Hub, DVP, Registry, AI Governance, and Sandboxes OpenAPI sources
in `content/reference/api/` and generates presentation data for Hugo. Engine
v1.57 uses the same presentation pipeline with its upstream specification;
older Engine versions retain their existing sources and ReDoc renderer.

## Commands

Site builds use the committed `data/api-reference.json`. To build the site,
use Hugo and the site's Node dependencies. From the repository root:

```console
$ npm ci
$ hugo server
```

For a static build with HTML/Markdown checks, run
`./hack/api-docs/run.sh build`. To serve that build on port 1314, run
`./hack/api-docs/run.sh serve`. Set `DOCS_URL` and `DOCS_PORT` when using another
address. These commands use the committed data without running the generator.

After changing a specification, source manifest entry, or generator code, use
the Go version declared in `hack/api-docs/go.mod` to regenerate the data.

Install the Vacuum version pinned in the [Dockerfile](../../Dockerfile) and make
the `vacuum` executable available on `PATH`. The script runs Vacuum and the Go
generator without installing tools.

```console
$ ./hack/api-docs/run.sh test
$ ./hack/api-docs/run.sh generate
```

Commit `data/api-reference.json` with the source changes. With
`docker compose watch`, the regenerated JSON syncs to the server and Hugo
rebuilds the reference. Docker builds and Netlify deploy previews also use the
committed JSON. Validation reports and local builds remain under
`tmp/api-reference/`.

## Processor inputs

`sources.json` registers input specifications and their product/manual
relationships. Both generation and the validation wrapper use this manifest.
Hugo reads the generated presentation data, rather than this file. Authentication
descriptions come from the specifications and linked guides.

Optional source metadata controls routing and overview ownership:

| Field | Purpose |
| --- | --- |
| `url` | Override the overview URL and the base route for operations and schemas. |
| `sourceURL` | Override the published specification link without changing the input `source`. |
| `authoredOverview` | Set to `true` to render an authored overview instead of creating one in the adapter. |

`validation/` contains the Vacuum rules and locked official
schema resources. `testdata/` contains validation fixtures. The presentation model
version and supported dialect are defined by the processor, not configurable
manifest fields.

### Authored overviews

Set `authoredOverview: true` in a source manifest entry when an authored page
owns the overview at that entry's `url`. The generator carries the flag into
the presentation model. Hugo skips overview creation for that entry and still
generates all operations and schemas. Without the flag, Hugo creates the
overview too.

On the authored page, select `api-docs`, set `view` to `overview`, and set
`apiID` to the manifest entry's ID:

```yaml
layout: api-docs
params:
  view: overview
  apiID: engine-1.57
```

Keep the page's title, URL, and aliases in its front matter. Engine v1.57 uses
`content/reference/api/engine/version/v1.57.md` to preserve its published path
and latest aliases. Hugo does not merge authored front matter into an overview
created by the adapter at the same content path.

The manifest's optional `url` and `sourceURL` overrides control the overview
route and published specification link, respectively. They do not transfer
ownership of the overview or its aliases to the adapter.

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

The published YAML URLs still serve the source files directly.

## Validation

`check` reports contract errors and example warnings:

```console
$ ./hack/api-docs/run.sh check
```

Generation runs the same validation. Schema, reference, dialect, and documentation
profile errors block generation; there is no exception baseline for these errors.
Invalid example values produce diagnostics with severity `warning` and do not
block generation. Reports include both errors and warnings and are written to
`tmp/api-reference/validation.json`.

The presentation model omits failed media and schema samples from the examples
used for rendering and cURL generation. Raw upstream annotations and published
specification downloads remain unchanged. Example warnings do not relax the
supported contract or dialect limits.

The build workflow's validation matrix runs API reference checks alongside the
other validation targets. It runs tests, validates the sources, regenerates the
JSON, and compares the result with the committed file to detect stale data.
This matrix entry is advisory: failures do not block site builds or deployment.

Run the same check locally with Docker:

```console
$ docker buildx bake validate-api-reference
```

This target runs generator tests, strict validation, and a byte-for-byte
comparison inside Docker. It leaves the working tree unchanged and is also
included in the `validate` Bake group.

## Tests and scope

Go fixtures cover dialects, references, recursion, boolean schemas, examples,
security overrides, server and parameter precedence, and request generation.
`verify-output.mjs` checks all generated HTML/Markdown pairs, retention of
Engine v1.40–v1.56 in ReDoc, and byte-identical published specifications. It
also checks local links and fragments from API pages, including links to
`docs.docker.com`. The site `htmltest` checks include migrated Hub, DVP,
Registry, and AI Governance references. The shared browser suite checks API
operation filtering and request copying. Run it with `docker buildx bake test-site`;
see [Testing](../../CONTRIBUTING.md#testing) for local setup.

Callbacks and webhook navigation are unsupported and fail validation. Request
examples are POSIX shell templates; they do not make service calls.

Examples are optional. Supplied examples are checked against their schemas;
values that fail validation produce warnings and are omitted from rendered
samples. When a request body has no valid example, the generated cURL command
reads from `request-body` and directs the reader to prepare that file from the
schema.
