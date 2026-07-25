---
title: "Introduction"
description: "Docker Agent is a multi-agent runtime that lets you build, run, and share AI agents with a YAML or HCL config — no glue code required."
keywords: docker agent, ai agents, getting started, introduction
weight: 10
canonical: https://docs.docker.com/ai/docker-agent/getting-started/introduction/
---

_Docker Agent is a multi-agent runtime that lets you build, run, and share AI agents with a YAML or HCL config — no glue code required._

## What is Docker Agent?

Docker Agent is an open-source tool from Docker that orchestrates AI
agents with specialized capabilities and tools. Instead of writing
code to wire up LLMs, tools, and workflows, you **declare** your
agents in YAML or HCL — their model, personality, tools, and how they
collaborate — and Docker Agent handles the rest.

- **Multi-Agent Architecture** — build hierarchical teams of agents that specialize in different tasks and delegate work to each other.
- **Rich Tool Ecosystem** — built-in tools for files, shell, memory, and todos. Extend with any MCP server from [Docker's MCP catalog](https://hub.docker.com/u/mcp).
- **Multi-Model Support** — OpenAI, Anthropic, Google Gemini, AWS Bedrock, Docker Model Runner, and reusable provider definitions with shared defaults.
- **Package & Share** — push agents to OCI registries and pull them anywhere, just like Docker images.
- **Multiple Interfaces** — interactive TUI, headless CLI, HTTP API server, MCP mode, and A2A protocol support.
- **Security-First Design** — tool confirmation prompts, containerized MCP tools via Docker, client isolation, and resource scoping.

## Why Docker Agent?

After spending years building AI agents using various frameworks, the
Docker team kept asking the same questions:

- **How do we make building agents less of a hassle?** — Most agents
  use the same building blocks. Docker Agent provides them out of the
  box.
- **Can we reuse those building blocks?** — Declarative YAML or HCL
  configs mean you can mix and match agents, models, and tools without
  rewriting code.
- **How can we share agents easily?** — Push agents to any OCI
  registry and run them anywhere with a single command.

Docker Agent is built in the open so the community can make use of
this work and contribute to its future.

## How it works

At its core, Docker Agent follows a simple loop:

1. **You define agents** in YAML or HCL — their model, instructions, tools, and sub-agents.
2. **You run an agent** via the TUI, CLI, or API.
3. **The agent processes your request** — calling tools, delegating to sub-agents, and reasoning step by step.
4. **Results stream back in real-time** via an event-driven architecture.

```yaml
# A minimal agent definition
agents:
  root:
    model: openai/gpt-5
    description: A helpful assistant
    instruction: You are a helpful assistant.
    toolsets:
      - type: think
```

```bash
# Run it
$ docker agent run agent.yaml
```

> [!TIP]
> Jump straight to the [Quick Start](../quickstart/index.md) if you want to build your first agent right away.

## What's next?

- [**Installation**](../installation/index.md) — install Docker Agent on macOS, Linux, or Windows.
- [**Quick Start**](../quickstart/index.md) — build your first agent in under 5 minutes.
