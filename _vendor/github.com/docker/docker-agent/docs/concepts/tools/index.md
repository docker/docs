---
title: "Tools"
description: "Tools give agents the ability to interact with the world — read files, run commands, search the web, query databases, and more."
keywords: docker agent, ai agents, concepts, tools
weight: 30
canonical: https://docs.docker.com/ai/docker-agent/concepts/tools/
---

_Tools give agents the ability to interact with the world — read files, run commands, search the web, query databases, and more._

## How Tools Work

When an agent needs to perform an action, it makes a **tool call**. The Docker Agent runtime executes the tool and returns the result to the agent, which can then use it to continue its work.

1. Agent receives a user message
2. Agent decides it needs to use a tool (e.g., read a file)
3. Docker Agent executes the tool and returns the result
4. Agent incorporates the result and responds

> [!NOTE]
> **Tool Confirmation**
>
> By default, Docker Agent asks for user confirmation before executing tools that have side effects (shell commands, file writes). Use `--yolo` to auto-approve all tool calls.

## Built-in Tools

Enable a built-in toolset by adding its `type` to the agent's `toolsets` list. For example:

```yaml
toolsets:
  - type: filesystem
  - type: shell
  - type: todo
```

Choose tools for the capabilities your agent needs:

- **Work with files and processes:** [Filesystem](../../tools/filesystem/index.md) and [Shell](../../tools/shell/index.md).
- **Plan and remember:** [Todo](../../tools/todo/index.md), [Memory](../../tools/memory/index.md), and [RAG](../../tools/rag/index.md).
- **Coordinate agents:** [Transfer Task](../../tools/transfer-task/index.md) and [Background Agents](../../tools/background-agents/index.md).

See [Built-in Tools](../../configuration/tools/index.md#built-in-tools) for the complete catalog, including each toolset's type and configuration reference.

## MCP Tools

Docker Agent supports the [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) for extending agents with external tools. There are three ways to connect MCP tools:

- **Docker MCP** (recommended) — Run MCP servers in Docker containers via the [MCP Gateway](https://github.com/docker/mcp-gateway). Browse the [Docker MCP Catalog](https://hub.docker.com/search?q=&type=mcp).
- **Local MCP (stdio)** — Run MCP servers as local processes communicating over stdin/stdout.
- **Remote MCP (Streamable HTTP / SSE)** — Connect to MCP servers running on a network. See [Remote MCP Servers](../../features/remote-mcp/index.md).

```yaml
toolsets:
  - type: mcp
    ref: docker:duckduckgo
```

See [Tool Config](../../configuration/tools/index.md#mcp-tools) for full MCP configuration reference.

> [!TIP]
> **See also**
>
> For full configuration reference, see [Tool Config](../../configuration/tools/index.md).
