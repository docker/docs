---
title: Get started with Docker Agentic Platform
linkTitle: Get started
description: Choose a kit and start a Docker Agentic Platform sandbox.
keywords: docker agentic platform, get started, kits, agents, sandbox, live terminal
weight: 10
aliases:
  - /agentic-platform/concepts/platform-model/
---

Start a sandbox from the Console to work with an agent in the cloud.

## Before you begin

You need a personal Docker account with an active
[Docker Agentic Platform subscription](signup.md#activate-cloud-access).

If your kit needs a model provider API key, add it under **Secrets** or enter
it when you launch the sandbox.

Browse [Kits](kits.md) to choose an agent such as Claude Code, Codex, or Hermes,
or use Shell to work without a pre-installed agent. The credentials you need
depend on the kit. For example, Hermes needs an Anthropic or OpenAI API key.

## Start a sandbox

1. Open [Docker Agentic Platform](https://agentic-platform.docker.com/) and
   select **New**.
2. Choose a kit.
3. Add any credentials requested by the launcher.
4. Review the [sandbox options](#sandbox-options), including the selected
   network policies, and adjust them for your task.
5. Select **Run**. If a required credential is missing, the launcher opens its
   input so you can add it.

When the sandbox is ready, its terminal opens. Use it to work with your agent,
or run commands if you chose Shell. After launch, you can't change the selected
credentials, tools, network policies, or compute size.

Return to **Sandboxes** to reopen the terminal.

## Sandbox options

Configure these options before selecting **Run**:

| Option | What to choose |
| --- | --- |
| Agent credentials | Add the API keys or tokens your kit needs. Saved credentials are reused by default; deselect optional credentials you don't want to include. See [Secrets](secrets.md). |
| GitHub token | Provide a token for Copilot or for cloning private repositories and pushing changes. Copilot uses the same token for both. Repository access is optional for other agents. See [GitHub credential](secrets.md#github-credential). |
| Network policies | Review the hosts and services the sandbox can reach. See [Network policies](policies.md). |
| MCP tools | Choose the servers your agent can use and authorize them if prompted. See [MCP](mcp.md). |
| Compute and platform | Choose CPU and memory resources. Under **Platform**, use **auto** or a platform supported by the kit's image. See [Sandbox platforms](sandboxes.md#choose-a-platform). |
| Timer | Under **Run for**, choose 1 to 24 hours. Under **When the time is up**, choose **Stop**, **Delete**, or **Restart** if offered. See [Lifecycle actions](sandboxes.md#manage-the-lifecycle). |

The launcher remembers your policy selections from the previous launch in the
same browser. Policies marked **Always applied** are selected automatically.
Review the selected policies before launching. Selecting **Balanced** alongside
**Open** doesn't restrict Open's access; see [How policies combine](policies.md#how-policies-combine).
