# API documentation prototype tooling

This toolchain converts copied API specifications, records upstream changes,
and builds an opt-in Hugo reference. See the
[prototype report](../../API-DOCUMENTATION-PROTOTYPE.md) for findings and limits.

Run from the repository root:

```console
./hack/api-docs/run.sh build
./hack/api-docs/run.sh serve
```

Open `http://localhost:1314/api-prototype/`. The local commands need Go 1.26.5,
Hugo 0.163.0, Node.js, npm, and standard shell utilities. The server command also
uses Python 3. The container build supplies its dependencies:

```console
docker buildx bake api-prototype
python3 -m http.server 1314 --bind 127.0.0.1 --directory tmp/api-prototype/site
```

The build exports the complete preview site, including the existing manuals,
under `tmp/api-prototype/site`. It never deploys or submits upstream changes.
The API pages exist only with the prototype Hugo config.

This draft's Netlify deploy previews run the same pipeline and publish
`tmp/api-prototype/site`, using `DEPLOY_PRIME_URL` as the site origin. Open
`/api-prototype/` on the deploy preview to explore the catalog. Production
builds retain their existing configuration.

## Commands

| Command | Behavior |
| --- | --- |
| `bootstrap` | Install pinned migration dependencies and build the Go command, Vacuum, and the contract comparison tool |
| `verify` | Verify source digests, replay all migration stages, and compare converted sources, ledgers, and textual diffs |
| `convert` | Recreate converted sources and migration records from immutable snapshots and recorded choices |
| `check` | Run source verification, Vacuum policy rules, and strict validation; recorded exceptions still fail strict validation |
| `generate` | Run the pipeline through prepared data, checked source downloads, and comparison reports, accepting only exact preview exceptions |
| `report` | Check source/model coverage and regenerate the migration summary |
| `test` | Run Go regression tests and migration reproduction |
| `build` | Generate data, build Hugo, preserve reference Markdown, index search, and check output |
| `serve` | Build and serve the preview on localhost |

The low-level Go command accepts `check`, `generate`, or `inspect`, followed by
the repository root. `inspect` writes diagnostics without declaring the input
acceptable. `--preview` accepts only recorded, matching exceptions. Use the
wrapper for the complete publication gates.

## Source and diagnostic records

The [catalog](../../prototypes/api-docs/catalog.json) describes the six snapshots.
The [source lock](../../prototypes/api-docs/source-lock.json) records their
origins and digests. Each migration ledger contains ordered before/after
changes; pointers refer to that stage's input/output document. A converter's
move can appear as a removal and an addition. Textual patches also retain
formatting changes.

Preserve existing operation IDs. Missing IDs come from the checked-in
[ID mapping](../../prototypes/api-docs/operation-ids.json). Add specific,
reviewed correction patches to `corrections.json` beside the catalog when
needed, with the same change classification and evidence fields as the ledger.
Do not edit the original snapshots or authoritative specifications.

Validation exceptions include the API, source digest, rule, pointer, diagnostic
hash, and reason. They are deliberate review records, not an automatically
refreshed allowlist. A changed source or diagnostic invalidates an exception.
`check` fails until the source is ready even when `build` produces a draft.

## Output and coverage

The Go command writes presentation model version 1 under
`tmp/api-prototype/data`. It contains effective operations and the original
schema graph. Hugo templates consume that interface. Source packages contain
both specifications, the migration ledger, textual diff, lock, and exceptions.
The six publication specifications have no external references; external
resource handling is exercised separately by fixtures.

Coverage reports compare operation identities, full operation objects, media
variants, headers, and named schemas. Output checks verify HTML/Markdown pairs
and preview-owned links. The browser checks export a function accepting a
Playwright `Page`; run it against the local preview with JavaScript enabled.

The source loader preserves YAML timestamp spelling and rejects duplicate
keys. JSON Schema evaluation uses the locked OAS dialect resources, with
format/content annotations left as annotations. Request/response annotations
are displayed; complete directional example evaluation and wire-format
validation remain implementation work.

## Search portability

The published `pagefind` 1.5.2 Linux binary fails on the tested 16 KB-page host.
The local build uses the optional Docker recipe to compile the same release
with `JEMALLOC_SYS_WITH_LG_PAGE=16`. Other local hosts use the published binary.
The prototype container includes that rebuild. This affects the opt-in target;
the production search build is unchanged.

Set `PAGEFIND_BIN` to use a prepared compatible binary. Set
`API_PROTOTYPE_URL` when building for a different preview origin. For the
local server, also set `API_PROTOTYPE_PORT` to the matching port.
