---
title: "Session Plan Tool"
description: "Per-session plan tracker for the draft, review, execute workflow."
keywords: docker agent, ai agents, tools, toolsets, session plan tool
linkTitle: "Session Plan"
weight: 160
canonical: https://docs.docker.com/ai/docker-agent/tools/session_plan/
---

_Per-session plan tracker for the "draft, review, execute" workflow._

## Overview

The `session_plan` toolset gives one agent a place to write a plan for the current session, signal that the plan is ready, and let the host route the next turn to an executing agent.

Different from the [`plan` toolset](../plan/index.md) — `plan` is for shared, named plans multiple agents collaborate on over many sessions. `session_plan` is for one ephemeral plan per session, scoped to that session by ID.

Plans live as Markdown files under:

```text
~/.cagent/session_plans/<session-id>.md
```

The tool surface is three tools:

| Tool                 | Description                                                                                          |
| -------------------- | ---------------------------------------------------------------------------------------------------- |
| `write_session_plan` | Create or replace this session's plan as markdown. There's exactly one plan per session.             |
| `read_session_plan`  | Read the plan written for the current session and return it as markdown.                             |
| `exit_plan_mode`     | Signal that the plan is ready for review. Does not switch agents on its own.                         |

## Configuration

```yaml
toolsets:
  - type: session_plan
```

No configuration options. The plan path is derived from the session ID; the agent does not name plans.

Restrict the toolset to a subset of tools the standard way:

```yaml
# An agent that consumes a plan but should not be able to (re)write or finalize one.
toolsets:
  - type: session_plan
    tools:
      - read_session_plan
```

## When to call exit_plan_mode

Call `exit_plan_mode` once the plan is complete and you do not intend to change it on the next turn. It validates that a plan exists for the session and returns a "ready for review" tool result. It does **not** switch agents or solicit user approval on its own — the host application owns the next-turn routing (for example, by reading the tool result, by a UI affordance the user toggles, or by a `handoff` declared on the agent).

This separation keeps the tool reusable across UIs: a CLI that prints tool results inline, a chat UI with a plan-mode toggle, and a server that auto-routes the next turn through a `handoff` can all consume the same signal without one stepping on another.

## Storage and cleanup

- Plans are markdown files written atomically (temp + rename), so concurrent readers — in this process or another — never observe a partial write.
- A best-effort sweep on first use of the toolset removes plan files older than 30 days under the plans directory. Stranded plans for long-gone sessions do not accumulate.
- The session ID identifies the file directly. There is no in-process mutex or revision counter, because two sessions cannot map to the same path.

## Events

A `session_plan_updated` event is emitted whenever `write_session_plan` succeeds:

```json
{
  "type": "session_plan_updated",
  "session_id": "...",
  "path": "/Users/.../.cagent/session_plans/<session-id>.md",
  "content": "# my plan\n...",
  "agent_name": "planner"
}
```

Embedders that render the plan inline can subscribe and update without re-reading the file.

## Example

A two-agent workflow: `root` executes, `planner` plans. `/plan` hands off to the planner; `exit_plan_mode` signals "ready", and the host decides what happens next.

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: Executes approved plans
    instruction: |
      You execute plans the planner has handed off. When you see a message
      that a plan has been approved, read it with read_session_plan and work
      through its steps in order.
    toolsets:
      - type: session_plan
        tools:
          - read_session_plan
      - type: filesystem
      - type: shell
    commands:
      plan:
        description: "Switch to the planner"
        agent: planner

  planner:
    model: anthropic/claude-sonnet-4-5
    description: Investigates and writes plans for review
    instruction: |
      Investigate the user's request, then write the plan with
      write_session_plan. Iterate with the user until the plan is complete,
      then call exit_plan_mode to mark it ready for review.
    toolsets:
      - type: session_plan
      - type: filesystem
        readonly: true
      - type: user_prompt
```

See [`examples/session_plan.yaml`](https://github.com/docker/docker-agent/blob/main/examples/session_plan.yaml) for a complete working example.

## Error Handling

- `read_session_plan` and `exit_plan_mode` return a "no plan written yet" error when called before `write_session_plan`.
- `write_session_plan` validates the session ID and refuses to write anything that could escape the plans directory; in practice the runtime generates UUIDs so this only triggers if an embedder supplies a hand-crafted ID.

> [!TIP]
> **session_plan vs. plan vs. todo vs. tasks**
>
> Use **session_plan** when one agent drafts an approach for the user to review before another agent executes it (ephemeral, one per session). Use [plan](../plan/index.md) for shared, named plans multiple agents collaborate on over many sessions. Use [todo](../todo/index.md) for lightweight in-session task lists. Use [tasks](../tasks/index.md) for a structured, persistent task database with priorities and dependencies.
