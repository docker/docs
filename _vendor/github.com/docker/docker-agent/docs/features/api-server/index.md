---
title: "API Server"
description: "Expose your agents via an HTTP API for programmatic access, web frontends, and integrations."
keywords: docker agent, ai agents, features, api server
weight: 80
canonical: https://docs.docker.com/ai/docker-agent/features/api-server/
---

_Expose your agents via an HTTP API for programmatic access, web frontends, and integrations._

## Overview

The `docker agent serve api` command starts an HTTP server that exposes your agents through a REST-style API with Server-Sent Events (SSE) streaming. Use it to build web UIs, integrate with CI/CD pipelines, or connect agents to other services.

```bash
# Start the API server
$ docker agent serve api agent.yaml

# Custom listen address
$ docker agent serve api agent.yaml --listen 0.0.0.0:8080 --auth-token "$AGENT_TOKEN"

# With session persistence
$ docker agent serve api agent.yaml --session-db ./sessions.db

# Auto-refresh from OCI registry every 10 minutes
$ docker agent serve api myorg/coder --pull-interval 10
```

> [!TIP]
> **When to use API server vs. chat server**
>
> Use the **API server** when you want full control over sessions, agent execution, tool-call confirmations, and streamed runtime events — this is Docker Agent's native protocol. Use the [Chat Server](../chat-server/index.md) when you want to plug Docker Agent into existing OpenAI-compatible tooling (chat UIs, IDE integrations, OpenAI SDK clients) instead.

## Endpoints

All endpoints are under the `/api` prefix.

### Agents

| Method | Path              | Description                       |
| ------ | ----------------- | --------------------------------- |
| `GET`  | `/api/agents`     | List all available agents         |
| `GET`  | `/api/agents/:id` | Get an agent's full configuration |

Each agent entry in the `GET /api/agents` response contains:

| Field        | Type            | Description                                                                                   |
| ------------ | --------------- | --------------------------------------------------------------------------------------------- |
| `name`       | string          | Agent identifier (config filename without `.yaml`).                                           |
| `description`| string          | The root agent's `description` field.                                                         |
| `multi`      | boolean         | `true` when the config defines more than one agent.                                           |
| `commands`   | array of string | Sorted list of named command keys defined on the root agent. Omitted when no commands exist.  |

### Remote agent sources

For an agent loaded from a remote HTTP(S) configuration source, endpoints that need to load that source return `502 Bad Gateway` when fetching it fails. A missing configured agent returns `404 Not Found`; invalid source URLs or configuration return `500 Internal Server Error`.

### Sessions

| Method   | Path                                | Description                                             |
| -------- | ----------------------------------- | ------------------------------------------------------- |
| `GET`    | `/api/sessions`                     | List all sessions. Pass `?active=true` to return only runtimes attached to this server, with lightweight `working_dir` and `streaming` status and no session-history read. |
| `POST`   | `/api/sessions`                     | Create a new session. Accepts an optional `title` field — when set, it is stored and LLM title generation is skipped. |
| `GET`    | `/api/sessions/:id`                 | Get a session by ID (messages, tokens, permissions)     |
| `GET`    | `/api/sessions/:id/status`          | Lightweight runtime state (streaming, title, agent, tokens). Requires an attached runtime. |
| `GET`    | `/api/sessions/:id/snapshot`        | Full state in one call (stored fields + runtime state + `last_event_seq`) for gapless resync — see [Reconnecting without gaps](#reconnecting-without-gaps). |
| `GET`    | `/api/sessions/:id/events`          | Live session event stream (SSE) with sequence numbers and replay. Available for a run attached via [`--listen`](#listen), or once a session has raised at least one out-of-band event (e.g. a background job's elicitation or an idle-session recall), which creates a session-scoped event log on demand carrying such out-of-band events — see [Session event stream](#session-event-stream-and-reconnection) for what each kind of log contains. |
| `DELETE` | `/api/sessions/:id`                 | Delete a session                                        |
| `PATCH`  | `/api/sessions/:id/title`           | Update session title                                    |
| `PATCH`  | `/api/sessions/:id/permissions`     | Update session permissions                              |
| `POST`   | `/api/sessions/:id/fork`            | Fork a session at a user message — creates a new session with messages `[0, message_index)` of the parent (see [Session Forking](#session-forking)) |
| `POST`   | `/api/sessions/:id/messages`        | Append a message directly to a session's history (bypasses the model). Returns `409 Conflict` while the session has an active run (see [Agent Execution](#agent-execution)). |
| `PATCH`  | `/api/sessions/:id/messages/:msg_id` | Update an existing message by ID. Returns `409 Conflict` while the session has an active run. |
| `POST`   | `/api/sessions/:id/resume`          | Resume a paused session (after tool confirmation)       |
| `POST`   | `/api/sessions/:id/tools/toggle`    | Toggle auto-approve (YOLO) mode                         |
| `POST`   | `/api/sessions/:id/elicitation`     | Respond to an MCP tool elicitation request. Pass the `elicitation_id` from the `elicitation_request` event to target a specific concurrent request; omitted, it resolves the sole pending one. |
| `POST`   | `/api/sessions/:id/steer`           | Inject messages into a running turn (pre-empts current) |
| `POST`   | `/api/sessions/:id/followup`        | Enqueue messages to run after the current turn finishes (supports an `Idempotency-Key` — see [Idempotent follow-ups](#idempotent-follow-ups)). |
| `GET`    | `/api/sessions/:id/models`          | List available models for the session's current agent   |

### Agent Execution

| Method | Path                                       | Description                                                                          |
| ------ | ------------------------------------------ | ------------------------------------------------------------------------------------ |
| `POST` | `/api/sessions/:id/agent/:agent`           | Run the root agent for a session (SSE stream)                                        |
| `POST` | `/api/sessions/:id/agent/:agent/:name`     | Run a specific named agent (SSE stream)                                              |
| `GET`  | `/api/agents/:id/:agent_name/tools/count`  | Count tools currently available to `:agent_name` (accounts for deferred toolsets).   |

**Path parameters:**

- **`:agent`** — The agent identifier, which is the **config filename without the `.yaml` extension**. This must match the filename passed to `docker agent serve api`. For example, if you start the server with `docker agent serve api my-assistant.yaml`, the agent identifier is `my-assistant`. When serving a directory of YAML files, each file becomes a separate agent identified by its filename without the extension.
- **`:name`** _(optional)_ — The name of a specific sub-agent defined in a multi-agent configuration. If omitted, the request targets the `root` agent. For example, in a config that defines agents named `root`, `coder`, and `reviewer`, use `/api/sessions/:id/agent/my-config/coder` to run the `coder` sub-agent directly.

**Examples:**

```bash
# Single-agent config: my-assistant.yaml
# Start: docker agent serve api my-assistant.yaml
# Run the root agent:
curl -N -X POST http://localhost:8080/api/sessions/$SID/agent/my-assistant \
  -H "Content-Type: application/json" \
  -d '{"messages":[{"role": "user", "content": "Hello!"}]}'

# Multi-agent config: team.yaml (defines agents: root, coder, reviewer)
# Start: docker agent serve api team.yaml
# Run the root agent:
curl -N -X POST http://localhost:8080/api/sessions/$SID/agent/team \
  -H "Content-Type: application/json" \
  -d '{"messages":[{"role": "user", "content": "Review this PR"}]}'

# Run a specific sub-agent (reviewer):
curl -N -X POST http://localhost:8080/api/sessions/$SID/agent/team/reviewer \
  -H "Content-Type: application/json" \
  -d '{"messages":[{"role": "user", "content": "Review this PR"}]}'
```

### Execution budgets

API-created native runtimes honor the manifest's `budget` and agent-selected
`budgets`. Consumption accumulates across requests and delegated work in the
same live runtime. Streams report `budget_usage` and `budget_exceeded`; once a
wallet is exhausted, further model calls charged to it stop without a
continuation prompt. Configurations without budgets remain unlimited.

Counters reset when the runtime is reconstructed, including after a server
restart. Stored transcript usage and cost totals are retained but do not seed
the new wallet. This is an execution limit, not a persistent quota, and does not
apply to external `harness:` execution. See [Budget](../../configuration/budget/index.md).

### Execution capabilities

`team_info.available_agents` entries carry optional `execution_capabilities`.
Local startup/refresh `agent_info` events include the same profile. The
Agent Inspector displays it without starting tools or harness drivers.

```json
{
  "mode": "native",
  "tool_execution": "docker-agent",
  "tool_approval": "docker-agent-policy",
  "budget_supported": true,
  "compaction_supported": true,
  "mid_turn_steering_supported": true,
  "prompt_hook_context_supported": true
}
```

These are integration capabilities, **not enabled settings or security
attestations**. Native approval can allow a call without prompting. Budget
support requires the caller to install budget configuration; CLI and HTTP API
currently do so, while ACP, MCP and A2A do not forward manifest wallets.
Compaction support does not mean automatic compaction is enabled. Profiles
refer to each agent's own execution: a native orchestrator's controls do not
cover delegated external harness actions.

Harness profiles use `mode: harness`, `tool_execution: harness-reported` and
`tool_approval: external`, with the four support flags false. Reported harness
tools have already run outside docker-agent's approval pipeline; `external`
does not guarantee the harness prompts or enforces approvals. Harness lifecycle
and before-LLM hooks still run, but native transient prompt-hook context is not
forwarded into the harness's conversation.

An absent profile means **unknown**, including on older servers and
manifest-only remote inspection. Use `team_info` as the roster authority;
model-only `agent_info` updates may omit the profile and must not clear a
previously received one. Remote clients forward server metadata without
inferring capabilities from model names or manifests.

### Health

| Method | Path        | Description                               |
| ------ | ----------- | ----------------------------------------- |
| `GET`  | `/api/ping` | Health check — returns `{"status": "ok"}` |

### OAuth

| Method | Path                     | Description                                                                                                                                                                                                                                          |
| ------ | ------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `POST` | `/api/mcp-oauth/callback` | Deliver an OAuth deeplink callback to a pending unmanaged OAuth flow. Success path: `?state=<state>&code=<code>`; authorization-server error path: `?state=<state>&error=<error>&error_description=<desc>`. Returns 400 if `state` is missing or neither `code` nor `error` is provided; 404 if no flow is awaiting that `state`. See [Remote MCP OAuth](../remote-mcp/index.md) for details. |

## Streaming Responses

The agent execution endpoints (`POST /api/sessions/:id/agent/:agent`) return **Server-Sent Events (SSE)**. The request body is a JSON object with a `messages` array and an optional `model` field. Setting `model` applies a persistent per-agent override on the session before the turn starts (subsequent turns reuse it). An empty or omitted `model` leaves the existing override untouched. (Each event is a JSON object representing a runtime event — remember that `:agent` is the config filename without the `.yaml` extension.):

```bash
# Send a message and stream the response
# (assuming the server was started with: docker agent serve api my-agent.yaml)
$ curl -N -X POST http://localhost:8080/api/sessions/$SID/agent/my-agent \
  -H "Content-Type: application/json" \
  -d '{"messages":[{"role": "user", "content": "Hello!"}]}'

# Same call, but switch the agent's model for this turn (and persist it):
$ curl -N -X POST http://localhost:8080/api/sessions/$SID/agent/my-agent \
  -H "Content-Type: application/json" \
  -d '{"messages":[{"role":"user","content":"Hello!"}],"model":"openai/gpt-4o"}'

# Response (SSE stream):
data: {"type":"stream_started","session_id":"...","agent_name":"root"}
data: {"type":"agent_choice","content":"Hello! How","agent_name":"root"}
data: {"type":"agent_choice","content":" can I help","agent_name":"root"}
data: {"type":"agent_choice","content":" you today?","agent_name":"root"}
data: {"type":"stream_stopped","session_id":"...","agent_name":"root"}
```

Streaming deltas (`agent_choice`, `agent_choice_reasoning`, `partial_tool_call` and `tool_call_output`) carry only the new bytes, so clients must concatenate them. The server merges consecutive deltas for the same message or tool call that arrive within `--stream-coalesce-window` (50 ms by default) into a single event. This applies to this endpoint and to `GET /api/sessions/:id/events`. Merged events carry the same content in the same order. Only the chunk boundaries change, so don't rely on one event per model token. The first delta after a quiet spell is sent at once, and any other event is sent only after the deltas that preceded it.

The `stream_stopped` event may include `finish_reason` with the terminal model result (such as `stop`, `length`, or `refusal`), even when that response has no content or usage. Its stream-level `reason` is `max_iterations` when execution stops at the iteration cap. These describe the root or sub-session identified by `session_id`; consumers should not combine child outcomes with the root result. Stream-stop delivery precedes session-end hooks, so draining the stream to closure remains the completion barrier.

Event types include:

- `stream_started` / `stream_stopped` — Agent execution lifecycle
- `agent_choice` — Streamed text content (partial responses)
- `agent_choice_reasoning` — Streamed reasoning content
- `partial_tool_call` — Streamed tool-call arguments, keyed by `tool_call.id`. The first one also carries `tool_definition`
- `tool_call` — Agent requesting tool execution
- `tool_call_output` — Streamed output of a running tool, keyed by `tool_call_id`
- `tool_call_confirmation` — Tool call waiting for user approval
- `tool_call_response` — Tool execution result
- `plan_changed` — A shared plan was created, updated, or deleted through the plan toolset. The payload carries the plan's `scope`, `name`, `action`, and `version` — never its content. Shared plans are deliberately process-global: every active stream served by the same process subscribes to the same shared plan notifier and receives the event regardless of which session performed the mutation, and the payload does not identify the mutating session (read the plan's `author` metadata for collaborative attribution).
- `error` — Error during execution

## Typical Workflow

1. **List agents** — `GET /api/agents` to discover available agents
2. **Create session** — `POST /api/sessions` to start a conversation
3. **Send message** — `POST /api/sessions/:id/agent/:agent` with user messages
4. **Stream response** — Read SSE events as the agent processes
5. **Handle confirmations** — If a tool call needs approval, `POST /api/sessions/:id/resume`
6. **Continue** — Send follow-up messages to the same session

```bash
# 1. List available agents
$ curl http://localhost:8080/api/agents
[{"name":"my-agent","multi":false,"description":"A helpful assistant","commands":["deploy","review"]}]

# 2. Create a session
$ curl -X POST http://localhost:8080/api/sessions \
  -H "Content-Type: application/json" -d '{}'
{"id":"abc-123","title":"","created_at":"..."}

# Create a session with a pre-supplied title (skips LLM title generation)
$ curl -X POST http://localhost:8080/api/sessions \
  -H "Content-Type: application/json" -d '{"title":"deploy check"}'
{"id":"def-456","title":"deploy check","created_at":"..."}
# title preserved; LLM title generation skipped

# 3. Run the agent with a message
$ curl -N -X POST http://localhost:8080/api/sessions/abc-123/agent/my-agent \
  -H "Content-Type: application/json" \
  -d '{"messages":[{"role":"user","content":"What files are in the current directory?"}]}'
```

## HTTP security

The REST API rejects cross-origin browser requests using Fetch Metadata and
Origin checks, including bodyless approval toggles and OAuth callback posts.
Reads are protected too: fetching an agent's tool count can start tool processes.
JSON endpoints require `Content-Type: application/json` (optional media-type
parameters are accepted); form and `text/plain` bodies are rejected with HTTP
415. Native clients without browser headers remain supported.

Unauthenticated TCP listeners accept only loopback IP addresses and `localhost`
in the Host header to prevent DNS rebinding. Unix sockets and named pipes retain
support for synthetic hostnames and rely on local access permissions. All
non-loopback TCP listeners require `--auth-token`; clients must send it as
`Authorization: Bearer <token>`. The `/health` and `/ready` probes remain exempt
from bearer authentication. Reverse proxies must preserve these protections;
forwarded headers do not authorize browser requests. For network deployments,
use TLS or a trusted TLS-terminating proxy to protect the bearer token in transit.

### Connecting with the CLI

Use `--remote-auth-token` with `run --remote`, or set
`DOCKER_AGENT_REMOTE_AUTH_TOKEN` to avoid putting the token in command-line
arguments. The flag takes precedence over the environment variable.

```bash
$ export DOCKER_AGENT_REMOTE_AUTH_TOKEN='<server-token>'
$ docker agent run --remote https://agent.example.com agent.yaml
```

## CLI Flags

```bash
docker agent serve api <agent-file>|<agents-dir> [flags]
```

See the [CLI reference](../cli/index.md#docker-agent-serve-api) for all flags, defaults, and shared runtime options.

> [!WARNING]
> `--auth-token` is required for non-loopback TCP listeners, including inherited TCP listeners. For multi-user or network-exposed deployments, also set `--session-workingdir-root` to confine the host directories callers can use for sessions.

> [!NOTE]
> **What `--max-request-size` does and doesn't cover**
>
> This is a finite, process-wide cap on one serialized inbound HTTP request body — it isn't a model context-window limit, and raising it doesn't increase what a provider/model accepts or how large a local attachment/prompt file can be. A larger cap also means the server buffers more memory per request from an unauthenticated or malicious client, so weigh that against your deployment's exposure. If a reverse proxy or gateway sits in front of this server, it may enforce its own, lower cap regardless of this flag. See [Troubleshooting: HTTP 413](../../guides/troubleshooting/index.md#http-413-request-body-too-large) for full diagnosis.

> [!TIP]
> **Live profiling (advanced)**
>
> For production diagnostics, set the `DOCKER_AGENT_PPROF_ADDR` environment variable (legacy `CAGENT_PPROF_ADDR` is also supported), or use the hidden `--pprof-addr` flag, with a TCP address such as `127.0.0.1:6060`. Docker Agent will start a Go pprof HTTP server at `/debug/pprof/`, which you can query with `go tool pprof`. Use a loopback address — a non-loopback binding logs a security warning. This flag is intentionally hidden from `--help`.

> [!TIP]
> **Multi-agent configs**
>
> You can point `docker agent serve api` at a directory containing multiple agent YAML files. Each becomes a separate agent accessible via `/api/agents`. Combine with `--pull-interval` to auto-refresh agents from an OCI registry.

## Session Persistence

Sessions are stored in a SQLite database (default: `session.db` in the current directory). This means:

- Sessions survive server restarts
- Multiple server instances can share a database
- Use `--session-db` to specify a custom path

## Tool Call Approval

By default, tool calls require approval. In the API workflow:

1. Agent makes a tool call → server emits a `tool_call_confirmation` event
2. Client reviews and sends `POST /api/sessions/:id/resume` with the decision
3. Execution continues based on approval/denial

Toggle auto-approve with `POST /api/sessions/:id/tools/toggle` for automated workflows.

## Driving a running TUI with `--listen` {#listen}

The same session API can be exposed by an **interactive run** so an external
process can drive it — send follow-up prompts, observe progress, read the
title — without scraping the terminal. Start a normal run and add `--listen`:

```bash
# Expose this run's control plane on a TCP port...
$ docker agent run agent.yaml --listen 127.0.0.1:8080

# ...or on a unix socket (no port to allocate; access is gated by file
# permissions). npipe:// (Windows) and fd:// are also accepted.
$ docker agent run agent.yaml --listen unix:///tmp/my-run.sock
```

The run keeps its interactive TUI; the control plane runs alongside it. A
follow-up delivered over HTTP is processed exactly as if it had been typed
into the TUI: it starts a turn even when the agent is idle, generates the
session title on the first turn, and streams the resulting events to both the
terminal and every connected API client.

```bash
# Send a follow-up to the attached run (SID is the --session id):
$ curl -X POST http://127.0.0.1:8080/api/sessions/$SID/followup \
    -H 'Content-Type: application/json' \
    -d '{"messages":[{"content":"Now add tests"}]}'
```

> [!NOTE]
> **Discovering a run**
>
> Each run started with `--listen` writes a discovery record to `<data-dir>/runs/<pid>.json` containing its address and initial session id, so a supervising process can find a live run by session id, pid, or address. TUI tabs opened later are separate sessions attached to the same control plane; use `GET /api/sessions?active=true` on that address to enumerate them and read their `streaming` state without loading session history.

> [!WARNING]
> **This control plane has a fixed 1 MiB request-body cap and no built-in authentication**
>
> Unlike the standalone `docker agent serve api` process above, an attached run's `--listen` control plane exposes neither `--max-request-size` nor `--auth-token`: every request is capped at a fixed, non-configurable 1 MiB body (returning HTTP 413 above it — see [Troubleshooting: HTTP 413](../../guides/troubleshooting/index.md#http-413-request-body-too-large)), and there is no bearer token to require. If you need a configurable cap, built-in bearer auth, or a network-reachable listener, run a standalone `docker agent serve api` deployment instead. `--listen` only supports loopback TCP addresses or local sockets; non-loopback TCP listeners are rejected. An authenticating proxy can connect to a local listener.

## Session event stream and reconnection

The Go remote client reconnects sequenced `/events` streams from the last received
ID. After a replay gap, the remote TUI waits for an idle snapshot and recovers
missing plain assistant text by message ID. Tool results, reasoning, and transient
interaction events are not reconstructed; a warning reports this limitation.
Ambiguous histories fail visibly instead of replaying answers twice.

Foreground POST run streams are not automatically resubmitted after interruption:
retrying a run could repeat tool side effects. A stream ending without a root
`stream_stopped` event reports an incomplete-response error. Streaming requests
are not subject to the metadata client's 30-second total timeout.

`GET /api/sessions/:id/events` is a **Server-Sent Events** stream of the
session's runtime events — `stream_started`, `agent_choice`, `tool_call`,
`session_title`, `token_usage`, `stream_stopped`, and so on. Unlike the
per-request stream returned by the agent-execution endpoint, it is
session-scoped and survives across turns, so a client can watch a session for
its whole lifetime. It is available for a run attached via
[`--listen`](#listen), and — since a session-scoped event log is created on
demand the first time a session runs an idle recall or raises an out-of-band event, such as an
`elicitation_request` from a background job — for any API-created session
that has produced at least one (see the
[Sessions endpoint table](#sessions) above). The two kinds of log differ in
coverage: a `--listen` run feeds its full runtime event stream into the log,
while an on-demand log for an API-created session carries the session's
out-of-band events — not necessarily all ordinary turn events, which flow on
the per-request SSE stream of the [agent-execution](#agent-execution)
request that runs the turn.

Each event carries a monotonic **sequence number** in the SSE `id:` field, and
the server buffers recent events. This makes the stream resumable:

- **Resume after a drop** — reconnect with the standard `Last-Event-ID` header
  (sent automatically by browser `EventSource` clients) or a `?since=<seq>`
  query parameter. Buffered events newer than that point are replayed before
  live tailing resumes, so nothing is missed.
- **Gap signal** — if the resume point has already fallen out of the buffer, the
  server sends a single `{"type":"gap"}` event (with no id) before the replay.
  The client should re-fetch the snapshot to resync, then continue tailing.
- **End of session** — when the session is ended server-side (for example via
  `DELETE /api/sessions/:id`) the server sends a terminal
  `{"type":"session_exited"}` event and closes the stream; a client that
  receives it should stop. A stream that closes **without** `session_exited` is
  a dropped connection — including the run process itself exiting — so reconnect
  with the last id; if the run is gone the reconnection simply fails.

### Reconnecting without gaps

`GET /api/sessions/:id/snapshot` returns the session's full state in one
response — stored fields (messages, tokens, permissions), live runtime state
(`streaming`, current `agent`), and `last_event_seq`: the sequence number of
the most recent event. Pair it with the event stream for an exact, gapless
resync:

```bash
# 1. Read the full state and the stream position it corresponds to.
$ SEQ=$(curl -s http://127.0.0.1:8080/api/sessions/$SID/snapshot | jq .last_event_seq)

# 2. Tail everything that happens after that point (replaying anything that
#    occurred between the two calls).
$ curl -N "http://127.0.0.1:8080/api/sessions/$SID/events?since=$SEQ"
```

This snapshot-then-tail pattern lets a client (or a client that just
restarted) rebuild a session's state and keep it correct without polling.

### Waiting for readiness

`GET /api/sessions/:id/status` reports a session's runtime state. Add
`?wait=<duration>` (e.g. `?wait=10s`) to block until that specific session's
runtime is attached and ready to accept follow-ups, then return its status, or
`503` on timeout. This is session-scoped, unlike `GET /api/ready`, which fires
as soon as any session is ready.

## Session Forking

`POST /api/sessions/:id/fork` creates a new session whose history is a copy of the parent up to (but **excluding**) a specified user message. This lets a client "branch" a conversation — e.g. rewind to an earlier question and try a different prompt — without losing the shared history that came before.

**Request body:**

```json
{ "user_message_index": 1 }
```

`user_message_index` is a **0-based ordinal** that counts only user-role messages in the parent's flat, user-visible message list. The targeted user message is **excluded** from the fork so clients can prefill it into their chat input for the user to edit and resubmit.

**Example:**

```bash
# Fork a session before the second user message (ordinal 1)
$ curl -X POST http://localhost:8080/api/sessions/$SID/fork \
  -H 'Content-Type: application/json' \
  -d '{"user_message_index": 1}'
# Returns: api.SessionResponse for the new forked session
# New session title: "<parent title> (fork 1)", "(fork 2)", etc.
```

**Validation:**

- Out-of-range ordinals (negative, or at/past the user-message count) return `400 Bad Request`.
- An ordinal that resolves to a user message inside a sub-session returns `400 Bad Request`. A sub-session is a nested session created when a multi-agent config delegates work to a child agent; its messages are embedded within the parent session's message list and cannot be used as a fork boundary.

## Idempotent follow-ups

`POST /api/sessions/:id/followup` accepts an optional `Idempotency-Key`
header, making the request safe to retry after a network timeout. A repeat
with a key already seen for the session is acknowledged without delivering the
follow-up again:

```bash
$ curl -X POST http://127.0.0.1:8080/api/sessions/$SID/followup \
    -H 'Content-Type: application/json' \
    -H 'Idempotency-Key: 7f3a-...' \
    -d '{"messages":[{"content":"Ship it"}]}'
# => {"status":"queued_streaming","duplicate":false}
# A retry with the same key => {"status":"duplicate","duplicate":true}
```

The response `status` is `queued_streaming` (a turn is running or starting),
`queued_idle` (delivered to an idle headless session, runs on the next turn),
or `duplicate`.

> [!NOTE]
> **See also**
>
> For interactive use, see the [Terminal UI](../tui/index.md). For agent-to-agent communication, see [A2A Protocol](../a2a/index.md) and [ACP](../acp/index.md). For MCP integration, see [MCP Mode](../mcp-mode/index.md). For an OpenAI-compatible chat-completions API, see the [Chat Server](../chat-server/index.md).
