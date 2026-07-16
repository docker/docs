---
title: "Quick Start"
description: "Get up and running with Docker Agent in under 5 minutes. Pick whichever path suits you best."
keywords: docker agent, ai agents, getting started, quick start
weight: 30
canonical: https://docs.docker.com/ai/docker-agent/getting-started/quickstart/
aliases:
  - /ai/docker-agent/tutorial/
---

_Get up and running with Docker Agent in under 5 minutes. Pick whichever path suits you best._

## Option A: Run the Default Agent

The fastest way to try Docker Agent — no config file needed:

```bash
# Launch the default agent with the interactive TUI
$ docker agent run
```

This starts a general-purpose assistant with sensible defaults. Just start chatting.

> [!NOTE]
> This needs a model: a cloud provider API key, or a local model pulled through Docker Model Runner. [Set Up a Model](../set-up-a-model/index.md) walks through both paths.

## Option B: Run a Pre-Built Agent from the Registry

Try a ready-made agent from the [agent catalog](https://hub.docker.com/u/agentcatalog) — no YAML needed:

```bash
# Run a pirate-themed assistant
$ docker agent run agentcatalog/pirate

# Run a coding agent
$ docker agent run agentcatalog/coder
```

## Option C: Generate a Config Interactively

Use the `docker agent new` command to scaffold a config file through prompts:

```bash
# Interactive wizard
$ docker agent new

# Or specify options directly
$ docker agent new --model openai/gpt-5

# Override iteration limits
$ docker agent new --model dmr/ai/gemma3-qat:12B --max-iterations 15
```

This generates an `agent.yaml` in the current directory. Then run it:

```bash
$ docker agent run agent.yaml
```

## Option D: Write Your Own Config

Create an `agent.yaml` by hand for full control. Here's a minimal example:

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: A helpful coding assistant
    instruction: |
      You are an expert software developer. Help users write
      clean, efficient code. Explain your reasoning.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```

This gives your agent:

- **Claude Sonnet 4.5** as the underlying model
- **Filesystem access** to read and write files
- **Shell access** to run commands
- **Think tool** for step-by-step reasoning

This config needs `ANTHROPIC_API_KEY` set. See [Set Up a Model](../set-up-a-model/index.md) for where to put it, or use a local `dmr/...` model that needs no key.

```bash
# Launch the interactive terminal UI
$ docker agent run agent.yaml
```

> [!TIP]
> **Prefer HCL?**
>
> You can write the same config as `agent.hcl` using labeled blocks and heredocs. See [HCL Configuration](../../configuration/hcl/index.md).

## Try It Out

Once your agent is running, try asking it to:

- _"List the files in the current directory"_
- _"Create a Python script that fetches weather data"_
- _"Explain what the code in main.go does"_

> [!TIP]
> Add `--yolo` to auto-approve all tool calls: `docker agent run agent.yaml --yolo`

## Non-Interactive Mode

Use `docker agent run --exec` for one-shot tasks:

```bash
# Ask a single question
$ docker agent run --exec agent.yaml "Create a Dockerfile for a Node.js app"

# Pipe input
$ cat error.log | docker agent run --exec agent.yaml "What's wrong in this log?"
```

## Add More Power

Give your agent persistent memory and web search:

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: Research assistant with memory
    instruction: |
      You are a research assistant. Search the web for information,
      remember important findings, and provide thorough analysis.
    toolsets:
      - type: think
      - type: memory
        path: ./research.db
      - type: mcp
        ref: docker:duckduckgo
```

> [!NOTE]
> **Docker MCP Tools**
>
> The `ref: docker:duckduckgo` syntax runs the DuckDuckGo MCP server in a Docker container. This is the recommended way to use MCP tools — secure, isolated, and easy to configure. Requires Docker Desktop.

## What's Next?

- [**Understand Agents**](../../concepts/agents/index.md) — learn how agents work and what you can configure.
- [**Multi-Agent Systems**](../../concepts/multi-agent/index.md) — build teams of collaborating agents.
- [**Configuration Reference**](../../configuration/overview/index.md) — full reference for all YAML and HCL options.
- [**Troubleshooting**](../../community/troubleshooting/index.md) — something not working? Debug tips and common fixes.
