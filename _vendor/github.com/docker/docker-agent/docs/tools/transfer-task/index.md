---
title: "Transfer Task Tool"
description: "Delegate tasks to sub-agents in multi-agent setups."
keywords: docker agent, ai agents, tools, toolsets, transfer task tool
linkTitle: "Transfer Task"
weight: 80
canonical: https://docs.docker.com/ai/docker-agent/tools/transfer-task/
---

_Delegate tasks to sub-agents in multi-agent setups._

## Overview

The `transfer_task` tool allows an agent to delegate tasks to specialized sub-agents and receive their results. This is the core mechanism for multi-agent orchestration.

**You don't need to add it manually** — it's automatically available when an agent has `sub_agents` configured.

## Configuration

The tool is enabled implicitly when `sub_agents` is set:

```yaml
agents:
  coordinator:
    model: openai/gpt-4o
    description: Coordinates work across specialists
    instruction: Analyze requests and delegate to the right specialist.
    sub_agents: [developer, researcher]

  developer:
    model: anthropic/claude-sonnet-4-5
    description: Expert software developer
    instruction: Write clean, production-ready code.
    toolsets:
      - type: filesystem
      - type: shell

  researcher:
    model: openai/gpt-4o
    description: Web researcher
    instruction: Search for information online.
    toolsets:
      - type: mcp
        ref: docker:duckduckgo
```

The coordinator agent automatically gets a `transfer_task` tool that can delegate to `developer` or `researcher`.

## Tool Interface

The `transfer_task` tool takes three parameters:

| Parameter         | Type   | Required | Description                                                                                 |
| ----------------- | ------ | -------- | ------------------------------------------------------------------------------------------- |
| `agent`           | string | ✓        | Name of the sub-agent to delegate to. Must be listed under the caller's `sub_agents`.        |
| `task`            | string | ✓        | Clear, concise description of the task the sub-agent should achieve.                        |
| `expected_output` | string | ✓        | Description of the result/format the caller expects back.                                   |

The call blocks until the sub-agent returns its result, which becomes the tool's response. For non-blocking parallel delegation, use [`background_agents`](../background-agents/index.md) instead.

## Delegation Limits

Sub-agents can have `sub_agents` of their own, so multi-level delegation chains are supported. Two runtime guards keep chains sane, applied to both `transfer_task` and `run_background_agent`:

- **Cycles are rejected.** A delegation targeting an agent that is already part of the active delegation chain (for example `a -> b -> a`) fails with an error naming the cycle.
- **Depth is capped at 10 nested delegations.** The root agent delegating to its first sub-agent counts as depth 1; a call that would exceed the cap fails with an error stating the attempted depth.

A rejected delegation returns a tool error to the calling agent and never starts the sub-agent.

> [!TIP]
> **See also**
>
> For parallel task delegation, see [Background Agents](../background-agents/index.md). For multi-agent patterns, see [Multi-Agent](../../concepts/multi-agent/index.md).
