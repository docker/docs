---
title: Sandboxes
description: Create and manage hosted agent environments in Docker Agentic Platform.
keywords: docker agentic platform, sandboxes, agents, cloud runtime, terminal, compute
weight: 20
aliases:
  - /agentic-platform/concepts/sandboxes/
  - /agentic-platform/guides/manage-sandboxes/
---

A sandbox is an isolated runtime for an agent or tool. Docker hosts and meters
the sandbox on managed cloud infrastructure and provides a live terminal for
interacting with it.

Docker Agentic Platform provides predefined sandbox types for Claude Code,
Codex, OpenCode, Copilot, and Gemini CLI. Each sandbox has its own compute,
filesystem, network access, and terminal. A sandbox continues running
independently of your connection to the Console until it is paused, stopped by
its lifecycle timer, or deleted.

## Source code and files

A sandbox starts with a fresh filesystem. Docker Agentic Platform does not
mount a repository, local directory, or workspace from your computer into the
sandbox by default.

A sandbox does not synchronize its filesystem with your computer or a remote
repository. Work remains only in the sandbox unless you commit and push it to a
remote repository. Push work that you want to keep before deleting the
sandbox; files that exist only in a deleted sandbox are not available from a
later sandbox.

Docker Agentic Platform supports GitHub repositories for bringing source into
a sandbox and preserving changes. Configure `GITHUB_TOKEN` under **Secrets** or
from the sandbox launcher to clone a private repository or perform write
operations such as pushing a branch or opening a pull request. The token must
have the required repository permissions. Public repositories can be cloned
without a GitHub credential, but writing to them still requires authentication.

## Open a sandbox

After you select **Run**, Docker creates the sandbox and opens its detail page.
Use the terminal to interact with the selected agent or tool.

Open **Sandboxes** to review each sandbox's name, type, status, hourly rate,
expiration, and age. Select a sandbox to reopen its detail page and terminal.

## Manage the lifecycle

A sandbox can be running or paused:

- Pause a running sandbox to stop its compute without deleting it.
- Resume a paused sandbox to continue working with it.
- Delete a sandbox when you no longer need it.

When you create a sandbox, set a lifecycle timer from 1 to 24 hours and choose
what happens when it expires. **Stop** stops the sandbox, while **Delete**
deletes the sandbox and its files.

The sandbox's authentication, tools, access policy, and compute size are fixed
when the sandbox is created and cannot be changed while it runs.

Docker bills sandbox compute per second while the sandbox runs. Model inference
uses your external provider credential and is metered by that provider. For
account, usage, and payment information, see [Docker Billing](/subscription-billing/billing/).

## Check sandbox configuration

If the agent cannot reach a service or use a tool, check the configuration that
applies to the request:

- Confirm that the network policies permit the destination.
- If the agent needs an MCP tool, confirm that its server is connected and
  authorized.
- If the destination requires authentication, confirm that the required
  service credential is configured under **Secrets**.
