---
title: "ACP (Agent Client Protocol)"
description: "Expose Docker Agent agents via the Agent Client Protocol for integration with ACP-compatible hosts like VS Code, IDEs, and other developer tools."
keywords: docker agent, ai agents, features, acp (agent client protocol)
linkTitle: "ACP"
weight: 70
canonical: https://docs.docker.com/ai/docker-agent/features/acp/
aliases:
  - /ai/docker-agent/integrations/acp/
---

_Expose Docker Agent agents via the Agent Client Protocol for integration with ACP-compatible hosts like VS Code, IDEs, and other developer tools._

## Overview

The `docker agent serve acp` command starts an ACP server that communicates over **stdio** (standard input/output). This makes it ideal for integration with editors, IDEs, and other tools that spawn agent processes — the host sends JSON-RPC messages to Docker Agent's stdin and reads responses from stdout.

ACP is built on the [ACP Go SDK](https://github.com/coder/acp-go-sdk) and provides a standardized way for client applications to interact with AI agents.

> [!NOTE]
> **ACP vs A2A vs MCP**
>
> **ACP** connects an agent to a *host application* (IDE, CLI tool) via stdio. **A2A** connects *agents to other agents* over HTTP. **MCP** exposes agents as *tools* for other MCP clients. Choose based on your integration target.

## Usage

```bash
# Start ACP server on stdio
$ docker agent serve acp ./agent.yaml

# With a multi-agent team config
$ docker agent serve acp ./team.yaml

# From an OCI registry
$ docker agent serve acp myorg/agent:tag

# With a custom session database
$ docker agent serve acp ./agent.yaml --session-db ./my-sessions.db
```

## How It Works

1. The host application spawns `docker agent serve acp agent.yaml` as a child process
2. Communication happens over **stdin/stdout** using the ACP protocol
3. The host sends user messages, Docker Agent processes them through the agent
4. Agent responses, tool calls, and events stream back to the host
5. Sessions are persisted in a SQLite database for continuity

```bash
# Conceptual flow:
Host Application
  └── spawns: docker agent serve acp agent.yaml
        ├── stdin  ← JSON-RPC requests from host
        └── stdout → JSON-RPC responses to host
```

## Features

- **Stdio transport** — No network ports needed; ideal for subprocess integration
- **Session persistence** — SQLite-backed sessions survive process restarts
- **Agent runtime support** — Supports configured tools, multi-agent delegation, and model fallbacks. Client-supplied stdio, Streamable HTTP, and SSE MCP servers are supported; audio prompt blocks are accepted (model/provider support varies); `session/load` replays persisted history, while `session/resume` reconnects without replay.
- **Multi-agent configs** — Team configurations with sub-agents work transparently
- **Filesystem operations** — Each session has its own toolsets; shell, filesystem, and Git tools resolve relative paths from that session's working directory
- **Tool permissions** — “Always allow this tool for this session” remembers approval for that tool only; it does not enable autonomous mode for other tools.

## Capability Matrix

The stable ACP v1 surface is covered by wire and behavior tests. `initialize` advertises the following agent capabilities; the client capability matrix exercises all 32 combinations of filesystem read/write, terminal support, and form/URL elicitation.

| Surface | Advertised/implemented behavior | Regression coverage |
| --- | --- | --- |
| Sessions | New, load with replay, resume, list, close, delete, additional directories | `TestACPCapabilityMatrix`, session lifecycle/load/delete suites |
| Authentication | `host-credentials` adoption and connection-local logout | Authentication wire/lifecycle suites |
| Prompt content | Text, resource links, embedded resources, images, audio | Content/audio suites; model/API limits still apply |
| Client MCP | Stdio, Streamable HTTP, legacy SSE | Client MCP and session-remote suites |
| Filesystem | Read tools require read capability; write requires write; edit requires both | `TestACPCapabilityMatrix`, filesystem capability suites |
| Shell terminals | Client-backed shell only when negotiated; otherwise native shell | `TestACPCapabilityMatrix`, terminal suite; no fallback after client failure |
| Elicitation | Restricted form/URL modes only when explicitly negotiated | `TestACPCapabilityMatrix`, elicitation suite; absent/null/empty capabilities do not enable a mode |
| Configuration | Advertised select options and legacy modes | Session configuration suite |
| Message/tool metadata | Persisted message IDs, bounded replay, programmatic tool names | Message-ID/tool-name suites |
| Tracing | Request-scoped W3C trace context; no tracing handshake | ACP tracing suite and MCP callback ownership tests |

Session forking, editor provider configuration, document synchronization, MCP-over-ACP, and next-edit suggestions are not advertised or implemented. Unsupported requests return method-not-found; unsupported notifications have no response. Terminal authentication and experimental environment-variable authentication remain unadvertised. This matrix is not a claim of ACP draft/v2 support or universal model modality support.

## Trace Propagation

Enable export using `docker agent serve acp agent.yaml --otel` and the [OpenTelemetry configuration](../../configuration/opentelemetry/index.md). ACP metadata propagation itself needs no capability negotiation or exporter.

Send W3C context in **request/notification `params._meta`**, not at the JSON-RPC envelope root:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "session/prompt",
  "params": {
    "sessionId": "your-session-id",
    "prompt": [{"type": "text", "text": "Explain this project"}],
    "_meta": {
      "traceparent": "00-11111111111111111111111111111111-2222222222222222-01",
      "tracestate": "vendor=value"
    }
  }
}
```

Each implemented agent request/notification starts an `acp.<method>` SERVER span, linked to valid incoming context and enclosing its handler. These boundary spans record only the static RPC method, not request content, paths, credential metadata, or raw errors. Invalid/non-string/oversized trace fields are ignored; each accepted field is bounded to 512 bytes and validated by the W3C propagator. A bad tracestate does not invalidate a valid traceparent. Client baggage and arbitrary `_meta` fields are not imported into tracing.

The request context flows into existing runtime/model/tool instrumentation. Outgoing ACP session updates, permission requests, filesystem operations, terminal operations (including detached cleanup), and elicitation carry only `traceparent`/`tracestate` in `params._meta`. Metadata is copied before injection; stale trace fields and baggage are removed. Response metadata is never adopted as the current trace. MCP elicitation and sampling callbacks use the sole owning tool call's trace rather than the connection's setup trace; existing ambiguous-owner restrictions remain in effect.

Traces are not stored on sessions: concurrent and later prompts use their own incoming context. Terminal final-close retries use the cleanup caller's context, not a remembered creation trace. Detached background/drain work may outlive the handler span while retaining its trace. This adds request SERVER spans and propagation, not separate ACP CLIENT spans for every outbound call, response metadata, or automatic metadata on unsupported SDK methods. Existing runtime/provider telemetry and its content-capture settings retain their own policy; the boundary-span restriction is not a global telemetry-redaction guarantee.

## Session Configuration

New/load/resume responses include stable select-based `configOptions` and legacy `modes`. `session/set_config_option` accepts exact advertised string IDs and returns the complete updated option list. Boolean variants, unknown IDs/values, arbitrary unlisted model references, and safety aliases are rejected. `session/set_mode` is a compatibility path to the same session safety setting.

| Option ID | Scope and values |
| --- | --- |
| `mode` | Docker Agent dispatcher safety: `default`, `strict`, `balanced`, `restricted`, `autonomous`. |
| `model` | Selected agent's configured named models (`model:<name>`) plus `default` to restore its configured providers. |
| `thought_level` | Supported explicit reasoning levels for a recognized single, non-routing model; `default` restores the selected model's configured budget. |

`default` safety preserves legacy behavior (read-only-annotated tools auto-approve; others ask), not balanced classification. Choosing it clears blanket autonomous approval. Deny rules, session-scoped ask rules, mandatory tool guards, and `preempt_yolo` hooks retain their existing precedence. Team-level ask rules and ordinary approval hooks do not override autonomous approval. Switching mode never clears remembered per-tool grants. These controls govern Docker Agent's dispatcher, not approvals inside external coding harnesses. Harness agents have no ACP model/reasoning controls.

Model options are built from local configuration only, without model-catalog/network discovery, provider construction, or tool startup. Nested alloys are not advertised; flat configured alloys can be selected but have no reasoning selector. Unknown/unsupported reasoning models, including models recognized only through a configured thinking budget, likewise omit that selector. A runtime model not representable by an advertised configured choice is shown as an inert `current` value; selecting it is a no-op rather than a new provider request.

Model and safety choices are persisted; failed writes preserve the previous live state and restore the exact previous providers. Reasoning changes are runtime-only: active load/resume keeps them, while cold reconstruction or model reset restores configured budgets. Explicit `none` remains distinct from an unset/adaptive/token-based configured budget. Cold reconstruction applies persisted model overrides before publishing the session; a stale or unavailable override fails setup instead of claiming a selection that was not applied.

Configuration changes require an idle session: running, queued, or draining foreground work and admitted background tasks reject mutations without canceling them. Prompts cannot enter during provider construction/persistence, and close/delete/shutdown cancel and join the operation. Load/resume during background work omits volatile model/reasoning selectors instead of displaying temporary overrides. Configuration notifications refresh at idle turn boundaries and agent-switch commands, without exposing temporary delegated selections. Once persistence succeeds, a notification failure does not undo the committed change; reconnect to recover current state.

## Slash Commands

ACP advertises supported commands when a session is created or resumed, at turn start, and when agent information changes. Discovery reads command metadata only; it does not start tools or expand command instructions.

| Command | Behavior |
| --- | --- |
| `/compact [instructions]` | Runs manual compaction of the root session using the existing runtime hooks and persistence. Reports applied, skipped, or failed rather than starting an ordinary model turn. |
| `/usage` | Reports current context-token counts and the session cost snapshot without adding history or calling the model. The context limit is the last value reported for the selected agent, not a fresh provider lookup; it is unknown until a usage event supplies it and may lag a model change. Recorded live child-session costs are included. |
| Configured literal-prompt commands | Replace the leading slash command with its instruction and append trailing arguments literally. Subsequent attachments are preserved. |
| Configured agent-switch commands | Switch the active agent after looking up the command in the original agent's table. A switch without text or attachments starts no model turn. The selection lasts for the active runtime; cold resume still starts its default agent. |

`compact`, `usage`, and `new` are reserved names. `/new` is not advertised or executed: clients must use `session/new` to create a fresh conversation without destroying existing history. `/usage` accepts no arguments, and both built-ins reject attachments rather than silently dropping them.

Only directly supplied leading text is parsed as a command. Text inside attached resources is never interpreted as a command. Unknown slash-prefixed messages pass through as ordinary chat. Known commands requiring URL opening, JavaScript `${...}`, or bang-tool expansion are not advertised and return an explicit error if invoked; those expansion paths are not yet integrated with ACP's permission flow.

Commands share the session's normal turn admission, cancellation, and cleanup. Manual compaction errors are fatal for that command (`compaction_failed`), even though the same diagnostic during automatic compaction can be recoverable for an ordinary prompt. A skipped compaction does not claim to have changed history. Built-in invocations and their status messages are not added to conversation history.

## Prompt Outcomes and Errors

`session/prompt` reports why the root session stopped:

| Stop reason | Meaning |
| --- | --- |
| `end_turn` | Normal completion |
| `max_tokens` | The final model response reached its output-token limit |
| `refusal` | The final model response was a refusal |
| `max_turn_requests` | Execution stopped at the iteration limit, rather than being approved to continue |
| `cancelled` | The prompt was canceled or its context expired |

Cancellation takes precedence, and responses wait until runtime events have fully drained. Child-session outcomes, recoverable compaction diagnostics, warnings, and model fallbacks do not by themselves fail the root prompt. A later successful response supersedes an earlier model stop during a continued turn.

Fatal root-runtime failures return JSON-RPC internal error `-32603` with `data.sessionId`, `data.runtimeCode`, and `data.error`. Budget termination is reported this way as `budget_exceeded`, not confused with an output-token or iteration limit. Diagnostic updates may already have streamed before the error response. Missing prompt/resume/load sessions return `-32002` (resource not found); invalid workspace parameters return `-32602` (invalid params).

## Authentication and Logout

ACP advertises the agent-managed `host-credentials` method (**Use configured host credentials**) and `agentCapabilities.auth.logout`. No keys are accepted in `authenticate`, `_meta`, or elicitation forms. This is explicit adoption of the agent host's configured credentials, not an interactive login or a check that a provider will accept a token.

Working ambient credentials and local models continue to work without an extra authentication call. If initialization encounters a recognized missing model-credential or Docker gateway login error, it still returns capabilities and authentication methods. Session creation/load/resume then returns ACP `Authentication required` (`-32000`), rather than preventing the editor from initializing. Invalid configuration and unrelated initialization failures remain errors.

Configure credentials **outside ACP, on the agent host**: for example, run `docker agent setup`, update the configured env file, or sign in to Docker Desktop for the Docker gateway. Then retry:

```json
{"method": "authenticate", "params": {"methodId": "host-credentials"}}
```

Authentication reloads credential sources and validates the configured team before reopening access. It can see newly created or changed credential files, while preserving source precedence and explicit provider overrides. It never changes the process environment, writes credentials, or substitutes another model. Changing the launching shell's environment cannot change an already-running agent; restart it to pick up those environment changes. Repeating authentication on a healthy connection is a no-op. Unknown method IDs are rejected. If a later session load detects missing credentials, close existing sessions before authenticating to refresh them; live sessions never have their clients silently replaced.

```json
{"method": "logout", "params": {}}
```

Logout is **connection-local**: it immediately blocks new work, cancels and drains session construction, prompts, background work, terminals, and MCP connections, then closes the validation team. Repeated logout joins the same cleanup. Canceling a logout request does not reopen access or abandon cleanup; shutdown also waits for it. Cleanup failure remains an error and prevents reauthentication.

Logout does **not** delete shared API keys, ChatGPT/Docker logins, revoke provider tokens, or erase conversation history. Other agent processes/connections remain unaffected, and a newly started agent can still use those ambient credentials. After successful logout, this connection requires `authenticate` before new sessions, load/resume, prompts, configuration changes, listing, or deletion. Close/cancel remain available for cleanup; a deletion already admitted before logout may have committed. Resume/load saved sessions after reauthentication to continue them.

Terminal authentication, browser/device-code login driven by the editor, and the deprecated experimental `env_var` authentication method are not advertised. Remote provider 401 responses are not treated as proof that the host credential snapshot should be replaced.

## Client-Supplied MCP Servers

Pass stdio, Streamable HTTP (`type: "http"`), or legacy SSE (`type: "sse"`) MCP servers in `mcpServers` on `session/new`, `session/resume`, or `session/load`:

```json
{
  "cwd": "/home/user/project",
  "mcpServers": [
    {
      "name": "project-tools",
      "command": "/usr/local/bin/project-mcp",
      "args": ["--stdio", "--project", "/home/user/project"],
      "env": [{"name": "PROJECT_MODE", "value": "development"}]
    }
  ]
}
```

The executable path must be absolute. These are **local subprocesses**, running with the agent's OS permissions and the session's working directory. Arguments and environment values are passed literally, without shell expansion or installation. Processes inherit the agent's environment (including credentials); supplied entries override inherited values, and the last duplicate wins. MCP-over-ACP configurations remain unsupported.

Remote servers can be mixed with stdio servers; names must be nonempty and unique across the list:

```json
{
  "cwd": "/home/user/project",
  "mcpServers": [
    {
      "type": "http",
      "name": "project-api",
      "url": "https://tools.example.com/mcp",
      "headers": [{"name": "Authorization", "value": "Bearer CLIENT_TOKEN"}]
    },
    {
      "type": "sse",
      "name": "legacy-tools",
      "url": "https://legacy.example.com/sse",
      "headers": []
    }
  ]
}
```

Remote connections originate from the **agent host**, not the editor. URLs must use HTTP(S), with no userinfo or fragment. They use the existing guarded outbound transport: direct connections to loopback, private, and link-local addresses are blocked; operator-configured proxies retain their existing egress policy. ACP cannot enable `allow_private_ips`. All requests, including redirects and SSE-discovered message endpoints, must remain on the configured origin (scheme, host, and effective port). HTTPS downgrades are rejected.

Headers are literal: no environment, upstream-header, JavaScript, or secret-provider expansion. Header names are case-insensitively unique and must be valid HTTP fields. Transport-controlled headers (including `Mcp-*`, `Host`, content negotiation/framing, hop-by-hop/proxy headers, `Last-Event-ID`, and idempotency keys) are rejected. Clients supply authentication headers themselves; these connections never consult the agent's OAuth/keyring tokens, discover OAuth metadata, open a browser, or prompt for authorization. Cookie jars are private to each server connection and are discarded on replacement. Prefer HTTPS for credentials.

Transport failures use generic diagnostics to avoid exposing URL/query credentials. Remote client connections retain host-only MCP tracing rather than full-URL HTTP spans.

The agent initializes the servers and lists their tools before accepting setup, within a 30-second setup budget. Client tools supplement configured tools and are exposed to the session's agents, honoring each agent's read-only filter. Model-facing names are bounded, generation-specific aliases; remembered per-tool approvals do not transfer to replacements. Existing global permission policies still apply.

Every successful resume supplies the complete replacement server list. An omitted or empty list removes all client servers. An idle session keeps its runtime and conversation; setup failure leaves its previous servers and additional roots unchanged. Cached calls to retired tools fail rather than being redirected, and calls already using the old servers are canceled and joined during retirement. Cleanup failure blocks further prompts/resumes and is retained by close/shutdown; failed cleanup remains an error for that agent process.

Servers are not automatically restarted or retried after disconnection; explicitly resume to create fresh connections. Closing a session retires its client servers, stops subprocesses, and closes/drains remote connections, including pending calls. Streamable HTTP teardown attempts session DELETE when a session ID exists, with forced local teardown on the cleanup deadline. This cannot guarantee that a remote server stops already-started work. Connection settings, headers, and environment values are not persisted as session configuration: clients must send them again on cold resume/load. This does not implement an OS sandbox or full ACP elicitation support.

## Elicitation

ACP clients can negotiate form and URL elicitation independently through explicit non-null `clientCapabilities.elicitation.form` and `.url` objects. An absent capability, `{}`, or a null mode does not advertise support. Unsupported modes continue to be declined without a client request; tool-permission and iteration-limit prompts remain separate and do not become auto-approved.

Supported requests use `elicitation/create` with the owning ACP session ID, including delegated/background work. Accepted form data is checked against the supported flat-object schema before returning to the tool. The adapter accepts primitive fields, string enums/titled enums, and enum arrays; nested structures, references, unknown keywords, and unsupported formats are declined rather than weakened. Invalid/missing response actions, invalid form values, and transport failures fail closed. Request/reply payloads are limited to 256 KiB, including normalized JSON encoding; oversized replies decline only the elicitation, leaving the connection usable. The connection permits at most 128 pending elicitation requests.

Form mode must not collect access-granting credentials. Internal OAuth flows, OAuth client-credential forms, sudo password forms, and explicit credential requests in names/messages/schema labels are declined. This screening is heuristic, not a guarantee against arbitrary or obfuscated free-text requests: only use trusted tools for elicitation, and use URL mode for sensitive interactions. An accepted URL request means consent to open the URL, not completion of authorization. Only HTTP(S) URLs without userinfo are forwarded; the agent does not open/fetch them, does not return URL-response content to tools, and does not currently send optional `elicitation/complete` notifications. Internal OAuth adaptation remains unsupported.

The pinned SDK lacks top-level elicitation session IDs and decodes generic numbers as float64. The ACP connection wrapper supplies the session field and rejects numerical values that would change through that SDK conversion; accepted form replies are inspected before SDK decoding. Integers with magnitude at least 2^53 and oversized numeric representations are declined. This protects the ACP boundary, not precision already lost by an upstream MCP server/client. Go embedders must use `Agent.NewConnection` to enable this compatibility path; externally bound SDK connections retain the headless-decline fallback.

Elicitation waits follow the originating operation's cancellation, not the foreground parent's lifetime for background work. MCP requests without one unambiguous owning call are not routed to an arbitrary session. Close, deletion, shutdown, and prompt cancellation retain their existing join guarantees; a blocked SDK transport writer still has no hard cancellation deadline. Client-supplied server setup outside a runtime request remains non-interactive. Elicitation request/reply bodies are not added to conversation history by the bridge, but tools can incorporate accepted form data in their results.

## Session Workspaces

New sessions and resumes that reconstruct a runtime load the agent configuration again and create independent teams and toolsets. Configuration changes affect subsequent loads, not teams already serving sessions. Config-relative paths, such as instruction files, remain relative to the agent configuration file.

The session's `cwd` is the execution directory, independent of where the ACP subprocess was launched. ACP `session/new` and `session/resume` wire requests require `cwd`; the SDK rejects omitted values. Direct Go callers retain the legacy empty-`cwd` fallback: new sessions use the configured or process directory without saving invented workspace provenance, and resumes use the saved directory when available. Workspace selection does not by itself sandbox tools or make explicitly shared storage private.

A resume must name the same directory as the saved workspace. Filesystem aliases of the same directory are accepted, but the saved path is not rewritten. A legacy session with no saved workspace cannot adopt an explicit `cwd`; ACP clients must create a new session instead.

Every successful resume replaces the complete `additionalDirectories` list. Omitting it or sending an empty array revokes all additional roots; previous roots are never implicitly restored. Invalid paths or workspace mismatches leave session state unchanged.

Resuming a registered session while a foreground prompt is running, queued, or draining, or another resume holds the setup reservation, returns a busy error without canceling that work or changing roots. Retry after it finishes. This guards foreground turns, not detached background work or already-issued client I/O; it is not an atomic revocation guarantee.

## Loading Session History

`session/load` reconnects to a saved session and sends its persisted conversation as `session/update` notifications before returning the load response. Unlike `session/resume`, it replays history. Both paths retain workspace identity validation and complete replacement of client MCP servers and additional roots. The load wire request requires `sessionId`, `cwd`, and a non-null `mcpServers` array (use `[]` for none). Cold load starts the configured default agent; loading an existing idle session keeps its runtime and agent selection.

Replay reads a captured copy of persisted history, including partial replies and errors saved before a failed turn. It does not call the model, execute historical tools, request permission, rerun hooks, or fetch attachment content. Requested MCP servers are set up as part of reconnection, as with resume. Foreground prompts and competing reconnects are rejected while a load holds the reservation; a load also rejects a running, queued, or draining foreground prompt without canceling it. Detached background work is not frozen by this reservation.

Visible user/assistant messages, stored reasoning, tool results, and errors retain stored item order. Nested sessions replay at their stored positions; original cross-session streaming interleaving is not recoverable. System/implicit messages, internal compaction summaries, evaluations, and provider-private state are excluded. Historical tools use fresh opaque IDs and terminal status from stored results. Missing results use `failed` with an explicit unknown-outcome notice, not an assertion that execution failed. Historical arguments and derived locations are omitted because stored arguments can predate input transforms/redaction. Tool-result attachments, transient edit diffs, and plan snapshots are not reconstructed from raw output.

Text is chunked into UTF-8-safe updates; tool output is limited to a 64 KiB prefix with a truncation notice. Stored inline binary attachments up to 512 KiB can be replayed; user audio within that bound is replayed as audio blocks with MIME parameters preserved. Larger blobs, external file/artifact references, remote image URLs, and assistant audio use unavailable-content notices without I/O. Replay notifications are capped at 1 MiB of JSON-encoded session payload, leaving room below the transport frame limit. This is a persisted transcript view, not lossless recovery of every live notification or original resource URI.

Replay errors return an error rather than a successful load response and close/unroute the session; already-delivered notifications cannot be rolled back. Retry with an explicit load after cleanup succeeds, resetting any partially rendered client history first. Request-context cancellation, close, and shutdown stop replay and join its operation; `session/cancel` remains a prompt cancellation method. The SDK cannot interrupt a blocked writer, so cleanup has no hard write deadline.

## Listing Sessions

`session/list` returns up to **50 sessions per page**, ordered by creation time (newest first), then by session ID ascending to break ties. To fetch the next page, pass the returned opaque `nextCursor` as `cursor` and keep the same `cwd` filter. An absent `nextCursor` marks the end; no matches return `sessions: []`. Invalid cursors or a changed filter return invalid params (`-32602`).

The optional `cwd` filter must be an absolute path. Matching uses cleaned, exact path strings, not prefixes, case folding, or symlink resolution. The directory need not still exist, so historical workspaces remain discoverable. Sessions with unknown or non-absolute workspace metadata are omitted rather than attributed to the server's current directory.

Listings use stored summaries, not conversation histories. The optional `updatedAt` is omitted because the store does not track last activity; creation time is used only for ordering, not reported as activity. Additional roots are reported from active session state and omitted for inactive sessions.

Pagination does not freeze a snapshot across requests: changes to history or metadata between pages can affect results. Page size bounds the response, while the store still retrieves all session summaries for filtering and ordering.

## Closing Sessions

A successful `session/close` response means the foreground turn has drained, the runtime has joined its background agents, and session toolset shutdown has completed through the existing toolset lifecycle. Client MCP shutdown starts before those joins so blocked subprocess I/O can be terminated rather than strand the drain. Close also cancels and joins in-flight resumes for that session so they cannot publish a replacement runtime after closure.

Concurrent close requests join the same cleanup. Canceling a close request only stops that caller's wait: cleanup continues, and resume remains blocked until it succeeds. After a successful close, an explicit resume may reconstruct the session. Cleanup failures are returned and retained; the same session cannot reopen in that agent process when cleanup is uncertain.

Server shutdown rejects new work, cancels admitted initialization/session/list operations, and drains them before closing the session store. Shutdown is a final join rather than a bounded timeout: an uncooperative runtime or tool can delay it. Toolset stop errors are surfaced, not treated as successful cleanup. These guarantees do not add disposal support to toolsets whose resources fall outside the existing lifecycle contract, nor undo already-issued client I/O.

## Deleting Sessions

`session/delete` permanently removes a root conversation's stored history, descendant sessions, and stored generated-media records/blobs. It does not remove files from the workspace. The agent advertises `sessionCapabilities.delete`; the pinned Go SDK still names the handler `UnstableDeleteSession`, but the wire method is `session/delete`. An empty session ID is invalid; already-deleted and unknown IDs succeed without creating a session. Subsequent load/resume requests for deleted history return not found.

Deletion closes the root session first, canceling its foreground work and joining loads/resumes, runtime/background tasks, and client MCP cleanup before deleting storage. Cleanup failure preserves history and remains an error for that lifecycle. An ordinary store error leaves the session closed and can be retried; SQLite subtree/media changes occur in one transaction.

Direct deletion of child sessions is rejected. Independently loaded descendants must be closed, with successful cleanup, before their root can be deleted. While a delete is admitted, new session construction and reconnects are temporarily rejected agent-wide to prevent a descendant from being loaded during the cascade. A delete also returns busy if another session already has an in-flight constructor whose ancestry is not yet known. Different-ID deletes are serialized by returning busy; concurrent requests for the same ID join one operation. Unrelated already-running sessions are not canceled.

Canceling a joining request stops only that caller's wait. Canceling the initiating request before storage deletion lets cleanup finish but skips deletion, keeping the admission barrier until the drain ends. Cancellation during a database call does not prove whether it committed; retrying deletion is safe. Shutdown joins any admitted deletion before the caller closes the store. These guarantees cover this agent instance, not other processes concurrently writing the same database; they do not provide an atomic filesystem revocation or a deadline for uncooperative cleanup.

## Plan Snapshots

Todo tool results with typed todo metadata produce complete ACP plan snapshots. An empty todo snapshot sends `entries: []`, replacing and clearing the previous plan; completed entries remain visible until the todo storage is cleared. Missing or unrelated metadata does not change the plan, and textual tool output is not parsed to infer one. This uses the stable v1 `plan` update, not ID-based plan-removal extensions.

## Context and Cost Accounting

The context gauge describes the root ACP conversation, not a delegated or background agent's separate context window. Child usage updates can change the displayed cost without replacing the root's token count or limit. Context changes reported for an in-place root handoff or compaction replace the root snapshot.

Cost is the sum of the latest cumulative runtime-session snapshots, keyed by session ID rather than agent name. Repeated snapshots replace previous values; parallel children using the same agent remain distinct. Restored descendant costs are already included in the root's initial contribution and are not added a second time. Compaction and evaluator costs follow the runtime's existing accounting; unpriced work is not estimated by ACP.

Background usage is recorded without unsolicited client writes. It appears on the next foreground usage update or `/usage`; the latter refreshes root counters and uses the same aggregate for its text and structured response. Accounting continues while canceled or failed streams drain, including nested background children, without changing prompt outcomes or adding teardown notifications.

No `usage_update` is emitted until a positive root context-window limit is known. `/usage` can still report text with an unknown limit and recorded cost. The last limit may lag a model change; switching the selected agent without a matching root usage observation makes it unknown. Agent names and child context windows are not inferred from each other's usage events.

## Message and Tool Identifiers

Message chunks carry agent-generated `messageId` values. Chunks and rich-content blocks from the same message share an ID; assistant thoughts and visible answers use distinct IDs. New model-stream attempts get fresh logical IDs before their first delta, including retries/fallbacks. Response-cache hits also emit an identified answer, once per hit; cache entries containing XML tool payloads are treated as misses rather than exposing text suppressed during their original stream.

Logical IDs are stored with messages and survive SQLite persistence and JSON round-trips. ACP derives opaque display UUIDs from the owning session, logical message, and content/thought channel. Persisted messages therefore keep their display IDs across active and cold `session/load`, while copied messages in a different owning session have different IDs. Nested histories use each child's owning session, not the root notification destination. User input is identified by the agent; the pinned SDK's obsolete draft `session/prompt.messageId` and response `userMessageId` fields are not adopted or echoed.

Older messages without logical IDs receive deterministic replay IDs based on their owning session and original history position. Those IDs remain stable only while the history's positions remain unchanged. Warnings, command results, and other transient display notices receive fresh IDs; replayed stored errors have position-based IDs and are not correlated to their original live notices. IDs do not guarantee durability: store failures and interrupted child output that was never persisted cannot be recovered by replay. Tool-call invocation IDs remain separate, opaque, and fresh during replay.

Tool calls also expose the exact programmatic `name`, separately from their human-readable `title`, on starts, permission requests, updates, and known completions. Client MCP names remain generation-specific aliases; names do not confer permission or transfer approvals. Unknown orphan names and invalid/oversized names are omitted, not guessed or truncated. Historical names do not enable replay of stored tool arguments or derived locations.

The pinned ACP SDK lacks the stable tool `name` field. `Agent.NewConnection` (used by `serve acp`) promotes a private `_meta` carrier only at protocol tool-call boundaries. Embedders binding the SDK connection directly instead receive `docker-agent/internal-acp-tool-name` in the tool call's `_meta`; use `NewConnection` for the stable top-level field. Arbitrary tool payloads and unrelated metadata are not rewritten.

## Prompt Attachments

ACP text remains text. Embedded text resources, binary resources, images, audio blocks, and successfully read file links become ordered document attachments with safe display names, MIME types, actual byte sizes, and inline payloads. Duplicate resources remain separate attachments. Text resources such as `application/json` are sent as text even when the model does not support that MIME as a binary format. Binary document support still depends on the selected model/provider.

Images are validated and normalized using the attachment pipeline; resizing includes a coordinate-mapping note. Per-attachment limits are 5 MiB for text, 20 MiB for decoded binary data, and 16 million pixels for locally decoded images. The transport's message-size limit still applies before these checks. Invalid or oversized attachments produce a bounded unavailable-content notice without exposing the payload or full source URI.

### Audio Prompts

The agent advertises `promptCapabilities.audio: true`. Supply raw base64 bytes and an audio MIME type in a `session/prompt` block:

```json
{"type": "audio", "mimeType": "audio/wav", "data": "<base64-encoded WAV bytes>"}
```

Audio is preserved byte-for-byte as an inline document, in prompt order; duplicate blocks remain separate. MIME type casing is normalized and format parameters such as `audio/pcm; rate=24000` are retained. MIME metadata is limited to 256 bytes. Empty audio, malformed base64, invalid/non-audio/wildcard MIME types, or oversized attachments become a fixed unavailable-content notice. These are MIME/base64/size checks, not audio codec validation. No URL/file lookup, playback, transcription, transcoding, subprocess, or automatic model switch is performed.

Every accepted audio block also adds a short text notice explaining that it is not a transcript and audio access depends on the selected model/provider. The notice remains if capability filtering or provider conversion drops the audio, including for audio-only prompts. Actual support requires both an audio-capable model (catalog metadata or an explicit `capabilities.audio` override) and a suitable API:

- **Gemini:** native inline audio bytes and MIME parameters, subject to the model/API's supported formats and limits.
- **OpenAI-compatible Chat Completions:** native `input_audio` for WAV (`audio/wav`, `audio/x-wav`, `audio/wave`, `audio/vnd.wave`) and MP3 (`audio/mpeg`, `audio/mp3`). Other audio formats are omitted, not relabeled or transcoded.
- **OpenAI Responses/ChatGPT and Anthropic:** audio bytes are not sent; the text notice remains. A capability override does not add an unsupported API format.

The 20 MiB decoded-binary attachment ceiling still applies, but the ACP transport rejects frames above 10 MiB **before** attachment handling. Base64 and JSON overhead therefore make the effective on-wire audio limit less than roughly 7.5 MiB, further reduced by other blocks. Such rejected frames cannot receive an attachment placeholder. Replay has the smaller 512 KiB inline-binary limit described above.

### Resource Links and Attachment Handling

File resource links are read only through a client that advertises `fs.readTextFile`, after session-root validation. File URIs are decoded once, so a literal `%20` in a filename is not changed into a space. Remote authorities and non-file URIs are not fetched, and unavailable links never fall back to host file contents. Resource links retain the separate session-root policy described below, not a filesystem toolset's allow/deny policy.

This preserves payloads rather than providing lossless protocol round-tripping: arbitrary ACP annotations, `_meta`, titles/descriptions, and original source URIs are not persisted as document metadata. Image encoding may change during normalization. The secret-redaction builtin scans document text and labels before model calls without changing stored history; it does not scan inside binary files, images, or audio.

Attachment-bearing prompts bypass lookup and storage in the agent response cache, whose keys contain only text. This does not disable provider-side prompt caching.

Tool-result presentation remains transformed text plus eligible edit diffs. Raw tool media, documents, and structured results are not forwarded to ACP, because the current output-transform contract covers text only.

## Client Terminals

When a client advertises `clientCapabilities.terminal: true`, configured `shell` tools execute through its `terminal/create`, `wait_for_exit`, `output`, `kill`, and `release` methods. Otherwise the existing native shell remains in use. Negotiated client execution never falls back to running locally after a client error. Tool permissions, safety classification, input hooks, and output transforms still surround the same `shell` call.

This adapter requires the client and agent to use the same platform and coherent workspace/interpreter paths: ACP does not negotiate a remote OS or shell. Client commands use `/bin/sh -c` on non-Windows hosts or the absolute system `cmd.exe /D /C` path on Windows, not host `$SHELL`, `ComSpec`, or PATH lookup. `get_environment_info` reports this expected interpreter. Other host-derived environment metadata describes the agent host and does not select the client shell. The client must enforce its own execution boundaries.

Only explicitly configured shell environment overrides are expanded and sent; the agent's full inherited environment is not transmitted. Explicitly referenced secrets are still sent if configured. The client supplies its own inherited environment. Working directories are absolute; relative `cwd` values resolve against the session toolset's workspace. This is not a shell sandbox: absolute directories outside workspace roots remain allowed, just as with native shell execution. `sudo_askpass` is unsupported for client-backed shell and returns an error without executing.

Calls retain the default 30-second timeout and can request a longer positive timeout. Output requests retain at most the last 64 KiB, with a truncation notice; ACP receives only the final transformed tool output, not raw streamed output or terminal references. Caller cancellation omits partial output because canceled output hooks cannot guarantee a rewrite. The client itself owns the process and necessarily sees its raw output. Nonzero exits and signals retain native-style error text; transport/protocol failures are tool errors.

Every received terminal ID is owned by the concrete session until release succeeds. Cancellation/timeout triggers kill and output collection, and release is attempted even if either fails. Session close signals terminal operations before joining runtime/background work, then makes one final release attempt for unresolved IDs. Unresolved cleanup blocks new work and prevents deletion from discarding history. Closed-session handlers cannot move ownership to a replacement session.

A create attempt is joined with a 10-second response budget even after tool cancellation, allowing late IDs to be released. If no ID arrives because of timeout, disconnection, or an ambiguous error, execution is unknown and the session retains a cleanup failure; the client must reclaim any orphaned process. Definitive invalid-request/invalid-params/method-not-found rejections do not poison the session. Five-second cleanup request budgets do not provide a hard deadline for a blocked SDK writer. No automatic create retry is performed.

`script`, background-job tools, Git commands, skill command execution, and filesystem post-edit hooks still run locally; this change does not route every subprocess through the client.

## Client Filesystem Capabilities

The ACP `filesystem` toolset exposes text operations according to the client's negotiated `fs` capabilities:

| Client support | Available client-backed tools |
| --- | --- |
| Neither read nor write | None |
| `readTextFile` only | `read_file`, `read_multiple_files` |
| `writeTextFile` only | `write_file` |
| Both | `read_file`, `read_multiple_files`, `write_file`, `edit_file` |

Configured tool filters and `readonly` restrictions still apply. Unsupported operations are omitted from normal and deferred discovery. These tools never fall back to host file contents when a capability is missing or a client request fails; edits require both client capabilities rather than mixing editor and disk content.

Single- and multi-file reads use client-provided text, including unsaved editor buffers and files not yet saved to disk. `read_multiple_files` makes one request per permitted path, in input order, preserving duplicate paths and the existing text/JSON output format. Ordinary errors are recorded per file without discarding other results. Cancellation stops further requests once observed; it cannot revoke an already-issued read or impose a hard timeout on the transport. These APIs are text-only, not image or binary reads.

Filesystem capabilities describe client methods, not authorization or a read-only sandbox. Search, directory listing, and directory mutation tools still run on the agent host with their existing policies. Disk-backed search does not include unsaved editor buffers, and large client reads remain subject to transport message-size limits.

## Filesystem Policies and Post-Edit Hooks

The ACP `filesystem` toolset's `read_file`, `read_multiple_files`, `write_file`, and `edit_file` operations enforce `allow_list`, `deny_list`, and `.agentsignore` before requesting client I/O. Session workspace roots remain an additional restriction: adding a workspace root does not override a deny rule or expand the configured allow list. Invalid allow/deny configuration disables these operations.

These are path checks, not an atomic sandbox around client I/O. The ACP client must enforce access boundaries when opening files; the agent cannot apply local `os.Root` protections to another process's file operations. User-supplied `resource_link` attachments follow a separate, session-root-checked path and are not governed by a filesystem toolset's policy.

Configured `post_edit` commands run **locally**, in the session toolset's working directory, only after a successful client write. They receive the checked target path in `${file}` and match patterns against that target, not a symlink alias. Hooks require the client and agent to share a coherent on-disk filesystem; editor-buffer-only writes are not mirrored to local disk. A hook failure reports that the write succeeded but the hook failed, without retrying or rolling back the write.

## Tool-Call Lifecycle

ACP reports a tool awaiting permission as `pending`. Actual execution is reported as `in_progress`, followed by `completed` or `failed`. Permission requests and subsequent updates refer to the same opaque ACP tool-call ID; permission approval is not inferred from display status.

Tool-call titles include the tool's display name, followed by a short detail when one is available: the command for shell-style tools (`shell` and `run_*`, `exec_*`, `stop_*` tools, read from `cmd`, `command`, or `script`), the path for file tools (`read_*`, `write_*`, `edit_*`, `delete_*`, `create_*`, and `directory_tree`), or the URL for `fetch` and `http_*` tools. For example, a permission prompt reads `Shell: docker ps -a` rather than just `Shell`. Whitespace in the detail is collapsed and it is truncated to 80 characters; if no detail is found the title is the display name alone. Other arguments, such as file content, never appear in the title.

Rejected, policy-denied, or unavailable tools can produce a terminal result without executing. These become failed tool items instead of aborting the conversation or inventing a running phase. Tool kinds use one display classifier across permission and execution updates: operation names take precedence over generic hints, and a destructive hint alone does not mean deletion.

ACP assigns fresh IDs rather than exposing runtime IDs directly. Correlation distinguishes agent names and supports sequential runtime-ID reuse, including across turns. The runtime does not yet provide enough identity to disambiguate simultaneous identical IDs from two sub-sessions of the same agent.

On interruption, ACP drains runtime events before finishing the prompt. Known results retain their actual status and output; visible calls with no result receive a best-effort failed update explaining that side effects may have occurred. Cancellation still determines the prompt's `cancelled` stop reason. Notification failure cannot undo tool execution, and the transport does not guarantee a hard timeout for blocked writes. Terminal write failures are not retried because delivery may have been partial.

Partial tool arguments and incremental tool output are not streamed by this adapter. Completed results, including eligible file diffs, remain available; this avoids presenting unapproved partial input or treating independently transformed output chunks as whole-output redaction.

## Tool Locations and File Diffs

Tool-start and permission updates report absolute file locations. Relative paths are resolved lexically against the session's known workspace, with home-directory expansion; this display metadata does not read files or grant access. URI-shaped values and relative paths without an absolute workspace are omitted.

Successful ACP `edit_file` operations can include a full-file diff captured from the client-read content and the exact content sent in the write. The completion retains its textual response as well. Snapshots preserve untouched text and newline conventions; they are not reconstructed from replacement snippets or reread when the result is displayed.

Diffs are omitted when filesystem `post_edit` commands are configured, the checked target changes between read and write, the snapshots contain invalid UTF-8, or the combined JSON-encoded diff exceeds 1 MiB. The edit still runs normally when its diff is omitted. These snapshots are not atomic filesystem history or a guarantee of final content after later runtime hooks or other writers.

`write_file` completions remain text-only: ACP offers no byte-bounded, atomic before-state read, and an optional full read could disconnect the client on a large file and prevent an otherwise valid overwrite. Unknown old content is never presented as a newly created file. Other tools without captured ACP edit metadata also remain text-only.

Captured file contents are presentation-only metadata, excluded from generic result/event JSON, model tool output, and persisted tool messages. The explicit ACP diff returns those unredacted file snapshots to the originating client; tool-response text transformations do not rewrite the diff. Raw output continues to contain only the normal transformed response text.

## CLI Flags

```bash
docker agent serve acp <agent-file>|<registry-ref> [flags]
```

See the [CLI reference](../cli/index.md#docker-agent-serve-acp) for all flags, defaults, and shared runtime options.

## Integration Example

A host application would spawn Docker Agent as a subprocess and communicate via the ACP protocol:

```javascript
// Pseudocode: request() sends newline-delimited JSON-RPC with unique IDs,
// correlates responses, and handles incoming client requests/notifications.
const client = spawnACP(["docker", "agent", "serve", "acp", "./agent.yaml"]);

await client.request("initialize", {
  protocolVersion: 1,
  clientCapabilities: {},
});
const session = await client.request("session/new", {
  cwd: process.cwd(),
  mcpServers: [],
});
await client.request("session/prompt", {
  sessionId: session.sessionId,
  prompt: [{ type: "text", text: "Explain this code" }],
});
```

> [!TIP]
> **When to use ACP**
>
> Use ACP when building **IDE integrations**, **editor plugins**, or any tool that wants to embed a Docker Agent agent as a subprocess. For HTTP-based integrations, use the [API Server](../api-server/index.md) instead.

> [!NOTE]
> **See also**
>
> For HTTP-based agent access, see the [API Server](../api-server/index.md). For agent-to-agent communication, see [A2A Protocol](../a2a/index.md). For exposing agents as MCP tools, see [MCP Mode](../mcp-mode/index.md).
