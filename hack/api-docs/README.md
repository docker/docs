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

The published YAML URLs still serve the source files directly. There is no
conversion step, snapshot dependency, Node migration package, or source archive
in the build.

## Validation

`check` is strict and reports every documentation profile failure:

```console
$ ./hack/api-docs/run.sh check
```

Generation runs the same strict validation. Any diagnostic fails the build;
there is no exception baseline. Reports are written to
`tmp/api-reference/validation.json`.

## Tests and scope

Go fixtures cover dialects, references, recursion, boolean schemas, examples,
security overrides, server and parameter precedence, and request generation.
`verify-output.mjs` checks all generated HTML/Markdown pairs and retention of
Engine v1.40–v1.56 in ReDoc, unchanged Governance rendering, and byte-identical
published specifications. `browser-checks.mjs` exports a Playwright check for
navigation, page aliases, filtering, requests, and narrow screens.

Callbacks and webhook navigation are unsupported and fail validation. Request
examples are POSIX shell templates; they do not make service calls.
