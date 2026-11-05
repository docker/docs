# API-change reporter

Compare two OpenAPI documents in JSON or YAML and report removed paths and HTTP
operations. This example supports inline Path Item objects. It does not resolve
references or establish complete API compatibility.

Use Python 3.11 or later. Install the pinned dependency and run the tests from
this directory:

```console
$ python3 -m venv .venv
$ .venv/bin/python -m pip install -r requirements.txt
$ .venv/bin/python -m unittest discover -s tests -v
$ .venv/bin/python reporter.py examples/before.yaml examples/after.json --html output/index.html
```

The examples contain an unchanged `GET /pets`, an added `POST /pets`, a removed
`DELETE /pets/{id}`, and a removed `/legacy` path. The report counts two removed
operations and one removed path. A successful comparison exits with code zero,
including when it finds removals. Invalid or missing input exits with code two.

Serve only the output directory to preview the report:

```console
$ .venv/bin/python -m http.server 8080 --bind 0.0.0.0 --directory output
```

The local-to-cloud guide uses this example as a starting point. Detection of
required query parameters is intentionally left for the agent's task.
