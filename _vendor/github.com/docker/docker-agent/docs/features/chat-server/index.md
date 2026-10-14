---
title: "Chat Server"
description: "Expose your agents through an OpenAI-compatible Chat Completions API so any tool that already speaks OpenAI can drive a Docker Agent agent."
keywords: docker agent, ai agents, features, chat server
weight: 90
canonical: https://docs.docker.com/ai/docker-agent/features/chat-server/
---

_Expose your agents through an OpenAI-compatible Chat Completions API so any tool that already speaks OpenAI can drive a Docker Agent agent._

## Overview

The `docker agent serve chat` command starts an HTTP server that exposes one or
more agents through an **OpenAI-compatible Chat Completions API** at
`/v1/chat/completions` and `/v1/models`. Any client that already speaks the
OpenAI protocol — for example
[Open WebUI](https://github.com/open-webui/open-webui), `curl`, the OpenAI
Python SDK, or LangChain — can drive a Docker Agent agent without any custom
integration.

```bash
# Single agent — exposed as the model `root`
$ docker agent serve chat agent.yaml

# Multi-agent config — every agent in the team becomes a model
$ docker agent serve chat ./team.yaml

# Pick a specific agent from a multi-agent config
$ docker agent serve chat ./team.yaml --agent reviewer

# Run an agent straight from the registry
$ docker agent serve chat myorg/agent:tag --listen 127.0.0.1:9090

# Require a Bearer token, sourced from an env var
$ docker agent serve chat agent.yaml --api-key-env CHAT_BEARER_TOKEN
```

> [!TIP]
> **When to use chat server vs. API server**
>
> Use the **chat server** when you want to plug Docker Agent into existing OpenAI-compatible tooling (chat UIs, IDE integrations, OpenAI SDK clients). Use the [API server](../api-server/index.md) when you want full control over sessions, agent execution, tool-call confirmations, and streamed runtime events.

## Endpoints

The OpenAI-compatible endpoints live under the `/v1` prefix to match the
OpenAI API surface. The OpenAPI specification is served at the top level so it
can be discovered without authentication.

| Method | Path                   | Description                                                            |
| ------ | ---------------------- | ---------------------------------------------------------------------- |
| `GET`  | `/v1/models`           | List the agents that this server exposes as models                     |
| `POST` | `/v1/chat/completions` | Send messages and receive a completion (regular or streaming)          |
| `GET`  | `/openapi.json`        | OpenAPI specification for the chat server                              |

The model identifier in `POST /v1/chat/completions` is the **agent name**.
For a single-agent config that's typically `root`; for a multi-agent config,
each named agent becomes its own selectable model.

## Quick Start

```bash
# 1. Start the server
$ docker agent serve chat agent.yaml
Listening on 127.0.0.1:8083
OpenAI-compatible chat completions endpoint: http://127.0.0.1:8083/v1/chat/completions

# 2. List exposed agents (models)
$ curl http://127.0.0.1:8083/v1/models
{"object":"list","data":[{"id":"root","object":"model","owned_by":"docker-agent"}]}

# 3. Send a chat request
$ curl http://127.0.0.1:8083/v1/chat/completions \
    -H 'Content-Type: application/json' \
    -d '{
      "model": "root",
      "messages": [{"role": "user", "content": "Hello!"}]
    }'
```

### Streaming

Set `"stream": true` in the request body to receive a Server-Sent Events
(SSE) stream of OpenAI-format `chat.completion.chunk` deltas:

```bash
$ curl -N http://127.0.0.1:8083/v1/chat/completions \
    -H 'Content-Type: application/json' \
    -d '{
      "model": "root",
      "stream": true,
      "messages": [{"role": "user", "content": "Stream a poem"}]
    }'
```

### Drive it from the OpenAI Python SDK

Because the wire format is OpenAI-compatible, point any OpenAI client at the
chat server's `base_url` and use the agent name as the model:

```python
from openai import OpenAI

client = OpenAI(
    base_url="http://127.0.0.1:8083/v1",
    api_key="not-needed-when-no-api-key-flag",  # required by the SDK, ignored if no auth
)

resp = client.chat.completions.create(
    model="root",
    messages=[{"role": "user", "content": "Hello!"}],
)
print(resp.choices[0].message.content)
```

## Server-side Conversation Caching

By default the server is **stateless**: every request must contain the full
message history, exactly like OpenAI's API. Enable server-side caching by
setting `--conversations-max` to a positive value, then send a stable
`X-Conversation-Id` header on each request:

```bash
$ docker agent serve chat agent.yaml --conversations-max 100 --conversation-ttl 30m
```

```bash
$ curl http://127.0.0.1:8083/v1/chat/completions \
    -H 'Content-Type: application/json' \
    -H 'X-Conversation-Id: my-thread-1' \
    -d '{
      "model": "root",
      "messages": [{"role": "user", "content": "Remember my name is Alice"}]
    }'

$ curl http://127.0.0.1:8083/v1/chat/completions \
    -H 'Content-Type: application/json' \
    -H 'X-Conversation-Id: my-thread-1' \
    -d '{
      "model": "root",
      "messages": [{"role": "user", "content": "What is my name?"}]
    }'
```

Cached conversations are evicted after `--conversation-ttl` of inactivity, or
when the cache hits `--conversations-max` items (oldest entries are evicted
first).

### Failure-safe caching

When a request fails — for example because the model returns an error or the `--request-timeout` expires — the conversation cache is **not updated**. The server clones the cached session before processing each request and only commits the updated session when the turn completes successfully. This means:

- A failed turn leaves the conversation in the same state it was before the request.
- Clients can safely retry with the same `X-Conversation-Id` after a failure.
- Transient errors do not corrupt the conversation history.

## Authentication

The chat server defaults to loopback binding. A non-loopback `--listen` address requires `--api-key`, `--api-key-env`, or the explicit `--insecure-no-auth` override. An environment variable selected by `--api-key-env` must be set and non-empty.

To require a Bearer
token, pass `--api-key` (literal value) or `--api-key-env` (name of an
environment variable that holds the value):

```bash
$ docker agent serve chat agent.yaml --api-key-env CHAT_BEARER_TOKEN
```

Clients must then send an `Authorization: Bearer <token>` header on every
request to `/v1/*`. Both `/v1/models` and `/v1/chat/completions` are
protected once a key is set.

> [!WARNING]
> **Public exposure**
>
> The default listen address is `127.0.0.1:8083`. Non-loopback binding is rejected unless `--api-key`, `--api-key-env`, or `--insecure-no-auth` is supplied. Use the insecure override only behind a trusted authentication boundary.

## Tool safety

The chat server resolves its safety policy in this order: `--safety`, agent configuration, runtime configuration, then `restricted`. Cached conversations retain the more restrictive of their prior policy and the server policy, so a continuation cannot regain permissions after the server policy becomes stricter.

## CORS

CORS is **disabled by default**. To allow a browser-based client to call the
server, set `--cors-origin` to the exact origin (scheme + host + port) that
should be allowed:

```bash
$ docker agent serve chat agent.yaml --cors-origin https://my-ui.example.com
```

## CLI Flags

```bash
docker agent serve chat <agent-file>|<registry-ref> [flags]
```

| Flag                          | Default            | Description                                                                                                       |
| ----------------------------- | ------------------ | ----------------------------------------------------------------------------------------------------------------- |
| `-a, --agent <name>`          | (all agents)       | Name of the agent to expose. If omitted, every agent in the config is exposed as a separate model.                |
| `-l, --listen <addr>`         | `127.0.0.1:8083`   | Address to listen on.                                                                                             |
| `--cors-origin <origin>`      | (none)             | Allowed CORS origin (e.g. `https://example.com`). Empty disables CORS.                                            |
| `--api-key <token>`           | (none)             | Required Bearer token clients must present (`Authorization: Bearer <token>`). Empty disables auth.                |
| `--api-key-env <name>`        | (none)             | Read the required API key from this non-empty environment variable.                                               |
| `--insecure-no-auth`          | `false`            | Permit unauthenticated non-loopback binding. Use only behind a trusted authentication boundary.                  |
| `--safety <policy>`           | `restricted`       | Tool safety policy. CLI value overrides agent/runtime configuration.                                               |
| `--max-request-size <bytes>`  | `1048576` (1 MiB)  | Maximum request body size in bytes. Requests whose body exceeds this limit are rejected with HTTP 413 (Request Entity Too Large) — see [Troubleshooting: HTTP 413](../../community/troubleshooting/index.md#http-413-request-body-too-large) if you hit this. |
| `--request-timeout <dur>`     | `5m`               | Per-request timeout (covers model + tool calls + streaming).                                                      |
| `--conversations-max <n>`     | `0`                | Cache up to N conversations server-side, keyed by `X-Conversation-Id`. `0` disables — clients must resend history. |
| `--conversation-ttl <dur>`    | `30m`              | Idle TTL after which a cached conversation is evicted.                                                            |
| `--max-idle-runtimes <n>`     | `4`                | Maximum number of idle runtimes pooled per agent. `0` disables pooling.                                           |

All [runtime configuration flags](../cli/index.md#runtime-configuration-flags)
(`--working-dir`, `--env-from-file`, `--models-gateway`, `--hook-*`, …) are
also accepted.

> [!NOTE]
> **What `--max-request-size` does and doesn't cover**
>
> This is a finite, process-wide cap on one serialized inbound HTTP request body — it isn't a model context-window limit, and raising it doesn't increase what a provider/model accepts or how large a local attachment/prompt file can be. A larger cap also means the server buffers more memory per request from an unauthenticated or malicious client, so weigh that against your deployment's exposure. If a reverse proxy or gateway sits in front of this server, it may enforce its own, lower cap regardless of this flag. See [Troubleshooting: HTTP 413](../../community/troubleshooting/index.md#http-413-request-body-too-large) for full diagnosis.

## Image Inputs

Messages can include OpenAI-style `image_url` content parts alongside `text`:

```json
{"type": "image_url", "image_url": {"url": "data:image/png;base64,..."}}
```

A `data:` URL embeds the image bytes directly in the JSON body, so its base64-encoded size counts toward `--max-request-size` above like any other request content. A remote `http(s)://` URL is passed through to the selected model provider rather than fetched by the chat server itself — whether it works depends on that provider and model: some accept a remote URL directly, others only accept `data:` URLs, and a provider/model without image support drops the part. Don't assume a remote image URL will be fetched or rendered universally; verify against the specific provider/model you've configured.

## Open WebUI Integration

Open WebUI can talk to any OpenAI-compatible endpoint. To plug Docker Agent
in:

1. Start the chat server, optionally with auth:

    ```bash
    $ docker agent serve chat agent.yaml \
        --listen 127.0.0.1:8083 \
        --cors-origin http://localhost:3000 \
        --api-key-env OPEN_WEBUI_TOKEN
    ```

2. In Open WebUI, add an OpenAI-compatible connection:

    - **API Base URL:** `http://127.0.0.1:8083/v1`
    - **API Key:** the value of `OPEN_WEBUI_TOKEN`

3. Each agent in your config appears as a selectable model.

> [!NOTE]
> **See also**
>
> For the Docker Agent–native HTTP API (sessions, tool-call confirmation, runtime events), see the [API Server](../api-server/index.md). For full CLI flag documentation, see the [CLI Reference](../cli/index.md#docker-agent-serve-chat).
