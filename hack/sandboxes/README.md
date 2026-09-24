# Sandboxes documentation import

The public cookbook export is published to the `generated` branch of
`docker/sbx-api` after its documentation CI passes. `source.json` pins both
the generated commit and the source commit recorded in that export's `SOURCE`
file. The guide list explicitly selects the recipes to publish.

Run the import from the repository root with Python 3 and `gh` authenticated for
`docker/sbx-api`:

```console
$ python3 hack/sync-sandboxes-docs.py
```

The import writes:

- Cookbook recipes to `content/manuals/ai/sandboxes-api/cookbook/`
- The public OpenAPI specification to
  `content/reference/api/sandboxes/api.yaml`

Treat these files as generated. Make prose, example, and API corrections in
`docker/sbx-api`, then import the updated export. The cookbook's `_index.md`
is maintained in docker/docs, along with the overview, installation,
authentication, concepts, errors, and limits pages.

The import excludes the upstream overview, installation page, and generated
Markdown API reference. It redirects recipe installation links to the local
installation page. Recipe code and prose are otherwise unchanged. The OpenAPI
YAML receives the temporary example correction recorded in
`source.json` as `apiPatch`. The patch supplies 18 missing examples and
corrects 13 invalid examples without changing API schemas or operations.
`patches/upstream-examples.patch` contains the owning sbx-api source fix;
`patches/api-examples.patch` contains its effect on the pinned export.
The import applies that patch before writing any files and rejects mismatches.

After sbx-api publishes the corrected export, update the pinned commits,
remove `apiPatch` and both patch files, and import again. The YAML can then
be copied byte for byte. Until then, running the import requires `git`
in addition to Python and `gh`.

The specification is registered in `hack/api-docs/sources.json`. The shared API
renderer generates an overview and separate operation and schema pages under
`/reference/api/sandboxes/latest/`, in HTML and Markdown. Run
`./hack/api-docs/run.sh generate` before Hugo. Strict validation checks the
specification and its examples. Fix failures upstream and re-import the export.

To update, change both commit IDs in `source.json`, review additions or removals
in the upstream guide list, and run the import. Removed guides require explicit
deletion and a link check. Review the diff and build the site before committing.
The import does not select the latest release or push changes.

Release automation can update this manifest and run the same command in a pull
request once upstream publishes a release-to-export mapping. Keep the source
and generated commits paired rather than importing a moving branch.
