---
title: "Tools"
description: "Tools give agents the ability to interact with the world — read files, run commands, search the web, query databases, and more."
keywords: docker agent, ai agents, concepts, tools
weight: 30
canonical: https://docs.docker.com/ai/docker-agent/concepts/tools/
---

_Tools give agents the ability to interact with the world — read files, run commands, search the web, query databases, and more._

## How Tools Work

When an agent needs to perform an action, it makes a **tool call**. The docker-agent runtime executes the tool and returns the result to the agent, which can then use it to continue its work.

1. Agent receives a user message
2. Agent decides it needs to use a tool (e.g., read a file)
3. docker-agent executes the tool and returns the result
4. Agent incorporates the result and responds

> [!NOTE]
> **Tool Confirmation**
>
> By default, docker-agent asks for user confirmation before executing tools that have side effects (shell commands, file writes). Use `--yolo` to auto-approve all tool calls.

## Built-in Tools

docker-agent ships with several built-in tools that require no external dependencies. Each is enabled by adding its `type` to the agent's `toolsets` list:

| Tool | Description |
| --- | --- |
| [Filesystem](../../tools/filesystem/index.md) | Read, write, list, search, and navigate files and directories |
| [Shell](../../tools/shell/index.md) | Execute shell commands synchronously |
| [Background Jobs](../../tools/background-jobs/index.md) | Run and manage long-running shell commands |
| [Think](../../tools/think/index.md) | Step-by-step reasoning scratchpad for planning and decision-making |
| [Todo](../../tools/todo/index.md) | Task list management for complex multi-step workflows |
| [Tasks](../../tools/tasks/index.md) | Persistent task database shared across sessions |
| [Memory](../../tools/memory/index.md) | Persistent key-value storage backed by SQLite |
| [Fetch](../../tools/fetch/index.md) | Read content from HTTP/HTTPS URLs (GET only) |
| [Script](../../tools/script/index.md) | Define custom shell scripts as named tools |
| [LSP](../../tools/lsp/index.md) | Connect to Language Server Protocol servers for code intelligence |
| [API](../../tools/api/index.md) | Create custom tools that call HTTP APIs without writing code |
| [OpenAPI](../../tools/openapi/index.md) | Generate tools from an OpenAPI 3.x document |
| [RAG](../../tools/rag/index.md) | Retrieval-augmented generation over indexed sources |
| [Model Picker](../../tools/model-picker/index.md) | Let the agent pick between several models per turn |
| [User Prompt](../../tools/user-prompt/index.md) | Ask users questions and collect interactive input |
| [Open URL](../../tools/open-url/index.md) | Open a fixed URL in the user's default browser |
| [Transfer Task](../../tools/transfer-task/index.md) | Delegate tasks to sub-agents (auto-enabled with `sub_agents`) |
| [Background Agents](../../tools/background-agents/index.md) | Dispatch work to sub-agents concurrently |
| [Handoff](../../tools/handoff/index.md) | Hand the conversation off to another local agent in the same config (auto-enabled with `handoffs:`) |
| [A2A](../../tools/a2a/index.md) | Connect to remote agents via the Agent-to-Agent protocol |
| [MCP Catalog](../../tools/mcp-catalog/index.md) | Discover and activate remote MCP servers from the Docker MCP Catalog on demand |

## MCP Tools

docker-agent supports the [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) for extending agents with external tools. There are three ways to connect MCP tools:

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
