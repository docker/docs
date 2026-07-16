---
title: "Agents"
description: "Agents are the core building blocks of docker-agent. Each agent is an AI-powered entity with a model, instructions, tools, and optional sub-agents."
keywords: docker agent, ai agents, concepts, agents
weight: 10
canonical: https://docs.docker.com/ai/docker-agent/concepts/agents/
---

_Agents are the core building blocks of docker-agent. Each agent is an AI-powered entity with a model, instructions, tools, and optional sub-agents._

## What is an Agent?

An agent in docker-agent is defined by:

- **Model** — The AI model powering it (e.g., Claude, GPT-5, Gemini). See [Models](../models/index.md).
- **Description** — A brief summary of what the agent does (used by other agents for delegation)
- **Instruction** — The system prompt that defines the agent's behavior and personality
- **Tools** — Capabilities like filesystem access, shell commands, or external APIs
- **Sub-agents** — Other agents it can delegate tasks to

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: Expert software developer
    instruction: |
      You are an expert developer. Write clean, efficient code
      and explain your reasoning step by step.
    toolsets:
      - type: filesystem
      - type: shell
      - type: think
```

## The Root Agent

Every docker-agent configuration has a **root agent** — the entry point that receives user messages. In a single-agent setup, this is the only agent. In a multi-agent setup, the root agent acts as a coordinator, delegating tasks to specialized sub-agents.

> [!NOTE]
> **Naming**
>
> The first agent defined in your YAML (or the one named `root`) is the root agent by default. You can also specify which agent to start with using `docker agent run config.yaml -a agent_name`.

## Agent Properties

| Property               | Type    | Required | Description                                                    |
| ---------------------- | ------- | -------- | -------------------------------------------------------------- |
| `model`                | string  | ✓        | Model reference (inline like `openai/gpt-5` or a named model) |
| `description`          | string  | ✓        | What the agent does — used by other agents for delegation      |
| `instruction`          | string  | ✓        | System prompt defining behavior                                |
| `toolsets`             | array   | ✗        | List of tool configurations                                    |
| `sub_agents`           | array   | ✗        | Names of agents this agent can delegate to                     |
| `fallback`             | object  | ✗        | Fallback model configuration for resilience                    |
| `add_date`             | boolean | ✗        | Include current date in context                                |
| `add_environment_info` | boolean | ✗        | Include OS, working directory, git info in context             |
| `max_iterations`       | int     | ✗        | Max tool-calling loops (default: unlimited)                    |
| `commands`             | object  | ✗        | Named prompts callable via `/command`                          |
| `skills`               | boolean \| list | ✗    | Enable skill discovery and loading. `true` = `["local"]`; list values may combine `"local"` with remote skill-server URLs. |

## Model Fallbacks

Agents can automatically fail over to alternative models when the primary model is unavailable:

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    fallback:
      models:
        - openai/gpt-5
        - google/gemini-3.5-flash
      retries: 2 # retries per model for 5xx errors
      cooldown: 1m # stick with fallback after 429
```

## Named Commands

Define reusable prompts that can be invoked as commands:

```yaml
agents:
  root:
    model: openai/gpt-5
    instruction: You are a helpful assistant.
    commands:
      df: "Check how much free space I have on my disk"
      greet: "Say hello to ${env.USER}"
```

```bash
# Run a named command
$ docker agent run agent.yaml /df
$ docker agent run agent.yaml /greet
```

Commands support environment variable interpolation using JavaScript template literal syntax. Undefined variables expand to empty strings.

## Default Agent

Running `docker agent run` without a config file uses a built-in default agent. This is a capable general-purpose agent for quick tasks without needing any configuration.

```bash
# Use the default agent
$ docker agent run

# Override the default with an alias
$ docker agent alias add default /path/to/my-agent.yaml
$ docker agent run  # now runs your custom agent
```

> [!TIP]
> **See also**
>
> For reusable task-specific instructions, see [Skills](../../features/skills/index.md). For multi-agent patterns, see [Multi-Agent](../multi-agent/index.md). For full config reference, see [Agent Config](../../configuration/agents/index.md).
