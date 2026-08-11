---
title: "Plan Tool"
description: "Shared persistent scratchpad for multi-agent collaboration."
keywords: docker agent, ai agents, tools, toolsets, plan tool
linkTitle: "Plan"
weight: 150
canonical: https://docs.docker.com/ai/docker-agent/tools/plan/
---

_Shared persistent scratchpad for multi-agent collaboration._

## Overview

The plan tool gives agents a shared, persistent scratchpad of named documents. Any agent in a multi-agent config that loads the `plan` toolset can read and write the same plans, and those plans survive across sessions. This makes it straightforward to wire a planner agent that sketches work and one or more executor agents that consume it without any custom tool wiring.

Plans are stored as JSON files in the docker-agent data directory (`~/.cagent/plans/` by default). All agents that share a process serialize on a single mutex so concurrent reads and writes are safe. Writes are atomic (temp file + rename), so a reader never observes partial content.

## Configuration

```yaml
toolsets:
  - type: plan
```

No additional options are required. All agents that include `type: plan` in their toolsets share the same plans.

## Available Tools

| Tool                    | Description                                                                                       |
| ----------------------- | ------------------------------------------------------------------------------------------------- |
| `write_plan`            | Create or update a shared plan by name. Replaces the entire plan content — read it first to preserve what you want to keep. Each write bumps the revision number. |
| `read_plan`             | Read a shared plan by name, including its title, content, author, status, revision number, and last-updated timestamp. |
| `list_plans`            | List all shared plans with their name, title, author, status, revision, and last-updated timestamp. |
| `delete_plan`           | Delete a shared plan by name.                                                                     |
| `update_plan_from_file` | Create or update a plan, taking the new content from a file on disk instead of inline. Use it with `export_plan_to_file` to edit a large plan without re-sending its whole body. |
| `export_plan_to_file`   | Write a plan's content to a file. The content goes to disk and is **not** returned as tool output, so materialising a plan costs no tokens. |
| `set_plan_status`       | Set a plan's free-form status without rewriting its body. The plan must already exist. |
| `get_plan_status`       | Read a plan's status and current revision without fetching its body.                  |

### Cheap edits with file-based revisions

Re-sending a whole plan on every revision is expensive. The file-based tools let
an agent edit a plan without paying input-token cost for its body:

1. `export_plan_to_file` writes the current plan content to a path. The content
   is written to disk and is **not** returned.
2. The agent edits that file in place with its filesystem tools.
3. `update_plan_from_file` commits the file's new contents as the next revision.

### Free-form status

Each plan carries a free-form `status` string. There is no fixed vocabulary:
define your own in the system prompt (e.g. `idle`, `in-progress`, `blocked`,
`done`, `canceled`). Read and write it independently of the body with
`get_plan_status` and `set_plan_status`, or pass `status` to `write_plan` and
`update_plan_from_file`. The TUI surfaces the status next to the plan title.

### Optimistic locking

When several sessions edit the same plan, concurrent writes could silently
overwrite each other. Every read returns a `revision` number; pass the value you
last read as `last_known_revision` to `write_plan`, `update_plan_from_file`,
`set_plan_status`, or `delete_plan`. If the plan changed since (its current
revision no longer matches), the write is rejected with a version-conflict
error and the caller should re-read the plan and retry. Omit
`last_known_revision` to write unconditionally (last writer wins).

### Plan Names

Plan names must match the pattern `[a-z0-9][a-z0-9_-]*` (lowercase letters, digits, `-`, `_`). This is enforced structurally so two different inputs can never collapse onto the same file and path-traversal is impossible by construction.

### Plan Fields

Each plan document contains:

| Field      | Description                                               |
| ---------- | --------------------------------------------------------- |
| `name`     | The plan's unique slug name                               |
| `title`    | A short human-readable title (optional)                   |
| `content`  | The full Markdown or free-form plan text                  |
| `author`   | Free-form label identifying who last wrote the plan       |
| `status`   | Free-form lifecycle label (optional), e.g. `in-progress`  |
| `revision` | Monotonically increasing version counter, bumped on every write |
| `updatedAt`| ISO 8601 timestamp of the last write                      |

## Example

Two agents collaborate on a shared plan — the architect drafts it and the builder refines it:

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: Coordinator
    instruction: |
      Route work between the architect and the builder.
    handoffs: [architect, builder]

  architect:
    model: anthropic/claude-sonnet-4-5
    description: Drafts high-level plans
    instruction: |
      Use list_plans and read_plan to inspect existing plans, then write_plan
      to create or revise one. Always read before writing. When done, hand off
      to the builder.
    toolsets:
      - type: plan
    handoffs: [builder]

  builder:
    model: openai/gpt-4o
    description: Adds implementation steps to plans
    instruction: |
      Read the architect's plan with read_plan, then use write_plan to append
      concrete implementation steps. Always read before writing. When done,
      hand off back to root.
    toolsets:
      - type: plan
    handoffs: [root]
```

See [`examples/shared_plan.yaml`](https://github.com/docker/docker-agent/blob/main/examples/shared_plan.yaml) for a complete working example.

## Error Handling

- `read_plan` returns a distinct "not found" error when a plan does not exist, as opposed to any other I/O error, so callers can tell "plan missing" from "plan unreadable."
- `list_plans` skips corrupt entries but reports them in a `warnings` field so an agent can detect and recover from a bad state (e.g., by calling `delete_plan`).
- `delete_plan` can remove a corrupt plan to recover from a bad state.

> [!TIP]
> **Plan vs. Todo vs. Tasks**
>
> Use **plan** for shared, free-form documents that multiple agents collaborate on (design docs, requirements, work items). Use [todo](../todo/index.md) for lightweight in-session task lists. Use [tasks](../tasks/index.md) for a structured, persistent task database with priorities and dependencies.
