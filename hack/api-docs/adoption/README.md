# Local API source adoption

This directory preserves the source-review evidence from
[the API documentation prototype](https://github.com/docker/docs/pull/26043).
The sibling implementation reads the authoritative local files directly.
No migration scripts or duplicate source snapshots are required to build it.

The JSON ledgers retain before/after values, source pointers, classifications,
evidence, and stable change IDs. Textual patches show the original conversions.
`source-lock.json` identifies the baseline revision and source digests. Historical
profile names and paths in these records describe that conversion, rather than
runtime inputs. Keep these records unchanged when making subsequent corrections;
record those corrections in separate reviewed commits.

| API | Recorded changes | Review before merge |
| --- | ---: | --- |
| [Hub](hub.json) | 157 | Confirm the provisional `team_repo` response schema. Review requiredness corrections, assigned operation IDs, and the inherited example mismatches. Preserve separate SCIM credentials. |
| [DVP](dvp.json) | 14 | Confirm operation-level security matches service behavior and the documented legacy login flow remains supported. Preserve authentication server overrides. |
| [Registry](registry.json) | 7 | Confirm `info.version: 2` and the bearer security declaration, including anonymous access alternatives. |

The runtime validation baseline is in `../known-issues.json`. Its 289 entries
remain visible in generated validation reports, without exposing migration
administration in reference pages. Resolve or explicitly approve the remaining
quality debt before production adoption. Contract assumptions require product
evidence; accepting documentation debt does not confirm those assumptions.

Published YAML retains its URL but changes to OpenAPI 3.2. Consumers of the YAML
must be considered before release: URL compatibility does not imply parser
compatibility. Existing page URLs, the DVP alias, and operation/tag/schema
fragment links remain supported. JavaScript follows historical operation
fragments to individual pages; without it, the overview supplies ordinary links.

Engine and Governance migration, historical Engine conversion, SDK catalogs,
and broader site information architecture are outside this implementation.
