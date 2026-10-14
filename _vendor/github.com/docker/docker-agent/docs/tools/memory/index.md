---
title: "Memory Tool"
description: "Persistent key-value storage backed by SQLite for cross-session recall."
keywords: docker agent, ai agents, tools, toolsets, memory tool
linkTitle: "Memory"
weight: 100
canonical: https://docs.docker.com/ai/docker-agent/tools/memory/
---

_Persistent key-value storage backed by SQLite for cross-session recall._

## Overview

The memory tool provides persistent key-value storage backed by SQLite. Data survives across sessions, allowing agents to remember facts, user preferences, project context, and past decisions. Memories can be organized with categories and searched by keyword.

By default, the database is stored at `~/.cagent/memory/<config-name>/memory.db`, where `<config-name>` is derived from the loaded configuration (typically the YAML file name) and falls back to `default` when unavailable. When the agent is loaded from an OCI reference (e.g. `docker/my-agent:latest`), characters that are reserved in filesystem paths (such as `:`) are sanitised in the `<config-name>` segment — the agent's display name elsewhere is unchanged. Agents declared in the same configuration share this database by default; set an explicit `path` per toolset to isolate them.

## Available Tools

| Tool              | Description                                                                      |
| ----------------- | -------------------------------------------------------------------------------- |
| `add_memory`      | Store a new memory with optional category                                        |
| `get_memories`    | Retrieve all stored memories                                                     |
| `delete_memory`   | Delete a specific memory by ID                                                   |
| `search_memories` | Search memories by keywords and/or category (more efficient than `get_memories`) |
| `update_memory`   | Update an existing memory's content and/or category by ID                        |

## Configuration

```yaml
toolsets:
  - type: memory
```

### Options

| Property | Type   | Default                                   | Description                      |
| -------- | ------ | ----------------------------------------- | -------------------------------- |
| `path`   | string | `~/.cagent/memory/<config-name>/memory.db` | Path to the SQLite database file |

### Custom Database Path

```yaml
toolsets:
  - type: memory
    path: ./agent_memory.db
```

## Categories

Memories support an optional `category` field for organization and filtering. Common categories include:

- `preference` — User preferences and settings
- `fact` — Factual information about the project or user
- `project` — Project-specific context
- `decision` — Past decisions and their rationale

> [!TIP]
> Memory is especially useful for long-running assistants that need to recall information across conversations — like coding preferences, project conventions, or context discovered during previous sessions.
