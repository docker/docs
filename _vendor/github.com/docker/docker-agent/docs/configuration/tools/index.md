---
title: "Tool Configuration"
description: "Complete reference for configuring built-in tools, MCP tools, and Docker-based tools."
keywords: docker agent, ai agents, configuration, yaml, tool configuration
linkTitle: "Tool Config"
weight: 50
canonical: https://docs.docker.com/ai/docker-agent/configuration/tools/
aliases:
  - /ai/docker-agent/reference/toolsets/
---

_Complete reference for configuring built-in tools, MCP tools, and Docker-based tools._

## Built-in Tools

Built-in tools are included with docker-agent and require no external dependencies. Add them to your agent's `toolsets` list by `type`. Each tool's dedicated page covers its full configuration options, available operations, and examples.

| Type | Description | Page |
| --- | --- | --- |
| `filesystem` | Read, write, list, search, navigate | [Filesystem](../../tools/filesystem/index.md) |
| `shell` | Execute shell commands synchronously | [Shell](../../tools/shell/index.md) |
| `background_jobs` | Run and manage long-running shell commands | [Background Jobs](../../tools/background-jobs/index.md) |
| `think` | Reasoning scratchpad | [Think](../../tools/think/index.md) |
| `plan` | Shared persistent scratchpad for multi-agent collaboration | [Plan](../../tools/plan/index.md) |
| `session_plan` | Per-session markdown plan for the draft-review-execute workflow | [Session Plan](../../tools/session_plan/index.md) |
| `session_context` | Reference a previous session as context (read-only) | [Session Context](../../tools/session_context/index.md) |
| `todo` | Task list management | [Todo](../../tools/todo/index.md) |
| `memory` | Persistent key-value storage (SQLite) | [Memory](../../tools/memory/index.md) |
| `tasks` | Persistent task database shared across sessions | [Tasks](../../tools/tasks/index.md) |
| `fetch` | HTTP `GET` requests with text/markdown/html output | [Fetch](../../tools/fetch/index.md) |
| `script` | Custom shell scripts as tools | [Script](../../tools/script/index.md) |
| `lsp` | Language Server Protocol integration | [LSP](../../tools/lsp/index.md) |
| `api` | Custom HTTP API tools | [API](../../tools/api/index.md) |
| `openapi` | Import every operation of an OpenAPI 3.x document as tools | [OpenAPI](../../tools/openapi/index.md) |
| `rag` | Retrieval-augmented generation over indexed sources | [RAG](../../tools/rag/index.md) |
| `model_picker` | Let the agent pick between several models per turn | [Model Picker](../../tools/model-picker/index.md) |
| `user_prompt` | Interactive user input | [User Prompt](../../tools/user-prompt/index.md) |
| `open_url` | Open a fixed URL in the user's default browser | [Open URL](../../tools/open-url/index.md) |
| `transfer_task` | Delegate to sub-agents (auto-enabled) | [Transfer Task](../../tools/transfer-task/index.md) |
| `background_agents` | Parallel sub-agent dispatch | [Background Agents](../../tools/background-agents/index.md) |
| `handoff` | Local conversation handoff to another agent in the same config (auto-enabled by `handoffs:`) | [Handoff](../../tools/handoff/index.md) |
| `a2a` | A2A remote agent connection | [A2A](../../tools/a2a/index.md) |
| `mcp_catalog` | Discover and activate remote MCP servers from the Docker MCP Catalog on demand | [MCP Catalog](../../tools/mcp-catalog/index.md) |

**Example:**

```yaml
toolsets:
  - type: filesystem
  - type: shell
  - type: background_jobs
  - type: think
  - type: todo
  - type: memory
    path: ./dev.db
```

## MCP Tools

Extend agents with external tools via the [Model Context Protocol](https://modelcontextprotocol.io/). For a standalone overview of the `mcp` toolset see the [MCP tool page](../../tools/mcp/index.md).

> [!TIP]
> **Reusable MCP definitions**
>
> Repeated MCP server definitions can be hoisted into the top-level `mcps:` section and referenced by name with `{type: mcp, ref: <name>}`. See [Reusable MCP Servers](../overview/index.md#reusable-mcp-servers-mcps).

### Docker MCP (Recommended)

Run MCP servers as secure Docker containers via the [MCP Gateway](https://github.com/docker/mcp-gateway):

```yaml
toolsets:
  - type: mcp
    ref: docker:duckduckgo # web search
  - type: mcp
    ref: docker:github-official # GitHub integration
```

Browse available tools at the [Docker MCP Catalog](https://hub.docker.com/search?q=&type=mcp).

| Property      | Type   | Description                                                      |
| ------------- | ------ | ---------------------------------------------------------------- |
| `ref`         | string | Docker MCP reference (`docker:name`)                             |
| `tools`       | array  | Optional: only expose these tools                                |
| `instruction` | string | Custom instructions injected into the agent's context            |
| `config`      | any    | MCP server-specific configuration (passed during initialization) |
| `working_dir` | string | Working directory for the MCP gateway subprocess. Only applies when the catalog entry runs as a local process (not remote). Relative paths are resolved against the agent's working directory. Supports `${env.VAR}` (canonical), plus `~` and shell-style `$VAR`/`${VAR}` expansion ([details](../overview/index.md#variable-expansion-in-config-fields)). |

### Local MCP (stdio)

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

| Property | Type | Description |
| --- | --- | --- |
| `command` | string | Command to execute the MCP server |
| `args` | array | Command arguments |
| `tools` | array | Optional: only expose these tools |
| `env` | object | Environment variables (key-value pairs) |
| `working_dir` | string | Working directory for the MCP server process. Relative paths are resolved against the agent's working directory. Defaults to the agent's working directory when omitted. Supports `${env.VAR}` (canonical), plus `~` and shell-style `$VAR`/`${VAR}` expansion ([details](../overview/index.md#variable-expansion-in-config-fields)). |
| `instruction` | string | Custom instructions injected into the agent's context |
| `version` | string | Package reference for [auto-installing](#auto-installing-tools) the command binary |

### Remote MCP (Streamable HTTP / SSE)

Connect to MCP servers over the network:

```yaml
toolsets:
  - type: mcp
    remote:
      url: "https://mcp-server.example.com"
      transport_type: "streamable"
      headers:
        Authorization: "Bearer your-token"
    # Optional: allow OAuth helper requests to reach private/internal IPs.
    allow_private_ips: true
    tools: ["search_web", "fetch_url"]
```

| Property                | Type    | Description                                                                                                           |
| ----------------------- | ------- | --------------------------------------------------------------------------------------------------------------------- |
| `remote.url`            | string  | URL of the MCP server. Accepts `https://`, `http://`, and `unix://` (Unix domain socket) schemes.                     |
| `remote.transport_type` | string  | `streamable` or `sse`                                                                                                 |
| `remote.headers`        | object  | HTTP headers (typically for auth)                                                                                     |
| `allow_private_ips`     | boolean | Permit remote MCP OAuth helper requests to dial non-public IP addresses. Use only for trusted internal servers.        |

## Auto-Installing Tools

When configuring MCP or LSP tools that require a binary command, docker agent can **automatically download and install** the command if it's not already available on your system. This uses the [aqua registry](https://github.com/aquaproj/aqua-registry) — a curated index of CLI tool packages.

### How It Works

1. When a toolset with a `command` is loaded, docker agent checks if the command is available in your `PATH`
2. If not found, it checks the docker agent tools directory (`~/.cagent/tools/bin/`)
3. If still not found, it looks up the command in the aqua registry and installs it automatically

### Explicit Package Reference

Use the `version` property to specify exactly which package to install:

```yaml
toolsets:
  - type: mcp
    command: gopls
    version: "golang/tools@v0.21.0"
    args: ["mcp"]
  - type: lsp
    command: rust-analyzer
    version: "rust-lang/rust-analyzer@2024-01-01"
    file_types: [".rs"]
```

The format is `owner/repo` or `owner/repo@version`. When a version is omitted, the latest release is used.

### Automatic Detection

If the `version` property is not set, docker agent tries to auto-detect the package from the command name by searching the aqua registry:

```yaml
toolsets:
  - type: mcp
    command: gopls  # auto-detected as golang/tools
    args: ["mcp"]
```

### Checksum Verification

Where the aqua registry includes a checksum manifest, downloaded binaries are verified against it before installation. Verification behaviour depends on the checksum type advertised:

- **Strong checksums (sha256, sha512, etc.)** — verified before the binary is installed. If the downloaded archive does not match, the install is aborted and an error is returned (fails closed).
- **Unsupported or weak checksum types (e.g. md5, sha1)** — skipped with a warning; installation proceeds without verification.
- **No manifest** — if no checksum is advertised in the registry entry, the binary is installed without verification.

### version_overrides Resolution

The auto-installer correctly resolves **`version_overrides`** entries in the aqua registry. Many common tools (for example, `fzf`) keep their package configuration — including download URLs and checksums — under `version_overrides` rather than at the top level of their registry entry. These tools previously failed to install silently; they are now handled correctly.

### Disabling Auto-Install

**Per toolset** — set `version` to `"false"` or `"off"`:

```yaml
toolsets:
  - type: mcp
    command: my-custom-server
    version: "false"
```

**Globally** — set the `DOCKER_AGENT_AUTO_INSTALL` environment variable:

```bash
export DOCKER_AGENT_AUTO_INSTALL=false
```

### Environment Variables

| Variable                     | Default            | Description                                      |
| ---------------------------- | ------------------ | ------------------------------------------------ |
| `DOCKER_AGENT_AUTO_INSTALL`  | (enabled)          | Set to `false` to disable all auto-installation  |
| `DOCKER_AGENT_TOOLS_DIR`     | `~/.cagent/tools/` | Base directory for installed tools               |
| `GITHUB_TOKEN`               | —                  | GitHub token to raise API rate limits (optional) |

Installed binaries are placed in `~/.cagent/tools/bin/` and cached so they are only downloaded once.

> [!TIP]
> Auto-install supports both Go packages (via `go install`) and GitHub release binaries (via archive download). The aqua registry metadata determines which method is used.

## Toolset Lifecycle

Long-running toolsets — local MCP servers (stdio), remote MCP servers (Streamable HTTP / SSE), and LSP servers — are managed by a single supervisor that can auto-reconnect them when they crash, time out, or drop their session. The `lifecycle` block on the toolset lets you tune that supervisor per toolset. It applies to every `type: mcp` and `type: lsp` toolset.

The simplest knob is `profile`, which picks a preset:

| Profile | Auto-restart | Use case |
| --- | --- | --- |
| `resilient` | Yes | Default. Exponential backoff on disconnect; the agent keeps running if the toolset is unavailable. Matches the historical docker-agent behaviour. |
| `strict` | No | Fail-fast. Marks the toolset as required. Intended for CI / headless runs where a missing dependency should be a hard error. |
| `best-effort` | No | Single attempt, no retries. Good for experimental MCPs whose flakiness should not amplify into a restart loop. |

```yaml
toolsets:
  - type: mcp
    ref: docker:duckduckgo
    lifecycle:
      profile: resilient   # default; shown here for clarity

  - type: lsp
    command: gopls
    file_types: [".go"]
    lifecycle:
      profile: strict

  - type: mcp
    ref: docker:openbnb-airbnb
    lifecycle:
      profile: best-effort
```

### Tuning the defaults

Any field set on `lifecycle` overrides the profile preset, so you can mix-and-match: pick a profile and only override the knobs you care about.

```yaml
toolsets:
  - type: mcp
    command: ["docker", "mcp", "gateway"]
    lifecycle:
      profile: resilient
      max_restarts: 10        # keep trying longer than the default of 5
      backoff:
        initial: 500ms
        max: 1m
        multiplier: 2
        jitter: 0.2           # 20% random offset to avoid thundering-herd retries
```

| Property | Type | Description |
| --- | --- | --- |
| `profile` | string | One of `resilient` (default), `strict`, `best-effort`. Picks defaults for every other field. |
| `restart` | string | When the supervisor should reconnect after a disconnect: `never`, `on_failure` (default), or `always`. For **remote** MCP toolsets (Streamable HTTP / SSE), `on_failure` is automatically promoted to `always` so idle-timeout closes reconnect gracefully — `never` is still honored. |
| `max_restarts` | int | Maximum consecutive restart attempts before the toolset is marked `Failed`. `0` uses the profile default (5); `-1` means unlimited. |
| `backoff.initial` | duration | First wait between attempts (Go duration: `500ms`, `1s`, …). Default: `1s`. |
| `backoff.max` | duration | Cap on the wait between attempts. Default: `32s`. |
| `backoff.multiplier` | number | Multiplier applied each attempt. Default: `2`. |
| `backoff.jitter` | number | Fraction (0..1) of the computed delay applied as a uniform random offset. `0` disables jitter (default). |
| `required` | boolean | Marks the toolset as critical. Today this is informational; a future eager-startup phase will refuse to start the agent when a required toolset cannot reach Ready. Defaults to `true` under `strict`, `false` otherwise. |
| `startup_timeout` | duration | Cap on the initial connect+initialize duration. Today this is informational; the eager-startup phase that enforces it ships in a follow-up. |
| `call_timeout` | duration | Documented per-call timeout. Informational; the runtime currently uses the caller's context for cancellation. |

> [!NOTE]
> **`required` and `startup_timeout` are not yet enforced**
>
> The schema validates these fields and the supervisor stores them, but no code path acts on them yet. They are documented now so config files written today keep working when the planned eager-startup phase lands. Picking the `strict` profile is forward-compatible — it will start enforcing `required=true` automatically.

### Inspecting and restarting toolsets at runtime

The TUI exposes the supervisor through two slash commands:

- `/tools` — the unified tools dialog. Its top section lists every toolset on the current agent with its lifecycle state (`Stopped`, `Starting`, `Ready`, `Degraded`, `Restarting`, `Failed`), restart count, and last error; its bottom section lists every tool the agent can call, grouped by category. Use this to answer both "what can the agent do?" and "is anything degraded?" with one command.
- `/toolset-restart <name>` — force the supervisor to reconnect the named toolset. Useful after completing OAuth, when a remote MCP server has been redeployed, or when an LSP like `gopls` is stuck.

See the [TUI reference](../../features/tui/index.md) for the full list of slash commands.

See [`examples/lifecycle.yaml`](https://github.com/docker/docker-agent/blob/main/examples/lifecycle.yaml) for a complete lifecycle configuration example.

## TOON-Encoded Tool Outputs

Many MCP servers return verbose JSON responses that consume a lot of context budget. The `toon` field on a toolset transparently re-encodes matching tools' JSON output as [TOON](https://github.com/alpkeskin/gotoon) — a compact, model-friendly key/value format — before the result is shown to the model.

```yaml
toolsets:
  - type: mcp
    ref: docker:github-official
    toon: ".*"          # toonify every tool from this MCP server
  - type: mcp
    command: my-server
    toon: "list_.*,get_.*" # only toonify list_/get_ tools
```

| Property | Type   | Description |
| -------- | ------ | ----------- |
| `toon`   | string | Comma-delimited list of regular expressions matching tool names whose JSON output should be re-encoded as TOON. Non-JSON outputs and non-matching tools are passed through untouched. |

When a tool's output is not valid JSON, it is returned unchanged — TOON encoding is best-effort and never breaks tools that emit plain text.

> [!NOTE]
> **When to use TOON**
>
> TOON typically yields 30-60% smaller payloads than equivalent JSON for MCP tools that return arrays of records (issue lists, search results, file listings, …). It works best when the schema is regular; one-off responses with deeply nested or heterogeneous shapes may benefit less.

## Per-Toolset Model Routing

The `model` field on a toolset overrides which LLM is invoked for the **next turn** after a tool from that toolset returns — letting you process simple tool results (file reads, knowledge-base lookups, shell stdout) with a cheaper or faster model while keeping the agent's primary model for reasoning.

```yaml
models:
  primary:
    provider: anthropic
    model: claude-sonnet-4-5
  fast:
    provider: anthropic
    model: claude-haiku-4-5

agents:
  root:
    model: primary
    toolsets:
      - type: filesystem
        model: fast            # process file reads with the fast model
      - type: shell
        model: fast            # ditto for shell stdout
      - type: mcp
        ref: docker:github-official
        model: openai/gpt-4o-mini  # inline provider/model also works
```

| Property | Type   | Description |
| -------- | ------ | ----------- |
| `model`  | string | Model used for the LLM turn that processes tool results from this toolset. Either a name from the `models:` section or an inline `provider/model` (e.g. `openai/gpt-4o-mini`). The override is **one-shot**: subsequent turns return to the agent's primary model. |

When multiple tool calls in a single turn come from toolsets with different `model` overrides, the runtime picks the override of the **first** tool call that has one set. See [`examples/per_tool_model_routing.yaml`](https://github.com/docker/docker-agent/blob/main/examples/per_tool_model_routing.yaml) for a complete configuration.

## Tool Filtering

Toolsets may expose many tools. Use the `tools` property to whitelist only the ones your agent needs. This works for any toolset type — not just MCP:

```yaml
toolsets:
  - type: mcp
    ref: docker:github-official
    tools: ["list_issues", "create_issue", "get_pull_request"]
  - type: filesystem
    tools: ["read_file", "search_files_content"]
  - type: shell
    tools: ["shell"]
```

> [!TIP]
> Filtering tools improves agent performance — fewer tools means less confusion for the model about which tool to use.

## Tool Instructions

Add context-specific instructions that get injected when a toolset is loaded:

```yaml
toolsets:
  - type: mcp
    ref: docker:github-official
    instruction: |
      Use these tools to manage GitHub issues.
      Always check for existing issues before creating new ones.
      Label new issues with 'triage' by default.
```

By default, the `instruction:` field **replaces** the toolset's built-in instructions (if any). To keep the built-in guidance and add your own rules on top, include the `{ORIGINAL_INSTRUCTIONS}` placeholder anywhere in your instruction text. At runtime it expands to the toolset's default instructions:

```yaml
toolsets:
  # Enrich: keep built-in instructions, then add your own rules
  - type: filesystem
    instruction: |
      {ORIGINAL_INSTRUCTIONS}

      ## Project-specific rules
      - Never modify files outside the `src/` directory.
      - Always create a backup before overwriting a file.

  # Enrich: prepend your rules before the built-in instructions
  - type: shell
    instruction: |
      Important: only run commands inside the project root.
      {ORIGINAL_INSTRUCTIONS}

  # Replace: omit the placeholder to discard built-in instructions entirely
  - type: mcp
    ref: docker:github-official
    instruction: |
      Only read GitHub issues. Never create, edit, or close anything.
```

Three patterns at a glance:

| Pattern | Description |
| --- | --- |
| `{ORIGINAL_INSTRUCTIONS}` then your text | Append your rules after the defaults |
| Your text then `{ORIGINAL_INSTRUCTIONS}` | Prepend your rules before the defaults |
| No placeholder | Replace the defaults entirely |

See [`examples/toolset_instructions.yaml`](https://github.com/docker/docker-agent/blob/main/examples/toolset_instructions.yaml) for a complete example.

## Deferred Tool Loading

Load tools on-demand to speed up agent startup. When a toolset is deferred, its tools are registered lazily — the tool server process is not started until the agent first calls one of its tools. This is useful for large toolsets (e.g., an MCP server with hundreds of tools) where startup time matters.

```yaml
toolsets:
  - type: mcp
    ref: docker:github-official
    defer: true
  - type: mcp
    ref: docker:slack
    defer: true
  - type: filesystem
```

Or defer specific tools within a toolset:

```yaml
toolsets:
  - type: mcp
    ref: docker:github-official
    defer:
      - "list_issues"
      - "search_repos"
```

When `defer` is a list of tool names, only those specific tools are deferred; all other tools in the toolset load eagerly. Setting `defer: true` defers the entire toolset.

See [`examples/deferred.yaml`](https://github.com/docker/docker-agent/blob/main/examples/deferred.yaml) for a complete example.

## Combined Example

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    description: Full-featured developer assistant
    instruction: You are an expert developer.
    toolsets:
      # Built-in tools
      - type: filesystem
      - type: shell
      - type: think
      - type: todo
      - type: memory
        path: ./dev.db
      - type: user_prompt
      # LSP for code intelligence
      - type: lsp
        command: gopls
        file_types: [".go"]
      # Custom scripts
      - type: script
        shell:
          run_tests:
            description: Run the test suite
            cmd: task test
          lint:
            description: Run the linter
            cmd: task lint
      # Custom API tool
      - type: api
        api_config:
          name: get_status
          method: GET
          endpoint: "https://api.example.com/status"
          instruction: Check service health
      # Docker MCP tools
      - type: mcp
        ref: docker:github-official
        tools: ["list_issues", "create_issue"]
      - type: mcp
        ref: docker:duckduckgo
      # Remote MCP
      - type: mcp
        remote:
          url: "https://internal-api.example.com/mcp"
          transport_type: "streamable"
          headers:
            Authorization: "Bearer ${env.INTERNAL_TOKEN}"
```

> [!WARNING]
> **Toolset Order Matters**
>
> If multiple toolsets provide a tool with the same name, the first one wins. Order your toolsets intentionally.
