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
- The generated REST reference and its OpenAPI specification to
  `content/reference/api/sandboxes/`

Treat these files as generated. Make prose, example, and API corrections in
`docker/sbx-api`, then import the updated export. The cookbook's `_index.md`
is maintained in docker/docs, along with the overview, installation,
authentication, concepts, errors, and limits pages.

The import excludes the upstream overview and installation page. It redirects
recipe installation links to the local installation page and adjusts the API
reference's title, layout, navigation, introductory links, and experimental
notice.
Recipe code and prose are otherwise unchanged. The OpenAPI YAML is copied
byte for byte.

The reference uses upstream's generated Markdown with a page layout that wraps
long paths and schema names on narrow screens. This preserves the OpenAPI 3.2
schemas, streaming response details, and operation-specific servers and
authentication. The existing
`api-reference` layout does not render all of those details. The generated
reference has its own schema disclosures and a downloadable specification;
it does not provide an interactive HTTP request console.

To update, change both commit IDs in `source.json`, review additions or removals
in the upstream guide list, and run the import. Removed guides require explicit
deletion and a link check. Review the diff and build the site before committing.
The import does not select the latest release or push changes.

Release automation can update this manifest and run the same command in a pull
request once upstream publishes a release-to-export mapping. Keep the source
and generated commits paired rather than importing a moving branch.
