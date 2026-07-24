---
title: "A2A Protocol"
description: "Expose Docker Agent agents via Google's Agent-to-Agent (A2A) protocol for interoperability with other agent frameworks."
keywords: docker agent, ai agents, features, a2a protocol
weight: 60
canonical: https://docs.docker.com/ai/docker-agent/features/a2a/
aliases:
  - /ai/docker-agent/integrations/a2a/
---

_Expose Docker Agent agents via Google's Agent-to-Agent (A2A) protocol for interoperability with other agent frameworks._

## Overview

The `docker agent serve a2a` command starts an A2A server that exposes your agents using the [A2A protocol](https://a2a-protocol.org/latest/). This enables communication between Docker Agent and other agent frameworks that support A2A.

> [!WARNING]
> **Early support**
>
> A2A support is functional but still evolving. Tool calls, artifacts, and memory features have limited A2A integration. See limitations below.

## Usage

```bash
# Start A2A server for an agent
$ docker agent serve a2a ./agent.yaml

# Specify a custom address
$ docker agent serve a2a ./agent.yaml --listen 127.0.0.1:9000

# Use an agent from the catalog
$ docker agent serve a2a agentcatalog/pirate
```

## Flags

| Flag                              | Default          | Description                                                                                                          |
| --------------------------------- | ---------------- | -------------------------------------------------------------------------------------------------------------------- |
| `-l, --listen <addr>`             | `127.0.0.1:8082` | Address to listen on.                                                                                                |
| `-a, --agent <name>`              | (first agent)    | Name of the agent to expose when the config contains multiple agents. Defaults to the team's first agent.            |
| `-s, --session-db <path>`         | `<data-dir>/session.db` | Path to the SQLite session database.                                                                          |
| `--working-dir <path>`            | current dir      | Working directory the agent runs in.                                                                                 |
| `--env-from-file <file>`          | (none)           | Load additional environment variables from a `.env` file (repeatable).                                               |
| `--models-gateway <url>`          | (none)           | Route all provider traffic through a models gateway URL.                                                             |
| `--code-mode-tools`               | `false`          | Expose tools as a single "code" toolset that accepts a JavaScript snippet to run.                                    |
| `--hook-pre-tool-use <cmd>`       | (none)           | Add a pre-tool-use hook (repeatable). See [Hooks](../../configuration/hooks/index.md).                     |
| `--hook-post-tool-use <cmd>`      | (none)           | Add a post-tool-use hook (repeatable).                                                                               |
| `--hook-session-start <cmd>`      | (none)           | Add a session-start hook (repeatable).                                                                               |
| `--hook-session-end <cmd>`        | (none)           | Add a session-end hook (repeatable).                                                                                 |
| `--hook-on-user-input <cmd>`      | (none)           | Add an on-user-input hook (repeatable).                                                                              |
| `--hook-stop <cmd>`               | (none)           | Add a stop hook, fired when the model finishes responding (repeatable).                                              |

## Features

- **Auto port selection** — Picks an available port if not specified
- **Agent card** — Provides standard A2A agent metadata
- **Full Docker Agent features** — Supports all tools, models, and gateway features
- **Multiple sources** — Load agents from files or the agent catalog

> [!TIP]
> **See also**
>
> For exposing agents via MCP instead, see [MCP Mode](../mcp-mode/index.md). For stdio-based integration, see [ACP](../acp/index.md). For the HTTP API, see [API Server](../api-server/index.md).

## Current Limitations

- Tool calls are handled internally, not exposed as separate A2A events
- A2A artifact support not yet integrated
- A2A memory features not yet integrated
- Multi-agent (sub-agent) scenarios need further work
