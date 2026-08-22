---
title: "Troubleshooting"
description: "Common issues and how to resolve them when working with Docker Agent."
keywords: docker agent, ai agents, community, troubleshooting
weight: 20
canonical: https://docs.docker.com/ai/docker-agent/community/troubleshooting/
---

_Common issues and how to resolve them when working with Docker Agent._

## Common Errors

### Context Window Exceeded

Error message: `context_length_exceeded` or similar.

- Use `/compact` in the TUI to summarize and reduce conversation history
- Set `num_history_items` in agent config to limit messages sent to the model
- Switch to a model with larger context (Claude Sonnet 4.5 supports 1M tokens, Gemini up to 2M)
- Break large tasks into smaller conversations

### Max Iterations Reached

The agent hit its `max_iterations` limit without completing the task.

- Increase `max_iterations` in agent config (default is unlimited, but many agents set 20-50)
- Check if the agent is stuck in a loop (enable `--debug` to see tool calls)
- Break complex tasks into smaller steps

### Model Fallback Triggered

When the primary model fails, Docker Agent automatically switches to fallback models. Look for log messages like `"Switching to fallback model"`.

- **429 errors:** Rate limited — the cooldown period keeps using the fallback
- **5xx errors:** Server issues — retries with exponential backoff first, then falls back
- **4xx errors:** Client errors — skips directly to next model

Configure fallback behavior in your agent config:

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    fallback:
      models: [openai/gpt-5-mini, openai/gpt-4o-mini]
      retries: 2 # retries per model for 5xx errors
      cooldown: 1m # how long to stick with fallback after 429
```

## Missing credentials or model errors

When Docker Agent can't find a usable model at startup, it fails fast with an actionable error. The message names the exact next step. `docker agent doctor` is the fastest way to see the full picture — which providers have credentials, whether Docker Model Runner is reachable, and which model `auto` would pick.

### Required environment variables not set

An agent (or a tool it uses) depends on environment variables that aren't configured:

```text
The following environment variables must be set:
 - ANTHROPIC_API_KEY

Provide them using any of these sources:
 - Shell environment:      export ANTHROPIC_API_KEY=<value>
 - Env file:               docker agent run --env-from-file <file> ...
 - Docker Agent env file:  docker agent setup (stores the key in ~/.config/cagent/.env)

See https://docs.docker.com/ai/docker-agent/guides/secrets/ for details.
```

Set the variable through any of the listed [secret sources](../../guides/secrets/index.md). When the missing variable is a model-provider API key, the error also suggests running a local model instead (`docker agent run --model dmr/ai/qwen3 ...`), which needs no API key, and links to the [Set Up a Model](../../getting-started/set-up-a-model/index.md) tutorial.

### No model available (`auto` selection failed)

The `auto` model selector found no configured cloud provider and no usable Docker Model Runner model:

```text
No model is currently available.

To fix this, you can:
  - Pull a Docker Model Runner model, e.g. `docker model pull ai/qwen3`
  - Install Docker Model Runner: https://docs.docker.com/ai/model-runner/get-started/
  - Configure an API key for a cloud provider:
    - anthropic: ANTHROPIC_API_KEY
    - openai: OPENAI_API_KEY
    ...
```

Either configure a cloud provider API key (see [API keys not set](#api-keys-not-set) below) or pull a local model. The [Set Up a Model](../../getting-started/set-up-a-model/index.md) tutorial walks through both paths. Run `docker agent doctor` to see which providers have credentials and whether Docker Model Runner is reachable.

### Docker Model Runner model not pulled

A `dmr/...` model was requested but isn't available locally:

```text
model ai/qwen3 is not pulled in Docker Model Runner

To resolve this, you can:
  - Pull it first: docker model pull ai/qwen3
  - Or choose a model that is already available (see `docker model ls`).
```

If instead you see `cannot query Docker Model Runner at <url>`, Docker Model Runner isn't installed or running — see the [Model Runner get-started guide](https://docs.docker.com/ai/model-runner/get-started/).

> [!TIP]
> **Diagnose before you run**
>
> Run `docker agent doctor` (or `docker agent doctor ./agent.yaml` to include a file's requirements) to check all three issues in one shot. It exits non-zero when something would block a run, making it useful as a CI preflight. See the [CLI reference](../../features/cli/index.md#docker-agent-doctor).

## Debug Mode

The first step for any issue is enabling debug logging. This provides detailed information about what Docker Agent is doing internally.

```bash
# Enable debug logging (writes to ~/.cagent/cagent.debug.log)
$ docker agent run config.yaml --debug

# Write debug logs to a custom file
$ docker agent run config.yaml --debug --log-file ./debug.log

# Enable OpenTelemetry tracing for deeper analysis
$ docker agent run config.yaml --otel
```

> [!TIP]
> Always enable `--debug` when reporting issues. The log file contains detailed traces of API calls, tool executions, and agent interactions.

## Agent Not Responding

### API keys not set

Each model provider requires its own API key as an environment variable:

| Provider      | Environment Variable                                |
| ------------- | --------------------------------------------------- |
| OpenAI        | `OPENAI_API_KEY`                                    |
| Anthropic     | `ANTHROPIC_API_KEY`                                 |
| Google Gemini | `GOOGLE_API_KEY` or `GEMINI_API_KEY`                |
| Mistral       | `MISTRAL_API_KEY`                                   |
| xAI           | `XAI_API_KEY`                                       |
| Nebius        | `NEBIUS_API_KEY`                                    |
| MiniMax       | `MINIMAX_API_KEY`                                   |
| Requesty      | `REQUESTY_API_KEY`                                  |
| OpenRouter    | `OPENROUTER_API_KEY`                                |
| GitHub Copilot | `GITHUB_TOKEN` (PAT with `copilot` scope)          |
| Azure OpenAI  | `AZURE_API_KEY` (override with `token_key`)         |
| AWS Bedrock   | `AWS_BEARER_TOKEN_BEDROCK` or AWS credentials chain |

```bash
# Verify your keys are set
$ env | grep API_KEY
```

### Incorrect model name

Model names must match the provider's naming exactly. Common mistakes:

- Using a deprecated model name (e.g. `gpt-4` instead of `gpt-5-mini` or `gpt-4o`)
- Model references are case-sensitive: `openai/gpt-5-mini` ≠ `openai/GPT-5-mini`

### Network connectivity

If the agent hangs or times out, check that you can reach the provider's API endpoint. Firewalls, VPNs, or proxy settings may block requests. Docker Agent does not evaluate PAC files or URLs directly. When Docker Desktop is running, eligible requests use its PAC adapter before environment proxy settings; `NO_PROXY` does not bypass that selection. Set `DOCKER_AGENT_DISABLE_DESKTOP_PROXY=1` (or `true`, `yes`, or `on`) to restore `HTTP_PROXY`, `HTTPS_PROXY`, `ALL_PROXY`, and `NO_PROXY` routing per request; see [Docker Desktop proxy](../../tools/fetch/index.md#docker-desktop-proxy) for scope and SSRF behavior.

## Tool Execution Failures

### MCP tools not found or failing

- Ensure the MCP tool command is installed and on your `PATH`
- Check file permissions — tools need to be executable
- Test MCP tools independently before integrating with Docker Agent
- For Docker-based MCP tools (`ref: docker:*`), ensure Docker Desktop is running

### Filesystem / shell tool errors

- Verify the agent has the correct toolset configured (`type: filesystem`, `type: shell`)
- Check that the working directory exists and is accessible
- On macOS, ensure terminal has the necessary permissions (e.g., Full Disk Access)

### Tool lifecycle issues

MCP and LSP toolsets are managed by a supervisor that auto-restarts them when they crash or drop their session. The TUI exposes that supervisor through two slash commands:

- `/tools` — the unified tools dialog. Its top section lists every toolset with its current state (`Stopped`, `Starting`, `Ready`, `Degraded`, `Restarting`, `Failed`), restart count, and last error; the bottom section lists every tool the agent can call. Start here whenever a tool seems missing or stuck.
- `/toolset-restart <name>` — force a supervisor-driven reconnect of the named toolset. Useful after completing OAuth, when a remote MCP server has been redeployed, or when a language server like `gopls` is unresponsive.

Remote MCP servers that return `401 invalid_token` (e.g. because the stored OAuth token was revoked or rotated) are now self-healing: Docker Agent silently exchanges the refresh token for a new one when possible, or surfaces an OAuth re-authentication prompt on your next message when refresh is not possible. No more stuck toolsets that require a process restart — but if you want to trigger re-auth immediately, `/toolset-restart <name>` forces it right away.

MCP tools using stdio transport must complete the initialization handshake before becoming available. If tools fail silently:

1. Run `/tools` to see whether the toolset is `Failed` or stuck in `Restarting`, and what the last error was.
2. Enable `--debug` and look for MCP protocol messages in the log
3. Check that the MCP server process starts and responds to `initialize`
4. Verify environment variables required by the tool are set (check `env` and `env_file` in the toolset config)

> [!NOTE]
> **Startup tool-listing timeout**
>
> At startup, Docker Agent queries each toolset for its tool list. If a toolset does not respond within 10 seconds (e.g. a wedged MCP stdio server that never answers `tools/list`), that toolset is skipped with a warning and the remaining toolsets load normally. The sidebar resolves showing whichever tools did load — no infinite spinner. Enable `--debug` to see the warning message, and use `/toolset-restart <name>` once the server becomes responsive.

If a toolset keeps crashing in a tight loop, tune the [`lifecycle`](../../configuration/tools/index.md#toolset-lifecycle) block on the toolset (e.g. raise `backoff.initial`, lower `max_restarts`, or switch to the `best-effort` profile) so a flaky dependency does not amplify into a restart storm.

## Configuration Errors

### YAML syntax issues

Docker Agent validates config at startup and reports errors with line numbers. Common problems:

- Incorrect indentation (YAML is whitespace-sensitive)
- Missing quotes around values containing special characters (`:`, `#`, `{`, `}`)
- Using tabs instead of spaces

### Missing references

- Local agents in `sub_agents` must be defined in the `agents` section (external OCI references like `myorg/agent:tag` are resolved from registries automatically)
- Named model references must exist in the `models` section (or use inline format like `openai/gpt-5`)
- RAG source names referenced by agents must be defined in the `rag` section

### Toolset validation

- The `path` field is valid for `memory` and `tasks` toolsets, and for the agent-level `cache` block
- MCP toolsets need either `command` (stdio), `remote` (Streamable HTTP/SSE), or `ref` (Docker)
- Provider names must be one of: `openai`, `anthropic`, `google`, `amazon-bedrock`, `dmr`, etc.

> [!NOTE]
> **Schema Validation**
>
> Use the [JSON schema](https://github.com/docker/docker-agent/blob/main/agent-schema.json) in your editor for real-time config validation and autocompletion.

## Session &amp; Connectivity Issues

### Downgrade fails with a newer-database error

If an older Docker Agent binary cannot open the session database after an upgrade,
the database may contain a schema migration that the older binary does not know.
Restore a database created by the older version, or use a binary that includes the
migration.

### Port conflicts

When running Docker Agent as an API server or MCP server, ensure the port is not already in use:

```bash
# Check if port 8080 is in use
$ lsof -i :8080

# Use a different port
$ docker agent serve api config.yaml --listen :9090
```

### MCP endpoint accessibility

For remote MCP servers, verify the endpoint is reachable:

```bash
# Test streamable HTTP endpoint
$ curl -v https://mcp-server.example.com/mcp
```

### Session isolation

The API server stores every conversation as a distinct session in the SQLite database (`session.db` by default). Each session is identified by its UUID and only mixes messages when the same session ID is reused. If conversations seem to bleed into each other:

- Make sure each client creates a fresh session via `POST /api/sessions` (don't reuse session IDs across users).
- Confirm `--session-db` points to the path you expect — a stale database from another run can resurface old sessions.
- Use `GET /api/sessions/:id` to inspect what is actually stored, and `DELETE /api/sessions/:id` to clear sessions you don't want anymore.

### HTTP 413: request body too large

Three kinds of process reject an oversized request body with `413 Request Entity Too Large`: `docker agent serve api`, `docker agent serve chat`, and an interactive run's control plane ([`docker agent run --listen`](../../features/api-server/index.md#listen)). They aren't configured the same way: `serve api` and `serve chat` each expose their own `--max-request-size` flag (1 MiB default). A `--listen` control plane has no such flag — it enforces a fixed, non-configurable 1 MiB limit — and no `--auth-token` either. Work through these layers in order:

1. **Identify which server is involved.** `serve api`, `serve chat`, and an attached run's `--listen` control plane are three separate kinds of process. `serve api` and `serve chat` each have their own `--max-request-size` flag and 1 MiB default — check the flags the process that returned the 413 was actually started with. A `--listen` control plane has no `--max-request-size` flag: its 1 MiB cap is fixed.
2. **Measure the serialized request body, not a source file's size.** JSON string escaping and, for any base64-encoded binary content, base64's ~33% expansion both inflate the wire size well past the original file size — a file just under the limit can still push the encoded request over it.
3. **Rule out an intermediary.** If a reverse proxy, gateway, or load balancer sits in front of Docker Agent, it usually enforces its own, independent body-size limit — often with a differently formatted error — and can reject the request before Docker Agent ever sees it.
4. **Confirm who actually returned the error.** A 413 (or a context-length error) can also come from the model provider itself once the request reaches it; that is a separate limit unrelated to `--max-request-size` — see [Context Window Exceeded](#context-window-exceeded) above.
5. **Resolve it.** Once you've confirmed Docker Agent's own server rejected the request, send less content — split it across turns. On `serve api` or `serve chat` you can also restart the server with a deliberately chosen, larger `--max-request-size` (see [API Server](../../features/api-server/index.md#cli-flags) or [Chat Server](../../features/chat-server/index.md#cli-flags)). A `--listen` control plane has no `--max-request-size` flag to raise — sending less content is the only fix.

A few things that catch people out:

- `--max-request-size` is set once at process startup and applies to every request that server handles — it isn't per-request or per-client. Only `serve api` and `serve chat` have it; a `--listen` control plane's 1 MiB cap can't be changed.
- `0` or a negative value falls back to the 1 MiB default; it does not mean "no limit".
- Retrying the same oversized body against the same server won't succeed — the limit doesn't change between requests.
- On all three, the body-size check runs ahead of request authentication, so an oversized request can come back as 413 even without valid credentials.
- Piping stdin into a **local** run (`docker agent run agent.yaml -`) never crosses Docker Agent's own inbound HTTP boundary — Docker Agent may still send that content onward to a model/provider over HTTP, but no request reaches Docker Agent's own server to be measured against a `--max-request-size` cap. Piping stdin into `docker agent run --remote ... -`, however, does cross that boundary: the CLI serializes that stdin text into a native API run request and sends it to whichever Docker Agent server the `--remote` address points at — a `serve api` process or another run's `--listen` control plane, never `serve chat`, which speaks a different protocol — so it's measured against that server's own limit like any other request (a configurable `--max-request-size` for `serve api`, or the fixed 1 MiB cap for a `--listen` control plane). That initial request carries only message text — conversion currently drops any attachment resolved locally (`@path`, `/attach`, `--attach`) for the first message, so it alone can't be the cause of a 413. But a `--remote` run doesn't stay text-only for its whole lifetime: a locally resolved attachment added to a *later* message while the agent is still busy — via the default steer behavior or an explicit follow-up (Alt+Enter) — is forwarded as part of that native API request, counts toward the same limit, and can trigger 413 just like any other oversized request.

> [!WARNING]
> Raising `--max-request-size` increases how much memory an unauthenticated or malicious client can force the server to buffer per request. Pick a value with your deployment's exposure in mind, and pair any non-loopback listener with `--auth-token` (API server) or `--api-key`/`--api-key-env` (chat server). A `--listen` control plane has neither flag — keep it on loopback, a unix socket, or behind an authenticating reverse proxy if it must be reachable from elsewhere.

## Performance Issues

### High memory usage

- Large context windows (64K+ tokens) consume significant memory — consider reducing `max_tokens`
- Use `num_history_items` in agent config to limit conversation history
- For DMR (local models), tune `runtime_flags` for your hardware (e.g., `--ngl` for GPU layers)

### Slow responses

- Check if MCP tools are adding latency (visible in debug logs)
- Use the `/cost` command in TUI to see token usage and identify expensive interactions
- For DMR, consider enabling [speculative decoding](../../providers/dmr/index.md) for faster inference

### Tool resource leaks

Monitor for tools that don't clean up properly — check debug logs for MCP server start/stop lifecycle events. Orphaned tool processes can consume system resources.

## Agent Store Issues

### Pull / push failures

```bash
# Test registry connectivity
$ docker pull docker.io/username/agent:latest

# Verify pulled agent content
$ docker agent share pull docker.io/username/agent:latest
```

### Agent content issues

- Ensure the pushed YAML is valid — run `docker agent run` locally before pushing
- Check that referenced resources (MCP tools, files) are available on the target machine
- For auto-refresh (`--pull-interval`), verify the registry is accessible from the server

## Log Analysis

When reviewing debug logs, search for these key patterns:

| Log Pattern                 | What It Indicates                                                                                |
| --------------------------- | ------------------------------------------------------------------------------------------------ |
| `"Starting runtime stream"` | Agent execution beginning                                                                        |
| `"Tool call"`               | A tool is being executed                                                                         |
| `"Tool call result"`        | Tool execution completed                                                                         |
| `"Stream stopped"`          | Agent finished processing                                                                        |
| `HTTP 429`                  | Rate limiting — consider adding a [fallback model](../../configuration/agents/index.md) |
| `context canceled`          | Operation was interrupted (timeout or user cancel)                                               |
| `[RAG Manager]`             | RAG retrieval operations                                                                         |
| `[Reranker]`                | Reranking operations                                                                             |

> [!WARNING]
> **Still stuck?**
>
> If these steps don't resolve your issue, file a bug on the [GitHub issue tracker](https://github.com/docker/docker-agent/issues) with your debug log attached, or ask on [Slack](https://dockercommunity.slack.com/archives/C09DASHHRU4).
