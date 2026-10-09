---
title: Build an API-change reporter in the cloud
description: Ask Claude Code to build and test an API-change reporter in a cloud sandbox, preview its HTML report in your browser, and download the source and tests.
keywords: cloud sandboxes, claude code, openapi, python, api changes
summary: |
  Turn an idea into a tested command-line tool, open its report through a
  sandbox port, and download the work to keep it.
params:
  tags: [ai]
  time: 20 minutes
---

<!-- Publication gate: run the Console journey with a recorded Claude Code kit
version. Verify the public report and ZIP after sandbox deletion. The screenshot
and test excerpt below come from the locally tested companion implementation.
See samples/cloud-sandbox-workflows/README.md for the remaining checks. -->

Build a small utility without setting up a development environment on your
computer. In this guide, you give Claude Code a task in a cloud sandbox, inspect
its tests, and open the result in your browser. You finish with a ZIP containing
the source, examples, tests, and report.

The utility compares two OpenAPI documents and reports removed paths and HTTP
operations. These checks can flag changes that affect API clients. They don't
establish complete API compatibility.

## Before you start

You need a browser, a personal Docker account with an active
[Docker Agentic Platform subscription](/manuals/agentic-platform/signup.md#activate-cloud-access),
and an Anthropic API key. No local installation or GitHub account is required.
Docker bills sandbox compute, and Anthropic bills model usage separately.

## Start the agent

1. Open the [Console](https://agentic-platform.docker.com/) and select **New**.
2. Select the Claude Code kit and add your Anthropic API key through the
   credential control. Keep the key out of the task prompt.
3. Review the selected [network policies](/manuals/agentic-platform/policies.md).
   Allow `api.anthropic.com:443`, `pypi.org:443`, and
   `files.pythonhosted.org:443` for the agent and Python dependency downloads.
   Add a custom policy if the selected policies don't cover those hosts.
4. Set **Run for** to two hours and **When the time is up** to **Stop**.
   Select **Run**.

The Console opens a terminal connected to the agent. Dependencies, files, and
processes created here belong to the sandbox. They don't change your computer.

## Give the agent a bounded task

Paste this request into the agent terminal:

```text
Build a Python CLI in /home/agent/workspace/api-change-reporter that compares
two OpenAPI documents in JSON or YAML. Detect removed paths and removed HTTP
operations only. Explain that this is not a complete API compatibility check.
Support inline Path Item objects; reject referenced Path Items with a clear error.

Use a project-local .venv and pin dependencies in requirements.txt. Create
reporter.py, a README with run instructions, and automated tests. Cover JSON
and YAML input, unchanged and added operations, removed paths, multiple removed
operations on one path, invalid input, and HTML escaping.

Create examples/before.yaml and examples/after.json with these changes:
- GET /pets and GET /pets/{id} stay unchanged.
- POST /pets is added.
- DELETE /pets/{id} is removed.
- /legacy, which had GET, is removed entirely.

Run the tests and compare the examples. Generate a self-contained HTML report
at output/index.html. Keep output/ limited to files intended for the browser.
Show the test command, results, and removed operations. Stop if tests fail.
```

The agent handles setup and implementation. Inspect its results before moving
on. The example should report one removed path and two removed operations.
Adding `POST /pets` should not count as a removal.

For comparison, the companion implementation produced this local test summary
and CLI output:

```text
Ran 8 tests
OK

Removed paths: 1
Removed operations: 2
  /legacy: GET
  /pets/{id}: DELETE
```

Your agent's test count can differ. Check that the tests cover the requested
behavior and that the report agrees with the inputs.

## Open the report

Ask the agent to serve the report:

```text
Serve only /home/agent/workspace/api-change-reporter/output on 0.0.0.0:8080
using the project's Python interpreter. Keep the server running while I
inspect the report. Verify that GET / returns the HTML report.
```

In the sandbox detail page, open **Connect**, find **Connect via a Port**, enter
`8080`, and select **Open Port**. Open the returned URL in your browser.

This is a public endpoint tied to the running sandbox. The server's dedicated
directory limits what you publish to the demonstration outputs. Use sample
specifications without private API details.

Request a small refinement:

```text
Group removed operations by path in the HTML report, with a methods column.
Keep the removed-path count and the scope limitation visible. Rerun the tests
and regenerate output/index.html so I can refresh the same URL.
```

The following report was rendered from the locally tested companion
implementation:

![API-change report showing GET removed from /legacy and DELETE removed from /pets/{id}](images/cloud-sandbox-api-reporter-report.png)

## Download the work

Ask the agent to package the project:

```text
Create output/api-change-reporter.zip containing only the source, README,
requirements.txt, tests, example inputs, and the generated HTML report. Exclude
.venv, caches, credentials, Git metadata, and unrelated sandbox files.

Inspect the ZIP's member list. Extract it to a temporary directory and verify
that its tests and documented example command work. Show the results and add
a relative download link to the ZIP on the served report page.
```

Refresh the report and follow the download link. Open the downloaded ZIP and
confirm that it contains the project and report. Saving this bundle is the step
that keeps the work after sandbox deletion.

Use the close action next to port `8080` in **Connect**, then delete the sandbox
from the Console. The report URL stops serving the application when the sandbox
is no longer running.

## Use your own inputs

Point the agent at two specification files in a repository it can access. Give
it the repository URL, revisions, and file paths, and allow the repository host
in the sandbox's network policy. For private GitHub files, select a
[GitHub credential](/manuals/agentic-platform/secrets.md#github-credential)
with repository read access before launching. Keep reports containing private
details out of a public port unless the service authenticates readers.

To continue developing this tool across local and cloud environments, see
[Send local work to the cloud](cloud-sandbox-local-handoff.md).
