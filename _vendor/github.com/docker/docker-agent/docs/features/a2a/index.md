---
title: "A2A Protocol"
description: "Expose Docker Agent agents via Google's Agent-to-Agent (A2A) protocol for interoperability with other agent frameworks."
keywords: docker agent, ai agents, features, a2a protocol
weight: 60
canonical: https://docs.docker.com/ai/docker-agent/features/a2a/
aliases:
  - /ai/docker-agent/integrations/a2a/
---

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

## Push notifications

The agent card advertises `pushNotifications: true`, so clients can follow a
task through a webhook instead of holding a stream open or polling. Register a
webhook with the `CreateTaskPushNotificationConfig` JSON-RPC method, or inline
in `SendMessage` through `configuration.taskPushNotificationConfig`. The
server then POSTs each task status and artifact update to the webhook URL as an
A2A `StreamResponse` JSON body. The config's `token` is sent in the `A2A-Notification-Token` header
and its `authentication` as an `Authorization: Bearer` or `Basic` header.

Delivery is asynchronous and best effort: a slow webhook never delays the task
or its streaming clients, but a callback can be lost. Clients that must not
miss the outcome should fall back to `GetTask`.

- Updates reach each webhook in order. Each attempt times out after 10 seconds
  and is not retried; a timeout or non-2xx response loses that update.
- Once 64 updates are waiting for a task's webhook, new intermediate updates
  are dropped. Final updates (completed, failed, canceled, rejected,
  input-required, auth-required) are still queued and evict the oldest
  intermediate one, so the queue may grow past 64 when final updates pile up.
- Pending updates live in memory. On shutdown the server waits up to 5 seconds
  for them, then drops the rest; they are not replayed after a restart.
- Only `307` and `308` redirects within the registered origin (same scheme,
  host and port) are followed, so the POST body is kept. Any other redirect
  fails the delivery.
- Failures are logged with the webhook's scheme and host and an error category
  or status code only, never its path, query, credentials or response text.

### Webhook destinations

By default, a webhook may only reach public addresses. Connections to loopback,
private (RFC 1918, IPv6 ULA), link-local (including cloud metadata endpoints
such as `169.254.169.254`), shared (`100.64.0.0/10`), multicast and unspecified
addresses are refused. The check runs on the resolved address of every
connection, so DNS rebinding and redirects are covered too. For the same
reason, `HTTP_PROXY`, `HTTPS_PROXY` and `NO_PROXY` are ignored for webhooks:
through a proxy, only the proxy's address could be checked.

`--push-allow-private` lifts both rules: webhooks may reach any address, and
environment proxies apply. Use it only for trusted internal deployments or
local testing.

```bash
$ docker agent serve a2a ./agent.yaml --push-allow-private
```

### Persistence and restarts

Tasks and webhook registrations are stored in `<session-db>-a2a/a2a.db`, next
to the session database, and survive a restart. Webhook tokens and credentials
are stored in plaintext, so the directory is private to the user running the
server (mode `0700` on Unix, an owner-only ACL on Windows). Each session
database has its own directory, so servers using different `--session-db`
files never share tasks. A server holds an exclusive lock on that directory
while it runs: a second `serve a2a` with the same `--session-db` fails to
start instead of failing the first server's running tasks.

A task still submitted or working when the server stops cannot finish. On the
next start it is marked failed, and its webhooks are sent that final status on
the same best-effort terms as any other update. Updates lost before the stop,
including the final one of a task that had already completed, are not
re-sent: read the task with `GetTask`. Tasks waiting for input or
authentication are left unchanged.

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
