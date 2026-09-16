# API documentation tooling review

This review replaces the component recommendation in the first
[architecture draft](API-DOCUMENTATION-ARCHITECTURE.md). It compares open-source
processing tools for Docker's proposed OpenAPI 3.2 source contract and tailored
presentation system. It does not compare reference UIs.
The [prototype report](API-DOCUMENTATION-PROTOTYPE.md) adds implementation
evidence and links to the reproducible pipeline and migration records.

Research date: September 7, 2026. Versions below are the examined releases.
Experiments ran in temporary directories against shared fixtures; no product
specifications, dependencies, or site implementation were changed.

## Recommendation

Use a Go preprocessing tool built on [libopenapi][libopenapi], with
[vacuum][vacuum] for structural and Docker policy checks and
[jsonschema/v6][go-jsonschema] for explicit schema compilation and example
validation. Retain the original source nodes and reference graph throughout.
Hugo consumes prepared presentation data and publishes HTML, Markdown,
navigation, search content, and source downloads.

Make [oasdiff][oasdiff] a review-report component. Keep runtime testing with
[Schemathesis][schemathesis] in product CI, where controlled services are
available. Neither belongs in the core rendering algorithm.

Publish a locked source package with an entry specification, its referenced
files, and provenance metadata as the canonical downloadable contract.
A single-file bundle can be an additional artifact after its transformation
passes fidelity checks. A common OpenAPI format does not require a single
physical source file.

The earlier `@redocly/cli` recommendation was based on too narrow a positive test.
It was intended to supply linting and bundling, not the presentation layer.
The broader tests do not support using it as the central validation gate
for this profile. Its bundling command remains useful in a narrower role.

## Why separate the components

“OpenAPI validation” can refer to several different operations:

| Responsibility | Question answered |
| --- | --- |
| Strict parsing | Is the YAML/JSON well formed, without duplicate keys or discarded values? |
| Document structure | Are OpenAPI objects, required fields, and references valid? |
| Docker source policy | Are operation IDs, navigation groups, descriptions, ownership, and examples present as required by the Docker profile? |
| Schema compilation | Is each embedded Schema Object valid under its declared dialect, including schemas without examples? |
| Example evaluation | Does a supplied example satisfy its schema and the selected request/response policy? |
| Contract comparison | What changed, and what could affect consumers or documentation links? |

A tool can perform one of these correctly while leaving another unchecked.
A parser can also repair input, or a convenience model can lose information
when serialized. Neither behavior is suitable for an authoritative
publication boundary without explicit handling.

The serious findings came from negative controls and semantic comparisons.
Checking that `query` or `itemSchema` appears in output would not have caught
a false schema becoming an unconstrained schema, or one reference constraint
overwriting another.

## Available components

The licenses below apply to the examined open-source packages. The proposed
core functions do not require a hosted account or paid service.

| Component and examined release | License | Assessment for this pipeline |
| --- | --- | --- |
| [libopenapi][libopenapi] 0.38.7 | MIT | Preferred source model, operation index, reference indexes, and source locations. Retain original schema nodes; do not reconstruct the contract through its schema rendering methods. |
| [vacuum][vacuum] 0.30.3 | MIT | Preferred policy/checking interface, with configurable rules, Go/JavaScript custom functions, and CI diagnostics. Use an explicit Docker ruleset and severity policy. |
| [libopenapi-validator][lib-validator] 0.14.0 | MIT | Supporting structural/request/response validation library. Its document convenience method does not replace explicit compilation of every schema. Its schema helpers use the same Go JSON Schema implementation recommended below. |
| [jsonschema/v6][go-jsonschema] 6.0.3 | Apache-2.0 | Preferred JSON Schema 2020-12 implementation for the Go tool. Supports resource registration, schema compilation, and instance validation. Configure the chosen dialect and resource loader deliberately. |
| [kin-openapi][kin] 0.149.0 | MIT | An active alternative with better 3.2 support than Hugo's pinned 0.139.0. Tests still found valid tag metadata rejected and reference-sibling information lost on serialization. |
| [Scalar parser][scalar-parser] 0.29.1, [document checker][scalar-validator] 0.1.1, and [json-magic][scalar-bundle] 0.13.4 | MIT | Separate components with different contracts. Strict document checking and narrow bundling are useful. Lenient parser validation and reference expansion are unsuitable as the authoritative semantic model. |
| [Spectral][spectral] CLI 6.16.3 / rulesets 1.22.7 | Apache-2.0 | A viable general policy engine. The built-in OAS ruleset does not handle the tested 3.2 document correctly; a Docker-owned ruleset would be necessary. Vacuum fits the proposed Go integration better. |
| [Redocly CLI/core][redocly] 2.51.2 | MIT | Tested bundling preserves important cases. Its structural lint rule rejects a valid boolean Schema Object, so it should not be the default source gate for the admitted profile. |
| [Swagger ApiDOM][apidom] 1.12.1 | Apache-2.0 | A credible typed syntax-tree foundation for a TypeScript alternative. Its 3.2 namespace recognizes the fixture, but the inspected 3.2 bundle strategy does no bundling. Full parser/resolver integration needs further work. |
| [APIDevTools swagger-parser][swagger-parser] 13.0.0 | MIT | Explicitly rejects OpenAPI 3.2. Not a candidate for the selected source format. |
| [APIDevTools JSON Schema ref parser][ref-parser] 15.3.6 | MIT | A separate generic reference/bundling tool that can process the tested 3.2 document. Reference expansion is not a complete JSON Schema semantic resolver. |
| [Ajv][ajv] 8.20.0 | MIT | A strong JSON Schema option if the compiler uses TypeScript. There is no need to add it alongside the Go evaluator by default. |
| [Python spec checker][python-spec] 0.9.0 / [schema checker][python-schema] 0.9.0 | Apache-2.0 / BSD-3-Clause | Useful independent diagnostics. The schema evaluator works when invoked directly; document-level checks missed important streaming schema/reference cases. Adding Python to the core build is unnecessary. |
| [oasdiff][oasdiff] 1.31.0 | Apache-2.0 | Preferred compatibility and changelog reporting. Tests detect several 3.2 changes that other comparison tools missed. Review full reports and apply Docker-specific severity policy. |
| [OpenAPITools openapi-diff][java-diff] 2.1.7 | Apache-2.0 | The examined documentation advertises OAS 3.0, without verified support for this 3.2 profile. Do not select it for this pipeline. |
| [Schemathesis][schemathesis] 4.26.0 | MIT | Suitable for optional product conformance tests. Its 3.2 operation and synthetic response handling passed targeted tests, but live stream testing requires a bounded test harness. |

This comparison does not establish that an entire ecosystem is better than
another. The Go recommendation follows from the useful combination of source
locations, reference indexes, custom rule integration, and a schema evaluator
that can work with intact resources. The repository already has a Go build
environment through Hugo. A TypeScript alternative remains possible, but
switching among JavaScript convenience parsers would not remove the semantic
work identified here.

## What the tests established

The shared valid fixture combined four operations, including `QUERY`, an
external streaming `itemSchema`, a boolean schema, recursive references,
a `$ref` sibling, a path-level server override, an unauthenticated operation,
and false/zero example values. Separate fixtures introduced missing metadata,
an unresolved streaming reference, an invalid schema type, an invalid example,
a duplicate YAML key, and a missing operation ID.

Missing `operationId` is an intentional Docker profile test; OpenAPI itself
does not require every operation to have one.

### Checks need explicit coverage

| Tested path | Observation | Consequence |
| --- | --- | --- |
| Vacuum recommended rules | Accepted the valid fixture; rejected missing version, unresolved streaming reference, and duplicate key. An invalid example was a warning with exit zero. Invalid `itemSchema.type` was not rejected. | Select rule severities explicitly and compile schemas independently. The recommended ruleset alone is insufficient. |
| Vacuum custom operation-ID rule | A rule including `QUERY` passed the valid fixture and rejected a missing ID with a source location. | The Docker profile can be enforced through custom rules; its selectors must cover all admitted operation collections. |
| libopenapi-validator `ValidateDocument()` | Accepted the invalid streaming schema type. Explicit schema compilation rejected it before any example was checked. | A capable schema evaluator is ineffective if the document traversal never sends it a schema. |
| `@redocly/cli` configured structural/example rules | Detected the invalid metadata, reference, type, and example, but rejected the valid boolean schema as an object-type error. | Its structural rule cannot be adopted unchanged for this source profile. The earlier positive probe did not exercise this case. |
| Scalar parser validation | Replaced missing `info.version` with `0.1.0` and accepted the document. Its separate strict checker rejected the missing field. | Avoid normalization/repair APIs at the authoritative boundary. Choose the exact package and function by role. |
| Scalar strict document checking | Accepted invalid `itemSchema.type`; the structural schema admits an object or boolean without validating its schema keywords. | Document checking and embedded schema compilation remain separate requirements. |
| Spectral built-in OAS rules | Treated the 3.2 fixture as 3.0 and rejected legal fields. A standalone custom rule worked. | Spectral is an alternative policy engine, not a verified 3.2 structure gate with its default OAS rules. |
| Python spec checking | Rejected missing version, but accepted the missing streaming reference and invalid streaming schema. Direct schema evaluation handled the tested JSON Schema constraints correctly. | Do not infer full schema/reference coverage from a claimed specification version. |

The build should not accumulate all these tools as redundant gates.
Vacuum, `libopenapi`, and the supporting Go validation library share
implementation dependencies; their agreement is not independent evidence
of correctness. Use focused regression fixtures and separate responsibilities.

### Preserve the graph and source nodes

Both vacuum's default inline bundling and composed bundling changed
`schema: false` to `{}`. The high-level schema `Render()` path did the same.
The first schema accepts no instances; the second accepts any instance.
This is a semantic change, not formatting.

However, `libopenapi` retained the literal false value in
`SchemaProxy.GetValueNode()` and in its source JSON representation. An
experiment registered intact resources with `jsonschema/v6`, selected schemas
by URI/pointer, and rejected unregistered resource access. It correctly
compiled recursive schemas, resolved the external streaming schema, rejected
the invalid type, and preserved false-schema validation.

That is the recommended processing path: use typed operations as convenient
views, with original schema nodes and reference identities as authority.
Do not serialize a convenience model back into the canonical specification.

The assembled prototype found a further limit: the convenience source JSON
view normalizes YAML timestamp strings, including fractional-second spelling.
The implementation reads original YAML scalar nodes for authoritative values
and uses the typed library model for operation source locations. The earlier
false-schema probe established preservation of that specific value, not every
source scalar.

Reference expansion also failed in other components. Scalar and the generic
APIDevTools reference parser lost the base constraint in this case:

```yaml
Base:
  type: object
  required: [a]
Sibling:
  $ref: '#/components/schemas/Base'
  required: [b]
```

The expanded object retained only `required: [b]`; both constraints must
apply. Additional named-anchor tests exposed incorrect or incomplete
resolution. ApiDOM's inspected [schema expansion implementation][apidom-expand]
also lets sibling keywords override constraints from the referenced schema.
A fully expanded object should therefore not
be the compiler's canonical schema representation.

### Bundling can remain optional

`@redocly/cli` bundling, Scalar `json-magic` bundling, and the generic APIDevTools
bundler preserved the tested false schema, recursive references, reference
siblings, `QUERY`, and external streaming schema. This is useful evidence
for an optional single-file derivative, rather than a reason to adopt their
other APIs as the compiler.

There are still conditions to verify: reference errors must fail the task,
network resolution must follow the dependency lock, and resource/anchor
identity must survive. Scalar's default error handling can log unresolved
references and return output, and object inputs can be mutated. ApiDOM's
3.2 [bundle strategy][apidom-bundle] returns its parsed input; its existence is not
evidence that it builds a portable bundle.

Retaining the source package avoids making any such rewrite a mandatory
step. The package must contain the entrypoint and all locked resources in
a layout where their references resolve. Publish that tree at stable URLs
and provide an archive for download. Acquiring dependencies and preserving
their identities remains required; an archive alone does not fix references
that point outside the package.

If a single-file download is needed, introduce one of the passing bundling
paths behind the same semantic fixtures. `@redocly/cli` is a reasonable candidate
for that limited role. It should not supply canonical rendering input or
be mistaken for the source validation system.

### Compare contracts with appropriate limits

`oasdiff` detected removal of a `QUERY` operation, a boolean schema changing
from true to false, and an externally referenced streaming schema changing
to an incompatible inline type. The `libopenapi` method `CompareDocuments` missed the
latter two cases in the tested end-to-end calls.

Use the full diff/changelog report from `oasdiff` for review. Its breaking report
treated an operation changing from `security: []` to bearer authentication
as informational, so Docker needs an explicit effective-authentication
comparison and severity policy. Its `breaking` command also exits zero by
default even when findings are present; enforcement requires a configured
`--fail-on ERR` or `WARN` policy.

The [documented caveats][oasdiff-caveats] include dynamic references and
component path items. Do not treat an empty report as proof of complete
equivalence. Avoid using its schema-flattening options as compiler input;
[documented flattening limitations][oasdiff-flattening] can change constraints.
Compare against the previous accepted source for the same API version, or
an explicitly selected release baseline. Intentional versioned contract
changes need review, rather than an unconditional prohibition on change.

## Proposed components and pipeline

| Stage | Component | Output and boundary |
| --- | --- | --- |
| 1. Acquire and lock | Repository import tooling | Immutable source tree, dependency lock, API manifest, source revisions, and digests. Only declared resources are available to later stages. |
| 2. Parse and index | `libopenapi` | Strictly parsed source nodes, operation model, reference indexes, and source locations. Keep duplicate-key detection enabled and restrict indexed directories. |
| 3. Check document and profile | Vacuum plus an explicit Docker rule package | Structural/profile diagnostics with source pointers. Required rules fail as errors; no generic quality score or silent repair. |
| 4. Compile schemas | `jsonschema/v6` | A compiled schema for every admitted schema position, including `itemSchema`, whether or not it has an example. Register intact resources and the exact chosen dialect. |
| 5. Check examples | The same compiled schemas plus Docker request/response policy | Errors tied to the example and schema locations. No value coercion, inserted defaults, or removed properties. Handle media-specific examples explicitly. |
| 6. Prepare reference data | Docker-owned Go compiler | Versioned operation/navigation data, effective settings, schema graph, guide links, source provenance, and coverage diagnostics. |
| 7. Review changes | `oasdiff` plus prepared-data comparisons | Full contract report and Docker-specific checks for effective security, stable IDs, response/media coverage, and unsupported features. |
| 8. Publish | Hugo content adapters and templates | HTML, Markdown, source tree/archive, navigation, search, and stable links. Optional separately checked single-file bundle. |

Stages 2–6 can live in one Go command and share loaded resources. Source
repositories can run its check-only mode with the same profile release;
the docs import runs those checks again on the exact revision it publishes.
Pin compatible library and rule versions together, rather than independently
upgrading a CLI and its semantic model.

The Docker-owned portion should remain bounded: acquisition policy, source
profile, schema-position enumeration, presentation data, product metadata,
and fidelity checks. Reuse the libraries for parsing, reference indexes,
schema evaluation, and contract comparison. Do not implement another JSON
Schema evaluator or attempt to flatten arbitrary schemas in templates.

### Required integration checks

Before production use, verify the assembled pipeline, not only each package:

- Register and test the exact OAS dialect and annotation vocabulary. The
  successful raw-schema Go probe used Draft 2020-12; it does not certify all
  OAS-specific annotation behavior. Preserve the source's declared dialect.
- Test `$id`, named anchors, dynamic references, cycles, and `$ref` siblings
  with the locked resource registry. Block any admitted feature the pipeline
  cannot handle faithfully.
- Define request/response, `format`, regular-expression, content, and streaming
  example policies. A passing JSON instance test does not prove correct wire
  encoding or actual service behavior.
- Keep original schema nodes through the prepared model and prove the
  renderer covers them. A correct raw source does not excuse a partial UI.
- Test diagnostics, CLI exit behavior, and coverage on the same fixture corpus
  whenever dependencies or the Docker source profile change.

### Components to keep out of the default build

Do not run Spectral alongside Vacuum, `ajv` alongside the Go schema evaluator,
or Python checks alongside both to add more checks. Add an independent
tool temporarily when investigating a discrepancy, with its scope recorded.

Avoid depending directly on undocumented convenience models such as
`pb33f/doctor` for Docker's compiler API. Use documented source/index
facilities and isolate them behind a small adapter.

Schemathesis is useful in product CI against controlled API instances. The
research only ran operation enumeration and synthetic response checks; no
HTTP calls were sent. Its [SSE documentation][schemathesis-sse] says it
buffers complete responses, so infinite streams require bounded fixtures or
another harness. Live API testing should not make docs builds dependent on
service availability or credentials.

## Effect on the architecture proposal

The tailored UI, Hugo publication, OpenAPI 3.2 target, and source ownership
requirements remain. The component recommendation changes from a
Redocly/TypeScript starting point to a source-preserving Go pipeline.

The canonical publication artifact becomes the validated source package
and its reference graph. Single-file bundling becomes an optional derivative.
Schema compilation and example evaluation become explicit mandatory stages.
Compatibility reports support review with known limits, and the pipeline
retains its own completeness checks.

These recommendations are grounded in the examined versions and fixtures.
They are a more specific implementation basis, not a claim of complete
OpenAPI conformance by any individual tool.

The prototype uses the locked OAI document schema for structural validation,
explicit compilation for Schema Objects, and Vacuum for Docker policy rules.
It also isolates generated Markdown from a site processing script that rewrites
URLs inside fenced JSON. These refinements retain the proposed component
boundaries while incorporating end-to-end evidence.

[libopenapi]: https://github.com/pb33f/libopenapi/releases/tag/v0.38.7
[vacuum]: https://github.com/daveshanley/vacuum/blob/v0.30.3/README.md
[lib-validator]: https://github.com/pb33f/libopenapi-validator/tree/v0.14.0
[go-jsonschema]: https://github.com/santhosh-tekuri/jsonschema/blob/v6.0.3/README.md
[kin]: https://github.com/getkin/kin-openapi/blob/v0.149.0/README.md
[scalar-parser]: https://github.com/scalar/scalar/tree/main/packages/openapi-parser
[scalar-validator]: https://github.com/scalar/scalar/blob/main/packages/openapi-validator/src/validate.ts
[scalar-bundle]: https://github.com/scalar/scalar/blob/main/packages/json-magic/src/bundle/bundle.ts
[spectral]: https://github.com/stoplightio/spectral/blob/0c2fa158147dcef0d8d577fbf753abbb1748f9ad/packages/rulesets/src/oas/functions/oasDocumentSchema.ts
[redocly]: https://github.com/Redocly/redocly-cli/blob/8938c898beaf8849054f6660427bfafe03b89a54/packages/core/src/types/oas3_2.ts
[apidom]: https://github.com/swagger-api/apidom/tree/219d365159f147e6068c56ff06d57ce2ce9297f6/packages/apidom-ns-openapi-3-2
[apidom-bundle]: https://github.com/swagger-api/apidom/blob/219d365159f147e6068c56ff06d57ce2ce9297f6/packages/apidom-reference/src/bundle/strategies/openapi-3-2/index.ts
[apidom-expand]: https://github.com/swagger-api/apidom/blob/219d365159f147e6068c56ff06d57ce2ce9297f6/packages/apidom-reference/src/dereference/strategies/openapi-3-2/visitor.ts
[swagger-parser]: https://github.com/APIDevTools/swagger-parser/blob/2e1672443dc4834d82abec1229c894d6cf8fc11a/lib/index.js
[ref-parser]: https://github.com/APIDevTools/json-schema-ref-parser
[ajv]: https://github.com/ajv-validator/ajv/blob/v8.20.0/docs/json-schema.md
[python-spec]: https://github.com/python-openapi/openapi-spec-validator/blob/0.9.0/README.rst
[python-schema]: https://github.com/python-openapi/openapi-schema-validator/blob/0.9.0/README.rst
[oasdiff]: https://github.com/oasdiff/oasdiff/releases/tag/v1.31.0
[oasdiff-caveats]: https://github.com/oasdiff/oasdiff/blob/v1.31.0/docs/OPENAPI-31.md
[oasdiff-flattening]: https://github.com/oasdiff/oasdiff/blob/v1.31.0/docs/ALLOF.md
[java-diff]: https://github.com/OpenAPITools/openapi-diff/blob/2.1.7/README.md
[schemathesis]: https://github.com/schemathesis/schemathesis/releases/tag/v4.26.0
[schemathesis-sse]: https://github.com/schemathesis/schemathesis/blob/v4.26.0/docs/guides/server-sent-events.md
