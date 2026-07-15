---
title: "Background Agents Tool"
description: "Dispatch work to sub-agents concurrently and collect results asynchronously."
keywords: docker agent, ai agents, tools, toolsets, background agents tool
linkTitle: "Background Agents"
weight: 90
canonical: https://docs.docker.com/ai/docker-agent/tools/background-agents/
---

_Dispatch work to sub-agents concurrently and collect results asynchronously._

## Overview

The background agents tool lets an orchestrator dispatch work to sub-agents concurrently and collect results asynchronously. Unlike [transfer_task](../transfer-task/index.md) (which blocks until the sub-agent finishes), background agent tasks run in parallel — the orchestrator can start several tasks, do other work, and check on them later.

## Available Tools

| Tool                     | Description                                                     |
| ------------------------ | --------------------------------------------------------------- |
| `run_background_agent`   | Start a sub-agent task in the background; returns a task ID     |
| `list_background_agents` | List all background tasks with their status and runtime         |
| `view_background_agent`  | View live output or final result of a task by ID                |
| `stop_background_agent`  | Cancel a running task by ID                                     |

### `run_background_agent` parameters

| Parameter         | Type   | Required | Description                                                                 |
| ----------------- | ------ | -------- | --------------------------------------------------------------------------- |
| `agent`           | string | ✓        | Name of the sub-agent to run. Must be listed under the caller's `sub_agents`. |
| `task`            | string | ✓        | Clear, concise description of the task the sub-agent should achieve.        |
| `expected_output` | string | ✗        | Optional description of the result format the caller expects.               |

`run_background_agent` returns a **task ID** string. Tools run by the sub-agent are pre-approved, so only dispatch to trusted sub-agents with well-scoped tasks.

### `view_background_agent` and `stop_background_agent` parameters

| Parameter | Type   | Required | Description                                                    |
| --------- | ------ | -------- | -------------------------------------------------------------- |
| `task_id` | string | ✓        | Task ID returned by `run_background_agent` or `list_background_agents`. |

`list_background_agents` takes no parameters.

## Configuration

```yaml
toolsets:
  - type: background_agents
```

No configuration options. Requires the agent to have `sub_agents` configured so the background tasks have agents to dispatch to.

## Example

```yaml
agents:
  coordinator:
    model: openai/gpt-4o
    description: Orchestrates parallel research
    instruction: Fan out research tasks and synthesize results.
    sub_agents: [researcher]
    toolsets:
      - type: background_agents
      - type: think

  researcher:
    model: openai/gpt-4o
    description: Web researcher
    instruction: Research topics thoroughly.
    toolsets:
      - type: mcp
        ref: docker:duckduckgo
```

> [!TIP]
> **When to Use**
>
> Use `background_agents` when your orchestrator needs to fan out work to multiple specialists in parallel — for example, researching several topics simultaneously or running independent code analyses side by side.

In the TUI, each background task's token usage is accounted for live: the sidebar's Agents panel shows the sub-agent's context usage percentage on its roster row, the Agent Inspector shows its exact token counts, and the task's cost joins the session total.

## Using Harness Sub-Agents

Background agents work equally well with [harness-backed sub-agents](../../features/harnesses/index.md) — sub-agents driven by external coding CLIs such as Claude Code or Codex. This lets you dispatch multiple independent coding tasks in parallel:

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: Orchestrator that fans out coding tasks
    instruction: |
      Dispatch the frontend and backend tasks in parallel,
      then collect results and produce a summary.
    sub_agents:
      - claude-coder
      - codex-coder
    toolsets:
      - type: background_agents

  claude-coder:
    description: Frontend specialist (Claude Code)
    harness:
      type: claude-code
      effort: medium

  codex-coder:
    description: Backend specialist (Codex)
    harness:
      type: codex
```

The orchestrator calls `run_background_agent` for each coding task, then uses `list_background_agents` and `view_background_agent` to collect results when they finish.

> [!NOTE]
> **Harness toolsets are ignored**
>
> Harness agents use the external CLI's own tools — any `toolsets:` configured on the harness agent are silently ignored. See [Coding Harnesses](../../features/harnesses/index.md) for details and caveats.

See [`examples/coding_harness_background_agents.yaml`](https://github.com/docker/docker-agent/blob/main/examples/coding_harness_background_agents.yaml) for a complete configuration.
