# Kit capability support data

The [kit capability page](../../content/manuals/ai/sandboxes/customize/author/capabilities.md)
uses two sources:

- `source.json` pins a public `docker/sandbox-kit-spec` commit and the checksum
  of `kit.schema.json`. Its documentation paths identify each capability page
  at that commit. The schema defines the capability inventory and config
  shapes; it does not state Docker Sandboxes support.
- `support.json` records Docker Sandboxes support for local and cloud sandboxes,
  the first `sbx` version where applicable, and the restrictions shown to
  readers. These entries need product owner review. `unverified` means the
  support status has not been confirmed; it does not mean unsupported.

`data/kit_capabilities.json` is generated from these files. Hugo renders the
table from that data for both HTML and Markdown output. Builds do not fetch
GitHub content.

## Update the catalog

Resolve the public spec branch to a full commit SHA, then run:

```console
$ python3 hack/kit-capabilities/sync.py sync --ref <SPEC_COMMIT_SHA>
```

The sync command requires the `gh` CLI to fetch the schema and documentation
tree through the GitHub API. The weekly sync workflow stops if the spec adds a
capability without an explicit `support.json` classification. Add the entry after
confirming its Docker Sandboxes behavior, or classify it as `unverified` until
that review is complete. Sync checks every entry before changing the pinned
files. Review the generated diff and the hand-written sections of the page,
especially settings that require sandbox creation.

After editing `support.json`, regenerate the data. The check and tests run
offline with Python's standard library:

```console
$ python3 hack/kit-capabilities/sync.py generate
$ python3 hack/kit-capabilities/sync.py check
$ python3 -m unittest discover -s hack/kit-capabilities -p 'test_*.py'
```
