---
title: "LSP Tool"
description: "Connect to Language Server Protocol servers for code intelligence."
keywords: docker agent, ai agents, tools, toolsets, lsp tool
linkTitle: "LSP"
weight: 220
canonical: https://docs.docker.com/ai/docker-agent/tools/lsp/
---

_Connect to Language Server Protocol servers for code intelligence._

## Overview

The LSP tool connects your agent to any Language Server Protocol (LSP) server, providing comprehensive code intelligence capabilities like go-to-definition, find references, diagnostics, and more.

> [!NOTE]
> **What is LSP?**
>
> The [Language Server Protocol](https://microsoft.github.io/language-server-protocol/) is a standard for providing language features like autocomplete, go-to-definition, and diagnostics. Most programming languages have LSP servers available.

## Available Tools

The LSP toolset provides these tools to the agent:

| Tool                    | Description                                   | Read-Only |
| ----------------------- | --------------------------------------------- | --------- |
| `lsp_workspace`         | Get workspace info and available capabilities | ✓         |
| `lsp_hover`             | Get type info and documentation for a symbol  | ✓         |
| `lsp_definition`        | Find where a symbol is defined                | ✓         |
| `lsp_references`        | Find all references to a symbol               | ✓         |
| `lsp_document_symbols`  | List all symbols in a file                    | ✓         |
| `lsp_workspace_symbols` | Search symbols across the workspace           | ✓         |
| `lsp_diagnostics`       | Get errors and warnings for a file            | ✓         |
| `lsp_code_actions`      | Get available quick fixes and refactorings    | ✓         |
| `lsp_rename`            | Rename a symbol across the workspace          | ✗         |
| `lsp_format`            | Format a file                                 | ✗         |
| `lsp_call_hierarchy`    | Find incoming/outgoing calls                  | ✓         |
| `lsp_type_hierarchy`    | Find supertypes/subtypes                      | ✓         |
| `lsp_implementations`   | Find interface implementations                | ✓         |
| `lsp_signature_help`    | Get function signature at call site           | ✓         |
| `lsp_inlay_hints`       | Get type annotations and parameter names      | ✓         |

## Configuration

```yaml
agents:
  developer:
    model: anthropic/claude-sonnet-4-5
    description: Code developer with LSP support
    instruction: You are a software developer.
    toolsets:
      - type: lsp
        command: gopls
        args: []
        file_types: [".go"]
      - type: filesystem
      - type: shell
```

## Properties

| Property      | Type   | Required | Description                                                                                                                  |
| ------------- | ------ | -------- | ---------------------------------------------------------------------------------------------------------------------------- |
| `command`     | string | ✓        | LSP server executable command                                                                                                |
| `args`        | array  | ✗        | Command-line arguments for the LSP server                                                                                    |
| `env`         | object | ✗        | Environment variables for the LSP process                                                                                    |
| `file_types`  | array  | ✗        | File extensions this LSP handles (e.g., `[".go", ".mod"]`)                                                                   |
| `working_dir` | string | ✗        | Working directory for the LSP server process. Relative paths are resolved against the agent's working directory. Defaults to the agent's working directory when omitted. |
| `version`     | string | ✗        | Package reference for [auto-installing](../../configuration/tools/index.md#auto-installing-tools) the command binary |

## Common LSP Servers

Here are configurations for popular languages:

### Go (gopls)

```yaml
toolsets:
  - type: lsp
    command: gopls
    version: "golang/tools@v0.21.0" # optional: auto-install if not in PATH
    file_types: [".go"]
```

If your Go module lives in a subdirectory (e.g. a monorepo where `go.mod` is under `./backend`), set `working_dir` so `gopls` is started from the module root:

```yaml
toolsets:
  - type: lsp
    command: gopls
    file_types: [".go"]
    working_dir: ./backend # gopls must be started from the module root
```

### TypeScript/JavaScript (typescript-language-server)

```yaml
toolsets:
  - type: lsp
    command: typescript-language-server
    args: ["--stdio"]
    file_types: [".ts", ".tsx", ".js", ".jsx"]
```

### Python (pylsp)

```yaml
toolsets:
  - type: lsp
    command: pylsp
    file_types: [".py"]
```

### Rust (rust-analyzer)

```yaml
toolsets:
  - type: lsp
    command: rust-analyzer
    file_types: [".rs"]
```

### C/C++ (clangd)

```yaml
toolsets:
  - type: lsp
    command: clangd
    file_types: [".c", ".cpp", ".h", ".hpp"]
```

## Multiple LSP Servers

You can configure multiple LSP servers for different file types:

```yaml
agents:
  polyglot:
    model: anthropic/claude-sonnet-4-5
    description: Multi-language developer
    instruction: You are a full-stack developer.
    toolsets:
      - type: lsp
        command: gopls
        file_types: [".go"]
      - type: lsp
        command: typescript-language-server
        args: ["--stdio"]
        file_types: [".ts", ".tsx", ".js", ".jsx"]
      - type: lsp
        command: pylsp
        file_types: [".py"]
      - type: filesystem
      - type: shell
```

## Workflow Instructions

The LSP tool includes built-in instructions that guide the agent on how to use it effectively. The agent learns to:

1. Start with `lsp_workspace` to understand available capabilities
2. Use `lsp_workspace_symbols` to find relevant code
3. Use `lsp_references` before modifying any symbol
4. Check `lsp_diagnostics` after every code change
5. Apply `lsp_format` after edits are complete

> [!TIP]
> **Best Practice**
>
> Always include the `filesystem` tool alongside LSP. The agent needs filesystem access to read and write code files, while LSP provides intelligence about the code.

## Capability Detection

Not all LSP servers support all features. During the `initialize` handshake, docker-agent reads the server's `ServerCapabilities` and **filters out the `lsp_*` tools the server does not advertise**. The model never sees, for example, `lsp_inlay_hints` against a server that doesn't support it, so it can't waste a turn calling a tool that would only fail.

The agent uses `lsp_workspace` to discover what's available:

```text
Workspace Information:
- Root: /path/to/project
- Server: gopls v0.14.0
- File types: .go

Available Capabilities:
- Hover: Yes
- Go to Definition: Yes
- Find References: Yes
- Rename: Yes
- Code Actions: Yes
- Formatting: Yes
- Call Hierarchy: Yes
- Type Hierarchy: Yes
...
```

## Auto-Restart and Lifecycle

LSP toolsets are managed by the same supervisor as MCP toolsets, so a crashed `gopls` (or any other language server) is reconnected automatically with exponential backoff. Use the [`lifecycle`](../../configuration/tools/index.md#toolset-lifecycle) block to tune the policy per toolset — for example, mark `gopls` as `strict` if your CI flow requires it to be available, or use `/toolset-restart gopls` from the TUI to force a reconnect when the server gets stuck.

```yaml
toolsets:
  - type: lsp
    command: gopls
    file_types: [".go"]
    lifecycle:
      profile: resilient # default: auto-restart on crash with exponential backoff
```

## Position Format

All LSP tools use **1-based** line and character positions:

- Line 1 is the first line of the file
- Character 1 is the first character on a line

```json
{
  "file": "/path/to/file.go",
  "line": 42,
  "character": 15
}
```

> [!TIP]
> **Auto-Installation**
>
> docker-agent can automatically download and install LSP servers if they are not found in your PATH. Use the `version` property to specify a package, or let docker-agent auto-detect it from the command name. See [Auto-Installing Tools](../../configuration/tools/index.md#auto-installing-tools) for details.
