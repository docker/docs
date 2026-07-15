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
> If the `command` is not in your `PATH`, docker-agent looks it up in the [aqua registry](https://github.com/aquaproj/aqua-registry) and installs it for you. Use `version: "false"` to opt out, or set `DOCKER_AGENT_AUTO_INSTALL=false` globally. See [Auto-Installing Tools](../../configuration/tools/index.md#auto-installing-tools).

## Remote MCP (Streamable HTTP / SSE)

Connect to MCP servers over the network. OAuth flows (including [Dynamic Client Registration](https://datatracker.ietf.org/doc/html/rfc7591)) are handled automatically — docker-agent opens your browser when authentication is required and caches tokens for subsequent sessions. Tokens are refreshed silently when they expire or are revoked server-side; if a silent refresh is not possible, the OAuth prompt reappears on the next message.

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
| `remote.url`            | string  | Base URL of the MCP server. |
| `remote.transport_type` | string  | `streamable` or `sse`. |
| `remote.headers`        | object  | HTTP headers (typically for static auth tokens). |
| `remote.oauth`          | object  | Explicit OAuth client credentials for servers that don't support DCR. See [Remote MCP Servers](../../features/remote-mcp/index.md#oauth-for-servers-without-dynamic-client-registration). |
| `allow_private_ips`     | boolean | Permit remote MCP OAuth helper requests to dial non-public IP addresses. Use only for trusted internal servers. |

For a curated list of public remote MCP endpoints (Linear, GitHub, Vercel, Notion, …) and full OAuth configuration details, see [Remote MCP Servers](../../features/remote-mcp/index.md).

## Embedded Resources

MCP tool results can include embedded resources — images, PDFs, and text files returned directly in the tool response. Docker Agent preserves these as attachments and forwards them to the model as native content blocks:

- **Anthropic** — images become `image` blocks in the `tool_result`; PDFs and other documents become `document` blocks.
- **OpenAI** — images are forwarded as `input_image` data URIs; PDFs as `input_file` data URIs in the tool result content.
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

```yaml
toolsets:
  - type: mcp
    ref: docker:github-official
    tools: ["list_issues", "create_issue", "get_pull_request"]
```

Whitelisting tools improves model accuracy — fewer choices means less confusion.

### Deferred loading

Skip the toolset's startup cost until its tools are actually called:

```yaml
toolsets:
  - type: mcp
    ref: docker:github-official
    defer: true
  # Or defer specific tools within a toolset:
  - type: mcp
    ref: docker:slack
    defer: ["list_channels", "search_messages"]
```

### Custom instructions

```yaml
toolsets:
  - type: mcp
    ref: docker:github-official
    instruction: |
      Use these tools to manage GitHub issues.
      Always check for existing issues before creating new ones.
      Label new issues with 'triage' by default.
```

### TOON-encoded outputs

Re-encode verbose JSON outputs as the compact [TOON](https://github.com/alpkeskin/gotoon) format to save context budget. Typically yields 30–60% smaller payloads on list/search tools.

`toon` is a regex string that is matched against tool names. Any tool whose name matches the pattern has its JSON output transparently re-encoded as TOON before it is shown to the model. The re-encoding reduces schema verbosity, which is especially useful when a model struggles with large or repetitive tool output.

```yaml
toolsets:
  - type: mcp
    ref: docker:github-official
    toon: ".*"            # toonify every tool from this server
  - type: mcp
    command: my-server
    toon: "list_.*,get_.*" # only toonify list_/get_ tools
```

The value is a comma-separated list of regexes (or a single regex). A tool name must match at least one pattern to be re-encoded. Setting `toon: ".*"` re-encodes all tools from that toolset.

See [`examples/github-toon.yaml`](https://github.com/docker/docker-agent/blob/main/examples/github-toon.yaml) for a practical example using the GitHub MCP server.

### Per-toolset model routing

Process tool results from this toolset with a different (typically cheaper / faster) model. The override is one-shot — subsequent turns return to the agent's primary model:

```yaml
toolsets:
  - type: mcp
    ref: docker:github-official
    model: openai/gpt-4o-mini
```

See [Per-Toolset Model Routing](../../configuration/tools/index.md#per-toolset-model-routing).

### Lifecycle (auto-restart, profiles)

Local stdio and remote MCP servers are supervised: crashed servers reconnect automatically with exponential backoff. **Remote** MCP servers (Streamable HTTP / SSE) also reconnect after idle/clean connection closes — services like Notion and Linear periodically close idle connections, and docker-agent reconnects transparently. Tune the policy with the `lifecycle` block:

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
