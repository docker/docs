# API reference build

This tool validates the Hub, DVP, Registry, AI Governance, and Sandboxes OpenAPI sources
in `content/reference/api/` and generates presentation data for Hugo. Engine
retains its existing sources and renderer.

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
the Go version declared in `hack/api-docs/go.mod` to regenerate the data:

```console
$ ./hack/api-docs/run.sh test
$ ./hack/api-docs/run.sh generate
```

Commit `data/api-reference.json` with the source changes. With
`docker compose watch`, the regenerated JSON syncs to the server and Hugo
rebuilds the reference. Docker builds and Netlify deploy previews also use the
committed JSON. Validation reports, binaries, and local builds remain under
`tmp/api-reference/`.

## Processor inputs

`sources.json` registers input specifications and their product/manual
relationships. Both generation and the validation wrapper use this manifest.
Hugo reads the generated presentation data, rather than this file. Authentication
descriptions come from the specifications and linked guides.

`validation/` contains the Vacuum rules and locked official
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

The published YAML URLs still serve the source files directly.

## Validation

`check` is strict and reports every documentation profile failure:

```console
$ ./hack/api-docs/run.sh check
```

Generation runs the same strict validation. Any diagnostic fails generation;
there is no exception baseline. Reports are written to
`tmp/api-reference/validation.json`.

The independent API reference data workflow runs tests, validates the sources,
and regenerates the JSON when API sources, generator files, or the committed
JSON change. It compares the result with the committed file and reports stale
data. This check is advisory: failures do not block site builds or deployment.

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
Registry, and AI Governance references. `browser-checks.mjs` exports a
Playwright check for navigation, page aliases, filtering, requests, and narrow
screens.

Callbacks and webhook navigation are unsupported and fail validation. Request
examples are POSIX shell templates; they do not make service calls.

Examples are optional. Supplied examples must validate against their schemas.
When a request body has no example, the generated cURL command reads from
`request-body` and directs the reader to prepare that file from the schema.
