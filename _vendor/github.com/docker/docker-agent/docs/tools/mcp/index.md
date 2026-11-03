---
title: "MCP Tool"
description: "Extend agents with external tools via the Model Context Protocol."
keywords: docker agent, ai agents, tools, toolsets, mcp tool
linkTitle: "MCP"
weight: 130
canonical: https://docs.docker.com/ai/docker-agent/tools/mcp/
aliases:
  - /ai/docker-agent/integrations/mcp/
---

_Extend agents with external tools via the Model Context Protocol (MCP)._

## Overview

The `mcp` toolset connects your agent to any MCP server — a process or remote service that exposes tools, resources, and prompts over the [Model Context Protocol](https://modelcontextprotocol.io/). Three flavours are supported:

| Flavour | Transport | Best for |
| --- | --- | --- |
| **Docker MCP** | Container via the [MCP Gateway](https://github.com/docker/mcp-gateway) | Curated, sandboxed servers from the [Docker MCP Catalog](https://hub.docker.com/u/mcp) |
| **Local stdio** | Subprocess over stdin/stdout | Custom or community MCP servers run from a binary or `npx`/`pip` package |
| **Remote** | Streamable HTTP or SSE | Cloud services with hosted MCP endpoints (Linear, Notion, Atlassian, …) |

> [!NOTE]
> **What is MCP?**
>
> The [Model Context Protocol](https://modelcontextprotocol.io/) is an open standard for connecting AI tools. Docker Agent can both _use_ MCP servers (this page) and _expose_ agents as MCP servers — see [MCP Mode](../../features/mcp-mode/index.md).

## Docker MCP (Recommended)

Run MCP servers as secure Docker containers via the MCP Gateway. The `ref: docker:<name>` syntax pulls a curated definition from the Docker MCP Catalog:

```yaml
toolsets:
  - type: mcp
    ref: docker:duckduckgo        # web search
  - type: mcp
    ref: docker:github-official   # GitHub integration
    tools: ["list_issues", "create_issue"]
```

Browse available servers at the [Docker MCP Catalog](https://hub.docker.com/u/mcp).

| Property      | Type   | Description                                                      |
| ------------- | ------ | ---------------------------------------------------------------- |
| `ref`         | string | Docker MCP reference (`docker:name`) or a name from the [reusable `mcps:`](../../configuration/overview/index.md#reusable-mcp-servers-mcps) block. |
| `tools`       | array  | Optional whitelist — only expose these tools to the model.       |
| `instruction` | string | Custom instructions injected into the agent's context.           |
| `config`      | any    | MCP server-specific configuration passed during initialization.  |
| `working_dir` | string | Working directory for the MCP gateway subprocess. Only applies when the catalog entry runs as a local process (not remote). Relative paths are resolved against the agent's working directory. Supports `${env.VAR}` (canonical), plus `~` and shell-style `$VAR`/`${VAR}` expansion ([details](../../configuration/overview/index.md#variable-expansion-in-config-fields)). |

## Local MCP (stdio)

Run MCP servers as local processes communicating over stdin/stdout:

```yaml
toolsets:
  - type: mcp
    command: python
    args: ["-m", "mcp_server"]
    tools: ["search", "fetch"]
    env:
      API_KEY: value
```

| Property      | Type   | Description |
| ------------- | ------ | ----------- |
| `command`     | string | Command to execute the MCP server. |
| `args`        | array  | Command arguments. |
| `tools`       | array  | Optional whitelist — only expose these tools. |
| `env`         | object | Environment variables (key-value pairs). |
| `working_dir` | string | Working directory for the MCP server process. Relative paths are resolved against the agent's working directory. Defaults to the agent's working directory when omitted. Supports `${env.VAR}` (canonical), plus `~` and shell-style `$VAR`/`${VAR}` expansion ([details](../../configuration/overview/index.md#variable-expansion-in-config-fields)). |
| `instruction` | string | Custom instructions injected into the agent's context. |
| `version`     | string | Package reference for [auto-installing](../../configuration/tools/index.md#auto-installing-tools) the command binary. |

> [!TIP]
> **Auto-installation**
>
> If the `command` is not in your `PATH`, Docker Agent looks it up in the [aqua registry](https://github.com/aquaproj/aqua-registry) and installs it for you. Use `version: "false"` to opt out, or set `DOCKER_AGENT_AUTO_INSTALL=false` globally. See [Auto-Installing Tools](../../configuration/tools/index.md#auto-installing-tools).

## Remote MCP (Streamable HTTP / SSE)

Connect to MCP servers over the network. OAuth flows (including [Dynamic Client Registration](https://datatracker.ietf.org/doc/html/rfc7591)) are handled automatically — Docker Agent opens your browser when authentication is required and caches tokens for subsequent sessions. Tokens are refreshed silently when they expire or are revoked server-side; if a silent refresh is not possible, the OAuth prompt reappears on the next message.

```yaml
toolsets:
  - type: mcp
    remote:
      url: "https://mcp.linear.app/mcp"
      transport_type: "streamable"               # or "sse" for legacy servers
      headers:
        Authorization: "Bearer ${env.LINEAR_TOKEN}"
    # Optional: allow OAuth helper requests to reach private/internal IPs.
    allow_private_ips: false
    tools: ["search_issues", "create_issue"]
```

| Property                | Type    | Description |
| ----------------------- | ------- | ----------- |
| `remote.url`            | string  | URL of the MCP server. Accepts `https://`, `http://`, and `unix://` (Unix domain socket) schemes. |
| `remote.transport_type` | string  | `streamable` or `sse`. |
| `remote.headers`        | object  | HTTP headers sent on every request. Values support `${env.VAR}` and `${headers.NAME}` placeholders, resolved per request. See [Remote MCP Servers](../../features/remote-mcp/index.md#configuration) for details. |
| `remote.oauth`          | object  | Explicit OAuth client credentials for servers that don't support DCR. See [Remote MCP Servers](../../features/remote-mcp/index.md#oauth-for-servers-without-dynamic-client-registration). |
| `allow_private_ips`     | boolean | Permit remote MCP OAuth helper requests to dial non-public IP addresses. Use only for trusted internal servers. |

For a curated list of public remote MCP endpoints (Linear, GitHub, Vercel, Notion, …) and full OAuth configuration details, see [Remote MCP Servers](../../features/remote-mcp/index.md).

When Docker Desktop is running, eligible MCP OAuth discovery, token, and helper requests use its PAC adapter before environment proxy settings, but remote MCP Streamable HTTP/SSE transport does not. Set `DOCKER_AGENT_DISABLE_DESKTOP_PROXY=1` (or `true`, `yes`, or `on`) to restore standard `HTTP_PROXY`, `HTTPS_PROXY`, `ALL_PROXY`, and `NO_PROXY` routing; `NO_PROXY` does not bypass Docker Desktop PAC selection. Docker Agent does not evaluate PAC files or URLs directly—see [Docker Desktop proxy](../fetch/index.md#docker-desktop-proxy).

## MCP Prompts

MCP servers can expose **prompts** — named, parameterized templates that the server provides via the `/prompts` endpoint. Docker Agent discovers these at toolset startup and registers them as **slash commands** in the TUI, so you can invoke them directly from the input box.

```text
# Type / to see available prompts alongside built-in commands
/review         # invoke an MCP prompt named "review"
/summarize My text here   # invoke with the first argument filled in
```

**How it works:**

- Each MCP prompt appears in the command palette (accessible via <kbd>Ctrl</kbd>+<kbd>K</kbd>) under the **MCP Prompts** category.
- Typing `/<prompt-name>` in the input box invokes the prompt immediately.
- If the prompt declares arguments and you provide text after the slash command, that text is mapped to the first declared argument.
- If a required argument is missing, Docker Agent opens the argument input dialog before running the prompt.
- When no argument is needed or all required arguments are supplied, the prompt runs immediately.

> [!NOTE]
> MCP prompt discovery requires a YAML-declared `mcp` toolset. Prompts from servers activated through the [Docker MCP Catalog](../../tools/mcp-catalog/index.md) (`ref: docker:<name>`) are not currently surfaced.

## Embedded Resources

MCP tool results can include embedded resources — images, PDFs, and text files returned directly in the tool response. Docker Agent preserves these as attachments and forwards them to the model as native content blocks:

- **Anthropic** — images become `image` blocks in the `tool_result`; PDFs and other documents become `document` blocks.
- **OpenAI** — images are forwarded as `input_image` data URIs; PDFs as `input_file` data URIs in the tool result content.
- **OpenAI Chat Completions** (`api_type: openai_chatcompletions`, the default, and OpenAI-compatible providers) — tool messages only accept text, so images (including image documents) are sent as text in the tool message as their data URI or URL. PDFs and other documents are attached in a follow-up user message labeled "Attached content from tool result:".
- **Bedrock** and **Gemini** — receive equivalent provider-native representations.

No configuration is required. When an MCP server returns an embedded resource alongside its text output, the resource is automatically attached and sent to the model on the next turn. This is useful for MCP servers that generate charts, export PDFs, or return binary data as part of their responses.

## Reusable Definitions (`mcps:`)

Repeated MCP server configurations can be hoisted into the top-level `mcps:` section and referenced by name with `{type: mcp, ref: <name>}`:

```yaml
mcps:
  github:
    remote:
      url: https://api.githubcopilot.com/mcp
      transport_type: sse
  playwright:
    command: npx
    args: ["-y", "@modelcontextprotocol/server-playwright"]

agents:
  root:
    model: openai/gpt-5
    toolsets:
      - type: mcp
        ref: github
      - type: mcp
        ref: playwright
```

See [Reusable MCP Servers](../../configuration/overview/index.md#reusable-mcp-servers-mcps) for the full reference.

## Common Options

These properties apply to every MCP toolset regardless of flavour:

### Tool filtering

Use `tools:` to expose only the operations your agent needs. See [Tool Filtering](../../configuration/tools/index.md#tool-filtering) for syntax and examples.

### Deferred loading

Use `defer: true` or a list of tool names to load tools on demand. See [Deferred Tool Loading](../../configuration/tools/index.md#deferred-tool-loading) for startup and discovery behavior.

### Custom instructions

Use `instruction:` to replace or extend a toolset's guidance. See [Tool Instructions](../../configuration/tools/index.md#tool-instructions), including the `{ORIGINAL_INSTRUCTIONS}` placeholder.

### TOON-encoded outputs

Use `toon:` to re-encode matching tools' JSON object outputs in a compact format. See [TOON-Encoded Tool Outputs](../../configuration/tools/index.md#toon-encoded-tool-outputs) for pattern syntax and which outputs pass through unchanged, or the [GitHub TOON example](https://github.com/docker/docker-agent/blob/main/examples/github-toon.yaml).

### Per-toolset model routing

Use `model:` to process tool results with a different model for one turn. See [Per-Toolset Model Routing](../../configuration/tools/index.md#per-toolset-model-routing) for examples and how competing overrides are resolved.

### Lifecycle (auto-restart, profiles)

Local stdio and remote MCP servers are supervised: crashed servers reconnect automatically with exponential backoff. **Remote** MCP servers (Streamable HTTP / SSE) also reconnect after idle/clean connection closes — services like Notion and Linear periodically close idle connections, and Docker Agent reconnects transparently. Tune the policy with the `lifecycle` block:

```yaml
toolsets:
  - type: mcp
    ref: docker:duckduckgo
    lifecycle:
      profile: resilient   # default; auto-restart with backoff
  - type: mcp
    command: docker
    args: ["mcp", "gateway"]
    lifecycle:
      profile: strict      # fail-fast: required, no retries
```

See [Toolset Lifecycle](../../configuration/tools/index.md#toolset-lifecycle) for all profiles and tuning knobs, and [`/toolset-restart`](../../features/tui/index.md) to force a reconnect from the TUI.

**Startup failure behaviour:** local MCP failures (missing binary, connection refused, bad auth) fail fast — each turn retries immediately with no artificial delay. Remote MCP servers (Streamable HTTP / SSE) that respond with one of a fixed set of retryable HTTP statuses — 429 Too Many Requests, 408 Request Timeout, 500/502/503/504, or 529 (Anthropic-style "overloaded") — are paced by the same [bounded exponential backoff gate](../rag/index.md#indexing-failures-retries-and-backoff) that RAG embedding calls use, so a temporarily-overloaded remote MCP server does not trigger a new connect attempt on every agent turn. This is a fixed enumeration, not a full 5xx range: less-common codes such as 501, 505, or the Cloudflare 520–527 family do not arm the gate. MCP paces every connection attempt on this set; RAG's own trigger set differs slightly — 429 arms the gate on the very first failure, while 408/5xx only arm it once every file in an indexing run has hit one of these retryable statuses with none indexed successfully (see [RAG's retry-policy table](../rag/index.md#what-triggers-backoff)). This pacing also applies when toolsets are wrapped in [code mode](../../features/code-mode/index.md#limits--security-notes): a retryable failure in the degraded subset paces that subset's retry the same way, while the composite's healthy tools stay available.

## Combined Example

```yaml
mcps:
  github:
    remote:
      url: https://api.githubcopilot.com/mcp
      transport_type: sse

agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: Full-featured developer assistant
    instruction: You are an expert developer.
    toolsets:
      # Docker MCP catalog entry
      - type: mcp
        ref: docker:duckduckgo

      # Reusable definition from the top-level mcps: block
      - type: mcp
        ref: github
        tools: ["list_issues", "create_issue"]
        toon: "list_.*"

      # Local stdio server with auto-install
      - type: mcp
        command: gopls
        version: "golang/tools@v0.21.0"
        args: ["mcp"]

      # Remote MCP with OAuth (handled automatically)
      - type: mcp
        remote:
          url: "https://mcp.linear.app/mcp"
          transport_type: "streamable"
        instruction: Use Linear for issue tracking.
```

> [!WARNING]
> **Toolset order matters**
>
> If multiple toolsets provide a tool with the same name, the first one wins: the duplicate from the later toolset is ignored and a warning identifies both toolsets. Order your toolsets intentionally. To keep both tools callable, give the MCP toolset a unique `name:` (its tools are then exposed as `<name>_<tool>`) or restrict the overlapping toolset with its `tools:` filter.

## See Also

- [Tool Configuration](../../configuration/tools/index.md) — full reference for every toolset type, plus shared options (lifecycle, TOON, model routing, …).
- [Reusable MCP Servers](../../configuration/overview/index.md#reusable-mcp-servers-mcps) — the top-level `mcps:` block.
- [Remote MCP Servers](../../features/remote-mcp/index.md) — catalog of public remote MCP endpoints + OAuth recipes.
- [MCP Mode](../../features/mcp-mode/index.md) — expose your own agents as MCP tools to Claude Desktop, Claude Code, etc.
- [Auto-Installing Tools](../../configuration/tools/index.md#auto-installing-tools) — automatic installation of MCP server binaries.
