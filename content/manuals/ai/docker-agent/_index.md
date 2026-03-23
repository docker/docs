---
title: Docker Agent
description: Docker Agent lets you build, orchestrate, and share AI agents that work together as a team.
weight: 60
aliases:
  - /ai/cagent/
  - /manuals/ai/cagent/
params:
  sidebar:
    group: Open source
    badge:
      color: violet
      text: Experimental
keywords: [ai, agent, docker agent, cagent]
---

{{< summary-bar feature_name="Docker Agent" >}}

[Docker Agent](https://github.com/docker/docker-agent) is an open-source framework
for building teams of specialized AI agents. Instead of prompting one
generalist model, you define agents with specific roles and instructions that
collaborate to solve problems. Run these agent teams from your terminal using
any LLM provider.

> [!NOTE]
> Docker Agent is a framework for building and running custom agent teams.
> For Docker's built-in AI assistant, see [Gordon](/ai/gordon/) (`docker ai`).

## Why agent teams

One agent handling complex work means constant context-switching. Split the work
across focused agents instead - each handles what it's best at. Docker Agent manages
the coordination.

Here's a two-agent team that debugs problems:

```yaml
agents:
  root:
    model: openai/gpt-5-mini # Change to the model that you want to use
    description: Bug investigator
    instruction: |
      Analyze error messages, stack traces, and code to find bug root causes.
      Explain what's wrong and why it's happening.
      Delegate fix implementation to the fixer agent.
    sub_agents: [fixer]
    toolsets:
      - type: filesystem
      - type: mcp
        ref: docker:duckduckgo

  fixer:
    model: anthropic/claude-sonnet-4-5 # Change to the model that you want to use
    description: Fix implementer
    instruction: |
      Write fixes for bugs diagnosed by the investigator.
      Make minimal, targeted changes and add tests to prevent regression.
    toolsets:
      - type: filesystem
      - type: shell
```

The root agent investigates and explains the problem. When it understands the
issue, it hands off to `fixer` for implementation. Each agent stays focused on
its specialty.

## Installation

Docker Agent is included in Docker Desktop 4.63 and later. In Docker Desktop versions 4.49 through 4.62, this feature was called cagent.

For Docker Engine users or custom installations:

- **Homebrew**: `brew install docker-agent`
- **Winget**: `winget install Docker.Agent`
- **Pre-built binaries**: [GitHub
  releases](https://github.com/docker/docker-agent/releases)
- **From source**: See the [Docker Agent
  repository](https://github.com/docker/docker-agent?tab=readme-ov-file#build-from-source)

The `docker-agent` binary should be copied to `~/.docker/cli-plugins` and then can be used with the `docker agent` command. Alternatively, it can be used as a standalone binary.

## Get started

Try the bug analyzer team:

1. Set your API key for the model provider you want to use:

   ```console
   $ export ANTHROPIC_API_KEY=<your_key>  # For Claude models
   $ export OPENAI_API_KEY=<your_key>     # For OpenAI models
   $ export GOOGLE_API_KEY=<your_key>     # For Gemini models
   ```

2. Save the [example configuration](#why-agent-teams) as `debugger.yaml`.

3. Run your agent team:

   ```console
   $ docker agent run debugger.yaml
   ```

You'll see a prompt where you can describe bugs or paste error messages. The
investigator analyzes the problem, then hands off to the fixer for
implementation.

## How it works

You interact with the _root agent_, which can delegate work to sub-agents you
define. Each agent:

- Uses its own model and parameters
- Has its own context (agents don't share knowledge)
- Can access built-in tools like todo lists, memory, and task delegation
- Can use external tools via [MCP
  servers](/manuals/ai/mcp-catalog-and-toolkit/mcp-gateway.md)

The root agent delegates tasks to agents listed under `sub_agents`. Sub-agents
can have their own sub-agents for deeper hierarchies.

## Configuration options

Agent configurations are YAML files. A basic structure looks like this:

```yaml
agents:
  root:
    model: claude-sonnet-4-0
    description: Brief role summary
    instruction: |
      Detailed instructions for this agent...
    sub_agents: [helper]

  helper:
    model: gpt-5-mini
    description: Specialist agent role
    instruction: |
      Instructions for the helper agent...
```

You can also configure model settings (like context limits), tools (including
MCP servers), and more. See the [configuration
reference](./reference/config.md)
for complete details.

## Share agent teams

Agent configurations are packaged as OCI artifacts. Push and pull them like
container images:

```console
$ docker agent share push ./debugger.yaml myusername/debugger
$ docker agent share pull myusername/debugger
```

Use Docker Hub or any OCI-compatible registry. Pushing creates the repository if
it doesn't exist yet.

## What's next

- Follow the [tutorial](./tutorial.md) to build your first coding agent
- Learn [best practices](./best-practices.md) for building effective agents
- Integrate Docker Agent with your [editor](./integrations/acp.md) or use agents as
  [tools in MCP clients](./integrations/mcp.md)
- Browse example agent configurations in the [Docker Agent
  repository](https://github.com/docker/docker-agent/tree/main/examples)
- Use `docker agent new` to generate agent teams with AI <!-- TODO: link to some page
  where we explain this, probably a CLI reference? -->
- Connect agents to external tools via the [Docker MCP
  Gateway](/manuals/ai/mcp-catalog-and-toolkit/mcp-gateway.md)
- Read the full [configuration
  reference](https://github.com/docker/docker-agent?tab=readme-ov-file#-configuration-reference)
  <!-- TODO: move to this site/repo -->
