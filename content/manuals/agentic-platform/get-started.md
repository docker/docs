---
title: Get started with Docker Agentic Platform
linkTitle: Get started
description: Choose an environment and start a Docker Agentic Platform sandbox.
keywords: docker agentic platform, get started, agents, sandbox, live terminal
weight: 10
aliases:
  - /agentic-platform/concepts/platform-model/
  - /agentic-platform/concepts/kits/
  - /agentic-platform/guides/manage-kits/
---

The Docker Agentic Platform launcher collects the configuration needed to start
an isolated sandbox.

## Before you begin

You need a Docker account and access to Docker Agentic Platform. The launcher
prompts for a model provider API key when needed. You can add the key under
**Secrets** before creating the sandbox or provide it in the launcher when
prompted.

The available sandbox types are Claude Code, Codex, OpenCode, Copilot, Gemini
CLI, and Shell. The supported model provider credentials are Anthropic, OpenAI,
GitHub Copilot, and Google.

## Start a sandbox

1. Open [Docker Agentic Platform](https://agentic-platform.docker.com/) and
   select **New**.
2. Choose a sandbox type and add any requested model credential. Copilot uses
   `GITHUB_TOKEN`; add the same GitHub secret to any other sandbox type that
   needs private repository access.
3. Configure the sandbox. The initial settings are **open access**, **no tools
   added**, and **medium compute**.
4. Choose whether Docker stops or deletes the sandbox when its timer expires,
   and set the timer from 1 to 24 hours.
5. Review the configuration and select **Run**.

You cannot change the sandbox's authentication, tools, access policy, or compute
size after it starts. Docker creates the sandbox, marks it as running, and opens
its terminal.

Use the terminal to interact with the sandbox.

Return to **Sandboxes** to find the running sandbox and reopen its terminal.

## Next steps

- [Manage sandboxes](/manuals/agentic-platform/sandboxes.md)
- [Connect MCP servers](/manuals/agentic-platform/mcp.md)
- [Manage secrets](/manuals/agentic-platform/secrets.md)
- [Manage network policies](/manuals/agentic-platform/policies.md)
