---
title: Investigate an upgrade and report to GitHub
description: Ask Claude Code to investigate an Express upgrade in a cloud sandbox, test a focused fix, and return an evidence-backed report to a GitHub issue.
keywords: cloud sandboxes, claude code, express, dependency upgrade, github issues
summary: |
  Reproduce a dependency upgrade failure, test the smallest fix, and return
  findings to the GitHub issue where the work was requested.
params:
  tags: [ai]
  time: 20 minutes
---

<!-- Publication gate: publish the companion repository and add its URL and
baseline revision. Verify Console authentication, successful and blocked issue
reports, and capture the resulting issue comment. Local evidence is in
samples/cloud-sandbox-workflows/output/. Email delivery has not been tested;
do not add delivery instructions until a test message reaches the author. -->

Ask an agent to investigate a dependency upgrade and put the findings where
you'll review them. This guide uses Claude Code in a cloud sandbox to establish
a passing baseline, reproduce a failure, and test a fix. The agent returns an
evidence-backed report to a GitHub issue so the findings outlive the sandbox.

The example upgrades a small application from Express 4.21.2 to 5.1.0. These
historical versions demonstrate a routing change. They aren't recommendations
for a production application.

## Before you start

You need a browser, an Anthropic API key, and a personal Docker account with an
active [Docker Agentic Platform subscription](/manuals/agentic-platform/signup.md#activate-cloud-access).
You also need a GitHub repository and an issue where you can post findings.
No local installation is required. Compute and model usage are billed
separately.

Use an existing project issue or create one in your fork of the companion
repository. Its `express-upgrade/` directory contains an application with a
fallback route and tests for `/`, `/projects/42/settings`, and `/api/status`.
Record the full baseline commit SHA before starting.

Create a fine-grained GitHub token scoped to that repository, with
**Contents: Read-only** and **Issues: Read and write**. If the source and issue
are in different repositories, include access to both. The investigation
doesn't need permission to push code.

## Launch with repository access

In the [Console](https://agentic-platform.docker.com/), select **New** and the
Claude Code kit. Configure these options before selecting **Run**:

- Add your Anthropic API key and select your token under **GitHub token**
- Review [network policies](/manuals/agentic-platform/policies.md) for
  `api.anthropic.com:443`, `github.com:443`, `api.github.com:443`, and
  `registry.npmjs.org:443`
- Set **Run for** to two hours and **When the time is up** to **Stop**

Managed credentials authenticate the agent's requests, and network policies
control which destinations it can reach. Configure both. Keep credential
values out of prompts and repository files.

## Request the investigation

Replace `REPOSITORY_URL`, `BASELINE_SHA`, and `ISSUE_URL` with your repository,
full commit SHA, and destination issue URL. Paste this request:

```text
Investigate upgrading express from 4.21.2 to exactly 5.1.0.
Repository: REPOSITORY_URL
Baseline revision: BASELINE_SHA
Project directory within the repository: express-upgrade
Test command: npm test
Destination issue: ISSUE_URL

Clone the repository inside the sandbox and check out the baseline revision.
Record the commit, Node.js version, and installed Express version. Run npm ci
and npm test before changing anything. If the baseline fails, stop the upgrade
and write a blocked-investigation report with the failure evidence.

If the baseline passes, install express@5.1.0 with an exact version, run the
same tests, and record the failure before editing application code. Diagnose
the cause, propose the smallest fix, and rerun all tests. Preserve the tests
for /, /projects/42/settings, and /api/status. Don't weaken their assertions.

Save report.md with the baseline revision, versions tested, before/after test
results, diagnosis, proposed changes, and unresolved work. Include the meaningful
code diff. Save the full dependency, lockfile, and application diff as proposed.patch.
Do not push changes. Show me the report and paths to both files.
```

## Read the evidence

The complete companion application was tested locally at each stage:

| Stage | Result |
| --- | --- |
| Express 4.21.2 | All three route tests pass |
| Express 5.1.0, unchanged fallback | Application fails to load with `Missing parameter name` |
| Express 5.1.0, `/*splat` fallback | Two tests pass; `/` returns 404 |
| Express 5.1.0, `/{*splat}` fallback | All three route tests pass |

Express 5 requires a named wildcard. The braces also include the root path,
as described in the [Express migration guide](https://expressjs.com/en/guide/migrating-5.html#path-route-matching-syntax).
For this fixture, the application change is:

```diff
-app.get("*", (req, res) => {
+app.get("/{*splat}", (req, res) => {
```

The dependency declaration and lockfile also change. The API route stays before
the fallback. Check that the agent's report accounts for all three routes,
including `/`, rather than accepting a fix based only on a nested client route.

A failing baseline needs a different conclusion: the investigation is blocked,
and the failure isn't evidence against Express 5. The local verifier exercises
this case by making a baseline assertion fail and checking that no upgrade is
attempted.

## Return the report to the issue

After reviewing `report.md`, ask the agent to deliver it:

```text
Post the complete contents of report.md as one comment on ISSUE_URL using
GitHub's REST API and the configured GitHub credential. Include the report even
if the investigation was blocked. Don't include secrets or unrelated files.
Verify that the request succeeded, read the comment back, and return its URL.
If delivery fails, keep report.md and tell me the error. Don't claim it was posted.
```

Replace `ISSUE_URL` again before sending. GitHub's
[create-comment endpoint](https://docs.github.com/en/rest/issues/comments#create-an-issue-comment)
accepts the Markdown report as the comment body. GitHub MCP is an alternative
tool interface; this workflow already has authentication and an HTTP API.

Open the returned comment URL. Confirm that the baseline revision, failure,
proposed diff, and final test results are visible. Save `proposed.patch` and any
other files you want to keep before deleting the sandbox. For this small example,
ask the agent to show the patch in the terminal and copy it to a local file.
The issue retains the report, but posting a comment doesn't preserve files left
in the sandbox.

## Keep the fix or request email

Optional. Ask the agent to commit the tested fix on a branch and open a draft
pull request. Configure **Contents: Read and write** and
**Pull requests: Read and write** on the repository token before launching a
sandbox for that task.

Optional. Use Resend for a separate completion email containing the issue link.
Before launch, configure a [custom secret](/manuals/agentic-platform/secrets.md#add-a-custom-secret)
for `api.resend.com` with the `Authorization` header and `Bearer %s` format,
and permit that host in the network policy. Resend's
[test sender restriction](https://resend.com/docs/knowledge-base/403-error-resend-dev-domain)
limits messages from `onboarding@resend.dev` to the account owner's address.
Other recipients require a verified domain. Validate delivery to your own
address before using this extension for completion messages.
