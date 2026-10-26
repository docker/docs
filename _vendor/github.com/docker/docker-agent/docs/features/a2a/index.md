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

# Use an agent from an OCI registry
$ docker agent serve a2a myorg/agent:tag
```

## Flags

See the [CLI reference](../cli/index.md#docker-agent-serve-a2a) for all flags, defaults, and shared runtime options.

## Authentication and network exposure

Loopback listeners may run without authentication. Non-loopback listeners require
`--auth-token` unless `--insecure-no-auth` explicitly acknowledges the exposure.
Clients must send `Authorization: Bearer <token>` for both agent-card discovery
and JSON-RPC invocation. Configure browser access with `--cors-origin`; it accepts
comma-separated literal origins or `~`-prefixed regular expressions and permits
credentials only for matching origins.

```bash
$ docker agent serve a2a ./agent.yaml --auth-token "$A2A_TOKEN" \
    --cors-origin http://localhost:3000
```

## Tool safety and migration

A2A sessions default to the `restricted` tool safety policy. Precedence is the
`--safety` flag, then agent YAML, then runtime YAML. YAML may select `strict`,
`balanced`, or `restricted`; `safety: autonomous` stops startup and directs the
operator to `--safety autonomous`. That CLI flag is the only deliberate opt-in
to autonomous tool execution.

Existing deployments should choose an explicit policy before upgrading. Migration
027 labels pre-existing sessions as `run`, so they cannot be resumed through
`/invoke`; clients must start new A2A contexts. An A2A context ID that collides
with another session is rejected without changing that session.

Downgrading to a binary that predates migration 027 fails because the session
database has a newer schema (`ErrNewerDatabase`). Restore an older database, or
use a binary that includes the migration. Revert changes without removing the
migration catalogue entry.

## Resuming conversations

Invocations with the same A2A context ID resume the stored conversation. Before
each resumed run, Docker Agent reapplies the selected agent's `max_iterations`,
`max_consecutive_tool_calls`, `max_old_tool_call_tokens` and
`max_tool_result_tokens`. This restores limits that are not persisted and uses
the currently loaded configuration rather than stale session values. Zero or
negative values retain their documented default/unlimited semantics.

Resuming retains conversation history, title and workspace. The existing safety
policy is preserved or tightened to the server's safety ceiling, never relaxed.

## Features

- **Configurable listener** — Listens on `127.0.0.1:8082` by default; use `--listen` to choose another address or port
- **Agent card** — Provides standard A2A agent metadata
- **Full Docker Agent features** — Supports all tools, models, and gateway features
- **Multiple sources** — Load agents from files or OCI registries

> [!TIP]
> **See also**
>
> For exposing agents via MCP instead, see [MCP Mode](../mcp-mode/index.md). For stdio-based integration, see [ACP](../acp/index.md). For the HTTP API, see [API Server](../api-server/index.md).

## Current Limitations

- Tool calls are handled internally, not exposed as separate A2A events
- A2A artifact support not yet integrated
- A2A memory features not yet integrated
- Multi-agent (sub-agent) scenarios need further work
