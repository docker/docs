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
$ docker agent serve api agent.yaml --listen 0.0.0.0:8080

# With session persistence
$ docker agent serve api agent.yaml --session-db ./sessions.db

# Auto-refresh from OCI registry every 10 minutes
$ docker agent serve api agentcatalog/coder --pull-interval 10
```

## Endpoints

All endpoints are under the `/api` prefix.

### Agents

| Method | Path              | Description                       |
| ------ | ----------------- | --------------------------------- |
| `GET`  | `/api/agents`     | List all available agents         |
| `GET`  | `/api/agents/:id` | Get an agent's full configuration |

### Sessions

| Method   | Path                                | Description                                             |
| -------- | ----------------------------------- | ------------------------------------------------------- |
| `GET`    | `/api/sessions`                     | List all sessions                                       |
| `POST`   | `/api/sessions`                     | Create a new session. Accepts an optional `title` field — when set, it is stored and LLM title generation is skipped. |
| `GET`    | `/api/sessions/:id`                 | Get a session by ID (messages, tokens, permissions)     |
| `GET`    | `/api/sessions/:id/status`          | Lightweight runtime state (streaming, title, agent, tokens). Requires an attached runtime. |
| `GET`    | `/api/sessions/:id/snapshot`        | Full state in one call (stored fields + runtime state + `last_event_seq`) for gapless resync — see [Reconnecting without gaps](#reconnecting-without-gaps). |
| `GET`    | `/api/sessions/:id/events`          | Live session event stream (SSE) with sequence numbers and replay. Available for a run attached via [`--listen`](#listen). |
| `DELETE` | `/api/sessions/:id`                 | Delete a session                                        |
| `PATCH`  | `/api/sessions/:id/title`           | Update session title                                    |
| `PATCH`  | `/api/sessions/:id/permissions`     | Update session permissions                              |
| `POST`   | `/api/sessions/:id/fork`            | Fork a session at a user message — creates a new session with messages `[0, message_index)` of the parent (see [Session Forking](#session-forking)) |
| `POST`   | `/api/sessions/:id/resume`          | Resume a paused session (after tool confirmation)       |
| `POST`   | `/api/sessions/:id/tools/toggle`    | Toggle auto-approve (YOLO) mode                         |
| `POST`   | `/api/sessions/:id/elicitation`     | Respond to an MCP tool elicitation request              |
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
data: {"type":"stream_started","session_id":"...","agent":"root"}
data: {"type":"agent_choice","content":"Hello! How","agent":"root"}
data: {"type":"agent_choice","content":" can I help","agent":"root"}
data: {"type":"agent_choice","content":" you today?","agent":"root"}
data: {"type":"stream_stopped","session_id":"...","agent":"root"}
```

Event types include:

- `stream_started` / `stream_stopped` — Agent execution lifecycle
- `agent_choice` — Streamed text content (partial responses)
- `tool_call` — Agent requesting tool execution
- `tool_call_confirmation` — Tool call waiting for user approval
- `tool_call_response` — Tool execution result
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
[{"name":"my-agent","multi":false,"description":"A helpful assistant"}]

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

## CLI Flags

```bash
docker agent serve api <agent-file>|<agents-dir> [flags]
```

| Flag               | Default          | Description                                      |
| ------------------ | ---------------- | ------------------------------------------------ |
| `-l, --listen`     | `127.0.0.1:8080` | Address to listen on                             |
| `--auth-token`     | (none)           | Bearer token required for all API requests. Leave empty to disable authentication (safe when listening on loopback interfaces only). Recommended when `--listen` binds to a network-reachable interface. |
| `-s, --session-db` | `session.db`     | Path to the SQLite session database              |
| `--pull-interval`  | `0` (disabled)   | Auto-pull OCI reference every N minutes          |
| `--fake`           | (none)           | Replay AI responses from cassette file (testing) |
| `--record`         | (none)           | Record AI API interactions to cassette file. Routes through `--models-gateway` when one is configured. |
| `--mcp-oauth-redirect-uri` | (none)   | Public HTTPS URL advertised as the OAuth `redirect_uri` for unmanaged MCP OAuth flows. When set, docker-agent drives PKCE and code exchange in-process and sends the full authorize URL to the client via elicitation. See [Remote MCP](../remote-mcp/index.md) for details. |

> [!TIP]
> **Live profiling (advanced)**
>
> For production diagnostics, set the `CAGENT_PPROF_ADDR` environment variable (or the hidden `--pprof-addr` flag) to a TCP address such as `127.0.0.1:6060`. docker-agent will start a Go pprof HTTP server at `/debug/pprof/`, which you can query with `go tool pprof`. Use a loopback address — a non-loopback binding logs a security warning. This flag is intentionally hidden from `--help`.

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
> Each run started with `--listen` writes a discovery record to `<data-dir>/runs/<pid>.json` containing its address and session id, so a supervising process can find a live run by session id, pid, or address.

## Session event stream and reconnection

`GET /api/sessions/:id/events` is a **Server-Sent Events** stream of the
session's runtime events — `stream_started`, `agent_choice`, `tool_call`,
`session_title`, `token_usage`, `stream_stopped`, and so on. Unlike the
per-request stream returned by the agent-execution endpoint, it is
session-scoped and survives across turns, so a client can watch a session for
its whole lifetime. It is available for a run attached via
[`--listen`](#listen).

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
