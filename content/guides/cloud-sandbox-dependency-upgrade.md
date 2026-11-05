---
title: Assess a dependency migration in the cloud
description: See how Claude Code investigates an Express upgrade in a cloud sandbox, tests compatibility changes, and creates a GitHub issue to plan the migration.
keywords: cloud sandboxes, claude code, express, dependency migration, github issues
summary: |
  Follow an agent's assessment of a dependency migration, from compatibility
  failures to a GitHub tracking issue with a recommendation and remaining work.
params:
  tags: [ai]
  time: 5 minutes
---

<!-- Publication gate: capture a cloud run of the expanded support-portal
sample, including the generated tracking issue. The compatibility evidence below
was reproduced locally. See samples/cloud-sandbox-workflows/README.md. -->

A dependency upgrade can pass one check and still change behavior that your
clients rely on. This demonstration follows Claude Code as it investigates an
Express migration in a cloud sandbox and creates a GitHub tracking issue. The
result is a recommendation with evidence and remaining work that a maintainer
can use to plan the release.

The sandbox provides the execution environment, agent credentials, and repository
access. The investigation runs there, and the issue keeps the findings accessible
after the sandbox is deleted.

## Assess a support portal

The example application serves a client interface and a ticket-search API.
It uses Express 4.21.2, and the assessment targets Express 5.1.0. These historical
versions reproduce compatibility changes. They aren't recommendations for a
production application.

Eight tests cover behavior that the upgrade needs to preserve:

- The client opens at `/` and a nested route
- API routes return JSON before the client fallback
- Ticket searches honor nested status and assignee filters
- JSON searches return matching tickets
- An empty search request returns all tickets

The decision is whether the application can migrate while preserving its client
and API contract. Finding a routing fix alone doesn't answer that question.

## Give the agent access

To try the workflow, you need a browser, an Anthropic API key, and a personal
Docker account with an active
[Docker Agentic Platform subscription](/manuals/agentic-platform/signup.md#activate-cloud-access).
No local development environment is required. Compute and model usage are
billed separately.

Use a GitHub repository with Issues enabled. Configure a fine-grained token
with **Contents: Read-only** and **Issues: Read and write** on that repository.
These permissions cover private repository access and creation of the tracking
issue.

In the [Console](https://agentic-platform.docker.com/), select **New** and the
Claude Code kit. Add the Anthropic key and select the token under **GitHub token**.
Select the **Balanced** network policy.

Optional. Select a [custom network policy](/manuals/agentic-platform/policies.md#create-a-policy)
for `api.anthropic.com:443`, `github.com:443`, `api.github.com:443`, and
`registry.npmjs.org:443` instead of **Balanced**. Deselect **Open** if it is
selected, since its rules permit broader access.

Select **Run** to launch the sandbox.

## Request a migration assessment

A useful request combines the investigation and its destination. Replace
`REPOSITORY_URL`, `BASELINE_SHA`, and `PROJECT_DIRECTORY` with the repository
URL, full source commit SHA, and project directory:

```text
Assess whether this support portal can migrate from Express 4.21.2 to 5.1.0
while preserving its client and API behavior.
Repository: REPOSITORY_URL
Baseline: BASELINE_SHA
Project: PROJECT_DIRECTORY

Clone the baseline. Run npm ci and npm test before changing dependencies.
If it passes, upgrade to exactly Express 5.1.0, record the failures, and test
minimal compatibility fixes. Keep existing test assertions. Don't push code.

Create a migration tracking issue in the same repository. Include a recommendation,
source revision, versions, test results, affected client behavior, proposed changes,
and remaining release checks. Keep it under 400 words; omit full logs and lockfile
contents. Save the proposed patch and issue body in the sandbox, and return the
issue URL. Reuse a matching tracking issue if one exists.

If the baseline fails, report the investigation as blocked without upgrading.
```

The issue is the output of the task. You don't need to create a destination
issue or make a separate request to post a report.

## Look beyond the first failure

The local reproduction found three compatibility problems:

| Behavior | Upgrade impact | Proposed change |
| --- | --- | --- |
| Client fallback | The unnamed wildcard prevents application startup | Name the wildcard and include the root path with `/{*splat}` |
| Nested query filters | A filtered request returns all tickets because the default parser no longer produces nested objects | Select the extended query parser to preserve the existing API |
| Empty search | Destructuring an unparsed request body returns HTTP 500 | Handle an undefined body as an empty search |

These changes follow the [Express migration guidance](https://expressjs.com/en/guide/migrating-5.html).
The root route matters: `/*splat` handles nested routes but excludes `/`.
After fixing startup, the tests expose incorrect filtering and the empty-search
failure.

The test results show why the assessment needs to continue after each fix:

| Stage | Result |
| --- | --- |
| Express 4.21.2 baseline | Eight tests pass |
| Express 5.1.0, unchanged application | Application fails to load |
| Fallback includes `/` | Five tests pass; three fail |
| Nested query parsing restored | Seven tests pass; one fails |
| Empty request body handled | Eight tests pass |

Restoring the extended parser preserves existing requests such as
`/api/tickets?filters[status]=open`. Changing the API to accept flat filters
would also require changes to clients. That choice belongs in the migration
plan rather than being hidden in a patch.

## Read the tracking issue

A useful issue leads with the decision and keeps the evidence short. For this
example, its recommendation could be:

> Preserve the existing API contract for this migration. The proposed routing,
> query parsing, and empty-body changes pass all eight tests. Validate real client
> requests, filter limits, and dependency advisories before scheduling a release.

The issue should identify the tested revision and distinguish confirmed failures
from untested behavior. Follow-up work might include validating malformed filters,
checking production client requests, and planning rollout and rollback. Eight
passing tests establish the covered behavior; they leave those release checks
for the maintainer.

The agent can create the issue through
[GitHub's issue API](https://docs.github.com/en/rest/issues/issues#create-an-issue).
The configured GitHub credential authenticates the request. An additional tool
connection isn't required for this workflow.

Open the returned URL and read the recommendation, evidence, and remaining work.
Save any proposed patch you want to keep, then select **Delete** in the Console.
The issue remains available for planning after the execution environment is gone.

Optional. If you decide to proceed, ask an agent to prepare a draft pull request
linked to the tracking issue. Before launching that sandbox, configure
**Contents: Read and write** and **Pull requests: Read and write** on the token.
The tracking issue gives that follow-up task its scope and acceptance checks.
