---
title: Dynamic MCP
linkTitle: Dynamic MCP
description: Discover and add MCP servers on-demand using natural language with Dynamic MCP servers
keywords: dynamic mcps, mcp discovery, mcp-find, mcp-add, code-mode, ai agents, model context protocol
weight: 35
params:
  sidebar:
    badge:
      color: green
      text: New
---

Dynamic MCP enables AI agents to discover and add MCP servers on-demand during
a conversation, without manual configuration. Instead of pre-configuring every
MCP server before starting your agent session, clients can search the
[MCP Catalog](/manuals/ai/mcp-catalog-and-toolkit/catalog.md) and add servers
as needed.

This capability is enabled automatically when you connect an MCP client to the
[MCP Toolkit](/manuals/ai/mcp-catalog-and-toolkit/toolkit.md). The gateway
provides a set of primordial tools that agents use to discover and manage
servers during runtime.

{{% experimental %}}

Dynamic MCP is an experimental feature in early development. While you're
welcome to try it out and explore its capabilities, you may encounter
unexpected behavior or limitations. Feedback is welcome via at [GitHub
issues](https://github.com/docker/mcp-gateway/issues) for bug reports and
[GitHub discussions](https://github.com/docker/mcp-gateway/discussions) for
general questions and feature requests.

{{% /experimental %}}

## How it works

When you connect a client to the MCP Gateway, the gateway exposes a small set
of management tools alongside any MCP servers you've already enabled. These
management tools let agents interact with the gateway's configuration:

| Tool             | Description                                                              |
| ---------------- | ------------------------------------------------------------------------ |
| `mcp-find`       | Search for MCP servers in the catalog by name or description             |
| `mcp-add`        | Add a new MCP server to the current session                              |
| `mcp-config-set` | Configure settings for an MCP server                                     |
| `mcp-remove`     | Remove an MCP server from the session                                    |
| `mcp-exec`       | Execute a tool by name that exists in the current session                |
| `code-mode`      | Create a JavaScript-enabled tool that combines multiple MCP server tools |

With these tools available, an agent can search the catalog, add servers,
handle authentication, and use newly added tools directly without requiring a
restart or manual configuration.

Dynamically added servers and tools are associated with your _current session
only_. When you start a new session, previously added servers are not
automatically included.

## Prerequisites

To use Dynamic MCP, you need:

- Docker Desktop version 4.50 or later, with [MCP Toolkit](/manuals/ai/mcp-catalog-and-toolkit/toolkit.md) enabled
- An LLM application that supports MCP (such as Claude Desktop, Visual Studio Code, or Claude Code)
- Your client configured to connect to the MCP Gateway

See [Get started with Docker MCP Toolkit](/manuals/ai/mcp-catalog-and-toolkit/get-started.md)
for setup instructions.

## Usage

Dynamic MCP is enabled automatically when you use the MCP Toolkit. Your
connected clients can now use `mcp-find`, `mcp-add`, and other management tools
during conversations.

To see Dynamic MCP in action, connect your AI client to the Docker MCP Toolkit
and try this prompt:

```plaintext
What MCP servers can I use for working with SQL databases?
```

Given this prompt, your agent will use the `mcp-find` tool provided by MCP
Toolkit to search for SQL-related servers in the [MCP Catalog](./catalog.md).

And to add a server to a session, simply write a prompt and the MCP Toolkit
takes care of installing and running the server:

```plaintext
Add the postgres mcp server
```

## Tool composition with code mode

The `code-mode` tool is available as an experimental capability for creating
custom JavaScript functions that combine multiple MCP server tools. The
intended use case is to enable workflows that coordinate multiple services
in a single operation.

> **Note**
>
> Code mode is in early development and is not yet reliable for general use.
> The documentation intentionally omits usage examples at this time.
>
> The core Dynamic MCP capabilities (`mcp-find`, `mcp-add`, `mcp-config-set`,
> `mcp-remove`) work as documented and are the recommended focus for current
> use.

The architecture works as follows:

1. The agent calls `code-mode` with a list of server names and a tool name
2. The gateway creates a sandbox with access to those servers' tools
3. A new tool is registered in the current session with the specified name
4. The agent calls the newly created tool
5. The code executes in the sandbox with access to the specified tools
6. Results are returned to the agent

The sandbox can only interact with the outside world through MCP tools,
which are already running in isolated containers with restricted privileges.

## Security considerations

Dynamic MCP maintains the same security model as static MCP server
configuration in MCP Toolkit:

- All servers in the MCP Catalog are built, signed, and maintained by Docker
- Servers run in isolated containers with restricted resources
- Code mode runs agent-written JavaScript in an isolated sandbox that can only
  interact through MCP tools
- Credentials are managed by the gateway and injected securely into containers

The key difference with dynamic capabilities is that agents can add new tools
during runtime.

## Disabling Dynamic MCP

Dynamic MCP is enabled by default in the MCP Toolkit. If you prefer to use only
statically configured MCP servers, you can disable the dynamic tools feature:

```console
$ docker mcp feature disable dynamic-tools
```

To re-enable the feature later:

```console
$ docker mcp feature enable dynamic-tools
```

After changing this setting, you may need to restart any connected MCP clients.

## Further reading

Check out the [Dynamic MCP servers with Docker](https://docker.com/blog) blog
post for more examples and inspiration on how you can use dynamic tools.
